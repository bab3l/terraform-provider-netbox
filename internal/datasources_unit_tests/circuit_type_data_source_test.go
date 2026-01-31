package datasources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/datasources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestCircuitTypeDataSourceSchema(t *testing.T) {
	t.Parallel()

	d := datasources.NewCircuitTypeDataSource()
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

func TestCircuitTypeDataSourceMetadata(t *testing.T) {
	t.Parallel()

	d := datasources.NewCircuitTypeDataSource()
	testutil.ValidateDataSourceMetadata(t, d, "netbox", "netbox_circuit_type")
}

func TestCircuitTypeDataSourceConfigure(t *testing.T) {
	t.Parallel()

	d := datasources.NewCircuitTypeDataSource()
	testutil.ValidateDataSourceConfigure(t, d)
}
