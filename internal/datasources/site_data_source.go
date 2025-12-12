// Package datasources contains Terraform data source implementations for the Netbox provider.
//
// This package integrates with the go-netbox OpenAPI client to provide
// read-only access to Netbox resources via Terraform data sources.
package datasources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bab3l/go-netbox"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &SiteDataSource{}

func NewSiteDataSource() datasource.DataSource {
	return &SiteDataSource{}
}

// SiteDataSource defines the data source implementation.
type SiteDataSource struct {
	client *netbox.APIClient
}

// SiteDataSourceModel describes the data source data model.
type SiteDataSourceModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Slug         types.String `tfsdk:"slug"`
	Status       types.String `tfsdk:"status"`
	Region       types.String `tfsdk:"region"`
	RegionID     types.String `tfsdk:"region_id"`
	Group        types.String `tfsdk:"group"`
	GroupID      types.String `tfsdk:"group_id"`
	Tenant       types.String `tfsdk:"tenant"`
	TenantID     types.String `tfsdk:"tenant_id"`
	Facility     types.String `tfsdk:"facility"`
	Description  types.String `tfsdk:"description"`
	Comments     types.String `tfsdk:"comments"`
	Tags         types.Set    `tfsdk:"tags"`
	CustomFields types.Set    `tfsdk:"custom_fields"`
}

func (d *SiteDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_site"
}

func (d *SiteDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get information about a site in Netbox. Sites represent physical locations such as data centers, offices, or other facilities where network infrastructure is deployed. You can identify the site using `id`, `slug`, or `name`.",

		Attributes: map[string]schema.Attribute{
			"id":            nbschema.DSIDAttribute("site"),
			"name":          nbschema.DSNameAttribute("site"),
			"slug":          nbschema.DSSlugAttribute("site"),
			"status":        nbschema.DSComputedStringAttribute("Operational status of the site (e.g., `planned`, `staging`, `active`, `decommissioning`, `retired`)."),
			"region":        nbschema.DSComputedStringAttribute("Name of the region where this site is located."),
			"region_id":     nbschema.DSComputedStringAttribute("ID of the region where this site is located."),
			"group":         nbschema.DSComputedStringAttribute("Name of the site group."),
			"group_id":      nbschema.DSComputedStringAttribute("ID of the site group."),
			"tenant":        nbschema.DSComputedStringAttribute("Name of the tenant that owns this site."),
			"tenant_id":     nbschema.DSComputedStringAttribute("ID of the tenant that owns this site."),
			"facility":      nbschema.DSComputedStringAttribute("Local facility identifier or description."),
			"description":   nbschema.DSComputedStringAttribute("Detailed description of the site."),
			"comments":      nbschema.DSComputedStringAttribute("Additional comments or notes about the site."),
			"tags":          nbschema.DSTagsAttribute(),
			"custom_fields": nbschema.DSCustomFieldsAttribute(),
		},
	}
}

func (d *SiteDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*netbox.APIClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *netbox.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	d.client = client
}

func (d *SiteDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data SiteDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var site *netbox.Site
	var err error
	var httpResp *http.Response

	// Determine if we're searching by ID, slug, or name
	if !data.ID.IsNull() {
		siteID := data.ID.ValueString()
		tflog.Debug(ctx, "Reading site by ID", map[string]interface{}{"id": siteID})

		var siteIDInt int32
		if _, parseErr := fmt.Sscanf(siteID, "%d", &siteIDInt); parseErr != nil {
			resp.Diagnostics.AddError("Invalid Site ID", fmt.Sprintf("Site ID must be a number, got: %s", siteID))
			return
		}
		site, httpResp, err = d.client.DcimAPI.DcimSitesRetrieve(ctx, siteIDInt).Execute()
		defer utils.CloseResponseBody(httpResp)
	} else if !data.Slug.IsNull() {
		siteSlug := data.Slug.ValueString()
		tflog.Debug(ctx, "Reading site by slug", map[string]interface{}{"slug": siteSlug})

		var sites *netbox.PaginatedSiteList
		sites, httpResp, err = d.client.DcimAPI.DcimSitesList(ctx).Slug([]string{siteSlug}).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError("Error reading site", utils.FormatAPIError("read site by slug", err, httpResp))
			return
		}
		if len(sites.GetResults()) == 0 {
			resp.Diagnostics.AddError("Site Not Found", fmt.Sprintf("No site found with slug: %s", siteSlug))
			return
		}
		if len(sites.GetResults()) > 1 {
			resp.Diagnostics.AddError("Multiple Sites Found", fmt.Sprintf("Multiple sites found with slug: %s. This should not happen as slugs should be unique.", siteSlug))
			return
		}
		site = &sites.GetResults()[0]
	} else if !data.Name.IsNull() {
		siteName := data.Name.ValueString()
		tflog.Debug(ctx, "Reading site by name", map[string]interface{}{"name": siteName})

		var sites *netbox.PaginatedSiteList
		sites, httpResp, err = d.client.DcimAPI.DcimSitesList(ctx).Name([]string{siteName}).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError("Error reading site", utils.FormatAPIError("read site by name", err, httpResp))
			return
		}
		if len(sites.GetResults()) == 0 {
			resp.Diagnostics.AddError("Site Not Found", fmt.Sprintf("No site found with name: %s", siteName))
			return
		}
		if len(sites.GetResults()) > 1 {
			resp.Diagnostics.AddError("Multiple Sites Found", fmt.Sprintf("Multiple sites found with name: %s. Site names may not be unique in Netbox.", siteName))
			return
		}
		site = &sites.GetResults()[0]
	} else {
		resp.Diagnostics.AddError("Missing Site Identifier", "Either 'id', 'slug', or 'name' must be specified to identify the site.")
		return
	}

	if err != nil {
		resp.Diagnostics.AddError("Error reading site", utils.FormatAPIError("read site", err, httpResp))
		return
	}
	if httpResp.StatusCode == 404 {
		resp.Diagnostics.AddError("Site Not Found", "The specified site was not found in Netbox.")
		return
	}
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError("Error reading site", fmt.Sprintf("Expected HTTP 200, got: %d", httpResp.StatusCode))
		return
	}

	// Map response to state using helper
	d.mapSiteToState(ctx, site, &data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// mapSiteToState maps API response to Terraform state for data sources.
func (d *SiteDataSource) mapSiteToState(ctx context.Context, site *netbox.Site, data *SiteDataSourceModel) {
	data.ID = types.StringValue(fmt.Sprintf("%d", site.GetId()))
	data.Name = types.StringValue(site.GetName())
	data.Slug = types.StringValue(site.GetSlug())

	// Handle status
	if site.HasStatus() {
		status := site.GetStatus()
		if status.HasValue() {
			statusValue, _ := status.GetValueOk()
			data.Status = types.StringValue(string(*statusValue))
		}
	} else {
		data.Status = types.StringNull()
	}

	// Handle region reference
	if site.HasRegion() {
		region := site.GetRegion()
		data.Region = types.StringValue(region.GetName())
		data.RegionID = types.StringValue(fmt.Sprintf("%d", region.GetId()))
	} else {
		data.Region = types.StringNull()
		data.RegionID = types.StringNull()
	}

	// Handle group reference
	if site.HasGroup() {
		group := site.GetGroup()
		data.Group = types.StringValue(group.GetName())
		data.GroupID = types.StringValue(fmt.Sprintf("%d", group.GetId()))
	} else {
		data.Group = types.StringNull()
		data.GroupID = types.StringNull()
	}

	// Handle tenant reference
	if site.HasTenant() {
		tenant := site.GetTenant()
		data.Tenant = types.StringValue(tenant.GetName())
		data.TenantID = types.StringValue(fmt.Sprintf("%d", tenant.GetId()))
	} else {
		data.Tenant = types.StringNull()
		data.TenantID = types.StringNull()
	}

	// Handle optional string fields
	data.Description = utils.StringFromAPIPreserveEmpty(site.HasDescription(), site.GetDescription, data.Description)
	data.Comments = utils.StringFromAPIPreserveEmpty(site.HasComments(), site.GetComments, data.Comments)
	data.Facility = utils.StringFromAPIPreserveEmpty(site.HasFacility(), site.GetFacility, data.Facility)

	// Handle tags
	if site.HasTags() {
		tags := utils.NestedTagsToTagModels(site.GetTags())
		tagsValue, tagDiags := types.SetValueFrom(ctx, utils.GetTagsAttributeType().ElemType, tags)
		if !tagDiags.HasError() {
			data.Tags = tagsValue
		}
	} else {
		data.Tags = types.SetNull(utils.GetTagsAttributeType().ElemType)
	}

	// Handle custom fields
	if site.HasCustomFields() {
		customFields := utils.MapToCustomFieldModels(site.GetCustomFields(), nil)
		customFieldsValue, cfDiags := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, customFields)
		if !cfDiags.HasError() {
			data.CustomFields = customFieldsValue
		}
	} else {
		data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)
	}
}
