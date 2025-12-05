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

var _ resource.Resource = &TagResource{}
var _ resource.ResourceWithImportState = &TagResource{}

func NewTagResource() resource.Resource {
	return &TagResource{}
}

// TagResource defines the tag resource implementation.
type TagResource struct {
	client *netbox.APIClient
}

// TagResourceModel describes the tag resource data model.
type TagResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Slug        types.String `tfsdk:"slug"`
	Color       types.String `tfsdk:"color"`
	Description types.String `tfsdk:"description"`
	ObjectTypes types.List   `tfsdk:"object_types"`
}

func (r *TagResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tag"
}

func (r *TagResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a tag in Netbox. Tags can be applied to most objects for categorization and filtering.",
		Attributes: map[string]schema.Attribute{
			"id":          nbschema.IDAttribute("tag"),
			"name":        nbschema.NameAttribute("tag", 100),
			"slug":        nbschema.SlugAttribute("tag"),
			"color":       nbschema.ComputedColorAttribute("tag"),
			"description": nbschema.DescriptionAttribute("tag"),
			"object_types": schema.ListAttribute{
				MarkdownDescription: "List of object types this tag can be applied to. If empty, the tag can be applied to any object type. Example: `[\"dcim.device\", \"dcim.site\"]`",
				Optional:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

func (r *TagResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *TagResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data TagResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tagRequest := netbox.NewTagRequest(
		data.Name.ValueString(),
		data.Slug.ValueString(),
	)

	// Set optional fields
	tagRequest.Color = utils.StringPtr(data.Color)
	tagRequest.Description = utils.StringPtr(data.Description)

	// Handle object_types list
	if !data.ObjectTypes.IsNull() && !data.ObjectTypes.IsUnknown() {
		var objectTypes []string
		diags := data.ObjectTypes.ElementsAs(ctx, &objectTypes, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		tagRequest.ObjectTypes = objectTypes
	}

	tag, httpResp, err := r.client.ExtrasAPI.ExtrasTagsCreate(ctx).TagRequest(*tagRequest).Execute()

	// Handle the case where the tag was created but the response parsing fails
	// (go-netbox expects tagged_items in response but Netbox CREATE doesn't return it)
	if err != nil && httpResp != nil && httpResp.StatusCode == 201 {
		// Tag was created - extract ID from response body and do a read
		tagID := utils.ExtractIDFromResponse(httpResp)
		if tagID > 0 {
			tflog.Debug(ctx, "Tag created but response parsing failed, fetching by ID", map[string]interface{}{"id": tagID})
			tag, httpResp, err = r.client.ExtrasAPI.ExtrasTagsRetrieve(ctx, tagID).Execute()
		}
	}

	if err != nil {
		resp.Diagnostics.AddError("Error creating tag", utils.FormatAPIError("create tag", err, httpResp))
		return
	}
	if httpResp.StatusCode != 201 && httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError("Error creating tag", fmt.Sprintf("Expected HTTP 201 or 200, got: %d", httpResp.StatusCode))
		return
	}
	if tag == nil {
		resp.Diagnostics.AddError("Tag API returned nil", "No tag object returned from Netbox API.")
		return
	}

	// Map response to state
	r.mapTagToState(ctx, tag, &data)

	tflog.Debug(ctx, "Created tag", map[string]interface{}{
		"id":   data.ID.ValueString(),
		"name": data.Name.ValueString(),
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TagResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data TagResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tagID, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid Tag ID", fmt.Sprintf("Tag ID must be a number, got: %s", data.ID.ValueString()))
		return
	}

	tag, httpResp, err := r.client.ExtrasAPI.ExtrasTagsRetrieve(ctx, tagID).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Error reading tag", utils.FormatAPIError(fmt.Sprintf("read tag ID %s", data.ID.ValueString()), err, httpResp))
		return
	}
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError("Error reading tag", fmt.Sprintf("Expected HTTP 200, got: %d", httpResp.StatusCode))
		return
	}

	// Map response to state
	r.mapTagToState(ctx, tag, &data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TagResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data TagResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tagID, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid Tag ID", fmt.Sprintf("Tag ID must be a number, got: %s", data.ID.ValueString()))
		return
	}

	tagRequest := netbox.NewTagRequest(
		data.Name.ValueString(),
		data.Slug.ValueString(),
	)

	// Set optional fields
	tagRequest.Color = utils.StringPtr(data.Color)
	tagRequest.Description = utils.StringPtr(data.Description)

	// Handle object_types list
	if !data.ObjectTypes.IsNull() && !data.ObjectTypes.IsUnknown() {
		var objectTypes []string
		diags := data.ObjectTypes.ElementsAs(ctx, &objectTypes, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		tagRequest.ObjectTypes = objectTypes
	}

	tag, httpResp, err := r.client.ExtrasAPI.ExtrasTagsUpdate(ctx, tagID).TagRequest(*tagRequest).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Error updating tag", utils.FormatAPIError(fmt.Sprintf("update tag ID %s", data.ID.ValueString()), err, httpResp))
		return
	}
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError("Error updating tag", fmt.Sprintf("Expected HTTP 200, got: %d", httpResp.StatusCode))
		return
	}

	// Map response to state
	r.mapTagToState(ctx, tag, &data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TagResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data TagResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tagID, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid Tag ID", fmt.Sprintf("Tag ID must be a number, got: %s", data.ID.ValueString()))
		return
	}

	httpResp, err := r.client.ExtrasAPI.ExtrasTagsDestroy(ctx, tagID).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Error deleting tag", utils.FormatAPIError(fmt.Sprintf("delete tag ID %s", data.ID.ValueString()), err, httpResp))
		return
	}
	if httpResp.StatusCode != 204 {
		resp.Diagnostics.AddError("Error deleting tag", fmt.Sprintf("Expected HTTP 204, got: %d", httpResp.StatusCode))
		return
	}
}

func (r *TagResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// mapTagToState maps API response to Terraform state.
func (r *TagResource) mapTagToState(ctx context.Context, tag *netbox.Tag, data *TagResourceModel) {
	data.ID = types.StringValue(fmt.Sprintf("%d", tag.GetId()))
	data.Name = types.StringValue(tag.GetName())
	data.Slug = types.StringValue(tag.GetSlug())
	data.Color = utils.StringFromAPI(tag.HasColor(), tag.GetColor, data.Color)
	data.Description = utils.StringFromAPI(tag.HasDescription(), tag.GetDescription, data.Description)

	// Handle object_types
	if tag.HasObjectTypes() && len(tag.GetObjectTypes()) > 0 {
		objectTypesValue, _ := types.ListValueFrom(ctx, types.StringType, tag.GetObjectTypes())
		data.ObjectTypes = objectTypesValue
	} else {
		data.ObjectTypes = types.ListNull(types.StringType)
	}
}
