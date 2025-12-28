package resources

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/bab3l/go-netbox"
	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &ConfigContextResource{}
var _ resource.ResourceWithImportState = &ConfigContextResource{}

func NewConfigContextResource() resource.Resource {
	return &ConfigContextResource{}
}

// ConfigContextResource defines the config context resource implementation.
type ConfigContextResource struct {
	client *netbox.APIClient
}

// ConfigContextResourceModel describes the config context resource data model.
type ConfigContextResourceModel struct {
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Description   types.String `tfsdk:"description"`
	Weight        types.Int64  `tfsdk:"weight"`
	IsActive      types.Bool   `tfsdk:"is_active"`
	Data          types.String `tfsdk:"data"`
	Regions       types.Set    `tfsdk:"regions"`
	SiteGroups    types.Set    `tfsdk:"site_groups"`
	Sites         types.Set    `tfsdk:"sites"`
	Locations     types.Set    `tfsdk:"locations"`
	DeviceTypes   types.Set    `tfsdk:"device_types"`
	Roles         types.Set    `tfsdk:"roles"`
	Platforms     types.Set    `tfsdk:"platforms"`
	ClusterTypes  types.Set    `tfsdk:"cluster_types"`
	ClusterGroups types.Set    `tfsdk:"cluster_groups"`
	Clusters      types.Set    `tfsdk:"clusters"`
	TenantGroups  types.Set    `tfsdk:"tenant_groups"`
	Tenants       types.Set    `tfsdk:"tenants"`
	Tags          types.Set    `tfsdk:"tags"`
}

func (r *ConfigContextResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_config_context"
}

func (r *ConfigContextResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a config context in Netbox. Config contexts allow you to define arbitrary JSON data that is automatically merged and applied to devices and virtual machines based on assignment criteria.",
		Attributes: map[string]schema.Attribute{
			"id":          nbschema.IDAttribute("config context"),
			"name":        nbschema.NameAttribute("config context", 100),
			"description": nbschema.DescriptionAttribute("config context"),
			"weight": schema.Int64Attribute{
				MarkdownDescription: "Weight determines the order in which config contexts are merged. Contexts with a higher weight override those with a lower weight. Defaults to `1000`.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(1000),
				Validators: []validator.Int64{
					int64validator.Between(0, 32767),
				},
			},
			"is_active": schema.BoolAttribute{
				MarkdownDescription: "Whether this config context is active. Defaults to `true`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"data": schema.StringAttribute{
				MarkdownDescription: "JSON-formatted configuration data. This data will be merged into the configuration context for matching devices/VMs.",
				Required:            true,
			},
			"regions": schema.SetAttribute{
				MarkdownDescription: "Set of region IDs to assign this config context to.",
				Optional:            true,
				ElementType:         types.Int64Type,
			},
			"site_groups": schema.SetAttribute{
				MarkdownDescription: "Set of site group IDs to assign this config context to.",
				Optional:            true,
				ElementType:         types.Int64Type,
			},
			"sites": schema.SetAttribute{
				MarkdownDescription: "Set of site IDs to assign this config context to.",
				Optional:            true,
				ElementType:         types.Int64Type,
			},
			"locations": schema.SetAttribute{
				MarkdownDescription: "Set of location IDs to assign this config context to.",
				Optional:            true,
				ElementType:         types.Int64Type,
			},
			"device_types": schema.SetAttribute{
				MarkdownDescription: "Set of device type IDs to assign this config context to.",
				Optional:            true,
				ElementType:         types.Int64Type,
			},
			"roles": schema.SetAttribute{
				MarkdownDescription: "Set of device role IDs to assign this config context to.",
				Optional:            true,
				ElementType:         types.Int64Type,
			},
			"platforms": schema.SetAttribute{
				MarkdownDescription: "Set of platform IDs to assign this config context to.",
				Optional:            true,
				ElementType:         types.Int64Type,
			},
			"cluster_types": schema.SetAttribute{
				MarkdownDescription: "Set of cluster type IDs to assign this config context to.",
				Optional:            true,
				ElementType:         types.Int64Type,
			},
			"cluster_groups": schema.SetAttribute{
				MarkdownDescription: "Set of cluster group IDs to assign this config context to.",
				Optional:            true,
				ElementType:         types.Int64Type,
			},
			"clusters": schema.SetAttribute{
				MarkdownDescription: "Set of cluster IDs to assign this config context to.",
				Optional:            true,
				ElementType:         types.Int64Type,
			},
			"tenant_groups": schema.SetAttribute{
				MarkdownDescription: "Set of tenant group IDs to assign this config context to.",
				Optional:            true,
				ElementType:         types.Int64Type,
			},
			"tenants": schema.SetAttribute{
				MarkdownDescription: "Set of tenant IDs to assign this config context to.",
				Optional:            true,
				ElementType:         types.Int64Type,
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "Set of tag slugs to assign this config context to. Devices/VMs with any of these tags will receive this config context.",
				Optional:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

func (r *ConfigContextResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*netbox.APIClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *netbox.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	r.client = client
}

func (r *ConfigContextResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ConfigContextResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Parse the JSON data
	var jsonData interface{}
	if err := json.Unmarshal([]byte(data.Data.ValueString()), &jsonData); err != nil {
		resp.Diagnostics.AddError("Invalid JSON Data", fmt.Sprintf("Could not parse JSON data: %s", err))
		return
	}

	// Create the API request
	request := netbox.NewConfigContextRequest(data.Name.ValueString(), jsonData)

	// Set optional fields
	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		desc := data.Description.ValueString()
		request.Description = &desc
	}
	if !data.Weight.IsNull() && !data.Weight.IsUnknown() {
		weight, err := utils.SafeInt32FromValue(data.Weight)
		if err != nil {
			resp.Diagnostics.AddError("Invalid value", fmt.Sprintf("Weight value overflow: %s", err))
			return
		}
		request.Weight = &weight
	}
	if !data.IsActive.IsNull() && !data.IsActive.IsUnknown() {
		isActive := data.IsActive.ValueBool()
		request.IsActive = &isActive
	}

	// Set assignment criteria
	if !data.Regions.IsNull() && !data.Regions.IsUnknown() {
		request.Regions = setToInt32Slice(ctx, data.Regions)
	}

	if !data.SiteGroups.IsNull() && !data.SiteGroups.IsUnknown() {
		request.SiteGroups = setToInt32Slice(ctx, data.SiteGroups)
	}

	if !data.Sites.IsNull() && !data.Sites.IsUnknown() {
		request.Sites = setToInt32Slice(ctx, data.Sites)
	}

	if !data.Locations.IsNull() && !data.Locations.IsUnknown() {
		request.Locations = setToInt32Slice(ctx, data.Locations)
	}

	if !data.DeviceTypes.IsNull() && !data.DeviceTypes.IsUnknown() {
		request.DeviceTypes = setToInt32Slice(ctx, data.DeviceTypes)
	}

	if !data.Roles.IsNull() && !data.Roles.IsUnknown() {
		request.Roles = setToInt32Slice(ctx, data.Roles)
	}

	if !data.Platforms.IsNull() && !data.Platforms.IsUnknown() {
		request.Platforms = setToInt32Slice(ctx, data.Platforms)
	}

	if !data.ClusterTypes.IsNull() && !data.ClusterTypes.IsUnknown() {
		request.ClusterTypes = setToInt32Slice(ctx, data.ClusterTypes)
	}

	if !data.ClusterGroups.IsNull() && !data.ClusterGroups.IsUnknown() {
		request.ClusterGroups = setToInt32Slice(ctx, data.ClusterGroups)
	}

	if !data.Clusters.IsNull() && !data.Clusters.IsUnknown() {
		request.Clusters = setToInt32Slice(ctx, data.Clusters)
	}

	if !data.TenantGroups.IsNull() && !data.TenantGroups.IsUnknown() {
		request.TenantGroups = setToInt32Slice(ctx, data.TenantGroups)
	}

	if !data.Tenants.IsNull() && !data.Tenants.IsUnknown() {
		request.Tenants = setToInt32Slice(ctx, data.Tenants)
	}

	if !data.Tags.IsNull() && !data.Tags.IsUnknown() {
		request.Tags = setToStringSlice(ctx, data.Tags)
	}

	tflog.Debug(ctx, "Creating config context", map[string]interface{}{
		"name": data.Name.ValueString(),
	})

	result, httpResp, err := r.client.ExtrasAPI.ExtrasConfigContextsCreate(ctx).
		ConfigContextRequest(*request).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating config context",
			utils.FormatAPIError("create config context", err, httpResp),
		)
		return
	}

	// Map response to state
	mapConfigContextResponseToModel(ctx, result, &data)
	tflog.Trace(ctx, "Created config context", map[string]interface{}{
		"id":   data.ID.ValueString(),
		"name": data.Name.ValueString(),
	})
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ConfigContextResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ConfigContextResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Could not parse ID %s: %s", data.ID.ValueString(), err))
		return
	}
	tflog.Debug(ctx, "Reading config context", map[string]interface{}{
		"id": id,
	})

	result, httpResp, err := r.client.ExtrasAPI.ExtrasConfigContextsRetrieve(ctx, id).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			tflog.Warn(ctx, "Config context not found, removing from state", map[string]interface{}{
				"id": id,
			})
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading config context",
			utils.FormatAPIError("read config context", err, httpResp),
		)
		return
	}

	// Map response to state
	mapConfigContextResponseToModel(ctx, result, &data)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ConfigContextResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ConfigContextResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Could not parse ID %s: %s", data.ID.ValueString(), err))
		return
	}

	// Parse the JSON data
	var jsonData interface{}
	if err := json.Unmarshal([]byte(data.Data.ValueString()), &jsonData); err != nil {
		resp.Diagnostics.AddError("Invalid JSON Data", fmt.Sprintf("Could not parse JSON data: %s", err))
		return
	}

	// Create the API request
	request := netbox.NewConfigContextRequest(data.Name.ValueString(), jsonData)

	// Set optional fields
	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		request.SetDescription(data.Description.ValueString())
	}

	if !data.Weight.IsNull() && !data.Weight.IsUnknown() {
		weight, err := utils.SafeInt32FromValue(data.Weight)
		if err != nil {
			resp.Diagnostics.AddError("Invalid value", fmt.Sprintf("Weight value overflow: %s", err))
			return
		}
		request.SetWeight(weight)
	}

	if !data.IsActive.IsNull() && !data.IsActive.IsUnknown() {
		isActive := data.IsActive.ValueBool()
		request.SetIsActive(isActive)
	}

	// Set assignment criteria - for update, always set them (even if empty)
	request.Regions = setToInt32Slice(ctx, data.Regions)
	request.SiteGroups = setToInt32Slice(ctx, data.SiteGroups)
	request.Sites = setToInt32Slice(ctx, data.Sites)
	request.Locations = setToInt32Slice(ctx, data.Locations)
	request.DeviceTypes = setToInt32Slice(ctx, data.DeviceTypes)
	request.Roles = setToInt32Slice(ctx, data.Roles)
	request.Platforms = setToInt32Slice(ctx, data.Platforms)
	request.ClusterTypes = setToInt32Slice(ctx, data.ClusterTypes)
	request.ClusterGroups = setToInt32Slice(ctx, data.ClusterGroups)
	request.Clusters = setToInt32Slice(ctx, data.Clusters)
	request.TenantGroups = setToInt32Slice(ctx, data.TenantGroups)
	request.Tenants = setToInt32Slice(ctx, data.Tenants)
	request.Tags = setToStringSlice(ctx, data.Tags)

	tflog.Debug(ctx, "Updating config context", map[string]interface{}{
		"id":   id,
		"name": data.Name.ValueString(),
	})

	result, httpResp, err := r.client.ExtrasAPI.ExtrasConfigContextsUpdate(ctx, id).
		ConfigContextRequest(*request).Execute()

	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating config context",
			utils.FormatAPIError("update config context", err, httpResp),
		)
		return
	}

	// Map response to state
	mapConfigContextResponseToModel(ctx, result, &data)
	tflog.Trace(ctx, "Updated config context", map[string]interface{}{
		"id":   data.ID.ValueString(),
		"name": data.Name.ValueString(),
	})
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ConfigContextResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ConfigContextResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Could not parse ID %s: %s", data.ID.ValueString(), err))
		return
	}
	tflog.Debug(ctx, "Deleting config context", map[string]interface{}{
		"id": id,
	})

	httpResp, err := r.client.ExtrasAPI.ExtrasConfigContextsDestroy(ctx, id).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		// If the resource was already deleted (404), consider it a success
		if httpResp != nil && httpResp.StatusCode == 404 {
			tflog.Debug(ctx, "Config context already deleted", map[string]interface{}{"id": id})
			return
		}
		resp.Diagnostics.AddError(
			"Error deleting config context",
			utils.FormatAPIError("delete config context", err, httpResp),
		)
		return
	}
	tflog.Trace(ctx, "Deleted config context", map[string]interface{}{
		"id": id,
	})
}

func (r *ConfigContextResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Helper function to convert types.Set of Int64 to []int32.
func setToInt32Slice(ctx context.Context, set types.Set) []int32 {
	if set.IsNull() || set.IsUnknown() {
		return []int32{}
	}
	var int64Values []int64
	set.ElementsAs(ctx, &int64Values, false)
	result := make([]int32, len(int64Values))
	for i, v := range int64Values {
		val32, err := utils.SafeInt32(v)
		if err != nil {
			// Return empty slice on error - caller should validate data
			return []int32{}
		}
		result[i] = val32
	}
	return result
}

// Helper function to convert types.Set of String to []string.
func setToStringSlice(ctx context.Context, set types.Set) []string {
	if set.IsNull() || set.IsUnknown() {
		return []string{}
	}
	var stringValues []string
	set.ElementsAs(ctx, &stringValues, false)
	return stringValues
}

// Helper function to map API response to model.
func mapConfigContextResponseToModel(ctx context.Context, result *netbox.ConfigContext, data *ConfigContextResourceModel) {
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
		data.Weight = types.Int64Value(1000) // Default
	}

	if result.HasIsActive() {
		data.IsActive = types.BoolValue(result.GetIsActive())
	} else {
		data.IsActive = types.BoolValue(true) // Default
	}

	// Data field - serialize JSON back to string
	if result.Data != nil {
		jsonBytes, err := json.Marshal(result.GetData())
		if err == nil {
			data.Data = types.StringValue(string(jsonBytes))
		}
	}

	// Assignment criteria - convert from response objects to ID sets
	data.Regions = regionsToSet(ctx, result.Regions)
	data.SiteGroups = siteGroupsToSet(ctx, result.SiteGroups)
	data.Sites = sitesToSet(ctx, result.Sites)
	data.Locations = locationsToSet(ctx, result.Locations)
	data.DeviceTypes = deviceTypesToSet(ctx, result.DeviceTypes)
	data.Roles = rolesToSet(ctx, result.Roles)
	data.Platforms = platformsToSet(ctx, result.Platforms)
	data.ClusterTypes = clusterTypesToSet(ctx, result.ClusterTypes)
	data.ClusterGroups = clusterGroupsToSet(ctx, result.ClusterGroups)
	data.Clusters = clustersToSet(ctx, result.Clusters)
	data.TenantGroups = tenantGroupsToSet(ctx, result.TenantGroups)
	data.Tenants = tenantsToSet(ctx, result.Tenants)
	data.Tags = tagsSlugToSet(ctx, result.Tags)
}

// Helper functions to convert API response arrays to types.Set.
func regionsToSet(ctx context.Context, regions []netbox.Region) types.Set {
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

func siteGroupsToSet(ctx context.Context, siteGroups []netbox.SiteGroup) types.Set {
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

func sitesToSet(ctx context.Context, sites []netbox.Site) types.Set {
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

func locationsToSet(ctx context.Context, locations []netbox.Location) types.Set {
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

func deviceTypesToSet(ctx context.Context, deviceTypes []netbox.DeviceType) types.Set {
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

func rolesToSet(ctx context.Context, roles []netbox.DeviceRole) types.Set {
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

func platformsToSet(ctx context.Context, platforms []netbox.Platform) types.Set {
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

func clusterTypesToSet(ctx context.Context, clusterTypes []netbox.ClusterType) types.Set {
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

func clusterGroupsToSet(ctx context.Context, clusterGroups []netbox.ClusterGroup) types.Set {
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

func clustersToSet(ctx context.Context, clusters []netbox.Cluster) types.Set {
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

func tenantGroupsToSet(ctx context.Context, tenantGroups []netbox.TenantGroup) types.Set {
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

func tenantsToSet(ctx context.Context, tenants []netbox.Tenant) types.Set {
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

func tagsSlugToSet(ctx context.Context, tags []string) types.Set {
	if len(tags) == 0 {
		return types.SetNull(types.StringType)
	}
	set, _ := types.SetValueFrom(ctx, types.StringType, tags)
	return set
}
