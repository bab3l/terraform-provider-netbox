package netboxlookup

import (
	"context"
	"fmt"

	"github.com/bab3l/go-netbox"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// LookupManufacturerBrief returns a BriefManufacturerRequest from an ID or slug
func LookupManufacturerBrief(ctx context.Context, client *netbox.APIClient, value string) (*netbox.BriefManufacturerRequest, diag.Diagnostics) {
	var id int32
	if _, err := fmt.Sscanf(value, "%d", &id); err == nil {
		resource, resp, err := client.DcimAPI.DcimManufacturersRetrieve(ctx, id).Execute()
		if err != nil || resp.StatusCode != 200 {
			return nil, diag.Diagnostics{diag.NewErrorDiagnostic("Manufacturer lookup failed", err.Error())}
		}
		return &netbox.BriefManufacturerRequest{
			Name: resource.GetName(),
			Slug: resource.GetSlug(),
		}, nil
	}
	// Lookup by slug
	list, resp, err := client.DcimAPI.DcimManufacturersList(ctx).Slug([]string{value}).Execute()
	if err != nil || resp.StatusCode != 200 {
		return nil, diag.Diagnostics{diag.NewErrorDiagnostic("Manufacturer lookup failed", fmt.Sprintf("Could not find manufacturer with slug '%s': %v", value, err))}
	}
	if list != nil && len(list.Results) > 0 {
		resource := list.Results[0]
		return &netbox.BriefManufacturerRequest{
			Name: resource.GetName(),
			Slug: resource.GetSlug(),
		}, nil
	}
	return nil, diag.Diagnostics{diag.NewErrorDiagnostic("Manufacturer lookup failed", fmt.Sprintf("No manufacturer found with slug '%s'", value))}
}

// LookupTenantBrief returns a BriefTenantRequest from an ID or slug
func LookupTenantBrief(ctx context.Context, client *netbox.APIClient, value string) (*netbox.BriefTenantRequest, diag.Diagnostics) {
	var id int32
	if _, err := fmt.Sscanf(value, "%d", &id); err == nil {
		resource, resp, err := client.TenancyAPI.TenancyTenantsRetrieve(ctx, id).Execute()
		if err != nil || resp.StatusCode != 200 {
			return nil, diag.Diagnostics{diag.NewErrorDiagnostic("Tenant lookup failed", err.Error())}
		}
		return &netbox.BriefTenantRequest{
			Name: resource.GetName(),
			Slug: resource.GetSlug(),
		}, nil
	}
	// Optionally, lookup by slug or name if not an ID
	return nil, diag.Diagnostics{diag.NewErrorDiagnostic("Tenant lookup failed", "Invalid input")}
}

// LookupTenantGroupBrief returns a BriefTenantGroupRequest from an ID or slug
func LookupTenantGroupBrief(ctx context.Context, client *netbox.APIClient, value string) (*netbox.BriefTenantGroupRequest, diag.Diagnostics) {
	var id int32
	if _, err := fmt.Sscanf(value, "%d", &id); err == nil {
		resource, resp, err := client.TenancyAPI.TenancyTenantGroupsRetrieve(ctx, id).Execute()
		if err != nil || resp.StatusCode != 200 {
			return nil, diag.Diagnostics{diag.NewErrorDiagnostic("Tenant group lookup failed", err.Error())}
		}
		return &netbox.BriefTenantGroupRequest{
			Name: resource.GetName(),
			Slug: resource.GetSlug(),
		}, nil
	}
	return nil, diag.Diagnostics{diag.NewErrorDiagnostic("Tenant group lookup failed", "Invalid input")}
}

// LookupRegionBrief returns a BriefRegionRequest from an ID or slug
func LookupRegionBrief(ctx context.Context, client *netbox.APIClient, value string) (*netbox.BriefRegionRequest, diag.Diagnostics) {
	var id int32
	if _, err := fmt.Sscanf(value, "%d", &id); err == nil {
		resource, resp, err := client.DcimAPI.DcimRegionsRetrieve(ctx, id).Execute()
		if err != nil || resp.StatusCode != 200 {
			return nil, diag.Diagnostics{diag.NewErrorDiagnostic("Region lookup failed", err.Error())}
		}
		return &netbox.BriefRegionRequest{
			Name: resource.GetName(),
			Slug: resource.GetSlug(),
		}, nil
	}
	return nil, diag.Diagnostics{diag.NewErrorDiagnostic("Region lookup failed", "Invalid input")}
}

// LookupSiteGroupBrief returns a BriefSiteGroupRequest from an ID or slug
func LookupSiteGroupBrief(ctx context.Context, client *netbox.APIClient, value string) (*netbox.BriefSiteGroupRequest, diag.Diagnostics) {
	var id int32
	if _, err := fmt.Sscanf(value, "%d", &id); err == nil {
		resource, resp, err := client.DcimAPI.DcimSiteGroupsRetrieve(ctx, id).Execute()
		if err != nil || resp.StatusCode != 200 {
			return nil, diag.Diagnostics{diag.NewErrorDiagnostic("Site group lookup failed", err.Error())}
		}
		return &netbox.BriefSiteGroupRequest{
			Name: resource.GetName(),
			Slug: resource.GetSlug(),
		}, nil
	}
	return nil, diag.Diagnostics{diag.NewErrorDiagnostic("Site group lookup failed", "Invalid input")}
}

// LookupSiteBrief returns a BriefSiteRequest from an ID or slug
func LookupSiteBrief(ctx context.Context, client *netbox.APIClient, value string) (*netbox.BriefSiteRequest, diag.Diagnostics) {
	var id int32
	if _, err := fmt.Sscanf(value, "%d", &id); err == nil {
		resource, resp, err := client.DcimAPI.DcimSitesRetrieve(ctx, id).Execute()
		if err != nil || resp.StatusCode != 200 {
			return nil, diag.Diagnostics{diag.NewErrorDiagnostic("Site lookup failed", err.Error())}
		}
		return &netbox.BriefSiteRequest{
			Name: resource.GetName(),
			Slug: resource.GetSlug(),
		}, nil
	}
	// Lookup by slug
	list, resp, err := client.DcimAPI.DcimSitesList(ctx).Slug([]string{value}).Execute()
	if err != nil || resp.StatusCode != 200 {
		return nil, diag.Diagnostics{diag.NewErrorDiagnostic("Site lookup failed", fmt.Sprintf("Could not find site with slug '%s': %v", value, err))}
	}
	if list != nil && len(list.Results) > 0 {
		resource := list.Results[0]
		return &netbox.BriefSiteRequest{
			Name: resource.GetName(),
			Slug: resource.GetSlug(),
		}, nil
	}
	return nil, diag.Diagnostics{diag.NewErrorDiagnostic("Site lookup failed", fmt.Sprintf("No site found with slug '%s'", value))}
}

// LookupLocationBrief returns a BriefLocationRequest from an ID or slug
func LookupLocationBrief(ctx context.Context, client *netbox.APIClient, value string) (*netbox.BriefLocationRequest, diag.Diagnostics) {
	var id int32
	if _, err := fmt.Sscanf(value, "%d", &id); err == nil {
		resource, resp, err := client.DcimAPI.DcimLocationsRetrieve(ctx, id).Execute()
		if err != nil || resp.StatusCode != 200 {
			return nil, diag.Diagnostics{diag.NewErrorDiagnostic("Location lookup failed", err.Error())}
		}
		return &netbox.BriefLocationRequest{
			Name: resource.GetName(),
			Slug: resource.GetSlug(),
		}, nil
	}
	// Lookup by slug
	list, resp, err := client.DcimAPI.DcimLocationsList(ctx).Slug([]string{value}).Execute()
	if err != nil || resp.StatusCode != 200 {
		return nil, diag.Diagnostics{diag.NewErrorDiagnostic("Location lookup failed", fmt.Sprintf("Could not find location with slug '%s': %v", value, err))}
	}
	if list != nil && len(list.Results) > 0 {
		resource := list.Results[0]
		return &netbox.BriefLocationRequest{
			Name: resource.GetName(),
			Slug: resource.GetSlug(),
		}, nil
	}
	return nil, diag.Diagnostics{diag.NewErrorDiagnostic("Location lookup failed", fmt.Sprintf("No location found with slug '%s'", value))}
}

// LookupRackRoleBrief returns a BriefRackRoleRequest from an ID or slug
func LookupRackRoleBrief(ctx context.Context, client *netbox.APIClient, value string) (*netbox.BriefRackRoleRequest, diag.Diagnostics) {
	var id int32
	if _, err := fmt.Sscanf(value, "%d", &id); err == nil {
		resource, resp, err := client.DcimAPI.DcimRackRolesRetrieve(ctx, id).Execute()
		if err != nil || resp.StatusCode != 200 {
			return nil, diag.Diagnostics{diag.NewErrorDiagnostic("Rack role lookup failed", err.Error())}
		}
		return &netbox.BriefRackRoleRequest{
			Name: resource.GetName(),
			Slug: resource.GetSlug(),
		}, nil
	}
	// Lookup by slug
	list, resp, err := client.DcimAPI.DcimRackRolesList(ctx).Slug([]string{value}).Execute()
	if err != nil || resp.StatusCode != 200 {
		return nil, diag.Diagnostics{diag.NewErrorDiagnostic("Rack role lookup failed", fmt.Sprintf("Could not find rack role with slug '%s': %v", value, err))}
	}
	if list != nil && len(list.Results) > 0 {
		resource := list.Results[0]
		return &netbox.BriefRackRoleRequest{
			Name: resource.GetName(),
			Slug: resource.GetSlug(),
		}, nil
	}
	return nil, diag.Diagnostics{diag.NewErrorDiagnostic("Rack role lookup failed", fmt.Sprintf("No rack role found with slug '%s'", value))}
}

// LookupRackTypeBrief returns a BriefRackTypeRequest from an ID or model name
func LookupRackTypeBrief(ctx context.Context, client *netbox.APIClient, value string) (*netbox.BriefRackTypeRequest, diag.Diagnostics) {
	var id int32
	if _, err := fmt.Sscanf(value, "%d", &id); err == nil {
		resource, resp, err := client.DcimAPI.DcimRackTypesRetrieve(ctx, id).Execute()
		if err != nil || resp.StatusCode != 200 {
			return nil, diag.Diagnostics{diag.NewErrorDiagnostic("Rack type lookup failed", err.Error())}
		}
		// Get manufacturer for the request
		manufacturer := resource.GetManufacturer()
		manufacturerRequest := netbox.BriefManufacturerRequest{
			Name: manufacturer.GetName(),
			Slug: manufacturer.GetSlug(),
		}
		return &netbox.BriefRackTypeRequest{
			Manufacturer: manufacturerRequest,
			Model:        resource.GetModel(),
			Slug:         resource.GetSlug(),
		}, nil
	}
	// Lookup by model name
	list, resp, err := client.DcimAPI.DcimRackTypesList(ctx).Model([]string{value}).Execute()
	if err != nil || resp.StatusCode != 200 {
		return nil, diag.Diagnostics{diag.NewErrorDiagnostic("Rack type lookup failed", fmt.Sprintf("Could not find rack type with model '%s': %v", value, err))}
	}
	if list != nil && len(list.Results) > 0 {
		resource := list.Results[0]
		manufacturer := resource.GetManufacturer()
		manufacturerRequest := netbox.BriefManufacturerRequest{
			Name: manufacturer.GetName(),
			Slug: manufacturer.GetSlug(),
		}
		return &netbox.BriefRackTypeRequest{
			Manufacturer: manufacturerRequest,
			Model:        resource.GetModel(),
			Slug:         resource.GetSlug(),
		}, nil
	}
	return nil, diag.Diagnostics{diag.NewErrorDiagnostic("Rack type lookup failed", fmt.Sprintf("No rack type found with model '%s'", value))}
}
