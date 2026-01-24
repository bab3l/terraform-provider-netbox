// Package resources provides Terraform resource implementations for NetBox objects.

package resources

import (
	"context"
	"fmt"
	"maps"
	"net/http"

	"github.com/bab3l/go-netbox"
	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &InventoryItemRoleResource{}
	_ resource.ResourceWithConfigure   = &InventoryItemRoleResource{}
	_ resource.ResourceWithImportState = &InventoryItemRoleResource{}
	_ resource.ResourceWithIdentity    = &InventoryItemRoleResource{}
)

// NewInventoryItemRoleResource returns a new resource implementing the inventory item role resource.
func NewInventoryItemRoleResource() resource.Resource {
	return &InventoryItemRoleResource{}
}

// InventoryItemRoleResource defines the resource implementation.
type InventoryItemRoleResource struct {
	client *netbox.APIClient
}

// InventoryItemRoleResourceModel describes the resource data model.
type InventoryItemRoleResourceModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Slug         types.String `tfsdk:"slug"`
	Color        types.String `tfsdk:"color"`
	Description  types.String `tfsdk:"description"`
	Tags         types.Set    `tfsdk:"tags"`
	CustomFields types.Set    `tfsdk:"custom_fields"`
}

// Metadata returns the resource type name.
func (r *InventoryItemRoleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_inventory_item_role"
}

// Schema defines the schema for the resource.
func (r *InventoryItemRoleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an inventory item role in NetBox. Inventory item roles define the functional purpose of inventory items within devices.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique numeric ID of the inventory item role.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the inventory item role.",
				Required:            true,
			},
			"slug": schema.StringAttribute{
				MarkdownDescription: "A unique slug identifier for the inventory item role.",
				Required:            true,
			},
			"color": nbschema.ComputedColorAttribute("inventory item role"),
		},
	}

	// Add description attribute
	maps.Copy(resp.Schema.Attributes, nbschema.DescriptionOnlyAttributes("inventory item role"))

	// Add tags and custom_fields
	resp.Schema.Attributes["tags"] = nbschema.TagsSlugAttribute()
	resp.Schema.Attributes["custom_fields"] = nbschema.CustomFieldsAttribute()
}

func (r *InventoryItemRoleResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = nbschema.ImportIdentityWithCustomFieldsSchema()
}

func (r *InventoryItemRoleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *InventoryItemRoleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data InventoryItemRoleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build request
	apiReq := netbox.NewInventoryItemRoleRequest(data.Name.ValueString(), data.Slug.ValueString())

	// Set optional fields
	if !data.Color.IsNull() && !data.Color.IsUnknown() {
		apiReq.SetColor(data.Color.ValueString())
	}

	// Handle description
	utils.ApplyDescription(apiReq, data.Description)

	// Handle tags and custom_fields
	utils.ApplyTagsFromSlugs(ctx, r.client, apiReq, data.Tags, &resp.Diagnostics)
	utils.ApplyCustomFields(ctx, apiReq, data.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Creating inventory item role", map[string]interface{}{
		"name": data.Name.ValueString(),
		"slug": data.Slug.ValueString(),
	})
	response, httpResp, err := r.client.DcimAPI.DcimInventoryItemRolesCreate(ctx).InventoryItemRoleRequest(*apiReq).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating inventory item role",
			utils.FormatAPIError(fmt.Sprintf("create inventory item role %s", data.Name.ValueString()), err, httpResp),
		)
		return
	}

	// Map response to model
	r.mapResponseToModel(ctx, response, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Trace(ctx, "Created inventory item role", map[string]interface{}{
		"id":   data.ID.ValueString(),
		"name": data.Name.ValueString(),
	})
	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(data.ID.ValueString()), data.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read refreshes the resource state.
func (r *InventoryItemRoleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data InventoryItemRoleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Preserve original custom_fields from state for potential restoration
	originalCustomFields := data.CustomFields

	roleID, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Inventory Item Role ID",
			fmt.Sprintf("Inventory Item Role ID must be a number, got: %s", data.ID.ValueString()),
		)
		return
	}
	tflog.Debug(ctx, "Reading inventory item role", map[string]interface{}{
		"id": roleID,
	})

	response, httpResp, err := r.client.DcimAPI.DcimInventoryItemRolesRetrieve(ctx, roleID).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading inventory item role",
			utils.FormatAPIError(fmt.Sprintf("read inventory item role ID %d", roleID), err, httpResp),
		)
		return
	}

	// Map response to model
	r.mapResponseToModel(ctx, response, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// If custom_fields was null or empty before (not managed or explicitly cleared),
	// restore that state after mapping.
	if originalCustomFields.IsNull() || (utils.IsSet(originalCustomFields) && len(originalCustomFields.Elements()) == 0) {
		tflog.Debug(ctx, "Custom fields unmanaged/cleared, preserving original state during Read", map[string]interface{}{
			"was_null":  originalCustomFields.IsNull(),
			"was_empty": !originalCustomFields.IsNull() && len(originalCustomFields.Elements()) == 0,
		})
		data.CustomFields = originalCustomFields
	}

	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(data.ID.ValueString()), data.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates the resource.
func (r *InventoryItemRoleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Read BOTH state and plan for merge-aware custom fields
	var state, plan InventoryItemRoleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	roleID, err := utils.ParseID(plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Inventory Item Role ID",
			fmt.Sprintf("Inventory Item Role ID must be a number, got: %s", plan.ID.ValueString()),
		)
		return
	}

	// Build request
	apiReq := netbox.NewInventoryItemRoleRequest(plan.Name.ValueString(), plan.Slug.ValueString())

	// Set optional fields
	if !plan.Color.IsNull() && !plan.Color.IsUnknown() {
		apiReq.SetColor(plan.Color.ValueString())
	}

	// Handle description
	utils.ApplyDescription(apiReq, plan.Description)

	// Handle tags and custom_fields with merge
	utils.ApplyTagsFromSlugs(ctx, r.client, apiReq, plan.Tags, &resp.Diagnostics)
	utils.ApplyCustomFieldsWithMerge(ctx, apiReq, plan.CustomFields, state.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Store plan tags/customfields for filter-to-owned population
	planTags := plan.Tags
	planCustomFields := plan.CustomFields

	tflog.Debug(ctx, "Updating inventory item role", map[string]interface{}{
		"id": roleID,
	})
	response, httpResp, err := r.client.DcimAPI.DcimInventoryItemRolesUpdate(ctx, roleID).InventoryItemRoleRequest(*apiReq).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating inventory item role",
			utils.FormatAPIError(fmt.Sprintf("update inventory item role ID %d", roleID), err, httpResp),
		)
		return
	}

	// Map response to model
	plan.Tags = planTags
	plan.CustomFields = planCustomFields
	r.mapResponseToModel(ctx, response, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(plan.ID.ValueString()), plan.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete deletes the resource.
func (r *InventoryItemRoleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data InventoryItemRoleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	roleID, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Inventory Item Role ID",
			fmt.Sprintf("Inventory Item Role ID must be a number, got: %s", data.ID.ValueString()),
		)
		return
	}
	tflog.Debug(ctx, "Deleting inventory item role", map[string]interface{}{
		"id": roleID,
	})
	httpResp, err := r.client.DcimAPI.DcimInventoryItemRolesDestroy(ctx, roleID).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			return
		}
		resp.Diagnostics.AddError(
			"Error deleting inventory item role",
			utils.FormatAPIError(fmt.Sprintf("delete inventory item role ID %d", roleID), err, httpResp),
		)
		return
	}
}

// ImportState imports an existing resource.
func (r *InventoryItemRoleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if parsed, ok := utils.ParseImportIdentityCustomFields(ctx, req.Identity, &resp.Diagnostics); ok {
		if resp.Diagnostics.HasError() {
			return
		}
		if parsed.ID == "" {
			resp.Diagnostics.AddError("Invalid import identity", "Identity id must be provided")
			return
		}

		roleID, err := utils.ParseID(parsed.ID)
		if err != nil {
			resp.Diagnostics.AddError("Invalid Import ID", fmt.Sprintf("Inventory Item Role ID must be a number, got: %s", parsed.ID))
			return
		}
		response, httpResp, err := r.client.DcimAPI.DcimInventoryItemRolesRetrieve(ctx, roleID).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError("Error importing inventory item role", utils.FormatAPIError(fmt.Sprintf("import inventory item role ID %d", roleID), err, httpResp))
			return
		}

		var data InventoryItemRoleResourceModel
		if response.HasTags() {
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

	utils.ImportStatePassthroughIDWithValidation(ctx, req, resp, path.Root("id"), true)
}

// mapResponseToModel maps the API response to the Terraform model.
func (r *InventoryItemRoleResource) mapResponseToModel(ctx context.Context, role *netbox.InventoryItemRole, data *InventoryItemRoleResourceModel, diags *diag.Diagnostics) {
	data.ID = types.StringValue(fmt.Sprintf("%d", role.GetId()))
	data.Name = types.StringValue(role.GetName())
	data.Slug = types.StringValue(role.GetSlug())

	// Map color
	if color, ok := role.GetColorOk(); ok && color != nil && *color != "" {
		data.Color = types.StringValue(*color)
	} else {
		data.Color = types.StringNull()
	}

	// Map description
	if desc, ok := role.GetDescriptionOk(); ok && desc != nil && *desc != "" {
		data.Description = types.StringValue(*desc)
	} else {
		data.Description = types.StringNull()
	}

	// Handle tags (filter-to-owned)
	planTags := data.Tags
	wasExplicitlyEmpty := !planTags.IsNull() && !planTags.IsUnknown() && len(planTags.Elements()) == 0
	switch {
	case role.HasTags() && len(role.GetTags()) > 0:
		tagSlugs := make([]string, 0, len(role.GetTags()))
		for _, tag := range role.GetTags() {
			tagSlugs = append(tagSlugs, tag.GetSlug())
		}
		data.Tags = utils.TagsSlugToSet(ctx, tagSlugs)
	case wasExplicitlyEmpty:
		data.Tags = types.SetValueMust(types.StringType, []attr.Value{})
	default:
		data.Tags = types.SetNull(types.StringType)
	}

	// Handle custom fields using filter-to-owned helper
	data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, data.CustomFields, role.GetCustomFields(), diags)
}
