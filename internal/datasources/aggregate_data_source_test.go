package datasources

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestAggregateDataSource(t *testing.T) {
	d := NewAggregateDataSource()
	if d == nil {
		t.Fatal("Aggregate data source should not be nil")
	}
}

func TestAggregateDataSourceSchema(t *testing.T) {
	ctx := context.Background()
	d := NewAggregateDataSource()

	req := datasource.SchemaRequest{}
	resp := &datasource.SchemaResponse{}

	d.Schema(ctx, req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("Schema() returned errors: %v", resp.Diagnostics.Errors())
	}

	schema := resp.Schema

	// Check that required attributes exist
	expectedAttrs := []string{"id", "prefix", "rir", "rir_name", "tenant", "tenant_name", "date_added", "description", "comments", "tags"}
	for _, attr := range expectedAttrs {
		if _, ok := schema.Attributes[attr]; !ok {
			t.Errorf("Schema should have '%s' attribute", attr)
		}
	}

	// Verify that lookup fields are optional
	idAttr := schema.Attributes["id"]
	if !idAttr.IsOptional() {
		t.Error("'id' attribute should be optional for lookup")
	}
	prefixAttr := schema.Attributes["prefix"]
	if !prefixAttr.IsOptional() {
		t.Error("'prefix' attribute should be optional for lookup")
	}
}

func TestAggregateDataSourceMetadata(t *testing.T) {
	ctx := context.Background()
	d := NewAggregateDataSource()

	req := datasource.MetadataRequest{
		ProviderTypeName: "netbox",
	}
	resp := &datasource.MetadataResponse{}

	d.Metadata(ctx, req, resp)

	expectedTypeName := "netbox_aggregate"
	if resp.TypeName != expectedTypeName {
		t.Errorf("Expected type name '%s', got '%s'", expectedTypeName, resp.TypeName)
	}
}
