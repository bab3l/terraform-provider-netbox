// Package resources provides Terraform resource implementations for NetBox objects.

package resources

import (
	"context"
	"fmt"

	"github.com/bab3l/go-netbox"
	lookup "github.com/bab3l/terraform-provider-netbox/internal/netboxlookup"
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
	_ resource.Resource = &ModuleBayTemplateResource{}

	_ resource.ResourceWithConfigure = &ModuleBayTemplateResource{}

	_ resource.ResourceWithImportState = &ModuleBayTemplateResource{}
)

// NewModuleBayTemplateResource returns a new resource implementing the module bay template resource.

func NewModuleBayTemplateResource() resource.Resource {

	return &ModuleBayTemplateResource{}

}

// ModuleBayTemplateResource defines the resource implementation.

type ModuleBayTemplateResource struct {
	client *netbox.APIClient
}

// ModuleBayTemplateResourceModel describes the resource data model.

type ModuleBayTemplateResourceModel struct {
	ID types.String `tfsdk:"id"`

	DeviceType types.String `tfsdk:"device_type"`

	ModuleType types.String `tfsdk:"module_type"`

	Name types.String `tfsdk:"name"`

	Label types.String `tfsdk:"label"`

	Position types.String `tfsdk:"position"`

	Description types.String `tfsdk:"description"`
}

// Metadata returns the resource type name.

func (r *ModuleBayTemplateResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {

	resp.TypeName = req.ProviderTypeName + "_module_bay_template"

}

// Schema defines the schema for the resource.

func (r *ModuleBayTemplateResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	resp.Schema = schema.Schema{

		MarkdownDescription: "Manages a module bay template in NetBox. Module bay templates define module bay configurations for device types or module types.",

		Attributes: map[string]schema.Attribute{

			"id": schema.StringAttribute{

				MarkdownDescription: "The unique numeric ID of the module bay template.",

				Computed: true,

				PlanModifiers: []planmodifier.String{

					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"device_type": schema.StringAttribute{

				MarkdownDescription: "The device type this module bay template belongs to (ID or model name). Either device_type or module_type is required.",

				Optional: true,
			},

			"module_type": schema.StringAttribute{

				MarkdownDescription: "The module type this module bay template belongs to (ID or model name). Either device_type or module_type is required.",

				Optional: true,
			},

			"name": schema.StringAttribute{

				MarkdownDescription: "The name of the module bay template. {module} is accepted as a substitution for the module bay position when attached to a module type.",

				Required: true,
			},

			"label": schema.StringAttribute{

				MarkdownDescription: "Physical label of the module bay template.",

				Optional: true,
			},

			"position": schema.StringAttribute{

				MarkdownDescription: "Identifier to reference when renaming installed components.",

				Optional: true,
			},

			"description": schema.StringAttribute{

				MarkdownDescription: "A description of the module bay template.",

				Optional: true,
			},
		},
	}

}

// Configure adds the provider configured client to the resource.

func (r *ModuleBayTemplateResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {

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

// Create creates the resource.

func (r *ModuleBayTemplateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var data ModuleBayTemplateResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	// Validate that at least one of device_type or module_type is set

	if data.DeviceType.IsNull() && data.ModuleType.IsNull() {

		resp.Diagnostics.AddError(

			"Missing Required Attribute",

			"Either device_type or module_type must be specified.",
		)

		return

	}

	// Build request

	apiReq := netbox.NewModuleBayTemplateRequest(data.Name.ValueString())

	// Set device type or module type

	if !data.DeviceType.IsNull() && !data.DeviceType.IsUnknown() {

		deviceType, diags := lookup.LookupDeviceType(ctx, r.client, data.DeviceType.ValueString())

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {

			return

		}

		apiReq.SetDeviceType(*deviceType)

	}

	if !data.ModuleType.IsNull() && !data.ModuleType.IsUnknown() {

		moduleType, diags := lookup.LookupModuleType(ctx, r.client, data.ModuleType.ValueString())

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {

			return

		}

		apiReq.SetModuleType(*moduleType)

	}

	// Set optional fields

	if !data.Label.IsNull() && !data.Label.IsUnknown() {

		apiReq.SetLabel(data.Label.ValueString())

	}

	if !data.Position.IsNull() && !data.Position.IsUnknown() {

		apiReq.SetPosition(data.Position.ValueString())

	}

	if !data.Description.IsNull() && !data.Description.IsUnknown() {

		apiReq.SetDescription(data.Description.ValueString())

	}

	tflog.Debug(ctx, "Creating module bay template", map[string]interface{}{

		"name": data.Name.ValueString(),

		"device_type": data.DeviceType.ValueString(),

		"module_type": data.ModuleType.ValueString(),
	})

	// Create the resource

	result, httpResp, err := r.client.DcimAPI.DcimModuleBayTemplatesCreate(ctx).ModuleBayTemplateRequest(*apiReq).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		resp.Diagnostics.AddError(

			"Error creating module bay template",

			utils.FormatAPIError("create module bay template", err, httpResp),
		)

		return

	}

	// Map response to state

	r.mapToState(ctx, result, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {

		return

	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

// Read reads the resource.

func (r *ModuleBayTemplateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	var data ModuleBayTemplateResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	// Parse ID

	var id int32

	_, err := fmt.Sscanf(data.ID.ValueString(), "%d", &id)

	if err != nil {

		resp.Diagnostics.AddError(

			"Error parsing module bay template ID",

			fmt.Sprintf("Could not parse ID '%s': %s", data.ID.ValueString(), err.Error()),
		)

		return

	}

	// Read from API

	result, httpResp, err := r.client.DcimAPI.DcimModuleBayTemplatesRetrieve(ctx, id).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		if httpResp != nil && httpResp.StatusCode == 404 {

			resp.State.RemoveResource(ctx)

			return

		}

		resp.Diagnostics.AddError(

			"Error reading module bay template",

			utils.FormatAPIError(fmt.Sprintf("read module bay template ID %d", id), err, httpResp),
		)

		return

	}

	// Map response to state

	r.mapToState(ctx, result, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {

		return

	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

// Update updates the resource.

func (r *ModuleBayTemplateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var data ModuleBayTemplateResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	// Parse ID

	var id int32

	_, err := fmt.Sscanf(data.ID.ValueString(), "%d", &id)

	if err != nil {

		resp.Diagnostics.AddError(

			"Error parsing module bay template ID",

			fmt.Sprintf("Could not parse ID '%s': %s", data.ID.ValueString(), err.Error()),
		)

		return

	}

	// Validate that at least one of device_type or module_type is set

	if data.DeviceType.IsNull() && data.ModuleType.IsNull() {

		resp.Diagnostics.AddError(

			"Missing Required Attribute",

			"Either device_type or module_type must be specified.",
		)

		return

	}

	// Build request

	apiReq := netbox.NewModuleBayTemplateRequest(data.Name.ValueString())

	// Set device type or module type

	if !data.DeviceType.IsNull() && !data.DeviceType.IsUnknown() {

		deviceType, diags := lookup.LookupDeviceType(ctx, r.client, data.DeviceType.ValueString())

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {

			return

		}

		apiReq.SetDeviceType(*deviceType)

	}

	if !data.ModuleType.IsNull() && !data.ModuleType.IsUnknown() {

		moduleType, diags := lookup.LookupModuleType(ctx, r.client, data.ModuleType.ValueString())

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {

			return

		}

		apiReq.SetModuleType(*moduleType)

	}

	// Set optional fields

	if !data.Label.IsNull() && !data.Label.IsUnknown() {

		apiReq.SetLabel(data.Label.ValueString())

	}

	if !data.Position.IsNull() && !data.Position.IsUnknown() {

		apiReq.SetPosition(data.Position.ValueString())

	}

	if !data.Description.IsNull() && !data.Description.IsUnknown() {

		apiReq.SetDescription(data.Description.ValueString())

	}

	tflog.Debug(ctx, "Updating module bay template", map[string]interface{}{

		"id": id,

		"name": data.Name.ValueString(),

		"device_type": data.DeviceType.ValueString(),

		"module_type": data.ModuleType.ValueString(),
	})

	// Update the resource

	result, httpResp, err := r.client.DcimAPI.DcimModuleBayTemplatesUpdate(ctx, id).ModuleBayTemplateRequest(*apiReq).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		resp.Diagnostics.AddError(

			"Error updating module bay template",

			utils.FormatAPIError(fmt.Sprintf("update module bay template ID %d", id), err, httpResp),
		)

		return

	}

	// Map response to state

	r.mapToState(ctx, result, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {

		return

	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

// Delete deletes the resource.

func (r *ModuleBayTemplateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var data ModuleBayTemplateResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	// Parse ID

	var id int32

	_, err := fmt.Sscanf(data.ID.ValueString(), "%d", &id)

	if err != nil {

		resp.Diagnostics.AddError(

			"Error parsing module bay template ID",

			fmt.Sprintf("Could not parse ID '%s': %s", data.ID.ValueString(), err.Error()),
		)

		return

	}

	tflog.Debug(ctx, "Deleting module bay template", map[string]interface{}{"id": id})

	// Delete the resource

	httpResp, err := r.client.DcimAPI.DcimModuleBayTemplatesDestroy(ctx, id).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		if httpResp != nil && httpResp.StatusCode == 404 {

			return

		}

		resp.Diagnostics.AddError(

			"Error deleting module bay template",

			utils.FormatAPIError(fmt.Sprintf("delete module bay template ID %d", id), err, httpResp),
		)

		return

	}

}

// ImportState imports the resource state.

func (r *ModuleBayTemplateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

}

// mapToState maps the API response to the Terraform state.

func (r *ModuleBayTemplateResource) mapToState(ctx context.Context, result *netbox.ModuleBayTemplate, data *ModuleBayTemplateResourceModel, diags *diag.Diagnostics) {

	data.ID = types.StringValue(fmt.Sprintf("%d", result.GetId()))

	data.Name = types.StringValue(result.GetName())

	// Map device type

	if result.HasDeviceType() && result.GetDeviceType().Id != 0 {

		deviceType := result.GetDeviceType()

		data.DeviceType = types.StringValue(fmt.Sprintf("%d", deviceType.GetId()))

	} else {

		data.DeviceType = types.StringNull()

	}

	// Map module type

	if result.HasModuleType() && result.GetModuleType().Id != 0 {

		moduleType := result.GetModuleType()

		data.ModuleType = types.StringValue(fmt.Sprintf("%d", moduleType.GetId()))

	} else {

		data.ModuleType = types.StringNull()

	}

	// Map label

	if result.HasLabel() && result.GetLabel() != "" {

		data.Label = types.StringValue(result.GetLabel())

	} else {

		data.Label = types.StringNull()

	}

	// Map position

	if result.HasPosition() && result.GetPosition() != "" {

		data.Position = types.StringValue(result.GetPosition())

	} else {

		data.Position = types.StringNull()

	}

	// Map description

	if result.HasDescription() && result.GetDescription() != "" {

		data.Description = types.StringValue(result.GetDescription())

	} else {

		data.Description = types.StringNull()

	}

}
