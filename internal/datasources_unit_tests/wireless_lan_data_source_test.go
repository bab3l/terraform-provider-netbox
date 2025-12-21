package datasources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/datasources"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestWirelessLANDataSource(t *testing.T) {
	d := datasources.NewWirelessLANDataSource()
	if d == nil {
		t.Fatal("WirelessLAN data source should not be nil")
	}
}

func TestWirelessLANDataSourceSchema(t *testing.T) {
	ctx := context.Background()
	d := datasources.NewWirelessLANDataSource()

	req := datasource.SchemaRequest{}
	resp := &datasource.SchemaResponse{}

	d.Schema(ctx, req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("Schema() returned errors: %v", resp.Diagnostics.Errors())
	}

	schema := resp.Schema

	// Check that required attributes exist
	expectedAttrs := []string{"id", "ssid", "description", "group_id", "group_name", "status", "vlan_id", "tenant_id", "auth_type", "auth_cipher", "tags"}
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
	ssidAttr := schema.Attributes["ssid"]
	if !ssidAttr.IsOptional() {
		t.Error("'ssid' attribute should be optional for lookup")
	}
}

func TestWirelessLANDataSourceMetadata(t *testing.T) {
	ctx := context.Background()
	d := datasources.NewWirelessLANDataSource()

	req := datasource.MetadataRequest{
		ProviderTypeName: "netbox",
	}
	resp := &datasource.MetadataResponse{}

	d.Metadata(ctx, req, resp)

	expectedTypeName := "netbox_wireless_lan"
	if resp.TypeName != expectedTypeName {
		t.Errorf("Expected type name '%s', got '%s'", expectedTypeName, resp.TypeName)
	}
}
