package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"netbox": providerserver.NewProtocol6WithError(New("test")()),
}

func TestProvider(t *testing.T) {
	// Test that the provider can be instantiated
	p := New("test")()
	if p == nil {
		t.Fatal("Provider should not be nil")
	}
}

func TestProviderSchema(t *testing.T) {
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
	ctx := context.Background()
	p := New("test")()

	// Test that the provider provides expected resources
	resources := p.Resources(ctx)

	// Verify we have a reasonable number of resources (at least 60)
	// The actual count will grow as more resources are implemented
	minResourceCount := 60
	if len(resources) < minResourceCount {
		t.Errorf("Provider should provide at least %d resources, got %d", minResourceCount, len(resources))
	}

	// Verify all resources can be instantiated
	for i, resourceFunc := range resources {
		resource := resourceFunc()
		if resource == nil {
			t.Errorf("Resource %d should not be nil", i)
		}
	}
}

func TestProviderDataSources(t *testing.T) {
	ctx := context.Background()
	p := New("test")()

	// Test that the provider provides data sources (even if empty for now)
	dataSources := p.DataSources(ctx)

	// This is expected to be empty currently, but the call shouldn't panic
	_ = dataSources
}
