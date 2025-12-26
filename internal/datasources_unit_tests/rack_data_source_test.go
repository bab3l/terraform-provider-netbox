package datasources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/datasources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestRackDataSourceSchema(t *testing.T) {

	t.Parallel()
	d := datasources.NewRackDataSource()

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

func TestRackDataSourceMetadata(t *testing.T) {

	t.Parallel()
	d := datasources.NewRackDataSource()
	testutil.ValidateDataSourceMetadata(t, d, "netbox", "netbox_rack")
}

func TestRackDataSourceConfigure(t *testing.T) {

	t.Parallel()
	d := datasources.NewRackDataSource()
	testutil.ValidateDataSourceConfigure(t, d)
}
