// Package resources provides Terraform resource implementations for NetBox objects.

package resources

import (
	"context"
	"fmt"
	"maps"

	"github.com/bab3l/go-netbox"
	lookup "github.com/bab3l/terraform-provider-netbox/internal/netboxlookup"
	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.

var (
	_ resource.Resource = &RearPortResource{}

	_ resource.ResourceWithConfigure = &RearPortResource{}

	_ resource.ResourceWithImportState = &RearPortResource{}
	_ resource.ResourceWithIdentity    = &RearPortResource{}
)

// NewRearPortResource returns a new resource implementing the rear port resource.

func NewRearPortResource() resource.Resource {
	return &RearPortResource{}
}

// RearPortResource defines the resource implementation.

type RearPortResource struct {
	client *netbox.APIClient
}

// RearPortResourceModel describes the resource data model.

type RearPortResourceModel struct {
	ID types.String `tfsdk:"id"`

	Device types.String `tfsdk:"device"`

	Name types.String `tfsdk:"name"`

	Label types.String `tfsdk:"label"`

	Type types.String `tfsdk:"type"`

	Color types.String `tfsdk:"color"`

	Positions types.Int32 `tfsdk:"positions"`

	Description types.String `tfsdk:"description"`

	MarkConnected types.Bool `tfsdk:"mark_connected"`

	Tags types.Set `tfsdk:"tags"`

	CustomFields types.Set `tfsdk:"custom_fields"`
}

// Metadata returns the resource type name.

func (r *RearPortResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_rear_port"
}

// Schema defines the schema for the resource.

func (r *RearPortResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a rear port in NetBox. Rear ports represent physical ports on the back of a device, typically used for patch panels and fiber distribution.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique numeric ID of the rear port.",

				Computed: true,

				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"device": nbschema.RequiredReferenceAttributeWithDiffSuppress(
				"device",
				"The device this rear port belongs to (ID or name).",
			),

			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the rear port.",

				Required: true,
			},

			"label": schema.StringAttribute{
				MarkdownDescription: "Physical label of the rear port.",

				Optional: true,
			},

			"type": schema.StringAttribute{
				MarkdownDescription: "The type of rear port (e.g., `8p8c`, `8p6c`, `110-punch`, `bnc`, `f`, `n`, `mrj21`, `fc`, `lc`, `lc-pc`, `lc-upc`, `lc-apc`, `lsh`, `mpo`, `mtrj`, `sc`, `sc-pc`, `sc-upc`, `sc-apc`, `st`, `cs`, `sn`, `splice`, `other`).",

				Required: true,
			},

			"color": schema.StringAttribute{
				MarkdownDescription: "Color of the rear port in hex format (e.g., `aa1409`).",

				Optional: true,
			},

			"positions": schema.Int32Attribute{
				MarkdownDescription: "Number of front ports that may be mapped to this rear port (1-1024). Default is 1.",

				Optional: true,

				Computed: true,

				Default: int32default.StaticInt32(1),
			},

			"mark_connected": schema.BoolAttribute{
				MarkdownDescription: "Treat as if a cable is connected.",

				Optional: true,

				Computed: true,

				Default: booldefault.StaticBool(false),
			},
		},
	}

	// Add description attribute
	maps.Copy(resp.Schema.Attributes, nbschema.DescriptionOnlyAttributes("rear port"))

	// Add common metadata attributes (tags, custom_fields)
	resp.Schema.Attributes["tags"] = nbschema.TagsSlugAttribute()
	resp.Schema.Attributes["custom_fields"] = nbschema.CustomFieldsAttribute()
}

func (r *RearPortResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = nbschema.ImportIdentityWithCustomFieldsSchema()
}

// Configure adds the provider configured client to the resource.

func (r *RearPortResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *RearPortResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data RearPortResourceModel

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

	apiReq := netbox.NewWritableRearPortRequest(*device, data.Name.ValueString(), netbox.FrontPortTypeValue(data.Type.ValueString()))

	// Set optional fields
	utils.ApplyLabel(apiReq, data.Label)

	if !data.Color.IsNull() && !data.Color.IsUnknown() {
		apiReq.SetColor(data.Color.ValueString())
	}

	if !data.Positions.IsNull() && !data.Positions.IsUnknown() {
		apiReq.SetPositions(data.Positions.ValueInt32())
	}

	if !data.MarkConnected.IsNull() && !data.MarkConnected.IsUnknown() {
		apiReq.SetMarkConnected(data.MarkConnected.ValueBool())
	}

	// Handle description, tags, and custom fields
	utils.ApplyDescription(apiReq, data.Description)
	utils.ApplyTagsFromSlugs(ctx, r.client, apiReq, data.Tags, &resp.Diagnostics)
	utils.ApplyCustomFields(ctx, apiReq, data.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating rear port", map[string]interface{}{
		"device": data.Device.ValueString(),

		"name": data.Name.ValueString(),
	})

	response, httpResp, err := r.client.DcimAPI.DcimRearPortsCreate(ctx).WritableRearPortRequest(*apiReq).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {
		resp.Diagnostics.AddError(

			"Error creating rear port",

			utils.FormatAPIError("create rear port", err, httpResp),
		)

		return
	}

	// Map response to model

	r.mapResponseToModel(ctx, response, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(data.ID.ValueString()), data.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "Created rear port", map[string]interface{}{
		"id": data.ID.ValueString(),
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read retrieves the resource.

func (r *RearPortResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data RearPortResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	portID, err := utils.ParseID(data.ID.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(

			"Invalid ID",

			fmt.Sprintf("Could not parse ID %q: %s", data.ID.ValueString(), err),
		)

		return
	}

	tflog.Debug(ctx, "Reading rear port", map[string]interface{}{
		"id": portID,
	})

	response, httpResp, err := r.client.DcimAPI.DcimRearPortsRetrieve(ctx, portID).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			tflog.Debug(ctx, "Rear port not found, removing from state", map[string]interface{}{
				"id": portID,
			})

			resp.State.RemoveResource(ctx)

			return
		}

		resp.Diagnostics.AddError(

			"Error reading rear port",

			utils.FormatAPIError(fmt.Sprintf("read rear port ID %d", portID), err, httpResp),
		)

		return
	}

	// Map response to model

	r.mapResponseToModel(ctx, response, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(data.ID.ValueString()), data.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates the resource.

func (r *RearPortResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Read both state and plan for merge-aware custom fields handling
	var state, plan RearPortResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	portID, err := utils.ParseID(plan.ID.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(

			"Invalid ID",

			fmt.Sprintf("Could not parse ID %q: %s", plan.ID.ValueString(), err),
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

	apiReq := netbox.NewWritableRearPortRequest(*device, plan.Name.ValueString(), netbox.FrontPortTypeValue(plan.Type.ValueString()))

	// Set optional fields
	utils.ApplyLabel(apiReq, plan.Label)

	// For nullable string fields, explicitly clear if null in plan
	if plan.Color.IsNull() {
		apiReq.SetColor("")
	} else if !plan.Color.IsUnknown() {
		apiReq.SetColor(plan.Color.ValueString())
	}

	if !plan.Positions.IsNull() && !plan.Positions.IsUnknown() {
		apiReq.SetPositions(plan.Positions.ValueInt32())
	}

	if !plan.MarkConnected.IsNull() && !plan.MarkConnected.IsUnknown() {
		apiReq.SetMarkConnected(plan.MarkConnected.ValueBool())
	}

	// Handle description, tags, and custom fields with merge-aware helpers
	utils.ApplyDescription(apiReq, plan.Description)
	utils.ApplyTagsFromSlugs(ctx, r.client, apiReq, plan.Tags, &resp.Diagnostics)
	utils.ApplyCustomFieldsWithMerge(ctx, apiReq, plan.CustomFields, state.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Updating rear port", map[string]interface{}{
		"id": portID,
	})

	response, httpResp, err := r.client.DcimAPI.DcimRearPortsUpdate(ctx, portID).WritableRearPortRequest(*apiReq).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {
		resp.Diagnostics.AddError(

			"Error updating rear port",

			utils.FormatAPIError(fmt.Sprintf("update rear port ID %d", portID), err, httpResp),
		)

		return
	}

	// Map response to model

	r.mapResponseToModel(ctx, response, &plan, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(plan.ID.ValueString()), plan.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete removes the resource.

func (r *RearPortResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data RearPortResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	portID, err := utils.ParseID(data.ID.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(

			"Invalid ID",

			fmt.Sprintf("Could not parse ID %q: %s", data.ID.ValueString(), err),
		)

		return
	}

	tflog.Debug(ctx, "Deleting rear port", map[string]interface{}{
		"id": portID,
	})

	httpResp, err := r.client.DcimAPI.DcimRearPortsDestroy(ctx, portID).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			// Resource already deleted

			return
		}

		resp.Diagnostics.AddError(

			"Error deleting rear port",

			utils.FormatAPIError(fmt.Sprintf("delete rear port ID %d", portID), err, httpResp),
		)

		return
	}
}

// ImportState imports the resource.

func (r *RearPortResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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
			resp.Diagnostics.AddError("Invalid Import ID", fmt.Sprintf("Rear Port ID must be a number, got: %s", parsed.ID))
			return
		}

		response, httpResp, err := r.client.DcimAPI.DcimRearPortsRetrieve(ctx, portID).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError("Error importing rear port", utils.FormatAPIError(fmt.Sprintf("import rear port ID %d", portID), err, httpResp))
			return
		}

		var data RearPortResourceModel
		if response.HasTags() {
			var tagSlugs []string
			for _, tag := range response.GetTags() {
				tagSlugs = append(tagSlugs, tag.GetSlug())
			}
			data.Tags = utils.TagsSlugToSet(ctx, tagSlugs)
		} else {
			data.Tags = types.SetNull(types.StringType)
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

	utils.ImportStatePassthroughIDWithValidation(ctx, req, resp, path.Root("id"), true)
}

// mapResponseToModel maps the API response to the Terraform model.

func (r *RearPortResource) mapResponseToModel(ctx context.Context, port *netbox.RearPort, data *RearPortResourceModel, diags *diag.Diagnostics) {
	data.ID = types.StringValue(fmt.Sprintf("%d", port.GetId()))

	data.Name = types.StringValue(port.GetName())

	// Map device - preserve user's input format

	if device := port.GetDevice(); device.Id != 0 {
		data.Device = utils.UpdateReferenceAttribute(data.Device, device.GetName(), "", device.GetId())
	}

	// Map type

	data.Type = types.StringValue(string(port.Type.GetValue()))

	// Map label

	if label, ok := port.GetLabelOk(); ok && label != nil && *label != "" {
		data.Label = types.StringValue(*label)
	} else {
		data.Label = types.StringNull()
	}

	// Map color

	if color, ok := port.GetColorOk(); ok && color != nil && *color != "" {
		data.Color = types.StringValue(*color)
	} else {
		data.Color = types.StringNull()
	}

	// Map positions - always set since it's computed (defaults to 1)

	if positions, ok := port.GetPositionsOk(); ok && positions != nil {
		data.Positions = types.Int32Value(*positions)
	} else {
		data.Positions = types.Int32Value(1)
	}

	// Map description

	if desc, ok := port.GetDescriptionOk(); ok && desc != nil && *desc != "" {
		data.Description = types.StringValue(*desc)
	} else {
		data.Description = types.StringNull()
	}

	// Map mark_connected

	if markConnected, ok := port.GetMarkConnectedOk(); ok && markConnected != nil {
		data.MarkConnected = types.BoolValue(*markConnected)
	} else {
		data.MarkConnected = types.BoolValue(false)
	}

	// Handle tags - filter to owned slugs only
	switch {
	case data.Tags.IsNull():
		data.Tags = types.SetNull(types.StringType)
	case len(data.Tags.Elements()) == 0:
		data.Tags = types.SetValueMust(types.StringType, []attr.Value{})
	case port.HasTags():
		var tagSlugs []string
		for _, tag := range port.GetTags() {
			tagSlugs = append(tagSlugs, tag.GetSlug())
		}
		data.Tags = utils.TagsSlugToSet(ctx, tagSlugs)
	default:
		data.Tags = types.SetValueMust(types.StringType, []attr.Value{})
	}

	// Handle custom fields with filter-to-owned pattern
	data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, data.CustomFields, port.GetCustomFields(), diags)
}
