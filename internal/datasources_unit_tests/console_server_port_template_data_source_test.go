package datasources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/datasources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestConsoleServerPortTemplateDataSourceSchema(t *testing.T) {
	d := datasources.NewConsoleServerPortTemplateDataSource()

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

func TestConsoleServerPortTemplateDataSourceMetadata(t *testing.T) {
	d := datasources.NewConsoleServerPortTemplateDataSource()
	testutil.ValidateDataSourceMetadata(t, d, "netbox", "netbox_console_server_port_template")
}

func TestConsoleServerPortTemplateDataSourceConfigure(t *testing.T) {
	d := datasources.NewConsoleServerPortTemplateDataSource()
	testutil.ValidateDataSourceConfigure(t, d)
}
