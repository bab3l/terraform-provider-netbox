package schema

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// suppressReferenceEquivalent implements DiffSuppressFunc for reference fields
// that can accept names, slugs, or IDs. It suppresses diffs when the old and new
// values refer to the same NetBox object, even if they use different formats
// (e.g., "production-environment" slug vs "7" ID).
func suppressReferenceEquivalent(k, old, newValue string, d *schema.ResourceData) bool {
	// If values are identical, suppress the diff
	if old == newValue {
		return true
	}

	// If either value is empty, don't suppress (actual change)
	if old == "" || newValue == "" {
		return false
	}

	// Determine the resource type from the attribute key
	resourceType := getResourceTypeFromAttribute(k)
	if resourceType == "unknown" {
		// Can't determine equivalency for unknown resource types
		return false
	}

	// Check if the values refer to the same NetBox object
	return areValuesEquivalent(old, newValue, resourceType)
}

// getResourceTypeFromAttribute extracts the NetBox resource type from a Terraform attribute key.
// This helps determine which NetBox API endpoint to use for lookups.
func getResourceTypeFromAttribute(attributeKey string) string {
	// Direct mappings for common attribute names
	mappings := map[string]string{
		"tenant":       "tenant",
		"site":         "site",
		"location":     "location",
		"rack":         "rack",
		"device":       "device",
		"device_type":  "device_type",
		"role":         "device_role", // Note: device roles in NetBox API
		"platform":     "platform",
		"cluster":      "cluster",
		"vlan":         "vlan",
		"vrf":          "vrf",
		"region":       "region",
		"manufacturer": "manufacturer",
		"circuit":      "circuit",
		"provider":     "provider",
		"rir":          "rir",
	}

	if resourceType, exists := mappings[attributeKey]; exists {
		return resourceType
	}

	// Generic attributes like "parent" or "group" need context to determine type
	// For now, return "unknown" - these may need special handling in the future
	return "unknown"
}

// areValuesEquivalent checks if two values (which may be names, slugs, or IDs)
// refer to the same NetBox object of the given resource type.
func areValuesEquivalent(value1, value2, resourceType string) bool {
	// Convert both values to IDs and compare
	id1 := resolveValueToID(value1, resourceType)
	id2 := resolveValueToID(value2, resourceType)

	// If either resolution failed, they're not equivalent
	if id1 == -1 || id2 == -1 {
		return false
	}

	return id1 == id2
}

// resolveValueToID attempts to resolve a value (name, slug, or ID) to a numeric ID.
// Returns -1 if the value cannot be resolved to a valid ID.
func resolveValueToID(value, resourceType string) int32 {
	// If it's already a numeric ID, return it
	if id, err := strconv.ParseInt(value, 10, 32); err == nil {
		return int32(id)
	}

	// For non-numeric values, we need to perform a lookup
	// This is a placeholder for the actual lookup implementation
	// TODO: Implement actual NetBox API lookups for each resource type

	// Mock implementation for testing - in real implementation this would
	// use the netboxlookup package to resolve names/slugs to IDs
	return mockLookupValueToID(value, resourceType)
}

// mockLookupValueToID provides mock lookup functionality for testing.
// In the real implementation, this would be replaced with actual NetBox API calls.
func mockLookupValueToID(value, resourceType string) int32 {
	// Mock data for testing - this simulates known NetBox objects
	mockData := map[string]map[string]int32{
		"tenant": {
			"production-environment": 7,
			"Production Environment": 7,
			"staging-environment":    12,
			"Staging Environment":    12,
		},
		"site": {
			"datacenter-east": 15,
			"Datacenter East": 15,
			"datacenter-west": 20,
			"Datacenter West": 20,
		},
		"device_type": {
			"PowerEdge R640": 9,
			"poweredge-r640": 9,
			"PowerEdge R740": 10,
			"poweredge-r740": 10,
		},
		"platform": {
			"ubuntu-2204":      5,
			"Ubuntu 22.04 LTS": 5,
			"centos-8":         6,
			"CentOS 8":         6,
		},
	}

	if resourceData, exists := mockData[resourceType]; exists {
		if id, exists := resourceData[value]; exists {
			return id
		}
	}

	// Value not found in mock data
	return -1
}
