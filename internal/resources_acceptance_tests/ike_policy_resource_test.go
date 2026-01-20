package resources_acceptance_tests

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var _ = testAccIKEPolicyResourceConfig_withoutVersion

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
			{
				Config:   testAccIKEPolicyResourceConfig_basic(name),
				PlanOnly: true,
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
	const testPresharedKey = "test-preshared-key"

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIKEPolicyCleanup(name)
	cleanup.RegisterIKEProposalCleanup(proposalName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckIKEPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIKEPolicyResourceConfig_withOptionalFields(proposalName, name, testDescription, testComments, testPresharedKey),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_ike_policy.test", "name", name),
					resource.TestCheckResourceAttr("netbox_ike_policy.test", "version", "1"),
					resource.TestCheckResourceAttr("netbox_ike_policy.test", "mode", "aggressive"),
					resource.TestCheckResourceAttr("netbox_ike_policy.test", "description", testDescription),
					resource.TestCheckResourceAttr("netbox_ike_policy.test", "comments", testComments),
					resource.TestCheckResourceAttr("netbox_ike_policy.test", "proposals.#", "1"),
					resource.TestCheckResourceAttr("netbox_ike_policy.test", "preshared_key", testPresharedKey),
				),
			},
			{
				Config: testAccIKEPolicyResourceConfig_withoutOptionalFields(proposalName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_ike_policy.test", "name", name),
					resource.TestCheckNoResourceAttr("netbox_ike_policy.test", "description"),
					resource.TestCheckNoResourceAttr("netbox_ike_policy.test", "comments"),
					resource.TestCheckResourceAttr("netbox_ike_policy.test", "mode", "aggressive"),
					resource.TestCheckNoResourceAttr("netbox_ike_policy.test", "proposals"),
					resource.TestCheckNoResourceAttr("netbox_ike_policy.test", "preshared_key"),
				),
			},
			{
				Config:      testAccIKEPolicyResourceConfig_versionOnly(proposalName, name, 2),
				ExpectError: regexp.MustCompile("Cannot change IKE policy to version=2"),
			},
		},
	})
}

func testAccIKEPolicyResourceConfig_withoutOptionalFields(proposalName, name string) string {
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
}
`, proposalName, name)
}

func testAccIKEPolicyResourceConfig_withOptionalFields(proposalName, name, description, comments, presharedKey string) string {
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
	preshared_key = %[5]q
}
`, proposalName, name, description, comments, presharedKey)
}

func testAccIKEPolicyResourceConfig_versionOnly(proposalName, name string, version int) string {
	return fmt.Sprintf(`
resource "netbox_ike_proposal" "test" {
	name                     = %[1]q
	authentication_method    = "preshared-keys"
	encryption_algorithm     = "aes-128-cbc"
	authentication_algorithm = "hmac-sha256"
	group                    = 14
}

resource "netbox_ike_policy" "test" {
	name    = %[2]q
	version = %[3]d
}
`, proposalName, name, version)
}

func testAccIKEPolicyResourceConfig_withoutVersion(proposalName, name string) string {
	return fmt.Sprintf(`
resource "netbox_ike_proposal" "test" {
	name                     = %[1]q
	authentication_method    = "preshared-keys"
	encryption_algorithm     = "aes-128-cbc"
	authentication_algorithm = "hmac-sha256"
	group                    = 14
}

resource "netbox_ike_policy" "test" {
	name = %[2]q
}
`, proposalName, name)
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

func TestAccIKEPolicyResource_tagLifecycle(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-ike-policy-tags")
	slug1 := testutil.RandomSlug("tag1")
	slug2 := testutil.RandomSlug("tag2")
	slug3 := testutil.RandomSlug("tag3")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIKEPolicyCleanup(name)
	cleanup.RegisterTagCleanup(slug1)
	cleanup.RegisterTagCleanup(slug2)
	cleanup.RegisterTagCleanup(slug3)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIKEPolicyResourceConfig_tags(name, slug1, slug2, caseTag1Tag2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_ike_policy.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemAttr("netbox_ike_policy.test", "tags.*", slug1),
					resource.TestCheckTypeSetElemAttr("netbox_ike_policy.test", "tags.*", slug2),
				),
			},
			{
				Config: testAccIKEPolicyResourceConfig_tags(name, slug1, slug2, caseTag1Uscore2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_ike_policy.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemAttr("netbox_ike_policy.test", "tags.*", slug1),
					resource.TestCheckTypeSetElemAttr("netbox_ike_policy.test", "tags.*", slug2),
				),
			},
			{
				Config: testAccIKEPolicyResourceConfig_tags(name, slug1, slug2, caseTag3),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_ike_policy.test", "tags.#", "1"),
					resource.TestCheckTypeSetElemAttr("netbox_ike_policy.test", "tags.*", "tag3"),
				),
			},
			{
				Config: testAccIKEPolicyResourceConfig_tags(name, slug1, slug2, tagsEmpty),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_ike_policy.test", "tags.#", "0"),
				),
			},
		},
	})
}

func TestAccIKEPolicyResource_tagOrderInvariance(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-ike-policy-tag-order")
	slug1 := testutil.RandomSlug("tag1")
	slug2 := testutil.RandomSlug("tag2")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIKEPolicyCleanup(name)
	cleanup.RegisterTagCleanup(slug1)
	cleanup.RegisterTagCleanup(slug2)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIKEPolicyResourceConfig_tagsOrder(name, slug1, slug2, caseTag1Tag2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_ike_policy.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemAttr("netbox_ike_policy.test", "tags.*", slug1),
					resource.TestCheckTypeSetElemAttr("netbox_ike_policy.test", "tags.*", slug2),
				),
			},
			{
				Config: testAccIKEPolicyResourceConfig_tagsOrder(name, slug1, slug2, caseTag2Uscore1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_ike_policy.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemAttr("netbox_ike_policy.test", "tags.*", slug1),
					resource.TestCheckTypeSetElemAttr("netbox_ike_policy.test", "tags.*", slug2),
				),
			},
		},
	})
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

func testAccIKEPolicyResourceConfig_tags(name, slug1, slug2, tagCase string) string {
	var tagsConfig string
	switch tagCase {
	case caseTag1Tag2:
		tagsConfig = tagsDoubleSlug
	case caseTag1Uscore2:
		tagsConfig = tagsDoubleSlug
	case caseTag3:
		tagsConfig = tagsSingleSlug
	case tagsEmpty:
		tagsConfig = tagsEmpty
	}

	return fmt.Sprintf(`
resource "netbox_tag" "tag1" {
  name = "Tag1-%[1]s"
  slug = %[1]q
}

resource "netbox_tag" "tag2" {
  name = "Tag2-%[2]s"
  slug = %[2]q
}

resource "netbox_tag" "tag3" {
  name = "Tag3"
  slug = "tag3"
}

resource "netbox_ike_policy" "test" {
  name    = %[3]q
  version = 2
  %[4]s
}
`, slug1, slug2, name, tagsConfig)
}

func testAccIKEPolicyResourceConfig_tagsOrder(name, slug1, slug2, tagCase string) string {
	var tagsConfig string
	switch tagCase {
	case caseTag1Tag2:
		tagsConfig = tagsDoubleSlug
	case caseTag2Uscore1:
		tagsConfig = tagsDoubleSlugReversed
	}

	return fmt.Sprintf(`
resource "netbox_tag" "tag1" {
  name = "Tag1-%[1]s"
  slug = %[1]q
}

resource "netbox_tag" "tag2" {
  name = "Tag2-%[2]s"
  slug = %[2]q
}

resource "netbox_ike_policy" "test" {
  name    = %[3]q
  version = 2
  %[4]s
}
`, slug1, slug2, name, tagsConfig)
}

func TestAccIKEPolicyResource_validationErrors(t *testing.T) {
	testutil.RunMultiValidationErrorTest(t, testutil.MultiValidationErrorTestConfig{
		ResourceName: "netbox_ike_policy",
		TestCases: map[string]testutil.ValidationErrorCase{
			"missing_name": {
				Config: func() string {
					return `
provider "netbox" {}

resource "netbox_ike_policy" "test" {
  # name missing
  version = 1
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
		},
	})
}
