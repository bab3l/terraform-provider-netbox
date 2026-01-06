// Package datasources contains Terraform data source implementations for the Netbox provider.

package datasources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bab3l/go-netbox"
	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &CircuitGroupAssignmentDataSource{}

func NewCircuitGroupAssignmentDataSource() datasource.DataSource {
	return &CircuitGroupAssignmentDataSource{}
}

// CircuitGroupAssignmentDataSource defines the data source implementation.
type CircuitGroupAssignmentDataSource struct {
	client *netbox.APIClient
}

// CircuitGroupAssignmentDataSourceModel describes the data source data model.
type CircuitGroupAssignmentDataSourceModel struct {
	ID           types.String `tfsdk:"id"`
	GroupID      types.String `tfsdk:"group_id"`
	GroupName    types.String `tfsdk:"group_name"`
	CircuitID    types.String `tfsdk:"circuit_id"`
	CircuitCID   types.String `tfsdk:"circuit_cid"`
	Priority     types.String `tfsdk:"priority"`
	PriorityName types.String `tfsdk:"priority_name"`
	DisplayName  types.String `tfsdk:"display_name"`
	Tags         types.Set    `tfsdk:"tags"`
}

func (d *CircuitGroupAssignmentDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_circuit_group_assignment"
}

func (d *CircuitGroupAssignmentDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get information about a circuit group assignment in Netbox.",
		Attributes: map[string]schema.Attribute{
			"id":            nbschema.DSIDAttribute("circuit group assignment"),
			"group_id":      nbschema.DSComputedStringAttribute("ID of the circuit group."),
			"group_name":    nbschema.DSComputedStringAttribute("Name of the circuit group."),
			"circuit_id":    nbschema.DSComputedStringAttribute("ID of the circuit."),
			"circuit_cid":   nbschema.DSComputedStringAttribute("Circuit ID (CID) of the circuit."),
			"priority":      nbschema.DSComputedStringAttribute("Priority value (primary, secondary, tertiary, inactive)."),
			"priority_name": nbschema.DSComputedStringAttribute("Display name for the priority."),
			"display_name":  nbschema.DSComputedStringAttribute("The display name of the circuit group assignment."),
			"tags":          nbschema.DSTagsAttribute(),
		},
	}
}

func (d *CircuitGroupAssignmentDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*netbox.APIClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *netbox.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	d.client = client
}

func (d *CircuitGroupAssignmentDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CircuitGroupAssignmentDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// ID is required for lookup
	if data.ID.IsNull() || data.ID.IsUnknown() {
		resp.Diagnostics.AddError(
			"Missing required identifier",
			"The 'id' attribute must be specified to lookup a circuit group assignment.",
		)
		return
	}
	var idInt int32
	if _, parseErr := fmt.Sscanf(data.ID.ValueString(), "%d", &idInt); parseErr != nil {
		resp.Diagnostics.AddError(
			"Invalid ID format",
			fmt.Sprintf("Could not parse circuit group assignment ID '%s': %s", data.ID.ValueString(), parseErr.Error()),
		)
		return
	}
	tflog.Debug(ctx, "Looking up circuit group assignment by ID", map[string]interface{}{
		"id": idInt,
	})
	result, httpResp, err := d.client.CircuitsAPI.CircuitsCircuitGroupAssignmentsRetrieve(ctx, idInt).Execute()
	defer utils.CloseResponseBody(httpResp)
	if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
		resp.Diagnostics.AddError(
			"Circuit Group Assignment Not Found",
			fmt.Sprintf("No circuit group assignment found with ID: %d", idInt),
		)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading circuit group assignment",
			utils.FormatAPIError("read circuit group assignment", err, httpResp),
		)
		return
	}

	// Map response to state
	d.mapResponseToState(ctx, result, &data, resp)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// mapResponseToState maps a CircuitGroupAssignment API response to the Terraform state model.
func (d *CircuitGroupAssignmentDataSource) mapResponseToState(ctx context.Context, assignment *netbox.CircuitGroupAssignment, data *CircuitGroupAssignmentDataSourceModel, resp *datasource.ReadResponse) {
	data.ID = types.StringValue(fmt.Sprintf("%d", assignment.GetId()))

	// Group (required field)
	group := assignment.GetGroup()
	data.GroupID = types.StringValue(fmt.Sprintf("%d", group.GetId()))
	data.GroupName = types.StringValue(group.GetName())

	// Circuit (required field)
	circuit := assignment.GetCircuit()
	data.CircuitID = types.StringValue(fmt.Sprintf("%d", circuit.GetId()))
	data.CircuitCID = types.StringValue(circuit.GetCid())

	// Priority
	if assignment.HasPriority() && assignment.Priority != nil {
		priority := assignment.GetPriority()
		if priority.Value != nil && string(*priority.Value) != "" {
			data.Priority = types.StringValue(string(*priority.Value))
		} else {
			data.Priority = types.StringNull()
		}
		if priority.Label != nil && string(*priority.Label) != "" {
			data.PriorityName = types.StringValue(string(*priority.Label))
		} else {
			data.PriorityName = types.StringNull()
		}
	} else {
		data.Priority = types.StringNull()
		data.PriorityName = types.StringNull()
	}

	// Tags
	if assignment.HasTags() && len(assignment.GetTags()) > 0 {
		tags := utils.NestedTagsToTagModels(assignment.GetTags())
		tagsValue, diags := types.SetValueFrom(ctx, utils.GetTagsAttributeType().ElemType, tags)
		resp.Diagnostics.Append(diags...)
		data.Tags = tagsValue
	} else {
		data.Tags = types.SetNull(utils.GetTagsAttributeType().ElemType)
	}

	// Map display name
	if assignment.GetDisplay() != "" {
		data.DisplayName = types.StringValue(assignment.GetDisplay())
	} else {
		data.DisplayName = types.StringNull()
	}
}
