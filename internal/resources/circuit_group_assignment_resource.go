// Package resources contains Terraform resource implementations for the Netbox provider.

package resources

import (
	"context"
	"fmt"
	"maps"

	"github.com/bab3l/go-netbox"
	"github.com/bab3l/terraform-provider-netbox/internal/netboxlookup"
	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &CircuitGroupAssignmentResource{}
var _ resource.ResourceWithImportState = &CircuitGroupAssignmentResource{}

func NewCircuitGroupAssignmentResource() resource.Resource {
	return &CircuitGroupAssignmentResource{}
}

// CircuitGroupAssignmentResource defines the resource implementation.
type CircuitGroupAssignmentResource struct {
	client *netbox.APIClient
}

// CircuitGroupAssignmentResourceModel describes the resource data model.
type CircuitGroupAssignmentResourceModel struct {
	ID           types.String `tfsdk:"id"`
	Group        types.String `tfsdk:"group_id"`
	Circuit      types.String `tfsdk:"circuit_id"`
	Priority     types.String `tfsdk:"priority"`
	Tags         types.Set    `tfsdk:"tags"`
	CustomFields types.Set    `tfsdk:"custom_fields"`
}

func (r *CircuitGroupAssignmentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_circuit_group_assignment"
}

func (r *CircuitGroupAssignmentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a circuit group assignment in Netbox. A circuit group assignment links a circuit to a circuit group with an optional priority.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The unique identifier of the circuit group assignment.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"group_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The ID of the circuit group to assign the circuit to.",
			},
			"circuit_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The ID of the circuit to assign to the group.",
			},
			"priority": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The priority of this circuit within the group. Valid values are: `primary`, `secondary`, `tertiary`, `inactive`.",
				Validators: []validator.String{
					stringvalidator.OneOf("primary", "secondary", "tertiary", "inactive", ""),
				},
			},
		},
	}

	// Add common metadata attributes (tags, custom_fields)
	maps.Copy(resp.Schema.Attributes, nbschema.CommonMetadataAttributes())
}

func (r *CircuitGroupAssignmentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *CircuitGroupAssignmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data CircuitGroupAssignmentResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Look up group by ID or slug
	circuitGroup, groupDiags := netboxlookup.LookupCircuitGroup(ctx, r.client, data.Group.ValueString())
	resp.Diagnostics.Append(groupDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Look up circuit by ID or CID
	circuit, circuitDiags := netboxlookup.LookupCircuit(ctx, r.client, data.Circuit.ValueString())
	resp.Diagnostics.Append(circuitDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Creating circuit group assignment", map[string]interface{}{
		"group_id":   circuitGroup.GetSlug(),
		"circuit_id": circuit.GetCid(),
	})

	// Convert CircuitGroup to BriefCircuitGroupRequest for API
	groupRequest := netbox.BriefCircuitGroupRequest{
		Name: circuitGroup.GetName(),
	}

	// Build the API request using the converted group and circuit
	assignmentRequest := netbox.NewWritableCircuitGroupAssignmentRequest(groupRequest, *circuit)

	// Set priority if provided
	if !data.Priority.IsNull() && !data.Priority.IsUnknown() && data.Priority.ValueString() != "" {
		priority := netbox.BriefCircuitGroupAssignmentSerializerPriorityValue(data.Priority.ValueString())
		assignmentRequest.Priority = &priority
	}

	// Handle tags
	utils.ApplyTags(ctx, assignmentRequest, data.Tags, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Call the API to create the circuit group assignment
	assignment, httpResp, err := r.client.CircuitsAPI.CircuitsCircuitGroupAssignmentsCreate(ctx).
		WritableCircuitGroupAssignmentRequest(*assignmentRequest).
		Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating circuit group assignment",
			utils.FormatAPIError("create circuit group assignment", err, httpResp),
		)
		return
	}

	// Map response to state
	r.mapResponseToState(ctx, assignment, &data, &resp.Diagnostics)
	tflog.Debug(ctx, "Created circuit group assignment", map[string]interface{}{
		"id": assignment.GetId(),
	})
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CircuitGroupAssignmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data CircuitGroupAssignmentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Parse the ID
	id, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid ID format",
			fmt.Sprintf("Could not parse circuit group assignment ID '%s': %s", data.ID.ValueString(), err),
		)
		return
	}
	tflog.Debug(ctx, "Reading circuit group assignment", map[string]interface{}{
		"id": id,
	})

	// Call the API to read the circuit group assignment
	assignment, httpResp, err := r.client.CircuitsAPI.CircuitsCircuitGroupAssignmentsRetrieve(ctx, id).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		// Check if resource was deleted outside of Terraform
		if httpResp != nil && httpResp.StatusCode == 404 {
			tflog.Debug(ctx, "Circuit group assignment not found, removing from state", map[string]interface{}{
				"id": id,
			})
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading circuit group assignment",
			utils.FormatAPIError("read circuit group assignment", err, httpResp),
		)
		return
	}

	// Map response to state
	r.mapResponseToState(ctx, assignment, &data, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CircuitGroupAssignmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data CircuitGroupAssignmentResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Parse the ID
	id, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid ID format",
			fmt.Sprintf("Could not parse circuit group assignment ID '%s': %s", data.ID.ValueString(), err),
		)
		return
	}

	// Look up group by ID or slug
	circuitGroup, groupDiags := netboxlookup.LookupCircuitGroup(ctx, r.client, data.Group.ValueString())
	resp.Diagnostics.Append(groupDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Look up circuit by ID or CID
	circuit, circuitDiags := netboxlookup.LookupCircuit(ctx, r.client, data.Circuit.ValueString())
	resp.Diagnostics.Append(circuitDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Updating circuit group assignment", map[string]interface{}{
		"id":         id,
		"group_id":   circuitGroup.GetSlug(),
		"circuit_id": circuit.GetCid(),
	})

	// Convert CircuitGroup to BriefCircuitGroupRequest for API
	groupRequest := netbox.BriefCircuitGroupRequest{
		Name: circuitGroup.GetName(),
	}

	// Build the API request
	assignmentRequest := netbox.NewWritableCircuitGroupAssignmentRequest(groupRequest, *circuit)

	// Set priority if provided
	if !data.Priority.IsNull() && !data.Priority.IsUnknown() && data.Priority.ValueString() != "" {
		priority := netbox.BriefCircuitGroupAssignmentSerializerPriorityValue(data.Priority.ValueString())
		assignmentRequest.Priority = &priority
	}

	// Handle tags
	utils.ApplyTags(ctx, assignmentRequest, data.Tags, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Call the API to update the circuit group assignment
	assignment, httpResp, err := r.client.CircuitsAPI.CircuitsCircuitGroupAssignmentsUpdate(ctx, id).
		WritableCircuitGroupAssignmentRequest(*assignmentRequest).
		Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating circuit group assignment",
			utils.FormatAPIError("update circuit group assignment", err, httpResp),
		)
		return
	}

	// Map response to state
	r.mapResponseToState(ctx, assignment, &data, &resp.Diagnostics)
	tflog.Debug(ctx, "Updated circuit group assignment", map[string]interface{}{
		"id": id,
	})
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CircuitGroupAssignmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CircuitGroupAssignmentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Parse the ID
	id, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid ID format",
			fmt.Sprintf("Could not parse circuit group assignment ID '%s': %s", data.ID.ValueString(), err),
		)
		return
	}
	tflog.Debug(ctx, "Deleting circuit group assignment", map[string]interface{}{
		"id": id,
	})

	// Call the API to delete the circuit group assignment
	httpResp, err := r.client.CircuitsAPI.CircuitsCircuitGroupAssignmentsDestroy(ctx, id).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		// Ignore 404 errors (resource already deleted)
		if httpResp != nil && httpResp.StatusCode == 404 {
			tflog.Debug(ctx, "Circuit group assignment already deleted", map[string]interface{}{
				"id": id,
			})
			return
		}
		resp.Diagnostics.AddError(
			"Error deleting circuit group assignment",
			utils.FormatAPIError("delete circuit group assignment", err, httpResp),
		)
		return
	}
	tflog.Debug(ctx, "Deleted circuit group assignment", map[string]interface{}{
		"id": id,
	})
}

func (r *CircuitGroupAssignmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := utils.ParseID(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid ID format",
			fmt.Sprintf("Could not parse circuit group assignment ID '%s': %s", req.ID, err),
		)
		return
	}
	tflog.Debug(ctx, "Importing circuit group assignment", map[string]interface{}{
		"id": id,
	})

	// Call the API to read the circuit group assignment
	assignment, httpResp, err := r.client.CircuitsAPI.CircuitsCircuitGroupAssignmentsRetrieve(ctx, id).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error importing circuit group assignment",
			utils.FormatAPIError("import circuit group assignment", err, httpResp),
		)
		return
	}

	// Create a new model and map the response
	var data CircuitGroupAssignmentResourceModel
	r.mapResponseToState(ctx, assignment, &data, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// mapResponseToState maps a CircuitGroupAssignment API response to the Terraform state model.
func (r *CircuitGroupAssignmentResource) mapResponseToState(ctx context.Context, assignment *netbox.CircuitGroupAssignment, data *CircuitGroupAssignmentResourceModel, diags *diag.Diagnostics) {
	data.ID = types.StringValue(fmt.Sprintf("%d", assignment.GetId()))

	// DisplayName
	if assignment.Display != "" {
	} else {
	}

	// Group (required field) - preserve user's input format
	group := assignment.GetGroup()
	data.Group = utils.UpdateReferenceAttribute(data.Group, group.GetName(), "", group.GetId())

	// Circuit (required field) - preserve user's input format
	circuit := assignment.GetCircuit()
	data.Circuit = utils.UpdateReferenceAttribute(data.Circuit, circuit.GetCid(), "", circuit.GetId())

	// Priority
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
		tagsValue, d := types.SetValueFrom(ctx, utils.GetTagsAttributeType().ElemType, tags)
		diags.Append(d...)
		data.Tags = tagsValue
	} else {
		data.Tags = types.SetNull(utils.GetTagsAttributeType().ElemType)
	}

	// Custom fields - circuit group assignments don't have custom fields in the response
	data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)
}
