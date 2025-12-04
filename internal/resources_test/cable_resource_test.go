package resources_test

import (
	"context"
	"testing"

	"github.com/bab3l/go-netbox"
	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestCableResource(t *testing.T) {
	t.Parallel()

	r := resources.NewCableResource()
	if r == nil {
		t.Fatal("Expected non-nil cable resource")
	}
}

func TestCableResourceSchema(t *testing.T) {
	t.Parallel()

	r := resources.NewCableResource()
	schemaRequest := fwresource.SchemaRequest{}
	schemaResponse := &fwresource.SchemaResponse{}

	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)
	}

	if schemaResponse.Schema.Attributes == nil {
		t.Fatal("Expected schema to have attributes")
	}

	// Required attributes - a_terminations and b_terminations are required
	requiredAttrs := []string{"a_terminations", "b_terminations"}
	for _, attr := range requiredAttrs {
		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {
			t.Errorf("Expected required attribute %s to exist in schema", attr)
		}
	}

	// Computed attributes
	computedAttrs := []string{"id"}
	for _, attr := range computedAttrs {
		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {
			t.Errorf("Expected computed attribute %s to exist in schema", attr)
		}
	}

	// Optional attributes
	optionalAttrs := []string{
		"type", "status", "tenant", "label", "color",
		"length", "length_unit", "description", "comments",
		"tags", "custom_fields",
	}
	for _, attr := range optionalAttrs {
		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {
			t.Errorf("Expected optional attribute %s to exist in schema", attr)
		}
	}
}

func TestCableResourceMetadata(t *testing.T) {
	t.Parallel()

	r := resources.NewCableResource()
	metadataRequest := fwresource.MetadataRequest{
		ProviderTypeName: "netbox",
	}
	metadataResponse := &fwresource.MetadataResponse{}

	r.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_cable"
	if metadataResponse.TypeName != expected {
		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)
	}
}

func TestCableResourceConfigure(t *testing.T) {
	t.Parallel()

	r := resources.NewCableResource().(*resources.CableResource)

	// Test with nil provider data (should not error)
	configureRequest := fwresource.ConfigureRequest{
		ProviderData: nil,
	}
	configureResponse := &fwresource.ConfigureResponse{}

	r.Configure(context.Background(), configureRequest, configureResponse)

	if configureResponse.Diagnostics.HasError() {
		t.Fatalf("Configure with nil provider data should not error: %+v", configureResponse.Diagnostics)
	}

	// Test with correct provider data type
	client := &netbox.APIClient{}
	configureRequest.ProviderData = client
	configureResponse = &fwresource.ConfigureResponse{}

	r.Configure(context.Background(), configureRequest, configureResponse)

	if configureResponse.Diagnostics.HasError() {
		t.Errorf("Expected no error with correct provider data, got: %+v", configureResponse.Diagnostics)
	}
}

func TestCableResourceConfigureWrongType(t *testing.T) {
	t.Parallel()

	r := resources.NewCableResource().(*resources.CableResource)

	// Test with wrong type provider data (should error)
	configureRequest := fwresource.ConfigureRequest{
		ProviderData: "wrong-type",
	}
	configureResponse := &fwresource.ConfigureResponse{}

	r.Configure(context.Background(), configureRequest, configureResponse)

	if !configureResponse.Diagnostics.HasError() {
		t.Fatal("Configure with wrong type provider data should error")
	}
}
