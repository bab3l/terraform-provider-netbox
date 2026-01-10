//go:build customfields
// +build customfields

package datasources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAggregateDataSource_customFields(t *testing.T) {
	customFieldName := testutil.RandomCustomFieldName("tf_test_agg_ds_cf")
	prefix := "192.0.2.0/24"
	rirName := testutil.RandomName("tf-test-rir-ds-cf")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAggregateDataSourceConfig_customFields(customFieldName, prefix, rirName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_aggregate.test", "prefix", prefix),
					resource.TestCheckResourceAttr("data.netbox_aggregate.test", "custom_fields.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_aggregate.test", "custom_fields.0.name", customFieldName),
					resource.TestCheckResourceAttr("data.netbox_aggregate.test", "custom_fields.0.type", "text"),
					resource.TestCheckResourceAttr("data.netbox_aggregate.test", "custom_fields.0.value", "datasource-test-value"),
				),
			},
		},
	})
}

func testAccAggregateDataSourceConfig_customFields(customFieldName, prefix, rirName string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "test" {
  name         = %q
  object_types = ["ipam.aggregate"]
  type         = "text"
}

resource "netbox_rir" "test" {
  name = %q
  slug = %q
}

resource "netbox_aggregate" "test" {
  prefix = %q
  rir    = netbox_rir.test.id

  custom_fields = [
    {
      name  = netbox_custom_field.test.name
      type  = "text"
      value = "datasource-test-value"
    }
  ]
}

data "netbox_aggregate" "test" {
  prefix = netbox_aggregate.test.prefix

  depends_on = [netbox_aggregate.test]
}
`, customFieldName, rirName, rirName, prefix)
}
