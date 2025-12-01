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
	"github.com/hashicorp/terraform-plugin-framework/diag"
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
var _ resource.Resource = &RackRoleResource{}
var _ resource.ResourceWithImportState = &RackRoleResource{}

func NewRackRoleResource() resource.Resource {
	return &RackRoleResource{}
}

// RackRoleResource defines the resource implementation.
type RackRoleResource struct {
	client *netbox.APIClient
}

// RackRoleResourceModel describes the resource data model.
type RackRoleResourceModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Slug         types.String `tfsdk:"slug"`
	Color        types.String `tfsdk:"color"`
	Description  types.String `tfsdk:"description"`
	Tags         types.Set    `tfsdk:"tags"`
	CustomFields types.Set    `tfsdk:"custom_fields"`
}

func (r *RackRoleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_rack_role"
}

func (r *RackRoleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a rack role in Netbox. Rack roles are used to categorize racks by their function or purpose within the data center (e.g., 'Network', 'Compute', 'Storage').",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Unique identifier for the rack role (assigned by Netbox).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Full name of the rack role. This is the human-readable display name.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 100),
				},
			},
			"slug": schema.StringAttribute{
				MarkdownDescription: "URL-friendly identifier for the rack role. Must be unique and contain only alphanumeric characters, hyphens, and underscores.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 100),
					validators.ValidSlug(),
				},
			},
			"color": schema.StringAttribute{
				MarkdownDescription: "Color code for the rack role in hexadecimal format (e.g., 'aa1409' for red). Used for visual identification in the Netbox UI. If not specified, Netbox assigns a default color.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(6, 6),
					stringvalidator.RegexMatches(
						validators.HexColorRegex(),
						"must be a valid 6-character hexadecimal color code (e.g., 'aa1409')",
					),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Detailed description of the rack role, its purpose, or other relevant information.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(200),
				},
			},
			"tags": schema.SetNestedAttribute{
				MarkdownDescription: "Tags assigned to this rack role. Tags provide a way to categorize and organize resources.",
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
				MarkdownDescription: "Custom fields assigned to this rack role. Custom fields allow you to store additional structured data.",
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

func (r *RackRoleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *RackRoleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data RackRoleResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create rack role using go-netbox client
	tflog.Debug(ctx, "Creating rack role", map[string]interface{}{
		"name": data.Name.ValueString(),
		"slug": data.Slug.ValueString(),
	})

	// Prepare the rack role request
	rackRoleRequest := netbox.RackRoleRequest{
		Name: data.Name.ValueString(),
		Slug: data.Slug.ValueString(),
	}

	// Set optional fields if provided
	if !data.Color.IsNull() && !data.Color.IsUnknown() {
		color := data.Color.ValueString()
		rackRoleRequest.Color = &color
	}

	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		description := data.Description.ValueString()
		rackRoleRequest.Description = &description
	}

	// Handle tags
	if !data.Tags.IsNull() && !data.Tags.IsUnknown() {
		var tags []utils.TagModel
		resp.Diagnostics.Append(data.Tags.ElementsAs(ctx, &tags, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		rackRoleRequest.Tags = utils.TagsToNestedTagRequests(tags)
	}

	// Handle custom fields
	if !data.CustomFields.IsNull() && !data.CustomFields.IsUnknown() {
		var customFields []utils.CustomFieldModel
		resp.Diagnostics.Append(data.CustomFields.ElementsAs(ctx, &customFields, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		rackRoleRequest.CustomFields = utils.CustomFieldsToMap(customFields)
	}

	// Create the rack role via API
	rackRole, httpResp, err := r.client.DcimAPI.DcimRackRolesCreate(ctx).RackRoleRequest(rackRoleRequest).Execute()
	if err != nil {
		// Use enhanced error handler that detects duplicates and provides import hints
		handler := utils.CreateErrorHandler{
			ResourceType: "netbox_rack_role",
			ResourceName: "this.rack_role",
			SlugValue:    data.Slug.ValueString(),
			LookupFunc: func(lookupCtx context.Context, slug string) (string, error) {
				// Try to look up existing rack role by slug
				list, _, lookupErr := r.client.DcimAPI.DcimRackRolesList(lookupCtx).Slug([]string{slug}).Execute()
				if lookupErr != nil {
					return "", lookupErr
				}
				if list != nil && len(list.Results) > 0 {
					return fmt.Sprintf("%d", list.Results[0].GetId()), nil
				}
				return "", nil
			},
		}
		handler.HandleCreateError(ctx, err, httpResp, &resp.Diagnostics)
		return
	}

	if httpResp.StatusCode != 201 {
		resp.Diagnostics.AddError(
			"Error creating rack role",
			fmt.Sprintf("Expected HTTP 201, got: %d", httpResp.StatusCode),
		)
		return
	}

	// Map response to state
	r.mapRackRoleToState(ctx, rackRole, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "created a rack role resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RackRoleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data RackRoleResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get the rack role ID from state
	rackRoleID := data.ID.ValueString()

	tflog.Debug(ctx, "Reading rack role", map[string]interface{}{
		"id": rackRoleID,
	})

	// Parse the rack role ID to int32 for the API call
	var rackRoleIDInt int32
	if _, err := fmt.Sscanf(rackRoleID, "%d", &rackRoleIDInt); err != nil {
		resp.Diagnostics.AddError(
			"Invalid Rack Role ID",
			fmt.Sprintf("Rack Role ID must be a number, got: %s", rackRoleID),
		)
		return
	}

	// Retrieve the rack role via API
	rackRole, httpResp, err := r.client.DcimAPI.DcimRackRolesRetrieve(ctx, rackRoleIDInt).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			// Rack role no longer exists, remove from state
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading rack role",
			utils.FormatAPIError(fmt.Sprintf("read rack role ID %s", rackRoleID), err, httpResp),
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

	// Map response to state
	r.mapRackRoleToState(ctx, rackRole, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RackRoleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data RackRoleResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get the rack role ID from state
	rackRoleID := data.ID.ValueString()

	tflog.Debug(ctx, "Updating rack role", map[string]interface{}{
		"id":   rackRoleID,
		"name": data.Name.ValueString(),
		"slug": data.Slug.ValueString(),
	})

	// Parse the rack role ID to int32 for the API call
	var rackRoleIDInt int32
	if _, err := fmt.Sscanf(rackRoleID, "%d", &rackRoleIDInt); err != nil {
		resp.Diagnostics.AddError(
			"Invalid Rack Role ID",
			fmt.Sprintf("Rack Role ID must be a number, got: %s", rackRoleID),
		)
		return
	}

	// Prepare the rack role update request
	rackRoleRequest := netbox.RackRoleRequest{
		Name: data.Name.ValueString(),
		Slug: data.Slug.ValueString(),
	}

	// Set optional fields if provided
	if !data.Color.IsNull() && !data.Color.IsUnknown() {
		color := data.Color.ValueString()
		rackRoleRequest.Color = &color
	}

	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		description := data.Description.ValueString()
		rackRoleRequest.Description = &description
	}

	// Handle tags
	if !data.Tags.IsNull() && !data.Tags.IsUnknown() {
		var tags []utils.TagModel
		resp.Diagnostics.Append(data.Tags.ElementsAs(ctx, &tags, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		rackRoleRequest.Tags = utils.TagsToNestedTagRequests(tags)
	}

	// Handle custom fields
	if !data.CustomFields.IsNull() && !data.CustomFields.IsUnknown() {
		var customFields []utils.CustomFieldModel
		resp.Diagnostics.Append(data.CustomFields.ElementsAs(ctx, &customFields, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		rackRoleRequest.CustomFields = utils.CustomFieldsToMap(customFields)
	}

	// Update the rack role via API
	rackRole, httpResp, err := r.client.DcimAPI.DcimRackRolesUpdate(ctx, rackRoleIDInt).RackRoleRequest(rackRoleRequest).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating rack role",
			utils.FormatAPIError(fmt.Sprintf("update rack role ID %s", rackRoleID), err, httpResp),
		)
		return
	}

	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Error updating rack role",
			fmt.Sprintf("Expected HTTP 200, got: %d", httpResp.StatusCode),
		)
		return
	}

	// Map response to state
	r.mapRackRoleToState(ctx, rackRole, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RackRoleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data RackRoleResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get the rack role ID from state
	rackRoleID := data.ID.ValueString()

	tflog.Debug(ctx, "Deleting rack role", map[string]interface{}{
		"id": rackRoleID,
	})

	// Parse the rack role ID to int32 for the API call
	var rackRoleIDInt int32
	if _, err := fmt.Sscanf(rackRoleID, "%d", &rackRoleIDInt); err != nil {
		resp.Diagnostics.AddError(
			"Invalid Rack Role ID",
			fmt.Sprintf("Rack Role ID must be a number, got: %s", rackRoleID),
		)
		return
	}

	// Delete the rack role via API
	httpResp, err := r.client.DcimAPI.DcimRackRolesDestroy(ctx, rackRoleIDInt).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting rack role",
			utils.FormatAPIError(fmt.Sprintf("delete rack role ID %s", rackRoleID), err, httpResp),
		)
		return
	}

	if httpResp.StatusCode != 204 {
		resp.Diagnostics.AddError(
			"Error deleting rack role",
			fmt.Sprintf("Expected HTTP 204, got: %d", httpResp.StatusCode),
		)
		return
	}

	tflog.Trace(ctx, "deleted a rack role resource")
}

func (r *RackRoleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// mapRackRoleToState maps a RackRole API response to the Terraform state model
func (r *RackRoleResource) mapRackRoleToState(ctx context.Context, rackRole *netbox.RackRole, data *RackRoleResourceModel, diags *diag.Diagnostics) {
	data.ID = types.StringValue(fmt.Sprintf("%d", rackRole.GetId()))
	data.Name = types.StringValue(rackRole.GetName())
	data.Slug = types.StringValue(rackRole.GetSlug())

	// Handle color
	if rackRole.HasColor() && rackRole.GetColor() != "" {
		data.Color = types.StringValue(rackRole.GetColor())
	} else if !data.Color.IsNull() {
		// Preserve null if originally null and API returns empty
		data.Color = types.StringNull()
	}

	// Handle description
	if rackRole.HasDescription() && rackRole.GetDescription() != "" {
		data.Description = types.StringValue(rackRole.GetDescription())
	} else if !data.Description.IsNull() {
		// Preserve null if originally null and API returns empty
		data.Description = types.StringNull()
	}

	// Handle tags
	if rackRole.HasTags() {
		tags := utils.NestedTagsToTagModels(rackRole.GetTags())
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
	if rackRole.HasCustomFields() && !data.CustomFields.IsNull() {
		var stateCustomFields []utils.CustomFieldModel
		cfDiags := data.CustomFields.ElementsAs(ctx, &stateCustomFields, false)
		diags.Append(cfDiags...)
		if diags.HasError() {
			return
		}

		customFields := utils.MapToCustomFieldModels(rackRole.GetCustomFields(), stateCustomFields)
		customFieldsValue, cfValueDiags := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, customFields)
		diags.Append(cfValueDiags...)
		if diags.HasError() {
			return
		}
		data.CustomFields = customFieldsValue
	} else if data.CustomFields.IsNull() {
		data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)
	}
}
