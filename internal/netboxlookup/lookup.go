// Package netboxlookup provides lookup utilities for Netbox resources.
//
// This file provides backward-compatible wrapper functions that delegate to
// the generic lookup implementation in generic_lookup.go.
//
// For new code, consider using the generic functions directly:
//   - GenericLookup[TFull, TBrief](ctx, value, config) for the most flexibility
//   - LookupX(ctx, client, value) convenience functions for common types
package netboxlookup

import (
	"context"

	"github.com/bab3l/go-netbox"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// =====================================================
// BACKWARD-COMPATIBLE WRAPPER FUNCTIONS
// =====================================================
// These functions maintain the old API (LookupXBrief) while delegating to
// the new generic implementation. New code should use LookupX instead.

// LookupManufacturerBrief returns a BriefManufacturerRequest from an ID or slug.
// Deprecated: Use LookupManufacturer instead.
func LookupManufacturerBrief(ctx context.Context, client *netbox.APIClient, value string) (*netbox.BriefManufacturerRequest, diag.Diagnostics) {
	return LookupManufacturer(ctx, client, value)
}

// LookupTenantBrief returns a BriefTenantRequest from an ID or slug.
// Deprecated: Use LookupTenant instead.
func LookupTenantBrief(ctx context.Context, client *netbox.APIClient, value string) (*netbox.BriefTenantRequest, diag.Diagnostics) {
	return LookupTenant(ctx, client, value)
}

// LookupTenantGroupBrief returns a BriefTenantGroupRequest from an ID or slug.
// Deprecated: Use LookupTenantGroup instead.
func LookupTenantGroupBrief(ctx context.Context, client *netbox.APIClient, value string) (*netbox.BriefTenantGroupRequest, diag.Diagnostics) {
	return LookupTenantGroup(ctx, client, value)
}

// LookupRegionBrief returns a BriefRegionRequest from an ID or slug.
// Deprecated: Use LookupRegion instead.
func LookupRegionBrief(ctx context.Context, client *netbox.APIClient, value string) (*netbox.BriefRegionRequest, diag.Diagnostics) {
	return LookupRegion(ctx, client, value)
}

// LookupSiteGroupBrief returns a BriefSiteGroupRequest from an ID or slug.
// Deprecated: Use LookupSiteGroup instead.
func LookupSiteGroupBrief(ctx context.Context, client *netbox.APIClient, value string) (*netbox.BriefSiteGroupRequest, diag.Diagnostics) {
	return LookupSiteGroup(ctx, client, value)
}

// LookupSiteBrief returns a BriefSiteRequest from an ID or slug.
// Deprecated: Use LookupSite instead.
func LookupSiteBrief(ctx context.Context, client *netbox.APIClient, value string) (*netbox.BriefSiteRequest, diag.Diagnostics) {
	return LookupSite(ctx, client, value)
}

// LookupLocationBrief returns a BriefLocationRequest from an ID or slug.
// Deprecated: Use LookupLocation instead.
func LookupLocationBrief(ctx context.Context, client *netbox.APIClient, value string) (*netbox.BriefLocationRequest, diag.Diagnostics) {
	return LookupLocation(ctx, client, value)
}

// LookupRackRoleBrief returns a BriefRackRoleRequest from an ID or slug.
// Deprecated: Use LookupRackRole instead.
func LookupRackRoleBrief(ctx context.Context, client *netbox.APIClient, value string) (*netbox.BriefRackRoleRequest, diag.Diagnostics) {
	return LookupRackRole(ctx, client, value)
}

// LookupRackTypeBrief returns a BriefRackTypeRequest from an ID or model name.
// Deprecated: Use LookupRackType instead.
func LookupRackTypeBrief(ctx context.Context, client *netbox.APIClient, value string) (*netbox.BriefRackTypeRequest, diag.Diagnostics) {
	return LookupRackType(ctx, client, value)
}

// LookupPlatformBrief returns a BriefPlatformRequest from an ID or slug.
// Deprecated: Use LookupPlatform instead.
func LookupPlatformBrief(ctx context.Context, client *netbox.APIClient, value string) (*netbox.BriefPlatformRequest, diag.Diagnostics) {
	return LookupPlatform(ctx, client, value)
}

// LookupDeviceTypeBrief returns a BriefDeviceTypeRequest from an ID or slug.
// Deprecated: Use LookupDeviceType instead.
func LookupDeviceTypeBrief(ctx context.Context, client *netbox.APIClient, value string) (*netbox.BriefDeviceTypeRequest, diag.Diagnostics) {
	return LookupDeviceType(ctx, client, value)
}

// LookupDeviceRoleBrief returns a BriefDeviceRoleRequest from an ID or slug.
// Deprecated: Use LookupDeviceRole instead.
func LookupDeviceRoleBrief(ctx context.Context, client *netbox.APIClient, value string) (*netbox.BriefDeviceRoleRequest, diag.Diagnostics) {
	return LookupDeviceRole(ctx, client, value)
}

// LookupRackBrief returns a BriefRackRequest from an ID or name.
// Deprecated: Use LookupRack instead.
func LookupRackBrief(ctx context.Context, client *netbox.APIClient, value string) (*netbox.BriefRackRequest, diag.Diagnostics) {
	return LookupRack(ctx, client, value)
}
