// Package resources contains Terraform resource implementations for NetBox objects.

package resources

import (
	"context"
	"fmt"
	"maps"

	"github.com/bab3l/go-netbox"
	"github.com/bab3l/terraform-provider-netbox/internal/netboxlookup"
	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.

var (
	_ resource.Resource = &RackTypeResource{}

	_ resource.ResourceWithConfigure = &RackTypeResource{}

	_ resource.ResourceWithImportState = &RackTypeResource{}
)

// NewRackTypeResource returns a new resource implementing the RackType resource.

func NewRackTypeResource() resource.Resource {
	return &RackTypeResource{}
}

// RackTypeResource defines the resource implementation.

type RackTypeResource struct {
	client *netbox.APIClient
}

// RackTypeResourceModel describes the resource data model.

type RackTypeResourceModel struct {
	ID types.String `tfsdk:"id"`

	Manufacturer types.String `tfsdk:"manufacturer"`

	Model types.String `tfsdk:"model"`

	Slug types.String `tfsdk:"slug"`

	Description types.String `tfsdk:"description"`

	FormFactor types.String `tfsdk:"form_factor"`

	Width types.Int64 `tfsdk:"width"`

	UHeight types.Int64 `tfsdk:"u_height"`

	StartingUnit types.Int64 `tfsdk:"starting_unit"`

	DescUnits types.Bool `tfsdk:"desc_units"`

	OuterWidth types.Int64 `tfsdk:"outer_width"`

	OuterDepth types.Int64 `tfsdk:"outer_depth"`

	OuterUnit types.String `tfsdk:"outer_unit"`

	Weight types.Float64 `tfsdk:"weight"`

	MaxWeight types.Int64 `tfsdk:"max_weight"`

	WeightUnit types.String `tfsdk:"weight_unit"`

	MountingDepth types.Int64 `tfsdk:"mounting_depth"`

	Comments types.String `tfsdk:"comments"`

	Tags types.Set `tfsdk:"tags"`

	CustomFields types.Set `tfsdk:"custom_fields"`
}

// Metadata returns the resource type name.

func (r *RackTypeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_rack_type"
}

// Schema defines the schema for the resource.

func (r *RackTypeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a rack type in NetBox. Rack types are templates that define the specifications for racks, including dimensions, capacity, and physical characteristics.",

		Attributes: map[string]schema.Attribute{
			"id": nbschema.IDAttribute("rack type"),

			"manufacturer": nbschema.RequiredReferenceAttribute("manufacturer", "The manufacturer of this rack type."),

			"model": nbschema.ModelAttribute("rack type", 100),

			"slug": nbschema.SlugAttribute("rack type"),

			"description": nbschema.DescriptionAttribute("rack type"),

			"form_factor": schema.StringAttribute{
				MarkdownDescription: "Form factor of the rack type. Valid values include: 2-post-frame, 4-post-frame, 4-post-cabinet, wall-frame, wall-frame-vertical, wall-cabinet, wall-cabinet-vertical.",

				Optional: true,
			},

			"width": schema.Int64Attribute{
				MarkdownDescription: "Rail-to-rail width in inches. Common values: 10, 19, 21, 23.",

				Optional: true,

				Computed: true,
			},

			"u_height": schema.Int64Attribute{
				MarkdownDescription: "Height in rack units (U). Default is 42.",

				Optional: true,

				Computed: true,
			},

			"starting_unit": schema.Int64Attribute{
				MarkdownDescription: "Starting unit number for the rack. Default is 1.",

				Optional: true,

				Computed: true,
			},

			"desc_units": schema.BoolAttribute{
				MarkdownDescription: "Whether units are numbered top-to-bottom (descending). Default is false.",

				Optional: true,

				Computed: true,
			},

			"outer_width": schema.Int64Attribute{
				MarkdownDescription: "Outer dimension of rack (width) in millimeters or inches.",

				Optional: true,
			},

			"outer_depth": schema.Int64Attribute{
				MarkdownDescription: "Outer dimension of rack (depth) in millimeters or inches.",

				Optional: true,
			},

			"outer_unit": schema.StringAttribute{
				MarkdownDescription: "Unit for outer dimensions. Valid values: mm (millimeters), in (inches).",

				Optional: true,
			},

			"weight": schema.Float64Attribute{
				MarkdownDescription: "Weight of the rack.",

				Optional: true,
			},

			"max_weight": schema.Int64Attribute{
				MarkdownDescription: "Maximum load capacity for the rack.",

				Optional: true,
			},

			"weight_unit": schema.StringAttribute{
				MarkdownDescription: "Unit for weight. Valid values: kg (kilograms), g (grams), lb (pounds), oz (ounces).",

				Optional: true,
			},

			"mounting_depth": schema.Int64Attribute{
				MarkdownDescription: "Maximum depth of a mounted device, in millimeters. For four-post racks, this is the distance between the front and rear rails.",

				Optional: true,
			},
		},
	}

	// Add common descriptive attributes (description, comments)
	maps.Copy(resp.Schema.Attributes, nbschema.CommonDescriptiveAttributes("rack type"))

	// Add common metadata attributes (tags, custom_fields)
	maps.Copy(resp.Schema.Attributes, nbschema.CommonMetadataAttributes())
}

// Configure adds the provider configured client to the resource.

func (r *RackTypeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create creates a new rack type resource.

func (r *RackTypeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data RackTypeResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Build the request

	// For Create, there is no prior state so pass empty state
	var emptyState RackTypeResourceModel
	rackTypeRequest, diags := r.buildRequest(ctx, &data, &emptyState)

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating rack type", map[string]interface{}{
		"model": data.Model.ValueString(),
	})

	rackType, httpResp, err := r.client.DcimAPI.DcimRackTypesCreate(ctx).WritableRackTypeRequest(*rackTypeRequest).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {
		resp.Diagnostics.AddError(

			"Error creating rack type",

			utils.FormatAPIError("create rack type", err, httpResp),
		)

		return
	}

	// Map response to state

	r.mapResponseToModel(ctx, rackType, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Created rack type", map[string]interface{}{
		"id": rackType.GetId(),

		"model": rackType.GetModel(),
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read reads the rack type resource.

func (r *RackTypeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data RackTypeResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Preserve original custom_fields to detect null/empty cases
	originalCustomFields := data.CustomFields

	rackTypeID, err := utils.ParseID(data.ID.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(

			"Invalid Rack Type ID",

			fmt.Sprintf("Could not parse rack type ID: %s", err),
		)

		return
	}

	tflog.Debug(ctx, "Reading rack type", map[string]interface{}{
		"id": rackTypeID,
	})

	rackType, httpResp, err := r.client.DcimAPI.DcimRackTypesRetrieve(ctx, rackTypeID).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			resp.State.RemoveResource(ctx)

			return
		}

		resp.Diagnostics.AddError(

			"Error reading rack type",

			utils.FormatAPIError(fmt.Sprintf("read rack type ID %d", rackTypeID), err, httpResp),
		)

		return
	}

	// Map response to state

	r.mapResponseToModel(ctx, rackType, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	// If custom_fields was explicitly null/empty in config, preserve that
	if originalCustomFields.IsNull() || len(originalCustomFields.Elements()) == 0 {
		data.CustomFields = originalCustomFields
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates the rack type resource.

func (r *RackTypeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, plan RackTypeResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	rackTypeID, err := utils.ParseID(plan.ID.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(

			"Invalid Rack Type ID",

			fmt.Sprintf("Could not parse rack type ID: %s", err),
		)

		return
	}

	// Build the request

	rackTypeRequest, diags := r.buildRequest(ctx, &plan, &state)

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Updating rack type", map[string]interface{}{
		"id": rackTypeID,
	})

	rackType, httpResp, err := r.client.DcimAPI.DcimRackTypesUpdate(ctx, rackTypeID).WritableRackTypeRequest(*rackTypeRequest).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {
		resp.Diagnostics.AddError(

			"Error updating rack type",

			utils.FormatAPIError(fmt.Sprintf("update rack type ID %d", rackTypeID), err, httpResp),
		)

		return
	}

	// Store plan's custom_fields to filter the response
	planCustomFields := plan.CustomFields

	// Map response to state

	r.mapResponseToModel(ctx, rackType, &plan, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	// Filter custom_fields to only those owned by this resource
	plan.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, planCustomFields, rackType.GetCustomFields(), &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete deletes the rack type resource.

func (r *RackTypeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data RackTypeResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	rackTypeID, err := utils.ParseID(data.ID.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(

			"Invalid Rack Type ID",

			fmt.Sprintf("Could not parse rack type ID: %s", err),
		)

		return
	}

	tflog.Debug(ctx, "Deleting rack type", map[string]interface{}{
		"id": rackTypeID,
	})

	httpResp, err := r.client.DcimAPI.DcimRackTypesDestroy(ctx, rackTypeID).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			return
		}

		resp.Diagnostics.AddError(

			"Error deleting rack type",

			utils.FormatAPIError(fmt.Sprintf("delete rack type ID %d", rackTypeID), err, httpResp),
		)

		return
	}
}

// ImportState imports an existing rack type resource.

func (r *RackTypeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// buildRequest builds the API request from the Terraform model.

func (r *RackTypeResource) buildRequest(ctx context.Context, plan, state *RackTypeResourceModel) (*netbox.WritableRackTypeRequest, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Look up manufacturer

	manufacturerRequest, lookupDiags := netboxlookup.LookupManufacturer(ctx, r.client, plan.Manufacturer.ValueString())

	diags.Append(lookupDiags...)

	if diags.HasError() {
		return nil, diags
	}

	// Form factor is required for WritableRackTypeRequest

	formFactor := netbox.PatchedWritableRackTypeRequestFormFactor("")

	if !plan.FormFactor.IsNull() && !plan.FormFactor.IsUnknown() {
		formFactor = netbox.PatchedWritableRackTypeRequestFormFactor(plan.FormFactor.ValueString())
	}

	rackTypeRequest := netbox.NewWritableRackTypeRequest(

		*manufacturerRequest,

		plan.Model.ValueString(),

		plan.Slug.ValueString(),

		formFactor,
	)

	utils.ApplyDescription(rackTypeRequest, plan.Description)

	// FormFactor is already set in the constructor, but update if explicitly provided
	// (no need to set again since we pass it in the constructor)

	if !plan.Width.IsNull() && !plan.Width.IsUnknown() {
		widthVal, err := utils.SafeInt32FromValue(plan.Width)

		if err != nil {
			diags.AddError("Invalid value", fmt.Sprintf("Width value overflow: %s", err))

			return nil, diags
		}

		width := netbox.PatchedWritableRackRequestWidth(widthVal)

		rackTypeRequest.SetWidth(width)
	}

	if !plan.UHeight.IsNull() && !plan.UHeight.IsUnknown() {
		uHeight, err := utils.SafeInt32FromValue(plan.UHeight)

		if err != nil {
			diags.AddError("Invalid value", fmt.Sprintf("UHeight value overflow: %s", err))

			return nil, diags
		}

		rackTypeRequest.SetUHeight(uHeight)
	}

	if !plan.StartingUnit.IsNull() && !plan.StartingUnit.IsUnknown() {
		startingUnit, err := utils.SafeInt32FromValue(plan.StartingUnit)

		if err != nil {
			diags.AddError("Invalid value", fmt.Sprintf("StartingUnit value overflow: %s", err))

			return nil, diags
		}

		rackTypeRequest.SetStartingUnit(startingUnit)
	}

	if !plan.DescUnits.IsNull() && !plan.DescUnits.IsUnknown() {
		rackTypeRequest.SetDescUnits(plan.DescUnits.ValueBool())
	}

	if !plan.OuterWidth.IsNull() && !plan.OuterWidth.IsUnknown() {
		outerWidth, err := utils.SafeInt32FromValue(plan.OuterWidth)

		if err != nil {
			diags.AddError("Invalid value", fmt.Sprintf("OuterWidth value overflow: %s", err))

			return nil, diags
		}

		rackTypeRequest.SetOuterWidth(outerWidth)
	} else if plan.OuterWidth.IsNull() {
		// Use AdditionalProperties to send null because of omitempty in the generated client
		if rackTypeRequest.AdditionalProperties == nil {
			rackTypeRequest.AdditionalProperties = make(map[string]interface{})
		}
		rackTypeRequest.AdditionalProperties["outer_width"] = nil
	}

	if !plan.OuterDepth.IsNull() && !plan.OuterDepth.IsUnknown() {
		outerDepth, err := utils.SafeInt32FromValue(plan.OuterDepth)

		if err != nil {
			diags.AddError("Invalid value", fmt.Sprintf("OuterDepth value overflow: %s", err))

			return nil, diags
		}

		rackTypeRequest.SetOuterDepth(outerDepth)
	} else if plan.OuterDepth.IsNull() {
		// Use AdditionalProperties to send null because of omitempty in the generated client
		if rackTypeRequest.AdditionalProperties == nil {
			rackTypeRequest.AdditionalProperties = make(map[string]interface{})
		}
		rackTypeRequest.AdditionalProperties["outer_depth"] = nil
	}

	if !plan.OuterUnit.IsNull() && !plan.OuterUnit.IsUnknown() {
		outerUnit := netbox.PatchedWritableRackRequestOuterUnit(plan.OuterUnit.ValueString())

		rackTypeRequest.SetOuterUnit(outerUnit)
	} else if plan.OuterUnit.IsNull() {
		// Use AdditionalProperties to send null because of omitempty in the generated client
		if rackTypeRequest.AdditionalProperties == nil {
			rackTypeRequest.AdditionalProperties = make(map[string]interface{})
		}
		rackTypeRequest.AdditionalProperties["outer_unit"] = nil
	}

	if !plan.Weight.IsNull() && !plan.Weight.IsUnknown() {
		rackTypeRequest.SetWeight(plan.Weight.ValueFloat64())
	} else if plan.Weight.IsNull() {
		// Use AdditionalProperties to send null because of omitempty in the generated client
		if rackTypeRequest.AdditionalProperties == nil {
			rackTypeRequest.AdditionalProperties = make(map[string]interface{})
		}
		rackTypeRequest.AdditionalProperties["weight"] = nil
	}

	if !plan.MaxWeight.IsNull() && !plan.MaxWeight.IsUnknown() {
		maxWeight, err := utils.SafeInt32FromValue(plan.MaxWeight)

		if err != nil {
			diags.AddError("Invalid value", fmt.Sprintf("MaxWeight value overflow: %s", err))

			return nil, diags
		}

		rackTypeRequest.SetMaxWeight(maxWeight)
	} else if plan.MaxWeight.IsNull() {
		// Use AdditionalProperties to send null because of omitempty in the generated client
		if rackTypeRequest.AdditionalProperties == nil {
			rackTypeRequest.AdditionalProperties = make(map[string]interface{})
		}
		rackTypeRequest.AdditionalProperties["max_weight"] = nil
	}

	if !plan.WeightUnit.IsNull() && !plan.WeightUnit.IsUnknown() {
		weightUnit := netbox.DeviceTypeWeightUnitValue(plan.WeightUnit.ValueString())

		rackTypeRequest.SetWeightUnit(weightUnit)
	}
	// Note: Don't send explicit null for weight_unit - NetBox has a default value (kg)
	// and sending null violates the NOT NULL constraint in NetBox 4.1.11+

	if !plan.MountingDepth.IsNull() && !plan.MountingDepth.IsUnknown() {
		mountingDepth, err := utils.SafeInt32FromValue(plan.MountingDepth)

		if err != nil {
			diags.AddError("Invalid value", fmt.Sprintf("MountingDepth value overflow: %s", err))

			return nil, diags
		}

		rackTypeRequest.SetMountingDepth(mountingDepth)
	} else if plan.MountingDepth.IsNull() {
		// Use AdditionalProperties to send null because of omitempty in the generated client
		if rackTypeRequest.AdditionalProperties == nil {
			rackTypeRequest.AdditionalProperties = make(map[string]interface{})
		}
		rackTypeRequest.AdditionalProperties["mounting_depth"] = nil
	}

	utils.ApplyComments(rackTypeRequest, plan.Comments)

	utils.ApplyTags(ctx, rackTypeRequest, plan.Tags, &diags)

	utils.ApplyCustomFieldsWithMerge(ctx, rackTypeRequest, plan.CustomFields, state.CustomFields, &diags)

	return rackTypeRequest, diags
}

// mapResponseToModel maps the API response to the Terraform model.

func (r *RackTypeResource) mapResponseToModel(ctx context.Context, rackType *netbox.RackType, data *RackTypeResourceModel, diags *diag.Diagnostics) {
	data.ID = types.StringValue(fmt.Sprintf("%d", rackType.GetId()))

	data.Model = types.StringValue(rackType.GetModel())

	data.Slug = types.StringValue(rackType.GetSlug())

	// Map manufacturer - return ID

	data.Manufacturer = utils.UpdateReferenceAttribute(data.Manufacturer, rackType.Manufacturer.Name, rackType.Manufacturer.Slug, rackType.Manufacturer.Id)

	// Map description

	if desc, ok := rackType.GetDescriptionOk(); ok && desc != nil && *desc != "" {
		data.Description = types.StringValue(*desc)
	} else {
		data.Description = types.StringNull()
	}

	// Map form_factor

	if ff, ok := rackType.GetFormFactorOk(); ok && ff != nil {
		data.FormFactor = types.StringValue(string(ff.GetValue()))
	} else {
		data.FormFactor = types.StringNull()
	}

	// Map width

	if width, ok := rackType.GetWidthOk(); ok && width != nil {
		data.Width = types.Int64Value(int64(width.GetValue()))
	} else {
		data.Width = types.Int64Null()
	}

	// Map u_height

	if uHeight, ok := rackType.GetUHeightOk(); ok && uHeight != nil {
		data.UHeight = types.Int64Value(int64(*uHeight))
	} else {
		data.UHeight = types.Int64Null()
	}

	// Map starting_unit

	if startingUnit, ok := rackType.GetStartingUnitOk(); ok && startingUnit != nil {
		data.StartingUnit = types.Int64Value(int64(*startingUnit))
	} else {
		data.StartingUnit = types.Int64Null()
	}

	// Map desc_units - API returns false by default, so we should preserve that

	if descUnits, ok := rackType.GetDescUnitsOk(); ok {
		data.DescUnits = types.BoolValue(*descUnits)
	} else {
		data.DescUnits = types.BoolNull()
	}

	// Map outer_width

	if outerWidth, ok := rackType.GetOuterWidthOk(); ok && outerWidth != nil {
		data.OuterWidth = types.Int64Value(int64(*outerWidth))
	} else {
		data.OuterWidth = types.Int64Null()
	}

	// Map outer_depth

	if outerDepth, ok := rackType.GetOuterDepthOk(); ok && outerDepth != nil {
		data.OuterDepth = types.Int64Value(int64(*outerDepth))
	} else {
		data.OuterDepth = types.Int64Null()
	}

	// Map outer_unit

	if outerUnit, ok := rackType.GetOuterUnitOk(); ok && outerUnit != nil {
		data.OuterUnit = types.StringValue(string(outerUnit.GetValue()))
	} else {
		data.OuterUnit = types.StringNull()
	}

	// Map weight

	if weight, ok := rackType.GetWeightOk(); ok && weight != nil {
		data.Weight = types.Float64Value(*weight)
	} else {
		data.Weight = types.Float64Null()
	}

	// Map max_weight

	if maxWeight, ok := rackType.GetMaxWeightOk(); ok && maxWeight != nil {
		data.MaxWeight = types.Int64Value(int64(*maxWeight))
	} else {
		data.MaxWeight = types.Int64Null()
	}

	// Map weight_unit

	if weightUnit, ok := rackType.GetWeightUnitOk(); ok && weightUnit != nil {
		data.WeightUnit = types.StringValue(string(weightUnit.GetValue()))
	} else {
		data.WeightUnit = types.StringNull()
	}

	// Map mounting_depth

	if mountingDepth, ok := rackType.GetMountingDepthOk(); ok && mountingDepth != nil {
		data.MountingDepth = types.Int64Value(int64(*mountingDepth))
	} else {
		data.MountingDepth = types.Int64Null()
	}

	// Map comments

	if comments, ok := rackType.GetCommentsOk(); ok && comments != nil && *comments != "" {
		data.Comments = types.StringValue(*comments)
	} else {
		data.Comments = types.StringNull()
	}

	// Handle tags
	data.Tags = utils.PopulateTagsFromAPI(ctx, rackType.HasTags(), rackType.GetTags(), data.Tags, diags)
	if diags.HasError() {
		return
	}

	// Handle custom fields
	data.CustomFields = utils.PopulateCustomFieldsFromAPI(ctx, rackType.HasCustomFields(), rackType.GetCustomFields(), data.CustomFields, diags)
	if diags.HasError() {
		return
	}
}
