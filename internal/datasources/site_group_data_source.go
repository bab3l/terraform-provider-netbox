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

	"github.com/bab3l/terraform-provider-netbox/internal/utils"
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
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier for the site group. Specify `id`, `slug`, or `name` to identify the site group.",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Full name of the site group. Can be used to identify the site group instead of `id` or `slug`.",
				Optional:            true,
				Computed:            true,
			},
			"slug": schema.StringAttribute{
				MarkdownDescription: "URL-friendly identifier for the site group. Specify `id`, `slug`, or `name` to identify the site group.",
				Optional:            true,
				Computed:            true,
			},
			"parent": schema.StringAttribute{
				MarkdownDescription: "Name or ID of the parent site group. Null if this is a top-level group.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Detailed description of the site group.",
				Computed:            true,
			},
			"tags": schema.SetNestedAttribute{
				MarkdownDescription: "Tags assigned to this site group.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "Name of the tag.",
							Computed:            true,
						},
						"slug": schema.StringAttribute{
							MarkdownDescription: "Slug of the tag.",
							Computed:            true,
						},
					},
				},
			},
			"custom_fields": schema.SetNestedAttribute{
				MarkdownDescription: "Custom fields assigned to this site group.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "Name of the custom field.",
							Computed:            true,
						},
						"type": schema.StringAttribute{
							MarkdownDescription: "Type of the custom field.",
							Computed:            true,
						},
						"value": schema.StringAttribute{
							MarkdownDescription: "Value of the custom field.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *SiteGroupDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *SiteGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data SiteGroupDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var siteGroup *netbox.SiteGroup
	var err error
	var httpResp *http.Response

	// Determine if we're searching by ID, slug, or name
	if !data.ID.IsNull() {
		// Search by ID
		siteGroupID := data.ID.ValueString()
		tflog.Debug(ctx, "Reading site group by ID", map[string]interface{}{
			"id": siteGroupID,
		})

		// Parse the site group ID to int32 for the API call
		var siteGroupIDInt int32
		if _, parseErr := fmt.Sscanf(siteGroupID, "%d", &siteGroupIDInt); parseErr != nil {
			resp.Diagnostics.AddError(
				"Invalid Site Group ID",
				fmt.Sprintf("Site Group ID must be a number, got: %s", siteGroupID),
			)
			return
		}

		// Retrieve the site group via API
		siteGroup, httpResp, err = d.client.DcimAPI.DcimSiteGroupsRetrieve(ctx, siteGroupIDInt).Execute()
	} else if !data.Slug.IsNull() {
		// Search by slug
		siteGroupSlug := data.Slug.ValueString()
		tflog.Debug(ctx, "Reading site group by slug", map[string]interface{}{
			"slug": siteGroupSlug,
		})

		// List site groups with slug filter
		var siteGroups *netbox.PaginatedSiteGroupList
		siteGroups, httpResp, err = d.client.DcimAPI.DcimSiteGroupsList(ctx).Slug([]string{siteGroupSlug}).Execute()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading site group",
				utils.FormatAPIError("read site group by slug", err, httpResp),
			)
			return
		}
		if len(siteGroups.GetResults()) == 0 {
			resp.Diagnostics.AddError(
				"Site Group Not Found",
				fmt.Sprintf("No site group found with slug: %s", siteGroupSlug),
			)
			return
		}
		if len(siteGroups.GetResults()) > 1 {
			resp.Diagnostics.AddError(
				"Multiple Site Groups Found",
				fmt.Sprintf("Multiple site groups found with slug: %s. This should not happen as slugs should be unique.", siteGroupSlug),
			)
			return
		}
		siteGroup = &siteGroups.GetResults()[0]
	} else if !data.Name.IsNull() {
		// Search by name
		siteGroupName := data.Name.ValueString()
		tflog.Debug(ctx, "Reading site group by name", map[string]interface{}{
			"name": siteGroupName,
		})

		// List site groups with name filter
		var siteGroups *netbox.PaginatedSiteGroupList
		siteGroups, httpResp, err = d.client.DcimAPI.DcimSiteGroupsList(ctx).Name([]string{siteGroupName}).Execute()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading site group",
				utils.FormatAPIError("read site group by name", err, httpResp),
			)
			return
		}
		if len(siteGroups.GetResults()) == 0 {
			resp.Diagnostics.AddError(
				"Site Group Not Found",
				fmt.Sprintf("No site group found with name: %s", siteGroupName),
			)
			return
		}
		if len(siteGroups.GetResults()) > 1 {
			resp.Diagnostics.AddError(
				"Multiple Site Groups Found",
				fmt.Sprintf("Multiple site groups found with name: %s. Site group names may not be unique in Netbox.", siteGroupName),
			)
			return
		}
		siteGroup = &siteGroups.GetResults()[0]
	} else {
		resp.Diagnostics.AddError(
			"Missing Site Group Identifier",
			"Either 'id', 'slug', or 'name' must be specified to identify the site group.",
		)
		return
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading site group",
			utils.FormatAPIError("read site group", err, httpResp),
		)
		return
	}

	if httpResp.StatusCode == 404 {
		resp.Diagnostics.AddError(
			"Site Group Not Found",
			"The specified site group was not found in Netbox.",
		)
		return
	}

	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Error reading site group",
			fmt.Sprintf("Expected HTTP 200, got: %d", httpResp.StatusCode),
		)
		return
	}

	// Update the model with the response from the API
	data.ID = types.StringValue(fmt.Sprintf("%d", siteGroup.GetId()))
	data.Name = types.StringValue(siteGroup.GetName())
	data.Slug = types.StringValue(siteGroup.GetSlug())

	// Handle parent site group
	if siteGroup.HasParent() {
		parent := siteGroup.GetParent()
		data.Parent = types.StringValue(parent.GetName())
	} else {
		data.Parent = types.StringNull()
	}

	if siteGroup.HasDescription() {
		data.Description = types.StringValue(siteGroup.GetDescription())
	} else {
		data.Description = types.StringNull()
	}

	// Handle tags
	if siteGroup.HasTags() {
		tags := utils.NestedTagsToTagModels(siteGroup.GetTags())
		tagsValue, diags := types.SetValueFrom(ctx, utils.GetTagsAttributeType().ElemType, tags)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		data.Tags = tagsValue
	} else {
		data.Tags = types.SetNull(utils.GetTagsAttributeType().ElemType)
	}

	// Handle custom fields
	if siteGroup.HasCustomFields() {
		// For data sources, we extract all available custom fields
		customFields := utils.MapToCustomFieldModels(siteGroup.GetCustomFields(), nil)
		customFieldsValue, diags := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, customFields)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		data.CustomFields = customFieldsValue
	} else {
		data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
