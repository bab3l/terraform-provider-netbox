// Package datasources contains Terraform data source implementations for the Netbox provider.
package datasources

import (
	"context"
	"fmt"

	"github.com/bab3l/go-netbox"
	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &ContactAssignmentDataSource{}

func NewContactAssignmentDataSource() datasource.DataSource {
	return &ContactAssignmentDataSource{}
}

// ContactAssignmentDataSource defines the data source implementation.
type ContactAssignmentDataSource struct {
	client *netbox.APIClient
}

// ContactAssignmentDataSourceModel describes the data source data model.
type ContactAssignmentDataSourceModel struct {
	ID           types.String `tfsdk:"id"`
	ObjectType   types.String `tfsdk:"object_type"`
	ObjectID     types.String `tfsdk:"object_id"`
	ContactID    types.String `tfsdk:"contact_id"`
	ContactName  types.String `tfsdk:"contact_name"`
	RoleID       types.String `tfsdk:"role_id"`
	RoleName     types.String `tfsdk:"role_name"`
	Priority     types.String `tfsdk:"priority"`
	PriorityName types.String `tfsdk:"priority_name"`
	Tags         types.Set    `tfsdk:"tags"`
}

func (d *ContactAssignmentDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_contact_assignment"
}

func (d *ContactAssignmentDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get information about a contact assignment in Netbox.",

		Attributes: map[string]schema.Attribute{
			"id":            nbschema.DSIDAttribute("contact assignment"),
			"object_type":   nbschema.DSComputedStringAttribute("Content type of the assigned object (e.g., dcim.site, dcim.device)."),
			"object_id":     nbschema.DSComputedStringAttribute("ID of the assigned object."),
			"contact_id":    nbschema.DSComputedStringAttribute("ID of the contact."),
			"contact_name":  nbschema.DSComputedStringAttribute("Name of the contact."),
			"role_id":       nbschema.DSComputedStringAttribute("ID of the contact role."),
			"role_name":     nbschema.DSComputedStringAttribute("Name of the contact role."),
			"priority":      nbschema.DSComputedStringAttribute("Priority value (primary, secondary, tertiary, inactive)."),
			"priority_name": nbschema.DSComputedStringAttribute("Display name for the priority."),
			"tags":          nbschema.DSTagsAttribute(),
		},
	}
}

func (d *ContactAssignmentDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ContactAssignmentDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ContactAssignmentDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// ID is required for lookup
	if data.ID.IsNull() || data.ID.IsUnknown() {
		resp.Diagnostics.AddError(
			"Missing required identifier",
			"The 'id' attribute must be specified to lookup a contact assignment.",
		)
		return
	}

	var idInt int32
	if _, parseErr := fmt.Sscanf(data.ID.ValueString(), "%d", &idInt); parseErr != nil {
		resp.Diagnostics.AddError(
			"Invalid ID format",
			fmt.Sprintf("Could not parse contact assignment ID '%s': %s", data.ID.ValueString(), parseErr.Error()),
		)
		return
	}

	tflog.Debug(ctx, "Looking up contact assignment by ID", map[string]interface{}{
		"id": idInt,
	})

	result, httpResp, err := d.client.TenancyAPI.TenancyContactAssignmentsRetrieve(ctx, idInt).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading contact assignment",
			utils.FormatAPIError("read contact assignment", err, httpResp),
		)
		return
	}

	// Map response to state
	d.mapResponseToState(ctx, result, &data, resp)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// mapResponseToState maps a ContactAssignment API response to the Terraform state model.
func (d *ContactAssignmentDataSource) mapResponseToState(ctx context.Context, assignment *netbox.ContactAssignment, data *ContactAssignmentDataSourceModel, resp *datasource.ReadResponse) {
	data.ID = types.StringValue(fmt.Sprintf("%d", assignment.GetId()))
	data.ObjectType = types.StringValue(assignment.GetObjectType())
	data.ObjectID = types.StringValue(fmt.Sprintf("%d", assignment.GetObjectId()))

	// Contact (required field)
	contact := assignment.GetContact()
	data.ContactID = types.StringValue(fmt.Sprintf("%d", contact.GetId()))
	data.ContactName = types.StringValue(contact.GetName())

	// Role (optional field)
	if assignment.HasRole() && assignment.Role.Get() != nil {
		role := assignment.GetRole()
		data.RoleID = types.StringValue(fmt.Sprintf("%d", role.GetId()))
		data.RoleName = types.StringValue(role.GetName())
	} else {
		data.RoleID = types.StringNull()
		data.RoleName = types.StringNull()
	}

	// Priority (optional field)
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
}
