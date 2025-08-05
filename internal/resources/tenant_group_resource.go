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
var _ resource.Resource = &TenantGroupResource{}
var _ resource.ResourceWithImportState = &TenantGroupResource{}

func NewTenantGroupResource() resource.Resource {
	return &TenantGroupResource{}
}

// TenantGroupResource defines the resource implementation.
type TenantGroupResource struct {
	client *netbox.APIClient
}

// TenantGroupResourceModel describes the resource data model.
type TenantGroupResourceModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Slug         types.String `tfsdk:"slug"`
	Parent       types.String `tfsdk:"parent"`
	Description  types.String `tfsdk:"description"`
	Tags         types.Set    `tfsdk:"tags"`
	CustomFields types.Set    `tfsdk:"custom_fields"`
}

func (r *TenantGroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tenant_group"
}

func (r *TenantGroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a tenant group in Netbox. Tenant groups provide a hierarchical way to organize tenants for multi-tenancy scenarios, allowing you to create nested organizational structures for better management and reporting.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Unique identifier for the tenant group (assigned by Netbox).",
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Full name of the tenant group. This is the human-readable display name.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 100),
				},
			},
			"slug": schema.StringAttribute{
				MarkdownDescription: "URL-friendly identifier for the tenant group. Must be unique and contain only alphanumeric characters, hyphens, and underscores.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 100),
					validators.ValidSlug(),
				},
			},
			"parent": schema.StringAttribute{
				MarkdownDescription: "ID of the parent tenant group. Leave empty for top-level groups.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						validators.IntegerRegex(),
						"must be a valid integer ID",
					),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Detailed description of the tenant group, its purpose, or other relevant information.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(200),
				},
			},
			"tags": schema.SetNestedAttribute{
				MarkdownDescription: "Tags assigned to this tenant group. Tags provide a way to categorize and organize resources.",
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
				MarkdownDescription: "Custom fields assigned to this tenant group. Custom fields allow you to store additional structured data.",
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

func (r *TenantGroupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *TenantGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data TenantGroupResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create tenant group using go-netbox client
	tflog.Debug(ctx, "Creating tenant group", map[string]interface{}{
		"name": data.Name.ValueString(),
		"slug": data.Slug.ValueString(),
	})

	// Prepare the tenant group request
	tenantGroupRequest := netbox.WritableTenantGroupRequest{
		Name: data.Name.ValueString(),
		Slug: data.Slug.ValueString(),
	}

	// Set optional fields if provided
	if !data.Parent.IsNull() {
		var parentIDInt int32
		if _, err := fmt.Sscanf(data.Parent.ValueString(), "%d", &parentIDInt); err != nil {
			resp.Diagnostics.AddError(
				"Invalid Parent ID",
				fmt.Sprintf("Parent ID must be a number, got: %s", data.Parent.ValueString()),
			)
			return
		}
		parentID := int32(parentIDInt)
		tenantGroupRequest.Parent = *netbox.NewNullableInt32(&parentID)
	}

	if !data.Description.IsNull() {
		description := data.Description.ValueString()
		tenantGroupRequest.Description = &description
	}

	// Handle tags
	if !data.Tags.IsNull() && !data.Tags.IsUnknown() {
		var tags []utils.TagModel
		resp.Diagnostics.Append(data.Tags.ElementsAs(ctx, &tags, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		tenantGroupRequest.Tags = utils.TagsToNestedTagRequests(tags)
	}

	// Handle custom fields
	if !data.CustomFields.IsNull() && !data.CustomFields.IsUnknown() {
		var customFields []utils.CustomFieldModel
		resp.Diagnostics.Append(data.CustomFields.ElementsAs(ctx, &customFields, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		tenantGroupRequest.CustomFields = utils.CustomFieldsToMap(customFields)
	}

	// Create the tenant group via API
	tenantGroup, httpResp, err := r.client.TenancyAPI.TenancyTenantGroupsCreate(ctx).WritableTenantGroupRequest(tenantGroupRequest).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating tenant group",
			fmt.Sprintf("Could not create tenant group, unexpected error: %s", err),
		)
		return
	}

	if httpResp.StatusCode != 201 {
		resp.Diagnostics.AddError(
			"Error creating tenant group",
			fmt.Sprintf("Expected HTTP 201, got: %d", httpResp.StatusCode),
		)
		return
	}

	// Update the model with the response from the API
	data.ID = types.StringValue(fmt.Sprintf("%d", tenantGroup.GetId()))
	data.Name = types.StringValue(tenantGroup.GetName())
	data.Slug = types.StringValue(tenantGroup.GetSlug())

	if tenantGroup.HasParent() {
		parent := tenantGroup.GetParent()
		data.Parent = types.StringValue(fmt.Sprintf("%d", parent.GetId()))
	} else {
		data.Parent = types.StringNull()
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

	tflog.Trace(ctx, "created a tenant group resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TenantGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data TenantGroupResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get the tenant group ID from state
	tenantGroupID := data.ID.ValueString()

	tflog.Debug(ctx, "Reading tenant group", map[string]interface{}{
		"id": tenantGroupID,
	})

	// Parse the tenant group ID to int32 for the API call
	var tenantGroupIDInt int32
	if _, err := fmt.Sscanf(tenantGroupID, "%d", &tenantGroupIDInt); err != nil {
		resp.Diagnostics.AddError(
			"Invalid Tenant Group ID",
			fmt.Sprintf("Tenant Group ID must be a number, got: %s", tenantGroupID),
		)
		return
	}

	// Retrieve the tenant group via API
	tenantGroup, httpResp, err := r.client.TenancyAPI.TenancyTenantGroupsRetrieve(ctx, tenantGroupIDInt).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading tenant group",
			fmt.Sprintf("Could not read tenant group ID %s: %s", tenantGroupID, err),
		)
		return
	}

	if httpResp.StatusCode == 404 {
		// Tenant group no longer exists, remove from state
		resp.State.RemoveResource(ctx)
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
		data.Parent = types.StringValue(fmt.Sprintf("%d", parent.GetId()))
	} else {
		data.Parent = types.StringNull()
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

	// Handle custom fields - we need to preserve the state structure
	if tenantGroup.HasCustomFields() && !data.CustomFields.IsNull() {
		var stateCustomFields []utils.CustomFieldModel
		resp.Diagnostics.Append(data.CustomFields.ElementsAs(ctx, &stateCustomFields, false)...)
		if resp.Diagnostics.HasError() {
			return
		}

		customFields := utils.MapToCustomFieldModels(tenantGroup.GetCustomFields(), stateCustomFields)
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

func (r *TenantGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data TenantGroupResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get the tenant group ID from state
	tenantGroupID := data.ID.ValueString()

	tflog.Debug(ctx, "Updating tenant group", map[string]interface{}{
		"id":   tenantGroupID,
		"name": data.Name.ValueString(),
		"slug": data.Slug.ValueString(),
	})

	// Parse the tenant group ID to int32 for the API call
	var tenantGroupIDInt int32
	if _, err := fmt.Sscanf(tenantGroupID, "%d", &tenantGroupIDInt); err != nil {
		resp.Diagnostics.AddError(
			"Invalid Tenant Group ID",
			fmt.Sprintf("Tenant Group ID must be a number, got: %s", tenantGroupID),
		)
		return
	}

	// Prepare the tenant group update request
	tenantGroupRequest := netbox.WritableTenantGroupRequest{
		Name: data.Name.ValueString(),
		Slug: data.Slug.ValueString(),
	}

	// Set optional fields if provided
	if !data.Parent.IsNull() {
		var parentIDInt int32
		if _, err := fmt.Sscanf(data.Parent.ValueString(), "%d", &parentIDInt); err != nil {
			resp.Diagnostics.AddError(
				"Invalid Parent ID",
				fmt.Sprintf("Parent ID must be a number, got: %s", data.Parent.ValueString()),
			)
			return
		}
		parentID := int32(parentIDInt)
		tenantGroupRequest.Parent = *netbox.NewNullableInt32(&parentID)
	}

	if !data.Description.IsNull() {
		description := data.Description.ValueString()
		tenantGroupRequest.Description = &description
	}

	// Handle tags
	if !data.Tags.IsNull() && !data.Tags.IsUnknown() {
		var tags []utils.TagModel
		resp.Diagnostics.Append(data.Tags.ElementsAs(ctx, &tags, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		tenantGroupRequest.Tags = utils.TagsToNestedTagRequests(tags)
	}

	// Handle custom fields
	if !data.CustomFields.IsNull() && !data.CustomFields.IsUnknown() {
		var customFields []utils.CustomFieldModel
		resp.Diagnostics.Append(data.CustomFields.ElementsAs(ctx, &customFields, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		tenantGroupRequest.CustomFields = utils.CustomFieldsToMap(customFields)
	}

	// Update the tenant group via API
	tenantGroup, httpResp, err := r.client.TenancyAPI.TenancyTenantGroupsUpdate(ctx, tenantGroupIDInt).WritableTenantGroupRequest(tenantGroupRequest).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating tenant group",
			fmt.Sprintf("Could not update tenant group ID %s: %s", tenantGroupID, err),
		)
		return
	}

	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Error updating tenant group",
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
		data.Parent = types.StringValue(fmt.Sprintf("%d", parent.GetId()))
	} else {
		data.Parent = types.StringNull()
	}

	if tenantGroup.HasDescription() {
		data.Description = types.StringValue(tenantGroup.GetDescription())
	} else {
		data.Description = types.StringNull()
	}

	// Handle tags in response
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

	// Handle custom fields in response
	if tenantGroup.HasCustomFields() && !data.CustomFields.IsNull() {
		var stateCustomFields []utils.CustomFieldModel
		resp.Diagnostics.Append(data.CustomFields.ElementsAs(ctx, &stateCustomFields, false)...)
		if resp.Diagnostics.HasError() {
			return
		}

		customFields := utils.MapToCustomFieldModels(tenantGroup.GetCustomFields(), stateCustomFields)
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

func (r *TenantGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data TenantGroupResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get the tenant group ID from state
	tenantGroupID := data.ID.ValueString()

	tflog.Debug(ctx, "Deleting tenant group", map[string]interface{}{
		"id": tenantGroupID,
	})

	// Parse the tenant group ID to int32 for the API call
	var tenantGroupIDInt int32
	if _, err := fmt.Sscanf(tenantGroupID, "%d", &tenantGroupIDInt); err != nil {
		resp.Diagnostics.AddError(
			"Invalid Tenant Group ID",
			fmt.Sprintf("Tenant Group ID must be a number, got: %s", tenantGroupID),
		)
		return
	}

	// Delete the tenant group via API
	httpResp, err := r.client.TenancyAPI.TenancyTenantGroupsDestroy(ctx, tenantGroupIDInt).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting tenant group",
			fmt.Sprintf("Could not delete tenant group ID %s: %s", tenantGroupID, err),
		)
		return
	}

	if httpResp.StatusCode != 204 {
		resp.Diagnostics.AddError(
			"Error deleting tenant group",
			fmt.Sprintf("Expected HTTP 204, got: %d", httpResp.StatusCode),
		)
		return
	}

	tflog.Trace(ctx, "deleted a tenant group resource")
}

func (r *TenantGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
