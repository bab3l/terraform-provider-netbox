package datasources

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestPowerPortDataSource(t *testing.T) {
	d := NewPowerPortDataSource()
	if d == nil {
		t.Fatal("PowerPort data source should not be nil")
	}
}

func TestPowerPortDataSourceSchema(t *testing.T) {
	ctx := context.Background()
	d := NewPowerPortDataSource()

	req := datasource.SchemaRequest{}
	resp := &datasource.SchemaResponse{}

	d.Schema(ctx, req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("Schema() returned errors: %v", resp.Diagnostics.Errors())
	}

	schema := resp.Schema

	// Check that required attributes exist
	expectedAttrs := []string{"id", "device_id", "device", "name", "label", "type", "maximum_draw", "allocated_draw", "description", "mark_connected"}
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
}

func TestPowerPortDataSourceMetadata(t *testing.T) {
	ctx := context.Background()
	d := NewPowerPortDataSource()

	req := datasource.MetadataRequest{
		ProviderTypeName: "netbox",
	}
	resp := &datasource.MetadataResponse{}

	d.Metadata(ctx, req, resp)

	expectedTypeName := "netbox_power_port"
	if resp.TypeName != expectedTypeName {
		t.Errorf("Expected type name '%s', got '%s'", expectedTypeName, resp.TypeName)
	}
}
