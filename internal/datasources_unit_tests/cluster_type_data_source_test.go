package datasources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/datasources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestClusterTypeDataSourceSchema(t *testing.T) {
	t.Parallel()

	d := datasources.NewClusterTypeDataSource()
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

func TestClusterTypeDataSourceMetadata(t *testing.T) {
	t.Parallel()

	d := datasources.NewClusterTypeDataSource()
	testutil.ValidateDataSourceMetadata(t, d, "netbox", "netbox_cluster_type")
}

func TestClusterTypeDataSourceConfigure(t *testing.T) {
	t.Parallel()

	d := datasources.NewClusterTypeDataSource()
	testutil.ValidateDataSourceConfigure(t, d)
}
