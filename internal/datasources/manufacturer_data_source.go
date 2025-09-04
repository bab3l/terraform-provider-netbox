package datasources

import (
	"context"

	"github.com/bab3l/go-netbox"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &ManufacturerDataSource{}

func NewManufacturerDataSource() datasource.DataSource {
	return &ManufacturerDataSource{}
}

type ManufacturerDataSource struct {
	client *netbox.APIClient
}

type ManufacturerDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Slug        types.String `tfsdk:"slug"`
	Description types.String `tfsdk:"description"`
}

func (d *ManufacturerDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_manufacturer"
}

func (d *ManufacturerDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get information about a manufacturer in Netbox. Manufacturers are used to group devices and platforms by vendor.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier for the manufacturer. Specify `id`, `slug`, or `name` to identify the manufacturer.",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Full name of the manufacturer. Can be used to identify the manufacturer instead of `id` or `slug`.",
				Optional:            true,
				Computed:            true,
			},
			"slug": schema.StringAttribute{
				MarkdownDescription: "URL-friendly identifier for the manufacturer. Specify `id`, `slug`, or `name` to identify the manufacturer.",
				Optional:            true,
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Detailed description of the manufacturer.",
				Computed:            true,
			},
		},
	}
}

// Implement Read method here
func (d *ManufacturerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ManufacturerDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// TODO: Implement lookup by id, slug, or name using go-netbox client
	// Update data with API response
}
