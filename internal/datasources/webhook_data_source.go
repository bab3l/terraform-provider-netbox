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
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &WebhookDataSource{}

func NewWebhookDataSource() datasource.DataSource {
	return &WebhookDataSource{}
}

// WebhookDataSource defines the webhook data source implementation.
type WebhookDataSource struct {
	client *netbox.APIClient
}

// WebhookDataSourceModel describes the webhook data source data model.
type WebhookDataSourceModel struct {
	ID                types.String `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	Description       types.String `tfsdk:"description"`
	DisplayName       types.String `tfsdk:"display_name"`
	PayloadURL        types.String `tfsdk:"payload_url"`
	HTTPMethod        types.String `tfsdk:"http_method"`
	HTTPContentType   types.String `tfsdk:"http_content_type"`
	AdditionalHeaders types.String `tfsdk:"additional_headers"`
	BodyTemplate      types.String `tfsdk:"body_template"`
	SSLVerification   types.Bool   `tfsdk:"ssl_verification"`
	CAFilePath        types.String `tfsdk:"ca_file_path"`
	Tags              types.List   `tfsdk:"tags"`
	CustomFields      types.Set    `tfsdk:"custom_fields"`
}

func (d *WebhookDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_webhook"
}

func (d *WebhookDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get information about a webhook in Netbox.",
		Attributes: map[string]schema.Attribute{
			"id":                 nbschema.DSIDAttribute("webhook"),
			"name":               nbschema.DSNameAttribute("webhook"),
			"description":        nbschema.DSComputedStringAttribute("Description of the webhook."),
			"display_name":       nbschema.DSComputedStringAttribute("Display name for the webhook."),
			"payload_url":        nbschema.DSComputedStringAttribute("The URL that will be called when the webhook is triggered."),
			"http_method":        nbschema.DSComputedStringAttribute("The HTTP method used when calling the webhook URL."),
			"http_content_type":  nbschema.DSComputedStringAttribute("The HTTP content type header."),
			"additional_headers": nbschema.DSComputedStringAttribute("Additional HTTP headers to include in the request."),
			"body_template":      nbschema.DSComputedStringAttribute("Jinja2 template for a custom request body."),
			"ssl_verification":   nbschema.DSComputedBoolAttribute("Whether SSL certificate verification is enabled."),
			"ca_file_path":       nbschema.DSComputedStringAttribute("The specific CA certificate file to use for SSL verification."),
			"tags": schema.ListAttribute{
				MarkdownDescription: "Tags assigned to this webhook.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"custom_fields": nbschema.DSCustomFieldsAttribute(),
		},
	}
}

func (d *WebhookDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *WebhookDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data WebhookDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var webhook *netbox.Webhook
	var err error
	var httpResp *http.Response

	// Lookup by id or name
	switch {
	case !data.ID.IsNull() && !data.ID.IsUnknown() && data.ID.ValueString() != "":
		webhookID, parseErr := utils.ParseID(data.ID.ValueString())
		if parseErr != nil {
			resp.Diagnostics.AddError("Invalid Webhook ID", "Webhook ID must be a number.")
			return
		}
		webhook, httpResp, err = d.client.ExtrasAPI.ExtrasWebhooksRetrieve(ctx, webhookID).Execute()
		defer utils.CloseResponseBody(httpResp)
		if err != nil {
			resp.Diagnostics.AddError("Error reading webhook", utils.FormatAPIError("read webhook", err, httpResp))
			return
		}
	case !data.Name.IsNull() && !data.Name.IsUnknown() && data.Name.ValueString() != "":
		name := data.Name.ValueString()
		webhooks, httpResp, listErr := d.client.ExtrasAPI.ExtrasWebhooksList(ctx).Name([]string{name}).Execute()
		defer utils.CloseResponseBody(httpResp)
		if listErr != nil {
			resp.Diagnostics.AddError("Error reading webhook", utils.FormatAPIError("read webhook by name", listErr, httpResp))
			return
		}
		if webhooks == nil || len(webhooks.GetResults()) == 0 {
			resp.Diagnostics.AddError("Webhook Not Found", fmt.Sprintf("No webhook found with name: %s", name))
			return
		}
		webhook = &webhooks.GetResults()[0]
	default:
		resp.Diagnostics.AddError("Missing Webhook Identifier", "Either 'id' or 'name' must be specified.")
		return
	}

	if webhook == nil {
		resp.Diagnostics.AddError("Webhook Not Found", "No webhook found with the specified identifier.")
		return
	}

	// Map API response to state
	d.mapWebhookToState(ctx, webhook, &data, &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// mapWebhookToState maps API response to Terraform state.
func (d *WebhookDataSource) mapWebhookToState(ctx context.Context, webhook *netbox.Webhook, data *WebhookDataSourceModel, diags *diag.Diagnostics) {
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

	// Map display name
	if displayName := webhook.GetDisplay(); displayName != "" {
		data.DisplayName = types.StringValue(displayName)
	} else {
		data.DisplayName = types.StringNull()
	}

	// Handle tags (slug list)
	data.Tags = utils.PopulateTagsSlugListFromAPI(ctx, webhook.HasTags(), webhook.GetTags(), diags)

	// Map custom fields
	if webhook.HasCustomFields() {
		customFields := utils.MapAllCustomFieldsToModels(webhook.GetCustomFields())
		customFieldsValue, cfDiags := types.SetValueFrom(ctx, utils.GetCustomFieldsAttributeType().ElemType, customFields)
		if !cfDiags.HasError() {
			data.CustomFields = customFieldsValue
		}
	} else {
		data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)
	}
}
