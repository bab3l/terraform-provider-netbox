// Package resources contains Terraform resource implementations for the Netbox provider.

package resources

import (
	"context"
	"fmt"
	"maps"

	"github.com/bab3l/go-netbox"
	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.

var (
	_ resource.Resource = &ExportTemplateResource{}

	_ resource.ResourceWithConfigure = &ExportTemplateResource{}

	_ resource.ResourceWithImportState = &ExportTemplateResource{}
)

// NewExportTemplateResource returns a new resource implementing the export template resource.

func NewExportTemplateResource() resource.Resource {
	return &ExportTemplateResource{}
}

// ExportTemplateResource defines the resource implementation.

type ExportTemplateResource struct {
	client *netbox.APIClient
}

// ExportTemplateResourceModel describes the resource data model.

type ExportTemplateResourceModel struct {
	ID types.String `tfsdk:"id"`

	ObjectTypes types.Set `tfsdk:"object_types"`

	Name types.String `tfsdk:"name"`

	Description types.String `tfsdk:"description"`

	DisplayName types.String `tfsdk:"display_name"`

	TemplateCode types.String `tfsdk:"template_code"`

	MimeType types.String `tfsdk:"mime_type"`

	FileExtension types.String `tfsdk:"file_extension"`

	AsAttachment types.Bool `tfsdk:"as_attachment"`
}

// Metadata returns the resource type name.

func (r *ExportTemplateResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_export_template"
}

// Schema defines the schema for the resource.

func (r *ExportTemplateResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an export template in NetBox. Export templates define Jinja2 templates for exporting data from NetBox.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique numeric ID of the export template.",

				Computed: true,

				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"object_types": schema.SetAttribute{
				MarkdownDescription: "Set of object types this template applies to (e.g., 'dcim.device', 'dcim.site').",

				Required: true,

				ElementType: types.StringType,
			},

			"name": nbschema.NameAttribute("export template", 100),

			"display_name": nbschema.DisplayNameAttribute("export template"),
			"template_code": schema.StringAttribute{
				MarkdownDescription: "Jinja2 template code. The list of objects being exported is passed as a context variable named `queryset`.",

				Required: true,
			},

			"mime_type": schema.StringAttribute{
				MarkdownDescription: "MIME type for the rendered output. Defaults to `text/plain; charset=utf-8`.",

				Optional: true,
			},

			"file_extension": schema.StringAttribute{
				MarkdownDescription: "Extension to append to the rendered filename.",

				Optional: true,
			},

			"as_attachment": schema.BoolAttribute{
				MarkdownDescription: "Download file as attachment. Defaults to `true`.",

				Optional: true,

				Computed: true,

				Default: booldefault.StaticBool(true),
			},
		},
	}

	// Add description attribute
	maps.Copy(resp.Schema.Attributes, nbschema.DescriptionOnlyAttributes("export template"))
}

// Configure adds the provider configured client to the resource.

func (r *ExportTemplateResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create creates a new export template.

func (r *ExportTemplateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ExportTemplateResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Convert object types to string slice

	var objectTypes []string

	if !data.ObjectTypes.IsNull() && !data.ObjectTypes.IsUnknown() {
		diags := data.ObjectTypes.ElementsAs(ctx, &objectTypes, false)

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {
			return
		}
	}

	// Build the API request

	exportTemplateRequest := netbox.NewExportTemplateRequest(

		objectTypes,

		data.Name.ValueString(),

		data.TemplateCode.ValueString(),
	)

	// Set optional fields

	utils.ApplyDescription(exportTemplateRequest, data.Description)

	exportTemplateRequest.MimeType = utils.StringPtr(data.MimeType)

	exportTemplateRequest.FileExtension = utils.StringPtr(data.FileExtension)

	if !data.AsAttachment.IsNull() && !data.AsAttachment.IsUnknown() {
		asAttachment := data.AsAttachment.ValueBool()

		exportTemplateRequest.AsAttachment = &asAttachment
	}

	tflog.Debug(ctx, "Creating export template", map[string]interface{}{
		"name": data.Name.ValueString(),

		"object_types": objectTypes,
	})

	// Call the API

	exportTemplate, httpResp, err := r.client.ExtrasAPI.ExtrasExportTemplatesCreate(ctx).
		ExportTemplateRequest(*exportTemplateRequest).
		Execute()

	if err != nil {
		resp.Diagnostics.AddError(

			"Error creating export template",

			utils.FormatAPIError("create export template", err, httpResp),
		)

		return
	}

	// Map response to state

	r.mapResponseToState(ctx, exportTemplate, &data, &resp.Diagnostics)

	tflog.Debug(ctx, "Created export template", map[string]interface{}{
		"id": exportTemplate.GetId(),
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read reads the export template.

func (r *ExportTemplateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ExportTemplateResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Parse ID

	id, err := utils.ParseID(data.ID.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(

			"Invalid ID format",

			fmt.Sprintf("Could not parse export template ID '%s': %s", data.ID.ValueString(), err),
		)

		return
	}

	tflog.Debug(ctx, "Reading export template", map[string]interface{}{
		"id": id,
	})

	// Call the API

	exportTemplate, httpResp, err := r.client.ExtrasAPI.ExtrasExportTemplatesRetrieve(ctx, id).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			resp.State.RemoveResource(ctx)

			return
		}

		resp.Diagnostics.AddError(

			"Error reading export template",

			utils.FormatAPIError("read export template", err, httpResp),
		)

		return
	}

	// Map response to state

	r.mapResponseToState(ctx, exportTemplate, &data, &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates the export template.

func (r *ExportTemplateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ExportTemplateResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Parse ID

	id, err := utils.ParseID(data.ID.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(

			"Invalid ID format",

			fmt.Sprintf("Could not parse export template ID '%s': %s", data.ID.ValueString(), err),
		)

		return
	}

	// Convert object types to string slice

	var objectTypes []string

	if !data.ObjectTypes.IsNull() && !data.ObjectTypes.IsUnknown() {
		diags := data.ObjectTypes.ElementsAs(ctx, &objectTypes, false)

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {
			return
		}
	}

	// Build the API request

	exportTemplateRequest := netbox.NewExportTemplateRequest(

		objectTypes,

		data.Name.ValueString(),

		data.TemplateCode.ValueString(),
	)

	// Set optional fields

	utils.ApplyDescription(exportTemplateRequest, data.Description)

	exportTemplateRequest.MimeType = utils.StringPtr(data.MimeType)

	exportTemplateRequest.FileExtension = utils.StringPtr(data.FileExtension)

	if !data.AsAttachment.IsNull() && !data.AsAttachment.IsUnknown() {
		asAttachment := data.AsAttachment.ValueBool()

		exportTemplateRequest.AsAttachment = &asAttachment
	}

	tflog.Debug(ctx, "Updating export template", map[string]interface{}{
		"id": id,

		"name": data.Name.ValueString(),
	})

	// Call the API

	exportTemplate, httpResp, err := r.client.ExtrasAPI.ExtrasExportTemplatesUpdate(ctx, id).
		ExportTemplateRequest(*exportTemplateRequest).
		Execute()

	if err != nil {
		resp.Diagnostics.AddError(

			"Error updating export template",

			utils.FormatAPIError("update export template", err, httpResp),
		)

		return
	}

	// Map response to state

	r.mapResponseToState(ctx, exportTemplate, &data, &resp.Diagnostics)

	tflog.Debug(ctx, "Updated export template", map[string]interface{}{
		"id": exportTemplate.GetId(),
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete deletes the export template.

func (r *ExportTemplateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ExportTemplateResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Parse ID

	id, err := utils.ParseID(data.ID.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(

			"Invalid ID format",

			fmt.Sprintf("Could not parse export template ID '%s': %s", data.ID.ValueString(), err),
		)

		return
	}

	tflog.Debug(ctx, "Deleting export template", map[string]interface{}{
		"id": id,
	})

	httpResp, err := r.client.ExtrasAPI.ExtrasExportTemplatesDestroy(ctx, id).Execute()

	defer utils.CloseResponseBody(httpResp)

	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			return
		}

		resp.Diagnostics.AddError(

			"Error deleting export template",

			utils.FormatAPIError("delete export template", err, httpResp),
		)

		return
	}

	tflog.Debug(ctx, "Deleted export template", map[string]interface{}{
		"id": id,
	})
}

// ImportState imports an existing export template.

func (r *ExportTemplateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// mapResponseToState maps the API response to the Terraform state.

func (r *ExportTemplateResource) mapResponseToState(ctx context.Context, exportTemplate *netbox.ExportTemplate, data *ExportTemplateResourceModel, diags *diag.Diagnostics) {
	data.ID = types.StringValue(fmt.Sprintf("%d", exportTemplate.GetId()))

	data.Name = types.StringValue(exportTemplate.GetName())

	data.TemplateCode = types.StringValue(exportTemplate.GetTemplateCode())

	// Handle object types

	if exportTemplate.ObjectTypes != nil && len(exportTemplate.ObjectTypes) > 0 {
		objectTypesSet, d := types.SetValueFrom(ctx, types.StringType, exportTemplate.ObjectTypes)

		diags.Append(d...)

		data.ObjectTypes = objectTypesSet
	} else {
		data.ObjectTypes = types.SetNull(types.StringType)
	}

	// Handle description

	if exportTemplate.HasDescription() && exportTemplate.GetDescription() != "" {
		data.Description = types.StringValue(exportTemplate.GetDescription())
	} else {
		data.Description = types.StringNull()
	}

	// Handle mime type

	if exportTemplate.HasMimeType() && exportTemplate.GetMimeType() != "" {
		data.MimeType = types.StringValue(exportTemplate.GetMimeType())
	} else {
		data.MimeType = types.StringNull()
	}

	// Handle file extension

	if exportTemplate.HasFileExtension() && exportTemplate.GetFileExtension() != "" {
		data.FileExtension = types.StringValue(exportTemplate.GetFileExtension())
	} else {
		data.FileExtension = types.StringNull()
	}

	// Handle as_attachment

	if exportTemplate.HasAsAttachment() {
		data.AsAttachment = types.BoolValue(exportTemplate.GetAsAttachment())
	} else {
		data.AsAttachment = types.BoolValue(true) // Default
	}

	// Handle display_name (computed field, always set a value)
	display := exportTemplate.GetDisplay()
	if display != "" {
		data.DisplayName = types.StringValue(display)
	} else {
		data.DisplayName = types.StringValue("")
	}
}
