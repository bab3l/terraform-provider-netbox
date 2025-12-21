package datasources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/datasources"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestSiteGroupDataSource(t *testing.T) {
	// Test that the site group data source can be instantiated
	d := datasources.NewSiteGroupDataSource()
	if d == nil {
		t.Fatal("Site group data source should not be nil")
	}
}

func TestSiteGroupDataSourceSchema(t *testing.T) {
	ctx := context.Background()
	d := datasources.NewSiteGroupDataSource()

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
	if _, ok := schema.Attributes["parent"]; !ok {
		t.Error("Schema should have 'parent' attribute")
	}
	if _, ok := schema.Attributes["description"]; !ok {
		t.Error("Schema should have 'description' attribute")
	}
	if _, ok := schema.Attributes["tags"]; !ok {
		t.Error("Schema should have 'tags' attribute")
	}
	if _, ok := schema.Attributes["custom_fields"]; !ok {
		t.Error("Schema should have 'custom_fields' attribute")
	}
}

func TestSiteGroupDataSourceMetadata(t *testing.T) {
	ctx := context.Background()
	d := datasources.NewSiteGroupDataSource()

	req := datasource.MetadataRequest{
		ProviderTypeName: "netbox",
	}
	resp := &datasource.MetadataResponse{}

	d.Metadata(ctx, req, resp)

	expectedTypeName := "netbox_site_group"
	if resp.TypeName != expectedTypeName {
		t.Errorf("Expected type name '%s', got '%s'", expectedTypeName, resp.TypeName)
	}
}
