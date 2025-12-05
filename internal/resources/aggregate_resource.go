// Package resources contains Terraform resource implementations for the Netbox provider.
package resources

import (
	"context"
	"fmt"

	"github.com/bab3l/go-netbox"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/bab3l/terraform-provider-netbox/internal/netboxlookup"
	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &AggregateResource{}
	_ resource.ResourceWithConfigure   = &AggregateResource{}
	_ resource.ResourceWithImportState = &AggregateResource{}
)

// NewAggregateResource returns a new Aggregate resource.
func NewAggregateResource() resource.Resource {
	return &AggregateResource{}
}

// AggregateResource defines the resource implementation.
type AggregateResource struct {
	client *netbox.APIClient
}

// AggregateResourceModel describes the resource data model.
type AggregateResourceModel struct {
	ID           types.String `tfsdk:"id"`
	Prefix       types.String `tfsdk:"prefix"`
	RIR          types.String `tfsdk:"rir"`
	Tenant       types.String `tfsdk:"tenant"`
	DateAdded    types.String `tfsdk:"date_added"`
	Description  types.String `tfsdk:"description"`
	Comments     types.String `tfsdk:"comments"`
	Tags         types.Set    `tfsdk:"tags"`
	CustomFields types.Set    `tfsdk:"custom_fields"`
}

// Metadata returns the resource type name.
func (r *AggregateResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_aggregate"
}

// Schema defines the schema for the resource.
func (r *AggregateResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an aggregate in Netbox. Aggregates are top-level IP address blocks that represent the entire address space available for allocation by an organization.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique numeric ID of the aggregate.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"prefix": schema.StringAttribute{
				MarkdownDescription: "The IP prefix in CIDR notation (e.g., 10.0.0.0/8, 2001:db8::/32).",
				Required:            true,
			},
			"rir": schema.StringAttribute{
				MarkdownDescription: "The name, slug, or ID of the Regional Internet Registry (RIR) this aggregate belongs to.",
				Required:            true,
			},
			"tenant": schema.StringAttribute{
				MarkdownDescription: "The name or ID of the tenant this aggregate is assigned to.",
				Optional:            true,
			},
			"date_added": schema.StringAttribute{
				MarkdownDescription: "The date this aggregate was added (YYYY-MM-DD format).",
				Optional:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "A description of the aggregate.",
				Optional:            true,
			},
			"comments": schema.StringAttribute{
				MarkdownDescription: "Additional comments about the aggregate.",
				Optional:            true,
			},
			"tags":          nbschema.TagsAttribute(),
			"custom_fields": nbschema.CustomFieldsAttribute(),
		},
	}
}

// Configure sets the client for the resource.
func (r *AggregateResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create creates a new aggregate resource.
func (r *AggregateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data AggregateResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the create request
	createReq, diags := r.buildCreateRequest(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating aggregate", map[string]interface{}{
		"prefix": data.Prefix.ValueString(),
	})

	// Call API to create aggregate
	aggregate, httpResp, err := r.client.IpamAPI.IpamAggregatesCreate(ctx).WritableAggregateRequest(*createReq).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating aggregate",
			fmt.Sprintf("Could not create aggregate: %s\nHTTP Response: %v", err.Error(), httpResp),
		)
		return
	}

	// Map response to model
	r.mapResponseToModel(ctx, aggregate, &data)

	tflog.Debug(ctx, "Created aggregate", map[string]interface{}{
		"id": data.ID.ValueString(),
	})

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read reads the aggregate resource.
func (r *AggregateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data AggregateResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid ID",
			fmt.Sprintf("Could not convert ID to integer: %s", err.Error()),
		)
		return
	}

	tflog.Debug(ctx, "Reading aggregate", map[string]interface{}{
		"id": id,
	})

	// Call API to read aggregate
	aggregate, httpResp, err := r.client.IpamAPI.IpamAggregatesRetrieve(ctx, id).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			tflog.Debug(ctx, "Aggregate not found, removing from state", map[string]interface{}{
				"id": id,
			})
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading aggregate",
			fmt.Sprintf("Could not read aggregate: %s\nHTTP Response: %v", err.Error(), httpResp),
		)
		return
	}

	// Map response to model
	r.mapResponseToModel(ctx, aggregate, &data)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates the aggregate resource.
func (r *AggregateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data AggregateResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid ID",
			fmt.Sprintf("Could not convert ID to integer: %s", err.Error()),
		)
		return
	}

	// Build the update request
	updateReq, diags := r.buildCreateRequest(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Updating aggregate", map[string]interface{}{
		"id": id,
	})

	// Call API to update aggregate
	aggregate, httpResp, err := r.client.IpamAPI.IpamAggregatesUpdate(ctx, id).WritableAggregateRequest(*updateReq).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating aggregate",
			fmt.Sprintf("Could not update aggregate: %s\nHTTP Response: %v", err.Error(), httpResp),
		)
		return
	}

	// Map response to model
	r.mapResponseToModel(ctx, aggregate, &data)

	tflog.Debug(ctx, "Updated aggregate", map[string]interface{}{
		"id": data.ID.ValueString(),
	})

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete deletes the aggregate resource.
func (r *AggregateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data AggregateResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid ID",
			fmt.Sprintf("Could not convert ID to integer: %s", err.Error()),
		)
		return
	}

	tflog.Debug(ctx, "Deleting aggregate", map[string]interface{}{
		"id": id,
	})

	// Call API to delete aggregate
	httpResp, err := r.client.IpamAPI.IpamAggregatesDestroy(ctx, id).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			tflog.Debug(ctx, "Aggregate already deleted", map[string]interface{}{
				"id": id,
			})
			return
		}
		resp.Diagnostics.AddError(
			"Error deleting aggregate",
			fmt.Sprintf("Could not delete aggregate: %s\nHTTP Response: %v", err.Error(), httpResp),
		)
		return
	}

	tflog.Debug(ctx, "Deleted aggregate", map[string]interface{}{
		"id": id,
	})
}

// ImportState imports an existing aggregate.
func (r *AggregateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// buildCreateRequest builds a WritableAggregateRequest from the model.
func (r *AggregateResource) buildCreateRequest(ctx context.Context, data *AggregateResourceModel) (*netbox.WritableAggregateRequest, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Look up RIR (required)
	rir, rirDiags := netboxlookup.LookupRIR(ctx, r.client, data.RIR.ValueString())
	diags.Append(rirDiags...)
	if diags.HasError() {
		return nil, diags
	}

	createReq := netbox.NewWritableAggregateRequest(data.Prefix.ValueString(), *rir)

	// Handle tenant (optional)
	if !data.Tenant.IsNull() && !data.Tenant.IsUnknown() {
		tenant, tenantDiags := netboxlookup.LookupTenant(ctx, r.client, data.Tenant.ValueString())
		diags.Append(tenantDiags...)
		if diags.HasError() {
			return nil, diags
		}
		createReq.SetTenant(*tenant)
	}

	// Handle date_added (optional)
	if !data.DateAdded.IsNull() && !data.DateAdded.IsUnknown() {
		createReq.SetDateAdded(data.DateAdded.ValueString())
	}

	// Handle description (optional)
	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		createReq.SetDescription(data.Description.ValueString())
	}

	// Handle comments (optional)
	if !data.Comments.IsNull() && !data.Comments.IsUnknown() {
		createReq.SetComments(data.Comments.ValueString())
	}

	// Handle tags
	if !data.Tags.IsNull() && !data.Tags.IsUnknown() {
		tags, tagDiags := utils.TagModelsToNestedTagRequests(ctx, data.Tags)
		diags.Append(tagDiags...)
		if diags.HasError() {
			return nil, diags
		}
		createReq.SetTags(tags)
	}

	// Handle custom fields
	if !data.CustomFields.IsNull() && !data.CustomFields.IsUnknown() {
		var customFields []utils.CustomFieldModel
		cfDiags := data.CustomFields.ElementsAs(ctx, &customFields, false)
		diags.Append(cfDiags...)
		if diags.HasError() {
			return nil, diags
		}
		createReq.SetCustomFields(utils.CustomFieldsToMap(customFields))
	}

	return createReq, diags
}

// mapResponseToModel maps the API response to the Terraform model.
func (r *AggregateResource) mapResponseToModel(ctx context.Context, aggregate *netbox.Aggregate, data *AggregateResourceModel) {
	data.ID = types.StringValue(fmt.Sprintf("%d", aggregate.GetId()))
	data.Prefix = types.StringValue(aggregate.GetPrefix())

	// Map RIR
	if rir := aggregate.GetRir(); rir.Id != 0 {
		data.RIR = types.StringValue(fmt.Sprintf("%d", rir.Id))
	}

	// Map tenant
	if tenant, ok := aggregate.GetTenantOk(); ok && tenant != nil && tenant.Id != 0 {
		data.Tenant = types.StringValue(fmt.Sprintf("%d", tenant.GetId()))
	} else if data.Tenant.IsNull() {
		// Keep null if it was null
	} else {
		data.Tenant = types.StringNull()
	}

	// Map date_added
	if dateAdded := aggregate.GetDateAdded(); dateAdded != "" {
		data.DateAdded = types.StringValue(dateAdded)
	} else if data.DateAdded.IsNull() {
		// Keep null if it was null
	} else {
		data.DateAdded = types.StringNull()
	}

	// Map description
	if description, ok := aggregate.GetDescriptionOk(); ok && description != nil {
		data.Description = types.StringValue(*description)
	} else if data.Description.IsNull() {
		// Keep null if it was null
	} else {
		data.Description = types.StringNull()
	}

	// Map comments
	if comments, ok := aggregate.GetCommentsOk(); ok && comments != nil {
		data.Comments = types.StringValue(*comments)
	} else if data.Comments.IsNull() {
		// Keep null if it was null
	} else {
		data.Comments = types.StringNull()
	}

	// Tags
	if len(aggregate.Tags) > 0 {
		tags := utils.NestedTagsToTagModels(aggregate.Tags)
		tagsValue, _ := types.SetValueFrom(ctx, utils.GetTagsAttributeType().ElemType, tags)
		data.Tags = tagsValue
	} else {
		data.Tags = types.SetNull(utils.GetTagsAttributeType().ElemType)
	}

	// Custom Fields
	if len(aggregate.CustomFields) > 0 && !data.CustomFields.IsNull() {
		var stateCustomFields []utils.CustomFieldModel
		data.CustomFields.ElementsAs(ctx, &stateCustomFields, false)
		customFields := utils.MapToCustomFieldModels(aggregate.CustomFields, stateCustomFields)
		customFieldsValue, _ := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, customFields)
		data.CustomFields = customFieldsValue
	} else if len(aggregate.CustomFields) > 0 {
		customFields := utils.MapToCustomFieldModels(aggregate.CustomFields, []utils.CustomFieldModel{})
		customFieldsValue, _ := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, customFields)
		data.CustomFields = customFieldsValue
	} else {
		data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)
	}
}
