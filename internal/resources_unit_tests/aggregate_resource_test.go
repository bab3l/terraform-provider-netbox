package resources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestAggregateResource(t *testing.T) {

	t.Parallel()

	t.Parallel()

	r := resources.NewAggregateResource()

	if r == nil {

		t.Fatal("Expected non-nil Aggregate resource")

	}

}

func TestAggregateResourceSchema(t *testing.T) {

	t.Parallel()

	t.Parallel()

	r := resources.NewAggregateResource()

	schemaRequest := fwresource.SchemaRequest{}

	schemaResponse := &fwresource.SchemaResponse{}

	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {

		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)

	}

	if schemaResponse.Schema.Attributes == nil {

		t.Fatal("Expected schema to have attributes")

	}

	// Required attributes

	requiredAttrs := []string{"prefix", "rir"}

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

	optionalAttrs := []string{"tenant", "date_added", "description", "comments", "tags", "custom_fields"}

	for _, attr := range optionalAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected optional attribute %s to exist in schema", attr)

		}

	}

}

func TestAggregateResourceMetadata(t *testing.T) {

	t.Parallel()

	t.Parallel()

	r := resources.NewAggregateResource()

	metadataRequest := fwresource.MetadataRequest{

		ProviderTypeName: "netbox",
	}

	metadataResponse := &fwresource.MetadataResponse{}

	r.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_aggregate"

	if metadataResponse.TypeName != expected {

		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)

	}

}

func TestAggregateResourceConfigure(t *testing.T) {

	t.Parallel()

	t.Parallel()

	r := resources.NewAggregateResource()

	// Type assert to access Configure method

	configurable, ok := r.(fwresource.ResourceWithConfigure)

	if !ok {

		t.Fatal("Resource does not implement ResourceWithConfigure")

	}

	configureRequest := fwresource.ConfigureRequest{

		ProviderData: nil,
	}

	configureResponse := &fwresource.ConfigureResponse{}

	configurable.Configure(context.Background(), configureRequest, configureResponse)

	if configureResponse.Diagnostics.HasError() {

		t.Errorf("Expected no error with nil provider data, got: %+v", configureResponse.Diagnostics)

	}

}
