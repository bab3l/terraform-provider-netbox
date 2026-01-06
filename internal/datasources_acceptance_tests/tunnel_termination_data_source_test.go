package datasources_acceptance_tests

import (
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Acceptance Tests

func TestAccTunnelTerminationDataSource_IDPreservation(t *testing.T) {

	t.Parallel()
	tunnelName := testutil.RandomName("tf-test-tunnel-term-ds-id")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTunnelCleanup(tunnelName)
	cleanup.RegisterTunnelTerminationCleanup(tunnelName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckTunnelTerminationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTunnelTerminationDataSourceConfig_byID(tunnelName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_tunnel_termination.test", "id"),
					resource.TestCheckResourceAttrPair("data.netbox_tunnel_termination.test", "tunnel", "netbox_tunnel_termination.test", "tunnel"),
				),
			},
		},
	})
}

func TestAccTunnelTerminationDataSource_byID(t *testing.T) {

	t.Parallel()
	tunnelName := testutil.RandomName("tf-test-tunnel-term-ds")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTunnelCleanup(tunnelName)
	cleanup.RegisterTunnelTerminationCleanup(tunnelName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckTunnelTerminationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTunnelTerminationDataSourceConfig_byID(tunnelName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.netbox_tunnel_termination.test", "id", "netbox_tunnel_termination.test", "id"),
					resource.TestCheckResourceAttrPair("data.netbox_tunnel_termination.test", "tunnel", "netbox_tunnel_termination.test", "tunnel"),
					resource.TestCheckResourceAttr("data.netbox_tunnel_termination.test", "termination_type", "dcim.device"),
					resource.TestCheckResourceAttr("data.netbox_tunnel_termination.test", "role", "peer"),
				),
			},
		},
	})
}

func TestAccTunnelTerminationDataSource_byTunnel(t *testing.T) {

	t.Parallel()
	tunnelName := testutil.RandomName("tf-test-tunnel-term-ds2")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTunnelCleanup(tunnelName)
	cleanup.RegisterTunnelTerminationCleanup(tunnelName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckTunnelTerminationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTunnelTerminationDataSourceConfig_byTunnel(tunnelName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.netbox_tunnel_termination.test", "tunnel", "netbox_tunnel.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_tunnel_termination.test", "termination_type", "dcim.device"),
				),
			},
		},
	})
}

func testAccTunnelTerminationDataSourceConfig_byID(tunnelName string) string {
	return `
resource "netbox_tunnel" "test" {
  name          = "` + tunnelName + `"
  encapsulation = "ipsec-tunnel"
}

resource "netbox_tunnel_termination" "test" {
  tunnel           = netbox_tunnel.test.id
  termination_type = "dcim.device"
  role             = "peer"
}

data "netbox_tunnel_termination" "test" {
  id = netbox_tunnel_termination.test.id
}
`
}

func testAccTunnelTerminationDataSourceConfig_byTunnel(tunnelName string) string {
	return `
resource "netbox_tunnel" "test" {
  name          = "` + tunnelName + `"
  encapsulation = "ipsec-tunnel"
}

resource "netbox_tunnel_termination" "test" {
  tunnel           = netbox_tunnel.test.id
  termination_type = "dcim.device"
  role             = "peer"
}

data "netbox_tunnel_termination" "test" {
  tunnel = netbox_tunnel.test.id
  depends_on = [netbox_tunnel_termination.test]
}
`
}
