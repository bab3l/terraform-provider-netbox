// Package datasources provides Terraform data source implementations for NetBox objects.

package datasources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bab3l/go-netbox"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ datasource.DataSource              = &InventoryItemRoleDataSource{}
	_ datasource.DataSourceWithConfigure = &InventoryItemRoleDataSource{}
)

// NewInventoryItemRoleDataSource returns a new data source implementing the inventory item role data source.
func NewInventoryItemRoleDataSource() datasource.DataSource {
	return &InventoryItemRoleDataSource{}
}

// InventoryItemRoleDataSource defines the data source implementation.
type InventoryItemRoleDataSource struct {
	client *netbox.APIClient
}

// InventoryItemRoleDataSourceModel describes the data source data model.
type InventoryItemRoleDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Slug        types.String `tfsdk:"slug"`
	Color       types.String `tfsdk:"color"`
	Description types.String `tfsdk:"description"`
	DisplayName types.String `tfsdk:"display_name"`
	Tags        types.Set    `tfsdk:"tags"`
}

// Metadata returns the data source type name.
func (d *InventoryItemRoleDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_inventory_item_role"
}

// Schema defines the schema for the data source.
func (d *InventoryItemRoleDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves information about an inventory item role in NetBox.",
		Attributes: map[string]schema.Attribute{
			// Filter attributes
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique numeric ID of the inventory item role. Use this to filter by ID.",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the inventory item role. Use this to filter by name.",
				Optional:            true,
				Computed:            true,
			},
			"slug": schema.StringAttribute{
				MarkdownDescription: "The slug of the inventory item role. Use this to filter by slug.",
				Optional:            true,
				Computed:            true,
			},
			// Computed attributes
			"color": schema.StringAttribute{
				MarkdownDescription: "The color associated with this role (6-character hex code).",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "A description of the inventory item role.",
				Computed:            true,
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "The display name of the inventory item role.",
				Computed:            true,
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "Tags associated with this inventory item role.",
				Computed:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *InventoryItemRoleDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *InventoryItemRoleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data InventoryItemRoleDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var role *netbox.InventoryItemRole

	// If ID is provided, look up directly
	if !data.ID.IsNull() && !data.ID.IsUnknown() {
		roleID, err := utils.ParseID(data.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Invalid Inventory Item Role ID",
				fmt.Sprintf("Inventory Item Role ID must be a number, got: %s", data.ID.ValueString()),
			)
			return
		}
		tflog.Debug(ctx, "Looking up inventory item role by ID", map[string]interface{}{
			"id": roleID,
		})
		response, httpResp, err := d.client.DcimAPI.DcimInventoryItemRolesRetrieve(ctx, roleID).Execute()
		defer utils.CloseResponseBody(httpResp)
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			resp.Diagnostics.AddError(
				"Inventory Item Role Not Found",
				fmt.Sprintf("No inventory item role found with ID: %d", roleID),
			)
			return
		}
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading inventory item role",
				utils.FormatAPIError(fmt.Sprintf("read inventory item role ID %d", roleID), err, httpResp),
			)
			return
		}
		role = response
	} else {
		// Search by filters
		tflog.Debug(ctx, "Searching for inventory item role", map[string]interface{}{
			"name": data.Name.ValueString(),
			"slug": data.Slug.ValueString(),
		})
		listReq := d.client.DcimAPI.DcimInventoryItemRolesList(ctx)
		if !data.Name.IsNull() && !data.Name.IsUnknown() {
			listReq = listReq.Name([]string{data.Name.ValueString()})
		}
		if !data.Slug.IsNull() && !data.Slug.IsUnknown() {
			listReq = listReq.Slug([]string{data.Slug.ValueString()})
		}
		response, httpResp, err := listReq.Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading inventory item roles",
				utils.FormatAPIError("list inventory item roles", err, httpResp),
			)
			return
		}
		if response.GetCount() == 0 {
			resp.Diagnostics.AddError(
				"No inventory item role found",
				"No inventory item role matching the specified criteria was found.",
			)
			return
		}
		if response.GetCount() > 1 {
			resp.Diagnostics.AddError(
				"Multiple inventory item roles found",
				fmt.Sprintf("Found %d inventory item roles matching the specified criteria. Please provide more specific filters.", response.GetCount()),
			)
			return
		}
		role = &response.GetResults()[0]
	}

	// Map response to model
	data.ID = types.StringValue(fmt.Sprintf("%d", role.GetId()))
	data.Name = types.StringValue(role.GetName())
	data.Slug = types.StringValue(role.GetSlug())

	// Map color
	if color, ok := role.GetColorOk(); ok && color != nil && *color != "" {
		data.Color = types.StringValue(*color)
	} else {
		data.Color = types.StringNull()
	}

	// Map description
	if desc, ok := role.GetDescriptionOk(); ok && desc != nil && *desc != "" {
		data.Description = types.StringValue(*desc)
	} else {
		data.Description = types.StringNull()
	}

	// Handle tags (simplified - just names)
	if role.HasTags() && len(role.GetTags()) > 0 {
		tagNames := make([]string, 0, len(role.GetTags()))
		for _, tag := range role.GetTags() {
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

	// Map display_name
	if role.GetDisplay() != "" {
		data.DisplayName = types.StringValue(role.GetDisplay())
	} else {
		data.DisplayName = types.StringNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
