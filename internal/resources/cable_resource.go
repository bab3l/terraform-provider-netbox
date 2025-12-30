// Package resources contains Terraform resource implementations for the Netbox provider.

package resources

import (
	"context"
	"fmt"
	"maps"

	"github.com/bab3l/go-netbox"
	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/float64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &CableResource{}
var _ resource.ResourceWithImportState = &CableResource{}

func NewCableResource() resource.Resource {
	return &CableResource{}
}

// CableResource defines the resource implementation.
type CableResource struct {
	client *netbox.APIClient
}

// CableResourceModel describes the resource data model.
type CableResourceModel struct {
	ID            types.String  `tfsdk:"id"`
	ATerminations types.List    `tfsdk:"a_terminations"`
	BTerminations types.List    `tfsdk:"b_terminations"`
	Type          types.String  `tfsdk:"type"`
	Status        types.String  `tfsdk:"status"`
	Tenant        types.String  `tfsdk:"tenant"`
	Label         types.String  `tfsdk:"label"`
	Color         types.String  `tfsdk:"color"`
	Length        types.Float64 `tfsdk:"length"`
	LengthUnit    types.String  `tfsdk:"length_unit"`
	Description   types.String  `tfsdk:"description"`
	Comments      types.String  `tfsdk:"comments"`
	Tags          types.Set     `tfsdk:"tags"`
	CustomFields  types.Set     `tfsdk:"custom_fields"`
}

// TerminationModel represents a cable termination point.
type TerminationModel struct {
	ObjectType types.String `tfsdk:"object_type"`
	ObjectID   types.Int64  `tfsdk:"object_id"`
}

func (r *CableResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cable"
}

func (r *CableResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	terminationNestedObject := schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"object_type": schema.StringAttribute{
				MarkdownDescription: "Content type of the termination object. Common values: `dcim.interface`, `dcim.frontport`, `dcim.rearport`, `dcim.powerport`, `dcim.poweroutlet`, `dcim.consoleport`, `dcim.consoleserverport`, `circuits.circuittermination`.",
				Required:            true,
			},
			"object_id": schema.Int64Attribute{
				MarkdownDescription: "ID of the termination object.",
				Required:            true,
			},
		},
	}
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a cable connection between two endpoints in Netbox. Cables represent physical connections between interfaces, ports, or circuit terminations.",
		Attributes: map[string]schema.Attribute{
			"id": nbschema.IDAttribute("cable"),
			"a_terminations": schema.ListNestedAttribute{
				MarkdownDescription: "A-side termination points for this cable. Each termination specifies an object type and ID.",
				Required:            true,
				NestedObject:        terminationNestedObject,
			},
			"b_terminations": schema.ListNestedAttribute{
				MarkdownDescription: "B-side termination points for this cable. Each termination specifies an object type and ID.",
				Required:            true,
				NestedObject:        terminationNestedObject,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Type of cable. Valid values: `cat3`, `cat5`, `cat5e`, `cat6`, `cat6a`, `cat7`, `cat7a`, `cat8`, `dac-active`, `dac-passive`, `mrj21-trunk`, `coaxial`, `mmf`, `mmf-om1`, `mmf-om2`, `mmf-om3`, `mmf-om4`, `mmf-om5`, `smf`, `smf-os1`, `smf-os2`, `aoc`, `usb`, `power`.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"cat3", "cat5", "cat5e", "cat6", "cat6a", "cat7", "cat7a", "cat8",
						"dac-active", "dac-passive", "mrj21-trunk", "coaxial",
						"mmf", "mmf-om1", "mmf-om2", "mmf-om3", "mmf-om4", "mmf-om5",
						"smf", "smf-os1", "smf-os2", "aoc", "usb", "power", "",
					),
				},
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "Connection status. Valid values: `connected`, `planned`, `decommissioning`. Defaults to `connected`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("connected"),
				Validators: []validator.String{
					stringvalidator.OneOf("connected", "planned", "decommissioning"),
				},
			},
			"tenant": schema.StringAttribute{
				MarkdownDescription: "ID of the tenant that owns this cable.",
				Optional:            true,
			},
			"label": schema.StringAttribute{
				MarkdownDescription: "Physical label attached to the cable.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(100),
				},
			},
			"color": nbschema.ColorAttribute("cable"),
			"length": schema.Float64Attribute{
				MarkdownDescription: "Length of the cable.",
				Optional:            true,
				Validators: []validator.Float64{
					float64validator.AtLeast(0),
				},
			},
			"length_unit": schema.StringAttribute{
				MarkdownDescription: "Unit for cable length. Valid values: `km`, `m`, `cm`, `mi`, `ft`, `in`.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("km", "m", "cm", "mi", "ft", "in", ""),
				},
			},
		},
	}

	// Add description and comments attributes
	maps.Copy(resp.Schema.Attributes, nbschema.CommonDescriptiveAttributes("cable"))

	// Add common metadata attributes (tags, custom_fields)
	maps.Copy(resp.Schema.Attributes, nbschema.CommonMetadataAttributes())
}

func (r *CableResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *CableResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data CableResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Creating cable resource")

	// Build the request
	cableRequest := netbox.NewWritableCableRequest()

	// Set A terminations
	aTerminations, diags := r.parseTerminations(ctx, data.ATerminations)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	cableRequest.ATerminations = aTerminations

	// Set B terminations
	bTerminations, diags := r.parseTerminations(ctx, data.BTerminations)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	cableRequest.BTerminations = bTerminations

	// Set optional fields
	if !data.Type.IsNull() && !data.Type.IsUnknown() && data.Type.ValueString() != "" {
		cableType := netbox.CableType(data.Type.ValueString())
		cableRequest.Type = &cableType
	}

	if !data.Status.IsNull() && !data.Status.IsUnknown() {
		status := netbox.CableStatusValue(data.Status.ValueString())
		cableRequest.Status = &status
	}

	if !data.Tenant.IsNull() && !data.Tenant.IsUnknown() {
		tenantID, err := utils.ParseID(data.Tenant.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Invalid Tenant ID", fmt.Sprintf("Could not parse tenant ID: %s", err))
			return
		}
		tenantRequest := netbox.BriefTenantRequest{Name: fmt.Sprintf("tenant-%d", tenantID)}
		cableRequest.Tenant = *netbox.NewNullableBriefTenantRequest(&tenantRequest)
	}

	if !data.Label.IsNull() && !data.Label.IsUnknown() {
		label := data.Label.ValueString()
		cableRequest.Label = &label
	}

	if !data.Color.IsNull() && !data.Color.IsUnknown() {
		color := data.Color.ValueString()
		cableRequest.Color = &color
	}

	if !data.Length.IsNull() && !data.Length.IsUnknown() {
		length := data.Length.ValueFloat64()
		cableRequest.Length = *netbox.NewNullableFloat64(&length)
	}

	if !data.LengthUnit.IsNull() && !data.LengthUnit.IsUnknown() && data.LengthUnit.ValueString() != "" {
		lengthUnit := netbox.CableLengthUnitValue(data.LengthUnit.ValueString())
		cableRequest.LengthUnit = &lengthUnit
	}

	// Set common fields (description, comments, tags, custom_fields)
	utils.ApplyCommonFields(ctx, cableRequest, data.Description, data.Comments, data.Tags, data.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the cable
	result, httpResp, err := r.client.DcimAPI.DcimCablesCreate(ctx).WritableCableRequest(*cableRequest).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating cable",
			utils.FormatAPIError("create cable", err, httpResp),
		)
		return
	}

	// Map response to state
	resp.Diagnostics.Append(r.mapResponseToState(ctx, result, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Trace(ctx, "Created cable resource", map[string]interface{}{"id": data.ID.ValueString()})
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CableResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data CableResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Could not parse cable ID: %s", err))
		return
	}

	tflog.Debug(ctx, "Reading cable resource", map[string]interface{}{"id": id})
	result, httpResp, err := r.client.DcimAPI.DcimCablesRetrieve(ctx, id).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			tflog.Debug(ctx, "Cable not found, removing from state", map[string]interface{}{"id": id})
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading cable",
			utils.FormatAPIError("read cable", err, httpResp),
		)
		return
	}

	resp.Diagnostics.Append(r.mapResponseToState(ctx, result, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CableResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data CableResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Could not parse cable ID: %s", err))
		return
	}
	tflog.Debug(ctx, "Updating cable resource", map[string]interface{}{"id": id})

	// Build the request
	cableRequest := netbox.NewWritableCableRequest()

	// Set A terminations
	aTerminations, diags := r.parseTerminations(ctx, data.ATerminations)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	cableRequest.ATerminations = aTerminations

	// Set B terminations
	bTerminations, diags := r.parseTerminations(ctx, data.BTerminations)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	cableRequest.BTerminations = bTerminations

	// Set optional fields
	if !data.Type.IsNull() && !data.Type.IsUnknown() && data.Type.ValueString() != "" {
		cableType := netbox.CableType(data.Type.ValueString())
		cableRequest.Type = &cableType
	}

	if !data.Status.IsNull() && !data.Status.IsUnknown() {
		status := netbox.CableStatusValue(data.Status.ValueString())
		cableRequest.Status = &status
	}

	if !data.Tenant.IsNull() && !data.Tenant.IsUnknown() {
		tenantID, err := utils.ParseID(data.Tenant.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Invalid Tenant ID", fmt.Sprintf("Could not parse tenant ID: %s", err))
			return
		}
		tenantRequest := netbox.BriefTenantRequest{Name: fmt.Sprintf("tenant-%d", tenantID)}
		cableRequest.Tenant = *netbox.NewNullableBriefTenantRequest(&tenantRequest)
	} else {
		cableRequest.Tenant = *netbox.NewNullableBriefTenantRequest(nil)
	}

	if !data.Label.IsNull() && !data.Label.IsUnknown() {
		label := data.Label.ValueString()
		cableRequest.Label = &label
	}

	if !data.Color.IsNull() && !data.Color.IsUnknown() {
		color := data.Color.ValueString()
		cableRequest.Color = &color
	}

	if !data.Length.IsNull() && !data.Length.IsUnknown() {
		length := data.Length.ValueFloat64()
		cableRequest.Length = *netbox.NewNullableFloat64(&length)
	} else {
		cableRequest.Length = *netbox.NewNullableFloat64(nil)
	}

	if !data.LengthUnit.IsNull() && !data.LengthUnit.IsUnknown() && data.LengthUnit.ValueString() != "" {
		lengthUnit := netbox.CableLengthUnitValue(data.LengthUnit.ValueString())
		cableRequest.LengthUnit = &lengthUnit
	} else {
		cableRequest.LengthUnit = nil
	}

	// Set common fields (description, comments, tags, custom_fields)
	utils.ApplyCommonFields(ctx, cableRequest, data.Description, data.Comments, data.Tags, data.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update the cable
	result, httpResp, err := r.client.DcimAPI.DcimCablesUpdate(ctx, id).WritableCableRequest(*cableRequest).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating cable",
			utils.FormatAPIError("update cable", err, httpResp),
		)
		return
	}

	// Map response to state
	resp.Diagnostics.Append(r.mapResponseToState(ctx, result, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Trace(ctx, "Updated cable resource", map[string]interface{}{"id": data.ID.ValueString()})
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CableResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CableResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Could not parse cable ID: %s", err))
		return
	}
	tflog.Debug(ctx, "Deleting cable resource", map[string]interface{}{"id": id})
	httpResp, err := r.client.DcimAPI.DcimCablesDestroy(ctx, id).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			tflog.Debug(ctx, "Cable already deleted", map[string]interface{}{"id": id})
			return
		}
		resp.Diagnostics.AddError(
			"Error deleting cable",
			utils.FormatAPIError("delete cable", err, httpResp),
		)
		return
	}
	tflog.Trace(ctx, "Deleted cable resource", map[string]interface{}{"id": id})
}

func (r *CableResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// parseTerminations converts Terraform termination list to API format.
func (r *CableResource) parseTerminations(ctx context.Context, terminations types.List) ([]netbox.GenericObjectRequest, diag.Diagnostics) {
	var diags diag.Diagnostics
	if terminations.IsNull() || terminations.IsUnknown() {
		return nil, diags
	}
	var models []TerminationModel
	diags.Append(terminations.ElementsAs(ctx, &models, false)...)
	if diags.HasError() {
		return nil, diags
	}
	result := make([]netbox.GenericObjectRequest, len(models))
	for i, m := range models {
		objectID, err := utils.SafeInt32FromValue(m.ObjectID)
		if err != nil {
			diags.AddError("Invalid value", fmt.Sprintf("ObjectID value overflow: %s", err))
			return nil, diags
		}
		result[i] = *netbox.NewGenericObjectRequest(
			m.ObjectType.ValueString(),
			objectID,
		)
	}
	return result, diags
}

// mapResponseToState maps API response to Terraform state.
func (r *CableResource) mapResponseToState(ctx context.Context, result *netbox.Cable, data *CableResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics
	data.ID = types.StringValue(fmt.Sprintf("%d", result.GetId()))

	// Map A terminations
	if result.HasATerminations() {
		aTerms, d := r.mapTerminationsToState(ctx, result.GetATerminations())
		diags.Append(d...)
		data.ATerminations = aTerms
	}

	// Map B terminations
	if result.HasBTerminations() {
		bTerms, d := r.mapTerminationsToState(ctx, result.GetBTerminations())
		diags.Append(d...)
		data.BTerminations = bTerms
	}

	// Type
	if result.HasType() && result.GetType() != "" {
		data.Type = types.StringValue(string(result.GetType()))
	} else {
		data.Type = types.StringNull()
	}

	// Status
	if result.HasStatus() {
		status := result.GetStatus()
		data.Status = types.StringValue(string(status.GetValue()))
	} else {
		data.Status = types.StringNull()
	}

	// Tenant - preserve user's input format
	if result.HasTenant() && result.GetTenant().Id != 0 {
		tenant := result.GetTenant()
		data.Tenant = utils.UpdateReferenceAttribute(data.Tenant, tenant.GetName(), tenant.GetSlug(), tenant.GetId())
	} else {
		data.Tenant = types.StringNull()
	}

	// Label
	if result.HasLabel() && result.GetLabel() != "" {
		data.Label = types.StringValue(result.GetLabel())
	} else {
		data.Label = types.StringNull()
	}

	// Color
	if result.HasColor() && result.GetColor() != "" {
		data.Color = types.StringValue(result.GetColor())
	} else {
		data.Color = types.StringNull()
	}

	// Length
	if result.HasLength() && result.GetLength() != 0 {
		data.Length = types.Float64Value(result.GetLength())
	} else {
		data.Length = types.Float64Null()
	}

	// Length unit
	if result.HasLengthUnit() {
		lengthUnit := result.GetLengthUnit()
		if lengthUnit.GetValue() != "" {
			data.LengthUnit = types.StringValue(string(lengthUnit.GetValue()))
		} else {
			data.LengthUnit = types.StringNull()
		}
	} else {
		data.LengthUnit = types.StringNull()
	}

	// Description
	if result.HasDescription() && result.GetDescription() != "" {
		data.Description = types.StringValue(result.GetDescription())
	} else {
		data.Description = types.StringNull()
	}

	// Comments
	if result.HasComments() && result.GetComments() != "" {
		data.Comments = types.StringValue(result.GetComments())
	} else {
		data.Comments = types.StringNull()
	}

	// Tags
	if result.HasTags() && len(result.GetTags()) > 0 {
		tags := utils.NestedTagsToTagModels(result.GetTags())
		tagsValue, d := types.SetValueFrom(ctx, utils.GetTagsAttributeType().ElemType, tags)
		diags.Append(d...)
		data.Tags = tagsValue
	} else {
		data.Tags = types.SetNull(utils.GetTagsAttributeType().ElemType)
	}

	// Custom fields
	if result.HasCustomFields() && len(result.GetCustomFields()) > 0 {
		var existingModels []utils.CustomFieldModel
		if !data.CustomFields.IsNull() {
			diags.Append(data.CustomFields.ElementsAs(ctx, &existingModels, false)...)
		}
		customFields := utils.MapToCustomFieldModels(result.GetCustomFields(), existingModels)
		customFieldsValue, d := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, customFields)
		diags.Append(d...)
		data.CustomFields = customFieldsValue
	} else {
		data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)
	}
	return diags
}

// mapTerminationsToState converts API terminations to Terraform state.
func (r *CableResource) mapTerminationsToState(ctx context.Context, terminations []netbox.GenericObject) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics
	if len(terminations) == 0 {
		return types.ListNull(getTerminationObjectType()), diags
	}
	models := make([]TerminationModel, len(terminations))
	for i, t := range terminations {
		models[i] = TerminationModel{
			ObjectType: types.StringValue(t.GetObjectType()),
			ObjectID:   types.Int64Value(int64(t.GetObjectId())),
		}
	}
	result, d := types.ListValueFrom(ctx, getTerminationObjectType(), models)
	diags.Append(d...)
	return result, diags
}

// getTerminationObjectType returns the Terraform object type for terminations.
func getTerminationObjectType() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"object_type": types.StringType,
			"object_id":   types.Int64Type,
		},
	}
}
