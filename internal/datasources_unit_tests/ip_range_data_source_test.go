package datasources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/datasources"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestIPRangeDataSource(t *testing.T) {
	d := datasources.NewIPRangeDataSource()
	if d == nil {
		t.Fatal("IPRange data source should not be nil")
	}
}

func TestIPRangeDataSourceSchema(t *testing.T) {
	ctx := context.Background()
	d := datasources.NewIPRangeDataSource()

	req := datasource.SchemaRequest{}
	resp := &datasource.SchemaResponse{}

	d.Schema(ctx, req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("Schema() returned errors: %v", resp.Diagnostics.Errors())
	}

	schema := resp.Schema

	// Check that required attributes exist
	expectedAttrs := []string{"id", "start_address", "end_address", "size", "vrf", "tenant", "status", "role", "description", "tags"}
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
	startAddrAttr := schema.Attributes["start_address"]
	if !startAddrAttr.IsOptional() {
		t.Error("'start_address' attribute should be optional for lookup")
	}
}

func TestIPRangeDataSourceMetadata(t *testing.T) {
	ctx := context.Background()
	d := datasources.NewIPRangeDataSource()

	req := datasource.MetadataRequest{
		ProviderTypeName: "netbox",
	}
	resp := &datasource.MetadataResponse{}

	d.Metadata(ctx, req, resp)

	expectedTypeName := "netbox_ip_range"
	if resp.TypeName != expectedTypeName {
		t.Errorf("Expected type name '%s', got '%s'", expectedTypeName, resp.TypeName)
	}
}
