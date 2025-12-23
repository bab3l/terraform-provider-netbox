package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRearPortTemplateDataSource_basic(t *testing.T) {

	t.Parallel()
	manufacturerSlug := testutil.RandomSlug("manufacturer")
	deviceTypeModel := testutil.RandomName("device-type")
	deviceTypeSlug := testutil.RandomSlug("device-type")
	portTemplateName := testutil.RandomName("rear-port-template")

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
				Config: testAccRearPortTemplateDataSourceConfig(portTemplateName, manufacturerSlug, deviceTypeModel, deviceTypeSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_rear_port_template.test", "name", portTemplateName),
					resource.TestCheckResourceAttr("data.netbox_rear_port_template.test", "type", "8p8c"),
					resource.TestCheckResourceAttrSet("data.netbox_rear_port_template.test", "device_type"),
				),
			},
		},
	})
}

func TestAccRearPortTemplateDataSource_byDeviceTypeAndName(t *testing.T) {

	t.Parallel()
	manufacturerSlug := testutil.RandomSlug("manufacturer")
	deviceTypeModel := testutil.RandomName("device-type")
	deviceTypeSlug := testutil.RandomSlug("device-type")
	portTemplateName := testutil.RandomName("rear-port-template")

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
				Config: testAccRearPortTemplateDataSourceConfigByDeviceTypeAndName(portTemplateName, manufacturerSlug, deviceTypeModel, deviceTypeSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_rear_port_template.test", "name", portTemplateName),
					resource.TestCheckResourceAttr("data.netbox_rear_port_template.test", "type", "8p8c"),
					resource.TestCheckResourceAttrSet("data.netbox_rear_port_template.test", "device_type"),
				),
			},
		},
	})
}

func testAccRearPortTemplateDataSourceConfig(name, manufacturerSlug, deviceTypeModel, deviceTypeSlug string) string {
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

resource "netbox_rear_port_template" "test" {
  device_type = netbox_device_type.test.id
  name        = "%s"
  type        = "8p8c"
}

data "netbox_rear_port_template" "test" {
  id = netbox_rear_port_template.test.id
}
`, manufacturerSlug, manufacturerSlug, deviceTypeModel, deviceTypeSlug, name)
}

func testAccRearPortTemplateDataSourceConfigByDeviceTypeAndName(name, manufacturerSlug, deviceTypeModel, deviceTypeSlug string) string {
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

resource "netbox_rear_port_template" "test" {
  device_type = netbox_device_type.test.id
  name        = "%s"
  type        = "8p8c"
}

data "netbox_rear_port_template" "test" {
  device_type = netbox_device_type.test.id
  name        = netbox_rear_port_template.test.name
}
`, manufacturerSlug, manufacturerSlug, deviceTypeModel, deviceTypeSlug, name)
}
