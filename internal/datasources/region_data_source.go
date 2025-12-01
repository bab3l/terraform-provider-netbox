// Package datasources contains Terraform data source implementations for the Netbox provider.
package datasources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bab3l/go-netbox"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/bab3l/terraform-provider-netbox/internal/utils"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &RegionDataSource{}

func NewRegionDataSource() datasource.DataSource {
	return &RegionDataSource{}
}

// RegionDataSource defines the data source implementation.
type RegionDataSource struct {
	client *netbox.APIClient
}

// RegionDataSourceModel describes the data source data model.
type RegionDataSourceModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Slug         types.String `tfsdk:"slug"`
	Parent       types.String `tfsdk:"parent"`
	ParentID     types.String `tfsdk:"parent_id"`
	Description  types.String `tfsdk:"description"`
	Tags         types.Set    `tfsdk:"tags"`
	CustomFields types.Set    `tfsdk:"custom_fields"`
}

func (d *RegionDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_region"
}

func (d *RegionDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get information about a region in Netbox. Regions provide hierarchical geographic organization for sites. You can identify the region using `id`, `slug`, or `name`.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier for the region. Specify `id`, `slug`, or `name` to identify the region.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 50),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Full name of the region. Can be used to identify the region instead of `id` or `slug`.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 100),
				},
			},
			"slug": schema.StringAttribute{
				MarkdownDescription: "URL-friendly identifier for the region. Specify `id`, `slug`, or `name` to identify the region.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 100),
				},
			},
			"parent": schema.StringAttribute{
				MarkdownDescription: "ID of the parent region. Null if this is a top-level region.",
				Computed:            true,
			},
			"parent_id": schema.StringAttribute{
				MarkdownDescription: "ID of the parent region (same as parent, for compatibility).",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Detailed description of the region.",
				Computed:            true,
			},
			"tags": schema.SetNestedAttribute{
				MarkdownDescription: "Tags assigned to this region.",
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
				MarkdownDescription: "Custom fields assigned to this region.",
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

func (d *RegionDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *RegionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data RegionDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var region *netbox.Region
	var err error
	var httpResp *http.Response

	// Determine if we're searching by ID, slug, or name
	if !data.ID.IsNull() {
		regionID := data.ID.ValueString()
		tflog.Debug(ctx, "Reading region by ID", map[string]interface{}{
			"id": regionID,
		})

		var regionIDInt int32
		if _, parseErr := fmt.Sscanf(regionID, "%d", &regionIDInt); parseErr != nil {
			resp.Diagnostics.AddError(
				"Invalid Region ID",
				fmt.Sprintf("Region ID must be a number, got: %s", regionID),
			)
			return
		}

		region, httpResp, err = d.client.DcimAPI.DcimRegionsRetrieve(ctx, regionIDInt).Execute()
	} else if !data.Slug.IsNull() {
		regionSlug := data.Slug.ValueString()
		tflog.Debug(ctx, "Reading region by slug", map[string]interface{}{
			"slug": regionSlug,
		})

		var regions *netbox.PaginatedRegionList
		regions, httpResp, err = d.client.DcimAPI.DcimRegionsList(ctx).Slug([]string{regionSlug}).Execute()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading region",
				utils.FormatAPIError("read region by slug", err, httpResp),
			)
			return
		}
		if len(regions.GetResults()) == 0 {
			resp.Diagnostics.AddError(
				"Region Not Found",
				fmt.Sprintf("No region found with slug: %s", regionSlug),
			)
			return
		}
		if len(regions.GetResults()) > 1 {
			resp.Diagnostics.AddError(
				"Multiple Regions Found",
				fmt.Sprintf("Multiple regions found with slug: %s. This should not happen as slugs should be unique.", regionSlug),
			)
			return
		}
		region = &regions.GetResults()[0]
	} else if !data.Name.IsNull() {
		regionName := data.Name.ValueString()
		tflog.Debug(ctx, "Reading region by name", map[string]interface{}{
			"name": regionName,
		})

		var regions *netbox.PaginatedRegionList
		regions, httpResp, err = d.client.DcimAPI.DcimRegionsList(ctx).Name([]string{regionName}).Execute()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading region",
				utils.FormatAPIError("read region by name", err, httpResp),
			)
			return
		}
		if len(regions.GetResults()) == 0 {
			resp.Diagnostics.AddError(
				"Region Not Found",
				fmt.Sprintf("No region found with name: %s", regionName),
			)
			return
		}
		if len(regions.GetResults()) > 1 {
			resp.Diagnostics.AddError(
				"Multiple Regions Found",
				fmt.Sprintf("Multiple regions found with name: %s. Region names may not be unique in Netbox.", regionName),
			)
			return
		}
		region = &regions.GetResults()[0]
	} else {
		resp.Diagnostics.AddError(
			"Missing Region Identifier",
			"Either 'id', 'slug', or 'name' must be specified to identify the region.",
		)
		return
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading region",
			utils.FormatAPIError("read region", err, httpResp),
		)
		return
	}

	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Error reading region",
			fmt.Sprintf("Expected HTTP 200, got: %d", httpResp.StatusCode),
		)
		return
	}

	// Map response to state
	data.ID = types.StringValue(fmt.Sprintf("%d", region.GetId()))
	data.Name = types.StringValue(region.GetName())
	data.Slug = types.StringValue(region.GetSlug())

	if region.HasParent() && region.GetParent().Id != 0 {
		parent := region.GetParent()
		parentID := fmt.Sprintf("%d", parent.GetId())
		data.Parent = types.StringValue(parentID)
		data.ParentID = types.StringValue(parentID)
	} else {
		data.Parent = types.StringNull()
		data.ParentID = types.StringNull()
	}

	if region.HasDescription() {
		data.Description = types.StringValue(region.GetDescription())
	} else {
		data.Description = types.StringNull()
	}

	// Handle tags
	if region.HasTags() {
		tags := utils.NestedTagsToTagModels(region.GetTags())
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
	if region.HasCustomFields() {
		customFields := utils.MapToCustomFieldModels(region.GetCustomFields(), []utils.CustomFieldModel{})
		customFieldsValue, diags := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, customFields)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		data.CustomFields = customFieldsValue
	} else {
		data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
