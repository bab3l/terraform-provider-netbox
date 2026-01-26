// Package resources provides Terraform resource implementations for NetBox objects.

package resources

import (
	"context"
	"fmt"
	"maps"
	"net/http"

	"github.com/bab3l/go-netbox"
	lookup "github.com/bab3l/terraform-provider-netbox/internal/netboxlookup"
	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.

var (
	_ resource.Resource = &ConsoleServerPortResource{}

	_ resource.ResourceWithConfigure = &ConsoleServerPortResource{}

	_ resource.ResourceWithImportState = &ConsoleServerPortResource{}
	_ resource.ResourceWithIdentity    = &ConsoleServerPortResource{}
)

// NewConsoleServerPortResource returns a new resource implementing the console server port resource.

func NewConsoleServerPortResource() resource.Resource {
	return &ConsoleServerPortResource{}
}

// ConsoleServerPortResource defines the resource implementation.
type ConsoleServerPortResource struct {
	client *netbox.APIClient
}

// ConsoleServerPortResourceModel describes the resource data model.
type ConsoleServerPortResourceModel struct {
	ID            types.String `tfsdk:"id"`
	Device        types.String `tfsdk:"device"`
	Name          types.String `tfsdk:"name"`
	Label         types.String `tfsdk:"label"`
	Type          types.String `tfsdk:"type"`
	Speed         types.Int32  `tfsdk:"speed"`
	Description   types.String `tfsdk:"description"`
	MarkConnected types.Bool   `tfsdk:"mark_connected"`
	Tags          types.Set    `tfsdk:"tags"`
	CustomFields  types.Set    `tfsdk:"custom_fields"`
}

// Metadata returns the resource type name.
func (r *ConsoleServerPortResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_console_server_port"
}

// Schema defines the schema for the resource.
func (r *ConsoleServerPortResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a console server port in NetBox. Console server ports are physical console connections on console servers that provide remote access to other devices.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique numeric ID of the console server port.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"device": nbschema.RequiredReferenceAttributeWithDiffSuppress(
				"device",
				"The device this console server port belongs to (ID or name).",
			),
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the console server port.",
				Required:            true,
			},
			"label": schema.StringAttribute{
				MarkdownDescription: "Physical label of the console server port.",
				Optional:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Console server port type. Valid values: `de-9`, `db-25`, `rj-11`, `rj-12`, `rj-45`, `mini-din-8`, `usb-a`, `usb-b`, `usb-c`, `usb-mini-a`, `usb-mini-b`, `usb-micro-a`, `usb-micro-b`, `usb-micro-ab`, `other`.",
				Optional:            true,
			},
			"speed": schema.Int32Attribute{
				MarkdownDescription: "Console server port speed in bps. Valid values: `1200`, `2400`, `4800`, `9600`, `19200`, `38400`, `57600`, `115200`.",
				Optional:            true,
			},
			"mark_connected": schema.BoolAttribute{
				MarkdownDescription: "Treat as if a cable is connected.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"tags":          nbschema.TagsSlugAttribute(),
			"custom_fields": nbschema.CustomFieldsAttribute(),
		},
	}

	// Add description attribute
	maps.Copy(resp.Schema.Attributes, nbschema.DescriptionOnlyAttributes("console server port"))

	// Tags and custom fields are defined directly in the schema above.
}

func (r *ConsoleServerPortResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = nbschema.ImportIdentityWithCustomFieldsSchema()
}

// Configure adds the provider configured client to the resource.
func (r *ConsoleServerPortResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *ConsoleServerPortResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ConsoleServerPortResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Lookup device
	device, diags := lookup.LookupDevice(ctx, r.client, data.Device.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build request
	apiReq := netbox.NewWritableConsoleServerPortRequest(*device, data.Name.ValueString())

	// Set optional fields
	if !data.Label.IsNull() && !data.Label.IsUnknown() {
		apiReq.SetLabel(data.Label.ValueString())
	}

	if !data.Type.IsNull() && !data.Type.IsUnknown() {
		portType := netbox.PatchedWritableConsolePortRequestType(data.Type.ValueString())
		apiReq.SetType(portType)
	}

	if !data.Speed.IsNull() && !data.Speed.IsUnknown() {
		speed := netbox.PatchedWritableConsolePortRequestSpeed(data.Speed.ValueInt32())
		apiReq.SetSpeed(speed)
	}
	if !data.MarkConnected.IsNull() && !data.MarkConnected.IsUnknown() {
		apiReq.SetMarkConnected(data.MarkConnected.ValueBool())
	}

	// Handle description, tags, and custom fields using helpers
	utils.ApplyDescription(apiReq, data.Description)
	utils.ApplyTagsFromSlugs(ctx, r.client, apiReq, data.Tags, &resp.Diagnostics)
	utils.ApplyCustomFields(ctx, apiReq, data.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Creating console server port", map[string]interface{}{
		"device": data.Device.ValueString(),
		"name":   data.Name.ValueString(),
	})

	response, httpResp, err := r.client.DcimAPI.DcimConsoleServerPortsCreate(ctx).WritableConsoleServerPortRequest(*apiReq).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating console server port",
			utils.FormatAPIError(fmt.Sprintf("create console server port %s", data.Name.ValueString()), err, httpResp),
		)
		return
	}

	// Map response to model
	r.mapResponseToModel(ctx, response, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(data.ID.ValueString()), data.CustomFields, &resp.Diagnostics)
	tflog.Trace(ctx, "Created console server port", map[string]interface{}{
		"id":   data.ID.ValueString(),
		"name": data.Name.ValueString(),
	})
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read refreshes the resource state.
func (r *ConsoleServerPortResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ConsoleServerPortResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	portID, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Console Server Port ID",
			fmt.Sprintf("Console Server Port ID must be a number, got: %s", data.ID.ValueString()),
		)
		return
	}
	tflog.Debug(ctx, "Reading console server port", map[string]interface{}{
		"id": portID,
	})
	response, httpResp, err := r.client.DcimAPI.DcimConsoleServerPortsRetrieve(ctx, portID).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading console server port",
			utils.FormatAPIError(fmt.Sprintf("read console server port ID %d", portID), err, httpResp),
		)
		return
	}

	// Preserve original custom_fields value from state
	originalCustomFields := data.CustomFields

	// Map response to model
	r.mapResponseToModel(ctx, response, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// If custom_fields was null or empty before, restore that state
	// This prevents drift when config doesn't declare custom_fields
	if originalCustomFields.IsNull() || (utils.IsSet(originalCustomFields) && len(originalCustomFields.Elements()) == 0) {
		data.CustomFields = originalCustomFields
	}
	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(data.ID.ValueString()), data.CustomFields, &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates the resource.
func (r *ConsoleServerPortResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, plan ConsoleServerPortResourceModel

	// Read both state and plan for merge-aware custom fields handling
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	portID, err := utils.ParseID(plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Console Server Port ID",
			fmt.Sprintf("Console Server Port ID must be a number, got: %s", plan.ID.ValueString()),
		)
		return
	}

	// Lookup device
	device, diags := lookup.LookupDevice(ctx, r.client, plan.Device.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build request
	apiReq := netbox.NewWritableConsoleServerPortRequest(*device, plan.Name.ValueString())

	// Set optional fields
	utils.ApplyLabel(apiReq, plan.Label)

	// Handle type (optional enum)
	if !plan.Type.IsNull() && !plan.Type.IsUnknown() {
		portType := netbox.PatchedWritableConsolePortRequestType(plan.Type.ValueString())
		apiReq.SetType(portType)
	} else if plan.Type.IsNull() && !state.Type.IsNull() {
		// Explicitly clear when removed from config, otherwise NetBox will keep the old value.
		apiReq.SetType(netbox.PATCHEDWRITABLECONSOLEPORTREQUESTTYPE_EMPTY)
	}

	// Handle speed (nullable enum)
	if !plan.Speed.IsNull() && !plan.Speed.IsUnknown() {
		speed := netbox.PatchedWritableConsolePortRequestSpeed(plan.Speed.ValueInt32())
		apiReq.SetSpeed(speed)
	} else if plan.Speed.IsNull() && !state.Speed.IsNull() {
		// Explicitly clear when removed from config.
		apiReq.SetSpeedNil()
	}

	if !plan.MarkConnected.IsNull() && !plan.MarkConnected.IsUnknown() {
		apiReq.SetMarkConnected(plan.MarkConnected.ValueBool())
	}

	// Handle description, tags, and custom fields using merge-aware helpers
	utils.ApplyDescription(apiReq, plan.Description)
	if utils.IsSet(plan.Tags) {
		utils.ApplyTagsFromSlugs(ctx, r.client, apiReq, plan.Tags, &resp.Diagnostics)
	} else if utils.IsSet(state.Tags) {
		utils.ApplyTagsFromSlugs(ctx, r.client, apiReq, state.Tags, &resp.Diagnostics)
	}
	// Apply custom fields with merge logic to preserve unmanaged fields
	utils.ApplyCustomFieldsWithMerge(ctx, apiReq, plan.CustomFields, state.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Updating console server port", map[string]interface{}{
		"id": portID,
	})

	response, httpResp, err := r.client.DcimAPI.DcimConsoleServerPortsUpdate(ctx, portID).WritableConsoleServerPortRequest(*apiReq).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating console server port",
			utils.FormatAPIError(fmt.Sprintf("update console server port ID %d", portID), err, httpResp),
		)
		return
	}

	// Save the plan's custom fields before mapping (for filter-to-owned pattern)
	planCustomFields := plan.CustomFields

	// Map response to model
	r.mapResponseToModel(ctx, response, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Apply filter-to-owned pattern for custom fields
	plan.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, planCustomFields, response.GetCustomFields(), &resp.Diagnostics)
	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(plan.ID.ValueString()), plan.CustomFields, &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete deletes the resource.
func (r *ConsoleServerPortResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ConsoleServerPortResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	portID, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Console Server Port ID",
			fmt.Sprintf("Console Server Port ID must be a number, got: %s", data.ID.ValueString()),
		)
		return
	}
	tflog.Debug(ctx, "Deleting console server port", map[string]interface{}{
		"id": portID,
	})

	httpResp, err := r.client.DcimAPI.DcimConsoleServerPortsDestroy(ctx, portID).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			return
		}
		resp.Diagnostics.AddError(
			"Error deleting console server port",
			utils.FormatAPIError(fmt.Sprintf("delete console server port ID %d", portID), err, httpResp),
		)
		return
	}
}

// ImportState imports an existing resource.
func (r *ConsoleServerPortResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if parsed, ok := utils.ParseImportIdentityCustomFields(ctx, req.Identity, &resp.Diagnostics); ok {
		if resp.Diagnostics.HasError() {
			return
		}
		if parsed.ID == "" {
			resp.Diagnostics.AddError("Invalid import identity", "Identity id must be provided")
			return
		}

		portID, err := utils.ParseID(parsed.ID)
		if err != nil {
			resp.Diagnostics.AddError(
				"Invalid Import ID",
				fmt.Sprintf("Console Server Port ID must be a number, got: %s", parsed.ID),
			)
			return
		}

		response, httpResp, err := r.client.DcimAPI.DcimConsoleServerPortsRetrieve(ctx, portID).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error importing console server port",
				utils.FormatAPIError(fmt.Sprintf("import console server port ID %d", portID), err, httpResp),
			)
			return
		}

		var data ConsoleServerPortResourceModel
		data.Tags = types.SetNull(types.StringType)
		if device := response.GetDevice(); device.Id != 0 {
			data.Device = utils.UpdateReferenceAttribute(data.Device, device.GetName(), "", device.GetId())
		}
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

		r.mapResponseToModel(ctx, response, &data, &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}
		if parsed.HasCustomFields {
			data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, data.CustomFields, response.GetCustomFields(), &resp.Diagnostics)
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

	var data ConsoleServerPortResourceModel
	data.ID = types.StringValue(req.ID)
	data.Tags = types.SetNull(types.StringType)
	data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// mapResponseToModel maps the API response to the Terraform model.
func (r *ConsoleServerPortResource) mapResponseToModel(ctx context.Context, consoleServerPort *netbox.ConsoleServerPort, data *ConsoleServerPortResourceModel, diags *diag.Diagnostics) {
	data.ID = types.StringValue(fmt.Sprintf("%d", consoleServerPort.GetId()))
	data.Name = types.StringValue(consoleServerPort.GetName())

	// Map device - preserve user's input format
	if device := consoleServerPort.GetDevice(); device.Id != 0 {
		data.Device = utils.UpdateReferenceAttribute(data.Device, device.GetName(), "", device.GetId())
	}

	// Map label
	if label, ok := consoleServerPort.GetLabelOk(); ok && label != nil && *label != "" {
		data.Label = types.StringValue(*label)
	} else {
		data.Label = types.StringNull()
	}

	// Map type
	if consoleServerPort.Type != nil {
		portType := string(consoleServerPort.Type.GetValue())
		if portType != "" {
			data.Type = types.StringValue(portType)
		} else {
			data.Type = types.StringNull()
		}
	} else {
		data.Type = types.StringNull()
	}

	// Map speed
	if consoleServerPort.Speed.IsSet() && consoleServerPort.Speed.Get() != nil {
		data.Speed = types.Int32Value(int32(consoleServerPort.Speed.Get().GetValue()))
	} else {
		data.Speed = types.Int32Null()
	}

	// Map description
	if desc, ok := consoleServerPort.GetDescriptionOk(); ok && desc != nil && *desc != "" {
		data.Description = types.StringValue(*desc)
	} else {
		data.Description = types.StringNull()
	}

	// Map mark_connected
	if mc, ok := consoleServerPort.GetMarkConnectedOk(); ok && mc != nil {
		data.MarkConnected = types.BoolValue(*mc)
	} else {
		data.MarkConnected = types.BoolValue(false)
	}

	// Tags (slug list)
	data.Tags = utils.PopulateTagsSlugFromAPI(ctx, consoleServerPort.HasTags(), consoleServerPort.GetTags(), data.Tags)

	// Handle custom fields - filter to only owned fields
	data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, data.CustomFields, consoleServerPort.GetCustomFields(), diags)
}
