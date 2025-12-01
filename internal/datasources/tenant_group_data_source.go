// Package datasources contains Terraform data source implementations for the Netbox provider.
//
// This package provides read-only access to Netbox resources for use in Terraform configurations.
package datasources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bab3l/go-netbox"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/bab3l/terraform-provider-netbox/internal/utils"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &TenantGroupDataSource{}

func NewTenantGroupDataSource() datasource.DataSource {
	return &TenantGroupDataSource{}
}

// TenantGroupDataSource defines the data source implementation.
type TenantGroupDataSource struct {
	client *netbox.APIClient
}

// TenantGroupDataSourceModel describes the data source data model.
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
		MarkdownDescription: "Use this data source to get information about a tenant group in Netbox. Tenant groups provide hierarchical organization of tenants for multi-tenancy scenarios. You can identify the tenant group using `id`, `slug`, or `name`.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier for the tenant group. Specify `id`, `slug`, or `name` to identify the tenant group.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 50),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Full name of the tenant group. Can be used to identify the tenant group instead of `id` or `slug`.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 100),
				},
			},
			"slug": schema.StringAttribute{
				MarkdownDescription: "URL-friendly identifier for the tenant group. Specify `id`, `slug`, or `name` to identify the tenant group.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 100),
				},
			},
			"parent": schema.StringAttribute{
				MarkdownDescription: "Name of the parent tenant group. Null if this is a top-level group.",
				Computed:            true,
			},
			"parent_id": schema.StringAttribute{
				MarkdownDescription: "ID of the parent tenant group. Null if this is a top-level group.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Detailed description of the tenant group.",
				Computed:            true,
			},
			"tags": schema.SetNestedAttribute{
				MarkdownDescription: "Tags assigned to this tenant group.",
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
				MarkdownDescription: "Custom fields assigned to this tenant group.",
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

func (d *TenantGroupDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *TenantGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data TenantGroupDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var tenantGroup *netbox.TenantGroup
	var err error
	var httpResp *http.Response

	// Determine if we're searching by ID, slug, or name
	if !data.ID.IsNull() {
		// Search by ID
		tenantGroupID := data.ID.ValueString()
		tflog.Debug(ctx, "Reading tenant group by ID", map[string]interface{}{
			"id": tenantGroupID,
		})

		// Parse the tenant group ID to int32 for the API call
		var tenantGroupIDInt int32
		if _, parseErr := fmt.Sscanf(tenantGroupID, "%d", &tenantGroupIDInt); parseErr != nil {
			resp.Diagnostics.AddError(
				"Invalid Tenant Group ID",
				fmt.Sprintf("Tenant Group ID must be a number, got: %s", tenantGroupID),
			)
			return
		}

		// Retrieve the tenant group via API
		tenantGroup, httpResp, err = d.client.TenancyAPI.TenancyTenantGroupsRetrieve(ctx, tenantGroupIDInt).Execute()
	} else if !data.Slug.IsNull() {
		// Search by slug
		tenantGroupSlug := data.Slug.ValueString()
		tflog.Debug(ctx, "Reading tenant group by slug", map[string]interface{}{
			"slug": tenantGroupSlug,
		})

		// List tenant groups with slug filter
		var tenantGroups *netbox.PaginatedTenantGroupList
		tenantGroups, httpResp, err = d.client.TenancyAPI.TenancyTenantGroupsList(ctx).Slug([]string{tenantGroupSlug}).Execute()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading tenant group",
				utils.FormatAPIError("read tenant group by slug", err, httpResp),
			)
			return
		}
		if len(tenantGroups.GetResults()) == 0 {
			resp.Diagnostics.AddError(
				"Tenant Group Not Found",
				fmt.Sprintf("No tenant group found with slug: %s", tenantGroupSlug),
			)
			return
		}
		if len(tenantGroups.GetResults()) > 1 {
			resp.Diagnostics.AddError(
				"Multiple Tenant Groups Found",
				fmt.Sprintf("Multiple tenant groups found with slug: %s. This should not happen as slugs should be unique.", tenantGroupSlug),
			)
			return
		}
		tenantGroup = &tenantGroups.GetResults()[0]
	} else if !data.Name.IsNull() {
		// Search by name
		tenantGroupName := data.Name.ValueString()
		tflog.Debug(ctx, "Reading tenant group by name", map[string]interface{}{
			"name": tenantGroupName,
		})

		// List tenant groups with name filter
		var tenantGroups *netbox.PaginatedTenantGroupList
		tenantGroups, httpResp, err = d.client.TenancyAPI.TenancyTenantGroupsList(ctx).Name([]string{tenantGroupName}).Execute()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading tenant group",
				utils.FormatAPIError("read tenant group by name", err, httpResp),
			)
			return
		}
		if len(tenantGroups.GetResults()) == 0 {
			resp.Diagnostics.AddError(
				"Tenant Group Not Found",
				fmt.Sprintf("No tenant group found with name: %s", tenantGroupName),
			)
			return
		}
		if len(tenantGroups.GetResults()) > 1 {
			resp.Diagnostics.AddError(
				"Multiple Tenant Groups Found",
				fmt.Sprintf("Multiple tenant groups found with name: %s. Tenant group names may not be unique in Netbox.", tenantGroupName),
			)
			return
		}
		tenantGroup = &tenantGroups.GetResults()[0]
	} else {
		resp.Diagnostics.AddError(
			"Missing Tenant Group Identifier",
			"Either 'id', 'slug', or 'name' must be specified to identify the tenant group.",
		)
		return
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading tenant group",
			utils.FormatAPIError("read tenant group", err, httpResp),
		)
		return
	}

	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Error reading tenant group",
			fmt.Sprintf("Expected HTTP 200, got: %d", httpResp.StatusCode),
		)
		return
	}

	// Update the model with the response from the API
	data.ID = types.StringValue(fmt.Sprintf("%d", tenantGroup.GetId()))
	data.Name = types.StringValue(tenantGroup.GetName())
	data.Slug = types.StringValue(tenantGroup.GetSlug())

	if tenantGroup.HasParent() {
		parent := tenantGroup.GetParent()
		// Set parent to the ID for consistency with the resource
		data.Parent = types.StringValue(fmt.Sprintf("%d", parent.GetId()))
		data.ParentID = types.StringValue(fmt.Sprintf("%d", parent.GetId()))
	} else {
		data.Parent = types.StringNull()
		data.ParentID = types.StringNull()
	}

	if tenantGroup.HasDescription() {
		data.Description = types.StringValue(tenantGroup.GetDescription())
	} else {
		data.Description = types.StringNull()
	}

	// Handle tags
	if tenantGroup.HasTags() {
		tags := utils.NestedTagsToTagModels(tenantGroup.GetTags())
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
	if tenantGroup.HasCustomFields() {
		customFields := utils.MapToCustomFieldModels(tenantGroup.GetCustomFields(), []utils.CustomFieldModel{})
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
