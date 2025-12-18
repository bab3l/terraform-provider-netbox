// Package resources provides Terraform resource implementations for NetBox objects.

package resources

import (
	"context"
	"fmt"

	"github.com/bab3l/go-netbox"
	"github.com/bab3l/terraform-provider-netbox/internal/netboxlookup"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.

var _ resource.Resource = &ConsolePortTemplateResource{}

var _ resource.ResourceWithImportState = &ConsolePortTemplateResource{}

// NewConsolePortTemplateResource returns a new resource implementing the console port template resource.

func NewConsolePortTemplateResource() resource.Resource {

	return &ConsolePortTemplateResource{}

}

// ConsolePortTemplateResource defines the resource implementation.

type ConsolePortTemplateResource struct {
	client *netbox.APIClient
}

// ConsolePortTemplateResourceModel describes the resource data model.

type ConsolePortTemplateResourceModel struct {
	ID types.Int32 `tfsdk:"id"`

	DeviceType types.String `tfsdk:"device_type"`

	ModuleType types.String `tfsdk:"module_type"`

	Name types.String `tfsdk:"name"`

	Label types.String `tfsdk:"label"`

	Type types.String `tfsdk:"type"`

	Description types.String `tfsdk:"description"`
}

// Metadata returns the resource type name.

func (r *ConsolePortTemplateResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {

	resp.TypeName = req.ProviderTypeName + "_console_port_template"

}

// Schema defines the schema for the resource.

func (r *ConsolePortTemplateResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	resp.Schema = schema.Schema{

		MarkdownDescription: "Manages a console port template in NetBox. Console port templates define the default console ports that will be created on new devices of a specific device type or modules of a module type.",

		Attributes: map[string]schema.Attribute{

			"id": schema.Int32Attribute{

				MarkdownDescription: "The unique numeric ID of the console port template.",

				Computed: true,

				PlanModifiers: []planmodifier.Int32{

					int32planmodifier.UseStateForUnknown(),
				},
			},

			"device_type": schema.StringAttribute{

				MarkdownDescription: "The device type ID or slug. Either device_type or module_type must be specified.",

				Optional: true,
			},

			"module_type": schema.StringAttribute{

				MarkdownDescription: "The module type ID or model name. Either device_type or module_type must be specified.",

				Optional: true,
			},

			"name": schema.StringAttribute{

				MarkdownDescription: "The name of the console port template. Use {module} as a substitution for the module bay position when attached to a module type.",

				Required: true,
			},

			"label": schema.StringAttribute{

				MarkdownDescription: "Physical label of the console port template.",

				Optional: true,

				Computed: true,

				Default: stringdefault.StaticString(""),
			},

			"type": schema.StringAttribute{

				MarkdownDescription: "The type of console port (e.g., de-9, db-25, rj-45, usb-a, usb-b, usb-c, usb-mini-a, usb-mini-b, usb-micro-a, usb-micro-b, usb-micro-ab, other).",

				Optional: true,
			},

			"description": schema.StringAttribute{

				MarkdownDescription: "A description of the console port template.",

				Optional: true,

				Computed: true,

				Default: stringdefault.StaticString(""),
			},
		},
	}

}

// Configure adds the provider configured client to the resource.

func (r *ConsolePortTemplateResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {

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

func (r *ConsolePortTemplateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var data ConsolePortTemplateResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	// Build the API request

	apiReq := netbox.NewWritableConsolePortTemplateRequest(data.Name.ValueString())

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

		apiReq.SetType(netbox.ConsolePortTypeValue(data.Type.ValueString()))

	}

	if !data.Description.IsNull() && !data.Description.IsUnknown() {

		apiReq.SetDescription(data.Description.ValueString())

	}

	tflog.Debug(ctx, "Creating console port template", map[string]interface{}{

		"name": data.Name.ValueString(),
	})

	response, httpResp, err := r.client.DcimAPI.DcimConsolePortTemplatesCreate(ctx).WritableConsolePortTemplateRequest(*apiReq).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		resp.Diagnostics.AddError(

			"Error creating console port template",

			utils.FormatAPIError("create console port template", err, httpResp),
		)

		return

	}

	// Map response to model

	r.mapResponseToModel(response, &data)

	tflog.Trace(ctx, "Created console port template", map[string]interface{}{

		"id": data.ID.ValueInt32(),
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

// Read refreshes the Terraform state with the latest data.

func (r *ConsolePortTemplateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	var data ConsolePortTemplateResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	templateID := data.ID.ValueInt32()

	tflog.Debug(ctx, "Reading console port template", map[string]interface{}{

		"id": templateID,
	})

	response, httpResp, err := r.client.DcimAPI.DcimConsolePortTemplatesRetrieve(ctx, templateID).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		if httpResp != nil && httpResp.StatusCode == 404 {

			tflog.Debug(ctx, "Console port template not found, removing from state", map[string]interface{}{

				"id": templateID,
			})

			resp.State.RemoveResource(ctx)

			return

		}

		resp.Diagnostics.AddError(

			"Error reading console port template",

			utils.FormatAPIError(fmt.Sprintf("read console port template ID %d", templateID), err, httpResp),
		)

		return

	}

	// Map response to model

	r.mapResponseToModel(response, &data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

// Update updates the resource and sets the updated Terraform state on success.

func (r *ConsolePortTemplateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var data ConsolePortTemplateResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	templateID := data.ID.ValueInt32()

	// Build the API request

	apiReq := netbox.NewWritableConsolePortTemplateRequest(data.Name.ValueString())

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

		apiReq.SetType(netbox.ConsolePortTypeValue(data.Type.ValueString()))

	}

	if !data.Description.IsNull() && !data.Description.IsUnknown() {

		apiReq.SetDescription(data.Description.ValueString())

	}

	tflog.Debug(ctx, "Updating console port template", map[string]interface{}{

		"id": templateID,
	})

	response, httpResp, err := r.client.DcimAPI.DcimConsolePortTemplatesUpdate(ctx, templateID).WritableConsolePortTemplateRequest(*apiReq).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		resp.Diagnostics.AddError(

			"Error updating console port template",

			utils.FormatAPIError(fmt.Sprintf("update console port template ID %d", templateID), err, httpResp),
		)

		return

	}

	// Map response to model

	r.mapResponseToModel(response, &data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

// Delete deletes the resource and removes the Terraform state on success.

func (r *ConsolePortTemplateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var data ConsolePortTemplateResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	templateID := data.ID.ValueInt32()

	tflog.Debug(ctx, "Deleting console port template", map[string]interface{}{

		"id": templateID,
	})

	httpResp, err := r.client.DcimAPI.DcimConsolePortTemplatesDestroy(ctx, templateID).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		if httpResp != nil && httpResp.StatusCode == 404 {

			// Resource already deleted

			return

		}

		resp.Diagnostics.AddError(

			"Error deleting console port template",

			utils.FormatAPIError(fmt.Sprintf("delete console port template ID %d", templateID), err, httpResp),
		)

		return

	}

}

// ImportState imports the resource state from Terraform.

func (r *ConsolePortTemplateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

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

func (r *ConsolePortTemplateResource) mapResponseToModel(template *netbox.ConsolePortTemplate, data *ConsolePortTemplateResourceModel) {

	data.ID = types.Int32Value(template.GetId())

	data.Name = types.StringValue(template.GetName())

	// Map device type - preserve user's input format

	if template.DeviceType.IsSet() && template.DeviceType.Get() != nil {

		dt := template.DeviceType.Get()

		data.DeviceType = utils.UpdateReferenceAttribute(data.DeviceType, dt.GetModel(), dt.GetSlug(), dt.Id)

	} else {

		data.DeviceType = types.StringNull()

	}

	// Map module type - preserve user's input format

	if template.ModuleType.IsSet() && template.ModuleType.Get() != nil {

		mt := template.ModuleType.Get()

		data.ModuleType = utils.UpdateReferenceAttribute(data.ModuleType, mt.GetModel(), "", mt.Id)

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

	if template.Type != nil {

		data.Type = types.StringValue(string(template.Type.GetValue()))

	} else {

		data.Type = types.StringNull()

	}

	// Map description

	if desc, ok := template.GetDescriptionOk(); ok && desc != nil {

		data.Description = types.StringValue(*desc)

	} else {

		data.Description = types.StringValue("")

	}

}
