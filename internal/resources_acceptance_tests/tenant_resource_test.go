package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// NOTE: Custom field tests for tenant resource are in resources_acceptance_tests_customfields package

func TestAccTenantResource_basic(t *testing.T) {
	t.Parallel()

	// Generate unique names to avoid conflicts between test runs
	name := testutil.RandomName("tf-test-tenant")
	slug := testutil.RandomSlug("tf-test-tenant")

	// Register cleanup to ensure resource is deleted even if test fails
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTenantCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckTenantDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTenantResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_tenant.test", "id"),
					resource.TestCheckResourceAttr("netbox_tenant.test", "name", name),
					resource.TestCheckResourceAttr("netbox_tenant.test", "slug", slug),
				),
			},
			{
				Config:   testAccTenantResourceConfig_basic(name, slug),
				PlanOnly: true,
			},
		},
	})
}

func TestAccTenantResource_full(t *testing.T) {
	t.Parallel()

	// Generate unique names
	name := testutil.RandomName("tf-test-tenant-full")
	slug := testutil.RandomSlug("tf-test-tenant-full")
	groupName := testutil.RandomName("tf-test-tenant-group")
	groupSlug := testutil.RandomSlug("tf-test-tenant-group")
	description := testutil.RandomName("description")
	updatedDescription := testutil.RandomName("description-updated")
	comments := testutil.RandomName("comments")
	updatedComments := testutil.RandomName("comments-updated")
	tagName1 := testutil.RandomName("tag1")
	tagSlug1 := testutil.RandomSlug("tag1")
	tagName2 := testutil.RandomName("tag2")
	tagSlug2 := testutil.RandomSlug("tag2")

	// Register cleanup
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTenantCleanup(slug)
	cleanup.RegisterTenantGroupCleanup(groupSlug)
	cleanup.RegisterTagCleanup(tagSlug1)
	cleanup.RegisterTagCleanup(tagSlug2)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckTenantDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTenantResourceConfig_full(name, slug, groupName, groupSlug, description, comments, tagName1, tagSlug1, tagName2, tagSlug2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_tenant.test", "id"),
					resource.TestCheckResourceAttr("netbox_tenant.test", "name", name),
					resource.TestCheckResourceAttr("netbox_tenant.test", "slug", slug),
					resource.TestCheckResourceAttr("netbox_tenant.test", "description", description),
					resource.TestCheckResourceAttr("netbox_tenant.test", "comments", comments),
					resource.TestCheckResourceAttr("netbox_tenant.test", "group", groupName),
					resource.TestCheckResourceAttr("netbox_tenant.test", "tags.#", "2"),
				),
			},
			{
				Config:   testAccTenantResourceConfig_full(name, slug, groupName, groupSlug, description, comments, tagName1, tagSlug1, tagName2, tagSlug2),
				PlanOnly: true,
			},
			{
				Config: testAccTenantResourceConfig_fullUpdate(name, slug, groupName, groupSlug, updatedDescription, updatedComments, tagName1, tagSlug1, tagName2, tagSlug2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_tenant.test", "description", updatedDescription),
					resource.TestCheckResourceAttr("netbox_tenant.test", "comments", updatedComments),
				),
			},
			{
				Config:   testAccTenantResourceConfig_fullUpdate(name, slug, groupName, groupSlug, updatedDescription, updatedComments, tagName1, tagSlug1, tagName2, tagSlug2),
				PlanOnly: true,
			},
		},
	})
}

func TestAccTenantResource_tagLifecycle(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-tenant-tags")
	slug := testutil.RandomSlug("tf-test-tenant-tags")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2Slug := testutil.RandomSlug("tag2")
	tag3Slug := testutil.RandomSlug("tag3")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTenantCleanup(slug)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)
	cleanup.RegisterTagCleanup(tag3Slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckTenantDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTenantResourceConfig_tags(name, slug, tag1Slug, tag2Slug, tag3Slug, caseTag1Tag2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_tenant.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemAttr("netbox_tenant.test", "tags.*", tag1Slug),
					resource.TestCheckTypeSetElemAttr("netbox_tenant.test", "tags.*", tag2Slug),
				),
			},
			{
				Config: testAccTenantResourceConfig_tags(name, slug, tag1Slug, tag2Slug, tag3Slug, caseTag1Uscore2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_tenant.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemAttr("netbox_tenant.test", "tags.*", tag1Slug),
					resource.TestCheckTypeSetElemAttr("netbox_tenant.test", "tags.*", tag2Slug),
				),
			},
			{
				Config: testAccTenantResourceConfig_tags(name, slug, tag1Slug, tag2Slug, tag3Slug, caseTag3),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_tenant.test", "tags.#", "1"),
					resource.TestCheckTypeSetElemAttr("netbox_tenant.test", "tags.*", tag3Slug),
				),
			},
			{
				Config: testAccTenantResourceConfig_tags(name, slug, tag1Slug, tag2Slug, tag3Slug, tagsEmpty),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_tenant.test", "tags.#", "0"),
				),
			},
		},
	})
}

func TestAccTenantResource_tagOrderInvariance(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-tenant-tag-order")
	slug := testutil.RandomSlug("tf-test-tenant-tag-order")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2Slug := testutil.RandomSlug("tag2")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTenantCleanup(slug)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckTenantDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTenantResourceConfig_tagsOrder(name, slug, tag1Slug, tag2Slug, caseTag1Tag2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_tenant.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemAttr("netbox_tenant.test", "tags.*", tag1Slug),
					resource.TestCheckTypeSetElemAttr("netbox_tenant.test", "tags.*", tag2Slug),
				),
			},
			{
				Config: testAccTenantResourceConfig_tagsOrder(name, slug, tag1Slug, tag2Slug, caseTag2Uscore1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_tenant.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemAttr("netbox_tenant.test", "tags.*", tag1Slug),
					resource.TestCheckTypeSetElemAttr("netbox_tenant.test", "tags.*", tag2Slug),
				),
			},
		},
	})
}

func TestAccTenantResource_update(t *testing.T) {
	t.Parallel()

	// Generate unique names
	name := testutil.RandomName("tf-test-tenant-update")
	slug := testutil.RandomSlug("tf-test-tenant-upd")
	updatedName := testutil.RandomName("tf-test-tenant-updated")

	// Register cleanup
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTenantCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckTenantDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTenantResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_tenant.test", "id"),
					resource.TestCheckResourceAttr("netbox_tenant.test", "name", name),
				),
			},
			{
				Config:   testAccTenantResourceConfig_basic(name, slug),
				PlanOnly: true,
			},
			{
				Config: testAccTenantResourceConfig_basic(updatedName, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_tenant.test", "id"),
					resource.TestCheckResourceAttr("netbox_tenant.test", "name", updatedName),
				),
			},
			{
				Config:   testAccTenantResourceConfig_basic(updatedName, slug),
				PlanOnly: true,
			},
		},
	})
}

// testAccTenantResourceConfig_basic returns a basic test configuration.
func testAccTenantResourceConfig_basic(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_tenant" "test" {
  name = %q
  slug = %q
}
`, name, slug)
}

// testAccTenantResourceConfig_full returns a test configuration with all fields.
func testAccTenantResourceConfig_full(name, slug, groupName, groupSlug, description, comments, tagName1, tagSlug1, tagName2, tagSlug2 string) string {
	return fmt.Sprintf(`
resource "netbox_tenant_group" "test" {
	name = %[3]q
	slug = %[4]q
}

resource "netbox_tag" "tag1" {
	name = %[7]q
	slug = %[8]q
}

resource "netbox_tag" "tag2" {
	name = %[9]q
	slug = %[10]q
}

resource "netbox_tenant" "test" {
	name        = %[1]q
	slug        = %[2]q
	group       = netbox_tenant_group.test.name
	description = %[5]q
	comments    = %[6]q

	tags = [
		netbox_tag.tag1.slug,
		netbox_tag.tag2.slug
	]
}
`, name, slug, groupName, groupSlug, description, comments, tagName1, tagSlug1, tagName2, tagSlug2)
}

func testAccTenantResourceConfig_fullUpdate(name, slug, groupName, groupSlug, description, comments, tagName1, tagSlug1, tagName2, tagSlug2 string) string {
	return fmt.Sprintf(`
resource "netbox_tenant_group" "test" {
	name = %[3]q
	slug = %[4]q
}

resource "netbox_tag" "tag1" {
	name = %[7]q
	slug = %[8]q
}

resource "netbox_tag" "tag2" {
	name = %[9]q
	slug = %[10]q
}

resource "netbox_tenant" "test" {
	name        = %[1]q
	slug        = %[2]q
	group       = netbox_tenant_group.test.name
	description = %[5]q
	comments    = %[6]q

	tags = [
		netbox_tag.tag1.slug,
		netbox_tag.tag2.slug
	]
}
`, name, slug, groupName, groupSlug, description, comments, tagName1, tagSlug1, tagName2, tagSlug2)
}

func TestAccTenantResource_import(t *testing.T) {
	t.Parallel()

	// Generate unique names to avoid conflicts between test runs
	name := testutil.RandomName("tf-test-tenant-import")
	slug := testutil.RandomSlug("tf-test-tenant-imp")

	// Register cleanup
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTenantCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckTenantDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTenantResourceConfig_import(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_tenant.test", "name", name),
					resource.TestCheckResourceAttr("netbox_tenant.test", "slug", slug),
				),
			},
			{
				ResourceName:      "netbox_tenant.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:   testAccTenantResourceConfig_import(name, slug),
				PlanOnly: true,
			},
		},
	})
}

func testAccTenantResourceConfig_import(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_tenant" "test" {
  name = %q
  slug = %q
}
`, name, slug)
}

func TestAccConsistency_Tenant(t *testing.T) {
	t.Parallel()

	tenantName := testutil.RandomName("tenant")
	tenantSlug := testutil.RandomSlug("tenant")
	groupName := testutil.RandomName("group")
	groupSlug := testutil.RandomSlug("group")
	tagName1 := testutil.RandomName("tag1")
	tagSlug1 := testutil.RandomSlug("tag1")
	tagName2 := testutil.RandomName("tag2")
	tagSlug2 := testutil.RandomSlug("tag2")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTenantCleanup(tenantSlug)
	cleanup.RegisterTenantGroupCleanup(groupSlug)
	cleanup.RegisterTagCleanup(tagSlug1)
	cleanup.RegisterTagCleanup(tagSlug2)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTenantConsistencyConfig(tenantName, tenantSlug, groupName, groupSlug, tagName1, tagSlug1, tagName2, tagSlug2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_tenant.test", "name", tenantName),
					resource.TestCheckResourceAttr("netbox_tenant.test", "group", groupName),
				),
			},
			{
				PlanOnly: true,
				Config:   testAccTenantConsistencyConfig(tenantName, tenantSlug, groupName, groupSlug, tagName1, tagSlug1, tagName2, tagSlug2),
			},
		},
	})
}

func testAccTenantConsistencyConfig(tenantName, tenantSlug, groupName, groupSlug, tagName1, tagSlug1, tagName2, tagSlug2 string) string {
	return fmt.Sprintf(`
resource "netbox_tenant_group" "test" {
  name = %q
  slug = %q
}

resource "netbox_tag" "tag1" {
	name = %[5]q
	slug = %[6]q
}

resource "netbox_tag" "tag2" {
	name = %[7]q
	slug = %[8]q
}

resource "netbox_tenant" "test" {
	name  = %[3]q
	slug  = %[4]q
  group = netbox_tenant_group.test.name

	tags = [
		netbox_tag.tag1.slug,
		netbox_tag.tag2.slug
	]
}
`, groupName, groupSlug, tenantName, tenantSlug, tagName1, tagSlug1, tagName2, tagSlug2)
}

func TestAccConsistency_Tenant_LiteralNames(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-tenant-lit")
	slug := testutil.RandomSlug("tf-test-tenant-lit")
	description := testutil.RandomName("description")
	comments := testutil.RandomName("comments")
	tagName1 := testutil.RandomName("tag1")
	tagSlug1 := testutil.RandomSlug("tag1")
	tagName2 := testutil.RandomName("tag2")
	tagSlug2 := testutil.RandomSlug("tag2")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTenantCleanup(slug)
	cleanup.RegisterTagCleanup(tagSlug1)
	cleanup.RegisterTagCleanup(tagSlug2)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckTenantDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTenantConsistencyLiteralNamesConfig(name, slug, description, comments, tagName1, tagSlug1, tagName2, tagSlug2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_tenant.test", "id"),
					resource.TestCheckResourceAttr("netbox_tenant.test", "name", name),
					resource.TestCheckResourceAttr("netbox_tenant.test", "slug", slug),
					resource.TestCheckResourceAttr("netbox_tenant.test", "description", description),
					resource.TestCheckResourceAttr("netbox_tenant.test", "comments", comments),
					resource.TestCheckResourceAttr("netbox_tenant.test", "tags.#", "2"),
				),
			},
			{
				Config:   testAccTenantConsistencyLiteralNamesConfig(name, slug, description, comments, tagName1, tagSlug1, tagName2, tagSlug2),
				PlanOnly: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_tenant.test", "id"),
				),
			},
		},
	})
}

func testAccTenantConsistencyLiteralNamesConfig(name, slug, description, comments, tagName1, tagSlug1, tagName2, tagSlug2 string) string {
	return fmt.Sprintf(`
resource "netbox_tag" "tag1" {
	name = %[4]q
	slug = %[5]q
}

resource "netbox_tag" "tag2" {
	name = %[6]q
	slug = %[7]q
}

resource "netbox_tenant" "test" {
	name        = %[1]q
	slug        = %[2]q
	description = %[3]q
	comments    = %[8]q

	tags = [
		netbox_tag.tag1.slug,
		netbox_tag.tag2.slug
	]
}
`, name, slug, description, tagName1, tagSlug1, tagName2, tagSlug2, comments)
}

func TestAccTenantResource_externalDeletion(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-tenant-del")
	slug := testutil.RandomSlug("tf-test-tenant-del")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTenantCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckTenantDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTenantResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_tenant.test", "id"),
					resource.TestCheckResourceAttr("netbox_tenant.test", "name", name),
					resource.TestCheckResourceAttr("netbox_tenant.test", "slug", slug),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					items, _, err := client.TenancyAPI.TenancyTenantsList(context.Background()).Slug([]string{slug}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find tenant for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.TenancyAPI.TenancyTenantsDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete tenant: %v", err)
					}
					t.Logf("Successfully externally deleted tenant with ID: %d", itemID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

// TestAccTenantResource_removeOptionalFields tests that optional nullable fields
// can be successfully removed from the configuration without causing inconsistent state.
// This verifies the bugfix for: "Provider produced inconsistent result after apply".
func TestAccTenantResource_removeOptionalFields(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-tenant-remove")
	slug := testutil.RandomSlug("tf-test-tenant-remove")
	groupName := testutil.RandomName("tf-test-tenant-group")
	groupSlug := testutil.RandomSlug("tf-test-tenant-group")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTenantCleanup(slug)
	cleanup.RegisterTenantGroupCleanup(groupSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckTenantDestroy,
		Steps: []resource.TestStep{
			// Step 1: Create tenant with group
			{
				Config: testAccTenantResourceConfig_withGroup(name, slug, groupName, groupSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_tenant.test", "id"),
					resource.TestCheckResourceAttr("netbox_tenant.test", "name", name),
					resource.TestCheckResourceAttr("netbox_tenant.test", "slug", slug),
					resource.TestCheckResourceAttrSet("netbox_tenant.test", "group"),
				),
			},
			// Step 2: Remove group - this should set group to null
			{
				Config: testAccTenantResourceConfig_withoutGroup(name, slug, groupName, groupSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_tenant.test", "id"),
					resource.TestCheckResourceAttr("netbox_tenant.test", "name", name),
					resource.TestCheckResourceAttr("netbox_tenant.test", "slug", slug),
					resource.TestCheckNoResourceAttr("netbox_tenant.test", "group"),
				),
			},
			// Step 3: Re-add group - verify field can be set again
			{
				Config: testAccTenantResourceConfig_withGroup(name, slug, groupName, groupSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_tenant.test", "id"),
					resource.TestCheckResourceAttr("netbox_tenant.test", "name", name),
					resource.TestCheckResourceAttr("netbox_tenant.test", "slug", slug),
					resource.TestCheckResourceAttrSet("netbox_tenant.test", "group"),
				),
			},
		},
	})
}

func testAccTenantResourceConfig_withGroup(name, slug, groupName, groupSlug string) string {
	return fmt.Sprintf(`
resource "netbox_tenant_group" "test" {
  name = %q
  slug = %q
}

resource "netbox_tenant" "test" {
	name  = %q
	slug  = %q
	group = netbox_tenant_group.test.id
}
`, groupName, groupSlug, name, slug)
}

func testAccTenantResourceConfig_withoutGroup(name, slug, groupName, groupSlug string) string {
	return fmt.Sprintf(`
resource "netbox_tenant_group" "test" {
  name = %q
  slug = %q
}

resource "netbox_tenant" "test" {
	name = %q
	slug = %q
}
`, groupName, groupSlug, name, slug)
}

// TestAccTenantResource_validationErrors tests validation error scenarios.
func TestAccTenantResource_validationErrors(t *testing.T) {
	testutil.RunMultiValidationErrorTest(t, testutil.MultiValidationErrorTestConfig{
		ResourceName: "netbox_tenant",
		TestCases: map[string]testutil.ValidationErrorCase{
			"missing_name": {
				Config: func() string {
					return `
resource "netbox_tenant" "test" {
  slug = "test-tenant"
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
			"missing_slug": {
				Config: func() string {
					return `
resource "netbox_tenant" "test" {
  name = "Test Tenant"
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
			"invalid_group_reference": {
				Config: func() string {
					return `
resource "netbox_tenant" "test" {
  name  = "Test Tenant"
  slug  = "test-tenant"
  group = "99999"
}
`
				},
				ExpectedError: testutil.ErrPatternNotFound,
			},
		},
	})
}

func testAccTenantResourceConfig_tags(name, slug, tag1Slug, tag2Slug, tag3Slug, tagCase string) string {
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
  name = "Tag1-%[3]s"
  slug = %[3]q
}

resource "netbox_tag" "tag2" {
  name = "Tag2-%[4]s"
  slug = %[4]q
}

resource "netbox_tag" "tag3" {
  name = "Tag3-%[5]s"
  slug = %[5]q
}

resource "netbox_tenant" "test" {
  name = %[1]q
  slug = %[2]q
  %[6]s
}
`, name, slug, tag1Slug, tag2Slug, tag3Slug, tagsConfig)
}

func testAccTenantResourceConfig_tagsOrder(name, slug, tag1Slug, tag2Slug, tagCase string) string {
	var tagsConfig string
	switch tagCase {
	case caseTag1Tag2:
		tagsConfig = tagsDoubleSlug
	case caseTag2Uscore1:
		tagsConfig = tagsDoubleSlugReversed
	}

	return fmt.Sprintf(`
resource "netbox_tag" "tag1" {
  name = "Tag1-%[3]s"
  slug = %[3]q
}

resource "netbox_tag" "tag2" {
  name = "Tag2-%[4]s"
  slug = %[4]q
}

resource "netbox_tenant" "test" {
  name = %[1]q
  slug = %[2]q
  %[5]s
}
`, name, slug, tag1Slug, tag2Slug, tagsConfig)
}
