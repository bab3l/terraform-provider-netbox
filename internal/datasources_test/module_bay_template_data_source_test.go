package datasources_test

import (
	"context"
	"testing"

	"github.com/bab3l/go-netbox"
	"github.com/bab3l/terraform-provider-netbox/internal/datasources"
	fwdatasource "github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestModuleBayTemplateDataSource(t *testing.T) {
	t.Parallel()

	d := datasources.NewModuleBayTemplateDataSource()
	if d == nil {
		t.Fatal("Expected non-nil ModuleBayTemplate data source")
	}
}

func TestModuleBayTemplateDataSourceSchema(t *testing.T) {
	t.Parallel()

	d := datasources.NewModuleBayTemplateDataSource()
	schemaRequest := fwdatasource.SchemaRequest{}
	schemaResponse := &fwdatasource.SchemaResponse{}

	d.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)
	}

	if schemaResponse.Schema.Attributes == nil {
		t.Fatal("Expected schema to have attributes")
	}

	// Lookup attributes
	lookupAttrs := []string{"id", "name"}
	for _, attr := range lookupAttrs {
		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {
			t.Errorf("Expected lookup attribute %s to exist in schema", attr)
		}
	}

	// Output attributes
	outputAttrs := []string{"device_type", "device_type_id", "module_type", "module_type_id", "label", "position", "description"}
	for _, attr := range outputAttrs {
		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {
			t.Errorf("Expected output attribute %s to exist in schema", attr)
		}
	}
}

func TestModuleBayTemplateDataSourceMetadata(t *testing.T) {
	t.Parallel()

	d := datasources.NewModuleBayTemplateDataSource()
	metadataRequest := fwdatasource.MetadataRequest{
		ProviderTypeName: "netbox",
	}
	metadataResponse := &fwdatasource.MetadataResponse{}

	d.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_module_bay_template"
	if metadataResponse.TypeName != expected {
		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)
	}
}

func TestModuleBayTemplateDataSourceConfigure(t *testing.T) {
	t.Parallel()

	d := datasources.NewModuleBayTemplateDataSource().(*datasources.ModuleBayTemplateDataSource)

	configureRequest := fwdatasource.ConfigureRequest{
		ProviderData: nil,
	}
	configureResponse := &fwdatasource.ConfigureResponse{}

	d.Configure(context.Background(), configureRequest, configureResponse)

	if configureResponse.Diagnostics.HasError() {
		t.Errorf("Expected no error with nil provider data, got: %+v", configureResponse.Diagnostics)
	}

	client := &netbox.APIClient{}
	configureRequest.ProviderData = client
	configureResponse = &fwdatasource.ConfigureResponse{}

	d.Configure(context.Background(), configureRequest, configureResponse)

	if configureResponse.Diagnostics.HasError() {
		t.Errorf("Expected no error with correct provider data, got: %+v", configureResponse.Diagnostics)
	}
}
