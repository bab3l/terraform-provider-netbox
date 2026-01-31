package datasources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/datasources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestASNRangeDataSourceSchema(t *testing.T) {
	t.Parallel()

	d := datasources.NewASNRangeDataSource()
	req := datasource.SchemaRequest{}
	resp := &datasource.SchemaResponse{}
	d.Schema(context.Background(), req, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Schema returned errors: %v", resp.Diagnostics)
	}

	testutil.ValidateDataSourceSchema(t, resp.Schema.Attributes, testutil.DataSourceValidation{
		LookupAttrs: []string{"id", "name", "slug"},
		ComputedAttrs: []string{
			"id",
			"name",
			"slug",
			"rir",
			"rir_name",
			"start",
			"end",
			"tenant",
			"tenant_name",
			"description",
			"asn_count",
			"display_name",
			"tags",
			"custom_fields",
		},
	})
}

func TestASNRangeDataSourceMetadata(t *testing.T) {
	t.Parallel()

	d := datasources.NewASNRangeDataSource()
	testutil.ValidateDataSourceMetadata(t, d, "netbox", "netbox_asn_range")
}

func TestASNRangeDataSourceConfigure(t *testing.T) {
	t.Parallel()

	d := datasources.NewASNRangeDataSource()
	testutil.ValidateDataSourceConfigure(t, d)
}
