package datasources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/go-netbox"
	"github.com/bab3l/terraform-provider-netbox/internal/datasources"
	fwdatasource "github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestVirtualDeviceContextDataSource(t *testing.T) {
	t.Parallel()

	d := datasources.NewVirtualDeviceContextDataSource()
	if d == nil {
		t.Fatal("Expected non-nil VirtualDeviceContext data source")
	}
}

func TestVirtualDeviceContextDataSourceSchema(t *testing.T) {
	t.Parallel()

	d := datasources.NewVirtualDeviceContextDataSource()
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
	outputAttrs := []string{"device", "device_id", "status", "identifier", "tenant", "tenant_id", "primary_ip4", "primary_ip6", "description", "comments", "tags", "custom_fields"}
	for _, attr := range outputAttrs {
		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {
			t.Errorf("Expected output attribute %s to exist in schema", attr)
		}
	}
}

func TestVirtualDeviceContextDataSourceMetadata(t *testing.T) {
	t.Parallel()

	d := datasources.NewVirtualDeviceContextDataSource()
	metadataRequest := fwdatasource.MetadataRequest{
		ProviderTypeName: "netbox",
	}
	metadataResponse := &fwdatasource.MetadataResponse{}

	d.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_virtual_device_context"
	if metadataResponse.TypeName != expected {
		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)
	}
}

func TestVirtualDeviceContextDataSourceConfigure(t *testing.T) {
	t.Parallel()

	d := datasources.NewVirtualDeviceContextDataSource().(*datasources.VirtualDeviceContextDataSource)

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
