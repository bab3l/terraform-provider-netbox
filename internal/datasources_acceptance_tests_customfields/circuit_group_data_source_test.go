//go:build customfields
// +build customfields

package datasources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCircuitGroupDataSource_customFields(t *testing.T) {
	customFieldName := testutil.RandomCustomFieldName("tf_test_circuitgroup_ds_cf")
	circuitGroupName := testutil.RandomName("tf-test-circuitgroup-ds-cf")
	circuitGroupSlug := testutil.RandomSlug("tf-test-circuitgroup-ds-cf")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCircuitGroupDataSourceConfig_customFields(customFieldName, circuitGroupName, circuitGroupSlug),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_circuit_group.test", "name", circuitGroupName),
					resource.TestCheckResourceAttr("data.netbox_circuit_group.test", "custom_fields.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_circuit_group.test", "custom_fields.0.name", customFieldName),
					resource.TestCheckResourceAttr("data.netbox_circuit_group.test", "custom_fields.0.type", "text"),
					resource.TestCheckResourceAttr("data.netbox_circuit_group.test", "custom_fields.0.value", "datasource-test-value"),
				),
			},
		},
	})
}

func testAccCircuitGroupDataSourceConfig_customFields(customFieldName, circuitGroupName, circuitGroupSlug string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "test" {
  name         = %q
  object_types = ["circuits.circuitgroup"]
  type         = "text"
}

resource "netbox_circuit_group" "test" {
  name = %q
  slug = %q

  custom_fields = [
    {
      name  = netbox_custom_field.test.name
      type  = "text"
      value = "datasource-test-value"
    }
  ]
}

data "netbox_circuit_group" "test" {
  name = %q

  depends_on = [netbox_circuit_group.test]
}
`, customFieldName, circuitGroupName, circuitGroupSlug, circuitGroupName)
}
