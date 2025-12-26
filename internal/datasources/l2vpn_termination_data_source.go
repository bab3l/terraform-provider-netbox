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
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.

var _ datasource.DataSource = &L2VPNTerminationDataSource{}

func NewL2VPNTerminationDataSource() datasource.DataSource {

	return &L2VPNTerminationDataSource{}

}

// L2VPNTerminationDataSource defines the data source implementation.

type L2VPNTerminationDataSource struct {
	client *netbox.APIClient
}

// L2VPNTerminationDataSourceModel describes the data source data model.

type L2VPNTerminationDataSourceModel struct {
	ID types.String `tfsdk:"id"`

	DisplayName types.String `tfsdk:"display_name"`

	L2VPN types.String `tfsdk:"l2vpn"`

	AssignedObjectType types.String `tfsdk:"assigned_object_type"`

	AssignedObjectID types.Int64 `tfsdk:"assigned_object_id"`

	Tags types.Set `tfsdk:"tags"`

	CustomFields types.Set `tfsdk:"custom_fields"`
}

func (d *L2VPNTerminationDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {

	resp.TypeName = req.ProviderTypeName + "_l2vpn_termination"

}

func (d *L2VPNTerminationDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {

	resp.Schema = schema.Schema{

		MarkdownDescription: "Use this data source to get information about a Layer 2 VPN termination in Netbox.",

		Attributes: map[string]schema.Attribute{

			"id": schema.StringAttribute{

				MarkdownDescription: "ID of the L2VPN termination. Required for lookup.",

				Required: true,
			},

			"display_name": schema.StringAttribute{

				MarkdownDescription: "The display name of the L2VPN termination.",

				Computed: true,
			},

			"l2vpn": schema.StringAttribute{

				MarkdownDescription: "ID of the L2VPN this termination belongs to.",

				Computed: true,
			},

			"assigned_object_type": schema.StringAttribute{

				MarkdownDescription: "Content type of the assigned object. Valid values: `dcim.interface`, `ipam.vlan`, `virtualization.vminterface`.",

				Computed: true,
			},

			"assigned_object_id": schema.Int64Attribute{

				MarkdownDescription: "ID of the assigned object (interface or VLAN).",

				Computed: true,
			},

			"tags": nbschema.DSTagsAttribute(),

			"custom_fields": nbschema.DSCustomFieldsAttribute(),
		},
	}

}

func (d *L2VPNTerminationDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {

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

func (d *L2VPNTerminationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var data L2VPNTerminationDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	// Parse ID

	var idInt int32

	if _, err := fmt.Sscanf(data.ID.ValueString(), "%d", &idInt); err != nil {

		resp.Diagnostics.AddError(

			"Invalid ID format",

			fmt.Sprintf("Could not parse L2VPN termination ID '%s': %s", data.ID.ValueString(), err.Error()),
		)

		return

	}

	tflog.Debug(ctx, "Reading L2VPN termination", map[string]interface{}{

		"id": idInt,
	})

	// Read from API

	termination, httpResp, err := d.client.VpnAPI.VpnL2vpnTerminationsRetrieve(ctx, idInt).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		resp.Diagnostics.AddError(

			"Error reading L2VPN termination",

			utils.FormatAPIError("read L2VPN termination", err, httpResp),
		)

		return

	}

	// Map response to state

	d.mapResponseToState(ctx, termination, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {

		return

	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

// mapResponseToState maps an L2VPNTermination API response to the Terraform state model.

func (d *L2VPNTerminationDataSource) mapResponseToState(ctx context.Context, termination *netbox.L2VPNTermination, data *L2VPNTerminationDataSourceModel, diags *diag.Diagnostics) {

	data.ID = types.StringValue(fmt.Sprintf("%d", termination.GetId()))

	// Display Name

	if termination.GetDisplay() != "" {

		data.DisplayName = types.StringValue(termination.GetDisplay())

	} else {

		data.DisplayName = types.StringNull()

	}

	// L2VPN

	l2vpn := termination.GetL2vpn()

	data.L2VPN = types.StringValue(fmt.Sprintf("%d", l2vpn.GetId()))

	// Assigned object

	data.AssignedObjectType = types.StringValue(termination.GetAssignedObjectType())

	data.AssignedObjectID = types.Int64Value(termination.GetAssignedObjectId())

	// Tags

	if termination.HasTags() && len(termination.GetTags()) > 0 {

		tags := utils.NestedTagsToTagModels(termination.GetTags())

		tagsValue, d := types.SetValueFrom(ctx, utils.GetTagsAttributeType().ElemType, tags)

		diags.Append(d...)

		data.Tags = tagsValue

	} else {

		data.Tags = types.SetNull(utils.GetTagsAttributeType().ElemType)

	}

	// Custom fields

	if termination.HasCustomFields() && len(termination.GetCustomFields()) > 0 {

		customFields := utils.MapToCustomFieldModels(termination.GetCustomFields(), nil)

		customFieldsValue, d := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, customFields)

		diags.Append(d...)

		data.CustomFields = customFieldsValue

	} else {

		data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)

	}

}
