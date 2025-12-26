// Package resources provides Terraform resource implementations for NetBox objects.

package resources

import (
	"context"
	"fmt"
	"maps"

	"github.com/bab3l/go-netbox"
	"github.com/bab3l/terraform-provider-netbox/internal/netboxlookup"
	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.

var _ resource.Resource = &FrontPortTemplateResource{}

var _ resource.ResourceWithImportState = &FrontPortTemplateResource{}

// NewFrontPortTemplateResource returns a new resource implementing the front port template resource.

func NewFrontPortTemplateResource() resource.Resource {

	return &FrontPortTemplateResource{}

}

// FrontPortTemplateResource defines the resource implementation.

type FrontPortTemplateResource struct {
	client *netbox.APIClient
}

// FrontPortTemplateResourceModel describes the resource data model.

type FrontPortTemplateResourceModel struct {
	ID types.Int32 `tfsdk:"id"`

	DeviceType types.String `tfsdk:"device_type"`

	ModuleType types.String `tfsdk:"module_type"`

	Name types.String `tfsdk:"name"`

	Label types.String `tfsdk:"label"`

	Type types.String `tfsdk:"type"`

	Color types.String `tfsdk:"color"`

	DisplayName types.String `tfsdk:"display_name"`

	RearPort types.String `tfsdk:"rear_port"`

	RearPortPosition types.Int32 `tfsdk:"rear_port_position"`

	Description types.String `tfsdk:"description"`
}

// Metadata returns the resource type name.

func (r *FrontPortTemplateResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {

	resp.TypeName = req.ProviderTypeName + "_front_port_template"

}

// Schema defines the schema for the resource.

func (r *FrontPortTemplateResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	resp.Schema = schema.Schema{

		MarkdownDescription: "Manages a front port template in NetBox. Front port templates define the default front ports that will be created on new devices of a specific device type or modules of a module type. Each front port must map to a rear port.",

		Attributes: map[string]schema.Attribute{

			"id": schema.Int32Attribute{

				MarkdownDescription: "The unique numeric ID of the front port template.",

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

				MarkdownDescription: "The name of the front port template. Use {module} as a substitution for the module bay position when attached to a module type.",

				Required: true,
			},

			"label": schema.StringAttribute{

				MarkdownDescription: "Physical label of the front port template.",

				Optional: true,

				Computed: true,

				Default: stringdefault.StaticString(""),
			},

			"type": schema.StringAttribute{

				MarkdownDescription: "The type of front port (e.g., `8p8c`, `8p6c`, `110-punch`, `bnc`, `f`, `n`, `mrj21`, `fc`, `lc`, `lc-pc`, `lc-upc`, `lc-apc`, `lsh`, `lsh-pc`, `lsh-upc`, `lsh-apc`, `mpo`, `mtrj`, `sc`, `sc-pc`, `sc-upc`, `sc-apc`, `st`, `cs`, `sn`, `sma-905`, `sma-906`, `splice`, `other`).",

				Required: true,
			},

			"color": schema.StringAttribute{

				MarkdownDescription: "Color of the front port in hex format (e.g., `aa1409`).",

				Optional: true,

				Computed: true,

				Default: stringdefault.StaticString(""),
			},

			"display_name": nbschema.DisplayNameAttribute("front port template"),

			"rear_port": schema.StringAttribute{

				MarkdownDescription: "The name of the rear port template on the same device type or module type that this front port maps to.",

				Required: true,
			},

			"rear_port_position": schema.Int32Attribute{

				MarkdownDescription: "Position on the rear port that this front port maps to (1-1024). Default is 1.",

				Optional: true,

				Computed: true,

				Default: int32default.StaticInt32(1),
			},
		},
	}

	// Add description attribute
	maps.Copy(resp.Schema.Attributes, nbschema.DescriptionOnlyAttributes("front port template"))
}

// Configure adds the provider configured client to the resource.

func (r *FrontPortTemplateResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {

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

func (r *FrontPortTemplateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var data FrontPortTemplateResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	// Build the rear port reference - required field referencing by name

	rearPortRef := netbox.NewBriefRearPortTemplateRequest(data.RearPort.ValueString())

	// Build the API request

	apiReq := netbox.NewWritableFrontPortTemplateRequest(

		data.Name.ValueString(),

		netbox.FrontPortTypeValue(data.Type.ValueString()),

		*rearPortRef,
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

	if !data.Color.IsNull() && !data.Color.IsUnknown() {

		apiReq.SetColor(data.Color.ValueString())

	}

	if !data.RearPortPosition.IsNull() && !data.RearPortPosition.IsUnknown() {

		apiReq.SetRearPortPosition(data.RearPortPosition.ValueInt32())

	}

	if !data.Description.IsNull() && !data.Description.IsUnknown() {

		apiReq.SetDescription(data.Description.ValueString())

	}

	tflog.Debug(ctx, "Creating front port template", map[string]interface{}{

		"name": data.Name.ValueString(),

		"rear_port": data.RearPort.ValueString(),
	})

	response, httpResp, err := r.client.DcimAPI.DcimFrontPortTemplatesCreate(ctx).WritableFrontPortTemplateRequest(*apiReq).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		resp.Diagnostics.AddError(

			"Error creating front port template",

			utils.FormatAPIError("create front port template", err, httpResp),
		)

		return

	}

	// Map response to model

	r.mapResponseToModel(response, &data)

	tflog.Trace(ctx, "Created front port template", map[string]interface{}{

		"id": data.ID.ValueInt32(),
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

// Read refreshes the Terraform state with the latest data.

func (r *FrontPortTemplateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	var data FrontPortTemplateResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	templateID := data.ID.ValueInt32()

	tflog.Debug(ctx, "Reading front port template", map[string]interface{}{

		"id": templateID,
	})

	response, httpResp, err := r.client.DcimAPI.DcimFrontPortTemplatesRetrieve(ctx, templateID).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		if httpResp != nil && httpResp.StatusCode == 404 {

			tflog.Debug(ctx, "Front port template not found, removing from state", map[string]interface{}{

				"id": templateID,
			})

			resp.State.RemoveResource(ctx)

			return

		}

		resp.Diagnostics.AddError(

			"Error reading front port template",

			utils.FormatAPIError(fmt.Sprintf("read front port template ID %d", templateID), err, httpResp),
		)

		return

	}

	// Map response to model

	r.mapResponseToModel(response, &data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

// Update updates the resource and sets the updated Terraform state on success.

func (r *FrontPortTemplateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var data FrontPortTemplateResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	templateID := data.ID.ValueInt32()

	// Build the rear port reference - required field referencing by name

	rearPortRef := netbox.NewBriefRearPortTemplateRequest(data.RearPort.ValueString())

	// Build the API request

	apiReq := netbox.NewWritableFrontPortTemplateRequest(

		data.Name.ValueString(),

		netbox.FrontPortTypeValue(data.Type.ValueString()),

		*rearPortRef,
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

	if !data.Color.IsNull() && !data.Color.IsUnknown() {

		apiReq.SetColor(data.Color.ValueString())

	}

	if !data.RearPortPosition.IsNull() && !data.RearPortPosition.IsUnknown() {

		apiReq.SetRearPortPosition(data.RearPortPosition.ValueInt32())

	}

	if !data.Description.IsNull() && !data.Description.IsUnknown() {

		apiReq.SetDescription(data.Description.ValueString())

	}

	tflog.Debug(ctx, "Updating front port template", map[string]interface{}{

		"id": templateID,
	})

	response, httpResp, err := r.client.DcimAPI.DcimFrontPortTemplatesUpdate(ctx, templateID).WritableFrontPortTemplateRequest(*apiReq).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		resp.Diagnostics.AddError(

			"Error updating front port template",

			utils.FormatAPIError(fmt.Sprintf("update front port template ID %d", templateID), err, httpResp),
		)

		return

	}

	// Map response to model

	r.mapResponseToModel(response, &data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

// Delete deletes the resource and removes the Terraform state on success.

func (r *FrontPortTemplateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var data FrontPortTemplateResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	templateID := data.ID.ValueInt32()

	tflog.Debug(ctx, "Deleting front port template", map[string]interface{}{

		"id": templateID,
	})

	httpResp, err := r.client.DcimAPI.DcimFrontPortTemplatesDestroy(ctx, templateID).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		if httpResp != nil && httpResp.StatusCode == 404 {

			// Resource already deleted

			return

		}

		resp.Diagnostics.AddError(

			"Error deleting front port template",

			utils.FormatAPIError(fmt.Sprintf("delete front port template ID %d", templateID), err, httpResp),
		)

		return

	}

}

// ImportState imports the resource state from Terraform.

func (r *FrontPortTemplateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

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

func (r *FrontPortTemplateResource) mapResponseToModel(template *netbox.FrontPortTemplate, data *FrontPortTemplateResourceModel) {

	data.ID = types.Int32Value(template.GetId())

	data.Name = types.StringValue(template.GetName())

	// DisplayName
	if template.Display != "" {
		data.DisplayName = types.StringValue(template.Display)
	} else {
		data.DisplayName = types.StringNull()
	}

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

	// Map type

	data.Type = types.StringValue(string(template.Type.GetValue()))

	// Map label

	if label, ok := template.GetLabelOk(); ok && label != nil {

		data.Label = types.StringValue(*label)

	} else {

		data.Label = types.StringValue("")

	}

	// Map color

	if color, ok := template.GetColorOk(); ok && color != nil {

		data.Color = types.StringValue(*color)

	} else {

		data.Color = types.StringValue("")

	}

	// Map rear port - store the name for reference

	data.RearPort = types.StringValue(template.RearPort.GetName())

	// Map rear port position

	if pos, ok := template.GetRearPortPositionOk(); ok && pos != nil {

		data.RearPortPosition = types.Int32Value(*pos)

	} else {

		data.RearPortPosition = types.Int32Value(1)

	}

	// Map description

	if desc, ok := template.GetDescriptionOk(); ok && desc != nil {

		data.Description = types.StringValue(*desc)

	} else {

		data.Description = types.StringValue("")

	}

}
