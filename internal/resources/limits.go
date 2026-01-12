package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jtimothystewart/dtiam/internal/client"
)

// LimitsHandler handles account limits resources.
type LimitsHandler struct {
	BaseHandler
}

// NewLimitsHandler creates a new limits handler.
func NewLimitsHandler(c *client.Client) *LimitsHandler {
	return &LimitsHandler{
		BaseHandler: BaseHandler{
			Client:    c,
			Name:      "limit",
			Path:      "/limits",
			ListKey:   "items",
			IDField:   "name",
			NameField: "name",
		},
	}
}

// List lists account limits.
func (h *LimitsHandler) List(ctx context.Context, params map[string]string) ([]map[string]any, error) {
	body, err := h.Client.Get(ctx, h.Path, params)
	if err != nil {
		return nil, h.handleError("list", err)
	}

	return h.extractList(body)
}

// Get gets a limit by name.
func (h *LimitsHandler) Get(ctx context.Context, name string) (map[string]any, error) {
	items, err := h.List(ctx, nil)
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		if itemName, ok := item["name"].(string); ok {
			if strings.EqualFold(itemName, name) {
				return item, nil
			}
		}
	}

	return nil, fmt.Errorf("limit %q not found", name)
}

// GetSummary returns a comprehensive summary of all limits.
func (h *LimitsHandler) GetSummary(ctx context.Context) (map[string]any, error) {
	items, err := h.List(ctx, nil)
	if err != nil {
		return nil, err
	}

	var nearCapacity, atCapacity int
	limits := make([]map[string]any, 0, len(items))

	for _, item := range items {
		current := h.getNumericValue(item, "current", "value")
		max := h.getNumericValue(item, "max", "limit")

		var usagePercent float64
		var status string
		var available int

		if max > 0 {
			usagePercent = float64(current) / float64(max) * 100
			available = max - current

			if usagePercent >= 100 {
				status = "at_capacity"
				atCapacity++
			} else if usagePercent >= 80 {
				status = "near_capacity"
				nearCapacity++
			} else {
				status = "ok"
			}
		} else {
			status = "unknown"
		}

		limits = append(limits, map[string]any{
			"name":          item["name"],
			"current":       current,
			"max":           max,
			"usage_percent": usagePercent,
			"available":     available,
			"status":        status,
		})
	}

	return map[string]any{
		"limits":              limits,
		"total_limits":        len(limits),
		"limits_near_capacity": nearCapacity,
		"limits_at_capacity":   atCapacity,
	}, nil
}

// CheckCapacity checks if there is capacity for additional resources.
func (h *LimitsHandler) CheckCapacity(ctx context.Context, limitName string, additional int) (map[string]any, error) {
	if additional <= 0 {
		additional = 1
	}

	limit, err := h.Get(ctx, limitName)
	if err != nil {
		return map[string]any{
			"limit_name":  limitName,
			"found":       false,
			"has_capacity": false,
			"message":     fmt.Sprintf("Limit %q not found", limitName),
		}, nil
	}

	current := h.getNumericValue(limit, "current", "value")
	max := h.getNumericValue(limit, "max", "limit")
	available := max - current
	hasCapacity := available >= additional

	var message string
	if hasCapacity {
		message = fmt.Sprintf("Capacity available: %d/%d used, %d available, requesting %d",
			current, max, available, additional)
	} else {
		message = fmt.Sprintf("Insufficient capacity: %d/%d used, %d available, requesting %d",
			current, max, available, additional)
	}

	return map[string]any{
		"limit_name":  limitName,
		"found":       true,
		"has_capacity": hasCapacity,
		"current":     current,
		"max":         max,
		"available":   available,
		"requested":   additional,
		"message":     message,
	}, nil
}

// getNumericValue gets a numeric value from a map, trying multiple keys.
func (h *LimitsHandler) getNumericValue(m map[string]any, keys ...string) int {
	for _, key := range keys {
		if val, ok := m[key]; ok {
			switch v := val.(type) {
			case float64:
				return int(v)
			case int:
				return v
			case int64:
				return int(v)
			}
		}
	}
	return 0
}

// extractList handles limit-specific response formats.
func (h *LimitsHandler) extractList(body []byte) ([]map[string]any, error) {
	var response map[string]any
	if err := json.Unmarshal(body, &response); err != nil {
		// Try parsing as array
		var items []map[string]any
		if err := json.Unmarshal(body, &items); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}
		return items, nil
	}

	// Try common list keys
	for _, key := range []string{"items", "limits"} {
		if items, ok := response[key]; ok {
			return toMapSlice(items)
		}
	}

	// Single item response - wrap in array
	if _, ok := response["name"]; ok {
		return []map[string]any{response}, nil
	}

	return []map[string]any{}, nil
}
