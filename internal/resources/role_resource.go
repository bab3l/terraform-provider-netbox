// Package resources contains Terraform resource implementations for the Netbox provider.

package resources

import (
	"context"
	"fmt"

	"github.com/bab3l/go-netbox"
	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/bab3l/terraform-provider-netbox/internal/validators"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.

var (
	_ resource.Resource = &RoleResource{}

	_ resource.ResourceWithConfigure = &RoleResource{}

	_ resource.ResourceWithImportState = &RoleResource{}
)

// NewRoleResource returns a new Role resource.

func NewRoleResource() resource.Resource {

	return &RoleResource{}

}

// RoleResource defines the resource implementation.

type RoleResource struct {
	client *netbox.APIClient
}

// RoleResourceModel describes the resource data model.

type RoleResourceModel struct {
	ID types.String `tfsdk:"id"`

	Name types.String `tfsdk:"name"`

	Slug types.String `tfsdk:"slug"`

	Weight types.Int64 `tfsdk:"weight"`

	Description types.String `tfsdk:"description"`

	Tags types.Set `tfsdk:"tags"`

	CustomFields types.Set `tfsdk:"custom_fields"`
}

// Metadata returns the resource type name.

func (r *RoleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {

	resp.TypeName = req.ProviderTypeName + "_role"

}

// Schema defines the schema for the resource.

func (r *RoleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	resp.Schema = schema.Schema{

		MarkdownDescription: "Manages an IPAM role in NetBox. Roles are used to categorize prefixes and VLANs by their functional purpose (e.g., Production, Development, Customer).",

		Attributes: map[string]schema.Attribute{

			"id": schema.StringAttribute{

				MarkdownDescription: "The unique numeric ID of the role.",

				Computed: true,

				PlanModifiers: []planmodifier.String{

					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"name": schema.StringAttribute{

				MarkdownDescription: "The name of the role.",

				Required: true,

				Validators: []validator.String{

					stringvalidator.LengthBetween(1, 100),
				},
			},

			"slug": schema.StringAttribute{

				MarkdownDescription: "URL-friendly unique identifier for the role.",

				Required: true,

				Validators: []validator.String{

					stringvalidator.LengthBetween(1, 100),

					validators.ValidSlug(),
				},
			},

			"weight": schema.Int64Attribute{

				MarkdownDescription: "Weight for sorting. Lower values appear first.",

				Optional: true,

				Computed: true,

				Default: int64default.StaticInt64(1000),
			},

			"description": schema.StringAttribute{

				MarkdownDescription: "A description of the role.",

				Optional: true,

				Validators: []validator.String{

					stringvalidator.LengthAtMost(200),
				},
			},

			"tags": nbschema.TagsAttribute(),

			"custom_fields": nbschema.CustomFieldsAttribute(),
		},
	}

}

// Configure adds the provider configured client to the resource.

func (r *RoleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {

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

// Create creates the resource and sets the initial Terraform state.

func (r *RoleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var data RoleResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	tflog.Debug(ctx, "Creating role", map[string]interface{}{

		"name": data.Name.ValueString(),

		"slug": data.Slug.ValueString(),
	})

	// Build the role request

	roleRequest, diags := r.buildRoleRequest(ctx, &data)

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {

		return

	}

	// Call the API

	role, httpResp, err := r.client.IpamAPI.IpamRolesCreate(ctx).RoleRequest(*roleRequest).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		resp.Diagnostics.AddError(

			"Error creating role",

			utils.FormatAPIError(fmt.Sprintf("create role %s", data.Name.ValueString()), err, httpResp),
		)

		return

	}

	tflog.Debug(ctx, "Created role", map[string]interface{}{

		"id": role.GetId(),

		"name": role.GetName(),
	})

	// Map response to state

	r.mapResponseToModel(ctx, role, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {

		return

	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

// Read refreshes the Terraform state with the latest data.

func (r *RoleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	var data RoleResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	// Parse the ID

	roleID, err := utils.ParseID(data.ID.ValueString())

	if err != nil {

		resp.Diagnostics.AddError(

			"Invalid Role ID",

			fmt.Sprintf("Role ID must be a number, got: %s", data.ID.ValueString()),
		)

		return

	}

	tflog.Debug(ctx, "Reading role", map[string]interface{}{

		"id": roleID,
	})

	// Call the API

	role, httpResp, err := r.client.IpamAPI.IpamRolesRetrieve(ctx, roleID).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		if httpResp != nil && httpResp.StatusCode == 404 {

			tflog.Debug(ctx, "Role not found, removing from state", map[string]interface{}{

				"id": roleID,
			})

			resp.State.RemoveResource(ctx)

			return

		}

		resp.Diagnostics.AddError(

			"Error reading role",

			utils.FormatAPIError(fmt.Sprintf("read role ID %d", roleID), err, httpResp),
		)

		return

	}

	// Map response to state

	r.mapResponseToModel(ctx, role, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {

		return

	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

// Update updates the resource and sets the updated Terraform state.

func (r *RoleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var data RoleResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	// Parse the ID

	roleID, err := utils.ParseID(data.ID.ValueString())

	if err != nil {

		resp.Diagnostics.AddError(

			"Invalid Role ID",

			fmt.Sprintf("Role ID must be a number, got: %s", data.ID.ValueString()),
		)

		return

	}

	tflog.Debug(ctx, "Updating role", map[string]interface{}{

		"id": roleID,

		"name": data.Name.ValueString(),
	})

	// Build the role request

	roleRequest, diags := r.buildRoleRequest(ctx, &data)

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {

		return

	}

	// Call the API

	role, httpResp, err := r.client.IpamAPI.IpamRolesUpdate(ctx, roleID).RoleRequest(*roleRequest).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		resp.Diagnostics.AddError(

			"Error updating role",

			utils.FormatAPIError(fmt.Sprintf("update role ID %d", roleID), err, httpResp),
		)

		return

	}

	tflog.Debug(ctx, "Updated role", map[string]interface{}{

		"id": role.GetId(),

		"name": role.GetName(),
	})

	// Map response to state

	r.mapResponseToModel(ctx, role, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {

		return

	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

// Delete deletes the resource and removes the Terraform state.

func (r *RoleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var data RoleResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	// Parse the ID

	roleID, err := utils.ParseID(data.ID.ValueString())

	if err != nil {

		resp.Diagnostics.AddError(

			"Invalid Role ID",

			fmt.Sprintf("Role ID must be a number, got: %s", data.ID.ValueString()),
		)

		return

	}

	tflog.Debug(ctx, "Deleting role", map[string]interface{}{

		"id": roleID,

		"name": data.Name.ValueString(),
	})

	// Call the API

	httpResp, err := r.client.IpamAPI.IpamRolesDestroy(ctx, roleID).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		if httpResp != nil && httpResp.StatusCode == 404 {

			// Resource already deleted

			return

		}

		resp.Diagnostics.AddError(

			"Error deleting role",

			utils.FormatAPIError(fmt.Sprintf("delete role ID %d", roleID), err, httpResp),
		)

		return

	}

	tflog.Debug(ctx, "Deleted role", map[string]interface{}{

		"id": roleID,
	})

}

// ImportState imports the resource state.

func (r *RoleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

}

// buildRoleRequest builds a RoleRequest from the Terraform model.

func (r *RoleResource) buildRoleRequest(ctx context.Context, data *RoleResourceModel) (*netbox.RoleRequest, diag.Diagnostics) {

	var diags diag.Diagnostics

	// Create the request with required fields

	roleRequest := netbox.NewRoleRequest(data.Name.ValueString(), data.Slug.ValueString())

	// Handle weight (optional)

	if !data.Weight.IsNull() && !data.Weight.IsUnknown() {

		weight, err := utils.SafeInt32FromValue(data.Weight)

		if err != nil {

			diags.AddError("Invalid value", fmt.Sprintf("Weight value overflow: %s", err))

			return nil, diags

		}

		roleRequest.Weight = &weight

	}

	// Handle description (optional)

	if !data.Description.IsNull() && !data.Description.IsUnknown() {

		desc := data.Description.ValueString()

		roleRequest.Description = &desc

	}

	// Handle tags

	if !data.Tags.IsNull() && !data.Tags.IsUnknown() {

		tags, tagDiags := utils.TagModelsToNestedTagRequests(ctx, data.Tags)

		diags.Append(tagDiags...)

		if diags.HasError() {

			return nil, diags

		}

		roleRequest.Tags = tags

	}

	// Handle custom fields

	if !data.CustomFields.IsNull() && !data.CustomFields.IsUnknown() {

		var customFieldModels []utils.CustomFieldModel

		cfDiags := data.CustomFields.ElementsAs(ctx, &customFieldModels, false)

		diags.Append(cfDiags...)

		if diags.HasError() {

			return nil, diags

		}

		roleRequest.CustomFields = utils.CustomFieldModelsToMap(customFieldModels)

	}

	return roleRequest, diags

}

// mapResponseToModel maps the API response to the Terraform model.

func (r *RoleResource) mapResponseToModel(ctx context.Context, role *netbox.Role, data *RoleResourceModel, diags *diag.Diagnostics) {

	data.ID = types.StringValue(fmt.Sprintf("%d", role.GetId()))

	data.Name = types.StringValue(role.GetName())

	data.Slug = types.StringValue(role.GetSlug())

	// Map weight

	if weight, ok := role.GetWeightOk(); ok && weight != nil {

		data.Weight = types.Int64Value(int64(*weight))

	} else {

		data.Weight = types.Int64Value(1000)

	}

	// Map description

	if desc, ok := role.GetDescriptionOk(); ok && desc != nil && *desc != "" {

		data.Description = types.StringValue(*desc)

	} else {

		data.Description = types.StringNull()

	}

	// Handle tags

	if role.HasTags() {

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

		if !data.CustomFields.IsNull() {

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
