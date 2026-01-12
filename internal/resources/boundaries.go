package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/jtimothystewart/dtiam/internal/client"
)

// BoundaryHandler handles boundary resources.
type BoundaryHandler struct {
	BaseHandler
}

// NewBoundaryHandler creates a new boundary handler.
func NewBoundaryHandler(c *client.Client) *BoundaryHandler {
	path := fmt.Sprintf("/repo/account/%s/boundaries", c.AccountUUID())
	return &BoundaryHandler{
		BaseHandler: BaseHandler{
			Client:    c,
			Name:      "boundary",
			Path:      path,
			ListKey:   "boundaries",
			IDField:   "uuid",
			NameField: "name",
		},
	}
}

// Create creates a new boundary.
func (h *BoundaryHandler) Create(ctx context.Context, name string, managementZones []string, boundaryQuery, description *string) (map[string]any, error) {
	data := map[string]any{
		"name": name,
	}

	if description != nil {
		data["description"] = *description
	}

	// Build boundary query from management zones if provided
	if len(managementZones) > 0 {
		query := h.buildZoneQuery(managementZones)
		data["boundaryQuery"] = query
	} else if boundaryQuery != nil {
		data["boundaryQuery"] = *boundaryQuery
	} else {
		return nil, fmt.Errorf("either managementZones or boundaryQuery is required")
	}

	return h.BaseHandler.Create(ctx, data)
}

// Update updates an existing boundary.
func (h *BoundaryHandler) Update(ctx context.Context, boundaryID string, name *string, managementZones []string, boundaryQuery, description *string) (map[string]any, error) {
	// Get existing boundary
	existing, err := h.Get(ctx, boundaryID)
	if err != nil {
		return nil, err
	}

	// Merge updates
	if name != nil {
		existing["name"] = *name
	}
	if description != nil {
		existing["description"] = *description
	}
	if len(managementZones) > 0 {
		existing["boundaryQuery"] = h.buildZoneQuery(managementZones)
	} else if boundaryQuery != nil {
		existing["boundaryQuery"] = *boundaryQuery
	}

	return h.BaseHandler.Update(ctx, boundaryID, existing)
}

// GetAttachedPolicies gets policies that use this boundary.
func (h *BoundaryHandler) GetAttachedPolicies(ctx context.Context, boundaryID string) ([]map[string]any, error) {
	// Get all bindings and filter by boundary
	bindingHandler := NewBindingHandler(h.Client)
	bindings, err := bindingHandler.List(ctx, nil)
	if err != nil {
		return nil, err
	}

	var attached []map[string]any
	for _, binding := range bindings {
		boundaries, ok := binding["boundaries"].([]string)
		if !ok {
			if boundariesAny, ok := binding["boundaries"].([]any); ok {
				for _, b := range boundariesAny {
					if bStr, ok := b.(string); ok {
						if bStr == boundaryID {
							attached = append(attached, map[string]any{
								"policyUuid": binding["policyUuid"],
								"groupUuid":  binding["groupUuid"],
							})
							break
						}
					}
				}
			}
			continue
		}

		for _, b := range boundaries {
			if b == boundaryID {
				attached = append(attached, map[string]any{
					"policyUuid": binding["policyUuid"],
					"groupUuid":  binding["groupUuid"],
				})
				break
			}
		}
	}

	return attached, nil
}

// buildZoneQuery builds a boundary query from management zone names.
func (h *BoundaryHandler) buildZoneQuery(managementZones []string) string {
	var conditions []string
	for _, zone := range managementZones {
		conditions = append(conditions, fmt.Sprintf(`managementZone.name = "%s"`, zone))
	}

	whereClause := strings.Join(conditions, " OR ")

	// Build query for each resource type
	resourceTypes := []string{"environment", "storage", "settings"}
	var queries []string
	for _, rt := range resourceTypes {
		queries = append(queries, fmt.Sprintf("ALLOW %s:* WHERE %s", rt, whereClause))
	}

	return strings.Join(queries, ";")
}
