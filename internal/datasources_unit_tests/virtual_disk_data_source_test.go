package datasources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/datasources"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestVirtualDiskDataSource(t *testing.T) {
	t.Parallel()

	d := datasources.NewVirtualDiskDataSource()
	if d == nil {
		t.Fatal("Expected non-nil VirtualDisk data source")
	}
}

func TestVirtualDiskDataSourceSchema(t *testing.T) {
	t.Parallel()

	d := datasources.NewVirtualDiskDataSource()
	schemaRequest := datasource.SchemaRequest{}
	schemaResponse := &datasource.SchemaResponse{}

	d.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)
	}

	if schemaResponse.Schema.Attributes == nil {
		t.Fatal("Expected schema to have attributes")
	}

	// Check for lookup attributes (optional)
	lookupAttrs := []string{"id", "name", "virtual_machine"}
	for _, attr := range lookupAttrs {
		a, exists := schemaResponse.Schema.Attributes[attr]
		if !exists {
			t.Errorf("Expected lookup attribute %s to exist in schema", attr)
			continue
		}
		if a.IsRequired() {
			t.Errorf("Expected lookup attribute %s to be optional, not required", attr)
		}
	}

	// Check computed attributes
	computedAttrs := []string{"virtual_machine_name", "size", "description", "tags"}
	for _, attr := range computedAttrs {
		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {
			t.Errorf("Expected computed attribute %s to exist in schema", attr)
		}
	}
}

func TestVirtualDiskDataSourceMetadata(t *testing.T) {
	t.Parallel()

	d := datasources.NewVirtualDiskDataSource()
	metadataRequest := datasource.MetadataRequest{
		ProviderTypeName: "netbox",
	}
	metadataResponse := &datasource.MetadataResponse{}

	d.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_virtual_disk"
	if metadataResponse.TypeName != expected {
		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)
	}
}
