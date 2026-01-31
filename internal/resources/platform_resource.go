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
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &PlatformResource{}
	_ resource.ResourceWithImportState = &PlatformResource{}
)

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
		MarkdownDescription: "Manages a platform in Netbox. Platforms represent the software running on a device, such as an operating system or firmware version.",
		Attributes: map[string]schema.Attribute{
			"id":           nbschema.IDAttribute("platform"),
			"name":         nbschema.NameAttribute("platform", 100),
			"slug":         nbschema.SlugAttribute("platform"),
			"manufacturer": nbschema.ReferenceAttributeWithDiffSuppress("manufacturer", "Reference to the manufacturer (ID or slug)."),
		},
	}

	// Add description attribute
	maps.Copy(resp.Schema.Attributes, nbschema.DescriptionOnlyAttributes("platform"))
}

// Implement Create, Read, Update, Delete, and ImportState methods here.
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
		manufacturerRef, diags := netboxlookup.LookupManufacturer(ctx, r.client, data.Manufacturer.ValueString())
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		platformRequest.Manufacturer = *netbox.NewNullableBriefManufacturerRequest(manufacturerRef)
	}

	// Apply description
	utils.ApplyDescription(platformRequest, data.Description)
	platform, httpResp, err := r.client.DcimAPI.DcimPlatformsCreate(ctx).PlatformRequest(*platformRequest).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError("Error creating platform", utils.FormatAPIError("create platform", err, httpResp))
		return
	}
	if httpResp.StatusCode != http.StatusCreated {
		resp.Diagnostics.AddError("Error creating platform", fmt.Sprintf("Expected HTTP 201, got: %d", httpResp.StatusCode))
		return
	}
	data.ID = types.StringValue(fmt.Sprintf("%d", platform.GetId()))
	data.Name = types.StringValue(platform.GetName())
	data.Slug = types.StringValue(platform.GetSlug())
	if platform.HasManufacturer() {
		m, ok := platform.GetManufacturerOk()
		if ok && m != nil {
			data.Manufacturer = utils.UpdateReferenceAttribute(data.Manufacturer, m.Slug, "", m.Id)
		}
	}
	if platform.HasDescription() {
		desc := platform.GetDescription()
		// Preserve null if original was null and API returns empty string
		if desc == "" && data.Description.IsNull() {
			data.Description = types.StringNull()
		} else {
			data.Description = types.StringValue(desc)
		}
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
	platformIDInt, err := utils.ParseID(platformID)
	if err != nil {
		resp.Diagnostics.AddError("Invalid Platform ID", fmt.Sprintf("Platform ID must be a number, got: %s", platformID))
		return
	}
	platform, httpResp, err := r.client.DcimAPI.DcimPlatformsRetrieve(ctx, platformIDInt).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading platform", utils.FormatAPIError(fmt.Sprintf("read platform ID %s", platformID), err, httpResp))
		return
	}
	if httpResp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError("Error reading platform", fmt.Sprintf("Expected HTTP 200, got: %d", httpResp.StatusCode))
		return
	}
	data.ID = types.StringValue(fmt.Sprintf("%d", platform.GetId()))
	data.Name = types.StringValue(platform.GetName())
	data.Slug = types.StringValue(platform.GetSlug())
	if platform.HasManufacturer() {
		m, ok := platform.GetManufacturerOk()
		if ok && m != nil {
			data.Manufacturer = utils.UpdateReferenceAttribute(data.Manufacturer, m.Slug, "", m.Id)
		} else {
			data.Manufacturer = types.StringNull()
		}
	} else {
		data.Manufacturer = types.StringNull()
	}
	if platform.HasDescription() {
		desc := platform.GetDescription()
		// Preserve null if original was null and API returns empty strings
		if desc == "" && data.Description.IsNull() {
			data.Description = types.StringNull()
		} else {
			data.Description = types.StringValue(desc)
		}
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
	platformIDInt, err := utils.ParseID(platformID)
	if err != nil {
		resp.Diagnostics.AddError("Invalid Platform ID", fmt.Sprintf("Platform ID must be a number, got: %s", platformID))
		return
	}
	platformRequest := netbox.PlatformRequest{
		Name: data.Name.ValueString(),
		Slug: data.Slug.ValueString(),
	}

	if !data.Manufacturer.IsNull() {
		manufacturerRef, diags := netboxlookup.LookupManufacturer(ctx, r.client, data.Manufacturer.ValueString())
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		platformRequest.Manufacturer = *netbox.NewNullableBriefManufacturerRequest(manufacturerRef)
	} else if data.Manufacturer.IsNull() {
		// Explicitly set to nil when removed from config
		platformRequest.SetManufacturerNil()
	}

	// Apply description
	utils.ApplyDescription(&platformRequest, data.Description)
	platform, httpResp, err := r.client.DcimAPI.DcimPlatformsUpdate(ctx, platformIDInt).PlatformRequest(platformRequest).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError("Error updating platform", utils.FormatAPIError(fmt.Sprintf("update platform ID %s", platformID), err, httpResp))
		return
	}
	if httpResp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError("Error updating platform", fmt.Sprintf("Expected HTTP 200, got: %d", httpResp.StatusCode))
		return
	}
	data.ID = types.StringValue(fmt.Sprintf("%d", platform.GetId()))
	data.Name = types.StringValue(platform.GetName())
	data.Slug = types.StringValue(platform.GetSlug())
	if platform.HasManufacturer() {
		m, ok := platform.GetManufacturerOk()
		if ok && m != nil {
			data.Manufacturer = utils.UpdateReferenceAttribute(data.Manufacturer, m.Slug, "", m.Id)
		} else {
			data.Manufacturer = types.StringNull()
		}
	} else {
		data.Manufacturer = types.StringNull()
	}
	if platform.HasDescription() {
		desc := platform.GetDescription()
		// Preserve null if original was null and API returns empty string
		if desc == "" && data.Description.IsNull() {
			data.Description = types.StringNull()
		} else {
			data.Description = types.StringValue(desc)
		}
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
	platformIDInt, err := utils.ParseID(platformID)
	if err != nil {
		resp.Diagnostics.AddError("Invalid Platform ID", fmt.Sprintf("Platform ID must be a number, got: %s", platformID))
		return
	}
	httpResp, err := r.client.DcimAPI.DcimPlatformsDestroy(ctx, platformIDInt).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			return // Already deleted
		}
		resp.Diagnostics.AddError("Error deleting platform", utils.FormatAPIError(fmt.Sprintf("delete platform ID %s", platformID), err, httpResp))
		return
	}
	if httpResp.StatusCode != http.StatusNoContent {
		resp.Diagnostics.AddError("Error deleting platform", fmt.Sprintf("Expected HTTP 204, got: %d", httpResp.StatusCode))
		return
	}
}

func (r *PlatformResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	utils.ImportStatePassthroughIDWithValidation(ctx, req, resp, path.Root("id"), true)
}
