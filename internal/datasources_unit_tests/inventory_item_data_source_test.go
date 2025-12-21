package datasources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/datasources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestInventoryItemDataSourceSchema(t *testing.T) {
	d := datasources.NewInventoryItemDataSource()

	req := datasource.SchemaRequest{}
	resp := &datasource.SchemaResponse{}

	d.Schema(context.Background(), req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("Schema returned errors: %v", resp.Diagnostics)
	}

	testutil.ValidateDataSourceSchema(t, resp.Schema.Attributes, testutil.DataSourceValidation{
		LookupAttrs:   []string{},
		ComputedAttrs: []string{},
	})
}

func TestInventoryItemDataSourceMetadata(t *testing.T) {
	d := datasources.NewInventoryItemDataSource()
	testutil.ValidateDataSourceMetadata(t, d, "netbox", "netbox_inventory_item")
}

func TestInventoryItemDataSourceConfigure(t *testing.T) {
	d := datasources.NewInventoryItemDataSource()
	testutil.ValidateDataSourceConfigure(t, d)
}
