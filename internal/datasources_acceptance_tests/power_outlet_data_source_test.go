package datasources_acceptance_tests

import (
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPowerOutletDataSource_basic(t *testing.T) {

	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPowerOutletDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_power_outlet.test", "name", "Test Power Outlet"),
					resource.TestCheckResourceAttr("data.netbox_power_outlet.test", "type", "iec-60320-c13"),
				),
			},
		},
	})
}

const testAccPowerOutletDataSourceConfig = `
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

resource "netbox_device" "test" {
  name        = "test-device"
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
}

resource "netbox_power_outlet" "test" {
  device = netbox_device.test.id
  name   = "Test Power Outlet"
  type   = "iec-60320-c13"
}

data "netbox_power_outlet" "test" {
  id = netbox_power_outlet.test.id
}
`
