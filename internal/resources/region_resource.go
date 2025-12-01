// Package resources contains Terraform resource implementations for the Netbox provider.
package resources

import (
	"context"
	"fmt"

	"github.com/bab3l/go-netbox"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/bab3l/terraform-provider-netbox/internal/validators"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &RegionResource{}
var _ resource.ResourceWithImportState = &RegionResource{}

func NewRegionResource() resource.Resource {
	return &RegionResource{}
}

// RegionResource defines the resource implementation.
type RegionResource struct {
	client *netbox.APIClient
}

// RegionResourceModel describes the resource data model.
type RegionResourceModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Slug         types.String `tfsdk:"slug"`
	Parent       types.String `tfsdk:"parent"`
	Description  types.String `tfsdk:"description"`
	Tags         types.Set    `tfsdk:"tags"`
	CustomFields types.Set    `tfsdk:"custom_fields"`
}

func (r *RegionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_region"
}

func (r *RegionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a region in Netbox. Regions provide a hierarchical way to organize sites geographically, such as continents, countries, states, or cities.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Unique identifier for the region (assigned by Netbox).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Full name of the region (e.g., 'North America', 'United States', 'California').",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 100),
				},
			},
			"slug": schema.StringAttribute{
				MarkdownDescription: "URL-friendly identifier for the region. Must be unique and contain only alphanumeric characters, hyphens, and underscores.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 100),
					validators.ValidSlug(),
				},
			},
			"parent": schema.StringAttribute{
				MarkdownDescription: "ID of the parent region. Leave empty for top-level regions. This enables hierarchical organization of geographic areas.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						validators.IntegerRegex(),
						"must be a valid integer ID",
					),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Detailed description of the region.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(200),
				},
			},
			"tags": schema.SetNestedAttribute{
				MarkdownDescription: "Tags assigned to this region.",
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
				MarkdownDescription: "Custom fields assigned to this region.",
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
							MarkdownDescription: "Type of the custom field.",
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

func (r *RegionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *RegionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data RegionResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating region", map[string]interface{}{
		"name": data.Name.ValueString(),
		"slug": data.Slug.ValueString(),
	})

	// Build the API request
	regionRequest := netbox.NewWritableRegionRequest(data.Name.ValueString(), data.Slug.ValueString())

	// Set optional parent
	if !data.Parent.IsNull() && !data.Parent.IsUnknown() {
		parentID := data.Parent.ValueString()
		var parentIDInt int32
		if _, err := fmt.Sscanf(parentID, "%d", &parentIDInt); err != nil {
			resp.Diagnostics.AddError(
				"Invalid Parent ID",
				fmt.Sprintf("Parent ID must be a number, got: %s", parentID),
			)
			return
		}
		regionRequest.Parent = *netbox.NewNullableInt32(&parentIDInt)
	}

	// Set optional description
	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		desc := data.Description.ValueString()
		regionRequest.Description = &desc
	}

	// Handle tags
	if !data.Tags.IsNull() && !data.Tags.IsUnknown() {
		tags, diags := utils.TagModelsToNestedTagRequests(ctx, data.Tags)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		regionRequest.Tags = tags
	}

	// Handle custom fields
	if !data.CustomFields.IsNull() && !data.CustomFields.IsUnknown() {
		var customFieldModels []utils.CustomFieldModel
		diags := data.CustomFields.ElementsAs(ctx, &customFieldModels, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		regionRequest.CustomFields = utils.CustomFieldModelsToMap(customFieldModels)
	}

	// Call the API
	region, httpResp, err := r.client.DcimAPI.DcimRegionsCreate(ctx).WritableRegionRequest(*regionRequest).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating region",
			utils.FormatAPIError("create region", err, httpResp),
		)
		return
	}

	if httpResp.StatusCode != 201 {
		resp.Diagnostics.AddError(
			"Error creating region",
			fmt.Sprintf("Expected HTTP 201, got: %d", httpResp.StatusCode),
		)
		return
	}

	// Map response to state
	data.ID = types.StringValue(fmt.Sprintf("%d", region.GetId()))
	data.Name = types.StringValue(region.GetName())
	data.Slug = types.StringValue(region.GetSlug())

	if region.HasParent() && region.GetParent().Id != 0 {
		parent := region.GetParent()
		data.Parent = types.StringValue(fmt.Sprintf("%d", parent.GetId()))
	} else {
		data.Parent = types.StringNull()
	}

	if region.HasDescription() {
		desc := region.GetDescription()
		if desc == "" && data.Description.IsNull() {
			data.Description = types.StringNull()
		} else if desc == "" {
			data.Description = types.StringNull()
		} else {
			data.Description = types.StringValue(desc)
		}
	} else {
		data.Description = types.StringNull()
	}

	// Handle tags in response
	if region.HasTags() {
		tags := utils.NestedTagsToTagModels(region.GetTags())
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
	if region.HasCustomFields() {
		var existingModels []utils.CustomFieldModel
		if !data.CustomFields.IsNull() {
			diags := data.CustomFields.ElementsAs(ctx, &existingModels, false)
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}
		}
		customFields := utils.MapToCustomFieldModels(region.GetCustomFields(), existingModels)
		customFieldsValue, diags := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, customFields)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		data.CustomFields = customFieldsValue
	} else {
		data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)
	}

	tflog.Trace(ctx, "created a region resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RegionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data RegionResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	regionID := data.ID.ValueString()
	tflog.Debug(ctx, "Reading region", map[string]interface{}{
		"id": regionID,
	})

	var regionIDInt int32
	if _, err := fmt.Sscanf(regionID, "%d", &regionIDInt); err != nil {
		resp.Diagnostics.AddError(
			"Invalid Region ID",
			fmt.Sprintf("Region ID must be a number, got: %s", regionID),
		)
		return
	}

	region, httpResp, err := r.client.DcimAPI.DcimRegionsRetrieve(ctx, regionIDInt).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading region",
			utils.FormatAPIError(fmt.Sprintf("read region ID %s", regionID), err, httpResp),
		)
		return
	}

	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Error reading region",
			fmt.Sprintf("Expected HTTP 200, got: %d", httpResp.StatusCode),
		)
		return
	}

	// Map response to state
	data.ID = types.StringValue(fmt.Sprintf("%d", region.GetId()))
	data.Name = types.StringValue(region.GetName())
	data.Slug = types.StringValue(region.GetSlug())

	if region.HasParent() && region.GetParent().Id != 0 {
		parent := region.GetParent()
		data.Parent = types.StringValue(fmt.Sprintf("%d", parent.GetId()))
	} else {
		data.Parent = types.StringNull()
	}

	if region.HasDescription() {
		desc := region.GetDescription()
		if desc == "" && data.Description.IsNull() {
			data.Description = types.StringNull()
		} else if desc == "" {
			data.Description = types.StringNull()
		} else {
			data.Description = types.StringValue(desc)
		}
	} else {
		data.Description = types.StringNull()
	}

	// Handle tags
	if region.HasTags() {
		tags := utils.NestedTagsToTagModels(region.GetTags())
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
	if region.HasCustomFields() {
		var existingModels []utils.CustomFieldModel
		if !data.CustomFields.IsNull() {
			diags := data.CustomFields.ElementsAs(ctx, &existingModels, false)
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}
		}
		customFields := utils.MapToCustomFieldModels(region.GetCustomFields(), existingModels)
		customFieldsValue, diags := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, customFields)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		data.CustomFields = customFieldsValue
	} else {
		data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RegionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data RegionResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	regionID := data.ID.ValueString()
	tflog.Debug(ctx, "Updating region", map[string]interface{}{
		"id": regionID,
	})

	var regionIDInt int32
	if _, err := fmt.Sscanf(regionID, "%d", &regionIDInt); err != nil {
		resp.Diagnostics.AddError(
			"Invalid Region ID",
			fmt.Sprintf("Region ID must be a number, got: %s", regionID),
		)
		return
	}

	// Build the API request
	regionRequest := netbox.NewWritableRegionRequest(data.Name.ValueString(), data.Slug.ValueString())

	// Set optional parent
	if !data.Parent.IsNull() && !data.Parent.IsUnknown() {
		parentID := data.Parent.ValueString()
		var parentIDInt int32
		if _, err := fmt.Sscanf(parentID, "%d", &parentIDInt); err != nil {
			resp.Diagnostics.AddError(
				"Invalid Parent ID",
				fmt.Sprintf("Parent ID must be a number, got: %s", parentID),
			)
			return
		}
		regionRequest.Parent = *netbox.NewNullableInt32(&parentIDInt)
	}

	// Set optional description
	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		desc := data.Description.ValueString()
		regionRequest.Description = &desc
	}

	// Handle tags
	if !data.Tags.IsNull() && !data.Tags.IsUnknown() {
		tags, diags := utils.TagModelsToNestedTagRequests(ctx, data.Tags)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		regionRequest.Tags = tags
	}

	// Handle custom fields
	if !data.CustomFields.IsNull() && !data.CustomFields.IsUnknown() {
		var customFieldModels []utils.CustomFieldModel
		diags := data.CustomFields.ElementsAs(ctx, &customFieldModels, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		regionRequest.CustomFields = utils.CustomFieldModelsToMap(customFieldModels)
	}

	// Call the API
	region, httpResp, err := r.client.DcimAPI.DcimRegionsUpdate(ctx, regionIDInt).WritableRegionRequest(*regionRequest).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating region",
			utils.FormatAPIError(fmt.Sprintf("update region ID %s", regionID), err, httpResp),
		)
		return
	}

	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Error updating region",
			fmt.Sprintf("Expected HTTP 200, got: %d", httpResp.StatusCode),
		)
		return
	}

	// Map response to state
	data.ID = types.StringValue(fmt.Sprintf("%d", region.GetId()))
	data.Name = types.StringValue(region.GetName())
	data.Slug = types.StringValue(region.GetSlug())

	if region.HasParent() && region.GetParent().Id != 0 {
		parent := region.GetParent()
		data.Parent = types.StringValue(fmt.Sprintf("%d", parent.GetId()))
	} else {
		data.Parent = types.StringNull()
	}

	if region.HasDescription() {
		desc := region.GetDescription()
		if desc == "" && data.Description.IsNull() {
			data.Description = types.StringNull()
		} else if desc == "" {
			data.Description = types.StringNull()
		} else {
			data.Description = types.StringValue(desc)
		}
	} else {
		data.Description = types.StringNull()
	}

	// Handle tags
	if region.HasTags() {
		tags := utils.NestedTagsToTagModels(region.GetTags())
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
	if region.HasCustomFields() {
		var existingModels []utils.CustomFieldModel
		if !data.CustomFields.IsNull() {
			diags := data.CustomFields.ElementsAs(ctx, &existingModels, false)
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}
		}
		customFields := utils.MapToCustomFieldModels(region.GetCustomFields(), existingModels)
		customFieldsValue, diags := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, customFields)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		data.CustomFields = customFieldsValue
	} else {
		data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)
	}

	tflog.Trace(ctx, "updated a region resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RegionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data RegionResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	regionID := data.ID.ValueString()
	tflog.Debug(ctx, "Deleting region", map[string]interface{}{
		"id": regionID,
	})

	var regionIDInt int32
	if _, err := fmt.Sscanf(regionID, "%d", &regionIDInt); err != nil {
		resp.Diagnostics.AddError(
			"Invalid Region ID",
			fmt.Sprintf("Region ID must be a number, got: %s", regionID),
		)
		return
	}

	httpResp, err := r.client.DcimAPI.DcimRegionsDestroy(ctx, regionIDInt).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting region",
			utils.FormatAPIError(fmt.Sprintf("delete region ID %s", regionID), err, httpResp),
		)
		return
	}

	if httpResp.StatusCode != 204 {
		resp.Diagnostics.AddError(
			"Error deleting region",
			fmt.Sprintf("Expected HTTP 204, got: %d", httpResp.StatusCode),
		)
		return
	}

	tflog.Trace(ctx, "deleted a region resource")
}

func (r *RegionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
