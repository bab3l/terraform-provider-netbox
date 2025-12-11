// Package datasources_test provides unit tests for the NetBox Terraform provider data sources.
package datasources_test

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/datasources"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestWirelessLinkDataSourceMetadata(t *testing.T) {
	d := datasources.NewWirelessLinkDataSource()
	req := datasource.MetadataRequest{ProviderTypeName: "netbox"}
	resp := &datasource.MetadataResponse{}
	d.Metadata(context.Background(), req, resp)

	expected := "netbox_wireless_link"
	if resp.TypeName != expected {
		t.Errorf("Expected type name %q, got %q", expected, resp.TypeName)
	}
}

func TestWirelessLinkDataSourceSchema(t *testing.T) {
	d := datasources.NewWirelessLinkDataSource()
	req := datasource.SchemaRequest{}
	resp := &datasource.SchemaResponse{}
	d.Schema(context.Background(), req, resp)

	// Check that response has no diagnostics errors
	if resp.Diagnostics.HasError() {
		t.Errorf("Schema returned errors: %v", resp.Diagnostics)
	}

	// Verify expected attributes exist
	expectedAttrs := []string{
		"id", "interface_a", "interface_b", "ssid", "status", "tenant", "tenant_id",
		"auth_type", "auth_cipher", "distance", "distance_unit", "description",
		"comments", "tags", "custom_fields",
	}
	for _, attr := range expectedAttrs {
		if _, ok := resp.Schema.Attributes[attr]; !ok {
			t.Errorf("Expected attribute %q not found in schema", attr)
		}
	}
}
