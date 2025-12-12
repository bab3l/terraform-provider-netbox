// Package datasources provides Terraform data source implementations for NetBox objects.
package datasources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bab3l/go-netbox"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &WirelessLinkDataSource{}
var _ datasource.DataSourceWithConfigure = &WirelessLinkDataSource{}

// NewWirelessLinkDataSource returns a new data source implementing wireless link lookup.
func NewWirelessLinkDataSource() datasource.DataSource {
	return &WirelessLinkDataSource{}
}

// WirelessLinkDataSource defines the data source implementation.
type WirelessLinkDataSource struct {
	client *netbox.APIClient
}

// WirelessLinkDataSourceModel describes the data source data model.
type WirelessLinkDataSourceModel struct {
	ID           types.String  `tfsdk:"id"`
	InterfaceA   types.String  `tfsdk:"interface_a"`
	InterfaceB   types.String  `tfsdk:"interface_b"`
	SSID         types.String  `tfsdk:"ssid"`
	Status       types.String  `tfsdk:"status"`
	Tenant       types.String  `tfsdk:"tenant"`
	TenantID     types.String  `tfsdk:"tenant_id"`
	AuthType     types.String  `tfsdk:"auth_type"`
	AuthCipher   types.String  `tfsdk:"auth_cipher"`
	Distance     types.Float64 `tfsdk:"distance"`
	DistanceUnit types.String  `tfsdk:"distance_unit"`
	Description  types.String  `tfsdk:"description"`
	Comments     types.String  `tfsdk:"comments"`
	Tags         types.Set     `tfsdk:"tags"`
	CustomFields types.Set     `tfsdk:"custom_fields"`
}

// Metadata returns the data source type name.
func (d *WirelessLinkDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_wireless_link"
}

// Schema defines the schema for the data source.
func (d *WirelessLinkDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get information about a wireless link in NetBox.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "ID of the wireless link to look up. Either `id` must be specified.",
				Optional:            true,
				Computed:            true,
			},
			"interface_a": schema.StringAttribute{
				MarkdownDescription: "ID of the first interface (A-side) of the wireless link.",
				Computed:            true,
			},
			"interface_b": schema.StringAttribute{
				MarkdownDescription: "ID of the second interface (B-side) of the wireless link.",
				Computed:            true,
			},
			"ssid": schema.StringAttribute{
				MarkdownDescription: "The SSID (network name) for the wireless link.",
				Computed:            true,
			},
			"status":      nbschema.DSComputedStringAttribute("Connection status of the wireless link."),
			"tenant":      nbschema.DSComputedStringAttribute("Name of the tenant that owns this wireless link."),
			"tenant_id":   nbschema.DSComputedStringAttribute("ID of the tenant that owns this wireless link."),
			"auth_type":   nbschema.DSComputedStringAttribute("Authentication type."),
			"auth_cipher": nbschema.DSComputedStringAttribute("Authentication cipher."),
			"distance": schema.Float64Attribute{
				MarkdownDescription: "Distance of the wireless link.",
				Computed:            true,
			},
			"distance_unit": nbschema.DSComputedStringAttribute("Unit for distance."),
			"description":   nbschema.DSComputedStringAttribute("Description of the wireless link."),
			"comments":      nbschema.DSComputedStringAttribute("Additional comments about the wireless link."),
			"tags":          nbschema.DSTagsAttribute(),
			"custom_fields": nbschema.DSCustomFieldsAttribute(),
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *WirelessLinkDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *WirelessLinkDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data WirelessLinkDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var result *netbox.WirelessLink
	var httpResp *http.Response
	var err error

	if !data.ID.IsNull() && !data.ID.IsUnknown() {
		// Lookup by ID
		id, parseErr := utils.ParseID(data.ID.ValueString())
		if parseErr != nil {
			resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("ID must be a number, got: %s", data.ID.ValueString()))
			return
		}

		tflog.Debug(ctx, "Looking up wireless link by ID", map[string]interface{}{"id": id})

		result, httpResp, err = d.client.WirelessAPI.WirelessWirelessLinksRetrieve(ctx, id).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError("Error Reading Wireless Link",
				utils.FormatAPIError(fmt.Sprintf("read wireless link ID %d", id), err, httpResp))
			return
		}
	} else {
		resp.Diagnostics.AddError("Missing Required Attribute",
			"The 'id' attribute must be specified to look up a wireless link.")
		return
	}

	// Map the response to state
	d.mapToState(ctx, result, &data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// mapToState maps the API response to the Terraform state.
func (d *WirelessLinkDataSource) mapToState(ctx context.Context, result *netbox.WirelessLink, data *WirelessLinkDataSourceModel) {
	data.ID = types.StringValue(fmt.Sprintf("%d", result.GetId()))

	// Map interface IDs
	interfaceA := result.GetInterfaceA()
	data.InterfaceA = types.StringValue(fmt.Sprintf("%d", interfaceA.GetId()))

	interfaceB := result.GetInterfaceB()
	data.InterfaceB = types.StringValue(fmt.Sprintf("%d", interfaceB.GetId()))

	// Map optional fields
	if result.HasSsid() && result.GetSsid() != "" {
		data.SSID = types.StringValue(result.GetSsid())
	} else {
		data.SSID = types.StringNull()
	}

	if result.HasStatus() {
		status := result.GetStatus()
		data.Status = types.StringValue(string(status.GetValue()))
	} else {
		data.Status = types.StringNull()
	}

	if result.HasTenant() && result.GetTenant().Id != 0 {
		tenant := result.GetTenant()
		data.Tenant = types.StringValue(tenant.GetName())
		data.TenantID = types.StringValue(fmt.Sprintf("%d", tenant.GetId()))
	} else {
		data.Tenant = types.StringNull()
		data.TenantID = types.StringNull()
	}

	if result.HasAuthType() {
		authType := result.GetAuthType()
		data.AuthType = types.StringValue(string(authType.GetValue()))
	} else {
		data.AuthType = types.StringNull()
	}

	if result.HasAuthCipher() {
		authCipher := result.GetAuthCipher()
		data.AuthCipher = types.StringValue(string(authCipher.GetValue()))
	} else {
		data.AuthCipher = types.StringNull()
	}

	if result.HasDistance() {
		distance, ok := result.GetDistanceOk()
		if ok && distance != nil {
			data.Distance = types.Float64Value(*distance)
		} else {
			data.Distance = types.Float64Null()
		}
	} else {
		data.Distance = types.Float64Null()
	}

	if result.HasDistanceUnit() {
		distanceUnit := result.GetDistanceUnit()
		if distanceUnit.Value != nil && *distanceUnit.Value != "" {
			data.DistanceUnit = types.StringValue(string(*distanceUnit.Value))
		} else {
			data.DistanceUnit = types.StringNull()
		}
	} else {
		data.DistanceUnit = types.StringNull()
	}

	if result.HasDescription() && result.GetDescription() != "" {
		data.Description = types.StringValue(result.GetDescription())
	} else {
		data.Description = types.StringNull()
	}

	if result.HasComments() && result.GetComments() != "" {
		data.Comments = types.StringValue(result.GetComments())
	} else {
		data.Comments = types.StringNull()
	}

	// Map tags
	if result.HasTags() {
		tags := utils.NestedTagsToTagModels(result.GetTags())
		tagsValue, _ := types.SetValueFrom(ctx, utils.GetTagsAttributeType().ElemType, tags)
		data.Tags = tagsValue
	} else {
		data.Tags = types.SetNull(utils.GetTagsAttributeType().ElemType)
	}

	// Map custom fields
	if result.HasCustomFields() {
		customFields := utils.MapToCustomFieldModels(result.GetCustomFields(), nil)
		customFieldsValue, _ := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, customFields)
		data.CustomFields = customFieldsValue
	} else {
		data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)
	}
}
