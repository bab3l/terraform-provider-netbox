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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &InventoryItemResource{}
	_ resource.ResourceWithConfigure   = &InventoryItemResource{}
	_ resource.ResourceWithImportState = &InventoryItemResource{}
	_ resource.ResourceWithIdentity    = &InventoryItemResource{}
)

// NewInventoryItemResource returns a new resource implementing the inventory item resource.
func NewInventoryItemResource() resource.Resource {
	return &InventoryItemResource{}
}

// InventoryItemResource defines the resource implementation.
type InventoryItemResource struct {
	client *netbox.APIClient
}

// InventoryItemResourceModel describes the resource data model.
type InventoryItemResourceModel struct {
	ID           types.String `tfsdk:"id"`
	Device       types.String `tfsdk:"device"`
	Name         types.String `tfsdk:"name"`
	Label        types.String `tfsdk:"label"`
	Parent       types.String `tfsdk:"parent"`
	Role         types.String `tfsdk:"role"`
	Manufacturer types.String `tfsdk:"manufacturer"`
	PartID       types.String `tfsdk:"part_id"`
	Serial       types.String `tfsdk:"serial"`
	AssetTag     types.String `tfsdk:"asset_tag"`
	Discovered   types.Bool   `tfsdk:"discovered"`
	Description  types.String `tfsdk:"description"`
	Tags         types.Set    `tfsdk:"tags"`
	CustomFields types.Set    `tfsdk:"custom_fields"`
}

// Metadata returns the resource type name.
func (r *InventoryItemResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_inventory_item"
}

// Schema defines the schema for the resource.
func (r *InventoryItemResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Manages an inventory item in NetBox. Inventory items represent hardware components installed within a device, such as power supplies, CPUs, or line cards.
~> **Deprecation Warning:** Beginning in NetBox v4.3, inventory items are deprecated and planned for removal in a future release. Users are strongly encouraged to use [modules](https://netboxlabs.com/docs/netbox/models/dcim/module/) and [module types](https://netboxlabs.com/docs/netbox/models/dcim/moduletype/) instead.`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique numeric ID of the inventory item.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"device": nbschema.RequiredReferenceAttributeWithDiffSuppress("device", "The device this inventory item belongs to (ID or name)."),
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the inventory item.",
				Required:            true,
			},
			"label": schema.StringAttribute{
				MarkdownDescription: "Physical label on the inventory item.",
				Optional:            true,
			},
			"parent":       nbschema.ReferenceAttributeWithDiffSuppress("inventory_item", "Parent inventory item (ID) for nested items."),
			"role":         nbschema.ReferenceAttributeWithDiffSuppress("inventory_item_role", "The functional role of the inventory item (ID or slug)."),
			"manufacturer": nbschema.ReferenceAttributeWithDiffSuppress("manufacturer", "The manufacturer of the inventory item (ID or slug)."),
			"part_id": schema.StringAttribute{
				MarkdownDescription: "Manufacturer-assigned part identifier.",
				Optional:            true,
			},
			"serial": schema.StringAttribute{
				MarkdownDescription: "Serial number of the inventory item.",
				Optional:            true,
			},
			"asset_tag": schema.StringAttribute{
				MarkdownDescription: "A unique tag used to identify this inventory item.",
				Optional:            true,
			},
			"discovered": schema.BoolAttribute{
				MarkdownDescription: "Whether this item was automatically discovered.",
				Optional:            true,
				Computed:            true,
			},
		},
	}

	// Add description attribute
	maps.Copy(resp.Schema.Attributes, nbschema.DescriptionOnlyAttributes("inventory item"))

	// Add common metadata attributes (tags, custom_fields)
	resp.Schema.Attributes["tags"] = nbschema.TagsSlugAttribute()
	resp.Schema.Attributes["custom_fields"] = nbschema.CustomFieldsAttribute()
}

func (r *InventoryItemResource) IdentitySchema(ctx context.Context, req resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = nbschema.ImportIdentityWithCustomFieldsSchema()
}

func (r *InventoryItemResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *InventoryItemResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data InventoryItemResourceModel
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
	apiReq := netbox.NewInventoryItemRequest(*device, data.Name.ValueString())

	// Set optional fields
	if !data.Label.IsNull() && !data.Label.IsUnknown() {
		apiReq.SetLabel(data.Label.ValueString())
	}

	if !data.Parent.IsNull() && !data.Parent.IsUnknown() {
		parentID, err := utils.ParseID(data.Parent.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Invalid Parent ID",
				fmt.Sprintf("Parent must be a numeric ID, got: %s", data.Parent.ValueString()),
			)
			return
		}
		apiReq.SetParent(parentID)
	}

	if !data.Role.IsNull() && !data.Role.IsUnknown() {
		role, diags := lookup.LookupInventoryItemRole(ctx, r.client, data.Role.ValueString())
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		apiReq.SetRole(*role)
	}

	if !data.Manufacturer.IsNull() && !data.Manufacturer.IsUnknown() {
		manufacturer, diags := lookup.LookupManufacturer(ctx, r.client, data.Manufacturer.ValueString())
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		apiReq.SetManufacturer(*manufacturer)
	}

	// Part ID
	if utils.IsSet(data.PartID) {
		apiReq.SetPartId(data.PartID.ValueString())
	} else if data.PartID.IsNull() {
		apiReq.SetPartId("")
	}

	// Serial
	if utils.IsSet(data.Serial) {
		apiReq.SetSerial(data.Serial.ValueString())
	} else if data.Serial.IsNull() {
		apiReq.SetSerial("")
	}

	// Asset tag
	if utils.IsSet(data.AssetTag) {
		apiReq.SetAssetTag(data.AssetTag.ValueString())
	} else if data.AssetTag.IsNull() {
		apiReq.AssetTag = *netbox.NewNullableString(nil)
	}

	// Discovered
	if utils.IsSet(data.Discovered) {
		apiReq.SetDiscovered(data.Discovered.ValueBool())
	} else if data.Discovered.IsNull() {
		apiReq.Discovered = nil
	}

	// Handle description, tags, and custom fields
	utils.ApplyDescription(apiReq, data.Description)
	utils.ApplyTagsFromSlugs(ctx, r.client, apiReq, data.Tags, &resp.Diagnostics)
	utils.ApplyCustomFields(ctx, apiReq, data.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Creating inventory item", map[string]interface{}{
		"device": data.Device.ValueString(),
		"name":   data.Name.ValueString(),
	})
	response, httpResp, err := r.client.DcimAPI.DcimInventoryItemsCreate(ctx).InventoryItemRequest(*apiReq).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating inventory item",
			utils.FormatAPIError(fmt.Sprintf("create inventory item %s", data.Name.ValueString()), err, httpResp),
		)
		return
	}

	// Map response to model
	r.mapResponseToModel(ctx, response, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(data.ID.ValueString()), data.CustomFields, &resp.Diagnostics)
	tflog.Trace(ctx, "Created inventory item", map[string]interface{}{
		"id":   data.ID.ValueString(),
		"name": data.Name.ValueString(),
	})
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read refreshes the resource state.
func (r *InventoryItemResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data InventoryItemResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	itemID, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Inventory Item ID",
			fmt.Sprintf("Inventory Item ID must be a number, got: %s", data.ID.ValueString()),
		)
		return
	}
	tflog.Debug(ctx, "Reading inventory item", map[string]interface{}{
		"id": itemID,
	})

	response, httpResp, err := r.client.DcimAPI.DcimInventoryItemsRetrieve(ctx, itemID).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading inventory item",
			utils.FormatAPIError(fmt.Sprintf("read inventory item ID %d", itemID), err, httpResp),
		)
		return
	}

	// Preserve original custom_fields state
	originalCustomFields := data.CustomFields

	// Map response to model
	r.mapResponseToModel(ctx, response, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Restore original custom_fields if it was null/empty and API returned none
	if !utils.IsSet(originalCustomFields) && !utils.IsSet(data.CustomFields) {
		data.CustomFields = originalCustomFields
	}
	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(data.ID.ValueString()), data.CustomFields, &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates the resource.
func (r *InventoryItemResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, data InventoryItemResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	itemID, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Inventory Item ID",
			fmt.Sprintf("Inventory Item ID must be a number, got: %s", data.ID.ValueString()),
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
	apiReq := netbox.NewInventoryItemRequest(*device, data.Name.ValueString())

	// Set optional fields
	utils.ApplyLabel(apiReq, data.Label)

	if !data.Parent.IsNull() && !data.Parent.IsUnknown() {
		parentID, err := utils.ParseID(data.Parent.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Invalid Parent ID",
				fmt.Sprintf("Parent must be a numeric ID, got: %s", data.Parent.ValueString()),
			)
			return
		}
		apiReq.SetParent(parentID)
	}

	if !data.Role.IsNull() && !data.Role.IsUnknown() {
		role, diags := lookup.LookupInventoryItemRole(ctx, r.client, data.Role.ValueString())
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		apiReq.SetRole(*role)
	}

	if !data.Manufacturer.IsNull() && !data.Manufacturer.IsUnknown() {
		manufacturer, diags := lookup.LookupManufacturer(ctx, r.client, data.Manufacturer.ValueString())
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		apiReq.SetManufacturer(*manufacturer)
	}

	// Part ID
	if utils.IsSet(data.PartID) {
		apiReq.SetPartId(data.PartID.ValueString())
	} else if data.PartID.IsNull() {
		apiReq.SetPartId("")
	}

	// Serial
	if utils.IsSet(data.Serial) {
		apiReq.SetSerial(data.Serial.ValueString())
	} else if data.Serial.IsNull() {
		apiReq.SetSerial("")
	}

	// Asset tag
	if utils.IsSet(data.AssetTag) {
		apiReq.SetAssetTag(data.AssetTag.ValueString())
	} else if data.AssetTag.IsNull() {
		apiReq.AssetTag = *netbox.NewNullableString(nil)
	}

	// Discovered
	if utils.IsSet(data.Discovered) {
		apiReq.SetDiscovered(data.Discovered.ValueBool())
	} else if data.Discovered.IsNull() {
		apiReq.Discovered = nil
	}

	// Handle description, tags, and custom fields with merge-aware behavior
	utils.ApplyDescription(apiReq, data.Description)

	// Handle tags - merge-aware: use plan if provided, else use state
	if utils.IsSet(data.Tags) {
		utils.ApplyTagsFromSlugs(ctx, r.client, apiReq, data.Tags, &resp.Diagnostics)
	} else if utils.IsSet(state.Tags) {
		utils.ApplyTagsFromSlugs(ctx, r.client, apiReq, state.Tags, &resp.Diagnostics)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	// Apply custom fields with merge logic (preserves unmanaged fields from state)
	utils.ApplyCustomFieldsWithMerge(ctx, apiReq, data.CustomFields, state.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Updating inventory item", map[string]interface{}{
		"id": itemID,
	})
	response, httpResp, err := r.client.DcimAPI.DcimInventoryItemsUpdate(ctx, itemID).InventoryItemRequest(*apiReq).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating inventory item",
			utils.FormatAPIError(fmt.Sprintf("update inventory item ID %d", itemID), err, httpResp),
		)
		return
	}

	r.mapResponseToModel(ctx, response, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	utils.SetIdentityCustomFields(ctx, resp.Identity, types.StringValue(data.ID.ValueString()), data.CustomFields, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete deletes the resource.

func (r *InventoryItemResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data InventoryItemResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	itemID, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Inventory Item ID",
			fmt.Sprintf("Inventory Item ID must be a number, got: %s", data.ID.ValueString()),
		)
		return
	}
	tflog.Debug(ctx, "Deleting inventory item", map[string]interface{}{
		"id": itemID,
	})

	httpResp, err := r.client.DcimAPI.DcimInventoryItemsDestroy(ctx, itemID).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			return
		}
		resp.Diagnostics.AddError(
			"Error deleting inventory item",
			utils.FormatAPIError(fmt.Sprintf("delete inventory item ID %d", itemID), err, httpResp),
		)
		return
	}
}

// ImportState imports an existing resource.
func (r *InventoryItemResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if parsed, ok := utils.ParseImportIdentityCustomFields(ctx, req.Identity, &resp.Diagnostics); ok {
		if resp.Diagnostics.HasError() {
			return
		}
		if parsed.ID == "" {
			resp.Diagnostics.AddError("Invalid import identity", "Identity id must be provided")
			return
		}

		itemID, err := utils.ParseID(parsed.ID)
		if err != nil {
			resp.Diagnostics.AddError(
				"Invalid Import ID",
				fmt.Sprintf("Inventory Item ID must be a number, got: %s", parsed.ID),
			)
			return
		}

		response, httpResp, err := r.client.DcimAPI.DcimInventoryItemsRetrieve(ctx, itemID).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error importing inventory item",
				utils.FormatAPIError(fmt.Sprintf("import inventory item ID %d", itemID), err, httpResp),
			)
			return
		}

		var data InventoryItemResourceModel
		if device := response.GetDevice(); device.Id != 0 {
			data.Device = types.StringValue(device.GetName())
		}
		if response.HasTags() && len(response.GetTags()) > 0 {
			tagSlugs := make([]string, 0, len(response.GetTags()))
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

	itemID, err := utils.ParseID(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Inventory Item ID must be a number, got: %s", req.ID),
		)
		return
	}

	response, httpResp, err := r.client.DcimAPI.DcimInventoryItemsRetrieve(ctx, itemID).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error importing inventory item",
			utils.FormatAPIError(fmt.Sprintf("import inventory item ID %d", itemID), err, httpResp),
		)
		return
	}

	var data InventoryItemResourceModel
	r.mapResponseToModel(ctx, response, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// mapResponseToModel maps the API response to the Terraform model.
func (r *InventoryItemResource) mapResponseToModel(ctx context.Context, item *netbox.InventoryItem, data *InventoryItemResourceModel, diags *diag.Diagnostics) {
	data.ID = types.StringValue(fmt.Sprintf("%d", item.GetId()))
	data.Name = types.StringValue(item.GetName())

	// Map device - preserve user's input format
	if device := item.GetDevice(); device.Id != 0 {
		data.Device = utils.UpdateReferenceAttribute(data.Device, device.GetName(), "", device.GetId())
	}

	// Map label
	if label, ok := item.GetLabelOk(); ok && label != nil && *label != "" {
		data.Label = types.StringValue(*label)
	} else {
		data.Label = types.StringNull()
	}

	// Map parent (NullableInt32 - just an ID, not a nested object)
	if item.Parent.IsSet() && item.Parent.Get() != nil {
		data.Parent = types.StringValue(fmt.Sprintf("%d", *item.Parent.Get()))
	} else {
		data.Parent = types.StringNull()
	}

	// Map role - preserve user's input format
	if item.Role.IsSet() && item.Role.Get() != nil {
		role := item.Role.Get()
		data.Role = utils.UpdateReferenceAttribute(data.Role, role.GetName(), role.GetSlug(), role.GetId())
	} else {
		data.Role = types.StringNull()
	}

	// Map manufacturer - preserve user's input format
	if item.Manufacturer.IsSet() && item.Manufacturer.Get() != nil {
		mfr := item.Manufacturer.Get()
		data.Manufacturer = utils.UpdateReferenceAttribute(data.Manufacturer, mfr.GetName(), mfr.GetSlug(), mfr.GetId())
	} else {
		data.Manufacturer = types.StringNull()
	}

	// Map part_id
	if partID, ok := item.GetPartIdOk(); ok && partID != nil && *partID != "" {
		data.PartID = types.StringValue(*partID)
	} else {
		data.PartID = types.StringNull()
	}

	// Map serial
	if serial, ok := item.GetSerialOk(); ok && serial != nil && *serial != "" {
		data.Serial = types.StringValue(*serial)
	} else {
		data.Serial = types.StringNull()
	}

	// Map asset_tag
	if item.AssetTag.IsSet() && item.AssetTag.Get() != nil && *item.AssetTag.Get() != "" {
		data.AssetTag = types.StringValue(*item.AssetTag.Get())
	} else {
		data.AssetTag = types.StringNull()
	}

	// Map discovered
	if discovered, ok := item.GetDiscoveredOk(); ok && discovered != nil {
		data.Discovered = types.BoolValue(*discovered)
	} else {
		data.Discovered = types.BoolNull()
	}

	// Map description
	if desc, ok := item.GetDescriptionOk(); ok && desc != nil && *desc != "" {
		data.Description = types.StringValue(*desc)
	} else {
		data.Description = types.StringNull()
	}

	// Handle tags with filter-to-owned pattern
	planTags := data.Tags
	wasExplicitlyEmpty := !planTags.IsNull() && !planTags.IsUnknown() && len(planTags.Elements()) == 0
	switch {
	case item.HasTags() && len(item.GetTags()) > 0:
		tagSlugs := make([]string, 0, len(item.GetTags()))
		for _, tag := range item.GetTags() {
			tagSlugs = append(tagSlugs, tag.GetSlug())
		}
		data.Tags = utils.TagsSlugToSet(ctx, tagSlugs)
	case wasExplicitlyEmpty:
		data.Tags = types.SetValueMust(types.StringType, []attr.Value{})
	default:
		data.Tags = types.SetNull(types.StringType)
	}

	// Handle custom fields - use filtered-to-owned for partial management
	if item.HasCustomFields() {
		data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, data.CustomFields, item.GetCustomFields(), diags)
	}
}
