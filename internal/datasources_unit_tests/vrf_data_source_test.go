package datasources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/datasources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestVRFDataSourceSchema(t *testing.T) {

	t.Parallel()
	d := datasources.NewVRFDataSource()

	req := datasource.SchemaRequest{}
	resp := &datasource.SchemaResponse{}

	d.Schema(context.Background(), req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("Schema returned errors: %v", resp.Diagnostics)
	}

	testutil.ValidateDataSourceSchema(t, resp.Schema.Attributes, testutil.DataSourceValidation{
		LookupAttrs:   []string{"id", "name"},
		ComputedAttrs: []string{"id", "name", "rd", "tenant", "enforce_unique", "import_targets", "export_targets", "description", "display_name", "comments", "tags", "custom_fields"},
	})
}

func TestVRFDataSourceMetadata(t *testing.T) {

	t.Parallel()
	d := datasources.NewVRFDataSource()
	testutil.ValidateDataSourceMetadata(t, d, "netbox", "netbox_vrf")
}

func TestVRFDataSourceConfigure(t *testing.T) {

	t.Parallel()
	d := datasources.NewVRFDataSource()
	testutil.ValidateDataSourceConfigure(t, d)
}
