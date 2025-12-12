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
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ datasource.DataSource = &TenantGroupDataSource{}

func NewTenantGroupDataSource() datasource.DataSource {
	return &TenantGroupDataSource{}
}

type TenantGroupDataSource struct {
	client *netbox.APIClient
}

type TenantGroupDataSourceModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Slug         types.String `tfsdk:"slug"`
	Parent       types.String `tfsdk:"parent"`
	ParentID     types.String `tfsdk:"parent_id"`
	Description  types.String `tfsdk:"description"`
	Tags         types.Set    `tfsdk:"tags"`
	CustomFields types.Set    `tfsdk:"custom_fields"`
}

func (d *TenantGroupDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tenant_group"
}

func (d *TenantGroupDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get information about a tenant group in Netbox.",
		Attributes: map[string]schema.Attribute{
			"id":            nbschema.DSIDAttribute("tenant group"),
			"name":          nbschema.DSNameAttribute("tenant group"),
			"slug":          nbschema.DSSlugAttribute("tenant group"),
			"parent":        nbschema.DSComputedStringAttribute("Name of the parent tenant group."),
			"parent_id":     nbschema.DSComputedStringAttribute("ID of the parent tenant group."),
			"description":   nbschema.DSComputedStringAttribute("Description of the tenant group."),
			"tags":          nbschema.DSTagsAttribute(),
			"custom_fields": nbschema.DSCustomFieldsAttribute(),
		},
	}
}

func (d *TenantGroupDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*netbox.APIClient)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Data Source Configure Type", fmt.Sprintf("Expected *netbox.APIClient, got: %T.", req.ProviderData))
		return
	}
	d.client = client
}

func (d *TenantGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data TenantGroupDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var tenantGroup *netbox.TenantGroup
	var err error
	var httpResp *http.Response

	switch {
	case !data.ID.IsNull():
		tflog.Debug(ctx, "Reading tenant group by ID", map[string]interface{}{"id": data.ID.ValueString()})
		tenantGroupIDInt := utils.ParseInt32FromString(data.ID.ValueString())
		if tenantGroupIDInt == 0 {
			resp.Diagnostics.AddError("Invalid Tenant Group ID", "Tenant Group ID must be a number.")
			return
		}
		tenantGroup, httpResp, err = d.client.TenancyAPI.TenancyTenantGroupsRetrieve(ctx, tenantGroupIDInt).Execute()
		defer utils.CloseResponseBody(httpResp)
	case !data.Slug.IsNull():
		tflog.Debug(ctx, "Reading tenant group by slug", map[string]interface{}{"slug": data.Slug.ValueString()})
		var tenantGroups *netbox.PaginatedTenantGroupList
		tenantGroups, httpResp, err = d.client.TenancyAPI.TenancyTenantGroupsList(ctx).Slug([]string{data.Slug.ValueString()}).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err == nil && len(tenantGroups.GetResults()) > 0 {
			tenantGroup = &tenantGroups.GetResults()[0]
		} else if err == nil {
			resp.Diagnostics.AddError("Tenant Group Not Found", fmt.Sprintf("No tenant group found with slug: %s", data.Slug.ValueString()))
			return
		}
	case !data.Name.IsNull():
		tflog.Debug(ctx, "Reading tenant group by name", map[string]interface{}{"name": data.Name.ValueString()})
		var tenantGroups *netbox.PaginatedTenantGroupList
		tenantGroups, httpResp, err = d.client.TenancyAPI.TenancyTenantGroupsList(ctx).Name([]string{data.Name.ValueString()}).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err == nil && len(tenantGroups.GetResults()) > 0 {
			tenantGroup = &tenantGroups.GetResults()[0]
		} else if err == nil {
			resp.Diagnostics.AddError("Tenant Group Not Found", fmt.Sprintf("No tenant group found with name: %s", data.Name.ValueString()))
			return
		}
	default:
		resp.Diagnostics.AddError("Missing Tenant Group Identifier", "Either 'id', 'slug', or 'name' must be specified.")
		return
	}

	if err != nil {
		resp.Diagnostics.AddError("Error reading tenant group", utils.FormatAPIError("read tenant group", err, httpResp))
		return
	}
	if httpResp == nil || httpResp.StatusCode != 200 || tenantGroup == nil {
		resp.Diagnostics.AddError("Tenant Group Not Found", "No tenant group found with the specified identifier.")
		return
	}

	d.mapTenantGroupToState(ctx, tenantGroup, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (d *TenantGroupDataSource) mapTenantGroupToState(ctx context.Context, tenantGroup *netbox.TenantGroup, data *TenantGroupDataSourceModel, diags *diag.Diagnostics) {
	data.ID = types.StringValue(fmt.Sprintf("%d", tenantGroup.GetId()))
	data.Name = types.StringValue(tenantGroup.GetName())
	data.Slug = types.StringValue(tenantGroup.GetSlug())
	data.Description = utils.StringFromAPI(tenantGroup.HasDescription(), tenantGroup.GetDescription, data.Description)

	if tenantGroup.HasParent() {
		parent := tenantGroup.GetParent()
		if parent.GetId() != 0 {
			data.Parent = types.StringValue(parent.GetName())
			data.ParentID = types.StringValue(fmt.Sprintf("%d", parent.GetId()))
		} else {
			data.Parent = types.StringNull()
			data.ParentID = types.StringNull()
		}
	} else {
		data.Parent = types.StringNull()
		data.ParentID = types.StringNull()
	}

	if tenantGroup.HasTags() {
		tags := utils.NestedTagsToTagModels(tenantGroup.GetTags())
		tagsValue, tagDiags := types.SetValueFrom(ctx, utils.GetTagsAttributeType().ElemType, tags)
		diags.Append(tagDiags...)
		if !diags.HasError() {
			data.Tags = tagsValue
		}
	} else {
		data.Tags = types.SetNull(utils.GetTagsAttributeType().ElemType)
	}

	if tenantGroup.HasCustomFields() {
		customFields := utils.MapToCustomFieldModels(tenantGroup.GetCustomFields(), []utils.CustomFieldModel{})
		customFieldsValue, cfDiags := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, customFields)
		diags.Append(cfDiags...)
		if !diags.HasError() {
			data.CustomFields = customFieldsValue
		}
	} else {
		data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)
	}
}
