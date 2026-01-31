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

func TestVLANDataSourceSchema(t *testing.T) {
	t.Parallel()

	d := datasources.NewVLANDataSource()
	req := datasource.SchemaRequest{}
	resp := &datasource.SchemaResponse{}
	d.Schema(context.Background(), req, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Schema returned errors: %v", resp.Diagnostics)
	}

	testutil.ValidateDataSourceSchema(t, resp.Schema.Attributes, testutil.DataSourceValidation{
		LookupAttrs: []string{"id", "vid", "name"},
		ComputedAttrs: []string{
			"id",
			"vid",
			"name",
			"site",
			"site_id",
			"group",
			"group_id",
			"tenant",
			"tenant_id",
			"status",
			"role",
			"description",
			"comments",
			"display_name",
			"tags",
			"custom_fields",
		},
	})

	testutil.ValidateDataSourceInt32AttributeHasValidatorType(
		t,
		resp.Schema.Attributes["vid"],
		"vid",
		reflect.TypeOf(validators.VLANIDInt32Validator{}),
	)
}

func TestVLANDataSourceMetadata(t *testing.T) {
	t.Parallel()

	d := datasources.NewVLANDataSource()
	testutil.ValidateDataSourceMetadata(t, d, "netbox", "netbox_vlan")
}

func TestVLANDataSourceConfigure(t *testing.T) {
	t.Parallel()

	d := datasources.NewVLANDataSource()
	testutil.ValidateDataSourceConfigure(t, d)
}
