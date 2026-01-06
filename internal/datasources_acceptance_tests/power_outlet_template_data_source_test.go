package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPowerOutletTemplateDataSource_IDPreservation(t *testing.T) {
	t.Parallel()

	manufacturerSlug := testutil.RandomSlug("manufacturer-id")
	deviceTypeModel := testutil.RandomName("device-type-id")
	deviceTypeSlug := testutil.RandomSlug("device-type-id")
	powerOutletTemplateName := testutil.RandomName("power-outlet-template-id")

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
				Config: testAccPowerOutletTemplateDataSourceConfig(manufacturerSlug, deviceTypeModel, deviceTypeSlug, powerOutletTemplateName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_power_outlet_template.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_power_outlet_template.test", "name", powerOutletTemplateName),
				),
			},
		},
	})
}

func TestAccPowerOutletTemplateDataSource_basic(t *testing.T) {
	t.Parallel()

	manufacturerSlug := testutil.RandomSlug("manufacturer")
	deviceTypeModel := testutil.RandomName("device-type")
	deviceTypeSlug := testutil.RandomSlug("device-type")
	powerOutletTemplateName := testutil.RandomName("power-outlet-template")

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
				Config: testAccPowerOutletTemplateDataSourceConfig(manufacturerSlug, deviceTypeModel, deviceTypeSlug, powerOutletTemplateName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_power_outlet_template.test", "name", powerOutletTemplateName),
					resource.TestCheckResourceAttr("data.netbox_power_outlet_template.test", "type", "iec-60320-c13"),
				),
			},
		},
	})
}

func TestAccPowerOutletTemplateDataSource_byDeviceTypeAndName(t *testing.T) {
	t.Parallel()

	manufacturerSlug := testutil.RandomSlug("manufacturer")
	deviceTypeModel := testutil.RandomName("device-type")
	deviceTypeSlug := testutil.RandomSlug("device-type")
	powerOutletTemplateName := testutil.RandomName("power-outlet-template")

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
				Config: testAccPowerOutletTemplateDataSourceConfigByDeviceTypeAndName(manufacturerSlug, deviceTypeModel, deviceTypeSlug, powerOutletTemplateName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_power_outlet_template.by_device_type", "name", powerOutletTemplateName),
					resource.TestCheckResourceAttrSet("data.netbox_power_outlet_template.by_device_type", "device_type"),
				),
			},
		},
	})
}

func testAccPowerOutletTemplateDataSourceConfig(manufacturerSlug, deviceTypeModel, deviceTypeSlug, powerOutletTemplateName string) string {
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

resource "netbox_power_outlet_template" "test" {
  device_type = netbox_device_type.test.id
  name        = "%s"
  type        = "iec-60320-c13"
}

data "netbox_power_outlet_template" "test" {
  id = netbox_power_outlet_template.test.id
}
`, manufacturerSlug, manufacturerSlug, deviceTypeModel, deviceTypeSlug, powerOutletTemplateName)
}

func testAccPowerOutletTemplateDataSourceConfigByDeviceTypeAndName(manufacturerSlug, deviceTypeModel, deviceTypeSlug, powerOutletTemplateName string) string {
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

resource "netbox_power_outlet_template" "test" {
  device_type = netbox_device_type.test.id
  name        = "%s"
  type        = "iec-60320-c13"
}

data "netbox_power_outlet_template" "by_device_type" {
  device_type = netbox_device_type.test.id
  name        = netbox_power_outlet_template.test.name
}
`, manufacturerSlug, manufacturerSlug, deviceTypeModel, deviceTypeSlug, powerOutletTemplateName)
}
