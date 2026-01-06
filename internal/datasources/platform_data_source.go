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
)

var _ datasource.DataSource = &PlatformDataSource{}

func NewPlatformDataSource() datasource.DataSource {
	return &PlatformDataSource{}
}

type PlatformDataSource struct {
	client *netbox.APIClient
}

type PlatformDataSourceModel struct {
	ID types.String `tfsdk:"id"`

	DisplayName types.String `tfsdk:"display_name"`

	Name types.String `tfsdk:"name"`

	Slug types.String `tfsdk:"slug"`

	Manufacturer types.String `tfsdk:"manufacturer"`

	Description types.String `tfsdk:"description"`
}

func (d *PlatformDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_platform"
}

func (d *PlatformDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get information about a platform type in Netbox. Platforms represent operating systems or firmware types for devices.",

		Attributes: map[string]schema.Attribute{
			"id": nbschema.DSIDAttribute("platform"),

			"display_name": nbschema.DSComputedStringAttribute("The display name of the platform."),

			"name": nbschema.DSNameAttribute("platform"),

			"slug": nbschema.DSSlugAttribute("platform"),

			"manufacturer": nbschema.DSComputedStringAttribute("Name or ID of the manufacturer for this platform."),

			"description": nbschema.DSComputedStringAttribute("Detailed description of the platform."),
		},
	}
}

func (d *PlatformDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *PlatformDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data PlatformDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var platform *netbox.Platform

	var err error

	var httpResp *http.Response

	// Lookup by id, slug, or name

	switch {
	case !data.ID.IsNull():

		platformID := data.ID.ValueString()

		var platformIDInt int32

		if _, parseErr := fmt.Sscanf(platformID, "%d", &platformIDInt); parseErr != nil {
			resp.Diagnostics.AddError("Invalid Platform ID", "Platform ID must be a number.")

			return
		}

		var p *netbox.Platform

		p, httpResp, err = d.client.DcimAPI.DcimPlatformsRetrieve(ctx, platformIDInt).Execute()

		defer utils.CloseResponseBody(httpResp)

		if err == nil && httpResp.StatusCode == http.StatusOK {
			platform = p
		}

	case !data.Slug.IsNull():

		slug := data.Slug.ValueString()

		var platforms *netbox.PaginatedPlatformList

		platforms, httpResp, err = d.client.DcimAPI.DcimPlatformsList(ctx).Slug([]string{slug}).Execute()

		defer utils.CloseResponseBody(httpResp)

		if err == nil && httpResp.StatusCode == http.StatusOK && len(platforms.GetResults()) > 0 {
			platform = &platforms.GetResults()[0]
		}

	case !data.Name.IsNull():

		name := data.Name.ValueString()

		var platforms *netbox.PaginatedPlatformList

		platforms, httpResp, err = d.client.DcimAPI.DcimPlatformsList(ctx).Name([]string{name}).Execute()

		defer utils.CloseResponseBody(httpResp)

		if err == nil && httpResp.StatusCode == http.StatusOK && len(platforms.GetResults()) > 0 {
			platform = &platforms.GetResults()[0]
		}

	default:

		resp.Diagnostics.AddError("Missing Platform Identifier", "Either 'id', 'slug', or 'name' must be specified.")

		return
	}

	if err != nil {
		resp.Diagnostics.AddError("Error reading platform", utils.FormatAPIError("read platform", err, httpResp))

		return
	}

	if httpResp == nil || httpResp.StatusCode != http.StatusOK || platform == nil {
		resp.Diagnostics.AddError("Platform Not Found", "No platform found with the specified identifier.")

		return
	}

	data.ID = types.StringValue(fmt.Sprintf("%d", platform.GetId()))

	// Display Name

	if platform.GetDisplay() != "" {
		data.DisplayName = types.StringValue(platform.GetDisplay())
	} else {
		data.DisplayName = types.StringNull()
	}

	data.Name = types.StringValue(platform.GetName())

	data.Slug = types.StringValue(platform.GetSlug())

	if platform.HasManufacturer() {
		manufacturer := platform.GetManufacturer()

		data.Manufacturer = types.StringValue(manufacturer.GetName())
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
