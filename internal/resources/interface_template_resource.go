// Package resources provides Terraform resource implementations for NetBox objects.
package resources

import (
	"context"
	"fmt"

	"github.com/bab3l/go-netbox"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/bab3l/terraform-provider-netbox/internal/netboxlookup"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &InterfaceTemplateResource{}
var _ resource.ResourceWithImportState = &InterfaceTemplateResource{}

// NewInterfaceTemplateResource returns a new resource implementing the interface template resource.
func NewInterfaceTemplateResource() resource.Resource {
	return &InterfaceTemplateResource{}
}

// InterfaceTemplateResource defines the resource implementation.
type InterfaceTemplateResource struct {
	client *netbox.APIClient
}

// InterfaceTemplateResourceModel describes the resource data model.
type InterfaceTemplateResourceModel struct {
	ID          types.Int32  `tfsdk:"id"`
	DeviceType  types.String `tfsdk:"device_type"`
	ModuleType  types.String `tfsdk:"module_type"`
	Name        types.String `tfsdk:"name"`
	Label       types.String `tfsdk:"label"`
	Type        types.String `tfsdk:"type"`
	Enabled     types.Bool   `tfsdk:"enabled"`
	MgmtOnly    types.Bool   `tfsdk:"mgmt_only"`
	Description types.String `tfsdk:"description"`
	Bridge      types.Int32  `tfsdk:"bridge"`
	PoeMode     types.String `tfsdk:"poe_mode"`
	PoeType     types.String `tfsdk:"poe_type"`
	RfRole      types.String `tfsdk:"rf_role"`
}

// Metadata returns the resource type name.
func (r *InterfaceTemplateResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_interface_template"
}

// Schema defines the schema for the resource.
func (r *InterfaceTemplateResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an interface template in NetBox. Interface templates define the default interfaces that will be created on new devices of a specific device type or modules of a module type.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int32Attribute{
				MarkdownDescription: "The unique numeric ID of the interface template.",
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
				MarkdownDescription: "The name of the interface template. Use {module} as a substitution for the module bay position when attached to a module type.",
				Required:            true,
			},
			"label": schema.StringAttribute{
				MarkdownDescription: "Physical label of the interface template.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "The type of interface (e.g., 1000base-t, 10gbase-x-sfpp, 25gbase-x-sfp28, 40gbase-x-qsfpp, 100gbase-x-qsfp28, virtual, lag, etc.).",
				Required:            true,
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Whether the interface is enabled by default.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"mgmt_only": schema.BoolAttribute{
				MarkdownDescription: "Whether the interface is for management only.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "A description of the interface template.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"bridge": schema.Int32Attribute{
				MarkdownDescription: "The ID of the bridge interface template this interface belongs to.",
				Optional:            true,
			},
			"poe_mode": schema.StringAttribute{
				MarkdownDescription: "PoE mode (pd or pse).",
				Optional:            true,
			},
			"poe_type": schema.StringAttribute{
				MarkdownDescription: "PoE type (type1-ieee802.3af, type2-ieee802.3at, type3-ieee802.3bt, type4-ieee802.3bt, passive-24v-2pair, passive-24v-4pair, passive-48v-2pair, passive-48v-4pair).",
				Optional:            true,
			},
			"rf_role": schema.StringAttribute{
				MarkdownDescription: "Wireless role (ap or station).",
				Optional:            true,
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *InterfaceTemplateResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *InterfaceTemplateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data InterfaceTemplateResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the API request
	apiReq := netbox.NewWritableInterfaceTemplateRequest(
		data.Name.ValueString(),
		netbox.InterfaceTypeValue(data.Type.ValueString()),
	)

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
	if !data.Enabled.IsNull() && !data.Enabled.IsUnknown() {
		apiReq.SetEnabled(data.Enabled.ValueBool())
	}
	if !data.MgmtOnly.IsNull() && !data.MgmtOnly.IsUnknown() {
		apiReq.SetMgmtOnly(data.MgmtOnly.ValueBool())
	}
	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		apiReq.SetDescription(data.Description.ValueString())
	}
	if !data.Bridge.IsNull() && !data.Bridge.IsUnknown() {
		apiReq.SetBridge(data.Bridge.ValueInt32())
	}
	if !data.PoeMode.IsNull() && !data.PoeMode.IsUnknown() {
		apiReq.SetPoeMode(netbox.InterfacePoeModeValue(data.PoeMode.ValueString()))
	}
	if !data.PoeType.IsNull() && !data.PoeType.IsUnknown() {
		apiReq.SetPoeType(netbox.InterfacePoeTypeValue(data.PoeType.ValueString()))
	}
	if !data.RfRole.IsNull() && !data.RfRole.IsUnknown() {
		apiReq.SetRfRole(netbox.WirelessRole(data.RfRole.ValueString()))
	}

	tflog.Debug(ctx, "Creating interface template", map[string]interface{}{
		"name": data.Name.ValueString(),
	})

	response, httpResp, err := r.client.DcimAPI.DcimInterfaceTemplatesCreate(ctx).WritableInterfaceTemplateRequest(*apiReq).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating interface template",
			utils.FormatAPIError("create interface template", err, httpResp),
		)
		return
	}

	// Map response to model
	r.mapResponseToModel(response, &data)

	tflog.Trace(ctx, "Created interface template", map[string]interface{}{
		"id": data.ID.ValueInt32(),
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *InterfaceTemplateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data InterfaceTemplateResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	templateID := data.ID.ValueInt32()

	tflog.Debug(ctx, "Reading interface template", map[string]interface{}{
		"id": templateID,
	})

	response, httpResp, err := r.client.DcimAPI.DcimInterfaceTemplatesRetrieve(ctx, templateID).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			tflog.Debug(ctx, "Interface template not found, removing from state", map[string]interface{}{
				"id": templateID,
			})
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading interface template",
			utils.FormatAPIError(fmt.Sprintf("read interface template ID %d", templateID), err, httpResp),
		)
		return
	}

	// Map response to model
	r.mapResponseToModel(response, &data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *InterfaceTemplateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data InterfaceTemplateResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	templateID := data.ID.ValueInt32()

	// Build the API request
	apiReq := netbox.NewWritableInterfaceTemplateRequest(
		data.Name.ValueString(),
		netbox.InterfaceTypeValue(data.Type.ValueString()),
	)

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
	if !data.Enabled.IsNull() && !data.Enabled.IsUnknown() {
		apiReq.SetEnabled(data.Enabled.ValueBool())
	}
	if !data.MgmtOnly.IsNull() && !data.MgmtOnly.IsUnknown() {
		apiReq.SetMgmtOnly(data.MgmtOnly.ValueBool())
	}
	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		apiReq.SetDescription(data.Description.ValueString())
	}
	if !data.Bridge.IsNull() && !data.Bridge.IsUnknown() {
		apiReq.SetBridge(data.Bridge.ValueInt32())
	}
	if !data.PoeMode.IsNull() && !data.PoeMode.IsUnknown() {
		apiReq.SetPoeMode(netbox.InterfacePoeModeValue(data.PoeMode.ValueString()))
	}
	if !data.PoeType.IsNull() && !data.PoeType.IsUnknown() {
		apiReq.SetPoeType(netbox.InterfacePoeTypeValue(data.PoeType.ValueString()))
	}
	if !data.RfRole.IsNull() && !data.RfRole.IsUnknown() {
		apiReq.SetRfRole(netbox.WirelessRole(data.RfRole.ValueString()))
	}

	tflog.Debug(ctx, "Updating interface template", map[string]interface{}{
		"id": templateID,
	})

	response, httpResp, err := r.client.DcimAPI.DcimInterfaceTemplatesUpdate(ctx, templateID).WritableInterfaceTemplateRequest(*apiReq).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating interface template",
			utils.FormatAPIError(fmt.Sprintf("update interface template ID %d", templateID), err, httpResp),
		)
		return
	}

	// Map response to model
	r.mapResponseToModel(response, &data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *InterfaceTemplateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data InterfaceTemplateResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	templateID := data.ID.ValueInt32()

	tflog.Debug(ctx, "Deleting interface template", map[string]interface{}{
		"id": templateID,
	})

	httpResp, err := r.client.DcimAPI.DcimInterfaceTemplatesDestroy(ctx, templateID).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			// Resource already deleted
			return
		}
		resp.Diagnostics.AddError(
			"Error deleting interface template",
			utils.FormatAPIError(fmt.Sprintf("delete interface template ID %d", templateID), err, httpResp),
		)
		return
	}
}

// ImportState imports the resource state from Terraform.
func (r *InterfaceTemplateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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
func (r *InterfaceTemplateResource) mapResponseToModel(template *netbox.InterfaceTemplate, data *InterfaceTemplateResourceModel) {
	data.ID = types.Int32Value(template.GetId())
	data.Name = types.StringValue(template.GetName())
	data.Type = types.StringValue(string(template.Type.GetValue()))

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

	// Map enabled
	if enabled, ok := template.GetEnabledOk(); ok && enabled != nil {
		data.Enabled = types.BoolValue(*enabled)
	} else {
		data.Enabled = types.BoolValue(true)
	}

	// Map mgmt_only
	if mgmtOnly, ok := template.GetMgmtOnlyOk(); ok && mgmtOnly != nil {
		data.MgmtOnly = types.BoolValue(*mgmtOnly)
	} else {
		data.MgmtOnly = types.BoolValue(false)
	}

	// Map description
	if desc, ok := template.GetDescriptionOk(); ok && desc != nil {
		data.Description = types.StringValue(*desc)
	} else {
		data.Description = types.StringValue("")
	}

	// Map bridge
	if template.Bridge.IsSet() && template.Bridge.Get() != nil {
		data.Bridge = types.Int32Value(template.Bridge.Get().Id)
	} else {
		data.Bridge = types.Int32Null()
	}

	// Map poe_mode
	if poeMode, ok := template.GetPoeModeOk(); ok && poeMode != nil && poeMode.Value != nil {
		data.PoeMode = types.StringValue(string(*poeMode.Value))
	} else {
		data.PoeMode = types.StringNull()
	}

	// Map poe_type
	if poeType, ok := template.GetPoeTypeOk(); ok && poeType != nil && poeType.Value != nil {
		data.PoeType = types.StringValue(string(*poeType.Value))
	} else {
		data.PoeType = types.StringNull()
	}

	// Map rf_role
	if rfRole, ok := template.GetRfRoleOk(); ok && rfRole != nil && rfRole.Value != nil {
		data.RfRole = types.StringValue(string(*rfRole.Value))
	} else {
		data.RfRole = types.StringNull()
	}
}
