package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVirtualDeviceContextDataSource_IDPreservation(t *testing.T) {

	t.Parallel()

	mfgName := testutil.RandomName("mfg-id")
	mfgSlug := testutil.RandomSlug("mfg-id")
	deviceTypeName := testutil.RandomName("device-type-id")
	deviceTypeSlug := testutil.RandomSlug("device-type-id")
	siteName := testutil.RandomName("site-id")
	siteSlug := testutil.RandomSlug("site-id")
	roleName := testutil.RandomName("device-role-id")
	roleSlug := testutil.RandomSlug("device-role-id")
	deviceName := testutil.RandomName("device-id")
	vdcName := testutil.RandomName("vdc-id")
	vdcID := int32(1)

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(mfgSlug)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterDeviceRoleCleanup(roleSlug)
	cleanup.RegisterDeviceCleanup(deviceName)
	cleanup.RegisterVirtualDeviceContextCleanup(vdcID)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckManufacturerDestroy,
			testutil.CheckSiteDestroy,
			testutil.CheckDeviceRoleDestroy,
			testutil.CheckDeviceDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualDeviceContextDataSourceConfig(mfgName, mfgSlug, deviceTypeName, deviceTypeSlug, siteName, siteSlug, roleName, roleSlug, deviceName, vdcName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_virtual_device_context.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_virtual_device_context.test", "name", vdcName),
				),
			},
		},
	})
}

func TestAccVirtualDeviceContextDataSource_byID(t *testing.T) {

	t.Parallel()

	mfgName := testutil.RandomName("mfg")
	mfgSlug := testutil.RandomSlug("mfg")
	deviceTypeName := testutil.RandomName("device-type")
	deviceTypeSlug := testutil.RandomSlug("device-type")
	siteName := testutil.RandomName("site")
	siteSlug := testutil.RandomSlug("site")
	roleName := testutil.RandomName("device-role")
	roleSlug := testutil.RandomSlug("device-role")
	deviceName := testutil.RandomName("device")
	vdcName := testutil.RandomName("vdc")
	vdcID := int32(1) // VDC ID will be auto-assigned

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(mfgSlug)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterDeviceRoleCleanup(roleSlug)
	cleanup.RegisterDeviceCleanup(deviceName)
	cleanup.RegisterVirtualDeviceContextCleanup(vdcID)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckManufacturerDestroy,
			testutil.CheckSiteDestroy,
			testutil.CheckDeviceRoleDestroy,
			testutil.CheckDeviceDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualDeviceContextDataSourceConfig(mfgName, mfgSlug, deviceTypeName, deviceTypeSlug, siteName, siteSlug, roleName, roleSlug, deviceName, vdcName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_virtual_device_context.test", "name", vdcName),
					resource.TestCheckResourceAttrSet("data.netbox_virtual_device_context.test", "device"),
				),
			},
		},
	})
}

func testAccVirtualDeviceContextDataSourceConfig(mfgName, mfgSlug, deviceTypeName, deviceTypeSlug, siteName, siteSlug, roleName, roleSlug, deviceName, vdcName string) string {
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

resource "netbox_site" "test" {
  name = "%s"
  slug = "%s"
}

resource "netbox_device_role" "test" {
  name = "%s"
  slug = "%s"
}

resource "netbox_device" "test" {
  name        = "%s"
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
}

resource "netbox_virtual_device_context" "test" {
  name   = "%s"
  device = netbox_device.test.id
  status = "active"
}

data "netbox_virtual_device_context" "test" {
  id = netbox_virtual_device_context.test.id
}
`, mfgName, mfgSlug, deviceTypeName, deviceTypeSlug, siteName, siteSlug, roleName, roleSlug, deviceName, vdcName)
}
