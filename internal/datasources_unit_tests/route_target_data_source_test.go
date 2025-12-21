package datasources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/datasources"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestRouteTargetDataSource(t *testing.T) {
	t.Parallel()

	d := datasources.NewRouteTargetDataSource()
	if d == nil {
		t.Fatal("Expected non-nil RouteTarget data source")
	}
}

func TestRouteTargetDataSourceSchema(t *testing.T) {
	t.Parallel()

	d := datasources.NewRouteTargetDataSource()
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
	lookupAttrs := []string{"id", "name"}
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
	computedAttrs := []string{"tenant", "tenant_name", "description", "comments", "tags"}
	for _, attr := range computedAttrs {
		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {
			t.Errorf("Expected computed attribute %s to exist in schema", attr)
		}
	}
}

func TestRouteTargetDataSourceMetadata(t *testing.T) {
	t.Parallel()

	d := datasources.NewRouteTargetDataSource()
	metadataRequest := datasource.MetadataRequest{
		ProviderTypeName: "netbox",
	}
	metadataResponse := &datasource.MetadataResponse{}

	d.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_route_target"
	if metadataResponse.TypeName != expected {
		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)
	}
}
