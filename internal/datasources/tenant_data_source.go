// Package datasources contains Terraform data source implementations for the Netbox provider.
//
// This package provides read-only access to Netbox resources for use in Terraform configurations.
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
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &TenantDataSource{}

func NewTenantDataSource() datasource.DataSource {
	return &TenantDataSource{}
}

// TenantDataSource defines the data source implementation.
type TenantDataSource struct {
	client *netbox.APIClient
}

// TenantDataSourceModel describes the data source data model.
type TenantDataSourceModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Slug         types.String `tfsdk:"slug"`
	Group        types.String `tfsdk:"group"`
	GroupID      types.String `tfsdk:"group_id"`
	Description  types.String `tfsdk:"description"`
	Comments     types.String `tfsdk:"comments"`
	Tags         types.Set    `tfsdk:"tags"`
	CustomFields types.Set    `tfsdk:"custom_fields"`
}

func (d *TenantDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tenant"
}

func (d *TenantDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get information about a tenant in Netbox. Tenants represent individual customers or organizational units in multi-tenancy scenarios. You can identify the tenant using `id`, `slug`, or `name`.",

		Attributes: map[string]schema.Attribute{
			"id":            nbschema.DSIDAttribute("tenant"),
			"name":          nbschema.DSNameAttribute("tenant"),
			"slug":          nbschema.DSSlugAttribute("tenant"),
			"group":         nbschema.DSComputedStringAttribute("Name of the tenant group that this tenant belongs to."),
			"group_id":      nbschema.DSComputedStringAttribute("ID of the tenant group that this tenant belongs to."),
			"description":   nbschema.DSComputedStringAttribute("Detailed description of the tenant."),
			"comments":      nbschema.DSComputedStringAttribute("Additional comments or notes about the tenant."),
			"tags":          nbschema.DSTagsAttribute(),
			"custom_fields": nbschema.DSCustomFieldsAttribute(),
		},
	}
}

func (d *TenantDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *TenantDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data TenantDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var tenant *netbox.Tenant
	var err error
	var httpResp *http.Response

	// Lookup by id, slug, or name
	switch {
	case !data.ID.IsNull():
		tenantID := data.ID.ValueString()
		var tenantIDInt int32
		if _, parseErr := fmt.Sscanf(tenantID, "%d", &tenantIDInt); parseErr != nil {
			resp.Diagnostics.AddError("Invalid Tenant ID", "Tenant ID must be a number.")
			return
		}
		var t *netbox.Tenant
		t, httpResp, err = d.client.TenancyAPI.TenancyTenantsRetrieve(ctx, tenantIDInt).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err == nil && httpResp.StatusCode == 200 {
			tenant = t
		}
	case !data.Slug.IsNull():
		slug := data.Slug.ValueString()
		var tenants *netbox.PaginatedTenantList
		tenants, httpResp, err = d.client.TenancyAPI.TenancyTenantsList(ctx).Slug([]string{slug}).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err == nil && httpResp.StatusCode == 200 && len(tenants.GetResults()) > 0 {
			tenant = &tenants.GetResults()[0]
		}
	case !data.Name.IsNull():
		name := data.Name.ValueString()
		var tenants *netbox.PaginatedTenantList
		tenants, httpResp, err = d.client.TenancyAPI.TenancyTenantsList(ctx).Name([]string{name}).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err == nil && httpResp.StatusCode == 200 && len(tenants.GetResults()) > 0 {
			tenant = &tenants.GetResults()[0]
		}
	default:
		resp.Diagnostics.AddError("Missing Tenant Identifier", "Either 'id', 'slug', or 'name' must be specified.")
		return
	}

	if err != nil {
		resp.Diagnostics.AddError("Error reading tenant", utils.FormatAPIError("read tenant", err, httpResp))
		return
	}

	if httpResp == nil || httpResp.StatusCode != 200 || tenant == nil {
		resp.Diagnostics.AddError("Tenant Not Found", "No tenant found with the specified identifier.")
		return
	}

	// Map API response to model using helper
	d.mapTenantToState(ctx, tenant, &data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// mapTenantToState maps API response to Terraform state using state helpers.
func (d *TenantDataSource) mapTenantToState(ctx context.Context, tenant *netbox.Tenant, data *TenantDataSourceModel) {
	data.ID = types.StringValue(fmt.Sprintf("%d", tenant.GetId()))
	data.Name = types.StringValue(tenant.GetName())
	data.Slug = types.StringValue(tenant.GetSlug())

	// Handle group reference - data source exposes both name and ID
	if tenant.HasGroup() {
		group := tenant.GetGroup()
		data.Group = types.StringValue(group.Name)
		data.GroupID = types.StringValue(fmt.Sprintf("%d", group.Id))
	} else {
		data.Group = types.StringNull()
		data.GroupID = types.StringNull()
	}

	// Handle optional string fields using helpers
	data.Description = utils.StringFromAPI(tenant.HasDescription(), tenant.GetDescription, data.Description)
	data.Comments = utils.StringFromAPI(tenant.HasComments(), tenant.GetComments, data.Comments)

	// Handle tags
	if tenant.HasTags() {
		tags := utils.NestedTagsToTagModels(tenant.GetTags())
		tagsValue, tagDiags := types.SetValueFrom(ctx, utils.GetTagsAttributeType().ElemType, tags)
		if !tagDiags.HasError() {
			data.Tags = tagsValue
		}
	} else {
		data.Tags = types.SetNull(utils.GetTagsAttributeType().ElemType)
	}

	// Handle custom fields
	if tenant.HasCustomFields() {
		customFields := utils.MapToCustomFieldModels(tenant.GetCustomFields(), []utils.CustomFieldModel{})
		customFieldsValue, cfDiags := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, customFields)
		if !cfDiags.HasError() {
			data.CustomFields = customFieldsValue
		}
	} else {
		data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)
	}
}
