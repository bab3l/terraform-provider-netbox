// Package datasources contains Terraform data source implementations for the Netbox provider.
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
var _ datasource.DataSource = &ScriptDataSource{}

// NewScriptDataSource returns a new data source implementing the script data source.
func NewScriptDataSource() datasource.DataSource {
	return &ScriptDataSource{}
}

// ScriptDataSource defines the data source implementation.
type ScriptDataSource struct {
	client *netbox.APIClient
}

// ScriptDataSourceModel describes the data source data model.
type ScriptDataSourceModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Module       types.Int64  `tfsdk:"module"`
	Description  types.String `tfsdk:"description"`
	IsExecutable types.Bool   `tfsdk:"is_executable"`
	Display      types.String `tfsdk:"display"`
}

// Metadata returns the data source type name.
func (d *ScriptDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_script"
}

// Schema defines the schema for the data source.
func (d *ScriptDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get information about a script in NetBox. Scripts are Python files loaded from the filesystem and can only be read, not created via the API.",

		Attributes: map[string]schema.Attribute{
			"id":            nbschema.DSIDAttribute("script"),
			"name":          nbschema.DSNameAttribute("script"),
			"module":        nbschema.DSComputedInt64Attribute("Module ID containing the script."),
			"description":   nbschema.DSComputedStringAttribute("Description of the script."),
			"is_executable": nbschema.DSComputedBoolAttribute("Whether the script is executable."),
			"display":       nbschema.DSComputedStringAttribute("Display name of the script."),
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *ScriptDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// Read reads the script data.
func (d *ScriptDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ScriptDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var script *netbox.Script
	var httpResp *http.Response
	var err error

	// Lookup by ID or name
	switch {
	case !data.ID.IsNull() && !data.ID.IsUnknown():
		tflog.Debug(ctx, "Reading script by ID", map[string]interface{}{
			"id": data.ID.ValueString(),
		})

		script, httpResp, err = d.client.ExtrasAPI.ExtrasScriptsRetrieve(ctx, data.ID.ValueString()).Execute()
		defer utils.CloseResponseBody(httpResp)
	case !data.Name.IsNull() && !data.Name.IsUnknown():
		// Lookup by name
		tflog.Debug(ctx, "Reading script by name", map[string]interface{}{
			"name": data.Name.ValueString(),
		})

		list, listResp, listErr := d.client.ExtrasAPI.ExtrasScriptsList(ctx).
			Name([]string{data.Name.ValueString()}).
			Execute()
		httpResp = listResp
		defer utils.CloseResponseBody(httpResp)
		err = listErr

		if err == nil {
			results := list.GetResults()
			if len(results) == 0 {
				resp.Diagnostics.AddError(
					"Not Found",
					fmt.Sprintf("No script found with name: %s", data.Name.ValueString()),
				)
				return
			}
			if len(results) > 1 {
				resp.Diagnostics.AddError(
					"Multiple Found",
					fmt.Sprintf("Multiple scripts found with name: %s. Please use id for a more specific lookup.", data.Name.ValueString()),
				)
				return
			}
			script = &results[0]
		}
	default:
		resp.Diagnostics.AddError(
			"Missing Identifier",
			"Either 'id' or 'name' must be specified to look up a script.",
		)
		return
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading script",
			utils.FormatAPIError("read script", err, httpResp),
		)
		return
	}

	// Map response to state
	d.mapResponseToState(ctx, script, &data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// mapResponseToState maps the API response to the Terraform state.
func (d *ScriptDataSource) mapResponseToState(ctx context.Context, script *netbox.Script, data *ScriptDataSourceModel) {
	data.ID = types.StringValue(fmt.Sprintf("%d", script.GetId()))
	data.Name = types.StringValue(script.GetName())
	data.Module = types.Int64Value(int64(script.GetModule()))
	data.Display = types.StringValue(script.GetDisplay())
	data.IsExecutable = types.BoolValue(script.GetIsExecutable())

	// Handle description (nullable string)
	desc, _ := script.GetDescriptionOk()
	if desc != nil && *desc != "" {
		data.Description = types.StringValue(*desc)
	} else {
		data.Description = types.StringNull()
	}
}
