package datasources_unit_tests

import (
	"context"
	"reflect"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/datasources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/bab3l/terraform-provider-netbox/internal/validators"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestASNDataSourceSchema(t *testing.T) {

	t.Parallel()
	d := datasources.NewASNDataSource()

	req := datasource.SchemaRequest{}
	resp := &datasource.SchemaResponse{}

	d.Schema(context.Background(), req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("Schema returned errors: %v", resp.Diagnostics)
	}

	testutil.ValidateDataSourceSchema(t, resp.Schema.Attributes, testutil.DataSourceValidation{
		LookupAttrs: []string{"id", "asn"},
		ComputedAttrs: []string{
			"id",
			"asn",
			"display_name",
			"rir",
			"rir_id",
			"tenant",
			"tenant_id",
			"description",
			"comments",
			"tags",
			"custom_fields",
		},
	})

	testutil.ValidateDataSourceInt64AttributeHasValidatorType(
		t,
		resp.Schema.Attributes["asn"],
		"asn",
		reflect.TypeOf(validators.ASNInt64Validator{}),
	)
}

func TestASNDataSourceMetadata(t *testing.T) {

	t.Parallel()
	d := datasources.NewASNDataSource()
	testutil.ValidateDataSourceMetadata(t, d, "netbox", "netbox_asn")
}

func TestASNDataSourceConfigure(t *testing.T) {

	t.Parallel()
	d := datasources.NewASNDataSource()
	testutil.ValidateDataSourceConfigure(t, d)
}
