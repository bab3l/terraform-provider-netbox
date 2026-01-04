// Package resources provides Terraform resource implementations for NetBox objects.
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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &EventRuleResource{}
	_ resource.ResourceWithConfigure   = &EventRuleResource{}
	_ resource.ResourceWithImportState = &EventRuleResource{}
)

// NewEventRuleResource returns a new resource implementing the event rule resource.
func NewEventRuleResource() resource.Resource {
	return &EventRuleResource{}
}

// EventRuleResource defines the resource implementation.
type EventRuleResource struct {
	client *netbox.APIClient
}

// EventRuleResourceModel describes the resource data model.
type EventRuleResourceModel struct {
	ID               types.String `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	ObjectTypes      types.Set    `tfsdk:"object_types"`
	EventTypes       types.Set    `tfsdk:"event_types"`
	Enabled          types.Bool   `tfsdk:"enabled"`
	Conditions       types.String `tfsdk:"conditions"`
	ActionType       types.String `tfsdk:"action_type"`
	ActionObjectType types.String `tfsdk:"action_object_type"`
	ActionObjectID   types.String `tfsdk:"action_object_id"`
	Description      types.String `tfsdk:"description"`
	Tags             types.Set    `tfsdk:"tags"`
	CustomFields     types.Set    `tfsdk:"custom_fields"`
}

// Metadata returns the resource type name.
func (r *EventRuleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_event_rule"
}

// Schema defines the schema for the resource.
func (r *EventRuleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an event rule in NetBox. Event rules define actions to be executed automatically when certain events occur on specific object types.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique numeric ID of the event rule.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": nbschema.NameAttribute("event rule", 150),
			"object_types": schema.SetAttribute{
				MarkdownDescription: "The object types that this event rule applies to (e.g., `dcim.device`, `ipam.ipaddress`).",
				Required:            true,
				ElementType:         types.StringType,
			},
			"event_types": schema.SetAttribute{
				MarkdownDescription: "The types of events which will trigger this rule. Valid values: `object_created`, `object_updated`, `object_deleted`, `job_started`, `job_completed`, `job_failed`, `job_errored`.",
				Required:            true,
				ElementType:         types.StringType,
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Whether the event rule is enabled. Defaults to `true`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"conditions": schema.StringAttribute{
				MarkdownDescription: "A JSON object defining conditions which determine whether the event will be generated. Leave empty for no conditions.",
				Optional:            true,
			},
			"action_type": schema.StringAttribute{
				MarkdownDescription: "The type of action to execute. Valid values: `webhook`, `script`, `notification`.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("webhook", "script", "notification"),
				},
			},
			"action_object_type": schema.StringAttribute{
				MarkdownDescription: "The content type of the action object (e.g., `extras.webhook`, `extras.script`, `extras.notificationgroup`).",
				Required:            true,
			},
			"action_object_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the action object (webhook, script, or notification group).",
				Optional:            true,
			},
		},
	}

	// Add description attribute
	maps.Copy(resp.Schema.Attributes, nbschema.DescriptionOnlyAttributes("event rule"))

	// Add common metadata attributes (tags, custom_fields)
	maps.Copy(resp.Schema.Attributes, nbschema.CommonMetadataAttributes())
}

// Configure adds the provider configured client to the resource.
func (r *EventRuleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *EventRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data EventRuleResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating event rule", map[string]interface{}{
		"name": data.Name.ValueString(),
	})

	// Extract object types
	var objectTypes []string
	resp.Diagnostics.Append(data.ObjectTypes.ElementsAs(ctx, &objectTypes, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Extract event types
	var eventTypeStrings []string
	resp.Diagnostics.Append(data.EventTypes.ElementsAs(ctx, &eventTypeStrings, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert event type strings to the proper type
	eventTypes := make([]netbox.EventRuleEventTypesInner, len(eventTypeStrings))
	for i, et := range eventTypeStrings {
		eventTypes[i] = netbox.EventRuleEventTypesInner(et)
	}

	// Build the request
	request := netbox.NewWritableEventRuleRequest(
		objectTypes,
		data.Name.ValueString(),
		eventTypes,
		data.ActionObjectType.ValueString(),
	)

	// Set optional fields
	if !data.Enabled.IsNull() && !data.Enabled.IsUnknown() {
		enabled := data.Enabled.ValueBool()
		request.Enabled = &enabled
	}

	if !data.Conditions.IsNull() && !data.Conditions.IsUnknown() {
		// Parse conditions as JSON
		conditionsStr := data.Conditions.ValueString()
		if conditionsStr != "" {
			request.Conditions = conditionsStr
		}
	}

	if !data.ActionType.IsNull() && !data.ActionType.IsUnknown() {
		actionType := netbox.EventRuleActionTypeValue(data.ActionType.ValueString())
		request.ActionType = &actionType
	}

	if !data.ActionObjectID.IsNull() && !data.ActionObjectID.IsUnknown() {
		actionObjectID, err := utils.ParseID64(data.ActionObjectID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Invalid Action Object ID", fmt.Sprintf("Action Object ID must be a number, got: %s", data.ActionObjectID.ValueString()))
			return
		}
		request.ActionObjectId = *netbox.NewNullableInt64(&actionObjectID)
	}

	// Apply common fields (description, tags, custom_fields)
	utils.ApplyDescription(request, data.Description)
	utils.ApplyTags(ctx, request, data.Tags, &resp.Diagnostics)
	utils.ApplyCustomFields(ctx, request, data.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the event rule
	result, httpResp, err := r.client.ExtrasAPI.ExtrasEventRulesCreate(ctx).
		WritableEventRuleRequest(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Error Creating Event Rule",
			utils.FormatAPIError("create event rule", err, httpResp))
		return
	}

	// Map the response to state
	r.mapToState(ctx, result, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Created event rule", map[string]interface{}{
		"id": data.ID.ValueString(),
	})
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *EventRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data EventRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("ID must be a number, got: %s", data.ID.ValueString()))
		return
	}
	result, httpResp, err := r.client.ExtrasAPI.ExtrasEventRulesRetrieve(ctx, id).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			tflog.Debug(ctx, "Event rule not found, removing from state", map[string]interface{}{"id": id})
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error Reading Event Rule",
			utils.FormatAPIError(fmt.Sprintf("read event rule ID %d", id), err, httpResp))
		return
	}
	r.mapToState(ctx, result, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates the resource.
func (r *EventRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data EventRuleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("ID must be a number, got: %s", data.ID.ValueString()))
		return
	}
	tflog.Debug(ctx, "Updating event rule", map[string]interface{}{
		"id":   id,
		"name": data.Name.ValueString(),
	})

	// Extract object types
	var objectTypes []string
	resp.Diagnostics.Append(data.ObjectTypes.ElementsAs(ctx, &objectTypes, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Extract event types
	var eventTypeStrings []string
	resp.Diagnostics.Append(data.EventTypes.ElementsAs(ctx, &eventTypeStrings, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert event type strings to the proper type
	eventTypes := make([]netbox.EventRuleEventTypesInner, len(eventTypeStrings))
	for i, et := range eventTypeStrings {
		eventTypes[i] = netbox.EventRuleEventTypesInner(et)
	}

	// Build the request
	request := netbox.NewWritableEventRuleRequest(
		objectTypes,
		data.Name.ValueString(),
		eventTypes,
		data.ActionObjectType.ValueString(),
	)

	// Set optional fields (same as Create)
	if !data.Enabled.IsNull() && !data.Enabled.IsUnknown() {
		enabled := data.Enabled.ValueBool()
		request.Enabled = &enabled
	}

	if !data.Conditions.IsNull() && !data.Conditions.IsUnknown() {
		conditionsStr := data.Conditions.ValueString()
		if conditionsStr != "" {
			request.Conditions = conditionsStr
		}
	}

	if !data.ActionType.IsNull() && !data.ActionType.IsUnknown() {
		actionType := netbox.EventRuleActionTypeValue(data.ActionType.ValueString())
		request.ActionType = &actionType
	}

	if !data.ActionObjectID.IsNull() && !data.ActionObjectID.IsUnknown() {
		actionObjectID, err := utils.ParseID64(data.ActionObjectID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Invalid Action Object ID", fmt.Sprintf("Action Object ID must be a number, got: %s", data.ActionObjectID.ValueString()))
			return
		}
		request.ActionObjectId = *netbox.NewNullableInt64(&actionObjectID)
	}

	// Apply common fields (description, tags, custom_fields)
	utils.ApplyDescription(request, data.Description)
	utils.ApplyTags(ctx, request, data.Tags, &resp.Diagnostics)
	utils.ApplyCustomFields(ctx, request, data.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update the event rule
	result, httpResp, err := r.client.ExtrasAPI.ExtrasEventRulesUpdate(ctx, id).
		WritableEventRuleRequest(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Error Updating Event Rule",
			utils.FormatAPIError(fmt.Sprintf("update event rule ID %d", id), err, httpResp))
		return
	}

	// Map the response to state
	r.mapToState(ctx, result, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete deletes the resource.
func (r *EventRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data EventRuleResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("ID must be a number, got: %s", data.ID.ValueString()))
		return
	}

	tflog.Debug(ctx, "Deleting event rule", map[string]interface{}{"id": id})

	httpResp, err := r.client.ExtrasAPI.ExtrasEventRulesDestroy(ctx, id).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError("Error Deleting Event Rule",
			utils.FormatAPIError(fmt.Sprintf("delete event rule ID %d", id), err, httpResp))
		return
	}
}

// ImportState imports the resource state.
func (r *EventRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// mapToState maps the API response to the Terraform state.
func (r *EventRuleResource) mapToState(ctx context.Context, result *netbox.EventRule, data *EventRuleResourceModel, diags *diag.Diagnostics) {
	data.ID = types.StringValue(fmt.Sprintf("%d", result.GetId()))
	data.Name = types.StringValue(result.GetName())

	// Map object types
	objectTypesValue, objDiags := types.SetValueFrom(ctx, types.StringType, result.GetObjectTypes())
	diags.Append(objDiags...)
	data.ObjectTypes = objectTypesValue

	// Map event types
	eventTypesRaw := result.GetEventTypes()
	eventTypeStrings := make([]string, len(eventTypesRaw))
	for i, et := range eventTypesRaw {
		eventTypeStrings[i] = string(et)
	}
	eventTypesValue, etDiags := types.SetValueFrom(ctx, types.StringType, eventTypeStrings)
	diags.Append(etDiags...)
	data.EventTypes = eventTypesValue

	// Map enabled
	if result.HasEnabled() {
		data.Enabled = types.BoolValue(result.GetEnabled())
	} else {
		data.Enabled = types.BoolValue(true)
	}

	// Map conditions
	if result.HasConditions() && result.GetConditions() != nil {
		// Conditions is an interface{}, serialize to JSON string
		conditionsJSON, err := utils.ToJSONString(result.GetConditions())
		if err == nil && conditionsJSON != "" && conditionsJSON != "null" {
			data.Conditions = types.StringValue(conditionsJSON)
		} else {
			data.Conditions = types.StringNull()
		}
	} else {
		data.Conditions = types.StringNull()
	}

	// Map action type - always present (not optional)
	actionType := result.GetActionType()
	data.ActionType = types.StringValue(string(actionType.GetValue()))

	// Map action object type
	data.ActionObjectType = types.StringValue(result.GetActionObjectType())

	// Map action object ID
	if result.HasActionObjectId() {
		actionObjID, ok := result.GetActionObjectIdOk()
		if ok && actionObjID != nil {
			data.ActionObjectID = types.StringValue(fmt.Sprintf("%d", *actionObjID))
		} else {
			data.ActionObjectID = types.StringNull()
		}
	} else {
		data.ActionObjectID = types.StringNull()
	}

	// Map description
	if result.HasDescription() && result.GetDescription() != "" {
		data.Description = types.StringValue(result.GetDescription())
	} else {
		data.Description = types.StringNull()
	}

	// Handle tags using consolidated helper
	data.Tags = utils.PopulateTagsFromAPI(ctx, result.HasTags(), result.GetTags(), data.Tags, diags)
	if diags.HasError() {
		return
	}

	// Handle custom fields using consolidated helper
	data.CustomFields = utils.PopulateCustomFieldsFromAPI(ctx, result.HasCustomFields(), result.GetCustomFields(), data.CustomFields, diags)
}
