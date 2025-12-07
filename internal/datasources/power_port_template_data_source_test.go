package datasources

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestPowerPortTemplateDataSource(t *testing.T) {
	d := NewPowerPortTemplateDataSource()
	if d == nil {
		t.Fatal("PowerPortTemplate data source should not be nil")
	}
}

func TestPowerPortTemplateDataSourceSchema(t *testing.T) {
	ctx := context.Background()
	d := NewPowerPortTemplateDataSource()

	req := datasource.SchemaRequest{}
	resp := &datasource.SchemaResponse{}

	d.Schema(ctx, req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("Schema() returned errors: %v", resp.Diagnostics.Errors())
	}

	schema := resp.Schema

	// Check that required attributes exist
	expectedAttrs := []string{"id", "device_type", "module_type", "name", "label", "type", "maximum_draw", "allocated_draw", "description"}
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

func TestPowerPortTemplateDataSourceMetadata(t *testing.T) {
	ctx := context.Background()
	d := NewPowerPortTemplateDataSource()

	req := datasource.MetadataRequest{
		ProviderTypeName: "netbox",
	}
	resp := &datasource.MetadataResponse{}

	d.Metadata(ctx, req, resp)

	expectedTypeName := "netbox_power_port_template"
	if resp.TypeName != expectedTypeName {
		t.Errorf("Expected type name '%s', got '%s'", expectedTypeName, resp.TypeName)
	}
}
