// Package datasources contains Terraform data source implementations for NetBox objects.

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
	_ datasource.DataSource = &PowerPanelDataSource{}

	_ datasource.DataSourceWithConfigure = &PowerPanelDataSource{}
)

// NewPowerPanelDataSource returns a new data source implementing the PowerPanel data source.

func NewPowerPanelDataSource() datasource.DataSource {

	return &PowerPanelDataSource{}

}

// PowerPanelDataSource defines the data source implementation.

type PowerPanelDataSource struct {
	client *netbox.APIClient
}

// PowerPanelDataSourceModel describes the data source data model.

type PowerPanelDataSourceModel struct {
	ID types.String `tfsdk:"id"`

	Site types.String `tfsdk:"site"`

	Location types.String `tfsdk:"location"`

	Name types.String `tfsdk:"name"`

	Description types.String `tfsdk:"description"`

	Comments types.String `tfsdk:"comments"`

	Tags types.Set `tfsdk:"tags"`

	CustomFields types.Set `tfsdk:"custom_fields"`
}

// Metadata returns the data source type name.

func (d *PowerPanelDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {

	resp.TypeName = req.ProviderTypeName + "_power_panel"

}

// Schema defines the schema for the data source.

func (d *PowerPanelDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {

	resp.Schema = schema.Schema{

		MarkdownDescription: "Retrieves information about a power panel in NetBox.",

		Attributes: map[string]schema.Attribute{

			"id": schema.StringAttribute{

				MarkdownDescription: "The unique numeric ID of the power panel. Use this to look up by ID.",

				Optional: true,

				Computed: true,
			},

			"site": schema.StringAttribute{

				MarkdownDescription: "The site this power panel belongs to (ID).",

				Optional: true,

				Computed: true,
			},

			"location": schema.StringAttribute{

				MarkdownDescription: "The location within the site (ID).",

				Computed: true,
			},

			"name": schema.StringAttribute{

				MarkdownDescription: "The name of the power panel. Use with site for lookup.",

				Optional: true,

				Computed: true,
			},

			"description": schema.StringAttribute{

				MarkdownDescription: "A description of the power panel.",

				Computed: true,
			},

			"comments": schema.StringAttribute{

				MarkdownDescription: "Additional comments or notes about the power panel.",

				Computed: true,
			},

			"tags": nbschema.DSTagsAttribute(),

			"custom_fields": nbschema.DSCustomFieldsAttribute(),
		},
	}

}

// Configure adds the provider configured client to the data source.

func (d *PowerPanelDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {

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

func (d *PowerPanelDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var data PowerPanelDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	var pp *netbox.PowerPanel

	// Look up by ID if provided

	switch {

	case !data.ID.IsNull() && !data.ID.IsUnknown():

		ppID, err := utils.ParseID(data.ID.ValueString())

		if err != nil {

			resp.Diagnostics.AddError(

				"Invalid Power Panel ID",

				fmt.Sprintf("Power panel ID must be a number, got: %s", data.ID.ValueString()),
			)

			return

		}

		tflog.Debug(ctx, "Reading power panel by ID", map[string]interface{}{

			"id": ppID,
		})

		result, httpResp, err := d.client.DcimAPI.DcimPowerPanelsRetrieve(ctx, ppID).Execute()

		defer utils.CloseResponseBody(httpResp)

		if err != nil {

			resp.Diagnostics.AddError(

				"Error reading power panel",

				utils.FormatAPIError(fmt.Sprintf("read power panel ID %d", ppID), err, httpResp),
			)

			return

		}

		pp = result

	case !data.Name.IsNull() && !data.Name.IsUnknown():

		// Look up by name (optionally filtered by site)

		tflog.Debug(ctx, "Reading power panel by name", map[string]interface{}{

			"name": data.Name.ValueString(),
		})

		listReq := d.client.DcimAPI.DcimPowerPanelsList(ctx).Name([]string{data.Name.ValueString()})

		// Filter by site if provided

		if !data.Site.IsNull() && !data.Site.IsUnknown() {

			siteID, err := utils.ParseID(data.Site.ValueString())

			if err != nil {

				resp.Diagnostics.AddError(

					"Invalid Site ID",

					fmt.Sprintf("Site ID must be a number, got: %s", data.Site.ValueString()),
				)

				return

			}

			listReq = listReq.SiteId([]int32{siteID})

		}

		listResp, httpResp, err := listReq.Execute()

		defer utils.CloseResponseBody(httpResp)

		if err != nil {

			resp.Diagnostics.AddError(

				"Error reading power panel",

				utils.FormatAPIError(fmt.Sprintf("read power panel by name %s", data.Name.ValueString()), err, httpResp),
			)

			return

		}

		if listResp.GetCount() == 0 {

			resp.Diagnostics.AddError(

				"Power panel not found",

				fmt.Sprintf("No power panel found with name: %s", data.Name.ValueString()),
			)

			return

		}

		if listResp.GetCount() > 1 {

			resp.Diagnostics.AddError(

				"Multiple power panels found",

				fmt.Sprintf("Found %d power panels with name: %s. Please specify the site to narrow results.", listResp.GetCount(), data.Name.ValueString()),
			)

			return

		}

		pp = &listResp.GetResults()[0]

	default:

		resp.Diagnostics.AddError(

			"Missing Required Attribute",

			"Either 'id' or 'name' must be specified to look up a power panel.",
		)

		return

	}

	// Map response to model

	d.mapResponseToModel(ctx, pp, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {

		return

	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

// mapResponseToModel maps the API response to the Terraform model.

func (d *PowerPanelDataSource) mapResponseToModel(ctx context.Context, pp *netbox.PowerPanel, data *PowerPanelDataSourceModel, diags *diag.Diagnostics) {

	data.ID = types.StringValue(fmt.Sprintf("%d", pp.GetId()))

	data.Name = types.StringValue(pp.GetName())

	// Map site

	data.Site = types.StringValue(fmt.Sprintf("%d", pp.Site.GetId()))

	// Map location

	if pp.Location.IsSet() && pp.Location.Get() != nil {

		data.Location = types.StringValue(fmt.Sprintf("%d", pp.Location.Get().GetId()))

	} else {

		data.Location = types.StringNull()

	}

	// Map description

	if desc, ok := pp.GetDescriptionOk(); ok && desc != nil && *desc != "" {

		data.Description = types.StringValue(*desc)

	} else {

		data.Description = types.StringNull()

	}

	// Map comments

	if comments, ok := pp.GetCommentsOk(); ok && comments != nil && *comments != "" {

		data.Comments = types.StringValue(*comments)

	} else {

		data.Comments = types.StringNull()

	}

	// Handle tags

	if pp.HasTags() && len(pp.GetTags()) > 0 {

		tags := utils.NestedTagsToTagModels(pp.GetTags())

		tagsValue, tagDiags := types.SetValueFrom(ctx, utils.GetTagsAttributeType().ElemType, tags)

		diags.Append(tagDiags...)

		if diags.HasError() {

			return

		}

		data.Tags = tagsValue

	} else {

		data.Tags = types.SetNull(utils.GetTagsAttributeType().ElemType)

	}

	// Handle custom fields

	if pp.HasCustomFields() {

		apiCustomFields := pp.GetCustomFields()

		customFields := utils.MapToCustomFieldModels(apiCustomFields, nil)

		customFieldsValue, cfDiags := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, customFields)

		diags.Append(cfDiags...)

		if diags.HasError() {

			return

		}

		data.CustomFields = customFieldsValue

	} else {

		data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)

	}

}
