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
	_ resource.Resource = &PowerPanelResource{}

	_ resource.ResourceWithConfigure = &PowerPanelResource{}

	_ resource.ResourceWithImportState = &PowerPanelResource{}
)

// NewPowerPanelResource returns a new resource implementing the power panel resource.

func NewPowerPanelResource() resource.Resource {
	return &PowerPanelResource{}
}

// PowerPanelResource defines the resource implementation.

type PowerPanelResource struct {
	client *netbox.APIClient
}

// PowerPanelResourceModel describes the resource data model.

type PowerPanelResourceModel struct {
	ID types.String `tfsdk:"id"`

	Site types.String `tfsdk:"site"`

	Location types.String `tfsdk:"location"`

	Name types.String `tfsdk:"name"`

	Description types.String `tfsdk:"description"`

	Comments types.String `tfsdk:"comments"`

	Tags types.Set `tfsdk:"tags"`

	CustomFields types.Set `tfsdk:"custom_fields"`
}

// Metadata returns the resource type name.

func (r *PowerPanelResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_power_panel"
}

// Schema defines the schema for the resource.

func (r *PowerPanelResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a power panel in NetBox. Power panels represent power distribution panels in data centers.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique numeric ID of the power panel.",

				Computed: true,

				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"site": schema.StringAttribute{
				MarkdownDescription: "The site this power panel belongs to (ID or slug).",

				Required: true,
			},

			"location": schema.StringAttribute{
				MarkdownDescription: "The location within the site (ID or slug).",

				Optional: true,
			},

			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the power panel.",

				Required: true,
			},

			"description": schema.StringAttribute{
				MarkdownDescription: "A description of the power panel.",

				Optional: true,
			},
		},
	}

	// Add common descriptive attributes (description, comments)
	maps.Copy(resp.Schema.Attributes, nbschema.CommonDescriptiveAttributes("power panel"))

	// Add common metadata attributes (tags, custom_fields)
	maps.Copy(resp.Schema.Attributes, nbschema.CommonMetadataAttributes())
}

// Configure adds the provider configured client to the resource.

func (r *PowerPanelResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *PowerPanelResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data PowerPanelResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Lookup site

	site, diags := lookup.LookupSite(ctx, r.client, data.Site.ValueString())

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Build request

	apiReq := netbox.NewPowerPanelRequest(*site, data.Name.ValueString())

	// Set optional fields

	if !data.Location.IsNull() && !data.Location.IsUnknown() {
		location, diags := lookup.LookupLocation(ctx, r.client, data.Location.ValueString())

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {
			return
		}

		apiReq.SetLocation(*location)
	}

	// Apply common fields (description, comments, tags, custom_fields)
	utils.ApplyCommonFields(ctx, apiReq, data.Description, data.Comments, data.Tags, data.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating power panel", map[string]interface{}{
		"name": data.Name.ValueString(),
	})

	response, httpResp, err := r.client.DcimAPI.DcimPowerPanelsCreate(ctx).PowerPanelRequest(*apiReq).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {
		resp.Diagnostics.AddError(

			"Error creating power panel",

			utils.FormatAPIError(fmt.Sprintf("create power panel %s", data.Name.ValueString()), err, httpResp),
		)

		return
	}

	// Map response to model

	r.mapResponseToModel(ctx, response, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "Created power panel", map[string]interface{}{
		"id": data.ID.ValueString(),

		"name": data.Name.ValueString(),
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read refreshes the resource state.

func (r *PowerPanelResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data PowerPanelResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ppID, err := utils.ParseID(data.ID.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(

			"Invalid Power Panel ID",

			fmt.Sprintf("Power panel ID must be a number, got: %s", data.ID.ValueString()),
		)

		return
	}

	tflog.Debug(ctx, "Reading power panel", map[string]interface{}{
		"id": ppID,
	})

	response, httpResp, err := r.client.DcimAPI.DcimPowerPanelsRetrieve(ctx, ppID).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			resp.State.RemoveResource(ctx)

			return
		}

		resp.Diagnostics.AddError(

			"Error reading power panel",

			utils.FormatAPIError(fmt.Sprintf("read power panel ID %d", ppID), err, httpResp),
		)

		return
	}

	// Map response to model

	r.mapResponseToModel(ctx, response, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates the resource.

func (r *PowerPanelResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, plan PowerPanelResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ppID, err := utils.ParseID(plan.ID.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(

			"Invalid Power Panel ID",

			fmt.Sprintf("Power panel ID must be a number, got: %s", plan.ID.ValueString()),
		)

		return
	}

	// Lookup site

	site, diags := lookup.LookupSite(ctx, r.client, plan.Site.ValueString())

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Build request

	apiReq := netbox.NewPowerPanelRequest(*site, plan.Name.ValueString())

	// Set optional fields

	if !plan.Location.IsNull() && !plan.Location.IsUnknown() {
		location, diags := lookup.LookupLocation(ctx, r.client, plan.Location.ValueString())

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {
			return
		}

		apiReq.SetLocation(*location)
	}

	// Apply common fields with merge-aware helpers
	utils.ApplyDescription(apiReq, plan.Description)
	utils.ApplyComments(apiReq, plan.Comments)
	utils.ApplyTags(ctx, apiReq, plan.Tags, &resp.Diagnostics)
	utils.ApplyCustomFieldsWithMerge(ctx, apiReq, plan.CustomFields, state.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Updating power panel", map[string]interface{}{
		"id": ppID,
	})

	response, httpResp, err := r.client.DcimAPI.DcimPowerPanelsUpdate(ctx, ppID).PowerPanelRequest(*apiReq).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {
		resp.Diagnostics.AddError(

			"Error updating power panel",

			utils.FormatAPIError(fmt.Sprintf("update power panel ID %d", ppID), err, httpResp),
		)

		return
	}

	// Map response to model

	r.mapResponseToModel(ctx, response, &plan, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete deletes the resource.

func (r *PowerPanelResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data PowerPanelResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ppID, err := utils.ParseID(data.ID.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(

			"Invalid Power Panel ID",

			fmt.Sprintf("Power panel ID must be a number, got: %s", data.ID.ValueString()),
		)

		return
	}

	tflog.Debug(ctx, "Deleting power panel", map[string]interface{}{
		"id": ppID,
	})

	httpResp, err := r.client.DcimAPI.DcimPowerPanelsDestroy(ctx, ppID).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			return
		}

		resp.Diagnostics.AddError(

			"Error deleting power panel",

			utils.FormatAPIError(fmt.Sprintf("delete power panel ID %d", ppID), err, httpResp),
		)

		return
	}
}

// ImportState imports an existing resource.

func (r *PowerPanelResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Try to parse as ID first

	ppID, err := utils.ParseID(req.ID)

	if err != nil {
		resp.Diagnostics.AddError(

			"Invalid Import ID",

			fmt.Sprintf("Power panel ID must be a number, got: %s", req.ID),
		)

		return
	}

	response, httpResp, err := r.client.DcimAPI.DcimPowerPanelsRetrieve(ctx, ppID).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {
		resp.Diagnostics.AddError(

			"Error importing power panel",

			utils.FormatAPIError(fmt.Sprintf("import power panel ID %d", ppID), err, httpResp),
		)

		return
	}

	var data PowerPanelResourceModel

	r.mapResponseToModel(ctx, response, &data, &resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// mapResponseToModel maps the API response to the Terraform model.

func (r *PowerPanelResource) mapResponseToModel(ctx context.Context, pp *netbox.PowerPanel, data *PowerPanelResourceModel, diags *diag.Diagnostics) {
	data.ID = types.StringValue(fmt.Sprintf("%d", pp.GetId()))

	data.Name = types.StringValue(pp.GetName())

	// Map site - preserve user's input format

	site := pp.Site

	data.Site = utils.UpdateReferenceAttribute(data.Site, site.GetName(), site.GetSlug(), site.GetId())

	// Map location - preserve user's input format

	if pp.Location.IsSet() && pp.Location.Get() != nil {
		loc := pp.Location.Get()

		data.Location = utils.UpdateReferenceAttribute(data.Location, loc.GetName(), loc.GetSlug(), loc.GetId())
	} else {
		data.Location = types.StringNull()
	}

	// Map description

	if desc, ok := pp.GetDescriptionOk(); ok && desc != nil && *desc != "" {
		data.Description = types.StringValue(*desc)
	} else {
		data.Description = types.StringNull()
	}

	// Map comments

	if comments, ok := pp.GetCommentsOk(); ok && comments != nil && *comments != "" {
		data.Comments = types.StringValue(*comments)
	} else {
		data.Comments = types.StringNull()
	}

	// Handle tags
	data.Tags = utils.PopulateTagsFromAPI(ctx, pp.HasTags(), pp.GetTags(), data.Tags, diags)
	if diags.HasError() {
		return
	}

	// Handle custom fields (filter to owned)
	data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, data.CustomFields, pp.GetCustomFields(), diags)
	if diags.HasError() {
		return
	}
}
