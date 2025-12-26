// Package resources provides Terraform resource implementations for NetBox objects.

package resources

import (
	"context"
	"fmt"
	"maps"

	"github.com/bab3l/go-netbox"
	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
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
	_ resource.Resource = &InventoryItemRoleResource{}

	_ resource.ResourceWithConfigure = &InventoryItemRoleResource{}

	_ resource.ResourceWithImportState = &InventoryItemRoleResource{}
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
	ID types.String `tfsdk:"id"`

	Name types.String `tfsdk:"name"`

	Slug types.String `tfsdk:"slug"`

	Color types.String `tfsdk:"color"`

	Description types.String `tfsdk:"description"`

	DisplayName types.String `tfsdk:"display_name"`

	Tags types.Set `tfsdk:"tags"`

	CustomFields types.Set `tfsdk:"custom_fields"`
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

				Computed: true,

				PlanModifiers: []planmodifier.String{

					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"name": schema.StringAttribute{

				MarkdownDescription: "The name of the inventory item role.",

				Required: true,
			},

			"slug": schema.StringAttribute{

				MarkdownDescription: "A unique slug identifier for the inventory item role.",

				Required: true,
			},

			"color": nbschema.ComputedColorAttribute("inventory item role"),

			"display_name": nbschema.DisplayNameAttribute("inventory item role"),
		},
	}

	// Add description attribute
	maps.Copy(resp.Schema.Attributes, nbschema.DescriptionOnlyAttributes("inventory item role"))

	// Add common metadata attributes (tags, custom_fields)
	maps.Copy(resp.Schema.Attributes, nbschema.CommonMetadataAttributes())
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

	if !data.Description.IsNull() && !data.Description.IsUnknown() {

		apiReq.SetDescription(data.Description.ValueString())

	}

	// Handle tags

	if !data.Tags.IsNull() && !data.Tags.IsUnknown() {

		tags, tagDiags := utils.TagModelsToNestedTagRequests(ctx, data.Tags)

		resp.Diagnostics.Append(tagDiags...)

		if resp.Diagnostics.HasError() {

			return

		}

		apiReq.SetTags(tags)

	}

	// Handle custom fields

	if !data.CustomFields.IsNull() && !data.CustomFields.IsUnknown() {

		var cfModels []utils.CustomFieldModel

		resp.Diagnostics.Append(data.CustomFields.ElementsAs(ctx, &cfModels, false)...)

		if resp.Diagnostics.HasError() {

			return

		}

		apiReq.SetCustomFields(utils.CustomFieldModelsToMap(cfModels))

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

		"id": data.ID.ValueString(),

		"name": data.Name.ValueString(),
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

// Read refreshes the resource state.

func (r *InventoryItemRoleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

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

	tflog.Debug(ctx, "Reading inventory item role", map[string]interface{}{

		"id": roleID,
	})

	response, httpResp, err := r.client.DcimAPI.DcimInventoryItemRolesRetrieve(ctx, roleID).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		if httpResp != nil && httpResp.StatusCode == 404 {

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

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

// Update updates the resource.

func (r *InventoryItemRoleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var data InventoryItemRoleResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

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

	// Build request

	apiReq := netbox.NewInventoryItemRoleRequest(data.Name.ValueString(), data.Slug.ValueString())

	// Set optional fields

	if !data.Color.IsNull() && !data.Color.IsUnknown() {

		apiReq.SetColor(data.Color.ValueString())

	}

	if !data.Description.IsNull() && !data.Description.IsUnknown() {

		apiReq.SetDescription(data.Description.ValueString())

	}

	// Handle tags

	if !data.Tags.IsNull() && !data.Tags.IsUnknown() {

		tags, tagDiags := utils.TagModelsToNestedTagRequests(ctx, data.Tags)

		resp.Diagnostics.Append(tagDiags...)

		if resp.Diagnostics.HasError() {

			return

		}

		apiReq.SetTags(tags)

	}

	// Handle custom fields

	if !data.CustomFields.IsNull() && !data.CustomFields.IsUnknown() {

		var cfModels []utils.CustomFieldModel

		resp.Diagnostics.Append(data.CustomFields.ElementsAs(ctx, &cfModels, false)...)

		if resp.Diagnostics.HasError() {

			return

		}

		apiReq.SetCustomFields(utils.CustomFieldModelsToMap(cfModels))

	}

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

	r.mapResponseToModel(ctx, response, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {

		return

	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

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

		if httpResp != nil && httpResp.StatusCode == 404 {

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

	roleID, err := utils.ParseID(req.ID)

	if err != nil {

		resp.Diagnostics.AddError(

			"Invalid Import ID",

			fmt.Sprintf("Inventory Item Role ID must be a number, got: %s", req.ID),
		)

		return

	}

	response, httpResp, err := r.client.DcimAPI.DcimInventoryItemRolesRetrieve(ctx, roleID).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		resp.Diagnostics.AddError(

			"Error importing inventory item role",

			utils.FormatAPIError(fmt.Sprintf("import inventory item role ID %d", roleID), err, httpResp),
		)

		return

	}

	var data InventoryItemRoleResourceModel

	r.mapResponseToModel(ctx, response, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {

		return

	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

// mapResponseToModel maps the API response to the Terraform model.

func (r *InventoryItemRoleResource) mapResponseToModel(ctx context.Context, role *netbox.InventoryItemRole, data *InventoryItemRoleResourceModel, diags *diag.Diagnostics) {

	data.ID = types.StringValue(fmt.Sprintf("%d", role.GetId()))

	data.Name = types.StringValue(role.GetName())

	data.Slug = types.StringValue(role.GetSlug())

	// DisplayName
	if role.Display != "" {
		data.DisplayName = types.StringValue(role.Display)
	} else {
		data.DisplayName = types.StringNull()
	}

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

	// Handle tags

	if role.HasTags() && len(role.GetTags()) > 0 {

		tags := utils.NestedTagsToTagModels(role.GetTags())

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

	if role.HasCustomFields() {

		apiCustomFields := role.GetCustomFields()

		var stateCustomFieldModels []utils.CustomFieldModel

		if !data.CustomFields.IsNull() && !data.CustomFields.IsUnknown() {

			data.CustomFields.ElementsAs(ctx, &stateCustomFieldModels, false)

		}

		customFields := utils.MapToCustomFieldModels(apiCustomFields, stateCustomFieldModels)

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
