// Package datasources contains Terraform data source implementations for the Netbox provider.
package datasources

import (
	"context"
	"fmt"

	"github.com/bab3l/go-netbox"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ datasource.DataSource              = &RIRDataSource{}
	_ datasource.DataSourceWithConfigure = &RIRDataSource{}
)

// NewRIRDataSource returns a new RIR data source.
func NewRIRDataSource() datasource.DataSource {
	return &RIRDataSource{}
}

// RIRDataSource defines the data source implementation.
type RIRDataSource struct {
	client *netbox.APIClient
}

// RIRDataSourceModel describes the data source data model.
type RIRDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Slug        types.String `tfsdk:"slug"`
	IsPrivate   types.Bool   `tfsdk:"is_private"`
	Description types.String `tfsdk:"description"`
	Tags        types.List   `tfsdk:"tags"`
}

// Metadata returns the data source type name.
func (d *RIRDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_rir"
}

// Schema defines the schema for the data source.
func (d *RIRDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get information about a Regional Internet Registry (RIR) in Netbox.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the RIR. Either `id`, `name`, or `slug` must be specified.",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the RIR.",
				Optional:            true,
				Computed:            true,
			},
			"slug": schema.StringAttribute{
				MarkdownDescription: "The slug of the RIR.",
				Optional:            true,
				Computed:            true,
			},
			"is_private": schema.BoolAttribute{
				MarkdownDescription: "Whether IP space managed by this RIR is considered private.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the RIR.",
				Computed:            true,
			},
			"tags": schema.ListAttribute{
				MarkdownDescription: "The tags assigned to this RIR.",
				Computed:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *RIRDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *RIRDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data RIRDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var rir *netbox.RIR

	// Check if we're looking up by ID, name, or slug
	switch {
	case utils.IsSet(data.ID):
		var idInt int
		_, err := fmt.Sscanf(data.ID.ValueString(), "%d", &idInt)
		if err != nil {
			resp.Diagnostics.AddError(
				"Invalid ID",
				fmt.Sprintf("Unable to parse ID %q: %s", data.ID.ValueString(), err.Error()),
			)
			return
		}

		tflog.Debug(ctx, "Reading RIR by ID", map[string]interface{}{
			"id": idInt,
		})

		id32, err := utils.SafeInt32(int64(idInt))
		if err != nil {
			resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("ID value overflow: %s", err))
			return
		}

		result, httpResp, err := d.client.IpamAPI.IpamRirsRetrieve(ctx, id32).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading RIR",
				utils.FormatAPIError(fmt.Sprintf("retrieve RIR ID %d", idInt), err, httpResp),
			)
			return
		}
		rir = result
	case utils.IsSet(data.Name):
		// Looking up by name
		tflog.Debug(ctx, "Reading RIR by name", map[string]interface{}{
			"name": data.Name.ValueString(),
		})

		listReq := d.client.IpamAPI.IpamRirsList(ctx)
		listReq = listReq.Name([]string{data.Name.ValueString()})

		results, httpResp, err := listReq.Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error listing RIRs",
				utils.FormatAPIError(fmt.Sprintf("list RIRs with name %q", data.Name.ValueString()), err, httpResp),
			)
			return
		}

		if results.Count == 0 {
			resp.Diagnostics.AddError(
				"RIR not found",
				fmt.Sprintf("No RIR found with name %q", data.Name.ValueString()),
			)
			return
		}

		if results.Count > 1 {
			resp.Diagnostics.AddError(
				"Multiple RIRs found",
				fmt.Sprintf("Found %d RIRs with name %q. Please use 'id' or 'slug' to specify the exact RIR.", results.Count, data.Name.ValueString()),
			)
			return
		}

		rir = &results.Results[0]
	case utils.IsSet(data.Slug):
		// Looking up by slug
		tflog.Debug(ctx, "Reading RIR by slug", map[string]interface{}{
			"slug": data.Slug.ValueString(),
		})

		listReq := d.client.IpamAPI.IpamRirsList(ctx)
		listReq = listReq.Slug([]string{data.Slug.ValueString()})

		results, httpResp, err := listReq.Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error listing RIRs",
				utils.FormatAPIError(fmt.Sprintf("list RIRs with slug %q", data.Slug.ValueString()), err, httpResp),
			)
			return
		}

		if results.Count == 0 {
			resp.Diagnostics.AddError(
				"RIR not found",
				fmt.Sprintf("No RIR found with slug %q", data.Slug.ValueString()),
			)
			return
		}

		rir = &results.Results[0]
	default:
		resp.Diagnostics.AddError(
			"Missing search criteria",
			"Either 'id', 'name', or 'slug' must be specified to look up a RIR.",
		)
		return
	}

	// Map response to model
	d.mapRIRToDataSourceModel(ctx, rir, &data)

	tflog.Debug(ctx, "Read RIR", map[string]interface{}{
		"id":   data.ID.ValueString(),
		"name": data.Name.ValueString(),
	})

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// mapRIRToDataSourceModel maps a Netbox RIR to the Terraform data source model.
func (d *RIRDataSource) mapRIRToDataSourceModel(ctx context.Context, rir *netbox.RIR, data *RIRDataSourceModel) {
	data.ID = types.StringValue(fmt.Sprintf("%d", rir.Id))
	data.Name = types.StringValue(rir.Name)
	data.Slug = types.StringValue(rir.Slug)

	// Is Private
	if rir.IsPrivate != nil {
		data.IsPrivate = types.BoolValue(*rir.IsPrivate)
	} else {
		data.IsPrivate = types.BoolValue(false)
	}

	// Description
	if rir.Description != nil && *rir.Description != "" {
		data.Description = types.StringValue(*rir.Description)
	} else {
		data.Description = types.StringNull()
	}

	// Tags - convert to list of strings (tag names)
	if len(rir.Tags) > 0 {
		tagNames := make([]string, len(rir.Tags))
		for i, tag := range rir.Tags {
			tagNames[i] = tag.Name
		}
		tagsList, _ := types.ListValueFrom(ctx, types.StringType, tagNames)
		data.Tags = tagsList
	} else {
		data.Tags = types.ListNull(types.StringType)
	}
}
