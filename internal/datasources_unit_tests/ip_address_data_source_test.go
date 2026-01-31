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

func TestIPAddressDataSourceSchema(t *testing.T) {
	t.Parallel()

	d := datasources.NewIPAddressDataSource()
	req := datasource.SchemaRequest{}
	resp := &datasource.SchemaResponse{}
	d.Schema(context.Background(), req, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Schema returned errors: %v", resp.Diagnostics)
	}

	testutil.ValidateDataSourceSchema(t, resp.Schema.Attributes, testutil.DataSourceValidation{
		LookupAttrs:   []string{"id", "address"},
		ComputedAttrs: []string{"id", "address", "display_name", "vrf", "vrf_id", "tenant", "tenant_id", "status", "role", "assigned_object_type", "assigned_object_id", "nat_inside", "dns_name", "description", "comments", "tags", "custom_fields"},
	})

	testutil.ValidateDataSourceStringAttributeHasValidatorType(
		t,
		resp.Schema.Attributes["address"],
		"address",
		reflect.TypeOf(validators.IPAddressWithPrefixValidator{}),
	)
}

func TestIPAddressDataSourceMetadata(t *testing.T) {
	t.Parallel()

	d := datasources.NewIPAddressDataSource()
	testutil.ValidateDataSourceMetadata(t, d, "netbox", "netbox_ip_address")
}

func TestIPAddressDataSourceConfigure(t *testing.T) {
	t.Parallel()

	d := datasources.NewIPAddressDataSource()
	testutil.ValidateDataSourceConfigure(t, d)
}
