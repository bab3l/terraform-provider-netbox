package netboxlookup

import (
	"context"
	"fmt"

	"github.com/bab3l/go-netbox"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

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
	return nil, diag.Diagnostics{diag.NewErrorDiagnostic("Site lookup failed", "Invalid input")}
}
