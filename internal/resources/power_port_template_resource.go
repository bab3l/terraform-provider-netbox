// Package resources provides Terraform resource implementations for NetBox objects.
package resources

import (
	"context"
	"fmt"

	"github.com/bab3l/go-netbox"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/bab3l/terraform-provider-netbox/internal/netboxlookup"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &PowerPortTemplateResource{}
var _ resource.ResourceWithImportState = &PowerPortTemplateResource{}

// NewPowerPortTemplateResource returns a new resource implementing the power port template resource.
func NewPowerPortTemplateResource() resource.Resource {
	return &PowerPortTemplateResource{}
}

// PowerPortTemplateResource defines the resource implementation.
type PowerPortTemplateResource struct {
	client *netbox.APIClient
}

// PowerPortTemplateResourceModel describes the resource data model.
type PowerPortTemplateResourceModel struct {
	ID            types.Int32  `tfsdk:"id"`
	DeviceType    types.String `tfsdk:"device_type"`
	ModuleType    types.String `tfsdk:"module_type"`
	Name          types.String `tfsdk:"name"`
	Label         types.String `tfsdk:"label"`
	Type          types.String `tfsdk:"type"`
	MaximumDraw   types.Int32  `tfsdk:"maximum_draw"`
	AllocatedDraw types.Int32  `tfsdk:"allocated_draw"`
	Description   types.String `tfsdk:"description"`
}

// Metadata returns the resource type name.
func (r *PowerPortTemplateResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_power_port_template"
}

// Schema defines the schema for the resource.
func (r *PowerPortTemplateResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a power port template in NetBox. Power port templates define the default power ports that will be created on new devices of a specific device type or modules of a module type.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int32Attribute{
				MarkdownDescription: "The unique numeric ID of the power port template.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.UseStateForUnknown(),
				},
			},
			"device_type": schema.StringAttribute{
				MarkdownDescription: "The device type ID or slug. Either device_type or module_type must be specified.",
				Optional:            true,
			},
			"module_type": schema.StringAttribute{
				MarkdownDescription: "The module type ID or model name. Either device_type or module_type must be specified.",
				Optional:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the power port template. Use {module} as a substitution for the module bay position when attached to a module type.",
				Required:            true,
			},
			"label": schema.StringAttribute{
				MarkdownDescription: "Physical label of the power port template.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "The type of power port (e.g., iec-60320-c6, iec-60320-c8, iec-60320-c14, nema-1-15p, nema-5-15p, nema-5-20p, etc.).",
				Optional:            true,
			},
			"maximum_draw": schema.Int32Attribute{
				MarkdownDescription: "Maximum power draw in watts.",
				Optional:            true,
			},
			"allocated_draw": schema.Int32Attribute{
				MarkdownDescription: "Allocated power draw in watts.",
				Optional:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "A description of the power port template.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *PowerPortTemplateResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *PowerPortTemplateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data PowerPortTemplateResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the API request
	apiReq := netbox.NewWritablePowerPortTemplateRequest(data.Name.ValueString())

	// Set device type or module type
	if !data.DeviceType.IsNull() && !data.DeviceType.IsUnknown() {
		deviceType, diags := netboxlookup.LookupDeviceType(ctx, r.client, data.DeviceType.ValueString())
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		apiReq.SetDeviceType(*deviceType)
	}
	if !data.ModuleType.IsNull() && !data.ModuleType.IsUnknown() {
		moduleType, diags := netboxlookup.LookupModuleType(ctx, r.client, data.ModuleType.ValueString())
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
	if !data.Type.IsNull() && !data.Type.IsUnknown() {
		apiReq.SetType(netbox.PatchedWritablePowerPortTemplateRequestType(data.Type.ValueString()))
	}
	if !data.MaximumDraw.IsNull() && !data.MaximumDraw.IsUnknown() {
		apiReq.SetMaximumDraw(data.MaximumDraw.ValueInt32())
	}
	if !data.AllocatedDraw.IsNull() && !data.AllocatedDraw.IsUnknown() {
		apiReq.SetAllocatedDraw(data.AllocatedDraw.ValueInt32())
	}
	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		apiReq.SetDescription(data.Description.ValueString())
	}

	tflog.Debug(ctx, "Creating power port template", map[string]interface{}{
		"name": data.Name.ValueString(),
	})

	response, httpResp, err := r.client.DcimAPI.DcimPowerPortTemplatesCreate(ctx).WritablePowerPortTemplateRequest(*apiReq).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating power port template",
			utils.FormatAPIError("create power port template", err, httpResp),
		)
		return
	}

	// Map response to model
	r.mapResponseToModel(response, &data)

	tflog.Trace(ctx, "Created power port template", map[string]interface{}{
		"id": data.ID.ValueInt32(),
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *PowerPortTemplateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data PowerPortTemplateResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	templateID := data.ID.ValueInt32()

	tflog.Debug(ctx, "Reading power port template", map[string]interface{}{
		"id": templateID,
	})

	response, httpResp, err := r.client.DcimAPI.DcimPowerPortTemplatesRetrieve(ctx, templateID).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			tflog.Debug(ctx, "Power port template not found, removing from state", map[string]interface{}{
				"id": templateID,
			})
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading power port template",
			utils.FormatAPIError(fmt.Sprintf("read power port template ID %d", templateID), err, httpResp),
		)
		return
	}

	// Map response to model
	r.mapResponseToModel(response, &data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *PowerPortTemplateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data PowerPortTemplateResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	templateID := data.ID.ValueInt32()

	// Build the API request
	apiReq := netbox.NewWritablePowerPortTemplateRequest(data.Name.ValueString())

	// Set device type or module type
	if !data.DeviceType.IsNull() && !data.DeviceType.IsUnknown() {
		deviceType, diags := netboxlookup.LookupDeviceType(ctx, r.client, data.DeviceType.ValueString())
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		apiReq.SetDeviceType(*deviceType)
	}
	if !data.ModuleType.IsNull() && !data.ModuleType.IsUnknown() {
		moduleType, diags := netboxlookup.LookupModuleType(ctx, r.client, data.ModuleType.ValueString())
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
	if !data.Type.IsNull() && !data.Type.IsUnknown() {
		apiReq.SetType(netbox.PatchedWritablePowerPortTemplateRequestType(data.Type.ValueString()))
	}
	if !data.MaximumDraw.IsNull() && !data.MaximumDraw.IsUnknown() {
		apiReq.SetMaximumDraw(data.MaximumDraw.ValueInt32())
	}
	if !data.AllocatedDraw.IsNull() && !data.AllocatedDraw.IsUnknown() {
		apiReq.SetAllocatedDraw(data.AllocatedDraw.ValueInt32())
	}
	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		apiReq.SetDescription(data.Description.ValueString())
	}

	tflog.Debug(ctx, "Updating power port template", map[string]interface{}{
		"id": templateID,
	})

	response, httpResp, err := r.client.DcimAPI.DcimPowerPortTemplatesUpdate(ctx, templateID).WritablePowerPortTemplateRequest(*apiReq).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating power port template",
			utils.FormatAPIError(fmt.Sprintf("update power port template ID %d", templateID), err, httpResp),
		)
		return
	}

	// Map response to model
	r.mapResponseToModel(response, &data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *PowerPortTemplateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data PowerPortTemplateResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	templateID := data.ID.ValueInt32()

	tflog.Debug(ctx, "Deleting power port template", map[string]interface{}{
		"id": templateID,
	})

	httpResp, err := r.client.DcimAPI.DcimPowerPortTemplatesDestroy(ctx, templateID).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			// Resource already deleted
			return
		}
		resp.Diagnostics.AddError(
			"Error deleting power port template",
			utils.FormatAPIError(fmt.Sprintf("delete power port template ID %d", templateID), err, httpResp),
		)
		return
	}
}

// ImportState imports the resource state from Terraform.
func (r *PowerPortTemplateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Parse the import ID as an integer
	id, err := utils.ParseInt32ID(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Could not parse import ID %q as integer: %s", req.ID, err),
		)
		return
	}

	// Set the ID in state
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

// mapResponseToModel maps the API response to the Terraform model.
func (r *PowerPortTemplateResource) mapResponseToModel(template *netbox.PowerPortTemplate, data *PowerPortTemplateResourceModel) {
	data.ID = types.Int32Value(template.GetId())
	data.Name = types.StringValue(template.GetName())

	// Map device type - store the ID as string for lookup compatibility
	if template.DeviceType.IsSet() && template.DeviceType.Get() != nil {
		data.DeviceType = types.StringValue(fmt.Sprintf("%d", template.DeviceType.Get().Id))
	} else {
		data.DeviceType = types.StringNull()
	}

	// Map module type - store the ID as string for lookup compatibility
	if template.ModuleType.IsSet() && template.ModuleType.Get() != nil {
		data.ModuleType = types.StringValue(fmt.Sprintf("%d", template.ModuleType.Get().Id))
	} else {
		data.ModuleType = types.StringNull()
	}

	// Map label
	if label, ok := template.GetLabelOk(); ok && label != nil {
		data.Label = types.StringValue(*label)
	} else {
		data.Label = types.StringValue("")
	}

	// Map type
	if template.Type.IsSet() && template.Type.Get() != nil {
		data.Type = types.StringValue(string(template.Type.Get().GetValue()))
	} else {
		data.Type = types.StringNull()
	}

	// Map maximum draw
	if template.MaximumDraw.IsSet() && template.MaximumDraw.Get() != nil {
		data.MaximumDraw = types.Int32Value(*template.MaximumDraw.Get())
	} else {
		data.MaximumDraw = types.Int32Null()
	}

	// Map allocated draw
	if template.AllocatedDraw.IsSet() && template.AllocatedDraw.Get() != nil {
		data.AllocatedDraw = types.Int32Value(*template.AllocatedDraw.Get())
	} else {
		data.AllocatedDraw = types.Int32Null()
	}

	// Map description
	if desc, ok := template.GetDescriptionOk(); ok && desc != nil {
		data.Description = types.StringValue(*desc)
	} else {
		data.Description = types.StringValue("")
	}
}
