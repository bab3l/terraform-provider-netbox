//go:build customfields
// +build customfields

package datasources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPrefixesDataSource_queryWithCustomFields(t *testing.T) {
	prefix := fmt.Sprintf("10.%d.%d.0/24", acctest.RandIntRange(1, 254), acctest.RandIntRange(1, 254))
	customFieldName := testutil.RandomCustomFieldName("tf_test_prefixes_q_cf")
	customFieldValue := "test-value-" + acctest.RandString(8)

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterPrefixCleanup(prefix)
	cleanup.RegisterCustomFieldCleanup(customFieldName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPrefixesDataSourceConfig_withCustomFields(prefix, customFieldName, customFieldValue),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_prefixes.test", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_prefixes.test", "cidrs.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_prefixes.test", "cidrs.0", prefix),
					resource.TestCheckResourceAttrPair("data.netbox_prefixes.test", "ids.0", "netbox_prefix.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_prefixes.test", "prefixes.#", "1"),
					resource.TestCheckResourceAttrPair("data.netbox_prefixes.test", "prefixes.0.id", "netbox_prefix.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_prefixes.test", "prefixes.0.prefix", prefix),
				),
			},
		},
	})
}

func testAccPrefixesDataSourceConfig_withCustomFields(prefix, customFieldName, customFieldValue string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "test" {
  name         = %[2]q
  object_types = ["ipam.prefix"]
  type         = "text"
}

resource "netbox_prefix" "test" {
  prefix = %[1]q
  status = "active"

  custom_fields = [
    {
      name  = netbox_custom_field.test.name
      type  = "text"
      value = %[3]q
    }
  ]
}

data "netbox_prefixes" "test" {
  filter {
    name   = "custom_field_value"
    values = ["${netbox_custom_field.test.name}=%[3]s"]
  }

  depends_on = [netbox_prefix.test]
}
`, prefix, customFieldName, customFieldValue)
}
