// Package datasources provides Terraform data source implementations for NetBox objects.

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

var (
	_ datasource.DataSource = &VirtualDeviceContextDataSource{}

	_ datasource.DataSourceWithConfigure = &VirtualDeviceContextDataSource{}
)

// NewVirtualDeviceContextDataSource returns a new data source implementing the virtual device context data source.

func NewVirtualDeviceContextDataSource() datasource.DataSource {

	return &VirtualDeviceContextDataSource{}

}

// VirtualDeviceContextDataSource defines the data source implementation.

type VirtualDeviceContextDataSource struct {
	client *netbox.APIClient
}

// VirtualDeviceContextDataSourceModel describes the data source data model.

type VirtualDeviceContextDataSourceModel struct {
	ID types.String `tfsdk:"id"`

	Name types.String `tfsdk:"name"`

	Device types.String `tfsdk:"device"`

	DeviceID types.String `tfsdk:"device_id"`

	Identifier types.Int64 `tfsdk:"identifier"`

	Tenant types.String `tfsdk:"tenant"`

	TenantID types.String `tfsdk:"tenant_id"`

	PrimaryIP4 types.String `tfsdk:"primary_ip4"`

	PrimaryIP6 types.String `tfsdk:"primary_ip6"`

	Status types.String `tfsdk:"status"`

	Description types.String `tfsdk:"description"`

	Comments types.String `tfsdk:"comments"`

	Tags types.Set `tfsdk:"tags"`

	CustomFields types.Set `tfsdk:"custom_fields"`
}

// Metadata returns the data source type name.

func (d *VirtualDeviceContextDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {

	resp.TypeName = req.ProviderTypeName + "_virtual_device_context"

}

// Schema defines the schema for the data source.

func (d *VirtualDeviceContextDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {

	resp.Schema = schema.Schema{

		MarkdownDescription: "Use this data source to get information about a virtual device context in NetBox.",

		Attributes: map[string]schema.Attribute{

			"id": schema.StringAttribute{

				MarkdownDescription: "The unique numeric ID of the virtual device context.",

				Required: true,
			},

			"name": schema.StringAttribute{

				MarkdownDescription: "The name of the virtual device context.",

				Computed: true,
			},

			"device": schema.StringAttribute{

				MarkdownDescription: "The name of the device this VDC belongs to.",

				Computed: true,
			},

			"device_id": schema.StringAttribute{

				MarkdownDescription: "The ID of the device this VDC belongs to.",

				Computed: true,
			},

			"identifier": schema.Int64Attribute{

				MarkdownDescription: "Numeric identifier unique to the parent device.",

				Computed: true,
			},

			"tenant": schema.StringAttribute{

				MarkdownDescription: "The name of the tenant associated with this VDC.",

				Computed: true,
			},

			"tenant_id": schema.StringAttribute{

				MarkdownDescription: "The ID of the tenant associated with this VDC.",

				Computed: true,
			},

			"primary_ip4": schema.StringAttribute{

				MarkdownDescription: "Primary IPv4 address assigned to this VDC.",

				Computed: true,
			},

			"primary_ip6": schema.StringAttribute{

				MarkdownDescription: "Primary IPv6 address assigned to this VDC.",

				Computed: true,
			},

			"status": schema.StringAttribute{

				MarkdownDescription: "Operational status of the VDC.",

				Computed: true,
			},

			"description": schema.StringAttribute{

				MarkdownDescription: "A description of the virtual device context.",

				Computed: true,
			},

			"comments": schema.StringAttribute{

				MarkdownDescription: "Additional comments about the VDC.",

				Computed: true,
			},

			"tags": nbschema.DSTagsAttribute(),

			"custom_fields": nbschema.DSCustomFieldsAttribute(),
		},
	}

}

// Configure adds the provider configured client to the data source.

func (d *VirtualDeviceContextDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {

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

func (d *VirtualDeviceContextDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var data VirtualDeviceContextDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	// Parse ID

	var id int32

	_, err := fmt.Sscanf(data.ID.ValueString(), "%d", &id)

	if err != nil {

		resp.Diagnostics.AddError(

			"Invalid ID",

			fmt.Sprintf("Could not parse ID '%s': %s", data.ID.ValueString(), err.Error()),
		)

		return

	}

	tflog.Debug(ctx, "Reading virtual device context", map[string]interface{}{"id": id})

	// Read from API

	result, httpResp, err := d.client.DcimAPI.DcimVirtualDeviceContextsRetrieve(ctx, id).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		resp.Diagnostics.AddError(

			"Error reading virtual device context",

			utils.FormatAPIError(fmt.Sprintf("read virtual device context ID %d", id), err, httpResp),
		)

		return

	}

	// Map response to state

	d.mapToState(ctx, result, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {

		return

	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

// mapToState maps the API response to the Terraform state.

func (d *VirtualDeviceContextDataSource) mapToState(ctx context.Context, result *netbox.VirtualDeviceContext, data *VirtualDeviceContextDataSourceModel, diags *diag.Diagnostics) {

	data.ID = types.StringValue(fmt.Sprintf("%d", result.GetId()))

	data.Name = types.StringValue(result.GetName())

	// Map device (required field)

	device := result.GetDevice()

	data.Device = types.StringValue(device.GetName())

	data.DeviceID = types.StringValue(fmt.Sprintf("%d", device.GetId()))

	// Map identifier

	if result.HasIdentifier() {

		identifierPtr, ok := result.GetIdentifierOk()

		if ok && identifierPtr != nil {

			data.Identifier = types.Int64Value(int64(*identifierPtr))

		} else {

			data.Identifier = types.Int64Null()

		}

	} else {

		data.Identifier = types.Int64Null()

	}

	// Map tenant

	if result.HasTenant() && result.GetTenant().Id != 0 {

		tenant := result.GetTenant()

		data.Tenant = types.StringValue(tenant.GetName())

		data.TenantID = types.StringValue(fmt.Sprintf("%d", tenant.GetId()))

	} else {

		data.Tenant = types.StringNull()

		data.TenantID = types.StringNull()

	}

	// Map primary IPs

	if result.HasPrimaryIp4() && result.GetPrimaryIp4().Id != 0 {

		ip := result.GetPrimaryIp4()

		data.PrimaryIP4 = types.StringValue(ip.GetAddress())

	} else {

		data.PrimaryIP4 = types.StringNull()

	}

	if result.HasPrimaryIp6() && result.GetPrimaryIp6().Id != 0 {

		ip := result.GetPrimaryIp6()

		data.PrimaryIP6 = types.StringValue(ip.GetAddress())

	} else {

		data.PrimaryIP6 = types.StringNull()

	}

	// Map status (required field)

	status := result.GetStatus()

	data.Status = types.StringValue(string(status.GetValue()))

	// Map description

	if result.HasDescription() && result.GetDescription() != "" {

		data.Description = types.StringValue(result.GetDescription())

	} else {

		data.Description = types.StringNull()

	}

	// Map comments

	if result.HasComments() && result.GetComments() != "" {

		data.Comments = types.StringValue(result.GetComments())

	} else {

		data.Comments = types.StringNull()

	}

	// Map tags

	if result.HasTags() && len(result.GetTags()) > 0 {

		tags := utils.NestedTagsToTagModels(result.GetTags())

		tagsValue, tagDiags := types.SetValueFrom(ctx, utils.GetTagsAttributeType().ElemType, tags)

		diags.Append(tagDiags...)

		data.Tags = tagsValue

	} else {

		data.Tags = types.SetNull(utils.GetTagsAttributeType().ElemType)

	}

	// Map custom fields

	if result.HasCustomFields() && len(result.GetCustomFields()) > 0 {

		customFields := utils.MapToCustomFieldModels(result.GetCustomFields(), nil)

		cfValue, cfDiags := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, customFields)

		diags.Append(cfDiags...)

		data.CustomFields = cfValue

	} else {

		data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)

	}

}
