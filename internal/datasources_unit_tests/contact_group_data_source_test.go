package datasources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/datasources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestContactGroupDataSourceSchema(t *testing.T) {
	d := datasources.NewContactGroupDataSource()

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

func TestContactGroupDataSourceMetadata(t *testing.T) {
	d := datasources.NewContactGroupDataSource()
	testutil.ValidateDataSourceMetadata(t, d, "netbox", "netbox_contact_group")
}

func TestContactGroupDataSourceConfigure(t *testing.T) {
	d := datasources.NewContactGroupDataSource()
	testutil.ValidateDataSourceConfigure(t, d)
}
