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
	_ resource.Resource                = &ConsolePortResource{}
	_ resource.ResourceWithConfigure   = &ConsolePortResource{}
	_ resource.ResourceWithImportState = &ConsolePortResource{}
	_ resource.ResourceWithIdentity    = &ConsolePortResource{}
)

// NewConsolePortResource returns a new resource implementing the console port resource.
func NewConsolePortResource() resource.Resource {
	return &ConsolePortResource{}
}

// ConsolePortResource defines the resource implementation.
type ConsolePortResource struct {
	client *netbox.APIClient
}

// ConsolePortResourceModel describes the resource data model.
type ConsolePortResourceModel struct {
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
func (r *ConsolePortResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_console_port"
}

// Schema defines the schema for the resource.
func (r *ConsolePortResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a console port in NetBox. Console ports are physical console connections on devices.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique numeric ID of the console port.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"device": nbschema.RequiredReferenceAttributeWithDiffSuppress("device", "ID or name of the device this console port belongs to. Required."),
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the console port.",
				Required:            true,
			},
			"label": schema.StringAttribute{
				MarkdownDescription: "Physical label of the console port.",
				Optional:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Console port type. Valid values: `de-9`, `db-25`, `rj-11`, `rj-12`, `rj-45`, `mini-din-8`, `usb-a`, `usb-b`, `usb-c`, `usb-mini-a`, `usb-mini-b`, `usb-micro-a`, `usb-micro-b`, `usb-micro-ab`, `other`.",
				Optional:            true,
			},
			"speed": schema.Int32Attribute{
				MarkdownDescription: "Console port speed in bps. Valid values: `1200`, `2400`, `4800`, `9600`, `19200`, `38400`, `57600`, `115200`.",
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
	maps.Copy(resp.Schema.Attributes, nbschema.DescriptionOnlyAttributes("console port"))
}

func (r *ConsolePortResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = nbschema.ImportIdentityWithCustomFieldsSchema()
}

// Configure adds the provider configured client to the resource.
func (r *ConsolePortResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *ConsolePortResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ConsolePortResourceModel
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
	apiReq := netbox.NewWritableConsolePortRequest(*device, data.Name.ValueString())

	// Set optional fields
	utils.ApplyLabel(apiReq, data.Label)

	// Type (*PatchedWritableConsolePortRequestType) - use empty string to clear
	if utils.IsSet(data.Type) {
		portType := netbox.PatchedWritableConsolePortRequestType(data.Type.ValueString())
		apiReq.SetType(portType)
	} else if data.Type.IsNull() {
		apiReq.SetType("")
	}

	// Speed (NullablePatchedWritableConsolePortRequestSpeed) - use Nullable wrapper to clear
	if utils.IsSet(data.Speed) {
		speed := netbox.PatchedWritableConsolePortRequestSpeed(data.Speed.ValueInt32())
		apiReq.SetSpeed(speed)
	} else if data.Speed.IsNull() {
		apiReq.Speed = *netbox.NewNullablePatchedWritableConsolePortRequestSpeed(nil)
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
	tflog.Debug(ctx, "Creating console port", map[string]interface{}{
		"device": data.Device.ValueString(),
		"name":   data.Name.ValueString(),
	})

	response, httpResp, err := r.client.DcimAPI.DcimConsolePortsCreate(ctx).WritableConsolePortRequest(*apiReq).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating console port",
			utils.FormatAPIError(fmt.Sprintf("create console port %s", data.Name.ValueString()), err, httpResp),
		)
		return
	}

	// Map response to model
	r.mapResponseToModel(ctx, response, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(data.ID.ValueString()), data.CustomFields, &resp.Diagnostics)
	tflog.Trace(ctx, "Created console port", map[string]interface{}{
		"id":   data.ID.ValueString(),
		"name": data.Name.ValueString(),
	})
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read refreshes the resource state.
func (r *ConsolePortResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ConsolePortResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	portID, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Console Port ID",
			fmt.Sprintf("Console Port ID must be a number, got: %s", data.ID.ValueString()),
		)
		return
	}
	tflog.Debug(ctx, "Reading console port", map[string]interface{}{
		"id": portID,
	})

	response, httpResp, err := r.client.DcimAPI.DcimConsolePortsRetrieve(ctx, portID).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading console port",
			utils.FormatAPIError(fmt.Sprintf("read console port ID %d", portID), err, httpResp),
		)
		return
	}

	// Map response to model
	r.mapResponseToModel(ctx, response, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(data.ID.ValueString()), data.CustomFields, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates the resource.
func (r *ConsolePortResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, data ConsolePortResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	portID, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Console Port ID",
			fmt.Sprintf("Console Port ID must be a number, got: %s", data.ID.ValueString()),
		)
		return
	}

	// Lookup device
	device, diags := lookup.LookupDevice(ctx, r.client, data.Device.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build request
	apiReq := netbox.NewWritableConsolePortRequest(*device, data.Name.ValueString())

	// Set optional fields
	utils.ApplyLabel(apiReq, data.Label)

	// Type (*PatchedWritableConsolePortRequestType) - use empty string to clear
	if utils.IsSet(data.Type) {
		portType := netbox.PatchedWritableConsolePortRequestType(data.Type.ValueString())
		apiReq.SetType(portType)
	} else if data.Type.IsNull() {
		apiReq.SetType("")
	}

	// Speed (NullablePatchedWritableConsolePortRequestSpeed) - use Nullable wrapper to clear
	if utils.IsSet(data.Speed) {
		speed := netbox.PatchedWritableConsolePortRequestSpeed(data.Speed.ValueInt32())
		apiReq.SetSpeed(speed)
	} else if data.Speed.IsNull() {
		apiReq.Speed = *netbox.NewNullablePatchedWritableConsolePortRequestSpeed(nil)
	}

	if !data.MarkConnected.IsNull() && !data.MarkConnected.IsUnknown() {
		apiReq.SetMarkConnected(data.MarkConnected.ValueBool())
	}

	// Handle description, tags, and custom fields with merge-aware behavior
	utils.ApplyDescription(apiReq, data.Description)

	// Apply tags - merge-aware: use plan if provided, else use state
	if utils.IsSet(data.Tags) {
		utils.ApplyTagsFromSlugs(ctx, r.client, apiReq, data.Tags, &resp.Diagnostics)
	} else if utils.IsSet(state.Tags) {
		utils.ApplyTagsFromSlugs(ctx, r.client, apiReq, state.Tags, &resp.Diagnostics)
	}

	// Apply custom fields with merge logic (preserves unmanaged fields)
	utils.ApplyCustomFieldsWithMerge(ctx, apiReq, data.CustomFields, state.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Updating console port", map[string]interface{}{
		"id": portID,
	})
	response, httpResp, err := r.client.DcimAPI.DcimConsolePortsUpdate(ctx, portID).WritableConsolePortRequest(*apiReq).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating console port",
			utils.FormatAPIError(fmt.Sprintf("update console port ID %d", portID), err, httpResp),
		)
		return
	}

	// Map response to model
	r.mapResponseToModel(ctx, response, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(data.ID.ValueString()), data.CustomFields, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete deletes the resource.
func (r *ConsolePortResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ConsolePortResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	portID, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Console Port ID",
			fmt.Sprintf("Console Port ID must be a number, got: %s", data.ID.ValueString()),
		)
		return
	}
	tflog.Debug(ctx, "Deleting console port", map[string]interface{}{
		"id": portID,
	})

	httpResp, err := r.client.DcimAPI.DcimConsolePortsDestroy(ctx, portID).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			return
		}
		resp.Diagnostics.AddError(
			"Error deleting console port",
			utils.FormatAPIError(fmt.Sprintf("delete console port ID %d", portID), err, httpResp),
		)
		return
	}
}

// ImportState imports an existing resource.
func (r *ConsolePortResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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
				fmt.Sprintf("Console Port ID must be a number, got: %s", parsed.ID),
			)
			return
		}

		response, httpResp, err := r.client.DcimAPI.DcimConsolePortsRetrieve(ctx, portID).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error importing console port",
				utils.FormatAPIError(fmt.Sprintf("import console port ID %d", portID), err, httpResp),
			)
			return
		}

		var data ConsolePortResourceModel
		data.Tags = types.SetNull(types.StringType)
		if device := response.GetDevice(); device.Id != 0 {
			data.Device = types.StringValue(device.GetName())
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

	var data ConsolePortResourceModel
	data.ID = types.StringValue(req.ID)
	data.Tags = types.SetNull(types.StringType)
	data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// mapResponseToModel maps the API response to the Terraform model.
func (r *ConsolePortResource) mapResponseToModel(ctx context.Context, consolePort *netbox.ConsolePort, data *ConsolePortResourceModel, diags *diag.Diagnostics) {
	data.ID = types.StringValue(fmt.Sprintf("%d", consolePort.GetId()))
	data.Name = types.StringValue(consolePort.GetName())

	// Map device - preserve user's input format
	if device := consolePort.GetDevice(); device.Id != 0 {
		data.Device = utils.UpdateReferenceAttribute(data.Device, device.GetName(), "", device.GetId())
	}

	// Map label
	if label, ok := consolePort.GetLabelOk(); ok && label != nil && *label != "" {
		data.Label = types.StringValue(*label)
	} else {
		data.Label = types.StringNull()
	}

	// Map type
	if consolePort.Type != nil {
		data.Type = types.StringValue(string(consolePort.Type.GetValue()))
	} else {
		data.Type = types.StringNull()
	}

	// Map speed
	if consolePort.Speed.IsSet() && consolePort.Speed.Get() != nil {
		data.Speed = types.Int32Value(int32(consolePort.Speed.Get().GetValue()))
	} else {
		data.Speed = types.Int32Null()
	}

	// Map description
	if desc, ok := consolePort.GetDescriptionOk(); ok && desc != nil && *desc != "" {
		data.Description = types.StringValue(*desc)
	} else {
		data.Description = types.StringNull()
	}

	// Map mark_connected
	if mc, ok := consolePort.GetMarkConnectedOk(); ok && mc != nil {
		data.MarkConnected = types.BoolValue(*mc)
	} else {
		data.MarkConnected = types.BoolNull()
	}

	// Tags (slug list)
	data.Tags = utils.PopulateTagsSlugFromAPI(ctx, consolePort.HasTags(), consolePort.GetTags(), data.Tags)

	// Handle custom fields
	if consolePort.HasCustomFields() {
		data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, data.CustomFields, consolePort.GetCustomFields(), diags)
	} else {
		data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)
	}
}
