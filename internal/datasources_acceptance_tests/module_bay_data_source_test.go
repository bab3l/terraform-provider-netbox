package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccModuleBayDataSource_IDPreservation(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-module-bay-site-ds-id")
	siteSlug := testutil.GenerateSlug(siteName)
	deviceRoleName := testutil.RandomName("tf-test-module-bay-role-ds-id")
	deviceRoleSlug := testutil.GenerateSlug(deviceRoleName)
	mfgName := testutil.RandomName("tf-test-module-bay-mfg-ds-id")
	mfgSlug := testutil.GenerateSlug(mfgName)
	deviceTypeModel := testutil.RandomName("tf-test-module-bay-dt-id")
	deviceTypeSlug := testutil.RandomSlug("module-bay-dt-id")
	deviceName := testutil.RandomName("tf-test-module-bay-device-ds-id")
	moduleBayName := testutil.RandomName("module-bay-id")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterDeviceRoleCleanup(deviceRoleSlug)
	cleanup.RegisterManufacturerCleanup(mfgSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)
	cleanup.RegisterDeviceCleanup(deviceName)
	cleanup.RegisterModuleBayCleanup(moduleBayName)

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
				Config: testAccModuleBayDataSourceConfig(siteName, siteSlug, deviceRoleName, deviceRoleSlug, mfgName, mfgSlug, deviceTypeModel, deviceTypeSlug, deviceName, moduleBayName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_module_bay.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_module_bay.test", "name", moduleBayName),
				),
			},
		},
	})
}

func TestAccModuleBayDataSource_basic(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-module-bay-site-ds")
	siteSlug := testutil.GenerateSlug(siteName)
	deviceRoleName := testutil.RandomName("tf-test-module-bay-role-ds")
	deviceRoleSlug := testutil.GenerateSlug(deviceRoleName)
	mfgName := testutil.RandomName("tf-test-module-bay-mfg-ds")
	mfgSlug := testutil.GenerateSlug(mfgName)
	deviceTypeModel := testutil.RandomName("tf-test-module-bay-dt")
	deviceTypeSlug := testutil.RandomSlug("module-bay-dt")
	deviceName := testutil.RandomName("tf-test-module-bay-device-ds")
	moduleBayName := testutil.RandomName("module-bay")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterDeviceRoleCleanup(deviceRoleSlug)
	cleanup.RegisterManufacturerCleanup(mfgSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)
	cleanup.RegisterDeviceCleanup(deviceName)
	cleanup.RegisterModuleBayCleanup(moduleBayName)

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
				Config: testAccModuleBayDataSourceConfig(siteName, siteSlug, deviceRoleName, deviceRoleSlug, mfgName, mfgSlug, deviceTypeModel, deviceTypeSlug, deviceName, moduleBayName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_module_bay.test", "name", moduleBayName),
					resource.TestCheckResourceAttrSet("data.netbox_module_bay.test", "device_id"),
				),
			},
		},
	})
}

func TestAccModuleBayDataSource_byDeviceAndName(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-module-bay-site-ds")
	siteSlug := testutil.GenerateSlug(siteName)
	deviceRoleName := testutil.RandomName("tf-test-module-bay-role-ds")
	deviceRoleSlug := testutil.GenerateSlug(deviceRoleName)
	mfgName := testutil.RandomName("tf-test-module-bay-mfg-ds")
	mfgSlug := testutil.GenerateSlug(mfgName)
	deviceTypeModel := testutil.RandomName("tf-test-module-bay-dt")
	deviceTypeSlug := testutil.RandomSlug("module-bay-dt")
	deviceName := testutil.RandomName("tf-test-module-bay-device-ds")
	moduleBayName := testutil.RandomName("module-bay")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterDeviceRoleCleanup(deviceRoleSlug)
	cleanup.RegisterManufacturerCleanup(mfgSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)
	cleanup.RegisterDeviceCleanup(deviceName)
	cleanup.RegisterModuleBayCleanup(moduleBayName)

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
				Config: testAccModuleBayDataSourceConfigByDeviceAndName(siteName, siteSlug, deviceRoleName, deviceRoleSlug, mfgName, mfgSlug, deviceTypeModel, deviceTypeSlug, deviceName, moduleBayName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_module_bay.by_name", "name", moduleBayName),
					resource.TestCheckResourceAttrSet("data.netbox_module_bay.by_name", "device_id"),
				),
			},
		},
	})
}

func testAccModuleBayDataSourceConfig(siteName, siteSlug, deviceRoleName, deviceRoleSlug, mfgName, mfgSlug, deviceTypeModel, deviceTypeSlug, deviceName, moduleBayName string) string {
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
  model        = %[7]q
  slug         = %[8]q
}

resource "netbox_device" "test" {
  name        = %[9]q
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
}

resource "netbox_module_bay" "test" {
  device = netbox_device.test.id
  name   = %[10]q
}

data "netbox_module_bay" "test" {
  id = netbox_module_bay.test.id
}
`, siteName, siteSlug, deviceRoleName, deviceRoleSlug, mfgName, mfgSlug, deviceTypeModel, deviceTypeSlug, deviceName, moduleBayName)
}

func testAccModuleBayDataSourceConfigByDeviceAndName(siteName, siteSlug, deviceRoleName, deviceRoleSlug, mfgName, mfgSlug, deviceTypeModel, deviceTypeSlug, deviceName, moduleBayName string) string {
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
  model        = %[7]q
  slug         = %[8]q
}

resource "netbox_device" "test" {
  name        = %[9]q
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
}

resource "netbox_module_bay" "test" {
  device = netbox_device.test.id
  name   = %[10]q
}

data "netbox_module_bay" "by_name" {
  device_id = netbox_device.test.id
  name      = netbox_module_bay.test.name
}
`, siteName, siteSlug, deviceRoleName, deviceRoleSlug, mfgName, mfgSlug, deviceTypeModel, deviceTypeSlug, deviceName, moduleBayName)
}
