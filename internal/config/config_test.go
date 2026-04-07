package config

import (
	"testing"
)

func TestNewConfig(t *testing.T) {
	cfg := NewConfig()

	if cfg.APIVersion != "v1" {
		t.Errorf("NewConfig() APIVersion = %v, want v1", cfg.APIVersion)
	}
	if cfg.Kind != "Config" {
		t.Errorf("NewConfig() Kind = %v, want Config", cfg.Kind)
	}
	if cfg.CurrentContext != "" {
		t.Errorf("NewConfig() CurrentContext = %v, want empty", cfg.CurrentContext)
	}
	if len(cfg.Contexts) != 0 {
		t.Errorf("NewConfig() Contexts = %v, want empty slice", cfg.Contexts)
	}
	if len(cfg.Credentials) != 0 {
		t.Errorf("NewConfig() Credentials = %v, want empty slice", cfg.Credentials)
	}
}

func TestConfig_GetContext(t *testing.T) {
	cfg := NewConfig()
	accountUUID := "test-uuid"
	credRef := "test-cred"
	if err := cfg.SetContext("test", &accountUUID, &credRef); err != nil {
		t.Fatalf("SetContext() unexpected error: %v", err)
	}

	// Test found
	ctx := cfg.GetContext("test")
	if ctx == nil {
		t.Fatal("GetContext() returned nil for existing context")
	}
	if ctx.AccountUUID != "test-uuid" {
		t.Errorf("GetContext() AccountUUID = %v, want test-uuid", ctx.AccountUUID)
	}
	if ctx.CredentialsRef != "test-cred" {
		t.Errorf("GetContext() CredentialsRef = %v, want test-cred", ctx.CredentialsRef)
	}

	// Test not found
	ctx = cfg.GetContext("nonexistent")
	if ctx != nil {
		t.Error("GetContext() should return nil for nonexistent context")
	}
}

func TestConfig_GetCurrentContext(t *testing.T) {
	cfg := NewConfig()

	// Test empty current context
	ctx := cfg.GetCurrentContext()
	if ctx != nil {
		t.Error("GetCurrentContext() should return nil when no context is set")
	}

	// Add and set context
	accountUUID := "test-uuid"
	if err := cfg.SetContext("test", &accountUUID, nil); err != nil {
		t.Fatalf("SetContext() unexpected error: %v", err)
	}
	if err := cfg.UseContext("test"); err != nil {
		t.Fatalf("UseContext() unexpected error: %v", err)
	}

	ctx = cfg.GetCurrentContext()
	if ctx == nil {
		t.Fatal("GetCurrentContext() returned nil for set context")
	}
	if ctx.AccountUUID != "test-uuid" {
		t.Errorf("GetCurrentContext() AccountUUID = %v, want test-uuid", ctx.AccountUUID)
	}
}

func TestConfig_SetCredential(t *testing.T) {
	cfg := NewConfig()

	// Add new credential
	cfg.SetCredential("test", "client-id", "client-secret")

	cred := cfg.GetCredential("test")
	if cred == nil {
		t.Fatal("GetCredential() returned nil for existing credential")
	}
	if cred.ClientID != "client-id" {
		t.Errorf("GetCredential() ClientID = %v, want client-id", cred.ClientID)
	}
	if cred.ClientSecret != "client-secret" {
		t.Errorf("GetCredential() ClientSecret = %v, want client-secret", cred.ClientSecret)
	}

	// Update existing credential
	cfg.SetCredential("test", "new-id", "new-secret")
	cred = cfg.GetCredential("test")
	if cred.ClientID != "new-id" {
		t.Errorf("GetCredential() after update ClientID = %v, want new-id", cred.ClientID)
	}
}

func TestConfig_DeleteContext(t *testing.T) {
	cfg := NewConfig()
	accountUUID := "test-uuid"
	if err := cfg.SetContext("test", &accountUUID, nil); err != nil {
		t.Fatalf("SetContext() unexpected error: %v", err)
	}
	if err := cfg.UseContext("test"); err != nil {
		t.Fatalf("UseContext() unexpected error: %v", err)
	}

	// Delete existing
	if !cfg.DeleteContext("test") {
		t.Error("DeleteContext() returned false for existing context")
	}
	if cfg.GetContext("test") != nil {
		t.Error("GetContext() should return nil after deletion")
	}
	if cfg.CurrentContext != "" {
		t.Error("CurrentContext should be cleared when active context is deleted")
	}

	// Delete nonexistent
	if cfg.DeleteContext("nonexistent") {
		t.Error("DeleteContext() returned true for nonexistent context")
	}
}

func TestConfig_DeleteCredential(t *testing.T) {
	cfg := NewConfig()
	cfg.SetCredential("test", "id", "secret")

	// Delete existing
	if !cfg.DeleteCredential("test") {
		t.Error("DeleteCredential() returned false for existing credential")
	}
	if cfg.GetCredential("test") != nil {
		t.Error("GetCredential() should return nil after deletion")
	}

	// Delete nonexistent
	if cfg.DeleteCredential("nonexistent") {
		t.Error("DeleteCredential() returned true for nonexistent credential")
	}
}

func TestConfig_UseContext(t *testing.T) {
	cfg := NewConfig()

	// Test error for nonexistent context
	err := cfg.UseContext("nonexistent")
	if err == nil {
		t.Error("UseContext() should return error for nonexistent context")
	}

	// Add and use context
	accountUUID := "test-uuid"
	if err := cfg.SetContext("test", &accountUUID, nil); err != nil {
		t.Fatalf("SetContext() unexpected error: %v", err)
	}
	err = cfg.UseContext("test")
	if err != nil {
		t.Errorf("UseContext() unexpected error: %v", err)
	}
	if cfg.CurrentContext != "test" {
		t.Errorf("CurrentContext = %v, want test", cfg.CurrentContext)
	}
}

func TestConfig_GetContextNames(t *testing.T) {
	cfg := NewConfig()

	// Empty config
	names := cfg.GetContextNames()
	if len(names) != 0 {
		t.Errorf("GetContextNames() = %v, want empty slice", names)
	}

	// Add contexts
	uuid1 := "uuid1"
	uuid2 := "uuid2"
	if err := cfg.SetContext("ctx1", &uuid1, nil); err != nil {
		t.Fatalf("SetContext() unexpected error: %v", err)
	}
	if err := cfg.SetContext("ctx2", &uuid2, nil); err != nil {
		t.Fatalf("SetContext() unexpected error: %v", err)
	}

	names = cfg.GetContextNames()
	if len(names) != 2 {
		t.Errorf("GetContextNames() len = %d, want 2", len(names))
	}
}

func TestConfig_GetCredentialNames(t *testing.T) {
	cfg := NewConfig()

	// Empty config
	names := cfg.GetCredentialNames()
	if len(names) != 0 {
		t.Errorf("GetCredentialNames() = %v, want empty slice", names)
	}

	// Add credentials
	cfg.SetCredential("cred1", "id1", "secret1")
	cfg.SetCredential("cred2", "id2", "secret2")

	names = cfg.GetCredentialNames()
	if len(names) != 2 {
		t.Errorf("GetCredentialNames() len = %d, want 2", len(names))
	}
}

func TestMaskSecret(t *testing.T) {
	tests := []struct {
		name   string
		secret string
		want   string
	}{
		{"short secret", "abc", "****"},
		{"8 char secret", "12345678", "****"},
		{"longer secret", "1234567890123456", "1234********3456"},
		{"very long secret", "abcdefghijklmnopqrstuvwxyz", "abcd******************wxyz"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MaskSecret(tt.secret)
			if got != tt.want {
				t.Errorf("MaskSecret(%q) = %q, want %q", tt.secret, got, tt.want)
			}
		})
	}
}

func TestConfig_Validate(t *testing.T) {
	// Valid config
	cfg := NewConfig()
	if err := cfg.Validate(); err != nil {
		t.Errorf("Validate() unexpected error for empty config: %v", err)
	}

	// Valid config with context and credential
	uuid := "test-uuid"
	credRef := "test-cred"
	cfg.SetCredential("test-cred", "id", "secret")
	if err := cfg.SetContext("test", &uuid, &credRef); err != nil {
		t.Fatalf("SetContext() unexpected error: %v", err)
	}
	if err := cfg.UseContext("test"); err != nil {
		t.Fatalf("UseContext() unexpected error: %v", err)
	}

	if err := cfg.Validate(); err != nil {
		t.Errorf("Validate() unexpected error for valid config: %v", err)
	}

	// Invalid - current context not found
	cfg2 := NewConfig()
	cfg2.CurrentContext = "nonexistent"
	if err := cfg2.Validate(); err == nil {
		t.Error("Validate() should return error for nonexistent current context")
	}

	// Invalid - credential ref not found
	cfg3 := NewConfig()
	uuid3 := "uuid"
	badCredRef := "nonexistent"
	if err := cfg3.SetContext("test", &uuid3, &badCredRef); err != nil {
		t.Fatalf("SetContext() unexpected error: %v", err)
	}
	if err := cfg3.Validate(); err == nil {
		t.Error("Validate() should return error for nonexistent credential ref")
	}
}
