// Package datasources contains Terraform data source implementations for NetBox objects.
package datasources

import (
	"context"
	"fmt"

	"github.com/bab3l/go-netbox"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ datasource.DataSource              = &RoleDataSource{}
	_ datasource.DataSourceWithConfigure = &RoleDataSource{}
)

// NewRoleDataSource returns a new data source implementing the Role data source.
func NewRoleDataSource() datasource.DataSource {
	return &RoleDataSource{}
}

// RoleDataSource defines the data source implementation.
type RoleDataSource struct {
	client *netbox.APIClient
}

// RoleDataSourceModel describes the data source data model.
type RoleDataSourceModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Slug         types.String `tfsdk:"slug"`
	Weight       types.Int64  `tfsdk:"weight"`
	Description  types.String `tfsdk:"description"`
	PrefixCount  types.Int64  `tfsdk:"prefix_count"`
	VlanCount    types.Int64  `tfsdk:"vlan_count"`
	Tags         types.Set    `tfsdk:"tags"`
	CustomFields types.Set    `tfsdk:"custom_fields"`
}

// Metadata returns the data source type name.
func (d *RoleDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_role"
}

// Schema defines the schema for the data source.
func (d *RoleDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves information about an IPAM role in NetBox.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique numeric ID of the role. Use this to look up by ID.",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the role. Use this to look up by name.",
				Optional:            true,
				Computed:            true,
			},
			"slug": schema.StringAttribute{
				MarkdownDescription: "URL-friendly unique identifier for the role. Use this to look up by slug.",
				Optional:            true,
				Computed:            true,
			},
			"weight": schema.Int64Attribute{
				MarkdownDescription: "Weight for sorting.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "A description of the role.",
				Computed:            true,
			},
			"prefix_count": schema.Int64Attribute{
				MarkdownDescription: "Number of prefixes assigned to this role.",
				Computed:            true,
			},
			"vlan_count": schema.Int64Attribute{
				MarkdownDescription: "Number of VLANs assigned to this role.",
				Computed:            true,
			},
			"tags":          nbschema.DSTagsAttribute(),
			"custom_fields": nbschema.DSCustomFieldsAttribute(),
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *RoleDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// Read refreshes the data source data.
func (d *RoleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data RoleDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var role *netbox.Role

	// Look up by ID, slug, or name
	switch {
	case !data.ID.IsNull() && !data.ID.IsUnknown():
		roleID, err := utils.ParseID(data.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Invalid Role ID",
				fmt.Sprintf("Role ID must be a number, got: %s", data.ID.ValueString()),
			)
			return
		}

		tflog.Debug(ctx, "Reading role by ID", map[string]interface{}{
			"id": roleID,
		})

		r, httpResp, err := d.client.IpamAPI.IpamRolesRetrieve(ctx, roleID).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading role",
				utils.FormatAPIError(fmt.Sprintf("read role ID %d", roleID), err, httpResp),
			)
			return
		}
		role = r
	case !data.Slug.IsNull() && !data.Slug.IsUnknown():
		// Look up by slug
		tflog.Debug(ctx, "Reading role by slug", map[string]interface{}{
			"slug": data.Slug.ValueString(),
		})

		listResp, httpResp, err := d.client.IpamAPI.IpamRolesList(ctx).Slug([]string{data.Slug.ValueString()}).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading role",
				utils.FormatAPIError(fmt.Sprintf("read role by slug %s", data.Slug.ValueString()), err, httpResp),
			)
			return
		}

		if listResp.GetCount() == 0 {
			resp.Diagnostics.AddError(
				"Role not found",
				fmt.Sprintf("No role found with slug: %s", data.Slug.ValueString()),
			)
			return
		}

		if listResp.GetCount() > 1 {
			resp.Diagnostics.AddError(
				"Multiple roles found",
				fmt.Sprintf("Found %d roles with slug: %s", listResp.GetCount(), data.Slug.ValueString()),
			)
			return
		}

		role = &listResp.GetResults()[0]
	case !data.Name.IsNull() && !data.Name.IsUnknown():
		// Look up by name
		tflog.Debug(ctx, "Reading role by name", map[string]interface{}{
			"name": data.Name.ValueString(),
		})

		listResp, httpResp, err := d.client.IpamAPI.IpamRolesList(ctx).Name([]string{data.Name.ValueString()}).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading role",
				utils.FormatAPIError(fmt.Sprintf("read role by name %s", data.Name.ValueString()), err, httpResp),
			)
			return
		}

		if listResp.GetCount() == 0 {
			resp.Diagnostics.AddError(
				"Role not found",
				fmt.Sprintf("No role found with name: %s", data.Name.ValueString()),
			)
			return
		}

		if listResp.GetCount() > 1 {
			resp.Diagnostics.AddError(
				"Multiple roles found",
				fmt.Sprintf("Found %d roles with name: %s", listResp.GetCount(), data.Name.ValueString()),
			)
			return
		}

		role = &listResp.GetResults()[0]
	default:
		resp.Diagnostics.AddError(
			"Missing Required Attribute",
			"Either 'id', 'slug', or 'name' must be specified to look up a role.",
		)
		return
	}

	// Map response to model
	d.mapResponseToModel(ctx, role, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// mapResponseToModel maps the API response to the Terraform model.
func (d *RoleDataSource) mapResponseToModel(ctx context.Context, role *netbox.Role, data *RoleDataSourceModel, diags *diag.Diagnostics) {
	data.ID = types.StringValue(fmt.Sprintf("%d", role.GetId()))
	data.Name = types.StringValue(role.GetName())
	data.Slug = types.StringValue(role.GetSlug())

	// Map weight
	if weight, ok := role.GetWeightOk(); ok && weight != nil {
		data.Weight = types.Int64Value(int64(*weight))
	} else {
		data.Weight = types.Int64Null()
	}

	// Map description
	if desc, ok := role.GetDescriptionOk(); ok && desc != nil && *desc != "" {
		data.Description = types.StringValue(*desc)
	} else {
		data.Description = types.StringNull()
	}

	// Map counts
	data.PrefixCount = types.Int64Value(role.GetPrefixCount())
	data.VlanCount = types.Int64Value(role.GetVlanCount())

	// Handle tags
	if role.HasTags() && len(role.GetTags()) > 0 {
		tags := utils.NestedTagsToTagModels(role.GetTags())
		tagsValue, tagDiags := types.SetValueFrom(ctx, utils.GetTagsAttributeType().ElemType, tags)
		diags.Append(tagDiags...)
		if diags.HasError() {
			return
		}
		data.Tags = tagsValue
	} else {
		data.Tags = types.SetNull(utils.GetTagsAttributeType().ElemType)
	}

	// Handle custom fields
	if role.HasCustomFields() {
		apiCustomFields := role.GetCustomFields()
		customFields := utils.MapToCustomFieldModels(apiCustomFields, nil)
		customFieldsValue, cfDiags := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, customFields)
		diags.Append(cfDiags...)
		if diags.HasError() {
			return
		}
		data.CustomFields = customFieldsValue
	} else {
		data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)
	}
}
