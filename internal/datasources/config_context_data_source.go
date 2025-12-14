package datasources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bab3l/go-netbox"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ datasource.DataSource = &ConfigContextDataSource{}

func NewConfigContextDataSource() datasource.DataSource {

	return &ConfigContextDataSource{}

}

// ConfigContextDataSource defines the config context data source implementation.

type ConfigContextDataSource struct {
	client *netbox.APIClient
}

// ConfigContextDataSourceModel describes the config context data source data model.

type ConfigContextDataSourceModel struct {
	ID types.String `tfsdk:"id"`

	Name types.String `tfsdk:"name"`

	Description types.String `tfsdk:"description"`

	Weight types.Int64 `tfsdk:"weight"`

	IsActive types.Bool `tfsdk:"is_active"`

	Data types.String `tfsdk:"data"`

	Regions types.Set `tfsdk:"regions"`

	SiteGroups types.Set `tfsdk:"site_groups"`

	Sites types.Set `tfsdk:"sites"`

	Locations types.Set `tfsdk:"locations"`

	DeviceTypes types.Set `tfsdk:"device_types"`

	Roles types.Set `tfsdk:"roles"`

	Platforms types.Set `tfsdk:"platforms"`

	ClusterTypes types.Set `tfsdk:"cluster_types"`

	ClusterGroups types.Set `tfsdk:"cluster_groups"`

	Clusters types.Set `tfsdk:"clusters"`

	TenantGroups types.Set `tfsdk:"tenant_groups"`

	Tenants types.Set `tfsdk:"tenants"`

	Tags types.Set `tfsdk:"tags"`
}

func (d *ConfigContextDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {

	resp.TypeName = req.ProviderTypeName + "_config_context"

}

func (d *ConfigContextDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {

	resp.Schema = schema.Schema{

		MarkdownDescription: "Use this data source to get information about a config context in Netbox.",

		Attributes: map[string]schema.Attribute{

			"id": schema.StringAttribute{

				MarkdownDescription: "The unique identifier of the config context. Specify either `id` or `name`.",

				Optional: true,

				Computed: true,
			},

			"name": schema.StringAttribute{

				MarkdownDescription: "The name of the config context. Specify either `id` or `name`.",

				Optional: true,

				Computed: true,
			},

			"description": schema.StringAttribute{

				MarkdownDescription: "A description of the config context.",

				Computed: true,
			},

			"weight": schema.Int64Attribute{

				MarkdownDescription: "The weight of the config context. Higher weight contexts override lower weight contexts.",

				Computed: true,
			},

			"is_active": schema.BoolAttribute{

				MarkdownDescription: "Whether the config context is active.",

				Computed: true,
			},

			"data": schema.StringAttribute{

				MarkdownDescription: "The JSON configuration data.",

				Computed: true,
			},

			"regions": schema.SetAttribute{

				MarkdownDescription: "Set of region IDs this config context is assigned to.",

				Computed: true,

				ElementType: types.Int64Type,
			},

			"site_groups": schema.SetAttribute{

				MarkdownDescription: "Set of site group IDs this config context is assigned to.",

				Computed: true,

				ElementType: types.Int64Type,
			},

			"sites": schema.SetAttribute{

				MarkdownDescription: "Set of site IDs this config context is assigned to.",

				Computed: true,

				ElementType: types.Int64Type,
			},

			"locations": schema.SetAttribute{

				MarkdownDescription: "Set of location IDs this config context is assigned to.",

				Computed: true,

				ElementType: types.Int64Type,
			},

			"device_types": schema.SetAttribute{

				MarkdownDescription: "Set of device type IDs this config context is assigned to.",

				Computed: true,

				ElementType: types.Int64Type,
			},

			"roles": schema.SetAttribute{

				MarkdownDescription: "Set of device role IDs this config context is assigned to.",

				Computed: true,

				ElementType: types.Int64Type,
			},

			"platforms": schema.SetAttribute{

				MarkdownDescription: "Set of platform IDs this config context is assigned to.",

				Computed: true,

				ElementType: types.Int64Type,
			},

			"cluster_types": schema.SetAttribute{

				MarkdownDescription: "Set of cluster type IDs this config context is assigned to.",

				Computed: true,

				ElementType: types.Int64Type,
			},

			"cluster_groups": schema.SetAttribute{

				MarkdownDescription: "Set of cluster group IDs this config context is assigned to.",

				Computed: true,

				ElementType: types.Int64Type,
			},

			"clusters": schema.SetAttribute{

				MarkdownDescription: "Set of cluster IDs this config context is assigned to.",

				Computed: true,

				ElementType: types.Int64Type,
			},

			"tenant_groups": schema.SetAttribute{

				MarkdownDescription: "Set of tenant group IDs this config context is assigned to.",

				Computed: true,

				ElementType: types.Int64Type,
			},

			"tenants": schema.SetAttribute{

				MarkdownDescription: "Set of tenant IDs this config context is assigned to.",

				Computed: true,

				ElementType: types.Int64Type,
			},

			"tags": schema.SetAttribute{

				MarkdownDescription: "Set of tag slugs this config context is assigned to.",

				Computed: true,

				ElementType: types.StringType,
			},
		},
	}

}

func (d *ConfigContextDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {

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

func (d *ConfigContextDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var data ConfigContextDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	// Validate that either id or name is provided

	if data.ID.IsNull() && data.Name.IsNull() {

		resp.Diagnostics.AddError(

			"Missing Required Attribute",

			"Either 'id' or 'name' must be specified to look up a config context.",
		)

		return

	}

	var result *netbox.ConfigContext

	var err error

	var httpResp *http.Response

	if !data.ID.IsNull() {

		// Lookup by ID

		id, parseErr := utils.ParseID(data.ID.ValueString())

		if parseErr != nil {

			resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Could not parse ID %s: %s", data.ID.ValueString(), parseErr))

			return

		}

		tflog.Debug(ctx, "Reading config context by ID", map[string]interface{}{

			"id": id,
		})

		result, httpResp, err = d.client.ExtrasAPI.ExtrasConfigContextsRetrieve(ctx, id).Execute()

		defer utils.CloseResponseBody(httpResp)

	} else {

		// Lookup by name

		tflog.Debug(ctx, "Reading config context by name", map[string]interface{}{

			"name": data.Name.ValueString(),
		})

		listResult, listHttpResp, listErr := d.client.ExtrasAPI.ExtrasConfigContextsList(ctx).
			Name([]string{data.Name.ValueString()}).
			Execute()

		defer utils.CloseResponseBody(listHttpResp)

		httpResp = listHttpResp

		err = listErr

		if err == nil {

			if listResult.GetCount() == 0 {

				resp.Diagnostics.AddError(

					"Config context not found",

					fmt.Sprintf("No config context found with name: %s", data.Name.ValueString()),
				)

				return

			}

			if listResult.GetCount() > 1 {

				resp.Diagnostics.AddError(

					"Multiple config contexts found",

					fmt.Sprintf("Found %d config contexts with name: %s. Please use 'id' for more specific lookup.", listResult.GetCount(), data.Name.ValueString()),
				)

				return

			}

			result = &listResult.GetResults()[0]

		}

	}

	if err != nil {

		if httpResp != nil && httpResp.StatusCode == 404 {

			resp.Diagnostics.AddError(

				"Config context not found",

				utils.FormatAPIError("find config context", err, httpResp),
			)

			return

		}

		resp.Diagnostics.AddError(

			"Error reading config context",

			utils.FormatAPIError("read config context", err, httpResp),
		)

		return

	}

	// Map response to state

	mapConfigContextDataSourceResponse(ctx, result, &data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

// Helper function to map API response to data source model.

func mapConfigContextDataSourceResponse(ctx context.Context, result *netbox.ConfigContext, data *ConfigContextDataSourceModel) {

	data.ID = types.StringValue(fmt.Sprintf("%d", result.GetId()))

	data.Name = types.StringValue(result.GetName())

	// Optional fields

	if result.HasDescription() && result.GetDescription() != "" {

		data.Description = types.StringValue(result.GetDescription())

	} else {

		data.Description = types.StringNull()

	}

	if result.HasWeight() {

		data.Weight = types.Int64Value(int64(result.GetWeight()))

	} else {

		data.Weight = types.Int64Null()

	}

	if result.HasIsActive() {

		data.IsActive = types.BoolValue(result.GetIsActive())

	} else {

		data.IsActive = types.BoolNull()

	}

	// Data field - serialize JSON back to string

	if result.Data != nil {

		jsonBytes, err := json.Marshal(result.GetData())

		if err == nil {

			data.Data = types.StringValue(string(jsonBytes))

		}

	}

	// Assignment criteria - convert from response objects to ID sets

	data.Regions = dsRegionsToSet(ctx, result.Regions)

	data.SiteGroups = dsSiteGroupsToSet(ctx, result.SiteGroups)

	data.Sites = dsSitesToSet(ctx, result.Sites)

	data.Locations = dsLocationsToSet(ctx, result.Locations)

	data.DeviceTypes = dsDeviceTypesToSet(ctx, result.DeviceTypes)

	data.Roles = dsRolesToSet(ctx, result.Roles)

	data.Platforms = dsPlatformsToSet(ctx, result.Platforms)

	data.ClusterTypes = dsClusterTypesToSet(ctx, result.ClusterTypes)

	data.ClusterGroups = dsClusterGroupsToSet(ctx, result.ClusterGroups)

	data.Clusters = dsClustersToSet(ctx, result.Clusters)

	data.TenantGroups = dsTenantGroupsToSet(ctx, result.TenantGroups)

	data.Tenants = dsTenantsToSet(ctx, result.Tenants)

	data.Tags = dsTagsSlugToSet(ctx, result.Tags)

}

// Helper functions to convert API response arrays to types.Set

func dsRegionsToSet(ctx context.Context, regions []netbox.Region) types.Set {

	if len(regions) == 0 {

		return types.SetNull(types.Int64Type)

	}

	values := make([]int64, len(regions))

	for i, r := range regions {

		values[i] = int64(r.GetId())

	}

	set, _ := types.SetValueFrom(ctx, types.Int64Type, values)

	return set

}

func dsSiteGroupsToSet(ctx context.Context, siteGroups []netbox.SiteGroup) types.Set {

	if len(siteGroups) == 0 {

		return types.SetNull(types.Int64Type)

	}

	values := make([]int64, len(siteGroups))

	for i, sg := range siteGroups {

		values[i] = int64(sg.GetId())

	}

	set, _ := types.SetValueFrom(ctx, types.Int64Type, values)

	return set

}

func dsSitesToSet(ctx context.Context, sites []netbox.Site) types.Set {

	if len(sites) == 0 {

		return types.SetNull(types.Int64Type)

	}

	values := make([]int64, len(sites))

	for i, s := range sites {

		values[i] = int64(s.GetId())

	}

	set, _ := types.SetValueFrom(ctx, types.Int64Type, values)

	return set

}

func dsLocationsToSet(ctx context.Context, locations []netbox.Location) types.Set {

	if len(locations) == 0 {

		return types.SetNull(types.Int64Type)

	}

	values := make([]int64, len(locations))

	for i, l := range locations {

		values[i] = int64(l.GetId())

	}

	set, _ := types.SetValueFrom(ctx, types.Int64Type, values)

	return set

}

func dsDeviceTypesToSet(ctx context.Context, deviceTypes []netbox.DeviceType) types.Set {

	if len(deviceTypes) == 0 {

		return types.SetNull(types.Int64Type)

	}

	values := make([]int64, len(deviceTypes))

	for i, dt := range deviceTypes {

		values[i] = int64(dt.GetId())

	}

	set, _ := types.SetValueFrom(ctx, types.Int64Type, values)

	return set

}

func dsRolesToSet(ctx context.Context, roles []netbox.DeviceRole) types.Set {

	if len(roles) == 0 {

		return types.SetNull(types.Int64Type)

	}

	values := make([]int64, len(roles))

	for i, r := range roles {

		values[i] = int64(r.GetId())

	}

	set, _ := types.SetValueFrom(ctx, types.Int64Type, values)

	return set

}

func dsPlatformsToSet(ctx context.Context, platforms []netbox.Platform) types.Set {

	if len(platforms) == 0 {

		return types.SetNull(types.Int64Type)

	}

	values := make([]int64, len(platforms))

	for i, p := range platforms {

		values[i] = int64(p.GetId())

	}

	set, _ := types.SetValueFrom(ctx, types.Int64Type, values)

	return set

}

func dsClusterTypesToSet(ctx context.Context, clusterTypes []netbox.ClusterType) types.Set {

	if len(clusterTypes) == 0 {

		return types.SetNull(types.Int64Type)

	}

	values := make([]int64, len(clusterTypes))

	for i, ct := range clusterTypes {

		values[i] = int64(ct.GetId())

	}

	set, _ := types.SetValueFrom(ctx, types.Int64Type, values)

	return set

}

func dsClusterGroupsToSet(ctx context.Context, clusterGroups []netbox.ClusterGroup) types.Set {

	if len(clusterGroups) == 0 {

		return types.SetNull(types.Int64Type)

	}

	values := make([]int64, len(clusterGroups))

	for i, cg := range clusterGroups {

		values[i] = int64(cg.GetId())

	}

	set, _ := types.SetValueFrom(ctx, types.Int64Type, values)

	return set

}

func dsClustersToSet(ctx context.Context, clusters []netbox.Cluster) types.Set {

	if len(clusters) == 0 {

		return types.SetNull(types.Int64Type)

	}

	values := make([]int64, len(clusters))

	for i, c := range clusters {

		values[i] = int64(c.GetId())

	}

	set, _ := types.SetValueFrom(ctx, types.Int64Type, values)

	return set

}

func dsTenantGroupsToSet(ctx context.Context, tenantGroups []netbox.TenantGroup) types.Set {

	if len(tenantGroups) == 0 {

		return types.SetNull(types.Int64Type)

	}

	values := make([]int64, len(tenantGroups))

	for i, tg := range tenantGroups {

		values[i] = int64(tg.GetId())

	}

	set, _ := types.SetValueFrom(ctx, types.Int64Type, values)

	return set

}

func dsTenantsToSet(ctx context.Context, tenants []netbox.Tenant) types.Set {

	if len(tenants) == 0 {

		return types.SetNull(types.Int64Type)

	}

	values := make([]int64, len(tenants))

	for i, t := range tenants {

		values[i] = int64(t.GetId())

	}

	set, _ := types.SetValueFrom(ctx, types.Int64Type, values)

	return set

}

func dsTagsSlugToSet(ctx context.Context, tags []string) types.Set {

	if len(tags) == 0 {

		return types.SetNull(types.StringType)

	}

	set, _ := types.SetValueFrom(ctx, types.StringType, tags)

	return set

}
