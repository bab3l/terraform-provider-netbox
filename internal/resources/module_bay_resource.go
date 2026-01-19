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
	_ resource.Resource                = &ModuleBayResource{}
	_ resource.ResourceWithConfigure   = &ModuleBayResource{}
	_ resource.ResourceWithImportState = &ModuleBayResource{}
	_ resource.ResourceWithIdentity    = &ModuleBayResource{}
)

// NewModuleBayResource returns a new resource implementing the module bay resource.
func NewModuleBayResource() resource.Resource {
	return &ModuleBayResource{}
}

// ModuleBayResource defines the resource implementation.
type ModuleBayResource struct {
	client *netbox.APIClient
}

// ModuleBayResourceModel describes the resource data model.
type ModuleBayResourceModel struct {
	ID           types.String `tfsdk:"id"`
	Device       types.String `tfsdk:"device"`
	Name         types.String `tfsdk:"name"`
	Label        types.String `tfsdk:"label"`
	Position     types.String `tfsdk:"position"`
	Description  types.String `tfsdk:"description"`
	Tags         types.Set    `tfsdk:"tags"`
	CustomFields types.Set    `tfsdk:"custom_fields"`
}

// Metadata returns the resource type name.
func (r *ModuleBayResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_module_bay"
}

// Schema defines the schema for the resource.
func (r *ModuleBayResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a module bay in NetBox. Module bays are slots within devices that can accept modules.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique numeric ID of the module bay.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"device": schema.StringAttribute{
				MarkdownDescription: "The device this module bay belongs to (ID or name).",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the module bay.",
				Required:            true,
			},
			"label": schema.StringAttribute{
				MarkdownDescription: "Physical label of the module bay.",
				Optional:            true,
			},
			"position": schema.StringAttribute{
				MarkdownDescription: "Identifier to reference when renaming installed components.",
				Optional:            true,
			},
		},
	}

	// Add description attribute
	maps.Copy(resp.Schema.Attributes, nbschema.DescriptionOnlyAttributes("module bay"))

	// Add tags and custom fields
	resp.Schema.Attributes["tags"] = nbschema.TagsSlugAttribute()
	resp.Schema.Attributes["custom_fields"] = nbschema.CustomFieldsAttribute()
}

func (r *ModuleBayResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = nbschema.ImportIdentityWithCustomFieldsSchema()
}

// Configure adds the provider configured client to the resource.
func (r *ModuleBayResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *ModuleBayResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ModuleBayResourceModel
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

	// Build request
	apiReq := netbox.NewModuleBayRequest(*device, data.Name.ValueString())

	// Set optional fields
	utils.ApplyLabel(apiReq, data.Label)
	if !data.Position.IsNull() && !data.Position.IsUnknown() {
		apiReq.SetPosition(data.Position.ValueString())
	}

	// Handle description and metadata fields
	utils.ApplyDescription(apiReq, data.Description)
	utils.ApplyTagsFromSlugs(ctx, r.client, apiReq, data.Tags, &resp.Diagnostics)
	utils.ApplyCustomFields(ctx, apiReq, data.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Creating module bay", map[string]interface{}{
		"device": data.Device.ValueString(),
		"name":   data.Name.ValueString(),
	})
	response, httpResp, err := r.client.DcimAPI.DcimModuleBaysCreate(ctx).ModuleBayRequest(*apiReq).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating module bay",
			utils.FormatAPIError(fmt.Sprintf("create module bay %s", data.Name.ValueString()), err, httpResp),
		)
		return
	}

	// Save plan state for filter-to-owned pattern
	planTags := data.Tags
	planCustomFields := data.CustomFields

	// Map response to model
	r.mapResponseToModel(ctx, response, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Apply filter-to-owned pattern for tags and custom_fields
	wasExplicitlyEmpty := !planTags.IsNull() && !planTags.IsUnknown() && len(planTags.Elements()) == 0
	switch {
	case response.HasTags() && len(response.GetTags()) > 0:
		tagSlugs := make([]string, 0, len(response.GetTags()))
		for _, tag := range response.GetTags() {
			tagSlugs = append(tagSlugs, tag.GetSlug())
		}
		data.Tags = utils.TagsSlugToSet(ctx, tagSlugs)
	case wasExplicitlyEmpty:
		data.Tags = types.SetValueMust(types.StringType, []attr.Value{})
	default:
		data.Tags = types.SetNull(types.StringType)
	}
	data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, planCustomFields, response.GetCustomFields(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(data.ID.ValueString()), data.CustomFields, &resp.Diagnostics)
	tflog.Trace(ctx, "Created module bay", map[string]interface{}{
		"id":   data.ID.ValueString(),
		"name": data.Name.ValueString(),
	})
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read refreshes the resource state.
func (r *ModuleBayResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ModuleBayResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	bayID, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Module Bay ID",
			fmt.Sprintf("Module Bay ID must be a number, got: %s", data.ID.ValueString()),
		)
		return
	}
	tflog.Debug(ctx, "Reading module bay", map[string]interface{}{
		"id": bayID,
	})
	response, httpResp, err := r.client.DcimAPI.DcimModuleBaysRetrieve(ctx, bayID).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading module bay",
			utils.FormatAPIError(fmt.Sprintf("read module bay ID %d", bayID), err, httpResp),
		)
		return
	}

	// Save state for filter-to-owned pattern
	stateTags := data.Tags
	stateCustomFields := data.CustomFields

	// Map response to model
	r.mapResponseToModel(ctx, response, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Preserve null/empty state values for tags and custom_fields
	wasExplicitlyEmpty := !stateTags.IsNull() && !stateTags.IsUnknown() && len(stateTags.Elements()) == 0
	switch {
	case response.HasTags() && len(response.GetTags()) > 0:
		tagSlugs := make([]string, 0, len(response.GetTags()))
		for _, tag := range response.GetTags() {
			tagSlugs = append(tagSlugs, tag.GetSlug())
		}
		data.Tags = utils.TagsSlugToSet(ctx, tagSlugs)
	case wasExplicitlyEmpty:
		data.Tags = types.SetValueMust(types.StringType, []attr.Value{})
	default:
		data.Tags = types.SetNull(types.StringType)
	}
	data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, stateCustomFields, response.GetCustomFields(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(data.ID.ValueString()), data.CustomFields, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates the resource.
func (r *ModuleBayResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Read both state and plan for merge-aware custom fields
	var state, plan ModuleBayResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	bayID, err := utils.ParseID(plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Module Bay ID",
			fmt.Sprintf("Module Bay ID must be a number, got: %s", plan.ID.ValueString()),
		)
		return
	}

	// Lookup device
	device, diags := lookup.LookupDevice(ctx, r.client, plan.Device.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build request
	apiReq := netbox.NewModuleBayRequest(*device, plan.Name.ValueString())

	// Set optional fields
	utils.ApplyLabel(apiReq, plan.Label)
	switch {
	case plan.Position.IsUnknown():
		// Leave position unset when unknown to avoid sending an invalid value.
	case plan.Position.IsNull():
		// NetBox PATCH semantics: omitted fields are preserved, so explicitly clear.
		apiReq.SetPosition("")
	default:
		apiReq.SetPosition(plan.Position.ValueString())
	}

	// Handle description and metadata fields (merge-aware)
	utils.ApplyDescription(apiReq, plan.Description)
	if utils.IsSet(plan.Tags) {
		utils.ApplyTagsFromSlugs(ctx, r.client, apiReq, plan.Tags, &resp.Diagnostics)
	} else if utils.IsSet(state.Tags) {
		utils.ApplyTagsFromSlugs(ctx, r.client, apiReq, state.Tags, &resp.Diagnostics)
	}
	if resp.Diagnostics.HasError() {
		return
	}
	utils.ApplyCustomFieldsWithMerge(ctx, apiReq, plan.CustomFields, state.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Updating module bay", map[string]interface{}{
		"id": bayID,
	})
	response, httpResp, err := r.client.DcimAPI.DcimModuleBaysUpdate(ctx, bayID).ModuleBayRequest(*apiReq).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating module bay",
			utils.FormatAPIError(fmt.Sprintf("update module bay ID %d", bayID), err, httpResp),
		)
		return
	}

	// Save plan state for filter-to-owned pattern
	planTags := plan.Tags
	planCustomFields := plan.CustomFields

	// Map response to model
	r.mapResponseToModel(ctx, response, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Apply filter-to-owned pattern for tags and custom_fields
	wasExplicitlyEmpty := !planTags.IsNull() && !planTags.IsUnknown() && len(planTags.Elements()) == 0
	switch {
	case response.HasTags() && len(response.GetTags()) > 0:
		tagSlugs := make([]string, 0, len(response.GetTags()))
		for _, tag := range response.GetTags() {
			tagSlugs = append(tagSlugs, tag.GetSlug())
		}
		plan.Tags = utils.TagsSlugToSet(ctx, tagSlugs)
	case wasExplicitlyEmpty:
		plan.Tags = types.SetValueMust(types.StringType, []attr.Value{})
	default:
		plan.Tags = types.SetNull(types.StringType)
	}
	plan.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, planCustomFields, response.GetCustomFields(), &resp.Diagnostics)
	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(plan.ID.ValueString()), plan.CustomFields, &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete deletes the resource.
func (r *ModuleBayResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ModuleBayResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	bayID, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Module Bay ID",
			fmt.Sprintf("Module Bay ID must be a number, got: %s", data.ID.ValueString()),
		)
		return
	}
	tflog.Debug(ctx, "Deleting module bay", map[string]interface{}{
		"id": bayID,
	})
	httpResp, err := r.client.DcimAPI.DcimModuleBaysDestroy(ctx, bayID).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			return
		}
		resp.Diagnostics.AddError(
			"Error deleting module bay",
			utils.FormatAPIError(fmt.Sprintf("delete module bay ID %d", bayID), err, httpResp),
		)
		return
	}
}

// ImportState imports an existing resource.
func (r *ModuleBayResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if parsed, ok := utils.ParseImportIdentityCustomFields(ctx, req.Identity, &resp.Diagnostics); ok {
		if resp.Diagnostics.HasError() {
			return
		}
		if parsed.ID == "" {
			resp.Diagnostics.AddError("Invalid import identity", "Identity id must be provided")
			return
		}

		bayID, err := utils.ParseID(parsed.ID)
		if err != nil {
			resp.Diagnostics.AddError(
				"Invalid Import ID",
				fmt.Sprintf("Module Bay ID must be a number, got: %s", parsed.ID),
			)
			return
		}
		response, httpResp, err := r.client.DcimAPI.DcimModuleBaysRetrieve(ctx, bayID).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error importing module bay",
				utils.FormatAPIError(fmt.Sprintf("import module bay ID %d", bayID), err, httpResp),
			)
			return
		}
		var data ModuleBayResourceModel
		if device := response.GetDevice(); device.Id != 0 {
			data.Device = types.StringValue(device.GetName())
		}
		if response.HasTags() && len(response.GetTags()) > 0 {
			tagSlugs := make([]string, 0, len(response.GetTags()))
			for _, tag := range response.GetTags() {
				tagSlugs = append(tagSlugs, tag.GetSlug())
			}
			data.Tags = utils.TagsSlugToSet(ctx, tagSlugs)
		} else {
			data.Tags = types.SetNull(types.StringType)
		}
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

	bayID, err := utils.ParseID(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Module Bay ID must be a number, got: %s", req.ID),
		)
		return
	}
	response, httpResp, err := r.client.DcimAPI.DcimModuleBaysRetrieve(ctx, bayID).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error importing module bay",
			utils.FormatAPIError(fmt.Sprintf("import module bay ID %d", bayID), err, httpResp),
		)
		return
	}
	var data ModuleBayResourceModel
	r.mapResponseToModel(ctx, response, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Ensure tags/custom_fields have concrete types for import
	if response.HasTags() && len(response.GetTags()) > 0 {
		tagSlugs := make([]string, 0, len(response.GetTags()))
		for _, tag := range response.GetTags() {
			tagSlugs = append(tagSlugs, tag.GetSlug())
		}
		data.Tags = utils.TagsSlugToSet(ctx, tagSlugs)
	} else {
		data.Tags = types.SetNull(types.StringType)
	}
	if response.HasCustomFields() {
		data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, data.CustomFields, response.GetCustomFields(), &resp.Diagnostics)
	} else {
		data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)
	}
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// mapResponseToModel maps the API response to the Terraform model.
func (r *ModuleBayResource) mapResponseToModel(ctx context.Context, moduleBay *netbox.ModuleBay, data *ModuleBayResourceModel, diags *diag.Diagnostics) {
	data.ID = types.StringValue(fmt.Sprintf("%d", moduleBay.GetId()))
	data.Name = types.StringValue(moduleBay.GetName())

	// Map device - preserve user's input format
	if device := moduleBay.GetDevice(); device.Id != 0 {
		data.Device = utils.UpdateReferenceAttribute(data.Device, device.GetName(), "", device.GetId())
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

	// Tags and custom fields are handled in Create/Read/Update methods using filter-to-owned pattern.
}
