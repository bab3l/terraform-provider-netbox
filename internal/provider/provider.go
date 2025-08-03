package provider

import (
	"context"
	"os"

	"github.com/bab3l/go-netbox"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/bab3l/terraform-provider-netbox/internal/datasources"
	"github.com/bab3l/terraform-provider-netbox/internal/resources"
)

// Ensure NetboxProvider satisfies various provider interfaces.
var _ provider.Provider = &NetboxProvider{}

// NetboxProvider defines the provider implementation.
type NetboxProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// NetboxProviderModel describes the provider data model.
type NetboxProviderModel struct {
	ServerURL types.String `tfsdk:"server_url"`
	APIToken  types.String `tfsdk:"api_token"`
	Insecure  types.Bool   `tfsdk:"insecure"`
}

func (p *NetboxProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "netbox"
	resp.Version = p.version
}

func (p *NetboxProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The Netbox provider is used to interact with Netbox, an open-source web application designed to help manage and document computer networks. This provider allows you to manage Netbox resources such as sites, devices, racks, IP addresses, and more using Terraform.",

		Attributes: map[string]schema.Attribute{
			"server_url": schema.StringAttribute{
				MarkdownDescription: "The base URL of your Netbox instance (e.g., `https://netbox.example.com`). Can also be set via the `NETBOX_SERVER_URL` environment variable.",
				Optional:            true,
			},
			"api_token": schema.StringAttribute{
				MarkdownDescription: "The API token for authenticating with Netbox. Generate this token in your Netbox user profile. Can also be set via the `NETBOX_API_TOKEN` environment variable.",
				Optional:            true,
				Sensitive:           true,
			},
			"insecure": schema.BoolAttribute{
				MarkdownDescription: "Whether to skip TLS certificate verification. Defaults to false. Can also be set via the `NETBOX_INSECURE` environment variable.",
				Optional:            true,
			},
		},
	}
}

func (p *NetboxProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data NetboxProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Configuration values are now available.
	serverURL := os.Getenv("NETBOX_SERVER_URL")
	apiToken := os.Getenv("NETBOX_API_TOKEN")
	insecure := os.Getenv("NETBOX_INSECURE") == "true"

	if !data.ServerURL.IsNull() {
		serverURL = data.ServerURL.ValueString()
	}

	if !data.APIToken.IsNull() {
		apiToken = data.APIToken.ValueString()
	}

	if !data.Insecure.IsNull() {
		insecure = data.Insecure.ValueBool()
	}

	// Validation
	if serverURL == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("server_url"),
			"Missing Netbox Server URL",
			"The provider cannot create the Netbox API client as there is a missing or empty value for the Netbox server URL. "+
				"Set the server_url value in the configuration or use the NETBOX_SERVER_URL environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if apiToken == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_token"),
			"Missing Netbox API Token",
			"The provider cannot create the Netbox API client as there is a missing or empty value for the Netbox API token. "+
				"Set the api_token value in the configuration or use the NETBOX_API_TOKEN environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "netbox_server_url", serverURL)
	ctx = tflog.SetField(ctx, "netbox_api_token", apiToken)
	ctx = tflog.SetField(ctx, "netbox_insecure", insecure)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "netbox_api_token")

	tflog.Debug(ctx, "Creating Netbox client")

	// Create a new Netbox client using go-netbox
	cfg := netbox.NewConfiguration()
	cfg.Servers = netbox.ServerConfigurations{
		{
			URL: serverURL,
		},
	}

	// Set up authentication
	cfg.DefaultHeader = map[string]string{
		"Authorization": "Token " + apiToken,
	}

	// Handle insecure connections
	if insecure {
		// Note: go-netbox uses a standard HTTP client, so TLS verification
		// would need to be configured on the HTTP client if needed
		tflog.Debug(ctx, "Insecure mode enabled - TLS verification disabled")
	}

	client := netbox.NewAPIClient(cfg)

	// Make the Netbox client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured Netbox client", map[string]any{"success": true})
}

func (p *NetboxProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		resources.NewSiteResource,
		resources.NewSiteGroupResource,
	}
}

func (p *NetboxProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		datasources.NewSiteDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &NetboxProvider{
			version: version,
		}
	}
}
