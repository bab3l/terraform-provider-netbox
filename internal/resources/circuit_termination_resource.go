// Package resources contains Terraform resource implementations for the Netbox provider.

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
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.

var (
	_ resource.Resource                = &CircuitTerminationResource{}
	_ resource.ResourceWithConfigure   = &CircuitTerminationResource{}
	_ resource.ResourceWithImportState = &CircuitTerminationResource{}
)

// NewCircuitTerminationResource returns a new Circuit Termination resource.
func NewCircuitTerminationResource() resource.Resource {
	return &CircuitTerminationResource{}
}

// CircuitTerminationResource defines the resource implementation.
type CircuitTerminationResource struct {
	client *netbox.APIClient
}

// CircuitTerminationResourceModel describes the resource data model.
type CircuitTerminationResourceModel struct {
	ID              types.String `tfsdk:"id"`
	Circuit         types.String `tfsdk:"circuit"`
	TermSide        types.String `tfsdk:"term_side"`
	Site            types.String `tfsdk:"site"`
	ProviderNetwork types.String `tfsdk:"provider_network"`
	PortSpeed       types.Int64  `tfsdk:"port_speed"`
	UpstreamSpeed   types.Int64  `tfsdk:"upstream_speed"`
	XconnectID      types.String `tfsdk:"xconnect_id"`
	PPInfo          types.String `tfsdk:"pp_info"`
	Description     types.String `tfsdk:"description"`
	MarkConnected   types.Bool   `tfsdk:"mark_connected"`
	Tags            types.Set    `tfsdk:"tags"`
	CustomFields    types.Set    `tfsdk:"custom_fields"`
}

// Metadata returns the resource type name.
func (r *CircuitTerminationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_circuit_termination"
}

// Schema defines the schema for the resource.
func (r *CircuitTerminationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a circuit termination in Netbox. Circuit terminations represent the physical endpoints of a circuit at either the A-side or Z-side.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique numeric ID of the circuit termination.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"circuit": schema.StringAttribute{
				MarkdownDescription: "The ID or CID (circuit identifier) of the circuit this termination belongs to.",
				Required:            true,
			},
			"term_side": schema.StringAttribute{
				MarkdownDescription: "The termination side. Valid values are `A` or `Z`.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("A", "Z"),
				},
			},
			"site": schema.StringAttribute{
				MarkdownDescription: "The name, slug, or ID of the site where this termination is located.",
				Optional:            true,
			},
			"provider_network": schema.StringAttribute{
				MarkdownDescription: "The ID of the provider network for this termination.",
				Optional:            true,
			},
			"port_speed": schema.Int64Attribute{
				MarkdownDescription: "The physical circuit speed in Kbps.",
				Optional:            true,
			},
			"upstream_speed": schema.Int64Attribute{
				MarkdownDescription: "The upstream speed in Kbps, if different from port speed.",
				Optional:            true,
			},
			"xconnect_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the local cross-connect.",
				Optional:            true,
			},
			"pp_info": schema.StringAttribute{
				MarkdownDescription: "Patch panel ID and port number(s).",
				Optional:            true,
			},
			"mark_connected": schema.BoolAttribute{
				MarkdownDescription: "Treat as if a cable is connected. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
		},
	}

	// Add description attribute
	maps.Copy(resp.Schema.Attributes, nbschema.DescriptionOnlyAttributes("circuit termination"))

	// Add metadata attributes (slug list tags, custom_fields)
	resp.Schema.Attributes["tags"] = nbschema.TagsSlugAttribute()
	resp.Schema.Attributes["custom_fields"] = nbschema.CustomFieldsAttribute()
}

// Configure sets the client for the resource.
func (r *CircuitTerminationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
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

// Create creates a new circuit termination resource.
func (r *CircuitTerminationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data CircuitTerminationResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the create request
	createReq, diags := r.buildCreateRequest(ctx, &data, nil)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Creating circuit termination", map[string]interface{}{
		"circuit":   data.Circuit.ValueString(),
		"term_side": data.TermSide.ValueString(),
	})

	// Call API to create circuit termination
	termination, httpResp, err := r.client.CircuitsAPI.CircuitsCircuitTerminationsCreate(ctx).CircuitTerminationRequest(*createReq).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating circuit termination",
			fmt.Sprintf("Could not create circuit termination: %s\nHTTP Response: %v", err.Error(), httpResp),
		)
		return
	}

	// Map response to model
	r.mapResponseToModel(ctx, termination, &data, &resp.Diagnostics)
	tflog.Debug(ctx, "Created circuit termination", map[string]interface{}{
		"id": data.ID.ValueString(),
	})

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read reads the circuit termination resource.
func (r *CircuitTerminationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data CircuitTerminationResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid ID",
			fmt.Sprintf("Could not convert ID to integer: %s", err.Error()),
		)
		return
	}
	tflog.Debug(ctx, "Reading circuit termination", map[string]interface{}{
		"id": id,
	})

	// Call API to read circuit termination
	termination, httpResp, err := r.client.CircuitsAPI.CircuitsCircuitTerminationsRetrieve(ctx, id).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			tflog.Debug(ctx, "Circuit termination not found, removing from state", map[string]interface{}{
				"id": id,
			})
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading circuit termination",
			fmt.Sprintf("Could not read circuit termination: %s\nHTTP Response: %v", err.Error(), httpResp),
		)
		return
	}

	// Map response to model
	r.mapResponseToModel(ctx, termination, &data, &resp.Diagnostics)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates the circuit termination resource.
func (r *CircuitTerminationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, data CircuitTerminationResourceModel

	// Read Terraform plan and state data into the models
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid ID",
			fmt.Sprintf("Could not convert ID to integer: %s", err.Error()),
		)
		return
	}

	// Build the update request (pass state for merge-aware custom fields)
	updateReq, diags := r.buildCreateRequest(ctx, &data, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Updating circuit termination", map[string]interface{}{
		"id": id,
	})

	// Call API to update circuit termination
	termination, httpResp, err := r.client.CircuitsAPI.CircuitsCircuitTerminationsUpdate(ctx, id).CircuitTerminationRequest(*updateReq).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating circuit termination",
			fmt.Sprintf("Could not update circuit termination: %s\nHTTP Response: %v", err.Error(), httpResp),
		)
		return
	}

	// Map response to model
	r.mapResponseToModel(ctx, termination, &data, &resp.Diagnostics)
	tflog.Debug(ctx, "Updated circuit termination", map[string]interface{}{
		"id": data.ID.ValueString(),
	})

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete deletes the circuit termination resource.
func (r *CircuitTerminationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CircuitTerminationResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid ID",
			fmt.Sprintf("Could not convert ID to integer: %s", err.Error()),
		)
		return
	}
	tflog.Debug(ctx, "Deleting circuit termination", map[string]interface{}{
		"id": id,
	})

	// Call API to delete circuit termination
	httpResp, err := r.client.CircuitsAPI.CircuitsCircuitTerminationsDestroy(ctx, id).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			tflog.Debug(ctx, "Circuit termination already deleted", map[string]interface{}{
				"id": id,
			})
			return
		}
		resp.Diagnostics.AddError(
			"Error deleting circuit termination",
			fmt.Sprintf("Could not delete circuit termination: %s\nHTTP Response: %v", err.Error(), httpResp),
		)
		return
	}
	tflog.Debug(ctx, "Deleted circuit termination", map[string]interface{}{
		"id": id,
	})
}

// ImportState imports an existing circuit termination.
func (r *CircuitTerminationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// buildCreateRequest builds a CircuitTerminationRequest from the model.
func (r *CircuitTerminationResource) buildCreateRequest(ctx context.Context, data *CircuitTerminationResourceModel, state *CircuitTerminationResourceModel) (*netbox.CircuitTerminationRequest, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Look up Circuit (required)
	circuit, circuitDiags := netboxlookup.LookupCircuit(ctx, r.client, data.Circuit.ValueString())
	diags.Append(circuitDiags...)
	if diags.HasError() {
		return nil, diags
	}

	// Parse term_side
	termSide, err := netbox.NewTermination1FromValue(data.TermSide.ValueString())
	if err != nil {
		diags.AddError("Invalid term_side", err.Error())
		return nil, diags
	}
	createReq := netbox.NewCircuitTerminationRequest(*circuit, *termSide)

	// Handle site (optional)
	if !data.Site.IsNull() && !data.Site.IsUnknown() {
		site, siteDiags := netboxlookup.LookupSite(ctx, r.client, data.Site.ValueString())
		diags.Append(siteDiags...)
		if diags.HasError() {
			return nil, diags
		}
		createReq.SetSite(*site)
	}

	// Handle provider_network (optional) - reference by name
	if !data.ProviderNetwork.IsNull() && !data.ProviderNetwork.IsUnknown() {
		// Provider network is referenced by name
		pnReq := netbox.NewBriefProviderNetworkRequest(data.ProviderNetwork.ValueString())
		createReq.SetProviderNetwork(*pnReq)
	} else if data.ProviderNetwork.IsNull() {
		// Explicitly clear provider_network
		createReq.SetProviderNetworkNil()
	}

	// Handle port_speed (optional)
	if !data.PortSpeed.IsNull() && !data.PortSpeed.IsUnknown() {
		portSpeed, err := utils.SafeInt32FromValue(data.PortSpeed)
		if err != nil {
			diags.AddError("Invalid value", fmt.Sprintf("PortSpeed value overflow: %s", err))
			return nil, diags
		}
		createReq.SetPortSpeed(portSpeed)
	} else if data.PortSpeed.IsNull() {
		// Explicitly clear port_speed
		createReq.SetPortSpeedNil()
	}

	// Handle upstream_speed (optional)
	if !data.UpstreamSpeed.IsNull() && !data.UpstreamSpeed.IsUnknown() {
		upstreamSpeed, err := utils.SafeInt32FromValue(data.UpstreamSpeed)
		if err != nil {
			diags.AddError("Invalid value", fmt.Sprintf("UpstreamSpeed value overflow: %s", err))
			return nil, diags
		}
		createReq.SetUpstreamSpeed(upstreamSpeed)
	} else if data.UpstreamSpeed.IsNull() {
		// Explicitly clear upstream_speed
		createReq.SetUpstreamSpeedNil()
	}

	// Handle xconnect_id (optional)
	if !data.XconnectID.IsNull() && !data.XconnectID.IsUnknown() {
		createReq.SetXconnectId(data.XconnectID.ValueString())
	} else if data.XconnectID.IsNull() {
		// Explicitly clear xconnect_id
		emptyString := ""
		createReq.XconnectId = &emptyString
	}

	// Handle pp_info (optional)
	if !data.PPInfo.IsNull() && !data.PPInfo.IsUnknown() {
		createReq.SetPpInfo(data.PPInfo.ValueString())
	} else if data.PPInfo.IsNull() {
		// Explicitly clear pp_info
		emptyString := ""
		createReq.PpInfo = &emptyString
	}

	// Handle description (optional)
	utils.ApplyDescription(createReq, data.Description)

	// Handle mark_connected (optional)
	if !data.MarkConnected.IsNull() && !data.MarkConnected.IsUnknown() {
		createReq.SetMarkConnected(data.MarkConnected.ValueBool())
	}

	// Handle tags with conditional logic (use plan if set, otherwise state)
	if !data.Tags.IsNull() && !data.Tags.IsUnknown() {
		utils.ApplyTagsFromSlugs(ctx, r.client, createReq, data.Tags, &diags)
	} else if state != nil && !state.Tags.IsNull() && !state.Tags.IsUnknown() {
		utils.ApplyTagsFromSlugs(ctx, r.client, createReq, state.Tags, &diags)
	}

	// Handle custom fields with merge-aware logic
	if state != nil {
		utils.ApplyCustomFieldsWithMerge(ctx, createReq, data.CustomFields, state.CustomFields, &diags)
	} else {
		utils.ApplyCustomFields(ctx, createReq, data.CustomFields, &diags)
	}
	return createReq, diags
}

// mapResponseToModel maps the API response to the Terraform model.
func (r *CircuitTerminationResource) mapResponseToModel(ctx context.Context, termination *netbox.CircuitTermination, data *CircuitTerminationResourceModel, diags *diag.Diagnostics) {
	data.ID = types.StringValue(fmt.Sprintf("%d", termination.GetId()))
	data.TermSide = types.StringValue(string(termination.GetTermSide()))

	// Map Circuit - preserve user's input format
	if circuit := termination.GetCircuit(); circuit.Id != 0 {
		data.Circuit = utils.UpdateReferenceAttribute(data.Circuit, circuit.GetCid(), "", circuit.Id)
	}

	// Map Site - preserve user's input format
	if site, ok := termination.GetSiteOk(); ok && site != nil && site.Id != 0 {
		data.Site = utils.UpdateReferenceAttribute(data.Site, site.GetName(), site.GetSlug(), site.Id)
	} else {
		data.Site = types.StringNull()
	}

	// Map ProviderNetwork - preserve user's input format
	if pn, ok := termination.GetProviderNetworkOk(); ok && pn != nil && pn.Id != 0 {
		data.ProviderNetwork = utils.UpdateReferenceAttribute(data.ProviderNetwork, pn.GetName(), "", pn.Id)
	} else {
		data.ProviderNetwork = types.StringNull()
	}

	// Map port_speed
	if portSpeed, ok := termination.GetPortSpeedOk(); ok && portSpeed != nil {
		data.PortSpeed = types.Int64Value(int64(*portSpeed))
	} else {
		data.PortSpeed = types.Int64Null()
	}

	// Map upstream_speed
	if upstreamSpeed, ok := termination.GetUpstreamSpeedOk(); ok && upstreamSpeed != nil {
		data.UpstreamSpeed = types.Int64Value(int64(*upstreamSpeed))
	} else {
		data.UpstreamSpeed = types.Int64Null()
	}

	// Map xconnect_id
	if xconnectID, ok := termination.GetXconnectIdOk(); ok && xconnectID != nil && *xconnectID != "" {
		data.XconnectID = types.StringValue(*xconnectID)
	} else {
		data.XconnectID = types.StringNull()
	}

	// Map pp_info
	if ppInfo, ok := termination.GetPpInfoOk(); ok && ppInfo != nil && *ppInfo != "" {
		data.PPInfo = types.StringValue(*ppInfo)
	} else {
		data.PPInfo = types.StringNull()
	}

	// Map description
	if description, ok := termination.GetDescriptionOk(); ok && description != nil && *description != "" {
		data.Description = types.StringValue(*description)
	} else {
		data.Description = types.StringNull()
	}

	// Map mark_connected
	if markConnected, ok := termination.GetMarkConnectedOk(); ok && markConnected != nil {
		data.MarkConnected = types.BoolValue(*markConnected)
	} else {
		data.MarkConnected = types.BoolValue(false)
	}

	// Populate tags using slug list format
	wasExplicitlyEmpty := !data.Tags.IsNull() && !data.Tags.IsUnknown() && len(data.Tags.Elements()) == 0
	switch {
	case len(termination.Tags) > 0:
		tagSlugs := make([]string, 0, len(termination.Tags))
		for _, tag := range termination.Tags {
			tagSlugs = append(tagSlugs, tag.GetSlug())
		}
		data.Tags = utils.TagsSlugToSet(ctx, tagSlugs)
	case wasExplicitlyEmpty:
		data.Tags = types.SetValueMust(types.StringType, []attr.Value{})
	default:
		data.Tags = types.SetNull(types.StringType)
	}
	if termination.HasCustomFields() {
		data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, data.CustomFields, termination.CustomFields, diags)
	}
}
