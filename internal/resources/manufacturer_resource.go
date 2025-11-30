package resources

import (
	"context"
	"fmt"

	"github.com/bab3l/go-netbox"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &ManufacturerResource{}
var _ resource.ResourceWithImportState = &ManufacturerResource{}

func NewManufacturerResource() resource.Resource {
	return &ManufacturerResource{}
}

type ManufacturerResource struct {
	client *netbox.APIClient
}

type ManufacturerResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Slug        types.String `tfsdk:"slug"`
	Description types.String `tfsdk:"description"`
}

func (r *ManufacturerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_manufacturer"
}

func (r *ManufacturerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a manufacturer in Netbox. Manufacturers are used to group devices and platforms by vendor.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Unique identifier for the manufacturer (assigned by Netbox).",
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Full name of the manufacturer.",
			},
			"slug": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "URL-friendly identifier for the manufacturer.",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Detailed description of the manufacturer.",
			},
		},
	}
}

// Implement Configure, Create, Read, Update, Delete, ImportState methods here
func (r *ManufacturerResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ManufacturerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ManufacturerResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	manufacturerRequest := netbox.ManufacturerRequest{
		Name: data.Name.ValueString(),
		Slug: data.Slug.ValueString(),
	}
	if !data.Description.IsNull() {
		desc := data.Description.ValueString()
		manufacturerRequest.Description = &desc
	}
	manufacturer, httpResp, err := r.client.DcimAPI.DcimManufacturersCreate(ctx).ManufacturerRequest(manufacturerRequest).Execute()
	tflog.Debug(ctx, "Manufacturer create API response", map[string]interface{}{
		"http_status":  httpResp.StatusCode,
		"manufacturer": manufacturer,
		"error":        err,
	})
	tflog.Debug(ctx, "Raw manufacturer API response", map[string]interface{}{
		"manufacturer_raw": fmt.Sprintf("%#v", manufacturer),
		"httpResp_raw":     fmt.Sprintf("%#v", httpResp),
		"error":            err,
	})
	if err != nil {
		resp.Diagnostics.AddError("Error creating manufacturer", utils.FormatAPIError("create manufacturer", err, httpResp))
		return
	}
	if httpResp.StatusCode != 201 {
		resp.Diagnostics.AddError("Error creating manufacturer", fmt.Sprintf("Expected HTTP 201, got: %d", httpResp.StatusCode))
		return
	}
	if manufacturer == nil {
		resp.Diagnostics.AddError("Manufacturer API returned nil", "No manufacturer object returned from Netbox API.")
		return
	}
	data.ID = types.StringValue(fmt.Sprintf("%d", manufacturer.GetId()))
	data.Name = types.StringValue(manufacturer.GetName())
	data.Slug = types.StringValue(manufacturer.GetSlug())
	if manufacturer.HasDescription() {
		data.Description = types.StringValue(manufacturer.GetDescription())
	} else {
		data.Description = types.StringNull()
	}
	tflog.Debug(ctx, "Setting resource state after manufacturer create", map[string]interface{}{
		"id":          data.ID.ValueString(),
		"name":        data.Name.ValueString(),
		"slug":        data.Slug.ValueString(),
		"description": data.Description.ValueString(),
	})
	diags := resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		tflog.Error(ctx, "Error setting resource state after manufacturer create", map[string]interface{}{
			"diagnostics": diags,
			"id":          data.ID.ValueString(),
			"name":        data.Name.ValueString(),
			"slug":        data.Slug.ValueString(),
			"description": data.Description.ValueString(),
		})
		return
	}
}

func (r *ManufacturerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ManufacturerResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	manufacturerID := data.ID.ValueString()
	var manufacturerIDInt int32
	if _, err := fmt.Sscanf(manufacturerID, "%d", &manufacturerIDInt); err != nil {
		resp.Diagnostics.AddError("Invalid Manufacturer ID", fmt.Sprintf("Manufacturer ID must be a number, got: %s", manufacturerID))
		return
	}
	manufacturer, httpResp, err := r.client.DcimAPI.DcimManufacturersRetrieve(ctx, manufacturerIDInt).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Error reading manufacturer", utils.FormatAPIError(fmt.Sprintf("read manufacturer ID %s", manufacturerID), err, httpResp))
		return
	}
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError("Error reading manufacturer", fmt.Sprintf("Expected HTTP 200, got: %d", httpResp.StatusCode))
		return
	}
	data.ID = types.StringValue(fmt.Sprintf("%d", manufacturer.GetId()))
	data.Name = types.StringValue(manufacturer.GetName())
	data.Slug = types.StringValue(manufacturer.GetSlug())
	if manufacturer.HasDescription() {
		data.Description = types.StringValue(manufacturer.GetDescription())
	} else {
		data.Description = types.StringNull()
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ManufacturerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ManufacturerResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	manufacturerID := data.ID.ValueString()
	var manufacturerIDInt int32
	if _, err := fmt.Sscanf(manufacturerID, "%d", &manufacturerIDInt); err != nil {
		resp.Diagnostics.AddError("Invalid Manufacturer ID", fmt.Sprintf("Manufacturer ID must be a number, got: %s", manufacturerID))
		return
	}
	manufacturerRequest := netbox.ManufacturerRequest{
		Name: data.Name.ValueString(),
		Slug: data.Slug.ValueString(),
	}
	if !data.Description.IsNull() {
		desc := data.Description.ValueString()
		manufacturerRequest.Description = &desc
	}
	manufacturer, httpResp, err := r.client.DcimAPI.DcimManufacturersUpdate(ctx, manufacturerIDInt).ManufacturerRequest(manufacturerRequest).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Error updating manufacturer", utils.FormatAPIError(fmt.Sprintf("update manufacturer ID %s", manufacturerID), err, httpResp))
		return
	}
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError("Error updating manufacturer", fmt.Sprintf("Expected HTTP 200, got: %d", httpResp.StatusCode))
		return
	}
	data.ID = types.StringValue(fmt.Sprintf("%d", manufacturer.GetId()))
	data.Name = types.StringValue(manufacturer.GetName())
	data.Slug = types.StringValue(manufacturer.GetSlug())
	if manufacturer.HasDescription() {
		data.Description = types.StringValue(manufacturer.GetDescription())
	} else {
		data.Description = types.StringNull()
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ManufacturerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ManufacturerResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	manufacturerID := data.ID.ValueString()
	var manufacturerIDInt int32
	if _, err := fmt.Sscanf(manufacturerID, "%d", &manufacturerIDInt); err != nil {
		resp.Diagnostics.AddError("Invalid Manufacturer ID", fmt.Sprintf("Manufacturer ID must be a number, got: %s", manufacturerID))
		return
	}
	httpResp, err := r.client.DcimAPI.DcimManufacturersDestroy(ctx, manufacturerIDInt).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Error deleting manufacturer", utils.FormatAPIError(fmt.Sprintf("delete manufacturer ID %s", manufacturerID), err, httpResp))
		return
	}
	if httpResp.StatusCode != 204 {
		resp.Diagnostics.AddError("Error deleting manufacturer", fmt.Sprintf("Expected HTTP 204, got: %d", httpResp.StatusCode))
		return
	}
}

func (r *ManufacturerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
