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

func TestPrefixDataSourceSchema(t *testing.T) {

	t.Parallel()
	d := datasources.NewPrefixDataSource()

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
			"site",
			"site_id",
			"vrf",
			"vrf_id",
			"tenant",
			"tenant_id",
			"vlan",
			"vlan_id",
			"status",
			"role",
			"description",
			"comments",
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

func TestPrefixDataSourceMetadata(t *testing.T) {

	t.Parallel()
	d := datasources.NewPrefixDataSource()
	testutil.ValidateDataSourceMetadata(t, d, "netbox", "netbox_prefix")
}

func TestPrefixDataSourceConfigure(t *testing.T) {

	t.Parallel()
	d := datasources.NewPrefixDataSource()
	testutil.ValidateDataSourceConfigure(t, d)
}
