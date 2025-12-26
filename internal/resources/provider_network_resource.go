// Package resources contains Terraform resource implementations for the Netbox provider.

package resources

import (
	"context"
	"fmt"

	"github.com/bab3l/go-netbox"
	"github.com/bab3l/terraform-provider-netbox/internal/netboxlookup"
	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
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

// Ensure provider defined types fully satisfy framework interfaces.

var (
	_ resource.Resource = &ProviderNetworkResource{}

	_ resource.ResourceWithConfigure = &ProviderNetworkResource{}

	_ resource.ResourceWithImportState = &ProviderNetworkResource{}
)

// NewProviderNetworkResource returns a new ProviderNetwork resource.

func NewProviderNetworkResource() resource.Resource {

	return &ProviderNetworkResource{}

}

// ProviderNetworkResource defines the resource implementation.

type ProviderNetworkResource struct {
	client *netbox.APIClient
}

// ProviderNetworkResourceModel describes the resource data model.

type ProviderNetworkResourceModel struct {
	ID types.String `tfsdk:"id"`

	CircuitProvider types.String `tfsdk:"circuit_provider"`

	Name types.String `tfsdk:"name"`

	ServiceID types.String `tfsdk:"service_id"`

	Description types.String `tfsdk:"description"`

	Comments types.String `tfsdk:"comments"`

	DisplayName types.String `tfsdk:"display_name"`

	Tags types.Set `tfsdk:"tags"`

	CustomFields types.Set `tfsdk:"custom_fields"`
}

// Metadata returns the resource type name.

func (r *ProviderNetworkResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {

	resp.TypeName = req.ProviderTypeName + "_provider_network"

}

// Schema defines the schema for the resource.

func (r *ProviderNetworkResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	resp.Schema = schema.Schema{

		MarkdownDescription: "Manages a provider network in NetBox. Provider networks represent the network infrastructure of circuit providers.",

		Attributes: map[string]schema.Attribute{

			"id": schema.StringAttribute{

				MarkdownDescription: "The unique numeric ID of the provider network.",

				Computed: true,

				PlanModifiers: []planmodifier.String{

					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"circuit_provider": schema.StringAttribute{

				MarkdownDescription: "The circuit provider that owns this network. Can be specified by name, slug, or ID.",

				Required: true,
			},

			"name": schema.StringAttribute{

				MarkdownDescription: "The name of the provider network.",

				Required: true,

				Validators: []validator.String{

					stringvalidator.LengthBetween(1, 100),
				},
			},

			"service_id": schema.StringAttribute{

				MarkdownDescription: "A unique identifier for this network provided by the circuit provider.",

				Optional: true,

				Validators: []validator.String{

					stringvalidator.LengthAtMost(100),
				},
			},

			"description": schema.StringAttribute{

				MarkdownDescription: "A description of the provider network.",

				Optional: true,

				Validators: []validator.String{

					stringvalidator.LengthAtMost(200),
				},
			},

			"comments": schema.StringAttribute{

				MarkdownDescription: "Additional comments or notes about this provider network.",

				Optional: true,
			},

			"display_name": nbschema.DisplayNameAttribute("provider network"),

			"tags": nbschema.TagsAttribute(),

			"custom_fields": nbschema.CustomFieldsAttribute(),
		},
	}

}

// Configure adds the provider configured client to the resource.

func (r *ProviderNetworkResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {

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

func (r *ProviderNetworkResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var data ProviderNetworkResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	tflog.Debug(ctx, "Creating provider network", map[string]interface{}{

		"circuit_provider": data.CircuitProvider.ValueString(),

		"name": data.Name.ValueString(),
	})

	// Build the provider network request

	pnRequest, diags := r.buildProviderNetworkRequest(ctx, &data)

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {

		return

	}

	// Call the API

	pn, httpResp, err := r.client.CircuitsAPI.CircuitsProviderNetworksCreate(ctx).ProviderNetworkRequest(*pnRequest).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		resp.Diagnostics.AddError(

			"Error creating provider network",

			utils.FormatAPIError(fmt.Sprintf("create provider network %s", data.Name.ValueString()), err, httpResp),
		)

		return

	}

	tflog.Debug(ctx, "Created provider network", map[string]interface{}{

		"id": pn.GetId(),

		"name": pn.GetName(),
	})

	// Map response to state

	r.mapResponseToModel(ctx, pn, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {

		return

	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

// Read refreshes the Terraform state with the latest data.

func (r *ProviderNetworkResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	var data ProviderNetworkResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	// Parse the ID

	pnID, err := utils.ParseID(data.ID.ValueString())

	if err != nil {

		resp.Diagnostics.AddError(

			"Invalid Provider Network ID",

			fmt.Sprintf("Provider network ID must be a number, got: %s", data.ID.ValueString()),
		)

		return

	}

	tflog.Debug(ctx, "Reading provider network", map[string]interface{}{

		"id": pnID,
	})

	// Call the API

	pn, httpResp, err := r.client.CircuitsAPI.CircuitsProviderNetworksRetrieve(ctx, pnID).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		if httpResp != nil && httpResp.StatusCode == 404 {

			tflog.Debug(ctx, "Provider network not found, removing from state", map[string]interface{}{

				"id": pnID,
			})

			resp.State.RemoveResource(ctx)

			return

		}

		resp.Diagnostics.AddError(

			"Error reading provider network",

			utils.FormatAPIError(fmt.Sprintf("read provider network ID %d", pnID), err, httpResp),
		)

		return

	}

	// Map response to state

	r.mapResponseToModel(ctx, pn, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {

		return

	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

// Update updates the resource and sets the updated Terraform state.

func (r *ProviderNetworkResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var data ProviderNetworkResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	// Parse the ID

	pnID, err := utils.ParseID(data.ID.ValueString())

	if err != nil {

		resp.Diagnostics.AddError(

			"Invalid Provider Network ID",

			fmt.Sprintf("Provider network ID must be a number, got: %s", data.ID.ValueString()),
		)

		return

	}

	tflog.Debug(ctx, "Updating provider network", map[string]interface{}{

		"id": pnID,

		"name": data.Name.ValueString(),
	})

	// Build the provider network request

	pnRequest, diags := r.buildProviderNetworkRequest(ctx, &data)

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {

		return

	}

	// Call the API

	pn, httpResp, err := r.client.CircuitsAPI.CircuitsProviderNetworksUpdate(ctx, pnID).ProviderNetworkRequest(*pnRequest).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		resp.Diagnostics.AddError(

			"Error updating provider network",

			utils.FormatAPIError(fmt.Sprintf("update provider network ID %d", pnID), err, httpResp),
		)

		return

	}

	tflog.Debug(ctx, "Updated provider network", map[string]interface{}{

		"id": pn.GetId(),

		"name": pn.GetName(),
	})

	// Map response to state

	r.mapResponseToModel(ctx, pn, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {

		return

	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

// Delete deletes the resource and removes the Terraform state.

func (r *ProviderNetworkResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var data ProviderNetworkResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	// Parse the ID

	pnID, err := utils.ParseID(data.ID.ValueString())

	if err != nil {

		resp.Diagnostics.AddError(

			"Invalid Provider Network ID",

			fmt.Sprintf("Provider network ID must be a number, got: %s", data.ID.ValueString()),
		)

		return

	}

	tflog.Debug(ctx, "Deleting provider network", map[string]interface{}{

		"id": pnID,

		"name": data.Name.ValueString(),
	})

	// Call the API

	httpResp, err := r.client.CircuitsAPI.CircuitsProviderNetworksDestroy(ctx, pnID).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		if httpResp != nil && httpResp.StatusCode == 404 {

			// Resource already deleted

			return

		}

		resp.Diagnostics.AddError(

			"Error deleting provider network",

			utils.FormatAPIError(fmt.Sprintf("delete provider network ID %d", pnID), err, httpResp),
		)

		return

	}

	tflog.Debug(ctx, "Deleted provider network", map[string]interface{}{

		"id": pnID,
	})

}

// ImportState imports the resource state.

func (r *ProviderNetworkResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

}

// buildProviderNetworkRequest builds a ProviderNetworkRequest from the Terraform model.

func (r *ProviderNetworkResource) buildProviderNetworkRequest(ctx context.Context, data *ProviderNetworkResourceModel) (*netbox.ProviderNetworkRequest, diag.Diagnostics) {

	var diags diag.Diagnostics

	// Lookup provider (required)

	provider, providerDiags := netboxlookup.LookupProvider(ctx, r.client, data.CircuitProvider.ValueString())

	diags.Append(providerDiags...)

	if diags.HasError() {

		return nil, diags

	}

	// Create the request with required fields

	pnRequest := netbox.NewProviderNetworkRequest(*provider, data.Name.ValueString())

	// Handle service_id (optional)

	if !data.ServiceID.IsNull() && !data.ServiceID.IsUnknown() {

		serviceID := data.ServiceID.ValueString()

		pnRequest.ServiceId = &serviceID

	}

	// Handle description (optional)

	if !data.Description.IsNull() && !data.Description.IsUnknown() {

		desc := data.Description.ValueString()

		pnRequest.Description = &desc

	}

	// Handle comments (optional)

	if !data.Comments.IsNull() && !data.Comments.IsUnknown() {

		comments := data.Comments.ValueString()

		pnRequest.Comments = &comments

	}

	// Handle tags

	if !data.Tags.IsNull() && !data.Tags.IsUnknown() {

		tags, tagDiags := utils.TagModelsToNestedTagRequests(ctx, data.Tags)

		diags.Append(tagDiags...)

		if diags.HasError() {

			return nil, diags

		}

		pnRequest.Tags = tags

	}

	// Handle custom fields

	if !data.CustomFields.IsNull() && !data.CustomFields.IsUnknown() {

		var customFieldModels []utils.CustomFieldModel

		cfDiags := data.CustomFields.ElementsAs(ctx, &customFieldModels, false)

		diags.Append(cfDiags...)

		if diags.HasError() {

			return nil, diags

		}

		pnRequest.CustomFields = utils.CustomFieldModelsToMap(customFieldModels)

	}

	return pnRequest, diags

}

// mapResponseToModel maps the API response to the Terraform model.

func (r *ProviderNetworkResource) mapResponseToModel(ctx context.Context, pn *netbox.ProviderNetwork, data *ProviderNetworkResourceModel, diags *diag.Diagnostics) {

	data.ID = types.StringValue(fmt.Sprintf("%d", pn.GetId()))

	data.Name = types.StringValue(pn.GetName())

	// Map Provider (use ID to match what was passed in)

	data.CircuitProvider = types.StringValue(fmt.Sprintf("%d", pn.Provider.GetId()))

	// Map service_id

	data.ServiceID = utils.StringFromAPI(pn.HasServiceId(), pn.GetServiceId, data.ServiceID)

	// Map description

	data.Description = utils.StringFromAPI(pn.HasDescription(), pn.GetDescription, data.Description)

	// Map comments

	data.Comments = utils.StringFromAPI(pn.HasComments(), pn.GetComments, data.Comments)

	// Map display_name
	if pn.Display != "" {
		data.DisplayName = types.StringValue(pn.Display)
	} else {
		data.DisplayName = types.StringNull()
	}

	// Handle tags

	if pn.HasTags() {

		tags := utils.NestedTagsToTagModels(pn.GetTags())

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

	if pn.HasCustomFields() {

		apiCustomFields := pn.GetCustomFields()

		var stateCustomFieldModels []utils.CustomFieldModel

		if !data.CustomFields.IsNull() {

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
