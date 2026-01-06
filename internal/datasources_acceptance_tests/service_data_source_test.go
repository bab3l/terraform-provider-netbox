package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccServiceDataSource_IDPreservation(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-service-site-ds-id")
	siteSlug := testutil.GenerateSlug(siteName)
	deviceRoleName := testutil.RandomName("tf-test-service-role-ds-id")
	deviceRoleSlug := testutil.GenerateSlug(deviceRoleName)
	mfgName := testutil.RandomName("tf-test-service-mfg-ds-id")
	mfgSlug := testutil.GenerateSlug(mfgName)
	deviceTypeModel := testutil.RandomName("tf-test-device-type-id")
	deviceTypeSlug := testutil.RandomSlug("device-type-id")
	deviceName := testutil.RandomName("tf-test-service-device-ds-id")
	serviceName := testutil.RandomName("tf-test-service-id")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterDeviceRoleCleanup(deviceRoleSlug)
	cleanup.RegisterManufacturerCleanup(mfgSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)
	cleanup.RegisterDeviceCleanup(deviceName)
	cleanup.RegisterServiceCleanup(serviceName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckSiteDestroy,
			testutil.CheckDeviceRoleDestroy,
			testutil.CheckManufacturerDestroy,
			testutil.CheckDeviceTypeDestroy,
			testutil.CheckDeviceDestroy,
			testutil.CheckServiceDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccServiceDataSourceConfig(siteName, siteSlug, deviceRoleName, deviceRoleSlug, mfgName, mfgSlug, deviceTypeModel, deviceTypeSlug, deviceName, serviceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_service.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_service.test", "name", serviceName),
				),
			},
		},
	})
}

func TestAccServiceDataSource_basic(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-service-site-ds")
	siteSlug := testutil.GenerateSlug(siteName)
	deviceRoleName := testutil.RandomName("tf-test-service-role-ds")
	deviceRoleSlug := testutil.GenerateSlug(deviceRoleName)
	mfgName := testutil.RandomName("tf-test-service-mfg-ds")
	mfgSlug := testutil.GenerateSlug(mfgName)
	deviceTypeModel := testutil.RandomName("tf-test-device-type")
	deviceTypeSlug := testutil.RandomSlug("device-type")
	deviceName := testutil.RandomName("tf-test-service-device-ds")
	serviceName := testutil.RandomName("tf-test-service")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterDeviceRoleCleanup(deviceRoleSlug)
	cleanup.RegisterManufacturerCleanup(mfgSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)
	cleanup.RegisterDeviceCleanup(deviceName)
	cleanup.RegisterServiceCleanup(serviceName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckSiteDestroy,
			testutil.CheckDeviceRoleDestroy,
			testutil.CheckManufacturerDestroy,
			testutil.CheckDeviceTypeDestroy,
			testutil.CheckDeviceDestroy,
			testutil.CheckServiceDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccServiceDataSourceConfig(siteName, siteSlug, deviceRoleName, deviceRoleSlug, mfgName, mfgSlug, deviceTypeModel, deviceTypeSlug, deviceName, serviceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_service.test", "name", serviceName),
					resource.TestCheckResourceAttr("data.netbox_service.test", "protocol", "tcp"),
				),
			},
		},
	})
}

func TestAccServiceDataSource_byNameAndDevice(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-service-site-ds")
	siteSlug := testutil.GenerateSlug(siteName)
	deviceRoleName := testutil.RandomName("tf-test-service-role-ds")
	deviceRoleSlug := testutil.GenerateSlug(deviceRoleName)
	mfgName := testutil.RandomName("tf-test-service-mfg-ds")
	mfgSlug := testutil.GenerateSlug(mfgName)
	deviceTypeModel := testutil.RandomName("tf-test-device-type")
	deviceTypeSlug := testutil.RandomSlug("device-type")
	deviceName := testutil.RandomName("tf-test-service-device-ds")
	serviceName := testutil.RandomName("tf-test-service")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterDeviceRoleCleanup(deviceRoleSlug)
	cleanup.RegisterManufacturerCleanup(mfgSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)
	cleanup.RegisterDeviceCleanup(deviceName)
	cleanup.RegisterServiceCleanup(serviceName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckSiteDestroy,
			testutil.CheckDeviceRoleDestroy,
			testutil.CheckManufacturerDestroy,
			testutil.CheckDeviceTypeDestroy,
			testutil.CheckDeviceDestroy,
			testutil.CheckServiceDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccServiceDataSourceConfigByNameAndDevice(siteName, siteSlug, deviceRoleName, deviceRoleSlug, mfgName, mfgSlug, deviceTypeModel, deviceTypeSlug, deviceName, serviceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_service.test", "name", serviceName),
					resource.TestCheckResourceAttr("data.netbox_service.test", "protocol", "tcp"),
				),
			},
		},
	})
}

func testAccServiceDataSourceConfig(siteName, siteSlug, deviceRoleName, deviceRoleSlug, mfgName, mfgSlug, deviceTypeModel, deviceTypeSlug, deviceName, serviceName string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_device_role" "test" {
  name = %[3]q
  slug = %[4]q
}

resource "netbox_manufacturer" "test" {
  name = %[5]q
  slug = %[6]q
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = %[7]q
  slug         = %[8]q
}

resource "netbox_device" "test" {
  name        = %[9]q
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
}

resource "netbox_service" "test" {
  device   = netbox_device.test.id
  name     = %[10]q
  protocol = "tcp"
  ports    = [80, 443]
}

data "netbox_service" "test" {
  id = netbox_service.test.id
}
`, siteName, siteSlug, deviceRoleName, deviceRoleSlug, mfgName, mfgSlug, deviceTypeModel, deviceTypeSlug, deviceName, serviceName)
}

func testAccServiceDataSourceConfigByNameAndDevice(siteName, siteSlug, deviceRoleName, deviceRoleSlug, mfgName, mfgSlug, deviceTypeModel, deviceTypeSlug, deviceName, serviceName string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_device_role" "test" {
  name = %[3]q
  slug = %[4]q
}

resource "netbox_manufacturer" "test" {
  name = %[5]q
  slug = %[6]q
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = %[7]q
  slug         = %[8]q
}

resource "netbox_device" "test" {
  name        = %[9]q
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
}

resource "netbox_service" "test" {
  device   = netbox_device.test.id
  name     = %[10]q
  protocol = "tcp"
  ports    = [80, 443]
}

data "netbox_service" "test" {
  name   = netbox_service.test.name
  device = netbox_device.test.id
}
`, siteName, siteSlug, deviceRoleName, deviceRoleSlug, mfgName, mfgSlug, deviceTypeModel, deviceTypeSlug, deviceName, serviceName)
}
