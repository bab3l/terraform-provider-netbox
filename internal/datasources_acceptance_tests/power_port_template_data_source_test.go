package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPowerPortTemplateDataSource_basic(t *testing.T) {

	t.Parallel()
	manufacturerSlug := testutil.RandomSlug("manufacturer")
	deviceTypeModel := testutil.RandomName("device-type")
	deviceTypeSlug := testutil.RandomSlug("device-type")
	powerPortTemplateName := testutil.RandomName("power-port-template")

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
				Config: testAccPowerPortTemplateDataSourceConfig(manufacturerSlug, deviceTypeModel, deviceTypeSlug, powerPortTemplateName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_power_port_template.test", "name", powerPortTemplateName),
					resource.TestCheckResourceAttr("data.netbox_power_port_template.test", "type", "iec-60320-c14"),
				),
			},
		},
	})
}

func TestAccPowerPortTemplateDataSource_byDeviceTypeAndName(t *testing.T) {

	t.Parallel()
	manufacturerSlug := testutil.RandomSlug("manufacturer")
	deviceTypeModel := testutil.RandomName("device-type")
	deviceTypeSlug := testutil.RandomSlug("device-type")
	powerPortTemplateName := testutil.RandomName("power-port-template")

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
				Config: testAccPowerPortTemplateDataSourceConfigByDeviceTypeAndName(manufacturerSlug, deviceTypeModel, deviceTypeSlug, powerPortTemplateName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_power_port_template.test", "name", powerPortTemplateName),
					resource.TestCheckResourceAttr("data.netbox_power_port_template.test", "type", "iec-60320-c14"),
				),
			},
		},
	})
}

func testAccPowerPortTemplateDataSourceConfig(manufacturerSlug, deviceTypeModel, deviceTypeSlug, powerPortTemplateName string) string {
	return fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
  name = "%s"
  slug = "%s"
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = "%s"
  slug         = "%s"
}

resource "netbox_power_port_template" "test" {
  device_type = netbox_device_type.test.id
  name        = "%s"
  type        = "iec-60320-c14"
}

data "netbox_power_port_template" "test" {
  id = netbox_power_port_template.test.id
}
`, manufacturerSlug, manufacturerSlug, deviceTypeModel, deviceTypeSlug, powerPortTemplateName)
}

func testAccPowerPortTemplateDataSourceConfigByDeviceTypeAndName(manufacturerSlug, deviceTypeModel, deviceTypeSlug, powerPortTemplateName string) string {
	return fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
  name = "%s"
  slug = "%s"
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = "%s"
  slug         = "%s"
}

resource "netbox_power_port_template" "test" {
  device_type = netbox_device_type.test.id
  name        = "%s"
  type        = "iec-60320-c14"
}

data "netbox_power_port_template" "test" {
  device_type = netbox_device_type.test.id
  name        = netbox_power_port_template.test.name
}
`, manufacturerSlug, manufacturerSlug, deviceTypeModel, deviceTypeSlug, powerPortTemplateName)
}
