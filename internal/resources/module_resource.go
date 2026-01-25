// Package resources provides Terraform resource implementations for NetBox objects.

package resources

import (
	"context"
	"fmt"
	"maps"
	"net/http"

	"github.com/bab3l/go-netbox"
	lookup "github.com/bab3l/terraform-provider-netbox/internal/netboxlookup"
	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &ModuleResource{}
	_ resource.ResourceWithConfigure   = &ModuleResource{}
	_ resource.ResourceWithImportState = &ModuleResource{}
	_ resource.ResourceWithIdentity    = &ModuleResource{}
)

// NewModuleResource returns a new resource implementing the module resource.
func NewModuleResource() resource.Resource {
	return &ModuleResource{}
}

// ModuleResource defines the resource implementation.
type ModuleResource struct {
	client *netbox.APIClient
}

// ModuleResourceModel describes the resource data model.
type ModuleResourceModel struct {
	ID           types.String `tfsdk:"id"`
	Device       types.String `tfsdk:"device"`
	ModuleBay    types.Int32  `tfsdk:"module_bay"`
	ModuleType   types.String `tfsdk:"module_type"`
	Status       types.String `tfsdk:"status"`
	Serial       types.String `tfsdk:"serial"`
	AssetTag     types.String `tfsdk:"asset_tag"`
	Description  types.String `tfsdk:"description"`
	Comments     types.String `tfsdk:"comments"`
	Tags         types.Set    `tfsdk:"tags"`
	CustomFields types.Set    `tfsdk:"custom_fields"`
}

// Metadata returns the resource type name.
func (r *ModuleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_module"
}

// Schema defines the schema for the resource.
func (r *ModuleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a module in NetBox. Modules are hardware components installed in module bays within devices. This is the recommended replacement for inventory items (deprecated in NetBox v4.3).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique numeric ID of the module.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"device": schema.StringAttribute{
				MarkdownDescription: "The device this module is installed in (ID or name).",
				Required:            true,
			},
			"module_bay": schema.Int32Attribute{
				MarkdownDescription: "The module bay ID where this module is installed.",
				Required:            true,
			},
			"module_type": schema.StringAttribute{
				MarkdownDescription: "The module type (ID or model name).",
				Required:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "Operational status. Valid values: `offline`, `active`, `planned`, `staged`, `failed`, `decommissioning`.",
				Optional:            true,
				Computed:            true,
			},
			"serial": schema.StringAttribute{
				MarkdownDescription: "Serial number of the module.",
				Optional:            true,
			},
			"asset_tag": schema.StringAttribute{
				MarkdownDescription: "A unique tag used to identify this module.",
				Optional:            true,
			},
		},
	}

	// Add common descriptive attributes (description, comments)
	maps.Copy(resp.Schema.Attributes, nbschema.CommonDescriptiveAttributes("module"))

	// Add tags and custom fields
	resp.Schema.Attributes["tags"] = nbschema.TagsSlugAttribute()
	resp.Schema.Attributes["custom_fields"] = nbschema.CustomFieldsAttribute()
}

func (r *ModuleResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = nbschema.ImportIdentityWithCustomFieldsSchema()
}

func (r *ModuleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create creates the resource.
func (r *ModuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ModuleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Lookup device
	device, diags := lookup.LookupDevice(ctx, r.client, data.Device.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Lookup module type
	moduleType, diags := lookup.LookupModuleType(ctx, r.client, data.ModuleType.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build request
	apiReq := netbox.NewWritableModuleRequest(*device, data.ModuleBay.ValueInt32(), *moduleType)

	// Set optional fields
	if !data.Status.IsNull() && !data.Status.IsUnknown() {
		status := netbox.ModuleStatusValue(data.Status.ValueString())
		apiReq.SetStatus(status)
	}

	if !data.Serial.IsNull() && !data.Serial.IsUnknown() {
		apiReq.SetSerial(data.Serial.ValueString())
	}

	if !data.AssetTag.IsNull() && !data.AssetTag.IsUnknown() {
		apiReq.SetAssetTag(data.AssetTag.ValueString())
	}

	// Set common fields (description, comments, tags, custom_fields)
	utils.ApplyDescription(apiReq, data.Description)
	utils.ApplyComments(apiReq, data.Comments)
	utils.ApplyTagsFromSlugs(ctx, r.client, apiReq, data.Tags, &resp.Diagnostics)
	utils.ApplyCustomFields(ctx, apiReq, data.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Creating module", map[string]interface{}{
		"device":      data.Device.ValueString(),
		"module_bay":  data.ModuleBay.ValueInt32(),
		"module_type": data.ModuleType.ValueString(),
	})
	response, httpResp, err := r.client.DcimAPI.DcimModulesCreate(ctx).WritableModuleRequest(*apiReq).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating module",
			utils.FormatAPIError("create module", err, httpResp),
		)
		return
	}

	// Map response to model
	r.mapResponseToModel(ctx, response, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(data.ID.ValueString()), data.CustomFields, &resp.Diagnostics)
	tflog.Trace(ctx, "Created module", map[string]interface{}{
		"id": data.ID.ValueString(),
	})
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read refreshes the resource state.
func (r *ModuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ModuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	moduleID, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Module ID",
			fmt.Sprintf("Module ID must be a number, got: %s", data.ID.ValueString()),
		)
		return
	}
	tflog.Debug(ctx, "Reading module", map[string]interface{}{
		"id": moduleID,
	})
	response, httpResp, err := r.client.DcimAPI.DcimModulesRetrieve(ctx, moduleID).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading module",
			utils.FormatAPIError(fmt.Sprintf("read module ID %d", moduleID), err, httpResp),
		)
		return
	}

	// Map response to model
	r.mapResponseToModel(ctx, response, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(data.ID.ValueString()), data.CustomFields, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates the resource.
func (r *ModuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, plan ModuleResourceModel

	// Read both state and plan for merge-aware custom fields handling
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	moduleID, err := utils.ParseID(plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Module ID",
			fmt.Sprintf("Module ID must be a number, got: %s", plan.ID.ValueString()),
		)
		return
	}

	// Lookup device
	device, diags := lookup.LookupDevice(ctx, r.client, plan.Device.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Lookup module type
	moduleType, diags := lookup.LookupModuleType(ctx, r.client, plan.ModuleType.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build request
	apiReq := netbox.NewWritableModuleRequest(*device, plan.ModuleBay.ValueInt32(), *moduleType)

	// Set optional fields
	if !plan.Status.IsNull() && !plan.Status.IsUnknown() {
		status := netbox.ModuleStatusValue(plan.Status.ValueString())
		apiReq.SetStatus(status)
	} else if plan.Status.IsNull() {
		// Explicitly clear status by setting empty string
		emptyStatus := netbox.ModuleStatusValue("")
		apiReq.SetStatus(emptyStatus)
	}

	if !plan.Serial.IsNull() && !plan.Serial.IsUnknown() {
		apiReq.SetSerial(plan.Serial.ValueString())
	} else if plan.Serial.IsNull() {
		// Explicitly clear serial
		apiReq.SetSerial("")
	}

	if !plan.AssetTag.IsNull() && !plan.AssetTag.IsUnknown() {
		apiReq.SetAssetTag(plan.AssetTag.ValueString())
	} else if plan.AssetTag.IsNull() {
		// Explicitly clear asset_tag
		apiReq.SetAssetTag("")
	}

	// Set common fields (description, comments, tags, custom_fields)
	utils.ApplyDescription(apiReq, plan.Description)
	utils.ApplyComments(apiReq, plan.Comments)
	utils.ApplyTagsFromSlugs(ctx, r.client, apiReq, plan.Tags, &resp.Diagnostics)
	// Apply custom fields with merge logic to preserve unmanaged fields
	utils.ApplyCustomFieldsWithMerge(ctx, apiReq, plan.CustomFields, state.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Updating module", map[string]interface{}{
		"id": moduleID,
	})
	response, httpResp, err := r.client.DcimAPI.DcimModulesUpdate(ctx, moduleID).WritableModuleRequest(*apiReq).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating module",
			utils.FormatAPIError(fmt.Sprintf("update module ID %d", moduleID), err, httpResp),
		)
		return
	}

	// Save the plan's custom fields before mapping (for filter-to-owned pattern)
	planCustomFields := plan.CustomFields

	// Map response to model
	r.mapResponseToModel(ctx, response, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Apply filter-to-owned pattern for custom fields
	plan.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, planCustomFields, response.GetCustomFields(), &resp.Diagnostics)
	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(plan.ID.ValueString()), plan.CustomFields, &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete deletes the resource.
func (r *ModuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ModuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	moduleID, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Module ID",
			fmt.Sprintf("Module ID must be a number, got: %s", data.ID.ValueString()),
		)
		return
	}
	tflog.Debug(ctx, "Deleting module", map[string]interface{}{
		"id": moduleID,
	})
	httpResp, err := r.client.DcimAPI.DcimModulesDestroy(ctx, moduleID).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			return
		}
		resp.Diagnostics.AddError(
			"Error deleting module",
			utils.FormatAPIError(fmt.Sprintf("delete module ID %d", moduleID), err, httpResp),
		)
		return
	}
}

// ImportState imports an existing resource.
func (r *ModuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if parsed, ok := utils.ParseImportIdentityCustomFields(ctx, req.Identity, &resp.Diagnostics); ok {
		if resp.Diagnostics.HasError() {
			return
		}
		if parsed.ID == "" {
			resp.Diagnostics.AddError("Invalid import identity", "Identity id must be provided")
			return
		}

		moduleID, err := utils.ParseID(parsed.ID)
		if err != nil {
			resp.Diagnostics.AddError(
				"Invalid Import ID",
				fmt.Sprintf("Module ID must be a number, got: %s", parsed.ID),
			)
			return
		}
		response, httpResp, err := r.client.DcimAPI.DcimModulesRetrieve(ctx, moduleID).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error importing module",
				utils.FormatAPIError(fmt.Sprintf("import module ID %d", moduleID), err, httpResp),
			)
			return
		}
		var data ModuleResourceModel
		if device := response.GetDevice(); device.Id != 0 {
			data.Device = types.StringValue(fmt.Sprintf("%d", device.GetId()))
		}
		if moduleBay := response.GetModuleBay(); moduleBay.Id != 0 {
			data.ModuleBay = types.Int32Value(moduleBay.Id)
		}
		if moduleType := response.GetModuleType(); moduleType.Id != 0 {
			data.ModuleType = types.StringValue(fmt.Sprintf("%d", moduleType.GetId()))
		}
		data.Tags = utils.PopulateTagsSlugFromAPI(ctx, response.HasTags(), response.GetTags(), data.Tags)
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

		r.mapResponseToModel(ctx, response, &data, &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}
		if parsed.HasCustomFields {
			data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, data.CustomFields, response.GetCustomFields(), &resp.Diagnostics)
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

	moduleID, err := utils.ParseID(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Module ID must be a number, got: %s", req.ID),
		)
		return
	}
	response, httpResp, err := r.client.DcimAPI.DcimModulesRetrieve(ctx, moduleID).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error importing module",
			utils.FormatAPIError(fmt.Sprintf("import module ID %d", moduleID), err, httpResp),
		)
		return
	}
	var data ModuleResourceModel
	r.mapResponseToModel(ctx, response, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// mapResponseToModel maps the API response to the Terraform model.
func (r *ModuleResource) mapResponseToModel(ctx context.Context, module *netbox.Module, data *ModuleResourceModel, diags *diag.Diagnostics) {
	data.ID = types.StringValue(fmt.Sprintf("%d", module.GetId()))

	// Map device - preserve user's input format
	if device := module.GetDevice(); device.Id != 0 {
		data.Device = utils.UpdateReferenceAttribute(data.Device, device.GetName(), "", device.GetId())
	}

	// Map module_bay
	moduleBay := module.GetModuleBay()
	data.ModuleBay = types.Int32Value(moduleBay.Id)

	// Map module_type - preserve user's input format
	if mt := module.GetModuleType(); mt.Id != 0 {
		data.ModuleType = utils.UpdateReferenceAttribute(data.ModuleType, mt.GetModel(), "", mt.GetId())
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

	// Populate tags - filter to only those managed in config
	data.Tags = utils.PopulateTagsSlugFilteredToOwned(ctx, module.HasTags(), module.GetTags(), data.Tags)

	// Filter custom fields to only those managed in config
	data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, data.CustomFields, module.GetCustomFields(), diags)
}
