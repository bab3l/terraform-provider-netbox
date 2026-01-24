package datasources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/datasources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestCircuitDataSourceSchema(t *testing.T) {

	t.Parallel()
	d := datasources.NewCircuitDataSource()

	req := datasource.SchemaRequest{}
	resp := &datasource.SchemaResponse{}

	d.Schema(context.Background(), req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("Schema returned errors: %v", resp.Diagnostics)
	}

	testutil.ValidateDataSourceSchema(t, resp.Schema.Attributes, testutil.DataSourceValidation{
		LookupAttrs:   []string{"id", "cid"},
		ComputedAttrs: []string{"id", "cid", "circuit_provider", "provider_account", "type", "status", "tenant", "install_date", "termination_date", "commit_rate", "description", "comments", "display_name", "tags", "custom_fields"},
	})
}

func TestCircuitDataSourceMetadata(t *testing.T) {

	t.Parallel()
	d := datasources.NewCircuitDataSource()
	testutil.ValidateDataSourceMetadata(t, d, "netbox", "netbox_circuit")
}

func TestCircuitDataSourceConfigure(t *testing.T) {

	t.Parallel()
	d := datasources.NewCircuitDataSource()
	testutil.ValidateDataSourceConfigure(t, d)
}
