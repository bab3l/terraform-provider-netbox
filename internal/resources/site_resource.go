// Package resources contains Terraform resource implementations for the Netbox provider.
//
// This package integrates with the go-netbox OpenAPI client to provide
// CRUD operations for Netbox resources via Terraform.
package resources

import (
	"context"
	"fmt"

	"github.com/bab3l/go-netbox"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &SiteResource{}
var _ resource.ResourceWithImportState = &SiteResource{}

func NewSiteResource() resource.Resource {
	return &SiteResource{}
}

// SiteResource defines the resource implementation.
type SiteResource struct {
	client *netbox.APIClient
}

// SiteResourceModel describes the resource data model.
type SiteResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Slug        types.String `tfsdk:"slug"`
	Status      types.String `tfsdk:"status"`
	Region      types.String `tfsdk:"region"`
	Group       types.String `tfsdk:"group"`
	Tenant      types.String `tfsdk:"tenant"`
	Facility    types.String `tfsdk:"facility"`
	Description types.String `tfsdk:"description"`
	Comments    types.String `tfsdk:"comments"`
}

func (r *SiteResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_site"
}

func (r *SiteResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a site in Netbox. Sites represent physical locations such as data centers, offices, or other facilities where network infrastructure is deployed.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Unique identifier for the site (assigned by Netbox).",
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Full name of the site. This is the human-readable display name.",
				Required:            true,
			},
			"slug": schema.StringAttribute{
				MarkdownDescription: "URL-friendly identifier for the site. Must be unique and contain only alphanumeric characters, hyphens, and underscores.",
				Required:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "Operational status of the site. Valid values include: `planned`, `staging`, `active`, `decommissioning`, `retired`.",
				Optional:            true,
				Computed:            true,
			},
			"region": schema.StringAttribute{
				MarkdownDescription: "Name or ID of the region where this site is located. Regions help organize sites geographically.",
				Optional:            true,
			},
			"group": schema.StringAttribute{
				MarkdownDescription: "Name or ID of the site group. Site groups provide an additional level of organization.",
				Optional:            true,
			},
			"tenant": schema.StringAttribute{
				MarkdownDescription: "Name or ID of the tenant that owns this site. Used for multi-tenancy scenarios.",
				Optional:            true,
			},
			"facility": schema.StringAttribute{
				MarkdownDescription: "Local facility identifier or description (e.g., building name, floor, room number).",
				Optional:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Detailed description of the site, its purpose, or other relevant information.",
				Optional:            true,
			},
			"comments": schema.StringAttribute{
				MarkdownDescription: "Additional comments or notes about the site. Supports Markdown formatting.",
				Optional:            true,
			},
		},
	}
}

func (r *SiteResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *SiteResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data SiteResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create site using go-netbox client
	tflog.Debug(ctx, "Creating site", map[string]interface{}{
		"name": data.Name.ValueString(),
		"slug": data.Slug.ValueString(),
	})

	// Prepare the site request
	siteRequest := netbox.WritableSiteRequest{
		Name: data.Name.ValueString(),
		Slug: data.Slug.ValueString(),
	}

	// Set optional fields if provided
	if !data.Status.IsNull() {
		// For now, let's skip status validation and use it as-is
		// In a production implementation, you'd want to validate against allowed values
		tflog.Debug(ctx, "Status field provided but skipped for now", map[string]interface{}{
			"status": data.Status.ValueString(),
		})
	}
	if !data.Description.IsNull() {
		description := data.Description.ValueString()
		siteRequest.Description = &description
	}
	if !data.Comments.IsNull() {
		comments := data.Comments.ValueString()
		siteRequest.Comments = &comments
	}
	if !data.Facility.IsNull() {
		facility := data.Facility.ValueString()
		siteRequest.Facility = &facility
	}

	// Create the site via API
	site, httpResp, err := r.client.DcimAPI.DcimSitesCreate(ctx).WritableSiteRequest(siteRequest).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating site",
			fmt.Sprintf("Could not create site, unexpected error: %s", err),
		)
		return
	}

	if httpResp.StatusCode != 201 {
		resp.Diagnostics.AddError(
			"Error creating site",
			fmt.Sprintf("Expected HTTP 201, got: %d", httpResp.StatusCode),
		)
		return
	}

	// Update the model with the response from the API
	data.ID = types.StringValue(fmt.Sprintf("%d", site.GetId()))
	data.Name = types.StringValue(site.GetName())
	data.Slug = types.StringValue(site.GetSlug())

	if site.HasStatus() {
		status := site.GetStatus()
		if status.HasValue() {
			statusValue, _ := status.GetValueOk()
			data.Status = types.StringValue(string(*statusValue))
		}
	} else {
		data.Status = types.StringValue("active") // default status
	}

	if site.HasDescription() {
		data.Description = types.StringValue(site.GetDescription())
	}

	if site.HasComments() {
		data.Comments = types.StringValue(site.GetComments())
	}

	if site.HasFacility() {
		data.Facility = types.StringValue(site.GetFacility())
	}

	tflog.Trace(ctx, "created a site resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SiteResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data SiteResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get the site ID from state
	siteID := data.ID.ValueString()

	tflog.Debug(ctx, "Reading site", map[string]interface{}{
		"id": siteID,
	})

	// Parse the site ID to int32 for the API call
	var siteIDInt int32
	if _, err := fmt.Sscanf(siteID, "%d", &siteIDInt); err != nil {
		resp.Diagnostics.AddError(
			"Invalid Site ID",
			fmt.Sprintf("Site ID must be a number, got: %s", siteID),
		)
		return
	}

	// Retrieve the site via API
	site, httpResp, err := r.client.DcimAPI.DcimSitesRetrieve(ctx, siteIDInt).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading site",
			fmt.Sprintf("Could not read site ID %s: %s", siteID, err),
		)
		return
	}

	if httpResp.StatusCode == 404 {
		// Site no longer exists, remove from state
		resp.State.RemoveResource(ctx)
		return
	}

	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Error reading site",
			fmt.Sprintf("Expected HTTP 200, got: %d", httpResp.StatusCode),
		)
		return
	}

	// Update the model with the response from the API
	data.ID = types.StringValue(fmt.Sprintf("%d", site.GetId()))
	data.Name = types.StringValue(site.GetName())
	data.Slug = types.StringValue(site.GetSlug())

	if site.HasStatus() {
		status := site.GetStatus()
		if status.HasValue() {
			statusValue, _ := status.GetValueOk()
			data.Status = types.StringValue(string(*statusValue))
		}
	}

	if site.HasDescription() {
		data.Description = types.StringValue(site.GetDescription())
	} else {
		data.Description = types.StringNull()
	}

	if site.HasComments() {
		data.Comments = types.StringValue(site.GetComments())
	} else {
		data.Comments = types.StringNull()
	}

	if site.HasFacility() {
		data.Facility = types.StringValue(site.GetFacility())
	} else {
		data.Facility = types.StringNull()
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SiteResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data SiteResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: Implement site update using go-netbox client
	tflog.Debug(ctx, "Updating site", map[string]interface{}{
		"id":   data.ID.ValueString(),
		"name": data.Name.ValueString(),
		"slug": data.Slug.ValueString(),
	})

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SiteResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data SiteResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: Implement site deletion using go-netbox client
	tflog.Debug(ctx, "Deleting site", map[string]interface{}{
		"id": data.ID.ValueString(),
	})
}

func (r *SiteResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
