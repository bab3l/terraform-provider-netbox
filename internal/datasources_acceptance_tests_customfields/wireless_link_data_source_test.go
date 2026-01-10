//go:build customfields
// +build customfields

package datasources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccWirelessLinkDataSource_customFields(t *testing.T) {
	customFieldName := testutil.RandomCustomFieldName("tf_test_wlink_ds_cf")
	ssid := testutil.RandomName("tf-test-wlink-ds-cf")
	siteName := testutil.RandomName("tf-test-site-ds-cf")
	deviceName1 := testutil.RandomName("tf-test-dev1-ds-cf")
	deviceName2 := testutil.RandomName("tf-test-dev2-ds-cf")
	interfaceName1 := "wlan0"
	interfaceName2 := "wlan0"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccWirelessLinkDataSourceConfig_customFields(customFieldName, ssid, siteName, deviceName1, deviceName2, interfaceName1, interfaceName2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_wireless_link.test", "ssid", ssid),
					resource.TestCheckResourceAttr("data.netbox_wireless_link.test", "custom_fields.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_wireless_link.test", "custom_fields.0.name", customFieldName),
					resource.TestCheckResourceAttr("data.netbox_wireless_link.test", "custom_fields.0.type", "text"),
					resource.TestCheckResourceAttr("data.netbox_wireless_link.test", "custom_fields.0.value", "datasource-test-value"),
				),
			},
		},
	})
}

func testAccWirelessLinkDataSourceConfig_customFields(customFieldName, ssid, siteName, deviceName1, deviceName2, interfaceName1, interfaceName2 string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "test" {
  name         = %[1]q
  object_types = ["wireless.wirelesslink"]
  type         = "text"
}

resource "netbox_site" "test" {
  name = %[3]q
  slug = %[3]q
}

resource "netbox_device_role" "test" {
  name  = "test-role-%[3]s"
  slug  = "test-role-%[3]s"
  color = "ff0000"
}

resource "netbox_manufacturer" "test" {
  name = "test-manufacturer-%[3]s"
  slug = "test-manufacturer-%[3]s"
}

resource "netbox_device_type" "test" {
  model        = "test-model-%[3]s"
  slug         = "test-model-%[3]s"
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_device" "test1" {
  name        = %[4]q
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
}

resource "netbox_device" "test2" {
  name        = %[5]q
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
}

resource "netbox_interface" "test1" {
  device = netbox_device.test1.id
  name   = %[6]q
  type   = "ieee802.11a"
}

resource "netbox_interface" "test2" {
  device = netbox_device.test2.id
  name   = %[7]q
  type   = "ieee802.11a"
}

resource "netbox_wireless_link" "test" {
  interface_a = netbox_interface.test1.id
  interface_b = netbox_interface.test2.id
  ssid        = %[2]q

  custom_fields = [
    {
      name  = netbox_custom_field.test.name
      type  = "text"
      value = "datasource-test-value"
    }
  ]
}

data "netbox_wireless_link" "test" {
  id = netbox_wireless_link.test.id

  depends_on = [netbox_wireless_link.test]
}
`, customFieldName, ssid, siteName, deviceName1, deviceName2, interfaceName1, interfaceName2)
}
