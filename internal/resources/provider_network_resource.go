// Package resources contains Terraform resource implementations for the Netbox provider.

package resources

import (
	"context"
	"fmt"
	"maps"

	"github.com/bab3l/go-netbox"
	"github.com/bab3l/terraform-provider-netbox/internal/netboxlookup"
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

// Ensure provider defined types fully satisfy framework interfaces.

var (
	_ resource.Resource = &ProviderNetworkResource{}

	_ resource.ResourceWithConfigure = &ProviderNetworkResource{}

	_ resource.ResourceWithImportState = &ProviderNetworkResource{}
	_ resource.ResourceWithIdentity    = &ProviderNetworkResource{}
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
		},
	}

	// Add common descriptive attributes (description, comments)
	maps.Copy(resp.Schema.Attributes, nbschema.CommonDescriptiveAttributes("provider network"))

	// Add tags and custom fields attributes
	resp.Schema.Attributes["tags"] = nbschema.TagsSlugAttribute()
	resp.Schema.Attributes["custom_fields"] = nbschema.CustomFieldsAttribute()
}

func (r *ProviderNetworkResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = nbschema.ImportIdentityWithCustomFieldsSchema()
}

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

	pnRequest, diags := r.buildProviderNetworkRequest(ctx, &data, nil)

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

	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(data.ID.ValueString()), data.CustomFields, &resp.Diagnostics)
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

	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(data.ID.ValueString()), data.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates the resource and sets the updated Terraform state.

func (r *ProviderNetworkResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, data ProviderNetworkResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
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

	// Build the provider network request (pass state for merge-aware custom fields)

	pnRequest, diags := r.buildProviderNetworkRequest(ctx, &data, &state)

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

	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(data.ID.ValueString()), data.CustomFields, &resp.Diagnostics)
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
	if parsed, ok := utils.ParseImportIdentityCustomFields(ctx, req.Identity, &resp.Diagnostics); ok {
		if resp.Diagnostics.HasError() {
			return
		}
		if parsed.ID == "" {
			resp.Diagnostics.AddError("Invalid import identity", "Identity id must be provided")
			return
		}

		pnID, err := utils.ParseID(parsed.ID)
		if err != nil {
			resp.Diagnostics.AddError("Invalid Provider Network ID", fmt.Sprintf("Provider network ID must be a number, got: %s", parsed.ID))
			return
		}
		pn, httpResp, err := r.client.CircuitsAPI.CircuitsProviderNetworksRetrieve(ctx, pnID).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError("Error importing provider network", utils.FormatAPIError(fmt.Sprintf("read provider network ID %d", pnID), err, httpResp))
			return
		}

		var data ProviderNetworkResourceModel
		if pn.GetProvider().Id != 0 {
			provider := pn.GetProvider()
			data.CircuitProvider = types.StringValue(fmt.Sprintf("%d", provider.GetId()))
		}
		if pn.HasTags() {
			tagSlugs := make([]string, 0, len(pn.GetTags()))
			for _, tag := range pn.GetTags() {
				tagSlugs = append(tagSlugs, tag.GetSlug())
			}
			data.Tags = utils.TagsSlugToSet(ctx, tagSlugs)
		} else {
			data.Tags = types.SetNull(types.StringType)
		}
		if parsed.HasCustomFields {
			if len(parsed.CustomFields) == 0 {
				data.CustomFields = types.SetValueMust(utils.GetCustomFieldsAttributeType().ElemType, []attr.Value{})
			} else {
				ownedSet, setDiags := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, parsed.CustomFields)
				resp.Diagnostics.Append(setDiags...)
				if resp.Diagnostics.HasError() {
					return
				}
				data.CustomFields = ownedSet
			}
		} else {
			data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)
		}

		r.mapResponseToModel(ctx, pn, &data, &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}
		if parsed.HasCustomFields {
			data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, data.CustomFields, pn.GetCustomFields(), &resp.Diagnostics)
		} else {
			data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)
		}
		if resp.Diagnostics.HasError() {
			return
		}

		if resp.Identity != nil {
			listValue, listDiags := types.ListValueFrom(ctx, types.StringType, parsed.CustomFieldItems)
			resp.Diagnostics.Append(listDiags...)
			if resp.Diagnostics.HasError() {
				return
			}
			resp.Diagnostics.Append(resp.Identity.Set(ctx, &utils.ImportIdentityCustomFieldsModel{
				ID:           types.StringValue(parsed.ID),
				CustomFields: listValue,
			})...)
		}

		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
		return
	}

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// buildProviderNetworkRequest builds a ProviderNetworkRequest from the Terraform model.

func (r *ProviderNetworkResource) buildProviderNetworkRequest(ctx context.Context, data *ProviderNetworkResourceModel, state *ProviderNetworkResourceModel) (*netbox.ProviderNetworkRequest, diag.Diagnostics) {
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
	} else if data.ServiceID.IsNull() {
		serviceID := ""
		pnRequest.ServiceId = &serviceID
	}

	// Handle description and comments - explicitly handle null values for removal
	if data.Description.IsNull() {
		pnRequest.SetDescription("")
	} else if !data.Description.IsUnknown() {
		pnRequest.SetDescription(data.Description.ValueString())
	}

	if data.Comments.IsNull() {
		pnRequest.SetComments("")
	} else if !data.Comments.IsUnknown() {
		pnRequest.SetComments(data.Comments.ValueString())
	}

	// Handle tags (from slug list)
	utils.ApplyTagsFromSlugs(ctx, r.client, pnRequest, data.Tags, &diags)
	if diags.HasError() {
		return nil, diags
	}

	// Handle custom fields - merge-aware for updates
	var stateCustomFields types.Set
	if state != nil {
		stateCustomFields = state.CustomFields
	}
	utils.ApplyCustomFieldsWithMerge(ctx, pnRequest, data.CustomFields, stateCustomFields, &diags)

	if diags.HasError() {
		return nil, diags
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

	// Filter tags to owned (slug list format)
	switch {
	case data.Tags.IsNull():
		data.Tags = types.SetNull(types.StringType)
	case len(data.Tags.Elements()) == 0:
		data.Tags = types.SetValueMust(types.StringType, []attr.Value{})
	case pn.HasTags():
		var tagSlugs []string
		for _, tag := range pn.GetTags() {
			tagSlugs = append(tagSlugs, tag.GetSlug())
		}
		data.Tags = utils.TagsSlugToSet(ctx, tagSlugs)
	default:
		data.Tags = types.SetValueMust(types.StringType, []attr.Value{})
	}

	// Populate custom fields
	if pn.HasCustomFields() {
		data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, data.CustomFields, pn.GetCustomFields(), diags)
	}
}
