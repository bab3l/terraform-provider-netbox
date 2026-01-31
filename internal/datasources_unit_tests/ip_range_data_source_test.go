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

func TestIPRangeDataSourceSchema(t *testing.T) {
	t.Parallel()

	d := datasources.NewIPRangeDataSource()
	req := datasource.SchemaRequest{}
	resp := &datasource.SchemaResponse{}
	d.Schema(context.Background(), req, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Schema returned errors: %v", resp.Diagnostics)
	}

	testutil.ValidateDataSourceSchema(t, resp.Schema.Attributes, testutil.DataSourceValidation{
		LookupAttrs: []string{"id", "start_address", "end_address"},
		ComputedAttrs: []string{
			"id",
			"display_name",
			"start_address",
			"end_address",
			"size",
			"vrf",
			"vrf_id",
			"tenant",
			"tenant_id",
			"status",
			"role",
			"mark_utilized",
			"description",
			"comments",
			"tags",
			"custom_fields",
		},
	})

	testutil.ValidateDataSourceStringAttributeHasValidatorType(
		t,
		resp.Schema.Attributes["start_address"],
		"start_address",
		reflect.TypeOf(validators.IPAddressValidator{}),
	)
	testutil.ValidateDataSourceStringAttributeHasValidatorType(
		t,
		resp.Schema.Attributes["end_address"],
		"end_address",
		reflect.TypeOf(validators.IPAddressValidator{}),
	)
}

func TestIPRangeDataSourceMetadata(t *testing.T) {
	t.Parallel()

	d := datasources.NewIPRangeDataSource()
	testutil.ValidateDataSourceMetadata(t, d, "netbox", "netbox_ip_range")
}

func TestIPRangeDataSourceConfigure(t *testing.T) {
	t.Parallel()

	d := datasources.NewIPRangeDataSource()
	testutil.ValidateDataSourceConfigure(t, d)
}
