// Package resources contains Terraform resource implementations for NetBox objects.
package resources

import (
	"context"
	"fmt"
	"maps"
	"net/http"

	"github.com/bab3l/go-netbox"
	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &CustomFieldResource{}
	_ resource.ResourceWithImportState = &CustomFieldResource{}
	_ resource.ResourceWithConfigure   = &CustomFieldResource{}
)

// NewCustomFieldResource returns a new resource implementing the CustomField resource.
func NewCustomFieldResource() resource.Resource {
	return &CustomFieldResource{}
}

// CustomFieldResource defines the resource implementation.
type CustomFieldResource struct {
	client *netbox.APIClient
}

// CustomFieldResourceModel describes the resource data model.
type CustomFieldResourceModel struct {
	ID                types.String `tfsdk:"id"`
	ObjectTypes       types.Set    `tfsdk:"object_types"`
	Type              types.String `tfsdk:"type"`
	RelatedObjectType types.String `tfsdk:"related_object_type"`
	Name              types.String `tfsdk:"name"`
	Label             types.String `tfsdk:"label"`
	GroupName         types.String `tfsdk:"group_name"`
	Description       types.String `tfsdk:"description"`
	Required          types.Bool   `tfsdk:"required"`
	SearchWeight      types.Int64  `tfsdk:"search_weight"`
	FilterLogic       types.String `tfsdk:"filter_logic"`
	UIVisible         types.String `tfsdk:"ui_visible"`
	UIEditable        types.String `tfsdk:"ui_editable"`
	IsCloneable       types.Bool   `tfsdk:"is_cloneable"`
	Default           types.String `tfsdk:"default"`
	Weight            types.Int64  `tfsdk:"weight"`
	ValidationMinimum types.Int64  `tfsdk:"validation_minimum"`
	ValidationMaximum types.Int64  `tfsdk:"validation_maximum"`
	ValidationRegex   types.String `tfsdk:"validation_regex"`
	ChoiceSet         types.String `tfsdk:"choice_set"`
	Comments          types.String `tfsdk:"comments"`
}

// Metadata returns the resource type name.
func (r *CustomFieldResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_custom_field"
}

// Schema defines the schema for the resource.
func (r *CustomFieldResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a custom field in NetBox. Custom fields allow extending NetBox objects with additional user-defined attributes.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique numeric ID of the custom field.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"object_types": schema.SetAttribute{
				MarkdownDescription: "The object types this custom field applies to (e.g., 'dcim.device', 'virtualization.virtualmachine').",
				Required:            true,
				ElementType:         types.StringType,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "The type of custom field. Valid values: `text`, `longtext`, `integer`, `decimal`, `boolean`, `date`, `datetime`, `url`, `json`, `select`, `multiselect`, `object`, `multiobject`.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("text", "longtext", "integer", "decimal", "boolean", "date", "datetime", "url", "json", "select", "multiselect", "object", "multiobject"),
				},
			},
			"related_object_type": schema.StringAttribute{
				MarkdownDescription: "The related object type for object and multiobject custom fields (e.g., 'dcim.device').",
				Optional:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Internal field name. Must be lowercase and contain only letters, numbers, and underscores.",
				Required:            true,
			},
			"label": schema.StringAttribute{
				MarkdownDescription: "Name of the field as displayed to users. If not provided, the field's name will be used.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"group_name": schema.StringAttribute{
				MarkdownDescription: "Custom fields within the same group will be displayed together.",
				Optional:            true,
			},
			"required": schema.BoolAttribute{
				MarkdownDescription: "If true, this field is required when creating new objects or editing an existing object.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"search_weight": schema.Int64Attribute{
				MarkdownDescription: "Weighting for search. Lower values are considered more important. Fields with a search weight of zero will be ignored.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(1000),
			},
			"filter_logic": schema.StringAttribute{
				MarkdownDescription: "Filter logic for the custom field. Valid values: `disabled`, `loose`, `exact`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("loose"),
				Validators: []validator.String{
					stringvalidator.OneOf("disabled", "loose", "exact"),
				},
			},
			"ui_visible": schema.StringAttribute{
				MarkdownDescription: "UI visibility setting. Valid values: `always`, `if-set`, `hidden`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("always"),
				Validators: []validator.String{
					stringvalidator.OneOf("always", "if-set", "hidden"),
				},
			},
			"ui_editable": schema.StringAttribute{
				MarkdownDescription: "UI editability setting. Valid values: `yes`, `no`, `hidden`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("yes"),
				Validators: []validator.String{
					stringvalidator.OneOf("yes", "no", "hidden"),
				},
			},
			"is_cloneable": schema.BoolAttribute{
				MarkdownDescription: "Replicate this value when cloning objects.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"default": schema.StringAttribute{
				MarkdownDescription: "Default value for the field (must be a JSON value). Encapsulate strings with double quotes.",
				Optional:            true,
			},
			"weight": schema.Int64Attribute{
				MarkdownDescription: "Fields with higher weights appear lower in a form.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(100),
			},
			"validation_minimum": schema.Int64Attribute{
				MarkdownDescription: "Minimum allowed value (for numeric fields).",
				Optional:            true,
			},
			"validation_maximum": schema.Int64Attribute{
				MarkdownDescription: "Maximum allowed value (for numeric fields).",
				Optional:            true,
			},
			"validation_regex": schema.StringAttribute{
				MarkdownDescription: "Regular expression to enforce on text field values. Use ^ and $ to force matching of entire string.",
				Optional:            true,
			},
			"choice_set": schema.StringAttribute{
				MarkdownDescription: "The choice set name for select and multiselect custom fields.",
				Optional:            true,
			},
		},
	}

	// Add description and comments attributes
	maps.Copy(resp.Schema.Attributes, nbschema.CommonDescriptiveAttributes("custom field"))
}

// Configure adds the provider configured client to the resource.
func (r *CustomFieldResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create creates the resource and sets the initial Terraform state.
func (r *CustomFieldResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data CustomFieldResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating custom field", map[string]interface{}{
		"name": data.Name.ValueString(),
		"type": data.Type.ValueString(),
	})

	// Build the custom field request
	customFieldRequest, diags := r.buildCustomFieldRequest(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Call the API
	customField, httpResp, err := r.client.ExtrasAPI.ExtrasCustomFieldsCreate(ctx).WritableCustomFieldRequest(*customFieldRequest).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating custom field",
			utils.FormatAPIError(fmt.Sprintf("create custom field %s", data.Name.ValueString()), err, httpResp),
		)
		return
	}

	tflog.Debug(ctx, "Created custom field", map[string]interface{}{
		"id":   customField.GetId(),
		"name": customField.GetName(),
	})

	// Map response to state
	r.mapResponseToModel(ctx, customField, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *CustomFieldResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data CustomFieldResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Parse the ID
	customFieldID, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Custom Field ID",
			fmt.Sprintf("Custom field ID must be a number, got: %s", data.ID.ValueString()),
		)
		return
	}

	tflog.Debug(ctx, "Reading custom field", map[string]interface{}{
		"id": customFieldID,
	})

	// Call the API
	customField, httpResp, err := r.client.ExtrasAPI.ExtrasCustomFieldsRetrieve(ctx, customFieldID).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			tflog.Debug(ctx, "Custom field not found, removing from state", map[string]interface{}{
				"id": customFieldID,
			})
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading custom field",
			utils.FormatAPIError(fmt.Sprintf("read custom field ID %d", customFieldID), err, httpResp),
		)
		return
	}

	// Map response to state
	r.mapResponseToModel(ctx, customField, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates the resource and sets the updated Terraform state.
func (r *CustomFieldResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data CustomFieldResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Parse the ID
	customFieldID, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Custom Field ID",
			fmt.Sprintf("Custom field ID must be a number, got: %s", data.ID.ValueString()),
		)
		return
	}

	tflog.Debug(ctx, "Updating custom field", map[string]interface{}{
		"id":   customFieldID,
		"name": data.Name.ValueString(),
	})

	// Build the custom field request
	customFieldRequest, diags := r.buildCustomFieldRequest(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Call the API
	customField, httpResp, err := r.client.ExtrasAPI.ExtrasCustomFieldsUpdate(ctx, customFieldID).WritableCustomFieldRequest(*customFieldRequest).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating custom field",
			utils.FormatAPIError(fmt.Sprintf("update custom field ID %d", customFieldID), err, httpResp),
		)
		return
	}

	tflog.Debug(ctx, "Updated custom field", map[string]interface{}{
		"id":   customField.GetId(),
		"name": customField.GetName(),
	})

	// Map response to state
	r.mapResponseToModel(ctx, customField, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete deletes the resource and removes the Terraform state.
func (r *CustomFieldResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CustomFieldResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Parse the ID
	customFieldID, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Custom Field ID",
			fmt.Sprintf("Custom field ID must be a number, got: %s", data.ID.ValueString()),
		)
		return
	}

	tflog.Debug(ctx, "Deleting custom field", map[string]interface{}{
		"id":   customFieldID,
		"name": data.Name.ValueString(),
	})

	// Call the API
	httpResp, err := r.client.ExtrasAPI.ExtrasCustomFieldsDestroy(ctx, customFieldID).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			// Resource already deleted
			return
		}
		resp.Diagnostics.AddError(
			"Error deleting custom field",
			utils.FormatAPIError(fmt.Sprintf("delete custom field ID %d", customFieldID), err, httpResp),
		)
		return
	}

	tflog.Debug(ctx, "Deleted custom field", map[string]interface{}{
		"id": customFieldID,
	})
}

// ImportState imports the resource state.
func (r *CustomFieldResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// buildCustomFieldRequest builds a WritableCustomFieldRequest from the Terraform model.
func (r *CustomFieldResource) buildCustomFieldRequest(ctx context.Context, data *CustomFieldResourceModel) (*netbox.WritableCustomFieldRequest, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Extract object types from set
	var objectTypes []string
	diags.Append(data.ObjectTypes.ElementsAs(ctx, &objectTypes, false)...)
	if diags.HasError() {
		return nil, diags
	}

	// Create the request
	createReq := netbox.NewWritableCustomFieldRequest(objectTypes, data.Name.ValueString())

	// Build the custom field type
	cfType := netbox.PatchedWritableCustomFieldRequestType(data.Type.ValueString())
	createReq.SetType(cfType)

	// Handle related_object_type (optional, for object/multiobject types)
	if !data.RelatedObjectType.IsNull() && !data.RelatedObjectType.IsUnknown() {
		createReq.SetRelatedObjectType(data.RelatedObjectType.ValueString())
	}

	// Handle label (optional)
	if !data.Label.IsNull() && !data.Label.IsUnknown() && data.Label.ValueString() != "" {
		createReq.SetLabel(data.Label.ValueString())
	}

	// Handle group_name (optional)
	if !data.GroupName.IsNull() && !data.GroupName.IsUnknown() {
		createReq.SetGroupName(data.GroupName.ValueString())
	}

	// Apply common descriptive fields (description, comments)
	utils.ApplyDescriptiveFields(createReq, data.Description, data.Comments)

	// Handle required (optional)
	if !data.Required.IsNull() && !data.Required.IsUnknown() {
		createReq.SetRequired(data.Required.ValueBool())
	}

	// Handle search_weight (optional)
	if !data.SearchWeight.IsNull() && !data.SearchWeight.IsUnknown() {
		searchWeight, err := utils.SafeInt32FromValue(data.SearchWeight)
		if err != nil {
			diags.AddError("Invalid value", fmt.Sprintf("SearchWeight value overflow: %s", err))
			return nil, diags
		}
		createReq.SetSearchWeight(searchWeight)
	}

	// Handle filter_logic (optional)
	if !data.FilterLogic.IsNull() && !data.FilterLogic.IsUnknown() {
		filterLogic := netbox.PatchedWritableCustomFieldRequestFilterLogic(data.FilterLogic.ValueString())
		createReq.SetFilterLogic(filterLogic)
	}

	// Handle ui_visible (optional)
	if !data.UIVisible.IsNull() && !data.UIVisible.IsUnknown() {
		uiVisible := netbox.PatchedWritableCustomFieldRequestUiVisible(data.UIVisible.ValueString())
		createReq.SetUiVisible(uiVisible)
	}

	// Handle ui_editable (optional)
	if !data.UIEditable.IsNull() && !data.UIEditable.IsUnknown() {
		uiEditable := netbox.PatchedWritableCustomFieldRequestUiEditable(data.UIEditable.ValueString())
		createReq.SetUiEditable(uiEditable)
	}

	// Handle is_cloneable (optional)
	if !data.IsCloneable.IsNull() && !data.IsCloneable.IsUnknown() {
		createReq.SetIsCloneable(data.IsCloneable.ValueBool())
	}

	// Handle default (optional) - stored as JSON string
	if !data.Default.IsNull() && !data.Default.IsUnknown() {
		// Default can be any JSON value, so we store it as-is
		createReq.SetDefault(data.Default.ValueString())
	}

	// Handle weight (optional)
	if !data.Weight.IsNull() && !data.Weight.IsUnknown() {
		weight, err := utils.SafeInt32FromValue(data.Weight)
		if err != nil {
			diags.AddError("Invalid value", fmt.Sprintf("Weight value overflow: %s", err))
			return nil, diags
		}
		createReq.SetWeight(weight)
	}

	// Handle validation_minimum (optional)
	if !data.ValidationMinimum.IsNull() && !data.ValidationMinimum.IsUnknown() {
		createReq.SetValidationMinimum(data.ValidationMinimum.ValueInt64())
	}

	// Handle validation_maximum (optional)
	if !data.ValidationMaximum.IsNull() && !data.ValidationMaximum.IsUnknown() {
		createReq.SetValidationMaximum(data.ValidationMaximum.ValueInt64())
	}

	// Handle validation_regex (optional)
	if !data.ValidationRegex.IsNull() && !data.ValidationRegex.IsUnknown() {
		createReq.SetValidationRegex(data.ValidationRegex.ValueString())
	}

	// Handle choice_set (optional) - lookup by name
	if !data.ChoiceSet.IsNull() && !data.ChoiceSet.IsUnknown() {
		choiceSetReq := netbox.NewBriefCustomFieldChoiceSetRequest(data.ChoiceSet.ValueString())
		createReq.SetChoiceSet(*choiceSetReq)
	}

	return createReq, diags
}

// mapResponseToModel maps the API response to the Terraform model.
func (r *CustomFieldResource) mapResponseToModel(ctx context.Context, customField *netbox.CustomField, data *CustomFieldResourceModel, diags *diag.Diagnostics) {
	data.ID = types.StringValue(fmt.Sprintf("%d", customField.GetId()))
	data.Name = types.StringValue(customField.GetName())
	data.Type = types.StringValue(string(customField.Type.GetValue()))

	// Map object_types
	objectTypesValue, objDiags := types.SetValueFrom(ctx, types.StringType, customField.GetObjectTypes())
	diags.Append(objDiags...)
	if diags.HasError() {
		return
	}
	data.ObjectTypes = objectTypesValue

	// Map related_object_type
	if rot, ok := customField.GetRelatedObjectTypeOk(); ok && rot != nil && *rot != "" {
		data.RelatedObjectType = types.StringValue(*rot)
	} else {
		data.RelatedObjectType = types.StringNull()
	}

	// Map label
	if label, ok := customField.GetLabelOk(); ok && label != nil && *label != "" {
		data.Label = types.StringValue(*label)
	} else {
		data.Label = types.StringValue("")
	}

	// Map group_name
	if groupName, ok := customField.GetGroupNameOk(); ok && groupName != nil && *groupName != "" {
		data.GroupName = types.StringValue(*groupName)
	} else {
		data.GroupName = types.StringNull()
	}

	// Map description
	if desc, ok := customField.GetDescriptionOk(); ok && desc != nil && *desc != "" {
		data.Description = types.StringValue(*desc)
	} else {
		data.Description = types.StringNull()
	}

	// Map required
	if required, ok := customField.GetRequiredOk(); ok && required != nil {
		data.Required = types.BoolValue(*required)
	} else {
		data.Required = types.BoolValue(false)
	}

	// Map search_weight
	if searchWeight, ok := customField.GetSearchWeightOk(); ok && searchWeight != nil {
		data.SearchWeight = types.Int64Value(int64(*searchWeight))
	} else {
		data.SearchWeight = types.Int64Value(1000)
	}

	// Map filter_logic
	if customField.HasFilterLogic() && customField.FilterLogic != nil {
		data.FilterLogic = types.StringValue(string(customField.FilterLogic.GetValue()))
	} else {
		data.FilterLogic = types.StringValue("loose")
	}

	// Map ui_visible
	if customField.HasUiVisible() && customField.UiVisible != nil {
		data.UIVisible = types.StringValue(string(customField.UiVisible.GetValue()))
	} else {
		data.UIVisible = types.StringValue("always")
	}

	// Map ui_editable
	if customField.HasUiEditable() && customField.UiEditable != nil {
		data.UIEditable = types.StringValue(string(customField.UiEditable.GetValue()))
	} else {
		data.UIEditable = types.StringValue("yes")
	}

	// Map is_cloneable
	if isCloneable, ok := customField.GetIsCloneableOk(); ok && isCloneable != nil {
		data.IsCloneable = types.BoolValue(*isCloneable)
	} else {
		data.IsCloneable = types.BoolValue(false)
	}

	// Map default - convert to string representation
	if defaultVal := customField.GetDefault(); defaultVal != nil {
		// Handle different types of default values
		switch v := defaultVal.(type) {
		case string:
			data.Default = types.StringValue(v)
		default:
			data.Default = types.StringValue(fmt.Sprintf("%v", v))
		}
	} else {
		data.Default = types.StringNull()
	}

	// Map weight
	if weight, ok := customField.GetWeightOk(); ok && weight != nil {
		data.Weight = types.Int64Value(int64(*weight))
	} else {
		data.Weight = types.Int64Value(100)
	}

	// Map validation_minimum
	if valMin, ok := customField.GetValidationMinimumOk(); ok && valMin != nil {
		data.ValidationMinimum = types.Int64Value(*valMin)
	} else {
		data.ValidationMinimum = types.Int64Null()
	}

	// Map validation_maximum
	if valMax, ok := customField.GetValidationMaximumOk(); ok && valMax != nil {
		data.ValidationMaximum = types.Int64Value(*valMax)
	} else {
		data.ValidationMaximum = types.Int64Null()
	}

	// Map validation_regex
	if valRegex, ok := customField.GetValidationRegexOk(); ok && valRegex != nil && *valRegex != "" {
		data.ValidationRegex = types.StringValue(*valRegex)
	} else {
		data.ValidationRegex = types.StringNull()
	}

	// Map choice_set
	if customField.ChoiceSet.IsSet() && customField.ChoiceSet.Get() != nil {
		data.ChoiceSet = types.StringValue(customField.ChoiceSet.Get().GetName())
	} else {
		data.ChoiceSet = types.StringNull()
	}

	// Map comments
	if comments, ok := customField.GetCommentsOk(); ok && comments != nil && *comments != "" {
		data.Comments = types.StringValue(*comments)
	} else {
		data.Comments = types.StringNull()
	}
}
