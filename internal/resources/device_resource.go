// Package resources contains Terraform resource implementations for the Netbox provider.
package resources

import (
	"context"
	"fmt"

	"github.com/bab3l/go-netbox"
	"github.com/bab3l/terraform-provider-netbox/internal/netboxlookup"
	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/float64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &DeviceResource{}
var _ resource.ResourceWithImportState = &DeviceResource{}

func NewDeviceResource() resource.Resource {
	return &DeviceResource{}
}

// DeviceResource defines the resource implementation.
type DeviceResource struct {
	client *netbox.APIClient
}

// DeviceResourceModel describes the resource data model.
type DeviceResourceModel struct {
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

func (r *DeviceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_device"
}

func (r *DeviceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a device in Netbox. Devices represent physical or virtual hardware in your infrastructure, such as servers, switches, routers, and other network equipment.",

		Attributes: map[string]schema.Attribute{
			"id":          nbschema.IDAttribute("device"),
			"name":        nbschema.OptionalNameAttribute("device", 64),
			"device_type": nbschema.RequiredReferenceAttribute("device type", "ID or slug of the device type for this device. Required."),
			"role":        nbschema.RequiredReferenceAttribute("device role", "ID or slug of the device role. Required."),
			"tenant":      nbschema.ReferenceAttribute("tenant", "ID or slug of the tenant that owns this device."),
			"platform":    nbschema.ReferenceAttribute("platform", "ID or slug of the platform (operating system/software) running on this device."),
			"serial":      nbschema.SerialAttribute(),
			"asset_tag":   nbschema.AssetTagAttribute(),
			"site":        nbschema.RequiredReferenceAttribute("site", "ID or slug of the site where this device is located. Required."),
			"location":    nbschema.ReferenceAttribute("location", "ID or slug of the location within the site where this device is installed."),
			"rack":        nbschema.ReferenceAttribute("rack", "ID or name of the rack where this device is mounted."),
			"position": schema.Float64Attribute{
				MarkdownDescription: "Position in the rack (in rack units from the bottom). Must be a positive number.",
				Optional:            true,
				Validators: []validator.Float64{
					float64validator.AtLeast(0),
				},
			},
			"face": schema.StringAttribute{
				MarkdownDescription: "Which face of the rack the device is mounted on. Valid values: 'front', 'rear'.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("front", "rear", ""),
				},
			},
			"latitude": schema.Float64Attribute{
				MarkdownDescription: "GPS latitude coordinate in decimal format (xx.yyyyyy).",
				Optional:            true,
				Validators: []validator.Float64{
					float64validator.Between(-90, 90),
				},
			},
			"longitude": schema.Float64Attribute{
				MarkdownDescription: "GPS longitude coordinate in decimal format (xx.yyyyyy).",
				Optional:            true,
				Validators: []validator.Float64{
					float64validator.Between(-180, 180),
				},
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "Operational status of the device. Valid values: 'offline', 'active', 'planned', 'staged', 'failed', 'inventory', 'decommissioning'.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("active"),
				Validators: []validator.String{
					stringvalidator.OneOf("offline", "active", "planned", "staged", "failed", "inventory", "decommissioning"),
				},
			},
			"airflow": schema.StringAttribute{
				MarkdownDescription: "Direction of airflow through the device. Valid values: 'front-to-rear', 'rear-to-front', 'left-to-right', 'right-to-left', 'side-to-rear', 'passive', 'mixed'.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("front-to-rear", "rear-to-front", "left-to-right", "right-to-left", "side-to-rear", "passive", "mixed", ""),
				},
			},
			"vc_position": schema.Int64Attribute{
				MarkdownDescription: "Position within a virtual chassis (0-255).",
				Optional:            true,
				Validators: []validator.Int64{
					int64validator.Between(0, 255),
				},
			},
			"vc_priority": schema.Int64Attribute{
				MarkdownDescription: "Virtual chassis master election priority (0-255).",
				Optional:            true,
				Validators: []validator.Int64{
					int64validator.Between(0, 255),
				},
			},
			"description":   nbschema.DescriptionAttribute("device"),
			"comments":      nbschema.CommentsAttribute("device"),
			"tags":          nbschema.TagsAttribute(),
			"custom_fields": nbschema.CustomFieldsAttribute(),
		},
	}
}

func (r *DeviceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*netbox.APIClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *netbox.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *DeviceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data DeviceResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating device", map[string]interface{}{
		"name":        data.Name.ValueString(),
		"device_type": data.DeviceType.ValueString(),
		"role":        data.Role.ValueString(),
		"site":        data.Site.ValueString(),
	})

	// Look up required references
	deviceType, diags := netboxlookup.LookupDeviceType(ctx, r.client, data.DeviceType.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	role, diags := netboxlookup.LookupDeviceRole(ctx, r.client, data.Role.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	site, diags := netboxlookup.LookupSite(ctx, r.client, data.Site.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the device request with required fields
	deviceRequest := netbox.WritableDeviceWithConfigContextRequest{
		DeviceType: *deviceType,
		Role:       *role,
		Site:       *site,
	}

	// Set optional fields
	if !data.Name.IsNull() && !data.Name.IsUnknown() {
		deviceRequest.SetName(data.Name.ValueString())
	}

	if !data.Tenant.IsNull() && !data.Tenant.IsUnknown() {
		tenant, diags := netboxlookup.LookupTenant(ctx, r.client, data.Tenant.ValueString())
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		deviceRequest.SetTenant(*tenant)
	}

	if !data.Platform.IsNull() && !data.Platform.IsUnknown() {
		platform, diags := netboxlookup.LookupPlatform(ctx, r.client, data.Platform.ValueString())
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		deviceRequest.SetPlatform(*platform)
	}

	if !data.Serial.IsNull() && !data.Serial.IsUnknown() {
		serial := data.Serial.ValueString()
		deviceRequest.Serial = &serial
	}

	if !data.AssetTag.IsNull() && !data.AssetTag.IsUnknown() {
		deviceRequest.SetAssetTag(data.AssetTag.ValueString())
	}

	if !data.Location.IsNull() && !data.Location.IsUnknown() {
		location, diags := netboxlookup.LookupLocation(ctx, r.client, data.Location.ValueString())
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		deviceRequest.SetLocation(*location)
	}

	if !data.Rack.IsNull() && !data.Rack.IsUnknown() {
		rack, diags := netboxlookup.LookupRack(ctx, r.client, data.Rack.ValueString())
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		deviceRequest.SetRack(*rack)
	}

	if !data.Position.IsNull() && !data.Position.IsUnknown() {
		position := data.Position.ValueFloat64()
		deviceRequest.SetPosition(position)
	}

	if !data.Face.IsNull() && !data.Face.IsUnknown() && data.Face.ValueString() != "" {
		face := netbox.RackFace1(data.Face.ValueString())
		deviceRequest.Face = &face
	}

	if !data.Latitude.IsNull() && !data.Latitude.IsUnknown() {
		latitude := data.Latitude.ValueFloat64()
		deviceRequest.SetLatitude(latitude)
	}

	if !data.Longitude.IsNull() && !data.Longitude.IsUnknown() {
		longitude := data.Longitude.ValueFloat64()
		deviceRequest.SetLongitude(longitude)
	}

	if !data.Status.IsNull() && !data.Status.IsUnknown() {
		status := netbox.DeviceStatusValue(data.Status.ValueString())
		deviceRequest.Status = &status
	}

	if !data.Airflow.IsNull() && !data.Airflow.IsUnknown() && data.Airflow.ValueString() != "" {
		airflow := netbox.DeviceAirflowValue(data.Airflow.ValueString())
		deviceRequest.Airflow = &airflow
	}

	if !data.VcPosition.IsNull() && !data.VcPosition.IsUnknown() {
		vcPosition, err := utils.SafeInt32FromValue(data.VcPosition)
		if err != nil {
			resp.Diagnostics.AddError("Invalid value", fmt.Sprintf("VcPosition value overflow: %s", err))
			return
		}
		deviceRequest.SetVcPosition(vcPosition)
	}

	if !data.VcPriority.IsNull() && !data.VcPriority.IsUnknown() {
		vcPriority, err := utils.SafeInt32FromValue(data.VcPriority)
		if err != nil {
			resp.Diagnostics.AddError("Invalid value", fmt.Sprintf("VcPriority value overflow: %s", err))
			return
		}
		deviceRequest.SetVcPriority(vcPriority)
	}

	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		description := data.Description.ValueString()
		deviceRequest.Description = &description
	}

	if !data.Comments.IsNull() && !data.Comments.IsUnknown() {
		comments := data.Comments.ValueString()
		deviceRequest.Comments = &comments
	}

	// Handle tags
	if !data.Tags.IsNull() && !data.Tags.IsUnknown() {
		var tags []utils.TagModel
		resp.Diagnostics.Append(data.Tags.ElementsAs(ctx, &tags, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		deviceRequest.Tags = utils.TagsToNestedTagRequests(tags)
	}

	// Handle custom fields
	if !data.CustomFields.IsNull() && !data.CustomFields.IsUnknown() {
		var customFields []utils.CustomFieldModel
		resp.Diagnostics.Append(data.CustomFields.ElementsAs(ctx, &customFields, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		deviceRequest.CustomFields = utils.CustomFieldsToMap(customFields)
	}

	// Call the API
	device, httpResp, err := r.client.DcimAPI.DcimDevicesCreate(ctx).WritableDeviceWithConfigContextRequest(deviceRequest).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating device",
			utils.FormatAPIError("create device", err, httpResp),
		)
		return
	}

	tflog.Debug(ctx, "Created device", map[string]interface{}{
		"id":   device.GetId(),
		"name": device.GetName(),
	})

	// Map response to state
	r.mapDeviceToState(ctx, device, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DeviceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data DeviceResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Parse the ID
	deviceID := data.ID.ValueString()
	var deviceIDInt int32
	deviceIDInt, err := utils.ParseID(deviceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Device ID",
			fmt.Sprintf("Device ID must be a number, got: %s", deviceID),
		)
		return
	}

	tflog.Debug(ctx, "Reading device", map[string]interface{}{
		"id": deviceID,
	})

	// Call the API
	device, httpResp, err := r.client.DcimAPI.DcimDevicesRetrieve(ctx, deviceIDInt).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			tflog.Debug(ctx, "Device not found, removing from state", map[string]interface{}{
				"id": deviceID,
			})
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading device",
			utils.FormatAPIError(fmt.Sprintf("read device ID %s", deviceID), err, httpResp),
		)
		return
	}

	// Map response to state
	r.mapDeviceToState(ctx, device, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DeviceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data DeviceResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Parse the ID
	deviceID := data.ID.ValueString()
	var deviceIDInt int32
	deviceIDInt, err := utils.ParseID(deviceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Device ID",
			fmt.Sprintf("Device ID must be a number, got: %s", deviceID),
		)
		return
	}

	tflog.Debug(ctx, "Updating device", map[string]interface{}{
		"id":   deviceID,
		"name": data.Name.ValueString(),
	})

	// Look up required references
	deviceType, diags := netboxlookup.LookupDeviceType(ctx, r.client, data.DeviceType.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	role, diags := netboxlookup.LookupDeviceRole(ctx, r.client, data.Role.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	site, diags := netboxlookup.LookupSite(ctx, r.client, data.Site.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the device request with required fields
	deviceRequest := netbox.WritableDeviceWithConfigContextRequest{
		DeviceType: *deviceType,
		Role:       *role,
		Site:       *site,
	}

	// Set optional fields (same as Create)
	if !data.Name.IsNull() && !data.Name.IsUnknown() {
		deviceRequest.SetName(data.Name.ValueString())
	}

	if !data.Tenant.IsNull() && !data.Tenant.IsUnknown() {
		tenant, diags := netboxlookup.LookupTenant(ctx, r.client, data.Tenant.ValueString())
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		deviceRequest.SetTenant(*tenant)
	}

	if !data.Platform.IsNull() && !data.Platform.IsUnknown() {
		platform, diags := netboxlookup.LookupPlatform(ctx, r.client, data.Platform.ValueString())
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		deviceRequest.SetPlatform(*platform)
	}

	if !data.Serial.IsNull() && !data.Serial.IsUnknown() {
		serial := data.Serial.ValueString()
		deviceRequest.Serial = &serial
	}

	if !data.AssetTag.IsNull() && !data.AssetTag.IsUnknown() {
		deviceRequest.SetAssetTag(data.AssetTag.ValueString())
	}

	if !data.Location.IsNull() && !data.Location.IsUnknown() {
		location, diags := netboxlookup.LookupLocation(ctx, r.client, data.Location.ValueString())
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		deviceRequest.SetLocation(*location)
	}

	if !data.Rack.IsNull() && !data.Rack.IsUnknown() {
		rack, diags := netboxlookup.LookupRack(ctx, r.client, data.Rack.ValueString())
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		deviceRequest.SetRack(*rack)
	}

	if !data.Position.IsNull() && !data.Position.IsUnknown() {
		position := data.Position.ValueFloat64()
		deviceRequest.SetPosition(position)
	}

	if !data.Face.IsNull() && !data.Face.IsUnknown() && data.Face.ValueString() != "" {
		face := netbox.RackFace1(data.Face.ValueString())
		deviceRequest.Face = &face
	}

	if !data.Latitude.IsNull() && !data.Latitude.IsUnknown() {
		latitude := data.Latitude.ValueFloat64()
		deviceRequest.SetLatitude(latitude)
	}

	if !data.Longitude.IsNull() && !data.Longitude.IsUnknown() {
		longitude := data.Longitude.ValueFloat64()
		deviceRequest.SetLongitude(longitude)
	}

	if !data.Status.IsNull() && !data.Status.IsUnknown() {
		status := netbox.DeviceStatusValue(data.Status.ValueString())
		deviceRequest.Status = &status
	}

	if !data.Airflow.IsNull() && !data.Airflow.IsUnknown() && data.Airflow.ValueString() != "" {
		airflow := netbox.DeviceAirflowValue(data.Airflow.ValueString())
		deviceRequest.Airflow = &airflow
	}

	if !data.VcPosition.IsNull() && !data.VcPosition.IsUnknown() {
		vcPosition, err := utils.SafeInt32FromValue(data.VcPosition)
		if err != nil {
			resp.Diagnostics.AddError("Invalid value", fmt.Sprintf("VcPosition value overflow: %s", err))
			return
		}
		deviceRequest.SetVcPosition(vcPosition)
	}

	if !data.VcPriority.IsNull() && !data.VcPriority.IsUnknown() {
		vcPriority, err := utils.SafeInt32FromValue(data.VcPriority)
		if err != nil {
			resp.Diagnostics.AddError("Invalid value", fmt.Sprintf("VcPriority value overflow: %s", err))
			return
		}
		deviceRequest.SetVcPriority(vcPriority)
	}

	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		description := data.Description.ValueString()
		deviceRequest.Description = &description
	}

	if !data.Comments.IsNull() && !data.Comments.IsUnknown() {
		comments := data.Comments.ValueString()
		deviceRequest.Comments = &comments
	}

	// Handle tags
	if !data.Tags.IsNull() && !data.Tags.IsUnknown() {
		var tags []utils.TagModel
		resp.Diagnostics.Append(data.Tags.ElementsAs(ctx, &tags, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		deviceRequest.Tags = utils.TagsToNestedTagRequests(tags)
	}

	// Handle custom fields
	if !data.CustomFields.IsNull() && !data.CustomFields.IsUnknown() {
		var customFields []utils.CustomFieldModel
		resp.Diagnostics.Append(data.CustomFields.ElementsAs(ctx, &customFields, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		deviceRequest.CustomFields = utils.CustomFieldsToMap(customFields)
	}

	// Call the API
	device, httpResp, err := r.client.DcimAPI.DcimDevicesUpdate(ctx, deviceIDInt).WritableDeviceWithConfigContextRequest(deviceRequest).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating device",
			utils.FormatAPIError(fmt.Sprintf("update device ID %s", deviceID), err, httpResp),
		)
		return
	}

	tflog.Debug(ctx, "Updated device", map[string]interface{}{
		"id":   device.GetId(),
		"name": device.GetName(),
	})

	// Map response to state
	r.mapDeviceToState(ctx, device, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DeviceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data DeviceResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Parse the ID
	deviceID := data.ID.ValueString()
	var deviceIDInt int32
	deviceIDInt, err := utils.ParseID(deviceID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Device ID",
			fmt.Sprintf("Device ID must be a number, got: %s", deviceID),
		)
		return
	}

	tflog.Debug(ctx, "Deleting device", map[string]interface{}{
		"id": deviceID,
	})

	// Call the API
	httpResp, err := r.client.DcimAPI.DcimDevicesDestroy(ctx, deviceIDInt).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			// Already deleted
			tflog.Debug(ctx, "Device already deleted", map[string]interface{}{
				"id": deviceID,
			})
			return
		}
		resp.Diagnostics.AddError(
			"Error deleting device",
			utils.FormatAPIError(fmt.Sprintf("delete device ID %s", deviceID), err, httpResp),
		)
		return
	}

	tflog.Debug(ctx, "Deleted device", map[string]interface{}{
		"id": deviceID,
	})
}

func (r *DeviceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// mapDeviceToState maps a DeviceWithConfigContext from the API to the Terraform state model.
func (r *DeviceResource) mapDeviceToState(ctx context.Context, device *netbox.DeviceWithConfigContext, data *DeviceResourceModel, diags *diag.Diagnostics) {
	data.ID = types.StringValue(fmt.Sprintf("%d", device.GetId()))

	// Handle name
	if device.HasName() && device.Name.Get() != nil && *device.Name.Get() != "" {
		data.Name = types.StringValue(*device.Name.Get())
	} else if !data.Name.IsNull() {
		data.Name = types.StringNull()
	}

	// Handle device_type - preserve the original input value (slug or ID)
	if data.DeviceType.IsNull() || data.DeviceType.IsUnknown() {
		data.DeviceType = types.StringValue(fmt.Sprintf("%d", device.DeviceType.GetId()))
	}
	// Otherwise keep the original value the user provided

	// Handle role - preserve the original input value (slug or ID)
	if data.Role.IsNull() || data.Role.IsUnknown() {
		data.Role = types.StringValue(fmt.Sprintf("%d", device.Role.GetId()))
	}
	// Otherwise keep the original value the user provided

	// Handle tenant - preserve the original input value
	switch {
	case device.HasTenant() && device.Tenant.Get() != nil:
		if data.Tenant.IsNull() || data.Tenant.IsUnknown() {
			data.Tenant = types.StringValue(fmt.Sprintf("%d", device.Tenant.Get().GetId()))
		}
	case !data.Tenant.IsNull() && !data.Tenant.IsUnknown():
		// User had a value but API says null - shouldn't happen normally
	default:
		data.Tenant = types.StringNull()
	}

	// Handle platform - preserve the original input value
	switch {
	case device.HasPlatform() && device.Platform.Get() != nil:
		if data.Platform.IsNull() || data.Platform.IsUnknown() {
			data.Platform = types.StringValue(fmt.Sprintf("%d", device.Platform.Get().GetId()))
		}
	case !data.Platform.IsNull() && !data.Platform.IsUnknown():
		// User had a value but API says null
	default:
		data.Platform = types.StringNull()
	}

	// Handle serial
	if device.HasSerial() && device.GetSerial() != "" {
		data.Serial = types.StringValue(device.GetSerial())
	} else if !data.Serial.IsNull() {
		data.Serial = types.StringNull()
	}

	// Handle asset_tag
	if device.HasAssetTag() && device.AssetTag.Get() != nil && *device.AssetTag.Get() != "" {
		data.AssetTag = types.StringValue(*device.AssetTag.Get())
	} else if !data.AssetTag.IsNull() {
		data.AssetTag = types.StringNull()
	}

	// Handle site - preserve the original input value
	if data.Site.IsNull() || data.Site.IsUnknown() {
		data.Site = types.StringValue(fmt.Sprintf("%d", device.Site.GetId()))
	}
	// Otherwise keep the original value the user provided

	// Handle location - preserve the original input value
	switch {
	case device.HasLocation() && device.Location.Get() != nil:
		if data.Location.IsNull() || data.Location.IsUnknown() {
			data.Location = types.StringValue(fmt.Sprintf("%d", device.Location.Get().GetId()))
		}
	case !data.Location.IsNull() && !data.Location.IsUnknown():
		// User had a value but API says null
	default:
		data.Location = types.StringNull()
	}

	// Handle rack - preserve the original input value
	switch {
	case device.HasRack() && device.Rack.Get() != nil:
		if data.Rack.IsNull() || data.Rack.IsUnknown() {
			data.Rack = types.StringValue(fmt.Sprintf("%d", device.Rack.Get().GetId()))
		}
	case !data.Rack.IsNull() && !data.Rack.IsUnknown():
		// User had a value but API says null
	default:
		data.Rack = types.StringNull()
	}

	// Handle position
	if device.HasPosition() && device.Position.Get() != nil {
		data.Position = types.Float64Value(*device.Position.Get())
	} else if !data.Position.IsNull() {
		data.Position = types.Float64Null()
	}

	// Handle face
	if device.HasFace() && device.Face != nil {
		data.Face = types.StringValue(string(device.Face.GetValue()))
	} else if !data.Face.IsNull() {
		data.Face = types.StringNull()
	}

	// Handle latitude
	if device.HasLatitude() && device.Latitude.Get() != nil {
		data.Latitude = types.Float64Value(*device.Latitude.Get())
	} else if !data.Latitude.IsNull() {
		data.Latitude = types.Float64Null()
	}

	// Handle longitude
	if device.HasLongitude() && device.Longitude.Get() != nil {
		data.Longitude = types.Float64Value(*device.Longitude.Get())
	} else if !data.Longitude.IsNull() {
		data.Longitude = types.Float64Null()
	}

	// Handle status
	if device.HasStatus() && device.Status != nil {
		data.Status = types.StringValue(string(device.Status.GetValue()))
	}

	// Handle airflow
	if device.HasAirflow() && device.Airflow != nil {
		data.Airflow = types.StringValue(string(device.Airflow.GetValue()))
	} else if !data.Airflow.IsNull() {
		data.Airflow = types.StringNull()
	}

	// Handle vc_position
	if device.HasVcPosition() && device.VcPosition.Get() != nil {
		data.VcPosition = types.Int64Value(int64(*device.VcPosition.Get()))
	} else if !data.VcPosition.IsNull() {
		data.VcPosition = types.Int64Null()
	}

	// Handle vc_priority
	if device.HasVcPriority() && device.VcPriority.Get() != nil {
		data.VcPriority = types.Int64Value(int64(*device.VcPriority.Get()))
	} else if !data.VcPriority.IsNull() {
		data.VcPriority = types.Int64Null()
	}

	// Handle description
	if device.HasDescription() && device.GetDescription() != "" {
		data.Description = types.StringValue(device.GetDescription())
	} else if !data.Description.IsNull() {
		data.Description = types.StringNull()
	}

	// Handle comments
	if device.HasComments() && device.GetComments() != "" {
		data.Comments = types.StringValue(device.GetComments())
	} else if !data.Comments.IsNull() {
		data.Comments = types.StringNull()
	}

	// Handle tags
	if device.HasTags() {
		tags := utils.NestedTagsToTagModels(device.GetTags())
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
	if device.HasCustomFields() && !data.CustomFields.IsNull() {
		var stateCustomFields []utils.CustomFieldModel
		cfDiags := data.CustomFields.ElementsAs(ctx, &stateCustomFields, false)
		diags.Append(cfDiags...)
		if diags.HasError() {
			return
		}

		customFields := utils.MapToCustomFieldModels(device.GetCustomFields(), stateCustomFields)
		customFieldsValue, cfValueDiags := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, customFields)
		diags.Append(cfValueDiags...)
		if diags.HasError() {
			return
		}
		data.CustomFields = customFieldsValue
	} else if data.CustomFields.IsNull() {
		data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)
	}
}
