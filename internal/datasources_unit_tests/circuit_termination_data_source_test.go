package datasources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/datasources"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestCircuitTerminationDataSource(t *testing.T) {
	d := datasources.NewCircuitTerminationDataSource()
	if d == nil {
		t.Fatal("CircuitTermination data source should not be nil")
	}
}

func TestCircuitTerminationDataSourceSchema(t *testing.T) {
	ctx := context.Background()
	d := datasources.NewCircuitTerminationDataSource()

	req := datasource.SchemaRequest{}
	resp := &datasource.SchemaResponse{}

	d.Schema(ctx, req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("Schema() returned errors: %v", resp.Diagnostics.Errors())
	}

	schema := resp.Schema

	// Check that required attributes exist
	expectedAttrs := []string{"id", "circuit", "circuit_cid", "term_side", "site", "site_name", "provider_network", "port_speed", "upstream_speed", "xconnect_id", "pp_info", "description", "mark_connected", "tags"}
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
	termSideAttr := schema.Attributes["term_side"]
	if !termSideAttr.IsOptional() {
		t.Error("'term_side' attribute should be optional for lookup")
	}
}

func TestCircuitTerminationDataSourceMetadata(t *testing.T) {
	ctx := context.Background()
	d := datasources.NewCircuitTerminationDataSource()

	req := datasource.MetadataRequest{
		ProviderTypeName: "netbox",
	}
	resp := &datasource.MetadataResponse{}

	d.Metadata(ctx, req, resp)

	expectedTypeName := "netbox_circuit_termination"
	if resp.TypeName != expectedTypeName {
		t.Errorf("Expected type name '%s', got '%s'", expectedTypeName, resp.TypeName)
	}
}
