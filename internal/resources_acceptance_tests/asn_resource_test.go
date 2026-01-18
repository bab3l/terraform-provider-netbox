package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccASNResource_basic(t *testing.T) {
	t.Parallel()

	rirName := testutil.RandomName("tf-test-rir")
	rirSlug := testutil.RandomSlug("tf-test-rir")
	// Generate a random ASN in the private range (64512-64711) - non-overlapping with other tests
	asn := int64(acctest.RandIntRange(64512, 64711))
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRIRCleanup(rirSlug)
	cleanup.RegisterASNCleanup(asn)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccASNResourceConfig_basic(rirName, rirSlug, asn),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_asn.test", "id"),
					resource.TestCheckResourceAttr("netbox_asn.test", "asn", fmt.Sprintf("%d", asn)),
				),
			},
			{
				ResourceName:            "netbox_asn.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"rir"},
			},
		},
	})
}

func TestAccASNResource_full(t *testing.T) {
	t.Parallel()

	rirName := testutil.RandomName("tf-test-rir")
	rirSlug := testutil.RandomSlug("tf-test-rir")
	tenantName := testutil.RandomName("tf-test-tenant")
	tenantSlug := testutil.RandomSlug("tf-test-tenant")

	// Generate a random ASN in the private range (64712-64911) - non-overlapping with other tests
	asn := int64(acctest.RandIntRange(64712, 64911))
	description := testutil.RandomName("description")
	updatedDescription := "Updated ASN description"

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRIRCleanup(rirSlug)
	cleanup.RegisterTenantCleanup(tenantSlug)
	cleanup.RegisterASNCleanup(asn)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccASNResourceConfig_full(rirName, rirSlug, tenantName, tenantSlug, asn, description, testutil.Comments),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_asn.test", "id"),
					resource.TestCheckResourceAttr("netbox_asn.test", "asn", fmt.Sprintf("%d", asn)),
					resource.TestCheckResourceAttr("netbox_asn.test", "description", description),
					resource.TestCheckResourceAttr("netbox_asn.test", "comments", testutil.Comments),
					resource.TestCheckResourceAttrSet("netbox_asn.test", "rir"),
					resource.TestCheckResourceAttrSet("netbox_asn.test", "tenant"),
				),
			},
			{
				Config: testAccASNResourceConfig_full(rirName, rirSlug, tenantName, tenantSlug, asn, updatedDescription, testutil.Comments),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_asn.test", "description", updatedDescription),
				),
			},
		},
	})
}

func TestAccASNResource_update(t *testing.T) {
	t.Parallel()

	rirName := testutil.RandomName("tf-test-rir-update")
	rirSlug := testutil.RandomSlug("tf-test-rir-update")
	// Generate a random ASN in the private range (65112-65311) - non-overlapping with other tests
	asn := int64(acctest.RandIntRange(65112, 65311))

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRIRCleanup(rirSlug)
	cleanup.RegisterASNCleanup(asn)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccASNResourceConfig_update(rirName, rirSlug, asn, testutil.Description1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_asn.test", "asn", fmt.Sprintf("%d", asn)),
					resource.TestCheckResourceAttr("netbox_asn.test", "description", testutil.Description1),
				),
			},
			{
				Config: testAccASNResourceConfig_update(rirName, rirSlug, asn, testutil.Description2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_asn.test", "asn", fmt.Sprintf("%d", asn)),
					resource.TestCheckResourceAttr("netbox_asn.test", "description", testutil.Description2),
				),
			},
		},
	})
}

func TestAccASNResource_removeOptionalFields(t *testing.T) {
	t.Parallel()

	rirName := testutil.RandomName("tf-test-rir-remove")
	rirSlug := testutil.RandomSlug("tf-test-rir-remove")
	tenantName := testutil.RandomName("tf-test-tenant-remove")
	tenantSlug := testutil.RandomSlug("tf-test-tenant-remove")
	asn := int64(acctest.RandIntRange(65312, 65400))

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRIRCleanup(rirSlug)
	cleanup.RegisterTenantCleanup(tenantSlug)
	cleanup.RegisterASNCleanup(asn)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Step 1: Create ASN with all optional fields populated
			{
				Config: testAccASNResourceConfig_full(rirName, rirSlug, tenantName, tenantSlug, asn, "Initial description", testutil.Comments),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_asn.test", "asn", fmt.Sprintf("%d", asn)),
					resource.TestCheckResourceAttrSet("netbox_asn.test", "tenant"),
					resource.TestCheckResourceAttrSet("netbox_asn.test", "rir"),
					resource.TestCheckResourceAttr("netbox_asn.test", "description", "Initial description"),
					resource.TestCheckResourceAttr("netbox_asn.test", "comments", testutil.Comments),
				),
			},
			// Step 2: Remove tenant and verify it's actually removed
			{
				Config: testAccASNResourceConfig_noTenant(rirName, rirSlug, tenantName, tenantSlug, asn, "Description after tenant removal"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_asn.test", "asn", fmt.Sprintf("%d", asn)),
					resource.TestCheckNoResourceAttr("netbox_asn.test", "tenant"),
					resource.TestCheckResourceAttrSet("netbox_asn.test", "rir"),
					resource.TestCheckResourceAttr("netbox_asn.test", "description", "Description after tenant removal"),
					resource.TestCheckResourceAttr("netbox_asn.test", "comments", testutil.Comments),
				),
			},
			// Step 3: Remove description and comments - this tests the null value handling bug
			{
				Config: testAccASNResourceConfig_noDescriptionOrComments(rirName, rirSlug, asn),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_asn.test", "asn", fmt.Sprintf("%d", asn)),
					resource.TestCheckResourceAttrSet("netbox_asn.test", "rir"),
					resource.TestCheckNoResourceAttr("netbox_asn.test", "description"),
					resource.TestCheckNoResourceAttr("netbox_asn.test", "comments"),
					resource.TestCheckNoResourceAttr("netbox_asn.test", "tenant"),
				),
			},
			// Step 4: Re-add all fields to verify they can be set again
			{
				Config: testAccASNResourceConfig_full(rirName, rirSlug, tenantName, tenantSlug, asn, "Final description", testutil.Comments),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_asn.test", "asn", fmt.Sprintf("%d", asn)),
					resource.TestCheckResourceAttrSet("netbox_asn.test", "tenant"),
					resource.TestCheckResourceAttrSet("netbox_asn.test", "rir"),
					resource.TestCheckResourceAttr("netbox_asn.test", "description", "Final description"),
					resource.TestCheckResourceAttr("netbox_asn.test", "comments", testutil.Comments),
				),
			},
		},
	})
}

func TestAccASNResource_external_deletion(t *testing.T) {
	t.Parallel()

	rirName := testutil.RandomName("tf-test-rir-ext-del")
	rirSlug := testutil.RandomSlug("tf-test-rir-ext-del")
	// Generate a random ASN in the private range (65312-65534) - non-overlapping with other tests
	asn := int64(acctest.RandIntRange(65312, 65534))

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRIRCleanup(rirSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccASNResourceConfig_basic(rirName, rirSlug, asn),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_asn.test", "id"),
					resource.TestCheckResourceAttr("netbox_asn.test", "asn", fmt.Sprintf("%d", asn)),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					// Safe conversion with bounds check
					if asn < 0 || asn > 4294967295 { // Max ASN value (32-bit)
						t.Fatalf("ASN value %d out of valid range", asn)
					}
					//nolint:gosec // Safe conversion - bounds checked above
					items, _, err := client.IpamAPI.IpamAsnsList(context.Background()).Asn([]int32{int32(asn)}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find asn for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.IpamAPI.IpamAsnsDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete asn: %v", err)
					}
					t.Logf("Successfully externally deleted asn with ID: %d", itemID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccASNResourceConfig_basic(rirName, rirSlug string, asn int64) string {
	return fmt.Sprintf(`
resource "netbox_rir" "test" {
  name = %q
  slug = %q
}

resource "netbox_asn" "test" {
  asn = %d
  rir = netbox_rir.test.id
}
`, rirName, rirSlug, asn)
}

func testAccASNResourceConfig_full(rirName, rirSlug, tenantName, tenantSlug string, asn int64, description, comments string) string {
	return fmt.Sprintf(`
resource "netbox_rir" "test" {
  name = %q
  slug = %q
}

resource "netbox_tenant" "test" {
  name = %q
  slug = %q
}

resource "netbox_asn" "test" {
  asn         = %d
  rir         = netbox_rir.test.id
  tenant      = netbox_tenant.test.id
  description = %q
  comments    = %q
}
`, rirName, rirSlug, tenantName, tenantSlug, asn, description, comments)
}

func testAccASNResourceConfig_update(rirName, rirSlug string, asn int64, description string) string {
	return fmt.Sprintf(`
resource "netbox_rir" "test" {
  name = %q
  slug = %q
}

resource "netbox_asn" "test" {
  asn         = %d
  rir         = netbox_rir.test.id
  description = %q
}
`, rirName, rirSlug, asn, description)
}

func testAccASNResourceConfig_noTenant(rirName, rirSlug, tenantName, tenantSlug string, asn int64, description string) string {
	return fmt.Sprintf(`
resource "netbox_rir" "test" {
  name = %q
  slug = %q
}

resource "netbox_tenant" "test" {
  name = %q
  slug = %q
}

resource "netbox_asn" "test" {
  asn         = %d
  rir         = netbox_rir.test.id
  description = %q
  comments    = %q
  # tenant intentionally omitted - should be null in state
}
`, rirName, rirSlug, tenantName, tenantSlug, asn, description, testutil.Comments)
}

func testAccASNResourceConfig_noDescriptionOrComments(rirName, rirSlug string, asn int64) string {
	return fmt.Sprintf(`
resource "netbox_rir" "test" {
  name = %q
  slug = %q
}

resource "netbox_asn" "test" {
  asn = %d
  rir = netbox_rir.test.id
  # description and comments intentionally omitted - should be null/empty in state
}
`, rirName, rirSlug, asn)
}

// TestAccConsistency_ASN_LiteralNames tests that reference attributes specified as literal string names
// are preserved and do not cause drift when the API returns numeric IDs.
func TestAccConsistency_ASN_LiteralNames(t *testing.T) {
	t.Parallel()

	asn := int64(65100)
	rirName := testutil.RandomName("rir")
	rirSlug := testutil.RandomSlug("rir")
	tenantName := testutil.RandomName("tenant")
	tenantSlug := testutil.RandomSlug("tenant")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRIRCleanup(rirSlug)
	cleanup.RegisterTenantCleanup(tenantSlug)
	cleanup.RegisterASNCleanup(asn)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccASNConsistencyLiteralNamesConfig(asn, rirName, rirSlug, tenantName, tenantSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_asn.test", "asn", fmt.Sprintf("%d", asn)),
					resource.TestCheckResourceAttr("netbox_asn.test", "rir", rirSlug),
					resource.TestCheckResourceAttr("netbox_asn.test", "tenant", tenantName),
				),
			},
			{
				// Critical: Verify no drift when refreshing state
				PlanOnly: true,
				Config:   testAccASNConsistencyLiteralNamesConfig(asn, rirName, rirSlug, tenantName, tenantSlug),
			},
		},
	})
}

func testAccASNConsistencyLiteralNamesConfig(asn int64, rirName, rirSlug, tenantName, tenantSlug string) string {
	return fmt.Sprintf(`
resource "netbox_rir" "test" {
  name = "%[2]s"
  slug = "%[3]s"
}

resource "netbox_tenant" "test" {
  name = "%[4]s"
  slug = "%[5]s"
}

resource "netbox_asn" "test" {
  asn = %[1]d
  # Use literal string names to mimic existing user state
  rir = "%[3]s"
  tenant = "%[4]s"
  depends_on = [netbox_rir.test, netbox_tenant.test]
}
`, asn, rirName, rirSlug, tenantName, tenantSlug)
}

// NOTE: Custom field tests for ASN resource are in resources_acceptance_tests_customfields package

// =============================================================================
// STANDARDIZED TAG TESTS (using helpers)
// =============================================================================

// TestAccASNResource_tagLifecycle tests the complete tag lifecycle using RunTagLifecycleTest helper.
func TestAccASNResource_tagLifecycle(t *testing.T) {
	t.Parallel()

	asn := int64(acctest.RandIntRange(64712, 64911))
	rirName := testutil.RandomName("rir-tag")
	rirSlug := testutil.RandomSlug("rir-tag")
	tag1Name := testutil.RandomName("tag1")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2Name := testutil.RandomName("tag2")
	tag2Slug := testutil.RandomSlug("tag2")
	tag3Name := testutil.RandomName("tag3")
	tag3Slug := testutil.RandomSlug("tag3")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRIRCleanup(rirSlug)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)
	cleanup.RegisterTagCleanup(tag3Slug)

	testutil.RunTagLifecycleTest(t, testutil.TagLifecycleTestConfig{
		ResourceName: "netbox_asn",
		ConfigWithoutTags: func() string {
			return testAccASNResourceConfig_tagLifecycle(asn, rirName, rirSlug, tag1Name, tag1Slug, tag2Name, tag2Slug, tag3Name, tag3Slug, "none")
		},
		ConfigWithTags: func() string {
			return testAccASNResourceConfig_tagLifecycle(asn, rirName, rirSlug, tag1Name, tag1Slug, tag2Name, tag2Slug, tag3Name, tag3Slug, "tag1_tag2")
		},
		ConfigWithDifferentTags: func() string {
			return testAccASNResourceConfig_tagLifecycle(asn, rirName, rirSlug, tag1Name, tag1Slug, tag2Name, tag2Slug, tag3Name, tag3Slug, "tag2_tag3")
		},
		ExpectedTagCount:          2,
		ExpectedDifferentTagCount: 2,
		CheckDestroy:              testutil.CheckASNDestroy,
	})
}

func testAccASNResourceConfig_tagLifecycle(asn int64, rirName, rirSlug, tag1Name, tag1Slug, tag2Name, tag2Slug, tag3Name, tag3Slug, tagSet string) string {
	baseConfig := fmt.Sprintf(`
resource "netbox_rir" "test" {
  name = %[2]q
  slug = %[3]q
}

resource "netbox_tag" "tag1" {
  name = %[4]q
  slug = %[5]q
}

resource "netbox_tag" "tag2" {
  name = %[6]q
  slug = %[7]q
}

resource "netbox_tag" "tag3" {
  name = %[8]q
  slug = %[9]q
}
`, asn, rirName, rirSlug, tag1Name, tag1Slug, tag2Name, tag2Slug, tag3Name, tag3Slug)

	//nolint:goconst // tagSet values are test-specific identifiers
	switch tagSet {
	case "tag1_tag2":
		return baseConfig + fmt.Sprintf(`
resource "netbox_asn" "test" {
  asn = %[1]d
  rir = netbox_rir.test.id
	tags = [netbox_tag.tag1.slug, netbox_tag.tag2.slug]
}
`, asn)
	case "tag2_tag3":
		return baseConfig + fmt.Sprintf(`
resource "netbox_asn" "test" {
  asn = %[1]d
  rir = netbox_rir.test.id
	tags = [netbox_tag.tag2.slug, netbox_tag.tag3.slug]
}
`, asn)
	default: // "none"
		return baseConfig + fmt.Sprintf(`
resource "netbox_asn" "test" {
  asn  = %[1]d
  rir  = netbox_rir.test.id
  tags = []
}
`, asn)
	}
}

// TestAccASNResource_tagOrderInvariance tests tag order using RunTagOrderTest helper.
func TestAccASNResource_tagOrderInvariance(t *testing.T) {
	t.Parallel()

	asn := int64(acctest.RandIntRange(64912, 65111))
	rirName := testutil.RandomName("rir-tag-order")
	rirSlug := testutil.RandomSlug("rir-tag-order")
	tag1Name := testutil.RandomName("tag1")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2Name := testutil.RandomName("tag2")
	tag2Slug := testutil.RandomSlug("tag2")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRIRCleanup(rirSlug)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)

	testutil.RunTagOrderTest(t, testutil.TagOrderTestConfig{
		ResourceName: "netbox_asn",
		ConfigWithTagsOrderA: func() string {
			return testAccASNResourceConfig_tagOrder(asn, rirName, rirSlug, tag1Name, tag1Slug, tag2Name, tag2Slug, true)
		},
		ConfigWithTagsOrderB: func() string {
			return testAccASNResourceConfig_tagOrder(asn, rirName, rirSlug, tag1Name, tag1Slug, tag2Name, tag2Slug, false)
		},
		ExpectedTagCount: 2,
		CheckDestroy:     testutil.CheckASNDestroy,
	})
}

func testAccASNResourceConfig_tagOrder(asn int64, rirName, rirSlug, tag1Name, tag1Slug, tag2Name, tag2Slug string, tag1First bool) string {
	baseConfig := fmt.Sprintf(`
resource "netbox_rir" "test" {
  name = %[2]q
  slug = %[3]q
}

resource "netbox_tag" "tag1" {
  name = %[4]q
  slug = %[5]q
}

resource "netbox_tag" "tag2" {
  name = %[6]q
  slug = %[7]q
}
`, asn, rirName, rirSlug, tag1Name, tag1Slug, tag2Name, tag2Slug)

	if tag1First {
		return baseConfig + fmt.Sprintf(`
resource "netbox_asn" "test" {
  asn = %[1]d
  rir = netbox_rir.test.id
	tags = [netbox_tag.tag1.slug, netbox_tag.tag2.slug]
}
`, asn)
	}

	return baseConfig + fmt.Sprintf(`
resource "netbox_asn" "test" {
  asn = %[1]d
  rir = netbox_rir.test.id
	tags = [netbox_tag.tag2.slug, netbox_tag.tag1.slug]
}
`, asn)
}

func TestAccASNResource_validationErrors(t *testing.T) {
	testutil.RunMultiValidationErrorTest(t, testutil.MultiValidationErrorTestConfig{
		ResourceName: "netbox_asn",
		TestCases: map[string]testutil.ValidationErrorCase{
			"missing_asn": {
				Config: func() string {
					return `
provider "netbox" {}

resource "netbox_asn" "test" {
  # asn missing
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
		},
	})
}
