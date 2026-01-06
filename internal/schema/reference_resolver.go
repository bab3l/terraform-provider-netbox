// Package schema provides reference resolution functions for diff suppression.
//
// This file extends the existing netboxlookup package with ID-returning functions
// specifically designed for diff suppression logic. These functions resolve
// names, slugs, or IDs to canonical IDs for equivalency checking.
package schema

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/bab3l/go-netbox"
	"github.com/bab3l/terraform-provider-netbox/internal/netboxlookup"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// ReferenceResolver provides functionality to resolve reference values to canonical IDs.
type ReferenceResolver struct {
	client *netbox.APIClient
}

// NewReferenceResolver creates a new reference resolver with the given NetBox client.
func NewReferenceResolver(client *netbox.APIClient) *ReferenceResolver {
	return &ReferenceResolver{client: client}
}

// ResolveToID resolves any reference (name, slug, ID) to canonical ID.
func (r *ReferenceResolver) ResolveToID(ctx context.Context, value string, resourceType string) (int32, error) {
	// Handle empty values
	if strings.TrimSpace(value) == "" {
		return 0, nil
	}

	// Try parsing as ID first
	if id, err := strconv.ParseInt(value, 10, 32); err == nil {
		return int32(id), nil
	}

	// Use existing lookup functions to resolve name/slug to ID
	var id int32
	var diags diag.Diagnostics

	switch resourceType {
	case "tenant":
		id, diags = r.resolveTenantToID(ctx, value)
	case "site":
		id, diags = r.resolveSiteToID(ctx, value)
	case "device_type":
		id, diags = r.resolveDeviceTypeToID(ctx, value)
	case "device_role":
		id, diags = r.resolveDeviceRoleToID(ctx, value)
	case "platform":
		id, diags = r.resolvePlatformToID(ctx, value)
	case "rack":
		id, diags = r.resolveRackToID(ctx, value)
	case "rack_role":
		id, diags = r.resolveRackRoleToID(ctx, value)
	case "rack_type":
		id, diags = r.resolveRackTypeToID(ctx, value)
	case "location":
		// Use existing ID function
		id, diags = netboxlookup.LookupLocationID(ctx, r.client, value)
	case "region":
		// Use existing ID function
		id, diags = netboxlookup.LookupRegionID(ctx, r.client, value)
	case "site_group":
		// Use existing ID function
		id, diags = netboxlookup.LookupSiteGroupID(ctx, r.client, value)
	case "tenant_group":
		// Use existing ID function
		id, diags = netboxlookup.LookupTenantGroupID(ctx, r.client, value)
	default:
		// Unknown resource type, return as-is and let normal comparison handle it
		return 0, nil
	}

	if diags.HasError() {
		return 0, convertDiagsToError(diags)
	}

	return id, nil
}

// Helper functions that use the generic lookup infrastructure to get IDs.
func (r *ReferenceResolver) resolveTenantToID(ctx context.Context, value string) (int32, diag.Diagnostics) {
	config := netboxlookup.TenantLookupConfig(r.client)
	return netboxlookup.GenericLookupID(ctx, value, config, func(t *netbox.Tenant) int32 {
		return t.GetId()
	})
}

func (r *ReferenceResolver) resolveSiteToID(ctx context.Context, value string) (int32, diag.Diagnostics) {
	config := netboxlookup.SiteLookupConfig(r.client)
	return netboxlookup.GenericLookupID(ctx, value, config, func(s *netbox.Site) int32 {
		return s.GetId()
	})
}

func (r *ReferenceResolver) resolveDeviceTypeToID(ctx context.Context, value string) (int32, diag.Diagnostics) {
	config := netboxlookup.DeviceTypeLookupConfig(r.client)
	return netboxlookup.GenericLookupID(ctx, value, config, func(dt *netbox.DeviceType) int32 {
		return dt.GetId()
	})
}

func (r *ReferenceResolver) resolveDeviceRoleToID(ctx context.Context, value string) (int32, diag.Diagnostics) {
	config := netboxlookup.DeviceRoleLookupConfig(r.client)
	return netboxlookup.GenericLookupID(ctx, value, config, func(dr *netbox.DeviceRole) int32 {
		return dr.GetId()
	})
}

func (r *ReferenceResolver) resolvePlatformToID(ctx context.Context, value string) (int32, diag.Diagnostics) {
	config := netboxlookup.PlatformLookupConfig(r.client)
	return netboxlookup.GenericLookupID(ctx, value, config, func(p *netbox.Platform) int32 {
		return p.GetId()
	})
}

func (r *ReferenceResolver) resolveRackToID(ctx context.Context, value string) (int32, diag.Diagnostics) {
	config := netboxlookup.RackLookupConfig(r.client)
	return netboxlookup.GenericLookupID(ctx, value, config, func(rack *netbox.Rack) int32 {
		return rack.GetId()
	})
}

func (r *ReferenceResolver) resolveRackRoleToID(ctx context.Context, value string) (int32, diag.Diagnostics) {
	config := netboxlookup.RackRoleLookupConfig(r.client)
	return netboxlookup.GenericLookupID(ctx, value, config, func(rr *netbox.RackRole) int32 {
		return rr.GetId()
	})
}

func (r *ReferenceResolver) resolveRackTypeToID(ctx context.Context, value string) (int32, diag.Diagnostics) {
	config := netboxlookup.RackTypeLookupConfig(r.client)
	return netboxlookup.GenericLookupID(ctx, value, config, func(rt *netbox.RackType) int32 {
		return rt.GetId()
	})
}

// convertDiagsToError converts terraform diagnostics to a simple error
// for use in diff suppression context where we can't return diagnostics.
func convertDiagsToError(diags diag.Diagnostics) error {
	if !diags.HasError() {
		return nil
	}

	var messages []string
	for _, d := range diags.Errors() {
		messages = append(messages, d.Summary())
	}

	if len(messages) == 0 {
		return nil
	}

	return fmt.Errorf("%s", strings.Join(messages, "; "))
}
