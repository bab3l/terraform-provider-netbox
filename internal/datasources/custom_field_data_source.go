// Package datasources contains Terraform data source implementations for NetBox objects.
package datasources

import (
	"context"
	"fmt"

	"github.com/bab3l/go-netbox"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ datasource.DataSource              = &CustomFieldDataSource{}
	_ datasource.DataSourceWithConfigure = &CustomFieldDataSource{}
)

// NewCustomFieldDataSource returns a new data source implementing the CustomField data source.
func NewCustomFieldDataSource() datasource.DataSource {
	return &CustomFieldDataSource{}
}

// CustomFieldDataSource defines the data source implementation.
type CustomFieldDataSource struct {
	client *netbox.APIClient
}

// CustomFieldDataSourceModel describes the data source data model.
type CustomFieldDataSourceModel struct {
	ID                types.String `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	ObjectTypes       types.Set    `tfsdk:"object_types"`
	Type              types.String `tfsdk:"type"`
	RelatedObjectType types.String `tfsdk:"related_object_type"`
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
	DataType          types.String `tfsdk:"data_type"`
	DisplayName       types.String `tfsdk:"display_name"`
}

// Metadata returns the data source type name.
func (d *CustomFieldDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_custom_field"
}

// Schema defines the schema for the data source.
func (d *CustomFieldDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves information about a custom field in NetBox.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique numeric ID of the custom field. Use this to look up by ID.",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The internal name of the custom field. Use this to look up by name.",
				Optional:            true,
				Computed:            true,
			},
			"object_types": schema.SetAttribute{
				MarkdownDescription: "The object types this custom field applies to.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "The type of custom field.",
				Computed:            true,
			},
			"related_object_type": schema.StringAttribute{
				MarkdownDescription: "The related object type for object and multiobject custom fields.",
				Computed:            true,
			},
			"label": schema.StringAttribute{
				MarkdownDescription: "Name of the field as displayed to users.",
				Computed:            true,
			},
			"group_name": schema.StringAttribute{
				MarkdownDescription: "Custom fields within the same group will be displayed together.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "A description of the custom field.",
				Computed:            true,
			},
			"required": schema.BoolAttribute{
				MarkdownDescription: "If true, this field is required when creating new objects or editing an existing object.",
				Computed:            true,
			},
			"search_weight": schema.Int64Attribute{
				MarkdownDescription: "Weighting for search.",
				Computed:            true,
			},
			"filter_logic": schema.StringAttribute{
				MarkdownDescription: "Filter logic for the custom field.",
				Computed:            true,
			},
			"ui_visible": schema.StringAttribute{
				MarkdownDescription: "UI visibility setting.",
				Computed:            true,
			},
			"ui_editable": schema.StringAttribute{
				MarkdownDescription: "UI editability setting.",
				Computed:            true,
			},
			"is_cloneable": schema.BoolAttribute{
				MarkdownDescription: "Replicate this value when cloning objects.",
				Computed:            true,
			},
			"default": schema.StringAttribute{
				MarkdownDescription: "Default value for the field (JSON string).",
				Computed:            true,
			},
			"weight": schema.Int64Attribute{
				MarkdownDescription: "Fields with higher weights appear lower in a form.",
				Computed:            true,
			},
			"validation_minimum": schema.Int64Attribute{
				MarkdownDescription: "Minimum allowed value (for numeric fields).",
				Computed:            true,
			},
			"validation_maximum": schema.Int64Attribute{
				MarkdownDescription: "Maximum allowed value (for numeric fields).",
				Computed:            true,
			},
			"validation_regex": schema.StringAttribute{
				MarkdownDescription: "Regular expression to enforce on text field values.",
				Computed:            true,
			},
			"choice_set": schema.StringAttribute{
				MarkdownDescription: "The choice set name for select and multiselect custom fields.",
				Computed:            true,
			},
			"comments": schema.StringAttribute{
				MarkdownDescription: "Comments or notes about the custom field.",
				Computed:            true,
			},
			"data_type": schema.StringAttribute{
				MarkdownDescription: "The data type of the custom field.",
				Computed:            true,
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "The display name of the custom field.",
				Computed:            true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *CustomFieldDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// Read refreshes the data source data.
func (d *CustomFieldDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CustomFieldDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var customField *netbox.CustomField

	// Look up by ID if provided
	switch {
	case !data.ID.IsNull() && !data.ID.IsUnknown():
		customFieldID, err := utils.ParseID(data.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Invalid Custom Field ID",
				fmt.Sprintf("Custom field ID must be a number, got: %s", data.ID.ValueString()),
			)
			return
		}

		tflog.Debug(ctx, "Reading custom field by ID", map[string]interface{}{
			"id": customFieldID,
		})

		cf, httpResp, err := d.client.ExtrasAPI.ExtrasCustomFieldsRetrieve(ctx, customFieldID).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading custom field",
				utils.FormatAPIError(fmt.Sprintf("read custom field ID %d", customFieldID), err, httpResp),
			)
			return
		}
		customField = cf
	case !data.Name.IsNull() && !data.Name.IsUnknown():
		// Look up by name
		tflog.Debug(ctx, "Reading custom field by name", map[string]interface{}{
			"name": data.Name.ValueString(),
		})

		listResp, httpResp, err := d.client.ExtrasAPI.ExtrasCustomFieldsList(ctx).Name([]string{data.Name.ValueString()}).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading custom field",
				utils.FormatAPIError(fmt.Sprintf("read custom field by name %s", data.Name.ValueString()), err, httpResp),
			)
			return
		}

		if listResp.GetCount() == 0 {
			resp.Diagnostics.AddError(
				"Custom field not found",
				fmt.Sprintf("No custom field found with name: %s", data.Name.ValueString()),
			)
			return
		}

		if listResp.GetCount() > 1 {
			resp.Diagnostics.AddError(
				"Multiple custom fields found",
				fmt.Sprintf("Found %d custom fields with name: %s", listResp.GetCount(), data.Name.ValueString()),
			)
			return
		}

		customField = &listResp.GetResults()[0]
	default:
		resp.Diagnostics.AddError(
			"Missing Required Attribute",
			"Either 'id' or 'name' must be specified to look up a custom field.",
		)
		return
	}

	// Map response to model
	d.mapResponseToModel(ctx, customField, &data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// mapResponseToModel maps the API response to the Terraform model.
func (d *CustomFieldDataSource) mapResponseToModel(ctx context.Context, customField *netbox.CustomField, data *CustomFieldDataSourceModel) {
	data.ID = types.StringValue(fmt.Sprintf("%d", customField.GetId()))
	data.Name = types.StringValue(customField.GetName())
	data.Type = types.StringValue(string(customField.Type.GetValue()))
	data.DataType = types.StringValue(customField.GetDataType())

	// Map object_types
	objectTypesValue, _ := types.SetValueFrom(ctx, types.StringType, customField.GetObjectTypes())
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
		data.Label = types.StringNull()
	}

	// Map group_name
	if groupName, ok := customField.GetGroupNameOk(); ok && groupName != nil && *groupName != "" {
		data.GroupName = types.StringValue(*groupName)
	} else {
		data.GroupName = types.StringNull()
	}

	// Map description
	if desc, ok := customField.GetDescriptionOk(); ok && desc != nil {
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
		data.SearchWeight = types.Int64Null()
	}

	// Map filter_logic
	if customField.HasFilterLogic() && customField.FilterLogic != nil {
		data.FilterLogic = types.StringValue(string(customField.FilterLogic.GetValue()))
	} else {
		data.FilterLogic = types.StringNull()
	}

	// Map ui_visible
	if customField.HasUiVisible() && customField.UiVisible != nil {
		data.UIVisible = types.StringValue(string(customField.UiVisible.GetValue()))
	} else {
		data.UIVisible = types.StringNull()
	}

	// Map ui_editable
	if customField.HasUiEditable() && customField.UiEditable != nil {
		data.UIEditable = types.StringValue(string(customField.UiEditable.GetValue()))
	} else {
		data.UIEditable = types.StringNull()
	}

	// Map is_cloneable
	if isCloneable, ok := customField.GetIsCloneableOk(); ok && isCloneable != nil {
		data.IsCloneable = types.BoolValue(*isCloneable)
	} else {
		data.IsCloneable = types.BoolValue(false)
	}

	// Map default - convert to string representation
	if defaultVal := customField.GetDefault(); defaultVal != nil {
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
		data.Weight = types.Int64Null()
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
	if comments, ok := customField.GetCommentsOk(); ok && comments != nil {
		data.Comments = types.StringValue(*comments)
	} else {
		data.Comments = types.StringNull()
	}

	// Map display_name
	if customField.GetDisplay() != "" {
		data.DisplayName = types.StringValue(customField.GetDisplay())
	} else {
		data.DisplayName = types.StringNull()
	}
}
