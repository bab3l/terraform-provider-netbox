package resources

import (
	"context"
	"fmt"

	"github.com/bab3l/go-netbox"
	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
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
	ID types.String `tfsdk:"id"`

	Name types.String `tfsdk:"name"`

	Slug types.String `tfsdk:"slug"`

	DisplayName types.String `tfsdk:"display_name"`

	Description types.String `tfsdk:"description"`
}

func (r *ManufacturerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {

	resp.TypeName = req.ProviderTypeName + "_manufacturer"

}

func (r *ManufacturerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	resp.Schema = schema.Schema{

		MarkdownDescription: "Manages a manufacturer in Netbox. Manufacturers are used to group devices and platforms by vendor.",

		Attributes: map[string]schema.Attribute{

			"id": nbschema.IDAttribute("manufacturer"),

			"name": nbschema.NameAttribute("manufacturer", 100),

			"slug": nbschema.SlugAttribute("manufacturer"),

			"display_name": nbschema.DisplayNameAttribute("manufacturer"),

			"description": nbschema.DescriptionAttribute("manufacturer"),
		},
	}

}

// Implement Configure, Create, Read, Update, Delete, ImportState methods here.

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

	// Use helper for optional string field

	manufacturerRequest.Description = utils.StringPtr(data.Description)

	manufacturer, httpResp, err := r.client.DcimAPI.DcimManufacturersCreate(ctx).ManufacturerRequest(manufacturerRequest).Execute()

	defer utils.CloseResponseBody(httpResp)

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

	// Map response to state using helpers

	r.mapManufacturerToState(manufacturer, &data)

	tflog.Debug(ctx, "Created manufacturer", map[string]interface{}{

		"id": data.ID.ValueString(),

		"name": data.Name.ValueString(),
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

func (r *ManufacturerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	var data ManufacturerResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return

	}

	manufacturerID := data.ID.ValueString()

	var manufacturerIDInt int32

	manufacturerIDInt, err := utils.ParseID(manufacturerID)

	if err != nil {

		resp.Diagnostics.AddError("Invalid Manufacturer ID", fmt.Sprintf("Manufacturer ID must be a number, got: %s", manufacturerID))

		return

	}

	manufacturer, httpResp, err := r.client.DcimAPI.DcimManufacturersRetrieve(ctx, manufacturerIDInt).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		resp.Diagnostics.AddError("Error reading manufacturer", utils.FormatAPIError(fmt.Sprintf("read manufacturer ID %s", manufacturerID), err, httpResp))

		return

	}

	if httpResp.StatusCode != 200 {

		resp.Diagnostics.AddError("Error reading manufacturer", fmt.Sprintf("Expected HTTP 200, got: %d", httpResp.StatusCode))

		return

	}

	// Map response to state using helpers

	r.mapManufacturerToState(manufacturer, &data)

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

	manufacturerIDInt, err := utils.ParseID(manufacturerID)

	if err != nil {

		resp.Diagnostics.AddError("Invalid Manufacturer ID", fmt.Sprintf("Manufacturer ID must be a number, got: %s", manufacturerID))

		return

	}

	manufacturerRequest := netbox.ManufacturerRequest{

		Name: data.Name.ValueString(),

		Slug: data.Slug.ValueString(),
	}

	// Use helper for optional string field

	manufacturerRequest.Description = utils.StringPtr(data.Description)

	manufacturer, httpResp, err := r.client.DcimAPI.DcimManufacturersUpdate(ctx, manufacturerIDInt).ManufacturerRequest(manufacturerRequest).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {

		resp.Diagnostics.AddError("Error updating manufacturer", utils.FormatAPIError(fmt.Sprintf("update manufacturer ID %s", manufacturerID), err, httpResp))

		return

	}

	if httpResp.StatusCode != 200 {

		resp.Diagnostics.AddError("Error updating manufacturer", fmt.Sprintf("Expected HTTP 200, got: %d", httpResp.StatusCode))

		return

	}

	// Map response to state using helpers

	r.mapManufacturerToState(manufacturer, &data)

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

	manufacturerIDInt, err := utils.ParseID(manufacturerID)

	if err != nil {

		resp.Diagnostics.AddError("Invalid Manufacturer ID", fmt.Sprintf("Manufacturer ID must be a number, got: %s", manufacturerID))

		return

	}

	httpResp, err := r.client.DcimAPI.DcimManufacturersDestroy(ctx, manufacturerIDInt).Execute()

	defer utils.CloseResponseBody(httpResp)

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

// mapManufacturerToState maps API response to Terraform state using state helpers.

func (r *ManufacturerResource) mapManufacturerToState(manufacturer *netbox.Manufacturer, data *ManufacturerResourceModel) {

	data.ID = types.StringValue(fmt.Sprintf("%d", manufacturer.GetId()))

	data.Name = types.StringValue(manufacturer.GetName())

	data.Slug = types.StringValue(manufacturer.GetSlug())

	data.DisplayName = types.StringValue(manufacturer.GetDisplay())

	data.Description = utils.StringFromAPI(manufacturer.HasDescription(), manufacturer.GetDescription, data.Description)

}
