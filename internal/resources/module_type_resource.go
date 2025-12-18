// Package resources provides Terraform resource implementations for NetBox objects.

package resources

import (
	"context"
	"fmt"

	"github.com/bab3l/go-netbox"
	lookup "github.com/bab3l/terraform-provider-netbox/internal/netboxlookup"
	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.

var (
	_ resource.Resource = &ModuleTypeResource{}

	_ resource.ResourceWithConfigure = &ModuleTypeResource{}

	_ resource.ResourceWithImportState = &ModuleTypeResource{}
)

// NewModuleTypeResource returns a new resource implementing the module type resource.

func NewModuleTypeResource() resource.Resource {

	return &ModuleTypeResource{}

}

// ModuleTypeResource defines the resource implementation.

type ModuleTypeResource struct {
	client *netbox.APIClient
}

// ModuleTypeResourceModel describes the resource data model.

type ModuleTypeResourceModel struct {
	ID types.String `tfsdk:"id"`

	Manufacturer types.String `tfsdk:"manufacturer"`

	Model types.String `tfsdk:"model"`

	PartNumber types.String `tfsdk:"part_number"`

	Airflow types.String `tfsdk:"airflow"`

	Weight types.Float64 `tfsdk:"weight"`

	WeightUnit types.String `tfsdk:"weight_unit"`

	Description types.String `tfsdk:"description"`

	Comments types.String `tfsdk:"comments"`

	Tags types.Set `tfsdk:"tags"`

	CustomFields types.Set `tfsdk:"custom_fields"`
}

// Metadata returns the resource type name.

func (r *ModuleTypeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {

	resp.TypeName = req.ProviderTypeName + "_module_type"

}

// Schema defines the schema for the resource.

func (r *ModuleTypeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	resp.Schema = schema.Schema{

		MarkdownDescription: "Manages a module type in NetBox. Module types define hardware module specifications (model, manufacturer, etc.) that can be instantiated as modules within devices.",

		Attributes: map[string]schema.Attribute{

			"id": schema.StringAttribute{

				MarkdownDescription: "The unique numeric ID of the module type.",

				Computed: true,

				PlanModifiers: []planmodifier.String{

					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"manufacturer": schema.StringAttribute{

				MarkdownDescription: "The manufacturer of the module type (ID or slug).",

				Required: true,
			},

			"model": schema.StringAttribute{

				MarkdownDescription: "The model name/number of the module type.",

				Required: true,
			},

			"part_number": schema.StringAttribute{

				MarkdownDescription: "Discrete part number (optional).",

				Optional: true,
			},

			"airflow": schema.StringAttribute{

				MarkdownDescription: "Airflow direction. Valid values: `front-to-rear`, `rear-to-front`, `left-to-right`, `right-to-left`, `side-to-rear`, `passive`, `mixed`.",

				Optional: true,
			},

			"weight": schema.Float64Attribute{

				MarkdownDescription: "Weight of the module.",

				Optional: true,
			},

			"weight_unit": schema.StringAttribute{

				MarkdownDescription: "Unit for weight measurement. Valid values: `kg`, `g`, `lb`, `oz`.",

				Optional: true,
			},

			"description": schema.StringAttribute{

				MarkdownDescription: "A description of the module type.",

				Optional: true,
			},

			"comments": schema.StringAttribute{

				MarkdownDescription: "Additional comments or notes.",

				Optional: true,
			},

			"tags": nbschema.TagsAttribute(),

			"custom_fields": nbschema.CustomFieldsAttribute(),
		},
	}

}

// Configure adds the provider configured client to the resource.

func (r *ModuleTypeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {

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

func (r *ModuleTypeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var data ModuleTypeResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	// Lookup manufacturer

	manufacturer, diags := lookup.LookupManufacturer(ctx, r.client, data.Manufacturer.ValueString())

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {

		return

	}

	// Build request

	apiReq := netbox.NewWritableModuleTypeRequest(*manufacturer, data.Model.ValueString())

	// Set optional fields

	if !data.PartNumber.IsNull() && !data.PartNumber.IsUnknown() {

		apiReq.SetPartNumber(data.PartNumber.ValueString())

	}

	if !data.Airflow.IsNull() && !data.Airflow.IsUnknown() {

		airflow := netbox.ModuleTypeAirflowValue(data.Airflow.ValueString())

		apiReq.SetAirflow(airflow)

	}

	if !data.Weight.IsNull() && !data.Weight.IsUnknown() {

		apiReq.SetWeight(data.Weight.ValueFloat64())

	}

	if !data.WeightUnit.IsNull() && !data.WeightUnit.IsUnknown() {

		weightUnit := netbox.DeviceTypeWeightUnitValue(data.WeightUnit.ValueString())

		apiReq.SetWeightUnit(weightUnit)

	}

	if !data.Description.IsNull() && !data.Description.IsUnknown() {

		apiReq.SetDescription(data.Description.ValueString())

	}

	if !data.Comments.IsNull() && !data.Comments.IsUnknown() {

		apiReq.SetComments(data.Comments.ValueString())

	}

	// Handle tags

	if !data.Tags.IsNull() && !data.Tags.IsUnknown() {

		tags, tagDiags := utils.TagModelsToNestedTagRequests(ctx, data.Tags)

		resp.Diagnostics.Append(tagDiags...)

		if resp.Diagnostics.HasError() {

			return

		}

		apiReq.SetTags(tags)

	}

	// Handle custom fields

	if !data.CustomFields.IsNull() && !data.CustomFields.IsUnknown() {

		var cfModels []utils.CustomFieldModel

		resp.Diagnostics.Append(data.CustomFields.ElementsAs(ctx, &cfModels, false)...)

		if resp.Diagnostics.HasError() {

			return

		}

		apiReq.SetCustomFields(utils.CustomFieldModelsToMap(cfModels))

	}

	tflog.Debug(ctx, "Creating module type", map[string]interface{}{

		"manufacturer": data.Manufacturer.ValueString(),

		"model": data.Model.ValueString(),
	})

	response, httpResp, err := r.client.DcimAPI.DcimModuleTypesCreate(ctx).WritableModuleTypeRequest(*apiReq).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		resp.Diagnostics.AddError(

			"Error creating module type",

			utils.FormatAPIError(fmt.Sprintf("create module type %s", data.Model.ValueString()), err, httpResp),
		)

		return

	}

	// Map response to model

	r.mapResponseToModel(ctx, response, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {

		return

	}

	tflog.Trace(ctx, "Created module type", map[string]interface{}{

		"id": data.ID.ValueString(),

		"model": data.Model.ValueString(),
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

// Read refreshes the resource state.

func (r *ModuleTypeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	var data ModuleTypeResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	typeID, err := utils.ParseID(data.ID.ValueString())

	if err != nil {

		resp.Diagnostics.AddError(

			"Invalid Module Type ID",

			fmt.Sprintf("Module Type ID must be a number, got: %s", data.ID.ValueString()),
		)

		return

	}

	tflog.Debug(ctx, "Reading module type", map[string]interface{}{

		"id": typeID,
	})

	response, httpResp, err := r.client.DcimAPI.DcimModuleTypesRetrieve(ctx, typeID).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		if httpResp != nil && httpResp.StatusCode == 404 {

			resp.State.RemoveResource(ctx)

			return

		}

		resp.Diagnostics.AddError(

			"Error reading module type",

			utils.FormatAPIError(fmt.Sprintf("read module type ID %d", typeID), err, httpResp),
		)

		return

	}

	// Map response to model

	r.mapResponseToModel(ctx, response, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {

		return

	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

// Update updates the resource.

func (r *ModuleTypeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var data ModuleTypeResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	typeID, err := utils.ParseID(data.ID.ValueString())

	if err != nil {

		resp.Diagnostics.AddError(

			"Invalid Module Type ID",

			fmt.Sprintf("Module Type ID must be a number, got: %s", data.ID.ValueString()),
		)

		return

	}

	// Lookup manufacturer

	manufacturer, diags := lookup.LookupManufacturer(ctx, r.client, data.Manufacturer.ValueString())

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {

		return

	}

	// Build request

	apiReq := netbox.NewWritableModuleTypeRequest(*manufacturer, data.Model.ValueString())

	// Set optional fields

	if !data.PartNumber.IsNull() && !data.PartNumber.IsUnknown() {

		apiReq.SetPartNumber(data.PartNumber.ValueString())

	}

	if !data.Airflow.IsNull() && !data.Airflow.IsUnknown() {

		airflow := netbox.ModuleTypeAirflowValue(data.Airflow.ValueString())

		apiReq.SetAirflow(airflow)

	}

	if !data.Weight.IsNull() && !data.Weight.IsUnknown() {

		apiReq.SetWeight(data.Weight.ValueFloat64())

	}

	if !data.WeightUnit.IsNull() && !data.WeightUnit.IsUnknown() {

		weightUnit := netbox.DeviceTypeWeightUnitValue(data.WeightUnit.ValueString())

		apiReq.SetWeightUnit(weightUnit)

	}

	if !data.Description.IsNull() && !data.Description.IsUnknown() {

		apiReq.SetDescription(data.Description.ValueString())

	}

	if !data.Comments.IsNull() && !data.Comments.IsUnknown() {

		apiReq.SetComments(data.Comments.ValueString())

	}

	// Handle tags

	if !data.Tags.IsNull() && !data.Tags.IsUnknown() {

		tags, tagDiags := utils.TagModelsToNestedTagRequests(ctx, data.Tags)

		resp.Diagnostics.Append(tagDiags...)

		if resp.Diagnostics.HasError() {

			return

		}

		apiReq.SetTags(tags)

	}

	// Handle custom fields

	if !data.CustomFields.IsNull() && !data.CustomFields.IsUnknown() {

		var cfModels []utils.CustomFieldModel

		resp.Diagnostics.Append(data.CustomFields.ElementsAs(ctx, &cfModels, false)...)

		if resp.Diagnostics.HasError() {

			return

		}

		apiReq.SetCustomFields(utils.CustomFieldModelsToMap(cfModels))

	}

	tflog.Debug(ctx, "Updating module type", map[string]interface{}{

		"id": typeID,
	})

	response, httpResp, err := r.client.DcimAPI.DcimModuleTypesUpdate(ctx, typeID).WritableModuleTypeRequest(*apiReq).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		resp.Diagnostics.AddError(

			"Error updating module type",

			utils.FormatAPIError(fmt.Sprintf("update module type ID %d", typeID), err, httpResp),
		)

		return

	}

	// Map response to model

	r.mapResponseToModel(ctx, response, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {

		return

	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

// Delete deletes the resource.

func (r *ModuleTypeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var data ModuleTypeResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	typeID, err := utils.ParseID(data.ID.ValueString())

	if err != nil {

		resp.Diagnostics.AddError(

			"Invalid Module Type ID",

			fmt.Sprintf("Module Type ID must be a number, got: %s", data.ID.ValueString()),
		)

		return

	}

	tflog.Debug(ctx, "Deleting module type", map[string]interface{}{

		"id": typeID,
	})

	httpResp, err := r.client.DcimAPI.DcimModuleTypesDestroy(ctx, typeID).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		if httpResp != nil && httpResp.StatusCode == 404 {

			return

		}

		resp.Diagnostics.AddError(

			"Error deleting module type",

			utils.FormatAPIError(fmt.Sprintf("delete module type ID %d", typeID), err, httpResp),
		)

		return

	}

}

// ImportState imports an existing resource.

func (r *ModuleTypeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	typeID, err := utils.ParseID(req.ID)

	if err != nil {

		resp.Diagnostics.AddError(

			"Invalid Import ID",

			fmt.Sprintf("Module Type ID must be a number, got: %s", req.ID),
		)

		return

	}

	response, httpResp, err := r.client.DcimAPI.DcimModuleTypesRetrieve(ctx, typeID).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		resp.Diagnostics.AddError(

			"Error importing module type",

			utils.FormatAPIError(fmt.Sprintf("import module type ID %d", typeID), err, httpResp),
		)

		return

	}

	var data ModuleTypeResourceModel

	r.mapResponseToModel(ctx, response, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {

		return

	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

// mapResponseToModel maps the API response to the Terraform model.

func (r *ModuleTypeResource) mapResponseToModel(ctx context.Context, moduleType *netbox.ModuleType, data *ModuleTypeResourceModel, diags *diag.Diagnostics) {

	data.ID = types.StringValue(fmt.Sprintf("%d", moduleType.GetId()))

	data.Model = types.StringValue(moduleType.GetModel())

	// Map manufacturer - preserve user's input format

	if mfr := moduleType.GetManufacturer(); mfr.Id != 0 {

		data.Manufacturer = utils.UpdateReferenceAttribute(data.Manufacturer, mfr.GetName(), mfr.GetSlug(), mfr.GetId())

	}

	// Map part_number

	if partNum, ok := moduleType.GetPartNumberOk(); ok && partNum != nil && *partNum != "" {

		data.PartNumber = types.StringValue(*partNum)

	} else {

		data.PartNumber = types.StringNull()

	}

	// Map airflow

	if moduleType.Airflow.IsSet() && moduleType.Airflow.Get() != nil {

		data.Airflow = types.StringValue(string(moduleType.Airflow.Get().GetValue()))

	} else {

		data.Airflow = types.StringNull()

	}

	// Map weight

	if moduleType.Weight.IsSet() && moduleType.Weight.Get() != nil {

		data.Weight = types.Float64Value(*moduleType.Weight.Get())

	} else {

		data.Weight = types.Float64Null()

	}

	// Map weight_unit

	if moduleType.WeightUnit.IsSet() && moduleType.WeightUnit.Get() != nil {

		data.WeightUnit = types.StringValue(string(moduleType.WeightUnit.Get().GetValue()))

	} else {

		data.WeightUnit = types.StringNull()

	}

	// Map description

	if desc, ok := moduleType.GetDescriptionOk(); ok && desc != nil && *desc != "" {

		data.Description = types.StringValue(*desc)

	} else {

		data.Description = types.StringNull()

	}

	// Map comments

	if comments, ok := moduleType.GetCommentsOk(); ok && comments != nil && *comments != "" {

		data.Comments = types.StringValue(*comments)

	} else {

		data.Comments = types.StringNull()

	}

	// Handle tags

	if moduleType.HasTags() && len(moduleType.GetTags()) > 0 {

		tags := utils.NestedTagsToTagModels(moduleType.GetTags())

		tagsValue, tagDiags := types.SetValueFrom(ctx, utils.GetTagsAttributeType().ElemType, tags)

		diags.Append(tagDiags...)

		if diags.HasError() {

			return

		}

		data.Tags = tagsValue

	} else {

		data.Tags = types.SetNull(utils.GetTagsAttributeType().ElemType)

	}

	// Handle custom fields

	if moduleType.HasCustomFields() {

		apiCustomFields := moduleType.GetCustomFields()

		var stateCustomFieldModels []utils.CustomFieldModel

		if !data.CustomFields.IsNull() && !data.CustomFields.IsUnknown() {

			data.CustomFields.ElementsAs(ctx, &stateCustomFieldModels, false)

		}

		customFields := utils.MapToCustomFieldModels(apiCustomFields, stateCustomFieldModels)

		customFieldsValue, cfDiags := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, customFields)

		diags.Append(cfDiags...)

		if diags.HasError() {

			return

		}

		data.CustomFields = customFieldsValue

	} else {

		data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)

	}

}
