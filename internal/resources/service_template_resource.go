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

	DisplayName types.String `tfsdk:"display_name"`

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

				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},

				Validators: []validator.String{
					stringvalidator.OneOf("tcp", "udp", "sctp"),
				},
			},

			"ports": schema.ListAttribute{
				MarkdownDescription: "List of port numbers the service listens on.",

				Required: true,

				ElementType: types.Int64Type,
			},

			"display_name": nbschema.DisplayNameAttribute("service template"),
		},
	}

	// Add description and comments attributes
	maps.Copy(resp.Schema.Attributes, nbschema.CommonDescriptiveAttributes("service template"))

	// Add common metadata attributes (tags, custom_fields)
	maps.Copy(resp.Schema.Attributes, nbschema.CommonMetadataAttributes())
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

	// Apply common fields (description, comments, tags, custom_fields)
	utils.ApplyCommonFields(ctx, serviceTemplateRequest, data.Description, data.Comments, data.Tags, data.CustomFields, &resp.Diagnostics)
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

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates the service template.

func (r *ServiceTemplateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ServiceTemplateResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

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

	// Apply common fields (description, comments, tags, custom_fields)
	utils.ApplyCommonFields(ctx, serviceTemplateRequest, data.Description, data.Comments, data.Tags, data.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Updating service template", map[string]interface{}{
		"id": id,

		"name": data.Name.ValueString(),
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

	r.mapResponseToState(ctx, serviceTemplate, &data, &resp.Diagnostics)

	tflog.Debug(ctx, "Updated service template", map[string]interface{}{
		"id": serviceTemplate.GetId(),
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
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
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
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

	// Map display_name
	if serviceTemplate.Display != "" {
		data.DisplayName = types.StringValue(serviceTemplate.Display)
	} else {
		data.DisplayName = types.StringNull()
	}

	// Handle tags

	if serviceTemplate.HasTags() && len(serviceTemplate.GetTags()) > 0 {
		tags := utils.NestedTagsToTagModels(serviceTemplate.GetTags())

		tagsValue, d := types.SetValueFrom(ctx, utils.GetTagsAttributeType().ElemType, tags)

		diags.Append(d...)

		data.Tags = tagsValue
	} else {
		data.Tags = types.SetNull(utils.GetTagsAttributeType().ElemType)
	}

	// Handle custom fields

	if serviceTemplate.HasCustomFields() {
		var existingModels []utils.CustomFieldModel

		if !data.CustomFields.IsNull() {
			d := data.CustomFields.ElementsAs(ctx, &existingModels, false)

			diags.Append(d...)
		}

		customFields := utils.MapToCustomFieldModels(serviceTemplate.GetCustomFields(), existingModels)

		if len(customFields) > 0 {
			customFieldsValue, d := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, customFields)

			diags.Append(d...)

			data.CustomFields = customFieldsValue
		} else {
			data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)
		}
	} else {
		data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)
	}
}
