package datasources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/datasources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestVirtualMachineDataSourceSchema(t *testing.T) {
	d := datasources.NewVirtualMachineDataSource()

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

func TestVirtualMachineDataSourceMetadata(t *testing.T) {
	d := datasources.NewVirtualMachineDataSource()
	testutil.ValidateDataSourceMetadata(t, d, "netbox", "netbox_virtual_machine")
}

func TestVirtualMachineDataSourceConfigure(t *testing.T) {
	d := datasources.NewVirtualMachineDataSource()
	testutil.ValidateDataSourceConfigure(t, d)
}
