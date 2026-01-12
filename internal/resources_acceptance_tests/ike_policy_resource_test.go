package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIKEPolicyResource_basic(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-ike-policy")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIKEPolicyCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIKEPolicyResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ike_policy.test", "id"),
					resource.TestCheckResourceAttr("netbox_ike_policy.test", "name", name),
					resource.TestCheckResourceAttr("netbox_ike_policy.test", "version", "2"),
				),
			},
		},
	})
}

func TestAccIKEPolicyResource_full(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-ike-policy-full")
	proposalName := testutil.RandomName("tf-test-ike-proposal-for-policy")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIKEPolicyCleanup(name)
	cleanup.RegisterIKEProposalCleanup(proposalName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIKEPolicyResourceConfig_full(name, proposalName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ike_policy.test", "id"),
					resource.TestCheckResourceAttr("netbox_ike_policy.test", "name", name),
					resource.TestCheckResourceAttr("netbox_ike_policy.test", "version", "1"),
					resource.TestCheckResourceAttr("netbox_ike_policy.test", "mode", "aggressive"),
					resource.TestCheckResourceAttr("netbox_ike_policy.test", "description", "Test IKE policy with full options"),
					resource.TestCheckResourceAttr("netbox_ike_policy.test", "comments", "Test comments for IKE policy"),
					resource.TestCheckResourceAttr("netbox_ike_policy.test", "proposals.#", "1"),
				),
			},
		},
	})
}

func TestAccIKEPolicyResource_update(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-ike-policy-update")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIKEPolicyCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIKEPolicyResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ike_policy.test", "id"),
					resource.TestCheckResourceAttr("netbox_ike_policy.test", "name", name),
					resource.TestCheckResourceAttr("netbox_ike_policy.test", "version", "2"),
				),
			},
			{
				Config: testAccIKEPolicyResourceConfig_updated(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ike_policy.test", "id"),
					resource.TestCheckResourceAttr("netbox_ike_policy.test", "name", name),
					resource.TestCheckResourceAttr("netbox_ike_policy.test", "version", "2"),
					resource.TestCheckResourceAttr("netbox_ike_policy.test", "description", "Updated description"),
				),
			},
		},
	})
}

func TestAccIKEPolicyResource_import(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-ike-policy")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIKEPolicyCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIKEPolicyResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ike_policy.test", "id"),
					resource.TestCheckResourceAttr("netbox_ike_policy.test", "name", name),
				),
			},
			{
				ResourceName:      "netbox_ike_policy.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccIKEPolicyResource_externalDeletion(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-ike-policy-extdel")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIKEPolicyCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIKEPolicyResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ike_policy.test", "id"),
					resource.TestCheckResourceAttr("netbox_ike_policy.test", "name", name),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					items, _, err := client.VpnAPI.VpnIkePoliciesList(context.Background()).Name([]string{name}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find IKE policy for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.VpnAPI.VpnIkePoliciesDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete IKE policy: %v", err)
					}
					t.Logf("Successfully externally deleted IKE policy with ID: %d", itemID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccIKEPolicyResource_IDPreservation(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-ike-policy-id")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIKEPolicyCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIKEPolicyResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ike_policy.test", "id"),
					resource.TestCheckResourceAttr("netbox_ike_policy.test", "name", name),
				),
			},
		},
	})
}

func testAccIKEPolicyResourceConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "netbox_ike_policy" "test" {
  name    = %q
  version = 2
}
`, name)
}

func testAccIKEPolicyResourceConfig_full(name, proposalName string) string {
	return fmt.Sprintf(`
resource "netbox_ike_proposal" "test" {
  name                     = %q
  authentication_method    = "preshared-keys"
  encryption_algorithm     = "aes-256-cbc"
  authentication_algorithm = "hmac-sha256"
  group                    = 14
}

resource "netbox_ike_policy" "test" {
  name        = %q
  version     = 1
  mode        = "aggressive"
  proposals   = [netbox_ike_proposal.test.id]
  description = "Test IKE policy with full options"
  comments    = "Test comments for IKE policy"
}
`, proposalName, name)
}

func TestAccIKEPolicyResource_removeOptionalFields(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-ike-policy-rem")
	proposalName := testutil.RandomName("tf-test-ike-proposal-rem")
	const testDescription = "Test Description"
	const testComments = "Test Comments"

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIKEPolicyCleanup(name)
	cleanup.RegisterIKEProposalCleanup(proposalName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckIKEPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIKEPolicyResourceConfig_withDescriptionComments(proposalName, name, testDescription, testComments),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_ike_policy.test", "name", name),
					resource.TestCheckResourceAttr("netbox_ike_policy.test", "description", testDescription),
					resource.TestCheckResourceAttr("netbox_ike_policy.test", "comments", testComments),
				),
			},
			{
				Config: testAccIKEPolicyResourceConfig_basicWithProposal(proposalName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_ike_policy.test", "name", name),
					resource.TestCheckNoResourceAttr("netbox_ike_policy.test", "description"),
					resource.TestCheckNoResourceAttr("netbox_ike_policy.test", "comments"),
				),
			},
		},
	})
}

func testAccIKEPolicyResourceConfig_basicWithProposal(proposalName, name string) string {
	return fmt.Sprintf(`
resource "netbox_ike_proposal" "test" {
  name                     = %[1]q
  authentication_method    = "preshared-keys"
  encryption_algorithm     = "aes-128-cbc"
  authentication_algorithm = "hmac-sha256"
  group                    = 14
}

resource "netbox_ike_policy" "test" {
  name      = %[2]q
  version   = 1
  mode      = "aggressive"
  proposals = [netbox_ike_proposal.test.id]
}
`, proposalName, name)
}

func testAccIKEPolicyResourceConfig_withDescriptionComments(proposalName, name, description, comments string) string {
	return fmt.Sprintf(`
resource "netbox_ike_proposal" "test" {
  name                     = %[1]q
  authentication_method    = "preshared-keys"
  encryption_algorithm     = "aes-128-cbc"
  authentication_algorithm = "hmac-sha256"
  group                    = 14
}

resource "netbox_ike_policy" "test" {
  name        = %[2]q
  version     = 1
  mode        = "aggressive"
  proposals   = [netbox_ike_proposal.test.id]
  description = %[3]q
  comments    = %[4]q
}
`, proposalName, name, description, comments)
}

func testAccIKEPolicyResourceConfig_updated(name string) string {
	return fmt.Sprintf(`
resource "netbox_ike_policy" "test" {
  name        = %q
  version     = 2
  description = "Updated description"
}
`, name)
}

func TestAccConsistency_IKEPolicy_LiteralNames(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-ike-policy-lit")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIKEPolicyCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIKEPolicyResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ike_policy.test", "id"),
					resource.TestCheckResourceAttr("netbox_ike_policy.test", "name", name),
					resource.TestCheckResourceAttr("netbox_ike_policy.test", "version", "2"),
				),
			},
			{
				Config:   testAccIKEPolicyResourceConfig_basic(name),
				PlanOnly: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ike_policy.test", "id"),
				),
			},
		},
	})
}
