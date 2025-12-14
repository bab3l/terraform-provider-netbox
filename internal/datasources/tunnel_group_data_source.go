// Package datasources contains Terraform data source implementations for the Netbox provider.

//

// This package integrates with the go-netbox OpenAPI client to provide

// read-only access to Netbox resources via Terraform data sources.

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
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.

var _ datasource.DataSource = &TunnelGroupDataSource{}

func NewTunnelGroupDataSource() datasource.DataSource {

	return &TunnelGroupDataSource{}

}

// TunnelGroupDataSource defines the data source implementation.

type TunnelGroupDataSource struct {
	client *netbox.APIClient
}

// TunnelGroupDataSourceModel describes the data source data model.

type TunnelGroupDataSourceModel struct {
	ID types.String `tfsdk:"id"`

	Name types.String `tfsdk:"name"`

	Slug types.String `tfsdk:"slug"`

	Description types.String `tfsdk:"description"`

	Tags types.Set `tfsdk:"tags"`

	CustomFields types.Set `tfsdk:"custom_fields"`
}

func (d *TunnelGroupDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {

	resp.TypeName = req.ProviderTypeName + "_tunnel_group"

}

func (d *TunnelGroupDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {

	resp.Schema = schema.Schema{

		MarkdownDescription: "Use this data source to get information about a tunnel group in Netbox. Tunnel groups are used to organize VPN tunnels. You can identify the tunnel group using `id`, `slug`, or `name`.",

		Attributes: map[string]schema.Attribute{

			"id": nbschema.DSIDAttribute("tunnel group"),

			"name": nbschema.DSNameAttribute("tunnel group"),

			"slug": nbschema.DSSlugAttribute("tunnel group"),

			"description": nbschema.DSComputedStringAttribute("Detailed description of the tunnel group."),

			"tags": nbschema.DSTagsAttribute(),

			"custom_fields": nbschema.DSCustomFieldsAttribute(),
		},
	}

}

func (d *TunnelGroupDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {

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

func (d *TunnelGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var data TunnelGroupDataSourceModel

	// Read Terraform configuration data into the model

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	var tunnelGroup *netbox.TunnelGroup

	var err error

	var httpResp *http.Response

	// Determine if we're searching by ID, slug, or name

	switch {

	case !data.ID.IsNull():

		// Search by ID

		tunnelGroupID := data.ID.ValueString()

		tflog.Debug(ctx, "Reading tunnel group by ID", map[string]interface{}{

			"id": tunnelGroupID,
		})

		// Parse the tunnel group ID to int32 for the API call

		var tunnelGroupIDInt int32

		if _, parseErr := fmt.Sscanf(tunnelGroupID, "%d", &tunnelGroupIDInt); parseErr != nil {

			resp.Diagnostics.AddError(

				"Invalid Tunnel Group ID",

				fmt.Sprintf("Tunnel Group ID must be a number, got: %s", tunnelGroupID),
			)

			return

		}

		// Retrieve the tunnel group via API

		tunnelGroup, httpResp, err = d.client.VpnAPI.VpnTunnelGroupsRetrieve(ctx, tunnelGroupIDInt).Execute()

		defer utils.CloseResponseBody(httpResp)

	case !data.Slug.IsNull():

		// Search by slug

		tunnelGroupSlug := data.Slug.ValueString()

		tflog.Debug(ctx, "Reading tunnel group by slug", map[string]interface{}{

			"slug": tunnelGroupSlug,
		})

		// List tunnel groups with slug filter

		var tunnelGroups *netbox.PaginatedTunnelGroupList

		tunnelGroups, httpResp, err = d.client.VpnAPI.VpnTunnelGroupsList(ctx).Slug([]string{tunnelGroupSlug}).Execute()

		defer utils.CloseResponseBody(httpResp)

		if err != nil {

			resp.Diagnostics.AddError(

				"Error reading tunnel group",

				utils.FormatAPIError("read tunnel group by slug", err, httpResp),
			)

			return

		}

		if len(tunnelGroups.GetResults()) == 0 {

			resp.Diagnostics.AddError(

				"Tunnel Group Not Found",

				fmt.Sprintf("No tunnel group found with slug: %s", tunnelGroupSlug),
			)

			return

		}

		if len(tunnelGroups.GetResults()) > 1 {

			resp.Diagnostics.AddError(

				"Multiple Tunnel Groups Found",

				fmt.Sprintf("Multiple tunnel groups found with slug: %s. This should not happen as slugs should be unique.", tunnelGroupSlug),
			)

			return

		}

		tunnelGroup = &tunnelGroups.GetResults()[0]

	case !data.Name.IsNull():

		// Search by name

		tunnelGroupName := data.Name.ValueString()

		tflog.Debug(ctx, "Reading tunnel group by name", map[string]interface{}{

			"name": tunnelGroupName,
		})

		// List tunnel groups with name filter

		var tunnelGroups *netbox.PaginatedTunnelGroupList

		tunnelGroups, httpResp, err = d.client.VpnAPI.VpnTunnelGroupsList(ctx).Name([]string{tunnelGroupName}).Execute()

		defer utils.CloseResponseBody(httpResp)

		if err != nil {

			resp.Diagnostics.AddError(

				"Error reading tunnel group",

				utils.FormatAPIError("read tunnel group by name", err, httpResp),
			)

			return

		}

		if len(tunnelGroups.GetResults()) == 0 {

			resp.Diagnostics.AddError(

				"Tunnel Group Not Found",

				fmt.Sprintf("No tunnel group found with name: %s", tunnelGroupName),
			)

			return

		}

		if len(tunnelGroups.GetResults()) > 1 {

			resp.Diagnostics.AddError(

				"Multiple Tunnel Groups Found",

				fmt.Sprintf("Multiple tunnel groups found with name: %s. Tunnel group names may not be unique in Netbox.", tunnelGroupName),
			)

			return

		}

		tunnelGroup = &tunnelGroups.GetResults()[0]

	default:

		resp.Diagnostics.AddError(

			"Missing Tunnel Group Identifier",

			"Either 'id', 'slug', or 'name' must be specified to identify the tunnel group.",
		)

		return

	}

	if err != nil {

		resp.Diagnostics.AddError(

			"Error reading tunnel group",

			utils.FormatAPIError("read tunnel group", err, httpResp),
		)

		return

	}

	if httpResp.StatusCode == 404 {

		resp.Diagnostics.AddError(

			"Tunnel Group Not Found",

			"The specified tunnel group was not found in Netbox.",
		)

		return

	}

	if httpResp.StatusCode != 200 {

		resp.Diagnostics.AddError(

			"Error reading tunnel group",

			fmt.Sprintf("Expected HTTP 200, got: %d", httpResp.StatusCode),
		)

		return

	}

	// Update the model with the response from the API

	data.ID = types.StringValue(fmt.Sprintf("%d", tunnelGroup.GetId()))

	data.Name = types.StringValue(tunnelGroup.GetName())

	data.Slug = types.StringValue(tunnelGroup.GetSlug())

	// Handle description

	if tunnelGroup.HasDescription() && tunnelGroup.GetDescription() != "" {

		data.Description = types.StringValue(tunnelGroup.GetDescription())

	} else {

		data.Description = types.StringNull()

	}

	// Handle tags

	if tunnelGroup.HasTags() {

		tags := utils.NestedTagsToTagModels(tunnelGroup.GetTags())

		tagsValue, diags := types.SetValueFrom(ctx, utils.GetTagsAttributeType().ElemType, tags)

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {

			return

		}

		data.Tags = tagsValue

	} else {

		data.Tags = types.SetNull(utils.GetTagsAttributeType().ElemType)

	}

	// Handle custom fields

	if tunnelGroup.HasCustomFields() {

		// For data sources, we extract all available custom fields

		customFields := utils.MapToCustomFieldModels(tunnelGroup.GetCustomFields(), nil)

		customFieldsValue, diags := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, customFields)

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {

			return

		}

		data.CustomFields = customFieldsValue

	} else {

		data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)

	}

	// Save data into Terraform state

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}
