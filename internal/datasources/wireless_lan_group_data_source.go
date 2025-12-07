// Package datasources provides Terraform data source implementations for NetBox objects.
package datasources

import (
	"context"
	"fmt"

	"github.com/bab3l/go-netbox"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/bab3l/terraform-provider-netbox/internal/utils"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ datasource.DataSource              = &WirelessLANGroupDataSource{}
	_ datasource.DataSourceWithConfigure = &WirelessLANGroupDataSource{}
)

// NewWirelessLANGroupDataSource returns a new data source implementing the wireless LAN group data source.
func NewWirelessLANGroupDataSource() datasource.DataSource {
	return &WirelessLANGroupDataSource{}
}

// WirelessLANGroupDataSource defines the data source implementation.
type WirelessLANGroupDataSource struct {
	client *netbox.APIClient
}

// WirelessLANGroupDataSourceModel describes the data source data model.
type WirelessLANGroupDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Slug        types.String `tfsdk:"slug"`
	Description types.String `tfsdk:"description"`
	ParentID    types.Int64  `tfsdk:"parent_id"`
	ParentName  types.String `tfsdk:"parent_name"`
	Tags        types.Set    `tfsdk:"tags"`
}

// Metadata returns the data source type name.
func (d *WirelessLANGroupDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_wireless_lan_group"
}

// Schema defines the schema for the data source.
func (d *WirelessLANGroupDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves information about a wireless LAN group in NetBox.",

		Attributes: map[string]schema.Attribute{
			// Filter attributes
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique numeric ID of the wireless LAN group. Use this to filter by ID.",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the wireless LAN group. Use this to filter by name.",
				Optional:            true,
				Computed:            true,
			},
			"slug": schema.StringAttribute{
				MarkdownDescription: "The slug of the wireless LAN group. Use this to filter by slug.",
				Optional:            true,
				Computed:            true,
			},

			// Computed attributes
			"description": schema.StringAttribute{
				MarkdownDescription: "A description of the wireless LAN group.",
				Computed:            true,
			},
			"parent_id": schema.Int64Attribute{
				MarkdownDescription: "The ID of the parent wireless LAN group.",
				Computed:            true,
			},
			"parent_name": schema.StringAttribute{
				MarkdownDescription: "The name of the parent wireless LAN group.",
				Computed:            true,
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "Tags associated with this wireless LAN group.",
				Computed:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *WirelessLANGroupDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// Read refreshes the data source data.
func (d *WirelessLANGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data WirelessLANGroupDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var group *netbox.WirelessLANGroup

	// If ID is provided, look up directly
	if !data.ID.IsNull() && !data.ID.IsUnknown() {
		groupID, err := utils.ParseID(data.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Invalid Wireless LAN Group ID",
				fmt.Sprintf("Wireless LAN Group ID must be a number, got: %s", data.ID.ValueString()),
			)
			return
		}

		tflog.Debug(ctx, "Looking up wireless LAN group by ID", map[string]interface{}{
			"id": groupID,
		})

		response, httpResp, err := d.client.WirelessAPI.WirelessWirelessLanGroupsRetrieve(ctx, groupID).Execute()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading wireless LAN group",
				utils.FormatAPIError(fmt.Sprintf("read wireless LAN group ID %d", groupID), err, httpResp),
			)
			return
		}
		group = response
	} else {
		// Search by filters
		tflog.Debug(ctx, "Searching for wireless LAN group", map[string]interface{}{
			"name": data.Name.ValueString(),
			"slug": data.Slug.ValueString(),
		})

		listReq := d.client.WirelessAPI.WirelessWirelessLanGroupsList(ctx)

		if !data.Name.IsNull() && !data.Name.IsUnknown() {
			listReq = listReq.Name([]string{data.Name.ValueString()})
		}

		if !data.Slug.IsNull() && !data.Slug.IsUnknown() {
			listReq = listReq.Slug([]string{data.Slug.ValueString()})
		}

		response, httpResp, err := listReq.Execute()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading wireless LAN groups",
				utils.FormatAPIError("list wireless LAN groups", err, httpResp),
			)
			return
		}

		if response.GetCount() == 0 {
			resp.Diagnostics.AddError(
				"No wireless LAN group found",
				"No wireless LAN group matching the specified criteria was found.",
			)
			return
		}

		if response.GetCount() > 1 {
			resp.Diagnostics.AddError(
				"Multiple wireless LAN groups found",
				fmt.Sprintf("Found %d wireless LAN groups matching the specified criteria. Please provide more specific filters.", response.GetCount()),
			)
			return
		}

		group = &response.GetResults()[0]
	}

	// Map response to model
	data.ID = types.StringValue(fmt.Sprintf("%d", group.GetId()))
	data.Name = types.StringValue(group.GetName())
	data.Slug = types.StringValue(group.GetSlug())

	// Map description
	if desc, ok := group.GetDescriptionOk(); ok && desc != nil && *desc != "" {
		data.Description = types.StringValue(*desc)
	} else {
		data.Description = types.StringNull()
	}

	// Map parent
	if group.Parent.IsSet() && group.Parent.Get() != nil {
		parent := group.Parent.Get()
		data.ParentID = types.Int64Value(int64(parent.GetId()))
		data.ParentName = types.StringValue(parent.GetName())
	} else {
		data.ParentID = types.Int64Null()
		data.ParentName = types.StringNull()
	}

	// Handle tags (simplified - just names)
	if group.HasTags() && len(group.GetTags()) > 0 {
		tagNames := make([]string, 0, len(group.GetTags()))
		for _, tag := range group.GetTags() {
			tagNames = append(tagNames, tag.GetName())
		}
		tagsValue, diags := types.SetValueFrom(ctx, types.StringType, tagNames)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		data.Tags = tagsValue
	} else {
		data.Tags = types.SetNull(types.StringType)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
