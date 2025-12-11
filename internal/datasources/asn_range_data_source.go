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
	_ datasource.DataSource              = &ASNRangeDataSource{}
	_ datasource.DataSourceWithConfigure = &ASNRangeDataSource{}
)

// NewASNRangeDataSource returns a new ASNRange data source.
func NewASNRangeDataSource() datasource.DataSource {
	return &ASNRangeDataSource{}
}

// ASNRangeDataSource defines the data source implementation.
type ASNRangeDataSource struct {
	client *netbox.APIClient
}

// ASNRangeDataSourceModel describes the data source data model.
type ASNRangeDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Slug        types.String `tfsdk:"slug"`
	RIR         types.String `tfsdk:"rir"`
	RIRName     types.String `tfsdk:"rir_name"`
	Start       types.String `tfsdk:"start"`
	End         types.String `tfsdk:"end"`
	Tenant      types.String `tfsdk:"tenant"`
	TenantName  types.String `tfsdk:"tenant_name"`
	Description types.String `tfsdk:"description"`
	ASNCount    types.Int64  `tfsdk:"asn_count"`
	Tags        types.List   `tfsdk:"tags"`
}

// Metadata returns the data source type name.
func (d *ASNRangeDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_asn_range"
}

// Schema defines the schema for the data source.
func (d *ASNRangeDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get information about an ASN Range in Netbox.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the ASN range. Either `id`, `name`, or `slug` must be specified.",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the ASN range.",
				Optional:            true,
				Computed:            true,
			},
			"slug": schema.StringAttribute{
				MarkdownDescription: "The slug of the ASN range.",
				Optional:            true,
				Computed:            true,
			},
			"rir": schema.StringAttribute{
				MarkdownDescription: "The ID of the RIR responsible for this ASN range.",
				Computed:            true,
			},
			"rir_name": schema.StringAttribute{
				MarkdownDescription: "The name of the RIR responsible for this ASN range.",
				Computed:            true,
			},
			"start": schema.StringAttribute{
				MarkdownDescription: "The starting ASN in this range.",
				Computed:            true,
			},
			"end": schema.StringAttribute{
				MarkdownDescription: "The ending ASN in this range.",
				Computed:            true,
			},
			"tenant": schema.StringAttribute{
				MarkdownDescription: "The ID of the tenant that owns this ASN range.",
				Computed:            true,
			},
			"tenant_name": schema.StringAttribute{
				MarkdownDescription: "The name of the tenant that owns this ASN range.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the ASN range.",
				Computed:            true,
			},
			"asn_count": schema.Int64Attribute{
				MarkdownDescription: "The number of ASNs allocated from this range.",
				Computed:            true,
			},
			"tags": schema.ListAttribute{
				MarkdownDescription: "The tags assigned to this ASN range.",
				Computed:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *ASNRangeDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// Read refreshes the Terraform state with the latest data.
func (d *ASNRangeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ASNRangeDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var asnRange *netbox.ASNRange

	// Check if we're looking up by ID
	if utils.IsSet(data.ID) {
		var idInt int
		_, err := fmt.Sscanf(data.ID.ValueString(), "%d", &idInt)
		if err != nil {
			resp.Diagnostics.AddError(
				"Invalid ID",
				fmt.Sprintf("Unable to parse ID %q: %s", data.ID.ValueString(), err.Error()),
			)
			return
		}

		tflog.Debug(ctx, "Reading ASNRange by ID", map[string]interface{}{
			"id": idInt,
		})

		id32, err := utils.SafeInt32(int64(idInt))
		if err != nil {
			resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("ID value overflow: %s", err))
			return
		}

		result, httpResp, err := d.client.IpamAPI.IpamAsnRangesRetrieve(ctx, id32).Execute()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading ASNRange",
				utils.FormatAPIError(fmt.Sprintf("retrieve ASNRange ID %d", idInt), err, httpResp),
			)
			return
		}
		asnRange = result
	} else if utils.IsSet(data.Name) {
		// Looking up by name
		tflog.Debug(ctx, "Reading ASNRange by name", map[string]interface{}{
			"name": data.Name.ValueString(),
		})

		listReq := d.client.IpamAPI.IpamAsnRangesList(ctx)
		listReq = listReq.Name([]string{data.Name.ValueString()})

		results, httpResp, err := listReq.Execute()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error listing ASNRanges",
				utils.FormatAPIError(fmt.Sprintf("list ASNRanges with name %q", data.Name.ValueString()), err, httpResp),
			)
			return
		}

		if results.Count == 0 {
			resp.Diagnostics.AddError(
				"ASNRange not found",
				fmt.Sprintf("No ASNRange found with name %q", data.Name.ValueString()),
			)
			return
		}

		if results.Count > 1 {
			resp.Diagnostics.AddError(
				"Multiple ASNRanges found",
				fmt.Sprintf("Found %d ASNRanges with name %q. Please use 'id' or 'slug' to specify the exact ASNRange.", results.Count, data.Name.ValueString()),
			)
			return
		}

		asnRange = &results.Results[0]
	} else if utils.IsSet(data.Slug) {
		// Looking up by slug
		tflog.Debug(ctx, "Reading ASNRange by slug", map[string]interface{}{
			"slug": data.Slug.ValueString(),
		})

		listReq := d.client.IpamAPI.IpamAsnRangesList(ctx)
		listReq = listReq.Slug([]string{data.Slug.ValueString()})

		results, httpResp, err := listReq.Execute()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error listing ASNRanges",
				utils.FormatAPIError(fmt.Sprintf("list ASNRanges with slug %q", data.Slug.ValueString()), err, httpResp),
			)
			return
		}

		if results.Count == 0 {
			resp.Diagnostics.AddError(
				"ASNRange not found",
				fmt.Sprintf("No ASNRange found with slug %q", data.Slug.ValueString()),
			)
			return
		}

		asnRange = &results.Results[0]
	} else {
		resp.Diagnostics.AddError(
			"Missing search criteria",
			"Either 'id', 'name', or 'slug' must be specified to look up an ASNRange.",
		)
		return
	}

	// Map response to model
	d.mapASNRangeToDataSourceModel(ctx, asnRange, &data)

	tflog.Debug(ctx, "Read ASNRange", map[string]interface{}{
		"id":   data.ID.ValueString(),
		"name": data.Name.ValueString(),
	})

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// mapASNRangeToDataSourceModel maps a Netbox ASNRange to the Terraform data source model.
func (d *ASNRangeDataSource) mapASNRangeToDataSourceModel(ctx context.Context, asnRange *netbox.ASNRange, data *ASNRangeDataSourceModel) {
	data.ID = types.StringValue(fmt.Sprintf("%d", asnRange.Id))
	data.Name = types.StringValue(asnRange.Name)
	data.Slug = types.StringValue(asnRange.Slug)
	data.RIR = types.StringValue(fmt.Sprintf("%d", asnRange.Rir.GetId()))
	data.RIRName = types.StringValue(asnRange.Rir.GetName())
	data.Start = types.StringValue(fmt.Sprintf("%d", asnRange.Start))
	data.End = types.StringValue(fmt.Sprintf("%d", asnRange.End))

	// Tenant
	if asnRange.HasTenant() && asnRange.Tenant.Get() != nil {
		tenant := asnRange.Tenant.Get()
		data.Tenant = types.StringValue(fmt.Sprintf("%d", tenant.GetId()))
		data.TenantName = types.StringValue(tenant.GetName())
	} else {
		data.Tenant = types.StringNull()
		data.TenantName = types.StringNull()
	}

	// Description
	if asnRange.Description != nil && *asnRange.Description != "" {
		data.Description = types.StringValue(*asnRange.Description)
	} else {
		data.Description = types.StringNull()
	}

	// ASN count - now optional, may be nil
	if asnRange.AsnCount != nil {
		data.ASNCount = types.Int64Value(int64(*asnRange.AsnCount))
	} else {
		data.ASNCount = types.Int64Value(0)
	}

	// Tags - convert to list of strings (tag names)
	if len(asnRange.Tags) > 0 {
		tagNames := make([]string, len(asnRange.Tags))
		for i, tag := range asnRange.Tags {
			tagNames[i] = tag.Name
		}
		tagsList, _ := types.ListValueFrom(ctx, types.StringType, tagNames)
		data.Tags = tagsList
	} else {
		data.Tags = types.ListNull(types.StringType)
	}
}
