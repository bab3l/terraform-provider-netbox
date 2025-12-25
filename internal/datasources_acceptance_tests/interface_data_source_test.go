package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccInterfaceDataSource_basic(t *testing.T) {

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
	interfaceName := testutil.RandomName("eth")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterDeviceRoleCleanup(roleSlug)
	cleanup.RegisterManufacturerCleanup(mfgSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)
	cleanup.RegisterDeviceCleanup(deviceName)
	cleanup.RegisterInterfaceCleanup(interfaceName, deviceName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckSiteDestroy,
			testutil.CheckDeviceRoleDestroy,
			testutil.CheckManufacturerDestroy,
			testutil.CheckDeviceTypeDestroy,
			testutil.CheckDeviceDestroy,
			testutil.CheckInterfaceDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccInterfaceDataSourceConfig(siteName, siteSlug, roleName, roleSlug, mfgName, mfgSlug, deviceTypeName, deviceTypeSlug, deviceName, interfaceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_interface.test", "name", interfaceName),
					resource.TestCheckResourceAttr("data.netbox_interface.test", "type", "1000base-t"),
					resource.TestCheckResourceAttrSet("data.netbox_interface.test", "device"),
				),
			},
		},
	})
}

func TestAccInterfaceDataSource_IDPreservation(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("site-id")
	siteSlug := testutil.RandomSlug("site-id")
	roleName := testutil.RandomName("device-role-id")
	roleSlug := testutil.RandomSlug("device-role-id")
	mfgName := testutil.RandomName("mfg-id")
	mfgSlug := testutil.RandomSlug("mfg-id")
	deviceTypeName := testutil.RandomName("device-type-id")
	deviceTypeSlug := testutil.RandomSlug("device-type-id")
	deviceName := testutil.RandomName("device-id")
	interfaceName := testutil.RandomName("eth-id")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterDeviceRoleCleanup(roleSlug)
	cleanup.RegisterManufacturerCleanup(mfgSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)
	cleanup.RegisterDeviceCleanup(deviceName)
	cleanup.RegisterInterfaceCleanup(interfaceName, deviceName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckSiteDestroy,
			testutil.CheckDeviceRoleDestroy,
			testutil.CheckManufacturerDestroy,
			testutil.CheckDeviceTypeDestroy,
			testutil.CheckDeviceDestroy,
			testutil.CheckInterfaceDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccInterfaceDataSourceConfig(siteName, siteSlug, roleName, roleSlug, mfgName, mfgSlug, deviceTypeName, deviceTypeSlug, deviceName, interfaceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_interface.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_interface.test", "name", interfaceName),
					resource.TestCheckResourceAttr("data.netbox_interface.test", "type", "1000base-t"),
				),
			},
		},
	})
}

func testAccInterfaceDataSourceConfig(siteName, siteSlug, roleName, roleSlug, mfgName, mfgSlug, deviceTypeName, deviceTypeSlug, deviceName, interfaceName string) string {
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

resource "netbox_interface" "test" {
  device = netbox_device.test.id
  name   = "%s"
  type   = "1000base-t"
}

data "netbox_interface" "test" {
  id = netbox_interface.test.id
}
`, siteName, siteSlug, roleName, roleSlug, mfgName, mfgSlug, deviceTypeName, deviceTypeSlug, deviceName, interfaceName)
}
