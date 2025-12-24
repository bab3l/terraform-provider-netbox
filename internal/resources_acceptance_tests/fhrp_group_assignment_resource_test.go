package resources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccFHRPGroupAssignmentResource_basic(t *testing.T) {

	t.Parallel()
	name := testutil.RandomName("test-fhrp-assign")
	interfaceName := testutil.RandomName("eth")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFHRPGroupAssignmentResourceConfig_basic(name, interfaceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_fhrp_group_assignment.test", "interface_type", "dcim.interface"),
					resource.TestCheckResourceAttr("netbox_fhrp_group_assignment.test", "priority", "100"),
					resource.TestCheckResourceAttrSet("netbox_fhrp_group_assignment.test", "id"),
					resource.TestCheckResourceAttrSet("netbox_fhrp_group_assignment.test", "group_id"),
					resource.TestCheckResourceAttrSet("netbox_fhrp_group_assignment.test", "interface_id"),
				),
			},
			{
				Config: testAccFHRPGroupAssignmentResourceConfig_updated(name, interfaceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_fhrp_group_assignment.test", "priority", "200"),
				),
			},
			{
				ResourceName:            "netbox_fhrp_group_assignment.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"group_id", "interface_id", "display_name"},
			},
		},
	})
}

func testAccFHRPGroupAssignmentResourceConfig_basic(name, interfaceName string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name = "%[1]s-site"
  slug = "%[1]s-site"
}

resource "netbox_manufacturer" "test" {
  name = "%[1]s-mfr"
  slug = "%[1]s-mfr"
}

resource "netbox_device_type" "test" {
  model        = "%[1]s-dt"
  slug         = "%[1]s-dt"
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_device_role" "test" {
  name  = "%[1]s-role"
  slug  = "%[1]s-role"
  color = "ff0000"
}

resource "netbox_device" "test" {
  name        = "%[1]s-device"
  site        = netbox_site.test.id
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
}

resource "netbox_interface" "test" {
  name   = %[2]q
  device = netbox_device.test.id
  type   = "virtual"
}

resource "netbox_fhrp_group" "test" {
  protocol = "vrrp2"
  group_id = 1
}

resource "netbox_fhrp_group_assignment" "test" {
  group_id       = netbox_fhrp_group.test.id
  interface_type = "dcim.interface"
  interface_id   = netbox_interface.test.id
  priority       = 100
}
`, name, interfaceName)
}

func testAccFHRPGroupAssignmentResourceConfig_updated(name, interfaceName string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name = "%[1]s-site"
  slug = "%[1]s-site"
}

resource "netbox_manufacturer" "test" {
  name = "%[1]s-mfr"
  slug = "%[1]s-mfr"
}

resource "netbox_device_type" "test" {
  model        = "%[1]s-dt"
  slug         = "%[1]s-dt"
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_device_role" "test" {
  name  = "%[1]s-role"
  slug  = "%[1]s-role"
  color = "ff0000"
}

resource "netbox_device" "test" {
  name        = "%[1]s-device"
  site        = netbox_site.test.id
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
}

resource "netbox_interface" "test" {
  name   = %[2]q
  device = netbox_device.test.id
  type   = "virtual"
}

resource "netbox_fhrp_group" "test" {
  protocol = "vrrp2"
  group_id = 1
}

resource "netbox_fhrp_group_assignment" "test" {
  group_id       = netbox_fhrp_group.test.id
  interface_type = "dcim.interface"
  interface_id   = netbox_interface.test.id
  priority       = 200
}
`, name, interfaceName)
}

func TestAccConsistency_FHRPGroupAssignment_LiteralNames(t *testing.T) {
	t.Parallel()
	name := testutil.RandomName("test-fhrp-assign-lit")
	interfaceName := testutil.RandomName("eth")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFHRPGroupAssignmentConsistencyLiteralNamesConfig(name, interfaceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_fhrp_group_assignment.test", "id"),
					resource.TestCheckResourceAttr("netbox_fhrp_group_assignment.test", "interface_type", "dcim.interface"),
					resource.TestCheckResourceAttr("netbox_fhrp_group_assignment.test", "priority", "100"),
				),
			},
			{
				Config:   testAccFHRPGroupAssignmentConsistencyLiteralNamesConfig(name, interfaceName),
				PlanOnly: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_fhrp_group_assignment.test", "id"),
				),
			},
		},
	})
}

func testAccFHRPGroupAssignmentConsistencyLiteralNamesConfig(name, interfaceName string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name   = "%[1]s-site"
  slug   = "%[1]s-site"
  status = "active"
}

resource "netbox_manufacturer" "test" {
  name = "%[1]s-mfg"
  slug = "%[1]s-mfg"
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = "%[1]s-dt"
  slug         = "%[1]s-dt"
}

resource "netbox_device_role" "test" {
  name = "%[1]s-role"
  slug = "%[1]s-role"
}

resource "netbox_device" "test" {
  name        = "%[1]s-device"
  site        = netbox_site.test.id
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  status      = "active"
}

resource "netbox_interface" "test" {
  name   = %[2]q
  device = netbox_device.test.id
  type   = "virtual"
}

resource "netbox_fhrp_group" "test" {
  protocol = "vrrp2"
  group_id = 1
}

resource "netbox_fhrp_group_assignment" "test" {
  group_id       = netbox_fhrp_group.test.id
  interface_type = "dcim.interface"
  interface_id   = netbox_interface.test.id
  priority       = 100
}
`, name, interfaceName)
}
