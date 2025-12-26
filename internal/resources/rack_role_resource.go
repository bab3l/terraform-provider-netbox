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

	DisplayName types.String `tfsdk:"display_name"`

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

			"display_name": nbschema.DisplayNameAttribute("rack role"),

			"name": nbschema.NameAttribute("rack role", 100),

			"slug": nbschema.SlugAttribute("rack role"),

			"color": nbschema.ComputedColorAttribute("rack role"),
		},
	}

	// Add description attribute
	maps.Copy(resp.Schema.Attributes, nbschema.DescriptionOnlyAttributes("rack role"))

	// Add common metadata attributes (tags, custom_fields)
	maps.Copy(resp.Schema.Attributes, nbschema.CommonMetadataAttributes())
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

	if !data.Description.IsNull() && !data.Description.IsUnknown() {

		description := data.Description.ValueString()

		rackRoleRequest.Description = &description

	}

	// Handle tags

	if !data.Tags.IsNull() && !data.Tags.IsUnknown() {

		var tags []utils.TagModel

		resp.Diagnostics.Append(data.Tags.ElementsAs(ctx, &tags, false)...)

		if resp.Diagnostics.HasError() {

			return

		}

		rackRoleRequest.Tags = utils.TagsToNestedTagRequests(tags)

	}

	// Handle custom fields

	if !data.CustomFields.IsNull() && !data.CustomFields.IsUnknown() {

		var customFields []utils.CustomFieldModel

		resp.Diagnostics.Append(data.CustomFields.ElementsAs(ctx, &customFields, false)...)

		if resp.Diagnostics.HasError() {

			return

		}

		rackRoleRequest.CustomFields = utils.CustomFieldsToMap(customFields)

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

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

func (r *RackRoleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var data RackRoleResourceModel

	// Read Terraform plan data into the model

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	// Get the rack role ID from state

	rackRoleID := data.ID.ValueString()

	tflog.Debug(ctx, "Updating rack role", map[string]interface{}{

		"id": rackRoleID,

		"name": data.Name.ValueString(),

		"slug": data.Slug.ValueString(),
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

		Name: data.Name.ValueString(),

		Slug: data.Slug.ValueString(),
	}

	// Set optional fields if provided

	if !data.Color.IsNull() && !data.Color.IsUnknown() {

		color := data.Color.ValueString()

		rackRoleRequest.Color = &color

	}

	if !data.Description.IsNull() && !data.Description.IsUnknown() {

		description := data.Description.ValueString()

		rackRoleRequest.Description = &description

	}

	// Handle tags

	if !data.Tags.IsNull() && !data.Tags.IsUnknown() {

		var tags []utils.TagModel

		resp.Diagnostics.Append(data.Tags.ElementsAs(ctx, &tags, false)...)

		if resp.Diagnostics.HasError() {

			return

		}

		rackRoleRequest.Tags = utils.TagsToNestedTagRequests(tags)

	}

	// Handle custom fields

	if !data.CustomFields.IsNull() && !data.CustomFields.IsUnknown() {

		var customFields []utils.CustomFieldModel

		resp.Diagnostics.Append(data.CustomFields.ElementsAs(ctx, &customFields, false)...)

		if resp.Diagnostics.HasError() {

			return

		}

		rackRoleRequest.CustomFields = utils.CustomFieldsToMap(customFields)

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

	r.mapRackRoleToState(ctx, rackRole, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {

		return

	}

	// Save updated data into Terraform state

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

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

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

}

// mapRackRoleToState maps a RackRole API response to the Terraform state model.

func (r *RackRoleResource) mapRackRoleToState(ctx context.Context, rackRole *netbox.RackRole, data *RackRoleResourceModel, diags *diag.Diagnostics) {

	data.ID = types.StringValue(fmt.Sprintf("%d", rackRole.GetId()))

	data.DisplayName = types.StringValue(rackRole.GetDisplay())

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

	// Handle tags

	if rackRole.HasTags() {

		tags := utils.NestedTagsToTagModels(rackRole.GetTags())

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

	if rackRole.HasCustomFields() && !data.CustomFields.IsNull() {

		var stateCustomFields []utils.CustomFieldModel

		cfDiags := data.CustomFields.ElementsAs(ctx, &stateCustomFields, false)

		diags.Append(cfDiags...)

		if diags.HasError() {

			return

		}

		customFields := utils.MapToCustomFieldModels(rackRole.GetCustomFields(), stateCustomFields)

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
