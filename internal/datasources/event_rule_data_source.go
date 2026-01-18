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
	_ datasource.DataSource              = &EventRuleDataSource{}
	_ datasource.DataSourceWithConfigure = &EventRuleDataSource{}
)

// NewEventRuleDataSource returns a new data source implementing the event rule data source.
func NewEventRuleDataSource() datasource.DataSource {
	return &EventRuleDataSource{}
}

// EventRuleDataSource defines the data source implementation.
type EventRuleDataSource struct {
	client *netbox.APIClient
}

// EventRuleDataSourceModel describes the data source model.
type EventRuleDataSourceModel struct {
	ID               types.String `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	ObjectTypes      types.Set    `tfsdk:"object_types"`
	EventTypes       types.Set    `tfsdk:"event_types"`
	Enabled          types.Bool   `tfsdk:"enabled"`
	Conditions       types.String `tfsdk:"conditions"`
	ActionType       types.String `tfsdk:"action_type"`
	ActionObjectType types.String `tfsdk:"action_object_type"`
	ActionObjectID   types.String `tfsdk:"action_object_id"`
	ActionObject     types.String `tfsdk:"action_object"`
	Description      types.String `tfsdk:"description"`
	DisplayName      types.String `tfsdk:"display_name"`
	Tags             types.Set    `tfsdk:"tags"`
	CustomFields     types.Set    `tfsdk:"custom_fields"`
	Created          types.String `tfsdk:"created"`
	LastUpdated      types.String `tfsdk:"last_updated"`
}

// Metadata returns the data source type name.
func (d *EventRuleDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_event_rule"
}

// Schema defines the schema for the data source.
func (d *EventRuleDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get information about an event rule in NetBox.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique numeric ID of the event rule to look up.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the event rule.",
				Computed:            true,
			},
			"object_types": schema.SetAttribute{
				MarkdownDescription: "The object types that this event rule applies to.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"event_types": schema.SetAttribute{
				MarkdownDescription: "The types of events which will trigger this rule.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Whether the event rule is enabled.",
				Computed:            true,
			},
			"conditions": schema.StringAttribute{
				MarkdownDescription: "A JSON object defining conditions which determine whether the event will be generated.",
				Computed:            true,
			},
			"action_type": schema.StringAttribute{
				MarkdownDescription: "The type of action to execute (webhook, script, notification).",
				Computed:            true,
			},
			"action_object_type": schema.StringAttribute{
				MarkdownDescription: "The content type of the action object.",
				Computed:            true,
			},
			"action_object_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the action object.",
				Computed:            true,
			},
			"action_object": schema.StringAttribute{
				MarkdownDescription: "The name of the action object.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "A description of the event rule.",
				Computed:            true,
			},
			"display_name":  nbschema.DSComputedStringAttribute("The display name of the event rule."),
			"tags":          nbschema.DSTagsAttribute(),
			"custom_fields": nbschema.DSCustomFieldsAttribute(),
			"created": schema.StringAttribute{
				MarkdownDescription: "The timestamp of when the event rule was created.",
				Computed:            true,
			},
			"last_updated": schema.StringAttribute{
				MarkdownDescription: "The timestamp of when the event rule was last updated.",
				Computed:            true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *EventRuleDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *EventRuleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data EventRuleDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.ID.IsNull() || data.ID.IsUnknown() || data.ID.ValueString() == "" {
		resp.Diagnostics.AddError("Missing required identifier", "The 'id' attribute must be specified to lookup an event rule.")
		return
	}

	id, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("ID must be a number, got: %s", data.ID.ValueString()))
		return
	}

	tflog.Debug(ctx, "Reading event rule", map[string]interface{}{"id": id})

	result, httpResp, err := d.client.ExtrasAPI.ExtrasEventRulesRetrieve(ctx, id).Execute()
	defer utils.CloseResponseBody(httpResp)

	if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
		resp.Diagnostics.AddError(
			"Event Rule Not Found",
			fmt.Sprintf("No event rule found with ID: %d", id),
		)
		return
	}

	if err != nil {
		resp.Diagnostics.AddError("Error Reading Event Rule",
			utils.FormatAPIError(fmt.Sprintf("read event rule ID %d", id), err, httpResp))
		return
	}

	d.mapToDataSourceModel(ctx, result, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// mapToDataSourceModel maps the API response to the data source model.
func (d *EventRuleDataSource) mapToDataSourceModel(ctx context.Context, result *netbox.EventRule, data *EventRuleDataSourceModel, diags *diag.Diagnostics) {
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
	data.Enabled = types.BoolValue(result.GetEnabled())

	// Map conditions
	if result.HasConditions() && result.GetConditions() != nil {
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

	// Map action object name - always present (not optional)
	actionObj := result.GetActionObject()
	if actionObj != nil {
		// This is a generic map, extract name if available
		if name, hasName := actionObj["name"]; hasName {
			if nameStr, isString := name.(string); isString {
				data.ActionObject = types.StringValue(nameStr)
			} else {
				data.ActionObject = types.StringNull()
			}
		} else {
			data.ActionObject = types.StringNull()
		}
	} else {
		data.ActionObject = types.StringNull()
	}

	// Map description
	if result.HasDescription() && result.GetDescription() != "" {
		data.Description = types.StringValue(result.GetDescription())
	} else {
		data.Description = types.StringNull()
	}

	// Map display_name
	if result.GetDisplay() != "" {
		data.DisplayName = types.StringValue(result.GetDisplay())
	} else {
		data.DisplayName = types.StringNull()
	}

	// Map tags
	if result.HasTags() {
		tags := utils.NestedTagsToTagModels(result.GetTags())
		tagsValue, tagDiags := types.SetValueFrom(ctx, utils.GetTagsAttributeType().ElemType, tags)
		diags.Append(tagDiags...)
		data.Tags = tagsValue
	} else {
		data.Tags = types.SetNull(utils.GetTagsAttributeType().ElemType)
	}

	// Map custom fields
	if result.HasCustomFields() {
		customFields := utils.MapAllCustomFieldsToModels(result.GetCustomFields())
		customFieldsValue, cfDiags := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, customFields)
		if !cfDiags.HasError() {
			data.CustomFields = customFieldsValue
		}
		diags.Append(cfDiags...)
	} else {
		data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)
	}

	// Map timestamps - Created and LastUpdated are NullableTime, use IsSet() pattern
	createdOk, createdIsSet := result.GetCreatedOk()
	if createdIsSet && createdOk != nil {
		data.Created = types.StringValue(createdOk.String())
	} else {
		data.Created = types.StringNull()
	}

	lastUpdatedOk, lastUpdatedIsSet := result.GetLastUpdatedOk()
	if lastUpdatedIsSet && lastUpdatedOk != nil {
		data.LastUpdated = types.StringValue(lastUpdatedOk.String())
	} else {
		data.LastUpdated = types.StringNull()
	}
}
