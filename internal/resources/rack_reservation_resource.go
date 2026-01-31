// Package resources provides Terraform resource implementations for NetBox objects.

package resources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bab3l/go-netbox"
	lookup "github.com/bab3l/terraform-provider-netbox/internal/netboxlookup"
	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/attr"
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
	_ resource.Resource                = &RackReservationResource{}
	_ resource.ResourceWithConfigure   = &RackReservationResource{}
	_ resource.ResourceWithImportState = &RackReservationResource{}
	_ resource.ResourceWithIdentity    = &RackReservationResource{}
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
	ID           types.String `tfsdk:"id"`
	Rack         types.String `tfsdk:"rack"`
	Units        types.Set    `tfsdk:"units"`
	User         types.String `tfsdk:"user"`
	Tenant       types.String `tfsdk:"tenant"`
	Description  types.String `tfsdk:"description"`
	Comments     types.String `tfsdk:"comments"`
	Tags         types.Set    `tfsdk:"tags"`
	CustomFields types.Set    `tfsdk:"custom_fields"`
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
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"rack": nbschema.RequiredReferenceAttributeWithDiffSuppress(
				"rack",
				"The rack containing the reserved units (ID or name).",
			),
			"units": schema.SetAttribute{
				MarkdownDescription: "The rack units (U positions) to reserve. Must be a set of integers.",
				Required:            true,
				ElementType:         types.Int64Type,
			},
			"user": nbschema.RequiredReferenceAttributeWithDiffSuppress(
				"user",
				"The user who owns this reservation (ID or username).",
			),
			"tenant": nbschema.ReferenceAttributeWithDiffSuppress(
				"tenant",
				"The tenant associated with this reservation (ID or slug).",
			),
			"description": schema.StringAttribute{
				MarkdownDescription: "A description of the reservation purpose.",
				Required:            true,
			},
			"comments": nbschema.CommentsAttribute("rack reservation"),
		},
	}

	// Add tags and custom fields attributes
	resp.Schema.Attributes["tags"] = nbschema.TagsSlugAttribute()
	resp.Schema.Attributes["custom_fields"] = nbschema.CustomFieldsAttribute()
}

func (r *RackReservationResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = nbschema.ImportIdentityWithCustomFieldsSchema()
}

// Configure the resource with the provider client.
func (r *RackReservationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
	utils.ApplyTagsFromSlugs(ctx, r.client, apiReq, data.Tags, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	utils.ApplyCustomFields(ctx, apiReq, data.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Creating rack reservation", map[string]interface{}{
		"rack":        data.Rack.ValueString(),
		"units":       units,
		"user":        data.User.ValueString(),
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
	if !utils.ValidateStatusCode(&resp.Diagnostics, "create rack reservation", httpResp, http.StatusCreated) {
		return
	}

	// Map response to state
	r.mapToState(ctx, result, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(data.ID.ValueString()), data.CustomFields, &resp.Diagnostics)
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
		if utils.HandleNotFound(httpResp, func() { resp.State.RemoveResource(ctx) }) {
			return
		}
		resp.Diagnostics.AddError(
			"Error reading rack reservation",
			utils.FormatAPIError(fmt.Sprintf("read rack reservation ID %d", id), err, httpResp),
		)
		return
	}
	if !utils.ValidateStatusCode(&resp.Diagnostics, "read rack reservation", httpResp, http.StatusOK) {
		return
	}

	// Map response to state
	r.mapToState(ctx, result, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(data.ID.ValueString()), data.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates the resource.
func (r *RackReservationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, plan RackReservationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Parse ID
	var id int32
	_, err := fmt.Sscanf(plan.ID.ValueString(), "%d", &id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error parsing rack reservation ID",
			fmt.Sprintf("Could not parse ID '%s': %s", plan.ID.ValueString(), err.Error()),
		)
		return
	}

	// Lookup rack
	rack, diags := lookup.LookupRack(ctx, r.client, plan.Rack.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Lookup user
	user, diags := lookup.LookupUser(ctx, r.client, plan.User.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert units to []int32
	var unitsInt64 []int64
	diags = plan.Units.ElementsAs(ctx, &unitsInt64, false)
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
	apiReq := netbox.NewRackReservationRequest(*rack, units, *user, plan.Description.ValueString())

	// Set optional fields
	if !plan.Tenant.IsNull() && !plan.Tenant.IsUnknown() {
		tenant, tenantDiags := lookup.LookupTenant(ctx, r.client, plan.Tenant.ValueString())
		resp.Diagnostics.Append(tenantDiags...)
		if resp.Diagnostics.HasError() {
			return
		}
		apiReq.SetTenant(*tenant)
	} else if !plan.Tenant.IsUnknown() {
		// Explicitly set nil to clear the field
		apiReq.SetTenantNil()
	}

	// Apply optional fields with merge-aware helpers
	utils.ApplyComments(apiReq, plan.Comments)
	utils.ApplyTagsFromSlugs(ctx, r.client, apiReq, plan.Tags, &resp.Diagnostics)
	utils.ApplyCustomFieldsWithMerge(ctx, apiReq, plan.CustomFields, state.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Updating rack reservation", map[string]interface{}{
		"id":          id,
		"rack":        plan.Rack.ValueString(),
		"units":       units,
		"user":        plan.User.ValueString(),
		"description": plan.Description.ValueString(),
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
	if !utils.ValidateStatusCode(&resp.Diagnostics, "update rack reservation", httpResp, http.StatusOK) {
		return
	}

	// Map response to state
	r.mapToState(ctx, result, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(plan.ID.ValueString()), plan.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
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
		if utils.HandleNotFound(httpResp, nil) {
			return
		}
		resp.Diagnostics.AddError(
			"Error deleting rack reservation",
			utils.FormatAPIError(fmt.Sprintf("delete rack reservation ID %d", id), err, httpResp),
		)
		return
	}
	if !utils.ValidateStatusCode(&resp.Diagnostics, "delete rack reservation", httpResp, http.StatusNoContent) {
		return
	}
}

// ImportState imports the resource state.
func (r *RackReservationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if parsed, ok := utils.ParseImportIdentityCustomFields(ctx, req.Identity, &resp.Diagnostics); ok {
		if resp.Diagnostics.HasError() {
			return
		}
		if parsed.ID == "" {
			resp.Diagnostics.AddError("Invalid import identity", "Identity id must be provided")
			return
		}

		id, err := utils.ParseID(parsed.ID)
		if err != nil {
			resp.Diagnostics.AddError("Invalid Import ID", fmt.Sprintf("Rack reservation ID must be a number, got: %s", parsed.ID))
			return
		}

		result, httpResp, err := r.client.DcimAPI.DcimRackReservationsRetrieve(ctx, id).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError("Error importing rack reservation", utils.FormatAPIError("import rack reservation", err, httpResp))
			return
		}
		if !utils.ValidateStatusCode(&resp.Diagnostics, "import rack reservation", httpResp, http.StatusOK) {
			return
		}

		var data RackReservationResourceModel
		data.Tags = utils.PopulateTagsSlugFromAPI(ctx, result.HasTags(), result.GetTags(), data.Tags)
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

		r.mapToState(ctx, result, &data, &resp.Diagnostics)
		if parsed.HasCustomFields {
			data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, data.CustomFields, result.GetCustomFields(), &resp.Diagnostics)
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

// mapToState maps the API response to the Terraform state.
func (r *RackReservationResource) mapToState(ctx context.Context, result *netbox.RackReservation, data *RackReservationResourceModel, diags *diag.Diagnostics) {
	data.ID = types.StringValue(fmt.Sprintf("%d", result.GetId()))

	// Map rack (required field) - preserve user's input format
	rack := result.GetRack()
	data.Rack = utils.UpdateReferenceAttribute(data.Rack, rack.GetName(), "", rack.GetId())

	// Map units
	apiUnits := result.GetUnits()
	units := make([]int64, len(apiUnits))
	for i, u := range apiUnits {
		units[i] = int64(u)
	}
	unitsValue, unitDiags := types.SetValueFrom(ctx, types.Int64Type, units)
	diags.Append(unitDiags...)
	data.Units = unitsValue

	// Map user (required field) - preserve user's input format (ID or username)
	user := result.GetUser()
	data.User = utils.UpdateReferenceAttribute(data.User, user.GetUsername(), "", user.GetId())

	// Map tenant - preserve user's input format
	if result.HasTenant() && result.GetTenant().Id != 0 {
		tenant := result.GetTenant()
		data.Tenant = utils.UpdateReferenceAttribute(data.Tenant, tenant.GetName(), tenant.GetSlug(), tenant.GetId())
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

	// Filter tags to owned (slug list format)
	data.Tags = utils.PopulateTagsSlugFilteredToOwned(ctx, result.HasTags(), result.GetTags(), data.Tags)

	// Map custom fields
	data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, data.CustomFields, result.GetCustomFields(), diags)
}
