package resources

import (
	"context"
	"fmt"

	"github.com/bab3l/go-netbox"
	nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
	"github.com/bab3l/terraform-provider-netbox/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &WebhookResource{}
var _ resource.ResourceWithImportState = &WebhookResource{}

func NewWebhookResource() resource.Resource {
	return &WebhookResource{}
}

// WebhookResource defines the webhook resource implementation.
type WebhookResource struct {
	client *netbox.APIClient
}

// WebhookResourceModel describes the webhook resource data model.
type WebhookResourceModel struct {
	ID                types.String `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	Description       types.String `tfsdk:"description"`
	PayloadURL        types.String `tfsdk:"payload_url"`
	HTTPMethod        types.String `tfsdk:"http_method"`
	HTTPContentType   types.String `tfsdk:"http_content_type"`
	AdditionalHeaders types.String `tfsdk:"additional_headers"`
	BodyTemplate      types.String `tfsdk:"body_template"`
	Secret            types.String `tfsdk:"secret"`
	SSLVerification   types.Bool   `tfsdk:"ssl_verification"`
	CAFilePath        types.String `tfsdk:"ca_file_path"`
	Tags              types.Set    `tfsdk:"tags"`
	CustomFields      types.Set    `tfsdk:"custom_fields"`
}

func (r *WebhookResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_webhook"
}

func (r *WebhookResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a webhook in Netbox. Webhooks allow Netbox to send HTTP requests to external systems when certain events occur.",
		Attributes: map[string]schema.Attribute{
			"id":            nbschema.IDAttribute("webhook"),
			"name":          nbschema.NameAttribute("webhook", 150),
			"description":   nbschema.DescriptionAttribute("webhook"),
			"tags":          nbschema.TagsAttribute(),
			"custom_fields": nbschema.CustomFieldsAttribute(),
			"payload_url": schema.StringAttribute{
				MarkdownDescription: "The URL that will be called when the webhook is triggered. Jinja2 template processing is supported.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(500),
				},
			},
			"http_method": schema.StringAttribute{
				MarkdownDescription: "The HTTP method used when calling the webhook URL. Valid values: `GET`, `POST`, `PUT`, `PATCH`, `DELETE`. Defaults to `POST`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("POST"),
				Validators: []validator.String{
					stringvalidator.OneOf("GET", "POST", "PUT", "PATCH", "DELETE"),
				},
			},
			"http_content_type": schema.StringAttribute{
				MarkdownDescription: "The HTTP content type header. Defaults to `application/json`.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("application/json"),
				Validators: []validator.String{
					stringvalidator.LengthAtMost(100),
				},
			},
			"additional_headers": schema.StringAttribute{
				MarkdownDescription: "Additional HTTP headers to include in the request. Headers should be defined in the format `Name: Value`. Jinja2 template processing is supported.",
				Optional:            true,
			},
			"body_template": schema.StringAttribute{
				MarkdownDescription: "Jinja2 template for a custom request body. If blank, a JSON object representing the change will be included.",
				Optional:            true,
			},
			"secret": schema.StringAttribute{
				MarkdownDescription: "Secret key for HMAC signature. When provided, the request will include an `X-Hook-Signature` header containing a HMAC hex digest of the payload body.",
				Optional:            true,
				Sensitive:           true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(255),
				},
			},
			"ssl_verification": schema.BoolAttribute{
				MarkdownDescription: "Enable SSL certificate verification. Disable with caution! Defaults to `true`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"ca_file_path": schema.StringAttribute{
				MarkdownDescription: "The specific CA certificate file to use for SSL verification. Leave blank to use the system defaults.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtMost(4096),
				},
			},
		},
	}
}

func (r *WebhookResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *WebhookResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data WebhookResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	webhookRequest := netbox.NewWebhookRequest(
		data.Name.ValueString(),
		data.PayloadURL.ValueString(),
	)

	// Set optional fields
	utils.ApplyDescription(webhookRequest, data.Description)

	if !data.AdditionalHeaders.IsNull() {
		webhookRequest.SetAdditionalHeaders(data.AdditionalHeaders.ValueString())
	} else {
		webhookRequest.SetAdditionalHeaders("")
	}

	if !data.BodyTemplate.IsNull() {
		webhookRequest.SetBodyTemplate(data.BodyTemplate.ValueString())
	} else {
		webhookRequest.SetBodyTemplate("")
	}

	if !data.Secret.IsNull() {
		webhookRequest.SetSecret(data.Secret.ValueString())
	} else {
		webhookRequest.SetSecret("")
	}

	// Set HTTP method
	if !data.HTTPMethod.IsNull() && !data.HTTPMethod.IsUnknown() {
		method := netbox.PatchedWebhookRequestHttpMethod(data.HTTPMethod.ValueString())
		webhookRequest.HttpMethod = &method
	}

	// Set HTTP content type
	if !data.HTTPContentType.IsNull() && !data.HTTPContentType.IsUnknown() {
		contentType := data.HTTPContentType.ValueString()
		webhookRequest.HttpContentType = &contentType
	}

	// Set SSL verification
	if !data.SSLVerification.IsNull() && !data.SSLVerification.IsUnknown() {
		sslVerify := data.SSLVerification.ValueBool()
		webhookRequest.SslVerification = &sslVerify
	}

	// Set CA file path
	if !data.CAFilePath.IsNull() && !data.CAFilePath.IsUnknown() {
		webhookRequest.CaFilePath = *netbox.NewNullableString(utils.StringPtr(data.CAFilePath))
	}

	// Apply metadata fields (tags and custom fields)
	// Note: We can't use ApplyCommonFields since webhook doesn't have comments
	utils.ApplyTags(ctx, webhookRequest, data.Tags, &resp.Diagnostics)
	utils.ApplyCustomFields(ctx, webhookRequest, data.CustomFields, &resp.Diagnostics)

	webhook, httpResp, err := r.client.ExtrasAPI.ExtrasWebhooksCreate(ctx).WebhookRequest(*webhookRequest).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError("Error creating webhook", utils.FormatAPIError("create webhook", err, httpResp))
		return
	}
	if httpResp.StatusCode != 201 {
		resp.Diagnostics.AddError("Error creating webhook", fmt.Sprintf("Expected HTTP 201, got: %d", httpResp.StatusCode))
		return
	}
	if webhook == nil {
		resp.Diagnostics.AddError("Webhook API returned nil", "No webhook object returned from Netbox API.")
		return
	}

	// Map response to state
	r.mapWebhookToState(ctx, webhook, &data, &resp.Diagnostics)

	tflog.Debug(ctx, "Created webhook", map[string]interface{}{
		"id":   data.ID.ValueString(),
		"name": data.Name.ValueString(),
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *WebhookResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data WebhookResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	webhookID, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid Webhook ID", fmt.Sprintf("Webhook ID must be a number, got: %s", data.ID.ValueString()))
		return
	}

	webhook, httpResp, err := r.client.ExtrasAPI.ExtrasWebhooksRetrieve(ctx, webhookID).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading webhook", utils.FormatAPIError("read webhook", err, httpResp))
		return
	}

	r.mapWebhookToState(ctx, webhook, &data, &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *WebhookResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data WebhookResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	webhookID, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid Webhook ID", fmt.Sprintf("Webhook ID must be a number, got: %s", data.ID.ValueString()))
		return
	}

	webhookRequest := netbox.NewWebhookRequest(
		data.Name.ValueString(),
		data.PayloadURL.ValueString(),
	)

	// Set optional fields
	utils.ApplyDescription(webhookRequest, data.Description)

	if !data.AdditionalHeaders.IsNull() {
		webhookRequest.SetAdditionalHeaders(data.AdditionalHeaders.ValueString())
	} else {
		webhookRequest.SetAdditionalHeaders("")
	}

	if !data.BodyTemplate.IsNull() {
		webhookRequest.SetBodyTemplate(data.BodyTemplate.ValueString())
	} else {
		webhookRequest.SetBodyTemplate("")
	}

	if !data.Secret.IsNull() {
		webhookRequest.SetSecret(data.Secret.ValueString())
	} else {
		webhookRequest.SetSecret("")
	}

	// Set HTTP method
	if !data.HTTPMethod.IsNull() && !data.HTTPMethod.IsUnknown() {
		method := netbox.PatchedWebhookRequestHttpMethod(data.HTTPMethod.ValueString())
		webhookRequest.HttpMethod = &method
	}

	// Set HTTP content type
	if !data.HTTPContentType.IsNull() && !data.HTTPContentType.IsUnknown() {
		contentType := data.HTTPContentType.ValueString()
		webhookRequest.HttpContentType = &contentType
	}

	// Set SSL verification
	if !data.SSLVerification.IsNull() && !data.SSLVerification.IsUnknown() {
		sslVerify := data.SSLVerification.ValueBool()
		webhookRequest.SslVerification = &sslVerify
	}

	// Set CA file path
	if !data.CAFilePath.IsNull() && !data.CAFilePath.IsUnknown() {
		webhookRequest.CaFilePath = *netbox.NewNullableString(utils.StringPtr(data.CAFilePath))
	}

	// Handle tags and custom fields
	utils.ApplyTags(ctx, webhookRequest, data.Tags, &resp.Diagnostics)
	utils.ApplyCustomFields(ctx, webhookRequest, data.CustomFields, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	webhook, httpResp, err := r.client.ExtrasAPI.ExtrasWebhooksUpdate(ctx, webhookID).WebhookRequest(*webhookRequest).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		resp.Diagnostics.AddError("Error updating webhook", utils.FormatAPIError("update webhook", err, httpResp))
		return
	}
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError("Error updating webhook", fmt.Sprintf("Expected HTTP 200, got: %d", httpResp.StatusCode))
		return
	}

	r.mapWebhookToState(ctx, webhook, &data, &resp.Diagnostics)

	tflog.Debug(ctx, "Updated webhook", map[string]interface{}{
		"id":   data.ID.ValueString(),
		"name": data.Name.ValueString(),
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *WebhookResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data WebhookResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	webhookID, err := utils.ParseID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid Webhook ID", fmt.Sprintf("Webhook ID must be a number, got: %s", data.ID.ValueString()))
		return
	}

	httpResp, err := r.client.ExtrasAPI.ExtrasWebhooksDestroy(ctx, webhookID).Execute()
	defer utils.CloseResponseBody(httpResp)
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			// Already deleted, nothing to do
			return
		}
		resp.Diagnostics.AddError("Error deleting webhook", utils.FormatAPIError("delete webhook", err, httpResp))
		return
	}

	tflog.Debug(ctx, "Deleted webhook", map[string]interface{}{
		"id":   data.ID.ValueString(),
		"name": data.Name.ValueString(),
	})
}

func (r *WebhookResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// mapWebhookToState maps a Netbox Webhook to the Terraform state model.
func (r *WebhookResource) mapWebhookToState(ctx context.Context, webhook *netbox.Webhook, data *WebhookResourceModel, diags *diag.Diagnostics) {
	data.ID = types.StringValue(fmt.Sprintf("%d", webhook.GetId()))
	data.Name = types.StringValue(webhook.GetName())
	data.PayloadURL = types.StringValue(webhook.GetPayloadUrl())

	// Map optional string fields
	if desc := webhook.GetDescription(); desc != "" {
		data.Description = types.StringValue(desc)
	} else {
		data.Description = types.StringNull()
	}

	// Map HTTP method
	if webhook.HttpMethod != nil {
		data.HTTPMethod = types.StringValue(string(*webhook.HttpMethod))
	} else {
		data.HTTPMethod = types.StringValue("POST")
	}

	// Map HTTP content type
	if contentType := webhook.GetHttpContentType(); contentType != "" {
		data.HTTPContentType = types.StringValue(contentType)
	} else {
		data.HTTPContentType = types.StringValue("application/json")
	}

	// Map additional headers
	if headers := webhook.GetAdditionalHeaders(); headers != "" {
		data.AdditionalHeaders = types.StringValue(headers)
	} else {
		data.AdditionalHeaders = types.StringNull()
	}

	// Map body template
	if body := webhook.GetBodyTemplate(); body != "" {
		data.BodyTemplate = types.StringValue(body)
	} else {
		data.BodyTemplate = types.StringNull()
	}

	// Secret is write-only - we can't read it back from the API
	// Keep existing state value if set
	if data.Secret.IsNull() || data.Secret.IsUnknown() {
		data.Secret = types.StringNull()
	}

	// Map SSL verification
	if webhook.SslVerification != nil {
		data.SSLVerification = types.BoolValue(*webhook.SslVerification)
	} else {
		data.SSLVerification = types.BoolValue(true)
	}

	// Map CA file path
	if webhook.CaFilePath.IsSet() && webhook.CaFilePath.Get() != nil && *webhook.CaFilePath.Get() != "" {
		data.CAFilePath = types.StringValue(*webhook.CaFilePath.Get())
	} else {
		data.CAFilePath = types.StringNull()
	}

	// Map display_name
	// Handle tags using consolidated helper
	data.Tags = utils.PopulateTagsFromAPI(ctx, webhook.HasTags(), webhook.GetTags(), data.Tags, diags)
	if diags.HasError() {
		return
	}

	// Map custom fields
	data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, data.CustomFields, webhook.GetCustomFields(), diags)
	if diags.HasError() {
		return
	}
}
