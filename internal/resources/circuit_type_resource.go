// Package resources contains Terraform resource implementations for the Netbox provider.

package resources

import (
	"context"
	"fmt"
	"regexp"

	"github.com/bab3l/go-netbox"
	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
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
)

// Ensure provider defined types fully satisfy framework interfaces.

var (
	_ resource.Resource                = &CircuitTypeResource{}
	_ resource.ResourceWithConfigure   = &CircuitTypeResource{}
	_ resource.ResourceWithImportState = &CircuitTypeResource{}
)

// NewCircuitTypeResource returns a new circuit type resource.
func NewCircuitTypeResource() resource.Resource {
	return &CircuitTypeResource{}
}

// CircuitTypeResource defines the circuit type resource implementation.
type CircuitTypeResource struct {
	client *netbox.APIClient
}

// CircuitTypeResourceModel describes the circuit type resource data model.
type CircuitTypeResourceModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Slug         types.String `tfsdk:"slug"`
	DisplayName  types.String `tfsdk:"display_name"`
	Description  types.String `tfsdk:"description"`
	Color        types.String `tfsdk:"color"`
	Tags         types.Set    `tfsdk:"tags"`
	CustomFields types.Set    `tfsdk:"custom_fields"`
}

// Metadata returns the resource type name.
func (r *CircuitTypeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_circuit_type"
}

// Schema defines the schema for the resource.
func (r *CircuitTypeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a circuit type in Netbox. Circuit types categorize the various types of circuits used by your organization (e.g., Internet Transit, MPLS, Point-to-Point, Metro Ethernet, etc.).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique numeric ID of the circuit type.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the circuit type. This should be descriptive and human-readable (e.g., 'Internet Transit', 'MPLS VPN', 'Dark Fiber').",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"slug": schema.StringAttribute{
				MarkdownDescription: "The URL-friendly slug for the circuit type. Must contain only lowercase letters, numbers, and hyphens.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[-a-z0-9_]+$`),
						"must contain only lowercase letters, numbers, underscores, and hyphens",
					),
				},
			},
			"display_name": nbschema.DisplayNameAttribute("circuit type"),
			"description": schema.StringAttribute{
				MarkdownDescription: "A description of the circuit type.",
				Optional:            true,
			},
			"color": schema.StringAttribute{
				MarkdownDescription: "The color to use when displaying this circuit type (6-character hex code without the leading #, e.g., 'aa1409').",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[0-9a-fA-F]{6}$`),
						"must be a 6-character hex color code (e.g., 'aa1409')",
					),
				},
			},
			"tags":          nbschema.TagsAttribute(),
			"custom_fields": nbschema.CustomFieldsAttribute(),
		},
	}
}

// Configure sets up the resource with the provider client.
func (r *CircuitTypeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create creates a new circuit type resource.
func (r *CircuitTypeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data CircuitTypeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the create request
	createReq := netbox.CircuitTypeRequest{
		Name: data.Name.ValueString(),
		Slug: data.Slug.ValueString(),
	}

	// Handle optional fields
	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		createReq.SetDescription(data.Description.ValueString())
	}
	if !data.Color.IsNull() && !data.Color.IsUnknown() {
		createReq.SetColor(data.Color.ValueString())
	}

	// Handle tags
	if !data.Tags.IsNull() && !data.Tags.IsUnknown() {
		var tags []utils.TagModel
		resp.Diagnostics.Append(data.Tags.ElementsAs(ctx, &tags, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		createReq.Tags = utils.TagsToNestedTagRequests(tags)
	}

	// Handle custom fields
	if !data.CustomFields.IsNull() && !data.CustomFields.IsUnknown() {
		var customFields []utils.CustomFieldModel
		resp.Diagnostics.Append(data.CustomFields.ElementsAs(ctx, &customFields, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		createReq.CustomFields = utils.CustomFieldsToMap(customFields)
	}
	tflog.Debug(ctx, "Creating circuit type", map[string]interface{}{
		"name": data.Name.ValueString(),
		"slug": data.Slug.ValueString(),
	})

	// Create the circuit type
	circuitType, httpResp, err := r.client.CircuitsAPI.CircuitsCircuitTypesCreate(ctx).CircuitTypeRequest(createReq).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating circuit type",
			utils.FormatAPIError("create circuit type", err, httpResp),
		)
		return
	}
	tflog.Debug(ctx, "Created circuit type", map[string]interface{}{
		"id":   circuitType.GetId(),
		"name": circuitType.GetName(),
	})
	// Map the response to state
	r.mapCircuitTypeToState(ctx, circuitType, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read reads the circuit type resource.
func (r *CircuitTypeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data CircuitTypeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Could not parse circuit type ID: %s", err))
		return
	}
	tflog.Debug(ctx, "Reading circuit type", map[string]interface{}{
		"id": id,
	})
	circuitType, httpResp, err := r.client.CircuitsAPI.CircuitsCircuitTypesRetrieve(ctx, id).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			tflog.Debug(ctx, "Circuit type not found, removing from state", map[string]interface{}{
				"id": id,
			})
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading circuit type",
			utils.FormatAPIError("read circuit type", err, httpResp),
		)
		return
	}

	// Map the response to state
	r.mapCircuitTypeToState(ctx, circuitType, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates the circuit type resource.
func (r *CircuitTypeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data CircuitTypeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Could not parse circuit type ID: %s", err))
		return
	}

	// Build the update request
	updateReq := netbox.CircuitTypeRequest{
		Name: data.Name.ValueString(),
		Slug: data.Slug.ValueString(),
	}

	// Handle optional fields
	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		updateReq.SetDescription(data.Description.ValueString())
	} else {
		updateReq.SetDescription("")
	}
	if !data.Color.IsNull() && !data.Color.IsUnknown() {
		updateReq.SetColor(data.Color.ValueString())
	} else {
		updateReq.SetColor("")
	}

	// Handle tags
	if !data.Tags.IsNull() && !data.Tags.IsUnknown() {
		var tags []utils.TagModel
		resp.Diagnostics.Append(data.Tags.ElementsAs(ctx, &tags, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		updateReq.Tags = utils.TagsToNestedTagRequests(tags)
	} else {
		updateReq.Tags = []netbox.NestedTagRequest{}
	}

	// Handle custom fields
	if !data.CustomFields.IsNull() && !data.CustomFields.IsUnknown() {
		var customFields []utils.CustomFieldModel
		resp.Diagnostics.Append(data.CustomFields.ElementsAs(ctx, &customFields, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		updateReq.CustomFields = utils.CustomFieldsToMap(customFields)
	}
	tflog.Debug(ctx, "Updating circuit type", map[string]interface{}{
		"id":   id,
		"name": data.Name.ValueString(),
	})

	// Update the circuit type
	circuitType, httpResp, err := r.client.CircuitsAPI.CircuitsCircuitTypesUpdate(ctx, id).CircuitTypeRequest(updateReq).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating circuit type",
			utils.FormatAPIError("update circuit type", err, httpResp),
		)
		return
	}
	tflog.Debug(ctx, "Updated circuit type", map[string]interface{}{
		"id":   circuitType.GetId(),
		"name": circuitType.GetName(),
	})

	// Map the response to state
	r.mapCircuitTypeToState(ctx, circuitType, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete deletes the circuit type resource.
func (r *CircuitTypeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CircuitTypeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Could not parse circuit type ID: %s", err))
		return
	}
	tflog.Debug(ctx, "Deleting circuit type", map[string]interface{}{
		"id": id,
	})
	httpResp, err := r.client.CircuitsAPI.CircuitsCircuitTypesDestroy(ctx, id).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			tflog.Debug(ctx, "Circuit type already deleted", map[string]interface{}{
				"id": id,
			})
			return
		}
		resp.Diagnostics.AddError(
			"Error deleting circuit type",
			utils.FormatAPIError("delete circuit type", err, httpResp),
		)
		return
	}
	tflog.Debug(ctx, "Deleted circuit type", map[string]interface{}{
		"id": id,
	})
}

// ImportState imports a circuit type resource.
func (r *CircuitTypeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// mapCircuitTypeToState maps a CircuitType to the Terraform state model.
func (r *CircuitTypeResource) mapCircuitTypeToState(ctx context.Context, circuitType *netbox.CircuitType, data *CircuitTypeResourceModel, diags *diag.Diagnostics) {
	data.ID = types.StringValue(fmt.Sprintf("%d", circuitType.GetId()))
	data.Name = types.StringValue(circuitType.GetName())
	data.Slug = types.StringValue(circuitType.GetSlug())
	data.DisplayName = types.StringValue(circuitType.GetDisplay())

	// Handle description
	if circuitType.HasDescription() && circuitType.GetDescription() != "" {
		data.Description = types.StringValue(circuitType.GetDescription())
	} else {
		data.Description = types.StringNull()
	}

	// Handle color
	if circuitType.HasColor() && circuitType.GetColor() != "" {
		data.Color = types.StringValue(circuitType.GetColor())
	} else {
		data.Color = types.StringNull()
	}

	// Handle tags
	data.Tags = utils.PopulateTagsFromNestedTags(ctx, circuitType.HasTags(), circuitType.GetTags(), diags)

	// Handle custom fields
	data.CustomFields = utils.PopulateCustomFieldsFromMap(ctx, circuitType.HasCustomFields(), circuitType.GetCustomFields(), data.CustomFields, diags)
}
