package datasources_test

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRearPortDataSource_basic(t *testing.T) {
	name := testutil.RandomName("test-rear-port")
	manufacturerName := testutil.RandomName("test-manufacturer-rp")
	manufacturerSlug := testutil.GenerateSlug(manufacturerName)
	deviceTypeName := testutil.RandomName("test-device-type-rp")
	deviceTypeSlug := testutil.GenerateSlug(deviceTypeName)
	deviceName := testutil.RandomName("test-device-rp")
	siteName := testutil.RandomName("test-site-rp")
	siteSlug := testutil.GenerateSlug(siteName)
	deviceRoleName := testutil.RandomName("test-device-role-rp")
	deviceRoleSlug := testutil.GenerateSlug(deviceRoleName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRearPortDataSourceConfig(name, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceName, siteName, siteSlug, deviceRoleName, deviceRoleSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_rear_port.test", "name", name),
					resource.TestCheckResourceAttr("data.netbox_rear_port.test", "type", "8p8c"),
					resource.TestCheckResourceAttr("data.netbox_rear_port.test", "label", "Test Label"),
					resource.TestCheckResourceAttr("data.netbox_rear_port.test", "description", "Test Description"),
				),
			},
		},
	})
}

func testAccRearPortDataSourceConfig(name, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceName, siteName, siteSlug, deviceRoleName, deviceRoleSlug string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name = %q
  slug = %q
  status = "active"
}

resource "netbox_manufacturer" "test" {
  name = %q
  slug = %q
}

resource "netbox_device_type" "test" {
  model        = %q
  slug         = %q
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_device_role" "test" {
  name = %q
  slug = %q
}

resource "netbox_device" "test" {
  name           = %q
  device_type    = netbox_device_type.test.id
  role           = netbox_device_role.test.id
  site           = netbox_site.test.id
}

resource "netbox_rear_port" "test" {
  name        = %q
  device      = netbox_device.test.id
  type        = "8p8c"
  label       = "Test Label"
  description = "Test Description"
}

data "netbox_rear_port" "test" {
  id = netbox_rear_port.test.id
}
`, siteName, siteSlug, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceRoleName, deviceRoleSlug, deviceName, name)
}
