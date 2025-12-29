// Package resources contains Terraform resource implementations for the Netbox provider.

package resources

import (
	"context"
	"fmt"
	"maps"

	"github.com/bab3l/go-netbox"
	"github.com/bab3l/terraform-provider-netbox/internal/netboxlookup"
	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.

var (
	_ resource.Resource = &RouteTargetResource{}

	_ resource.ResourceWithConfigure = &RouteTargetResource{}

	_ resource.ResourceWithImportState = &RouteTargetResource{}
)

// NewRouteTargetResource returns a new RouteTarget resource.

func NewRouteTargetResource() resource.Resource {
	return &RouteTargetResource{}
}

// RouteTargetResource defines the resource implementation.

type RouteTargetResource struct {
	client *netbox.APIClient
}

// RouteTargetResourceModel describes the resource data model.

type RouteTargetResourceModel struct {
	ID types.String `tfsdk:"id"`

	Name types.String `tfsdk:"name"`

	Tenant types.String `tfsdk:"tenant"`

	TenantID types.String `tfsdk:"tenant_id"`

	Description types.String `tfsdk:"description"`

	Comments types.String `tfsdk:"comments"`

	Tags types.Set `tfsdk:"tags"`

	CustomFields types.Set `tfsdk:"custom_fields"`
}

// Metadata returns the resource type name.

func (r *RouteTargetResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_route_target"
}

// Schema defines the schema for the resource.

func (r *RouteTargetResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Route Target in Netbox. Route targets are used to control the distribution of routes in VRFs (Virtual Routing and Forwarding) for BGP/MPLS VPN configurations.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique numeric ID of the route target.",

				Computed: true,

				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"name": schema.StringAttribute{
				MarkdownDescription: "The route target value (formatted in accordance with RFC 4360). Required.",

				Required: true,
			},

			"tenant": nbschema.ReferenceAttribute("tenant", "ID or slug of the tenant that owns this route target."),

			"tenant_id": nbschema.ComputedIDAttribute("tenant"),
		},
	}

	// Add description and comments attributes
	maps.Copy(resp.Schema.Attributes, nbschema.CommonDescriptiveAttributes("route target"))

	// Add common metadata attributes (tags, custom_fields)
	maps.Copy(resp.Schema.Attributes, nbschema.CommonMetadataAttributes())
}

// Configure adds the provider configured client to the resource.

func (r *RouteTargetResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create creates the resource and sets the initial Terraform state.

func (r *RouteTargetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data RouteTargetResourceModel

	// Read Terraform plan data into the model

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create the RouteTarget request

	rtRequest := netbox.NewRouteTargetRequest(data.Name.ValueString())

	// Set optional fields

	r.setOptionalFields(ctx, rtRequest, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating RouteTarget", map[string]interface{}{
		"name": data.Name.ValueString(),
	})

	// Create the RouteTarget

	rt, httpResp, err := r.client.IpamAPI.IpamRouteTargetsCreate(ctx).RouteTargetRequest(*rtRequest).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {
		resp.Diagnostics.AddError(

			"Error creating RouteTarget",

			utils.FormatAPIError("create RouteTarget", err, httpResp),
		)

		return
	}

	// Map response to model

	r.mapRouteTargetToState(ctx, rt, &data, &resp.Diagnostics)

	tflog.Debug(ctx, "Created RouteTarget", map[string]interface{}{
		"id": data.ID.ValueString(),

		"name": data.Name.ValueString(),
	})

	// Save data into Terraform state

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read refreshes the Terraform state with the latest data.

func (r *RouteTargetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data RouteTargetResourceModel

	// Read Terraform prior state data into the model

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Parse the ID

	id, err := utils.ParseID(data.ID.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(

			"Invalid ID",

			fmt.Sprintf("Unable to parse ID %q: %s", data.ID.ValueString(), err.Error()),
		)

		return
	}

	tflog.Debug(ctx, "Reading RouteTarget", map[string]interface{}{
		"id": id,
	})

	// Get the RouteTarget from Netbox

	rt, httpResp, err := r.client.IpamAPI.IpamRouteTargetsRetrieve(ctx, id).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			resp.State.RemoveResource(ctx)

			return
		}

		resp.Diagnostics.AddError(

			"Error reading RouteTarget",

			utils.FormatAPIError(fmt.Sprintf("read RouteTarget ID %d", id), err, httpResp),
		)

		return
	}

	// Map response to model

	r.mapRouteTargetToState(ctx, rt, &data, &resp.Diagnostics)

	tflog.Debug(ctx, "Read RouteTarget", map[string]interface{}{
		"id": data.ID.ValueString(),

		"name": data.Name.ValueString(),
	})

	// Save updated data into Terraform state

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates the resource and sets the updated Terraform state on success.

func (r *RouteTargetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data RouteTargetResourceModel

	// Read Terraform plan data into the model

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Parse the ID

	id, err := utils.ParseID(data.ID.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(

			"Invalid ID",

			fmt.Sprintf("Unable to parse ID %q: %s", data.ID.ValueString(), err.Error()),
		)

		return
	}

	// Create the RouteTarget request

	rtRequest := netbox.NewRouteTargetRequest(data.Name.ValueString())

	// Set optional fields

	r.setOptionalFields(ctx, rtRequest, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Updating RouteTarget", map[string]interface{}{
		"id": id,

		"name": data.Name.ValueString(),
	})

	// Update the RouteTarget

	rt, httpResp, err := r.client.IpamAPI.IpamRouteTargetsUpdate(ctx, id).RouteTargetRequest(*rtRequest).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {
		resp.Diagnostics.AddError(

			"Error updating RouteTarget",

			utils.FormatAPIError(fmt.Sprintf("update RouteTarget ID %d", id), err, httpResp),
		)

		return
	}

	// Map response to model

	r.mapRouteTargetToState(ctx, rt, &data, &resp.Diagnostics)

	tflog.Debug(ctx, "Updated RouteTarget", map[string]interface{}{
		"id": data.ID.ValueString(),

		"name": data.Name.ValueString(),
	})

	// Save updated data into Terraform state

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete deletes the resource and removes the Terraform state on success.

func (r *RouteTargetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data RouteTargetResourceModel

	// Read Terraform prior state data into the model

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Parse the ID

	id, err := utils.ParseID(data.ID.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(

			"Invalid ID",

			fmt.Sprintf("Unable to parse ID %q: %s", data.ID.ValueString(), err.Error()),
		)

		return
	}

	tflog.Debug(ctx, "Deleting RouteTarget", map[string]interface{}{
		"id": id,
	})

	// Delete the RouteTarget

	httpResp, err := r.client.IpamAPI.IpamRouteTargetsDestroy(ctx, id).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {
		resp.Diagnostics.AddError(

			"Error deleting RouteTarget",

			utils.FormatAPIError(fmt.Sprintf("delete RouteTarget ID %d", id), err, httpResp),
		)

		return
	}

	tflog.Debug(ctx, "Deleted RouteTarget", map[string]interface{}{
		"id": id,
	})
}

func (r *RouteTargetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// setOptionalFields sets optional fields on the RouteTarget request from the resource model.

func (r *RouteTargetResource) setOptionalFields(ctx context.Context, rtRequest *netbox.RouteTargetRequest, data *RouteTargetResourceModel, diags *diag.Diagnostics) {
	// Tenant

	if utils.IsSet(data.Tenant) {
		tenantRef, tenantDiags := netboxlookup.LookupTenant(ctx, r.client, data.Tenant.ValueString())

		diags.Append(tenantDiags...)

		if diags.HasError() {
			return
		}

		rtRequest.Tenant = *netbox.NewNullableBriefTenantRequest(tenantRef)
	}

	// Set common fields (description, comments, tags, custom_fields)
	utils.ApplyCommonFields(ctx, rtRequest, data.Description, data.Comments, data.Tags, data.CustomFields, diags)
	if diags.HasError() {
		return
	}
}

// mapRouteTargetToState maps a Netbox RouteTarget to the Terraform state model.

func (r *RouteTargetResource) mapRouteTargetToState(ctx context.Context, rt *netbox.RouteTarget, data *RouteTargetResourceModel, diags *diag.Diagnostics) {
	data.ID = types.StringValue(fmt.Sprintf("%d", rt.Id))

	data.Name = types.StringValue(rt.Name)

	// Tenant

	if rt.HasTenant() && rt.Tenant.Get() != nil {
		data.Tenant = utils.UpdateReferenceAttribute(data.Tenant, rt.Tenant.Get().Name, rt.Tenant.Get().Slug, rt.Tenant.Get().Id)
		data.TenantID = types.StringValue(fmt.Sprintf("%d", rt.Tenant.Get().Id))
	} else {
		data.Tenant = types.StringNull()
		data.TenantID = types.StringNull()
	}

	// Description

	if rt.Description != nil && *rt.Description != "" {
		data.Description = types.StringValue(*rt.Description)
	} else {
		data.Description = types.StringNull()
	}

	// Comments

	if rt.Comments != nil && *rt.Comments != "" {
		data.Comments = types.StringValue(*rt.Comments)
	} else {
		data.Comments = types.StringNull()
	}

	// Tags

	if len(rt.Tags) > 0 {
		tags := utils.NestedTagsToTagModels(rt.Tags)

		tagsValue, _ := types.SetValueFrom(ctx, utils.GetTagsAttributeType().ElemType, tags)

		data.Tags = tagsValue
	} else {
		data.Tags = types.SetNull(utils.GetTagsAttributeType().ElemType)
	}

	// Custom Fields

	switch {
	case len(rt.CustomFields) > 0 && !data.CustomFields.IsNull():

		var stateCustomFields []utils.CustomFieldModel

		data.CustomFields.ElementsAs(ctx, &stateCustomFields, false)

		customFields := utils.MapToCustomFieldModels(rt.CustomFields, stateCustomFields)

		customFieldsValue, _ := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, customFields)

		data.CustomFields = customFieldsValue

	case len(rt.CustomFields) > 0:

		customFields := utils.MapToCustomFieldModels(rt.CustomFields, []utils.CustomFieldModel{})

		customFieldsValue, _ := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, customFields)

		data.CustomFields = customFieldsValue

	default:

		data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)
	}
}
