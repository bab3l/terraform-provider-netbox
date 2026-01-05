package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccConsolePortDataSource_IDPreservation(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("cp-ds-id-site")
	siteSlug := testutil.RandomSlug("cp-ds-id-site")
	roleName := testutil.RandomName("cp-ds-id-role")
	roleSlug := testutil.RandomSlug("cp-ds-id-role")
	mfgName := testutil.RandomName("cp-ds-id-mfg")
	mfgSlug := testutil.RandomSlug("cp-ds-id-mfg")
	deviceTypeName := testutil.RandomName("cp-ds-id-device-type")
	deviceTypeSlug := testutil.RandomSlug("cp-ds-id-device-type")
	deviceName := testutil.RandomName("cp-ds-id-device")
	portName := testutil.RandomName("cp-ds-id-console")

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
				Config: testAccConsolePortDataSourceConfig(siteName, siteSlug, roleName, roleSlug, mfgName, mfgSlug, deviceTypeName, deviceTypeSlug, deviceName, portName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_console_port.by_id", "id"),
					resource.TestCheckResourceAttr("data.netbox_console_port.by_id", "name", portName),
					resource.TestCheckResourceAttr("data.netbox_console_port.by_id", "type", "de-9"),
					resource.TestCheckResourceAttrSet("data.netbox_console_port.by_id", "device"),
				),
			},
		},
	})
}

func TestAccConsolePortDataSource_byID(t *testing.T) {
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
	portName := testutil.RandomName("console")

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
				Config: testAccConsolePortDataSourceConfig(siteName, siteSlug, roleName, roleSlug, mfgName, mfgSlug, deviceTypeName, deviceTypeSlug, deviceName, portName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_console_port.by_id", "name", portName),
					resource.TestCheckResourceAttr("data.netbox_console_port.by_id", "type", "de-9"),
					resource.TestCheckResourceAttrSet("data.netbox_console_port.by_id", "id"),
					resource.TestCheckResourceAttrSet("data.netbox_console_port.by_id", "device"),
				),
			},
		},
	})
}

func TestAccConsolePortDataSource_byDeviceAndName(t *testing.T) {
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
	portName := testutil.RandomName("console")

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
				Config: testAccConsolePortDataSourceConfig(siteName, siteSlug, roleName, roleSlug, mfgName, mfgSlug, deviceTypeName, deviceTypeSlug, deviceName, portName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_console_port.by_device_and_name", "name", portName),
					resource.TestCheckResourceAttr("data.netbox_console_port.by_device_and_name", "type", "de-9"),
					resource.TestCheckResourceAttrSet("data.netbox_console_port.by_device_and_name", "id"),
					resource.TestCheckResourceAttrSet("data.netbox_console_port.by_device_and_name", "device"),
				),
			},
		},
	})
}

func testAccConsolePortDataSourceConfig(siteName, siteSlug, roleName, roleSlug, mfgName, mfgSlug, deviceTypeName, deviceTypeSlug, deviceName, portName string) string {
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

resource "netbox_console_port" "test" {
  device = netbox_device.test.id
  name   = "%s"
  type   = "de-9"
}

data "netbox_console_port" "by_id" {
  id = netbox_console_port.test.id
}

data "netbox_console_port" "by_device_and_name" {
  device_id = netbox_device.test.id
  name      = netbox_console_port.test.name
}
`, siteName, siteSlug, roleName, roleSlug, mfgName, mfgSlug, deviceTypeName, deviceTypeSlug, deviceName, portName)
}
