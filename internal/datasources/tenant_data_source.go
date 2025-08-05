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
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier for the tenant. Specify `id`, `slug`, or `name` to identify the tenant.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 50),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Full name of the tenant. Can be used to identify the tenant instead of `id` or `slug`.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 100),
				},
			},
			"slug": schema.StringAttribute{
				MarkdownDescription: "URL-friendly identifier for the tenant. Specify `id`, `slug`, or `name` to identify the tenant.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 100),
				},
			},
			"group": schema.StringAttribute{
				MarkdownDescription: "Name of the tenant group that this tenant belongs to.",
				Computed:            true,
			},
			"group_id": schema.StringAttribute{
				MarkdownDescription: "ID of the tenant group that this tenant belongs to.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Detailed description of the tenant.",
				Computed:            true,
			},
			"comments": schema.StringAttribute{
				MarkdownDescription: "Additional comments or notes about the tenant.",
				Computed:            true,
			},
			"tags": schema.SetNestedAttribute{
				MarkdownDescription: "Tags assigned to this tenant.",
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
				MarkdownDescription: "Custom fields assigned to this tenant.",
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

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var tenant *netbox.Tenant
	var err error
	var httpResp *http.Response

	// Determine if we're searching by ID, slug, or name
	if !data.ID.IsNull() {
		// Search by ID
		tenantID := data.ID.ValueString()
		tflog.Debug(ctx, "Reading tenant by ID", map[string]interface{}{
			"id": tenantID,
		})

		// Parse the tenant ID to int32 for the API call
		var tenantIDInt int32
		if _, parseErr := fmt.Sscanf(tenantID, "%d", &tenantIDInt); parseErr != nil {
			resp.Diagnostics.AddError(
				"Invalid Tenant ID",
				fmt.Sprintf("Tenant ID must be a number, got: %s", tenantID),
			)
			return
		}

		// Retrieve the tenant via API
		tenant, httpResp, err = d.client.TenancyAPI.TenancyTenantsRetrieve(ctx, tenantIDInt).Execute()
	} else if !data.Slug.IsNull() {
		// Search by slug
		tenantSlug := data.Slug.ValueString()
		tflog.Debug(ctx, "Reading tenant by slug", map[string]interface{}{
			"slug": tenantSlug,
		})

		// List tenants with slug filter
		tenants, httpResp, err := d.client.TenancyAPI.TenancyTenantsList(ctx).Slug([]string{tenantSlug}).Execute()
		if err == nil && httpResp.StatusCode == 200 {
			if len(tenants.GetResults()) == 0 {
				resp.Diagnostics.AddError(
					"Tenant Not Found",
					fmt.Sprintf("No tenant found with slug: %s", tenantSlug),
				)
				return
			}
			if len(tenants.GetResults()) > 1 {
				resp.Diagnostics.AddError(
					"Multiple Tenants Found",
					fmt.Sprintf("Multiple tenants found with slug: %s. This should not happen as slugs should be unique.", tenantSlug),
				)
				return
			}
			tenant = &tenants.GetResults()[0]
		}
	} else if !data.Name.IsNull() {
		// Search by name
		tenantName := data.Name.ValueString()
		tflog.Debug(ctx, "Reading tenant by name", map[string]interface{}{
			"name": tenantName,
		})

		// List tenants with name filter
		tenants, httpResp, err := d.client.TenancyAPI.TenancyTenantsList(ctx).Name([]string{tenantName}).Execute()
		if err == nil && httpResp.StatusCode == 200 {
			if len(tenants.GetResults()) == 0 {
				resp.Diagnostics.AddError(
					"Tenant Not Found",
					fmt.Sprintf("No tenant found with name: %s", tenantName),
				)
				return
			}
			if len(tenants.GetResults()) > 1 {
				resp.Diagnostics.AddError(
					"Multiple Tenants Found",
					fmt.Sprintf("Multiple tenants found with name: %s. Tenant names may not be unique in Netbox.", tenantName),
				)
				return
			}
			tenant = &tenants.GetResults()[0]
		}
	} else {
		resp.Diagnostics.AddError(
			"Missing Tenant Identifier",
			"Either 'id', 'slug', or 'name' must be specified to identify the tenant.",
		)
		return
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading tenant",
			fmt.Sprintf("Could not read tenant: %s", err),
		)
		return
	}

	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Error reading tenant",
			fmt.Sprintf("Expected HTTP 200, got: %d", httpResp.StatusCode),
		)
		return
	}

	// Update the model with the response from the API
	data.ID = types.StringValue(fmt.Sprintf("%d", tenant.GetId()))
	data.Name = types.StringValue(tenant.GetName())
	data.Slug = types.StringValue(tenant.GetSlug())

	if tenant.HasGroup() {
		group := tenant.GetGroup()
		data.Group = types.StringValue(group.GetName())
		data.GroupID = types.StringValue(fmt.Sprintf("%d", group.GetId()))
	} else {
		data.Group = types.StringNull()
		data.GroupID = types.StringNull()
	}

	if tenant.HasDescription() {
		data.Description = types.StringValue(tenant.GetDescription())
	} else {
		data.Description = types.StringNull()
	}

	if tenant.HasComments() {
		data.Comments = types.StringValue(tenant.GetComments())
	} else {
		data.Comments = types.StringNull()
	}

	// Handle tags
	if tenant.HasTags() {
		tags := utils.NestedTagsToTagModels(tenant.GetTags())
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
	if tenant.HasCustomFields() {
		customFields := utils.MapToCustomFieldModels(tenant.GetCustomFields(), []utils.CustomFieldModel{})
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
