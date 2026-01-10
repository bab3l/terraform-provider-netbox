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

func TestAccIPAddressesDataSource_queryWithCustomFields(t *testing.T) {
	ipAddress := fmt.Sprintf("192.168.%d.%d/24", acctest.RandIntRange(1, 254), acctest.RandIntRange(1, 254))
	customFieldName := testutil.RandomCustomFieldName("tf_test_ipaddresses_q_cf")
	customFieldValue := "test-value-" + acctest.RandString(8)

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIPAddressCleanup(ipAddress)
	cleanup.RegisterCustomFieldCleanup(customFieldName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIPAddressesDataSourceConfig_withCustomFields(ipAddress, customFieldName, customFieldValue),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_ip_addresses.test", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_ip_addresses.test", "addresses.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_ip_addresses.test", "addresses.0", ipAddress),
					resource.TestCheckResourceAttrPair("data.netbox_ip_addresses.test", "ids.0", "netbox_ip_address.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_ip_addresses.test", "ip_addresses.#", "1"),
					resource.TestCheckResourceAttrPair("data.netbox_ip_addresses.test", "ip_addresses.0.id", "netbox_ip_address.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_ip_addresses.test", "ip_addresses.0.address", ipAddress),
				),
			},
		},
	})
}

func testAccIPAddressesDataSourceConfig_withCustomFields(ipAddress, customFieldName, customFieldValue string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "test" {
  name         = %[2]q
  object_types = ["ipam.ipaddress"]
  type         = "text"
}

resource "netbox_ip_address" "test" {
  address = %[1]q
  status  = "active"

  custom_fields = [
    {
      name  = netbox_custom_field.test.name
      type  = "text"
      value = %[3]q
    }
  ]
}

data "netbox_ip_addresses" "test" {
  filter {
    name   = "custom_field_value"
    values = ["${netbox_custom_field.test.name}=%[3]s"]
  }

  depends_on = [netbox_ip_address.test]
}
`, ipAddress, customFieldName, customFieldValue)
}
