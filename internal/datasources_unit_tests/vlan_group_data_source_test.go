package datasources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/datasources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestVLANGroupDataSourceSchema(t *testing.T) {
	d := datasources.NewVLANGroupDataSource()

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

func TestVLANGroupDataSourceMetadata(t *testing.T) {
	d := datasources.NewVLANGroupDataSource()
	testutil.ValidateDataSourceMetadata(t, d, "netbox", "netbox_vlan_group")
}

func TestVLANGroupDataSourceConfigure(t *testing.T) {
	d := datasources.NewVLANGroupDataSource()
	testutil.ValidateDataSourceConfigure(t, d)
}
