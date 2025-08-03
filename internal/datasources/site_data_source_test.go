package datasources

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestSiteDataSource(t *testing.T) {
	// Test that the site data source can be instantiated
	d := NewSiteDataSource()
	if d == nil {
		t.Fatal("Site data source should not be nil")
	}
}

func TestSiteDataSourceSchema(t *testing.T) {
	ctx := context.Background()
	d := NewSiteDataSource()

	req := datasource.SchemaRequest{}
	resp := &datasource.SchemaResponse{}

	d.Schema(ctx, req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("Schema() returned errors: %v", resp.Diagnostics.Errors())
	}

	// Verify that the schema has the expected attributes
	schema := resp.Schema

	// Check that required attributes exist
	if _, ok := schema.Attributes["id"]; !ok {
		t.Error("Schema should have 'id' attribute")
	}
	if _, ok := schema.Attributes["name"]; !ok {
		t.Error("Schema should have 'name' attribute")
	}
	if _, ok := schema.Attributes["slug"]; !ok {
		t.Error("Schema should have 'slug' attribute")
	}
	if _, ok := schema.Attributes["status"]; !ok {
		t.Error("Schema should have 'status' attribute")
	}
}

func TestSiteDataSourceMetadata(t *testing.T) {
	ctx := context.Background()
	d := NewSiteDataSource()

	req := datasource.MetadataRequest{
		ProviderTypeName: "netbox",
	}
	resp := &datasource.MetadataResponse{}

	d.Metadata(ctx, req, resp)

	expectedTypeName := "netbox_site"
	if resp.TypeName != expectedTypeName {
		t.Errorf("Expected type name '%s', got '%s'", expectedTypeName, resp.TypeName)
	}
}
