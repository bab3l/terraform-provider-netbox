package datasources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/datasources"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestIPAddressDataSource(t *testing.T) {
	d := datasources.NewIPAddressDataSource()
	if d == nil {
		t.Fatal("IPAddress data source should not be nil")
	}
}

func TestIPAddressDataSourceSchema(t *testing.T) {
	ctx := context.Background()
	d := datasources.NewIPAddressDataSource()

	req := datasource.SchemaRequest{}
	resp := &datasource.SchemaResponse{}

	d.Schema(ctx, req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("Schema() returned errors: %v", resp.Diagnostics.Errors())
	}

	schema := resp.Schema

	// Check that required attributes exist
	expectedAttrs := []string{"id", "address", "status", "vrf", "tenant", "description", "tags"}
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
	addressAttr := schema.Attributes["address"]
	if !addressAttr.IsOptional() {
		t.Error("'address' attribute should be optional for lookup")
	}
}

func TestIPAddressDataSourceMetadata(t *testing.T) {
	ctx := context.Background()
	d := datasources.NewIPAddressDataSource()

	req := datasource.MetadataRequest{
		ProviderTypeName: "netbox",
	}
	resp := &datasource.MetadataResponse{}

	d.Metadata(ctx, req, resp)

	expectedTypeName := "netbox_ip_address"
	if resp.TypeName != expectedTypeName {
		t.Errorf("Expected type name '%s', got '%s'", expectedTypeName, resp.TypeName)
	}
}
