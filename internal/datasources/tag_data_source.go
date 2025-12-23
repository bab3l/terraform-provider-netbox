package datasources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bab3l/go-netbox"
	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &TagDataSource{}

func NewTagDataSource() datasource.DataSource {
	return &TagDataSource{}
}

// TagDataSource defines the tag data source implementation.
type TagDataSource struct {
	client *netbox.APIClient
}

// TagDataSourceModel describes the tag data source data model.
type TagDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Slug        types.String `tfsdk:"slug"`
	Color       types.String `tfsdk:"color"`
	Description types.String `tfsdk:"description"`
	DisplayName types.String `tfsdk:"display_name"`
	ObjectTypes types.List   `tfsdk:"object_types"`
}

func (d *TagDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tag"
}

func (d *TagDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get information about a tag in Netbox.",
		Attributes: map[string]schema.Attribute{
			"id":          nbschema.DSIDAttribute("tag"),
			"name":        nbschema.DSNameAttribute("tag"),
			"slug":        nbschema.DSSlugAttribute("tag"),
			"color":       nbschema.DSComputedStringAttribute("Color of the tag (6-character hex code)."),
			"description": nbschema.DSComputedStringAttribute("Description of the tag."), "display_name": nbschema.DSComputedStringAttribute("The display name of the tag."), "object_types": schema.ListAttribute{
				MarkdownDescription: "List of object types this tag can be applied to.",
				Computed:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

func (d *TagDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*netbox.APIClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *netbox.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = client
}

func (d *TagDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data TagDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var tag *netbox.Tag
	var err error
	var httpResp *http.Response

	// Lookup by id, slug, or name
	switch {
	case !data.ID.IsNull() && !data.ID.IsUnknown():
		tagID, parseErr := utils.ParseID(data.ID.ValueString())
		if parseErr != nil {
			resp.Diagnostics.AddError("Invalid Tag ID", "Tag ID must be a number.")
			return
		}
		tag, httpResp, err = d.client.ExtrasAPI.ExtrasTagsRetrieve(ctx, tagID).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError("Error reading tag", utils.FormatAPIError("read tag", err, httpResp))
			return
		}
	case !data.Slug.IsNull() && !data.Slug.IsUnknown():
		slug := data.Slug.ValueString()
		tags, httpResp, listErr := d.client.ExtrasAPI.ExtrasTagsList(ctx).Slug([]string{slug}).Execute()
		defer utils.CloseResponseBody(httpResp)
		if listErr != nil {
			resp.Diagnostics.AddError("Error reading tag", utils.FormatAPIError("read tag by slug", listErr, httpResp))
			return
		}
		if tags == nil || len(tags.GetResults()) == 0 {
			resp.Diagnostics.AddError("Tag Not Found", fmt.Sprintf("No tag found with slug: %s", slug))
			return
		}
		tag = &tags.GetResults()[0]
	case !data.Name.IsNull() && !data.Name.IsUnknown():
		name := data.Name.ValueString()
		tags, httpResp, listErr := d.client.ExtrasAPI.ExtrasTagsList(ctx).Name([]string{name}).Execute()
		defer utils.CloseResponseBody(httpResp)
		if listErr != nil {
			resp.Diagnostics.AddError("Error reading tag", utils.FormatAPIError("read tag by name", listErr, httpResp))
			return
		}
		if tags == nil || len(tags.GetResults()) == 0 {
			resp.Diagnostics.AddError("Tag Not Found", fmt.Sprintf("No tag found with name: %s", name))
			return
		}
		tag = &tags.GetResults()[0]
	default:
		resp.Diagnostics.AddError("Missing Tag Identifier", "Either 'id', 'slug', or 'name' must be specified.")
		return
	}

	if tag == nil {
		resp.Diagnostics.AddError("Tag Not Found", "No tag found with the specified identifier.")
		return
	}

	// Map API response to state
	d.mapTagToState(ctx, tag, &data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// mapTagToState maps API response to Terraform state.
func (d *TagDataSource) mapTagToState(ctx context.Context, tag *netbox.Tag, data *TagDataSourceModel) {
	data.ID = types.StringValue(fmt.Sprintf("%d", tag.GetId()))
	data.Name = types.StringValue(tag.GetName())
	data.Slug = types.StringValue(tag.GetSlug())
	data.Color = utils.StringFromAPI(tag.HasColor(), tag.GetColor, data.Color)
	data.Description = utils.StringFromAPI(tag.HasDescription(), tag.GetDescription, data.Description)
	data.DisplayName = utils.StringFromAPI(tag.GetDisplay() != "", tag.GetDisplay, data.DisplayName)

	// Handle object_types
	if tag.HasObjectTypes() && len(tag.GetObjectTypes()) > 0 {
		objectTypesValue, _ := types.ListValueFrom(ctx, types.StringType, tag.GetObjectTypes())
		data.ObjectTypes = objectTypesValue
	} else {
		data.ObjectTypes = types.ListNull(types.StringType)
	}
}
