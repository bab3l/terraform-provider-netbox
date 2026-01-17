package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDevicesDataSource_byNameFilter(t *testing.T) {
	t.Parallel()

	deviceName := testutil.RandomName("tf-test-device-q")
	manufacturerName := testutil.RandomName("tf-test-manufacturer")
	manufacturerSlug := testutil.RandomSlug("tf-test-mfr")
	deviceTypeModel := testutil.RandomName("tf-test-device-type")
	deviceTypeSlug := testutil.RandomSlug("tf-test-dt")
	deviceRoleName := testutil.RandomName("tf-test-device-role")
	deviceRoleSlug := testutil.RandomSlug("tf-test-dr")
	siteName := testutil.RandomName("tf-test-site")
	siteSlug := testutil.RandomSlug("tf-test-site")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterDeviceCleanup(deviceName)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)
	cleanup.RegisterDeviceRoleCleanup(deviceRoleSlug)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckDeviceDestroy,
			testutil.CheckDeviceTypeDestroy,
			testutil.CheckDeviceRoleDestroy,
			testutil.CheckManufacturerDestroy,
			testutil.CheckSiteDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccDevicesDataSourceConfig_byName(deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_devices.test", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_devices.test", "names.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_devices.test", "names.0", deviceName),
					resource.TestCheckResourceAttrPair("data.netbox_devices.test", "ids.0", "netbox_device.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_devices.test", "devices.#", "1"),
					resource.TestCheckResourceAttrPair("data.netbox_devices.test", "devices.0.id", "netbox_device.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_devices.test", "devices.0.name", deviceName),
				),
			},
		},
	})
}

func TestAccDevicesDataSource_byTagFilter(t *testing.T) {
	t.Parallel()

	deviceName := testutil.RandomName("tf-test-device-q-tag")
	manufacturerName := testutil.RandomName("tf-test-manufacturer-tag")
	manufacturerSlug := testutil.RandomSlug("tf-test-mfr-tag")
	deviceTypeModel := testutil.RandomName("tf-test-device-type-tag")
	deviceTypeSlug := testutil.RandomSlug("tf-test-dt-tag")
	deviceRoleName := testutil.RandomName("tf-test-device-role-tag")
	deviceRoleSlug := testutil.RandomSlug("tf-test-dr-tag")
	siteName := testutil.RandomName("tf-test-site-tag")
	siteSlug := testutil.RandomSlug("tf-test-site-tag")
	tagName := testutil.RandomName("tf-test-tag")
	tagSlug := testutil.RandomSlug("tf-test-tag")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterDeviceCleanup(deviceName)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)
	cleanup.RegisterDeviceRoleCleanup(deviceRoleSlug)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterTagCleanup(tagSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckDeviceDestroy,
			testutil.CheckDeviceTypeDestroy,
			testutil.CheckDeviceRoleDestroy,
			testutil.CheckManufacturerDestroy,
			testutil.CheckSiteDestroy,
			testutil.CheckTagDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccDevicesDataSourceConfig_byTag(deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug, tagName, tagSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_devices.test", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_devices.test", "names.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_devices.test", "names.0", deviceName),
					resource.TestCheckResourceAttrPair("data.netbox_devices.test", "ids.0", "netbox_device.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_devices.test", "devices.#", "1"),
					resource.TestCheckResourceAttrPair("data.netbox_devices.test", "devices.0.id", "netbox_device.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_devices.test", "devices.0.name", deviceName),
				),
			},
		},
	})
}

func TestAccDevicesDataSource_byNameAndTagFilters(t *testing.T) {
	t.Parallel()

	deviceName := testutil.RandomName("tf-test-device-q-multi")
	manufacturerName := testutil.RandomName("tf-test-manufacturer-multi")
	manufacturerSlug := testutil.RandomSlug("tf-test-mfr-multi")
	deviceTypeModel := testutil.RandomName("tf-test-device-type-multi")
	deviceTypeSlug := testutil.RandomSlug("tf-test-dt-multi")
	deviceRoleName := testutil.RandomName("tf-test-device-role-multi")
	deviceRoleSlug := testutil.RandomSlug("tf-test-dr-multi")
	siteName := testutil.RandomName("tf-test-site-multi")
	siteSlug := testutil.RandomSlug("tf-test-site-multi")
	tagName := testutil.RandomName("tf-test-tag-multi")
	tagSlug := testutil.RandomSlug("tf-test-tag-multi")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterDeviceCleanup(deviceName)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)
	cleanup.RegisterDeviceRoleCleanup(deviceRoleSlug)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterTagCleanup(tagSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckDeviceDestroy,
			testutil.CheckDeviceTypeDestroy,
			testutil.CheckDeviceRoleDestroy,
			testutil.CheckManufacturerDestroy,
			testutil.CheckSiteDestroy,
			testutil.CheckTagDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccDevicesDataSourceConfig_byNameAndTag(deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug, tagName, tagSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_devices.test", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_devices.test", "names.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_devices.test", "names.0", deviceName),
					resource.TestCheckResourceAttrPair("data.netbox_devices.test", "ids.0", "netbox_device.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_devices.test", "devices.#", "1"),
					resource.TestCheckResourceAttrPair("data.netbox_devices.test", "devices.0.id", "netbox_device.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_devices.test", "devices.0.name", deviceName),
				),
			},
		},
	})
}

func testAccDevicesDataSourceConfig_byName(deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug string) string {
	return fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.slug
  model        = %[3]q
  slug         = %[4]q
}

resource "netbox_device_role" "test" {
  name = %[5]q
  slug = %[6]q
}

resource "netbox_site" "test" {
  name   = %[7]q
  slug   = %[8]q
  status = "active"
}

resource "netbox_device" "test" {
  name        = %[9]q
  device_type = netbox_device_type.test.slug
  role        = netbox_device_role.test.slug
  site        = netbox_site.test.slug
}

data "netbox_devices" "test" {
  filter {
    name   = "name"
    values = [netbox_device.test.name]
  }
}
`, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug, deviceName)
}

func testAccDevicesDataSourceConfig_byTag(deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug, tagName, tagSlug string) string {
	return fmt.Sprintf(`
resource "netbox_tag" "test" {
	name = %[1]q
	slug = %[2]q
}

resource "netbox_manufacturer" "test" {
	name = %[3]q
	slug = %[4]q
}

resource "netbox_device_type" "test" {
	manufacturer = netbox_manufacturer.test.slug
	model        = %[5]q
	slug         = %[6]q
}

resource "netbox_device_role" "test" {
	name = %[7]q
	slug = %[8]q
}

resource "netbox_site" "test" {
	name   = %[9]q
	slug   = %[10]q
	status = "active"
}

resource "netbox_device" "test" {
	name        = %[11]q
	device_type = netbox_device_type.test.slug
	role        = netbox_device_role.test.slug
	site        = netbox_site.test.slug

	tags = [netbox_tag.test.slug]
}

data "netbox_devices" "test" {
	filter {
		name   = "tag"
		values = [netbox_tag.test.slug]
	}

	depends_on = [netbox_device.test]
}
`, tagName, tagSlug, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug, deviceName)
}

func testAccDevicesDataSourceConfig_byNameAndTag(deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug, tagName, tagSlug string) string {
	return fmt.Sprintf(`
resource "netbox_tag" "test" {
	name = %[1]q
	slug = %[2]q
}

resource "netbox_manufacturer" "test" {
	name = %[3]q
	slug = %[4]q
}

resource "netbox_device_type" "test" {
	manufacturer = netbox_manufacturer.test.slug
	model        = %[5]q
	slug         = %[6]q
}

resource "netbox_device_role" "test" {
	name = %[7]q
	slug = %[8]q
}

resource "netbox_site" "test" {
	name   = %[9]q
	slug   = %[10]q
	status = "active"
}

resource "netbox_device" "test" {
	name        = %[11]q
	device_type = netbox_device_type.test.slug
	role        = netbox_device_role.test.slug
	site        = netbox_site.test.slug

	tags = [netbox_tag.test.slug]
}

data "netbox_devices" "test" {
	filter {
		name   = "name"
		values = [netbox_device.test.name]
	}

	filter {
		name   = "tag"
		values = [netbox_tag.test.slug]
	}
}
`, tagName, tagSlug, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug, deviceName)
}
