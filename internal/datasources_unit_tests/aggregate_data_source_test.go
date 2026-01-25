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

func TestAggregateDataSourceSchema(t *testing.T) {
	t.Parallel()

	d := datasources.NewAggregateDataSource()
	req := datasource.SchemaRequest{}
	resp := &datasource.SchemaResponse{}
	d.Schema(context.Background(), req, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Schema returned errors: %v", resp.Diagnostics)
	}
	testutil.ValidateDataSourceSchema(t, resp.Schema.Attributes, testutil.DataSourceValidation{
		LookupAttrs: []string{"id", "prefix"},
		ComputedAttrs: []string{
			"id",
			"prefix",
			"rir",
			"rir_name",
			"tenant",
			"tenant_name",
			"date_added",
			"description",
			"comments",
			"display_name",
			"tags",
			"custom_fields",
		},
	})

	testutil.ValidateDataSourceStringAttributeHasValidatorType(
		t,
		resp.Schema.Attributes["prefix"],
		"prefix",
		reflect.TypeOf(validators.IPPrefixValidator{}),
	)
}

func TestAggregateDataSourceMetadata(t *testing.T) {
	t.Parallel()

	d := datasources.NewAggregateDataSource()
	testutil.ValidateDataSourceMetadata(t, d, "netbox", "netbox_aggregate")
}

func TestAggregateDataSourceConfigure(t *testing.T) {
	t.Parallel()

	d := datasources.NewAggregateDataSource()
	testutil.ValidateDataSourceConfigure(t, d)
}
