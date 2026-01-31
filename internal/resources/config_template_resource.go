// Package resources provides Terraform resource implementations for NetBox objects.

package resources

import (
	"context"
	"fmt"
	"maps"
	"net/http"

	"github.com/bab3l/go-netbox"
	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &ConfigTemplateResource{}
	_ resource.ResourceWithImportState = &ConfigTemplateResource{}
)

// NewConfigTemplateResource returns a new resource implementing the config template resource.
func NewConfigTemplateResource() resource.Resource {
	return &ConfigTemplateResource{}
}

// ConfigTemplateResource defines the resource implementation.
type ConfigTemplateResource struct {
	client *netbox.APIClient
}

// ConfigTemplateResourceModel describes the resource data model.
type ConfigTemplateResourceModel struct {
	ID           types.Int32  `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Description  types.String `tfsdk:"description"`
	TemplateCode types.String `tfsdk:"template_code"`
	DataPath     types.String `tfsdk:"data_path"`
	Tags         types.Set    `tfsdk:"tags"`
}

// Metadata returns the resource type name.
func (r *ConfigTemplateResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_config_template"
}

// Schema defines the schema for the resource.
func (r *ConfigTemplateResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a config template in NetBox. Config templates are Jinja2 templates used to render configuration for devices.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int32Attribute{
				MarkdownDescription: "The unique numeric ID of the config template.",
				Computed:            true,
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the config template.",
				Required:            true,
			},
			"template_code": schema.StringAttribute{
				MarkdownDescription: "Jinja2 template code.",
				Required:            true,
			},
			"data_path": schema.StringAttribute{
				MarkdownDescription: "Path to remote file (relative to data source root). Read-only.",
				Computed:            true,
			},
			"tags": nbschema.TagsSlugAttribute(),
		},
	}

	// Add description attribute
	maps.Copy(resp.Schema.Attributes, nbschema.DescriptionOnlyAttributes("config template"))
}

// Configure adds the provider configured client to the resource.
func (r *ConfigTemplateResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *ConfigTemplateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ConfigTemplateResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the API request
	apiReq := netbox.NewConfigTemplateRequest(
		data.Name.ValueString(),
		data.TemplateCode.ValueString(),
	)

	// Set optional fields
	utils.ApplyDescription(apiReq, data.Description)

	utils.ApplyTagsFromSlugs(ctx, r.client, apiReq, data.Tags, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating config template", map[string]interface{}{
		"name": data.Name.ValueString(),
	})
	response, httpResp, err := r.client.ExtrasAPI.ExtrasConfigTemplatesCreate(ctx).ConfigTemplateRequest(*apiReq).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating config template",
			utils.FormatAPIError("create config template", err, httpResp),
		)
		return
	}

	// Map response to model
	r.mapResponseToModel(response, &data)
	tflog.Trace(ctx, "Created config template", map[string]interface{}{
		"id": data.ID.ValueInt32(),
	})
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *ConfigTemplateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ConfigTemplateResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	templateID := data.ID.ValueInt32()
	tflog.Debug(ctx, "Reading config template", map[string]interface{}{
		"id": templateID,
	})
	response, httpResp, err := r.client.ExtrasAPI.ExtrasConfigTemplatesRetrieve(ctx, templateID).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			tflog.Debug(ctx, "Config template not found, removing from state", map[string]interface{}{
				"id": templateID,
			})
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading config template",
			utils.FormatAPIError(fmt.Sprintf("read config template ID %d", templateID), err, httpResp),
		)
		return
	}

	// Map response to model
	r.mapResponseToModel(response, &data)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *ConfigTemplateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ConfigTemplateResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	templateID := data.ID.ValueInt32()

	// Build the API request
	apiReq := netbox.NewConfigTemplateRequest(
		data.Name.ValueString(),
		data.TemplateCode.ValueString(),
	)

	// Set optional fields
	utils.ApplyDescription(apiReq, data.Description)

	utils.ApplyTagsFromSlugs(ctx, r.client, apiReq, data.Tags, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Updating config template", map[string]interface{}{
		"id": templateID,
	})
	response, httpResp, err := r.client.ExtrasAPI.ExtrasConfigTemplatesUpdate(ctx, templateID).ConfigTemplateRequest(*apiReq).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating config template",
			utils.FormatAPIError(fmt.Sprintf("update config template ID %d", templateID), err, httpResp),
		)
		return
	}

	// Map response to model
	r.mapResponseToModel(response, &data)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *ConfigTemplateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ConfigTemplateResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	templateID := data.ID.ValueInt32()
	tflog.Debug(ctx, "Deleting config template", map[string]interface{}{
		"id": templateID,
	})
	httpResp, err := r.client.ExtrasAPI.ExtrasConfigTemplatesDestroy(ctx, templateID).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			// Resource already deleted
			return
		}
		resp.Diagnostics.AddError(
			"Error deleting config template",
			utils.FormatAPIError(fmt.Sprintf("delete config template ID %d", templateID), err, httpResp),
		)
		return
	}
}

// ImportState imports the resource state from Terraform.
func (r *ConfigTemplateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Parse the import ID as an integer
	id, err := utils.ParseInt32ID(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Could not parse import ID %q as integer: %s", req.ID, err),
		)
		return
	}

	// Set the ID in state
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

// mapResponseToModel maps the API response to the Terraform model.
func (r *ConfigTemplateResource) mapResponseToModel(template *netbox.ConfigTemplate, data *ConfigTemplateResourceModel) {
	data.ID = types.Int32Value(template.GetId())
	data.Name = types.StringValue(template.GetName())

	data.TemplateCode = types.StringValue(template.GetTemplateCode())

	// Map description
	data.Description = utils.StringFromAPI(template.HasDescription(), template.GetDescription, data.Description)

	// Map data path (read-only)
	data.DataPath = types.StringValue(template.GetDataPath())

	// Tags (slug list)
	data.Tags = utils.PopulateTagsSlugFromAPI(context.Background(), template.HasTags(), template.GetTags(), data.Tags)
}
