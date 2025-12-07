package datasources

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestRoleDataSource(t *testing.T) {
	d := NewRoleDataSource()
	if d == nil {
		t.Fatal("Role data source should not be nil")
	}
}

func TestRoleDataSourceSchema(t *testing.T) {
	ctx := context.Background()
	d := NewRoleDataSource()

	req := datasource.SchemaRequest{}
	resp := &datasource.SchemaResponse{}

	d.Schema(ctx, req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("Schema() returned errors: %v", resp.Diagnostics.Errors())
	}

	schema := resp.Schema

	// Check that required attributes exist
	expectedAttrs := []string{"id", "name", "slug", "weight", "description", "prefix_count", "vlan_count", "tags", "custom_fields"}
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
	nameAttr := schema.Attributes["name"]
	if !nameAttr.IsOptional() {
		t.Error("'name' attribute should be optional for lookup")
	}
	slugAttr := schema.Attributes["slug"]
	if !slugAttr.IsOptional() {
		t.Error("'slug' attribute should be optional for lookup")
	}
}

func TestRoleDataSourceMetadata(t *testing.T) {
	ctx := context.Background()
	d := NewRoleDataSource()

	req := datasource.MetadataRequest{
		ProviderTypeName: "netbox",
	}
	resp := &datasource.MetadataResponse{}

	d.Metadata(ctx, req, resp)

	expectedTypeName := "netbox_role"
	if resp.TypeName != expectedTypeName {
		t.Errorf("Expected type name '%s', got '%s'", expectedTypeName, resp.TypeName)
	}
}
