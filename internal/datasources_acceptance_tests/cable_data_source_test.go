package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCableDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCableDataSourceConfig("cat6"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_cable.test", "type", "cat6"),
					resource.TestCheckResourceAttr("data.netbox_cable.test", "status", "connected"),
				),
			},
		},
	})
}

func testAccCableDataSourceConfig(cableType string) string {
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
  manufacturer = netbox_manufacturer.test.id
  model        = "Test Device Type"
  slug         = "test-device-type"
}

resource "netbox_device" "test_a" {
  name        = "test-device-a"
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
}

resource "netbox_device" "test_b" {
  name        = "test-device-b"
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
}

resource "netbox_interface" "test_a" {
  device = netbox_device.test_a.id
  name   = "eth0"
  type   = "1000base-t"
}

resource "netbox_interface" "test_b" {
  device = netbox_device.test_b.id
  name   = "eth0"
  type   = "1000base-t"
}

resource "netbox_cable" "test" {
  a_terminations = [
    {
      object_type = "dcim.interface"
      object_id   = netbox_interface.test_a.id
    }
  ]
  b_terminations = [
    {
      object_type = "dcim.interface"
      object_id   = netbox_interface.test_b.id
    }
  ]
  status = "connected"
  type   = "%s"
}

data "netbox_cable" "test" {
  id = netbox_cable.test.id
}
`, cableType)
}
