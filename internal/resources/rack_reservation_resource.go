// Package resources provides Terraform resource implementations for NetBox objects.

package resources

import (
	"context"
	"fmt"
	"maps"

	"github.com/bab3l/go-netbox"
	lookup "github.com/bab3l/terraform-provider-netbox/internal/netboxlookup"
	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.

var (
	_ resource.Resource = &RackReservationResource{}

	_ resource.ResourceWithConfigure = &RackReservationResource{}

	_ resource.ResourceWithImportState = &RackReservationResource{}
)

// NewRackReservationResource returns a new resource implementing the rack reservation resource.

func NewRackReservationResource() resource.Resource {
	return &RackReservationResource{}
}

// RackReservationResource defines the resource implementation.

type RackReservationResource struct {
	client *netbox.APIClient
}

// RackReservationResourceModel describes the resource data model.

type RackReservationResourceModel struct {
	ID types.String `tfsdk:"id"`

	Rack types.String `tfsdk:"rack"`

	Units types.Set `tfsdk:"units"`

	User types.String `tfsdk:"user"`

	Tenant types.String `tfsdk:"tenant"`

	Description types.String `tfsdk:"description"`

	Comments types.String `tfsdk:"comments"`

	Tags types.Set `tfsdk:"tags"`

	CustomFields types.Set `tfsdk:"custom_fields"`
}

// Metadata returns the resource type name.

func (r *RackReservationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_rack_reservation"
}

// Schema defines the schema for the resource.

func (r *RackReservationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a rack reservation in NetBox. Rack reservations allow you to designate specific units within a rack for a particular purpose or user.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique numeric ID of the rack reservation.",

				Computed: true,

				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"rack": schema.StringAttribute{
				MarkdownDescription: "The rack containing the reserved units (ID or name).",

				Required: true,
			},

			"units": schema.SetAttribute{
				MarkdownDescription: "The rack units (U positions) to reserve. Must be a set of integers.",

				Required: true,

				ElementType: types.Int64Type,
			},

			"user": schema.StringAttribute{
				MarkdownDescription: "The user who owns this reservation (ID or username).",

				Required: true,
			},

			"tenant": schema.StringAttribute{
				MarkdownDescription: "The tenant associated with this reservation (ID or slug).",

				Optional: true,
			},

			"description": schema.StringAttribute{
				MarkdownDescription: "A description of the reservation purpose.",

				Required: true,
			},

			"comments": nbschema.CommentsAttribute("rack reservation"),
		},
	}

	// Add common metadata attributes (tags, custom_fields)
	maps.Copy(resp.Schema.Attributes, nbschema.CommonMetadataAttributes())
}

func (r *RackReservationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create creates the resource.

func (r *RackReservationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data RackReservationResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Lookup rack

	rack, diags := lookup.LookupRack(ctx, r.client, data.Rack.ValueString())

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Lookup user

	user, diags := lookup.LookupUser(ctx, r.client, data.User.ValueString())

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Convert units to []int32

	var unitsInt64 []int64

	diags = data.Units.ElementsAs(ctx, &unitsInt64, false)

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	units := make([]int32, len(unitsInt64))

	for i, u := range unitsInt64 {
		val, err := utils.SafeInt32(u)

		if err != nil {
			resp.Diagnostics.AddError("Invalid value", fmt.Sprintf("Units value overflow: %s", err))

			return
		}

		units[i] = val
	}

	// Build request

	apiReq := netbox.NewRackReservationRequest(*rack, units, *user, data.Description.ValueString())

	// Set optional fields

	if !data.Tenant.IsNull() && !data.Tenant.IsUnknown() {
		tenant, tenantDiags := lookup.LookupTenant(ctx, r.client, data.Tenant.ValueString())

		resp.Diagnostics.Append(tenantDiags...)

		if resp.Diagnostics.HasError() {
			return
		}

		apiReq.SetTenant(*tenant)
	}

	// Apply optional fields (comments, tags, custom_fields)

	utils.ApplyComments(apiReq, data.Comments)

	utils.ApplyTags(ctx, apiReq, data.Tags, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	utils.ApplyCustomFields(ctx, apiReq, data.CustomFields, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating rack reservation", map[string]interface{}{
		"rack": data.Rack.ValueString(),

		"units": units,

		"user": data.User.ValueString(),

		"description": data.Description.ValueString(),
	})

	// Create the resource

	result, httpResp, err := r.client.DcimAPI.DcimRackReservationsCreate(ctx).RackReservationRequest(*apiReq).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {
		resp.Diagnostics.AddError(

			"Error creating rack reservation",

			utils.FormatAPIError("create rack reservation", err, httpResp),
		)

		return
	}

	// Map response to state

	r.mapToState(ctx, result, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read reads the resource.

func (r *RackReservationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data RackReservationResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Parse ID

	var id int32

	_, err := fmt.Sscanf(data.ID.ValueString(), "%d", &id)

	if err != nil {
		resp.Diagnostics.AddError(

			"Error parsing rack reservation ID",

			fmt.Sprintf("Could not parse ID '%s': %s", data.ID.ValueString(), err.Error()),
		)

		return
	}

	// Read from API

	result, httpResp, err := r.client.DcimAPI.DcimRackReservationsRetrieve(ctx, id).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			resp.State.RemoveResource(ctx)

			return
		}

		resp.Diagnostics.AddError(

			"Error reading rack reservation",

			utils.FormatAPIError(fmt.Sprintf("read rack reservation ID %d", id), err, httpResp),
		)

		return
	}

	// Map response to state

	r.mapToState(ctx, result, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates the resource.

func (r *RackReservationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data RackReservationResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Parse ID

	var id int32

	_, err := fmt.Sscanf(data.ID.ValueString(), "%d", &id)

	if err != nil {
		resp.Diagnostics.AddError(

			"Error parsing rack reservation ID",

			fmt.Sprintf("Could not parse ID '%s': %s", data.ID.ValueString(), err.Error()),
		)

		return
	}

	// Lookup rack

	rack, diags := lookup.LookupRack(ctx, r.client, data.Rack.ValueString())

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Lookup user

	user, diags := lookup.LookupUser(ctx, r.client, data.User.ValueString())

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Convert units to []int32

	var unitsInt64 []int64

	diags = data.Units.ElementsAs(ctx, &unitsInt64, false)

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	units := make([]int32, len(unitsInt64))

	for i, u := range unitsInt64 {
		val, err := utils.SafeInt32(u)

		if err != nil {
			resp.Diagnostics.AddError("Invalid value", fmt.Sprintf("Units value overflow: %s", err))

			return
		}

		units[i] = val
	}

	// Build request

	apiReq := netbox.NewRackReservationRequest(*rack, units, *user, data.Description.ValueString())

	// Set optional fields

	if !data.Tenant.IsNull() && !data.Tenant.IsUnknown() {
		tenant, tenantDiags := lookup.LookupTenant(ctx, r.client, data.Tenant.ValueString())

		resp.Diagnostics.Append(tenantDiags...)

		if resp.Diagnostics.HasError() {
			return
		}

		apiReq.SetTenant(*tenant)
	}

	// Apply optional fields (comments, tags, custom_fields)

	utils.ApplyComments(apiReq, data.Comments)

	utils.ApplyTags(ctx, apiReq, data.Tags, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	utils.ApplyCustomFields(ctx, apiReq, data.CustomFields, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Updating rack reservation", map[string]interface{}{
		"id": id,

		"rack": data.Rack.ValueString(),

		"units": units,

		"user": data.User.ValueString(),

		"description": data.Description.ValueString(),
	})

	// Update the resource

	result, httpResp, err := r.client.DcimAPI.DcimRackReservationsUpdate(ctx, id).RackReservationRequest(*apiReq).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {
		resp.Diagnostics.AddError(

			"Error updating rack reservation",

			utils.FormatAPIError(fmt.Sprintf("update rack reservation ID %d", id), err, httpResp),
		)

		return
	}

	// Map response to state

	r.mapToState(ctx, result, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete deletes the resource.

func (r *RackReservationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data RackReservationResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Parse ID

	var id int32

	_, err := fmt.Sscanf(data.ID.ValueString(), "%d", &id)

	if err != nil {
		resp.Diagnostics.AddError(

			"Error parsing rack reservation ID",

			fmt.Sprintf("Could not parse ID '%s': %s", data.ID.ValueString(), err.Error()),
		)

		return
	}

	tflog.Debug(ctx, "Deleting rack reservation", map[string]interface{}{"id": id})

	// Delete the resource

	httpResp, err := r.client.DcimAPI.DcimRackReservationsDestroy(ctx, id).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			return
		}

		resp.Diagnostics.AddError(

			"Error deleting rack reservation",

			utils.FormatAPIError(fmt.Sprintf("delete rack reservation ID %d", id), err, httpResp),
		)

		return
	}
}

// ImportState imports the resource state.

func (r *RackReservationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// mapToState maps the API response to the Terraform state.

func (r *RackReservationResource) mapToState(ctx context.Context, result *netbox.RackReservation, data *RackReservationResourceModel, diags *diag.Diagnostics) {
	data.ID = types.StringValue(fmt.Sprintf("%d", result.GetId()))

	// Map rack (required field) - preserve user's input format, but always set something for import
	rack := result.GetRack()
	rackValue := utils.UpdateReferenceAttribute(data.Rack, rack.GetName(), "", rack.GetId())

	// If null (happens during import), set the ID as fallback
	if rackValue.IsNull() {
		rackValue = types.StringValue(fmt.Sprintf("%d", rack.GetId()))
	}

	data.Rack = rackValue

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

	data.User = types.StringValue(fmt.Sprintf("%d", user.GetId()))

	// Map tenant

	if result.HasTenant() && result.GetTenant().Id != 0 {
		tenant := result.GetTenant()

		data.Tenant = types.StringValue(fmt.Sprintf("%d", tenant.GetId()))
	} else {
		data.Tenant = types.StringNull()
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
		var existingModels []utils.CustomFieldModel

		if !data.CustomFields.IsNull() {
			elemDiags := data.CustomFields.ElementsAs(ctx, &existingModels, false)

			diags.Append(elemDiags...)
		}

		customFields := utils.MapToCustomFieldModels(result.GetCustomFields(), existingModels)

		cfValue, cfDiags := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, customFields)

		diags.Append(cfDiags...)

		data.CustomFields = cfValue
	} else {
		data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)
	}

	// Map display_name
}
