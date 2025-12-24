// Package resources contains Terraform resource implementations for the Netbox provider.

package resources

import (
	"context"
	"fmt"

	"github.com/bab3l/go-netbox"
	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.

var _ resource.Resource = &ContactAssignmentResource{}

var _ resource.ResourceWithImportState = &ContactAssignmentResource{}

func NewContactAssignmentResource() resource.Resource {

	return &ContactAssignmentResource{}

}

// ContactAssignmentResource defines the resource implementation.

type ContactAssignmentResource struct {
	client *netbox.APIClient
}

// ContactAssignmentResourceModel describes the resource data model.

type ContactAssignmentResourceModel struct {
	ID types.String `tfsdk:"id"`

	ObjectType types.String `tfsdk:"object_type"`

	ObjectID types.String `tfsdk:"object_id"`

	Contact types.String `tfsdk:"contact_id"`

	Role types.String `tfsdk:"role_id"`

	Priority types.String `tfsdk:"priority"`

	DisplayName types.String `tfsdk:"display_name"`

	Tags types.Set `tfsdk:"tags"`

	CustomFields types.Set `tfsdk:"custom_fields"`
}

func (r *ContactAssignmentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {

	resp.TypeName = req.ProviderTypeName + "_contact_assignment"

}

func (r *ContactAssignmentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	resp.Schema = schema.Schema{

		MarkdownDescription: "Manages a contact assignment in Netbox. A contact assignment links a contact to any Netbox object (site, device, circuit, etc.) with an optional role and priority.",

		Attributes: map[string]schema.Attribute{

			"id": schema.StringAttribute{

				Computed: true,

				MarkdownDescription: "The unique identifier of the contact assignment.",

				PlanModifiers: []planmodifier.String{

					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"object_type": schema.StringAttribute{

				Required: true,

				MarkdownDescription: "The content type of the object to assign the contact to (e.g., `dcim.site`, `dcim.device`, `circuits.circuit`).",

				PlanModifiers: []planmodifier.String{

					stringplanmodifier.RequiresReplace(),
				},
			},

			"object_id": schema.StringAttribute{

				Required: true,

				MarkdownDescription: "The ID of the object to assign the contact to.",

				PlanModifiers: []planmodifier.String{

					stringplanmodifier.RequiresReplace(),
				},
			},

			"contact_id": schema.StringAttribute{

				Required: true,

				MarkdownDescription: "The ID of the contact to assign.",
			},

			"role_id": schema.StringAttribute{

				Optional: true,

				MarkdownDescription: "The ID of the contact role for this assignment.",
			},

			"priority": schema.StringAttribute{

				Optional: true,

				MarkdownDescription: "The priority of this contact assignment. Valid values are: `primary`, `secondary`, `tertiary`, `inactive`.",

				Validators: []validator.String{

					stringvalidator.OneOf("primary", "secondary", "tertiary", "inactive", ""),
				},
			},

			"display_name": nbschema.DisplayNameAttribute("contact assignment"),

			"tags": nbschema.TagsAttribute(),

			"custom_fields": nbschema.CustomFieldsAttribute(),
		},
	}

}

func (r *ContactAssignmentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {

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

func (r *ContactAssignmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var data ContactAssignmentResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	// Parse object ID

	objectID, err := utils.ParseID64(data.ObjectID.ValueString())

	if err != nil {

		resp.Diagnostics.AddError(

			"Invalid Object ID",

			fmt.Sprintf("Unable to parse object ID: %s", err),
		)

		return

	}

	// Parse contact ID

	contactID, err := utils.ParseID(data.Contact.ValueString())

	if err != nil {

		resp.Diagnostics.AddError(

			"Invalid Contact ID",

			fmt.Sprintf("Unable to parse contact ID: %s", err),
		)

		return

	}

	tflog.Debug(ctx, "Creating contact assignment", map[string]interface{}{

		"object_type": data.ObjectType.ValueString(),

		"object_id": objectID,

		"contact_id": contactID,
	})

	// Build the API request - need to create placeholder Brief objects

	briefContact := *netbox.NewBriefContactRequest("placeholder")

	assignmentRequest := netbox.NewWritableContactAssignmentRequest(

		data.ObjectType.ValueString(),

		objectID,

		briefContact,
	)

	// Use AdditionalProperties to pass the contact ID

	assignmentRequest.AdditionalProperties = make(map[string]interface{})

	assignmentRequest.AdditionalProperties["contact"] = int(contactID)

	// Set role if provided, otherwise unset it entirely

	if !data.Role.IsNull() && !data.Role.IsUnknown() && data.Role.ValueString() != "" {

		roleID, err := utils.ParseID(data.Role.ValueString())

		if err != nil {

			resp.Diagnostics.AddError(

				"Invalid Role ID",

				fmt.Sprintf("Unable to parse role ID: %s", err),
			)

			return

		}

		assignmentRequest.AdditionalProperties["role"] = int(roleID)

	} else {

		// Unset role entirely when not provided

		assignmentRequest.UnsetRole()

	}

	// Set priority if provided

	if !data.Priority.IsNull() && !data.Priority.IsUnknown() && data.Priority.ValueString() != "" {

		priority := netbox.BriefCircuitGroupAssignmentSerializerPriorityValue(data.Priority.ValueString())

		assignmentRequest.Priority = &priority

	}

	// Handle tags

	if !data.Tags.IsNull() {

		var tagModels []utils.TagModel

		diags := data.Tags.ElementsAs(ctx, &tagModels, false)

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {

			return

		}

		assignmentRequest.Tags = utils.TagsToNestedTagRequests(tagModels)

	}

	// Handle custom fields

	if !data.CustomFields.IsNull() {

		var customFieldModels []utils.CustomFieldModel

		diags := data.CustomFields.ElementsAs(ctx, &customFieldModels, false)

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {

			return

		}

		assignmentRequest.CustomFields = utils.CustomFieldModelsToMap(customFieldModels)

	}

	// Call the API to create the contact assignment

	assignment, httpResp, err := r.client.TenancyAPI.TenancyContactAssignmentsCreate(ctx).
		WritableContactAssignmentRequest(*assignmentRequest).
		Execute()

	if err != nil {

		resp.Diagnostics.AddError(

			"Error creating contact assignment",

			utils.FormatAPIError("create contact assignment", err, httpResp),
		)

		return

	}

	// Map response to state

	r.mapResponseToState(ctx, assignment, &data, &resp.Diagnostics)

	tflog.Debug(ctx, "Created contact assignment", map[string]interface{}{

		"id": assignment.GetId(),
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

func (r *ContactAssignmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	var data ContactAssignmentResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	// Parse the ID

	id, err := utils.ParseID(data.ID.ValueString())

	if err != nil {

		resp.Diagnostics.AddError(

			"Invalid ID format",

			fmt.Sprintf("Could not parse contact assignment ID '%s': %s", data.ID.ValueString(), err),
		)

		return

	}

	tflog.Debug(ctx, "Reading contact assignment", map[string]interface{}{

		"id": id,
	})

	assignment, httpResp, err := r.client.TenancyAPI.TenancyContactAssignmentsRetrieve(ctx, id).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		if httpResp != nil && httpResp.StatusCode == 404 {

			tflog.Info(ctx, "Contact assignment not found, removing from state", map[string]interface{}{

				"id": id,
			})

			resp.State.RemoveResource(ctx)

			return

		}

		resp.Diagnostics.AddError(

			"Error reading contact assignment",

			utils.FormatAPIError("read contact assignment", err, httpResp),
		)

		return

	}

	// Map response to state

	r.mapResponseToState(ctx, assignment, &data, &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

func (r *ContactAssignmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var data ContactAssignmentResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	// Parse the ID

	id, err := utils.ParseID(data.ID.ValueString())

	if err != nil {

		resp.Diagnostics.AddError(

			"Invalid ID format",

			fmt.Sprintf("Could not parse contact assignment ID '%s': %s", data.ID.ValueString(), err),
		)

		return

	}

	// Parse object ID

	objectID, err := utils.ParseID64(data.ObjectID.ValueString())

	if err != nil {

		resp.Diagnostics.AddError(

			"Invalid Object ID",

			fmt.Sprintf("Unable to parse object ID: %s", err),
		)

		return

	}

	// Parse contact ID

	contactID, err := utils.ParseID(data.Contact.ValueString())

	if err != nil {

		resp.Diagnostics.AddError(

			"Invalid Contact ID",

			fmt.Sprintf("Unable to parse contact ID: %s", err),
		)

		return

	}

	tflog.Debug(ctx, "Updating contact assignment", map[string]interface{}{

		"id": id,

		"object_type": data.ObjectType.ValueString(),

		"object_id": objectID,

		"contact_id": contactID,
	})

	// Build the API request

	briefContact := *netbox.NewBriefContactRequest("placeholder")

	assignmentRequest := netbox.NewWritableContactAssignmentRequest(

		data.ObjectType.ValueString(),

		objectID,

		briefContact,
	)

	// Use AdditionalProperties to pass the contact ID

	assignmentRequest.AdditionalProperties = make(map[string]interface{})

	assignmentRequest.AdditionalProperties["contact"] = int(contactID)

	// Set role if provided, otherwise unset it entirely

	if !data.Role.IsNull() && !data.Role.IsUnknown() && data.Role.ValueString() != "" {

		roleID, err := utils.ParseID(data.Role.ValueString())

		if err != nil {

			resp.Diagnostics.AddError(

				"Invalid Role ID",

				fmt.Sprintf("Unable to parse role ID: %s", err),
			)

			return

		}

		assignmentRequest.AdditionalProperties["role"] = int(roleID)

	} else {

		// Unset role entirely when not provided

		assignmentRequest.UnsetRole()

	}

	// Set priority if provided

	if !data.Priority.IsNull() && !data.Priority.IsUnknown() && data.Priority.ValueString() != "" {

		priority := netbox.BriefCircuitGroupAssignmentSerializerPriorityValue(data.Priority.ValueString())

		assignmentRequest.Priority = &priority

	}

	// Handle tags

	if !data.Tags.IsNull() {

		var tagModels []utils.TagModel

		diags := data.Tags.ElementsAs(ctx, &tagModels, false)

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {

			return

		}

		assignmentRequest.Tags = utils.TagsToNestedTagRequests(tagModels)

	}

	// Handle custom fields

	if !data.CustomFields.IsNull() {

		var customFieldModels []utils.CustomFieldModel

		diags := data.CustomFields.ElementsAs(ctx, &customFieldModels, false)

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {

			return

		}

		assignmentRequest.CustomFields = utils.CustomFieldModelsToMap(customFieldModels)

	}

	// Call the API to update the contact assignment

	assignment, httpResp, err := r.client.TenancyAPI.TenancyContactAssignmentsUpdate(ctx, id).
		WritableContactAssignmentRequest(*assignmentRequest).
		Execute()

	if err != nil {

		resp.Diagnostics.AddError(

			"Error updating contact assignment",

			utils.FormatAPIError("update contact assignment", err, httpResp),
		)

		return

	}

	// Map response to state, preserving computed display_name to avoid inconsistent result error
	displayNameBeforeUpdate := data.DisplayName
	r.mapResponseToState(ctx, assignment, &data, &resp.Diagnostics)
	data.DisplayName = displayNameBeforeUpdate

	tflog.Debug(ctx, "Updated contact assignment", map[string]interface{}{

		"id": assignment.GetId(),
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

func (r *ContactAssignmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var data ContactAssignmentResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	// Parse the ID

	id, err := utils.ParseID(data.ID.ValueString())

	if err != nil {

		resp.Diagnostics.AddError(

			"Invalid ID format",

			fmt.Sprintf("Could not parse contact assignment ID '%s': %s", data.ID.ValueString(), err),
		)

		return

	}

	tflog.Debug(ctx, "Deleting contact assignment", map[string]interface{}{

		"id": id,
	})

	httpResp, err := r.client.TenancyAPI.TenancyContactAssignmentsDestroy(ctx, id).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		if httpResp != nil && httpResp.StatusCode == 404 {

			tflog.Info(ctx, "Contact assignment already deleted", map[string]interface{}{

				"id": id,
			})

			return

		}

		resp.Diagnostics.AddError(

			"Error deleting contact assignment",

			utils.FormatAPIError("delete contact assignment", err, httpResp),
		)

		return

	}

	tflog.Debug(ctx, "Deleted contact assignment", map[string]interface{}{

		"id": id,
	})

}

func (r *ContactAssignmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

}

// mapResponseToState maps a ContactAssignment API response to the Terraform state model.

func (r *ContactAssignmentResource) mapResponseToState(ctx context.Context, assignment *netbox.ContactAssignment, data *ContactAssignmentResourceModel, diags *diag.Diagnostics) {

	data.ID = types.StringValue(fmt.Sprintf("%d", assignment.GetId()))

	data.ObjectType = types.StringValue(assignment.GetObjectType())

	data.ObjectID = types.StringValue(fmt.Sprintf("%d", assignment.GetObjectId()))

	// DisplayName
	if assignment.GetDisplay() != "" {
		data.DisplayName = types.StringValue(assignment.GetDisplay())
	} else {
		data.DisplayName = types.StringNull()
	}

	// Contact (required field) - preserve user's input format

	contact := assignment.GetContact()

	data.Contact = utils.UpdateReferenceAttribute(data.Contact, contact.GetName(), "", contact.GetId())

	// Role (optional field) - preserve user's input format

	if assignment.HasRole() && assignment.Role.Get() != nil {

		role := assignment.GetRole()

		data.Role = utils.UpdateReferenceAttribute(data.Role, role.GetName(), role.GetSlug(), role.GetId())

	} else {

		data.Role = types.StringNull()

	}

	// Priority (optional field)

	if assignment.HasPriority() && assignment.Priority != nil {

		priority := assignment.GetPriority()

		if priority.Value != nil && string(*priority.Value) != "" {

			data.Priority = types.StringValue(string(*priority.Value))

		} else {

			data.Priority = types.StringNull()

		}

	} else {

		data.Priority = types.StringNull()

	}

	// Tags

	if assignment.HasTags() && len(assignment.GetTags()) > 0 {

		tags := utils.NestedTagsToTagModels(assignment.GetTags())

		tagsValue, tagDiags := types.SetValueFrom(ctx, utils.GetTagsAttributeType().ElemType, tags)

		diags.Append(tagDiags...)

		data.Tags = tagsValue

	} else {

		data.Tags = types.SetNull(utils.GetTagsAttributeType().ElemType)

	}

	// Custom fields

	if assignment.HasCustomFields() && len(assignment.GetCustomFields()) > 0 {

		var existingModels []utils.CustomFieldModel

		if !data.CustomFields.IsNull() {

			cfDiags := data.CustomFields.ElementsAs(ctx, &existingModels, false)

			diags.Append(cfDiags...)

		}

		customFields := utils.MapToCustomFieldModels(assignment.GetCustomFields(), existingModels)

		customFieldsValue, cfDiags := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, customFields)

		diags.Append(cfDiags...)

		data.CustomFields = customFieldsValue

	} else {

		data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)

	}

}
