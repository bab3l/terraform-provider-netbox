package datasources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/datasources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestInventoryItemTemplateDataSourceSchema(t *testing.T) {
	t.Parallel()

	d := datasources.NewInventoryItemTemplateDataSource()
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

func TestInventoryItemTemplateDataSourceMetadata(t *testing.T) {
	t.Parallel()

	d := datasources.NewInventoryItemTemplateDataSource()
	testutil.ValidateDataSourceMetadata(t, d, "netbox", "netbox_inventory_item_template")
}

func TestInventoryItemTemplateDataSourceConfigure(t *testing.T) {
	t.Parallel()

	d := datasources.NewInventoryItemTemplateDataSource()
	testutil.ValidateDataSourceConfigure(t, d)
}
