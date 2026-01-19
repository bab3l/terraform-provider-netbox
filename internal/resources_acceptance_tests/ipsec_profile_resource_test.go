package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIPSECProfileResource_basic(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-ipsec-prof")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIPSecProfileCleanup(name)
	cleanup.RegisterIKEPolicyCleanup(name + "-ike-policy")
	cleanup.RegisterIPSecPolicyCleanup(name + "-ipsec-policy")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIPSECProfileResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ipsec_profile.test", "id"),
					resource.TestCheckResourceAttr("netbox_ipsec_profile.test", "name", name),
					resource.TestCheckResourceAttr("netbox_ipsec_profile.test", "mode", "esp"),
				),
			},
		},
	})
}

func TestAccIPSECProfileResource_full(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-ipsec-prof-full")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIPSecProfileCleanup(name)
	cleanup.RegisterIKEPolicyCleanup(name + "-ike-policy")
	cleanup.RegisterIPSecPolicyCleanup(name + "-ipsec-policy")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIPSECProfileResourceConfig_full(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ipsec_profile.test", "id"),
					resource.TestCheckResourceAttr("netbox_ipsec_profile.test", "name", name),
					resource.TestCheckResourceAttr("netbox_ipsec_profile.test", "mode", "ah"),
					resource.TestCheckResourceAttr("netbox_ipsec_profile.test", "description", "Test IPsec profile"),
				),
			},
		},
	})
}

func TestAccIPSECProfileResource_update(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-ipsec-prof-update")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIPSecProfileCleanup(name)
	cleanup.RegisterIKEPolicyCleanup(name + "-ike-policy")
	cleanup.RegisterIPSecPolicyCleanup(name + "-ipsec-policy")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIPSECProfileResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ipsec_profile.test", "id"),
					resource.TestCheckResourceAttr("netbox_ipsec_profile.test", "name", name),
				),
			},
			{
				Config: testAccIPSECProfileResourceConfig_full(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ipsec_profile.test", "id"),
					resource.TestCheckResourceAttr("netbox_ipsec_profile.test", "description", "Test IPsec profile"),
				),
			},
		},
	})
}

func TestAccIPSECProfileResource_import(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-ipsec-prof")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIPSecProfileCleanup(name)
	cleanup.RegisterIKEPolicyCleanup(name + "-ike-policy")
	cleanup.RegisterIPSecPolicyCleanup(name + "-ipsec-policy")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIPSECProfileResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ipsec_profile.test", "id"),
				),
			},
			{
				ResourceName:      "netbox_ipsec_profile.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:   testAccIPSECProfileResourceConfig_basic(name),
				PlanOnly: true,
			},
		},
	})
}

func TestAccIPSECProfileResource_tagLifecycle(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-ipsec-prof-tags")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2Slug := testutil.RandomSlug("tag2")
	tag3Slug := testutil.RandomSlug("tag3")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIPSecProfileCleanup(name)
	cleanup.RegisterIKEPolicyCleanup(name + "-ike-policy")
	cleanup.RegisterIPSecPolicyCleanup(name + "-ipsec-policy")
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)
	cleanup.RegisterTagCleanup(tag3Slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIPSECProfileResourceConfig_tags(name, tag1Slug, tag2Slug, tag3Slug, caseTag1Tag2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_ipsec_profile.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemAttr("netbox_ipsec_profile.test", "tags.*", tag1Slug),
					resource.TestCheckTypeSetElemAttr("netbox_ipsec_profile.test", "tags.*", tag2Slug),
				),
			},
			{
				Config: testAccIPSECProfileResourceConfig_tags(name, tag1Slug, tag2Slug, tag3Slug, caseTag1Uscore2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_ipsec_profile.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemAttr("netbox_ipsec_profile.test", "tags.*", tag1Slug),
					resource.TestCheckTypeSetElemAttr("netbox_ipsec_profile.test", "tags.*", tag2Slug),
				),
			},
			{
				Config: testAccIPSECProfileResourceConfig_tags(name, tag1Slug, tag2Slug, tag3Slug, caseTag3),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_ipsec_profile.test", "tags.#", "1"),
					resource.TestCheckTypeSetElemAttr("netbox_ipsec_profile.test", "tags.*", tag3Slug),
				),
			},
			{
				Config: testAccIPSECProfileResourceConfig_tags(name, tag1Slug, tag2Slug, tag3Slug, tagsEmpty),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_ipsec_profile.test", "tags.#", "0"),
				),
			},
		},
	})
}

func TestAccIPSECProfileResource_tagOrderInvariance(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-ipsec-prof-tag-order")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2Slug := testutil.RandomSlug("tag2")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIPSecProfileCleanup(name)
	cleanup.RegisterIKEPolicyCleanup(name + "-ike-policy")
	cleanup.RegisterIPSecPolicyCleanup(name + "-ipsec-policy")
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIPSECProfileResourceConfig_tagsOrder(name, tag1Slug, tag2Slug, caseTag1Tag2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_ipsec_profile.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemAttr("netbox_ipsec_profile.test", "tags.*", tag1Slug),
					resource.TestCheckTypeSetElemAttr("netbox_ipsec_profile.test", "tags.*", tag2Slug),
				),
			},
			{
				Config: testAccIPSECProfileResourceConfig_tagsOrder(name, tag1Slug, tag2Slug, caseTag2Uscore1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_ipsec_profile.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemAttr("netbox_ipsec_profile.test", "tags.*", tag1Slug),
					resource.TestCheckTypeSetElemAttr("netbox_ipsec_profile.test", "tags.*", tag2Slug),
				),
			},
		},
	})
}

func TestAccIPSECProfileResource_externalDeletion(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-ipsec-prof-extdel")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIPSecProfileCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIPSECProfileResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ipsec_profile.test", "id"),
					resource.TestCheckResourceAttr("netbox_ipsec_profile.test", "name", name),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					items, _, err := client.VpnAPI.VpnIpsecProfilesList(context.Background()).Name([]string{name}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find IPSec profile for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.VpnAPI.VpnIpsecProfilesDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete IPSec profile: %v", err)
					}
					t.Logf("Successfully externally deleted IPSec profile with ID: %d", itemID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccIPSECProfileResourceConfig_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "netbox_ipsec_profile" "test" {
  name           = %q
  mode           = "esp"
  ike_policy     = netbox_ike_policy.test.id
  ipsec_policy   = netbox_ipsec_policy.test.id
}
`, testAccIPSECProfileResourcePrereqs(name), name)
}

func testAccIPSECProfileResourceConfig_full(name string) string {
	return fmt.Sprintf(`
%s

resource "netbox_ipsec_profile" "test" {
  name           = %q
  mode           = "ah"
  ike_policy     = netbox_ike_policy.test.id
  ipsec_policy   = netbox_ipsec_policy.test.id
  description    = "Test IPsec profile"
}
`, testAccIPSECProfileResourcePrereqs(name), name)

}

func testAccIPSECProfileResourceConfig_tags(name, tag1Slug, tag2Slug, tag3Slug, tagCase string) string {
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

%[6]s

resource "netbox_ipsec_profile" "test" {
  name         = %[1]q
  mode         = "esp"
  ike_policy   = netbox_ike_policy.test.id
  ipsec_policy = netbox_ipsec_policy.test.id
  %[5]s
}
`, name, tag1Slug, tag2Slug, tag3Slug, tagsConfig, testAccIPSECProfileResourcePrereqs(name))
}

func testAccIPSECProfileResourceConfig_tagsOrder(name, tag1Slug, tag2Slug, tagCase string) string {
	var tagsConfig string
	switch tagCase {
	case caseTag1Tag2:
		tagsConfig = tagsDoubleSlug
	case caseTag2Uscore1:
		tagsConfig = tagsDoubleSlugReversed
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

%[5]s

resource "netbox_ipsec_profile" "test" {
  name         = %[1]q
  mode         = "esp"
  ike_policy   = netbox_ike_policy.test.id
  ipsec_policy = netbox_ipsec_policy.test.id
  %[4]s
}
`, name, tag1Slug, tag2Slug, tagsConfig, testAccIPSECProfileResourcePrereqs(name))
}

func testAccIPSECProfileResourcePrereqs(name string) string {
	return fmt.Sprintf(`
resource "netbox_ike_policy" "test" {
  name    = %q
  version = "1"
  mode    = "main"
}

resource "netbox_ipsec_policy" "test" {
  name = %q
}
`, name+"-ike-policy", name+"-ipsec-policy")
}

func TestAccIPSecProfileResource_removeOptionalFields(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-ipsec-profile-rem")
	ikePolicyName := testutil.RandomName("tf-test-ike-policy-rem")
	ipsecPolicyName := testutil.RandomName("tf-test-ipsec-policy-rem")
	const testDescription = "Test Description"
	const testComments = "Test Comments"

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIPSecProfileCleanup(name)
	cleanup.RegisterIKEPolicyCleanup(ikePolicyName)
	cleanup.RegisterIPSecPolicyCleanup(ipsecPolicyName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckIPSecProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIPSECProfileResourceConfig_withDescriptionComments(ikePolicyName, ipsecPolicyName, name, testDescription, testComments),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_ipsec_profile.test", "name", name),
					resource.TestCheckResourceAttr("netbox_ipsec_profile.test", "description", testDescription),
					resource.TestCheckResourceAttr("netbox_ipsec_profile.test", "comments", testComments),
				),
			},
			{
				Config: testAccIPSECProfileResourceConfig_basicWithPrereqs(ikePolicyName, ipsecPolicyName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_ipsec_profile.test", "name", name),
					resource.TestCheckNoResourceAttr("netbox_ipsec_profile.test", "description"),
					resource.TestCheckNoResourceAttr("netbox_ipsec_profile.test", "comments"),
				),
			},
		},
	})
}

func testAccIPSECProfileResourceConfig_basicWithPrereqs(ikePolicyName, ipsecPolicyName, name string) string {
	return fmt.Sprintf(`
resource "netbox_ike_proposal" "test" {
  name                     = "ike-proposal-for-profile"
  authentication_method    = "preshared-keys"
  encryption_algorithm     = "aes-128-cbc"
  authentication_algorithm = "hmac-sha256"
  group                    = 14
}

resource "netbox_ike_policy" "test" {
  name      = %[1]q
  version   = 1
  mode      = "main"
  proposals = [netbox_ike_proposal.test.id]
}

resource "netbox_ipsec_proposal" "test" {
  name                     = "ipsec-proposal-for-profile"
  encryption_algorithm     = "aes-128-cbc"
  authentication_algorithm = "hmac-sha256"
}

resource "netbox_ipsec_policy" "test" {
  name      = %[2]q
  proposals = [netbox_ipsec_proposal.test.id]
}

resource "netbox_ipsec_profile" "test" {
  name         = %[3]q
  mode         = "esp"
  ike_policy   = netbox_ike_policy.test.id
  ipsec_policy = netbox_ipsec_policy.test.id
}
`, ikePolicyName, ipsecPolicyName, name)
}

func testAccIPSECProfileResourceConfig_withDescriptionComments(ikePolicyName, ipsecPolicyName, name, description, comments string) string {
	return fmt.Sprintf(`
resource "netbox_ike_proposal" "test" {
  name                     = "ike-proposal-for-profile"
  authentication_method    = "preshared-keys"
  encryption_algorithm     = "aes-128-cbc"
  authentication_algorithm = "hmac-sha256"
  group                    = 14
}

resource "netbox_ike_policy" "test" {
  name      = %[1]q
  version   = 1
  mode      = "main"
  proposals = [netbox_ike_proposal.test.id]
}

resource "netbox_ipsec_proposal" "test" {
  name                     = "ipsec-proposal-for-profile"
  encryption_algorithm     = "aes-128-cbc"
  authentication_algorithm = "hmac-sha256"
}

resource "netbox_ipsec_policy" "test" {
  name      = %[2]q
  proposals = [netbox_ipsec_proposal.test.id]
}

resource "netbox_ipsec_profile" "test" {
  name         = %[3]q
  mode         = "esp"
  ike_policy   = netbox_ike_policy.test.id
  ipsec_policy = netbox_ipsec_policy.test.id
  description  = %[4]q
  comments     = %[5]q
}
`, ikePolicyName, ipsecPolicyName, name, description, comments)
}

func TestAccConsistency_IPSECProfile_LiteralNames(t *testing.T) {
	t.Parallel()

	ikePolicyName := testutil.RandomName("tf-test-ike-policy")
	ipsecPolicyName := testutil.RandomName("tf-test-ipsec-policy")
	ipsecProfileName := testutil.RandomName("tf-test-ipsec-profile-lit")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIPSecProfileCleanup(ipsecProfileName)
	cleanup.RegisterIKEPolicyCleanup(ikePolicyName)
	cleanup.RegisterIPSecPolicyCleanup(ipsecPolicyName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIPSECProfileConsistencyLiteralNamesConfig(ikePolicyName, ipsecPolicyName, ipsecProfileName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ipsec_profile.test", "id"),
					resource.TestCheckResourceAttr("netbox_ipsec_profile.test", "name", ipsecProfileName),
				),
			},
			{
				Config:   testAccIPSECProfileConsistencyLiteralNamesConfig(ikePolicyName, ipsecPolicyName, ipsecProfileName),
				PlanOnly: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ipsec_profile.test", "id"),
				),
			},
		},
	})
}

func testAccIPSECProfileConsistencyLiteralNamesConfig(ikePolicyName, ipsecPolicyName, profileName string) string {
	return fmt.Sprintf(`
resource "netbox_ike_policy" "test" {
  name    = %q
  version = "1"
  mode    = "main"
}

resource "netbox_ipsec_policy" "test" {
  name = %q
}

resource "netbox_ipsec_profile" "test" {
  name         = %q
  mode         = "esp"
  ike_policy   = netbox_ike_policy.test.id
  ipsec_policy = netbox_ipsec_policy.test.id
}
`, ikePolicyName, ipsecPolicyName, profileName)
}

func TestAccIPSecProfileResource_validationErrors(t *testing.T) {
	testutil.RunMultiValidationErrorTest(t, testutil.MultiValidationErrorTestConfig{
		ResourceName: "netbox_ipsec_profile",
		TestCases: map[string]testutil.ValidationErrorCase{
			"missing_name": {
				Config: func() string {
					return `
provider "netbox" {}

resource "netbox_ipsec_profile" "test" {
  # name missing
  mode = "esp"
  ike_policy = "test-ike"
  ipsec_policy = "test-ipsec"
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
			"missing_mode": {
				Config: func() string {
					return `
provider "netbox" {}

resource "netbox_ipsec_profile" "test" {
  name = "Test Profile"
  # mode missing
  ike_policy = "test-ike"
  ipsec_policy = "test-ipsec"
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
			"missing_ike_policy": {
				Config: func() string {
					return `
provider "netbox" {}

resource "netbox_ipsec_profile" "test" {
  name = "Test Profile"
  mode = "esp"
  # ike_policy missing
  ipsec_policy = "test-ipsec"
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
			"missing_ipsec_policy": {
				Config: func() string {
					return `
provider "netbox" {}

resource "netbox_ipsec_profile" "test" {
  name = "Test Profile"
  mode = "esp"
  ike_policy = "test-ike"
  # ipsec_policy missing
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
		},
	})
}
