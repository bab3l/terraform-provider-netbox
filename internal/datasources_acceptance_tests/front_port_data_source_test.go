package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccFrontPortDataSource_basic(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("test-front-port")
	manufacturerName := testutil.RandomName("test-manufacturer-fp")
	manufacturerSlug := testutil.GenerateSlug(manufacturerName)
	deviceTypeName := testutil.RandomName("test-device-type-fp")
	deviceTypeSlug := testutil.GenerateSlug(deviceTypeName)
	deviceName := testutil.RandomName("test-device-fp")
	siteName := testutil.RandomName("test-site-fp")
	siteSlug := testutil.GenerateSlug(siteName)
	deviceRoleName := testutil.RandomName("test-device-role-fp")
	deviceRoleSlug := testutil.GenerateSlug(deviceRoleName)
	rearPortName := testutil.RandomName("test-rear-port-fp")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterDeviceRoleCleanup(deviceRoleSlug)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
	cleanup.RegisterDeviceCleanup(deviceName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckSiteDestroy,
			testutil.CheckDeviceRoleDestroy,
			testutil.CheckManufacturerDestroy,
			testutil.CheckDeviceDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccFrontPortDataSourceConfig(name, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceName, siteName, siteSlug, deviceRoleName, deviceRoleSlug, rearPortName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_front_port.test", "name", name),
					resource.TestCheckResourceAttr("data.netbox_front_port.test", "type", "8p8c"),
					resource.TestCheckResourceAttr("data.netbox_front_port.test", "label", "Test Label"),
					resource.TestCheckResourceAttr("data.netbox_front_port.test", "description", "Test Description"),
				),
			},
		},
	})
}

func TestAccFrontPortDataSource_IDPreservation(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("test-front-port-id")
	manufacturerName := testutil.RandomName("test-manufacturer-fp-id")
	manufacturerSlug := testutil.GenerateSlug(manufacturerName)
	deviceTypeName := testutil.RandomName("test-device-type-fp-id")
	deviceTypeSlug := testutil.GenerateSlug(deviceTypeName)
	deviceName := testutil.RandomName("test-device-fp-id")
	siteName := testutil.RandomName("test-site-fp-id")
	siteSlug := testutil.GenerateSlug(siteName)
	deviceRoleName := testutil.RandomName("test-device-role-fp-id")
	deviceRoleSlug := testutil.GenerateSlug(deviceRoleName)
	rearPortName := testutil.RandomName("test-rear-port-fp-id")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterDeviceRoleCleanup(deviceRoleSlug)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
	cleanup.RegisterDeviceCleanup(deviceName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckSiteDestroy,
			testutil.CheckDeviceRoleDestroy,
			testutil.CheckManufacturerDestroy,
			testutil.CheckDeviceDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccFrontPortDataSourceConfig(name, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceName, siteName, siteSlug, deviceRoleName, deviceRoleSlug, rearPortName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_front_port.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_front_port.test", "name", name),
					resource.TestCheckResourceAttr("data.netbox_front_port.test", "type", "8p8c"),
				),
			},
		},
	})
}

func testAccFrontPortDataSourceConfig(name, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceName, siteName, siteSlug, deviceRoleName, deviceRoleSlug, rearPortName string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name = %q
  slug = %q
  status = "active"
}

resource "netbox_manufacturer" "test" {
  name = %q
  slug = %q
}

resource "netbox_device_type" "test" {
  model        = %q
  slug         = %q
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_device_role" "test" {
  name = %q
  slug = %q
}

resource "netbox_device" "test" {
  name           = %q
  device_type    = netbox_device_type.test.id
  role           = netbox_device_role.test.id
  site           = netbox_site.test.id
}

resource "netbox_rear_port" "test" {
  name   = %q
  device = netbox_device.test.id
  type   = "8p8c"
}

resource "netbox_front_port" "test" {
  name        = %q
  device      = netbox_device.test.id
  type        = "8p8c"
  rear_port   = netbox_rear_port.test.id
  rear_port_position = 1
  label       = "Test Label"
  description = "Test Description"
}

data "netbox_front_port" "test" {
  id = netbox_front_port.test.id
}
`, siteName, siteSlug, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceRoleName, deviceRoleSlug, deviceName, rearPortName, name)
}
func TestAccFrontPortDataSource_byDeviceAndName(t *testing.T) {

	t.Parallel()
	name := testutil.RandomName("test-front-port")
	manufacturerName := testutil.RandomName("test-manufacturer-fp")
	manufacturerSlug := testutil.GenerateSlug(manufacturerName)
	deviceTypeName := testutil.RandomName("test-device-type-fp")
	deviceTypeSlug := testutil.GenerateSlug(deviceTypeName)
	deviceName := testutil.RandomName("test-device-fp")
	siteName := testutil.RandomName("test-site-fp")
	siteSlug := testutil.GenerateSlug(siteName)
	deviceRoleName := testutil.RandomName("test-device-role-fp")
	deviceRoleSlug := testutil.GenerateSlug(deviceRoleName)
	rearPortName := testutil.RandomName("test-rear-port-fp")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterDeviceRoleCleanup(deviceRoleSlug)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
	cleanup.RegisterDeviceCleanup(deviceName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckSiteDestroy,
			testutil.CheckDeviceRoleDestroy,
			testutil.CheckManufacturerDestroy,
			testutil.CheckDeviceDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccFrontPortDataSourceConfigByDeviceAndName(name, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceName, siteName, siteSlug, deviceRoleName, deviceRoleSlug, rearPortName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_front_port.test", "name", name),
					resource.TestCheckResourceAttr("data.netbox_front_port.test", "type", "8p8c"),
				),
			},
		},
	})
}

func testAccFrontPortDataSourceConfigByDeviceAndName(name, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceName, siteName, siteSlug, deviceRoleName, deviceRoleSlug, rearPortName string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name = %q
  slug = %q
  status = "active"
}

resource "netbox_manufacturer" "test" {
  name = %q
  slug = %q
}

resource "netbox_device_type" "test" {
  model        = %q
  slug         = %q
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_device_role" "test" {
  name = %q
  slug = %q
}

resource "netbox_device" "test" {
  name           = %q
  device_type    = netbox_device_type.test.id
  role           = netbox_device_role.test.id
  site           = netbox_site.test.id
}

resource "netbox_rear_port" "test" {
  name   = %q
  device = netbox_device.test.id
  type   = "8p8c"
}

resource "netbox_front_port" "test" {
  name        = %q
  device      = netbox_device.test.id
  type        = "8p8c"
  rear_port   = netbox_rear_port.test.id
  rear_port_position = 1
  label       = "Test Label"
  description = "Test Description"
}

data "netbox_front_port" "test" {
  device_id = netbox_device.test.id
  name      = netbox_front_port.test.name
}
`, siteName, siteSlug, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceRoleName, deviceRoleSlug, deviceName, rearPortName, name)
}
