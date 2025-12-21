package datasources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/datasources"
	fwdatasource "github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestASNRangeDataSource(t *testing.T) {
	t.Parallel()

	d := datasources.NewASNRangeDataSource()
	if d == nil {
		t.Fatal("Expected non-nil ASNRange data source")
	}
}

func TestASNRangeDataSourceSchema(t *testing.T) {
	t.Parallel()

	d := datasources.NewASNRangeDataSource()
	schemaRequest := fwdatasource.SchemaRequest{}
	schemaResponse := &fwdatasource.SchemaResponse{}

	d.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)
	}

	if schemaResponse.Schema.Attributes == nil {
		t.Fatal("Expected schema to have attributes")
	}

	// Required lookup attributes (at least one must be provided, but all are Optional in schema)
	lookupAttrs := []string{"id", "name", "slug"}
	for _, attr := range lookupAttrs {
		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {
			t.Errorf("Expected lookup attribute %s to exist in schema", attr)
		}
	}

	// Output-only attributes
	outputAttrs := []string{"rir", "start", "end", "asn_count", "tenant", "description", "tags"}
	for _, attr := range outputAttrs {
		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {
			t.Errorf("Expected output attribute %s to exist in schema", attr)
		}
	}
}
