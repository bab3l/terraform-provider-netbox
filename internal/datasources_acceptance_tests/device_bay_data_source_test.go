package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDeviceBayDataSource_basic(t *testing.T) {

	t.Parallel()

	siteSlug := testutil.RandomSlug("site")
	deviceRoleSlug := testutil.RandomSlug("device-role")
	manufacturerSlug := testutil.RandomSlug("manufacturer")
	deviceName := testutil.RandomName("device")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterDeviceRoleCleanup(deviceRoleSlug)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
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
				Config: testAccDeviceBayDataSourceConfig("Bay 1", siteSlug, deviceRoleSlug, manufacturerSlug, deviceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_device_bay.test", "name", "Bay 1"),
					resource.TestCheckResourceAttrSet("data.netbox_device_bay.test", "device"),
				),
			},
		},
	})
}

func testAccDeviceBayDataSourceConfig(name, siteSlug, deviceRoleSlug, manufacturerSlug, deviceName string) string {
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
  manufacturer   = netbox_manufacturer.test.id
  model          = "Test Device Type"
  slug           = "test-device-type"
  subdevice_role = "parent"
}

resource "netbox_device" "test" {
  name        = "%s"
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
}

resource "netbox_device_bay" "test" {
  device = netbox_device.test.id
  name   = "%s"
}

data "netbox_device_bay" "test" {
  id = netbox_device_bay.test.id
}
`, siteSlug, siteSlug, deviceRoleSlug, deviceRoleSlug, manufacturerSlug, manufacturerSlug, deviceName, name)
}
