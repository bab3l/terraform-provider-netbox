// Package netboxlookup provides generic lookup utilities for Netbox resources.
//
// This package uses Go generics to provide a consistent pattern for looking up
// Netbox resources by ID or slug, reducing code duplication across resources.
package netboxlookup

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bab3l/go-netbox"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// =====================================================
// GENERIC LOOKUP INFRASTRUCTURE
// =====================================================

// LookupConfig defines how to look up a specific Netbox resource type.
// Use this with GenericLookup to create type-safe lookup functions.
type LookupConfig[TFull any, TBrief any] struct {
	// ResourceName is used in error messages (e.g., "Manufacturer", "Site")
	ResourceName string

	// RetrieveByID fetches a resource by its numeric ID
	RetrieveByID func(ctx context.Context, id int32) (TFull, *http.Response, error)

	// ListBySlug fetches resources matching a slug filter
	// Returns a slice of results (may be empty)
	ListBySlug func(ctx context.Context, slug string) ([]TFull, *http.Response, error)

	// ToBriefRequest converts a full resource to a Brief*Request for API calls
	ToBriefRequest func(resource TFull) TBrief
}

// GenericLookup performs a lookup by ID or slug using the provided config.
// It first tries to parse the value as an integer ID. If that fails,
// it falls back to looking up by slug.
//
// Returns the Brief request type suitable for use in API create/update calls.
func GenericLookup[TFull any, TBrief any](
	ctx context.Context,
	value string,
	config LookupConfig[TFull, TBrief],
) (*TBrief, diag.Diagnostics) {
	var id int32
	if _, err := fmt.Sscanf(value, "%d", &id); err == nil {
		// Lookup by ID
		resource, resp, err := config.RetrieveByID(ctx, id)
		if err != nil || resp.StatusCode != 200 {
			errMsg := "unknown error"
			if err != nil {
				errMsg = err.Error()
			}
			return nil, diag.Diagnostics{diag.NewErrorDiagnostic(
				config.ResourceName+" lookup failed",
				fmt.Sprintf("Could not find %s with ID %d: %s", config.ResourceName, id, errMsg),
			)}
		}
		result := config.ToBriefRequest(resource)
		return &result, nil
	}

	// Lookup by slug
	resources, resp, err := config.ListBySlug(ctx, value)
	if err != nil || resp.StatusCode != 200 {
		errMsg := "unknown error"
		if err != nil {
			errMsg = err.Error()
		}
		return nil, diag.Diagnostics{diag.NewErrorDiagnostic(
			config.ResourceName+" lookup failed",
			fmt.Sprintf("Could not find %s with slug '%s': %s", config.ResourceName, value, errMsg),
		)}
	}
	if len(resources) == 0 {
		return nil, diag.Diagnostics{diag.NewErrorDiagnostic(
			config.ResourceName+" lookup failed",
			fmt.Sprintf("No %s found with slug '%s'", config.ResourceName, value),
		)}
	}

	result := config.ToBriefRequest(resources[0])
	return &result, nil
}

// =====================================================
// LOOKUP CONFIGURATIONS
// =====================================================
// Pre-configured lookups for each Netbox resource type.

// ManufacturerLookupConfig returns the lookup configuration for Manufacturers.
func ManufacturerLookupConfig(client *netbox.APIClient) LookupConfig[*netbox.Manufacturer, netbox.BriefManufacturerRequest] {
	return LookupConfig[*netbox.Manufacturer, netbox.BriefManufacturerRequest]{
		ResourceName: "Manufacturer",
		RetrieveByID: func(ctx context.Context, id int32) (*netbox.Manufacturer, *http.Response, error) {
			return client.DcimAPI.DcimManufacturersRetrieve(ctx, id).Execute()
		},
		ListBySlug: func(ctx context.Context, slug string) ([]*netbox.Manufacturer, *http.Response, error) {
			list, resp, err := client.DcimAPI.DcimManufacturersList(ctx).Slug([]string{slug}).Execute()
			if err != nil {
				return nil, resp, err
			}
			// Convert to pointer slice
			results := make([]*netbox.Manufacturer, len(list.Results))
			for i := range list.Results {
				results[i] = &list.Results[i]
			}
			return results, resp, nil
		},
		ToBriefRequest: func(m *netbox.Manufacturer) netbox.BriefManufacturerRequest {
			return netbox.BriefManufacturerRequest{
				Name: m.GetName(),
				Slug: m.GetSlug(),
			}
		},
	}
}

// TenantLookupConfig returns the lookup configuration for Tenants.
func TenantLookupConfig(client *netbox.APIClient) LookupConfig[*netbox.Tenant, netbox.BriefTenantRequest] {
	return LookupConfig[*netbox.Tenant, netbox.BriefTenantRequest]{
		ResourceName: "Tenant",
		RetrieveByID: func(ctx context.Context, id int32) (*netbox.Tenant, *http.Response, error) {
			return client.TenancyAPI.TenancyTenantsRetrieve(ctx, id).Execute()
		},
		ListBySlug: func(ctx context.Context, slug string) ([]*netbox.Tenant, *http.Response, error) {
			list, resp, err := client.TenancyAPI.TenancyTenantsList(ctx).Slug([]string{slug}).Execute()
			if err != nil {
				return nil, resp, err
			}
			results := make([]*netbox.Tenant, len(list.Results))
			for i := range list.Results {
				results[i] = &list.Results[i]
			}
			return results, resp, nil
		},
		ToBriefRequest: func(t *netbox.Tenant) netbox.BriefTenantRequest {
			return netbox.BriefTenantRequest{
				Name: t.GetName(),
				Slug: t.GetSlug(),
			}
		},
	}
}

// TenantGroupLookupConfig returns the lookup configuration for Tenant Groups.
func TenantGroupLookupConfig(client *netbox.APIClient) LookupConfig[*netbox.TenantGroup, netbox.BriefTenantGroupRequest] {
	return LookupConfig[*netbox.TenantGroup, netbox.BriefTenantGroupRequest]{
		ResourceName: "Tenant group",
		RetrieveByID: func(ctx context.Context, id int32) (*netbox.TenantGroup, *http.Response, error) {
			return client.TenancyAPI.TenancyTenantGroupsRetrieve(ctx, id).Execute()
		},
		ListBySlug: func(ctx context.Context, slug string) ([]*netbox.TenantGroup, *http.Response, error) {
			list, resp, err := client.TenancyAPI.TenancyTenantGroupsList(ctx).Slug([]string{slug}).Execute()
			if err != nil {
				return nil, resp, err
			}
			results := make([]*netbox.TenantGroup, len(list.Results))
			for i := range list.Results {
				results[i] = &list.Results[i]
			}
			return results, resp, nil
		},
		ToBriefRequest: func(tg *netbox.TenantGroup) netbox.BriefTenantGroupRequest {
			return netbox.BriefTenantGroupRequest{
				Name: tg.GetName(),
				Slug: tg.GetSlug(),
			}
		},
	}
}

// RegionLookupConfig returns the lookup configuration for Regions.
func RegionLookupConfig(client *netbox.APIClient) LookupConfig[*netbox.Region, netbox.BriefRegionRequest] {
	return LookupConfig[*netbox.Region, netbox.BriefRegionRequest]{
		ResourceName: "Region",
		RetrieveByID: func(ctx context.Context, id int32) (*netbox.Region, *http.Response, error) {
			return client.DcimAPI.DcimRegionsRetrieve(ctx, id).Execute()
		},
		ListBySlug: func(ctx context.Context, slug string) ([]*netbox.Region, *http.Response, error) {
			list, resp, err := client.DcimAPI.DcimRegionsList(ctx).Slug([]string{slug}).Execute()
			if err != nil {
				return nil, resp, err
			}
			results := make([]*netbox.Region, len(list.Results))
			for i := range list.Results {
				results[i] = &list.Results[i]
			}
			return results, resp, nil
		},
		ToBriefRequest: func(r *netbox.Region) netbox.BriefRegionRequest {
			return netbox.BriefRegionRequest{
				Name: r.GetName(),
				Slug: r.GetSlug(),
			}
		},
	}
}

// SiteGroupLookupConfig returns the lookup configuration for Site Groups.
func SiteGroupLookupConfig(client *netbox.APIClient) LookupConfig[*netbox.SiteGroup, netbox.BriefSiteGroupRequest] {
	return LookupConfig[*netbox.SiteGroup, netbox.BriefSiteGroupRequest]{
		ResourceName: "Site group",
		RetrieveByID: func(ctx context.Context, id int32) (*netbox.SiteGroup, *http.Response, error) {
			return client.DcimAPI.DcimSiteGroupsRetrieve(ctx, id).Execute()
		},
		ListBySlug: func(ctx context.Context, slug string) ([]*netbox.SiteGroup, *http.Response, error) {
			list, resp, err := client.DcimAPI.DcimSiteGroupsList(ctx).Slug([]string{slug}).Execute()
			if err != nil {
				return nil, resp, err
			}
			results := make([]*netbox.SiteGroup, len(list.Results))
			for i := range list.Results {
				results[i] = &list.Results[i]
			}
			return results, resp, nil
		},
		ToBriefRequest: func(sg *netbox.SiteGroup) netbox.BriefSiteGroupRequest {
			return netbox.BriefSiteGroupRequest{
				Name: sg.GetName(),
				Slug: sg.GetSlug(),
			}
		},
	}
}

// SiteLookupConfig returns the lookup configuration for Sites.
func SiteLookupConfig(client *netbox.APIClient) LookupConfig[*netbox.Site, netbox.BriefSiteRequest] {
	return LookupConfig[*netbox.Site, netbox.BriefSiteRequest]{
		ResourceName: "Site",
		RetrieveByID: func(ctx context.Context, id int32) (*netbox.Site, *http.Response, error) {
			return client.DcimAPI.DcimSitesRetrieve(ctx, id).Execute()
		},
		ListBySlug: func(ctx context.Context, slug string) ([]*netbox.Site, *http.Response, error) {
			list, resp, err := client.DcimAPI.DcimSitesList(ctx).Slug([]string{slug}).Execute()
			if err != nil {
				return nil, resp, err
			}
			results := make([]*netbox.Site, len(list.Results))
			for i := range list.Results {
				results[i] = &list.Results[i]
			}
			return results, resp, nil
		},
		ToBriefRequest: func(s *netbox.Site) netbox.BriefSiteRequest {
			return netbox.BriefSiteRequest{
				Name: s.GetName(),
				Slug: s.GetSlug(),
			}
		},
	}
}

// LocationLookupConfig returns the lookup configuration for Locations.
func LocationLookupConfig(client *netbox.APIClient) LookupConfig[*netbox.Location, netbox.BriefLocationRequest] {
	return LookupConfig[*netbox.Location, netbox.BriefLocationRequest]{
		ResourceName: "Location",
		RetrieveByID: func(ctx context.Context, id int32) (*netbox.Location, *http.Response, error) {
			return client.DcimAPI.DcimLocationsRetrieve(ctx, id).Execute()
		},
		ListBySlug: func(ctx context.Context, slug string) ([]*netbox.Location, *http.Response, error) {
			list, resp, err := client.DcimAPI.DcimLocationsList(ctx).Slug([]string{slug}).Execute()
			if err != nil {
				return nil, resp, err
			}
			results := make([]*netbox.Location, len(list.Results))
			for i := range list.Results {
				results[i] = &list.Results[i]
			}
			return results, resp, nil
		},
		ToBriefRequest: func(l *netbox.Location) netbox.BriefLocationRequest {
			return netbox.BriefLocationRequest{
				Name: l.GetName(),
				Slug: l.GetSlug(),
			}
		},
	}
}

// RackRoleLookupConfig returns the lookup configuration for Rack Roles.
func RackRoleLookupConfig(client *netbox.APIClient) LookupConfig[*netbox.RackRole, netbox.BriefRackRoleRequest] {
	return LookupConfig[*netbox.RackRole, netbox.BriefRackRoleRequest]{
		ResourceName: "Rack role",
		RetrieveByID: func(ctx context.Context, id int32) (*netbox.RackRole, *http.Response, error) {
			return client.DcimAPI.DcimRackRolesRetrieve(ctx, id).Execute()
		},
		ListBySlug: func(ctx context.Context, slug string) ([]*netbox.RackRole, *http.Response, error) {
			list, resp, err := client.DcimAPI.DcimRackRolesList(ctx).Slug([]string{slug}).Execute()
			if err != nil {
				return nil, resp, err
			}
			results := make([]*netbox.RackRole, len(list.Results))
			for i := range list.Results {
				results[i] = &list.Results[i]
			}
			return results, resp, nil
		},
		ToBriefRequest: func(rr *netbox.RackRole) netbox.BriefRackRoleRequest {
			return netbox.BriefRackRoleRequest{
				Name: rr.GetName(),
				Slug: rr.GetSlug(),
			}
		},
	}
}

// PlatformLookupConfig returns the lookup configuration for Platforms.
func PlatformLookupConfig(client *netbox.APIClient) LookupConfig[*netbox.Platform, netbox.BriefPlatformRequest] {
	return LookupConfig[*netbox.Platform, netbox.BriefPlatformRequest]{
		ResourceName: "Platform",
		RetrieveByID: func(ctx context.Context, id int32) (*netbox.Platform, *http.Response, error) {
			return client.DcimAPI.DcimPlatformsRetrieve(ctx, id).Execute()
		},
		ListBySlug: func(ctx context.Context, slug string) ([]*netbox.Platform, *http.Response, error) {
			list, resp, err := client.DcimAPI.DcimPlatformsList(ctx).Slug([]string{slug}).Execute()
			if err != nil {
				return nil, resp, err
			}
			results := make([]*netbox.Platform, len(list.Results))
			for i := range list.Results {
				results[i] = &list.Results[i]
			}
			return results, resp, nil
		},
		ToBriefRequest: func(p *netbox.Platform) netbox.BriefPlatformRequest {
			return netbox.BriefPlatformRequest{
				Name: p.GetName(),
				Slug: p.GetSlug(),
			}
		},
	}
}

// DeviceRoleLookupConfig returns the lookup configuration for Device Roles.
func DeviceRoleLookupConfig(client *netbox.APIClient) LookupConfig[*netbox.DeviceRole, netbox.BriefDeviceRoleRequest] {
	return LookupConfig[*netbox.DeviceRole, netbox.BriefDeviceRoleRequest]{
		ResourceName: "Device role",
		RetrieveByID: func(ctx context.Context, id int32) (*netbox.DeviceRole, *http.Response, error) {
			return client.DcimAPI.DcimDeviceRolesRetrieve(ctx, id).Execute()
		},
		ListBySlug: func(ctx context.Context, slug string) ([]*netbox.DeviceRole, *http.Response, error) {
			list, resp, err := client.DcimAPI.DcimDeviceRolesList(ctx).Slug([]string{slug}).Execute()
			if err != nil {
				return nil, resp, err
			}
			results := make([]*netbox.DeviceRole, len(list.Results))
			for i := range list.Results {
				results[i] = &list.Results[i]
			}
			return results, resp, nil
		},
		ToBriefRequest: func(dr *netbox.DeviceRole) netbox.BriefDeviceRoleRequest {
			return netbox.BriefDeviceRoleRequest{
				Name: dr.GetName(),
				Slug: dr.GetSlug(),
			}
		},
	}
}

// =====================================================
// SPECIAL CASE LOOKUPS
// =====================================================
// Some resources need special handling (e.g., DeviceType needs Manufacturer).

// DeviceTypeLookupConfig returns the lookup configuration for Device Types.
func DeviceTypeLookupConfig(client *netbox.APIClient) LookupConfig[*netbox.DeviceType, netbox.BriefDeviceTypeRequest] {
	return LookupConfig[*netbox.DeviceType, netbox.BriefDeviceTypeRequest]{
		ResourceName: "Device type",
		RetrieveByID: func(ctx context.Context, id int32) (*netbox.DeviceType, *http.Response, error) {
			return client.DcimAPI.DcimDeviceTypesRetrieve(ctx, id).Execute()
		},
		ListBySlug: func(ctx context.Context, slug string) ([]*netbox.DeviceType, *http.Response, error) {
			list, resp, err := client.DcimAPI.DcimDeviceTypesList(ctx).Slug([]string{slug}).Execute()
			if err != nil {
				return nil, resp, err
			}
			results := make([]*netbox.DeviceType, len(list.Results))
			for i := range list.Results {
				results[i] = &list.Results[i]
			}
			return results, resp, nil
		},
		ToBriefRequest: func(dt *netbox.DeviceType) netbox.BriefDeviceTypeRequest {
			manufacturer := dt.GetManufacturer()
			return netbox.BriefDeviceTypeRequest{
				Manufacturer: netbox.BriefManufacturerRequest{
					Name: manufacturer.GetName(),
					Slug: manufacturer.GetSlug(),
				},
				Model: dt.GetModel(),
				Slug:  dt.GetSlug(),
			}
		},
	}
}

// RackTypeLookupConfig returns the lookup configuration for Rack Types.
func RackTypeLookupConfig(client *netbox.APIClient) LookupConfig[*netbox.RackType, netbox.BriefRackTypeRequest] {
	return LookupConfig[*netbox.RackType, netbox.BriefRackTypeRequest]{
		ResourceName: "Rack type",
		RetrieveByID: func(ctx context.Context, id int32) (*netbox.RackType, *http.Response, error) {
			return client.DcimAPI.DcimRackTypesRetrieve(ctx, id).Execute()
		},
		ListBySlug: func(ctx context.Context, slug string) ([]*netbox.RackType, *http.Response, error) {
			// RackType uses model instead of slug for list filtering
			list, resp, err := client.DcimAPI.DcimRackTypesList(ctx).Model([]string{slug}).Execute()
			if err != nil {
				return nil, resp, err
			}
			results := make([]*netbox.RackType, len(list.Results))
			for i := range list.Results {
				results[i] = &list.Results[i]
			}
			return results, resp, nil
		},
		ToBriefRequest: func(rt *netbox.RackType) netbox.BriefRackTypeRequest {
			manufacturer := rt.GetManufacturer()
			return netbox.BriefRackTypeRequest{
				Manufacturer: netbox.BriefManufacturerRequest{
					Name: manufacturer.GetName(),
					Slug: manufacturer.GetSlug(),
				},
				Model: rt.GetModel(),
				Slug:  rt.GetSlug(),
			}
		},
	}
}

// RackLookupConfig returns the lookup configuration for Racks.
// Note: Racks use name instead of slug for lookups.
func RackLookupConfig(client *netbox.APIClient) LookupConfig[*netbox.Rack, netbox.BriefRackRequest] {
	return LookupConfig[*netbox.Rack, netbox.BriefRackRequest]{
		ResourceName: "Rack",
		RetrieveByID: func(ctx context.Context, id int32) (*netbox.Rack, *http.Response, error) {
			return client.DcimAPI.DcimRacksRetrieve(ctx, id).Execute()
		},
		ListBySlug: func(ctx context.Context, name string) ([]*netbox.Rack, *http.Response, error) {
			// Rack uses name instead of slug
			list, resp, err := client.DcimAPI.DcimRacksList(ctx).Name([]string{name}).Execute()
			if err != nil {
				return nil, resp, err
			}
			results := make([]*netbox.Rack, len(list.Results))
			for i := range list.Results {
				results[i] = &list.Results[i]
			}
			return results, resp, nil
		},
		ToBriefRequest: func(r *netbox.Rack) netbox.BriefRackRequest {
			return netbox.BriefRackRequest{
				Name: r.GetName(),
			}
		},
	}
}

// =====================================================
// CONVENIENCE WRAPPER FUNCTIONS
// =====================================================
// These provide the same API as the old lookup.go functions for backward compatibility.

// LookupManufacturer looks up a Manufacturer by ID or slug.
func LookupManufacturer(ctx context.Context, client *netbox.APIClient, value string) (*netbox.BriefManufacturerRequest, diag.Diagnostics) {
	return GenericLookup(ctx, value, ManufacturerLookupConfig(client))
}

// LookupTenant looks up a Tenant by ID or slug.
func LookupTenant(ctx context.Context, client *netbox.APIClient, value string) (*netbox.BriefTenantRequest, diag.Diagnostics) {
	return GenericLookup(ctx, value, TenantLookupConfig(client))
}

// LookupTenantGroup looks up a Tenant Group by ID or slug.
func LookupTenantGroup(ctx context.Context, client *netbox.APIClient, value string) (*netbox.BriefTenantGroupRequest, diag.Diagnostics) {
	return GenericLookup(ctx, value, TenantGroupLookupConfig(client))
}

// LookupRegion looks up a Region by ID or slug.
func LookupRegion(ctx context.Context, client *netbox.APIClient, value string) (*netbox.BriefRegionRequest, diag.Diagnostics) {
	return GenericLookup(ctx, value, RegionLookupConfig(client))
}

// LookupSiteGroup looks up a Site Group by ID or slug.
func LookupSiteGroup(ctx context.Context, client *netbox.APIClient, value string) (*netbox.BriefSiteGroupRequest, diag.Diagnostics) {
	return GenericLookup(ctx, value, SiteGroupLookupConfig(client))
}

// LookupSite looks up a Site by ID or slug.
func LookupSite(ctx context.Context, client *netbox.APIClient, value string) (*netbox.BriefSiteRequest, diag.Diagnostics) {
	return GenericLookup(ctx, value, SiteLookupConfig(client))
}

// LookupLocation looks up a Location by ID or slug.
func LookupLocation(ctx context.Context, client *netbox.APIClient, value string) (*netbox.BriefLocationRequest, diag.Diagnostics) {
	return GenericLookup(ctx, value, LocationLookupConfig(client))
}

// LookupRackRole looks up a Rack Role by ID or slug.
func LookupRackRole(ctx context.Context, client *netbox.APIClient, value string) (*netbox.BriefRackRoleRequest, diag.Diagnostics) {
	return GenericLookup(ctx, value, RackRoleLookupConfig(client))
}

// LookupPlatform looks up a Platform by ID or slug.
func LookupPlatform(ctx context.Context, client *netbox.APIClient, value string) (*netbox.BriefPlatformRequest, diag.Diagnostics) {
	return GenericLookup(ctx, value, PlatformLookupConfig(client))
}

// LookupDeviceRole looks up a Device Role by ID or slug.
func LookupDeviceRole(ctx context.Context, client *netbox.APIClient, value string) (*netbox.BriefDeviceRoleRequest, diag.Diagnostics) {
	return GenericLookup(ctx, value, DeviceRoleLookupConfig(client))
}

// LookupDeviceType looks up a Device Type by ID or slug.
func LookupDeviceType(ctx context.Context, client *netbox.APIClient, value string) (*netbox.BriefDeviceTypeRequest, diag.Diagnostics) {
	return GenericLookup(ctx, value, DeviceTypeLookupConfig(client))
}

// LookupRackType looks up a Rack Type by ID or model name.
func LookupRackType(ctx context.Context, client *netbox.APIClient, value string) (*netbox.BriefRackTypeRequest, diag.Diagnostics) {
	return GenericLookup(ctx, value, RackTypeLookupConfig(client))
}

// LookupRack looks up a Rack by ID or name.
func LookupRack(ctx context.Context, client *netbox.APIClient, value string) (*netbox.BriefRackRequest, diag.Diagnostics) {
	return GenericLookup(ctx, value, RackLookupConfig(client))
}

// DeviceLookupConfig returns the lookup configuration for Devices.
func DeviceLookupConfig(client *netbox.APIClient) LookupConfig[*netbox.DeviceWithConfigContext, netbox.BriefDeviceRequest] {
	return LookupConfig[*netbox.DeviceWithConfigContext, netbox.BriefDeviceRequest]{
		ResourceName: "Device",
		RetrieveByID: func(ctx context.Context, id int32) (*netbox.DeviceWithConfigContext, *http.Response, error) {
			return client.DcimAPI.DcimDevicesRetrieve(ctx, id).Execute()
		},
		ListBySlug: func(ctx context.Context, nameOrSlug string) ([]*netbox.DeviceWithConfigContext, *http.Response, error) {
			// Devices use name, not slug
			list, resp, err := client.DcimAPI.DcimDevicesList(ctx).Name([]string{nameOrSlug}).Execute()
			if err != nil {
				return nil, resp, err
			}
			results := make([]*netbox.DeviceWithConfigContext, len(list.Results))
			for i := range list.Results {
				results[i] = &list.Results[i]
			}
			return results, resp, nil
		},
		ToBriefRequest: func(d *netbox.DeviceWithConfigContext) netbox.BriefDeviceRequest {
			req := netbox.NewBriefDeviceRequest()
			name := d.GetName()
			req.Name = *netbox.NewNullableString(&name)
			return *req
		},
	}
}

// LookupDevice looks up a Device by ID or name.
func LookupDevice(ctx context.Context, client *netbox.APIClient, value string) (*netbox.BriefDeviceRequest, diag.Diagnostics) {
	return GenericLookup(ctx, value, DeviceLookupConfig(client))
}

// LookupDeviceBrief returns a BriefDeviceRequest from an ID or name.
// Deprecated: Use LookupDevice instead.
func LookupDeviceBrief(ctx context.Context, client *netbox.APIClient, value string) (*netbox.BriefDeviceRequest, diag.Diagnostics) {
	return LookupDevice(ctx, client, value)
}

// =====================================================
// IPAM LOOKUP CONFIGS
// =====================================================

// VLANGroupLookupConfig returns the lookup configuration for VLAN Groups.
func VLANGroupLookupConfig(client *netbox.APIClient) LookupConfig[*netbox.VLANGroup, netbox.BriefVLANGroupRequest] {
	return LookupConfig[*netbox.VLANGroup, netbox.BriefVLANGroupRequest]{
		ResourceName: "VLAN Group",
		RetrieveByID: func(ctx context.Context, id int32) (*netbox.VLANGroup, *http.Response, error) {
			return client.IpamAPI.IpamVlanGroupsRetrieve(ctx, id).Execute()
		},
		ListBySlug: func(ctx context.Context, slug string) ([]*netbox.VLANGroup, *http.Response, error) {
			list, resp, err := client.IpamAPI.IpamVlanGroupsList(ctx).Slug([]string{slug}).Execute()
			if err != nil {
				return nil, resp, err
			}
			results := make([]*netbox.VLANGroup, len(list.Results))
			for i := range list.Results {
				results[i] = &list.Results[i]
			}
			return results, resp, nil
		},
		ToBriefRequest: func(vg *netbox.VLANGroup) netbox.BriefVLANGroupRequest {
			return netbox.BriefVLANGroupRequest{
				Name: vg.GetName(),
				Slug: vg.GetSlug(),
			}
		},
	}
}

// LookupVLANGroup looks up a VLAN Group by ID or slug.
func LookupVLANGroup(ctx context.Context, client *netbox.APIClient, value string) (*netbox.BriefVLANGroupRequest, diag.Diagnostics) {
	return GenericLookup(ctx, value, VLANGroupLookupConfig(client))
}

// RoleLookupConfig returns the lookup configuration for Roles (IPAM roles used for VLANs, Prefixes).
func RoleLookupConfig(client *netbox.APIClient) LookupConfig[*netbox.Role, netbox.BriefRoleRequest] {
	return LookupConfig[*netbox.Role, netbox.BriefRoleRequest]{
		ResourceName: "Role",
		RetrieveByID: func(ctx context.Context, id int32) (*netbox.Role, *http.Response, error) {
			return client.IpamAPI.IpamRolesRetrieve(ctx, id).Execute()
		},
		ListBySlug: func(ctx context.Context, slug string) ([]*netbox.Role, *http.Response, error) {
			list, resp, err := client.IpamAPI.IpamRolesList(ctx).Slug([]string{slug}).Execute()
			if err != nil {
				return nil, resp, err
			}
			results := make([]*netbox.Role, len(list.Results))
			for i := range list.Results {
				results[i] = &list.Results[i]
			}
			return results, resp, nil
		},
		ToBriefRequest: func(r *netbox.Role) netbox.BriefRoleRequest {
			return netbox.BriefRoleRequest{
				Name: r.GetName(),
				Slug: r.GetSlug(),
			}
		},
	}
}

// LookupRole looks up an IPAM Role by ID or slug.
func LookupRole(ctx context.Context, client *netbox.APIClient, value string) (*netbox.BriefRoleRequest, diag.Diagnostics) {
	return GenericLookup(ctx, value, RoleLookupConfig(client))
}

// VRFLookupConfig returns the lookup configuration for VRFs.
func VRFLookupConfig(client *netbox.APIClient) LookupConfig[*netbox.VRF, netbox.BriefVRFRequest] {
	return LookupConfig[*netbox.VRF, netbox.BriefVRFRequest]{
		ResourceName: "VRF",
		RetrieveByID: func(ctx context.Context, id int32) (*netbox.VRF, *http.Response, error) {
			return client.IpamAPI.IpamVrfsRetrieve(ctx, id).Execute()
		},
		ListBySlug: func(ctx context.Context, slug string) ([]*netbox.VRF, *http.Response, error) {
			// VRF doesn't have slug, lookup by name instead
			list, resp, err := client.IpamAPI.IpamVrfsList(ctx).Name([]string{slug}).Execute()
			if err != nil {
				return nil, resp, err
			}
			results := make([]*netbox.VRF, len(list.Results))
			for i := range list.Results {
				results[i] = &list.Results[i]
			}
			return results, resp, nil
		},
		ToBriefRequest: func(v *netbox.VRF) netbox.BriefVRFRequest {
			return netbox.BriefVRFRequest{
				Name: v.GetName(),
			}
		},
	}
}

// LookupVRF looks up a VRF by ID or name.
func LookupVRF(ctx context.Context, client *netbox.APIClient, value string) (*netbox.BriefVRFRequest, diag.Diagnostics) {
	return GenericLookup(ctx, value, VRFLookupConfig(client))
}

// VLANLookupConfig returns the lookup configuration for VLANs.
func VLANLookupConfig(client *netbox.APIClient) LookupConfig[*netbox.VLAN, netbox.BriefVLANRequest] {
	return LookupConfig[*netbox.VLAN, netbox.BriefVLANRequest]{
		ResourceName: "VLAN",
		RetrieveByID: func(ctx context.Context, id int32) (*netbox.VLAN, *http.Response, error) {
			return client.IpamAPI.IpamVlansRetrieve(ctx, id).Execute()
		},
		ListBySlug: func(ctx context.Context, slug string) ([]*netbox.VLAN, *http.Response, error) {
			// VLAN doesn't have slug, lookup by name instead
			list, resp, err := client.IpamAPI.IpamVlansList(ctx).Name([]string{slug}).Execute()
			if err != nil {
				return nil, resp, err
			}
			results := make([]*netbox.VLAN, len(list.Results))
			for i := range list.Results {
				results[i] = &list.Results[i]
			}
			return results, resp, nil
		},
		ToBriefRequest: func(v *netbox.VLAN) netbox.BriefVLANRequest {
			return netbox.BriefVLANRequest{
				Vid:  v.GetVid(),
				Name: v.GetName(),
			}
		},
	}
}

// LookupVLAN looks up a VLAN by ID or name.
func LookupVLAN(ctx context.Context, client *netbox.APIClient, value string) (*netbox.BriefVLANRequest, diag.Diagnostics) {
	return GenericLookup(ctx, value, VLANLookupConfig(client))
}
