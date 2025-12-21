package datasources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/go-netbox"
	"github.com/bab3l/terraform-provider-netbox/internal/datasources"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestTenantGroupDataSource(t *testing.T) {
	t.Parallel()

	d := datasources.NewTenantGroupDataSource()
	if d == nil {
		t.Fatal("Expected non-nil tenant group data source")
	}
}

func TestTenantGroupDataSourceSchema(t *testing.T) {
	t.Parallel()

	d := datasources.NewTenantGroupDataSource()
	schemaRequest := datasource.SchemaRequest{}
	schemaResponse := &datasource.SchemaResponse{}

	d.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)
	}

	if schemaResponse.Schema.Attributes == nil {
		t.Fatal("Expected schema to have attributes")
	}

	// Check that key attributes exist
	expectedAttrs := []string{"id", "name", "slug", "parent", "description", "tags", "custom_fields"}
	for _, attr := range expectedAttrs {
		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {
			t.Errorf("Expected attribute %s to exist in schema", attr)
		}
	}
}

func TestTenantGroupDataSourceMetadata(t *testing.T) {
	t.Parallel()

	d := datasources.NewTenantGroupDataSource()
	metadataRequest := datasource.MetadataRequest{
		ProviderTypeName: "netbox",
	}
	metadataResponse := &datasource.MetadataResponse{}

	d.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_tenant_group"
	if metadataResponse.TypeName != expected {
		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)
	}
}

func TestTenantGroupDataSourceConfigure(t *testing.T) {
	t.Parallel()

	d := datasources.NewTenantGroupDataSource().(*datasources.TenantGroupDataSource)

	// Test with nil provider data
	configureRequest := datasource.ConfigureRequest{
		ProviderData: nil,
	}
	configureResponse := &datasource.ConfigureResponse{}

	d.Configure(context.Background(), configureRequest, configureResponse)

	if configureResponse.Diagnostics.HasError() {
		t.Errorf("Expected no error with nil provider data, got: %+v", configureResponse.Diagnostics)
	}

	// Test with correct provider data
	client := &netbox.APIClient{}
	configureRequest.ProviderData = client
	configureResponse = &datasource.ConfigureResponse{}

	d.Configure(context.Background(), configureRequest, configureResponse)

	if configureResponse.Diagnostics.HasError() {
		t.Errorf("Expected no error with correct provider data, got: %+v", configureResponse.Diagnostics)
	}

	// Note: Cannot access unexported d.client field from a different package
	// The Configure method was successful if no errors are present above

	// Test with incorrect provider data
	configureRequest.ProviderData = "invalid"
	configureResponse = &datasource.ConfigureResponse{}

	d.Configure(context.Background(), configureRequest, configureResponse)

	if !configureResponse.Diagnostics.HasError() {
		t.Error("Expected error with incorrect provider data")
	}
}
