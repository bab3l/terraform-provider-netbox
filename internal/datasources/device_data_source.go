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
var _ datasource.DataSource = &DeviceDataSource{}

func NewDeviceDataSource() datasource.DataSource {
	return &DeviceDataSource{}
}

// DeviceDataSource defines the data source implementation.
type DeviceDataSource struct {
	client *netbox.APIClient
}

// DeviceDataSourceModel describes the data source data model.
type DeviceDataSourceModel struct {
	ID           types.String  `tfsdk:"id"`
	Name         types.String  `tfsdk:"name"`
	DeviceType   types.String  `tfsdk:"device_type"`
	Role         types.String  `tfsdk:"role"`
	Tenant       types.String  `tfsdk:"tenant"`
	Platform     types.String  `tfsdk:"platform"`
	Serial       types.String  `tfsdk:"serial"`
	AssetTag     types.String  `tfsdk:"asset_tag"`
	Site         types.String  `tfsdk:"site"`
	Location     types.String  `tfsdk:"location"`
	Rack         types.String  `tfsdk:"rack"`
	Position     types.Float64 `tfsdk:"position"`
	Face         types.String  `tfsdk:"face"`
	Latitude     types.Float64 `tfsdk:"latitude"`
	Longitude    types.Float64 `tfsdk:"longitude"`
	Status       types.String  `tfsdk:"status"`
	Airflow      types.String  `tfsdk:"airflow"`
	VcPosition   types.Int64   `tfsdk:"vc_position"`
	VcPriority   types.Int64   `tfsdk:"vc_priority"`
	Description  types.String  `tfsdk:"description"`
	Comments     types.String  `tfsdk:"comments"`
	Tags         types.Set     `tfsdk:"tags"`
	CustomFields types.Set     `tfsdk:"custom_fields"`
}

func (d *DeviceDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_device"
}

func (d *DeviceDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get information about a device in Netbox. Devices represent physical or virtual hardware. You can identify the device using `id`, `name`, or `serial`.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier for the device. Specify `id`, `name`, or `serial` to identify the device.",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the device. Can be used to identify the device instead of `id` or `serial`.",
				Optional:            true,
				Computed:            true,
			},
			"device_type": schema.StringAttribute{
				MarkdownDescription: "The device type. Returns the device type slug.",
				Computed:            true,
			},
			"role": schema.StringAttribute{
				MarkdownDescription: "The device role. Returns the device role slug.",
				Computed:            true,
			},
			"tenant": schema.StringAttribute{
				MarkdownDescription: "The tenant that owns this device. Returns the tenant slug.",
				Computed:            true,
			},
			"platform": schema.StringAttribute{
				MarkdownDescription: "The platform running on this device. Returns the platform slug.",
				Computed:            true,
			},
			"serial": schema.StringAttribute{
				MarkdownDescription: "Chassis serial number. Can be used to identify the device instead of `id` or `name`.",
				Optional:            true,
				Computed:            true,
			},
			"asset_tag": schema.StringAttribute{
				MarkdownDescription: "A unique tag used to identify this device.",
				Computed:            true,
			},
			"site": schema.StringAttribute{
				MarkdownDescription: "The site where this device is located. Returns the site slug.",
				Computed:            true,
			},
			"location": schema.StringAttribute{
				MarkdownDescription: "The location within the site. Returns the location slug.",
				Computed:            true,
			},
			"rack": schema.StringAttribute{
				MarkdownDescription: "The rack where this device is mounted. Returns the rack name.",
				Computed:            true,
			},
			"position": schema.Float64Attribute{
				MarkdownDescription: "Position in the rack (in rack units from the bottom).",
				Computed:            true,
			},
			"face": schema.StringAttribute{
				MarkdownDescription: "Which face of the rack the device is mounted on (front or rear).",
				Computed:            true,
			},
			"latitude": schema.Float64Attribute{
				MarkdownDescription: "GPS latitude coordinate in decimal format.",
				Computed:            true,
			},
			"longitude": schema.Float64Attribute{
				MarkdownDescription: "GPS longitude coordinate in decimal format.",
				Computed:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "Operational status of the device.",
				Computed:            true,
			},
			"airflow": schema.StringAttribute{
				MarkdownDescription: "Direction of airflow through the device.",
				Computed:            true,
			},
			"vc_position": schema.Int64Attribute{
				MarkdownDescription: "Position within a virtual chassis.",
				Computed:            true,
			},
			"vc_priority": schema.Int64Attribute{
				MarkdownDescription: "Virtual chassis master election priority.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Brief description of the device.",
				Computed:            true,
			},
			"comments": schema.StringAttribute{
				MarkdownDescription: "Comments about the device (supports Markdown).",
				Computed:            true,
			},
			"tags": schema.SetNestedAttribute{
				MarkdownDescription: "Tags assigned to this device.",
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
				MarkdownDescription: "Custom fields assigned to this device.",
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

func (d *DeviceDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *DeviceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data DeviceDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var device *netbox.DeviceWithConfigContext
	var httpResp *http.Response
	var err error

	// Look up device by ID, name, or serial
	if !data.ID.IsNull() && data.ID.ValueString() != "" {
		// Look up by ID
		var id int32
		if _, parseErr := fmt.Sscanf(data.ID.ValueString(), "%d", &id); parseErr != nil {
			resp.Diagnostics.AddError(
				"Invalid Device ID",
				fmt.Sprintf("Device ID must be a number, got: %s", data.ID.ValueString()),
			)
			return
		}

		tflog.Debug(ctx, "Looking up device by ID", map[string]interface{}{
			"id": id,
		})

		device, httpResp, err = d.client.DcimAPI.DcimDevicesRetrieve(ctx, id).Execute()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading device",
				utils.FormatAPIError(fmt.Sprintf("read device ID %d", id), err, httpResp),
			)
			return
		}
	} else if !data.Name.IsNull() && data.Name.ValueString() != "" {
		// Look up by name
		tflog.Debug(ctx, "Looking up device by name", map[string]interface{}{
			"name": data.Name.ValueString(),
		})

		list, listResp, listErr := d.client.DcimAPI.DcimDevicesList(ctx).Name([]string{data.Name.ValueString()}).Execute()
		if listErr != nil {
			resp.Diagnostics.AddError(
				"Error reading device",
				utils.FormatAPIError(fmt.Sprintf("list devices with name '%s'", data.Name.ValueString()), listErr, listResp),
			)
			return
		}

		if len(list.Results) == 0 {
			resp.Diagnostics.AddError(
				"Device not found",
				fmt.Sprintf("No device found with name '%s'", data.Name.ValueString()),
			)
			return
		}

		if len(list.Results) > 1 {
			resp.Diagnostics.AddError(
				"Multiple devices found",
				fmt.Sprintf("Multiple devices found with name '%s'. Please use ID or serial to identify the device uniquely.", data.Name.ValueString()),
			)
			return
		}

		device = &list.Results[0]
	} else if !data.Serial.IsNull() && data.Serial.ValueString() != "" {
		// Look up by serial
		tflog.Debug(ctx, "Looking up device by serial", map[string]interface{}{
			"serial": data.Serial.ValueString(),
		})

		list, listResp, listErr := d.client.DcimAPI.DcimDevicesList(ctx).Serial([]string{data.Serial.ValueString()}).Execute()
		if listErr != nil {
			resp.Diagnostics.AddError(
				"Error reading device",
				utils.FormatAPIError(fmt.Sprintf("list devices with serial '%s'", data.Serial.ValueString()), listErr, listResp),
			)
			return
		}

		if len(list.Results) == 0 {
			resp.Diagnostics.AddError(
				"Device not found",
				fmt.Sprintf("No device found with serial '%s'", data.Serial.ValueString()),
			)
			return
		}

		if len(list.Results) > 1 {
			resp.Diagnostics.AddError(
				"Multiple devices found",
				fmt.Sprintf("Multiple devices found with serial '%s'. Please use ID to identify the device uniquely.", data.Serial.ValueString()),
			)
			return
		}

		device = &list.Results[0]
	} else {
		resp.Diagnostics.AddError(
			"Missing identifier",
			"You must specify at least one of: id, name, or serial",
		)
		return
	}

	// Map the device to the data source model
	d.mapDeviceToDataSource(ctx, device, &data, resp)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// mapDeviceToDataSource maps a DeviceWithConfigContext from the API to the data source model.
func (d *DeviceDataSource) mapDeviceToDataSource(ctx context.Context, device *netbox.DeviceWithConfigContext, data *DeviceDataSourceModel, resp *datasource.ReadResponse) {
	data.ID = types.StringValue(fmt.Sprintf("%d", device.GetId()))

	// Handle name
	if device.HasName() && device.Name.Get() != nil && *device.Name.Get() != "" {
		data.Name = types.StringValue(*device.Name.Get())
	} else {
		data.Name = types.StringNull()
	}

	// Handle device_type (return slug)
	data.DeviceType = types.StringValue(device.DeviceType.GetSlug())

	// Handle role (return slug)
	data.Role = types.StringValue(device.Role.GetSlug())

	// Handle tenant
	if device.HasTenant() && device.Tenant.Get() != nil {
		data.Tenant = types.StringValue(device.Tenant.Get().GetSlug())
	} else {
		data.Tenant = types.StringNull()
	}

	// Handle platform
	if device.HasPlatform() && device.Platform.Get() != nil {
		data.Platform = types.StringValue(device.Platform.Get().GetSlug())
	} else {
		data.Platform = types.StringNull()
	}

	// Handle serial
	if device.HasSerial() && device.GetSerial() != "" {
		data.Serial = types.StringValue(device.GetSerial())
	} else {
		data.Serial = types.StringNull()
	}

	// Handle asset_tag
	if device.HasAssetTag() && device.AssetTag.Get() != nil && *device.AssetTag.Get() != "" {
		data.AssetTag = types.StringValue(*device.AssetTag.Get())
	} else {
		data.AssetTag = types.StringNull()
	}

	// Handle site (return slug)
	data.Site = types.StringValue(device.Site.GetSlug())

	// Handle location
	if device.HasLocation() && device.Location.Get() != nil {
		data.Location = types.StringValue(device.Location.Get().GetSlug())
	} else {
		data.Location = types.StringNull()
	}

	// Handle rack
	if device.HasRack() && device.Rack.Get() != nil {
		data.Rack = types.StringValue(device.Rack.Get().GetName())
	} else {
		data.Rack = types.StringNull()
	}

	// Handle position
	if device.HasPosition() && device.Position.Get() != nil {
		data.Position = types.Float64Value(*device.Position.Get())
	} else {
		data.Position = types.Float64Null()
	}

	// Handle face
	if device.HasFace() && device.Face != nil {
		data.Face = types.StringValue(string(device.Face.GetValue()))
	} else {
		data.Face = types.StringNull()
	}

	// Handle latitude
	if device.HasLatitude() && device.Latitude.Get() != nil {
		data.Latitude = types.Float64Value(*device.Latitude.Get())
	} else {
		data.Latitude = types.Float64Null()
	}

	// Handle longitude
	if device.HasLongitude() && device.Longitude.Get() != nil {
		data.Longitude = types.Float64Value(*device.Longitude.Get())
	} else {
		data.Longitude = types.Float64Null()
	}

	// Handle status
	if device.HasStatus() && device.Status != nil {
		data.Status = types.StringValue(string(device.Status.GetValue()))
	} else {
		data.Status = types.StringNull()
	}

	// Handle airflow
	if device.HasAirflow() && device.Airflow != nil {
		data.Airflow = types.StringValue(string(device.Airflow.GetValue()))
	} else {
		data.Airflow = types.StringNull()
	}

	// Handle vc_position
	if device.HasVcPosition() && device.VcPosition.Get() != nil {
		data.VcPosition = types.Int64Value(int64(*device.VcPosition.Get()))
	} else {
		data.VcPosition = types.Int64Null()
	}

	// Handle vc_priority
	if device.HasVcPriority() && device.VcPriority.Get() != nil {
		data.VcPriority = types.Int64Value(int64(*device.VcPriority.Get()))
	} else {
		data.VcPriority = types.Int64Null()
	}

	// Handle description
	if device.HasDescription() && device.GetDescription() != "" {
		data.Description = types.StringValue(device.GetDescription())
	} else {
		data.Description = types.StringNull()
	}

	// Handle comments
	if device.HasComments() && device.GetComments() != "" {
		data.Comments = types.StringValue(device.GetComments())
	} else {
		data.Comments = types.StringNull()
	}

	// Handle tags
	if device.HasTags() && len(device.GetTags()) > 0 {
		tags := utils.NestedTagsToTagModels(device.GetTags())
		tagsValue, tagDiags := types.SetValueFrom(ctx, utils.GetTagsAttributeType().ElemType, tags)
		resp.Diagnostics.Append(tagDiags...)
		if resp.Diagnostics.HasError() {
			return
		}
		data.Tags = tagsValue
	} else {
		data.Tags = types.SetNull(utils.GetTagsAttributeType().ElemType)
	}

	// Handle custom fields
	if device.HasCustomFields() && len(device.GetCustomFields()) > 0 {
		customFields := utils.MapToCustomFieldModels(device.GetCustomFields(), nil)
		customFieldsValue, cfDiags := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, customFields)
		resp.Diagnostics.Append(cfDiags...)
		if resp.Diagnostics.HasError() {
			return
		}
		data.CustomFields = customFieldsValue
	} else {
		data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)
	}
}
