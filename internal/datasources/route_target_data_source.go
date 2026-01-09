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
	_ datasource.DataSource              = &RouteTargetDataSource{}
	_ datasource.DataSourceWithConfigure = &RouteTargetDataSource{}
)

// NewRouteTargetDataSource returns a new RouteTarget data source.
func NewRouteTargetDataSource() datasource.DataSource {
	return &RouteTargetDataSource{}
}

// RouteTargetDataSource defines the data source implementation.
type RouteTargetDataSource struct {
	client *netbox.APIClient
}

// RouteTargetDataSourceModel describes the data source data model.
type RouteTargetDataSourceModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Tenant       types.String `tfsdk:"tenant"`
	TenantName   types.String `tfsdk:"tenant_name"`
	Description  types.String `tfsdk:"description"`
	Comments     types.String `tfsdk:"comments"`
	Tags         types.List   `tfsdk:"tags"`
	CustomFields types.Set    `tfsdk:"custom_fields"`
	DisplayName  types.String `tfsdk:"display_name"`
}

// Metadata returns the data source type name.
func (d *RouteTargetDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_route_target"
}

// Schema defines the schema for the data source.
func (d *RouteTargetDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get information about a Route Target in Netbox.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the route target. Either `id` or `name` must be specified.",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The route target value (formatted in accordance with RFC 4360).",
				Optional:            true,
				Computed:            true,
			},
			"tenant": schema.StringAttribute{
				MarkdownDescription: "The ID of the tenant that owns this route target.",
				Computed:            true,
			},
			"tenant_name": schema.StringAttribute{
				MarkdownDescription: "The name of the tenant that owns this route target.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the route target.",
				Computed:            true,
			},
			"comments": schema.StringAttribute{
				MarkdownDescription: "Comments about the route target.",
				Computed:            true,
			},
			"tags": schema.ListAttribute{
				MarkdownDescription: "The tags assigned to this route target.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"custom_fields": nbschema.DSCustomFieldsAttribute(),
			"display_name": schema.StringAttribute{
				MarkdownDescription: "The display name of the route target.",
				Computed:            true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *RouteTargetDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// Read refreshes the Terraform state with the latest data.
func (d *RouteTargetDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data RouteTargetDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var rt *netbox.RouteTarget

	// Check if we're looking up by ID
	switch {
	case utils.IsSet(data.ID):
		var idInt int
		_, err := fmt.Sscanf(data.ID.ValueString(), "%d", &idInt)
		if err != nil {
			resp.Diagnostics.AddError(
				"Invalid ID",
				fmt.Sprintf("Unable to parse ID %q: %s", data.ID.ValueString(), err.Error()),
			)
			return
		}
		tflog.Debug(ctx, "Reading RouteTarget by ID", map[string]interface{}{
			"id": idInt,
		})
		id32, err := utils.SafeInt32(int64(idInt))
		if err != nil {
			resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("ID value overflow: %s", err))
			return
		}
		result, httpResp, err := d.client.IpamAPI.IpamRouteTargetsRetrieve(ctx, id32).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading RouteTarget",
				utils.FormatAPIError(fmt.Sprintf("retrieve RouteTarget ID %d", idInt), err, httpResp),
			)
			return
		}
		rt = result

	case utils.IsSet(data.Name):
		// Looking up by name
		tflog.Debug(ctx, "Reading RouteTarget by name", map[string]interface{}{
			"name": data.Name.ValueString(),
		})
		listReq := d.client.IpamAPI.IpamRouteTargetsList(ctx)
		listReq = listReq.Name([]string{data.Name.ValueString()})
		results, httpResp, err := listReq.Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error listing RouteTargets",
				utils.FormatAPIError(fmt.Sprintf("list RouteTargets with name %q", data.Name.ValueString()), err, httpResp),
			)
			return
		}
		if results.Count == 0 {
			resp.Diagnostics.AddError(
				"RouteTarget not found",
				fmt.Sprintf("No RouteTarget found with name %q", data.Name.ValueString()),
			)
			return
		}
		if results.Count > 1 {
			resp.Diagnostics.AddError(
				"Multiple RouteTargets found",
				fmt.Sprintf("Found %d RouteTargets with name %q. Please use 'id' to specify the exact RouteTarget.", results.Count, data.Name.ValueString()),
			)
			return
		}
		rt = &results.Results[0]

	default:
		resp.Diagnostics.AddError(
			"Missing search criteria",
			"Either 'id' or 'name' must be specified to look up a RouteTarget.",
		)
		return
	}

	// Map response to model
	d.mapRouteTargetToDataSourceModel(ctx, rt, &data)
	tflog.Debug(ctx, "Read RouteTarget", map[string]interface{}{
		"id":   data.ID.ValueString(),
		"name": data.Name.ValueString(),
	})

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// mapRouteTargetToDataSourceModel maps a Netbox RouteTarget to the Terraform data source model.
func (d *RouteTargetDataSource) mapRouteTargetToDataSourceModel(ctx context.Context, rt *netbox.RouteTarget, data *RouteTargetDataSourceModel) {
	data.ID = types.StringValue(fmt.Sprintf("%d", rt.Id))
	data.Name = types.StringValue(rt.Name)

	// Tenant
	if rt.HasTenant() && rt.Tenant.Get() != nil {
		tenant := rt.Tenant.Get()
		data.Tenant = types.StringValue(fmt.Sprintf("%d", tenant.GetId()))
		data.TenantName = types.StringValue(tenant.GetName())
	} else {
		data.Tenant = types.StringNull()
		data.TenantName = types.StringNull()
	}

	// Description
	if rt.Description != nil && *rt.Description != "" {
		data.Description = types.StringValue(*rt.Description)
	} else {
		data.Description = types.StringNull()
	}

	// Comments
	if rt.Comments != nil && *rt.Comments != "" {
		data.Comments = types.StringValue(*rt.Comments)
	} else {
		data.Comments = types.StringNull()
	}

	// Tags - convert to list of strings (tag names)
	if len(rt.Tags) > 0 {
		tagNames := make([]string, len(rt.Tags))
		for i, tag := range rt.Tags {
			tagNames[i] = tag.Name
		}
		tagsList, _ := types.ListValueFrom(ctx, types.StringType, tagNames)
		data.Tags = tagsList
	} else {
		data.Tags = types.ListNull(types.StringType)
	}

	// Handle custom fields - datasources return ALL fields
	if rt.HasCustomFields() {
		customFields := utils.MapAllCustomFieldsToModels(rt.GetCustomFields())
		customFieldsValue, cfDiags := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, customFields)
		if !cfDiags.HasError() {
			data.CustomFields = customFieldsValue
		}
	} else {
		data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)
	}

	// Map display_name
	if rt.Display != "" {
		data.DisplayName = types.StringValue(rt.Display)
	} else {
		data.DisplayName = types.StringNull()
	}
}
