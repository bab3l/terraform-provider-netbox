package resources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccFHRPGroupAssignmentResource_basic(t *testing.T) {
	name := testutil.RandomName("test-fhrp-assign")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFHRPGroupAssignmentResourceConfig_basic(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_fhrp_group_assignment.test", "interface_type", "dcim.interface"),
					resource.TestCheckResourceAttr("netbox_fhrp_group_assignment.test", "priority", "100"),
					resource.TestCheckResourceAttrSet("netbox_fhrp_group_assignment.test", "id"),
					resource.TestCheckResourceAttrSet("netbox_fhrp_group_assignment.test", "group_id"),
					resource.TestCheckResourceAttrSet("netbox_fhrp_group_assignment.test", "interface_id"),
				),
			},
			{
				Config: testAccFHRPGroupAssignmentResourceConfig_updated(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_fhrp_group_assignment.test", "priority", "200"),
				),
			},
			{
				ResourceName:            "netbox_fhrp_group_assignment.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"group_id", "interface_id"},
			},
		},
	})
}

func testAccFHRPGroupAssignmentResourceConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name = "%s-site"
  slug = "%s-site"
}

resource "netbox_manufacturer" "test" {
  name = "%s-mfr"
  slug = "%s-mfr"
}

resource "netbox_device_type" "test" {
  model           = "%s-dt"
  slug            = "%s-dt"
  manufacturer_id = netbox_manufacturer.test.id
}

resource "netbox_device_role" "test" {
  name  = "%s-role"
  slug  = "%s-role"
  color = "ff0000"
}

resource "netbox_device" "test" {
  name           = "%s-device"
  site_id        = netbox_site.test.id
  device_type_id = netbox_device_type.test.id
  role_id        = netbox_device_role.test.id
}

resource "netbox_interface" "test" {
  name      = "eth0"
  device_id = netbox_device.test.id
  type      = "virtual"
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
`, name, name, name, name, name, name, name, name, name)
}

func testAccFHRPGroupAssignmentResourceConfig_updated(name string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name = "%s-site"
  slug = "%s-site"
}

resource "netbox_manufacturer" "test" {
  name = "%s-mfr"
  slug = "%s-mfr"
}

resource "netbox_device_type" "test" {
  model           = "%s-dt"
  slug            = "%s-dt"
  manufacturer_id = netbox_manufacturer.test.id
}

resource "netbox_device_role" "test" {
  name  = "%s-role"
  slug  = "%s-role"
  color = "ff0000"
}

resource "netbox_device" "test" {
  name           = "%s-device"
  site_id        = netbox_site.test.id
  device_type_id = netbox_device_type.test.id
  role_id        = netbox_device_role.test.id
}

resource "netbox_interface" "test" {
  name      = "eth0"
  device_id = netbox_device.test.id
  type      = "virtual"
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
`, name, name, name, name, name, name, name, name, name)
}
