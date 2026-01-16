package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIPSECPolicyResource_basic(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-ipsec-policy")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIPSecPolicyCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIPSECPolicyResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ipsec_policy.test", "id"),
					resource.TestCheckResourceAttr("netbox_ipsec_policy.test", "name", name),
				),
			},
		},
	})
}

func TestAccIPSECPolicyResource_full(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-ipsec-policy-full")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIPSecPolicyCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIPSECPolicyResourceConfig_full(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ipsec_policy.test", "id"),
					resource.TestCheckResourceAttr("netbox_ipsec_policy.test", "name", name),
					resource.TestCheckResourceAttr("netbox_ipsec_policy.test", "pfs_group", "14"),
					resource.TestCheckResourceAttr("netbox_ipsec_policy.test", "description", "Test IPsec policy"),
				),
			},
		},
	})
}

func TestAccIPSECPolicyResource_update(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-ipsec-policy-update")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIPSecPolicyCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIPSECPolicyResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ipsec_policy.test", "id"),
					resource.TestCheckResourceAttr("netbox_ipsec_policy.test", "name", name),
				),
			},
			{
				Config: testAccIPSECPolicyResourceConfig_full(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ipsec_policy.test", "id"),
					resource.TestCheckResourceAttr("netbox_ipsec_policy.test", "pfs_group", "14"),
					resource.TestCheckResourceAttr("netbox_ipsec_policy.test", "description", "Test IPsec policy"),
				),
			},
		},
	})
}

func TestAccIPSECPolicyResource_import(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-ipsec-policy")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIPSecPolicyCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIPSECPolicyResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ipsec_policy.test", "id"),
				),
			},
			{
				ResourceName:      "netbox_ipsec_policy.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccIPSECPolicyResource_tagLifecycle(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-ipsec-policy-tags")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2Slug := testutil.RandomSlug("tag2")
	tag3Slug := testutil.RandomSlug("tag3")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIPSecPolicyCleanup(name)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)
	cleanup.RegisterTagCleanup(tag3Slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIPSECPolicyResourceConfig_tags(name, tag1Slug, tag2Slug, tag3Slug, caseTag1Tag2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_ipsec_policy.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_ipsec_policy.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag1-%s", tag1Slug),
						"slug": tag1Slug,
					}),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_ipsec_policy.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag2-%s", tag2Slug),
						"slug": tag2Slug,
					}),
				),
			},
			{
				Config: testAccIPSECPolicyResourceConfig_tags(name, tag1Slug, tag2Slug, tag3Slug, caseTag1Uscore2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_ipsec_policy.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_ipsec_policy.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag1-%s", tag1Slug),
						"slug": tag1Slug,
					}),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_ipsec_policy.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag2-%s", tag2Slug),
						"slug": tag2Slug,
					}),
				),
			},
			{
				Config: testAccIPSECPolicyResourceConfig_tags(name, tag1Slug, tag2Slug, tag3Slug, caseTag3),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_ipsec_policy.test", "tags.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_ipsec_policy.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag3-%s", tag3Slug),
						"slug": tag3Slug,
					}),
				),
			},
			{
				Config: testAccIPSECPolicyResourceConfig_tags(name, tag1Slug, tag2Slug, tag3Slug, tagsEmpty),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_ipsec_policy.test", "tags.#", "0"),
				),
			},
		},
	})
}

func TestAccIPSECPolicyResource_tagOrderInvariance(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-ipsec-policy-tag-order")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2Slug := testutil.RandomSlug("tag2")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIPSecPolicyCleanup(name)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIPSECPolicyResourceConfig_tagsOrder(name, tag1Slug, tag2Slug, caseTag1Tag2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_ipsec_policy.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_ipsec_policy.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag1-%s", tag1Slug),
						"slug": tag1Slug,
					}),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_ipsec_policy.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag2-%s", tag2Slug),
						"slug": tag2Slug,
					}),
				),
			},
			{
				Config: testAccIPSECPolicyResourceConfig_tagsOrder(name, tag1Slug, tag2Slug, caseTag2Uscore1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_ipsec_policy.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_ipsec_policy.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag1-%s", tag1Slug),
						"slug": tag1Slug,
					}),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_ipsec_policy.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag2-%s", tag2Slug),
						"slug": tag2Slug,
					}),
				),
			},
		},
	})
}

func TestAccIPSECPolicyResource_externalDeletion(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-ipsec-policy-extdel")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIPSecPolicyCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIPSECPolicyResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ipsec_policy.test", "id"),
					resource.TestCheckResourceAttr("netbox_ipsec_policy.test", "name", name),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					items, _, err := client.VpnAPI.VpnIpsecPoliciesList(context.Background()).Name([]string{name}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find IPSec policy for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.VpnAPI.VpnIpsecPoliciesDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete IPSec policy: %v", err)
					}
					t.Logf("Successfully externally deleted IPSec policy with ID: %d", itemID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccIPSECPolicyResourceConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "netbox_ipsec_proposal" "test" {
  name                 = "%s-proposal"
  encryption_algorithm = "aes-128-cbc"
}

resource "netbox_ipsec_policy" "test" {
  name      = %q
  proposals = [netbox_ipsec_proposal.test.id]
}
`, name, name)
}

func testAccIPSECPolicyResourceConfig_full(name string) string {
	return fmt.Sprintf(`
resource "netbox_ipsec_policy" "test" {
  name        = %q
  pfs_group   = 14
  description = "Test IPsec policy"
}
`, name)
}

func testAccIPSECPolicyResourceConfig_tags(name, tag1Slug, tag2Slug, tag3Slug, tagCase string) string {
	var tagsConfig string
	switch tagCase {
	case caseTag1Tag2:
		tagsConfig = tagsDoubleNested
	case caseTag1Uscore2:
		tagsConfig = tagsDoubleNested
	case caseTag3:
		tagsConfig = tagsSingleNested
	case tagsEmpty:
		tagsConfig = tagsEmpty
	}

	return fmt.Sprintf(`
resource "netbox_tag" "tag1" {
  name = "Tag1-%[2]s"
  slug = %[2]q
}

resource "netbox_tag" "tag2" {
  name = "Tag2-%[3]s"
  slug = %[3]q
}

resource "netbox_tag" "tag3" {
  name = "Tag3-%[4]s"
  slug = %[4]q
}

resource "netbox_ipsec_policy" "test" {
  name = %[1]q
  %[5]s
}
`, name, tag1Slug, tag2Slug, tag3Slug, tagsConfig)
}

func testAccIPSECPolicyResourceConfig_tagsOrder(name, tag1Slug, tag2Slug, tagCase string) string {
	var tagsConfig string
	switch tagCase {
	case caseTag1Tag2:
		tagsConfig = tagsDoubleNested
	case caseTag2Uscore1:
		tagsConfig = tagsDoubleNestedReversed
	}

	return fmt.Sprintf(`
resource "netbox_tag" "tag1" {
  name = "Tag1-%[2]s"
  slug = %[2]q
}

resource "netbox_tag" "tag2" {
  name = "Tag2-%[3]s"
  slug = %[3]q
}

resource "netbox_ipsec_policy" "test" {
  name = %[1]q
  %[4]s
}
`, name, tag1Slug, tag2Slug, tagsConfig)
}

func TestAccIPSecPolicyResource_removeOptionalFields(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-ipsec-policy-rem")
	const testDescription = "Test Description"
	const testComments = "Test Comments"

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIPSecPolicyCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckIPSecPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIPSECPolicyResourceConfig_withDescriptionCommentsPFSGroup(name, testDescription, testComments, 14),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_ipsec_policy.test", "name", name),
					resource.TestCheckResourceAttr("netbox_ipsec_policy.test", "description", testDescription),
					resource.TestCheckResourceAttr("netbox_ipsec_policy.test", "comments", testComments),
					resource.TestCheckResourceAttr("netbox_ipsec_policy.test", "pfs_group", "14"),
					resource.TestCheckResourceAttr("netbox_ipsec_policy.test", "proposals.#", "1"),
				),
			},
			{
				Config: testAccIPSECPolicyResourceConfig_nameOnly(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_ipsec_policy.test", "name", name),
					resource.TestCheckNoResourceAttr("netbox_ipsec_policy.test", "description"),
					resource.TestCheckNoResourceAttr("netbox_ipsec_policy.test", "comments"),
					resource.TestCheckResourceAttr("netbox_ipsec_policy.test", "pfs_group", "14"),
					resource.TestCheckNoResourceAttr("netbox_ipsec_policy.test", "proposals"),
				),
			},
		},
	})
}

func testAccIPSECPolicyResourceConfig_withDescriptionCommentsPFSGroup(name, description, comments string, pfsGroup int) string {
	return fmt.Sprintf(`
resource "netbox_ipsec_proposal" "test" {
  name                 = "%s-proposal"
  encryption_algorithm = "aes-128-cbc"
}

resource "netbox_ipsec_policy" "test" {
  name        = %[1]q
  proposals   = [netbox_ipsec_proposal.test.id]
	pfs_group    = %[4]d
  description = %[2]q
  comments    = %[3]q
}
`, name, description, comments, pfsGroup)
}

func testAccIPSECPolicyResourceConfig_nameOnly(name string) string {
	return fmt.Sprintf(`
resource "netbox_ipsec_proposal" "test" {
  name                 = "%s-proposal"
  encryption_algorithm = "aes-128-cbc"
}

resource "netbox_ipsec_policy" "test" {
  name = %q
}
`, name, name)
}

func TestAccConsistency_IPSECPolicy_LiteralNames(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-ipsec-policy-lit")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIPSecPolicyCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIPSECPolicyResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ipsec_policy.test", "id"),
					resource.TestCheckResourceAttr("netbox_ipsec_policy.test", "name", name),
				),
			},
			{
				Config:   testAccIPSECPolicyResourceConfig_basic(name),
				PlanOnly: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ipsec_policy.test", "id"),
				),
			},
		},
	})
}
