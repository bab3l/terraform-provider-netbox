package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDeviceBayDataSource_basic(t *testing.T) {

	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDeviceBayDataSourceConfig("Bay 1"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_device_bay.test", "name", "Bay 1"),
					resource.TestCheckResourceAttrSet("data.netbox_device_bay.test", "device"),
				),
			},
		},
	})
}

func testAccDeviceBayDataSourceConfig(name string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name = "Test Site"
  slug = "test-site"
}

resource "netbox_device_role" "test" {
  name = "Test Device Role"
  slug = "test-device-role"
}

resource "netbox_manufacturer" "test" {
  name = "Test Manufacturer"
  slug = "test-manufacturer"
}

resource "netbox_device_type" "test" {
  manufacturer   = netbox_manufacturer.test.id
  model          = "Test Device Type"
  slug           = "test-device-type"
  subdevice_role = "parent"
}

resource "netbox_device" "test" {
  name        = "test-device"
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
`, name)
}
