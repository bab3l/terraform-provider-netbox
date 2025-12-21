package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccFrontPortTemplateDataSource_basic(t *testing.T) {

	t.Parallel()
	name := testutil.RandomName("test-front-port-template")
	manufacturerName := testutil.RandomName("test-manufacturer-fpt")
	manufacturerSlug := testutil.GenerateSlug(manufacturerName)
	deviceTypeName := testutil.RandomName("test-device-type-fpt")
	deviceTypeSlug := testutil.GenerateSlug(deviceTypeName)
	rearPortName := testutil.RandomName("test-rear-port-fpt")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckManufacturerDestroy,
			testutil.CheckDeviceTypeDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccFrontPortTemplateDataSourceConfig(name, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, rearPortName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_front_port_template.test", "name", name),
					resource.TestCheckResourceAttr("data.netbox_front_port_template.test", "type", "8p8c"),
					resource.TestCheckResourceAttr("data.netbox_front_port_template.test", "label", "Test Label"),
					resource.TestCheckResourceAttr("data.netbox_front_port_template.test", "description", "Test Description"),
				),
			},
		},
	})
}

func testAccFrontPortTemplateDataSourceConfig(name, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, rearPortName string) string {
	return fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
  name = %q
  slug = %q
}

resource "netbox_device_type" "test" {
  model        = %q
  slug         = %q
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_rear_port_template" "test" {
  name         = %q
  device_type  = netbox_device_type.test.id
  type         = "8p8c"
}

resource "netbox_front_port_template" "test" {
  name               = %q
  device_type        = netbox_device_type.test.id
  type               = "8p8c"
  rear_port          = netbox_rear_port_template.test.name
  rear_port_position = 1
  label              = "Test Label"
  description        = "Test Description"
}

data "netbox_front_port_template" "test" {
  id = netbox_front_port_template.test.id
}
`, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, rearPortName, name)
}
