// Package datasources contains Terraform data source implementations for the Netbox provider.
package datasources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bab3l/go-netbox"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &FHRPGroupAssignmentDataSource{}

// NewFHRPGroupAssignmentDataSource returns a new data source implementing the FHRP group assignment data source.
func NewFHRPGroupAssignmentDataSource() datasource.DataSource {
	return &FHRPGroupAssignmentDataSource{}
}

// FHRPGroupAssignmentDataSource defines the data source implementation.
type FHRPGroupAssignmentDataSource struct {
	client *netbox.APIClient
}

// FHRPGroupAssignmentDataSourceModel describes the data source data model.
type FHRPGroupAssignmentDataSourceModel struct {
	ID            types.String `tfsdk:"id"`
	GroupID       types.String `tfsdk:"group_id"`
	GroupName     types.String `tfsdk:"group_name"`
	InterfaceType types.String `tfsdk:"interface_type"`
	InterfaceID   types.String `tfsdk:"interface_id"`
	Priority      types.Int64  `tfsdk:"priority"`
}

// Metadata returns the data source type name.
func (d *FHRPGroupAssignmentDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_fhrp_group_assignment"
}

// Schema defines the schema for the data source.
func (d *FHRPGroupAssignmentDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get information about an FHRP group assignment in NetBox.",

		Attributes: map[string]schema.Attribute{
			"id":             nbschema.DSIDAttribute("FHRP group assignment"),
			"group_id":       nbschema.DSComputedStringAttribute("ID of the FHRP group."),
			"group_name":     nbschema.DSComputedStringAttribute("Name of the FHRP group."),
			"interface_type": nbschema.DSComputedStringAttribute("Type of interface (dcim.interface or virtualization.vminterface)."),
			"interface_id":   nbschema.DSComputedStringAttribute("ID of the interface."),
			"priority":       nbschema.DSComputedInt64Attribute("Priority of this assignment."),
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *FHRPGroupAssignmentDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// Read reads the FHRP group assignment data.
func (d *FHRPGroupAssignmentDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data FHRPGroupAssignmentDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var assignment *netbox.FHRPGroupAssignment
	var httpResp *http.Response
	var err error

	// Lookup by ID only - assignments don't have a name
	if !data.ID.IsNull() && !data.ID.IsUnknown() {
		id, parseErr := utils.ParseID(data.ID.ValueString())
		if parseErr != nil {
			resp.Diagnostics.AddError(
				"Invalid ID",
				fmt.Sprintf("Could not parse ID '%s': %s", data.ID.ValueString(), parseErr),
			)
			return
		}

		tflog.Debug(ctx, "Reading FHRP group assignment by ID", map[string]interface{}{
			"id": id,
		})

		assignment, httpResp, err = d.client.IpamAPI.IpamFhrpGroupAssignmentsRetrieve(ctx, id).Execute()
	} else {
		resp.Diagnostics.AddError(
			"Missing Identifier",
			"'id' must be specified to look up an FHRP group assignment.",
		)
		return
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading FHRP group assignment",
			utils.FormatAPIError("read FHRP group assignment", err, httpResp),
		)
		return
	}

	// Map response to state
	d.mapResponseToState(ctx, assignment, &data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// mapResponseToState maps the API response to the Terraform state.
func (d *FHRPGroupAssignmentDataSource) mapResponseToState(ctx context.Context, assignment *netbox.FHRPGroupAssignment, data *FHRPGroupAssignmentDataSourceModel) {
	data.ID = types.StringValue(fmt.Sprintf("%d", assignment.GetId()))

	// Get group info - access Id field directly since GetId() is a pointer receiver method
	group := assignment.GetGroup()
	data.GroupID = types.StringValue(fmt.Sprintf("%d", group.Id))

	// BriefFHRPGroup has Display, not Name
	if group.Display != "" {
		data.GroupName = types.StringValue(group.Display)
	} else {
		data.GroupName = types.StringNull()
	}

	data.InterfaceType = types.StringValue(assignment.GetInterfaceType())
	data.InterfaceID = types.StringValue(fmt.Sprintf("%d", assignment.GetInterfaceId()))
	data.Priority = types.Int64Value(int64(assignment.GetPriority()))
}
