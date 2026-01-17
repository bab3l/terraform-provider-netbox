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
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.

var (
	_ resource.Resource = &ProviderAccountResource{}

	_ resource.ResourceWithConfigure = &ProviderAccountResource{}

	_ resource.ResourceWithImportState = &ProviderAccountResource{}
)

// NewProviderAccountResource returns a new Provider Account resource.

func NewProviderAccountResource() resource.Resource {
	return &ProviderAccountResource{}
}

// ProviderAccountResource defines the resource implementation.

type ProviderAccountResource struct {
	client *netbox.APIClient
}

// ProviderAccountResourceModel describes the resource data model.

type ProviderAccountResourceModel struct {
	ID types.String `tfsdk:"id"`

	CircuitProvider types.String `tfsdk:"circuit_provider"`

	Name types.String `tfsdk:"name"`

	Account types.String `tfsdk:"account"`

	Description types.String `tfsdk:"description"`

	Comments types.String `tfsdk:"comments"`

	Tags types.Set `tfsdk:"tags"`

	CustomFields types.Set `tfsdk:"custom_fields"`
}

// Metadata returns the resource type name.

func (r *ProviderAccountResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_provider_account"
}

// Schema defines the schema for the resource.

func (r *ProviderAccountResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a provider account in Netbox. Provider accounts represent accounts with circuit providers.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique numeric ID of the provider account.",

				Computed: true,

				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"circuit_provider": schema.StringAttribute{
				MarkdownDescription: "The name, slug, or ID of the circuit provider this account belongs to.",

				Required: true,
			},

			"name": schema.StringAttribute{
				MarkdownDescription: "An optional name for this provider account.",

				Optional: true,
			},

			"account": schema.StringAttribute{
				MarkdownDescription: "The account identifier (e.g., account number or ID).",

				Required: true,
			},
		},
	}

	// Add common descriptive attributes (description, comments)
	maps.Copy(resp.Schema.Attributes, nbschema.CommonDescriptiveAttributes("provider account"))

	// Add common metadata attributes (tags, custom_fields)
	maps.Copy(resp.Schema.Attributes, nbschema.CommonMetadataAttributes())
}

func (r *ProviderAccountResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.

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

// Create creates a new provider account resource.

func (r *ProviderAccountResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ProviderAccountResourceModel

	// Read Terraform plan data into the model

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Build the create request

	createReq, diags := r.buildCreateRequest(ctx, &data, nil)

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating provider account", map[string]interface{}{
		"account": data.Account.ValueString(),
	})

	// Call API to create provider account

	providerAccount, httpResp, err := r.client.CircuitsAPI.CircuitsProviderAccountsCreate(ctx).ProviderAccountRequest(*createReq).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {
		resp.Diagnostics.AddError(

			"Error creating provider account",

			fmt.Sprintf("Could not create provider account: %s\nHTTP Response: %v", err.Error(), httpResp),
		)

		return
	}

	// Map response to model

	r.mapResponseToModel(ctx, providerAccount, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Created provider account", map[string]interface{}{
		"id": data.ID.ValueString(),
	})

	// Save data into Terraform state

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read reads the provider account resource.

func (r *ProviderAccountResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ProviderAccountResourceModel

	// Read Terraform prior state data into the model

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	id, err := utils.ParseID(data.ID.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(

			"Invalid ID",

			fmt.Sprintf("Could not convert ID to integer: %s", err.Error()),
		)

		return
	}

	tflog.Debug(ctx, "Reading provider account", map[string]interface{}{
		"id": id,
	})

	// Call API to read provider account

	providerAccount, httpResp, err := r.client.CircuitsAPI.CircuitsProviderAccountsRetrieve(ctx, id).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			tflog.Debug(ctx, "Provider account not found, removing from state", map[string]interface{}{
				"id": id,
			})

			resp.State.RemoveResource(ctx)

			return
		}

		resp.Diagnostics.AddError(

			"Error reading provider account",

			fmt.Sprintf("Could not read provider account: %s\nHTTP Response: %v", err.Error(), httpResp),
		)

		return
	}

	// Map response to model

	r.mapResponseToModel(ctx, providerAccount, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates the provider account resource.

func (r *ProviderAccountResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, data ProviderAccountResourceModel

	// Read Terraform plan and state data into the models

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	id, err := utils.ParseID(data.ID.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(

			"Invalid ID",

			fmt.Sprintf("Could not convert ID to integer: %s", err.Error()),
		)

		return
	}

	// Build the update request (pass state for merge-aware custom fields)

	updateReq, diags := r.buildCreateRequest(ctx, &data, &state)

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Updating provider account", map[string]interface{}{
		"id": id,
	})

	// Call API to update provider account

	providerAccount, httpResp, err := r.client.CircuitsAPI.CircuitsProviderAccountsUpdate(ctx, id).ProviderAccountRequest(*updateReq).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {
		resp.Diagnostics.AddError(

			"Error updating provider account",

			fmt.Sprintf("Could not update provider account: %s\nHTTP Response: %v", err.Error(), httpResp),
		)

		return
	}

	// Map response to model

	r.mapResponseToModel(ctx, providerAccount, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Updated provider account", map[string]interface{}{
		"id": data.ID.ValueString(),
	})

	// Save updated data into Terraform state

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete deletes the provider account resource.

func (r *ProviderAccountResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ProviderAccountResourceModel

	// Read Terraform prior state data into the model

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	id, err := utils.ParseID(data.ID.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(

			"Invalid ID",

			fmt.Sprintf("Could not convert ID to integer: %s", err.Error()),
		)

		return
	}

	tflog.Debug(ctx, "Deleting provider account", map[string]interface{}{
		"id": id,
	})

	// Call API to delete provider account

	httpResp, err := r.client.CircuitsAPI.CircuitsProviderAccountsDestroy(ctx, id).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			tflog.Debug(ctx, "Provider account already deleted", map[string]interface{}{
				"id": id,
			})

			return
		}

		resp.Diagnostics.AddError(

			"Error deleting provider account",

			fmt.Sprintf("Could not delete provider account: %s\nHTTP Response: %v", err.Error(), httpResp),
		)

		return
	}

	tflog.Debug(ctx, "Deleted provider account", map[string]interface{}{
		"id": id,
	})
}

// ImportState imports an existing provider account.

func (r *ProviderAccountResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// buildCreateRequest builds a ProviderAccountRequest from the model.

func (r *ProviderAccountResource) buildCreateRequest(ctx context.Context, data *ProviderAccountResourceModel, state *ProviderAccountResourceModel) (*netbox.ProviderAccountRequest, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Look up Provider (required)

	provider, providerDiags := netboxlookup.LookupProvider(ctx, r.client, data.CircuitProvider.ValueString())

	diags.Append(providerDiags...)

	if diags.HasError() {
		return nil, diags
	}

	createReq := netbox.NewProviderAccountRequest(*provider, data.Account.ValueString())

	// Handle name (optional)

	if !data.Name.IsNull() && !data.Name.IsUnknown() {
		createReq.SetName(data.Name.ValueString())
	}

	// Handle description and comments, tags and custom fields - merge-aware
	var stateTags, stateCustomFields types.Set
	if state != nil {
		stateTags = state.Tags
		stateCustomFields = state.CustomFields
	}

	utils.ApplyCommonFieldsWithMerge(ctx, createReq, data.Description, data.Comments, data.Tags, stateTags, data.CustomFields, stateCustomFields, &diags)

	if diags.HasError() {
		return nil, diags
	}

	return createReq, diags
}

// mapResponseToModel maps the API response to the Terraform model.

func (r *ProviderAccountResource) mapResponseToModel(ctx context.Context, providerAccount *netbox.ProviderAccount, data *ProviderAccountResourceModel, diags *diag.Diagnostics) {
	data.ID = types.StringValue(fmt.Sprintf("%d", providerAccount.GetId()))

	data.Account = types.StringValue(providerAccount.GetAccount())

	// Map Provider

	if provider := providerAccount.GetProvider(); provider.Id != 0 {
		data.CircuitProvider = types.StringValue(fmt.Sprintf("%d", provider.Id))
	}

	// Map name

	data.Name = utils.StringFromAPI(providerAccount.HasName(), providerAccount.GetName, data.Name)

	// Map description

	data.Description = utils.StringFromAPI(providerAccount.HasDescription(), providerAccount.GetDescription, data.Description)

	// Map comments

	data.Comments = utils.StringFromAPI(providerAccount.HasComments(), providerAccount.GetComments, data.Comments)

	// Populate tags and custom fields using unified helpers
	data.Tags = utils.PopulateTagsFromAPI(ctx, providerAccount.HasTags(), providerAccount.GetTags(), data.Tags, diags)
	if providerAccount.HasCustomFields() {
		data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, data.CustomFields, providerAccount.GetCustomFields(), diags)
	}
}
