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
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.

var (
	_ datasource.DataSource = &CircuitTerminationDataSource{}

	_ datasource.DataSourceWithConfigure = &CircuitTerminationDataSource{}
)

// NewCircuitTerminationDataSource returns a new Circuit Termination data source.

func NewCircuitTerminationDataSource() datasource.DataSource {

	return &CircuitTerminationDataSource{}

}

// CircuitTerminationDataSource defines the data source implementation.

type CircuitTerminationDataSource struct {
	client *netbox.APIClient
}

// CircuitTerminationDataSourceModel describes the data source data model.

type CircuitTerminationDataSourceModel struct {
	ID types.String `tfsdk:"id"`

	Circuit types.String `tfsdk:"circuit"`

	CircuitCID types.String `tfsdk:"circuit_cid"`

	TermSide types.String `tfsdk:"term_side"`

	Site types.String `tfsdk:"site"`

	SiteName types.String `tfsdk:"site_name"`

	ProviderNetwork types.String `tfsdk:"provider_network"`

	PortSpeed types.Int64 `tfsdk:"port_speed"`

	UpstreamSpeed types.Int64 `tfsdk:"upstream_speed"`

	XconnectID types.String `tfsdk:"xconnect_id"`

	PPInfo types.String `tfsdk:"pp_info"`

	Description types.String `tfsdk:"description"`

	DisplayName types.String `tfsdk:"display_name"`

	MarkConnected types.Bool `tfsdk:"mark_connected"`

	Tags types.List `tfsdk:"tags"`
}

// Metadata returns the data source type name.

func (d *CircuitTerminationDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {

	resp.TypeName = req.ProviderTypeName + "_circuit_termination"

}

// Schema defines the schema for the data source.

func (d *CircuitTerminationDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {

	resp.Schema = schema.Schema{

		MarkdownDescription: "Use this data source to retrieve information about a circuit termination in Netbox.",

		Attributes: map[string]schema.Attribute{

			"id": schema.StringAttribute{

				MarkdownDescription: "The unique numeric ID of the circuit termination. Use this to look up a termination by ID.",

				Optional: true,

				Computed: true,
			},

			"circuit": schema.StringAttribute{

				MarkdownDescription: "The ID of the circuit this termination belongs to.",

				Computed: true,
			},

			"circuit_cid": schema.StringAttribute{

				MarkdownDescription: "The CID (circuit identifier) of the circuit this termination belongs to.",

				Computed: true,
			},

			"term_side": schema.StringAttribute{

				MarkdownDescription: "The termination side (A or Z).",

				Optional: true,

				Computed: true,
			},

			"site": schema.StringAttribute{

				MarkdownDescription: "The ID of the site where this termination is located.",

				Computed: true,
			},

			"site_name": schema.StringAttribute{

				MarkdownDescription: "The name of the site where this termination is located.",

				Computed: true,
			},

			"provider_network": schema.StringAttribute{

				MarkdownDescription: "The ID of the provider network for this termination.",

				Computed: true,
			},

			"port_speed": schema.Int64Attribute{

				MarkdownDescription: "The physical circuit speed in Kbps.",

				Computed: true,
			},

			"upstream_speed": schema.Int64Attribute{

				MarkdownDescription: "The upstream speed in Kbps, if different from port speed.",

				Computed: true,
			},

			"xconnect_id": schema.StringAttribute{

				MarkdownDescription: "The ID of the local cross-connect.",

				Computed: true,
			},

			"pp_info": schema.StringAttribute{

				MarkdownDescription: "Patch panel ID and port number(s).",

				Computed: true,
			},

			"description": schema.StringAttribute{

				MarkdownDescription: "A description of the circuit termination.",

				Computed: true,
			},

			"display_name": nbschema.DSComputedStringAttribute("The display name of the circuit termination."),

			"mark_connected": schema.BoolAttribute{

				MarkdownDescription: "Whether the termination is treated as if a cable is connected.",

				Computed: true,
			},

			"tags": schema.ListAttribute{

				MarkdownDescription: "Tags assigned to this circuit termination.",

				Computed: true,

				ElementType: types.StringType,
			},
		},
	}

}

// Configure sets the client for the data source.

func (d *CircuitTerminationDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {

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

// Read reads the circuit termination data source.

func (d *CircuitTerminationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var data CircuitTerminationDataSourceModel

	// Read Terraform configuration data into the model

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	var termination *netbox.CircuitTermination

	// Look up by ID if provided

	if !data.ID.IsNull() && !data.ID.IsUnknown() {

		id, err := utils.ParseID(data.ID.ValueString())

		if err != nil {

			resp.Diagnostics.AddError(

				"Invalid ID",

				fmt.Sprintf("Could not convert ID to integer: %s", err.Error()),
			)

			return

		}

		tflog.Debug(ctx, "Looking up circuit termination by ID", map[string]interface{}{

			"id": id,
		})

		result, httpResp, err := d.client.CircuitsAPI.CircuitsCircuitTerminationsRetrieve(ctx, id).Execute()

		defer utils.CloseResponseBody(httpResp)

		if err != nil {

			resp.Diagnostics.AddError(

				"Error reading circuit termination",

				fmt.Sprintf("Could not read circuit termination with ID %d: %s\nHTTP Response: %v", id, err.Error(), httpResp),
			)

			return

		}

		termination = result

	} else {

		resp.Diagnostics.AddError(

			"Missing required attribute",

			"'id' must be specified to look up a circuit termination.",
		)

		return

	}

	// Map response to model

	d.mapResponseToModel(ctx, termination, &data)

	tflog.Debug(ctx, "Read circuit termination", map[string]interface{}{

		"id": data.ID.ValueString(),

		"term_side": data.TermSide.ValueString(),
	})

	// Save data into Terraform state

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

// mapResponseToModel maps the API response to the Terraform model.

func (d *CircuitTerminationDataSource) mapResponseToModel(ctx context.Context, termination *netbox.CircuitTermination, data *CircuitTerminationDataSourceModel) {

	data.ID = types.StringValue(fmt.Sprintf("%d", termination.GetId()))

	data.TermSide = types.StringValue(string(termination.GetTermSide()))

	// Map Circuit

	if circuit := termination.GetCircuit(); circuit.Id != 0 {

		data.Circuit = types.StringValue(fmt.Sprintf("%d", circuit.Id))

		data.CircuitCID = types.StringValue(circuit.GetCid())

	}

	// Map Site

	if site, ok := termination.GetSiteOk(); ok && site != nil && site.Id != 0 {

		data.Site = types.StringValue(fmt.Sprintf("%d", site.Id))

		data.SiteName = types.StringValue(site.GetName())

	} else {

		data.Site = types.StringNull()

		data.SiteName = types.StringNull()

	}

	// Map ProviderNetwork

	if pn, ok := termination.GetProviderNetworkOk(); ok && pn != nil && pn.Id != 0 {

		data.ProviderNetwork = types.StringValue(fmt.Sprintf("%d", pn.Id))

	} else {

		data.ProviderNetwork = types.StringNull()

	}

	// Map port_speed

	if portSpeed, ok := termination.GetPortSpeedOk(); ok && portSpeed != nil {

		data.PortSpeed = types.Int64Value(int64(*portSpeed))

	} else {

		data.PortSpeed = types.Int64Null()

	}

	// Map upstream_speed

	if upstreamSpeed, ok := termination.GetUpstreamSpeedOk(); ok && upstreamSpeed != nil {

		data.UpstreamSpeed = types.Int64Value(int64(*upstreamSpeed))

	} else {

		data.UpstreamSpeed = types.Int64Null()

	}

	// Map xconnect_id

	if xconnectID, ok := termination.GetXconnectIdOk(); ok && xconnectID != nil {

		data.XconnectID = types.StringValue(*xconnectID)

	} else {

		data.XconnectID = types.StringNull()

	}

	// Map pp_info

	if ppInfo, ok := termination.GetPpInfoOk(); ok && ppInfo != nil {

		data.PPInfo = types.StringValue(*ppInfo)

	} else {

		data.PPInfo = types.StringNull()

	}

	// Map description

	if description, ok := termination.GetDescriptionOk(); ok && description != nil {

		data.Description = types.StringValue(*description)

	} else {

		data.Description = types.StringNull()

	}

	// Map mark_connected

	if markConnected, ok := termination.GetMarkConnectedOk(); ok && markConnected != nil {

		data.MarkConnected = types.BoolValue(*markConnected)

	} else {

		data.MarkConnected = types.BoolValue(false)

	}

	// Map tags

	if tags := termination.GetTags(); len(tags) > 0 {

		tagNames := make([]string, len(tags))

		for i, tag := range tags {

			tagNames[i] = tag.Name

		}

		data.Tags, _ = types.ListValueFrom(ctx, types.StringType, tagNames)

	} else {

		data.Tags = types.ListNull(types.StringType)

	}

	// Map display name

	if termination.GetDisplay() != "" {
		data.DisplayName = types.StringValue(termination.GetDisplay())
	} else {
		data.DisplayName = types.StringNull()
	}

}
