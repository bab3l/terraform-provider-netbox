package datasources_test

import (
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPowerOutletTemplateDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPowerOutletTemplateDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_power_outlet_template.test", "name", "Test Power Outlet Template"),
					resource.TestCheckResourceAttr("data.netbox_power_outlet_template.test", "type", "iec-60320-c13"),
				),
			},
		},
	})
}

const testAccPowerOutletTemplateDataSourceConfig = `
resource "netbox_manufacturer" "test" {
  name = "Test Manufacturer"
  slug = "test-manufacturer"
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = "Test Device Type"
  slug         = "test-device-type"
}

resource "netbox_power_outlet_template" "test" {
  device_type = netbox_device_type.test.id
  name        = "Test Power Outlet Template"
  type        = "iec-60320-c13"
}

data "netbox_power_outlet_template" "test" {
  id = netbox_power_outlet_template.test.id
}
`
