// Package resources contains Terraform resource implementations for the Netbox provider.
//
// This package integrates with the go-netbox OpenAPI client to provide
// CRUD operations for Netbox resources via Terraform.
package resources

import (
	"context"
	"fmt"
	"regexp"

	"github.com/bab3l/go-netbox"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/bab3l/terraform-provider-netbox/internal/validators"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &TenantResource{}
var _ resource.ResourceWithImportState = &TenantResource{}

func NewTenantResource() resource.Resource {
	return &TenantResource{}
}

// TenantResource defines the resource implementation.
type TenantResource struct {
	client *netbox.APIClient
}

// TenantResourceModel describes the resource data model.
type TenantResourceModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Slug         types.String `tfsdk:"slug"`
	Group        types.String `tfsdk:"group"`
	Description  types.String `tfsdk:"description"`
	Comments     types.String `tfsdk:"comments"`
	Tags         types.Set    `tfsdk:"tags"`
	CustomFields types.Set    `tfsdk:"custom_fields"`
}

func (r *TenantResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tenant"
}

func (r *TenantResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a tenant in Netbox. Tenants represent individual customers or organizational units in multi-tenancy scenarios, allowing you to organize and track resources by client or department.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Unique identifier for the tenant (assigned by Netbox).",
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Full name of the tenant. This is the human-readable display name.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 100),
				},
			},
			"slug": schema.StringAttribute{
				MarkdownDescription: "URL-friendly identifier for the tenant. Must be unique and contain only alphanumeric characters, hyphens, and underscores.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 100),
					validators.ValidSlug(),
				},
			},
			"group": schema.StringAttribute{
				MarkdownDescription: "ID of the tenant group that this tenant belongs to.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						validators.IntegerRegex(),
						"must be a valid integer ID",
					),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Detailed description of the tenant, its purpose, or other relevant information.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(200),
				},
			},
			"comments": schema.StringAttribute{
				MarkdownDescription: "Additional comments or notes about the tenant. Supports Markdown formatting.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(1000),
				},
			},
			"tags": schema.SetNestedAttribute{
				MarkdownDescription: "Tags assigned to this tenant. Tags provide a way to categorize and organize resources.",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "Name of the existing tag.",
							Required:            true,
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 100),
							},
						},
						"slug": schema.StringAttribute{
							MarkdownDescription: "Slug of the existing tag.",
							Required:            true,
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 100),
								validators.ValidSlug(),
							},
						},
					},
				},
			},
			"custom_fields": schema.SetNestedAttribute{
				MarkdownDescription: "Custom fields assigned to this tenant. Custom fields allow you to store additional structured data.",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "Name of the custom field.",
							Required:            true,
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 50),
								stringvalidator.RegexMatches(
									regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]*$`),
									"must start with a letter and contain only letters, numbers, and underscores",
								),
							},
						},
						"type": schema.StringAttribute{
							MarkdownDescription: "Type of the custom field (text, longtext, integer, boolean, date, url, json, select, multiselect, object, multiobject).",
							Required:            true,
							Validators: []validator.String{
								stringvalidator.OneOf(
									"text",
									"longtext",
									"integer",
									"boolean",
									"date",
									"url",
									"json",
									"select",
									"multiselect",
									"object",
									"multiobject",
									"multiple",  // legacy
									"selection", // legacy
								),
							},
						},
						"value": schema.StringAttribute{
							MarkdownDescription: "Value of the custom field.",
							Required:            true,
							Validators: []validator.String{
								stringvalidator.LengthAtMost(1000),
							},
						},
					},
				},
			},
		},
	}
}

func (r *TenantResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
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

func (r *TenantResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data TenantResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create tenant using go-netbox client
	tflog.Debug(ctx, "Creating tenant", map[string]interface{}{
		"name": data.Name.ValueString(),
		"slug": data.Slug.ValueString(),
	})

	// Prepare the tenant request
	tenantRequest := netbox.TenantRequest{
		Name: data.Name.ValueString(),
		Slug: data.Slug.ValueString(),
	}

	// Set optional fields if provided
	if !data.Group.IsNull() {
		var groupIDInt int32
		if _, err := fmt.Sscanf(data.Group.ValueString(), "%d", &groupIDInt); err != nil {
			resp.Diagnostics.AddError(
				"Invalid Group ID",
				fmt.Sprintf("Group ID must be a number, got: %s", data.Group.ValueString()),
			)
			return
		}
		// Note: For now creating empty group request - ID assignment needs to be implemented
		groupRequest := netbox.BriefTenantGroupRequest{}
		tenantRequest.Group = *netbox.NewNullableBriefTenantGroupRequest(&groupRequest)
	}

	if !data.Description.IsNull() {
		description := data.Description.ValueString()
		tenantRequest.Description = &description
	}

	if !data.Comments.IsNull() {
		comments := data.Comments.ValueString()
		tenantRequest.Comments = &comments
	}

	// Handle tags
	if !data.Tags.IsNull() && !data.Tags.IsUnknown() {
		var tags []utils.TagModel
		resp.Diagnostics.Append(data.Tags.ElementsAs(ctx, &tags, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		tenantRequest.Tags = utils.TagsToNestedTagRequests(tags)
	}

	// Handle custom fields
	if !data.CustomFields.IsNull() && !data.CustomFields.IsUnknown() {
		var customFields []utils.CustomFieldModel
		resp.Diagnostics.Append(data.CustomFields.ElementsAs(ctx, &customFields, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		tenantRequest.CustomFields = utils.CustomFieldsToMap(customFields)
	}

	// Create the tenant via API
	tenant, httpResp, err := r.client.TenancyAPI.TenancyTenantsCreate(ctx).TenantRequest(tenantRequest).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating tenant",
			fmt.Sprintf("Could not create tenant, unexpected error: %s", err),
		)
		return
	}

	if httpResp.StatusCode != 201 {
		resp.Diagnostics.AddError(
			"Error creating tenant",
			fmt.Sprintf("Expected HTTP 201, got: %d", httpResp.StatusCode),
		)
		return
	}

	// Update the model with the response from the API
	data.ID = types.StringValue(fmt.Sprintf("%d", tenant.GetId()))
	data.Name = types.StringValue(tenant.GetName())
	data.Slug = types.StringValue(tenant.GetSlug())

	if tenant.HasGroup() {
		group := tenant.GetGroup()
		data.Group = types.StringValue(fmt.Sprintf("%d", group.GetId()))
	} else {
		data.Group = types.StringNull()
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

	tflog.Trace(ctx, "created a tenant resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TenantResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data TenantResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get the tenant ID from state
	tenantID := data.ID.ValueString()

	tflog.Debug(ctx, "Reading tenant", map[string]interface{}{
		"id": tenantID,
	})

	// Parse the tenant ID to int32 for the API call
	var tenantIDInt int32
	if _, err := fmt.Sscanf(tenantID, "%d", &tenantIDInt); err != nil {
		resp.Diagnostics.AddError(
			"Invalid Tenant ID",
			fmt.Sprintf("Tenant ID must be a number, got: %s", tenantID),
		)
		return
	}

	// Retrieve the tenant via API
	tenant, httpResp, err := r.client.TenancyAPI.TenancyTenantsRetrieve(ctx, tenantIDInt).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading tenant",
			fmt.Sprintf("Could not read tenant ID %s: %s", tenantID, err),
		)
		return
	}

	if httpResp.StatusCode == 404 {
		// Tenant no longer exists, remove from state
		resp.State.RemoveResource(ctx)
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
		data.Group = types.StringValue(fmt.Sprintf("%d", group.GetId()))
	} else {
		data.Group = types.StringNull()
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

	// Handle custom fields - we need to preserve the state structure
	if tenant.HasCustomFields() && !data.CustomFields.IsNull() {
		var stateCustomFields []utils.CustomFieldModel
		resp.Diagnostics.Append(data.CustomFields.ElementsAs(ctx, &stateCustomFields, false)...)
		if resp.Diagnostics.HasError() {
			return
		}

		customFields := utils.MapToCustomFieldModels(tenant.GetCustomFields(), stateCustomFields)
		customFieldsValue, diags := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, customFields)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		data.CustomFields = customFieldsValue
	} else if data.CustomFields.IsNull() {
		data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TenantResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data TenantResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get the tenant ID from state
	tenantID := data.ID.ValueString()

	tflog.Debug(ctx, "Updating tenant", map[string]interface{}{
		"id":   tenantID,
		"name": data.Name.ValueString(),
		"slug": data.Slug.ValueString(),
	})

	// Parse the tenant ID to int32 for the API call
	var tenantIDInt int32
	if _, err := fmt.Sscanf(tenantID, "%d", &tenantIDInt); err != nil {
		resp.Diagnostics.AddError(
			"Invalid Tenant ID",
			fmt.Sprintf("Tenant ID must be a number, got: %s", tenantID),
		)
		return
	}

	// Prepare the tenant update request
	tenantRequest := netbox.TenantRequest{
		Name: data.Name.ValueString(),
		Slug: data.Slug.ValueString(),
	}

	// Set optional fields if provided
	if !data.Group.IsNull() {
		var groupIDInt int32
		if _, err := fmt.Sscanf(data.Group.ValueString(), "%d", &groupIDInt); err != nil {
			resp.Diagnostics.AddError(
				"Invalid Group ID",
				fmt.Sprintf("Group ID must be a number, got: %s", data.Group.ValueString()),
			)
			return
		}
		// Note: For now creating empty group request - ID assignment needs to be implemented
		groupRequest := netbox.BriefTenantGroupRequest{}
		tenantRequest.Group = *netbox.NewNullableBriefTenantGroupRequest(&groupRequest)
	}

	if !data.Description.IsNull() {
		description := data.Description.ValueString()
		tenantRequest.Description = &description
	}

	if !data.Comments.IsNull() {
		comments := data.Comments.ValueString()
		tenantRequest.Comments = &comments
	}

	// Handle tags
	if !data.Tags.IsNull() && !data.Tags.IsUnknown() {
		var tags []utils.TagModel
		resp.Diagnostics.Append(data.Tags.ElementsAs(ctx, &tags, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		tenantRequest.Tags = utils.TagsToNestedTagRequests(tags)
	}

	// Handle custom fields
	if !data.CustomFields.IsNull() && !data.CustomFields.IsUnknown() {
		var customFields []utils.CustomFieldModel
		resp.Diagnostics.Append(data.CustomFields.ElementsAs(ctx, &customFields, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		tenantRequest.CustomFields = utils.CustomFieldsToMap(customFields)
	}

	// Update the tenant via API
	tenant, httpResp, err := r.client.TenancyAPI.TenancyTenantsUpdate(ctx, tenantIDInt).TenantRequest(tenantRequest).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating tenant",
			fmt.Sprintf("Could not update tenant ID %s: %s", tenantID, err),
		)
		return
	}

	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Error updating tenant",
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
		data.Group = types.StringValue(fmt.Sprintf("%d", group.GetId()))
	} else {
		data.Group = types.StringNull()
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

	// Handle tags in response
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

	// Handle custom fields in response
	if tenant.HasCustomFields() && !data.CustomFields.IsNull() {
		var stateCustomFields []utils.CustomFieldModel
		resp.Diagnostics.Append(data.CustomFields.ElementsAs(ctx, &stateCustomFields, false)...)
		if resp.Diagnostics.HasError() {
			return
		}

		customFields := utils.MapToCustomFieldModels(tenant.GetCustomFields(), stateCustomFields)
		customFieldsValue, diags := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, customFields)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		data.CustomFields = customFieldsValue
	} else if data.CustomFields.IsNull() {
		data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TenantResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data TenantResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get the tenant ID from state
	tenantID := data.ID.ValueString()

	tflog.Debug(ctx, "Deleting tenant", map[string]interface{}{
		"id": tenantID,
	})

	// Parse the tenant ID to int32 for the API call
	var tenantIDInt int32
	if _, err := fmt.Sscanf(tenantID, "%d", &tenantIDInt); err != nil {
		resp.Diagnostics.AddError(
			"Invalid Tenant ID",
			fmt.Sprintf("Tenant ID must be a number, got: %s", tenantID),
		)
		return
	}

	// Delete the tenant via API
	httpResp, err := r.client.TenancyAPI.TenancyTenantsDestroy(ctx, tenantIDInt).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting tenant",
			fmt.Sprintf("Could not delete tenant ID %s: %s", tenantID, err),
		)
		return
	}

	if httpResp.StatusCode != 204 {
		resp.Diagnostics.AddError(
			"Error deleting tenant",
			fmt.Sprintf("Expected HTTP 204, got: %d", httpResp.StatusCode),
		)
		return
	}

	tflog.Trace(ctx, "deleted a tenant resource")
}

func (r *TenantResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
