package datasources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/datasources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestTagDataSourceSchema(t *testing.T) {

	t.Parallel()
	d := datasources.NewTagDataSource()

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

func TestTagDataSourceMetadata(t *testing.T) {

	t.Parallel()
	d := datasources.NewTagDataSource()
	testutil.ValidateDataSourceMetadata(t, d, "netbox", "netbox_tag")
}

func TestTagDataSourceConfigure(t *testing.T) {

	t.Parallel()
	d := datasources.NewTagDataSource()
	testutil.ValidateDataSourceConfigure(t, d)
}
