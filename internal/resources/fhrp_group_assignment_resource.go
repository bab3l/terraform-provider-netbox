// Package resources contains Terraform resource implementations for the Netbox provider.

package resources

import (
	"context"
	"fmt"

	"github.com/bab3l/go-netbox"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/bab3l/terraform-provider-netbox/internal/validators"
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

var (
	_ resource.Resource = &FHRPGroupAssignmentResource{}

	_ resource.ResourceWithConfigure = &FHRPGroupAssignmentResource{}

	_ resource.ResourceWithImportState = &FHRPGroupAssignmentResource{}
)

// NewFHRPGroupAssignmentResource returns a new resource implementing the FHRP group assignment resource.

func NewFHRPGroupAssignmentResource() resource.Resource {
	return &FHRPGroupAssignmentResource{}
}

// FHRPGroupAssignmentResource defines the resource implementation.

type FHRPGroupAssignmentResource struct {
	client *netbox.APIClient
}

// FHRPGroupAssignmentResourceModel describes the resource data model.

type FHRPGroupAssignmentResourceModel struct {
	ID types.String `tfsdk:"id"`

	GroupID types.String `tfsdk:"group_id"`

	InterfaceType types.String `tfsdk:"interface_type"`

	InterfaceID types.String `tfsdk:"interface_id"`

	Priority types.Int64 `tfsdk:"priority"`
}

// Metadata returns the resource type name.

func (r *FHRPGroupAssignmentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_fhrp_group_assignment"
}

// Schema defines the schema for the resource.

func (r *FHRPGroupAssignmentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an FHRP group assignment in NetBox. FHRP group assignments link FHRP groups to interfaces.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique numeric ID of the FHRP group assignment.",

				Computed: true,

				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"group_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the FHRP group to assign.",

				Required: true,

				Validators: []validator.String{
					stringvalidator.RegexMatches(

						validators.IntegerRegex(),

						"must be a valid integer ID",
					),
				},
			},

			"interface_type": schema.StringAttribute{
				MarkdownDescription: "The type of interface. Valid values: `dcim.interface`, `virtualization.vminterface`.",

				Required: true,

				Validators: []validator.String{
					stringvalidator.OneOf("dcim.interface", "virtualization.vminterface"),
				},
			},

			"interface_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the interface to assign the FHRP group to.",

				Required: true,

				Validators: []validator.String{
					stringvalidator.RegexMatches(

						validators.IntegerRegex(),

						"must be a valid integer ID",
					),
				},
			},

			"priority": schema.Int64Attribute{
				MarkdownDescription: "The priority of this assignment (0-255).",

				Required: true,
			},
		},
	}
}

// Configure adds the provider configured client to the resource.

func (r *FHRPGroupAssignmentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create creates a new FHRP group assignment.

func (r *FHRPGroupAssignmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data FHRPGroupAssignmentResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Parse group ID

	groupID, err := utils.ParseID(data.GroupID.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(

			"Invalid Group ID",

			fmt.Sprintf("Could not parse group ID '%s': %s", data.GroupID.ValueString(), err),
		)

		return
	}

	// Parse interface ID

	interfaceID, err := utils.ParseID(data.InterfaceID.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(

			"Invalid Interface ID",

			fmt.Sprintf("Could not parse interface ID '%s': %s", data.InterfaceID.ValueString(), err),
		)

		return
	}

	tflog.Debug(ctx, "Creating FHRP group assignment", map[string]interface{}{
		"group_id": groupID,

		"interface_type": data.InterfaceType.ValueString(),

		"interface_id": interfaceID,

		"priority": data.Priority.ValueInt64(),
	})

	// Build the API request - we need to look up the FHRP group first to get its details

	fhrpGroup, httpResp, err := r.client.IpamAPI.IpamFhrpGroupsRetrieve(ctx, groupID).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {
		resp.Diagnostics.AddError(

			"Error looking up FHRP group",

			utils.FormatAPIError("lookup FHRP group", err, httpResp),
		)

		return
	}

	// Create BriefFHRPGroupRequest from the looked up group
	// Protocol is already a BriefFHRPGroupProtocol string type

	briefGroup := netbox.NewBriefFHRPGroupRequest(fhrpGroup.GetProtocol(), fhrpGroup.GetGroupId())

	priority, err := utils.SafeInt32FromValue(data.Priority)

	if err != nil {
		resp.Diagnostics.AddError("Invalid value", fmt.Sprintf("Priority value overflow: %s", err))

		return
	}

	assignmentRequest := netbox.NewFHRPGroupAssignmentRequest(

		*briefGroup,

		data.InterfaceType.ValueString(),

		int64(interfaceID),

		priority,
	)

	// Use AdditionalProperties to pass the group ID

	assignmentRequest.AdditionalProperties = make(map[string]interface{})

	assignmentRequest.AdditionalProperties["group"] = int(groupID)

	// Call the API

	assignment, httpResp, err := r.client.IpamAPI.IpamFhrpGroupAssignmentsCreate(ctx).
		FHRPGroupAssignmentRequest(*assignmentRequest).
		Execute()

	if err != nil {
		resp.Diagnostics.AddError(

			"Error creating FHRP group assignment",

			utils.FormatAPIError("create FHRP group assignment", err, httpResp),
		)

		return
	}

	// Map response to state

	r.mapResponseToState(ctx, assignment, &data, &resp.Diagnostics)

	tflog.Debug(ctx, "Created FHRP group assignment", map[string]interface{}{
		"id": assignment.GetId(),
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read reads the FHRP group assignment.

func (r *FHRPGroupAssignmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data FHRPGroupAssignmentResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Parse ID

	id, err := utils.ParseID(data.ID.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(

			"Invalid ID format",

			fmt.Sprintf("Could not parse FHRP group assignment ID '%s': %s", data.ID.ValueString(), err),
		)

		return
	}

	tflog.Debug(ctx, "Reading FHRP group assignment", map[string]interface{}{
		"id": id,
	})

	// Call the API

	assignment, httpResp, err := r.client.IpamAPI.IpamFhrpGroupAssignmentsRetrieve(ctx, id).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			resp.State.RemoveResource(ctx)

			return
		}

		resp.Diagnostics.AddError(

			"Error reading FHRP group assignment",

			utils.FormatAPIError("read FHRP group assignment", err, httpResp),
		)

		return
	}

	// Map response to state

	r.mapResponseToState(ctx, assignment, &data, &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates the FHRP group assignment.

func (r *FHRPGroupAssignmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data FHRPGroupAssignmentResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Parse ID

	id, err := utils.ParseID(data.ID.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(

			"Invalid ID format",

			fmt.Sprintf("Could not parse FHRP group assignment ID '%s': %s", data.ID.ValueString(), err),
		)

		return
	}

	// Parse group ID

	groupID, err := utils.ParseID(data.GroupID.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(

			"Invalid Group ID",

			fmt.Sprintf("Could not parse group ID '%s': %s", data.GroupID.ValueString(), err),
		)

		return
	}

	// Parse interface ID

	interfaceID, err := utils.ParseID(data.InterfaceID.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(

			"Invalid Interface ID",

			fmt.Sprintf("Could not parse interface ID '%s': %s", data.InterfaceID.ValueString(), err),
		)

		return
	}

	tflog.Debug(ctx, "Updating FHRP group assignment", map[string]interface{}{
		"id": id,

		"group_id": groupID,

		"interface_type": data.InterfaceType.ValueString(),

		"interface_id": interfaceID,

		"priority": data.Priority.ValueInt64(),
	})

	// Build the API request - we need to look up the FHRP group first to get its details

	fhrpGroup, httpResp, err := r.client.IpamAPI.IpamFhrpGroupsRetrieve(ctx, groupID).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {
		resp.Diagnostics.AddError(

			"Error looking up FHRP group",

			utils.FormatAPIError("lookup FHRP group", err, httpResp),
		)

		return
	}

	// Create BriefFHRPGroupRequest from the looked up group
	// Protocol is already a BriefFHRPGroupProtocol string type

	briefGroup := netbox.NewBriefFHRPGroupRequest(fhrpGroup.GetProtocol(), fhrpGroup.GetGroupId())

	priority, err := utils.SafeInt32FromValue(data.Priority)

	if err != nil {
		resp.Diagnostics.AddError("Invalid value", fmt.Sprintf("Priority value overflow: %s", err))

		return
	}

	assignmentRequest := netbox.NewFHRPGroupAssignmentRequest(

		*briefGroup,

		data.InterfaceType.ValueString(),

		int64(interfaceID),

		priority,
	)

	// Use AdditionalProperties to pass the group ID

	assignmentRequest.AdditionalProperties = make(map[string]interface{})

	assignmentRequest.AdditionalProperties["group"] = int(groupID)

	// Call the API

	assignment, httpResp, err := r.client.IpamAPI.IpamFhrpGroupAssignmentsUpdate(ctx, id).
		FHRPGroupAssignmentRequest(*assignmentRequest).
		Execute()

	if err != nil {
		resp.Diagnostics.AddError(

			"Error updating FHRP group assignment",

			utils.FormatAPIError("update FHRP group assignment", err, httpResp),
		)

		return
	}

	// Map response to state

	r.mapResponseToState(ctx, assignment, &data, &resp.Diagnostics)

	tflog.Debug(ctx, "Updated FHRP group assignment", map[string]interface{}{
		"id": assignment.GetId(),
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete deletes the FHRP group assignment.

func (r *FHRPGroupAssignmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data FHRPGroupAssignmentResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Parse ID

	id, err := utils.ParseID(data.ID.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(

			"Invalid ID format",

			fmt.Sprintf("Could not parse FHRP group assignment ID '%s': %s", data.ID.ValueString(), err),
		)

		return
	}

	tflog.Debug(ctx, "Deleting FHRP group assignment", map[string]interface{}{
		"id": id,
	})

	httpResp, err := r.client.IpamAPI.IpamFhrpGroupAssignmentsDestroy(ctx, id).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			return
		}

		resp.Diagnostics.AddError(

			"Error deleting FHRP group assignment",

			utils.FormatAPIError("delete FHRP group assignment", err, httpResp),
		)

		return
	}

	tflog.Debug(ctx, "Deleted FHRP group assignment", map[string]interface{}{
		"id": id,
	})
}

// ImportState imports an existing FHRP group assignment.

func (r *FHRPGroupAssignmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// mapResponseToState maps the API response to the Terraform state.

func (r *FHRPGroupAssignmentResource) mapResponseToState(ctx context.Context, assignment *netbox.FHRPGroupAssignment, data *FHRPGroupAssignmentResourceModel, diags *diag.Diagnostics) {
	data.ID = types.StringValue(fmt.Sprintf("%d", assignment.GetId()))

	group := assignment.GetGroup()

	data.GroupID = types.StringValue(fmt.Sprintf("%d", group.Id))

	data.InterfaceType = types.StringValue(assignment.GetInterfaceType())

	data.InterfaceID = types.StringValue(fmt.Sprintf("%d", assignment.GetInterfaceId()))

	data.Priority = types.Int64Value(int64(assignment.GetPriority()))

	// Map display_name (computed field, always set a value)
	if assignment.Display != "" {
	} else {
	}
}
