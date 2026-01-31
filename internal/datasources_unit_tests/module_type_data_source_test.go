package datasources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/datasources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestModuleTypeDataSourceSchema(t *testing.T) {
	t.Parallel()

	d := datasources.NewModuleTypeDataSource()
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

func TestModuleTypeDataSourceMetadata(t *testing.T) {
	t.Parallel()

	d := datasources.NewModuleTypeDataSource()
	testutil.ValidateDataSourceMetadata(t, d, "netbox", "netbox_module_type")
}

func TestModuleTypeDataSourceConfigure(t *testing.T) {
	t.Parallel()

	d := datasources.NewModuleTypeDataSource()
	testutil.ValidateDataSourceConfigure(t, d)
}
