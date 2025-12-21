package datasources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/datasources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestPowerPanelDataSourceSchema(t *testing.T) {
	d := datasources.NewPowerPanelDataSource()

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

func TestPowerPanelDataSourceMetadata(t *testing.T) {
	d := datasources.NewPowerPanelDataSource()
	testutil.ValidateDataSourceMetadata(t, d, "netbox", "netbox_power_panel")
}

func TestPowerPanelDataSourceConfigure(t *testing.T) {
	d := datasources.NewPowerPanelDataSource()
	testutil.ValidateDataSourceConfigure(t, d)
}
