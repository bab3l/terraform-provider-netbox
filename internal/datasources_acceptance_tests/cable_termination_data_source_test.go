package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCableTerminationDataSource_IDPreservation(t *testing.T) {
	t.Parallel()

	t.Skip("Skipping cable_termination data source test because netbox_cable resource does not export termination IDs")
}

func TestAccCableTerminationDataSource_byID(t *testing.T) {

	t.Parallel()
	t.Skip("Skipping cable_termination data source test because netbox_cable resource does not export termination IDs")

	siteName := testutil.RandomName("test-site-cable-term")
	siteSlug := testutil.GenerateSlug(siteName)
	deviceName := testutil.RandomName("test-device-cable-term")
	interfaceNameA := "eth0"
	interfaceNameB := "eth1"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCableTerminationDataSourceConfig(siteName, siteSlug, deviceName, interfaceNameA, interfaceNameB),
				Check:  resource.ComposeTestCheckFunc(
				// resource.TestCheckResourceAttr("data.netbox_cable_termination.test", "termination_type", "dcim.interface"),
				),
			},
		},
	})
}

func testAccCableTerminationDataSourceConfig(siteName, siteSlug, deviceName, interfaceNameA, interfaceNameB string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name = %q
  slug = %q
  status = "active"
}

resource "netbox_manufacturer" "test" {
  name = "Test Manufacturer Cable Term"
  slug = "test-manufacturer-cable-term"
}

resource "netbox_device_role" "test" {
  name = "Test Device Role Cable Term"
  slug = "test-device-role-cable-term"
}

resource "netbox_device_type" "test" {
  model = "Test Device Type Cable Term"
  slug  = "test-device-type-cable-term"
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_device" "test_a" {
  name           = "%s-a"
  device_type    = netbox_device_type.test.id
  role           = netbox_device_role.test.id
  site           = netbox_site.test.id
}

resource "netbox_device" "test_b" {
  name           = "%s-b"
  device_type    = netbox_device_type.test.id
  role           = netbox_device_role.test.id
  site           = netbox_site.test.id
}

resource "netbox_interface" "test_a" {
  name      = %q
  device    = netbox_device.test_a.id
  type      = "1000base-t"
}

resource "netbox_interface" "test_b" {
  name      = %q
  device    = netbox_device.test_b.id
  type      = "1000base-t"
}

resource "netbox_cable" "test" {
  status = "connected"
  type   = "cat6"
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
}

# data "netbox_cable_termination" "test" {
#   id = netbox_cable.test.a_terminations[0].id
# }
`, siteName, siteSlug, deviceName, deviceName, interfaceNameA, interfaceNameB)
}
