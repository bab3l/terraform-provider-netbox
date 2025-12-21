package datasources_acceptance_tests

import (
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPowerPortTemplateDataSource_basic(t *testing.T) {

	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPowerPortTemplateDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_power_port_template.test", "name", "Test Power Port Template"),
					resource.TestCheckResourceAttr("data.netbox_power_port_template.test", "type", "iec-60320-c14"),
				),
			},
		},
	})
}

const testAccPowerPortTemplateDataSourceConfig = `
resource "netbox_manufacturer" "test" {
  name = "Test Manufacturer"
  slug = "test-manufacturer"
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = "Test Device Type"
  slug         = "test-device-type"
}

resource "netbox_power_port_template" "test" {
  device_type = netbox_device_type.test.id
  name        = "Test Power Port Template"
  type        = "iec-60320-c14"
}

data "netbox_power_port_template" "test" {
  id = netbox_power_port_template.test.id
}
`
