package datasources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/datasources"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestConsolePortDataSource(t *testing.T) {
	d := datasources.NewConsolePortDataSource()
	if d == nil {
		t.Fatal("ConsolePort data source should not be nil")
	}
}

func TestConsolePortDataSourceSchema(t *testing.T) {
	ctx := context.Background()
	d := datasources.NewConsolePortDataSource()

	req := datasource.SchemaRequest{}
	resp := &datasource.SchemaResponse{}

	d.Schema(ctx, req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("Schema() returned errors: %v", resp.Diagnostics.Errors())
	}

	schema := resp.Schema

	// Check that required attributes exist
	expectedAttrs := []string{"id", "device_id", "device", "name", "label", "type", "speed", "description", "mark_connected"}
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
	deviceIDAttr := schema.Attributes["device_id"]
	if !deviceIDAttr.IsOptional() {
		t.Error("'device_id' attribute should be optional for lookup")
	}
}

func TestConsolePortDataSourceMetadata(t *testing.T) {
	ctx := context.Background()
	d := datasources.NewConsolePortDataSource()

	req := datasource.MetadataRequest{
		ProviderTypeName: "netbox",
	}
	resp := &datasource.MetadataResponse{}

	d.Metadata(ctx, req, resp)

	expectedTypeName := "netbox_console_port"
	if resp.TypeName != expectedTypeName {
		t.Errorf("Expected type name '%s', got '%s'", expectedTypeName, resp.TypeName)
	}
}
