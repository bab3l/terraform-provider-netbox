package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/provider"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccFHRPGroupAssignmentDataSource_byID(t *testing.T) {

	name := acctest.RandomWithPrefix("test-fhrp-assign-ds")
	siteSlug := acctest.RandomWithPrefix("site")
	deviceRoleSlug := acctest.RandomWithPrefix("role")
	manufacturerSlug := acctest.RandomWithPrefix("mfg")
	deviceSlug := acctest.RandomWithPrefix("device")
	interfaceName := acctest.RandomWithPrefix("eth")
	groupID := acctest.RandIntRange(1, 4094)

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterDeviceRoleCleanup(deviceRoleSlug)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
	cleanup.RegisterDeviceCleanup(deviceSlug)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckSiteDestroy,
			testutil.CheckDeviceRoleDestroy,
			testutil.CheckManufacturerDestroy,
			testutil.CheckDeviceDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccFHRPGroupAssignmentDataSourceConfig_byID(name, interfaceName, groupID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_fhrp_group_assignment.test", "interface_type", "dcim.interface"),
					resource.TestCheckResourceAttr("data.netbox_fhrp_group_assignment.test", "priority", "100"),
					resource.TestCheckResourceAttrSet("data.netbox_fhrp_group_assignment.test", "group_id"),
					resource.TestCheckResourceAttrSet("data.netbox_fhrp_group_assignment.test", "interface_id"),
				),
			},
		},
	})
}

func testAccFHRPGroupAssignmentDataSourceConfig_byID(name, interfaceName string, groupID int) string {
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
  model        = "%s-dt"
  slug         = "%s-dt"
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_device_role" "test" {
  name  = "%s-role"
  slug  = "%s-role"
  color = "ff0000"
}

resource "netbox_device" "test" {
  name        = "%s-device"
  site        = netbox_site.test.id
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
}

resource "netbox_interface" "test" {
  name   = "%s"
  device = netbox_device.test.id
  type   = "virtual"
}

resource "netbox_fhrp_group" "test" {
  protocol = "vrrp2"
  group_id = %d
}

resource "netbox_fhrp_group_assignment" "test" {
  group_id       = netbox_fhrp_group.test.id
  interface_type = "dcim.interface"
  interface_id   = netbox_interface.test.id
  priority       = 100
}

data "netbox_fhrp_group_assignment" "test" {
  id = netbox_fhrp_group_assignment.test.id
}
`, name, name, name, name, name, name, name, name, name, interfaceName, groupID)
}
