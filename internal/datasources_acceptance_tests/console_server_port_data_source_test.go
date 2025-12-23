package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccConsoleServerPortDataSource_basic(t *testing.T) {

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
				Config: testAccConsoleServerPortDataSourceConfig("console-server-port-0", siteSlug, deviceRoleSlug, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_console_server_port.test", "name", "console-server-port-0"),
					resource.TestCheckResourceAttr("data.netbox_console_server_port.test", "type", "de-9"),
					resource.TestCheckResourceAttrSet("data.netbox_console_server_port.test", "device"),
				),
			},
		},
	})
}

func testAccConsoleServerPortDataSourceConfig(name, siteSlug, deviceRoleSlug, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceName string) string {
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

data "netbox_console_server_port" "test" {
  id = netbox_console_server_port.test.id
}
`, siteSlug, siteSlug, deviceRoleSlug, deviceRoleSlug, manufacturerSlug, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceName, name)
}
