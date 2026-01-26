// Package schema provides schema utilities for Terraform Provider NetBox.
//
// This file implements diff suppression logic for reference fields to ensure
// consistent plan output when users specify names, slugs, or IDs for resources.
// Based on AWS provider patterns for handling equivalent resource references.
package schema

import (
	"context"
	"strings"

	"github.com/bab3l/go-netbox"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// suppressReferenceEquivalent suppresses diffs when old and new values
// refer to the same NetBox object but use different representations
// (name vs ID vs slug). This is used as a plan modifier for reference fields.
func suppressReferenceEquivalent(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	// If the user explicitly configured a value, honor their chosen format.
	// We only suppress diffs for unmanaged (null/unknown) config values.
	if !req.ConfigValue.IsNull() && !req.ConfigValue.IsUnknown() {
		return
	}

	// Skip if no prior state (new resource)
	if req.StateValue.IsNull() || req.StateValue.IsUnknown() {
		return
	}

	// Skip if planned value is unknown
	if req.PlanValue.IsUnknown() {
		return
	}

	oldValue := req.StateValue.ValueString()
	newValue := req.PlanValue.ValueString()

	// If values are identical, no suppression needed
	if oldValue == newValue {
		return
	}

	// Get NetBox client from provider - this will be passed via context
	client, ok := ctx.Value("netbox_client").(*netbox.APIClient)
	if !ok {
		// No client available, can't resolve values - let diff show
		return
	}

	// Detect resource type from attribute path
	resourceType := getResourceTypeFromAttribute(req.Path.String())
	if resourceType == "" {
		// Can't determine type, let Terraform show diff
		return
	}

	// Check if values represent the same resource
	if areValuesEquivalent(ctx, client, oldValue, newValue, resourceType) {
		// Values are equivalent, suppress the diff
		resp.PlanValue = req.StateValue
	}
}

// ReferenceEquivalencePlanModifier creates a plan modifier that suppresses diffs
// for equivalent reference values (name vs ID vs slug).
func ReferenceEquivalencePlanModifier() planmodifier.String {
	return suppressReferenceEquivalentModifier{}
}

// suppressReferenceEquivalentModifier implements the plan modifier interface.
type suppressReferenceEquivalentModifier struct{}

func (m suppressReferenceEquivalentModifier) Description(ctx context.Context) string {
	return "Suppresses plan differences when reference values are equivalent (name/slug/ID refer to the same object)"
}

func (m suppressReferenceEquivalentModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m suppressReferenceEquivalentModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	suppressReferenceEquivalent(ctx, req, resp)
}

// getResourceTypeFromAttribute extracts the NetBox resource type from a Terraform attribute path.
// This helps determine which NetBox API endpoint to use for lookups.
func getResourceTypeFromAttribute(attributePath string) string {
	// Extract the last component of the path for attribute name matching
	parts := strings.Split(attributePath, ".")
	if len(parts) == 0 {
		return ""
	}

	attributeName := parts[len(parts)-1]

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
		"rack_role":    "rack_role",
		"rack_type":    "rack_type",
		"site_group":   "site_group",
		"tenant_group": "tenant_group",
	}

	if resourceType, exists := mappings[attributeName]; exists {
		return resourceType
	}

	// Generic attributes like "parent" or "group" need context to determine type
	// For now, return empty - these may need special handling in the future
	return ""
}

// areValuesEquivalent checks if two values (which may be names, slugs, or IDs)
// refer to the same NetBox object of the given resource type.
func areValuesEquivalent(ctx context.Context, client *netbox.APIClient, value1, value2, resourceType string) bool {
	resolver := NewReferenceResolver(client)

	// Convert both values to IDs and compare
	id1, err1 := resolver.ResolveToID(ctx, value1, resourceType)
	id2, err2 := resolver.ResolveToID(ctx, value2, resourceType)

	// If either resolution failed, they're not equivalent
	if err1 != nil || err2 != nil {
		return false
	}

	// Both resolved successfully - compare IDs
	return id1 == id2
}
