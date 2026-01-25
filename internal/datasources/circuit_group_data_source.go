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
var _ datasource.DataSource = &CircuitGroupDataSource{}

func NewCircuitGroupDataSource() datasource.DataSource {
	return &CircuitGroupDataSource{}
}

// CircuitGroupDataSource defines the data source implementation.
type CircuitGroupDataSource struct {
	client *netbox.APIClient
}

// CircuitGroupDataSourceModel describes the data source data model.
type CircuitGroupDataSourceModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Slug         types.String `tfsdk:"slug"`
	Description  types.String `tfsdk:"description"`
	Tenant       types.String `tfsdk:"tenant"`
	TenantID     types.String `tfsdk:"tenant_id"`
	CircuitCount types.Int64  `tfsdk:"circuit_count"`
	DisplayName  types.String `tfsdk:"display_name"`
	Tags         types.Set    `tfsdk:"tags"`
	CustomFields types.Set    `tfsdk:"custom_fields"`
}

func (d *CircuitGroupDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_circuit_group"
}

func (d *CircuitGroupDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get information about a circuit group in Netbox. You can identify the circuit group using `id`, `slug`, or `name`.",
		Attributes: map[string]schema.Attribute{
			"id":            nbschema.DSIDAttribute("circuit group"),
			"name":          nbschema.DSNameAttribute("circuit group"),
			"slug":          nbschema.DSSlugAttribute("circuit group"),
			"description":   nbschema.DSComputedStringAttribute("Description of the circuit group."),
			"tenant":        nbschema.DSComputedStringAttribute("Name of the tenant."),
			"tenant_id":     nbschema.DSComputedStringAttribute("ID of the tenant."),
			"circuit_count": nbschema.DSComputedInt64Attribute("Number of circuits in this group."),
			"display_name":  nbschema.DSComputedStringAttribute("The display name of the circuit group."),
			"tags":          nbschema.DSTagsAttribute(),
			"custom_fields": nbschema.DSCustomFieldsAttribute(),
		},
	}
}

func (d *CircuitGroupDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *CircuitGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CircuitGroupDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var group *netbox.CircuitGroup

	// Lookup by ID
	switch {
	case !data.ID.IsNull() && !data.ID.IsUnknown() && data.ID.ValueString() != "":
		var idInt int32
		if _, parseErr := fmt.Sscanf(data.ID.ValueString(), "%d", &idInt); parseErr != nil {
			resp.Diagnostics.AddError(
				"Invalid ID format",
				fmt.Sprintf("Could not parse circuit group ID '%s': %s", data.ID.ValueString(), parseErr.Error()),
			)
			return
		}
		tflog.Debug(ctx, "Looking up circuit group by ID", map[string]interface{}{
			"id": idInt,
		})
		result, httpResp, err := d.client.CircuitsAPI.CircuitsCircuitGroupsRetrieve(ctx, idInt).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading circuit group",
				utils.FormatAPIError("read circuit group", err, httpResp),
			)
			return
		}
		group = result

	case !data.Slug.IsNull() && !data.Slug.IsUnknown() && data.Slug.ValueString() != "":

		// Lookup by slug
		tflog.Debug(ctx, "Looking up circuit group by slug", map[string]interface{}{
			"slug": data.Slug.ValueString(),
		})
		list, httpResp, err := d.client.CircuitsAPI.CircuitsCircuitGroupsList(ctx).
			Slug([]string{data.Slug.ValueString()}).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading circuit group",
				utils.FormatAPIError("find circuit group by slug", err, httpResp),
			)
			return
		}
		if list == nil || len(list.GetResults()) == 0 {
			resp.Diagnostics.AddError(
				"Circuit group not found",
				fmt.Sprintf("No circuit group found with slug: %s", data.Slug.ValueString()),
			)
			return
		}
		result := list.GetResults()[0]
		group = &result

	case !data.Name.IsNull() && !data.Name.IsUnknown() && data.Name.ValueString() != "":
		// Lookup by name
		tflog.Debug(ctx, "Looking up circuit group by name", map[string]interface{}{
			"name": data.Name.ValueString(),
		})
		list, httpResp, err := d.client.CircuitsAPI.CircuitsCircuitGroupsList(ctx).
			Name([]string{data.Name.ValueString()}).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading circuit group",
				utils.FormatAPIError("find circuit group by name", err, httpResp),
			)
			return
		}
		if list == nil || len(list.GetResults()) == 0 {
			resp.Diagnostics.AddError(
				"Circuit group not found",
				fmt.Sprintf("No circuit group found with name: %s", data.Name.ValueString()),
			)
			return
		}
		if len(list.GetResults()) > 1 {
			resp.Diagnostics.AddError(
				"Multiple circuit groups found",
				fmt.Sprintf("Found %d circuit groups with name '%s'. Please use a more specific identifier like 'id' or 'slug'.",
					len(list.GetResults()), data.Name.ValueString()),
			)
			return
		}
		result := list.GetResults()[0]
		group = &result

	default:
		resp.Diagnostics.AddError(
			"Missing required identifier",
			"Either 'id', 'slug', or 'name' must be specified to lookup a circuit group.",
		)
		return
	}

	// Map response to state
	d.mapResponseToState(ctx, group, &data, resp)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// mapResponseToState maps a CircuitGroup API response to the Terraform state model.
func (d *CircuitGroupDataSource) mapResponseToState(ctx context.Context, group *netbox.CircuitGroup, data *CircuitGroupDataSourceModel, resp *datasource.ReadResponse) {
	data.ID = types.StringValue(fmt.Sprintf("%d", group.GetId()))
	data.Name = types.StringValue(group.GetName())
	data.Slug = types.StringValue(group.GetSlug())

	// Description
	if group.HasDescription() && group.GetDescription() != "" {
		data.Description = types.StringValue(group.GetDescription())
	} else {
		data.Description = types.StringNull()
	}

	// Tenant
	if group.HasTenant() && group.Tenant.IsSet() && group.Tenant.Get() != nil {
		tenant := group.Tenant.Get()
		data.Tenant = types.StringValue(tenant.GetName())
		data.TenantID = types.StringValue(fmt.Sprintf("%d", tenant.GetId()))
	} else {
		data.Tenant = types.StringNull()
		data.TenantID = types.StringNull()
	}

	// Circuit count
	data.CircuitCount = types.Int64Value(group.GetCircuitCount())

	// Tags
	if group.HasTags() && len(group.GetTags()) > 0 {
		tags := utils.NestedTagsToTagModels(group.GetTags())
		tagsValue, diags := types.SetValueFrom(ctx, utils.GetTagsAttributeType().ElemType, tags)
		resp.Diagnostics.Append(diags...)
		data.Tags = tagsValue
	} else {
		data.Tags = types.SetNull(utils.GetTagsAttributeType().ElemType)
	}

	// Custom fields
	data.CustomFields = utils.CustomFieldsSetFromAPI(ctx, group.HasCustomFields(), group.GetCustomFields(), &resp.Diagnostics)

	// Map display name
	if group.GetDisplay() != "" {
		data.DisplayName = types.StringValue(group.GetDisplay())
	} else {
		data.DisplayName = types.StringNull()
	}
}
