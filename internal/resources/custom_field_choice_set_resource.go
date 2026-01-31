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
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &CustomFieldChoiceSetResource{}
	_ resource.ResourceWithImportState = &CustomFieldChoiceSetResource{}
)

func NewCustomFieldChoiceSetResource() resource.Resource {
	return &CustomFieldChoiceSetResource{}
}

// CustomFieldChoiceSetResource defines the resource implementation.
type CustomFieldChoiceSetResource struct {
	client *netbox.APIClient
}

// CustomFieldChoiceSetResourceModel describes the resource data model.
type CustomFieldChoiceSetResourceModel struct {
	ID                  types.String `tfsdk:"id"`
	Name                types.String `tfsdk:"name"`
	Description         types.String `tfsdk:"description"`
	BaseChoices         types.String `tfsdk:"base_choices"`
	ExtraChoices        types.List   `tfsdk:"extra_choices"`
	OrderAlphabetically types.Bool   `tfsdk:"order_alphabetically"`
}

// ChoicePairModel represents a key-value pair for choices.
type ChoicePairModel struct {
	Value types.String `tfsdk:"value"`
	Label types.String `tfsdk:"label"`
}

func (r *CustomFieldChoiceSetResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_custom_field_choice_set"
}

func (r *CustomFieldChoiceSetResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a custom field choice set in Netbox. Choice sets define the allowed values for selection custom fields.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier (assigned by Netbox).",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the choice set.",
				Required:            true,
			},
			"base_choices": schema.StringAttribute{
				MarkdownDescription: "Base choice set to inherit from. Valid values: `IATA` (Airport codes), `ISO_3166` (Country codes), `UN_LOCODE` (Location codes).",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("IATA", "ISO_3166", "UN_LOCODE"),
				},
			},
			"extra_choices": schema.ListNestedAttribute{
				MarkdownDescription: "List of extra choices. Each choice has a value and a label.",
				Required:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"value": schema.StringAttribute{
							MarkdownDescription: "The internal value stored when this choice is selected.",
							Required:            true,
						},
						"label": schema.StringAttribute{
							MarkdownDescription: "The display label shown to users.",
							Required:            true,
						},
					},
				},
			},
			"order_alphabetically": schema.BoolAttribute{
				MarkdownDescription: "Whether to order choices alphabetically. Defaults to false.",
				Optional:            true,
				Computed:            true,
			}},
	}

	// Add description attribute
	maps.Copy(resp.Schema.Attributes, nbschema.DescriptionOnlyAttributes("custom field choice set"))
}
func (r *CustomFieldChoiceSetResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *CustomFieldChoiceSetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data CustomFieldChoiceSetResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build extra_choices
	extraChoices, diags := r.buildExtraChoices(ctx, data.ExtraChoices)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	request := netbox.NewWritableCustomFieldChoiceSetRequest(
		data.Name.ValueString(),
		extraChoices,
	)

	// Set optional fields
	utils.ApplyDescription(request, data.Description)
	if !data.BaseChoices.IsNull() && !data.BaseChoices.IsUnknown() {
		baseChoices := netbox.PatchedWritableCustomFieldChoiceSetRequestBaseChoices(data.BaseChoices.ValueString())
		request.BaseChoices = &baseChoices
	}
	if !data.OrderAlphabetically.IsNull() && !data.OrderAlphabetically.IsUnknown() {
		orderAlpha := data.OrderAlphabetically.ValueBool()
		request.OrderAlphabetically = &orderAlpha
	} else if data.OrderAlphabetically.IsNull() {
		// Use AdditionalProperties to send null because of omitempty in the generated client
		if request.AdditionalProperties == nil {
			request.AdditionalProperties = make(map[string]interface{})
		}
		request.AdditionalProperties["order_alphabetically"] = nil
	}

	result, httpResp, err := r.client.ExtrasAPI.ExtrasCustomFieldChoiceSetsCreate(ctx).
		WritableCustomFieldChoiceSetRequest(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Error creating custom field choice set",
			utils.FormatAPIError("create custom field choice set", err, httpResp))
		return
	}
	if httpResp.StatusCode != http.StatusCreated {
		resp.Diagnostics.AddError("Error creating custom field choice set",
			fmt.Sprintf("Expected HTTP 201, got: %d", httpResp.StatusCode))
		return
	}
	r.mapToState(ctx, result, &data)
	tflog.Debug(ctx, "Created custom field choice set", map[string]interface{}{
		"id":   data.ID.ValueString(),
		"name": data.Name.ValueString(),
	})
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CustomFieldChoiceSetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data CustomFieldChoiceSetResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID",
			fmt.Sprintf("ID must be a number, got: %s", data.ID.ValueString()))
		return
	}
	result, httpResp, err := r.client.ExtrasAPI.ExtrasCustomFieldChoiceSetsRetrieve(ctx, id).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading custom field choice set",
			utils.FormatAPIError(fmt.Sprintf("read custom field choice set ID %s", data.ID.ValueString()), err, httpResp))
		return
	}
	r.mapToState(ctx, result, &data)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CustomFieldChoiceSetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data CustomFieldChoiceSetResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID",
			fmt.Sprintf("ID must be a number, got: %s", data.ID.ValueString()))
		return
	}

	// Build extra_choices
	extraChoices, diags := r.buildExtraChoices(ctx, data.ExtraChoices)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	request := netbox.NewWritableCustomFieldChoiceSetRequest(
		data.Name.ValueString(),
		extraChoices,
	)

	// Set optional fields
	utils.ApplyDescription(request, data.Description)
	if !data.BaseChoices.IsNull() && !data.BaseChoices.IsUnknown() {
		baseChoices := netbox.PatchedWritableCustomFieldChoiceSetRequestBaseChoices(data.BaseChoices.ValueString())
		request.BaseChoices = &baseChoices
	}
	// Note: base_choices cannot be cleared once set - API requires it to be non-empty or omitted
	if !data.OrderAlphabetically.IsNull() && !data.OrderAlphabetically.IsUnknown() {
		orderAlpha := data.OrderAlphabetically.ValueBool()
		request.OrderAlphabetically = &orderAlpha
	} else if data.OrderAlphabetically.IsNull() {
		// Use AdditionalProperties to send null because of omitempty in the generated client
		if request.AdditionalProperties == nil {
			request.AdditionalProperties = make(map[string]interface{})
		}
		request.AdditionalProperties["order_alphabetically"] = nil
	}
	result, httpResp, err := r.client.ExtrasAPI.ExtrasCustomFieldChoiceSetsUpdate(ctx, id).
		WritableCustomFieldChoiceSetRequest(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Error updating custom field choice set",
			utils.FormatAPIError(fmt.Sprintf("update custom field choice set ID %s", data.ID.ValueString()), err, httpResp))
		return
	}
	if httpResp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError("Error updating custom field choice set",
			fmt.Sprintf("Expected HTTP 200, got: %d", httpResp.StatusCode))
		return
	}
	r.mapToState(ctx, result, &data)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CustomFieldChoiceSetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CustomFieldChoiceSetResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID",
			fmt.Sprintf("ID must be a number, got: %s", data.ID.ValueString()))
		return
	}
	httpResp, err := r.client.ExtrasAPI.ExtrasCustomFieldChoiceSetsDestroy(ctx, id).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		// If the resource was already deleted (404), consider it a success
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			tflog.Debug(ctx, "Custom field choice set already deleted", map[string]interface{}{"id": id})
			return
		}
		resp.Diagnostics.AddError("Error deleting custom field choice set",
			utils.FormatAPIError(fmt.Sprintf("delete custom field choice set ID %s", data.ID.ValueString()), err, httpResp))
		return
	}
	if httpResp.StatusCode != http.StatusNoContent {
		resp.Diagnostics.AddError("Error deleting custom field choice set",
			fmt.Sprintf("Expected HTTP 204, got: %d", httpResp.StatusCode))
		return
	}
}

func (r *CustomFieldChoiceSetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	utils.ImportStatePassthroughIDWithValidation(ctx, req, resp, path.Root("id"), true)
}

// buildExtraChoices converts the Terraform list of choice pairs to the API format.
func (r *CustomFieldChoiceSetResource) buildExtraChoices(ctx context.Context, choicesList types.List) ([][]interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics
	if choicesList.IsNull() || choicesList.IsUnknown() {
		return [][]interface{}{}, diags
	}
	var choicePairs []ChoicePairModel
	d := choicesList.ElementsAs(ctx, &choicePairs, false)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}
	result := make([][]interface{}, len(choicePairs))
	for i, pair := range choicePairs {
		result[i] = []interface{}{pair.Value.ValueString(), pair.Label.ValueString()}
	}
	return result, diags
}

// mapToState maps API response to Terraform state.
func (r *CustomFieldChoiceSetResource) mapToState(ctx context.Context, result *netbox.CustomFieldChoiceSet, data *CustomFieldChoiceSetResourceModel) {
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

	// Map extra_choices from [][]interface{} to list of ChoicePairModel
	extraChoices := result.GetExtraChoices()
	if len(extraChoices) > 0 {
		choicePairs := make([]ChoicePairModel, len(extraChoices))
		for i, pair := range extraChoices {
			if len(pair) >= 2 {
				choicePairs[i] = ChoicePairModel{
					Value: types.StringValue(fmt.Sprintf("%v", pair[0])),
					Label: types.StringValue(fmt.Sprintf("%v", pair[1])),
				}
			}
		}
		choicesValue, _ := types.ListValueFrom(ctx, types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"value": types.StringType,
				"label": types.StringType,
			},
		}, choicePairs)
		data.ExtraChoices = choicesValue
	} else {
		data.ExtraChoices = types.ListNull(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"value": types.StringType,
				"label": types.StringType,
			},
		})
	}
	if result.HasOrderAlphabetically() {
		data.OrderAlphabetically = types.BoolValue(result.GetOrderAlphabetically())
	} else {
		data.OrderAlphabetically = types.BoolNull()
	}
}
