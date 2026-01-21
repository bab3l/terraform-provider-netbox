package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAggregateResource_basic(t *testing.T) {
	t.Parallel()

	rirName := testutil.RandomName("tf-test-rir")
	rirSlug := testutil.RandomSlug("tf-test-rir")
	prefix := testutil.RandomIPv4Prefix()

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRIRCleanup(rirSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAggregateResourceConfig_basic(rirName, rirSlug, prefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_aggregate.test", "id"),
					resource.TestCheckResourceAttr("netbox_aggregate.test", "prefix", prefix),
					resource.TestCheckResourceAttrSet("netbox_aggregate.test", "rir"),
				),
			},
			{
				ResourceName:            "netbox_aggregate.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"rir"},
			},
			{
				Config:             testAccAggregateResourceConfig_basic(rirName, rirSlug, prefix),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func TestAccAggregateResource_update(t *testing.T) {
	t.Parallel()

	rirName := testutil.RandomName("tf-test-rir-update")
	rirSlug := testutil.RandomSlug("tf-test-rir-update")
	prefix := testutil.RandomIPv4Prefix()

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRIRCleanup(rirSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAggregateResourceConfig_withDescription(rirName, rirSlug, prefix, testutil.Description1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_aggregate.test", "id"),
					resource.TestCheckResourceAttr("netbox_aggregate.test", "prefix", prefix),
					resource.TestCheckResourceAttr("netbox_aggregate.test", "description", testutil.Description1),
				),
			},
			{
				Config: testAccAggregateResourceConfig_withDescription(rirName, rirSlug, prefix, testutil.Description2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_aggregate.test", "description", testutil.Description2),
				),
			},
		},
	})
}

func TestAccAggregateResource_full(t *testing.T) {
	t.Parallel()

	rirName := testutil.RandomName("tf-test-rir-full")
	rirSlug := testutil.RandomSlug("tf-test-rir-full")
	tenantName := testutil.RandomName("tf-test-tenant-full")
	tenantSlug := testutil.RandomSlug("tf-test-tenant-full")
	prefix := testutil.RandomIPv4Prefix()
	description := testutil.RandomName("description")
	updatedDescription := "Updated aggregate description"
	dateAdded := "2024-01-15"
	updatedDate := "2024-06-20"

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRIRCleanup(rirSlug)
	cleanup.RegisterTenantCleanup(tenantSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAggregateResourceConfig_full(rirName, rirSlug, tenantName, tenantSlug, prefix, description, testutil.Comments, dateAdded),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_aggregate.test", "id"),
					resource.TestCheckResourceAttr("netbox_aggregate.test", "prefix", prefix),
					resource.TestCheckResourceAttr("netbox_aggregate.test", "description", description),
					resource.TestCheckResourceAttr("netbox_aggregate.test", "comments", testutil.Comments),
					resource.TestCheckResourceAttrSet("netbox_aggregate.test", "tenant"),
					resource.TestCheckResourceAttr("netbox_aggregate.test", "date_added", dateAdded),
				),
			},
			{
				Config: testAccAggregateResourceConfig_full(rirName, rirSlug, tenantName, tenantSlug, prefix, updatedDescription, testutil.Comments, updatedDate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_aggregate.test", "description", updatedDescription),
					resource.TestCheckResourceAttr("netbox_aggregate.test", "date_added", updatedDate),
				),
			},
		},
	})
}

func TestAccAggregateResource_removeOptionalFields(t *testing.T) {
	t.Parallel()

	rirName := testutil.RandomName("tf-test-rir-optional")
	rirSlug := testutil.RandomSlug("tf-test-rir-optional")
	tenantName := testutil.RandomName("tf-test-tenant-optional")
	tenantSlug := testutil.RandomSlug("tf-test-tenant-optional")
	prefix := testutil.RandomIPv4Prefix()

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRIRCleanup(rirSlug)
	cleanup.RegisterTenantCleanup(tenantSlug)

	testutil.TestRemoveOptionalFields(t, testutil.MultiFieldOptionalTestConfig{
		ResourceName: "netbox_aggregate",
		BaseConfig: func() string {
			return testAccAggregateResourceConfig_withTenant(rirName, rirSlug, tenantName, tenantSlug, prefix)
		},
		ConfigWithFields: func() string {
			return testAccAggregateResourceConfig_full(
				rirName, rirSlug, tenantName, tenantSlug,
				prefix,
				"Test description",
				"Test comments",
				"2024-01-15",
			)
		},
		OptionalFields: map[string]string{
			"description": "Test description",
			"comments":    "Test comments",
			"date_added":  "2024-01-15",
			// Note: tenant is not included as it requires TestCheckResourceAttrSet verification
			// which the test helper doesn't support for ID-based references
		},
		RequiredFields: map[string]string{
			"prefix": prefix,
		},
	})
}

func testAccAggregateResourceConfig_basic(rirName, rirSlug, prefix string) string {
	return fmt.Sprintf(`
resource "netbox_rir" "test" {
  name = %q
  slug = %q
}

resource "netbox_aggregate" "test" {
  prefix = %q
  rir    = netbox_rir.test.id
}
`, rirName, rirSlug, prefix)
}

func testAccAggregateResourceConfig_withDescription(rirName, rirSlug, prefix, description string) string {
	return fmt.Sprintf(`
resource "netbox_rir" "test" {
  name = %q
  slug = %q
}

resource "netbox_aggregate" "test" {
  prefix      = %q
  rir         = netbox_rir.test.id
  description = %q
}
`, rirName, rirSlug, prefix, description)
}

func testAccAggregateResourceConfig_withTenant(rirName, rirSlug, tenantName, tenantSlug, prefix string) string {
	return fmt.Sprintf(`
resource "netbox_rir" "test" {
  name = %q
  slug = %q
}

resource "netbox_tenant" "test" {
  name = %q
  slug = %q
}

resource "netbox_aggregate" "test" {
  prefix = %q
  rir    = netbox_rir.test.id
}
`, rirName, rirSlug, tenantName, tenantSlug, prefix)
}

func testAccAggregateResourceConfig_full(rirName, rirSlug, tenantName, tenantSlug, prefix, description, comments, dateAdded string) string {
	return fmt.Sprintf(`
resource "netbox_rir" "test" {
  name = %q
  slug = %q
}

resource "netbox_tenant" "test" {
  name = %q
  slug = %q
}

resource "netbox_aggregate" "test" {
  prefix      = %q
  rir         = netbox_rir.test.id
  tenant      = netbox_tenant.test.id
  description = %q
  comments    = %q
  date_added  = %q
}
`, rirName, rirSlug, tenantName, tenantSlug, prefix, description, comments, dateAdded)
}

func TestAccConsistency_Aggregate(t *testing.T) {
	t.Parallel()

	prefix := testutil.RandomIPv4Prefix()
	rirName := testutil.RandomName("rir")
	rirSlug := testutil.RandomSlug("rir")
	tenantName := testutil.RandomName("tenant")
	tenantSlug := testutil.RandomSlug("tenant")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRIRCleanup(rirSlug)
	cleanup.RegisterTenantCleanup(tenantSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAggregateConsistencyConfig(prefix, rirName, rirSlug, tenantName, tenantSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_aggregate.test", "prefix", prefix),
					resource.TestCheckResourceAttrSet("netbox_aggregate.test", "rir"),
					resource.TestCheckResourceAttrSet("netbox_aggregate.test", "tenant"),
				),
			},
			{
				PlanOnly: true,
				Config:   testAccAggregateConsistencyConfig(prefix, rirName, rirSlug, tenantName, tenantSlug),
			},
		},
	})
}

func testAccAggregateConsistencyConfig(prefix, rirName, rirSlug, tenantName, tenantSlug string) string {
	return fmt.Sprintf(`
resource "netbox_rir" "test" {
  name = "%[2]s"
  slug = "%[3]s"
}

resource "netbox_tenant" "test" {
  name = "%[4]s"
  slug = "%[5]s"
}

resource "netbox_aggregate" "test" {
  prefix = "%[1]s"
  rir = netbox_rir.test.id
  tenant = netbox_tenant.test.name
}
`, prefix, rirName, rirSlug, tenantName, tenantSlug)
}

// TestAccConsistency_Aggregate_LiteralNames tests that reference attributes specified as literal string names
// are preserved and do not cause drift when the API returns numeric IDs.
func TestAccConsistency_Aggregate_LiteralNames(t *testing.T) {
	t.Parallel()

	prefix := testutil.RandomIPv4Prefix()
	rirName := testutil.RandomName("rir")
	rirSlug := testutil.RandomSlug("rir")
	tenantName := testutil.RandomName("tenant")
	tenantSlug := testutil.RandomSlug("tenant")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRIRCleanup(rirSlug)
	cleanup.RegisterTenantCleanup(tenantSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAggregateConsistencyLiteralNamesConfig(prefix, rirName, rirSlug, tenantName, tenantSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_aggregate.test", "prefix", prefix),
					resource.TestCheckResourceAttrSet("netbox_aggregate.test", "rir"),
					resource.TestCheckResourceAttrSet("netbox_aggregate.test", "tenant"),
				),
			},
			{
				// Critical: Verify no drift when refreshing state
				PlanOnly: true,
				Config:   testAccAggregateConsistencyLiteralNamesConfig(prefix, rirName, rirSlug, tenantName, tenantSlug),
			},
		},
	})
}

func testAccAggregateConsistencyLiteralNamesConfig(prefix, rirName, rirSlug, tenantName, tenantSlug string) string {
	return fmt.Sprintf(`
resource "netbox_rir" "test" {
  name = "%[2]s"
  slug = "%[3]s"
}

resource "netbox_tenant" "test" {
  name = "%[4]s"
  slug = "%[5]s"
}

resource "netbox_aggregate" "test" {
  prefix = "%[1]s"
  # Use literal string names to mimic existing user state
  rir = "%[3]s"
  tenant = "%[4]s"
  depends_on = [netbox_rir.test, netbox_tenant.test]
}

`, prefix, rirName, rirSlug, tenantName, tenantSlug)

}

func TestAccAggregateResource_externalDeletion(t *testing.T) {
	t.Parallel()
	prefix := testutil.RandomIPv4Prefix()
	rirName := testutil.RandomName("tf-test-rir-ext-del")
	rirSlug := testutil.RandomSlug("tf-test-rir-ext-del")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRIRCleanup(rirSlug)
	cleanup.RegisterAggregateCleanup(prefix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAggregateResourceConfig_basic(rirName, rirSlug, prefix),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_aggregate.test", "prefix", prefix),
					resource.TestCheckResourceAttrSet("netbox_aggregate.test", "id"),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					items, _, err := client.IpamAPI.IpamAggregatesList(context.Background()).Prefix(prefix).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find aggregate for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.IpamAPI.IpamAggregatesDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete aggregate: %v", err)
					}
					t.Logf("Successfully externally deleted aggregate with ID: %d", itemID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

// NOTE: Custom field tests for aggregate resource are in resources_acceptance_tests_customfields package

// Enhancement 2: Test import with full optional fields
// This verifies that import correctly populates all fields from the API.
func TestAccAggregateResource_importPreservesOptionalFields(t *testing.T) {
	t.Parallel()

	prefix := testutil.RandomIPv4Prefix()
	rirName := testutil.RandomName("rir")
	rirSlug := testutil.RandomSlug("rir")
	tenantName := testutil.RandomName("tenant")
	tenantSlug := testutil.RandomSlug("tenant")
	description := testutil.RandomName("description")
	comments := testutil.RandomName("comments")
	dateAdded := "2024-01-01"

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRIRCleanup(rirSlug)
	cleanup.RegisterTenantCleanup(tenantSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Step 1: Create resource with ALL optional fields
			{
				Config: testAccAggregateResourceConfig_full(rirName, rirSlug, tenantName, tenantSlug, prefix, description, comments, dateAdded),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_aggregate.test", "id"),
					resource.TestCheckResourceAttr("netbox_aggregate.test", "prefix", prefix),
					resource.TestCheckResourceAttrSet("netbox_aggregate.test", "tenant"),
					resource.TestCheckResourceAttr("netbox_aggregate.test", "description", description),
					resource.TestCheckResourceAttr("netbox_aggregate.test", "comments", comments),
					resource.TestCheckResourceAttr("netbox_aggregate.test", "date_added", dateAdded),
				),
			},
			// Step 2: Import with SAME full config - verifies import preserves all fields
			{
				Config:            testAccAggregateResourceConfig_full(rirName, rirSlug, tenantName, tenantSlug, prefix, description, comments, dateAdded),
				ResourceName:      "netbox_aggregate.test",
				ImportState:       true,
				ImportStateVerify: true,
				// Note: tenant format may change (ID->name/slug), rir is always lookup
				ImportStateVerifyIgnore: []string{"rir", "tenant"},
				Check: resource.ComposeTestCheckFunc(
					testutil.ReferenceFieldCheck("netbox_aggregate.test", "rir"),
					testutil.ReferenceFieldCheck("netbox_aggregate.test", "tenant"),
				),
			},
			// Step 3: Verify no changes after import
			{
				Config:   testAccAggregateResourceConfig_full(rirName, rirSlug, tenantName, tenantSlug, prefix, description, comments, dateAdded),
				PlanOnly: true,
			},
		},
	})
}

// =============================================================================
// STANDARDIZED TAG TESTS (using helpers)
// =============================================================================

// TestAccAggregateResource_tagLifecycle tests the complete tag lifecycle using RunTagLifecycleTest helper.
func TestAccAggregateResource_tagLifecycle(t *testing.T) {
	t.Parallel()

	rirName := testutil.RandomName("rir-tag")
	rirSlug := testutil.RandomSlug("rir-tag")
	prefix := testutil.RandomIPv4Prefix()
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
		ResourceName: "netbox_aggregate",
		ConfigWithoutTags: func() string {
			return testAccAggregateResourceConfig_tagLifecycle(rirName, rirSlug, prefix, tag1Name, tag1Slug, tag2Name, tag2Slug, tag3Name, tag3Slug, "none")
		},
		ConfigWithTags: func() string {
			return testAccAggregateResourceConfig_tagLifecycle(rirName, rirSlug, prefix, tag1Name, tag1Slug, tag2Name, tag2Slug, tag3Name, tag3Slug, "tag1_tag2")
		},
		ConfigWithDifferentTags: func() string {
			return testAccAggregateResourceConfig_tagLifecycle(rirName, rirSlug, prefix, tag1Name, tag1Slug, tag2Name, tag2Slug, tag3Name, tag3Slug, "tag2_tag3")
		},
		ExpectedTagCount:          2,
		ExpectedDifferentTagCount: 2,
		CheckDestroy:              testutil.CheckAggregateDestroy,
	})
}

func testAccAggregateResourceConfig_tagLifecycle(rirName, rirSlug, prefix, tag1Name, tag1Slug, tag2Name, tag2Slug, tag3Name, tag3Slug, tagSet string) string {
	baseConfig := fmt.Sprintf(`
resource "netbox_rir" "test" {
  name = %[1]q
  slug = %[2]q
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
`, rirName, rirSlug, prefix, tag1Name, tag1Slug, tag2Name, tag2Slug, tag3Name, tag3Slug)

	//nolint:goconst // tagSet values are test-specific identifiers
	switch tagSet {
	case "tag1_tag2":
		return baseConfig + fmt.Sprintf(`
resource "netbox_aggregate" "test" {
  prefix = %[1]q
  rir    = netbox_rir.test.id
	tags = [netbox_tag.tag1.slug, netbox_tag.tag2.slug]
}
`, prefix)
	case "tag2_tag3":
		return baseConfig + fmt.Sprintf(`
resource "netbox_aggregate" "test" {
  prefix = %[1]q
  rir    = netbox_rir.test.id
	tags = [netbox_tag.tag2.slug, netbox_tag.tag3.slug]
}
`, prefix)
	default: // "none"
		return baseConfig + fmt.Sprintf(`
resource "netbox_aggregate" "test" {
  prefix = %[1]q
  rir    = netbox_rir.test.id
  tags   = []
}
`, prefix)
	}
}

// TestAccAggregateResource_tagOrderInvariance tests tag order using RunTagOrderTest helper.
func TestAccAggregateResource_tagOrderInvariance(t *testing.T) {
	t.Parallel()

	rirName := testutil.RandomName("rir-tag-order")
	rirSlug := testutil.RandomSlug("rir-tag-order")
	prefix := testutil.RandomIPv4Prefix()
	tag1Name := testutil.RandomName("tag1")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2Name := testutil.RandomName("tag2")
	tag2Slug := testutil.RandomSlug("tag2")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRIRCleanup(rirSlug)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)

	testutil.RunTagOrderTest(t, testutil.TagOrderTestConfig{
		ResourceName: "netbox_aggregate",
		ConfigWithTagsOrderA: func() string {
			return testAccAggregateResourceConfig_tagOrder(rirName, rirSlug, prefix, tag1Name, tag1Slug, tag2Name, tag2Slug, true)
		},
		ConfigWithTagsOrderB: func() string {
			return testAccAggregateResourceConfig_tagOrder(rirName, rirSlug, prefix, tag1Name, tag1Slug, tag2Name, tag2Slug, false)
		},
		ExpectedTagCount: 2,
		CheckDestroy:     testutil.CheckAggregateDestroy,
	})
}

func testAccAggregateResourceConfig_tagOrder(rirName, rirSlug, prefix, tag1Name, tag1Slug, tag2Name, tag2Slug string, tag1First bool) string {
	baseConfig := fmt.Sprintf(`
resource "netbox_rir" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_tag" "tag1" {
  name = %[4]q
  slug = %[5]q
}

resource "netbox_tag" "tag2" {
  name = %[6]q
  slug = %[7]q
}
`, rirName, rirSlug, prefix, tag1Name, tag1Slug, tag2Name, tag2Slug)

	if tag1First {
		return baseConfig + fmt.Sprintf(`
resource "netbox_aggregate" "test" {
  prefix = %[1]q
  rir    = netbox_rir.test.id
	tags = [netbox_tag.tag1.slug, netbox_tag.tag2.slug]
}
`, prefix)
	}

	return baseConfig + fmt.Sprintf(`
resource "netbox_aggregate" "test" {
  prefix = %[1]q
  rir    = netbox_rir.test.id
	tags = [netbox_tag.tag2.slug, netbox_tag.tag1.slug]
}
`, prefix)
}

func TestAccAggregateResource_validationErrors(t *testing.T) {
	testutil.RunMultiValidationErrorTest(t, testutil.MultiValidationErrorTestConfig{
		ResourceName: "netbox_aggregate",
		TestCases: map[string]testutil.ValidationErrorCase{
			"missing_prefix": {
				Config: func() string {
					return `
provider "netbox" {}

resource "netbox_rir" "test" {
  name = "test-rir"
  slug = "test-rir"
}

resource "netbox_aggregate" "test" {
  # prefix missing
  rir = netbox_rir.test.id
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
			"missing_rir": {
				Config: func() string {
					return `
provider "netbox" {}

resource "netbox_aggregate" "test" {
  prefix = "10.0.0.0/8"
  # rir missing
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
			"invalid_rir_reference": {
				Config: func() string {
					return `
provider "netbox" {}

resource "netbox_aggregate" "test" {
  prefix = "10.0.0.0/8"
  rir = "nonexistent-rir"
}
`
				},
				ExpectedError: testutil.ErrPatternNotFound,
			},
		},
	})
}
