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

var _ datasource.DataSource = &ContactRoleDataSource{}

func NewContactRoleDataSource() datasource.DataSource {
	return &ContactRoleDataSource{}
}

type ContactRoleDataSource struct {
	client *netbox.APIClient
}

type ContactRoleDataSourceModel struct {
	ID           types.String             `tfsdk:"id"`
	Name         types.String             `tfsdk:"name"`
	Slug         types.String             `tfsdk:"slug"`
	Description  types.String             `tfsdk:"description"`
	DisplayName  types.String             `tfsdk:"display_name"`
	CustomFields []utils.CustomFieldModel `tfsdk:"custom_fields"`
}

func (d *ContactRoleDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_contact_role"
}

func (d *ContactRoleDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get information about a contact role in Netbox.",

		Attributes: map[string]schema.Attribute{
			"id":            nbschema.DSIDAttribute("contact role"),
			"name":          nbschema.DSNameAttribute("contact role"),
			"slug":          nbschema.DSSlugAttribute("contact role"),
			"description":   nbschema.DSComputedStringAttribute("Description of the contact role."),
			"display_name":  nbschema.DSComputedStringAttribute("The display name of the contact role."),
			"custom_fields": nbschema.DSCustomFieldsAttribute(),
		},
	}
}

func (d *ContactRoleDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ContactRoleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ContactRoleDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var contactRole *netbox.ContactRole
	var httpResp *http.Response
	var err error

	// Lookup by ID first
	switch {
	case !data.ID.IsNull() && !data.ID.IsUnknown() && data.ID.ValueString() != "":
		id := utils.ParseInt32FromString(data.ID.ValueString())
		if id == 0 {
			resp.Diagnostics.AddError("Invalid ID", "ID must be a number")
			return
		}
		tflog.Debug(ctx, "Looking up contact role by ID", map[string]interface{}{"id": id})
		contactRole, httpResp, err = d.client.TenancyAPI.TenancyContactRolesRetrieve(ctx, id).Execute()
		defer utils.CloseResponseBody(httpResp)
	case !data.Slug.IsNull() && !data.Slug.IsUnknown() && data.Slug.ValueString() != "":
		// Lookup by slug
		slug := data.Slug.ValueString()
		tflog.Debug(ctx, "Looking up contact role by slug", map[string]interface{}{"slug": slug})
		list, listResp, listErr := d.client.TenancyAPI.TenancyContactRolesList(ctx).Slug([]string{slug}).Execute()
		defer utils.CloseResponseBody(listResp)
		httpResp = listResp
		err = listErr
		if err == nil && list != nil {
			contactRoleResult, ok := utils.ExpectSingleResult(
				list.Results,
				"Contact role not found",
				fmt.Sprintf("No contact role found with slug: %s", slug),
				"Multiple contact roles found",
				fmt.Sprintf("Found %d contact roles with slug: %s. Use ID for unique lookup.", len(list.Results), slug),
				&resp.Diagnostics,
			)
			if !ok {
				return
			}
			contactRole = contactRoleResult
		}
	case !data.Name.IsNull() && !data.Name.IsUnknown() && data.Name.ValueString() != "":
		// Lookup by name
		name := data.Name.ValueString()
		tflog.Debug(ctx, "Looking up contact role by name", map[string]interface{}{"name": name})
		list, listResp, listErr := d.client.TenancyAPI.TenancyContactRolesList(ctx).Name([]string{name}).Execute()
		defer utils.CloseResponseBody(listResp)
		httpResp = listResp
		err = listErr
		if err == nil && list != nil {
			contactRoleResult, ok := utils.ExpectSingleResult(
				list.Results,
				"Contact role not found",
				fmt.Sprintf("No contact role found with name: %s", name),
				"Multiple contact roles found",
				fmt.Sprintf("Found %d contact roles with name: %s. Use slug or ID for unique lookup.", len(list.Results), name),
				&resp.Diagnostics,
			)
			if !ok {
				return
			}
			contactRole = contactRoleResult
		}
	default:
		resp.Diagnostics.AddError("Missing identifier", "Either 'id', 'slug', or 'name' must be specified")
		return
	}

	if err != nil {
		resp.Diagnostics.AddError("Error reading contact role", utils.FormatAPIError("read contact role", err, httpResp))
		return
	}

	// Map response to state
	data.ID = types.StringValue(fmt.Sprintf("%d", contactRole.GetId()))
	data.Name = types.StringValue(contactRole.GetName())
	data.Slug = types.StringValue(contactRole.GetSlug())

	if contactRole.HasDescription() && contactRole.GetDescription() != "" {
		data.Description = types.StringValue(contactRole.GetDescription())
	} else {
		data.Description = types.StringNull()
	}

	// Map display name
	if contactRole.GetDisplay() != "" {
		data.DisplayName = types.StringValue(contactRole.GetDisplay())
	} else {
		data.DisplayName = types.StringNull()
	}

	// Handle custom fields
	if contactRole.HasCustomFields() && len(contactRole.GetCustomFields()) > 0 {
		data.CustomFields = utils.MapAllCustomFieldsToModels(contactRole.GetCustomFields())
	} else {
		data.CustomFields = nil
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
