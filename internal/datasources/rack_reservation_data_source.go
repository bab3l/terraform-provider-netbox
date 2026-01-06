// Package datasources provides Terraform data source implementations for NetBox objects.

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
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ datasource.DataSource              = &RackReservationDataSource{}
	_ datasource.DataSourceWithConfigure = &RackReservationDataSource{}
)

// NewRackReservationDataSource returns a new data source implementing the rack reservation data source.
func NewRackReservationDataSource() datasource.DataSource {
	return &RackReservationDataSource{}
}

// RackReservationDataSource defines the data source implementation.
type RackReservationDataSource struct {
	client *netbox.APIClient
}

// RackReservationDataSourceModel describes the data source data model.
type RackReservationDataSourceModel struct {
	ID           types.String `tfsdk:"id"`
	Rack         types.String `tfsdk:"rack"`
	RackID       types.String `tfsdk:"rack_id"`
	Units        types.Set    `tfsdk:"units"`
	User         types.String `tfsdk:"user"`
	UserID       types.String `tfsdk:"user_id"`
	Tenant       types.String `tfsdk:"tenant"`
	TenantID     types.String `tfsdk:"tenant_id"`
	Description  types.String `tfsdk:"description"`
	Comments     types.String `tfsdk:"comments"`
	Tags         types.Set    `tfsdk:"tags"`
	CustomFields types.Set    `tfsdk:"custom_fields"`
	DisplayName  types.String `tfsdk:"display_name"`
}

// Metadata returns the data source type name.
func (d *RackReservationDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_rack_reservation"
}

// Schema defines the schema for the data source.
func (d *RackReservationDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get information about a rack reservation in NetBox.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique numeric ID of the rack reservation.",
				Required:            true,
			},
			"rack": schema.StringAttribute{
				MarkdownDescription: "The name of the rack.",
				Computed:            true,
			},
			"rack_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the rack.",
				Computed:            true,
			},
			"units": schema.SetAttribute{
				MarkdownDescription: "The rack units (U positions) reserved.",
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"user": schema.StringAttribute{
				MarkdownDescription: "The username of the user who owns this reservation.",
				Computed:            true,
			},
			"user_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the user who owns this reservation.",
				Computed:            true,
			},
			"tenant": schema.StringAttribute{
				MarkdownDescription: "The name of the tenant associated with this reservation.",
				Computed:            true,
			},
			"tenant_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the tenant associated with this reservation.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "A description of the reservation purpose.",
				Computed:            true,
			},
			"comments": schema.StringAttribute{
				MarkdownDescription: "Additional comments about the reservation.",
				Computed:            true,
			},
			"tags":          nbschema.DSTagsAttribute(),
			"custom_fields": nbschema.DSCustomFieldsAttribute(),
			"display_name":  nbschema.DSComputedStringAttribute("The display name of the rack reservation."),
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *RackReservationDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// Read reads the data source.
func (d *RackReservationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data RackReservationDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Parse ID
	var id int32
	_, err := fmt.Sscanf(data.ID.ValueString(), "%d", &id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid ID",
			fmt.Sprintf("Could not parse ID '%s': %s", data.ID.ValueString(), err.Error()),
		)
		return
	}
	tflog.Debug(ctx, "Reading rack reservation", map[string]interface{}{"id": id})

	// Read from API
	result, httpResp, err := d.client.DcimAPI.DcimRackReservationsRetrieve(ctx, id).Execute()
	defer utils.CloseResponseBody(httpResp)
	if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
		resp.Diagnostics.AddError(
			"Rack Reservation Not Found",
			fmt.Sprintf("No rack reservation found with ID: %d", id),
		)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading rack reservation",
			utils.FormatAPIError(fmt.Sprintf("read rack reservation ID %d", id), err, httpResp),
		)
		return
	}

	// Map response to state
	d.mapToState(ctx, result, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// mapToState maps the API response to the Terraform state.
func (d *RackReservationDataSource) mapToState(ctx context.Context, result *netbox.RackReservation, data *RackReservationDataSourceModel, diags *diag.Diagnostics) {
	data.ID = types.StringValue(fmt.Sprintf("%d", result.GetId()))

	// Map rack (required field)
	rack := result.GetRack()
	data.Rack = types.StringValue(rack.GetName())
	data.RackID = types.StringValue(fmt.Sprintf("%d", rack.GetId()))

	// Map units
	apiUnits := result.GetUnits()
	units := make([]int64, len(apiUnits))
	for i, u := range apiUnits {
		units[i] = int64(u)
	}
	unitsValue, unitDiags := types.SetValueFrom(ctx, types.Int64Type, units)
	diags.Append(unitDiags...)
	data.Units = unitsValue

	// Map user (required field)
	user := result.GetUser()
	data.User = types.StringValue(user.GetUsername())
	data.UserID = types.StringValue(fmt.Sprintf("%d", user.GetId()))

	// Map tenant
	if result.HasTenant() && result.GetTenant().Id != 0 {
		tenant := result.GetTenant()
		data.Tenant = types.StringValue(tenant.GetName())
		data.TenantID = types.StringValue(fmt.Sprintf("%d", tenant.GetId()))
	} else {
		data.Tenant = types.StringNull()
		data.TenantID = types.StringNull()
	}

	// Map description
	data.Description = types.StringValue(result.GetDescription())

	// Map comments
	if result.HasComments() && result.GetComments() != "" {
		data.Comments = types.StringValue(result.GetComments())
	} else {
		data.Comments = types.StringNull()
	}

	// Map tags
	if result.HasTags() && len(result.GetTags()) > 0 {
		tags := utils.NestedTagsToTagModels(result.GetTags())
		tagsValue, tagDiags := types.SetValueFrom(ctx, utils.GetTagsAttributeType().ElemType, tags)
		diags.Append(tagDiags...)
		data.Tags = tagsValue
	} else {
		data.Tags = types.SetNull(utils.GetTagsAttributeType().ElemType)
	}

	// Map custom fields
	if result.HasCustomFields() && len(result.GetCustomFields()) > 0 {
		customFields := utils.MapToCustomFieldModels(result.GetCustomFields(), nil)
		cfValue, cfDiags := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, customFields)
		diags.Append(cfDiags...)
		data.CustomFields = cfValue
	} else {
		data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)
	}

	// Map display_name
	if result.GetDisplay() != "" {
		data.DisplayName = types.StringValue(result.GetDisplay())
	} else {
		data.DisplayName = types.StringNull()
	}
}
