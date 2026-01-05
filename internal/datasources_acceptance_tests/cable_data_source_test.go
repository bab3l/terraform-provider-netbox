package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCableDataSource_byID(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("site")
	siteSlug := testutil.RandomSlug("site")
	deviceRoleName := testutil.RandomName("role")
	deviceRoleSlug := testutil.RandomSlug("role")
	manufacturerName := testutil.RandomName("mfg")
	manufacturerSlug := testutil.RandomSlug("mfg")
	deviceTypeName := testutil.RandomName("dt")
	deviceTypeSlug := testutil.RandomSlug("dt")
	deviceAName := testutil.RandomName("device-a")
	deviceBName := testutil.RandomName("device-b")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterDeviceRoleCleanup(deviceRoleSlug)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckSiteDestroy,
			testutil.CheckDeviceRoleDestroy,
			testutil.CheckManufacturerDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccCableDataSourceConfig(siteName, siteSlug, deviceRoleName, deviceRoleSlug, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceAName, deviceBName, "cat6"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_cable.test", "type", "cat6"),
					resource.TestCheckResourceAttr("data.netbox_cable.test", "status", "connected"),
				),
			},
		},
	})
}

func TestAccCableDataSource_IDPreservation(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("cable-ds-id-site")
	siteSlug := testutil.RandomSlug("cable-ds-id-site")
	deviceRoleName := testutil.RandomName("cable-ds-id-role")
	deviceRoleSlug := testutil.RandomSlug("cable-ds-id-role")
	manufacturerName := testutil.RandomName("cable-ds-id-mfg")
	manufacturerSlug := testutil.RandomSlug("cable-ds-id-mfg")
	deviceTypeName := testutil.RandomName("cable-ds-id-dt")
	deviceTypeSlug := testutil.RandomSlug("cable-ds-id-dt")
	deviceAName := testutil.RandomName("cable-ds-id-a")
	deviceBName := testutil.RandomName("cable-ds-id-b")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterDeviceRoleCleanup(deviceRoleSlug)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckSiteDestroy,
			testutil.CheckDeviceRoleDestroy,
			testutil.CheckManufacturerDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccCableDataSourceConfig(siteName, siteSlug, deviceRoleName, deviceRoleSlug, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceAName, deviceBName, "cat6"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_cable.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_cable.test", "type", "cat6"),
				),
			},
		},
	})
}

func testAccCableDataSourceConfig(siteName, siteSlug, deviceRoleName, deviceRoleSlug, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceAName, deviceBName, cableType string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name = "%s"
  slug = "%s"
}

resource "netbox_device_role" "test" {
  name = "%s"
  slug = "%s"
}

resource "netbox_manufacturer" "test" {
  name = "%s"
  slug = "%s"
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = "%s"
  slug         = "%s"
}

resource "netbox_device" "test_a" {
  name        = "%s"
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
}

resource "netbox_device" "test_b" {
  name        = "%s"
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
}

resource "netbox_interface" "test_a" {
  device = netbox_device.test_a.id
  name   = "eth0"
  type   = "1000base-t"
}

resource "netbox_interface" "test_b" {
  device = netbox_device.test_b.id
  name   = "eth0"
  type   = "1000base-t"
}

resource "netbox_cable" "test" {
  a_terminations = [
    {
      object_type = "dcim.interface"
      object_id   = netbox_interface.test_a.id
    }
  ]
  b_terminations = [
    {
      object_type = "dcim.interface"
      object_id   = netbox_interface.test_b.id
    }
  ]
  status = "connected"
  type   = "%s"
}

data "netbox_cable" "test" {
  id = netbox_cable.test.id
}
`, siteName, siteSlug, deviceRoleName, deviceRoleSlug, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceAName, deviceBName, cableType)
}
