// Package resources contains Terraform resource implementations for the Netbox provider.

package resources

import (
	"context"
	"fmt"
	"maps"

	"github.com/bab3l/go-netbox"
	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
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
	_ resource.Resource = &ServiceTemplateResource{}

	_ resource.ResourceWithConfigure = &ServiceTemplateResource{}

	_ resource.ResourceWithImportState = &ServiceTemplateResource{}
	_ resource.ResourceWithIdentity    = &ServiceTemplateResource{}
)

// NewServiceTemplateResource returns a new resource implementing the service template resource.

func NewServiceTemplateResource() resource.Resource {
	return &ServiceTemplateResource{}
}

// ServiceTemplateResource defines the resource implementation.

type ServiceTemplateResource struct {
	client *netbox.APIClient
}

// ServiceTemplateResourceModel describes the resource data model.

type ServiceTemplateResourceModel struct {
	ID types.String `tfsdk:"id"`

	Name types.String `tfsdk:"name"`

	Protocol types.String `tfsdk:"protocol"`

	Ports types.List `tfsdk:"ports"`

	Description types.String `tfsdk:"description"`

	Comments types.String `tfsdk:"comments"`

	Tags types.Set `tfsdk:"tags"`

	CustomFields types.Set `tfsdk:"custom_fields"`
}

// Metadata returns the resource type name.

func (r *ServiceTemplateResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service_template"
}

// Schema defines the schema for the resource.

func (r *ServiceTemplateResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a service template in NetBox. Service templates define reusable service configurations that can be applied to devices or virtual machines.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique numeric ID of the service template.",

				Computed: true,

				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the service template (e.g., 'ssh', 'http', 'https').",

				Required: true,
			},

			"protocol": schema.StringAttribute{
				MarkdownDescription: "The protocol used by the service. Valid values: `tcp`, `udp`, `sctp`. Defaults to `tcp` if not specified.",

				Optional: true,

				Computed: true,

				Validators: []validator.String{
					stringvalidator.OneOf("tcp", "udp", "sctp"),
				},
			},

			"ports": schema.ListAttribute{
				MarkdownDescription: "List of port numbers the service listens on.",

				Required: true,

				ElementType: types.Int64Type,
			},
		},
	}

	// Add description and comments attributes
	maps.Copy(resp.Schema.Attributes, nbschema.CommonDescriptiveAttributes("service template"))

	// Add metadata attributes (slug list tags, custom_fields)
	resp.Schema.Attributes["tags"] = nbschema.TagsSlugAttribute()
	resp.Schema.Attributes["custom_fields"] = nbschema.CustomFieldsAttribute()
}

func (r *ServiceTemplateResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = nbschema.ImportIdentityWithCustomFieldsSchema()
}

// Configure adds the provider configured client to the resource.

func (r *ServiceTemplateResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create creates a new service template.

func (r *ServiceTemplateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ServiceTemplateResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Convert ports to int32 slice

	var ports []int32

	if !data.Ports.IsNull() && !data.Ports.IsUnknown() {
		var portValues []int64

		diags := data.Ports.ElementsAs(ctx, &portValues, false)

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {
			return
		}

		for _, p := range portValues {
			p32, err := utils.SafeInt32(p)

			if err != nil {
				resp.Diagnostics.AddError("Invalid port number", fmt.Sprintf("Port number overflow: %s", err))

				return
			}

			ports = append(ports, p32)
		}
	}

	// Build the API request - Protocol is required by WritableServiceTemplateRequest

	protocol := netbox.PATCHEDWRITABLESERVICEREQUESTPROTOCOL_TCP // Default to TCP

	if !data.Protocol.IsNull() && !data.Protocol.IsUnknown() {
		protocol = netbox.PatchedWritableServiceRequestProtocol(data.Protocol.ValueString())
	}

	serviceTemplateRequest := netbox.NewWritableServiceTemplateRequest(data.Name.ValueString(), protocol, ports)

	// Store plan values before mapping for filter-to-owned pattern
	planTags := data.Tags
	planCustomFields := data.CustomFields

	// Apply description and comments
	utils.ApplyDescriptiveFields(serviceTemplateRequest, data.Description, data.Comments)

	// Apply tags from slugs
	utils.ApplyTagsFromSlugs(ctx, r.client, serviceTemplateRequest, data.Tags, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Apply custom fields
	utils.ApplyCustomFields(ctx, serviceTemplateRequest, data.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating service template", map[string]interface{}{
		"name": data.Name.ValueString(),

		"ports": ports,
	})

	// Call the API

	serviceTemplate, httpResp, err := r.client.IpamAPI.IpamServiceTemplatesCreate(ctx).
		WritableServiceTemplateRequest(*serviceTemplateRequest).
		Execute()

	if err != nil {
		resp.Diagnostics.AddError(

			"Error creating service template",

			utils.FormatAPIError("create service template", err, httpResp),
		)

		return
	}

	// Map response to state

	r.mapResponseToState(ctx, serviceTemplate, &data, &resp.Diagnostics)

	// Populate tags and custom fields filtered to owned fields only
	data.Tags = utils.PopulateTagsSlugFilteredToOwned(ctx, serviceTemplate.HasTags(), serviceTemplate.GetTags(), planTags)
	data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, planCustomFields, serviceTemplate.GetCustomFields(), &resp.Diagnostics)

	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(data.ID.ValueString()), data.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Created service template", map[string]interface{}{
		"id": serviceTemplate.GetId(),
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read reads the service template.

func (r *ServiceTemplateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ServiceTemplateResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Store state values before mapping for filter-to-owned pattern
	stateTags := data.Tags
	stateCustomFields := data.CustomFields

	// Parse ID

	id, err := utils.ParseID(data.ID.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(

			"Invalid ID format",

			fmt.Sprintf("Could not parse service template ID '%s': %s", data.ID.ValueString(), err),
		)

		return
	}

	tflog.Debug(ctx, "Reading service template", map[string]interface{}{
		"id": id,
	})

	// Call the API

	serviceTemplate, httpResp, err := r.client.IpamAPI.IpamServiceTemplatesRetrieve(ctx, id).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			resp.State.RemoveResource(ctx)

			return
		}

		resp.Diagnostics.AddError(

			"Error reading service template",

			utils.FormatAPIError("read service template", err, httpResp),
		)

		return
	}

	// Map response to state

	r.mapResponseToState(ctx, serviceTemplate, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	// Override with filter-to-owned pattern: only show fields that were in original state
	data.Tags = utils.PopulateTagsSlugFilteredToOwned(ctx, serviceTemplate.HasTags(), serviceTemplate.GetTags(), stateTags)
	data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, stateCustomFields, serviceTemplate.GetCustomFields(), &resp.Diagnostics)

	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(data.ID.ValueString()), data.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates the service template.

func (r *ServiceTemplateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, plan ServiceTemplateResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Store plan values before mapping for filter-to-owned pattern
	planTags := plan.Tags
	planCustomFields := plan.CustomFields

	// Parse ID

	id, err := utils.ParseID(plan.ID.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(

			"Invalid ID format",

			fmt.Sprintf("Could not parse service template ID '%s': %s", plan.ID.ValueString(), err),
		)

		return
	}

	// Convert ports to int32 slice

	var ports []int32

	if !plan.Ports.IsNull() && !plan.Ports.IsUnknown() {
		var portValues []int64

		diags := plan.Ports.ElementsAs(ctx, &portValues, false)

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {
			return
		}

		for _, p := range portValues {
			p32, err := utils.SafeInt32(p)

			if err != nil {
				resp.Diagnostics.AddError("Invalid port number", fmt.Sprintf("Port number overflow: %s", err))

				return
			}

			ports = append(ports, p32)
		}
	}

	// Build the API request - Protocol is required by WritableServiceTemplateRequest

	protocol := netbox.PATCHEDWRITABLESERVICEREQUESTPROTOCOL_TCP // Default to TCP

	if !plan.Protocol.IsNull() && !plan.Protocol.IsUnknown() {
		protocol = netbox.PatchedWritableServiceRequestProtocol(plan.Protocol.ValueString())
	}

	serviceTemplateRequest := netbox.NewWritableServiceTemplateRequest(plan.Name.ValueString(), protocol, ports)

	// Apply description and comments
	utils.ApplyDescriptiveFields(serviceTemplateRequest, plan.Description, plan.Comments)

	// Handle tags and custom fields - merge-aware for partial management
	// If tags are in plan, use plan. If not, preserve state tags.
	if utils.IsSet(plan.Tags) {
		utils.ApplyTagsFromSlugs(ctx, r.client, serviceTemplateRequest, plan.Tags, &resp.Diagnostics)
	} else if utils.IsSet(state.Tags) {
		utils.ApplyTagsFromSlugs(ctx, r.client, serviceTemplateRequest, state.Tags, &resp.Diagnostics)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	// Apply custom fields with merge logic (preserves unmanaged fields from state)
	utils.ApplyCustomFieldsWithMerge(ctx, serviceTemplateRequest, plan.CustomFields, state.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Updating service template", map[string]interface{}{
		"id": id,

		"name": plan.Name.ValueString(),
	})

	// Call the API

	serviceTemplate, httpResp, err := r.client.IpamAPI.IpamServiceTemplatesUpdate(ctx, id).
		WritableServiceTemplateRequest(*serviceTemplateRequest).
		Execute()

	if err != nil {
		resp.Diagnostics.AddError(

			"Error updating service template",

			utils.FormatAPIError("update service template", err, httpResp),
		)

		return
	}

	// Map response to state

	r.mapResponseToState(ctx, serviceTemplate, &plan, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	// Override with filter-to-owned pattern
	plan.Tags = utils.PopulateTagsSlugFilteredToOwned(ctx, serviceTemplate.HasTags(), serviceTemplate.GetTags(), planTags)
	plan.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, planCustomFields, serviceTemplate.GetCustomFields(), &resp.Diagnostics)

	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(plan.ID.ValueString()), plan.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Updated service template", map[string]interface{}{
		"id": serviceTemplate.GetId(),
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete deletes the service template.

func (r *ServiceTemplateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ServiceTemplateResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Parse ID

	id, err := utils.ParseID(data.ID.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(

			"Invalid ID format",

			fmt.Sprintf("Could not parse service template ID '%s': %s", data.ID.ValueString(), err),
		)

		return
	}

	tflog.Debug(ctx, "Deleting service template", map[string]interface{}{
		"id": id,
	})

	httpResp, err := r.client.IpamAPI.IpamServiceTemplatesDestroy(ctx, id).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			return
		}

		resp.Diagnostics.AddError(

			"Error deleting service template",

			utils.FormatAPIError("delete service template", err, httpResp),
		)

		return
	}

	tflog.Debug(ctx, "Deleted service template", map[string]interface{}{
		"id": id,
	})
}

// ImportState imports an existing service template.

func (r *ServiceTemplateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if parsed, ok := utils.ParseImportIdentityCustomFields(ctx, req.Identity, &resp.Diagnostics); ok {
		if resp.Diagnostics.HasError() {
			return
		}
		if parsed.ID == "" {
			resp.Diagnostics.AddError("Invalid import identity", "Identity id must be provided")
			return
		}

		id, err := utils.ParseID(parsed.ID)
		if err != nil {
			resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Could not parse service template ID: %s", err))
			return
		}

		serviceTemplate, httpResp, err := r.client.IpamAPI.IpamServiceTemplatesRetrieve(ctx, id).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError("Error importing service template", utils.FormatAPIError("read service template", err, httpResp))
			return
		}

		var data ServiceTemplateResourceModel
		data.Tags = utils.PopulateTagsSlugFromAPI(ctx, serviceTemplate.HasTags(), serviceTemplate.GetTags(), data.Tags)
		if parsed.HasCustomFields {
			if len(parsed.CustomFields) == 0 {
				data.CustomFields = types.SetValueMust(utils.GetCustomFieldsAttributeType().ElemType, []attr.Value{})
			} else {
				ownedSet, setDiags := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, parsed.CustomFields)
				resp.Diagnostics.Append(setDiags...)
				if resp.Diagnostics.HasError() {
					return
				}
				data.CustomFields = ownedSet
			}
		} else {
			data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)
		}

		r.mapResponseToState(ctx, serviceTemplate, &data, &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}
		if parsed.HasCustomFields {
			data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, data.CustomFields, serviceTemplate.GetCustomFields(), &resp.Diagnostics)
		} else {
			data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)
		}
		if resp.Diagnostics.HasError() {
			return
		}

		if resp.Identity != nil {
			listValue, listDiags := types.ListValueFrom(ctx, types.StringType, parsed.CustomFieldItems)
			resp.Diagnostics.Append(listDiags...)
			if resp.Diagnostics.HasError() {
				return
			}
			resp.Diagnostics.Append(resp.Identity.Set(ctx, &utils.ImportIdentityCustomFieldsModel{
				ID:           types.StringValue(parsed.ID),
				CustomFields: listValue,
			})...)
		}

		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
		return
	}

	utils.ImportStatePassthroughIDWithValidation(ctx, req, resp, path.Root("id"), true)
}

// mapResponseToState maps the API response to the Terraform state.

func (r *ServiceTemplateResource) mapResponseToState(ctx context.Context, serviceTemplate *netbox.ServiceTemplate, data *ServiceTemplateResourceModel, diags *diag.Diagnostics) {
	data.ID = types.StringValue(fmt.Sprintf("%d", serviceTemplate.GetId()))

	data.Name = types.StringValue(serviceTemplate.GetName())

	// Handle protocol

	if serviceTemplate.HasProtocol() {
		protocol := serviceTemplate.GetProtocol()

		data.Protocol = types.StringValue(string(protocol.GetValue()))
	} else {
		data.Protocol = types.StringNull()
	}

	// Handle ports

	if serviceTemplate.Ports != nil {
		var ports []int64

		for _, p := range serviceTemplate.Ports {
			ports = append(ports, int64(p))
		}

		portsList, d := types.ListValueFrom(ctx, types.Int64Type, ports)

		diags.Append(d...)

		data.Ports = portsList
	} else {
		data.Ports = types.ListNull(types.Int64Type)
	}

	// Handle description

	if serviceTemplate.HasDescription() && serviceTemplate.GetDescription() != "" {
		data.Description = types.StringValue(serviceTemplate.GetDescription())
	} else {
		data.Description = types.StringNull()
	}

	// Handle comments

	if serviceTemplate.HasComments() && serviceTemplate.GetComments() != "" {
		data.Comments = types.StringValue(serviceTemplate.GetComments())
	} else {
		data.Comments = types.StringNull()
	}

	// Handle tags using consolidated helper
	data.Tags = utils.PopulateTagsSlugFilteredToOwned(ctx, serviceTemplate.HasTags(), serviceTemplate.GetTags(), data.Tags)
	if diags.HasError() {
		return
	}

	// Handle custom fields using consolidated helper
	data.CustomFields = utils.PopulateCustomFieldsFromAPI(ctx, serviceTemplate.HasCustomFields(), serviceTemplate.GetCustomFields(), data.CustomFields, diags)
}
