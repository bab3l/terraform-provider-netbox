package datasources

import (
	"context"

	"fmt"

	"github.com/bab3l/go-netbox"
	"github.com/bab3l/terraform-provider-netbox/internal/netboxlookup"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &PlatformDataSource{}

func NewPlatformDataSource() datasource.DataSource {
	return &PlatformDataSource{}
}

type PlatformDataSource struct {
	client *netbox.APIClient
}

type PlatformDataSourceModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Slug         types.String `tfsdk:"slug"`
	Manufacturer types.String `tfsdk:"manufacturer"`
	Description  types.String `tfsdk:"description"`
}

func (d *PlatformDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_platform"
}

func (d *PlatformDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get information about a platform type in Netbox. Platforms represent operating systems or firmware types for devices.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier for the platform. Specify `id`, `slug`, or `name` to identify the platform.",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Full name of the platform. Can be used to identify the platform instead of `id` or `slug`.",
				Optional:            true,
				Computed:            true,
			},
			"slug": schema.StringAttribute{
				MarkdownDescription: "URL-friendly identifier for the platform. Specify `id`, `slug`, or `name` to identify the platform.",
				Optional:            true,
				Computed:            true,
			},
			"manufacturer": schema.StringAttribute{
				MarkdownDescription: "Name or ID of the manufacturer for this platform.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Detailed description of the platform.",
				Computed:            true,
			},
		},
	}
}

// Implement Read method here
func (d *PlatformDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data PlatformDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var platform *netbox.Platform
	var err error
	var httpResp *netbox.APIResponse
	// Lookup by id, slug, or name
	if !data.ID.IsNull() {
		platformID := data.ID.ValueString()
		var platformIDInt int32
		if _, parseErr := fmt.Sscanf(platformID, "%d", &platformIDInt); parseErr != nil {
			resp.Diagnostics.AddError("Invalid Platform ID", "Platform ID must be a number.")
			return
		}
	} else if !data.Slug.IsNull() {
		slug := data.Slug.ValueString()
		platforms, httpResp, err := d.client.DcimAPI.DcimPlatformsList(ctx).Slug([]string{slug}).Execute()
		if err == nil && httpResp.StatusCode == 200 && len(platforms.GetResults()) > 0 {
			platform = &platforms.GetResults()[0]
		}
	} else if !data.Name.IsNull() {
		name := data.Name.ValueString()
		platforms, httpResp, err := d.client.DcimAPI.DcimPlatformsList(ctx).Name([]string{name}).Execute()
		if err == nil && httpResp.StatusCode == 200 && len(platforms.GetResults()) > 0 {
			platform = &platforms.GetResults()[0]
		}
	} else {
		resp.Diagnostics.AddError("Missing Platform Identifier", "Either 'id', 'slug', or 'name' must be specified.")
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Error reading platform", err.Error())
		return
	}
	if httpResp == nil || httpResp.StatusCode != 200 || platform == nil {
		resp.Diagnostics.AddError("Platform Not Found", "No platform found with the specified identifier.")
		return
	}
	data.ID = types.StringValue(fmt.Sprintf("%d", platform.GetId()))
	data.Name = types.StringValue(platform.GetName())
	data.Slug = types.StringValue(platform.GetSlug())
	if platform.HasManufacturer() {
		manufacturerID := platform.GetManufacturer()
		// Use helper to resolve manufacturer to name/slug
		manufacturerRef, diags := netboxlookup.LookupManufacturerBrief(ctx, d.client, fmt.Sprintf("%d", manufacturerID))
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		data.Manufacturer = types.StringValue(manufacturerRef.Name)
	} else {
		data.Manufacturer = types.StringNull()
	}
	if platform.HasDescription() {
		data.Description = types.StringValue(platform.GetDescription())
	} else {
		data.Description = types.StringNull()
	}
	// Comments field may not exist; set to null
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
