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
			"id":          nbschema.DSIDAttribute("manufacturer"),
			"name":        nbschema.DSNameAttribute("manufacturer"),
			"slug":        nbschema.DSSlugAttribute("manufacturer"),
			"description": nbschema.DSComputedStringAttribute("Description of the manufacturer."),
		},
	}
}

func (d *ManufacturerDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ManufacturerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ManufacturerDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var manufacturer *netbox.Manufacturer
	var err error
	var httpResp *http.Response

	// Lookup by id, slug, or name
	if !data.ID.IsNull() {
		manufacturerID := data.ID.ValueString()
		var manufacturerIDInt int32
		if _, parseErr := fmt.Sscanf(manufacturerID, "%d", &manufacturerIDInt); parseErr != nil {
			resp.Diagnostics.AddError("Invalid Manufacturer ID", "Manufacturer ID must be a number.")
			return
		}
		var m *netbox.Manufacturer
		m, httpResp, err = d.client.DcimAPI.DcimManufacturersRetrieve(ctx, manufacturerIDInt).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err == nil && httpResp.StatusCode == 200 {
			manufacturer = m
		}
	} else if !data.Slug.IsNull() {
		slug := data.Slug.ValueString()
		var manufacturers *netbox.PaginatedManufacturerList
		manufacturers, httpResp, err = d.client.DcimAPI.DcimManufacturersList(ctx).Slug([]string{slug}).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err == nil && httpResp.StatusCode == 200 && len(manufacturers.GetResults()) > 0 {
			manufacturer = &manufacturers.GetResults()[0]
		}
	} else if !data.Name.IsNull() {
		name := data.Name.ValueString()
		var manufacturers *netbox.PaginatedManufacturerList
		manufacturers, httpResp, err = d.client.DcimAPI.DcimManufacturersList(ctx).Name([]string{name}).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err == nil && httpResp.StatusCode == 200 && len(manufacturers.GetResults()) > 0 {
			manufacturer = &manufacturers.GetResults()[0]
		}
	} else {
		resp.Diagnostics.AddError("Missing Manufacturer Identifier", "Either 'id', 'slug', or 'name' must be specified.")
		return
	}

	if err != nil {
		resp.Diagnostics.AddError("Error reading manufacturer", utils.FormatAPIError("read manufacturer", err, httpResp))
		return
	}

	if httpResp == nil || httpResp.StatusCode != 200 || manufacturer == nil {
		resp.Diagnostics.AddError("Manufacturer Not Found", "No manufacturer found with the specified identifier.")
		return
	}

	// Map API response to model using helpers
	d.mapManufacturerToState(manufacturer, &data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// mapManufacturerToState maps API response to Terraform state using state helpers.
func (d *ManufacturerDataSource) mapManufacturerToState(manufacturer *netbox.Manufacturer, data *ManufacturerDataSourceModel) {
	data.ID = types.StringValue(fmt.Sprintf("%d", manufacturer.GetId()))
	data.Name = types.StringValue(manufacturer.GetName())
	data.Slug = types.StringValue(manufacturer.GetSlug())
	data.Description = utils.StringFromAPI(manufacturer.HasDescription(), manufacturer.GetDescription, data.Description)
}
