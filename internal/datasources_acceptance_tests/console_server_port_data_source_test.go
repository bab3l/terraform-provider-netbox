package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccConsoleServerPortDataSource_IDPreservation(t *testing.T) {
	t.Parallel()

	siteSlug := testutil.RandomSlug("site-id")
	deviceRoleSlug := testutil.RandomSlug("device-role-id")
	manufacturerSlug := testutil.RandomSlug("manufacturer-id")
	deviceTypeName := testutil.RandomName("tf-test-dt-id")
	deviceTypeSlug := testutil.RandomSlug("tf-test-dt-id")
	deviceName := testutil.RandomName("device-id")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterDeviceRoleCleanup(deviceRoleSlug)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
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
				Config: testAccConsoleServerPortDataSourceConfigByID("eth0", siteSlug, deviceRoleSlug, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_console_server_port.by_id", "id"),
					resource.TestCheckResourceAttr("data.netbox_console_server_port.by_id", "name", "eth0"),
				),
			},
		},
	})
}

func TestAccConsoleServerPortDataSource_byID(t *testing.T) {
	t.Parallel()

	siteSlug := testutil.RandomSlug("site")
	deviceRoleSlug := testutil.RandomSlug("device-role")
	manufacturerSlug := testutil.RandomSlug("manufacturer")
	deviceTypeName := testutil.RandomName("tf-test-dt")
	deviceTypeSlug := testutil.RandomSlug("tf-test-dt")
	deviceName := testutil.RandomName("device")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterDeviceRoleCleanup(deviceRoleSlug)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
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
				Config: testAccConsoleServerPortDataSourceConfigByID("console-server-port-0", siteSlug, deviceRoleSlug, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_console_server_port.by_id", "name", "console-server-port-0"),
					resource.TestCheckResourceAttr("data.netbox_console_server_port.by_id", "type", "de-9"),
					resource.TestCheckResourceAttrSet("data.netbox_console_server_port.by_id", "device"),
				),
			},
		},
	})
}

func TestAccConsoleServerPortDataSource_byDeviceAndName(t *testing.T) {
	t.Parallel()

	siteSlug := testutil.RandomSlug("site")
	deviceRoleSlug := testutil.RandomSlug("device-role")
	manufacturerSlug := testutil.RandomSlug("manufacturer")
	deviceTypeName := testutil.RandomName("tf-test-dt")
	deviceTypeSlug := testutil.RandomSlug("tf-test-dt")
	deviceName := testutil.RandomName("device")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterDeviceRoleCleanup(deviceRoleSlug)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
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
				Config: testAccConsoleServerPortDataSourceConfigByDeviceAndName("console-server-port-0", siteSlug, deviceRoleSlug, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_console_server_port.by_device_and_name", "name", "console-server-port-0"),
					resource.TestCheckResourceAttr("data.netbox_console_server_port.by_device_and_name", "type", "de-9"),
					resource.TestCheckResourceAttrSet("data.netbox_console_server_port.by_device_and_name", "device"),
				),
			},
		},
	})
}

func testAccConsoleServerPortDataSourceConfigByID(name, siteSlug, deviceRoleSlug, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceName string) string {
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
  model        = %q
  slug         = %q
}

resource "netbox_device" "test" {
  name        = "%s"
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
}

resource "netbox_console_server_port" "test" {
  device = netbox_device.test.id
  name   = "%s"
  type   = "de-9"
}

data "netbox_console_server_port" "by_id" {
  id = netbox_console_server_port.test.id
}
`, siteSlug, siteSlug, deviceRoleSlug, deviceRoleSlug, manufacturerSlug, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceName, name)
}

func testAccConsoleServerPortDataSourceConfigByDeviceAndName(name, siteSlug, deviceRoleSlug, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceName string) string {
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
  model        = %q
  slug         = %q
}

resource "netbox_device" "test" {
  name        = "%s"
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
}

resource "netbox_console_server_port" "test" {
  device = netbox_device.test.id
  name   = "%s"
  type   = "de-9"
}

data "netbox_console_server_port" "by_device_and_name" {
  device_id = netbox_device.test.id
  name      = netbox_console_server_port.test.name
}
`, siteSlug, siteSlug, deviceRoleSlug, deviceRoleSlug, manufacturerSlug, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceName, name)
}
