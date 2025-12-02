// Package resources contains Terraform resource implementations for the Netbox provider.
package resources

import (
	"context"
	"fmt"

	"github.com/bab3l/go-netbox"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/bab3l/terraform-provider-netbox/internal/validators"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &DeviceRoleResource{}
var _ resource.ResourceWithImportState = &DeviceRoleResource{}

func NewDeviceRoleResource() resource.Resource {
	return &DeviceRoleResource{}
}

// DeviceRoleResource defines the resource implementation.
type DeviceRoleResource struct {
	client *netbox.APIClient
}

// DeviceRoleResourceModel describes the resource data model.
type DeviceRoleResourceModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Slug         types.String `tfsdk:"slug"`
	Color        types.String `tfsdk:"color"`
	VMRole       types.Bool   `tfsdk:"vm_role"`
	Description  types.String `tfsdk:"description"`
	Tags         types.Set    `tfsdk:"tags"`
	CustomFields types.Set    `tfsdk:"custom_fields"`
}

func (r *DeviceRoleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_device_role"
}

func (r *DeviceRoleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a device role in Netbox. Device roles are used to categorize devices by their function within the network infrastructure (e.g., 'Router', 'Switch', 'Server', 'Firewall').",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Unique identifier for the device role (assigned by Netbox).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Full name of the device role. This is the human-readable display name.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 100),
				},
			},
			"slug": schema.StringAttribute{
				MarkdownDescription: "URL-friendly identifier for the device role. Must be unique and contain only alphanumeric characters, hyphens, and underscores.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 100),
					validators.ValidSlug(),
				},
			},
			"color": schema.StringAttribute{
				MarkdownDescription: "Color code for the device role in hexadecimal format (e.g., 'aa1409' for red). Used for visual identification in the Netbox UI. If not specified, Netbox assigns a default color.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(6, 6),
					stringvalidator.RegexMatches(
						validators.HexColorRegex(),
						"must be a valid 6-character hexadecimal color code (e.g., 'aa1409')",
					),
				},
			},
			"vm_role": schema.BoolAttribute{
				MarkdownDescription: "Whether virtual machines may be assigned to this role. Set to true to allow VMs to use this role, false otherwise. Defaults to true.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Detailed description of the device role, its purpose, or other relevant information.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(200),
				},
			},
			"tags": schema.SetNestedAttribute{
				MarkdownDescription: "Tags assigned to this device role. Tags provide a way to categorize and organize resources.",
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
				MarkdownDescription: "Custom fields assigned to this device role. Custom fields allow you to store additional structured data.",
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
							MarkdownDescription: "Type of the custom field (text, longtext, integer, boolean, date, url, json, select, multiselect, object, multiobject).",
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

func (r *DeviceRoleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// mapDeviceRoleToState maps a DeviceRole from the API to the Terraform state model.
func (r *DeviceRoleResource) mapDeviceRoleToState(ctx context.Context, deviceRole *netbox.DeviceRole, data *DeviceRoleResourceModel, diags *diag.Diagnostics) {
	data.ID = types.StringValue(fmt.Sprintf("%d", deviceRole.GetId()))
	data.Name = types.StringValue(deviceRole.GetName())
	data.Slug = types.StringValue(deviceRole.GetSlug())

	// Handle color - use value from API if available
	if deviceRole.HasColor() && deviceRole.GetColor() != "" {
		data.Color = types.StringValue(deviceRole.GetColor())
	} else if !data.Color.IsNull() {
		// Preserve null if originally null and API returns empty
		data.Color = types.StringNull()
	}

	// Handle vm_role
	if deviceRole.HasVmRole() {
		data.VMRole = types.BoolValue(deviceRole.GetVmRole())
	} else {
		data.VMRole = types.BoolValue(true) // Default to true per Netbox API
	}

	// Handle description
	if deviceRole.HasDescription() && deviceRole.GetDescription() != "" {
		data.Description = types.StringValue(deviceRole.GetDescription())
	} else if !data.Description.IsNull() {
		// Preserve null if originally null and API returns empty
		data.Description = types.StringNull()
	}

	// Handle tags
	if deviceRole.HasTags() {
		tags := utils.NestedTagsToTagModels(deviceRole.GetTags())
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
	if deviceRole.HasCustomFields() && !data.CustomFields.IsNull() {
		var stateCustomFields []utils.CustomFieldModel
		cfDiags := data.CustomFields.ElementsAs(ctx, &stateCustomFields, false)
		diags.Append(cfDiags...)
		if diags.HasError() {
			return
		}

		customFields := utils.MapToCustomFieldModels(deviceRole.GetCustomFields(), stateCustomFields)
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

func (r *DeviceRoleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data DeviceRoleResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating device role", map[string]interface{}{
		"name": data.Name.ValueString(),
		"slug": data.Slug.ValueString(),
	})

	// Build the device role request
	deviceRoleRequest := netbox.DeviceRoleRequest{
		Name: data.Name.ValueString(),
		Slug: data.Slug.ValueString(),
	}

	// Set optional fields if provided
	if !data.Color.IsNull() && !data.Color.IsUnknown() {
		color := data.Color.ValueString()
		deviceRoleRequest.Color = &color
	}

	if !data.VMRole.IsNull() && !data.VMRole.IsUnknown() {
		vmRole := data.VMRole.ValueBool()
		deviceRoleRequest.VmRole = &vmRole
	}

	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		description := data.Description.ValueString()
		deviceRoleRequest.Description = &description
	}

	// Handle tags
	if !data.Tags.IsNull() && !data.Tags.IsUnknown() {
		var tags []utils.TagModel
		resp.Diagnostics.Append(data.Tags.ElementsAs(ctx, &tags, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		deviceRoleRequest.Tags = utils.TagsToNestedTagRequests(tags)
	}

	// Handle custom fields
	if !data.CustomFields.IsNull() && !data.CustomFields.IsUnknown() {
		var customFields []utils.CustomFieldModel
		resp.Diagnostics.Append(data.CustomFields.ElementsAs(ctx, &customFields, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		deviceRoleRequest.CustomFields = utils.CustomFieldsToMap(customFields)
	}

	// Call the API
	deviceRole, httpResp, err := r.client.DcimAPI.DcimDeviceRolesCreate(ctx).DeviceRoleRequest(deviceRoleRequest).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating device role",
			utils.FormatAPIError("create device role", err, httpResp),
		)
		return
	}

	tflog.Debug(ctx, "Created device role", map[string]interface{}{
		"id":   deviceRole.GetId(),
		"name": deviceRole.GetName(),
	})

	// Map response to state
	r.mapDeviceRoleToState(ctx, deviceRole, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DeviceRoleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data DeviceRoleResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Parse the ID
	deviceRoleID := data.ID.ValueString()
	var deviceRoleIDInt int32
	if _, err := fmt.Sscanf(deviceRoleID, "%d", &deviceRoleIDInt); err != nil {
		resp.Diagnostics.AddError(
			"Invalid Device Role ID",
			fmt.Sprintf("Device Role ID must be a number, got: %s", deviceRoleID),
		)
		return
	}

	tflog.Debug(ctx, "Reading device role", map[string]interface{}{
		"id": deviceRoleID,
	})

	// Call the API
	deviceRole, httpResp, err := r.client.DcimAPI.DcimDeviceRolesRetrieve(ctx, deviceRoleIDInt).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			tflog.Debug(ctx, "Device role not found, removing from state", map[string]interface{}{
				"id": deviceRoleID,
			})
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading device role",
			utils.FormatAPIError(fmt.Sprintf("read device role ID %s", deviceRoleID), err, httpResp),
		)
		return
	}

	// Map response to state
	r.mapDeviceRoleToState(ctx, deviceRole, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DeviceRoleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data DeviceRoleResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Parse the ID
	deviceRoleID := data.ID.ValueString()
	var deviceRoleIDInt int32
	if _, err := fmt.Sscanf(deviceRoleID, "%d", &deviceRoleIDInt); err != nil {
		resp.Diagnostics.AddError(
			"Invalid Device Role ID",
			fmt.Sprintf("Device Role ID must be a number, got: %s", deviceRoleID),
		)
		return
	}

	tflog.Debug(ctx, "Updating device role", map[string]interface{}{
		"id":   deviceRoleID,
		"name": data.Name.ValueString(),
		"slug": data.Slug.ValueString(),
	})

	// Build the device role request
	deviceRoleRequest := netbox.DeviceRoleRequest{
		Name: data.Name.ValueString(),
		Slug: data.Slug.ValueString(),
	}

	// Set optional fields if provided
	if !data.Color.IsNull() && !data.Color.IsUnknown() {
		color := data.Color.ValueString()
		deviceRoleRequest.Color = &color
	}

	if !data.VMRole.IsNull() && !data.VMRole.IsUnknown() {
		vmRole := data.VMRole.ValueBool()
		deviceRoleRequest.VmRole = &vmRole
	}

	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		description := data.Description.ValueString()
		deviceRoleRequest.Description = &description
	}

	// Handle tags
	if !data.Tags.IsNull() && !data.Tags.IsUnknown() {
		var tags []utils.TagModel
		resp.Diagnostics.Append(data.Tags.ElementsAs(ctx, &tags, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		deviceRoleRequest.Tags = utils.TagsToNestedTagRequests(tags)
	}

	// Handle custom fields
	if !data.CustomFields.IsNull() && !data.CustomFields.IsUnknown() {
		var customFields []utils.CustomFieldModel
		resp.Diagnostics.Append(data.CustomFields.ElementsAs(ctx, &customFields, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		deviceRoleRequest.CustomFields = utils.CustomFieldsToMap(customFields)
	}

	// Call the API
	deviceRole, httpResp, err := r.client.DcimAPI.DcimDeviceRolesUpdate(ctx, deviceRoleIDInt).DeviceRoleRequest(deviceRoleRequest).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating device role",
			utils.FormatAPIError(fmt.Sprintf("update device role ID %s", deviceRoleID), err, httpResp),
		)
		return
	}

	tflog.Debug(ctx, "Updated device role", map[string]interface{}{
		"id":   deviceRole.GetId(),
		"name": deviceRole.GetName(),
	})

	// Map response to state
	r.mapDeviceRoleToState(ctx, deviceRole, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DeviceRoleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data DeviceRoleResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Parse the ID
	deviceRoleID := data.ID.ValueString()
	var deviceRoleIDInt int32
	if _, err := fmt.Sscanf(deviceRoleID, "%d", &deviceRoleIDInt); err != nil {
		resp.Diagnostics.AddError(
			"Invalid Device Role ID",
			fmt.Sprintf("Device Role ID must be a number, got: %s", deviceRoleID),
		)
		return
	}

	tflog.Debug(ctx, "Deleting device role", map[string]interface{}{
		"id": deviceRoleID,
	})

	// Call the API
	httpResp, err := r.client.DcimAPI.DcimDeviceRolesDestroy(ctx, deviceRoleIDInt).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			// Already deleted, consider success
			tflog.Debug(ctx, "Device role already deleted", map[string]interface{}{
				"id": deviceRoleID,
			})
			return
		}
		resp.Diagnostics.AddError(
			"Error deleting device role",
			utils.FormatAPIError(fmt.Sprintf("delete device role ID %s", deviceRoleID), err, httpResp),
		)
		return
	}

	tflog.Debug(ctx, "Deleted device role", map[string]interface{}{
		"id": deviceRoleID,
	})
}

func (r *DeviceRoleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
