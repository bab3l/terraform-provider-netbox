// Package datasources provides Terraform data source implementations for NetBox objects.

package datasources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bab3l/go-netbox"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.

var (
	_ datasource.DataSource = &WirelessLANDataSource{}

	_ datasource.DataSourceWithConfigure = &WirelessLANDataSource{}
)

// NewWirelessLANDataSource returns a new data source implementing the wireless LAN data source.

func NewWirelessLANDataSource() datasource.DataSource {

	return &WirelessLANDataSource{}

}

// WirelessLANDataSource defines the data source implementation.

type WirelessLANDataSource struct {
	client *netbox.APIClient
}

// WirelessLANDataSourceModel describes the data source data model.

type WirelessLANDataSourceModel struct {
	ID types.String `tfsdk:"id"`

	SSID types.String `tfsdk:"ssid"`

	Description types.String `tfsdk:"description"`

	DisplayName types.String `tfsdk:"display_name"`

	GroupID types.Int64 `tfsdk:"group_id"`

	GroupName types.String `tfsdk:"group_name"`

	Status types.String `tfsdk:"status"`

	VLANID types.Int64 `tfsdk:"vlan_id"`

	VLANName types.String `tfsdk:"vlan_name"`

	TenantID types.Int64 `tfsdk:"tenant_id"`

	TenantName types.String `tfsdk:"tenant_name"`

	AuthType types.String `tfsdk:"auth_type"`

	AuthCipher types.String `tfsdk:"auth_cipher"`

	Comments types.String `tfsdk:"comments"`

	Tags types.Set `tfsdk:"tags"`
}

// Metadata returns the data source type name.

func (d *WirelessLANDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {

	resp.TypeName = req.ProviderTypeName + "_wireless_lan"

}

// Schema defines the schema for the data source.

func (d *WirelessLANDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {

	resp.Schema = schema.Schema{

		MarkdownDescription: "Retrieves information about a wireless LAN (WiFi network) in NetBox.",

		Attributes: map[string]schema.Attribute{

			// Filter attributes

			"id": schema.StringAttribute{

				MarkdownDescription: "The unique numeric ID of the wireless LAN. Use this to filter by ID.",

				Optional: true,

				Computed: true,
			},

			"ssid": schema.StringAttribute{

				MarkdownDescription: "The SSID (network name) of the wireless LAN. Use this to filter by SSID.",

				Optional: true,

				Computed: true,
			},

			// Computed attributes

			"description": schema.StringAttribute{

				MarkdownDescription: "A description of the wireless LAN.",

				Computed: true,
			},

			"display_name": schema.StringAttribute{

				MarkdownDescription: "Display name for the wireless LAN.",

				Computed: true,
			},

			"group_id": schema.Int64Attribute{

				MarkdownDescription: "The ID of the wireless LAN group. Use this to filter by group.",

				Optional: true,

				Computed: true,
			},

			"group_name": schema.StringAttribute{

				MarkdownDescription: "The name of the wireless LAN group.",

				Computed: true,
			},

			"status": schema.StringAttribute{

				MarkdownDescription: "Status of the wireless LAN (active, reserved, disabled, deprecated).",

				Computed: true,
			},

			"vlan_id": schema.Int64Attribute{

				MarkdownDescription: "The ID of the associated VLAN.",

				Computed: true,
			},

			"vlan_name": schema.StringAttribute{

				MarkdownDescription: "The name of the associated VLAN.",

				Computed: true,
			},

			"tenant_id": schema.Int64Attribute{

				MarkdownDescription: "The ID of the tenant this wireless LAN belongs to.",

				Computed: true,
			},

			"tenant_name": schema.StringAttribute{

				MarkdownDescription: "The name of the tenant this wireless LAN belongs to.",

				Computed: true,
			},

			"auth_type": schema.StringAttribute{

				MarkdownDescription: "Authentication type (open, wep, wpa-personal, wpa-enterprise).",

				Computed: true,
			},

			"auth_cipher": schema.StringAttribute{

				MarkdownDescription: "Authentication cipher (auto, tkip, aes).",

				Computed: true,
			},

			"comments": schema.StringAttribute{

				MarkdownDescription: "Additional comments or notes about the wireless LAN.",

				Computed: true,
			},

			"tags": schema.SetAttribute{

				MarkdownDescription: "Tags associated with this wireless LAN.",

				Computed: true,

				ElementType: types.StringType,
			},
		},
	}

}

// Configure adds the provider configured client to the data source.

func (d *WirelessLANDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {

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

// Read refreshes the data source data.

func (d *WirelessLANDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var data WirelessLANDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	var wlan *netbox.WirelessLAN

	// If ID is provided, look up directly

	if !data.ID.IsNull() && !data.ID.IsUnknown() {

		wlanID, err := utils.ParseID(data.ID.ValueString())

		if err != nil {

			resp.Diagnostics.AddError(

				"Invalid Wireless LAN ID",

				fmt.Sprintf("Wireless LAN ID must be a number, got: %s", data.ID.ValueString()),
			)

			return

		}

		tflog.Debug(ctx, "Looking up wireless LAN by ID", map[string]interface{}{

			"id": wlanID,
		})

		response, httpResp, err := d.client.WirelessAPI.WirelessWirelessLansRetrieve(ctx, wlanID).Execute()

		defer utils.CloseResponseBody(httpResp)

		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			resp.Diagnostics.AddError(
				"Wireless LAN Not Found",
				fmt.Sprintf("No wireless LAN found with ID: %d", wlanID),
			)
			return
		}

		if err != nil {

			resp.Diagnostics.AddError(

				"Error reading wireless LAN",

				utils.FormatAPIError(fmt.Sprintf("read wireless LAN ID %d", wlanID), err, httpResp),
			)

			return

		}

		wlan = response

	} else {

		// Search by filters

		tflog.Debug(ctx, "Searching for wireless LAN", map[string]interface{}{

			"ssid": data.SSID.ValueString(),

			"group_id": data.GroupID.ValueInt64(),
		})

		listReq := d.client.WirelessAPI.WirelessWirelessLansList(ctx)

		if !data.SSID.IsNull() && !data.SSID.IsUnknown() {

			listReq = listReq.Ssid([]string{data.SSID.ValueString()})

		}

		if !data.GroupID.IsNull() && !data.GroupID.IsUnknown() {

			listReq = listReq.GroupId([]string{fmt.Sprintf("%d", data.GroupID.ValueInt64())})

		}

		response, httpResp, err := listReq.Execute()

		defer utils.CloseResponseBody(httpResp)

		if err != nil {

			resp.Diagnostics.AddError(

				"Error reading wireless LANs",

				utils.FormatAPIError("list wireless LANs", err, httpResp),
			)

			return

		}

		if response.GetCount() == 0 {

			resp.Diagnostics.AddError(

				"No wireless LAN found",

				"No wireless LAN matching the specified criteria was found.",
			)

			return

		}

		if response.GetCount() > 1 {

			resp.Diagnostics.AddError(

				"Multiple wireless LANs found",

				fmt.Sprintf("Found %d wireless LANs matching the specified criteria. Please provide more specific filters.", response.GetCount()),
			)

			return

		}

		wlan = &response.GetResults()[0]

	}

	// Map response to model

	data.ID = types.StringValue(fmt.Sprintf("%d", wlan.GetId()))

	data.SSID = types.StringValue(wlan.GetSsid())

	// Map description

	if desc, ok := wlan.GetDescriptionOk(); ok && desc != nil && *desc != "" {

		data.Description = types.StringValue(*desc)

	} else {

		data.Description = types.StringNull()

	}

	// Map group

	if wlan.Group.IsSet() && wlan.Group.Get() != nil {

		group := wlan.Group.Get()

		data.GroupID = types.Int64Value(int64(group.GetId()))

		data.GroupName = types.StringValue(group.GetName())

	} else {

		data.GroupID = types.Int64Null()

		data.GroupName = types.StringNull()

	}

	// Map status

	if status, ok := wlan.GetStatusOk(); ok && status != nil {

		data.Status = types.StringValue(string(status.GetValue()))

	} else {

		data.Status = types.StringNull()

	}

	// Map VLAN

	if wlan.Vlan.IsSet() && wlan.Vlan.Get() != nil {

		vlan := wlan.Vlan.Get()

		data.VLANID = types.Int64Value(int64(vlan.GetId()))

		data.VLANName = types.StringValue(vlan.GetName())

	} else {

		data.VLANID = types.Int64Null()

		data.VLANName = types.StringNull()

	}

	// Map tenant

	if wlan.Tenant.IsSet() && wlan.Tenant.Get() != nil {

		tenant := wlan.Tenant.Get()

		data.TenantID = types.Int64Value(int64(tenant.GetId()))

		data.TenantName = types.StringValue(tenant.GetName())

	} else {

		data.TenantID = types.Int64Null()

		data.TenantName = types.StringNull()

	}

	// Map auth_type

	if authType, ok := wlan.GetAuthTypeOk(); ok && authType != nil {

		data.AuthType = types.StringValue(string(authType.GetValue()))

	} else {

		data.AuthType = types.StringNull()

	}

	// Map auth_cipher

	if authCipher, ok := wlan.GetAuthCipherOk(); ok && authCipher != nil {

		data.AuthCipher = types.StringValue(string(authCipher.GetValue()))

	} else {

		data.AuthCipher = types.StringNull()

	}

	// Map comments

	if comments, ok := wlan.GetCommentsOk(); ok && comments != nil && *comments != "" {

		data.Comments = types.StringValue(*comments)

	} else {

		data.Comments = types.StringNull()

	}

	// Map display name

	if displayName := wlan.GetDisplay(); displayName != "" {

		data.DisplayName = types.StringValue(displayName)

	} else {

		data.DisplayName = types.StringNull()

	}

	// Handle tags (simplified - just names)

	if wlan.HasTags() && len(wlan.GetTags()) > 0 {

		tagNames := make([]string, 0, len(wlan.GetTags()))

		for _, tag := range wlan.GetTags() {

			tagNames = append(tagNames, tag.GetName())

		}

		tagsValue, diags := types.SetValueFrom(ctx, types.StringType, tagNames)

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {

			return

		}

		data.Tags = tagsValue

	} else {

		data.Tags = types.SetNull(types.StringType)

	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}
