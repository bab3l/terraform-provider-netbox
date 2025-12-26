package resources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccWirelessLinkResource_basic(t *testing.T) {

	t.Parallel()
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

func TestAccWirelessLinkResource_IDPreservation(t *testing.T) {
	t.Parallel()
	siteName := testutil.RandomName("site-wl-id")
	siteSlug := testutil.GenerateSlug(siteName)
	deviceName := testutil.RandomName("device-wl-id")
	interfaceNameA := "wlan0"
	interfaceNameB := "wlan1"

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			{
				Config: testAccWirelessLinkResourceConfig(siteName, siteSlug, deviceName, interfaceNameA, interfaceNameB),

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_wireless_link.test", "id"),
					resource.TestCheckResourceAttr("netbox_wireless_link.test", "status", "connected"),
					resource.TestCheckResourceAttr("netbox_wireless_link.test", "ssid", "Test SSID"),
				),
			},
		},
	})
}

func testAccWirelessLinkResourceConfig(siteName, siteSlug, deviceName, interfaceNameA, interfaceNameB string) string {
	manufacturerName := testutil.RandomName("mfr")
	manufacturerSlug := testutil.GenerateSlug(manufacturerName)
	deviceRoleName := testutil.RandomName("role")
	deviceRoleSlug := testutil.GenerateSlug(deviceRoleName)
	deviceTypeName := testutil.RandomName("dtype")
	deviceTypeSlug := testutil.GenerateSlug(deviceTypeName)

	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name = %[1]q
  slug = %[2]q
  status = "active"
}

resource "netbox_manufacturer" "test" {
  name = %[3]q
  slug = %[4]q
}

resource "netbox_device_role" "test" {
  name = %[5]q
  slug = %[6]q
}

resource "netbox_device_type" "test" {
  model = %[7]q
  slug  = %[8]q
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_device" "test_a" {
  name           = "%[9]s-a"
  device_type    = netbox_device_type.test.id
  role           = netbox_device_role.test.id
  site           = netbox_site.test.id
}

resource "netbox_device" "test_b" {
  name           = "%[9]s-b"
  device_type    = netbox_device_type.test.id
  role           = netbox_device_role.test.id
  site           = netbox_site.test.id
}

resource "netbox_interface" "test_a" {
  name      = %[10]q
  device    = netbox_device.test_a.id
  type      = "ieee802.11ac"
}

resource "netbox_interface" "test_b" {
  name      = %[11]q
  device    = netbox_device.test_b.id
  type      = "ieee802.11ac"
}

resource "netbox_wireless_link" "test" {
  interface_a = netbox_interface.test_a.id
  interface_b = netbox_interface.test_b.id
  ssid        = "Test SSID"
  status      = "connected"
}
`, siteName, siteSlug, manufacturerName, manufacturerSlug, deviceRoleName, deviceRoleSlug, deviceTypeName, deviceTypeSlug, deviceName, interfaceNameA, interfaceNameB)
}

func TestAccConsistency_WirelessLink_LiteralNames(t *testing.T) {
	t.Parallel()
	ssid := testutil.RandomName("tf-test-ssid-lit")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccWirelessLinkConsistencyLiteralNamesConfig(ssid),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_wireless_link.test", "id"),
					resource.TestCheckResourceAttr("netbox_wireless_link.test", "ssid", ssid),
				),
			},
		},
	})
}

func testAccWirelessLinkConsistencyLiteralNamesConfig(ssid string) string {
	manufacturerName := testutil.RandomName("mfr")
	manufacturerSlug := testutil.GenerateSlug(manufacturerName)
	deviceRoleName := testutil.RandomName("role")
	deviceRoleSlug := testutil.GenerateSlug(deviceRoleName)
	deviceTypeName := testutil.RandomName("dtype")
	deviceTypeSlug := testutil.GenerateSlug(deviceTypeName)
	deviceName := testutil.RandomName("device")
	siteName := testutil.RandomName("site")
	siteSlug := testutil.GenerateSlug(siteName)

	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name   = "%[10]s"
  slug   = "%[11]s"
  status = "active"
}

resource "netbox_manufacturer" "test" {
  name = "%[3]s"
  slug = "%[4]s"
}

resource "netbox_device_role" "test" {
  name = "%[5]s"
  slug = "%[6]s"
}

resource "netbox_device_type" "test" {
  model        = "%[7]s"
  slug         = "%[8]s"
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_device" "test_a" {
  name        = "%[9]s-a"
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
}

resource "netbox_device" "test_b" {
  name        = "%[9]s-b"
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
}

resource "netbox_interface" "test_a" {
  name   = "wlan0"
  device = netbox_device.test_a.id
  type   = "ieee802.11ac"
}

resource "netbox_interface" "test_b" {
  name   = "wlan1"
  device = netbox_device.test_b.id
  type   = "ieee802.11ac"
}

resource "netbox_wireless_link" "test" {
  interface_a = netbox_interface.test_a.id
  interface_b = netbox_interface.test_b.id
  ssid        = "%[1]s"
  status      = "connected"
}
`, ssid, siteName, manufacturerName, manufacturerSlug, deviceRoleName, deviceRoleSlug, deviceTypeName, deviceTypeSlug, deviceName, siteName, siteSlug)
}
