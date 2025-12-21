package datasources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/datasources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestInventoryItemRoleDataSourceSchema(t *testing.T) {
	d := datasources.NewInventoryItemRoleDataSource()

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

func TestInventoryItemRoleDataSourceMetadata(t *testing.T) {
	d := datasources.NewInventoryItemRoleDataSource()
	testutil.ValidateDataSourceMetadata(t, d, "netbox", "netbox_inventory_item_role")
}

func TestInventoryItemRoleDataSourceConfigure(t *testing.T) {
	d := datasources.NewInventoryItemRoleDataSource()
	testutil.ValidateDataSourceConfigure(t, d)
}
