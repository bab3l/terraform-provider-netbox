package datasources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/datasources"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestDeviceTypeDataSource(t *testing.T) {
	d := datasources.NewDeviceTypeDataSource()
	if d == nil {
		t.Fatal("DeviceType data source should not be nil")
	}
}

func TestDeviceTypeDataSourceSchema(t *testing.T) {
	ctx := context.Background()
	d := datasources.NewDeviceTypeDataSource()

	req := datasource.SchemaRequest{}
	resp := &datasource.SchemaResponse{}

	d.Schema(ctx, req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("Schema() returned errors: %v", resp.Diagnostics.Errors())
	}

	schema := resp.Schema

	// Check that required attributes exist
	expectedAttrs := []string{"id", "manufacturer", "model", "slug", "description", "tags", "custom_fields"}
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
	slugAttr := schema.Attributes["slug"]
	if !slugAttr.IsOptional() {
		t.Error("'slug' attribute should be optional for lookup")
	}
}

func TestDeviceTypeDataSourceMetadata(t *testing.T) {
	ctx := context.Background()
	d := datasources.NewDeviceTypeDataSource()

	req := datasource.MetadataRequest{
		ProviderTypeName: "netbox",
	}
	resp := &datasource.MetadataResponse{}

	d.Metadata(ctx, req, resp)

	expectedTypeName := "netbox_device_type"
	if resp.TypeName != expectedTypeName {
		t.Errorf("Expected type name '%s', got '%s'", expectedTypeName, resp.TypeName)
	}
}
