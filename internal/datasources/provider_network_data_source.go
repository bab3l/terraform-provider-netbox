// Package datasources contains Terraform data source implementations for NetBox objects.
package datasources

import (
	"context"
	"fmt"

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
	_ datasource.DataSource              = &ProviderNetworkDataSource{}
	_ datasource.DataSourceWithConfigure = &ProviderNetworkDataSource{}
)

// NewProviderNetworkDataSource returns a new data source implementing the ProviderNetwork data source.
func NewProviderNetworkDataSource() datasource.DataSource {
	return &ProviderNetworkDataSource{}
}

// ProviderNetworkDataSource defines the data source implementation.
type ProviderNetworkDataSource struct {
	client *netbox.APIClient
}

// ProviderNetworkDataSourceModel describes the data source data model.
type ProviderNetworkDataSourceModel struct {
	ID              types.String `tfsdk:"id"`
	CircuitProvider types.String `tfsdk:"circuit_provider"`
	Name            types.String `tfsdk:"name"`
	ServiceID       types.String `tfsdk:"service_id"`
	Description     types.String `tfsdk:"description"`
	Comments        types.String `tfsdk:"comments"`
	Tags            types.Set    `tfsdk:"tags"`
	CustomFields    types.Set    `tfsdk:"custom_fields"`
}

// Metadata returns the data source type name.
func (d *ProviderNetworkDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_provider_network"
}

// Schema defines the schema for the data source.
func (d *ProviderNetworkDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves information about a provider network in NetBox.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique numeric ID of the provider network. Use this to look up by ID.",
				Optional:            true,
				Computed:            true,
			},
			"circuit_provider": schema.StringAttribute{
				MarkdownDescription: "The circuit provider that owns this network. Can be used with name to filter.",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the provider network. Use this to look up by name.",
				Optional:            true,
				Computed:            true,
			},
			"service_id": schema.StringAttribute{
				MarkdownDescription: "A unique identifier for this network provided by the circuit provider.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "A description of the provider network.",
				Computed:            true,
			},
			"comments": schema.StringAttribute{
				MarkdownDescription: "Additional comments or notes about this provider network.",
				Computed:            true,
			},
			"tags":          nbschema.DSTagsAttribute(),
			"custom_fields": nbschema.DSCustomFieldsAttribute(),
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *ProviderNetworkDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *ProviderNetworkDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ProviderNetworkDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var pn *netbox.ProviderNetwork

	// Look up by ID if provided
	switch {
	case !data.ID.IsNull() && !data.ID.IsUnknown():
		pnID, err := utils.ParseID(data.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Invalid Provider Network ID",
				fmt.Sprintf("Provider network ID must be a number, got: %s", data.ID.ValueString()),
			)
			return
		}

		tflog.Debug(ctx, "Reading provider network by ID", map[string]interface{}{
			"id": pnID,
		})

		p, httpResp, err := d.client.CircuitsAPI.CircuitsProviderNetworksRetrieve(ctx, pnID).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading provider network",
				utils.FormatAPIError(fmt.Sprintf("read provider network ID %d", pnID), err, httpResp),
			)
			return
		}
		pn = p
	case !data.Name.IsNull() && !data.Name.IsUnknown():
		// Look up by name
		tflog.Debug(ctx, "Reading provider network by name", map[string]interface{}{
			"name": data.Name.ValueString(),
		})

		listReq := d.client.CircuitsAPI.CircuitsProviderNetworksList(ctx).Name([]string{data.Name.ValueString()})

		// Optionally filter by provider as well
		if !data.CircuitProvider.IsNull() && !data.CircuitProvider.IsUnknown() {
			listReq = listReq.Provider([]string{data.CircuitProvider.ValueString()})
		}

		listResp, httpResp, err := listReq.Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading provider network",
				utils.FormatAPIError(fmt.Sprintf("read provider network by name %s", data.Name.ValueString()), err, httpResp),
			)
			return
		}

		if listResp.GetCount() == 0 {
			resp.Diagnostics.AddError(
				"Provider network not found",
				fmt.Sprintf("No provider network found with name: %s", data.Name.ValueString()),
			)
			return
		}

		if listResp.GetCount() > 1 {
			resp.Diagnostics.AddError(
				"Multiple provider networks found",
				fmt.Sprintf("Found %d provider networks with name: %s. Consider filtering by provider as well.", listResp.GetCount(), data.Name.ValueString()),
			)
			return
		}

		pn = &listResp.GetResults()[0]
	default:
		resp.Diagnostics.AddError(
			"Missing Required Attribute",
			"Either 'id' or 'name' must be specified to look up a provider network.",
		)
		return
	}

	// Map response to model
	d.mapResponseToModel(ctx, pn, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// mapResponseToModel maps the API response to the Terraform model.
func (d *ProviderNetworkDataSource) mapResponseToModel(ctx context.Context, pn *netbox.ProviderNetwork, data *ProviderNetworkDataSourceModel, diags *diag.Diagnostics) {
	data.ID = types.StringValue(fmt.Sprintf("%d", pn.GetId()))
	data.Name = types.StringValue(pn.GetName())

	// Map Provider
	data.CircuitProvider = types.StringValue(pn.Provider.GetName())

	// Map service_id
	if serviceID, ok := pn.GetServiceIdOk(); ok && serviceID != nil && *serviceID != "" {
		data.ServiceID = types.StringValue(*serviceID)
	} else {
		data.ServiceID = types.StringNull()
	}

	// Map description
	if desc, ok := pn.GetDescriptionOk(); ok && desc != nil && *desc != "" {
		data.Description = types.StringValue(*desc)
	} else {
		data.Description = types.StringNull()
	}

	// Map comments
	if comments, ok := pn.GetCommentsOk(); ok && comments != nil && *comments != "" {
		data.Comments = types.StringValue(*comments)
	} else {
		data.Comments = types.StringNull()
	}

	// Handle tags
	if pn.HasTags() && len(pn.GetTags()) > 0 {
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
		customFields := utils.MapToCustomFieldModels(apiCustomFields, nil)
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
