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

func TestAccL2VPNTerminationDataSource_IDPreservation(t *testing.T) {
	name := acctest.RandomWithPrefix("test-l2vpn-term-ds-id")
	siteSlug := acctest.RandomWithPrefix("site-id")
	deviceRoleSlug := acctest.RandomWithPrefix("role-id")
	manufacturerSlug := acctest.RandomWithPrefix("mfg-id")
	deviceSlug := acctest.RandomWithPrefix("device-id")

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
				Config: testAccL2VPNTerminationDataSourceConfig_byID(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_l2vpn_termination.test", "id"),
				),
			},
		},
	})
}

func TestAccL2VPNTerminationDataSource_byID(t *testing.T) {
	name := acctest.RandomWithPrefix("test-l2vpn-term-ds")
	siteSlug := acctest.RandomWithPrefix("site")
	deviceRoleSlug := acctest.RandomWithPrefix("role")
	manufacturerSlug := acctest.RandomWithPrefix("mfg")
	deviceSlug := acctest.RandomWithPrefix("device")

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
				Config: testAccL2VPNTerminationDataSourceConfig_byID(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_l2vpn_termination.test", "assigned_object_type", "ipam.vlan"),
					resource.TestCheckResourceAttrSet("data.netbox_l2vpn_termination.test", "l2vpn"),
					resource.TestCheckResourceAttrSet("data.netbox_l2vpn_termination.test", "assigned_object_id"),
				),
			},
		},
	})
}

func testAccL2VPNTerminationDataSourceConfig_byID(name string) string {
	return fmt.Sprintf(`
resource "netbox_l2vpn" "test" {
  name = "%s"
  slug = "%s"
  type = "vxlan"
}

resource "netbox_vlan" "test" {
  name    = "%s-vlan"
  vid     = 100
}

resource "netbox_l2vpn_termination" "test" {
  l2vpn                = netbox_l2vpn.test.id
  assigned_object_type = "ipam.vlan"
  assigned_object_id   = netbox_vlan.test.id
}

data "netbox_l2vpn_termination" "test" {
  id = netbox_l2vpn_termination.test.id
}
`, name, name, name)
}
