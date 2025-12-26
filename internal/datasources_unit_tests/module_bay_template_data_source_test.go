package datasources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/datasources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestModuleBayTemplateDataSourceSchema(t *testing.T) {

	t.Parallel()
	d := datasources.NewModuleBayTemplateDataSource()

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

func TestModuleBayTemplateDataSourceMetadata(t *testing.T) {

	t.Parallel()
	d := datasources.NewModuleBayTemplateDataSource()
	testutil.ValidateDataSourceMetadata(t, d, "netbox", "netbox_module_bay_template")
}

func TestModuleBayTemplateDataSourceConfigure(t *testing.T) {

	t.Parallel()
	d := datasources.NewModuleBayTemplateDataSource()
	testutil.ValidateDataSourceConfigure(t, d)
}
