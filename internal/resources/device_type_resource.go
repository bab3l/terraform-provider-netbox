// Package resources contains Terraform resource implementations for the Netbox provider.
package resources

import (
	"context"
	"fmt"

	"github.com/bab3l/go-netbox"
	"github.com/hashicorp/terraform-plugin-framework-validators/float64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/bab3l/terraform-provider-netbox/internal/netboxlookup"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/bab3l/terraform-provider-netbox/internal/validators"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &DeviceTypeResource{}
var _ resource.ResourceWithImportState = &DeviceTypeResource{}

func NewDeviceTypeResource() resource.Resource {
	return &DeviceTypeResource{}
}

// DeviceTypeResource defines the resource implementation.
type DeviceTypeResource struct {
	client *netbox.APIClient
}

// DeviceTypeResourceModel describes the resource data model.
type DeviceTypeResourceModel struct {
	ID                     types.String  `tfsdk:"id"`
	Manufacturer           types.String  `tfsdk:"manufacturer"`
	Model                  types.String  `tfsdk:"model"`
	Slug                   types.String  `tfsdk:"slug"`
	DefaultPlatform        types.String  `tfsdk:"default_platform"`
	PartNumber             types.String  `tfsdk:"part_number"`
	UHeight                types.Float64 `tfsdk:"u_height"`
	ExcludeFromUtilization types.Bool    `tfsdk:"exclude_from_utilization"`
	IsFullDepth            types.Bool    `tfsdk:"is_full_depth"`
	SubdeviceRole          types.String  `tfsdk:"subdevice_role"`
	Airflow                types.String  `tfsdk:"airflow"`
	Weight                 types.Float64 `tfsdk:"weight"`
	WeightUnit             types.String  `tfsdk:"weight_unit"`
	Description            types.String  `tfsdk:"description"`
	Comments               types.String  `tfsdk:"comments"`
	Tags                   types.Set     `tfsdk:"tags"`
	CustomFields           types.Set     `tfsdk:"custom_fields"`
}

func (r *DeviceTypeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_device_type"
}

func (r *DeviceTypeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a device type in Netbox. Device types define the make and model of physical hardware, including specifications like rack height, airflow direction, and weight. Device types serve as templates when creating new devices.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Unique identifier for the device type (assigned by Netbox).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"manufacturer": schema.StringAttribute{
				MarkdownDescription: "ID or slug of the manufacturer of this device type. Required.",
				Required:            true,
			},
			"model": schema.StringAttribute{
				MarkdownDescription: "Model name/number of the device type (e.g., 'Catalyst 9300-48P', 'PowerEdge R640').",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 100),
				},
			},
			"slug": schema.StringAttribute{
				MarkdownDescription: "URL-friendly identifier for the device type. Must be unique and contain only alphanumeric characters, hyphens, and underscores.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 100),
					validators.ValidSlug(),
				},
			},
			"default_platform": schema.StringAttribute{
				MarkdownDescription: "ID or slug of the default platform for devices of this type.",
				Optional:            true,
			},
			"part_number": schema.StringAttribute{
				MarkdownDescription: "Discrete manufacturer part number for ordering/inventory purposes.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(50),
				},
			},
			"u_height": schema.Float64Attribute{
				MarkdownDescription: "Height of the device in rack units (U). Defaults to 1.0. Use 0 for devices that don't consume rack space.",
				Optional:            true,
				Computed:            true,
				Default:             float64default.StaticFloat64(1.0),
				Validators: []validator.Float64{
					float64validator.AtLeast(0),
					float64validator.AtMost(32767),
				},
			},
			"exclude_from_utilization": schema.BoolAttribute{
				MarkdownDescription: "If true, devices of this type are excluded when calculating rack utilization. Defaults to false.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"is_full_depth": schema.BoolAttribute{
				MarkdownDescription: "If true, device consumes both front and rear rack faces. Defaults to true.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"subdevice_role": schema.StringAttribute{
				MarkdownDescription: "Role of this device type as a subdevice. Valid values: 'parent' or 'child'. Leave empty if not a subdevice.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("parent", "child", ""),
				},
			},
			"airflow": schema.StringAttribute{
				MarkdownDescription: "Direction of airflow through the device. Valid values: 'front-to-rear', 'rear-to-front', 'left-to-right', 'right-to-left', 'side-to-rear', 'passive', 'mixed'.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("front-to-rear", "rear-to-front", "left-to-right", "right-to-left", "side-to-rear", "passive", "mixed", ""),
				},
			},
			"weight": schema.Float64Attribute{
				MarkdownDescription: "Weight of the device (use with weight_unit).",
				Optional:            true,
				Validators: []validator.Float64{
					float64validator.AtLeast(0),
				},
			},
			"weight_unit": schema.StringAttribute{
				MarkdownDescription: "Unit for weight measurement. Valid values: 'kg', 'g', 'lb', 'oz'.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("kg", "g", "lb", "oz", ""),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Brief description of the device type.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(200),
				},
			},
			"comments": schema.StringAttribute{
				MarkdownDescription: "Additional comments or notes about this device type. Supports Markdown formatting.",
				Optional:            true,
			},
			"tags": schema.SetNestedAttribute{
				MarkdownDescription: "Tags assigned to this device type.",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "Name of the existing tag.",
							Required:            true,
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 100),
							},
						},
						"slug": schema.StringAttribute{
							MarkdownDescription: "Slug of the existing tag.",
							Required:            true,
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 100),
								validators.ValidSlug(),
							},
						},
					},
				},
			},
			"custom_fields": schema.SetNestedAttribute{
				MarkdownDescription: "Custom fields assigned to this device type.",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "Name of the custom field.",
							Required:            true,
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 50),
								validators.ValidCustomFieldName(),
							},
						},
						"type": schema.StringAttribute{
							MarkdownDescription: "Type of the custom field.",
							Required:            true,
							Validators: []validator.String{
								validators.ValidCustomFieldType(),
							},
						},
						"value": schema.StringAttribute{
							MarkdownDescription: "Value of the custom field.",
							Required:            true,
							Validators: []validator.String{
								stringvalidator.LengthAtMost(1000),
							},
						},
					},
				},
			},
		},
	}
}

func (r *DeviceTypeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *DeviceTypeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data DeviceTypeResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating device type", map[string]interface{}{
		"model": data.Model.ValueString(),
		"slug":  data.Slug.ValueString(),
	})

	// Look up manufacturer
	manufacturer, diags := netboxlookup.LookupManufacturerBrief(ctx, r.client, data.Manufacturer.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the device type request
	deviceTypeRequest := netbox.WritableDeviceTypeRequest{
		Manufacturer: *manufacturer,
		Model:        data.Model.ValueString(),
		Slug:         data.Slug.ValueString(),
	}

	// Set optional fields
	if !data.DefaultPlatform.IsNull() && !data.DefaultPlatform.IsUnknown() {
		platform, diags := netboxlookup.LookupPlatformBrief(ctx, r.client, data.DefaultPlatform.ValueString())
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		deviceTypeRequest.SetDefaultPlatform(*platform)
	}

	if !data.PartNumber.IsNull() && !data.PartNumber.IsUnknown() {
		partNumber := data.PartNumber.ValueString()
		deviceTypeRequest.PartNumber = &partNumber
	}

	if !data.UHeight.IsNull() && !data.UHeight.IsUnknown() {
		uHeight := data.UHeight.ValueFloat64()
		deviceTypeRequest.UHeight = &uHeight
	}

	if !data.ExcludeFromUtilization.IsNull() && !data.ExcludeFromUtilization.IsUnknown() {
		excludeFromUtilization := data.ExcludeFromUtilization.ValueBool()
		deviceTypeRequest.ExcludeFromUtilization = &excludeFromUtilization
	}

	if !data.IsFullDepth.IsNull() && !data.IsFullDepth.IsUnknown() {
		isFullDepth := data.IsFullDepth.ValueBool()
		deviceTypeRequest.IsFullDepth = &isFullDepth
	}

	if !data.SubdeviceRole.IsNull() && !data.SubdeviceRole.IsUnknown() && data.SubdeviceRole.ValueString() != "" {
		subdeviceRole := netbox.ParentChildStatus1(data.SubdeviceRole.ValueString())
		deviceTypeRequest.SubdeviceRole = &subdeviceRole
	}

	if !data.Airflow.IsNull() && !data.Airflow.IsUnknown() && data.Airflow.ValueString() != "" {
		airflow := netbox.DeviceAirflowValue(data.Airflow.ValueString())
		deviceTypeRequest.Airflow = &airflow
	}

	if !data.Weight.IsNull() && !data.Weight.IsUnknown() {
		weight := data.Weight.ValueFloat64()
		deviceTypeRequest.Weight = *netbox.NewNullableFloat64(&weight)
	}

	if !data.WeightUnit.IsNull() && !data.WeightUnit.IsUnknown() && data.WeightUnit.ValueString() != "" {
		weightUnit := netbox.DeviceTypeWeightUnitValue(data.WeightUnit.ValueString())
		deviceTypeRequest.WeightUnit = &weightUnit
	}

	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		description := data.Description.ValueString()
		deviceTypeRequest.Description = &description
	}

	if !data.Comments.IsNull() && !data.Comments.IsUnknown() {
		comments := data.Comments.ValueString()
		deviceTypeRequest.Comments = &comments
	}

	// Handle tags
	if !data.Tags.IsNull() && !data.Tags.IsUnknown() {
		var tags []utils.TagModel
		resp.Diagnostics.Append(data.Tags.ElementsAs(ctx, &tags, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		deviceTypeRequest.Tags = utils.TagsToNestedTagRequests(tags)
	}

	// Handle custom fields
	if !data.CustomFields.IsNull() && !data.CustomFields.IsUnknown() {
		var customFields []utils.CustomFieldModel
		resp.Diagnostics.Append(data.CustomFields.ElementsAs(ctx, &customFields, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		deviceTypeRequest.CustomFields = utils.CustomFieldsToMap(customFields)
	}

	// Call the API
	deviceType, httpResp, err := r.client.DcimAPI.DcimDeviceTypesCreate(ctx).WritableDeviceTypeRequest(deviceTypeRequest).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating device type",
			utils.FormatAPIError("create device type", err, httpResp),
		)
		return
	}

	tflog.Debug(ctx, "Created device type", map[string]interface{}{
		"id":    deviceType.GetId(),
		"model": deviceType.GetModel(),
	})

	// Map response to state
	r.mapDeviceTypeToState(ctx, deviceType, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DeviceTypeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data DeviceTypeResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Parse the ID
	deviceTypeID := data.ID.ValueString()
	var deviceTypeIDInt int32
	if _, err := fmt.Sscanf(deviceTypeID, "%d", &deviceTypeIDInt); err != nil {
		resp.Diagnostics.AddError(
			"Invalid Device Type ID",
			fmt.Sprintf("Device Type ID must be a number, got: %s", deviceTypeID),
		)
		return
	}

	tflog.Debug(ctx, "Reading device type", map[string]interface{}{
		"id": deviceTypeID,
	})

	// Call the API
	deviceType, httpResp, err := r.client.DcimAPI.DcimDeviceTypesRetrieve(ctx, deviceTypeIDInt).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			tflog.Debug(ctx, "Device type not found, removing from state", map[string]interface{}{
				"id": deviceTypeID,
			})
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading device type",
			utils.FormatAPIError(fmt.Sprintf("read device type ID %s", deviceTypeID), err, httpResp),
		)
		return
	}

	// Map response to state
	r.mapDeviceTypeToState(ctx, deviceType, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DeviceTypeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data DeviceTypeResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Parse the ID
	deviceTypeID := data.ID.ValueString()
	var deviceTypeIDInt int32
	if _, err := fmt.Sscanf(deviceTypeID, "%d", &deviceTypeIDInt); err != nil {
		resp.Diagnostics.AddError(
			"Invalid Device Type ID",
			fmt.Sprintf("Device Type ID must be a number, got: %s", deviceTypeID),
		)
		return
	}

	tflog.Debug(ctx, "Updating device type", map[string]interface{}{
		"id":    deviceTypeID,
		"model": data.Model.ValueString(),
		"slug":  data.Slug.ValueString(),
	})

	// Look up manufacturer
	manufacturer, diags := netboxlookup.LookupManufacturerBrief(ctx, r.client, data.Manufacturer.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the device type request
	deviceTypeRequest := netbox.WritableDeviceTypeRequest{
		Manufacturer: *manufacturer,
		Model:        data.Model.ValueString(),
		Slug:         data.Slug.ValueString(),
	}

	// Set optional fields (same as Create)
	if !data.DefaultPlatform.IsNull() && !data.DefaultPlatform.IsUnknown() {
		platform, diags := netboxlookup.LookupPlatformBrief(ctx, r.client, data.DefaultPlatform.ValueString())
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		deviceTypeRequest.SetDefaultPlatform(*platform)
	}

	if !data.PartNumber.IsNull() && !data.PartNumber.IsUnknown() {
		partNumber := data.PartNumber.ValueString()
		deviceTypeRequest.PartNumber = &partNumber
	}

	if !data.UHeight.IsNull() && !data.UHeight.IsUnknown() {
		uHeight := data.UHeight.ValueFloat64()
		deviceTypeRequest.UHeight = &uHeight
	}

	if !data.ExcludeFromUtilization.IsNull() && !data.ExcludeFromUtilization.IsUnknown() {
		excludeFromUtilization := data.ExcludeFromUtilization.ValueBool()
		deviceTypeRequest.ExcludeFromUtilization = &excludeFromUtilization
	}

	if !data.IsFullDepth.IsNull() && !data.IsFullDepth.IsUnknown() {
		isFullDepth := data.IsFullDepth.ValueBool()
		deviceTypeRequest.IsFullDepth = &isFullDepth
	}

	if !data.SubdeviceRole.IsNull() && !data.SubdeviceRole.IsUnknown() && data.SubdeviceRole.ValueString() != "" {
		subdeviceRole := netbox.ParentChildStatus1(data.SubdeviceRole.ValueString())
		deviceTypeRequest.SubdeviceRole = &subdeviceRole
	}

	if !data.Airflow.IsNull() && !data.Airflow.IsUnknown() && data.Airflow.ValueString() != "" {
		airflow := netbox.DeviceAirflowValue(data.Airflow.ValueString())
		deviceTypeRequest.Airflow = &airflow
	}

	if !data.Weight.IsNull() && !data.Weight.IsUnknown() {
		weight := data.Weight.ValueFloat64()
		deviceTypeRequest.Weight = *netbox.NewNullableFloat64(&weight)
	}

	if !data.WeightUnit.IsNull() && !data.WeightUnit.IsUnknown() && data.WeightUnit.ValueString() != "" {
		weightUnit := netbox.DeviceTypeWeightUnitValue(data.WeightUnit.ValueString())
		deviceTypeRequest.WeightUnit = &weightUnit
	}

	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		description := data.Description.ValueString()
		deviceTypeRequest.Description = &description
	}

	if !data.Comments.IsNull() && !data.Comments.IsUnknown() {
		comments := data.Comments.ValueString()
		deviceTypeRequest.Comments = &comments
	}

	// Handle tags
	if !data.Tags.IsNull() && !data.Tags.IsUnknown() {
		var tags []utils.TagModel
		resp.Diagnostics.Append(data.Tags.ElementsAs(ctx, &tags, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		deviceTypeRequest.Tags = utils.TagsToNestedTagRequests(tags)
	}

	// Handle custom fields
	if !data.CustomFields.IsNull() && !data.CustomFields.IsUnknown() {
		var customFields []utils.CustomFieldModel
		resp.Diagnostics.Append(data.CustomFields.ElementsAs(ctx, &customFields, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		deviceTypeRequest.CustomFields = utils.CustomFieldsToMap(customFields)
	}

	// Call the API
	deviceType, httpResp, err := r.client.DcimAPI.DcimDeviceTypesUpdate(ctx, deviceTypeIDInt).WritableDeviceTypeRequest(deviceTypeRequest).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating device type",
			utils.FormatAPIError(fmt.Sprintf("update device type ID %s", deviceTypeID), err, httpResp),
		)
		return
	}

	tflog.Debug(ctx, "Updated device type", map[string]interface{}{
		"id":    deviceType.GetId(),
		"model": deviceType.GetModel(),
	})

	// Map response to state
	r.mapDeviceTypeToState(ctx, deviceType, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DeviceTypeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data DeviceTypeResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Parse the ID
	deviceTypeID := data.ID.ValueString()
	var deviceTypeIDInt int32
	if _, err := fmt.Sscanf(deviceTypeID, "%d", &deviceTypeIDInt); err != nil {
		resp.Diagnostics.AddError(
			"Invalid Device Type ID",
			fmt.Sprintf("Device Type ID must be a number, got: %s", deviceTypeID),
		)
		return
	}

	tflog.Debug(ctx, "Deleting device type", map[string]interface{}{
		"id": deviceTypeID,
	})

	// Call the API
	httpResp, err := r.client.DcimAPI.DcimDeviceTypesDestroy(ctx, deviceTypeIDInt).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			// Already deleted
			tflog.Debug(ctx, "Device type already deleted", map[string]interface{}{
				"id": deviceTypeID,
			})
			return
		}
		resp.Diagnostics.AddError(
			"Error deleting device type",
			utils.FormatAPIError(fmt.Sprintf("delete device type ID %s", deviceTypeID), err, httpResp),
		)
		return
	}

	tflog.Debug(ctx, "Deleted device type", map[string]interface{}{
		"id": deviceTypeID,
	})
}

func (r *DeviceTypeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// mapDeviceTypeToState maps a DeviceType from the API to the Terraform state model.
func (r *DeviceTypeResource) mapDeviceTypeToState(ctx context.Context, deviceType *netbox.DeviceType, data *DeviceTypeResourceModel, diags *diag.Diagnostics) {
	data.ID = types.StringValue(fmt.Sprintf("%d", deviceType.GetId()))
	data.Model = types.StringValue(deviceType.GetModel())
	data.Slug = types.StringValue(deviceType.GetSlug())

	// Handle manufacturer - preserve the original input value (slug or ID)
	// The user provides either an ID or slug, and we should keep what they provided
	// Only update if the value was null/unknown (shouldn't happen for required field)
	if data.Manufacturer.IsNull() || data.Manufacturer.IsUnknown() {
		data.Manufacturer = types.StringValue(fmt.Sprintf("%d", deviceType.Manufacturer.GetId()))
	}
	// Otherwise keep the original value the user provided

	// Handle default platform - preserve the original input value
	if deviceType.HasDefaultPlatform() && deviceType.DefaultPlatform.Get() != nil {
		// Only update if the value was null/unknown
		if data.DefaultPlatform.IsNull() || data.DefaultPlatform.IsUnknown() {
			data.DefaultPlatform = types.StringValue(fmt.Sprintf("%d", deviceType.DefaultPlatform.Get().GetId()))
		}
	} else if !data.DefaultPlatform.IsNull() && !data.DefaultPlatform.IsUnknown() {
		// API says no platform but user had one configured - keep user's value
		// This shouldn't happen in normal operation
	} else {
		data.DefaultPlatform = types.StringNull()
	}

	// Handle part number
	if deviceType.HasPartNumber() && deviceType.GetPartNumber() != "" {
		data.PartNumber = types.StringValue(deviceType.GetPartNumber())
	} else if !data.PartNumber.IsNull() {
		data.PartNumber = types.StringNull()
	}

	// Handle u_height
	if deviceType.HasUHeight() {
		data.UHeight = types.Float64Value(*deviceType.UHeight)
	}

	// Handle exclude_from_utilization
	if deviceType.HasExcludeFromUtilization() {
		data.ExcludeFromUtilization = types.BoolValue(deviceType.GetExcludeFromUtilization())
	}

	// Handle is_full_depth
	if deviceType.HasIsFullDepth() {
		data.IsFullDepth = types.BoolValue(deviceType.GetIsFullDepth())
	}

	// Handle subdevice_role
	if deviceType.HasSubdeviceRole() && deviceType.SubdeviceRole.Get() != nil {
		data.SubdeviceRole = types.StringValue(string(deviceType.SubdeviceRole.Get().GetValue()))
	} else if !data.SubdeviceRole.IsNull() {
		data.SubdeviceRole = types.StringNull()
	}

	// Handle airflow
	if deviceType.HasAirflow() && deviceType.Airflow.Get() != nil {
		data.Airflow = types.StringValue(string(deviceType.Airflow.Get().GetValue()))
	} else if !data.Airflow.IsNull() {
		data.Airflow = types.StringNull()
	}

	// Handle weight
	if deviceType.HasWeight() && deviceType.Weight.Get() != nil {
		data.Weight = types.Float64Value(*deviceType.Weight.Get())
	} else if !data.Weight.IsNull() {
		data.Weight = types.Float64Null()
	}

	// Handle weight_unit
	if deviceType.HasWeightUnit() && deviceType.WeightUnit.Get() != nil {
		data.WeightUnit = types.StringValue(string(deviceType.WeightUnit.Get().GetValue()))
	} else if !data.WeightUnit.IsNull() {
		data.WeightUnit = types.StringNull()
	}

	// Handle description
	if deviceType.HasDescription() && deviceType.GetDescription() != "" {
		data.Description = types.StringValue(deviceType.GetDescription())
	} else if !data.Description.IsNull() {
		data.Description = types.StringNull()
	}

	// Handle comments
	if deviceType.HasComments() && deviceType.GetComments() != "" {
		data.Comments = types.StringValue(deviceType.GetComments())
	} else if !data.Comments.IsNull() {
		data.Comments = types.StringNull()
	}

	// Handle tags
	if deviceType.HasTags() {
		tags := utils.NestedTagsToTagModels(deviceType.GetTags())
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
	if deviceType.HasCustomFields() && !data.CustomFields.IsNull() {
		var stateCustomFields []utils.CustomFieldModel
		cfDiags := data.CustomFields.ElementsAs(ctx, &stateCustomFields, false)
		diags.Append(cfDiags...)
		if diags.HasError() {
			return
		}

		customFields := utils.MapToCustomFieldModels(deviceType.GetCustomFields(), stateCustomFields)
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
