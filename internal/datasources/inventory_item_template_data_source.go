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
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.

var (
	_ datasource.DataSource = &InventoryItemTemplateDataSource{}

	_ datasource.DataSourceWithConfigure = &InventoryItemTemplateDataSource{}
)

// NewInventoryItemTemplateDataSource returns a new data source implementing the inventory item template data source.

func NewInventoryItemTemplateDataSource() datasource.DataSource {

	return &InventoryItemTemplateDataSource{}

}

// InventoryItemTemplateDataSource defines the data source implementation.

type InventoryItemTemplateDataSource struct {
	client *netbox.APIClient
}

// InventoryItemTemplateDataSourceModel describes the data source data model.

type InventoryItemTemplateDataSourceModel struct {
	ID types.String `tfsdk:"id"`

	DeviceType types.String `tfsdk:"device_type"`

	DeviceTypeID types.String `tfsdk:"device_type_id"`

	Parent types.String `tfsdk:"parent"`

	ParentID types.String `tfsdk:"parent_id"`

	Name types.String `tfsdk:"name"`

	Label types.String `tfsdk:"label"`

	Role types.String `tfsdk:"role"`

	RoleID types.String `tfsdk:"role_id"`

	Manufacturer types.String `tfsdk:"manufacturer"`

	ManufacturerID types.String `tfsdk:"manufacturer_id"`

	PartID types.String `tfsdk:"part_id"`

	Description types.String `tfsdk:"description"`

	DisplayName types.String `tfsdk:"display_name"`

	ComponentType types.String `tfsdk:"component_type"`

	ComponentID types.String `tfsdk:"component_id"`
}

// Metadata returns the data source type name.

func (d *InventoryItemTemplateDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {

	resp.TypeName = req.ProviderTypeName + "_inventory_item_template"

}

// Schema defines the schema for the data source.

func (d *InventoryItemTemplateDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {

	resp.Schema = schema.Schema{

		MarkdownDescription: "Use this data source to get information about an inventory item template in NetBox.",

		Attributes: map[string]schema.Attribute{

			"id": schema.StringAttribute{

				MarkdownDescription: "The unique numeric ID of the inventory item template.",

				Required: true,
			},

			"device_type": schema.StringAttribute{

				MarkdownDescription: "The model name of the device type this inventory item template belongs to.",

				Computed: true,
			},

			"device_type_id": schema.StringAttribute{

				MarkdownDescription: "The ID of the device type this inventory item template belongs to.",

				Computed: true,
			},

			"parent": schema.StringAttribute{

				MarkdownDescription: "The name of the parent inventory item template.",

				Computed: true,
			},

			"parent_id": schema.StringAttribute{

				MarkdownDescription: "The ID of the parent inventory item template.",

				Computed: true,
			},

			"name": schema.StringAttribute{

				MarkdownDescription: "The name of the inventory item template.",

				Computed: true,
			},

			"label": schema.StringAttribute{

				MarkdownDescription: "Physical label of the inventory item template.",

				Computed: true,
			},

			"role": schema.StringAttribute{

				MarkdownDescription: "The name of the inventory item role.",

				Computed: true,
			},

			"role_id": schema.StringAttribute{

				MarkdownDescription: "The ID of the inventory item role.",

				Computed: true,
			},

			"manufacturer": schema.StringAttribute{

				MarkdownDescription: "The name of the manufacturer.",

				Computed: true,
			},

			"manufacturer_id": schema.StringAttribute{

				MarkdownDescription: "The ID of the manufacturer.",

				Computed: true,
			},

			"part_id": schema.StringAttribute{

				MarkdownDescription: "Manufacturer-assigned part identifier.",

				Computed: true,
			},

			"description": schema.StringAttribute{

				MarkdownDescription: "A description of the inventory item template.",

				Computed: true,
			},

			"display_name": schema.StringAttribute{

				MarkdownDescription: "The display name of the inventory item template.",

				Computed: true,
			},

			"component_type": schema.StringAttribute{

				MarkdownDescription: "The type of component this inventory item represents.",

				Computed: true,
			},

			"component_id": schema.StringAttribute{

				MarkdownDescription: "The ID of the component this inventory item represents.",

				Computed: true,
			},
		},
	}

}

// Configure adds the provider configured client to the data source.

func (d *InventoryItemTemplateDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {

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

// Read reads the data source.

func (d *InventoryItemTemplateDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var data InventoryItemTemplateDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	// Parse ID

	var id int32

	_, err := fmt.Sscanf(data.ID.ValueString(), "%d", &id)

	if err != nil {

		resp.Diagnostics.AddError(

			"Invalid ID",

			fmt.Sprintf("Could not parse ID '%s': %s", data.ID.ValueString(), err.Error()),
		)

		return

	}

	tflog.Debug(ctx, "Reading inventory item template", map[string]interface{}{"id": id})

	// Read from API

	result, httpResp, err := d.client.DcimAPI.DcimInventoryItemTemplatesRetrieve(ctx, id).Execute()

	defer utils.CloseResponseBody(httpResp)

	if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
		resp.Diagnostics.AddError(
			"Inventory Item Template Not Found",
			fmt.Sprintf("No inventory item template found with ID: %d", id),
		)
		return
	}

	if err != nil {

		resp.Diagnostics.AddError(

			"Error reading inventory item template",

			utils.FormatAPIError(fmt.Sprintf("read inventory item template ID %d", id), err, httpResp),
		)

		return

	}

	// Map response to state

	d.mapToState(ctx, result, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {

		return

	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

// mapToState maps the API response to the Terraform state.

func (d *InventoryItemTemplateDataSource) mapToState(ctx context.Context, result *netbox.InventoryItemTemplate, data *InventoryItemTemplateDataSourceModel, diags *diag.Diagnostics) {

	data.ID = types.StringValue(fmt.Sprintf("%d", result.GetId()))

	data.Name = types.StringValue(result.GetName())

	// Map device type (required field)

	deviceType := result.GetDeviceType()

	data.DeviceType = types.StringValue(deviceType.GetModel())

	data.DeviceTypeID = types.StringValue(fmt.Sprintf("%d", deviceType.GetId()))

	// Map parent (NullableInt32 - just the ID, not a nested object)

	// Check if the value is set and non-nil (parent: null in JSON means no parent)

	if parentID, ok := result.GetParentOk(); ok && parentID != nil {

		data.ParentID = types.StringValue(fmt.Sprintf("%d", *parentID))

		// Parent name requires a separate lookup - we only have the ID

		data.Parent = types.StringNull()

	} else {

		data.Parent = types.StringNull()

		data.ParentID = types.StringNull()

	}

	// Map label

	if result.HasLabel() && result.GetLabel() != "" {

		data.Label = types.StringValue(result.GetLabel())

	} else {

		data.Label = types.StringNull()

	}

	// Map role

	if result.HasRole() && result.GetRole().Id != 0 {

		role := result.GetRole()

		data.Role = types.StringValue(role.GetName())

		data.RoleID = types.StringValue(fmt.Sprintf("%d", role.GetId()))

	} else {

		data.Role = types.StringNull()

		data.RoleID = types.StringNull()

	}

	// Map manufacturer

	if result.HasManufacturer() && result.GetManufacturer().Id != 0 {

		manufacturer := result.GetManufacturer()

		data.Manufacturer = types.StringValue(manufacturer.GetName())

		data.ManufacturerID = types.StringValue(fmt.Sprintf("%d", manufacturer.GetId()))

	} else {

		data.Manufacturer = types.StringNull()

		data.ManufacturerID = types.StringNull()

	}

	// Map part ID

	if result.HasPartId() && result.GetPartId() != "" {

		data.PartID = types.StringValue(result.GetPartId())

	} else {

		data.PartID = types.StringNull()

	}

	// Map description

	if result.HasDescription() && result.GetDescription() != "" {

		data.Description = types.StringValue(result.GetDescription())

	} else {

		data.Description = types.StringNull()

	}

	// Map component type

	if result.HasComponentType() {

		componentTypePtr, ok := result.GetComponentTypeOk()

		if ok && componentTypePtr != nil && *componentTypePtr != "" {

			data.ComponentType = types.StringValue(*componentTypePtr)

		} else {

			data.ComponentType = types.StringNull()

		}

	} else {

		data.ComponentType = types.StringNull()

	}

	// Map component ID

	if result.HasComponentId() {

		componentIDPtr, ok := result.GetComponentIdOk()

		if ok && componentIDPtr != nil {

			data.ComponentID = types.StringValue(fmt.Sprintf("%d", *componentIDPtr))

		} else {

			data.ComponentID = types.StringNull()

		}

	} else {

		data.ComponentID = types.StringNull()

	}

	// Map display_name

	if result.GetDisplay() != "" {

		data.DisplayName = types.StringValue(result.GetDisplay())

	} else {

		data.DisplayName = types.StringNull()

	}

}
