package resources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccWirelessLinkResource_basic(t *testing.T) {
	siteName := testutil.RandomName("test-site-wireless")
	siteSlug := testutil.GenerateSlug(siteName)
	deviceName := testutil.RandomName("test-device-wireless")
	interfaceNameA := "wlan0"
	interfaceNameB := "wlan1"

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			{
				Config: testAccWirelessLinkResourceConfig(siteName, siteSlug, deviceName, interfaceNameA, interfaceNameB),

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_wireless_link.test", "status", "connected"),
					resource.TestCheckResourceAttr("netbox_wireless_link.test", "ssid", "Test SSID"),
				),
			},
			{
				ResourceName:            "netbox_wireless_link.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"interface_a", "interface_b"},
			},
		},
	})
}

func testAccWirelessLinkResourceConfig(siteName, siteSlug, deviceName, interfaceNameA, interfaceNameB string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name = %q
  slug = %q
  status = "active"
}

resource "netbox_manufacturer" "test" {
  name = "Test Manufacturer Wireless"
  slug = "test-manufacturer-wireless"
}

resource "netbox_device_role" "test" {
  name = "Test Device Role Wireless"
  slug = "test-device-role-wireless"
}

resource "netbox_device_type" "test" {
  model = "Test Device Type Wireless"
  slug  = "test-device-type-wireless"
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_device" "test_a" {
  name           = "%s-a"
  device_type    = netbox_device_type.test.id
  role           = netbox_device_role.test.id
  site           = netbox_site.test.id
}

resource "netbox_device" "test_b" {
  name           = "%s-b"
  device_type    = netbox_device_type.test.id
  role           = netbox_device_role.test.id
  site           = netbox_site.test.id
}

resource "netbox_interface" "test_a" {
  name      = %q
  device    = netbox_device.test_a.id
  type      = "ieee802.11ac"
}

resource "netbox_interface" "test_b" {
  name      = %q
  device    = netbox_device.test_b.id
  type      = "ieee802.11ac"
}

resource "netbox_wireless_link" "test" {
  interface_a = netbox_interface.test_a.id
  interface_b = netbox_interface.test_b.id
  ssid        = "Test SSID"
  status      = "connected"
}
`, siteName, siteSlug, deviceName, deviceName, interfaceNameA, interfaceNameB)
}
