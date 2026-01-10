//go:build customfields
// +build customfields

package datasources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccFHRPGroupDataSource_customFields(t *testing.T) {
	customFieldName := testutil.RandomCustomFieldName("tf_test_fhrp_ds_cf")
	groupName := testutil.RandomName("tf-test-fhrp-ds-cf")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFHRPGroupDataSourceConfig_customFields(customFieldName, groupName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_fhrp_group.test", "name", groupName),
					resource.TestCheckResourceAttr("data.netbox_fhrp_group.test", "custom_fields.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_fhrp_group.test", "custom_fields.0.name", customFieldName),
					resource.TestCheckResourceAttr("data.netbox_fhrp_group.test", "custom_fields.0.type", "text"),
					resource.TestCheckResourceAttr("data.netbox_fhrp_group.test", "custom_fields.0.value", "datasource-test-value"),
				),
			},
		},
	})
}

func testAccFHRPGroupDataSourceConfig_customFields(customFieldName, groupName string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "test" {
  name         = %q
  object_types = ["ipam.fhrpgroup"]
  type         = "text"
}

resource "netbox_fhrp_group" "test" {
  name     = %q
  protocol = "vrrp2"
  group_id = 10

  custom_fields = [
    {
      name  = netbox_custom_field.test.name
      type  = "text"
      value = "datasource-test-value"
    }
  ]
}

data "netbox_fhrp_group" "test" {
  id = netbox_fhrp_group.test.id

  depends_on = [netbox_fhrp_group.test]
}
`, customFieldName, groupName)
}
