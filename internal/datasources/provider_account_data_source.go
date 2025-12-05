// Package datasources contains Terraform data source implementations for the Netbox provider.
package datasources

import (
	"context"
	"fmt"

	"github.com/bab3l/go-netbox"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/bab3l/terraform-provider-netbox/internal/utils"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ datasource.DataSource              = &ProviderAccountDataSource{}
	_ datasource.DataSourceWithConfigure = &ProviderAccountDataSource{}
)

// NewProviderAccountDataSource returns a new Provider Account data source.
func NewProviderAccountDataSource() datasource.DataSource {
	return &ProviderAccountDataSource{}
}

// ProviderAccountDataSource defines the data source implementation.
type ProviderAccountDataSource struct {
	client *netbox.APIClient
}

// ProviderAccountDataSourceModel describes the data source data model.
type ProviderAccountDataSourceModel struct {
	ID           types.String `tfsdk:"id"`
	Provider     types.String `tfsdk:"provider"`
	ProviderName types.String `tfsdk:"provider_name"`
	Name         types.String `tfsdk:"name"`
	Account      types.String `tfsdk:"account"`
	Description  types.String `tfsdk:"description"`
	Comments     types.String `tfsdk:"comments"`
	Tags         types.List   `tfsdk:"tags"`
}

// Metadata returns the data source type name.
func (d *ProviderAccountDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_provider_account"
}

// Schema defines the schema for the data source.
func (d *ProviderAccountDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to retrieve information about a provider account in Netbox.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique numeric ID of the provider account. Use this to look up a provider account by ID.",
				Optional:            true,
				Computed:            true,
			},
			"provider": schema.StringAttribute{
				MarkdownDescription: "The ID of the circuit provider this account belongs to.",
				Computed:            true,
			},
			"provider_name": schema.StringAttribute{
				MarkdownDescription: "The name of the circuit provider this account belongs to.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the provider account. Can be used with provider to look up an account.",
				Optional:            true,
				Computed:            true,
			},
			"account": schema.StringAttribute{
				MarkdownDescription: "The account identifier. Can be used with provider to look up an account.",
				Optional:            true,
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "A description of the provider account.",
				Computed:            true,
			},
			"comments": schema.StringAttribute{
				MarkdownDescription: "Additional comments about the provider account.",
				Computed:            true,
			},
			"tags": schema.ListAttribute{
				MarkdownDescription: "Tags assigned to this provider account.",
				Computed:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

// Configure sets the client for the data source.
func (d *ProviderAccountDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
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

// Read reads the provider account data source.
func (d *ProviderAccountDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ProviderAccountDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var providerAccount *netbox.ProviderAccount

	// Look up by ID if provided
	if !data.ID.IsNull() && !data.ID.IsUnknown() {
		id, err := utils.ParseID(data.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Invalid ID",
				fmt.Sprintf("Could not convert ID to integer: %s", err.Error()),
			)
			return
		}

		tflog.Debug(ctx, "Looking up provider account by ID", map[string]interface{}{
			"id": id,
		})

		result, httpResp, err := d.client.CircuitsAPI.CircuitsProviderAccountsRetrieve(ctx, id).Execute()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading provider account",
				fmt.Sprintf("Could not read provider account with ID %d: %s\nHTTP Response: %v", id, err.Error(), httpResp),
			)
			return
		}
		providerAccount = result
	} else if !data.Account.IsNull() && !data.Account.IsUnknown() {
		// Look up by account identifier
		account := data.Account.ValueString()

		tflog.Debug(ctx, "Looking up provider account by account identifier", map[string]interface{}{
			"account": account,
		})

		list, httpResp, err := d.client.CircuitsAPI.CircuitsProviderAccountsList(ctx).Account([]string{account}).Execute()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading provider account",
				fmt.Sprintf("Could not find provider account with account %s: %s\nHTTP Response: %v", account, err.Error(), httpResp),
			)
			return
		}

		if len(list.Results) == 0 {
			resp.Diagnostics.AddError(
				"Provider account not found",
				fmt.Sprintf("No provider account found with account identifier %s", account),
			)
			return
		}

		if len(list.Results) > 1 {
			resp.Diagnostics.AddError(
				"Multiple provider accounts found",
				fmt.Sprintf("Found %d provider accounts with account identifier %s, expected exactly one. Consider filtering by provider as well.", len(list.Results), account),
			)
			return
		}

		providerAccount = &list.Results[0]
	} else {
		resp.Diagnostics.AddError(
			"Missing required attribute",
			"Either 'id' or 'account' must be specified to look up a provider account.",
		)
		return
	}

	// Map response to model
	d.mapResponseToModel(ctx, providerAccount, &data)

	tflog.Debug(ctx, "Read provider account", map[string]interface{}{
		"id":      data.ID.ValueString(),
		"account": data.Account.ValueString(),
	})

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// mapResponseToModel maps the API response to the Terraform model.
func (d *ProviderAccountDataSource) mapResponseToModel(ctx context.Context, providerAccount *netbox.ProviderAccount, data *ProviderAccountDataSourceModel) {
	data.ID = types.StringValue(fmt.Sprintf("%d", providerAccount.GetId()))
	data.Account = types.StringValue(providerAccount.GetAccount())

	// Map Provider
	if provider := providerAccount.GetProvider(); provider.Id != 0 {
		data.Provider = types.StringValue(fmt.Sprintf("%d", provider.Id))
		data.ProviderName = types.StringValue(provider.GetName())
	}

	// Map name
	if name, ok := providerAccount.GetNameOk(); ok && name != nil && *name != "" {
		data.Name = types.StringValue(*name)
	} else {
		data.Name = types.StringNull()
	}

	// Map description
	if description, ok := providerAccount.GetDescriptionOk(); ok && description != nil {
		data.Description = types.StringValue(*description)
	} else {
		data.Description = types.StringNull()
	}

	// Map comments
	if comments, ok := providerAccount.GetCommentsOk(); ok && comments != nil {
		data.Comments = types.StringValue(*comments)
	} else {
		data.Comments = types.StringNull()
	}

	// Map tags
	if tags := providerAccount.GetTags(); len(tags) > 0 {
		tagNames := make([]string, len(tags))
		for i, tag := range tags {
			tagNames[i] = tag.Name
		}
		data.Tags, _ = types.ListValueFrom(ctx, types.StringType, tagNames)
	} else {
		data.Tags = types.ListNull(types.StringType)
	}
}
