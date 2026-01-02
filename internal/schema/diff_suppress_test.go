package schema

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// TestSuppressReferenceEquivalent tests the core diff suppression logic
// for reference fields that can accept names, slugs, or IDs.
func TestSuppressReferenceEquivalent(t *testing.T) {
	tests := []struct {
		name           string
		key            string
		oldValue       string
		newValue       string
		resourceType   string
		shouldSuppress bool
		description    string
	}{
		// Exact matches - should always suppress
		{
			name:           "exact_id_match",
			key:            "tenant",
			oldValue:       "7",
			newValue:       "7",
			resourceType:   "tenant",
			shouldSuppress: true,
			description:    "Identical ID values should suppress diff",
		},
		{
			name:           "exact_slug_match",
			key:            "tenant",
			oldValue:       "production-environment",
			newValue:       "production-environment",
			resourceType:   "tenant",
			shouldSuppress: true,
			description:    "Identical slug values should suppress diff",
		},
		{
			name:           "exact_name_match",
			key:            "tenant",
			oldValue:       "Production Environment",
			newValue:       "Production Environment",
			resourceType:   "tenant",
			shouldSuppress: true,
			description:    "Identical name values should suppress diff",
		},

		// ID vs Slug equivalency - should suppress when they refer to same object
		{
			name:           "id_to_slug_equivalent",
			key:            "tenant",
			oldValue:       "7",
			newValue:       "production-environment",
			resourceType:   "tenant",
			shouldSuppress: true,
			description:    "ID and equivalent slug should suppress diff",
		},
		{
			name:           "slug_to_id_equivalent",
			key:            "tenant",
			oldValue:       "production-environment",
			newValue:       "7",
			resourceType:   "tenant",
			shouldSuppress: true,
			description:    "Slug and equivalent ID should suppress diff",
		},

		// ID vs Name equivalency - should suppress when they refer to same object
		{
			name:           "id_to_name_equivalent",
			key:            "tenant",
			oldValue:       "7",
			newValue:       "Production Environment",
			resourceType:   "tenant",
			shouldSuppress: true,
			description:    "ID and equivalent name should suppress diff",
		},
		{
			name:           "name_to_id_equivalent",
			key:            "tenant",
			oldValue:       "Production Environment",
			newValue:       "7",
			resourceType:   "tenant",
			shouldSuppress: true,
			description:    "Name and equivalent ID should suppress diff",
		},

		// Slug vs Name equivalency - should suppress when they refer to same object
		{
			name:           "slug_to_name_equivalent",
			key:            "tenant",
			oldValue:       "production-environment",
			newValue:       "Production Environment",
			resourceType:   "tenant",
			shouldSuppress: true,
			description:    "Slug and equivalent name should suppress diff",
		},
		{
			name:           "name_to_slug_equivalent",
			key:            "tenant",
			oldValue:       "Production Environment",
			newValue:       "production-environment",
			resourceType:   "tenant",
			shouldSuppress: true,
			description:    "Name and equivalent slug should suppress diff",
		},

		// Non-equivalent values - should NOT suppress
		{
			name:           "different_ids",
			key:            "tenant",
			oldValue:       "7",
			newValue:       "12",
			resourceType:   "tenant",
			shouldSuppress: false,
			description:    "Different IDs should not suppress diff",
		},
		{
			name:           "different_slugs",
			key:            "tenant",
			oldValue:       "production-environment",
			newValue:       "staging-environment",
			resourceType:   "tenant",
			shouldSuppress: false,
			description:    "Different slugs should not suppress diff",
		},
		{
			name:           "different_names",
			key:            "tenant",
			oldValue:       "Production Environment",
			newValue:       "Staging Environment",
			resourceType:   "tenant",
			shouldSuppress: false,
			description:    "Different names should not suppress diff",
		},
		{
			name:           "id_to_different_slug",
			key:            "tenant",
			oldValue:       "7",
			newValue:       "staging-environment",
			resourceType:   "tenant",
			shouldSuppress: false,
			description:    "ID and non-equivalent slug should not suppress diff",
		},

		// Different resource types to test versatility
		{
			name:           "site_id_to_slug",
			key:            "site",
			oldValue:       "15",
			newValue:       "datacenter-east",
			resourceType:   "site",
			shouldSuppress: true,
			description:    "Site ID and equivalent slug should suppress diff",
		},
		{
			name:           "device_type_id_to_name",
			key:            "device_type",
			oldValue:       "9",
			newValue:       "PowerEdge R640",
			resourceType:   "device_type",
			shouldSuppress: true,
			description:    "Device type ID and equivalent name should suppress diff",
		},
		{
			name:           "platform_slug_to_name",
			key:            "platform",
			oldValue:       "ubuntu-2204",
			newValue:       "Ubuntu 22.04 LTS",
			resourceType:   "platform",
			shouldSuppress: true,
			description:    "Platform slug and equivalent name should suppress diff",
		},

		// Edge cases
		{
			name:           "empty_values",
			key:            "tenant",
			oldValue:       "",
			newValue:       "",
			resourceType:   "tenant",
			shouldSuppress: true,
			description:    "Empty values should suppress diff",
		},
		{
			name:           "empty_to_value",
			key:            "tenant",
			oldValue:       "",
			newValue:       "7",
			resourceType:   "tenant",
			shouldSuppress: false,
			description:    "Empty to value should not suppress diff",
		},
		{
			name:           "value_to_empty",
			key:            "tenant",
			oldValue:       "7",
			newValue:       "",
			resourceType:   "tenant",
			shouldSuppress: false,
			description:    "Value to empty should not suppress diff",
		},

		// Invalid/non-existent references - should not suppress
		{
			name:           "invalid_id",
			key:            "tenant",
			oldValue:       "99999",
			newValue:       "production-environment",
			resourceType:   "tenant",
			shouldSuppress: false,
			description:    "Invalid ID should not suppress diff",
		},
		{
			name:           "invalid_slug",
			key:            "tenant",
			oldValue:       "7",
			newValue:       "non-existent-tenant",
			resourceType:   "tenant",
			shouldSuppress: false,
			description:    "Invalid slug should not suppress diff",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock ResourceData for testing
			d := &schema.ResourceData{}

			result := suppressReferenceEquivalent(tt.key, tt.oldValue, tt.newValue, d)

			if result != tt.shouldSuppress {
				t.Errorf("suppressReferenceEquivalent(%s, %s, %s) = %v, want %v\nDescription: %s",
					tt.key, tt.oldValue, tt.newValue, result, tt.shouldSuppress, tt.description)
			}
		})
	}
}

// TestResourceTypeDetection tests the logic for determining resource type from attribute key.
func TestResourceTypeDetection(t *testing.T) {
	tests := []struct {
		attributeKey         string
		expectedResourceType string
		description          string
	}{
		{
			attributeKey:         "tenant",
			expectedResourceType: "tenant",
			description:          "Simple tenant reference should map to tenant resource type",
		},
		{
			attributeKey:         "site",
			expectedResourceType: "site",
			description:          "Simple site reference should map to site resource type",
		},
		{
			attributeKey:         "device_type",
			expectedResourceType: "device_type",
			description:          "Underscore resource type should map correctly",
		},
		{
			attributeKey:         "parent",
			expectedResourceType: "unknown",
			description:          "Generic parent reference should return unknown (requires context)",
		},
		{
			attributeKey:         "group",
			expectedResourceType: "unknown",
			description:          "Generic group reference should return unknown (requires context)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.attributeKey, func(t *testing.T) {
			result := getResourceTypeFromAttribute(tt.attributeKey)

			if result != tt.expectedResourceType {
				t.Errorf("getResourceTypeFromAttribute(%s) = %s, want %s\nDescription: %s",
					tt.attributeKey, result, tt.expectedResourceType, tt.description)
			}
		})
	}
}

// Benchmark tests to ensure performance is acceptable.
func BenchmarkSuppressReferenceEquivalent(b *testing.B) {
	d := &schema.ResourceData{}

	for i := 0; i < b.N; i++ {
		suppressReferenceEquivalent("tenant", "7", "production-environment", d)
	}
}
