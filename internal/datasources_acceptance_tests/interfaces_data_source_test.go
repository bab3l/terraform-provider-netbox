package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccInterfacesDataSource_byDeviceAndNameFilters(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("site-ifaces-q")
	siteSlug := testutil.RandomSlug("site-ifaces-q")
	roleName := testutil.RandomName("device-role-ifaces-q")
	roleSlug := testutil.RandomSlug("device-role-ifaces-q")
	mfgName := testutil.RandomName("mfg-ifaces-q")
	mfgSlug := testutil.RandomSlug("mfg-ifaces-q")
	deviceTypeName := testutil.RandomName("device-type-ifaces-q")
	deviceTypeSlug := testutil.RandomSlug("device-type-ifaces-q")
	deviceName := testutil.RandomName("device-ifaces-q")
	interfaceName := testutil.RandomName("eth-ifaces-q")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterDeviceRoleCleanup(roleSlug)
	cleanup.RegisterManufacturerCleanup(mfgSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)
	cleanup.RegisterDeviceCleanup(deviceName)
	cleanup.RegisterInterfaceCleanup(interfaceName, deviceName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckSiteDestroy,
			testutil.CheckDeviceRoleDestroy,
			testutil.CheckManufacturerDestroy,
			testutil.CheckDeviceTypeDestroy,
			testutil.CheckDeviceDestroy,
			testutil.CheckInterfaceDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccInterfacesDataSourceConfig_byDeviceAndName(siteName, siteSlug, roleName, roleSlug, mfgName, mfgSlug, deviceTypeName, deviceTypeSlug, deviceName, interfaceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_interfaces.test", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_interfaces.test", "names.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_interfaces.test", "names.0", interfaceName),
					resource.TestCheckResourceAttrPair("data.netbox_interfaces.test", "ids.0", "netbox_interface.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_interfaces.test", "interfaces.#", "1"),
					resource.TestCheckResourceAttrPair("data.netbox_interfaces.test", "interfaces.0.id", "netbox_interface.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_interfaces.test", "interfaces.0.name", interfaceName),
				),
			},
		},
	})
}

func TestAccInterfacesDataSource_byTagFilter(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("site-ifaces-q-tag")
	siteSlug := testutil.RandomSlug("site-ifaces-q-tag")
	roleName := testutil.RandomName("device-role-ifaces-q-tag")
	roleSlug := testutil.RandomSlug("device-role-ifaces-q-tag")
	mfgName := testutil.RandomName("mfg-ifaces-q-tag")
	mfgSlug := testutil.RandomSlug("mfg-ifaces-q-tag")
	deviceTypeName := testutil.RandomName("device-type-ifaces-q-tag")
	deviceTypeSlug := testutil.RandomSlug("device-type-ifaces-q-tag")
	deviceName := testutil.RandomName("device-ifaces-q-tag")
	interfaceName := testutil.RandomName("eth-ifaces-q-tag")
	tagName := testutil.RandomName("tf-test-tag-ifaces-q")
	tagSlug := testutil.RandomSlug("tf-test-tag-ifaces-q")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterDeviceRoleCleanup(roleSlug)
	cleanup.RegisterManufacturerCleanup(mfgSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)
	cleanup.RegisterDeviceCleanup(deviceName)
	cleanup.RegisterInterfaceCleanup(interfaceName, deviceName)
	cleanup.RegisterTagCleanup(tagSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckSiteDestroy,
			testutil.CheckDeviceRoleDestroy,
			testutil.CheckManufacturerDestroy,
			testutil.CheckDeviceTypeDestroy,
			testutil.CheckDeviceDestroy,
			testutil.CheckInterfaceDestroy,
			testutil.CheckTagDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccInterfacesDataSourceConfig_byTag(siteName, siteSlug, roleName, roleSlug, mfgName, mfgSlug, deviceTypeName, deviceTypeSlug, deviceName, interfaceName, tagName, tagSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_interfaces.test", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_interfaces.test", "names.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_interfaces.test", "names.0", interfaceName),
					resource.TestCheckResourceAttrPair("data.netbox_interfaces.test", "ids.0", "netbox_interface.test", "id"),
				),
			},
		},
	})
}

func TestAccInterfacesDataSource_byNameIcEnabledAndDeviceFilters(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("site-ifaces-q-multi")
	siteSlug := testutil.RandomSlug("site-ifaces-q-multi")
	roleName := testutil.RandomName("device-role-ifaces-q-multi")
	roleSlug := testutil.RandomSlug("device-role-ifaces-q-multi")
	mfgName := testutil.RandomName("mfg-ifaces-q-multi")
	mfgSlug := testutil.RandomSlug("mfg-ifaces-q-multi")
	deviceTypeName := testutil.RandomName("device-type-ifaces-q-multi")
	deviceTypeSlug := testutil.RandomSlug("device-type-ifaces-q-multi")
	deviceName := testutil.RandomName("device-ifaces-q-multi")
	interfaceName := testutil.RandomName("eth-ifaces-q-multi")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterDeviceRoleCleanup(roleSlug)
	cleanup.RegisterManufacturerCleanup(mfgSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)
	cleanup.RegisterDeviceCleanup(deviceName)
	cleanup.RegisterInterfaceCleanup(interfaceName, deviceName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckSiteDestroy,
			testutil.CheckDeviceRoleDestroy,
			testutil.CheckManufacturerDestroy,
			testutil.CheckDeviceTypeDestroy,
			testutil.CheckDeviceDestroy,
			testutil.CheckInterfaceDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccInterfacesDataSourceConfig_byNameIcEnabledAndDevice(siteName, siteSlug, roleName, roleSlug, mfgName, mfgSlug, deviceTypeName, deviceTypeSlug, deviceName, interfaceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_interfaces.test", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_interfaces.test", "names.#", "1"),
					resource.TestCheckResourceAttrPair("data.netbox_interfaces.test", "ids.0", "netbox_interface.test", "id"),
				),
			},
		},
	})
}

func testAccInterfacesDataSourceConfigBase(siteName, siteSlug, roleName, roleSlug, mfgName, mfgSlug, deviceTypeName, deviceTypeSlug, deviceName, interfaceName string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name = %q
  slug = %q
}

resource "netbox_device_role" "test" {
  name = %q
  slug = %q
}

resource "netbox_manufacturer" "test" {
  name = %q
  slug = %q
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = %q
  slug         = %q
}

resource "netbox_device" "test" {
  name        = %q
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
}
`, siteName, siteSlug, roleName, roleSlug, mfgName, mfgSlug, deviceTypeName, deviceTypeSlug, deviceName)
}

func testAccInterfacesDataSourceConfig_byDeviceAndName(siteName, siteSlug, roleName, roleSlug, mfgName, mfgSlug, deviceTypeName, deviceTypeSlug, deviceName, interfaceName string) string {
	return testAccInterfacesDataSourceConfigBase(siteName, siteSlug, roleName, roleSlug, mfgName, mfgSlug, deviceTypeName, deviceTypeSlug, deviceName, interfaceName) + fmt.Sprintf(`

resource "netbox_interface" "test" {
	device = netbox_device.test.id
	name   = %q
	type   = "1000base-t"
}

data "netbox_interfaces" "test" {
  filter {
    name   = "device_id"
    values = [netbox_device.test.id]
  }

  filter {
    name   = "name"
    values = [netbox_interface.test.name]
  }
}
`, interfaceName)
}

func testAccInterfacesDataSourceConfig_byTag(siteName, siteSlug, roleName, roleSlug, mfgName, mfgSlug, deviceTypeName, deviceTypeSlug, deviceName, interfaceName, tagName, tagSlug string) string {
	return fmt.Sprintf(`
resource "netbox_tag" "test" {
  name = %q
  slug = %q
}
`, tagName, tagSlug) + testAccInterfacesDataSourceConfigBase(siteName, siteSlug, roleName, roleSlug, mfgName, mfgSlug, deviceTypeName, deviceTypeSlug, deviceName, interfaceName) + fmt.Sprintf(`

resource "netbox_interface" "test" {
  device = netbox_device.test.id
  name   = %q
  type   = "1000base-t"

	tags = [netbox_tag.test.slug]
}

data "netbox_interfaces" "test" {
  filter {
    name   = "tag"
    values = [netbox_tag.test.slug]
  }

  depends_on = [netbox_interface.test]
}
`, interfaceName)
}

func testAccInterfacesDataSourceConfig_byNameIcEnabledAndDevice(siteName, siteSlug, roleName, roleSlug, mfgName, mfgSlug, deviceTypeName, deviceTypeSlug, deviceName, interfaceName string) string {
	return testAccInterfacesDataSourceConfigBase(siteName, siteSlug, roleName, roleSlug, mfgName, mfgSlug, deviceTypeName, deviceTypeSlug, deviceName, interfaceName) + fmt.Sprintf(`

resource "netbox_interface" "test" {
	device = netbox_device.test.id
	name   = %q
	type   = "1000base-t"
}

data "netbox_interfaces" "test" {
  filter {
    name   = "device_id"
    values = [netbox_device.test.id]
  }

  filter {
    name   = "name__ic"
    values = [upper(netbox_interface.test.name)]
  }

  filter {
    name   = "enabled"
    values = ["true"]
  }
}
`, interfaceName)
}
