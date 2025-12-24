package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccWirelessLinkDataSource_byID(t *testing.T) {

	t.Parallel()

	mfgName := testutil.RandomName("mfg")
	mfgSlug := testutil.RandomSlug("mfg")
	deviceTypeName := testutil.RandomName("device-type")
	deviceTypeSlug := testutil.RandomSlug("device-type")
	siteName := testutil.RandomName("site")
	siteSlug := testutil.RandomSlug("site")
	roleName := testutil.RandomName("device-role")
	roleSlug := testutil.RandomSlug("device-role")
	deviceAName := testutil.RandomName("device-a")
	deviceBName := testutil.RandomName("device-b")
	wirelessLinkName := testutil.RandomName("wl")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(mfgSlug)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterDeviceRoleCleanup(roleSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckManufacturerDestroy,
			testutil.CheckSiteDestroy,
			testutil.CheckDeviceRoleDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccWirelessLinkDataSourceByIDConfig(mfgName, mfgSlug, deviceTypeName, deviceTypeSlug, siteName, siteSlug, roleName, roleSlug, deviceAName, deviceBName, wirelessLinkName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_wireless_link.test", "ssid", wirelessLinkName),
					resource.TestCheckResourceAttrSet("data.netbox_wireless_link.test", "interface_a"),
					resource.TestCheckResourceAttrSet("data.netbox_wireless_link.test", "interface_b"),
					resource.TestCheckResourceAttrSet("data.netbox_wireless_link.test", "id"),
				),
			},
		},
	})
}

func testAccWirelessLinkDataSourceByIDConfig(mfgName, mfgSlug, deviceTypeName, deviceTypeSlug, siteName, siteSlug, roleName, roleSlug, deviceAName, deviceBName, wirelessLinkName string) string {
	return fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
  name = "%s"
  slug = "%s"
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = "%s"
  slug         = "%s"
}

resource "netbox_site" "test" {
  name = "%s"
  slug = "%s"
}

resource "netbox_device_role" "test" {
  name = "%s"
  slug = "%s"
}

resource "netbox_device" "test_a" {
  name        = "%s"
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
}

resource "netbox_device" "test_b" {
  name        = "%s"
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
`, mfgName, mfgSlug, deviceTypeName, deviceTypeSlug, siteName, siteSlug, roleName, roleSlug, deviceAName, deviceBName, wirelessLinkName)
}
