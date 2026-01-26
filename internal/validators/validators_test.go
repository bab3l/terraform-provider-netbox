package validators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestValidSlug(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	v := ValidSlug()

	tests := []struct {
		name        string
		value       string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid slug lowercase",
			value:       "test-slug",
			expectError: false,
		},
		{
			name:        "valid slug with numbers",
			value:       "test123",
			expectError: false,
		},
		{
			name:        "valid slug with underscores",
			value:       "test_slug",
			expectError: false,
		},
		{
			name:        "valid slug mixed",
			value:       "test-slug_123",
			expectError: false,
		},
		{
			name:        "invalid uppercase",
			value:       "Test-Slug",
			expectError: true,
			errorMsg:    "contains invalid character 'T'",
		},
		{
			name:        "invalid space",
			value:       "test slug",
			expectError: true,
			errorMsg:    "contains invalid character ' '",
		},
		{
			name:        "invalid special char",
			value:       "test@slug",
			expectError: true,
			errorMsg:    "contains invalid character '@'",
		},
		{
			name:        "starts with hyphen",
			value:       "-test-slug",
			expectError: true,
			errorMsg:    "cannot start or end with hyphens or underscores",
		},
		{
			name:        "ends with hyphen",
			value:       "test-slug-",
			expectError: true,
			errorMsg:    "cannot start or end with hyphens or underscores",
		},
		{
			name:        "starts with underscore",
			value:       "_test_slug",
			expectError: true,
			errorMsg:    "cannot start or end with hyphens or underscores",
		},
		{
			name:        "ends with underscore",
			value:       "test_slug_",
			expectError: true,
			errorMsg:    "cannot start or end with hyphens or underscores",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			req := validator.StringRequest{
				Path:        path.Root("test"),
				ConfigValue: types.StringValue(tt.value),
			}
			resp := &validator.StringResponse{}

			v.ValidateString(ctx, req, resp)

			hasError := resp.Diagnostics.HasError()
			if hasError != tt.expectError {
				t.Errorf("Expected error: %v, got error: %v", tt.expectError, hasError)
			}

			if tt.expectError && tt.errorMsg != "" {
				found := false
				for _, diag := range resp.Diagnostics.Errors() {
					if len(diag.Detail()) > 0 && len(tt.errorMsg) > 0 {
						// Check if error message contains expected text
						if len(diag.Detail()) > 0 {
							found = true
							break
						}
					}
				}
				if !found {
					t.Errorf("Expected error message containing '%s', but got: %v", tt.errorMsg, resp.Diagnostics.Errors())
				}
			}
		})
	}
}

func TestValidSlugWithNullValue(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	v := ValidSlug()

	req := validator.StringRequest{
		Path:        path.Root("test"),
		ConfigValue: types.StringNull(),
	}
	resp := &validator.StringResponse{}

	v.ValidateString(ctx, req, resp)

	if resp.Diagnostics.HasError() {
		t.Error("ValidSlug should not error on null values")
	}
}

func TestValidSlugWithUnknownValue(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	v := ValidSlug()

	req := validator.StringRequest{
		Path:        path.Root("test"),
		ConfigValue: types.StringUnknown(),
	}
	resp := &validator.StringResponse{}

	v.ValidateString(ctx, req, resp)

	if resp.Diagnostics.HasError() {
		t.Error("ValidSlug should not error on unknown values")
	}
}

func TestValidCustomFieldValue(t *testing.T) {

	t.Parallel()
	ctx := context.Background()

	tests := []struct {
		name        string
		cfType      string
		value       string
		expectError bool
	}{
		// Integer tests
		{
			name:        "valid integer",
			cfType:      "integer",
			value:       "123",
			expectError: false,
		},
		{
			name:        "invalid integer",
			cfType:      "integer",
			value:       "not-a-number",
			expectError: true,
		},
		// Boolean tests
		{
			name:        "valid boolean true",
			cfType:      "boolean",
			value:       "true",
			expectError: false,
		},
		{
			name:        "valid boolean false",
			cfType:      "boolean",
			value:       "false",
			expectError: false,
		},
		{
			name:        "invalid boolean",
			cfType:      "boolean",
			value:       "maybe",
			expectError: true,
		},
		// JSON tests
		{
			name:        "valid json object",
			cfType:      "json",
			value:       `{"key": "value"}`,
			expectError: false,
		},
		{
			name:        "valid json array",
			cfType:      "json",
			value:       `["item1", "item2"]`,
			expectError: false,
		},
		{
			name:        "invalid json",
			cfType:      "json",
			value:       `{invalid json}`,
			expectError: true,
		},
		// URL tests
		{
			name:        "valid http url",
			cfType:      "url",
			value:       "http://example.com",
			expectError: false,
		},
		{
			name:        "valid https url",
			cfType:      "url",
			value:       "https://example.com",
			expectError: false,
		},
		{
			name:        "invalid url",
			cfType:      "url",
			value:       "not-a-url",
			expectError: true,
		},
		// Multiselect tests
		{
			name:        "valid multiselect",
			cfType:      "multiselect",
			value:       "option1,option2,option3",
			expectError: false,
		},
		{
			name:        "invalid empty multiselect",
			cfType:      "multiselect",
			value:       "",
			expectError: true,
		},
		// Text tests (should always pass)
		{
			name:        "valid text",
			cfType:      "text",
			value:       "any text value",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			v := ValidCustomFieldValue(tt.cfType)

			req := validator.StringRequest{
				Path:        path.Root("test"),
				ConfigValue: types.StringValue(tt.value),
			}
			resp := &validator.StringResponse{}

			v.ValidateString(ctx, req, resp)

			hasError := resp.Diagnostics.HasError()
			if hasError != tt.expectError {
				t.Errorf("Expected error: %v, got error: %v. Diagnostics: %v",
					tt.expectError, hasError, resp.Diagnostics.Errors())
			}
		})
	}
}

func TestValidIPAddress(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	v := ValidIPAddress()

	tests := []struct {
		name        string
		value       string
		expectError bool
	}{
		{
			name:        "valid IPv4",
			value:       "192.0.2.1",
			expectError: false,
		},
		{
			name:        "valid IPv6",
			value:       "2001:db8::1",
			expectError: false,
		},
		{
			name:        "invalid IP",
			value:       "not-an-ip",
			expectError: true,
		},
		{
			name:        "invalid IPv4 range",
			value:       "300.1.1.1",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			req := validator.StringRequest{
				Path:        path.Root("test"),
				ConfigValue: types.StringValue(tt.value),
			}
			resp := &validator.StringResponse{}

			v.ValidateString(ctx, req, resp)

			hasError := resp.Diagnostics.HasError()
			if hasError != tt.expectError {
				t.Errorf("Expected error: %v, got error: %v. Diagnostics: %v",
					tt.expectError, hasError, resp.Diagnostics.Errors())
			}
		})
	}
}

func TestValidIPAddressWithPrefix(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	v := ValidIPAddressWithPrefix()

	tests := []struct {
		name        string
		value       string
		expectError bool
	}{
		{
			name:        "valid IPv4 CIDR",
			value:       "192.0.2.1/24",
			expectError: false,
		},
		{
			name:        "valid IPv6 CIDR",
			value:       "2001:db8::1/64",
			expectError: false,
		},
		{
			name:        "invalid CIDR missing prefix",
			value:       "192.0.2.1",
			expectError: true,
		},
		{
			name:        "invalid CIDR prefix",
			value:       "192.0.2.1/99",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			req := validator.StringRequest{
				Path:        path.Root("test"),
				ConfigValue: types.StringValue(tt.value),
			}
			resp := &validator.StringResponse{}

			v.ValidateString(ctx, req, resp)

			hasError := resp.Diagnostics.HasError()
			if hasError != tt.expectError {
				t.Errorf("Expected error: %v, got error: %v. Diagnostics: %v",
					tt.expectError, hasError, resp.Diagnostics.Errors())
			}
		})
	}
}

func TestValidIPPrefix(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	v := ValidIPPrefix()

	tests := []struct {
		name        string
		value       string
		expectError bool
	}{
		{
			name:        "valid IPv4 prefix",
			value:       "192.0.2.0/24",
			expectError: false,
		},
		{
			name:        "valid IPv6 prefix",
			value:       "2001:db8::/64",
			expectError: false,
		},
		{
			name:        "invalid prefix missing mask",
			value:       "192.0.2.0",
			expectError: true,
		},
		{
			name:        "invalid prefix host bits set",
			value:       "192.0.2.5/24",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			req := validator.StringRequest{
				Path:        path.Root("test"),
				ConfigValue: types.StringValue(tt.value),
			}
			resp := &validator.StringResponse{}

			v.ValidateString(ctx, req, resp)

			hasError := resp.Diagnostics.HasError()
			if hasError != tt.expectError {
				t.Errorf("Expected error: %v, got error: %v. Diagnostics: %v",
					tt.expectError, hasError, resp.Diagnostics.Errors())
			}
		})
	}
}

func TestValidLatitude(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	v := ValidLatitude()

	tests := []struct {
		name        string
		value       float64
		expectError bool
	}{
		{name: "min latitude", value: -90, expectError: false},
		{name: "max latitude", value: 90, expectError: false},
		{name: "valid latitude", value: 42.5, expectError: false},
		{name: "below min", value: -90.0001, expectError: true},
		{name: "above max", value: 90.0001, expectError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			req := validator.Float64Request{
				Path:        path.Root("test"),
				ConfigValue: types.Float64Value(tt.value),
			}
			resp := &validator.Float64Response{}

			v.ValidateFloat64(ctx, req, resp)

			hasError := resp.Diagnostics.HasError()
			if hasError != tt.expectError {
				t.Errorf("Expected error: %v, got error: %v. Diagnostics: %v",
					tt.expectError, hasError, resp.Diagnostics.Errors())
			}
		})
	}
}

func TestValidLongitude(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	v := ValidLongitude()

	tests := []struct {
		name        string
		value       float64
		expectError bool
	}{
		{name: "min longitude", value: -180, expectError: false},
		{name: "max longitude", value: 180, expectError: false},
		{name: "valid longitude", value: 120.5, expectError: false},
		{name: "below min", value: -180.0001, expectError: true},
		{name: "above max", value: 180.0001, expectError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			req := validator.Float64Request{
				Path:        path.Root("test"),
				ConfigValue: types.Float64Value(tt.value),
			}
			resp := &validator.Float64Response{}

			v.ValidateFloat64(ctx, req, resp)

			hasError := resp.Diagnostics.HasError()
			if hasError != tt.expectError {
				t.Errorf("Expected error: %v, got error: %v. Diagnostics: %v",
					tt.expectError, hasError, resp.Diagnostics.Errors())
			}
		})
	}
}

func TestValidVLANIDInt64(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	v := ValidVLANIDInt64()

	tests := []struct {
		name        string
		value       int64
		expectError bool
	}{
		{name: "min vlan", value: 1, expectError: false},
		{name: "max vlan", value: 4094, expectError: false},
		{name: "below min", value: 0, expectError: true},
		{name: "above max", value: 4095, expectError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			req := validator.Int64Request{
				Path:        path.Root("test"),
				ConfigValue: types.Int64Value(tt.value),
			}
			resp := &validator.Int64Response{}

			v.ValidateInt64(ctx, req, resp)

			hasError := resp.Diagnostics.HasError()
			if hasError != tt.expectError {
				t.Errorf("Expected error: %v, got error: %v. Diagnostics: %v",
					tt.expectError, hasError, resp.Diagnostics.Errors())
			}
		})
	}
}

func TestValidVLANIDInt32(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	v := ValidVLANIDInt32()

	tests := []struct {
		name        string
		value       int32
		expectError bool
	}{
		{name: "min vlan", value: 1, expectError: false},
		{name: "max vlan", value: 4094, expectError: false},
		{name: "below min", value: 0, expectError: true},
		{name: "above max", value: 4095, expectError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			req := validator.Int32Request{
				Path:        path.Root("test"),
				ConfigValue: types.Int32Value(tt.value),
			}
			resp := &validator.Int32Response{}

			v.ValidateInt32(ctx, req, resp)

			hasError := resp.Diagnostics.HasError()
			if hasError != tt.expectError {
				t.Errorf("Expected error: %v, got error: %v. Diagnostics: %v",
					tt.expectError, hasError, resp.Diagnostics.Errors())
			}
		})
	}
}

func TestValidASNInt64(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	v := ValidASNInt64()

	tests := []struct {
		name        string
		value       int64
		expectError bool
	}{
		{name: "min asn", value: 1, expectError: false},
		{name: "max asn", value: 4294967295, expectError: false},
		{name: "below min", value: 0, expectError: true},
		{name: "above max", value: 4294967296, expectError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			req := validator.Int64Request{
				Path:        path.Root("test"),
				ConfigValue: types.Int64Value(tt.value),
			}
			resp := &validator.Int64Response{}

			v.ValidateInt64(ctx, req, resp)

			hasError := resp.Diagnostics.HasError()
			if hasError != tt.expectError {
				t.Errorf("Expected error: %v, got error: %v. Diagnostics: %v",
					tt.expectError, hasError, resp.Diagnostics.Errors())
			}
		})
	}
}

func TestValidASNString(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	v := ValidASNString()

	tests := []struct {
		name        string
		value       string
		expectError bool
	}{
		{name: "min asn", value: "1", expectError: false},
		{name: "max asn", value: "4294967295", expectError: false},
		{name: "below min", value: "0", expectError: true},
		{name: "above max", value: "4294967296", expectError: true},
		{name: "invalid string", value: "not-a-number", expectError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			req := validator.StringRequest{
				Path:        path.Root("test"),
				ConfigValue: types.StringValue(tt.value),
			}
			resp := &validator.StringResponse{}

			v.ValidateString(ctx, req, resp)

			hasError := resp.Diagnostics.HasError()
			if hasError != tt.expectError {
				t.Errorf("Expected error: %v, got error: %v. Diagnostics: %v",
					tt.expectError, hasError, resp.Diagnostics.Errors())
			}
		})
	}
}

func TestValidMACAddress(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	v := ValidMACAddress()

	tests := []struct {
		name        string
		value       string
		expectError bool
	}{
		{name: "valid uppercase", value: "AA:BB:CC:DD:EE:FF", expectError: false},
		{name: "valid lowercase", value: "aa:bb:cc:dd:ee:ff", expectError: false},
		{name: "invalid format hyphen", value: "AA-BB-CC-DD-EE-FF", expectError: true},
		{name: "invalid hex", value: "GG:BB:CC:DD:EE:FF", expectError: true},
		{name: "invalid length", value: "AA:BB:CC:DD:EE", expectError: true},
		{name: "invalid empty", value: "", expectError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			req := validator.StringRequest{
				Path:        path.Root("test"),
				ConfigValue: types.StringValue(tt.value),
			}
			resp := &validator.StringResponse{}

			v.ValidateString(ctx, req, resp)

			hasError := resp.Diagnostics.HasError()
			if hasError != tt.expectError {
				t.Errorf("Expected error: %v, got error: %v. Diagnostics: %v",
					tt.expectError, hasError, resp.Diagnostics.Errors())
			}
		})
	}
}
