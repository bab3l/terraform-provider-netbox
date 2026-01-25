// Package datasources provides Terraform data source implementations for NetBox objects.

package datasources

import (
	"context"
	"fmt"

	"github.com/bab3l/go-netbox"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &ModuleDataSource{}

// NewModuleDataSource returns a new data source implementing the module data source.
func NewModuleDataSource() datasource.DataSource {
	return &ModuleDataSource{}
}

// ModuleDataSource defines the data source implementation.
type ModuleDataSource struct {
	client *netbox.APIClient
}

// ModuleDataSourceModel describes the data source data model.
type ModuleDataSourceModel struct {
	ID           types.Int32  `tfsdk:"id"`
	DisplayName  types.String `tfsdk:"display_name"`
	DeviceID     types.Int32  `tfsdk:"device_id"`
	Device       types.String `tfsdk:"device"`
	ModuleBayID  types.Int32  `tfsdk:"module_bay_id"`
	ModuleBay    types.String `tfsdk:"module_bay"`
	ModuleTypeID types.Int32  `tfsdk:"module_type_id"`
	ModuleType   types.String `tfsdk:"module_type"`
	Status       types.String `tfsdk:"status"`
	Serial       types.String `tfsdk:"serial"`
	AssetTag     types.String `tfsdk:"asset_tag"`
	Description  types.String `tfsdk:"description"`
	Comments     types.String `tfsdk:"comments"`
	CustomFields types.Set    `tfsdk:"custom_fields"`
}

// Metadata returns the data source type name.
func (d *ModuleDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_module"
}

// Schema defines the schema for the data source.
func (d *ModuleDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves information about a module in NetBox. Modules are hardware components installed in module bays within devices.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int32Attribute{
				MarkdownDescription: "The unique numeric ID of the module.",
				Optional:            true,
				Computed:            true,
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "The display name of the module.",
				Computed:            true,
			},
			"device_id": schema.Int32Attribute{
				MarkdownDescription: "The numeric ID of the device. Used with serial or module_bay_id for lookup when ID is not provided.",
				Optional:            true,
				Computed:            true,
			},
			"device": schema.StringAttribute{
				MarkdownDescription: "The name of the device.",
				Computed:            true,
			},
			"module_bay_id": schema.Int32Attribute{
				MarkdownDescription: "The numeric ID of the module bay. Used with device_id for lookup when ID is not provided.",
				Optional:            true,
				Computed:            true,
			},
			"module_bay": schema.StringAttribute{
				MarkdownDescription: "The name of the module bay.",
				Computed:            true,
			},
			"module_type_id": schema.Int32Attribute{
				MarkdownDescription: "The numeric ID of the module type.",
				Computed:            true,
			},
			"module_type": schema.StringAttribute{
				MarkdownDescription: "The model name of the module type.",
				Computed:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "Operational status.",
				Computed:            true,
			},
			"serial": schema.StringAttribute{
				MarkdownDescription: "Serial number of the module. Can be used for lookup with device_id.",
				Optional:            true,
				Computed:            true,
			},
			"asset_tag": schema.StringAttribute{
				MarkdownDescription: "A unique tag used to identify this module.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "A description of the module.",
				Computed:            true,
			},
			"comments": schema.StringAttribute{
				MarkdownDescription: "Additional comments or notes.",
				Computed:            true,
			},
			"custom_fields": schema.SetAttribute{
				MarkdownDescription: "Custom fields associated with this module.",
				Computed:            true,
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"name":  types.StringType,
						"type":  types.StringType,
						"value": types.StringType,
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *ModuleDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// Read retrieves the data source data.
func (d *ModuleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ModuleDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var module *netbox.Module
	switch {
	case !data.ID.IsNull() && !data.ID.IsUnknown():
		// Lookup by ID
		moduleID := data.ID.ValueInt32()
		tflog.Debug(ctx, "Reading module by ID", map[string]interface{}{
			"id": moduleID,
		})
		response, httpResp, err := d.client.DcimAPI.DcimModulesRetrieve(ctx, moduleID).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading module",
				utils.FormatAPIError(fmt.Sprintf("read module ID %d", moduleID), err, httpResp),
			)
			return
		}
		module = response

	case !data.DeviceID.IsNull() && !data.DeviceID.IsUnknown():
		// Lookup by device_id and optionally module_bay_id or serial
		deviceID := data.DeviceID.ValueInt32()
		tflog.Debug(ctx, "Reading module by device", map[string]interface{}{
			"device_id": deviceID,
		})
		listReq := d.client.DcimAPI.DcimModulesList(ctx).DeviceId([]int32{deviceID})
		if !data.ModuleBayID.IsNull() && !data.ModuleBayID.IsUnknown() {
			listReq = listReq.ModuleBayId([]string{fmt.Sprintf("%d", data.ModuleBayID.ValueInt32())})
		}
		if !data.Serial.IsNull() && !data.Serial.IsUnknown() {
			listReq = listReq.Serial([]string{data.Serial.ValueString()})
		}
		response, httpResp, err := listReq.Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading module",
				utils.FormatAPIError(fmt.Sprintf("read module by device ID %d", deviceID), err, httpResp),
			)
			return
		}
		count := int(response.GetCount())
		if count == 0 {
			resp.Diagnostics.AddError(
				"Module Not Found",
				fmt.Sprintf("No module found for device ID: %d", deviceID),
			)
			return
		}
		if count > 1 {
			resp.Diagnostics.AddError(
				"Multiple Modules Found",
				fmt.Sprintf("Found %d modules for device ID %d. Please provide module_bay_id, serial, or use ID to select a specific one.", count, deviceID),
			)
			return
		}
		module = &response.GetResults()[0]

	default:
		resp.Diagnostics.AddError(
			"Missing Required Attribute",
			"Either 'id' or 'device_id' must be specified to lookup a module.",
		)
		return
	}

	// Map response to model
	d.mapResponseToModel(module, &data)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// mapResponseToModel maps the API response to the Terraform model.
func (d *ModuleDataSource) mapResponseToModel(module *netbox.Module, data *ModuleDataSourceModel) {
	data.ID = types.Int32Value(module.GetId())

	// Display Name
	if module.GetDisplay() != "" {
		data.DisplayName = types.StringValue(module.GetDisplay())
	} else {
		data.DisplayName = types.StringNull()
	}

	// Map device
	if device := module.GetDevice(); device.Id != 0 {
		data.DeviceID = types.Int32Value(device.Id)
		data.Device = types.StringValue(device.GetName())
	}

	// Map module_bay
	moduleBay := module.GetModuleBay()
	data.ModuleBayID = types.Int32Value(moduleBay.Id)
	data.ModuleBay = types.StringValue(moduleBay.GetName())

	// Map module_type
	if mt := module.GetModuleType(); mt.Id != 0 {
		data.ModuleTypeID = types.Int32Value(mt.Id)
		data.ModuleType = types.StringValue(mt.GetModel())
	}

	// Map status
	if module.Status != nil {
		data.Status = types.StringValue(string(module.Status.GetValue()))
	} else {
		data.Status = types.StringNull()
	}

	// Map serial
	if serial, ok := module.GetSerialOk(); ok && serial != nil && *serial != "" {
		data.Serial = types.StringValue(*serial)
	} else {
		data.Serial = types.StringNull()
	}

	// Map asset_tag
	if module.AssetTag.IsSet() && module.AssetTag.Get() != nil && *module.AssetTag.Get() != "" {
		data.AssetTag = types.StringValue(*module.AssetTag.Get())
	} else {
		data.AssetTag = types.StringNull()
	}

	// Map description
	if desc, ok := module.GetDescriptionOk(); ok && desc != nil && *desc != "" {
		data.Description = types.StringValue(*desc)
	} else {
		data.Description = types.StringNull()
	}

	// Map comments
	if comments, ok := module.GetCommentsOk(); ok && comments != nil && *comments != "" {
		data.Comments = types.StringValue(*comments)
	} else {
		data.Comments = types.StringNull()
	}

	// Handle custom fields
	data.CustomFields = utils.CustomFieldsSetFromAPI(context.Background(), module.HasCustomFields(), module.GetCustomFields(), nil)
}
