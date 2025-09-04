package resources

import (
	"context"
	"fmt"

	"github.com/bab3l/go-netbox"
	"github.com/bab3l/terraform-provider-netbox/internal/netboxlookup"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &PlatformResource{}
var _ resource.ResourceWithImportState = &PlatformResource{}

func NewPlatformResource() resource.Resource {
	return &PlatformResource{}
}

type PlatformResource struct {
	client *netbox.APIClient
}

type PlatformResourceModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Slug         types.String `tfsdk:"slug"`
	Manufacturer types.String `tfsdk:"manufacturer"`
	Description  types.String `tfsdk:"description"`
}

func (r *PlatformResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "netbox_platform"
}

func (r *PlatformResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Unique identifier for the platform (assigned by Netbox).",
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Full name of the platform.",
			},
			"slug": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "URL-friendly identifier for the platform.",
			},
			"manufacturer": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Reference to the manufacturer (ID or slug).",
			},
		},
	}
}

// Implement Create, Read, Update, Delete, and ImportState methods here
func (r *PlatformResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*netbox.APIClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *netbox.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
	}
	r.client = client
}

func (r *PlatformResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data PlatformResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Creating platform", map[string]interface{}{
		"name": data.Name.ValueString(),
		"slug": data.Slug.ValueString(),
	})
	platformRequest := netbox.NewPlatformRequest(data.Name.ValueString(), data.Slug.ValueString())
	if !data.Manufacturer.IsNull() {
		manufacturerRef, diags := netboxlookup.LookupManufacturerBrief(ctx, r.client, data.Manufacturer.ValueString())
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		platformRequest.Manufacturer = *netbox.NewNullableBriefManufacturerRequest(manufacturerRef)
	}
	if !data.Description.IsNull() {
		desc := data.Description.ValueString()
		platformRequest.Description = &desc
	}
	// Tags and custom fields can be added here if needed

	platform, httpResp, err := r.client.DcimAPI.DcimPlatformsCreate(ctx).PlatformRequest(*platformRequest).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Error creating platform", fmt.Sprintf("Could not create platform: %s", err))
		return
	}
	if httpResp.StatusCode != 201 {
		resp.Diagnostics.AddError("Error creating platform", fmt.Sprintf("Expected HTTP 201, got: %d", httpResp.StatusCode))
		return
	}
	data.ID = types.StringValue(fmt.Sprintf("%d", platform.GetId()))
	data.Name = types.StringValue(platform.GetName())
	data.Slug = types.StringValue(platform.GetSlug())
	if platform.HasManufacturer() {
		data.Manufacturer = types.StringValue(platform.GetManufacturer().Name)
	} else {
		data.Manufacturer = types.StringNull()
	}
	if platform.HasDescription() {
		data.Description = types.StringValue(platform.GetDescription())
	} else {
		data.Description = types.StringNull()
	}
	tflog.Trace(ctx, "created a platform resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PlatformResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data PlatformResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	platformID := data.ID.ValueString()
	tflog.Debug(ctx, "Reading platform", map[string]interface{}{"id": platformID})
	var platformIDInt int32
	if _, err := fmt.Sscanf(platformID, "%d", &platformIDInt); err != nil {
		resp.Diagnostics.AddError("Invalid Platform ID", fmt.Sprintf("Platform ID must be a number, got: %s", platformID))
		return
	}
	platform, httpResp, err := r.client.DcimAPI.DcimPlatformsRetrieve(ctx, platformIDInt).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Error reading platform", fmt.Sprintf("Could not read platform: %s", err))
		return
	}
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError("Error reading platform", fmt.Sprintf("Expected HTTP 200, got: %d", httpResp.StatusCode))
		return
	}
	data.ID = types.StringValue(fmt.Sprintf("%d", platform.GetId()))
	data.Name = types.StringValue(platform.GetName())
	data.Slug = types.StringValue(platform.GetSlug())
	if platform.HasManufacturer() {
		// Use slug if available, else ID
		m := platform.GetManufacturer()
		if m.Slug != "" {
			data.Manufacturer = types.StringValue(m.Slug)
		} else {
			data.Manufacturer = types.StringValue(fmt.Sprintf("%d", m.Id))
		}
	} else {
		data.Manufacturer = types.StringNull()
	}
	if platform.HasDescription() {
		data.Description = types.StringValue(platform.GetDescription())
	} else {
		data.Description = types.StringNull()
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PlatformResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data PlatformResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	platformID := data.ID.ValueString()
	var platformIDInt int32
	if _, err := fmt.Sscanf(platformID, "%d", &platformIDInt); err != nil {
		resp.Diagnostics.AddError("Invalid Platform ID", fmt.Sprintf("Platform ID must be a number, got: %s", platformID))
		return
	}
	platformRequest := netbox.PlatformRequest{
		Name: data.Name.ValueString(),
		Slug: data.Slug.ValueString(),
	}
	if !data.Manufacturer.IsNull() {
		manufacturerRef, diags := netboxlookup.LookupManufacturerBrief(ctx, r.client, data.Manufacturer.ValueString())
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		platformRequest.Manufacturer = *netbox.NewNullableBriefManufacturerRequest(manufacturerRef)
	}
	if !data.Description.IsNull() {
		desc := data.Description.ValueString()
		platformRequest.Description = &desc
	}
	platform, httpResp, err := r.client.DcimAPI.DcimPlatformsUpdate(ctx, platformIDInt).PlatformRequest(platformRequest).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Error updating platform", fmt.Sprintf("Could not update platform: %s", err))
		return
	}
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError("Error updating platform", fmt.Sprintf("Expected HTTP 200, got: %d", httpResp.StatusCode))
		return
	}
	data.ID = types.StringValue(fmt.Sprintf("%d", platform.GetId()))
	data.Name = types.StringValue(platform.GetName())
	data.Slug = types.StringValue(platform.GetSlug())
	if platform.HasManufacturer() {
		m := platform.GetManufacturer()
		if m.Slug != "" {
			data.Manufacturer = types.StringValue(m.Slug)
		} else {
			data.Manufacturer = types.StringValue(fmt.Sprintf("%d", m.Id))
		}
	} else {
		data.Manufacturer = types.StringNull()
	}
	if platform.HasDescription() {
		data.Description = types.StringValue(platform.GetDescription())
	} else {
		data.Description = types.StringNull()
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PlatformResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data PlatformResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	platformID := data.ID.ValueString()
	var platformIDInt int32
	if _, err := fmt.Sscanf(platformID, "%d", &platformIDInt); err != nil {
		resp.Diagnostics.AddError("Invalid Platform ID", fmt.Sprintf("Platform ID must be a number, got: %s", platformID))
		return
	}
	httpResp, err := r.client.DcimAPI.DcimPlatformsDestroy(ctx, platformIDInt).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Error deleting platform", fmt.Sprintf("Could not delete platform: %s", err))
		return
	}
	if httpResp.StatusCode != 204 {
		resp.Diagnostics.AddError("Error deleting platform", fmt.Sprintf("Expected HTTP 204, got: %d", httpResp.StatusCode))
		return
	}
}

func (r *PlatformResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
