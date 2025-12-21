package datasources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/datasources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestCircuitTerminationDataSourceSchema(t *testing.T) {
	d := datasources.NewCircuitTerminationDataSource()

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

func TestCircuitTerminationDataSourceMetadata(t *testing.T) {
	d := datasources.NewCircuitTerminationDataSource()
	testutil.ValidateDataSourceMetadata(t, d, "netbox", "netbox_circuit_termination")
}

func TestCircuitTerminationDataSourceConfigure(t *testing.T) {
	d := datasources.NewCircuitTerminationDataSource()
	testutil.ValidateDataSourceConfigure(t, d)
}
