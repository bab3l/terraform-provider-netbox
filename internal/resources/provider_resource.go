// Package resources contains Terraform resource implementations for the Netbox provider.

package resources

import (
	"context"
	"fmt"
	"maps"

	"github.com/bab3l/go-netbox"
	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.

var (
	_ resource.Resource = &ProviderResource{}

	_ resource.ResourceWithConfigure = &ProviderResource{}

	_ resource.ResourceWithImportState = &ProviderResource{}
)

// NewProviderResource returns a new Provider resource (circuit provider, not Terraform provider).

func NewProviderResource() resource.Resource {
	return &ProviderResource{}
}

// ProviderResource defines the resource implementation for circuit providers.

type ProviderResource struct {
	client *netbox.APIClient
}

// ProviderResourceModel describes the resource data model.

type ProviderResourceModel struct {
	ID types.String `tfsdk:"id"`

	Name types.String `tfsdk:"name"`

	Slug types.String `tfsdk:"slug"`

	Description types.String `tfsdk:"description"`

	Comments types.String `tfsdk:"comments"`

	Tags types.Set `tfsdk:"tags"`

	CustomFields types.Set `tfsdk:"custom_fields"`
}

// Metadata returns the resource type name.

func (r *ProviderResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_provider"
}

// Schema defines the schema for the resource.

func (r *ProviderResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a circuit provider in Netbox. Providers represent the organizations (ISPs, carriers, etc.) that provide circuit connectivity services.",

		Attributes: map[string]schema.Attribute{
			"id": nbschema.IDAttribute("circuit provider"),

			"name": nbschema.NameAttribute("circuit provider", 100),

			"slug": nbschema.SlugAttribute("circuit provider"),
		},
	}

	// Add common descriptive attributes (description, comments)
	maps.Copy(resp.Schema.Attributes, nbschema.CommonDescriptiveAttributes("circuit provider"))

	// Add common metadata attributes (tags, custom_fields)
	maps.Copy(resp.Schema.Attributes, nbschema.CommonMetadataAttributes())
}

// Configure sets up the resource with the provider client.

func (r *ProviderResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// mapProviderToState maps a Provider from the API to the Terraform state model.

func (r *ProviderResource) mapProviderToState(ctx context.Context, provider *netbox.Provider, data *ProviderResourceModel, diags *diag.Diagnostics) {
	data.ID = types.StringValue(fmt.Sprintf("%d", provider.GetId()))

	data.Name = types.StringValue(provider.GetName())

	data.Slug = types.StringValue(provider.GetSlug())

	// Handle description

	if provider.HasDescription() && provider.GetDescription() != "" {
		data.Description = types.StringValue(provider.GetDescription())
	} else if !data.Description.IsNull() {
		data.Description = types.StringNull()
	}

	// Handle comments

	if provider.HasComments() && provider.GetComments() != "" {
		data.Comments = types.StringValue(provider.GetComments())
	} else if !data.Comments.IsNull() {
		data.Comments = types.StringNull()
	}

	// Handle tags

	if provider.HasTags() {
		tags := utils.NestedTagsToTagModels(provider.GetTags())

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

	if provider.HasCustomFields() && !data.CustomFields.IsNull() {
		var stateCustomFields []utils.CustomFieldModel

		cfDiags := data.CustomFields.ElementsAs(ctx, &stateCustomFields, false)

		diags.Append(cfDiags...)

		if diags.HasError() {
			return
		}

		customFields := utils.MapToCustomFieldModels(provider.GetCustomFields(), stateCustomFields)

		customFieldsValue, cfValueDiags := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, customFields)

		diags.Append(cfValueDiags...)

		if diags.HasError() {
			return
		}

		data.CustomFields = customFieldsValue
	}
}

// Create creates a new provider resource.

func (r *ProviderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ProviderResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating circuit provider", map[string]interface{}{
		"name": data.Name.ValueString(),

		"slug": data.Slug.ValueString(),
	})

	// Build the provider request

	providerRequest := netbox.ProviderRequest{
		Name: data.Name.ValueString(),

		Slug: data.Slug.ValueString(),
	}

	// Apply common fields (description, comments, tags, custom_fields)

	utils.ApplyCommonFields(ctx, &providerRequest, data.Description, data.Comments, data.Tags, data.CustomFields, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	// Call the API

	provider, httpResp, err := r.client.CircuitsAPI.CircuitsProvidersCreate(ctx).ProviderRequest(providerRequest).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {
		resp.Diagnostics.AddError(

			"Error creating circuit provider",

			utils.FormatAPIError("create circuit provider", err, httpResp),
		)

		return
	}

	tflog.Debug(ctx, "Created circuit provider", map[string]interface{}{
		"id": provider.GetId(),

		"name": provider.GetName(),
	})

	// Map response to state

	r.mapProviderToState(ctx, provider, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read reads the provider resource.

func (r *ProviderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ProviderResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	id, err := utils.ParseID(data.ID.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Could not parse provider ID: %s", err))

		return
	}

	tflog.Debug(ctx, "Reading circuit provider", map[string]interface{}{
		"id": id,
	})

	provider, httpResp, err := r.client.CircuitsAPI.CircuitsProvidersRetrieve(ctx, id).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			tflog.Debug(ctx, "Circuit provider not found, removing from state", map[string]interface{}{
				"id": id,
			})

			resp.State.RemoveResource(ctx)

			return
		}

		resp.Diagnostics.AddError(

			"Error reading circuit provider",

			utils.FormatAPIError("read circuit provider", err, httpResp),
		)

		return
	}

	// Map response to state

	r.mapProviderToState(ctx, provider, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates the provider resource.

func (r *ProviderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ProviderResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	id, err := utils.ParseID(data.ID.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Could not parse provider ID: %s", err))

		return
	}

	tflog.Debug(ctx, "Updating circuit provider", map[string]interface{}{
		"id": id,

		"name": data.Name.ValueString(),
	})

	// Build the provider request

	providerRequest := netbox.ProviderRequest{
		Name: data.Name.ValueString(),

		Slug: data.Slug.ValueString(),
	}

	// Apply common fields (description, comments, tags, custom_fields)

	utils.ApplyCommonFields(ctx, &providerRequest, data.Description, data.Comments, data.Tags, data.CustomFields, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	// Call the API

	provider, httpResp, err := r.client.CircuitsAPI.CircuitsProvidersUpdate(ctx, id).ProviderRequest(providerRequest).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {
		resp.Diagnostics.AddError(

			"Error updating circuit provider",

			utils.FormatAPIError("update circuit provider", err, httpResp),
		)

		return
	}

	tflog.Debug(ctx, "Updated circuit provider", map[string]interface{}{
		"id": provider.GetId(),

		"name": provider.GetName(),
	})

	// Map response to state

	r.mapProviderToState(ctx, provider, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete deletes the provider resource.

func (r *ProviderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ProviderResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	id, err := utils.ParseID(data.ID.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Could not parse provider ID: %s", err))

		return
	}

	tflog.Debug(ctx, "Deleting circuit provider", map[string]interface{}{
		"id": id,
	})

	httpResp, err := r.client.CircuitsAPI.CircuitsProvidersDestroy(ctx, id).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			tflog.Debug(ctx, "Circuit provider already deleted", map[string]interface{}{
				"id": id,
			})

			return
		}

		resp.Diagnostics.AddError(

			"Error deleting circuit provider",

			utils.FormatAPIError("delete circuit provider", err, httpResp),
		)

		return
	}

	tflog.Debug(ctx, "Deleted circuit provider", map[string]interface{}{
		"id": id,
	})
}

// ImportState imports the resource state.

func (r *ProviderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
