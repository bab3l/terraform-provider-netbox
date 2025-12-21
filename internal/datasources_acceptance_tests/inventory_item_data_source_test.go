package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccInventoryItemDataSource_basic(t *testing.T) {

	t.Parallel()

	siteName := testutil.RandomName("site")
	siteSlug := testutil.RandomSlug("site")
	roleName := testutil.RandomName("device-role")
	roleSlug := testutil.RandomSlug("device-role")
	mfgName := testutil.RandomName("mfg")
	mfgSlug := testutil.RandomSlug("mfg")
	deviceTypeName := testutil.RandomName("device-type")
	deviceTypeSlug := testutil.RandomSlug("device-type")
	deviceName := testutil.RandomName("device")
	itemName := testutil.RandomName("item")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterDeviceRoleCleanup(roleSlug)
	cleanup.RegisterManufacturerCleanup(mfgSlug)
	cleanup.RegisterDeviceCleanup(deviceName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckSiteDestroy,
			testutil.CheckDeviceRoleDestroy,
			testutil.CheckManufacturerDestroy,
			testutil.CheckDeviceDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccInventoryItemDataSourceConfig(siteName, siteSlug, roleName, roleSlug, mfgName, mfgSlug, deviceTypeName, deviceTypeSlug, deviceName, itemName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_inventory_item.test", "name", itemName),
					resource.TestCheckResourceAttrSet("data.netbox_inventory_item.test", "device_id"),
				),
			},
		},
	})
}

func testAccInventoryItemDataSourceConfig(siteName, siteSlug, roleName, roleSlug, mfgName, mfgSlug, deviceTypeName, deviceTypeSlug, deviceName, itemName string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name = "%s"
  slug = "%s"
}

resource "netbox_device_role" "test" {
  name = "%s"
  slug = "%s"
}

resource "netbox_manufacturer" "test" {
  name = "%s"
  slug = "%s"
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = "%s"
  slug         = "%s"
}

resource "netbox_device" "test" {
  name        = "%s"
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
}

resource "netbox_inventory_item" "test" {
  device = netbox_device.test.id
  name   = "%s"
}

data "netbox_inventory_item" "test" {
  id = netbox_inventory_item.test.id
}
`, siteName, siteSlug, roleName, roleSlug, mfgName, mfgSlug, deviceTypeName, deviceTypeSlug, deviceName, itemName)
}
