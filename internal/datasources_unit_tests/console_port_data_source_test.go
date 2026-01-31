package datasources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/datasources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestConsolePortDataSourceSchema(t *testing.T) {
	t.Parallel()

	d := datasources.NewConsolePortDataSource()
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

func TestConsolePortDataSourceMetadata(t *testing.T) {
	t.Parallel()

	d := datasources.NewConsolePortDataSource()
	testutil.ValidateDataSourceMetadata(t, d, "netbox", "netbox_console_port")
}

func TestConsolePortDataSourceConfigure(t *testing.T) {
	t.Parallel()

	d := datasources.NewConsolePortDataSource()
	testutil.ValidateDataSourceConfigure(t, d)
}
