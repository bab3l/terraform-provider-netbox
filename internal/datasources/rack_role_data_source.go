// Package datasources contains Terraform data source implementations for the Netbox provider.
//
// This package integrates with the go-netbox OpenAPI client to provide
// read-only access to Netbox resources via Terraform data sources.
package datasources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bab3l/go-netbox"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/bab3l/terraform-provider-netbox/internal/utils"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &RackRoleDataSource{}

func NewRackRoleDataSource() datasource.DataSource {
	return &RackRoleDataSource{}
}

// RackRoleDataSource defines the data source implementation.
type RackRoleDataSource struct {
	client *netbox.APIClient
}

// RackRoleDataSourceModel describes the data source data model.
type RackRoleDataSourceModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Slug         types.String `tfsdk:"slug"`
	Color        types.String `tfsdk:"color"`
	Description  types.String `tfsdk:"description"`
	Tags         types.Set    `tfsdk:"tags"`
	CustomFields types.Set    `tfsdk:"custom_fields"`
}

func (d *RackRoleDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_rack_role"
}

func (d *RackRoleDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get information about a rack role in Netbox. Rack roles categorize racks by their function (e.g., Production, Testing, Storage). You can identify the rack role using `id`, `slug`, or `name`.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier for the rack role. Specify `id`, `slug`, or `name` to identify the rack role.",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Full name of the rack role. Can be used to identify the rack role instead of `id` or `slug`.",
				Optional:            true,
				Computed:            true,
			},
			"slug": schema.StringAttribute{
				MarkdownDescription: "URL-friendly identifier for the rack role. Specify `id`, `slug`, or `name` to identify the rack role.",
				Optional:            true,
				Computed:            true,
			},
			"color": schema.StringAttribute{
				MarkdownDescription: "Color for the rack role in 6-character hexadecimal format (e.g., 'aa1409').",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Detailed description of the rack role.",
				Computed:            true,
			},
			"tags": schema.SetNestedAttribute{
				MarkdownDescription: "Tags assigned to this rack role.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "Name of the tag.",
							Computed:            true,
						},
						"slug": schema.StringAttribute{
							MarkdownDescription: "Slug of the tag.",
							Computed:            true,
						},
					},
				},
			},
			"custom_fields": schema.SetNestedAttribute{
				MarkdownDescription: "Custom fields assigned to this rack role.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "Name of the custom field.",
							Computed:            true,
						},
						"type": schema.StringAttribute{
							MarkdownDescription: "Type of the custom field.",
							Computed:            true,
						},
						"value": schema.StringAttribute{
							MarkdownDescription: "Value of the custom field.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *RackRoleDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *RackRoleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data RackRoleDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var rackRole *netbox.RackRole
	var err error
	var httpResp *http.Response

	// Determine if we're searching by ID, slug, or name
	if !data.ID.IsNull() {
		// Search by ID
		rackRoleID := data.ID.ValueString()
		tflog.Debug(ctx, "Reading rack role by ID", map[string]interface{}{
			"id": rackRoleID,
		})

		// Parse the rack role ID to int32 for the API call
		var rackRoleIDInt int32
		if _, parseErr := fmt.Sscanf(rackRoleID, "%d", &rackRoleIDInt); parseErr != nil {
			resp.Diagnostics.AddError(
				"Invalid Rack Role ID",
				fmt.Sprintf("Rack Role ID must be a number, got: %s", rackRoleID),
			)
			return
		}

		// Retrieve the rack role via API
		rackRole, httpResp, err = d.client.DcimAPI.DcimRackRolesRetrieve(ctx, rackRoleIDInt).Execute()
	} else if !data.Slug.IsNull() {
		// Search by slug
		rackRoleSlug := data.Slug.ValueString()
		tflog.Debug(ctx, "Reading rack role by slug", map[string]interface{}{
			"slug": rackRoleSlug,
		})

		// List rack roles with slug filter
		var rackRoles *netbox.PaginatedRackRoleList
		rackRoles, httpResp, err = d.client.DcimAPI.DcimRackRolesList(ctx).Slug([]string{rackRoleSlug}).Execute()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading rack role",
				utils.FormatAPIError("read rack role by slug", err, httpResp),
			)
			return
		}
		if len(rackRoles.GetResults()) == 0 {
			resp.Diagnostics.AddError(
				"Rack Role Not Found",
				fmt.Sprintf("No rack role found with slug: %s", rackRoleSlug),
			)
			return
		}
		if len(rackRoles.GetResults()) > 1 {
			resp.Diagnostics.AddError(
				"Multiple Rack Roles Found",
				fmt.Sprintf("Multiple rack roles found with slug: %s. This should not happen as slugs should be unique.", rackRoleSlug),
			)
			return
		}
		rackRole = &rackRoles.GetResults()[0]
	} else if !data.Name.IsNull() {
		// Search by name
		rackRoleName := data.Name.ValueString()
		tflog.Debug(ctx, "Reading rack role by name", map[string]interface{}{
			"name": rackRoleName,
		})

		// List rack roles with name filter
		var rackRoles *netbox.PaginatedRackRoleList
		rackRoles, httpResp, err = d.client.DcimAPI.DcimRackRolesList(ctx).Name([]string{rackRoleName}).Execute()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading rack role",
				utils.FormatAPIError("read rack role by name", err, httpResp),
			)
			return
		}
		if len(rackRoles.GetResults()) == 0 {
			resp.Diagnostics.AddError(
				"Rack Role Not Found",
				fmt.Sprintf("No rack role found with name: %s", rackRoleName),
			)
			return
		}
		if len(rackRoles.GetResults()) > 1 {
			resp.Diagnostics.AddError(
				"Multiple Rack Roles Found",
				fmt.Sprintf("Multiple rack roles found with name: %s. Rack role names may not be unique in Netbox.", rackRoleName),
			)
			return
		}
		rackRole = &rackRoles.GetResults()[0]
	} else {
		resp.Diagnostics.AddError(
			"Missing Rack Role Identifier",
			"Either 'id', 'slug', or 'name' must be specified to identify the rack role.",
		)
		return
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading rack role",
			utils.FormatAPIError("read rack role", err, httpResp),
		)
		return
	}

	if httpResp.StatusCode == 404 {
		resp.Diagnostics.AddError(
			"Rack Role Not Found",
			"The specified rack role was not found in Netbox.",
		)
		return
	}

	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Error reading rack role",
			fmt.Sprintf("Expected HTTP 200, got: %d", httpResp.StatusCode),
		)
		return
	}

	// Update the model with the response from the API
	data.ID = types.StringValue(fmt.Sprintf("%d", rackRole.GetId()))
	data.Name = types.StringValue(rackRole.GetName())
	data.Slug = types.StringValue(rackRole.GetSlug())

	// Handle color
	if rackRole.HasColor() && rackRole.GetColor() != "" {
		data.Color = types.StringValue(rackRole.GetColor())
	} else {
		data.Color = types.StringNull()
	}

	// Handle description
	if rackRole.HasDescription() && rackRole.GetDescription() != "" {
		data.Description = types.StringValue(rackRole.GetDescription())
	} else {
		data.Description = types.StringNull()
	}

	// Handle tags
	if rackRole.HasTags() {
		tags := utils.NestedTagsToTagModels(rackRole.GetTags())
		tagsValue, diags := types.SetValueFrom(ctx, utils.GetTagsAttributeType().ElemType, tags)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		data.Tags = tagsValue
	} else {
		data.Tags = types.SetNull(utils.GetTagsAttributeType().ElemType)
	}

	// Handle custom fields
	if rackRole.HasCustomFields() {
		// For data sources, we extract all available custom fields
		customFields := utils.MapToCustomFieldModels(rackRole.GetCustomFields(), nil)
		customFieldsValue, diags := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, customFields)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		data.CustomFields = customFieldsValue
	} else {
		data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
