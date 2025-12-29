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
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.

var (
	_ datasource.DataSource = &CableTerminationDataSource{}

	_ datasource.DataSourceWithConfigure = &CableTerminationDataSource{}
)

// NewCableTerminationDataSource returns a new data source implementing the cable termination data source.

func NewCableTerminationDataSource() datasource.DataSource {
	return &CableTerminationDataSource{}
}

// CableTerminationDataSource defines the data source implementation.

type CableTerminationDataSource struct {
	client *netbox.APIClient
}

// CableTerminationDataSourceModel describes the data source data model.

type CableTerminationDataSourceModel struct {
	ID types.String `tfsdk:"id"`

	Cable types.String `tfsdk:"cable"`

	CableEnd types.String `tfsdk:"cable_end"`

	TerminationType types.String `tfsdk:"termination_type"`

	TerminationID types.String `tfsdk:"termination_id"`

	Termination types.String `tfsdk:"termination"`
}

// Metadata returns the data source type name.

func (d *CableTerminationDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cable_termination"
}

// Schema defines the schema for the data source.

func (d *CableTerminationDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get information about a cable termination in NetBox.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique numeric ID of the cable termination.",

				Required: true,
			},

			"cable": schema.StringAttribute{
				MarkdownDescription: "The ID of the cable this termination belongs to.",

				Computed: true,
			},

			"cable_end": schema.StringAttribute{
				MarkdownDescription: "Which end of the cable this termination is on (A or B).",

				Computed: true,
			},

			"termination_type": schema.StringAttribute{
				MarkdownDescription: "The type of object this termination connects to.",

				Computed: true,
			},

			"termination_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the object this termination connects to.",

				Computed: true,
			},

			"termination": schema.StringAttribute{
				MarkdownDescription: "The display name of the termination object.",

				Computed: true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.

func (d *CableTerminationDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *CableTerminationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CableTerminationDataSourceModel

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

	tflog.Debug(ctx, "Reading cable termination", map[string]interface{}{"id": id})

	// Read from API

	result, httpResp, err := d.client.DcimAPI.DcimCableTerminationsRetrieve(ctx, id).Execute()

	defer utils.CloseResponseBody(httpResp)

	if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
		resp.Diagnostics.AddError(
			"Cable Termination Not Found",
			fmt.Sprintf("No cable termination found with ID: %d", id),
		)
		return
	}

	if err != nil {
		resp.Diagnostics.AddError(

			"Error reading cable termination",

			utils.FormatAPIError(fmt.Sprintf("read cable termination ID %d", id), err, httpResp),
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

func (d *CableTerminationDataSource) mapToState(ctx context.Context, result *netbox.CableTermination, data *CableTerminationDataSourceModel, diags *diag.Diagnostics) {
	data.ID = types.StringValue(fmt.Sprintf("%d", result.GetId()))

	// Map cable (required field)

	data.Cable = types.StringValue(fmt.Sprintf("%d", result.GetCable()))

	// Map cable end (required field)

	data.CableEnd = types.StringValue(string(result.GetCableEnd()))

	// Map termination type (required field)

	data.TerminationType = types.StringValue(result.GetTerminationType())

	// Map termination ID (required field)

	data.TerminationID = types.StringValue(fmt.Sprintf("%d", result.GetTerminationId()))

	// Map termination display (required interface{} field)

	termination := result.GetTermination()

	if termMap, ok := termination.(map[string]interface{}); ok {
		if display, hasDisplay := termMap["display"]; hasDisplay {
			if displayStr, isString := display.(string); isString {
				data.Termination = types.StringValue(displayStr)
			} else {
				data.Termination = types.StringNull()
			}
		} else {
			data.Termination = types.StringNull()
		}
	} else {
		data.Termination = types.StringNull()
	}
}
