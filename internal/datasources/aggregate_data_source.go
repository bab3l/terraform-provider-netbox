// Package datasources contains Terraform data source implementations for the Netbox provider.

package datasources

import (
	"context"
	"fmt"

	"github.com/bab3l/go-netbox"
	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ datasource.DataSource              = &AggregateDataSource{}
	_ datasource.DataSourceWithConfigure = &AggregateDataSource{}
)

// NewAggregateDataSource returns a new Aggregate data source.
func NewAggregateDataSource() datasource.DataSource {
	return &AggregateDataSource{}
}

// AggregateDataSource defines the data source implementation.
type AggregateDataSource struct {
	client *netbox.APIClient
}

// AggregateDataSourceModel describes the data source data model.
type AggregateDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Prefix      types.String `tfsdk:"prefix"`
	RIR         types.String `tfsdk:"rir"`
	RIRName     types.String `tfsdk:"rir_name"`
	Tenant      types.String `tfsdk:"tenant"`
	TenantName  types.String `tfsdk:"tenant_name"`
	DateAdded   types.String `tfsdk:"date_added"`
	Description types.String `tfsdk:"description"`
	Comments    types.String `tfsdk:"comments"`
	DisplayName types.String `tfsdk:"display_name"`
	Tags        types.List   `tfsdk:"tags"`
}

// Metadata returns the data source type name.
func (d *AggregateDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_aggregate"
}

// Schema defines the schema for the data source.
func (d *AggregateDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to retrieve information about an aggregate in Netbox. You can identify the aggregate using `id` or `prefix`.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique numeric ID of the aggregate. Use this to look up an aggregate by ID.",
				Optional:            true,
				Computed:            true,
			},
			"prefix": schema.StringAttribute{
				MarkdownDescription: "The IP prefix in CIDR notation. Use this to look up an aggregate by prefix.",
				Optional:            true,
				Computed:            true,
			},
			"rir": schema.StringAttribute{
				MarkdownDescription: "The ID of the Regional Internet Registry (RIR) this aggregate belongs to.",
				Computed:            true,
			},
			"rir_name": schema.StringAttribute{
				MarkdownDescription: "The name of the Regional Internet Registry (RIR) this aggregate belongs to.",
				Computed:            true,
			},
			"tenant": schema.StringAttribute{
				MarkdownDescription: "The ID of the tenant this aggregate is assigned to.",
				Computed:            true,
			},
			"tenant_name": schema.StringAttribute{
				MarkdownDescription: "The name of the tenant this aggregate is assigned to.",
				Computed:            true,
			},
			"date_added": schema.StringAttribute{
				MarkdownDescription: "The date this aggregate was added (YYYY-MM-DD format).",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "A description of the aggregate.",
				Computed:            true,
			},
			"comments": schema.StringAttribute{
				MarkdownDescription: "Additional comments about the aggregate.",
				Computed:            true,
			},
			"display_name": nbschema.DSComputedStringAttribute("The display name of the aggregate."),
			"tags": schema.ListAttribute{
				MarkdownDescription: "Tags assigned to this aggregate.",
				Computed:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

// Configure sets the client for the data source.
func (d *AggregateDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// Read reads the aggregate data source.
func (d *AggregateDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data AggregateDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var aggregate *netbox.Aggregate

	// Look up by ID if provided
	switch {
	case !data.ID.IsNull() && !data.ID.IsUnknown():
		id, err := utils.ParseID(data.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Invalid ID",
				fmt.Sprintf("Could not convert ID to integer: %s", err.Error()),
			)
			return
		}
		tflog.Debug(ctx, "Looking up aggregate by ID", map[string]interface{}{
			"id": id,
		})
		result, httpResp, err := d.client.IpamAPI.IpamAggregatesRetrieve(ctx, id).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading aggregate",
				fmt.Sprintf("Could not read aggregate with ID %d: %s\nHTTP Response: %v", id, err.Error(), httpResp),
			)
			return
		}
		aggregate = result

	case !data.Prefix.IsNull() && !data.Prefix.IsUnknown():
		// Look up by prefix
		prefix := data.Prefix.ValueString()
		tflog.Debug(ctx, "Looking up aggregate by prefix", map[string]interface{}{
			"prefix": prefix,
		})
		list, httpResp, err := d.client.IpamAPI.IpamAggregatesList(ctx).Prefix(prefix).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading aggregate",
				fmt.Sprintf("Could not find aggregate with prefix %s: %s\nHTTP Response: %v", prefix, err.Error(), httpResp),
			)
			return
		}
		if len(list.Results) == 0 {
			resp.Diagnostics.AddError(
				"Aggregate not found",
				fmt.Sprintf("No aggregate found with prefix %s", prefix),
			)
			return
		}
		if len(list.Results) > 1 {
			resp.Diagnostics.AddError(
				"Multiple aggregates found",
				fmt.Sprintf("Found %d aggregates with prefix %s, expected exactly one", len(list.Results), prefix),
			)
			return
		}
		aggregate = &list.Results[0]

	default:
		resp.Diagnostics.AddError(
			"Missing required attribute",
			"Either 'id' or 'prefix' must be specified to look up an aggregate.",
		)
		return
	}

	// Map response to model
	d.mapResponseToModel(ctx, aggregate, &data)
	tflog.Debug(ctx, "Read aggregate", map[string]interface{}{
		"id":     data.ID.ValueString(),
		"prefix": data.Prefix.ValueString(),
	})

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// mapResponseToModel maps the API response to the Terraform model.
func (d *AggregateDataSource) mapResponseToModel(ctx context.Context, aggregate *netbox.Aggregate, data *AggregateDataSourceModel) {
	data.ID = types.StringValue(fmt.Sprintf("%d", aggregate.GetId()))
	data.Prefix = types.StringValue(aggregate.GetPrefix())

	// Map RIR
	if rir := aggregate.GetRir(); rir.Id != 0 {
		data.RIR = types.StringValue(fmt.Sprintf("%d", rir.Id))
		data.RIRName = types.StringValue(rir.GetName())
	}

	// Map tenant
	if tenant, ok := aggregate.GetTenantOk(); ok && tenant != nil && tenant.Id != 0 {
		data.Tenant = types.StringValue(fmt.Sprintf("%d", tenant.GetId()))
		data.TenantName = types.StringValue(tenant.GetName())
	} else {
		data.Tenant = types.StringNull()
		data.TenantName = types.StringNull()
	}

	// Map date_added
	if dateAdded := aggregate.GetDateAdded(); dateAdded != "" {
		data.DateAdded = types.StringValue(dateAdded)
	} else {
		data.DateAdded = types.StringNull()
	}

	// Map description
	if description, ok := aggregate.GetDescriptionOk(); ok && description != nil {
		data.Description = types.StringValue(*description)
	} else {
		data.Description = types.StringNull()
	}

	// Map comments
	if comments, ok := aggregate.GetCommentsOk(); ok && comments != nil {
		data.Comments = types.StringValue(*comments)
	} else {
		data.Comments = types.StringNull()
	}

	// Map tags
	if tags := aggregate.GetTags(); len(tags) > 0 {
		tagNames := make([]string, len(tags))
		for i, tag := range tags {
			tagNames[i] = tag.Name
		}
		data.Tags, _ = types.ListValueFrom(ctx, types.StringType, tagNames)
	} else {
		data.Tags = types.ListNull(types.StringType)
	}

	// Map display name
	if aggregate.GetDisplay() != "" {
		data.DisplayName = types.StringValue(aggregate.GetDisplay())
	} else {
		data.DisplayName = types.StringNull()
	}
}
