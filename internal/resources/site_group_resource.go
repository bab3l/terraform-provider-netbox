// Package resources contains Terraform resource implementations for the Netbox provider.
//
// This package integrates with the go-netbox OpenAPI client to provide
// CRUD operations for Netbox resources via Terraform.
package resources

import (
	"context"
	"fmt"

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
var _ resource.Resource = &SiteGroupResource{}
var _ resource.ResourceWithImportState = &SiteGroupResource{}

func NewSiteGroupResource() resource.Resource {
	return &SiteGroupResource{}
}

// SiteGroupResource defines the resource implementation.
type SiteGroupResource struct {
	client *netbox.APIClient
}

// SiteGroupResourceModel describes the resource data model.
type SiteGroupResourceModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Slug         types.String `tfsdk:"slug"`
	Parent       types.String `tfsdk:"parent"`
	Description  types.String `tfsdk:"description"`
	Tags         types.Set    `tfsdk:"tags"`
	CustomFields types.Set    `tfsdk:"custom_fields"`
}

func (r *SiteGroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_site_group"
}

func (r *SiteGroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a site group in Netbox. Site groups provide a hierarchical way to organize sites, allowing you to create nested organizational structures for better management and reporting of your physical locations.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Unique identifier for the site group (assigned by Netbox).",
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Full name of the site group. This is the human-readable display name.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 100),
				},
			},
			"slug": schema.StringAttribute{
				MarkdownDescription: "URL-friendly identifier for the site group. Must be unique and contain only alphanumeric characters, hyphens, and underscores.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 100),
					validators.ValidSlug(),
				},
			},
			"parent": schema.StringAttribute{
				MarkdownDescription: "ID of the parent site group. Leave empty for top-level site groups. This enables hierarchical organization of site groups.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						validators.IntegerRegex(),
						"must be a valid integer ID",
					),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Detailed description of the site group, its purpose, or other relevant information.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(200),
				},
			},
			"tags": schema.SetNestedAttribute{
				MarkdownDescription: "Tags assigned to this site group. Tags provide a way to categorize and organize resources.",
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
				MarkdownDescription: "Custom fields assigned to this site group. Custom fields allow you to store additional structured data.",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "Name of the custom field.",
							Required:            true,
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 50),
								validators.ValidCustomFieldName(),
							},
						},
						"type": schema.StringAttribute{
							MarkdownDescription: "Type of the custom field (text, longtext, integer, boolean, date, url, json, select, multiselect, object, multiobject).",
							Required:            true,
							Validators: []validator.String{
								validators.ValidCustomFieldType(),
							},
						},
						"value": schema.StringAttribute{
							MarkdownDescription: "Value of the custom field.",
							Required:            true,
							Validators: []validator.String{
								stringvalidator.LengthAtMost(1000),
								validators.SimpleValidCustomFieldValue(),
							},
						},
					},
				},
			},
		},
	}
}

func (r *SiteGroupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *SiteGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data SiteGroupResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create site group using go-netbox client
	tflog.Debug(ctx, "Creating site group", map[string]interface{}{
		"name": data.Name.ValueString(),
		"slug": data.Slug.ValueString(),
	})

	// Prepare the site group request
	siteGroupRequest := netbox.WritableSiteGroupRequest{
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
		siteGroupRequest.Parent = *netbox.NewNullableInt32(&parentID)
	}

	if !data.Description.IsNull() {
		description := data.Description.ValueString()
		siteGroupRequest.Description = &description
	}

	// Handle tags
	if !data.Tags.IsNull() && !data.Tags.IsUnknown() {
		var tags []utils.TagModel
		resp.Diagnostics.Append(data.Tags.ElementsAs(ctx, &tags, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		siteGroupRequest.Tags = utils.TagsToNestedTagRequests(tags)
	}

	// Handle custom fields
	if !data.CustomFields.IsNull() && !data.CustomFields.IsUnknown() {
		var customFields []utils.CustomFieldModel
		resp.Diagnostics.Append(data.CustomFields.ElementsAs(ctx, &customFields, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		siteGroupRequest.CustomFields = utils.CustomFieldsToMap(customFields)
	}

	// Create the site group via API
	siteGroup, httpResp, err := r.client.DcimAPI.DcimSiteGroupsCreate(ctx).WritableSiteGroupRequest(siteGroupRequest).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating site group",
			utils.FormatAPIError("create site group", err, httpResp),
		)
		return
	}

	if httpResp.StatusCode != 201 {
		resp.Diagnostics.AddError(
			"Error creating site group",
			fmt.Sprintf("Expected HTTP 201, got: %d", httpResp.StatusCode),
		)
		return
	}

	// Update the model with the response from the API
	data.ID = types.StringValue(fmt.Sprintf("%d", siteGroup.GetId()))
	data.Name = types.StringValue(siteGroup.GetName())
	data.Slug = types.StringValue(siteGroup.GetSlug())

	// Handle parent - check both HasParent and that ID is non-zero
	if siteGroup.HasParent() {
		parent := siteGroup.GetParent()
		if parent.GetId() != 0 {
			data.Parent = types.StringValue(fmt.Sprintf("%d", parent.GetId()))
		} else {
			data.Parent = types.StringNull()
		}
	} else {
		data.Parent = types.StringNull()
	}

	if siteGroup.HasDescription() {
		data.Description = types.StringValue(siteGroup.GetDescription())
	} else {
		data.Description = types.StringNull()
	}

	// Handle tags in response
	if siteGroup.HasTags() {
		tags := utils.NestedTagsToTagModels(siteGroup.GetTags())
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
	if siteGroup.HasCustomFields() {
		var stateCustomFields []utils.CustomFieldModel
		if !data.CustomFields.IsNull() {
			resp.Diagnostics.Append(data.CustomFields.ElementsAs(ctx, &stateCustomFields, false)...)
			if resp.Diagnostics.HasError() {
				return
			}
		}

		customFields := utils.MapToCustomFieldModels(siteGroup.GetCustomFields(), stateCustomFields)
		customFieldsValue, diags := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, customFields)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		data.CustomFields = customFieldsValue
	} else {
		data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)
	}

	tflog.Trace(ctx, "created a site group resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SiteGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data SiteGroupResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get the site group ID from state
	siteGroupID := data.ID.ValueString()

	tflog.Debug(ctx, "Reading site group", map[string]interface{}{
		"id": siteGroupID,
	})

	// Parse the site group ID to int32 for the API call
	var siteGroupIDInt int32
	if _, err := fmt.Sscanf(siteGroupID, "%d", &siteGroupIDInt); err != nil {
		resp.Diagnostics.AddError(
			"Invalid Site Group ID",
			fmt.Sprintf("Site Group ID must be a number, got: %s", siteGroupID),
		)
		return
	}

	// Retrieve the site group via API
	siteGroup, httpResp, err := r.client.DcimAPI.DcimSiteGroupsRetrieve(ctx, siteGroupIDInt).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading site group",
			utils.FormatAPIError(fmt.Sprintf("read site group ID %s", siteGroupID), err, httpResp),
		)
		return
	}

	if httpResp.StatusCode == 404 {
		// Site group no longer exists, remove from state
		resp.State.RemoveResource(ctx)
		return
	}

	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Error reading site group",
			fmt.Sprintf("Expected HTTP 200, got: %d", httpResp.StatusCode),
		)
		return
	}

	// Update the model with the response from the API
	data.ID = types.StringValue(fmt.Sprintf("%d", siteGroup.GetId()))
	data.Name = types.StringValue(siteGroup.GetName())
	data.Slug = types.StringValue(siteGroup.GetSlug())

	// Handle parent - check both HasParent and that ID is non-zero
	if siteGroup.HasParent() {
		parent := siteGroup.GetParent()
		if parent.GetId() != 0 {
			data.Parent = types.StringValue(fmt.Sprintf("%d", parent.GetId()))
		} else {
			data.Parent = types.StringNull()
		}
	} else {
		data.Parent = types.StringNull()
	}

	if siteGroup.HasDescription() {
		data.Description = types.StringValue(siteGroup.GetDescription())
	} else {
		data.Description = types.StringNull()
	}

	// Handle tags
	if siteGroup.HasTags() {
		tags := utils.NestedTagsToTagModels(siteGroup.GetTags())
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
	if siteGroup.HasCustomFields() && !data.CustomFields.IsNull() {
		var stateCustomFields []utils.CustomFieldModel
		resp.Diagnostics.Append(data.CustomFields.ElementsAs(ctx, &stateCustomFields, false)...)
		if resp.Diagnostics.HasError() {
			return
		}

		customFields := utils.MapToCustomFieldModels(siteGroup.GetCustomFields(), stateCustomFields)
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

func (r *SiteGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data SiteGroupResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get the site group ID from state
	siteGroupID := data.ID.ValueString()

	tflog.Debug(ctx, "Updating site group", map[string]interface{}{
		"id":   siteGroupID,
		"name": data.Name.ValueString(),
		"slug": data.Slug.ValueString(),
	})

	// Parse the site group ID to int32 for the API call
	var siteGroupIDInt int32
	if _, err := fmt.Sscanf(siteGroupID, "%d", &siteGroupIDInt); err != nil {
		resp.Diagnostics.AddError(
			"Invalid Site Group ID",
			fmt.Sprintf("Site Group ID must be a number, got: %s", siteGroupID),
		)
		return
	}

	// Prepare the site group update request
	siteGroupRequest := netbox.WritableSiteGroupRequest{
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
		siteGroupRequest.Parent = *netbox.NewNullableInt32(&parentID)
	}

	if !data.Description.IsNull() {
		description := data.Description.ValueString()
		siteGroupRequest.Description = &description
	}

	// Handle tags
	if !data.Tags.IsNull() && !data.Tags.IsUnknown() {
		var tags []utils.TagModel
		resp.Diagnostics.Append(data.Tags.ElementsAs(ctx, &tags, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		siteGroupRequest.Tags = utils.TagsToNestedTagRequests(tags)
	}

	// Handle custom fields
	if !data.CustomFields.IsNull() && !data.CustomFields.IsUnknown() {
		var customFields []utils.CustomFieldModel
		resp.Diagnostics.Append(data.CustomFields.ElementsAs(ctx, &customFields, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		siteGroupRequest.CustomFields = utils.CustomFieldsToMap(customFields)
	}

	// Update the site group via API
	siteGroup, httpResp, err := r.client.DcimAPI.DcimSiteGroupsUpdate(ctx, siteGroupIDInt).WritableSiteGroupRequest(siteGroupRequest).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating site group",
			utils.FormatAPIError(fmt.Sprintf("update site group ID %s", siteGroupID), err, httpResp),
		)
		return
	}

	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Error updating site group",
			fmt.Sprintf("Expected HTTP 200, got: %d", httpResp.StatusCode),
		)
		return
	}

	// Update the model with the response from the API
	data.ID = types.StringValue(fmt.Sprintf("%d", siteGroup.GetId()))
	data.Name = types.StringValue(siteGroup.GetName())
	data.Slug = types.StringValue(siteGroup.GetSlug())

	// Handle parent - check both HasParent and that ID is non-zero
	if siteGroup.HasParent() {
		parent := siteGroup.GetParent()
		if parent.GetId() != 0 {
			data.Parent = types.StringValue(fmt.Sprintf("%d", parent.GetId()))
		} else {
			data.Parent = types.StringNull()
		}
	} else {
		data.Parent = types.StringNull()
	}

	if siteGroup.HasDescription() {
		data.Description = types.StringValue(siteGroup.GetDescription())
	} else {
		data.Description = types.StringNull()
	}

	// Handle tags in response
	if siteGroup.HasTags() {
		tags := utils.NestedTagsToTagModels(siteGroup.GetTags())
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
	if siteGroup.HasCustomFields() && !data.CustomFields.IsNull() {
		var stateCustomFields []utils.CustomFieldModel
		resp.Diagnostics.Append(data.CustomFields.ElementsAs(ctx, &stateCustomFields, false)...)
		if resp.Diagnostics.HasError() {
			return
		}

		customFields := utils.MapToCustomFieldModels(siteGroup.GetCustomFields(), stateCustomFields)
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

func (r *SiteGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data SiteGroupResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get the site group ID from state
	siteGroupID := data.ID.ValueString()

	tflog.Debug(ctx, "Deleting site group", map[string]interface{}{
		"id": siteGroupID,
	})

	// Parse the site group ID to int32 for the API call
	var siteGroupIDInt int32
	if _, err := fmt.Sscanf(siteGroupID, "%d", &siteGroupIDInt); err != nil {
		resp.Diagnostics.AddError(
			"Invalid Site Group ID",
			fmt.Sprintf("Site Group ID must be a number, got: %s", siteGroupID),
		)
		return
	}

	// Delete the site group via API
	httpResp, err := r.client.DcimAPI.DcimSiteGroupsDestroy(ctx, siteGroupIDInt).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting site group",
			utils.FormatAPIError(fmt.Sprintf("delete site group ID %s", siteGroupID), err, httpResp),
		)
		return
	}

	if httpResp.StatusCode != 204 {
		resp.Diagnostics.AddError(
			"Error deleting site group",
			fmt.Sprintf("Expected HTTP 204, got: %d", httpResp.StatusCode),
		)
		return
	}

	tflog.Trace(ctx, "deleted a site group resource")
}

func (r *SiteGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
