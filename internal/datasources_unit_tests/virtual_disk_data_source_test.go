package datasources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/datasources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestVirtualDiskDataSourceSchema(t *testing.T) {

	t.Parallel()
	d := datasources.NewVirtualDiskDataSource()

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

func TestVirtualDiskDataSourceMetadata(t *testing.T) {

	t.Parallel()
	d := datasources.NewVirtualDiskDataSource()
	testutil.ValidateDataSourceMetadata(t, d, "netbox", "netbox_virtual_disk")
}

func TestVirtualDiskDataSourceConfigure(t *testing.T) {

	t.Parallel()
	d := datasources.NewVirtualDiskDataSource()
	testutil.ValidateDataSourceConfigure(t, d)
}
