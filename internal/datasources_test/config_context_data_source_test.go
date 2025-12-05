package datasources_test

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/datasources"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestConfigContextDataSource_Metadata(t *testing.T) {
	d := datasources.NewConfigContextDataSource()

	req := datasource.MetadataRequest{
		ProviderTypeName: "netbox",
	}
	resp := &datasource.MetadataResponse{}

	d.Metadata(context.Background(), req, resp)

	if resp.TypeName != "netbox_config_context" {
		t.Errorf("Expected type name 'netbox_config_context', got '%s'", resp.TypeName)
	}
}

func TestConfigContextDataSource_Schema(t *testing.T) {
	d := datasources.NewConfigContextDataSource()

	req := datasource.SchemaRequest{}
	resp := &datasource.SchemaResponse{}

	d.Schema(context.Background(), req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("Schema returned errors: %v", resp.Diagnostics)
	}

	// Verify lookup attributes exist
	lookupAttrs := []string{"id", "name"}
	for _, attr := range lookupAttrs {
		if _, ok := resp.Schema.Attributes[attr]; !ok {
			t.Errorf("Expected lookup attribute '%s' to exist in schema", attr)
		}
	}

	// Verify computed attributes exist
	computedAttrs := []string{
		"description", "weight", "is_active", "data",
		"regions", "site_groups", "sites", "locations",
		"device_types", "roles", "platforms",
		"cluster_types", "cluster_groups", "clusters",
		"tenant_groups", "tenants", "tags",
	}
	for _, attr := range computedAttrs {
		if _, ok := resp.Schema.Attributes[attr]; !ok {
			t.Errorf("Expected computed attribute '%s' to exist in schema", attr)
		}
	}
}

func TestConfigContextDataSource_SchemaDescription(t *testing.T) {
	d := datasources.NewConfigContextDataSource()

	req := datasource.SchemaRequest{}
	resp := &datasource.SchemaResponse{}

	d.Schema(context.Background(), req, resp)

	if resp.Schema.MarkdownDescription == "" {
		t.Error("Expected schema to have a markdown description")
	}
}

func TestConfigContextDataSource_Configure(t *testing.T) {
	d := datasources.NewConfigContextDataSource()

	// Verify the data source implements the configurable interface
	configurable, ok := d.(datasource.DataSourceWithConfigure)
	if !ok {
		t.Skip("Data source does not implement DataSourceWithConfigure")
	}

	// Test with nil provider data - should not error
	req := datasource.ConfigureRequest{
		ProviderData: nil,
	}
	resp := &datasource.ConfigureResponse{}

	configurable.Configure(context.Background(), req, resp)

	if resp.Diagnostics.HasError() {
		t.Errorf("Configure with nil provider data should not error: %v", resp.Diagnostics)
	}
}
