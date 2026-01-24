// Package resources contains Terraform resource implementations for the Netbox provider.

package resources

import (
	"context"
	"fmt"
	"maps"
	"net/http"

	"github.com/bab3l/go-netbox"
	"github.com/bab3l/terraform-provider-netbox/internal/netboxlookup"
	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &DeviceRoleResource{}
var _ resource.ResourceWithImportState = &DeviceRoleResource{}
var _ resource.ResourceWithIdentity = &DeviceRoleResource{}

func NewDeviceRoleResource() resource.Resource {
	return &DeviceRoleResource{}
}

// DeviceRoleResource defines the resource implementation.
type DeviceRoleResource struct {
	client *netbox.APIClient
}

// DeviceRoleResourceModel describes the resource data model.
type DeviceRoleResourceModel struct {
	ID             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	Slug           types.String `tfsdk:"slug"`
	Color          types.String `tfsdk:"color"`
	VMRole         types.Bool   `tfsdk:"vm_role"`
	ConfigTemplate types.String `tfsdk:"config_template"`
	Description    types.String `tfsdk:"description"`
	Tags           types.Set    `tfsdk:"tags"`
	CustomFields   types.Set    `tfsdk:"custom_fields"`
}

func (r *DeviceRoleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_device_role"
}

func (r *DeviceRoleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a device role in Netbox. Device roles are used to categorize devices by their function within the network infrastructure (e.g., 'Router', 'Switch', 'Server', 'Firewall').",
		Attributes: map[string]schema.Attribute{
			"id":              nbschema.IDAttribute("device role"),
			"name":            nbschema.NameAttribute("device role", 100),
			"slug":            nbschema.SlugAttribute("device role"),
			"color":           nbschema.ComputedColorAttribute("device role"),
			"vm_role":         nbschema.BoolAttributeWithDefault("Whether virtual machines may be assigned to this role. Set to true to allow VMs to use this role, false otherwise. Defaults to true.", true),
			"config_template": nbschema.ReferenceAttributeWithDiffSuppress("config template", "ID or name of the config template assigned to this device role."),
			"tags":            nbschema.TagsSlugAttribute(),
			"custom_fields":   nbschema.CustomFieldsAttribute(),
		},
	}

	// Add description attribute
	maps.Copy(resp.Schema.Attributes, nbschema.DescriptionOnlyAttributes("device role"))

	// Tags and custom fields are defined directly in the schema above.
}

func (r *DeviceRoleResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = nbschema.ImportIdentityWithCustomFieldsSchema()
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

	// Handle config template - preserve the original input value
	switch {
	case deviceRole.ConfigTemplate.IsSet() && deviceRole.ConfigTemplate.Get() != nil:
		templateObj := deviceRole.ConfigTemplate.Get()
		data.ConfigTemplate = utils.UpdateReferenceAttribute(data.ConfigTemplate, templateObj.GetName(), "", templateObj.GetId())

	case !data.ConfigTemplate.IsNull() && !data.ConfigTemplate.IsUnknown():
		// User had a value but API says null

	default:
		data.ConfigTemplate = types.StringNull()
	}

	// Handle description
	if deviceRole.HasDescription() && deviceRole.GetDescription() != "" {
		data.Description = types.StringValue(deviceRole.GetDescription())
	} else if !data.Description.IsNull() {
		// Preserve null if originally null and API returns empty
		data.Description = types.StringNull()
	}

	// Tags and custom fields are handled in Create/Read/Update methods using filter-to-owned pattern.
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
	if !data.ConfigTemplate.IsNull() && !data.ConfigTemplate.IsUnknown() {
		configTemplate, diags := netboxlookup.LookupConfigTemplate(ctx, r.client, data.ConfigTemplate.ValueString())
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		deviceRoleRequest.SetConfigTemplate(*configTemplate)
	}

	// Apply description
	utils.ApplyDescription(&deviceRoleRequest, data.Description)

	// Store plan values for filter-to-owned pattern
	planTags := data.Tags
	planCustomFields := data.CustomFields

	// Handle tags
	utils.ApplyTagsFromSlugs(ctx, r.client, &deviceRoleRequest, data.Tags, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Handle custom fields (no merge needed for Create)
	utils.ApplyCustomFields(ctx, &deviceRoleRequest, data.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Call the API
	deviceRole, httpResp, err := r.client.DcimAPI.DcimDeviceRolesCreate(ctx).DeviceRoleRequest(deviceRoleRequest).Execute()
	defer utils.CloseResponseBody(httpResp)
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

	// Apply filter-to-owned pattern for tags
	wasExplicitlyEmpty := !planTags.IsNull() && !planTags.IsUnknown() && len(planTags.Elements()) == 0
	switch {
	case deviceRole.HasTags() && len(deviceRole.GetTags()) > 0:
		tagSlugs := make([]string, 0, len(deviceRole.GetTags()))
		for _, tag := range deviceRole.GetTags() {
			tagSlugs = append(tagSlugs, tag.GetSlug())
		}
		data.Tags = utils.TagsSlugToSet(ctx, tagSlugs)
	case wasExplicitlyEmpty:
		data.Tags = types.SetValueMust(types.StringType, []attr.Value{})
	default:
		data.Tags = types.SetNull(types.StringType)
	}
	data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, planCustomFields, deviceRole.GetCustomFields(), &resp.Diagnostics)
	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(data.ID.ValueString()), data.CustomFields, &resp.Diagnostics)
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
	deviceRoleIDInt, err := utils.ParseID(deviceRoleID)
	if err != nil {
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
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
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

	// Save original tags/custom_fields state before mapping
	originalTags := data.Tags
	originalCustomFields := data.CustomFields

	// Map response to state
	r.mapDeviceRoleToState(ctx, deviceRole, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Apply filter-to-owned pattern for tags
	wasExplicitlyEmpty := !originalTags.IsNull() && !originalTags.IsUnknown() && len(originalTags.Elements()) == 0
	switch {
	case deviceRole.HasTags() && len(deviceRole.GetTags()) > 0:
		tagSlugs := make([]string, 0, len(deviceRole.GetTags()))
		for _, tag := range deviceRole.GetTags() {
			tagSlugs = append(tagSlugs, tag.GetSlug())
		}
		data.Tags = utils.TagsSlugToSet(ctx, tagSlugs)
	case wasExplicitlyEmpty:
		data.Tags = types.SetValueMust(types.StringType, []attr.Value{})
	default:
		data.Tags = types.SetNull(types.StringType)
	}
	data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, originalCustomFields, deviceRole.GetCustomFields(), &resp.Diagnostics)

	// Preserve original custom_fields state if it was null or empty
	// This prevents unmanaged/cleared fields from reappearing in state
	if originalCustomFields.IsNull() || (utils.IsSet(originalCustomFields) && len(originalCustomFields.Elements()) == 0) {
		data.CustomFields = originalCustomFields
	}

	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(data.ID.ValueString()), data.CustomFields, &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DeviceRoleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan DeviceRoleResourceModel
	var state DeviceRoleResourceModel

	// Read both plan and state for merge-aware custom fields handling
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use plan as the data source
	data := plan

	// Parse the ID
	deviceRoleID := data.ID.ValueString()
	var deviceRoleIDInt int32
	deviceRoleIDInt, err := utils.ParseID(deviceRoleID)
	if err != nil {
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
	if !data.ConfigTemplate.IsNull() && !data.ConfigTemplate.IsUnknown() {
		configTemplate, diags := netboxlookup.LookupConfigTemplate(ctx, r.client, data.ConfigTemplate.ValueString())
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		deviceRoleRequest.SetConfigTemplate(*configTemplate)
	} else if data.ConfigTemplate.IsNull() {
		deviceRoleRequest.SetConfigTemplateNil()
	}

	// Apply description
	utils.ApplyDescription(&deviceRoleRequest, data.Description)

	// Handle tags - merge-aware: use plan if provided, else use state
	if utils.IsSet(data.Tags) {
		utils.ApplyTagsFromSlugs(ctx, r.client, &deviceRoleRequest, data.Tags, &resp.Diagnostics)
	} else if utils.IsSet(state.Tags) {
		utils.ApplyTagsFromSlugs(ctx, r.client, &deviceRoleRequest, state.Tags, &resp.Diagnostics)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	// Handle custom fields with merge-aware logic
	utils.ApplyCustomFieldsWithMerge(ctx, &deviceRoleRequest, plan.CustomFields, state.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Call the API
	deviceRole, httpResp, err := r.client.DcimAPI.DcimDeviceRolesUpdate(ctx, deviceRoleIDInt).DeviceRoleRequest(deviceRoleRequest).Execute()
	defer utils.CloseResponseBody(httpResp)
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

	// Apply filter-to-owned pattern for tags
	wasExplicitlyEmpty := !plan.Tags.IsNull() && !plan.Tags.IsUnknown() && len(plan.Tags.Elements()) == 0
	switch {
	case deviceRole.HasTags() && len(deviceRole.GetTags()) > 0:
		tagSlugs := make([]string, 0, len(deviceRole.GetTags()))
		for _, tag := range deviceRole.GetTags() {
			tagSlugs = append(tagSlugs, tag.GetSlug())
		}
		data.Tags = utils.TagsSlugToSet(ctx, tagSlugs)
	case wasExplicitlyEmpty:
		data.Tags = types.SetValueMust(types.StringType, []attr.Value{})
	default:
		data.Tags = types.SetNull(types.StringType)
	}
	data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, plan.CustomFields, deviceRole.GetCustomFields(), &resp.Diagnostics)
	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(data.ID.ValueString()), data.CustomFields, &resp.Diagnostics)
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
	deviceRoleIDInt, err := utils.ParseID(deviceRoleID)
	if err != nil {
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
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
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
	if parsed, ok := utils.ParseImportIdentityCustomFields(ctx, req.Identity, &resp.Diagnostics); ok {
		if resp.Diagnostics.HasError() {
			return
		}
		if parsed.ID == "" {
			resp.Diagnostics.AddError("Invalid import identity", "Identity id must be provided")
			return
		}

		deviceRoleIDInt, err := utils.ParseID(parsed.ID)
		if err != nil {
			resp.Diagnostics.AddError(
				"Invalid Device Role ID",
				fmt.Sprintf("Device Role ID must be a number, got: %s", parsed.ID),
			)
			return
		}

		deviceRole, httpResp, err := r.client.DcimAPI.DcimDeviceRolesRetrieve(ctx, deviceRoleIDInt).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error importing device role",
				utils.FormatAPIError(fmt.Sprintf("read device role ID %s", parsed.ID), err, httpResp),
			)
			return
		}

		var data DeviceRoleResourceModel
		r.mapDeviceRoleToState(ctx, deviceRole, &data, &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}

		data.Tags = types.SetNull(types.StringType)
		if parsed.HasCustomFields {
			if len(parsed.CustomFields) == 0 {
				data.CustomFields = types.SetValueMust(utils.GetCustomFieldsAttributeType().ElemType, []attr.Value{})
			} else {
				ownedSet, setDiags := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, parsed.CustomFields)
				resp.Diagnostics.Append(setDiags...)
				if resp.Diagnostics.HasError() {
					return
				}
				data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, ownedSet, deviceRole.GetCustomFields(), &resp.Diagnostics)
			}
		} else {
			data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)
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

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
