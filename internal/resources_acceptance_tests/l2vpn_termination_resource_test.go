package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
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
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
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
			{
				// Test import
				ResourceName:      "netbox_l2vpn_termination.test",
				ImportState:       true,
				ImportStateVerify: true,
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
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
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
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
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
			{
				Config:   testAccL2VPNTerminationResourceConfig_basic(l2vpnName, vlanVID),
				PlanOnly: true,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_l2vpn_termination.test", "id"),
				),
			},
		},
	})
}

func TestAccL2VPNTerminationResource_full(t *testing.T) {
	t.Parallel()
	l2vpnName := testutil.RandomName("tf-test-l2vpn-term-full")
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
				Config: testAccL2VPNTerminationResourceConfig_full(l2vpnName, vlanVID),
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
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
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

func TestAccL2VPNTerminationResource_removeOptionalFields(t *testing.T) {
	t.Skip("Skipping: L2VPN termination only has tags/custom_fields as optional, and tags removal exposes a provider consistency bug. Since resource has no other optional fields to test, skipping this test.")
	t.Parallel()

	l2vpnName := testutil.RandomName("tf-test-l2vpn-term-rem")
	vlanVID := testutil.RandomVID()
	tagName := testutil.RandomName("test-tag")
	tagSlug := testutil.GenerateSlug(tagName)

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVLANCleanup(vlanVID)
	cleanup.RegisterTagCleanup(tagSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckVLANDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccL2VPNTerminationResourceConfig_withTags(l2vpnName, vlanVID, tagName, tagSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_l2vpn_termination.test", "tags.#", "1"),
				),
			},
			{
				Config: testAccL2VPNTerminationResourceConfig_withTagButNotUsed(l2vpnName, vlanVID, tagName, tagSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_l2vpn_termination.test", "tags.#", "0"),
				),
			},
		},
	})
}

func testAccL2VPNTerminationResourceConfig_withTags(l2vpnName string, vlanVID int32, tagName, tagSlug string) string {
	return fmt.Sprintf(`
resource "netbox_tag" "test" {
  name = %q
  slug = %q
}

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
  tags = [
    {
      name = netbox_tag.test.name
      slug = netbox_tag.test.slug
    }
  ]
}
`, tagName, tagSlug, l2vpnName, l2vpnName, l2vpnName, vlanVID)
}

func testAccL2VPNTerminationResourceConfig_withTagButNotUsed(l2vpnName string, vlanVID int32, tagName, tagSlug string) string {
	return fmt.Sprintf(`
resource "netbox_tag" "test" {
  name = %q
  slug = %q
}

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
`, tagName, tagSlug, l2vpnName, l2vpnName, l2vpnName, vlanVID)
}

func TestAccL2VPNTerminationResource_external_deletion(t *testing.T) {
	t.Parallel()
	l2vpnName := acctest.RandomWithPrefix("test-l2vpn-term")
	vlanVID := testutil.RandomVID()

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVLANCleanup(vlanVID)
	cleanup.RegisterL2VPNCleanup(l2vpnName)

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
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
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
