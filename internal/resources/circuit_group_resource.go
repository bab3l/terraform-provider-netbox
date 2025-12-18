// Package resources contains Terraform resource implementations for the Netbox provider.

package resources

import (
	"context"
	"fmt"

	"github.com/bab3l/go-netbox"
	"github.com/bab3l/terraform-provider-netbox/internal/netboxlookup"
	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.

var (
	_ resource.Resource = &CircuitGroupResource{}

	_ resource.ResourceWithConfigure = &CircuitGroupResource{}

	_ resource.ResourceWithImportState = &CircuitGroupResource{}
)

// NewCircuitGroupResource returns a new circuit group resource.

func NewCircuitGroupResource() resource.Resource {

	return &CircuitGroupResource{}

}

// CircuitGroupResource defines the circuit group resource implementation.

type CircuitGroupResource struct {
	client *netbox.APIClient
}

// CircuitGroupResourceModel describes the circuit group resource data model.

type CircuitGroupResourceModel struct {
	ID types.String `tfsdk:"id"`

	Name types.String `tfsdk:"name"`

	Slug types.String `tfsdk:"slug"`

	Description types.String `tfsdk:"description"`

	Tenant types.String `tfsdk:"tenant"`

	TenantID types.String `tfsdk:"tenant_id"`

	Tags types.Set `tfsdk:"tags"`

	CustomFields types.Set `tfsdk:"custom_fields"`
}

// Metadata returns the resource type name.

func (r *CircuitGroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {

	resp.TypeName = req.ProviderTypeName + "_circuit_group"

}

// Schema defines the schema for the resource.

func (r *CircuitGroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	resp.Schema = schema.Schema{

		MarkdownDescription: "Manages a circuit group in Netbox. Circuit groups allow you to organize related circuits together for management and reporting purposes.",

		Attributes: map[string]schema.Attribute{

			"id": nbschema.IDAttribute("circuit group"),

			"name": nbschema.NameAttribute("circuit group", 100),

			"slug": nbschema.SlugAttribute("circuit group"),

			"description": nbschema.DescriptionAttribute("circuit group"),

			"tenant": nbschema.ReferenceAttribute("tenant", "ID or slug of the tenant."),

			"tenant_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The numeric ID of the tenant.",
			},

			"tags": nbschema.TagsAttribute(),

			"custom_fields": nbschema.CustomFieldsAttribute(),
		},
	}

}

// Configure adds the provider configured client to the resource.

func (r *CircuitGroupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {

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

func (r *CircuitGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var data CircuitGroupResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	// Build the API request

	groupRequest := netbox.NewCircuitGroupRequest(

		data.Name.ValueString(),

		data.Slug.ValueString(),
	)

	// Set optional fields

	if !data.Description.IsNull() && data.Description.ValueString() != "" {

		groupRequest.Description = netbox.PtrString(data.Description.ValueString())

	}

	// Handle tenant

	if !data.Tenant.IsNull() && data.Tenant.ValueString() != "" {

		tenant, tenantDiags := netboxlookup.LookupTenant(ctx, r.client, data.Tenant.ValueString())

		resp.Diagnostics.Append(tenantDiags...)

		if resp.Diagnostics.HasError() {

			return

		}

		groupRequest.Tenant = *netbox.NewNullableBriefTenantRequest(tenant)

	}

	// Handle tags

	if !data.Tags.IsNull() {

		var tagModels []utils.TagModel

		diags := data.Tags.ElementsAs(ctx, &tagModels, false)

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {

			return

		}

		groupRequest.Tags = utils.TagsToNestedTagRequests(tagModels)

	}

	// Handle custom fields

	if !data.CustomFields.IsNull() {

		var customFieldModels []utils.CustomFieldModel

		diags := data.CustomFields.ElementsAs(ctx, &customFieldModels, false)

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {

			return

		}

		groupRequest.CustomFields = utils.CustomFieldsToMap(customFieldModels)

	}

	tflog.Debug(ctx, "Creating circuit group", map[string]interface{}{

		"name": data.Name.ValueString(),

		"slug": data.Slug.ValueString(),
	})

	// Call the API

	group, httpResp, err := r.client.CircuitsAPI.CircuitsCircuitGroupsCreate(ctx).CircuitGroupRequest(*groupRequest).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		resp.Diagnostics.AddError(

			"Error creating circuit group",

			utils.FormatAPIError("create circuit group", err, httpResp),
		)

		return

	}

	// Map response to state

	r.mapResponseToState(ctx, group, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {

		return

	}

	tflog.Trace(ctx, "Created circuit group resource", map[string]interface{}{

		"id": data.ID.ValueString(),
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

// Read refreshes the Terraform state with the latest data.

func (r *CircuitGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	var data CircuitGroupResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	// Parse ID

	var idInt int32

	if _, err := fmt.Sscanf(data.ID.ValueString(), "%d", &idInt); err != nil {

		resp.Diagnostics.AddError(

			"Invalid ID format",

			fmt.Sprintf("Could not parse circuit group ID '%s': %s", data.ID.ValueString(), err.Error()),
		)

		return

	}

	// Read from API

	group, httpResp, err := r.client.CircuitsAPI.CircuitsCircuitGroupsRetrieve(ctx, idInt).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		if httpResp != nil && httpResp.StatusCode == 404 {

			resp.State.RemoveResource(ctx)

			return

		}

		resp.Diagnostics.AddError(

			"Error reading circuit group",

			utils.FormatAPIError("read circuit group", err, httpResp),
		)

		return

	}

	// Map response to state

	r.mapResponseToState(ctx, group, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {

		return

	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

// Update updates the resource and sets the updated Terraform state on success.

func (r *CircuitGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var data CircuitGroupResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	// Parse ID

	var idInt int32

	if _, err := fmt.Sscanf(data.ID.ValueString(), "%d", &idInt); err != nil {

		resp.Diagnostics.AddError(

			"Invalid ID format",

			fmt.Sprintf("Could not parse circuit group ID '%s': %s", data.ID.ValueString(), err.Error()),
		)

		return

	}

	// Build the API request

	groupRequest := netbox.NewCircuitGroupRequest(

		data.Name.ValueString(),

		data.Slug.ValueString(),
	)

	// Set optional fields

	if !data.Description.IsNull() && data.Description.ValueString() != "" {

		groupRequest.Description = netbox.PtrString(data.Description.ValueString())

	} else {

		groupRequest.Description = netbox.PtrString("")

	}

	// Handle tenant

	if !data.Tenant.IsNull() && data.Tenant.ValueString() != "" {

		tenant, tenantDiags := netboxlookup.LookupTenant(ctx, r.client, data.Tenant.ValueString())

		resp.Diagnostics.Append(tenantDiags...)

		if resp.Diagnostics.HasError() {

			return

		}

		groupRequest.Tenant = *netbox.NewNullableBriefTenantRequest(tenant)

	} else {

		groupRequest.Tenant = *netbox.NewNullableBriefTenantRequest(nil)

	}

	// Handle tags

	if !data.Tags.IsNull() {

		var tagModels []utils.TagModel

		diags := data.Tags.ElementsAs(ctx, &tagModels, false)

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {

			return

		}

		groupRequest.Tags = utils.TagsToNestedTagRequests(tagModels)

	} else {

		groupRequest.Tags = []netbox.NestedTagRequest{}

	}

	// Handle custom fields

	if !data.CustomFields.IsNull() {

		var customFieldModels []utils.CustomFieldModel

		diags := data.CustomFields.ElementsAs(ctx, &customFieldModels, false)

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {

			return

		}

		groupRequest.CustomFields = utils.CustomFieldsToMap(customFieldModels)

	}

	tflog.Debug(ctx, "Updating circuit group", map[string]interface{}{

		"id": idInt,

		"name": data.Name.ValueString(),
	})

	// Call the API

	group, httpResp, err := r.client.CircuitsAPI.CircuitsCircuitGroupsUpdate(ctx, idInt).CircuitGroupRequest(*groupRequest).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		resp.Diagnostics.AddError(

			"Error updating circuit group",

			utils.FormatAPIError("update circuit group", err, httpResp),
		)

		return

	}

	// Map response to state

	r.mapResponseToState(ctx, group, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {

		return

	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

// Delete deletes the resource and removes the Terraform state on success.

func (r *CircuitGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var data CircuitGroupResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	// Parse ID

	var idInt int32

	if _, err := fmt.Sscanf(data.ID.ValueString(), "%d", &idInt); err != nil {

		resp.Diagnostics.AddError(

			"Invalid ID format",

			fmt.Sprintf("Could not parse circuit group ID '%s': %s", data.ID.ValueString(), err.Error()),
		)

		return

	}

	tflog.Debug(ctx, "Deleting circuit group", map[string]interface{}{

		"id": idInt,
	})

	// Call the API

	httpResp, err := r.client.CircuitsAPI.CircuitsCircuitGroupsDestroy(ctx, idInt).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		if httpResp != nil && httpResp.StatusCode == 404 {

			// Already deleted

			return

		}

		resp.Diagnostics.AddError(

			"Error deleting circuit group",

			utils.FormatAPIError("delete circuit group", err, httpResp),
		)

		return

	}

}

// ImportState imports the resource state from an existing Netbox object.

func (r *CircuitGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

}

// mapResponseToState maps a CircuitGroup API response to the Terraform state model.

func (r *CircuitGroupResource) mapResponseToState(ctx context.Context, group *netbox.CircuitGroup, data *CircuitGroupResourceModel, diags *diag.Diagnostics) {

	data.ID = types.StringValue(fmt.Sprintf("%d", group.GetId()))

	data.Name = types.StringValue(group.GetName())

	data.Slug = types.StringValue(group.GetSlug())

	// Description

	if group.HasDescription() && group.GetDescription() != "" {

		data.Description = types.StringValue(group.GetDescription())

	} else {

		data.Description = types.StringNull()

	}

	// Tenant

	if group.HasTenant() && group.Tenant.IsSet() && group.Tenant.Get() != nil {

		tenant := group.Tenant.Get()

		data.TenantID = types.StringValue(fmt.Sprintf("%d", tenant.GetId()))

		userTenant := data.Tenant.ValueString()

		if userTenant == tenant.GetName() || userTenant == tenant.GetSlug() || userTenant == tenant.GetDisplay() || userTenant == fmt.Sprintf("%d", tenant.GetId()) {

			// Keep user's original value

		} else {

			data.Tenant = types.StringValue(tenant.GetName())

		}

	} else {

		data.Tenant = types.StringNull()
		data.TenantID = types.StringNull()

	}

	// Tags

	if group.HasTags() && len(group.GetTags()) > 0 {

		tags := utils.NestedTagsToTagModels(group.GetTags())

		tagsValue, d := types.SetValueFrom(ctx, utils.GetTagsAttributeType().ElemType, tags)

		diags.Append(d...)

		data.Tags = tagsValue

	} else {

		data.Tags = types.SetNull(utils.GetTagsAttributeType().ElemType)

	}

	// Custom fields

	if group.HasCustomFields() && len(group.GetCustomFields()) > 0 {

		var existingModels []utils.CustomFieldModel

		if !data.CustomFields.IsNull() {

			d := data.CustomFields.ElementsAs(ctx, &existingModels, false)

			diags.Append(d...)

		}

		customFields := utils.MapToCustomFieldModels(group.GetCustomFields(), existingModels)

		customFieldsValue, d := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, customFields)

		diags.Append(d...)

		data.CustomFields = customFieldsValue

	} else {

		data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)

	}

}
