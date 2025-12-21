package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccModuleBayDataSource_basic(t *testing.T) {

	t.Parallel()

	cleanup := testutil.NewCleanupResource(t)

	siteName := testutil.RandomName("tf-test-module-bay-site-ds")
	siteSlug := testutil.GenerateSlug(siteName)
	deviceRoleName := testutil.RandomName("tf-test-module-bay-role-ds")
	deviceRoleSlug := testutil.GenerateSlug(deviceRoleName)
	mfgName := testutil.RandomName("tf-test-module-bay-mfg-ds")
	mfgSlug := testutil.GenerateSlug(mfgName)
	deviceName := testutil.RandomName("tf-test-module-bay-device-ds")

	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterDeviceRoleCleanup(deviceRoleSlug)
	cleanup.RegisterManufacturerCleanup(mfgSlug)
	cleanup.RegisterDeviceTypeCleanup("test-device-type-module-bay-ds")
	cleanup.RegisterDeviceCleanup(deviceName)
	cleanup.RegisterModuleBayCleanup("Test Module Bay")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckSiteDestroy,
			testutil.CheckDeviceRoleDestroy,
			testutil.CheckManufacturerDestroy,
			testutil.CheckDeviceTypeDestroy,
			testutil.CheckDeviceDestroy,
			testutil.CheckModuleBayDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccModuleBayDataSourceConfig(siteName, siteSlug, deviceRoleName, deviceRoleSlug, mfgName, mfgSlug, deviceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_module_bay.test", "name", "Test Module Bay"),
					resource.TestCheckResourceAttrSet("data.netbox_module_bay.test", "device_id"),
				),
			},
		},
	})
}

func testAccModuleBayDataSourceConfig(siteName, siteSlug, deviceRoleName, deviceRoleSlug, mfgName, mfgSlug, deviceName string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_device_role" "test" {
  name = %[3]q
  slug = %[4]q
}

resource "netbox_manufacturer" "test" {
  name = %[5]q
  slug = %[6]q
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = "Test Device Type"
  slug         = "test-device-type-module-bay-ds"
}

resource "netbox_device" "test" {
  name        = %[7]q
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
}

resource "netbox_module_bay" "test" {
  device = netbox_device.test.id
  name   = "Test Module Bay"
}

data "netbox_module_bay" "test" {
  id = netbox_module_bay.test.id
}
`, siteName, siteSlug, deviceRoleName, deviceRoleSlug, mfgName, mfgSlug, deviceName)
}
