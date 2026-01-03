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
	"github.com/bab3l/terraform-provider-netbox/internal/validators"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &ASNRangeResource{}
	_ resource.ResourceWithConfigure   = &ASNRangeResource{}
	_ resource.ResourceWithImportState = &ASNRangeResource{}
)

// NewASNRangeResource returns a new ASNRange resource.
func NewASNRangeResource() resource.Resource {
	return &ASNRangeResource{}
}

// ASNRangeResource defines the resource implementation.
type ASNRangeResource struct {
	client *netbox.APIClient
}

// ASNRangeResourceModel describes the resource data model.
type ASNRangeResourceModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Slug         types.String `tfsdk:"slug"`
	RIR          types.String `tfsdk:"rir"`
	Start        types.String `tfsdk:"start"`
	End          types.String `tfsdk:"end"`
	Tenant       types.String `tfsdk:"tenant"`
	Description  types.String `tfsdk:"description"`
	Tags         types.Set    `tfsdk:"tags"`
	CustomFields types.Set    `tfsdk:"custom_fields"`
}

// Metadata returns the resource type name.
func (r *ASNRangeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_asn_range"
}

// Schema defines the schema for the resource.
func (r *ASNRangeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an ASN Range in Netbox. ASN ranges define a contiguous range of Autonomous System Numbers that can be allocated.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique numeric ID of the ASN range.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": nbschema.NameAttribute("ASN range", 100),
			"slug": nbschema.SlugAttribute("ASN range"),
			"rir":  nbschema.RequiredReferenceAttribute("RIR", "ID or slug of the Regional Internet Registry (RIR) responsible for this ASN range. Required."),
			"start": schema.StringAttribute{
				MarkdownDescription: "The starting ASN in this range. Required.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						validators.IntegerRegex(),
						"must be a valid integer",
					),
				},
			},
			"end": schema.StringAttribute{
				MarkdownDescription: "The ending ASN in this range. Required.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						validators.IntegerRegex(),
						"must be a valid integer",
					),
				},
			},
			"tenant": nbschema.ReferenceAttribute("tenant", "ID or slug of the tenant that owns this ASN range."),
		},
	}

	// Add description attribute
	maps.Copy(resp.Schema.Attributes, nbschema.DescriptionOnlyAttributes("ASN range"))

	// Add common metadata attributes (tags, custom_fields)
	maps.Copy(resp.Schema.Attributes, nbschema.CommonMetadataAttributes())
}

// Configure adds the provider configured client to the resource.
func (r *ASNRangeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *ASNRangeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ASNRangeResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Lookup RIR
	rirRef, rirDiags := netboxlookup.LookupRIR(ctx, r.client, data.RIR.ValueString())
	resp.Diagnostics.Append(rirDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Parse start and end to int64
	var start, end int64
	if _, err := fmt.Sscanf(data.Start.ValueString(), "%d", &start); err != nil {
		resp.Diagnostics.AddError(
			"Invalid Start",
			fmt.Sprintf("Unable to parse start %q: %s", data.Start.ValueString(), err.Error()),
		)
		return
	}
	if _, err := fmt.Sscanf(data.End.ValueString(), "%d", &end); err != nil {
		resp.Diagnostics.AddError(
			"Invalid End",
			fmt.Sprintf("Unable to parse end %q: %s", data.End.ValueString(), err.Error()),
		)
		return
	}

	// Create the ASNRange request
	asnRangeRequest := netbox.NewASNRangeRequest(
		data.Name.ValueString(),
		data.Slug.ValueString(),
		*rirRef,
		start,
		end,
	)

	// Set optional fields
	r.setOptionalFields(ctx, asnRangeRequest, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating ASNRange", map[string]interface{}{
		"name":  data.Name.ValueString(),
		"start": start,
		"end":   end,
	})

	// Create the ASNRange
	asnRange, httpResp, err := r.client.IpamAPI.IpamAsnRangesCreate(ctx).ASNRangeRequest(*asnRangeRequest).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating ASNRange",
			utils.FormatAPIError("create ASNRange", err, httpResp),
		)
		return
	}

	// Map response to model
	r.mapASNRangeToState(ctx, asnRange, &data, &resp.Diagnostics)
	tflog.Debug(ctx, "Created ASNRange", map[string]interface{}{
		"id":   data.ID.ValueString(),
		"name": data.Name.ValueString(),
	})

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *ASNRangeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ASNRangeResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Parse the ID
	id, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid ID",
			fmt.Sprintf("Unable to parse ID %q: %s", data.ID.ValueString(), err.Error()),
		)
		return
	}
	tflog.Debug(ctx, "Reading ASNRange", map[string]interface{}{
		"id": id,
	})

	// Get the ASNRange from Netbox
	asnRange, httpResp, err := r.client.IpamAPI.IpamAsnRangesRetrieve(ctx, id).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading ASNRange",
			utils.FormatAPIError(fmt.Sprintf("read ASNRange ID %d", id), err, httpResp),
		)
		return
	}

	// Map response to model
	r.mapASNRangeToState(ctx, asnRange, &data, &resp.Diagnostics)
	tflog.Debug(ctx, "Read ASNRange", map[string]interface{}{
		"id":   data.ID.ValueString(),
		"name": data.Name.ValueString(),
	})

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *ASNRangeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ASNRangeResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Parse the ID
	id, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid ID",
			fmt.Sprintf("Unable to parse ID %q: %s", data.ID.ValueString(), err.Error()),
		)
		return
	}

	// Lookup RIR
	rirRef, rirDiags := netboxlookup.LookupRIR(ctx, r.client, data.RIR.ValueString())
	resp.Diagnostics.Append(rirDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Parse start and end to int64
	var start, end int64
	if _, err := fmt.Sscanf(data.Start.ValueString(), "%d", &start); err != nil {
		resp.Diagnostics.AddError(
			"Invalid Start",
			fmt.Sprintf("Unable to parse start %q: %s", data.Start.ValueString(), err.Error()),
		)
		return
	}
	if _, err := fmt.Sscanf(data.End.ValueString(), "%d", &end); err != nil {
		resp.Diagnostics.AddError(
			"Invalid End",
			fmt.Sprintf("Unable to parse end %q: %s", data.End.ValueString(), err.Error()),
		)
		return
	}

	// Create the ASNRange request
	asnRangeRequest := netbox.NewASNRangeRequest(
		data.Name.ValueString(),
		data.Slug.ValueString(),
		*rirRef,
		start,
		end,
	)

	// Set optional fields
	r.setOptionalFields(ctx, asnRangeRequest, &data, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Updating ASNRange", map[string]interface{}{
		"id":   id,
		"name": data.Name.ValueString(),
	})

	// Update the ASNRange
	asnRange, httpResp, err := r.client.IpamAPI.IpamAsnRangesUpdate(ctx, id).ASNRangeRequest(*asnRangeRequest).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating ASNRange",
			utils.FormatAPIError(fmt.Sprintf("update ASNRange ID %d", id), err, httpResp),
		)
		return
	}

	// Map response to model
	r.mapASNRangeToState(ctx, asnRange, &data, &resp.Diagnostics)
	tflog.Debug(ctx, "Updated ASNRange", map[string]interface{}{
		"id":   data.ID.ValueString(),
		"name": data.Name.ValueString(),
	})

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *ASNRangeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ASNRangeResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Parse the ID
	id, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid ID",
			fmt.Sprintf("Unable to parse ID %q: %s", data.ID.ValueString(), err.Error()),
		)
		return
	}
	tflog.Debug(ctx, "Deleting ASNRange", map[string]interface{}{
		"id": id,
	})

	// Delete the ASNRange
	httpResp, err := r.client.IpamAPI.IpamAsnRangesDestroy(ctx, id).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			// Resource already deleted
			return
		}
		resp.Diagnostics.AddError(
			"Error deleting ASNRange",
			utils.FormatAPIError(fmt.Sprintf("delete ASNRange ID %d", id), err, httpResp),
		)
		return
	}
	tflog.Debug(ctx, "Deleted ASNRange", map[string]interface{}{
		"id": id,
	})
}

func (r *ASNRangeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// setOptionalFields sets optional fields on the ASNRange request from the resource model.
func (r *ASNRangeResource) setOptionalFields(ctx context.Context, asnRangeRequest *netbox.ASNRangeRequest, data *ASNRangeResourceModel, diags *diag.Diagnostics) {
	// Tenant
	if utils.IsSet(data.Tenant) {
		tenantRef, tenantDiags := netboxlookup.LookupTenant(ctx, r.client, data.Tenant.ValueString())
		diags.Append(tenantDiags...)
		if diags.HasError() {
			return
		}
		asnRangeRequest.Tenant = *netbox.NewNullableBriefTenantRequest(tenantRef)
	}

	// Description
	utils.ApplyDescription(asnRangeRequest, data.Description)

	// Handle tags and custom_fields
	utils.ApplyMetadataFields(ctx, asnRangeRequest, data.Tags, data.CustomFields, diags)
}

// mapASNRangeToState maps a Netbox ASNRange to the Terraform state model.
func (r *ASNRangeResource) mapASNRangeToState(ctx context.Context, asnRange *netbox.ASNRange, data *ASNRangeResourceModel, diags *diag.Diagnostics) {
	data.ID = types.StringValue(fmt.Sprintf("%d", asnRange.Id))
	data.Name = types.StringValue(asnRange.Name)
	data.Slug = types.StringValue(asnRange.Slug)

	// RIR - preserve user's input format
	rir := asnRange.Rir
	data.RIR = utils.UpdateReferenceAttribute(data.RIR, rir.GetName(), rir.GetSlug(), rir.GetId())
	data.Start = types.StringValue(fmt.Sprintf("%d", asnRange.Start))
	data.End = types.StringValue(fmt.Sprintf("%d", asnRange.End))

	// Tenant - preserve user's input format
	if asnRange.HasTenant() && asnRange.Tenant.Get() != nil {
		tenant := asnRange.Tenant.Get()
		data.Tenant = utils.UpdateReferenceAttribute(data.Tenant, tenant.GetName(), tenant.GetSlug(), tenant.GetId())
	} else {
		data.Tenant = types.StringNull()
	}

	// Description
	if asnRange.Description != nil && *asnRange.Description != "" {
		data.Description = types.StringValue(*asnRange.Description)
	} else {
		data.Description = types.StringNull()
	}

	// Tags
	data.Tags = utils.PopulateTagsFromAPI(ctx, len(asnRange.Tags) > 0, asnRange.Tags, data.Tags, diags)
	if diags.HasError() {
		return
	}

	// Custom Fields
	data.CustomFields = utils.PopulateCustomFieldsFromAPI(ctx, len(asnRange.CustomFields) > 0, asnRange.CustomFields, data.CustomFields, diags)
}
