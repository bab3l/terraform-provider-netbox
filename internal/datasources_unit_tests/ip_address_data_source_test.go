package datasources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/datasources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestIPAddressDataSourceSchema(t *testing.T) {
	d := datasources.NewIPAddressDataSource()

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

func TestIPAddressDataSourceMetadata(t *testing.T) {
	d := datasources.NewIPAddressDataSource()
	testutil.ValidateDataSourceMetadata(t, d, "netbox", "netbox_ip_address")
}

func TestIPAddressDataSourceConfigure(t *testing.T) {
	d := datasources.NewIPAddressDataSource()
	testutil.ValidateDataSourceConfigure(t, d)
}
