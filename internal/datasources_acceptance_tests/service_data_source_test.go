package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccServiceDataSource_basic(t *testing.T) {

	t.Parallel()

	cleanup := testutil.NewCleanupResource(t)

	siteName := testutil.RandomName("tf-test-service-site-ds")
	siteSlug := testutil.GenerateSlug(siteName)
	deviceRoleName := testutil.RandomName("tf-test-service-role-ds")
	deviceRoleSlug := testutil.GenerateSlug(deviceRoleName)
	mfgName := testutil.RandomName("tf-test-service-mfg-ds")
	mfgSlug := testutil.GenerateSlug(mfgName)
	deviceName := testutil.RandomName("tf-test-service-device-ds")

	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterDeviceRoleCleanup(deviceRoleSlug)
	cleanup.RegisterManufacturerCleanup(mfgSlug)
	cleanup.RegisterDeviceTypeCleanup("test-device-type-ds")
	cleanup.RegisterDeviceCleanup(deviceName)
	cleanup.RegisterServiceCleanup("Test Service")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckSiteDestroy,
			testutil.CheckDeviceRoleDestroy,
			testutil.CheckManufacturerDestroy,
			testutil.CheckDeviceTypeDestroy,
			testutil.CheckDeviceDestroy,
			testutil.CheckServiceDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccServiceDataSourceConfig(siteName, siteSlug, deviceRoleName, deviceRoleSlug, mfgName, mfgSlug, deviceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_service.test", "name", "Test Service"),
					resource.TestCheckResourceAttr("data.netbox_service.test", "protocol", "tcp"),
				),
			},
		},
	})
}

func testAccServiceDataSourceConfig(siteName, siteSlug, deviceRoleName, deviceRoleSlug, mfgName, mfgSlug, deviceName string) string {
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
  slug         = "test-device-type-ds"
}

resource "netbox_device" "test" {
  name        = %[7]q
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
}

resource "netbox_service" "test" {
  device   = netbox_device.test.id
  name     = "Test Service"
  protocol = "tcp"
  ports    = [80, 443]
}

data "netbox_service" "test" {
  id = netbox_service.test.id
}
`, siteName, siteSlug, deviceRoleName, deviceRoleSlug, mfgName, mfgSlug, deviceName)
}
