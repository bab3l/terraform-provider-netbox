package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccWirelessLinkDataSource_basic(t *testing.T) {

	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccWirelessLinkDataSourceConfig("test-wireless-link"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_wireless_link.test", "ssid", "test-wireless-link"),
					resource.TestCheckResourceAttrSet("data.netbox_wireless_link.test", "interface_a"),
					resource.TestCheckResourceAttrSet("data.netbox_wireless_link.test", "interface_b"),
				),
			},
		},
	})
}

func testAccWirelessLinkDataSourceConfig(name string) string {
	return fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
  name = "Test Manufacturer"
  slug = "test-manufacturer"
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = "Test Device Type"
  slug         = "test-device-type"
}

resource "netbox_site" "test" {
  name = "Test Site"
  slug = "test-site"
}

resource "netbox_device_role" "test" {
  name = "Test Device Role"
  slug = "test-device-role"
}

resource "netbox_device" "test_a" {
  name        = "test-device-a"
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
}

resource "netbox_device" "test_b" {
  name        = "test-device-b"
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
}

resource "netbox_interface" "test_a" {
  device = netbox_device.test_a.id
  name   = "wlan0"
  type   = "ieee802.11ac"
}

resource "netbox_interface" "test_b" {
  device = netbox_device.test_b.id
  name   = "wlan0"
  type   = "ieee802.11ac"
}

resource "netbox_wireless_link" "test" {
  interface_a = netbox_interface.test_a.id
  interface_b = netbox_interface.test_b.id
  ssid        = "%s"
}

data "netbox_wireless_link" "test" {
  id = netbox_wireless_link.test.id
}
`, name)
}
