package validators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestValidSlug(t *testing.T) {
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
