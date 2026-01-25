package schema

import (
	"context"
	"testing"

	"github.com/bab3l/go-netbox"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestReferenceResolver_ResolveToID tests the reference resolver functionality.
func TestReferenceResolver_ResolveToID(t *testing.T) {
	// Note: This test requires a running NetBox instance for integration testing
	// For unit testing, we'll mock the client behavior

	tests := []struct {
		name         string
		value        string
		resourceType string
		expectError  bool
		description  string
	}{
		{
			name:         "empty_value",
			value:        "",
			resourceType: "tenant",
			expectError:  false,
			description:  "Empty values should return 0 without error",
		},
		{
			name:         "numeric_id",
			value:        "123",
			resourceType: "tenant",
			expectError:  false,
			description:  "Numeric values should parse as ID directly",
		},
		{
			name:         "unknown_resource_type",
			value:        "test-value",
			resourceType: "unknown_type",
			expectError:  false,
			description:  "Unknown resource types should return 0 without error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock client - in real usage this would be a real NetBox client
			var client *netbox.APIClient = nil
			resolver := NewReferenceResolver(client)

			ctx := context.Background()
			id, err := resolver.ResolveToID(ctx, tt.value, tt.resourceType)

			if tt.expectError {
				require.Error(t, err, tt.description)
			} else {
				require.NoError(t, err, tt.description)
				switch tt.value {
				case "":
					assert.Equal(t, int32(0), id, "Empty value should resolve to 0")
				case "123":
					assert.Equal(t, int32(123), id, "Numeric value should parse directly")
				}
			}
		})
	}
}

// TestGetResourceTypeFromAttribute tests attribute path parsing.
func TestGetResourceTypeFromAttribute(t *testing.T) {
	tests := []struct {
		name          string
		attributePath string
		expected      string
		description   string
	}{
		{
			name:          "simple_tenant",
			attributePath: "tenant",
			expected:      "tenant",
			description:   "Simple tenant attribute should map correctly",
		},
		{
			name:          "nested_site",
			attributePath: "root.nested.site",
			expected:      "site",
			description:   "Nested site attribute should extract correctly",
		},
		{
			name:          "device_role",
			attributePath: "role",
			expected:      "device_role",
			description:   "Role should map to device_role for API compatibility",
		},
		{
			name:          "unknown_attribute",
			attributePath: "unknown_field",
			expected:      "",
			description:   "Unknown attributes should return empty string",
		},
		{
			name:          "rack_specific",
			attributePath: "config.rack_role",
			expected:      "rack_role",
			description:   "Rack-specific attributes should be detected",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getResourceTypeFromAttribute(tt.attributePath)
			assert.Equal(t, tt.expected, result, tt.description)
		})
	}
}

// TestReferenceEquivalencePlanModifier tests the plan modifier behavior.
func TestReferenceEquivalencePlanModifier(t *testing.T) {
	modifier := ReferenceEquivalencePlanModifier()

	// Test description methods
	ctx := context.Background()
	description := modifier.Description(ctx)
	assert.NotEmpty(t, description, "Plan modifier should have a description")

	markdownDesc := modifier.MarkdownDescription(ctx)
	assert.Equal(t, description, markdownDesc, "Markdown description should match regular description")
}

// TestSuppressReferenceEquivalent_NoClient tests behavior when no NetBox client is available.
func TestSuppressReferenceEquivalent_NoClient(t *testing.T) {
	// Create a plan modifier request without NetBox client in context
	ctx := context.Background()
	req := planmodifier.StringRequest{
		Path:       path.Root("tenant"),
		StateValue: types.StringValue("old-tenant"),
		PlanValue:  types.StringValue("new-tenant"),
	}
	resp := &planmodifier.StringResponse{}

	modifier := ReferenceEquivalencePlanModifier()
	modifier.PlanModifyString(ctx, req, resp)

	// Without client, no suppression should occur - plan value should not be modified
	assert.True(t, resp.PlanValue.IsNull(), "Plan value should not be modified without client")
}

// TestSuppressReferenceEquivalent_IdenticalValues tests that identical values are not modified.
func TestSuppressReferenceEquivalent_IdenticalValues(t *testing.T) {
	ctx := context.Background()
	req := planmodifier.StringRequest{
		Path:       path.Root("tenant"),
		StateValue: types.StringValue("same-value"),
		PlanValue:  types.StringValue("same-value"),
	}
	resp := &planmodifier.StringResponse{}

	modifier := ReferenceEquivalencePlanModifier()
	modifier.PlanModifyString(ctx, req, resp)

	// Identical values should not trigger any modification
	assert.True(t, resp.PlanValue.IsNull(), "Plan value should remain unmodified for identical values")
}

// TestSuppressReferenceEquivalent_NullValues tests handling of null values.
func TestSuppressReferenceEquivalent_NullValues(t *testing.T) {
	tests := []struct {
		name        string
		stateValue  types.String
		planValue   types.String
		description string
	}{
		{
			name:        "null_state",
			stateValue:  types.StringNull(),
			planValue:   types.StringValue("new-value"),
			description: "Null state value should not trigger suppression",
		},
		{
			name:        "unknown_plan",
			stateValue:  types.StringValue("old-value"),
			planValue:   types.StringUnknown(),
			description: "Unknown plan value should not trigger suppression",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			req := planmodifier.StringRequest{
				Path:       path.Root("tenant"),
				StateValue: tt.stateValue,
				PlanValue:  tt.planValue,
			}
			resp := &planmodifier.StringResponse{}

			modifier := ReferenceEquivalencePlanModifier()
			modifier.PlanModifyString(ctx, req, resp)

			// Should not modify the response
			assert.True(t, resp.PlanValue.IsNull(), tt.description)
		})
	}
}

// BenchmarkGetResourceTypeFromAttribute benchmarks attribute type detection.
func BenchmarkGetResourceTypeFromAttribute(b *testing.B) {
	attributePaths := []string{
		"tenant",
		"site",
		"device_type",
		"platform",
		"rack.nested.role",
		"config.location",
		"unknown_attribute",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, path := range attributePaths {
			getResourceTypeFromAttribute(path)
		}
	}
}
