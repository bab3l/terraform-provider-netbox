package datasources

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestConsoleServerPortDataSource(t *testing.T) {
	d := NewConsoleServerPortDataSource()
	if d == nil {
		t.Fatal("ConsoleServerPort data source should not be nil")
	}
}

func TestConsoleServerPortDataSourceSchema(t *testing.T) {
	ctx := context.Background()
	d := NewConsoleServerPortDataSource()

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

func TestConsoleServerPortDataSourceMetadata(t *testing.T) {
	ctx := context.Background()
	d := NewConsoleServerPortDataSource()

	req := datasource.MetadataRequest{
		ProviderTypeName: "netbox",
	}
	resp := &datasource.MetadataResponse{}

	d.Metadata(ctx, req, resp)

	expectedTypeName := "netbox_console_server_port"
	if resp.TypeName != expectedTypeName {
		t.Errorf("Expected type name '%s', got '%s'", expectedTypeName, resp.TypeName)
	}
}
