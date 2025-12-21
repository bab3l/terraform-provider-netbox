package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccModuleDataSource_basic(t *testing.T) {

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
	moduleTypeName := testutil.RandomName("module-type")
	moduleBayName := testutil.RandomName("module-bay")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterDeviceRoleCleanup(roleSlug)
	cleanup.RegisterManufacturerCleanup(mfgSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)
	cleanup.RegisterDeviceCleanup(deviceName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckSiteDestroy,
			testutil.CheckDeviceRoleDestroy,
			testutil.CheckManufacturerDestroy,
			testutil.CheckDeviceTypeDestroy,
			testutil.CheckDeviceDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccModuleDataSourceConfig(siteName, siteSlug, roleName, roleSlug, mfgName, mfgSlug, deviceTypeName, deviceTypeSlug, deviceName, moduleTypeName, moduleBayName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_module.test", "status", "active"),
					resource.TestCheckResourceAttrSet("data.netbox_module.test", "device_id"),
				),
			},
		},
	})
}

func testAccModuleDataSourceConfig(siteName, siteSlug, roleName, roleSlug, mfgName, mfgSlug, deviceTypeName, deviceTypeSlug, deviceName, moduleTypeName, moduleBayName string) string {
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

resource "netbox_module_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = "%s"
}

resource "netbox_module_bay" "test" {
  device = netbox_device.test.id
  name   = "%s"
}

resource "netbox_module" "test" {
  device      = netbox_device.test.id
  module_bay  = netbox_module_bay.test.id
  module_type = netbox_module_type.test.id
  status      = "active"
}

data "netbox_module" "test" {
  id = netbox_module.test.id
}
`, siteName, siteSlug, roleName, roleSlug, mfgName, mfgSlug, deviceTypeName, deviceTypeSlug, deviceName, moduleTypeName, moduleBayName)
}
