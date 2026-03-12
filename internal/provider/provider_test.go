package provider

import (
	"context"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

const (
	expectedResourceCount   = 104
	expectedDataSourceCount = 109
)

// TestAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
// Note: Exported for use by resource and datasource acceptance tests.
//
//	var TestAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
//		"netbox": providerserver.NewProtocol6WithError(New("test")()),
//	}
func TestProvider(t *testing.T) {
	t.Parallel()

	// Test that the provider can be instantiated
	p := New("test")()
	if p == nil {
		t.Fatal("Provider should not be nil")
	}
}

func TestProviderSchema(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	p := New("test")()

	// Test that the provider schema can be retrieved
	schemaReq := provider.SchemaRequest{}
	schemaResp := &provider.SchemaResponse{}
	p.Schema(ctx, schemaReq, schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("Provider schema should not have errors: %v", schemaResp.Diagnostics.Errors())
	}

	// Verify essential attributes exist
	attrs := schemaResp.Schema.Attributes
	if _, ok := attrs["server_url"]; !ok {
		t.Error("Provider schema should include server_url attribute")
	}
	if _, ok := attrs["api_token"]; !ok {
		t.Error("Provider schema should include api_token attribute")
	}
	if _, ok := attrs["insecure"]; !ok {
		t.Error("Provider schema should include insecure attribute")
	}
}

func TestProviderResources(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	p := New("test")()

	resources := p.Resources(ctx)
	if len(resources) != expectedResourceCount {
		t.Fatalf("Provider should provide exactly %d resources, got %d", expectedResourceCount, len(resources))
	}

	resourceNames := make(map[string]struct{}, len(resources))
	for i, resourceFunc := range resources {
		instance := resourceFunc()
		if instance == nil {
			t.Fatalf("Resource %d should not be nil", i)
		}

		metadataResp := &resource.MetadataResponse{}
		instance.Metadata(ctx, resource.MetadataRequest{ProviderTypeName: "netbox"}, metadataResp)
		if metadataResp.TypeName == "" {
			t.Fatalf("Resource %d returned an empty type name", i)
		}
		if _, exists := resourceNames[metadataResp.TypeName]; exists {
			t.Fatalf("Duplicate resource type registered: %s", metadataResp.TypeName)
		}
		resourceNames[metadataResp.TypeName] = struct{}{}
	}

	for _, name := range []string{
		"netbox_site",
		"netbox_device",
		"netbox_virtual_machine",
		"netbox_ip_address",
		"netbox_event_rule",
		"netbox_export_template",
	} {
		if _, ok := resourceNames[name]; !ok {
			t.Errorf("Provider should register resource %s", name)
		}
	}
}

func TestProviderDataSources(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	p := New("test")()

	dataSources := p.DataSources(ctx)
	if len(dataSources) != expectedDataSourceCount {
		t.Fatalf("Provider should provide exactly %d data sources, got %d", expectedDataSourceCount, len(dataSources))
	}

	dataSourceNames := make(map[string]struct{}, len(dataSources))
	for i, dataSourceFunc := range dataSources {
		instance := dataSourceFunc()
		if instance == nil {
			t.Fatalf("Data source %d should not be nil", i)
		}

		metadataResp := &datasource.MetadataResponse{}
		instance.Metadata(ctx, datasource.MetadataRequest{ProviderTypeName: "netbox"}, metadataResp)
		if metadataResp.TypeName == "" {
			t.Fatalf("Data source %d returned an empty type name", i)
		}
		if _, exists := dataSourceNames[metadataResp.TypeName]; exists {
			t.Fatalf("Duplicate data source type registered: %s", metadataResp.TypeName)
		}
		dataSourceNames[metadataResp.TypeName] = struct{}{}
	}

	for _, name := range []string{
		"netbox_site",
		"netbox_device",
		"netbox_devices",
		"netbox_virtual_device_context",
		"netbox_script",
		"netbox_export_template",
	} {
		if _, ok := dataSourceNames[name]; !ok {
			t.Errorf("Provider should register data source %s", name)
		}
	}
}

func TestNewNetboxClient_InsecureTLS(t *testing.T) {
	t.Parallel()

	cfg := newNetboxConfiguration("https://netbox.example.com", "token", true)
	if cfg == nil {
		t.Fatal("Expected configuration to be created")
	}

	if cfg.HTTPClient == nil {
		t.Fatal("Expected configured HTTP client")
	}

	transport, ok := cfg.HTTPClient.Transport.(*http.Transport)
	if !ok {
		t.Fatalf("Expected HTTP transport to be *http.Transport, got %T", cfg.HTTPClient.Transport)
	}
	if transport.TLSClientConfig == nil {
		t.Fatal("Expected TLS client config to be set")
	}
	if !transport.TLSClientConfig.InsecureSkipVerify {
		t.Fatal("Expected insecure TLS verification to be enabled")
	}

	if got := cfg.DefaultHeader["Authorization"]; got != "Token token" {
		t.Fatalf("Expected authorization header to be configured, got %q", got)
	}

	defaultTransport, ok := http.DefaultTransport.(*http.Transport)
	if !ok {
		t.Fatal("Expected default transport to be *http.Transport")
	}
	if defaultTransport.TLSClientConfig != nil && defaultTransport.TLSClientConfig.InsecureSkipVerify {
		t.Fatal("Default transport TLS config should not be mutated")
	}
}

func TestNewNetboxClient_DefaultTLS(t *testing.T) {
	t.Parallel()

	cfg := newNetboxConfiguration("https://netbox.example.com", "token", false)
	if cfg == nil {
		t.Fatal("Expected configuration to be created")
	}

	if cfg.HTTPClient == nil {
		t.Fatal("Expected configured HTTP client")
	}

	if cfg.HTTPClient != http.DefaultClient {
		t.Fatal("Expected secure client to reuse http.DefaultClient")
	}
}
