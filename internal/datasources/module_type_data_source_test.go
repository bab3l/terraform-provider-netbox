package datasources

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestModuleTypeDataSource(t *testing.T) {
	d := NewModuleTypeDataSource()
	if d == nil {
		t.Fatal("ModuleType data source should not be nil")
	}
}

func TestModuleTypeDataSourceSchema(t *testing.T) {
	ctx := context.Background()
	d := NewModuleTypeDataSource()

	req := datasource.SchemaRequest{}
	resp := &datasource.SchemaResponse{}

	d.Schema(ctx, req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("Schema() returned errors: %v", resp.Diagnostics.Errors())
	}

	schema := resp.Schema

	// Check that required attributes exist
	expectedAttrs := []string{"id", "model", "manufacturer_id", "manufacturer", "part_number", "airflow", "weight", "weight_unit", "description", "comments"}
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
	modelAttr := schema.Attributes["model"]
	if !modelAttr.IsOptional() {
		t.Error("'model' attribute should be optional for lookup")
	}
}

func TestModuleTypeDataSourceMetadata(t *testing.T) {
	ctx := context.Background()
	d := NewModuleTypeDataSource()

	req := datasource.MetadataRequest{
		ProviderTypeName: "netbox",
	}
	resp := &datasource.MetadataResponse{}

	d.Metadata(ctx, req, resp)

	expectedTypeName := "netbox_module_type"
	if resp.TypeName != expectedTypeName {
		t.Errorf("Expected type name '%s', got '%s'", expectedTypeName, resp.TypeName)
	}
}
