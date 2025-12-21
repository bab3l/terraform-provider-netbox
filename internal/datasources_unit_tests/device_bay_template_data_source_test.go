package datasources_unit_tests

import (
	"context"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/datasources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestDeviceBayTemplateDataSourceSchema(t *testing.T) {
	d := datasources.NewDeviceBayTemplateDataSource()

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

func TestDeviceBayTemplateDataSourceMetadata(t *testing.T) {
	d := datasources.NewDeviceBayTemplateDataSource()
	testutil.ValidateDataSourceMetadata(t, d, "netbox", "netbox_device_bay_template")
}

func TestDeviceBayTemplateDataSourceConfigure(t *testing.T) {
	d := datasources.NewDeviceBayTemplateDataSource()
	testutil.ValidateDataSourceConfigure(t, d)
}
