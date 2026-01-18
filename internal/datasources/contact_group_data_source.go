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

var _ datasource.DataSource = &ContactGroupDataSource{}

func NewContactGroupDataSource() datasource.DataSource {
	return &ContactGroupDataSource{}
}

type ContactGroupDataSource struct {
	client *netbox.APIClient
}

type ContactGroupDataSourceModel struct {
	ID           types.String             `tfsdk:"id"`
	Name         types.String             `tfsdk:"name"`
	Slug         types.String             `tfsdk:"slug"`
	Parent       types.String             `tfsdk:"parent"`
	ParentID     types.String             `tfsdk:"parent_id"`
	Description  types.String             `tfsdk:"description"`
	DisplayName  types.String             `tfsdk:"display_name"`
	CustomFields []utils.CustomFieldModel `tfsdk:"custom_fields"`
}

func (d *ContactGroupDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_contact_group"
}

func (d *ContactGroupDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get information about a contact group in Netbox. You can identify the contact group using `id`, `slug`, or `name`.",

		Attributes: map[string]schema.Attribute{
			"id":            nbschema.DSIDAttribute("contact group"),
			"name":          nbschema.DSNameAttribute("contact group"),
			"slug":          nbschema.DSSlugAttribute("contact group"),
			"parent":        nbschema.DSComputedStringAttribute("Name of the parent contact group."),
			"parent_id":     nbschema.DSComputedStringAttribute("ID of the parent contact group."),
			"description":   nbschema.DSComputedStringAttribute("Description of the contact group."),
			"display_name":  nbschema.DSComputedStringAttribute("The display name of the contact group."),
			"custom_fields": nbschema.DSCustomFieldsAttribute(),
		},
	}
}

func (d *ContactGroupDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ContactGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ContactGroupDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var contactGroup *netbox.ContactGroup
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
		tflog.Debug(ctx, "Looking up contact group by ID", map[string]interface{}{"id": id})
		contactGroup, httpResp, err = d.client.TenancyAPI.TenancyContactGroupsRetrieve(ctx, id).Execute()
		defer utils.CloseResponseBody(httpResp)
	case !data.Slug.IsNull() && !data.Slug.IsUnknown() && data.Slug.ValueString() != "":
		// Lookup by slug
		slug := data.Slug.ValueString()
		tflog.Debug(ctx, "Looking up contact group by slug", map[string]interface{}{"slug": slug})
		list, listResp, listErr := d.client.TenancyAPI.TenancyContactGroupsList(ctx).Slug([]string{slug}).Execute()
		defer utils.CloseResponseBody(listResp)
		httpResp = listResp
		err = listErr
		if err == nil && list != nil && len(list.Results) > 0 {
			contactGroup = &list.Results[0]
		} else if err == nil {
			resp.Diagnostics.AddError("Contact group not found", fmt.Sprintf("No contact group found with slug: %s", slug))
			return
		}
	case !data.Name.IsNull() && !data.Name.IsUnknown() && data.Name.ValueString() != "":
		// Lookup by name
		name := data.Name.ValueString()
		tflog.Debug(ctx, "Looking up contact group by name", map[string]interface{}{"name": name})
		list, listResp, listErr := d.client.TenancyAPI.TenancyContactGroupsList(ctx).Name([]string{name}).Execute()
		defer utils.CloseResponseBody(listResp)
		httpResp = listResp
		err = listErr
		if err == nil && list != nil {
			if len(list.Results) == 0 {
				resp.Diagnostics.AddError("Contact group not found", fmt.Sprintf("No contact group found with name: %s", name))
				return
			}
			if len(list.Results) > 1 {
				resp.Diagnostics.AddError("Multiple contact groups found", fmt.Sprintf("Found %d contact groups with name: %s. Use slug or ID for unique lookup.", len(list.Results), name))
				return
			}
			contactGroup = &list.Results[0]
		}
	default:
		resp.Diagnostics.AddError("Missing identifier", "Either 'id', 'slug', or 'name' must be specified")
		return
	}

	if err != nil {
		resp.Diagnostics.AddError("Error reading contact group", utils.FormatAPIError("read contact group", err, httpResp))
		return
	}

	// Map response to state
	data.ID = types.StringValue(fmt.Sprintf("%d", contactGroup.GetId()))
	data.Name = types.StringValue(contactGroup.GetName())
	data.Slug = types.StringValue(contactGroup.GetSlug())

	if contactGroup.HasDescription() && contactGroup.GetDescription() != "" {
		data.Description = types.StringValue(contactGroup.GetDescription())
	} else {
		data.Description = types.StringNull()
	}

	// Handle parent reference
	if contactGroup.HasParent() {
		parent := contactGroup.GetParent()
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

	// Map display name
	if contactGroup.GetDisplay() != "" {
		data.DisplayName = types.StringValue(contactGroup.GetDisplay())
	} else {
		data.DisplayName = types.StringNull()
	}

	// Handle custom fields
	if contactGroup.HasCustomFields() && len(contactGroup.GetCustomFields()) > 0 {
		data.CustomFields = utils.MapAllCustomFieldsToModels(contactGroup.GetCustomFields())
	} else {
		data.CustomFields = nil
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
