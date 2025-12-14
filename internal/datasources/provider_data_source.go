// Package datasources contains Terraform data source implementations for the Netbox provider.
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
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &ProviderDataSource{}

// NewProviderDataSource returns a new Provider data source (circuit provider, not Terraform provider).
func NewProviderDataSource() datasource.DataSource {
	return &ProviderDataSource{}
}

// ProviderDataSource defines the data source implementation for circuit providers.
type ProviderDataSource struct {
	client *netbox.APIClient
}

// ProviderDataSourceModel describes the data source data model.
type ProviderDataSourceModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Slug         types.String `tfsdk:"slug"`
	Description  types.String `tfsdk:"description"`
	Comments     types.String `tfsdk:"comments"`
	Tags         types.Set    `tfsdk:"tags"`
	CustomFields types.Set    `tfsdk:"custom_fields"`
}

// Metadata returns the data source type name.
func (d *ProviderDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_provider"
}

// Schema defines the schema for the data source.
func (d *ProviderDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves information about a circuit provider in Netbox. Providers represent the organizations (ISPs, carriers, etc.) that provide circuit connectivity services.",

		Attributes: map[string]schema.Attribute{
			"id":            nbschema.DSIDAttribute("circuit provider"),
			"name":          nbschema.DSNameAttribute("circuit provider"),
			"slug":          nbschema.DSSlugAttribute("circuit provider"),
			"description":   nbschema.DSComputedStringAttribute("Description of the circuit provider."),
			"comments":      nbschema.DSComputedStringAttribute("Additional comments or notes about the circuit provider."),
			"tags":          nbschema.DSTagsAttribute(),
			"custom_fields": nbschema.DSCustomFieldsAttribute(),
		},
	}
}

// Configure sets up the data source with the provider client.
func (d *ProviderDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *ProviderDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ProviderDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var provider *netbox.Provider
	var err error
	var httpResp *http.Response

	// Determine if we're searching by ID, slug, or name
	switch {
	case !data.ID.IsNull():
		// Search by ID
		providerID := data.ID.ValueString()
		tflog.Debug(ctx, "Reading circuit provider by ID", map[string]interface{}{
			"id": providerID,
		})

		var providerIDInt int32
		if _, parseErr := fmt.Sscanf(providerID, "%d", &providerIDInt); parseErr != nil {
			resp.Diagnostics.AddError(
				"Invalid Provider ID",
				fmt.Sprintf("Provider ID must be a number, got: %s", providerID),
			)
			return
		}

		provider, httpResp, err = d.client.CircuitsAPI.CircuitsProvidersRetrieve(ctx, providerIDInt).Execute()
		defer utils.CloseResponseBody(httpResp)
	case !data.Slug.IsNull():
		// Search by slug
		providerSlug := data.Slug.ValueString()
		tflog.Debug(ctx, "Reading circuit provider by slug", map[string]interface{}{
			"slug": providerSlug,
		})

		var providers *netbox.PaginatedProviderList
		providers, httpResp, err = d.client.CircuitsAPI.CircuitsProvidersList(ctx).Slug([]string{providerSlug}).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading circuit provider",
				utils.FormatAPIError("read circuit provider by slug", err, httpResp),
			)
			return
		}
		if len(providers.GetResults()) == 0 {
			resp.Diagnostics.AddError(
				"Provider Not Found",
				fmt.Sprintf("No circuit provider found with slug: %s", providerSlug),
			)
			return
		}
		provider = &providers.GetResults()[0]
	case !data.Name.IsNull():
		// Search by name
		providerName := data.Name.ValueString()
		tflog.Debug(ctx, "Reading circuit provider by name", map[string]interface{}{
			"name": providerName,
		})

		var providers *netbox.PaginatedProviderList
		providers, httpResp, err = d.client.CircuitsAPI.CircuitsProvidersList(ctx).Name([]string{providerName}).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading circuit provider",
				utils.FormatAPIError("read circuit provider by name", err, httpResp),
			)
			return
		}
		if len(providers.GetResults()) == 0 {
			resp.Diagnostics.AddError(
				"Provider Not Found",
				fmt.Sprintf("No circuit provider found with name: %s", providerName),
			)
			return
		}
		provider = &providers.GetResults()[0]
	default:
		resp.Diagnostics.AddError(
			"Missing Required Attribute",
			"At least one of 'id', 'name', or 'slug' must be specified to look up a circuit provider.",
		)
		return
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading circuit provider",
			utils.FormatAPIError("read circuit provider", err, httpResp),
		)
		return
	}

	// Map the provider to state
	data.ID = types.StringValue(fmt.Sprintf("%d", provider.GetId()))
	data.Name = types.StringValue(provider.GetName())
	data.Slug = types.StringValue(provider.GetSlug())

	// Handle description
	if provider.HasDescription() && provider.GetDescription() != "" {
		data.Description = types.StringValue(provider.GetDescription())
	} else {
		data.Description = types.StringNull()
	}

	// Handle comments
	if provider.HasComments() && provider.GetComments() != "" {
		data.Comments = types.StringValue(provider.GetComments())
	} else {
		data.Comments = types.StringNull()
	}

	// Handle tags
	if provider.HasTags() {
		tags := utils.NestedTagsToTagModels(provider.GetTags())
		tagsValue, tagDiags := types.SetValueFrom(ctx, utils.GetTagsAttributeType().ElemType, tags)
		resp.Diagnostics.Append(tagDiags...)
		if resp.Diagnostics.HasError() {
			return
		}
		data.Tags = tagsValue
	} else {
		data.Tags = types.SetNull(utils.GetTagsAttributeType().ElemType)
	}

	// Handle custom fields
	if provider.HasCustomFields() {
		customFields := utils.MapToCustomFieldModels(provider.GetCustomFields(), nil)
		customFieldsValue, cfDiags := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, customFields)
		resp.Diagnostics.Append(cfDiags...)
		if resp.Diagnostics.HasError() {
			return
		}
		data.CustomFields = customFieldsValue
	} else {
		data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)
	}

	tflog.Debug(ctx, "Read circuit provider", map[string]interface{}{
		"id":   provider.GetId(),
		"name": provider.GetName(),
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
