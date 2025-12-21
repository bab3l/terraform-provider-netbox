package datasources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/datasources"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestCableDataSource(t *testing.T) {
	d := datasources.NewCableDataSource()
	if d == nil {
		t.Fatal("Cable data source should not be nil")
	}
}

func TestCableDataSourceSchema(t *testing.T) {
	ctx := context.Background()
	d := datasources.NewCableDataSource()

	req := datasource.SchemaRequest{}
	resp := &datasource.SchemaResponse{}

	d.Schema(ctx, req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("Schema() returned errors: %v", resp.Diagnostics.Errors())
	}

	schema := resp.Schema

	// Check that expected attributes exist
	expectedAttrs := []string{
		"id", "a_terminations", "b_terminations",
		"type", "status", "tenant", "tenant_id",
		"label", "color", "length", "length_unit",
		"description", "comments", "tags",
	}
	for _, attr := range expectedAttrs {
		if _, ok := schema.Attributes[attr]; !ok {
			t.Errorf("Schema should have '%s' attribute", attr)
		}
	}

	// Verify that id is optional for lookup
	idAttr := schema.Attributes["id"]
	if !idAttr.IsOptional() {
		t.Error("'id' attribute should be optional for lookup")
	}
}

func TestCableDataSourceMetadata(t *testing.T) {
	ctx := context.Background()
	d := datasources.NewCableDataSource()

	req := datasource.MetadataRequest{
		ProviderTypeName: "netbox",
	}
	resp := &datasource.MetadataResponse{}

	d.Metadata(ctx, req, resp)

	expected := "netbox_cable"
	if resp.TypeName != expected {
		t.Errorf("Expected type name %s, got %s", expected, resp.TypeName)
	}
}
