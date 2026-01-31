// Package resources contains Terraform resource implementations for the Netbox provider.

package resources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bab3l/go-netbox"
	"github.com/bab3l/terraform-provider-netbox/internal/netboxlookup"
	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/float64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64default"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &DeviceTypeResource{}
	_ resource.ResourceWithImportState = &DeviceTypeResource{}
	_ resource.ResourceWithIdentity    = &DeviceTypeResource{}
)

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
			"id":               nbschema.IDAttribute("device type"),
			"manufacturer":     nbschema.RequiredReferenceAttributeWithDiffSuppress("manufacturer", "ID or slug of the manufacturer of this device type. Required."),
			"model":            nbschema.ModelAttribute("device type", 100),
			"slug":             nbschema.SlugAttribute("device type"),
			"default_platform": nbschema.ReferenceAttributeWithDiffSuppress("platform", "ID or slug of the default platform for devices of this type."),
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
			"exclude_from_utilization": nbschema.BoolAttributeWithDefault("If true, devices of this type are excluded when calculating rack utilization. Defaults to false.", false),
			"is_full_depth":            nbschema.BoolAttributeWithDefault("If true, device consumes both front and rear rack faces. Defaults to true.", true),
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
			"description":   nbschema.DescriptionAttribute("device type"),
			"comments":      nbschema.CommentsAttribute("device type"),
			"tags":          nbschema.TagsSlugAttribute(),
			"custom_fields": nbschema.CustomFieldsAttribute(),
		},
	}
}

func (r *DeviceTypeResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = nbschema.ImportIdentityWithCustomFieldsSchema()
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
	manufacturer, diags := netboxlookup.LookupManufacturer(ctx, r.client, data.Manufacturer.ValueString())
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
		platform, diags := netboxlookup.LookupPlatform(ctx, r.client, data.DefaultPlatform.ValueString())
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		deviceTypeRequest.SetDefaultPlatform(*platform)
	}

	if !data.PartNumber.IsNull() && !data.PartNumber.IsUnknown() {
		partNumber := data.PartNumber.ValueString()
		deviceTypeRequest.PartNumber = &partNumber
	} else if data.PartNumber.IsNull() {
		// NetBox may reject explicit nulls for some fields, but for updates we need to
		// clear fields when removed from config. The OpenAPI client represents nullable
		// strings as pointers, so we use an empty string and map it back to null in state.
		empty := ""
		deviceTypeRequest.PartNumber = &empty
	}

	if !data.UHeight.IsNull() && !data.UHeight.IsUnknown() {
		uHeight := data.UHeight.ValueFloat64()
		deviceTypeRequest.UHeight = &uHeight
	} else if data.UHeight.IsNull() {
		// NetBox defaults u_height to 1.0 when omitted; restore that default on update.
		defaultUHeight := 1.0
		deviceTypeRequest.UHeight = &defaultUHeight
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
	} else if data.SubdeviceRole.IsNull() {
		empty := netbox.ParentChildStatus1("")
		deviceTypeRequest.SubdeviceRole = &empty
	}

	if !data.Airflow.IsNull() && !data.Airflow.IsUnknown() && data.Airflow.ValueString() != "" {
		airflow := netbox.DeviceAirflowValue(data.Airflow.ValueString())
		deviceTypeRequest.Airflow = &airflow
	} else if data.Airflow.IsNull() {
		empty := netbox.DeviceAirflowValue("")
		deviceTypeRequest.Airflow = &empty
	}

	if !data.Weight.IsNull() && !data.Weight.IsUnknown() {
		weight := data.Weight.ValueFloat64()
		deviceTypeRequest.Weight = *netbox.NewNullableFloat64(&weight)
	} else if data.Weight.IsNull() {
		deviceTypeRequest.SetWeightNil()
	}

	if !data.WeightUnit.IsNull() && !data.WeightUnit.IsUnknown() && data.WeightUnit.ValueString() != "" {
		weightUnit := netbox.DeviceTypeWeightUnitValue(data.WeightUnit.ValueString())
		deviceTypeRequest.WeightUnit = &weightUnit
	} else if data.WeightUnit.IsNull() {
		empty := netbox.DeviceTypeWeightUnitValue("")
		deviceTypeRequest.WeightUnit = &empty
	}

	// Apply description and comments
	utils.ApplyDescription(&deviceTypeRequest, data.Description)
	utils.ApplyComments(&deviceTypeRequest, data.Comments)

	// Store plan values for filter-to-owned pattern
	planTags := data.Tags
	planCustomFields := data.CustomFields

	// Apply tags and custom fields
	utils.ApplyTagsFromSlugs(ctx, r.client, &deviceTypeRequest, data.Tags, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	utils.ApplyCustomFields(ctx, &deviceTypeRequest, data.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Call the API
	deviceType, httpResp, err := r.client.DcimAPI.DcimDeviceTypesCreate(ctx).WritableDeviceTypeRequest(deviceTypeRequest).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating device type",
			utils.FormatAPIError("create device type", err, httpResp),
		)
		return
	}
	if !utils.ValidateStatusCode(&resp.Diagnostics, "create device type", httpResp, http.StatusCreated) {
		return
	}
	tflog.Debug(ctx, "Created device type", map[string]interface{}{
		"id":    deviceType.GetId(),
		"model": deviceType.GetModel(),
	})

	// Map response to state
	r.mapDeviceTypeToState(deviceType, &data)
	if resp.Diagnostics.HasError() {
		return
	}

	// Apply filter-to-owned pattern for tags
	data.Tags = utils.PopulateTagsSlugFilteredToOwned(ctx, deviceType.HasTags(), deviceType.GetTags(), planTags)
	data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, planCustomFields, deviceType.GetCustomFields(), &resp.Diagnostics)
	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(data.ID.ValueString()), data.CustomFields, &resp.Diagnostics)
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
	deviceTypeIDInt, err := utils.ParseID(deviceTypeID)
	if err != nil {
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
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if utils.HandleNotFound(httpResp, func() {
			tflog.Debug(ctx, "Device type not found, removing from state", map[string]interface{}{
				"id": deviceTypeID,
			})
			resp.State.RemoveResource(ctx)
		}) {
			return
		}
		resp.Diagnostics.AddError(
			"Error reading device type",
			utils.FormatAPIError(fmt.Sprintf("read device type ID %s", deviceTypeID), err, httpResp),
		)
		return
	}
	if !utils.ValidateStatusCode(&resp.Diagnostics, "read device type", httpResp, http.StatusOK) {
		return
	}

	// Store state values for filter-to-owned pattern
	stateTags := data.Tags
	stateCustomFields := data.CustomFields

	// Map response to state
	r.mapDeviceTypeToState(deviceType, &data)
	if resp.Diagnostics.HasError() {
		return
	}

	// Apply filter-to-owned pattern for tags
	data.Tags = utils.PopulateTagsSlugFilteredToOwned(ctx, deviceType.HasTags(), deviceType.GetTags(), stateTags)
	data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, stateCustomFields, deviceType.GetCustomFields(), &resp.Diagnostics)
	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(data.ID.ValueString()), data.CustomFields, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DeviceTypeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data DeviceTypeResourceModel
	var state DeviceTypeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Parse the ID
	deviceTypeID := data.ID.ValueString()
	var deviceTypeIDInt int32
	deviceTypeIDInt, err := utils.ParseID(deviceTypeID)
	if err != nil {
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
	manufacturer, diags := netboxlookup.LookupManufacturer(ctx, r.client, data.Manufacturer.ValueString())
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
		platform, diags := netboxlookup.LookupPlatform(ctx, r.client, data.DefaultPlatform.ValueString())
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		deviceTypeRequest.SetDefaultPlatform(*platform)
	}

	if !data.PartNumber.IsNull() && !data.PartNumber.IsUnknown() {
		partNumber := data.PartNumber.ValueString()
		deviceTypeRequest.PartNumber = &partNumber
	} else if data.PartNumber.IsNull() {
		empty := ""
		deviceTypeRequest.PartNumber = &empty
	}

	if !data.UHeight.IsNull() && !data.UHeight.IsUnknown() {
		uHeight := data.UHeight.ValueFloat64()
		deviceTypeRequest.UHeight = &uHeight
	} else if data.UHeight.IsNull() {
		defaultUHeight := 1.0
		deviceTypeRequest.UHeight = &defaultUHeight
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
	// Note: subdevice_role has a NOT NULL constraint in Netbox DB and cannot be explicitly cleared

	if !data.Airflow.IsNull() && !data.Airflow.IsUnknown() && data.Airflow.ValueString() != "" {
		airflow := netbox.DeviceAirflowValue(data.Airflow.ValueString())
		deviceTypeRequest.Airflow = &airflow
	}
	// Note: airflow has a NOT NULL constraint in Netbox DB and cannot be explicitly cleared

	if !data.Weight.IsNull() && !data.Weight.IsUnknown() {
		weight := data.Weight.ValueFloat64()
		deviceTypeRequest.Weight = *netbox.NewNullableFloat64(&weight)
	} else if data.Weight.IsNull() {
		deviceTypeRequest.SetWeightNil()
	}

	if !data.WeightUnit.IsNull() && !data.WeightUnit.IsUnknown() && data.WeightUnit.ValueString() != "" {
		weightUnit := netbox.DeviceTypeWeightUnitValue(data.WeightUnit.ValueString())
		deviceTypeRequest.WeightUnit = &weightUnit
	}
	// Note: weight_unit has a DB NOT NULL constraint and cannot be cleared once set

	// Apply description and comments
	utils.ApplyDescription(&deviceTypeRequest, data.Description)
	utils.ApplyComments(&deviceTypeRequest, data.Comments)

	// Apply tags - merge-aware: use plan if provided, else use state
	if utils.IsSet(data.Tags) {
		utils.ApplyTagsFromSlugs(ctx, r.client, &deviceTypeRequest, data.Tags, &resp.Diagnostics)
	} else if utils.IsSet(state.Tags) {
		utils.ApplyTagsFromSlugs(ctx, r.client, &deviceTypeRequest, state.Tags, &resp.Diagnostics)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	// Apply custom fields with merge-aware logic
	utils.ApplyCustomFieldsWithMerge(ctx, &deviceTypeRequest, data.CustomFields, state.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Call the API
	deviceType, httpResp, err := r.client.DcimAPI.DcimDeviceTypesUpdate(ctx, deviceTypeIDInt).WritableDeviceTypeRequest(deviceTypeRequest).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating device type",
			utils.FormatAPIError(fmt.Sprintf("update device type ID %s", deviceTypeID), err, httpResp),
		)
		return
	}
	if !utils.ValidateStatusCode(&resp.Diagnostics, "update device type", httpResp, http.StatusOK) {
		return
	}
	tflog.Debug(ctx, "Updated device type", map[string]interface{}{
		"id":    deviceType.GetId(),
		"model": deviceType.GetModel(),
	})

	// Map response to state
	r.mapDeviceTypeToState(deviceType, &data)
	if resp.Diagnostics.HasError() {
		return
	}

	// Apply filter-to-owned pattern for tags
	data.Tags = utils.PopulateTagsSlugFilteredToOwned(ctx, deviceType.HasTags(), deviceType.GetTags(), data.Tags)
	data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, data.CustomFields, deviceType.GetCustomFields(), &resp.Diagnostics)
	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(data.ID.ValueString()), data.CustomFields, &resp.Diagnostics)
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
	deviceTypeIDInt, err := utils.ParseID(deviceTypeID)
	if err != nil {
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
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if utils.HandleNotFound(httpResp, func() {
			// Already deleted
			tflog.Debug(ctx, "Device type already deleted", map[string]interface{}{
				"id": deviceTypeID,
			})
		}) {
			return
		}
		resp.Diagnostics.AddError(
			"Error deleting device type",
			utils.FormatAPIError(fmt.Sprintf("delete device type ID %s", deviceTypeID), err, httpResp),
		)
		return
	}
	if !utils.ValidateStatusCode(&resp.Diagnostics, "delete device type", httpResp, http.StatusNoContent) {
		return
	}
	tflog.Debug(ctx, "Deleted device type", map[string]interface{}{
		"id": deviceTypeID,
	})
}

func (r *DeviceTypeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if parsed, ok := utils.ParseImportIdentityCustomFields(ctx, req.Identity, &resp.Diagnostics); ok {
		if resp.Diagnostics.HasError() {
			return
		}
		if parsed.ID == "" {
			resp.Diagnostics.AddError("Invalid import identity", "Identity id must be provided")
			return
		}

		deviceTypeIDInt, err := utils.ParseID(parsed.ID)
		if err != nil {
			resp.Diagnostics.AddError(
				"Invalid Device Type ID",
				fmt.Sprintf("Device Type ID must be a number, got: %s", parsed.ID),
			)
			return
		}

		deviceType, httpResp, err := r.client.DcimAPI.DcimDeviceTypesRetrieve(ctx, deviceTypeIDInt).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error importing device type",
				utils.FormatAPIError("read device type", err, httpResp),
			)
			return
		}
		if !utils.ValidateStatusCode(&resp.Diagnostics, "import device type", httpResp, http.StatusOK) {
			return
		}

		var data DeviceTypeResourceModel
		manufacturer := deviceType.GetManufacturer()
		data.Manufacturer = types.StringValue(fmt.Sprintf("%d", (&manufacturer).GetId()))
		if deviceType.HasDefaultPlatform() && deviceType.DefaultPlatform.Get() != nil {
			data.DefaultPlatform = types.StringValue(fmt.Sprintf("%d", deviceType.DefaultPlatform.Get().GetId()))
		}
		data.Tags = utils.PopulateTagsSlugFromAPI(ctx, deviceType.HasTags(), deviceType.GetTags(), data.Tags)
		if parsed.HasCustomFields {
			if len(parsed.CustomFields) == 0 {
				data.CustomFields = types.SetValueMust(utils.GetCustomFieldsAttributeType().ElemType, []attr.Value{})
			} else {
				ownedSet, setDiags := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, parsed.CustomFields)
				resp.Diagnostics.Append(setDiags...)
				if resp.Diagnostics.HasError() {
					return
				}
				data.CustomFields = ownedSet
			}
		} else {
			data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)
		}

		r.mapDeviceTypeToState(deviceType, &data)
		if resp.Diagnostics.HasError() {
			return
		}
		if parsed.HasCustomFields {
			data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, data.CustomFields, deviceType.GetCustomFields(), &resp.Diagnostics)
		} else {
			data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)
		}
		if resp.Diagnostics.HasError() {
			return
		}

		if resp.Identity != nil {
			listValue, listDiags := types.ListValueFrom(ctx, types.StringType, parsed.CustomFieldItems)
			resp.Diagnostics.Append(listDiags...)
			if resp.Diagnostics.HasError() {
				return
			}
			resp.Diagnostics.Append(resp.Identity.Set(ctx, &utils.ImportIdentityCustomFieldsModel{
				ID:           types.StringValue(parsed.ID),
				CustomFields: listValue,
			})...)
		}

		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
		return
	}

	utils.ImportStatePassthroughIDWithValidation(ctx, req, resp, path.Root("id"), true)
}

// mapDeviceTypeToState maps a DeviceType from the API to the Terraform state model.
func (r *DeviceTypeResource) mapDeviceTypeToState(deviceType *netbox.DeviceType, data *DeviceTypeResourceModel) {
	data.ID = types.StringValue(fmt.Sprintf("%d", deviceType.GetId()))
	data.Model = types.StringValue(deviceType.GetModel())
	data.Slug = types.StringValue(deviceType.GetSlug())

	// Handle manufacturer - preserve the original input value (slug, name, or ID)
	manufacturer := deviceType.Manufacturer
	data.Manufacturer = utils.UpdateReferenceAttribute(data.Manufacturer, manufacturer.GetName(), manufacturer.GetSlug(), manufacturer.GetId())

	// Handle default platform - preserve the original input value
	switch {
	case deviceType.HasDefaultPlatform() && deviceType.DefaultPlatform.Get() != nil:
		platform := deviceType.DefaultPlatform.Get()
		data.DefaultPlatform = utils.UpdateReferenceAttribute(data.DefaultPlatform, platform.GetName(), platform.GetSlug(), platform.GetId())

	case !data.DefaultPlatform.IsNull() && !data.DefaultPlatform.IsUnknown():
		// API says no platform but user had one configured - keep user's value
		// This shouldn't happen in normal operation

	default:
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
	if deviceType.HasSubdeviceRole() && deviceType.SubdeviceRole.Get() != nil && string(deviceType.SubdeviceRole.Get().GetValue()) != "" {
		data.SubdeviceRole = types.StringValue(string(deviceType.SubdeviceRole.Get().GetValue()))
	} else if !data.SubdeviceRole.IsNull() {
		data.SubdeviceRole = types.StringNull()
	}

	// Handle airflow
	if deviceType.HasAirflow() && deviceType.Airflow.Get() != nil && string(deviceType.Airflow.Get().GetValue()) != "" {
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
	if deviceType.HasWeightUnit() && deviceType.WeightUnit.Get() != nil && string(deviceType.WeightUnit.Get().GetValue()) != "" {
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

	// Tags and custom fields are handled in Create/Read/Update methods using filter-to-owned pattern.
}
