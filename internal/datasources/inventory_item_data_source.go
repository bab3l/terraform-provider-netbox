// Package datasources provides Terraform data source implementations for NetBox objects.

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
var (
	_ datasource.DataSource              = &InventoryItemDataSource{}
	_ datasource.DataSourceWithConfigure = &InventoryItemDataSource{}
)

// NewInventoryItemDataSource returns a new data source implementing the inventory item data source.
func NewInventoryItemDataSource() datasource.DataSource {
	return &InventoryItemDataSource{}
}

// InventoryItemDataSource defines the data source implementation.
type InventoryItemDataSource struct {
	client *netbox.APIClient
}

// InventoryItemDataSourceModel describes the data source data model.
type InventoryItemDataSourceModel struct {
	ID               types.String `tfsdk:"id"`
	DeviceID         types.Int64  `tfsdk:"device_id"`
	DeviceName       types.String `tfsdk:"device_name"`
	Name             types.String `tfsdk:"name"`
	Label            types.String `tfsdk:"label"`
	ParentID         types.Int64  `tfsdk:"parent_id"`
	RoleID           types.Int64  `tfsdk:"role_id"`
	RoleName         types.String `tfsdk:"role_name"`
	ManufacturerID   types.Int64  `tfsdk:"manufacturer_id"`
	ManufacturerName types.String `tfsdk:"manufacturer_name"`
	PartID           types.String `tfsdk:"part_id"`
	Serial           types.String `tfsdk:"serial"`
	AssetTag         types.String `tfsdk:"asset_tag"`
	Discovered       types.Bool   `tfsdk:"discovered"`
	Description      types.String `tfsdk:"description"`
	DisplayName      types.String `tfsdk:"display_name"`
	Tags             types.Set    `tfsdk:"tags"`
	CustomFields     types.Set    `tfsdk:"custom_fields"`
}

// Metadata returns the data source type name.
func (d *InventoryItemDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_inventory_item"
}

// Schema defines the schema for the data source.
func (d *InventoryItemDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Retrieves information about an inventory item in NetBox. Inventory items represent hardware components installed within a device.

~> **Deprecation Warning:** Beginning in NetBox v4.3, inventory items are deprecated and planned for removal in a future release. Users are strongly encouraged to use modules and module types instead.`,

		Attributes: map[string]schema.Attribute{
			// Filter attributes
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique numeric ID of the inventory item. Use this to filter by ID.",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the inventory item. Use this to filter by name.",
				Optional:            true,
				Computed:            true,
			},
			"device_id": schema.Int64Attribute{
				MarkdownDescription: "The ID of the device. Use this to filter by device.",
				Optional:            true,
				Computed:            true,
			},

			// Computed attributes
			"device_name": schema.StringAttribute{
				MarkdownDescription: "The name of the device.",
				Computed:            true,
			},

			"label": schema.StringAttribute{
				MarkdownDescription: "Physical label on the inventory item.",
				Computed:            true,
			},

			"parent_id": schema.Int64Attribute{
				MarkdownDescription: "The ID of the parent inventory item.",
				Computed:            true,
			},

			"role_id": schema.Int64Attribute{
				MarkdownDescription: "The ID of the inventory item role.",
				Computed:            true,
			},

			"role_name": schema.StringAttribute{
				MarkdownDescription: "The name of the inventory item role.",
				Computed:            true,
			},

			"manufacturer_id": schema.Int64Attribute{
				MarkdownDescription: "The ID of the manufacturer.",
				Computed:            true,
			},

			"manufacturer_name": schema.StringAttribute{
				MarkdownDescription: "The name of the manufacturer.",
				Computed:            true,
			},

			"part_id": schema.StringAttribute{
				MarkdownDescription: "Manufacturer-assigned part identifier.",
				Computed:            true,
			},

			"serial": schema.StringAttribute{
				MarkdownDescription: "Serial number of the inventory item.",
				Computed:            true,
			},

			"asset_tag": schema.StringAttribute{
				MarkdownDescription: "A unique tag used to identify this inventory item.",
				Computed:            true,
			},

			"discovered": schema.BoolAttribute{
				MarkdownDescription: "Whether this item was automatically discovered.",
				Computed:            true,
			},

			"description": schema.StringAttribute{
				MarkdownDescription: "A description of the inventory item.",
				Computed:            true,
			},

			"display_name": schema.StringAttribute{
				MarkdownDescription: "The display name of the inventory item.",
				Computed:            true,
			},

			"tags": schema.SetAttribute{
				MarkdownDescription: "Tags associated with this inventory item.",
				Computed:            true,
				ElementType:         types.StringType,
			},

			"custom_fields": nbschema.DSCustomFieldsAttribute(),
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *InventoryItemDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *InventoryItemDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data InventoryItemDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var item *netbox.InventoryItem

	// If ID is provided, look up directly
	if !data.ID.IsNull() && !data.ID.IsUnknown() {
		itemID, err := utils.ParseID(data.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Invalid Inventory Item ID",
				fmt.Sprintf("Inventory Item ID must be a number, got: %s", data.ID.ValueString()),
			)
			return
		}
		tflog.Debug(ctx, "Looking up inventory item by ID", map[string]interface{}{
			"id": itemID,
		})
		response, httpResp, err := d.client.DcimAPI.DcimInventoryItemsRetrieve(ctx, itemID).Execute()
		defer utils.CloseResponseBody(httpResp)
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			resp.Diagnostics.AddError(
				"Inventory Item Not Found",
				fmt.Sprintf("No inventory item found with ID: %d", itemID),
			)
			return
		}
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading inventory item",
				utils.FormatAPIError(fmt.Sprintf("read inventory item ID %d", itemID), err, httpResp),
			)
			return
		}
		item = response
	} else {
		// Search by filters
		tflog.Debug(ctx, "Searching for inventory item", map[string]interface{}{
			"name":      data.Name.ValueString(),
			"device_id": data.DeviceID.ValueInt64(),
		})
		listReq := d.client.DcimAPI.DcimInventoryItemsList(ctx)
		if !data.Name.IsNull() && !data.Name.IsUnknown() {
			listReq = listReq.Name([]string{data.Name.ValueString()})
		}
		if !data.DeviceID.IsNull() && !data.DeviceID.IsUnknown() {
			deviceID32, err := utils.SafeInt32FromValue(data.DeviceID)
			if err != nil {
				resp.Diagnostics.AddError("Invalid Device ID", fmt.Sprintf("Device ID value overflow: %s", err))
				return
			}
			listReq = listReq.DeviceId([]int32{deviceID32})
		}
		response, httpResp, err := listReq.Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading inventory items",
				utils.FormatAPIError("list inventory items", err, httpResp),
			)
			return
		}
		result, ok := utils.ExpectSingleResult(
			response.GetResults(),
			"No inventory item found",
			"No inventory item matching the specified criteria was found.",
			"Multiple inventory items found",
			fmt.Sprintf("Found %d inventory items matching the specified criteria. Please provide more specific filters.", response.GetCount()),
			&resp.Diagnostics,
		)
		if !ok {
			return
		}
		item = result
	}

	// Map response to model
	data.ID = types.StringValue(fmt.Sprintf("%d", item.GetId()))
	data.Name = types.StringValue(item.GetName())

	// Map device
	device := item.GetDevice()
	data.DeviceID = types.Int64Value(int64(device.GetId()))
	data.DeviceName = types.StringValue(device.GetName())

	// Map label
	if label, ok := item.GetLabelOk(); ok && label != nil && *label != "" {
		data.Label = types.StringValue(*label)
	} else {
		data.Label = types.StringNull()
	}

	// Map parent
	if item.Parent.IsSet() && item.Parent.Get() != nil {
		data.ParentID = types.Int64Value(int64(*item.Parent.Get()))
	} else {
		data.ParentID = types.Int64Null()
	}

	// Map role
	if item.Role.IsSet() && item.Role.Get() != nil {
		role := item.Role.Get()
		data.RoleID = types.Int64Value(int64(role.GetId()))
		data.RoleName = types.StringValue(role.GetName())
	} else {
		data.RoleID = types.Int64Null()
		data.RoleName = types.StringNull()
	}

	// Map manufacturer
	if item.Manufacturer.IsSet() && item.Manufacturer.Get() != nil {
		mfr := item.Manufacturer.Get()
		data.ManufacturerID = types.Int64Value(int64(mfr.GetId()))
		data.ManufacturerName = types.StringValue(mfr.GetName())
	} else {
		data.ManufacturerID = types.Int64Null()
		data.ManufacturerName = types.StringNull()
	}

	// Map part_id
	if partID, ok := item.GetPartIdOk(); ok && partID != nil && *partID != "" {
		data.PartID = types.StringValue(*partID)
	} else {
		data.PartID = types.StringNull()
	}

	// Map serial
	if serial, ok := item.GetSerialOk(); ok && serial != nil && *serial != "" {
		data.Serial = types.StringValue(*serial)
	} else {
		data.Serial = types.StringNull()
	}

	// Map asset_tag
	if item.AssetTag.IsSet() && item.AssetTag.Get() != nil && *item.AssetTag.Get() != "" {
		data.AssetTag = types.StringValue(*item.AssetTag.Get())
	} else {
		data.AssetTag = types.StringNull()
	}

	// Map discovered
	if discovered, ok := item.GetDiscoveredOk(); ok && discovered != nil {
		data.Discovered = types.BoolValue(*discovered)
	} else {
		data.Discovered = types.BoolNull()
	}

	// Map description
	if desc, ok := item.GetDescriptionOk(); ok && desc != nil && *desc != "" {
		data.Description = types.StringValue(*desc)
	} else {
		data.Description = types.StringNull()
	}

	// Handle tags (simplified - just names)
	if item.HasTags() && len(item.GetTags()) > 0 {
		tagNames := make([]string, 0, len(item.GetTags()))
		for _, tag := range item.GetTags() {
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
	if item.GetDisplay() != "" {
		data.DisplayName = types.StringValue(item.GetDisplay())
	} else {
		data.DisplayName = types.StringNull()
	}

	// Handle custom fields - datasources return ALL fields
	if item.HasCustomFields() {
		customFields := utils.MapAllCustomFieldsToModels(item.GetCustomFields())
		customFieldsValue, cfDiags := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, customFields)
		if !cfDiags.HasError() {
			data.CustomFields = customFieldsValue
		}
	} else {
		data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
