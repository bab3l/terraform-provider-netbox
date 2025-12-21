package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPowerOutletTemplateDataSource_basic(t *testing.T) {

	t.Parallel()
	manufacturerSlug := testutil.RandomSlug("manufacturer")
	deviceTypeSlug := testutil.RandomSlug("device-type")

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
				Config: testAccPowerOutletTemplateDataSourceConfig(manufacturerSlug, deviceTypeSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_power_outlet_template.test", "name", "Test Power Outlet Template"),
					resource.TestCheckResourceAttr("data.netbox_power_outlet_template.test", "type", "iec-60320-c13"),
				),
			},
		},
	})
}

func testAccPowerOutletTemplateDataSourceConfig(manufacturerSlug, deviceTypeSlug string) string {
	return fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
  name = "%s"
  slug = "%s"
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = "Test Device Type"
  slug         = "%s"
}

resource "netbox_power_outlet_template" "test" {
  device_type = netbox_device_type.test.id
  name        = "Test Power Outlet Template"
  type        = "iec-60320-c13"
}

data "netbox_power_outlet_template" "test" {
  id = netbox_power_outlet_template.test.id
}
`, manufacturerSlug, manufacturerSlug, deviceTypeSlug)
}
