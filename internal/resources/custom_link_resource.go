package resources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bab3l/go-netbox"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &CustomLinkResource{}
var _ resource.ResourceWithImportState = &CustomLinkResource{}

func NewCustomLinkResource() resource.Resource {
	return &CustomLinkResource{}
}

// CustomLinkResource defines the resource implementation.
type CustomLinkResource struct {
	client *netbox.APIClient
}

// CustomLinkResourceModel describes the resource data model.
type CustomLinkResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	ObjectTypes types.List   `tfsdk:"object_types"`
	Enabled     types.Bool   `tfsdk:"enabled"`
	LinkText    types.String `tfsdk:"link_text"`
	LinkURL     types.String `tfsdk:"link_url"`
	Weight      types.Int64  `tfsdk:"weight"`
	GroupName   types.String `tfsdk:"group_name"`
	ButtonClass types.String `tfsdk:"button_class"`
	NewWindow   types.Bool   `tfsdk:"new_window"`
}

func (r *CustomLinkResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_custom_link"
}

func (r *CustomLinkResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a custom link in Netbox. Custom links allow you to add dynamic links to object detail pages.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier (assigned by Netbox).",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the custom link.",
				Required:            true,
			},
			"object_types": schema.ListAttribute{
				MarkdownDescription: "List of object types this link applies to (e.g., `[\"dcim.device\", \"dcim.site\"]`).",
				Required:            true,
				ElementType:         types.StringType,
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Whether the custom link is enabled. Defaults to true.",
				Optional:            true,
				Computed:            true,
			},
			"link_text": schema.StringAttribute{
				MarkdownDescription: "Jinja2 template code for the link text.",
				Required:            true,
			},
			"link_url": schema.StringAttribute{
				MarkdownDescription: "Jinja2 template code for the link URL.",
				Required:            true,
			},
			"weight": schema.Int64Attribute{
				MarkdownDescription: "Weight for ordering. Lower values appear first. Defaults to 100.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.Int64{
					int64validator.Between(0, 32767),
				},
			},
			"group_name": schema.StringAttribute{
				MarkdownDescription: "Links with the same group name will appear as a dropdown menu.",
				Optional:            true,
			},
			"button_class": schema.StringAttribute{
				MarkdownDescription: "CSS class for the button. Valid values: `default`, `blue`, `indigo`, `purple`, `pink`, `red`, `orange`, `yellow`, `green`, `teal`, `cyan`, `gray`, `black`, `white`, `ghost-dark`.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"default", "blue", "indigo", "purple", "pink", "red",
						"orange", "yellow", "green", "teal", "cyan", "gray",
						"black", "white", "ghost-dark",
					),
				},
			},
			"new_window": schema.BoolAttribute{
				MarkdownDescription: "Whether to open the link in a new window. Defaults to false.",
				Optional:            true,
				Computed:            true,
			},
		},
	}
}

func (r *CustomLinkResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *CustomLinkResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data CustomLinkResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build object_types list
	var objectTypes []string
	diags := data.ObjectTypes.ElementsAs(ctx, &objectTypes, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	request := netbox.NewCustomLinkRequest(
		objectTypes,
		data.Name.ValueString(),
		data.LinkText.ValueString(),
		data.LinkURL.ValueString(),
	)

	// Set optional fields
	if !data.Enabled.IsNull() && !data.Enabled.IsUnknown() {
		request.SetEnabled(data.Enabled.ValueBool())
	}

	if !data.Weight.IsNull() && !data.Weight.IsUnknown() {
		weight, err := utils.SafeInt32FromValue(data.Weight)
		if err != nil {
			resp.Diagnostics.AddError("Invalid value", fmt.Sprintf("Weight value overflow: %s", err))
			return
		}
		request.SetWeight(weight)
	}

	if !data.GroupName.IsNull() && !data.GroupName.IsUnknown() {
		request.SetGroupName(data.GroupName.ValueString())
	}

	if !data.ButtonClass.IsNull() && !data.ButtonClass.IsUnknown() {
		buttonClass := netbox.CustomLinkButtonClass(data.ButtonClass.ValueString())

		request.SetButtonClass(buttonClass)
	}

	if !data.NewWindow.IsNull() && !data.NewWindow.IsUnknown() {
		request.SetNewWindow(data.NewWindow.ValueBool())
	}

	result, httpResp, err := r.client.ExtrasAPI.ExtrasCustomLinksCreate(ctx).
		CustomLinkRequest(*request).Execute()

	if err != nil {
		resp.Diagnostics.AddError("Error creating custom link",
			utils.FormatAPIError("create custom link", err, httpResp))
		return
	}

	if httpResp.StatusCode != http.StatusCreated {
		resp.Diagnostics.AddError("Error creating custom link",
			fmt.Sprintf("Expected HTTP 201, got: %d", httpResp.StatusCode))
		return
	}

	r.mapToState(ctx, result, &data)
	tflog.Debug(ctx, "Created custom link", map[string]interface{}{
		"id":   data.ID.ValueString(),
		"name": data.Name.ValueString(),
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CustomLinkResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data CustomLinkResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID",
			fmt.Sprintf("ID must be a number, got: %s", data.ID.ValueString()))
		return
	}

	result, httpResp, err := r.client.ExtrasAPI.ExtrasCustomLinksRetrieve(ctx, id).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError("Error reading custom link",
			utils.FormatAPIError(fmt.Sprintf("read custom link ID %s", data.ID.ValueString()), err, httpResp))
		return
	}

	r.mapToState(ctx, result, &data)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CustomLinkResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data CustomLinkResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID",
			fmt.Sprintf("ID must be a number, got: %s", data.ID.ValueString()))
		return
	}

	// Build object_types list
	var objectTypes []string
	diags := data.ObjectTypes.ElementsAs(ctx, &objectTypes, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := netbox.NewCustomLinkRequest(
		objectTypes,
		data.Name.ValueString(),
		data.LinkText.ValueString(),
		data.LinkURL.ValueString(),
	)

	// Set optional fields
	if !data.Enabled.IsNull() && !data.Enabled.IsUnknown() {
		enabled := data.Enabled.ValueBool()
		request.Enabled = &enabled
	}

	if !data.Weight.IsNull() && !data.Weight.IsUnknown() {
		weight, err := utils.SafeInt32FromValue(data.Weight)
		if err != nil {
			resp.Diagnostics.AddError("Invalid value", fmt.Sprintf("Weight value overflow: %s", err))
			return
		}
		request.Weight = &weight
	}

	if !data.GroupName.IsNull() && !data.GroupName.IsUnknown() {
		groupName := data.GroupName.ValueString()
		request.GroupName = &groupName
	}

	if !data.ButtonClass.IsNull() && !data.ButtonClass.IsUnknown() {
		buttonClass := netbox.CustomLinkButtonClass(data.ButtonClass.ValueString())
		request.SetButtonClass(buttonClass)
	}

	if !data.NewWindow.IsNull() && !data.NewWindow.IsUnknown() {
		request.SetNewWindow(data.NewWindow.ValueBool())
	}

	result, httpResp, err := r.client.ExtrasAPI.ExtrasCustomLinksUpdate(ctx, id).
		CustomLinkRequest(*request).Execute()

	if err != nil {
		resp.Diagnostics.AddError("Error updating custom link",
			utils.FormatAPIError(fmt.Sprintf("update custom link ID %s", data.ID.ValueString()), err, httpResp))
		return
	}

	if httpResp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError("Error updating custom link",
			fmt.Sprintf("Expected HTTP 200, got: %d", httpResp.StatusCode))
		return
	}

	r.mapToState(ctx, result, &data)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CustomLinkResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CustomLinkResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID",
			fmt.Sprintf("ID must be a number, got: %s", data.ID.ValueString()))
		return
	}
	httpResp, err := r.client.ExtrasAPI.ExtrasCustomLinksDestroy(ctx, id).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		// If the resource was already deleted (404), consider it a success
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			tflog.Debug(ctx, "Custom link already deleted", map[string]interface{}{"id": id})
			return
		}
		resp.Diagnostics.AddError("Error deleting custom link",
			utils.FormatAPIError(fmt.Sprintf("delete custom link ID %s", data.ID.ValueString()), err, httpResp))
		return
	}
	if httpResp.StatusCode != http.StatusNoContent {
		resp.Diagnostics.AddError("Error deleting custom link",
			fmt.Sprintf("Expected HTTP 204, got: %d", httpResp.StatusCode))
		return
	}
}

func (r *CustomLinkResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// mapToState maps API response to Terraform state.
func (r *CustomLinkResource) mapToState(ctx context.Context, result *netbox.CustomLink, data *CustomLinkResourceModel) {
	data.ID = types.StringValue(fmt.Sprintf("%d", result.GetId()))
	data.Name = types.StringValue(result.GetName())
	data.LinkText = types.StringValue(result.GetLinkText())
	data.LinkURL = types.StringValue(result.GetLinkUrl())

	// Handle object_types
	if len(result.GetObjectTypes()) > 0 {
		objectTypesValue, _ := types.ListValueFrom(ctx, types.StringType, result.GetObjectTypes())
		data.ObjectTypes = objectTypesValue
	} else {
		data.ObjectTypes = types.ListNull(types.StringType)
	}

	if result.HasEnabled() {
		data.Enabled = types.BoolValue(result.GetEnabled())
	} else {
		data.Enabled = types.BoolNull()
	}

	if result.HasWeight() {
		data.Weight = types.Int64Value(int64(result.GetWeight()))
	} else {
		data.Weight = types.Int64Null()
	}

	if result.HasGroupName() && result.GetGroupName() != "" {
		data.GroupName = types.StringValue(result.GetGroupName())
	} else {
		data.GroupName = types.StringNull()
	}

	if result.HasButtonClass() {
		data.ButtonClass = types.StringValue(string(*result.ButtonClass))
	} else {
		data.ButtonClass = types.StringNull()
	}

	if result.HasNewWindow() {
		data.NewWindow = types.BoolValue(result.GetNewWindow())
	} else {
		data.NewWindow = types.BoolNull()
	}
}
