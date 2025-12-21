package datasources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/datasources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestTenantDataSourceSchema(t *testing.T) {

	t.Parallel()
	d := datasources.NewTenantDataSource()

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

func TestTenantDataSourceMetadata(t *testing.T) {

	t.Parallel()
	d := datasources.NewTenantDataSource()
	testutil.ValidateDataSourceMetadata(t, d, "netbox", "netbox_tenant")
}

func TestTenantDataSourceConfigure(t *testing.T) {

	t.Parallel()
	d := datasources.NewTenantDataSource()
	testutil.ValidateDataSourceConfigure(t, d)
}
