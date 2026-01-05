package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIPSECProposalResource_basic(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-ipsec-prop")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIPSECProposalResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ipsec_proposal.test", "id"),
					resource.TestCheckResourceAttr("netbox_ipsec_proposal.test", "name", name),
				),
			},
		},
	})
}

func TestAccIPSECProposalResource_full(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-ipsec-prop-full")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIPSECProposalResourceConfig_full(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ipsec_proposal.test", "id"),
					resource.TestCheckResourceAttr("netbox_ipsec_proposal.test", "name", name),
					resource.TestCheckResourceAttr("netbox_ipsec_proposal.test", "encryption_algorithm", "aes-128-cbc"),
					resource.TestCheckResourceAttr("netbox_ipsec_proposal.test", "authentication_algorithm", "hmac-sha256"),
					resource.TestCheckResourceAttr("netbox_ipsec_proposal.test", "description", "Test IPsec proposal"),
				),
			},
		},
	})
}

func TestAccIPSECProposalResource_update(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-ipsec-prop-update")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIPSECProposalResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ipsec_proposal.test", "id"),
					resource.TestCheckResourceAttr("netbox_ipsec_proposal.test", "name", name),
				),
			},
			{
				Config: testAccIPSECProposalResourceConfig_full(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ipsec_proposal.test", "id"),
					resource.TestCheckResourceAttr("netbox_ipsec_proposal.test", "name", name),
					resource.TestCheckResourceAttr("netbox_ipsec_proposal.test", "encryption_algorithm", "aes-128-cbc"),
					resource.TestCheckResourceAttr("netbox_ipsec_proposal.test", "authentication_algorithm", "hmac-sha256"),
					resource.TestCheckResourceAttr("netbox_ipsec_proposal.test", "description", "Test IPsec proposal"),
				),
			},
		},
	})
}

func TestAccIPSECProposalResource_import(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-ipsec-prop")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIPSECProposalResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ipsec_proposal.test", "id"),
				),
			},
			{
				ResourceName:      "netbox_ipsec_proposal.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccIPSECProposalResource_externalDeletion(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-ipsec-prop-extdel")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIPSecProposalCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIPSECProposalResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ipsec_proposal.test", "id"),
					resource.TestCheckResourceAttr("netbox_ipsec_proposal.test", "name", name),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					items, _, err := client.VpnAPI.VpnIpsecProposalsList(context.Background()).Name([]string{name}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find IPSec proposal for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.VpnAPI.VpnIpsecProposalsDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete IPSec proposal: %v", err)
					}
					t.Logf("Successfully externally deleted IPSec proposal with ID: %d", itemID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}
func TestAccIPSecProposalResource_IDPreservation(t *testing.T) {
	t.Parallel()
	name := testutil.RandomName("tf-test-ipsec-proposal-id")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIPSECProposalResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ipsec_proposal.test", "id"),
					resource.TestCheckResourceAttr("netbox_ipsec_proposal.test", "name", name),
				),
			},
		},
	})
}
func testAccIPSECProposalResourceConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "netbox_ipsec_proposal" "test" {
  name                     = %q
  encryption_algorithm     = "aes-128-cbc"
}
`, name)
}

func testAccIPSECProposalResourceConfig_full(name string) string {
	return fmt.Sprintf(`
resource "netbox_ipsec_proposal" "test" {
  name                     = %q
  encryption_algorithm     = "aes-128-cbc"
  authentication_algorithm = "hmac-sha256"
  description              = "Test IPsec proposal"
}
`, name)
}

func TestAccConsistency_IPSECProposal_LiteralNames(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-ipsec-proposal-lit")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIPSECProposalResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ipsec_proposal.test", "id"),
					resource.TestCheckResourceAttr("netbox_ipsec_proposal.test", "name", name),
				),
			},
			{
				Config:   testAccIPSECProposalResourceConfig_basic(name),
				PlanOnly: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ipsec_proposal.test", "id"),
				),
			},
		},
	})
}
