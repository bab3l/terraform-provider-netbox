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

func TestAccIPAddressDataSource_customFields(t *testing.T) {
	ipAddress := fmt.Sprintf("192.168.%d.%d/24", acctest.RandIntRange(1, 254), acctest.RandIntRange(1, 254))
	customFieldName := testutil.RandomCustomFieldName("tf_test_ipaddress_ds_cf")
	customFieldValue := "test-value-" + acctest.RandString(8)

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIPAddressCleanup(ipAddress)
	cleanup.RegisterCustomFieldCleanup(customFieldName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIPAddressDataSourceConfigWithCustomFields(ipAddress, customFieldName, customFieldValue),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify datasource returns the custom field
					resource.TestCheckResourceAttr("data.netbox_ip_address.test", "custom_fields.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_ip_address.test", "custom_fields.0.name", customFieldName),
					resource.TestCheckResourceAttr("data.netbox_ip_address.test", "custom_fields.0.type", "text"),
					resource.TestCheckResourceAttr("data.netbox_ip_address.test", "custom_fields.0.value", customFieldValue),
				),
			},
		},
	})
}

func testAccIPAddressDataSourceConfigWithCustomFields(ipAddress, customFieldName, customFieldValue string) string {
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

data "netbox_ip_address" "test" {
  address = netbox_ip_address.test.address
  depends_on = [netbox_ip_address.test]
}
`, ipAddress, customFieldName, customFieldValue)
}
