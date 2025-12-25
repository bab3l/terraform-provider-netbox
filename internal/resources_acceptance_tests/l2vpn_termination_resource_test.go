package resources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/provider"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccL2VPNTerminationResource_basic(t *testing.T) {

	t.Parallel()
	l2vpnName := testutil.RandomName("tf-test-l2vpn-term")
	vlanVID := testutil.RandomVID()

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVLANCleanup(vlanVID)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckVLANDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccL2VPNTerminationResourceConfig_basic(l2vpnName, vlanVID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_l2vpn_termination.test", "id"),
					resource.TestCheckResourceAttr("netbox_l2vpn_termination.test", "assigned_object_type", "ipam.vlan"),
					resource.TestCheckResourceAttrSet("netbox_l2vpn_termination.test", "l2vpn"),
					resource.TestCheckResourceAttrSet("netbox_l2vpn_termination.test", "assigned_object_id"),
				),
			},
		},
	})
}

func TestAccL2VPNTerminationResource_IDPreservation(t *testing.T) {
	t.Parallel()

	l2vpnName := testutil.RandomName("tf-test-l2vpn-term-id")
	vlanVID := testutil.RandomVID()

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVLANCleanup(vlanVID)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckVLANDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccL2VPNTerminationResourceConfig_basic(l2vpnName, vlanVID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_l2vpn_termination.test", "id"),
					resource.TestCheckResourceAttr("netbox_l2vpn_termination.test", "assigned_object_type", "ipam.vlan"),
				),
			},
		},
	})
}

func testAccL2VPNTerminationResourceConfig_basic(l2vpnName string, vlanVID int32) string {
	return fmt.Sprintf(`
resource "netbox_l2vpn" "test" {
  name = %q
  slug = %q
  type = "vxlan"
}

resource "netbox_vlan" "test" {
  name = %q
  vid  = %d
}

resource "netbox_l2vpn_termination" "test" {
  l2vpn                = netbox_l2vpn.test.id
  assigned_object_type = "ipam.vlan"
  assigned_object_id   = netbox_vlan.test.id
}
`, l2vpnName, l2vpnName, l2vpnName, vlanVID)
}

func TestAccConsistency_L2VPNTermination_LiteralNames(t *testing.T) {

	t.Parallel()
	l2vpnName := testutil.RandomName("tf-test-l2vpn-lit")
	vlanVID := testutil.RandomVID()

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVLANCleanup(vlanVID)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckVLANDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccL2VPNTerminationConsistencyLiteralNamesConfig(l2vpnName, vlanVID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_l2vpn_termination.test", "id"),
					resource.TestCheckResourceAttr("netbox_l2vpn_termination.test", "assigned_object_type", "ipam.vlan"),
				),
			},
			{
				Config:   testAccL2VPNTerminationConsistencyLiteralNamesConfig(l2vpnName, vlanVID),
				PlanOnly: true,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_l2vpn_termination.test", "id"),
				),
			},
		},
	})
}

func testAccL2VPNTerminationConsistencyLiteralNamesConfig(l2vpnName string, vlanVID int32) string {
	return fmt.Sprintf(`
resource "netbox_l2vpn" "test" {
  name = %q
  slug = %q
  type = "vxlan"
}

resource "netbox_vlan" "test" {
  name = %q
  vid  = %d
}

resource "netbox_l2vpn_termination" "test" {
  l2vpn                = netbox_l2vpn.test.id
  assigned_object_type = "ipam.vlan"
  assigned_object_id   = netbox_vlan.test.id
}
`, l2vpnName, l2vpnName, l2vpnName, vlanVID)
}
