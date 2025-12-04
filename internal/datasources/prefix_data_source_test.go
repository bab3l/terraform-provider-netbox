package datasources

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestPrefixDataSource(t *testing.T) {
	d := NewPrefixDataSource()
	if d == nil {
		t.Fatal("Prefix data source should not be nil")
	}
}

func TestPrefixDataSourceSchema(t *testing.T) {
	ctx := context.Background()
	d := NewPrefixDataSource()

	req := datasource.SchemaRequest{}
	resp := &datasource.SchemaResponse{}

	d.Schema(ctx, req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("Schema() returned errors: %v", resp.Diagnostics.Errors())
	}

	schema := resp.Schema

	// Check that required attributes exist
	expectedAttrs := []string{"id", "prefix", "status", "vrf", "site", "tenant", "description", "tags"}
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

func TestPrefixDataSourceMetadata(t *testing.T) {
	ctx := context.Background()
	d := NewPrefixDataSource()

	req := datasource.MetadataRequest{
		ProviderTypeName: "netbox",
	}
	resp := &datasource.MetadataResponse{}

	d.Metadata(ctx, req, resp)

	expectedTypeName := "netbox_prefix"
	if resp.TypeName != expectedTypeName {
		t.Errorf("Expected type name '%s', got '%s'", expectedTypeName, resp.TypeName)
	}
}
