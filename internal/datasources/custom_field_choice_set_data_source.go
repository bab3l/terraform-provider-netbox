package datasources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bab3l/go-netbox"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ datasource.DataSource = &CustomFieldChoiceSetDataSource{}

func NewCustomFieldChoiceSetDataSource() datasource.DataSource {
	return &CustomFieldChoiceSetDataSource{}
}

// CustomFieldChoiceSetDataSource defines the data source implementation.
type CustomFieldChoiceSetDataSource struct {
	client *netbox.APIClient
}

// CustomFieldChoiceSetDataSourceModel describes the data source data model.
type CustomFieldChoiceSetDataSourceModel struct {
	ID                  types.String `tfsdk:"id"`
	Name                types.String `tfsdk:"name"`
	Description         types.String `tfsdk:"description"`
	BaseChoices         types.String `tfsdk:"base_choices"`
	ExtraChoices        types.List   `tfsdk:"extra_choices"`
	OrderAlphabetically types.Bool   `tfsdk:"order_alphabetically"`
	ChoicesCount        types.Int64  `tfsdk:"choices_count"`
}

func (d *CustomFieldChoiceSetDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_custom_field_choice_set"
}

func (d *CustomFieldChoiceSetDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get information about a custom field choice set in Netbox.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier for the choice set. Use to look up by ID.",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the choice set. Use to look up by name.",
				Optional:            true,
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the choice set.",
				Computed:            true,
			},
			"base_choices": schema.StringAttribute{
				MarkdownDescription: "Base choice set. Values: IATA, ISO_3166, UN_LOCODE.",
				Computed:            true,
			},
			"extra_choices": schema.ListNestedAttribute{
				MarkdownDescription: "List of extra choices.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"value": schema.StringAttribute{
							MarkdownDescription: "The internal value.",
							Computed:            true,
						},
						"label": schema.StringAttribute{
							MarkdownDescription: "The display label.",
							Computed:            true,
						},
					},
				},
			},
			"order_alphabetically": schema.BoolAttribute{
				MarkdownDescription: "Whether choices are ordered alphabetically.",
				Computed:            true,
			},
			"choices_count": schema.Int64Attribute{
				MarkdownDescription: "Total number of choices available.",
				Computed:            true,
			},
		},
	}
}

func (d *CustomFieldChoiceSetDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *CustomFieldChoiceSetDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CustomFieldChoiceSetDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var result *netbox.CustomFieldChoiceSet
	var httpResp *http.Response
	var err error

	switch {
	case !data.ID.IsNull() && data.ID.ValueString() != "":
		// Lookup by ID
		id, parseErr := utils.ParseID(data.ID.ValueString())
		if parseErr != nil {
			resp.Diagnostics.AddError("Invalid ID", "ID must be a number")
			return
		}
		result, httpResp, err = d.client.ExtrasAPI.ExtrasCustomFieldChoiceSetsRetrieve(ctx, id).Execute()
		defer utils.CloseResponseBody(httpResp)
	case !data.Name.IsNull() && data.Name.ValueString() != "":
		// Lookup by name
		list, listResp, listErr := d.client.ExtrasAPI.ExtrasCustomFieldChoiceSetsList(ctx).
			Name([]string{data.Name.ValueString()}).Execute()
		defer utils.CloseResponseBody(listResp)
		httpResp = listResp
		err = listErr
		if err == nil && list != nil {
			results := list.GetResults()
			if len(results) == 0 {
				resp.Diagnostics.AddError("Not Found",
					fmt.Sprintf("No custom field choice set found with name: %s", data.Name.ValueString()))
				return
			}
			if len(results) > 1 {
				resp.Diagnostics.AddError("Multiple Found",
					fmt.Sprintf("Multiple custom field choice sets found with name: %s. Please use ID instead.", data.Name.ValueString()))
				return
			}
			result = &results[0]
		}
	default:
		resp.Diagnostics.AddError("Missing Identifier", "Either 'id' or 'name' must be specified")
		return
	}

	if err != nil {
		resp.Diagnostics.AddError("Error reading custom field choice set",
			utils.FormatAPIError("read custom field choice set", err, httpResp))
		return
	}

	d.mapToState(ctx, result, &data)

	tflog.Debug(ctx, "Read custom field choice set", map[string]interface{}{
		"id":   data.ID.ValueString(),
		"name": data.Name.ValueString(),
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// mapToState maps API response to Terraform state
func (d *CustomFieldChoiceSetDataSource) mapToState(ctx context.Context, result *netbox.CustomFieldChoiceSet, data *CustomFieldChoiceSetDataSourceModel) {
	data.ID = types.StringValue(fmt.Sprintf("%d", result.GetId()))
	data.Name = types.StringValue(result.GetName())

	if result.HasDescription() && result.GetDescription() != "" {
		data.Description = types.StringValue(result.GetDescription())
	} else {
		data.Description = types.StringNull()
	}

	if result.HasBaseChoices() {
		baseChoices := result.GetBaseChoices()
		if baseChoices.Value != nil {
			data.BaseChoices = types.StringValue(string(*baseChoices.Value))
		} else {
			data.BaseChoices = types.StringNull()
		}
	} else {
		data.BaseChoices = types.StringNull()
	}

	// Map extra_choices
	choiceObjectType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"value": types.StringType,
			"label": types.StringType,
		},
	}
	extraChoices := result.GetExtraChoices()
	if len(extraChoices) > 0 {
		choiceValues := make([]attr.Value, len(extraChoices))
		for i, pair := range extraChoices {
			if len(pair) >= 2 {
				choiceValues[i], _ = types.ObjectValue(
					choiceObjectType.AttrTypes,
					map[string]attr.Value{
						"value": types.StringValue(fmt.Sprintf("%v", pair[0])),
						"label": types.StringValue(fmt.Sprintf("%v", pair[1])),
					},
				)
			}
		}
		choicesValue, _ := types.ListValue(choiceObjectType, choiceValues)
		data.ExtraChoices = choicesValue
	} else {
		data.ExtraChoices = types.ListNull(choiceObjectType)
	}

	if result.HasOrderAlphabetically() {
		data.OrderAlphabetically = types.BoolValue(result.GetOrderAlphabetically())
	} else {
		data.OrderAlphabetically = types.BoolNull()
	}

	data.ChoicesCount = types.Int64Value(int64(result.GetChoicesCount()))
}
