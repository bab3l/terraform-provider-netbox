package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/provider"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
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

func TestAccL2VPNTerminationResource_full(t *testing.T) {
	t.Parallel()
	l2vpnName := testutil.RandomName("tf-test-l2vpn-term-full")
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
				Config: testAccL2VPNTerminationResourceConfig_full(l2vpnName, vlanVID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_l2vpn_termination.test", "id"),
					resource.TestCheckResourceAttr("netbox_l2vpn_termination.test", "assigned_object_type", "ipam.vlan"),
					resource.TestCheckResourceAttrSet("netbox_l2vpn_termination.test", "l2vpn"),
					resource.TestCheckResourceAttrSet("netbox_l2vpn_termination.test", "assigned_object_id"),
					resource.TestCheckResourceAttrSet("netbox_l2vpn_termination.test", "display_name"),
				),
			},
		},
	})
}

func testAccL2VPNTerminationResourceConfig_full(l2vpnName string, vlanVID int32) string {
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

func TestAccL2VPNTerminationResource_update(t *testing.T) {
	t.Parallel()
	l2vpnName := testutil.RandomName("tf-test-l2vpn-term-update")
	l2vpnNameUpdated := testutil.RandomName("tf-test-l2vpn-term-update-new")
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
				Config: testAccL2VPNTerminationResourceConfig_updateInitial(l2vpnName, vlanVID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_l2vpn_termination.test", "id"),
					resource.TestCheckResourceAttr("netbox_l2vpn_termination.test", "assigned_object_type", "ipam.vlan"),
				),
			},
			{
				Config: testAccL2VPNTerminationResourceConfig_updateModified(l2vpnNameUpdated, vlanVID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_l2vpn_termination.test", "id"),
					resource.TestCheckResourceAttr("netbox_l2vpn_termination.test", "assigned_object_type", "ipam.vlan"),
				),
			},
		},
	})
}

func testAccL2VPNTerminationResourceConfig_updateInitial(l2vpnName string, vlanVID int32) string {
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

func testAccL2VPNTerminationResourceConfig_updateModified(l2vpnName string, vlanVID int32) string {
	return fmt.Sprintf(`
resource "netbox_l2vpn" "test" {
  name = %q
  slug = %q
  type = "vpws"
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

func TestAccL2VPNTerminationResource_external_deletion(t *testing.T) {
	t.Parallel()
	l2vpnName := acctest.RandomWithPrefix("test-l2vpn-term")
	vlanVID := testutil.RandomVID()

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVLANCleanup(vlanVID)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckVLANDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccL2VPNTerminationResourceConfig_externalDeletion(l2vpnName, vlanVID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_l2vpn_termination.test", "id"),
					resource.TestCheckResourceAttr("netbox_l2vpn_termination.test", "assigned_object_type", "ipam.vlan"),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}

					// Find the L2VPN termination by L2VPN name
					l2vpns, _, err := client.VpnAPI.VpnL2vpnsList(context.Background()).Name([]string{l2vpnName}).Execute()
					if err != nil || l2vpns == nil || len(l2vpns.Results) == 0 {
						t.Fatalf("Failed to find l2vpn for termination: %v", err)
					}
					l2vpnID := l2vpns.Results[0].Id

					// Find the L2VPN termination by L2VPN ID
					items, _, err := client.VpnAPI.VpnL2vpnTerminationsList(context.Background()).L2vpnId([]int32{l2vpnID}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find l2vpn termination for external deletion: %v", err)
					}
					itemID := items.Results[0].Id

					// Delete the L2VPN termination
					_, err = client.VpnAPI.VpnL2vpnTerminationsDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete l2vpn termination: %v", err)
					}
					t.Logf("Successfully externally deleted l2vpn termination with ID: %d", itemID)
				},
				Config: testAccL2VPNTerminationResourceConfig_externalDeletion(l2vpnName, vlanVID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_l2vpn_termination.test", "id"),
				),
			},
		},
	})
}

func testAccL2VPNTerminationResourceConfig_externalDeletion(l2vpnName string, vlanVID int32) string {
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
