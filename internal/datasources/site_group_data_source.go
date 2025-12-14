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
	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &SiteGroupDataSource{}

func NewSiteGroupDataSource() datasource.DataSource {
	return &SiteGroupDataSource{}
}

// SiteGroupDataSource defines the data source implementation.
type SiteGroupDataSource struct {
	client *netbox.APIClient
}

// SiteGroupDataSourceModel describes the data source data model.
type SiteGroupDataSourceModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Slug         types.String `tfsdk:"slug"`
	Parent       types.String `tfsdk:"parent"`
	Description  types.String `tfsdk:"description"`
	Tags         types.Set    `tfsdk:"tags"`
	CustomFields types.Set    `tfsdk:"custom_fields"`
}

func (d *SiteGroupDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_site_group"
}

func (d *SiteGroupDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get information about a site group in Netbox. Site groups provide hierarchical organization of sites with support for parent-child relationships. You can identify the site group using `id`, `slug`, or `name`.",

		Attributes: map[string]schema.Attribute{
			"id":            nbschema.DSIDAttribute("site group"),
			"name":          nbschema.DSNameAttribute("site group"),
			"slug":          nbschema.DSSlugAttribute("site group"),
			"parent":        nbschema.DSComputedStringAttribute("ID of the parent site group. Null if this is a top-level group."),
			"description":   nbschema.DSComputedStringAttribute("Detailed description of the site group."),
			"tags":          nbschema.DSTagsAttribute(),
			"custom_fields": nbschema.DSCustomFieldsAttribute(),
		},
	}
}

func (d *SiteGroupDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *SiteGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data SiteGroupDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var siteGroup *netbox.SiteGroup
	var err error
	var httpResp *http.Response

	// Determine if we're searching by ID, slug, or name
	switch {
	case !data.ID.IsNull():
		siteGroupID := data.ID.ValueString()
		tflog.Debug(ctx, "Reading site group by ID", map[string]interface{}{"id": siteGroupID})

		var siteGroupIDInt int32
		if _, parseErr := fmt.Sscanf(siteGroupID, "%d", &siteGroupIDInt); parseErr != nil {
			resp.Diagnostics.AddError("Invalid Site Group ID", fmt.Sprintf("Site Group ID must be a number, got: %s", siteGroupID))
			return
		}
		siteGroup, httpResp, err = d.client.DcimAPI.DcimSiteGroupsRetrieve(ctx, siteGroupIDInt).Execute()
		defer utils.CloseResponseBody(httpResp)
	case !data.Slug.IsNull():
		siteGroupSlug := data.Slug.ValueString()
		tflog.Debug(ctx, "Reading site group by slug", map[string]interface{}{"slug": siteGroupSlug})

		var siteGroups *netbox.PaginatedSiteGroupList
		siteGroups, httpResp, err = d.client.DcimAPI.DcimSiteGroupsList(ctx).Slug([]string{siteGroupSlug}).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError("Error reading site group", utils.FormatAPIError("read site group by slug", err, httpResp))
			return
		}
		if len(siteGroups.GetResults()) == 0 {
			resp.Diagnostics.AddError("Site Group Not Found", fmt.Sprintf("No site group found with slug: %s", siteGroupSlug))
			return
		}
		if len(siteGroups.GetResults()) > 1 {
			resp.Diagnostics.AddError("Multiple Site Groups Found", fmt.Sprintf("Multiple site groups found with slug: %s. This should not happen as slugs should be unique.", siteGroupSlug))
			return
		}
		siteGroup = &siteGroups.GetResults()[0]
	case !data.Name.IsNull():
		siteGroupName := data.Name.ValueString()
		tflog.Debug(ctx, "Reading site group by name", map[string]interface{}{"name": siteGroupName})

		var siteGroups *netbox.PaginatedSiteGroupList
		siteGroups, httpResp, err = d.client.DcimAPI.DcimSiteGroupsList(ctx).Name([]string{siteGroupName}).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError("Error reading site group", utils.FormatAPIError("read site group by name", err, httpResp))
			return
		}
		if len(siteGroups.GetResults()) == 0 {
			resp.Diagnostics.AddError("Site Group Not Found", fmt.Sprintf("No site group found with name: %s", siteGroupName))
			return
		}
		if len(siteGroups.GetResults()) > 1 {
			resp.Diagnostics.AddError("Multiple Site Groups Found", fmt.Sprintf("Multiple site groups found with name: %s. Site group names may not be unique in Netbox.", siteGroupName))
			return
		}
		siteGroup = &siteGroups.GetResults()[0]
	default:
		resp.Diagnostics.AddError("Missing Site Group Identifier", "Either 'id', 'slug', or 'name' must be specified to identify the site group.")
		return
	}

	if err != nil {
		resp.Diagnostics.AddError("Error reading site group", utils.FormatAPIError("read site group", err, httpResp))
		return
	}
	if httpResp.StatusCode == 404 {
		resp.Diagnostics.AddError("Site Group Not Found", "The specified site group was not found in Netbox.")
		return
	}
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError("Error reading site group", fmt.Sprintf("Expected HTTP 200, got: %d", httpResp.StatusCode))
		return
	}

	// Map response to state using helper
	d.mapSiteGroupToState(ctx, siteGroup, &data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// mapSiteGroupToState maps API response to Terraform state for data sources.
func (d *SiteGroupDataSource) mapSiteGroupToState(ctx context.Context, siteGroup *netbox.SiteGroup, data *SiteGroupDataSourceModel) {
	data.ID = types.StringValue(fmt.Sprintf("%d", siteGroup.GetId()))
	data.Name = types.StringValue(siteGroup.GetName())
	data.Slug = types.StringValue(siteGroup.GetSlug())

	// Handle parent site group
	if siteGroup.HasParent() {
		parent := siteGroup.GetParent()
		data.Parent = types.StringValue(fmt.Sprintf("%d", parent.GetId()))
	} else {
		data.Parent = types.StringNull()
	}

	// Handle description
	data.Description = utils.StringFromAPIPreserveEmpty(siteGroup.HasDescription(), siteGroup.GetDescription, data.Description)

	// Handle tags
	if siteGroup.HasTags() {
		tags := utils.NestedTagsToTagModels(siteGroup.GetTags())
		tagsValue, tagDiags := types.SetValueFrom(ctx, utils.GetTagsAttributeType().ElemType, tags)
		if !tagDiags.HasError() {
			data.Tags = tagsValue
		}
	} else {
		data.Tags = types.SetNull(utils.GetTagsAttributeType().ElemType)
	}

	// Handle custom fields
	if siteGroup.HasCustomFields() {
		customFields := utils.MapToCustomFieldModels(siteGroup.GetCustomFields(), nil)
		customFieldsValue, cfDiags := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, customFields)
		if !cfDiags.HasError() {
			data.CustomFields = customFieldsValue
		}
	} else {
		data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)
	}
}
