package datasources

import (
	"context"
	"testing"

	"github.com/bab3l/go-netbox"
	fwdatasource "github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestTagDataSource(t *testing.T) {
	t.Parallel()

	ds := NewTagDataSource()
	if ds == nil {
		t.Fatal("Expected non-nil tag data source")
	}
}

func TestTagDataSourceSchema(t *testing.T) {
	t.Parallel()

	ds := NewTagDataSource()
	schemaRequest := fwdatasource.SchemaRequest{}
	schemaResponse := &fwdatasource.SchemaResponse{}

	ds.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)
	}

	if schemaResponse.Schema.Attributes == nil {
		t.Fatal("Expected schema to have attributes")
	}

	// Check lookup attributes
	lookupAttrs := []string{"id", "name", "slug"}
	for _, attr := range lookupAttrs {
		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {
			t.Errorf("Expected lookup attribute %s to exist in schema", attr)
		}
	}

	// Check computed attributes
	computedAttrs := []string{"color", "description", "object_types"}
	for _, attr := range computedAttrs {
		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {
			t.Errorf("Expected computed attribute %s to exist in schema", attr)
		}
	}
}

func TestTagDataSourceMetadata(t *testing.T) {
	t.Parallel()

	ds := NewTagDataSource()
	metadataRequest := fwdatasource.MetadataRequest{
		ProviderTypeName: "netbox",
	}
	metadataResponse := &fwdatasource.MetadataResponse{}

	ds.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_tag"
	if metadataResponse.TypeName != expected {
		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)
	}
}

func TestTagDataSourceConfigure(t *testing.T) {
	t.Parallel()

	ds := NewTagDataSource().(*TagDataSource)

	// Test with nil provider data (should not error)
	configureRequest := fwdatasource.ConfigureRequest{
		ProviderData: nil,
	}
	configureResponse := &fwdatasource.ConfigureResponse{}

	ds.Configure(context.Background(), configureRequest, configureResponse)

	if configureResponse.Diagnostics.HasError() {
		t.Errorf("Expected no error with nil provider data, got: %+v", configureResponse.Diagnostics)
	}

	// Test with correct provider data
	client := &netbox.APIClient{}
	configureRequest.ProviderData = client
	configureResponse = &fwdatasource.ConfigureResponse{}

	ds.Configure(context.Background(), configureRequest, configureResponse)

	if configureResponse.Diagnostics.HasError() {
		t.Errorf("Expected no error with correct provider data, got: %+v", configureResponse.Diagnostics)
	}
}
