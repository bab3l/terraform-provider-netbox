package datasources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/datasources"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestProviderAccountDataSource(t *testing.T) {
	d := datasources.NewProviderAccountDataSource()
	if d == nil {
		t.Fatal("ProviderAccount data source should not be nil")
	}
}

func TestProviderAccountDataSourceSchema(t *testing.T) {
	ctx := context.Background()
	d := datasources.NewProviderAccountDataSource()

	req := datasource.SchemaRequest{}
	resp := &datasource.SchemaResponse{}

	d.Schema(ctx, req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("Schema() returned errors: %v", resp.Diagnostics.Errors())
	}

	schema := resp.Schema

	// Check that required attributes exist
	expectedAttrs := []string{"id", "circuit_provider", "provider_name", "name", "account", "description", "comments", "tags"}
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
	accountAttr := schema.Attributes["account"]
	if !accountAttr.IsOptional() {
		t.Error("'account' attribute should be optional for lookup")
	}
}

func TestProviderAccountDataSourceMetadata(t *testing.T) {
	ctx := context.Background()
	d := datasources.NewProviderAccountDataSource()

	req := datasource.MetadataRequest{
		ProviderTypeName: "netbox",
	}
	resp := &datasource.MetadataResponse{}

	d.Metadata(ctx, req, resp)

	expectedTypeName := "netbox_provider_account"
	if resp.TypeName != expectedTypeName {
		t.Errorf("Expected type name '%s', got '%s'", expectedTypeName, resp.TypeName)
	}
}
