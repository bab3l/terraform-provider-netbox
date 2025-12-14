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

var _ datasource.DataSource = &CustomLinkDataSource{}

func NewCustomLinkDataSource() datasource.DataSource {

	return &CustomLinkDataSource{}

}

// CustomLinkDataSource defines the data source implementation.

type CustomLinkDataSource struct {
	client *netbox.APIClient
}

// CustomLinkDataSourceModel describes the data source data model.

type CustomLinkDataSourceModel struct {
	ID types.String `tfsdk:"id"`

	Name types.String `tfsdk:"name"`

	ObjectTypes types.List `tfsdk:"object_types"`

	Enabled types.Bool `tfsdk:"enabled"`

	LinkText types.String `tfsdk:"link_text"`

	LinkURL types.String `tfsdk:"link_url"`

	Weight types.Int64 `tfsdk:"weight"`

	GroupName types.String `tfsdk:"group_name"`

	ButtonClass types.String `tfsdk:"button_class"`

	NewWindow types.Bool `tfsdk:"new_window"`
}

func (d *CustomLinkDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {

	resp.TypeName = req.ProviderTypeName + "_custom_link"

}

func (d *CustomLinkDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {

	resp.Schema = schema.Schema{

		MarkdownDescription: "Use this data source to get information about a custom link in Netbox.",

		Attributes: map[string]schema.Attribute{

			"id": schema.StringAttribute{

				MarkdownDescription: "Unique identifier for the custom link. Use to look up by ID.",

				Optional: true,

				Computed: true,
			},

			"name": schema.StringAttribute{

				MarkdownDescription: "Name of the custom link. Use to look up by name.",

				Optional: true,

				Computed: true,
			},

			"object_types": schema.ListAttribute{

				MarkdownDescription: "List of object types this link applies to.",

				Computed: true,

				ElementType: types.StringType,
			},

			"enabled": schema.BoolAttribute{

				MarkdownDescription: "Whether the custom link is enabled.",

				Computed: true,
			},

			"link_text": schema.StringAttribute{

				MarkdownDescription: "Jinja2 template code for the link text.",

				Computed: true,
			},

			"link_url": schema.StringAttribute{

				MarkdownDescription: "Jinja2 template code for the link URL.",

				Computed: true,
			},

			"weight": schema.Int64Attribute{

				MarkdownDescription: "Weight for ordering.",

				Computed: true,
			},

			"group_name": schema.StringAttribute{

				MarkdownDescription: "Group name for dropdown menus.",

				Computed: true,
			},

			"button_class": schema.StringAttribute{

				MarkdownDescription: "CSS class for the button.",

				Computed: true,
			},

			"new_window": schema.BoolAttribute{

				MarkdownDescription: "Whether to open the link in a new window.",

				Computed: true,
			},
		},
	}

}

func (d *CustomLinkDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {

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

func (d *CustomLinkDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var data CustomLinkDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	var result *netbox.CustomLink

	var httpResp *http.Response

	var err error

	switch {

	case !data.ID.IsNull() && data.ID.ValueString() != "":

		// Lookup by ID

		id, parseErr := utils.ParseID(data.ID.ValueString())

		if parseErr != nil {

			resp.Diagnostics.AddError("Invalid ID", "ID must be a number")

			return

		}

		result, httpResp, err = d.client.ExtrasAPI.ExtrasCustomLinksRetrieve(ctx, id).Execute()

		defer utils.CloseResponseBody(httpResp)

	case !data.Name.IsNull() && data.Name.ValueString() != "":

		// Lookup by name

		list, listResp, listErr := d.client.ExtrasAPI.ExtrasCustomLinksList(ctx).
			Name([]string{data.Name.ValueString()}).Execute()

		defer utils.CloseResponseBody(listResp)

		httpResp = listResp

		err = listErr

		if err == nil && list != nil {

			results := list.GetResults()

			if len(results) == 0 {

				resp.Diagnostics.AddError("Not Found",

					fmt.Sprintf("No custom link found with name: %s", data.Name.ValueString()))

				return

			}

			if len(results) > 1 {

				resp.Diagnostics.AddError("Multiple Found",

					fmt.Sprintf("Multiple custom links found with name: %s. Please use ID instead.", data.Name.ValueString()))

				return

			}

			result = &results[0]

		}

	default:

		resp.Diagnostics.AddError("Missing Identifier", "Either 'id' or 'name' must be specified")

		return

	}

	if err != nil {

		resp.Diagnostics.AddError("Error reading custom link",

			utils.FormatAPIError("read custom link", err, httpResp))

		return

	}

	d.mapToState(ctx, result, &data)

	tflog.Debug(ctx, "Read custom link", map[string]interface{}{

		"id": data.ID.ValueString(),

		"name": data.Name.ValueString(),
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

// mapToState maps API response to Terraform state.

func (d *CustomLinkDataSource) mapToState(ctx context.Context, result *netbox.CustomLink, data *CustomLinkDataSourceModel) {

	data.ID = types.StringValue(fmt.Sprintf("%d", result.GetId()))

	data.Name = types.StringValue(result.GetName())

	data.LinkText = types.StringValue(result.GetLinkText())

	data.LinkURL = types.StringValue(result.GetLinkUrl())

	// Handle object_types

	if len(result.GetObjectTypes()) > 0 {

		objectTypesValue, _ := types.ListValueFrom(ctx, types.StringType, result.GetObjectTypes())

		data.ObjectTypes = objectTypesValue

	} else {

		data.ObjectTypes = types.ListNull(types.StringType)

	}

	if result.HasEnabled() {

		data.Enabled = types.BoolValue(result.GetEnabled())

	} else {

		data.Enabled = types.BoolNull()

	}

	if result.HasWeight() {

		data.Weight = types.Int64Value(int64(result.GetWeight()))

	} else {

		data.Weight = types.Int64Null()

	}

	if result.HasGroupName() && result.GetGroupName() != "" {

		data.GroupName = types.StringValue(result.GetGroupName())

	} else {

		data.GroupName = types.StringNull()

	}

	if result.HasButtonClass() {

		data.ButtonClass = types.StringValue(string(*result.ButtonClass))

	} else {

		data.ButtonClass = types.StringNull()

	}

	if result.HasNewWindow() {

		data.NewWindow = types.BoolValue(result.GetNewWindow())

	} else {

		data.NewWindow = types.BoolNull()

	}

}
