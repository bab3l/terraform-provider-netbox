// Package resources contains Terraform resource implementations for NetBox objects.

package resources

import (
	"context"
	"fmt"
	"maps"
	"net/http"

	"github.com/bab3l/go-netbox"
	"github.com/bab3l/terraform-provider-netbox/internal/netboxlookup"
	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &DeviceBayResource{}
	_ resource.ResourceWithConfigure   = &DeviceBayResource{}
	_ resource.ResourceWithImportState = &DeviceBayResource{}
)

// NewDeviceBayResource returns a new resource implementing the DeviceBay resource.
func NewDeviceBayResource() resource.Resource {
	return &DeviceBayResource{}
}

// DeviceBayResource defines the resource implementation.
type DeviceBayResource struct {
	client *netbox.APIClient
}

// DeviceBayResourceModel describes the resource data model.
type DeviceBayResourceModel struct {
	ID              types.String `tfsdk:"id"`
	Device          types.String `tfsdk:"device"`
	Name            types.String `tfsdk:"name"`
	Label           types.String `tfsdk:"label"`
	Description     types.String `tfsdk:"description"`
	InstalledDevice types.String `tfsdk:"installed_device"`
	Tags            types.Set    `tfsdk:"tags"`
	CustomFields    types.Set    `tfsdk:"custom_fields"`
}

// Metadata returns the resource type name.
func (r *DeviceBayResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_device_bay"
}

// Schema defines the schema for the resource.
func (r *DeviceBayResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a device bay in NetBox. A device bay is a slot within a parent device where a child device can be installed, such as a blade server chassis slot or a modular switch slot.",
		Attributes: map[string]schema.Attribute{
			"id":     nbschema.IDAttribute("device bay"),
			"device": nbschema.RequiredReferenceAttribute("device", "The parent device containing this device bay. Accepts ID or name."),
			"name":   nbschema.NameAttribute("device bay", 64),
			"label": schema.StringAttribute{
				MarkdownDescription: "Physical label for the device bay.",
				Optional:            true,
			},
			"installed_device": schema.StringAttribute{
				MarkdownDescription: "The child device installed in this bay. Accepts ID or name.",
				Optional:            true,
			},
			"tags":          nbschema.TagsSlugAttribute(),
			"custom_fields": nbschema.CustomFieldsAttribute(),
		},
	}

	// Add description attribute
	maps.Copy(resp.Schema.Attributes, nbschema.DescriptionOnlyAttributes("device bay"))

	// Tags and custom fields are defined directly in the schema above.
}

// Configure adds the provider configured client to the resource.
func (r *DeviceBayResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create creates a new device bay resource.
func (r *DeviceBayResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data DeviceBayResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the request
	dbRequest, diags := r.buildRequest(ctx, &data, nil)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Creating device bay", map[string]interface{}{
		"name":   data.Name.ValueString(),
		"device": data.Device.ValueString(),
	})
	db, httpResp, err := r.client.DcimAPI.DcimDeviceBaysCreate(ctx).DeviceBayRequest(*dbRequest).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating device bay",
			utils.FormatAPIError("create device bay", err, httpResp),
		)
		return
	}

	// Map response to state
	r.mapResponseToModel(ctx, db, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Created device bay", map[string]interface{}{
		"id":   db.GetId(),
		"name": db.GetName(),
	})
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read reads the device bay resource.
func (r *DeviceBayResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data DeviceBayResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	dbID, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Device Bay ID",
			fmt.Sprintf("Could not parse device bay ID: %s", err),
		)
		return
	}
	tflog.Debug(ctx, "Reading device bay", map[string]interface{}{
		"id": dbID,
	})

	db, httpResp, err := r.client.DcimAPI.DcimDeviceBaysRetrieve(ctx, dbID).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading device bay",
			utils.FormatAPIError(fmt.Sprintf("read device bay ID %d", dbID), err, httpResp),
		)
		return
	}

	// Map response to state
	r.mapResponseToModel(ctx, db, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates the device bay resource.
func (r *DeviceBayResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, data DeviceBayResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	dbID, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Device Bay ID",
			fmt.Sprintf("Could not parse device bay ID: %s", err),
		)
		return
	}

	// Build the request
	dbRequest, diags := r.buildRequest(ctx, &data, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Updating device bay", map[string]interface{}{
		"id": dbID,
	})
	db, httpResp, err := r.client.DcimAPI.DcimDeviceBaysUpdate(ctx, dbID).DeviceBayRequest(*dbRequest).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating device bay",
			utils.FormatAPIError(fmt.Sprintf("update device bay ID %d", dbID), err, httpResp),
		)
		return
	}

	// Map response to state
	r.mapResponseToModel(ctx, db, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete deletes the device bay resource.
func (r *DeviceBayResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data DeviceBayResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	dbID, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Device Bay ID",
			fmt.Sprintf("Could not parse device bay ID: %s", err),
		)
		return
	}
	tflog.Debug(ctx, "Deleting device bay", map[string]interface{}{
		"id": dbID,
	})

	httpResp, err := r.client.DcimAPI.DcimDeviceBaysDestroy(ctx, dbID).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			return
		}
		resp.Diagnostics.AddError(
			"Error deleting device bay",
			utils.FormatAPIError(fmt.Sprintf("delete device bay ID %d", dbID), err, httpResp),
		)
		return
	}
}

// ImportState imports an existing device bay resource.
func (r *DeviceBayResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// buildRequest builds the API request from the Terraform model with merge-aware custom fields.
// state may be nil for Create operations, in which case only plan fields are used.
func (r *DeviceBayResource) buildRequest(ctx context.Context, data, state *DeviceBayResourceModel) (*netbox.DeviceBayRequest, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Look up device
	deviceRef, lookupDiags := netboxlookup.LookupDevice(ctx, r.client, data.Device.ValueString())
	diags.Append(lookupDiags...)
	if diags.HasError() {
		return nil, diags
	}
	dbRequest := netbox.NewDeviceBayRequest(*deviceRef, data.Name.ValueString())

	// Set optional fields
	utils.ApplyLabel(dbRequest, data.Label)

	if !data.InstalledDevice.IsNull() && !data.InstalledDevice.IsUnknown() {
		installedDeviceRef, lookupDiags := netboxlookup.LookupDevice(ctx, r.client, data.InstalledDevice.ValueString())
		diags.Append(lookupDiags...)
		if diags.HasError() {
			return nil, diags
		}
		dbRequest.InstalledDevice = *netbox.NewNullableBriefDeviceRequest(installedDeviceRef)
	} else if data.InstalledDevice.IsNull() {
		// Explicitly set to nil when removed from config
		dbRequest.SetInstalledDeviceNil()
	}

	// Handle description, tags, and custom fields with merge-aware behavior
	utils.ApplyDescription(dbRequest, data.Description)

	// Handle tags - merge-aware: use plan if provided, else use state (if available)
	if utils.IsSet(data.Tags) {
		utils.ApplyTagsFromSlugs(ctx, r.client, dbRequest, data.Tags, &diags)
	} else if state != nil && utils.IsSet(state.Tags) {
		utils.ApplyTagsFromSlugs(ctx, r.client, dbRequest, state.Tags, &diags)
	}
	if diags.HasError() {
		return nil, diags
	}

	// Apply custom fields with merge logic (preserves unmanaged fields from state)
	var stateCustomFields types.Set
	if state != nil {
		stateCustomFields = state.CustomFields
	} else {
		// For Create, there's no state to merge with
		stateCustomFields = types.SetNull(types.StringType)
	}
	utils.ApplyCustomFieldsWithMerge(ctx, dbRequest, data.CustomFields, stateCustomFields, &diags)
	if diags.HasError() {
		return nil, diags
	}
	return dbRequest, diags
}

// mapResponseToModel maps the API response to the Terraform model.
func (r *DeviceBayResource) mapResponseToModel(ctx context.Context, db *netbox.DeviceBay, data *DeviceBayResourceModel, diags *diag.Diagnostics) {
	data.ID = types.StringValue(fmt.Sprintf("%d", db.GetId()))
	data.Name = types.StringValue(db.GetName())

	// Map device - preserve user's input format
	data.Device = utils.UpdateReferenceAttribute(data.Device, db.Device.GetName(), "", db.Device.GetId())

	// Map label
	if label, ok := db.GetLabelOk(); ok && label != nil && *label != "" {
		data.Label = types.StringValue(*label)
	} else {
		data.Label = types.StringNull()
	}

	// Map description
	if desc, ok := db.GetDescriptionOk(); ok && desc != nil && *desc != "" {
		data.Description = types.StringValue(*desc)
	} else {
		data.Description = types.StringNull()
	}

	// Map installed_device - preserve user's input format
	if db.InstalledDevice.IsSet() && db.InstalledDevice.Get() != nil {
		installedDevice := db.InstalledDevice.Get()

		data.InstalledDevice = utils.UpdateReferenceAttribute(data.InstalledDevice, installedDevice.GetName(), "", installedDevice.GetId())
	} else {
		data.InstalledDevice = types.StringNull()
	}

	// Tags (slug list)
	wasExplicitlyEmpty := !data.Tags.IsNull() && !data.Tags.IsUnknown() && len(data.Tags.Elements()) == 0
	switch {
	case db.HasTags() && len(db.GetTags()) > 0:
		tagSlugs := make([]string, 0, len(db.GetTags()))
		for _, tag := range db.GetTags() {
			tagSlugs = append(tagSlugs, tag.GetSlug())
		}
		data.Tags = utils.TagsSlugToSet(ctx, tagSlugs)
	case wasExplicitlyEmpty:
		data.Tags = types.SetValueMust(types.StringType, []attr.Value{})
	default:
		data.Tags = types.SetNull(types.StringType)
	}

	// Handle custom fields - use filtered-to-owned for partial management
	if db.HasCustomFields() {
		data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, data.CustomFields, db.GetCustomFields(), diags)
	}
}
