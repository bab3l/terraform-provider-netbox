package datasources

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestCircuitDataSource(t *testing.T) {
	d := NewCircuitDataSource()
	if d == nil {
		t.Fatal("Circuit data source should not be nil")
	}
}

func TestCircuitDataSourceSchema(t *testing.T) {
	ctx := context.Background()
	d := NewCircuitDataSource()

	req := datasource.SchemaRequest{}
	resp := &datasource.SchemaResponse{}

	d.Schema(ctx, req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("Schema() returned errors: %v", resp.Diagnostics.Errors())
	}

	schema := resp.Schema

	// Check that required attributes exist
	expectedAttrs := []string{"id", "cid", "circuit_provider", "type", "status", "tenant", "description", "tags", "custom_fields"}
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
	cidAttr := schema.Attributes["cid"]
	if !cidAttr.IsOptional() {
		t.Error("'cid' attribute should be optional for lookup")
	}
}

func TestCircuitDataSourceMetadata(t *testing.T) {
	ctx := context.Background()
	d := NewCircuitDataSource()

	req := datasource.MetadataRequest{
		ProviderTypeName: "netbox",
	}
	resp := &datasource.MetadataResponse{}

	d.Metadata(ctx, req, resp)

	expectedTypeName := "netbox_circuit"
	if resp.TypeName != expectedTypeName {
		t.Errorf("Expected type name '%s', got '%s'", expectedTypeName, resp.TypeName)
	}
}
