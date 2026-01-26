// Package resources contains Terraform resource implementations for the Netbox provider.

package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/bab3l/go-netbox"
	"github.com/bab3l/terraform-provider-netbox/internal/netboxlookup"
	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &SiteASNAssignmentResource{}
	_ resource.ResourceWithConfigure   = &SiteASNAssignmentResource{}
	_ resource.ResourceWithImportState = &SiteASNAssignmentResource{}
)

// NewSiteASNAssignmentResource returns a new site ASN assignment resource.
func NewSiteASNAssignmentResource() resource.Resource {
	return &SiteASNAssignmentResource{}
}

// SiteASNAssignmentResource defines the resource implementation.
type SiteASNAssignmentResource struct {
	client *netbox.APIClient
}

// SiteASNAssignmentResourceModel describes the resource data model.
type SiteASNAssignmentResourceModel struct {
	ID   types.String `tfsdk:"id"`
	Site types.String `tfsdk:"site"`
	ASN  types.String `tfsdk:"asn"`
}

func (r *SiteASNAssignmentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_site_asn_assignment"
}

func (r *SiteASNAssignmentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a site ASN association in NetBox. This resource associates a site with an ASN using the site ASN list.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Resource ID in the format <site_id>:<asn_id>.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"site": schema.StringAttribute{
				MarkdownDescription: "ID or slug of the site to associate with the ASN.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					nbschema.ReferenceEquivalencePlanModifier(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"asn": schema.StringAttribute{
				MarkdownDescription: "The unique numeric ID of the ASN to associate with the site.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *SiteASNAssignmentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *SiteASNAssignmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data SiteASNAssignmentResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	siteID, diags := netboxlookup.GenericLookupID(ctx, data.Site.ValueString(), netboxlookup.SiteLookupConfig(r.client), func(s *netbox.Site) int32 {
		return s.GetId()
	})
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	asnID, err := utils.ParseID(data.ASN.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ASN ID", fmt.Sprintf("ASN ID must be a number, got: %s", data.ASN.ValueString()))
		return
	}

	r.updateSiteASNs(ctx, siteID, asnID, true, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	site, httpResp, err := r.client.DcimAPI.DcimSitesRetrieve(ctx, siteID).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError("Error reading site", utils.FormatAPIError(fmt.Sprintf("read site ID %d", siteID), err, httpResp))
		return
	}
	data.Site = utils.UpdateReferenceAttribute(data.Site, site.GetName(), site.GetSlug(), site.GetId())

	data.ID = types.StringValue(fmt.Sprintf("%d:%d", siteID, asnID))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SiteASNAssignmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data SiteASNAssignmentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	siteID, asnID, ok := r.parseIDs(ctx, data, &resp.Diagnostics)
	if !ok {
		return
	}

	site, httpResp, err := r.client.DcimAPI.DcimSitesRetrieve(ctx, siteID).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading site", utils.FormatAPIError(fmt.Sprintf("read site ID %d", siteID), err, httpResp))
		return
	}

	if !siteHasASN(site, asnID) {
		resp.State.RemoveResource(ctx)
		return
	}
	data.Site = utils.UpdateReferenceAttribute(data.Site, site.GetName(), site.GetSlug(), site.GetId())

	data.ID = types.StringValue(fmt.Sprintf("%d:%d", siteID, asnID))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SiteASNAssignmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, plan SiteASNAssignmentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	oldSiteID, oldAsnID, ok := r.parseIDs(ctx, state, &resp.Diagnostics)
	if !ok {
		return
	}
	newSiteID, newAsnID, ok := r.parseIDs(ctx, plan, &resp.Diagnostics)
	if !ok {
		return
	}

	if oldSiteID != newSiteID || oldAsnID != newAsnID {
		r.updateSiteASNs(ctx, oldSiteID, oldAsnID, false, &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}
		r.updateSiteASNs(ctx, newSiteID, newAsnID, true, &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	site, httpResp, err := r.client.DcimAPI.DcimSitesRetrieve(ctx, newSiteID).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError("Error reading site", utils.FormatAPIError(fmt.Sprintf("read site ID %d", newSiteID), err, httpResp))
		return
	}
	plan.Site = utils.UpdateReferenceAttribute(plan.Site, site.GetName(), site.GetSlug(), site.GetId())

	plan.ID = types.StringValue(fmt.Sprintf("%d:%d", newSiteID, newAsnID))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *SiteASNAssignmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data SiteASNAssignmentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	siteID, asnID, ok := r.parseIDs(ctx, data, &resp.Diagnostics)
	if !ok {
		return
	}

	r.updateSiteASNs(ctx, siteID, asnID, false, &resp.Diagnostics)
}

func (r *SiteASNAssignmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, ":")
	if len(parts) != 2 {
		resp.Diagnostics.AddError("Invalid import ID", "Expected import ID in the format <site_id>:<asn_id>")
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("site"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("asn"), parts[1])...)
}

func (r *SiteASNAssignmentResource) parseIDs(ctx context.Context, data SiteASNAssignmentResourceModel, diags *diag.Diagnostics) (int32, int32, bool) {
	var siteID int32
	var asnID int32

	if !data.Site.IsNull() && !data.Site.IsUnknown() {
		resolvedID, lookupDiags := netboxlookup.GenericLookupID(ctx, data.Site.ValueString(), netboxlookup.SiteLookupConfig(r.client), func(s *netbox.Site) int32 {
			return s.GetId()
		})
		diags.Append(lookupDiags...)
		if diags.HasError() {
			return 0, 0, false
		}
		siteID = resolvedID
	} else if !data.ID.IsNull() && !data.ID.IsUnknown() {
		parts := strings.Split(data.ID.ValueString(), ":")
		if len(parts) != 2 {
			diags.AddError("Invalid ID", "Expected ID in the format <site_id>:<asn_id>")
			return 0, 0, false
		}
		parsed, err := utils.ParseID(parts[0])
		if err != nil {
			diags.AddError("Invalid Site ID", fmt.Sprintf("Site ID must be a number, got: %s", parts[0]))
			return 0, 0, false
		}
		siteID = parsed
	}

	if !data.ASN.IsNull() && !data.ASN.IsUnknown() {
		parsed, err := utils.ParseID(data.ASN.ValueString())
		if err != nil {
			diags.AddError("Invalid ASN ID", fmt.Sprintf("ASN ID must be a number, got: %s", data.ASN.ValueString()))
			return 0, 0, false
		}
		asnID = parsed
	} else if !data.ID.IsNull() && !data.ID.IsUnknown() {
		parts := strings.Split(data.ID.ValueString(), ":")
		if len(parts) != 2 {
			diags.AddError("Invalid ID", "Expected ID in the format <site_id>:<asn_id>")
			return 0, 0, false
		}
		parsed, err := utils.ParseID(parts[1])
		if err != nil {
			diags.AddError("Invalid ASN ID", fmt.Sprintf("ASN ID must be a number, got: %s", parts[1]))
			return 0, 0, false
		}
		asnID = parsed
	}

	return siteID, asnID, true
}

func (r *SiteASNAssignmentResource) updateSiteASNs(ctx context.Context, siteID int32, asnID int32, add bool, diags *diag.Diagnostics) {
	site, httpResp, err := r.client.DcimAPI.DcimSitesRetrieve(ctx, siteID).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			return
		}
		diags.AddError("Error reading site", utils.FormatAPIError(fmt.Sprintf("read site ID %d", siteID), err, httpResp))
		return
	}

	asnIDs := make([]int32, 0, len(site.GetAsns()))
	found := false
	for _, asn := range site.GetAsns() {
		id := asn.GetId()
		if id == asnID {
			found = true
			if !add {
				continue
			}
		}
		asnIDs = append(asnIDs, id)
	}

	if add && !found {
		asnIDs = append(asnIDs, asnID)
	}

	siteRequest := netbox.WritableSiteRequest{
		Name: site.GetName(),
		Slug: site.GetSlug(),
	}
	siteRequest.SetAsns(asnIDs)

	updated, updateResp, updateErr := r.client.DcimAPI.DcimSitesUpdate(ctx, siteID).WritableSiteRequest(siteRequest).Execute()
	defer utils.CloseResponseBody(updateResp)
	if updateErr != nil {
		diags.AddError("Error updating site ASN list", utils.FormatAPIError(fmt.Sprintf("update site ID %d", siteID), updateErr, updateResp))
		return
	}

	tflog.Debug(ctx, "Updated site ASN assignments", map[string]interface{}{
		"site_id": siteID,
		"asn_id":  asnID,
		"count":   len(updated.GetAsns()),
	})
}

func siteHasASN(site *netbox.Site, asnID int32) bool {
	if site == nil {
		return false
	}
	for _, asn := range site.GetAsns() {
		if asn.GetId() == asnID {
			return true
		}
	}
	return false
}
