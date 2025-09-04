package resources_test

import (
	"context"
	"testing"

	"github.com/bab3l/go-netbox"
	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestTenantGroupResource(t *testing.T) {
	t.Parallel()

	r := resources.NewTenantGroupResource()
	if r == nil {
		t.Fatal("Expected non-nil tenant group resource")
	}
}

func TestTenantGroupResourceSchema(t *testing.T) {
	t.Parallel()

	r := resources.NewTenantGroupResource()
	schemaRequest := resource.SchemaRequest{}
	schemaResponse := &resource.SchemaResponse{}

	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)
	}

	if schemaResponse.Schema.Attributes == nil {
		t.Fatal("Expected schema to have attributes")
	}

	// Check required attributes
	requiredAttrs := []string{"name", "slug"}
	for _, attr := range requiredAttrs {
		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {
			t.Errorf("Expected required attribute %s to exist in schema", attr)
		}
	}

	// Check optional attributes
	optionalAttrs := []string{"parent", "description", "tags", "custom_fields"}
	for _, attr := range optionalAttrs {
		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {
			t.Errorf("Expected optional attribute %s to exist in schema", attr)
		}
	}

	// Check computed attributes
	computedAttrs := []string{"id"}
	for _, attr := range computedAttrs {
		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {
			t.Errorf("Expected computed attribute %s to exist in schema", attr)
		}
	}
}

func TestTenantGroupResourceMetadata(t *testing.T) {
	t.Parallel()

	r := resources.NewTenantGroupResource()
	metadataRequest := resource.MetadataRequest{
		ProviderTypeName: "netbox",
	}
	metadataResponse := &resource.MetadataResponse{}

	r.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_tenant_group"
	if metadataResponse.TypeName != expected {
		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)
	}
}

func TestTenantGroupResourceConfigure(t *testing.T) {
	t.Parallel()

	r := resources.NewTenantGroupResource().(*resources.TenantGroupResource)

	// Test with nil provider data
	configureRequest := resource.ConfigureRequest{
		ProviderData: nil,
	}
	configureResponse := &resource.ConfigureResponse{}

	r.Configure(context.Background(), configureRequest, configureResponse)

	if configureResponse.Diagnostics.HasError() {
		t.Errorf("Expected no error with nil provider data, got: %+v", configureResponse.Diagnostics)
	}

	// Test with correct provider data
	client := &netbox.APIClient{}
	configureRequest.ProviderData = client
	configureResponse = &resource.ConfigureResponse{}

	r.Configure(context.Background(), configureRequest, configureResponse)

	if configureResponse.Diagnostics.HasError() {
		t.Errorf("Expected no error with correct provider data, got: %+v", configureResponse.Diagnostics)
	}

	if r.client != client {
		t.Error("Expected client to be set")
	}

	// Test with incorrect provider data
	configureRequest.ProviderData = "invalid"
	configureResponse = &resource.ConfigureResponse{}

	r.Configure(context.Background(), configureRequest, configureResponse)

	if !configureResponse.Diagnostics.HasError() {
		t.Error("Expected error with incorrect provider data")
	}
}
