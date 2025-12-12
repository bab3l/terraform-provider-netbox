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
var _ datasource.DataSource = &ModuleBayDataSource{}

// NewModuleBayDataSource returns a new data source implementing the module bay data source.
func NewModuleBayDataSource() datasource.DataSource {
	return &ModuleBayDataSource{}
}

// ModuleBayDataSource defines the data source implementation.
type ModuleBayDataSource struct {
	client *netbox.APIClient
}

// ModuleBayDataSourceModel describes the data source data model.
type ModuleBayDataSourceModel struct {
	ID              types.Int32  `tfsdk:"id"`
	DeviceID        types.Int32  `tfsdk:"device_id"`
	Device          types.String `tfsdk:"device"`
	Name            types.String `tfsdk:"name"`
	Label           types.String `tfsdk:"label"`
	Position        types.String `tfsdk:"position"`
	Description     types.String `tfsdk:"description"`
	InstalledModule types.Int32  `tfsdk:"installed_module"`
}

// Metadata returns the data source type name.
func (d *ModuleBayDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_module_bay"
}

// Schema defines the schema for the data source.
func (d *ModuleBayDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves information about a module bay in NetBox. Module bays are slots within devices that can accept modules.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int32Attribute{
				MarkdownDescription: "The unique numeric ID of the module bay.",
				Optional:            true,
				Computed:            true,
			},
			"device_id": schema.Int32Attribute{
				MarkdownDescription: "The numeric ID of the device. Used with name for lookup when ID is not provided.",
				Optional:            true,
				Computed:            true,
			},
			"device": schema.StringAttribute{
				MarkdownDescription: "The name of the device.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the module bay. Used with device_id for lookup when ID is not provided.",
				Optional:            true,
				Computed:            true,
			},
			"label": schema.StringAttribute{
				MarkdownDescription: "Physical label of the module bay.",
				Computed:            true,
			},
			"position": schema.StringAttribute{
				MarkdownDescription: "Identifier to reference when renaming installed components.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "A description of the module bay.",
				Computed:            true,
			},
			"installed_module": schema.Int32Attribute{
				MarkdownDescription: "The ID of the installed module, if any.",
				Computed:            true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *ModuleBayDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *ModuleBayDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ModuleBayDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var moduleBay *netbox.ModuleBay

	if !data.ID.IsNull() && !data.ID.IsUnknown() {
		// Lookup by ID
		bayID := data.ID.ValueInt32()

		tflog.Debug(ctx, "Reading module bay by ID", map[string]interface{}{
			"id": bayID,
		})

		response, httpResp, err := d.client.DcimAPI.DcimModuleBaysRetrieve(ctx, bayID).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading module bay",
				utils.FormatAPIError(fmt.Sprintf("read module bay ID %d", bayID), err, httpResp),
			)
			return
		}
		moduleBay = response
	} else if !data.DeviceID.IsNull() && !data.DeviceID.IsUnknown() && !data.Name.IsNull() && !data.Name.IsUnknown() {
		// Lookup by device_id and name
		deviceID := data.DeviceID.ValueInt32()
		name := data.Name.ValueString()

		tflog.Debug(ctx, "Reading module bay by device and name", map[string]interface{}{
			"device_id": deviceID,
			"name":      name,
		})

		response, httpResp, err := d.client.DcimAPI.DcimModuleBaysList(ctx).DeviceId([]int32{deviceID}).Name([]string{name}).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading module bay",
				utils.FormatAPIError(fmt.Sprintf("read module bay by device ID %d and name %s", deviceID, name), err, httpResp),
			)
			return
		}

		count := int(response.GetCount())
		if count == 0 {
			resp.Diagnostics.AddError(
				"Module Bay Not Found",
				fmt.Sprintf("No module bay found for device ID %d with name: %s", deviceID, name),
			)
			return
		}
		if count > 1 {
			resp.Diagnostics.AddError(
				"Multiple Module Bays Found",
				fmt.Sprintf("Found %d module bays for device ID %d with name %s. Please use ID to select a specific one.", count, deviceID, name),
			)
			return
		}

		moduleBay = &response.GetResults()[0]
	} else {
		resp.Diagnostics.AddError(
			"Missing Required Attribute",
			"Either 'id' or both 'device_id' and 'name' must be specified to lookup a module bay.",
		)
		return
	}

	// Map response to model
	d.mapResponseToModel(moduleBay, &data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// mapResponseToModel maps the API response to the Terraform model.
func (d *ModuleBayDataSource) mapResponseToModel(moduleBay *netbox.ModuleBay, data *ModuleBayDataSourceModel) {
	data.ID = types.Int32Value(moduleBay.GetId())
	data.Name = types.StringValue(moduleBay.GetName())

	// Map device
	if device := moduleBay.GetDevice(); device.Id != 0 {
		data.DeviceID = types.Int32Value(device.Id)
		data.Device = types.StringValue(device.GetName())
	}

	// Map label
	if label, ok := moduleBay.GetLabelOk(); ok && label != nil && *label != "" {
		data.Label = types.StringValue(*label)
	} else {
		data.Label = types.StringNull()
	}

	// Map position
	if pos, ok := moduleBay.GetPositionOk(); ok && pos != nil && *pos != "" {
		data.Position = types.StringValue(*pos)
	} else {
		data.Position = types.StringNull()
	}

	// Map description
	if desc, ok := moduleBay.GetDescriptionOk(); ok && desc != nil && *desc != "" {
		data.Description = types.StringValue(*desc)
	} else {
		data.Description = types.StringNull()
	}

	// Map installed_module
	if moduleBay.InstalledModule.IsSet() && moduleBay.InstalledModule.Get() != nil {
		data.InstalledModule = types.Int32Value(moduleBay.InstalledModule.Get().Id)
	} else {
		data.InstalledModule = types.Int32Null()
	}
}
