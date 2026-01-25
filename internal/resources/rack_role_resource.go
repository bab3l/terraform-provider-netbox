// Package resources contains Terraform resource implementations for the Netbox provider.
//

// This package integrates with the go-netbox OpenAPI client to provide
// CRUD operations for Netbox resources via Terraform.

package resources

import (
	"context"
	"fmt"
	"maps"

	"github.com/bab3l/go-netbox"
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

var _ resource.Resource = &RackRoleResource{}

var _ resource.ResourceWithImportState = &RackRoleResource{}
var _ resource.ResourceWithIdentity = &RackRoleResource{}

func NewRackRoleResource() resource.Resource {
	return &RackRoleResource{}
}

// RackRoleResource defines the resource implementation.

type RackRoleResource struct {
	client *netbox.APIClient
}

// RackRoleResourceModel describes the resource data model.

type RackRoleResourceModel struct {
	ID types.String `tfsdk:"id"`

	Name types.String `tfsdk:"name"`

	Slug types.String `tfsdk:"slug"`

	Color types.String `tfsdk:"color"`

	Description types.String `tfsdk:"description"`

	Tags types.Set `tfsdk:"tags"`

	CustomFields types.Set `tfsdk:"custom_fields"`
}

func (r *RackRoleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_rack_role"
}

func (r *RackRoleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a rack role in Netbox. Rack roles are used to categorize racks by their function or purpose within the data center (e.g., 'Network', 'Compute', 'Storage').",

		Attributes: map[string]schema.Attribute{
			"id": nbschema.IDAttribute("rack role"),

			"name": nbschema.NameAttribute("rack role", 100),

			"slug": nbschema.SlugAttribute("rack role"),

			"color": nbschema.ComputedColorAttribute("rack role"),
		},
	}

	// Add description attribute
	maps.Copy(resp.Schema.Attributes, nbschema.DescriptionOnlyAttributes("rack role"))

	// Add tags and custom fields attributes
	resp.Schema.Attributes["tags"] = nbschema.TagsSlugAttribute()
	resp.Schema.Attributes["custom_fields"] = nbschema.CustomFieldsAttribute()
}

func (r *RackRoleResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = nbschema.ImportIdentityWithCustomFieldsSchema()
}

func (r *RackRoleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.

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

func (r *RackRoleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data RackRoleResourceModel

	// Read Terraform plan data into the model

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create rack role using go-netbox client

	tflog.Debug(ctx, "Creating rack role", map[string]interface{}{
		"name": data.Name.ValueString(),

		"slug": data.Slug.ValueString(),
	})

	// Prepare the rack role request

	rackRoleRequest := netbox.RackRoleRequest{
		Name: data.Name.ValueString(),

		Slug: data.Slug.ValueString(),
	}

	// Set optional fields if provided

	if !data.Color.IsNull() && !data.Color.IsUnknown() {
		color := data.Color.ValueString()
		rackRoleRequest.Color = &color
	}

	// Use utils helper for description
	utils.ApplyDescription(&rackRoleRequest, data.Description)

	// Apply tags and custom fields
	utils.ApplyTagsFromSlugs(ctx, r.client, &rackRoleRequest, data.Tags, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	utils.ApplyCustomFields(ctx, &rackRoleRequest, data.CustomFields, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create the rack role via API

	rackRole, httpResp, err := r.client.DcimAPI.DcimRackRolesCreate(ctx).RackRoleRequest(rackRoleRequest).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {
		// Use enhanced error handler that detects duplicates and provides import hints

		handler := utils.CreateErrorHandler{
			ResourceType: "netbox_rack_role",

			ResourceName: "this.rack_role",

			SlugValue: data.Slug.ValueString(),

			LookupFunc: func(lookupCtx context.Context, slug string) (string, error) {
				// Try to look up existing rack role by slug

				list, _, lookupErr := r.client.DcimAPI.DcimRackRolesList(lookupCtx).Slug([]string{slug}).Execute()

				if lookupErr != nil {
					return "", lookupErr
				}

				if list != nil && len(list.Results) > 0 {
					return fmt.Sprintf("%d", list.Results[0].GetId()), nil
				}

				return "", nil
			},
		}

		handler.HandleCreateError(ctx, err, httpResp, &resp.Diagnostics)

		return
	}

	if httpResp.StatusCode != 201 {
		resp.Diagnostics.AddError(

			"Error creating rack role",

			fmt.Sprintf("Expected HTTP 201, got: %d", httpResp.StatusCode),
		)

		return
	}

	// Map response to state

	r.mapRackRoleToState(ctx, rackRole, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "created a rack role resource")

	// Save data into Terraform state
	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(data.ID.ValueString()), data.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RackRoleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data RackRoleResourceModel

	// Read Terraform prior state data into the model

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get the rack role ID from state

	rackRoleID := data.ID.ValueString()

	tflog.Debug(ctx, "Reading rack role", map[string]interface{}{
		"id": rackRoleID,
	})

	// Parse the rack role ID to int32 for the API call

	var rackRoleIDInt int32

	rackRoleIDInt, err := utils.ParseID(rackRoleID)

	if err != nil {
		resp.Diagnostics.AddError(

			"Invalid Rack Role ID",

			fmt.Sprintf("Rack Role ID must be a number, got: %s", rackRoleID),
		)

		return
	}

	// Retrieve the rack role via API

	rackRole, httpResp, err := r.client.DcimAPI.DcimRackRolesRetrieve(ctx, rackRoleIDInt).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			// Rack role no longer exists, remove from state

			resp.State.RemoveResource(ctx)

			return
		}

		resp.Diagnostics.AddError(

			"Error reading rack role",

			utils.FormatAPIError(fmt.Sprintf("read rack role ID %s", rackRoleID), err, httpResp),
		)

		return
	}

	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(

			"Error reading rack role",

			fmt.Sprintf("Expected HTTP 200, got: %d", httpResp.StatusCode),
		)

		return
	}

	// Map response to state

	r.mapRackRoleToState(ctx, rackRole, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(data.ID.ValueString()), data.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RackRoleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, plan RackRoleResourceModel

	// Read both state and plan

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get the rack role ID from plan

	rackRoleID := plan.ID.ValueString()

	tflog.Debug(ctx, "Updating rack role", map[string]interface{}{
		"id": rackRoleID,

		"name": plan.Name.ValueString(),

		"slug": plan.Slug.ValueString(),
	})

	// Parse the rack role ID to int32 for the API call

	var rackRoleIDInt int32

	rackRoleIDInt, err := utils.ParseID(rackRoleID)

	if err != nil {
		resp.Diagnostics.AddError(

			"Invalid Rack Role ID",

			fmt.Sprintf("Rack Role ID must be a number, got: %s", rackRoleID),
		)

		return
	}

	// Prepare the rack role update request

	rackRoleRequest := netbox.RackRoleRequest{
		Name: plan.Name.ValueString(),

		Slug: plan.Slug.ValueString(),
	}

	// Set optional fields if provided

	if !plan.Color.IsNull() && !plan.Color.IsUnknown() {
		color := plan.Color.ValueString()
		rackRoleRequest.Color = &color
	}

	// Use utils helper for description
	utils.ApplyDescription(&rackRoleRequest, plan.Description)

	// Apply tags and custom fields
	utils.ApplyTagsFromSlugs(ctx, r.client, &rackRoleRequest, plan.Tags, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	utils.ApplyCustomFieldsWithMerge(ctx, &rackRoleRequest, plan.CustomFields, state.CustomFields, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update the rack role via API

	rackRole, httpResp, err := r.client.DcimAPI.DcimRackRolesUpdate(ctx, rackRoleIDInt).RackRoleRequest(rackRoleRequest).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {
		resp.Diagnostics.AddError(

			"Error updating rack role",

			utils.FormatAPIError(fmt.Sprintf("update rack role ID %s", rackRoleID), err, httpResp),
		)

		return
	}

	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(

			"Error updating rack role",

			fmt.Sprintf("Expected HTTP 200, got: %d", httpResp.StatusCode),
		)

		return
	}

	// Map response to state

	r.mapRackRoleToState(ctx, rackRole, &plan, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(plan.ID.ValueString()), plan.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *RackRoleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data RackRoleResourceModel

	// Read Terraform prior state data into the model

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get the rack role ID from state

	rackRoleID := data.ID.ValueString()

	tflog.Debug(ctx, "Deleting rack role", map[string]interface{}{
		"id": rackRoleID,
	})

	// Parse the rack role ID to int32 for the API call

	var rackRoleIDInt int32

	rackRoleIDInt, err := utils.ParseID(rackRoleID)

	if err != nil {
		resp.Diagnostics.AddError(

			"Invalid Rack Role ID",

			fmt.Sprintf("Rack Role ID must be a number, got: %s", rackRoleID),
		)

		return
	}

	// Delete the rack role via API

	httpResp, err := r.client.DcimAPI.DcimRackRolesDestroy(ctx, rackRoleIDInt).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			return
		}

		resp.Diagnostics.AddError(

			"Error deleting rack role",

			utils.FormatAPIError(fmt.Sprintf("delete rack role ID %s", rackRoleID), err, httpResp),
		)

		return
	}

	if httpResp.StatusCode != 204 {
		resp.Diagnostics.AddError(

			"Error deleting rack role",

			fmt.Sprintf("Expected HTTP 204, got: %d", httpResp.StatusCode),
		)

		return
	}

	tflog.Trace(ctx, "deleted a rack role resource")
}

func (r *RackRoleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if parsed, ok := utils.ParseImportIdentityCustomFields(ctx, req.Identity, &resp.Diagnostics); ok {
		if resp.Diagnostics.HasError() {
			return
		}
		if parsed.ID == "" {
			resp.Diagnostics.AddError("Invalid import identity", "Identity id must be provided")
			return
		}

		rackRoleIDInt, err := utils.ParseID(parsed.ID)
		if err != nil {
			resp.Diagnostics.AddError("Invalid Rack Role ID", fmt.Sprintf("Rack Role ID must be a number, got: %s", parsed.ID))
			return
		}
		rackRole, httpResp, err := r.client.DcimAPI.DcimRackRolesRetrieve(ctx, rackRoleIDInt).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError("Error importing rack role", utils.FormatAPIError(fmt.Sprintf("read rack role ID %s", parsed.ID), err, httpResp))
			return
		}

		var data RackRoleResourceModel
		data.Tags = utils.PopulateTagsSlugFromAPI(ctx, rackRole.HasTags(), rackRole.GetTags(), data.Tags)
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

		r.mapRackRoleToState(ctx, rackRole, &data, &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}
		if parsed.HasCustomFields {
			data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, data.CustomFields, rackRole.GetCustomFields(), &resp.Diagnostics)
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

// mapRackRoleToState maps a RackRole API response to the Terraform state model.

func (r *RackRoleResource) mapRackRoleToState(ctx context.Context, rackRole *netbox.RackRole, data *RackRoleResourceModel, diags *diag.Diagnostics) {
	data.ID = types.StringValue(fmt.Sprintf("%d", rackRole.GetId()))

	data.Name = types.StringValue(rackRole.GetName())

	data.Slug = types.StringValue(rackRole.GetSlug())

	// Handle color

	if rackRole.HasColor() && rackRole.GetColor() != "" {
		data.Color = types.StringValue(rackRole.GetColor())
	} else if !data.Color.IsNull() {
		// Preserve null if originally null and API returns empty

		data.Color = types.StringNull()
	}

	// Handle description

	if rackRole.HasDescription() && rackRole.GetDescription() != "" {
		data.Description = types.StringValue(rackRole.GetDescription())
	} else if !data.Description.IsNull() {
		// Preserve null if originally null and API returns empty

		data.Description = types.StringNull()
	}

	// Handle tags (filter to owned, slug list format)
	data.Tags = utils.PopulateTagsSlugFilteredToOwned(ctx, rackRole.HasTags(), rackRole.GetTags(), data.Tags)

	// Handle custom fields (filter to owned)
	data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, data.CustomFields, rackRole.GetCustomFields(), diags)
	if diags.HasError() {
		return
	}
}
