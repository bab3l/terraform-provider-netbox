// Package datasources contains Terraform data source implementations for NetBox objects.
package datasources

import (
	"context"
	"fmt"

	"github.com/bab3l/go-netbox"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ datasource.DataSource              = &DeviceBayDataSource{}
	_ datasource.DataSourceWithConfigure = &DeviceBayDataSource{}
)

// NewDeviceBayDataSource returns a new data source implementing the DeviceBay data source.
func NewDeviceBayDataSource() datasource.DataSource {
	return &DeviceBayDataSource{}
}

// DeviceBayDataSource defines the data source implementation.
type DeviceBayDataSource struct {
	client *netbox.APIClient
}

// DeviceBayDataSourceModel describes the data source data model.
type DeviceBayDataSourceModel struct {
	ID              types.String `tfsdk:"id"`
	Device          types.String `tfsdk:"device"`
	Name            types.String `tfsdk:"name"`
	Label           types.String `tfsdk:"label"`
	Description     types.String `tfsdk:"description"`
	InstalledDevice types.String `tfsdk:"installed_device"`
	Tags            types.Set    `tfsdk:"tags"`
	CustomFields    types.Set    `tfsdk:"custom_fields"`
}

// Metadata returns the data source type name.
func (d *DeviceBayDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_device_bay"
}

// Schema defines the schema for the data source.
func (d *DeviceBayDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves information about a device bay in NetBox.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique numeric ID of the device bay. Use this to look up by ID.",
				Optional:            true,
				Computed:            true,
			},
			"device": schema.StringAttribute{
				MarkdownDescription: "ID of the parent device. Use with name for lookup.",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the device bay. Use with device for lookup.",
				Optional:            true,
				Computed:            true,
			},
			"label": schema.StringAttribute{
				MarkdownDescription: "Physical label for the device bay.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "A description of the device bay.",
				Computed:            true,
			},
			"installed_device": schema.StringAttribute{
				MarkdownDescription: "ID of the child device installed in this bay.",
				Computed:            true,
			},
			"tags":          nbschema.DSTagsAttribute(),
			"custom_fields": nbschema.DSCustomFieldsAttribute(),
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *DeviceBayDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *DeviceBayDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data DeviceBayDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var db *netbox.DeviceBay

	// Look up by ID if provided
	switch {
	case !data.ID.IsNull() && !data.ID.IsUnknown():
		dbID, err := utils.ParseID(data.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Invalid Device Bay ID",
				fmt.Sprintf("Device bay ID must be a number, got: %s", data.ID.ValueString()),
			)
			return
		}

		tflog.Debug(ctx, "Reading device bay by ID", map[string]interface{}{
			"id": dbID,
		})

		result, httpResp, err := d.client.DcimAPI.DcimDeviceBaysRetrieve(ctx, dbID).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading device bay",
				utils.FormatAPIError(fmt.Sprintf("read device bay ID %d", dbID), err, httpResp),
			)
			return
		}
		db = result
	case !data.Device.IsNull() && !data.Device.IsUnknown() && !data.Name.IsNull() && !data.Name.IsUnknown():
		// Look up by device and name
		tflog.Debug(ctx, "Reading device bay by device and name", map[string]interface{}{
			"device": data.Device.ValueString(),
			"name":   data.Name.ValueString(),
		})

		listReq := d.client.DcimAPI.DcimDeviceBaysList(ctx).Name([]string{data.Name.ValueString()})

		// Parse device ID
		deviceID, err := utils.ParseID(data.Device.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Invalid Device ID",
				fmt.Sprintf("Device ID must be a number, got: %s", data.Device.ValueString()),
			)
			return
		}
		listReq = listReq.DeviceId([]int32{deviceID})

		listResp, httpResp, err := listReq.Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading device bay",
				utils.FormatAPIError(fmt.Sprintf("read device bay by name %s", data.Name.ValueString()), err, httpResp),
			)
			return
		}

		if listResp.GetCount() == 0 {
			resp.Diagnostics.AddError(
				"Device bay not found",
				fmt.Sprintf("No device bay found with name: %s on device: %s", data.Name.ValueString(), data.Device.ValueString()),
			)
			return
		}

		db = &listResp.GetResults()[0]
	default:
		resp.Diagnostics.AddError(
			"Missing Required Attribute",
			"Either 'id' or both 'device' and 'name' must be specified to look up a device bay.",
		)
		return
	}

	// Map response to model
	d.mapResponseToModel(ctx, db, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// mapResponseToModel maps the API response to the Terraform model.
func (d *DeviceBayDataSource) mapResponseToModel(ctx context.Context, db *netbox.DeviceBay, data *DeviceBayDataSourceModel, diags *diag.Diagnostics) {
	data.ID = types.StringValue(fmt.Sprintf("%d", db.GetId()))
	data.Name = types.StringValue(db.GetName())

	// Map device
	data.Device = types.StringValue(fmt.Sprintf("%d", db.Device.GetId()))

	// Map label
	if label, ok := db.GetLabelOk(); ok && label != nil && *label != "" {
		data.Label = types.StringValue(*label)
	} else {
		data.Label = types.StringNull()
	}

	// Map description
	if desc, ok := db.GetDescriptionOk(); ok && desc != nil && *desc != "" {
		data.Description = types.StringValue(*desc)
	} else {
		data.Description = types.StringNull()
	}

	// Map installed_device
	if db.InstalledDevice.IsSet() && db.InstalledDevice.Get() != nil {
		data.InstalledDevice = types.StringValue(fmt.Sprintf("%d", db.InstalledDevice.Get().GetId()))
	} else {
		data.InstalledDevice = types.StringNull()
	}

	// Handle tags
	if db.HasTags() && len(db.GetTags()) > 0 {
		tags := utils.NestedTagsToTagModels(db.GetTags())
		tagsValue, tagDiags := types.SetValueFrom(ctx, utils.GetTagsAttributeType().ElemType, tags)
		diags.Append(tagDiags...)
		if diags.HasError() {
			return
		}
		data.Tags = tagsValue
	} else {
		data.Tags = types.SetNull(utils.GetTagsAttributeType().ElemType)
	}

	// Handle custom fields
	if db.HasCustomFields() {
		apiCustomFields := db.GetCustomFields()
		customFields := utils.MapToCustomFieldModels(apiCustomFields, nil)
		customFieldsValue, cfDiags := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, customFields)
		diags.Append(cfDiags...)
		if diags.HasError() {
			return
		}
		data.CustomFields = customFieldsValue
	} else {
		data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)
	}
}
