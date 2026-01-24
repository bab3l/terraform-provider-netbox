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
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

const unknownErrorMsg = "unknown error"

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

		defer utils.CloseResponseBody(resp)

		if err != nil || resp.StatusCode != 200 {

			errMsg := unknownErrorMsg

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

	defer utils.CloseResponseBody(resp)

	if err != nil || resp.StatusCode != 200 {

		errMsg := unknownErrorMsg

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

// GenericLookupID performs a lookup by ID or slug and returns the ID.

func GenericLookupID[TFull any, TBrief any](

	ctx context.Context,

	value string,

	config LookupConfig[TFull, TBrief],

	getID func(TFull) int32,

) (int32, diag.Diagnostics) {

	var id int32

	if _, err := fmt.Sscanf(value, "%d", &id); err == nil {

		// Lookup by ID

		_, resp, err := config.RetrieveByID(ctx, id)

		defer utils.CloseResponseBody(resp)

		if err != nil || resp.StatusCode != 200 {

			errMsg := unknownErrorMsg

			if err != nil {

				errMsg = err.Error()

			}

			return 0, diag.Diagnostics{diag.NewErrorDiagnostic(

				config.ResourceName+" lookup failed",

				fmt.Sprintf("Could not find %s with ID %d: %s", config.ResourceName, id, errMsg),
			)}

		}

		return id, nil

	}

	// Lookup by slug

	resources, resp, err := config.ListBySlug(ctx, value)

	defer utils.CloseResponseBody(resp)

	if err != nil {

		return 0, diag.Diagnostics{diag.NewErrorDiagnostic(

			config.ResourceName+" lookup failed",

			fmt.Sprintf("Error searching for %s with slug '%s': %s", config.ResourceName, value, err.Error()),
		)}

	}

	if len(resources) == 0 {

		return 0, diag.Diagnostics{diag.NewErrorDiagnostic(

			config.ResourceName+" lookup failed",

			fmt.Sprintf("No %s found with slug '%s'", config.ResourceName, value),
		)}

	}

	return getID(resources[0]), nil

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

			// Try lookup by slug first

			list, resp, err := client.DcimAPI.DcimManufacturersList(ctx).Slug([]string{slug}).Execute()

			if err != nil {

				return nil, resp, err

			}

			if len(list.Results) > 0 {

				results := make([]*netbox.Manufacturer, len(list.Results))

				for i := range list.Results {

					results[i] = &list.Results[i]

				}

				return results, resp, nil

			}

			// If no results by slug, try lookup by name

			list, resp, err = client.DcimAPI.DcimManufacturersList(ctx).Name([]string{slug}).Execute()

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

			// Try lookup by slug first

			list, resp, err := client.TenancyAPI.TenancyTenantsList(ctx).Slug([]string{slug}).Execute()

			if err != nil {

				return nil, resp, err

			}

			if len(list.Results) > 0 {

				results := make([]*netbox.Tenant, len(list.Results))

				for i := range list.Results {

					results[i] = &list.Results[i]

				}

				return results, resp, nil

			}

			// If no results by slug, try lookup by name

			list, resp, err = client.TenancyAPI.TenancyTenantsList(ctx).Name([]string{slug}).Execute()

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

			// Try lookup by slug first

			list, resp, err := client.TenancyAPI.TenancyTenantGroupsList(ctx).Slug([]string{slug}).Execute()

			if err != nil {

				return nil, resp, err

			}

			if len(list.Results) > 0 {

				results := make([]*netbox.TenantGroup, len(list.Results))

				for i := range list.Results {

					results[i] = &list.Results[i]

				}

				return results, resp, nil

			}

			// If no results by slug, try lookup by name

			list, resp, err = client.TenancyAPI.TenancyTenantGroupsList(ctx).Name([]string{slug}).Execute()

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

			// Try lookup by slug first

			list, resp, err := client.DcimAPI.DcimRegionsList(ctx).Slug([]string{slug}).Execute()

			if err != nil {

				return nil, resp, err

			}

			if len(list.Results) > 0 {

				results := make([]*netbox.Region, len(list.Results))

				for i := range list.Results {

					results[i] = &list.Results[i]

				}

				return results, resp, nil

			}

			// If no results by slug, try lookup by name

			list, resp, err = client.DcimAPI.DcimRegionsList(ctx).Name([]string{slug}).Execute()

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

			// Try lookup by slug first

			list, resp, err := client.DcimAPI.DcimSiteGroupsList(ctx).Slug([]string{slug}).Execute()

			if err != nil {

				return nil, resp, err

			}

			if len(list.Results) > 0 {

				results := make([]*netbox.SiteGroup, len(list.Results))

				for i := range list.Results {

					results[i] = &list.Results[i]

				}

				return results, resp, nil

			}

			// If no results by slug, try lookup by name

			list, resp, err = client.DcimAPI.DcimSiteGroupsList(ctx).Name([]string{slug}).Execute()

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

			// Try lookup by slug first

			list, resp, err := client.DcimAPI.DcimSitesList(ctx).Slug([]string{slug}).Execute()

			if err != nil {

				return nil, resp, err

			}

			if len(list.Results) > 0 {

				results := make([]*netbox.Site, len(list.Results))

				for i := range list.Results {

					results[i] = &list.Results[i]

				}

				return results, resp, nil

			}

			// If no results by slug, try lookup by name

			list, resp, err = client.DcimAPI.DcimSitesList(ctx).Name([]string{slug}).Execute()

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

			// Try lookup by slug first

			list, resp, err := client.DcimAPI.DcimLocationsList(ctx).Slug([]string{slug}).Execute()

			if err != nil {

				return nil, resp, err

			}

			if len(list.Results) > 0 {

				results := make([]*netbox.Location, len(list.Results))

				for i := range list.Results {

					results[i] = &list.Results[i]

				}

				return results, resp, nil

			}

			// If no results by slug, try lookup by name

			list, resp, err = client.DcimAPI.DcimLocationsList(ctx).Name([]string{slug}).Execute()

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

			// Try lookup by slug first

			list, resp, err := client.DcimAPI.DcimRackRolesList(ctx).Slug([]string{slug}).Execute()

			if err != nil {

				return nil, resp, err

			}

			if len(list.Results) > 0 {

				results := make([]*netbox.RackRole, len(list.Results))

				for i := range list.Results {

					results[i] = &list.Results[i]

				}

				return results, resp, nil

			}

			// If no results by slug, try lookup by name

			list, resp, err = client.DcimAPI.DcimRackRolesList(ctx).Name([]string{slug}).Execute()

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

			// Try lookup by slug first

			list, resp, err := client.DcimAPI.DcimPlatformsList(ctx).Slug([]string{slug}).Execute()

			if err != nil {

				return nil, resp, err

			}

			if len(list.Results) > 0 {

				results := make([]*netbox.Platform, len(list.Results))

				for i := range list.Results {

					results[i] = &list.Results[i]

				}

				return results, resp, nil

			}

			// If no results by slug, try lookup by name

			list, resp, err = client.DcimAPI.DcimPlatformsList(ctx).Name([]string{slug}).Execute()

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

			// Try lookup by slug first

			list, resp, err := client.DcimAPI.DcimDeviceRolesList(ctx).Slug([]string{slug}).Execute()

			if err != nil {

				return nil, resp, err

			}

			if len(list.Results) > 0 {

				results := make([]*netbox.DeviceRole, len(list.Results))

				for i := range list.Results {

					results[i] = &list.Results[i]

				}

				return results, resp, nil

			}

			// If no results by slug, try lookup by name

			list, resp, err = client.DcimAPI.DcimDeviceRolesList(ctx).Name([]string{slug}).Execute()

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

			// Try lookup by slug first

			list, resp, err := client.DcimAPI.DcimDeviceTypesList(ctx).Slug([]string{slug}).Execute()

			if err != nil {

				return nil, resp, err

			}

			if len(list.Results) > 0 {

				results := make([]*netbox.DeviceType, len(list.Results))

				for i := range list.Results {

					results[i] = &list.Results[i]

				}

				return results, resp, nil

			}

			// If no results by slug, try lookup by model

			list, resp, err = client.DcimAPI.DcimDeviceTypesList(ctx).Model([]string{slug}).Execute()

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

				Slug: dt.GetSlug(),
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

			// Try lookup by slug first

			list, resp, err := client.DcimAPI.DcimRackTypesList(ctx).Slug([]string{slug}).Execute()

			if err != nil {

				return nil, resp, err

			}

			if len(list.Results) > 0 {

				results := make([]*netbox.RackType, len(list.Results))

				for i := range list.Results {

					results[i] = &list.Results[i]

				}

				return results, resp, nil

			}

			// If no results by slug, try lookup by model

			list, resp, err = client.DcimAPI.DcimRackTypesList(ctx).Model([]string{slug}).Execute()

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

				Slug: rt.GetSlug(),
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

// PowerPanelLookupConfig returns the lookup configuration for Power Panels.

func PowerPanelLookupConfig(client *netbox.APIClient) LookupConfig[*netbox.PowerPanel, netbox.BriefPowerPanelRequest] {

	return LookupConfig[*netbox.PowerPanel, netbox.BriefPowerPanelRequest]{

		ResourceName: "Power Panel",

		RetrieveByID: func(ctx context.Context, id int32) (*netbox.PowerPanel, *http.Response, error) {

			return client.DcimAPI.DcimPowerPanelsRetrieve(ctx, id).Execute()

		},

		ListBySlug: func(ctx context.Context, name string) ([]*netbox.PowerPanel, *http.Response, error) {

			// PowerPanel uses name instead of slug

			list, resp, err := client.DcimAPI.DcimPowerPanelsList(ctx).Name([]string{name}).Execute()

			if err != nil {

				return nil, resp, err

			}

			results := make([]*netbox.PowerPanel, len(list.Results))

			for i := range list.Results {

				results[i] = &list.Results[i]

			}

			return results, resp, nil

		},

		ToBriefRequest: func(p *netbox.PowerPanel) netbox.BriefPowerPanelRequest {

			return netbox.BriefPowerPanelRequest{

				Name: p.GetName(),
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

// LookupTenantGroupID looks up a Tenant Group by ID or slug and returns the ID.

func LookupTenantGroupID(ctx context.Context, client *netbox.APIClient, value string) (int32, diag.Diagnostics) {

	return GenericLookupID(ctx, value, TenantGroupLookupConfig(client), func(tg *netbox.TenantGroup) int32 {

		return tg.GetId()

	})

}

// LookupRegionID looks up a Region by ID or slug and returns the ID.

func LookupRegionID(ctx context.Context, client *netbox.APIClient, value string) (int32, diag.Diagnostics) {

	return GenericLookupID(ctx, value, RegionLookupConfig(client), func(r *netbox.Region) int32 {

		return r.GetId()

	})

}

// LookupSiteGroupID looks up a Site Group by ID or slug and returns the ID.

func LookupSiteGroupID(ctx context.Context, client *netbox.APIClient, value string) (int32, diag.Diagnostics) {

	return GenericLookupID(ctx, value, SiteGroupLookupConfig(client), func(sg *netbox.SiteGroup) int32 {

		return sg.GetId()

	})

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

// LookupLocationID looks up a Location by ID or slug and returns the ID.

func LookupLocationID(ctx context.Context, client *netbox.APIClient, value string) (int32, diag.Diagnostics) {

	return GenericLookupID(ctx, value, LocationLookupConfig(client), func(l *netbox.Location) int32 {

		return l.GetId()

	})

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

// LookupPowerPanel looks up a Power Panel by ID or name.

func LookupPowerPanel(ctx context.Context, client *netbox.APIClient, value string) (*netbox.BriefPowerPanelRequest, diag.Diagnostics) {

	return GenericLookup(ctx, value, PowerPanelLookupConfig(client))

}

// DeviceLookupConfig returns the lookup configuration for Devices.

func DeviceLookupConfig(client *netbox.APIClient) LookupConfig[*netbox.Device, netbox.BriefDeviceRequest] {

	return LookupConfig[*netbox.Device, netbox.BriefDeviceRequest]{

		ResourceName: "Device",

		RetrieveByID: func(ctx context.Context, id int32) (*netbox.Device, *http.Response, error) {

			return client.DcimAPI.DcimDevicesRetrieve(ctx, id).Execute()

		},

		ListBySlug: func(ctx context.Context, nameOrSlug string) ([]*netbox.Device, *http.Response, error) {

			// Devices use name, not slug

			list, resp, err := client.DcimAPI.DcimDevicesList(ctx).Name([]string{nameOrSlug}).Execute()

			if err != nil {

				return nil, resp, err

			}

			results := make([]*netbox.Device, len(list.Results))

			for i := range list.Results {

				results[i] = &list.Results[i]

			}

			return results, resp, nil

		},

		ToBriefRequest: func(d *netbox.Device) netbox.BriefDeviceRequest {

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

//

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

				Vid: v.GetVid(),

				Name: v.GetName(),
			}

		},
	}

}

// LookupVLAN looks up a VLAN by ID or name.

func LookupVLAN(ctx context.Context, client *netbox.APIClient, value string) (*netbox.BriefVLANRequest, diag.Diagnostics) {

	return GenericLookup(ctx, value, VLANLookupConfig(client))

}

// =====================================================

// VIRTUALIZATION LOOKUPS

// =====================================================

// ClusterTypeLookupConfig returns the lookup configuration for Cluster Types.

func ClusterTypeLookupConfig(client *netbox.APIClient) LookupConfig[*netbox.ClusterType, netbox.BriefClusterTypeRequest] {

	return LookupConfig[*netbox.ClusterType, netbox.BriefClusterTypeRequest]{

		ResourceName: "Cluster Type",

		RetrieveByID: func(ctx context.Context, id int32) (*netbox.ClusterType, *http.Response, error) {

			return client.VirtualizationAPI.VirtualizationClusterTypesRetrieve(ctx, id).Execute()

		},

		ListBySlug: func(ctx context.Context, slug string) ([]*netbox.ClusterType, *http.Response, error) {

			list, resp, err := client.VirtualizationAPI.VirtualizationClusterTypesList(ctx).Slug([]string{slug}).Execute()

			if err != nil {

				return nil, resp, err

			}

			results := make([]*netbox.ClusterType, len(list.Results))

			for i := range list.Results {

				results[i] = &list.Results[i]

			}

			return results, resp, nil

		},

		ToBriefRequest: func(ct *netbox.ClusterType) netbox.BriefClusterTypeRequest {

			return netbox.BriefClusterTypeRequest{

				Name: ct.GetName(),

				Slug: ct.GetSlug(),
			}

		},
	}

}

// LookupClusterType looks up a Cluster Type by ID or slug.

func LookupClusterType(ctx context.Context, client *netbox.APIClient, value string) (*netbox.BriefClusterTypeRequest, diag.Diagnostics) {

	return GenericLookup(ctx, value, ClusterTypeLookupConfig(client))

}

// ClusterGroupLookupConfig returns the lookup configuration for Cluster Groups.

func ClusterGroupLookupConfig(client *netbox.APIClient) LookupConfig[*netbox.ClusterGroup, netbox.BriefClusterGroupRequest] {

	return LookupConfig[*netbox.ClusterGroup, netbox.BriefClusterGroupRequest]{

		ResourceName: "Cluster Group",

		RetrieveByID: func(ctx context.Context, id int32) (*netbox.ClusterGroup, *http.Response, error) {

			return client.VirtualizationAPI.VirtualizationClusterGroupsRetrieve(ctx, id).Execute()

		},

		ListBySlug: func(ctx context.Context, slug string) ([]*netbox.ClusterGroup, *http.Response, error) {

			list, resp, err := client.VirtualizationAPI.VirtualizationClusterGroupsList(ctx).Slug([]string{slug}).Execute()

			if err != nil {

				return nil, resp, err

			}

			results := make([]*netbox.ClusterGroup, len(list.Results))

			for i := range list.Results {

				results[i] = &list.Results[i]

			}

			return results, resp, nil

		},

		ToBriefRequest: func(cg *netbox.ClusterGroup) netbox.BriefClusterGroupRequest {

			return netbox.BriefClusterGroupRequest{

				Name: cg.GetName(),

				Slug: cg.GetSlug(),
			}

		},
	}

}

// LookupClusterGroup looks up a Cluster Group by ID or slug.

func LookupClusterGroup(ctx context.Context, client *netbox.APIClient, value string) (*netbox.BriefClusterGroupRequest, diag.Diagnostics) {

	return GenericLookup(ctx, value, ClusterGroupLookupConfig(client))

}

// ClusterLookupConfig returns the lookup configuration for Clusters.

func ClusterLookupConfig(client *netbox.APIClient) LookupConfig[*netbox.Cluster, netbox.BriefClusterRequest] {

	return LookupConfig[*netbox.Cluster, netbox.BriefClusterRequest]{

		ResourceName: "Cluster",

		RetrieveByID: func(ctx context.Context, id int32) (*netbox.Cluster, *http.Response, error) {

			return client.VirtualizationAPI.VirtualizationClustersRetrieve(ctx, id).Execute()

		},

		ListBySlug: func(ctx context.Context, slug string) ([]*netbox.Cluster, *http.Response, error) {

			// Cluster doesn't have slug, lookup by name instead

			list, resp, err := client.VirtualizationAPI.VirtualizationClustersList(ctx).Name([]string{slug}).Execute()

			if err != nil {

				return nil, resp, err

			}

			results := make([]*netbox.Cluster, len(list.Results))

			for i := range list.Results {

				results[i] = &list.Results[i]

			}

			return results, resp, nil

		},

		ToBriefRequest: func(c *netbox.Cluster) netbox.BriefClusterRequest {

			return netbox.BriefClusterRequest{

				Name: c.GetName(),
			}

		},
	}

}

// LookupCluster looks up a Cluster by ID or name.

func LookupCluster(ctx context.Context, client *netbox.APIClient, value string) (*netbox.BriefClusterRequest, diag.Diagnostics) {

	return GenericLookup(ctx, value, ClusterLookupConfig(client))

}

// ConfigTemplateLookupConfig returns the lookup configuration for Config Templates.

func ConfigTemplateLookupConfig(client *netbox.APIClient) LookupConfig[*netbox.ConfigTemplate, netbox.BriefConfigTemplateRequest] {

	return LookupConfig[*netbox.ConfigTemplate, netbox.BriefConfigTemplateRequest]{

		ResourceName: "Config Template",

		RetrieveByID: func(ctx context.Context, id int32) (*netbox.ConfigTemplate, *http.Response, error) {

			return client.ExtrasAPI.ExtrasConfigTemplatesRetrieve(ctx, id).Execute()

		},

		ListBySlug: func(ctx context.Context, name string) ([]*netbox.ConfigTemplate, *http.Response, error) {

			list, resp, err := client.ExtrasAPI.ExtrasConfigTemplatesList(ctx).Name([]string{name}).Execute()

			if err != nil {

				return nil, resp, err

			}

			results := make([]*netbox.ConfigTemplate, len(list.Results))

			for i := range list.Results {

				results[i] = &list.Results[i]

			}

			return results, resp, nil

		},

		ToBriefRequest: func(t *netbox.ConfigTemplate) netbox.BriefConfigTemplateRequest {

			return netbox.BriefConfigTemplateRequest{

				Name: t.GetName(),
			}

		},
	}

}

// LookupConfigTemplate looks up a Config Template by ID or name.

func LookupConfigTemplate(ctx context.Context, client *netbox.APIClient, value string) (*netbox.BriefConfigTemplateRequest, diag.Diagnostics) {

	return GenericLookup(ctx, value, ConfigTemplateLookupConfig(client))

}

// VirtualMachineLookupConfig returns the lookup configuration for Virtual Machines.

func VirtualMachineLookupConfig(client *netbox.APIClient) LookupConfig[*netbox.VirtualMachineWithConfigContext, netbox.BriefVirtualMachineRequest] {

	return LookupConfig[*netbox.VirtualMachineWithConfigContext, netbox.BriefVirtualMachineRequest]{

		ResourceName: "Virtual Machine",

		RetrieveByID: func(ctx context.Context, id int32) (*netbox.VirtualMachineWithConfigContext, *http.Response, error) {

			return client.VirtualizationAPI.VirtualizationVirtualMachinesRetrieve(ctx, id).Execute()

		},

		ListBySlug: func(ctx context.Context, slug string) ([]*netbox.VirtualMachineWithConfigContext, *http.Response, error) {

			// Virtual Machine doesn't have slug, lookup by name instead

			list, resp, err := client.VirtualizationAPI.VirtualizationVirtualMachinesList(ctx).Name([]string{slug}).Execute()

			if err != nil {

				return nil, resp, err

			}

			results := make([]*netbox.VirtualMachineWithConfigContext, len(list.Results))

			for i := range list.Results {

				results[i] = &list.Results[i]

			}

			return results, resp, nil

		},

		ToBriefRequest: func(vm *netbox.VirtualMachineWithConfigContext) netbox.BriefVirtualMachineRequest {

			return netbox.BriefVirtualMachineRequest{

				Name: vm.GetName(),
			}

		},
	}

}

// LookupVirtualMachine looks up a Virtual Machine by ID or name.

func LookupVirtualMachine(ctx context.Context, client *netbox.APIClient, value string) (*netbox.BriefVirtualMachineRequest, diag.Diagnostics) {

	return GenericLookup(ctx, value, VirtualMachineLookupConfig(client))

}

// =====================================================

// CIRCUITS LOOKUP CONFIGURATIONS

// =====================================================

// ProviderLookupConfig returns the lookup configuration for circuit Providers.

func ProviderLookupConfig(client *netbox.APIClient) LookupConfig[*netbox.Provider, netbox.BriefProviderRequest] {

	return LookupConfig[*netbox.Provider, netbox.BriefProviderRequest]{

		ResourceName: "Provider",

		RetrieveByID: func(ctx context.Context, id int32) (*netbox.Provider, *http.Response, error) {

			return client.CircuitsAPI.CircuitsProvidersRetrieve(ctx, id).Execute()

		},

		ListBySlug: func(ctx context.Context, slug string) ([]*netbox.Provider, *http.Response, error) {

			list, resp, err := client.CircuitsAPI.CircuitsProvidersList(ctx).Slug([]string{slug}).Execute()

			if err != nil {

				return nil, resp, err

			}

			results := make([]*netbox.Provider, len(list.Results))

			for i := range list.Results {

				results[i] = &list.Results[i]

			}

			return results, resp, nil

		},

		ToBriefRequest: func(p *netbox.Provider) netbox.BriefProviderRequest {

			return netbox.BriefProviderRequest{

				Name: p.GetName(),

				Slug: p.GetSlug(),
			}

		},
	}

}

// LookupProvider looks up a circuit Provider by ID or slug.

func LookupProvider(ctx context.Context, client *netbox.APIClient, value string) (*netbox.BriefProviderRequest, diag.Diagnostics) {

	return GenericLookup(ctx, value, ProviderLookupConfig(client))

}

// LookupProviderAccount looks up a provider account by ID or account string, scoped to a provider.

func LookupProviderAccount(ctx context.Context, client *netbox.APIClient, providerID int32, value string) (*netbox.BriefProviderAccountRequest, diag.Diagnostics) {
	var id int32

	if _, err := fmt.Sscanf(value, "%d", &id); err == nil {
		account, resp, err := client.CircuitsAPI.CircuitsProviderAccountsRetrieve(ctx, id).Execute()
		defer utils.CloseResponseBody(resp)
		if err != nil || resp == nil || resp.StatusCode != http.StatusOK {
			errMsg := unknownErrorMsg
			if err != nil {
				errMsg = err.Error()
			}
			return nil, diag.Diagnostics{diag.NewErrorDiagnostic(
				"Provider Account lookup failed",
				fmt.Sprintf("Could not find Provider Account with ID %d: %s", id, errMsg),
			)}
		}
		brief := netbox.BriefProviderAccountRequest{Account: account.GetAccount()}
		if name := account.GetName(); name != "" {
			brief.Name = &name
		}
		return &brief, nil
	}

	listReq := client.CircuitsAPI.CircuitsProviderAccountsList(ctx).
		Account([]string{value}).
		ProviderId([]int32{providerID})
	list, resp, err := listReq.Execute()
	defer utils.CloseResponseBody(resp)
	if err != nil {
		return nil, diag.Diagnostics{diag.NewErrorDiagnostic(
			"Provider Account lookup failed",
			fmt.Sprintf("Could not list Provider Accounts with account '%s': %s", value, err.Error()),
		)}
	}
	if list == nil || len(list.Results) == 0 {
		return nil, diag.Diagnostics{diag.NewErrorDiagnostic(
			"Provider Account lookup failed",
			fmt.Sprintf("No Provider Account found with account '%s' for provider ID %d", value, providerID),
		)}
	}
	if len(list.Results) > 1 {
		return nil, diag.Diagnostics{diag.NewErrorDiagnostic(
			"Provider Account lookup failed",
			fmt.Sprintf("Multiple Provider Accounts found with account '%s' for provider ID %d", value, providerID),
		)}
	}

	account := list.Results[0]
	brief := netbox.BriefProviderAccountRequest{Account: account.GetAccount()}
	if name := account.GetName(); name != "" {
		brief.Name = &name
	}
	return &brief, nil
}

// CircuitTypeLookupConfig returns the lookup configuration for Circuit Types.

func CircuitTypeLookupConfig(client *netbox.APIClient) LookupConfig[*netbox.CircuitType, netbox.BriefCircuitTypeRequest] {

	return LookupConfig[*netbox.CircuitType, netbox.BriefCircuitTypeRequest]{

		ResourceName: "Circuit Type",

		RetrieveByID: func(ctx context.Context, id int32) (*netbox.CircuitType, *http.Response, error) {

			return client.CircuitsAPI.CircuitsCircuitTypesRetrieve(ctx, id).Execute()

		},

		ListBySlug: func(ctx context.Context, slug string) ([]*netbox.CircuitType, *http.Response, error) {

			list, resp, err := client.CircuitsAPI.CircuitsCircuitTypesList(ctx).Slug([]string{slug}).Execute()

			if err != nil {

				return nil, resp, err

			}

			results := make([]*netbox.CircuitType, len(list.Results))

			for i := range list.Results {

				results[i] = &list.Results[i]

			}

			return results, resp, nil

		},

		ToBriefRequest: func(ct *netbox.CircuitType) netbox.BriefCircuitTypeRequest {

			return netbox.BriefCircuitTypeRequest{

				Name: ct.GetName(),

				Slug: ct.GetSlug(),
			}

		},
	}

}

// LookupCircuitType looks up a Circuit Type by ID or slug.

func LookupCircuitType(ctx context.Context, client *netbox.APIClient, value string) (*netbox.BriefCircuitTypeRequest, diag.Diagnostics) {

	return GenericLookup(ctx, value, CircuitTypeLookupConfig(client))

}

// ContactGroupLookupConfig returns the lookup configuration for Contact Groups.

func ContactGroupLookupConfig(client *netbox.APIClient) LookupConfig[*netbox.ContactGroup, netbox.BriefContactGroupRequest] {

	return LookupConfig[*netbox.ContactGroup, netbox.BriefContactGroupRequest]{

		ResourceName: "Contact Group",

		RetrieveByID: func(ctx context.Context, id int32) (*netbox.ContactGroup, *http.Response, error) {

			return client.TenancyAPI.TenancyContactGroupsRetrieve(ctx, id).Execute()

		},

		ListBySlug: func(ctx context.Context, slug string) ([]*netbox.ContactGroup, *http.Response, error) {

			// Try lookup by slug first

			list, resp, err := client.TenancyAPI.TenancyContactGroupsList(ctx).Slug([]string{slug}).Execute()

			if err != nil {

				return nil, resp, err

			}

			if len(list.Results) > 0 {

				results := make([]*netbox.ContactGroup, len(list.Results))

				for i := range list.Results {

					results[i] = &list.Results[i]

				}

				return results, resp, nil

			}

			// If no results by slug, try lookup by name

			list, resp, err = client.TenancyAPI.TenancyContactGroupsList(ctx).Name([]string{slug}).Execute()

			if err != nil {

				return nil, resp, err

			}

			results := make([]*netbox.ContactGroup, len(list.Results))

			for i := range list.Results {

				results[i] = &list.Results[i]

			}

			return results, resp, nil

		},

		ToBriefRequest: func(cg *netbox.ContactGroup) netbox.BriefContactGroupRequest {

			return netbox.BriefContactGroupRequest{

				Name: cg.GetName(),

				Slug: cg.GetSlug(),
			}

		},
	}

}

// LookupContactGroup looks up a Contact Group by ID or slug.

func LookupContactGroup(ctx context.Context, client *netbox.APIClient, value string) (*netbox.BriefContactGroupRequest, diag.Diagnostics) {

	return GenericLookup(ctx, value, ContactGroupLookupConfig(client))

}

// LookupContactGroupID looks up a Contact Group by ID or slug and returns the ID.

func LookupContactGroupID(ctx context.Context, client *netbox.APIClient, value string) (int32, diag.Diagnostics) {

	return GenericLookupID(ctx, value, ContactGroupLookupConfig(client), func(cg *netbox.ContactGroup) int32 {

		return cg.GetId()

	})

}

// =====================================================

// IPAM RIR LOOKUPS

// =====================================================

// RIRLookupConfig returns the lookup configuration for RIRs (Regional Internet Registries).

func RIRLookupConfig(client *netbox.APIClient) LookupConfig[*netbox.RIR, netbox.BriefRIRRequest] {

	return LookupConfig[*netbox.RIR, netbox.BriefRIRRequest]{

		ResourceName: "RIR",

		RetrieveByID: func(ctx context.Context, id int32) (*netbox.RIR, *http.Response, error) {

			return client.IpamAPI.IpamRirsRetrieve(ctx, id).Execute()

		},

		ListBySlug: func(ctx context.Context, slug string) ([]*netbox.RIR, *http.Response, error) {

			list, resp, err := client.IpamAPI.IpamRirsList(ctx).Slug([]string{slug}).Execute()

			if err != nil {

				return nil, resp, err

			}

			results := make([]*netbox.RIR, len(list.Results))

			for i := range list.Results {

				results[i] = &list.Results[i]

			}

			return results, resp, nil

		},

		ToBriefRequest: func(r *netbox.RIR) netbox.BriefRIRRequest {

			return netbox.BriefRIRRequest{

				Name: r.GetName(),

				Slug: r.GetSlug(),
			}

		},
	}

}

// LookupRIR looks up a RIR by ID or slug.

func LookupRIR(ctx context.Context, client *netbox.APIClient, value string) (*netbox.BriefRIRRequest, diag.Diagnostics) {

	return GenericLookup(ctx, value, RIRLookupConfig(client))

}

// =====================================================

// CIRCUIT LOOKUPS

// =====================================================

// CircuitLookupConfig returns the lookup configuration for Circuits.

func CircuitLookupConfig(client *netbox.APIClient) LookupConfig[*netbox.Circuit, netbox.BriefCircuitRequest] {

	return LookupConfig[*netbox.Circuit, netbox.BriefCircuitRequest]{

		ResourceName: "Circuit",

		RetrieveByID: func(ctx context.Context, id int32) (*netbox.Circuit, *http.Response, error) {

			return client.CircuitsAPI.CircuitsCircuitsRetrieve(ctx, id).Execute()

		},

		ListBySlug: func(ctx context.Context, cid string) ([]*netbox.Circuit, *http.Response, error) {

			// Circuits use CID (circuit ID) instead of slug

			list, resp, err := client.CircuitsAPI.CircuitsCircuitsList(ctx).Cid([]string{cid}).Execute()

			if err != nil {

				return nil, resp, err

			}

			results := make([]*netbox.Circuit, len(list.Results))

			for i := range list.Results {

				results[i] = &list.Results[i]

			}

			return results, resp, nil

		},

		ToBriefRequest: func(c *netbox.Circuit) netbox.BriefCircuitRequest {

			return netbox.BriefCircuitRequest{

				Cid: c.GetCid(),

				Provider: *netbox.NewBriefProviderRequest(c.Provider.GetName(), c.Provider.GetSlug()),
			}

		},
	}

}

// LookupCircuit looks up a Circuit by ID or CID.

func LookupCircuit(ctx context.Context, client *netbox.APIClient, value string) (*netbox.BriefCircuitRequest, diag.Diagnostics) {

	return GenericLookup(ctx, value, CircuitLookupConfig(client))

}

// CircuitGroupLookupConfig returns the lookup configuration for Circuit Groups.

func CircuitGroupLookupConfig(client *netbox.APIClient) LookupConfig[*netbox.CircuitGroup, netbox.BriefCircuitGroupRequest] {

	return LookupConfig[*netbox.CircuitGroup, netbox.BriefCircuitGroupRequest]{

		ResourceName: "Circuit Group",

		RetrieveByID: func(ctx context.Context, id int32) (*netbox.CircuitGroup, *http.Response, error) {

			return client.CircuitsAPI.CircuitsCircuitGroupsRetrieve(ctx, id).Execute()

		},

		ListBySlug: func(ctx context.Context, slug string) ([]*netbox.CircuitGroup, *http.Response, error) {

			list, resp, err := client.CircuitsAPI.CircuitsCircuitGroupsList(ctx).Slug([]string{slug}).Execute()

			if err != nil {

				return nil, resp, err

			}

			results := make([]*netbox.CircuitGroup, len(list.Results))

			for i := range list.Results {

				results[i] = &list.Results[i]

			}

			return results, resp, nil

		},

		ToBriefRequest: func(cg *netbox.CircuitGroup) netbox.BriefCircuitGroupRequest {

			return netbox.BriefCircuitGroupRequest{

				Name: cg.GetName(),
			}

		},
	}

}

// LookupCircuitGroup looks up a Circuit Group by ID or slug.

func LookupCircuitGroup(ctx context.Context, client *netbox.APIClient, value string) (*netbox.CircuitGroup, diag.Diagnostics) {

	var id int32

	if _, err := fmt.Sscanf(value, "%d", &id); err == nil {

		// Lookup by ID

		resource, resp, err := client.CircuitsAPI.CircuitsCircuitGroupsRetrieve(ctx, id).Execute()

		defer utils.CloseResponseBody(resp)

		if err != nil || resp.StatusCode != 200 {

			errMsg := unknownErrorMsg

			if err != nil {

				errMsg = err.Error()

			}

			return nil, diag.Diagnostics{diag.NewErrorDiagnostic(

				"Circuit Group lookup failed",

				fmt.Sprintf("Could not find Circuit Group with ID %d: %s", id, errMsg),
			)}

		}

		return resource, nil

	}

	// Lookup by slug first

	list, resp, err := client.CircuitsAPI.CircuitsCircuitGroupsList(ctx).Slug([]string{value}).Execute()

	defer utils.CloseResponseBody(resp)

	if err == nil && resp.StatusCode == 200 && len(list.Results) > 0 {

		return &list.Results[0], nil

	}

	// Try lookup by name

	list, resp, err = client.CircuitsAPI.CircuitsCircuitGroupsList(ctx).Name([]string{value}).Execute()

	defer utils.CloseResponseBody(resp)

	if err != nil || resp.StatusCode != 200 {

		errMsg := unknownErrorMsg

		if err != nil {

			errMsg = err.Error()

		}

		return nil, diag.Diagnostics{diag.NewErrorDiagnostic(

			"Circuit Group lookup failed",

			fmt.Sprintf("Could not find Circuit Group with slug or name '%s': %s", value, errMsg),
		)}

	}

	if len(list.Results) == 0 {

		return nil, diag.Diagnostics{diag.NewErrorDiagnostic(

			"Circuit Group lookup failed",

			fmt.Sprintf("No Circuit Group found with slug or name '%s'", value),
		)}

	}

	return &list.Results[0], nil

}

// =====================================================

// WIRELESS LOOKUPS

// =====================================================

// WirelessLANGroupLookupConfig returns the lookup configuration for Wireless LAN Groups.

func WirelessLANGroupLookupConfig(client *netbox.APIClient) LookupConfig[*netbox.WirelessLANGroup, netbox.BriefWirelessLANGroupRequest] {

	return LookupConfig[*netbox.WirelessLANGroup, netbox.BriefWirelessLANGroupRequest]{

		ResourceName: "Wireless LAN Group",

		RetrieveByID: func(ctx context.Context, id int32) (*netbox.WirelessLANGroup, *http.Response, error) {

			return client.WirelessAPI.WirelessWirelessLanGroupsRetrieve(ctx, id).Execute()

		},

		ListBySlug: func(ctx context.Context, slug string) ([]*netbox.WirelessLANGroup, *http.Response, error) {

			list, resp, err := client.WirelessAPI.WirelessWirelessLanGroupsList(ctx).Slug([]string{slug}).Execute()

			if err != nil {

				return nil, resp, err

			}

			results := make([]*netbox.WirelessLANGroup, len(list.Results))

			for i := range list.Results {

				results[i] = &list.Results[i]

			}

			return results, resp, nil

		},

		ToBriefRequest: func(g *netbox.WirelessLANGroup) netbox.BriefWirelessLANGroupRequest {

			return netbox.BriefWirelessLANGroupRequest{

				Name: g.GetName(),

				Slug: g.GetSlug(),
			}

		},
	}

}

// LookupWirelessLANGroup looks up a Wireless LAN Group by ID or slug.

func LookupWirelessLANGroup(ctx context.Context, client *netbox.APIClient, value string) (*netbox.BriefWirelessLANGroupRequest, diag.Diagnostics) {

	return GenericLookup(ctx, value, WirelessLANGroupLookupConfig(client))

}

// =====================================================

// INVENTORY ITEM ROLE LOOKUPS

// =====================================================

// InventoryItemRoleLookupConfig returns the lookup configuration for Inventory Item Roles.

func InventoryItemRoleLookupConfig(client *netbox.APIClient) LookupConfig[*netbox.InventoryItemRole, netbox.BriefInventoryItemRoleRequest] {

	return LookupConfig[*netbox.InventoryItemRole, netbox.BriefInventoryItemRoleRequest]{

		ResourceName: "Inventory Item Role",

		RetrieveByID: func(ctx context.Context, id int32) (*netbox.InventoryItemRole, *http.Response, error) {

			return client.DcimAPI.DcimInventoryItemRolesRetrieve(ctx, id).Execute()

		},

		ListBySlug: func(ctx context.Context, slug string) ([]*netbox.InventoryItemRole, *http.Response, error) {

			list, resp, err := client.DcimAPI.DcimInventoryItemRolesList(ctx).Slug([]string{slug}).Execute()

			if err != nil {

				return nil, resp, err

			}

			results := make([]*netbox.InventoryItemRole, len(list.Results))

			for i := range list.Results {

				results[i] = &list.Results[i]

			}

			return results, resp, nil

		},

		ToBriefRequest: func(r *netbox.InventoryItemRole) netbox.BriefInventoryItemRoleRequest {

			return netbox.BriefInventoryItemRoleRequest{

				Name: r.GetName(),

				Slug: r.GetSlug(),
			}

		},
	}

}

// LookupInventoryItemRole looks up an Inventory Item Role by ID or slug.

func LookupInventoryItemRole(ctx context.Context, client *netbox.APIClient, value string) (*netbox.BriefInventoryItemRoleRequest, diag.Diagnostics) {

	return GenericLookup(ctx, value, InventoryItemRoleLookupConfig(client))

}

// =====================================================

// MODULE TYPE LOOKUPS

// =====================================================

// ModuleTypeLookupConfig returns the lookup configuration for Module Types.

func ModuleTypeLookupConfig(client *netbox.APIClient) LookupConfig[*netbox.ModuleType, netbox.BriefModuleTypeRequest] {

	return LookupConfig[*netbox.ModuleType, netbox.BriefModuleTypeRequest]{

		ResourceName: "Module Type",

		RetrieveByID: func(ctx context.Context, id int32) (*netbox.ModuleType, *http.Response, error) {

			return client.DcimAPI.DcimModuleTypesRetrieve(ctx, id).Execute()

		},

		ListBySlug: func(ctx context.Context, model string) ([]*netbox.ModuleType, *http.Response, error) {

			// Module types don't have slugs, so we search by model name

			list, resp, err := client.DcimAPI.DcimModuleTypesList(ctx).Model([]string{model}).Execute()

			if err != nil {

				return nil, resp, err

			}

			results := make([]*netbox.ModuleType, len(list.Results))

			for i := range list.Results {

				results[i] = &list.Results[i]

			}

			return results, resp, nil

		},

		ToBriefRequest: func(mt *netbox.ModuleType) netbox.BriefModuleTypeRequest {

			mfr := mt.GetManufacturer()

			return netbox.BriefModuleTypeRequest{

				Manufacturer: netbox.BriefManufacturerRequest{

					Name: mfr.GetName(),

					Slug: mfr.GetSlug(),
				},

				Model: mt.GetModel(),
			}

		},
	}

}

// LookupModuleType looks up a Module Type by ID or model name.

func LookupModuleType(ctx context.Context, client *netbox.APIClient, value string) (*netbox.BriefModuleTypeRequest, diag.Diagnostics) {

	return GenericLookup(ctx, value, ModuleTypeLookupConfig(client))

}

// =====================================================

// USER LOOKUPS

// =====================================================

// UserLookupConfig returns the lookup configuration for Users.

func UserLookupConfig(client *netbox.APIClient) LookupConfig[*netbox.User, netbox.BriefUserRequest] {

	return LookupConfig[*netbox.User, netbox.BriefUserRequest]{

		ResourceName: "User",

		RetrieveByID: func(ctx context.Context, id int32) (*netbox.User, *http.Response, error) {

			return client.UsersAPI.UsersUsersRetrieve(ctx, id).Execute()

		},

		ListBySlug: func(ctx context.Context, username string) ([]*netbox.User, *http.Response, error) {

			// Users don't have slugs, so we search by username

			list, resp, err := client.UsersAPI.UsersUsersList(ctx).Username([]string{username}).Execute()

			if err != nil {

				return nil, resp, err

			}

			results := make([]*netbox.User, len(list.Results))

			for i := range list.Results {

				results[i] = &list.Results[i]

			}

			return results, resp, nil

		},

		ToBriefRequest: func(u *netbox.User) netbox.BriefUserRequest {

			return netbox.BriefUserRequest{

				Username: u.GetUsername(),
			}

		},
	}

}

// LookupUser looks up a User by ID or username.

func LookupUser(ctx context.Context, client *netbox.APIClient, value string) (*netbox.BriefUserRequest, diag.Diagnostics) {

	return GenericLookup(ctx, value, UserLookupConfig(client))

}

// =====================================================

// IP ADDRESS LOOKUPS

// =====================================================

// IPAddressLookupConfig returns the lookup configuration for IP Addresses.

func IPAddressLookupConfig(client *netbox.APIClient) LookupConfig[*netbox.IPAddress, netbox.BriefIPAddressRequest] {

	return LookupConfig[*netbox.IPAddress, netbox.BriefIPAddressRequest]{

		ResourceName: "IP Address",

		RetrieveByID: func(ctx context.Context, id int32) (*netbox.IPAddress, *http.Response, error) {

			return client.IpamAPI.IpamIpAddressesRetrieve(ctx, id).Execute()

		},

		ListBySlug: func(ctx context.Context, address string) ([]*netbox.IPAddress, *http.Response, error) {

			// IP addresses don't have slugs, so we search by address

			list, resp, err := client.IpamAPI.IpamIpAddressesList(ctx).Address([]string{address}).Execute()

			if err != nil {

				return nil, resp, err

			}

			results := make([]*netbox.IPAddress, len(list.Results))

			for i := range list.Results {

				results[i] = &list.Results[i]

			}

			return results, resp, nil

		},

		ToBriefRequest: func(ip *netbox.IPAddress) netbox.BriefIPAddressRequest {

			return netbox.BriefIPAddressRequest{

				Address: ip.GetAddress(),
			}

		},
	}

}

// LookupIPAddress looks up an IP Address by ID or address string.

func LookupIPAddress(ctx context.Context, client *netbox.APIClient, value string) (*netbox.BriefIPAddressRequest, diag.Diagnostics) {

	return GenericLookup(ctx, value, IPAddressLookupConfig(client))

}

// PowerPortLookupConfig returns the lookup configuration for Power Ports.

func PowerPortLookupConfig(client *netbox.APIClient) LookupConfig[*netbox.PowerPort, netbox.BriefPowerPortRequest] {

	return LookupConfig[*netbox.PowerPort, netbox.BriefPowerPortRequest]{

		ResourceName: "Power Port",

		RetrieveByID: func(ctx context.Context, id int32) (*netbox.PowerPort, *http.Response, error) {

			return client.DcimAPI.DcimPowerPortsRetrieve(ctx, id).Execute()

		},

		ListBySlug: func(ctx context.Context, name string) ([]*netbox.PowerPort, *http.Response, error) {

			// PowerPort uses name instead of slug

			list, resp, err := client.DcimAPI.DcimPowerPortsList(ctx).Name([]string{name}).Execute()

			if err != nil {

				return nil, resp, err

			}

			results := make([]*netbox.PowerPort, len(list.Results))

			for i := range list.Results {

				results[i] = &list.Results[i]

			}

			return results, resp, nil

		},

		ToBriefRequest: func(p *netbox.PowerPort) netbox.BriefPowerPortRequest {

			device := p.GetDevice()
			deviceName := device.GetName()
			powerPortName := p.GetName()

			deviceReq := netbox.BriefDeviceRequest{}
			deviceReq.Name.Set(&deviceName)

			return netbox.BriefPowerPortRequest{
				Device: deviceReq,
				Name:   powerPortName,
			}

		},
	}

}

// LookupPowerPort looks up a Power Port by ID or name.

func LookupPowerPort(ctx context.Context, client *netbox.APIClient, value string) (*netbox.BriefPowerPortRequest, diag.Diagnostics) {

	return GenericLookup(ctx, value, PowerPortLookupConfig(client))

}
