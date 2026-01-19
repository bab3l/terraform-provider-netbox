package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// NOTE: Custom field tests for contact assignment resource are in resources_acceptance_tests_customfields package

func TestAccContactAssignmentResource_full(t *testing.T) {
	t.Parallel()

	testutil.TestAccPreCheck(t)
	randomName := testutil.RandomName("test-contact-assign-full")
	randomSlug := testutil.RandomSlug("test-ca-full")
	contactEmail := fmt.Sprintf("%s@example.com", testutil.RandomSlug("ca-full"))
	tagName := testutil.RandomName("tf-test-tag")
	tagSlug := testutil.RandomSlug("tf-test-tag")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(randomSlug + "-site")
	cleanup.RegisterContactCleanup(contactEmail)
	cleanup.RegisterContactRoleCleanup(randomSlug + "-role")
	cleanup.RegisterTagCleanup(tagSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccContactAssignmentResourceConfig_full(randomName, randomSlug, contactEmail, tagName, tagSlug),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_contact_assignment.test", "id"),
					resource.TestCheckResourceAttr("netbox_contact_assignment.test", "object_type", "dcim.site"),
					resource.TestCheckResourceAttrSet("netbox_contact_assignment.test", "object_id"),
					resource.TestCheckResourceAttrSet("netbox_contact_assignment.test", "contact_id"),
					resource.TestCheckResourceAttrSet("netbox_contact_assignment.test", "role_id"),
					resource.TestCheckResourceAttr("netbox_contact_assignment.test", "priority", "primary"),
					resource.TestCheckResourceAttr("netbox_contact_assignment.test", "tags.#", "1"),
					resource.TestCheckTypeSetElemAttr("netbox_contact_assignment.test", "tags.*", tagSlug),
				),
			},
		},
	})
}

func TestAccContactAssignmentResource_basic(t *testing.T) {
	t.Parallel()

	testutil.TestAccPreCheck(t)
	randomName := testutil.RandomName("test-contact-assign")
	randomSlug := testutil.RandomSlug("test-ca")
	contactEmail := fmt.Sprintf("%s@example.com", testutil.RandomSlug("ca-basic"))

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(randomSlug + "-site")
	cleanup.RegisterContactCleanup(contactEmail)
	cleanup.RegisterContactRoleCleanup(randomSlug + "-role")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccContactAssignmentResourceBasicWithEmail(randomName, randomSlug, contactEmail),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_contact_assignment.test", "object_type", "dcim.site"),
					resource.TestCheckResourceAttrSet("netbox_contact_assignment.test", "id"),
					resource.TestCheckResourceAttrSet("netbox_contact_assignment.test", "contact_id"),
					resource.TestCheckResourceAttrSet("netbox_contact_assignment.test", "object_id"),
				),
			},
			{
				ResourceName:            "netbox_contact_assignment.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"contact_id", "role_id"},
			},
			{
				Config:   testAccContactAssignmentResourceBasicWithEmail(randomName, randomSlug, contactEmail),
				PlanOnly: true,
			},
		},
	})
}

func TestAccContactAssignmentResource_withRole(t *testing.T) {
	t.Parallel()

	testutil.TestAccPreCheck(t)
	randomName := testutil.RandomName("test-contact-assign")
	randomSlug := testutil.RandomSlug("test-ca")
	contactEmail := fmt.Sprintf("%s@example.com", testutil.RandomSlug("ca-role"))

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(randomSlug + "-site")
	cleanup.RegisterContactCleanup(contactEmail)
	cleanup.RegisterContactRoleCleanup(randomSlug + "-role")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccContactAssignmentResourceWithRoleEmail(randomName, randomSlug, contactEmail),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_contact_assignment.test", "object_type", "dcim.site"),
					resource.TestCheckResourceAttr("netbox_contact_assignment.test", "priority", "primary"),
					resource.TestCheckResourceAttrSet("netbox_contact_assignment.test", "role_id"),
				),
			},
		},
	})
}

func TestAccContactAssignmentResource_tagLifecycle(t *testing.T) {
	t.Parallel()

	randomName := testutil.RandomName("test-contact-assign-tags")
	randomSlug := testutil.RandomSlug("test-ca-tags")
	contactEmail := fmt.Sprintf("%s@example.com", testutil.RandomSlug("ca-tags"))
	tag1Name := testutil.RandomName("tag1")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2Name := testutil.RandomName("tag2")
	tag2Slug := testutil.RandomSlug("tag2")
	tag3Name := testutil.RandomName("tag3")
	tag3Slug := testutil.RandomSlug("tag3")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(randomSlug + "-site")
	cleanup.RegisterContactCleanup(contactEmail)
	cleanup.RegisterContactRoleCleanup(randomSlug + "-role")
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)
	cleanup.RegisterTagCleanup(tag3Slug)

	testutil.RunTagLifecycleTest(t, testutil.TagLifecycleTestConfig{
		ResourceName: "netbox_contact_assignment",
		ConfigWithoutTags: func() string {
			return testAccContactAssignmentResourceConfig_tagLifecycle(randomName, randomSlug, contactEmail, "", "", "", "", "", "", "")
		},
		ConfigWithTags: func() string {
			return testAccContactAssignmentResourceConfig_tagLifecycle(randomName, randomSlug, contactEmail, tag1Name, tag1Slug, tag2Name, tag2Slug, "tag1,tag2", "", "")
		},
		ConfigWithDifferentTags: func() string {
			return testAccContactAssignmentResourceConfig_tagLifecycle(randomName, randomSlug, contactEmail, tag1Name, tag1Slug, tag2Name, tag2Slug, "tag3", tag3Name, tag3Slug)
		},
		ExpectedTagCount:          2,
		ExpectedDifferentTagCount: 1,
	})
}

func testAccContactAssignmentResourceConfig_tagLifecycle(name, slug, email, tag1Name, tag1Slug, tag2Name, tag2Slug, tagSet, tag3Name, tag3Slug string) string {
	tagResources := ""
	tagsList := ""

	if tag1Name != "" {
		tagResources += fmt.Sprintf(`
resource "netbox_tag" "tag1" {
  name = %q
  slug = %q
}
`, tag1Name, tag1Slug)
	}
	if tag2Name != "" {
		tagResources += fmt.Sprintf(`
resource "netbox_tag" "tag2" {
  name = %q
  slug = %q
}
`, tag2Name, tag2Slug)
	}
	if tag3Name != "" {
		tagResources += fmt.Sprintf(`
resource "netbox_tag" "tag3" {
  name = %q
  slug = %q
}
`, tag3Name, tag3Slug)
	}

	if tagSet != "" {
		switch tagSet {
		case caseTag1Tag2:
			tagsList = tagsDoubleSlug
		case caseTag3:
			tagsList = tagsSingleSlug
		default:
			tagsList = tagsEmpty
		}
	} else {
		tagsList = tagsEmpty
	}

	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name   = "%s-site"
  slug   = "%s-site"
  status = "active"
}

resource "netbox_contact" "test" {
  name  = "%s-contact"
  email = "%s"
}

resource "netbox_contact_role" "test" {
  name = "%s-role"
  slug = "%s-role"
}

%s

resource "netbox_contact_assignment" "test" {
  object_type = "dcim.site"
  object_id   = netbox_site.test.id
  contact_id  = netbox_contact.test.id
  role_id     = netbox_contact_role.test.id
  priority    = "primary"
%s
}
`, name, slug, name, email, name, slug, tagResources, tagsList)
}

func TestAccContactAssignmentResource_tagOrderInvariance(t *testing.T) {
	t.Parallel()

	randomName := testutil.RandomName("test-contact-assign-tag-order")
	randomSlug := testutil.RandomSlug("test-ca-tag-order")
	contactEmail := fmt.Sprintf("%s@example.com", testutil.RandomSlug("ca-tag-order"))
	tag1Name := testutil.RandomName("tag1")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2Name := testutil.RandomName("tag2")
	tag2Slug := testutil.RandomSlug("tag2")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(randomSlug + "-site")
	cleanup.RegisterContactCleanup(contactEmail)
	cleanup.RegisterContactRoleCleanup(randomSlug + "-role")
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)

	testutil.RunTagOrderTest(t, testutil.TagOrderTestConfig{
		ResourceName: "netbox_contact_assignment",
		ConfigWithTagsOrderA: func() string {
			return testAccContactAssignmentResourceConfig_tagOrder(randomName, randomSlug, contactEmail, tag1Name, tag1Slug, tag2Name, tag2Slug, true)
		},
		ConfigWithTagsOrderB: func() string {
			return testAccContactAssignmentResourceConfig_tagOrder(randomName, randomSlug, contactEmail, tag1Name, tag1Slug, tag2Name, tag2Slug, false)
		},
		ExpectedTagCount: 2,
	})
}

func testAccContactAssignmentResourceConfig_tagOrder(name, slug, email, tag1Name, tag1Slug, tag2Name, tag2Slug string, tag1First bool) string {
	tagsOrder := tagsDoubleSlug
	if !tag1First {
		tagsOrder = tagsDoubleSlugReversed
	}

	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name   = "%s-site"
  slug   = "%s-site"
  status = "active"
}

resource "netbox_contact" "test" {
  name  = "%s-contact"
  email = "%s"
}

resource "netbox_contact_role" "test" {
  name = "%s-role"
  slug = "%s-role"
}

resource "netbox_tag" "tag1" {
  name = %q
  slug = %q
}

resource "netbox_tag" "tag2" {
  name = %q
  slug = %q
}

resource "netbox_contact_assignment" "test" {
  object_type = "dcim.site"
  object_id   = netbox_site.test.id
  contact_id  = netbox_contact.test.id
  role_id     = netbox_contact_role.test.id
  priority    = "primary"
  %s
}
`, name, slug, name, email, name, slug, tag1Name, tag1Slug, tag2Name, tag2Slug, tagsOrder)
}

func TestAccContactAssignmentResource_update(t *testing.T) {
	t.Parallel()

	testutil.TestAccPreCheck(t)
	randomName := testutil.RandomName("test-contact-assign")
	randomSlug := testutil.RandomSlug("test-ca")
	contactEmail := fmt.Sprintf("%s@example.com", testutil.RandomSlug("ca-update"))

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(randomSlug + "-site")
	cleanup.RegisterContactCleanup(contactEmail)
	cleanup.RegisterContactRoleCleanup(randomSlug + "-role")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccContactAssignmentResourceBasicWithEmail(randomName, randomSlug, contactEmail),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_contact_assignment.test", "object_type", "dcim.site"),
				),
			},
			{
				Config: testAccContactAssignmentResourceWithPriorityEmail(randomName, randomSlug, contactEmail, "secondary"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_contact_assignment.test", "priority", "secondary"),
				),
			},
		},
	})
}

func TestAccConsistency_ContactAssignment_LiteralNames(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("test-ca")
	slug := testutil.RandomSlug("test-ca")
	contactEmail := fmt.Sprintf("%s@example.com", testutil.RandomSlug("ca-consistency"))

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(slug + "-site")
	cleanup.RegisterContactCleanup(contactEmail)
	cleanup.RegisterContactRoleCleanup(slug + "-role")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccContactAssignmentConsistencyLiteralNamesConfigWithEmail(name, slug, contactEmail),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_contact_assignment.test", "object_type", "dcim.site"),
				),
			},
			{
				PlanOnly: true,
				Config:   testAccContactAssignmentConsistencyLiteralNamesConfigWithEmail(name, slug, contactEmail),
			},
		},
	})
}

func testAccContactAssignmentResourceConfig_full(name, slug, email, tagName, tagSlug string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name   = "%s-site"
  slug   = "%s-site"
  status = "active"
}

resource "netbox_contact" "test" {
  name  = "%s-contact"
  email = "%s"
}

resource "netbox_contact_role" "test" {
  name = "%s-role"
  slug = "%s-role"
}

resource "netbox_tag" "test" {
  name = %q
  slug = %q
}

resource "netbox_contact_assignment" "test" {
  object_type = "dcim.site"
  object_id   = netbox_site.test.id
  contact_id  = netbox_contact.test.id
  role_id     = netbox_contact_role.test.id
  priority    = "primary"
	tags = [netbox_tag.test.slug]
}
`, name, slug, name, email, name, slug, tagName, tagSlug)
}

func testAccContactAssignmentResourceBasicWithEmail(name, slug, email string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name   = "%s-site"
  slug   = "%s-site"
  status = "active"
}

resource "netbox_contact" "test" {
  name  = "%s-contact"
  email = "%s"
}

resource "netbox_contact_role" "test" {
  name = "%s-role"
  slug = "%s-role"
}

resource "netbox_contact_assignment" "test" {
  object_type = "dcim.site"
  object_id   = netbox_site.test.id
  contact_id  = netbox_contact.test.id
  role_id     = netbox_contact_role.test.id
}
`, name, slug, name, email, name, slug)
}

func testAccContactAssignmentResourceWithRoleEmail(name, slug, email string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name   = "%s-site"
  slug   = "%s-site"
  status = "active"
}

resource "netbox_contact" "test" {
  name  = "%s-contact"
  email = "%s"
}

resource "netbox_contact_role" "test" {
  name = "%s-role"
  slug = "%s-role"
}

resource "netbox_contact_assignment" "test" {
  object_type = "dcim.site"
  object_id   = netbox_site.test.id
  contact_id  = netbox_contact.test.id
  role_id     = netbox_contact_role.test.id
  priority    = "primary"
}
`, name, slug, name, email, name, slug)
}

func testAccContactAssignmentResourceWithPriorityEmail(name, slug, email, priority string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name   = "%s-site"
  slug   = "%s-site"
  status = "active"
}

resource "netbox_contact" "test" {
  name  = "%s-contact"
  email = "%s"
}

resource "netbox_contact_role" "test" {
  name = "%s-role"
  slug = "%s-role"
}

resource "netbox_contact_assignment" "test" {
  object_type = "dcim.site"
  object_id   = netbox_site.test.id
  contact_id  = netbox_contact.test.id
  role_id     = netbox_contact_role.test.id
  priority    = "%s"
}
`, name, slug, name, email, name, slug, priority)
}

func testAccContactAssignmentConsistencyLiteralNamesConfigWithEmail(name, slug, email string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name   = %[1]q
  slug   = %[2]q
  status = "active"
}

resource "netbox_contact" "test" {
  name  = %[1]q
  email = %[3]q
}

resource "netbox_contact_role" "test" {
  name = %[2]q
  slug = %[2]q
}

resource "netbox_contact_assignment" "test" {
  object_type = "dcim.site"
  object_id   = netbox_site.test.id
  contact_id  = netbox_contact.test.id
  role_id     = netbox_contact_role.test.id
}
`, name, slug, email)
}

func TestAccContactAssignmentResource_externalDeletion(t *testing.T) {
	t.Parallel()

	testutil.TestAccPreCheck(t)

	name := testutil.RandomName("tf-test-site-del")
	slug := testutil.RandomSlug("tf-test-site-del")
	contactEmail := fmt.Sprintf("%s@example.com", testutil.RandomSlug("ca-del"))

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(slug + "-site")
	cleanup.RegisterContactCleanup(contactEmail)
	cleanup.RegisterContactRoleCleanup(slug + "-role")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccContactAssignmentResourceBasicWithEmail(name, slug, contactEmail),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_contact_assignment.test", "id"),
					resource.TestCheckResourceAttr("netbox_contact_assignment.test", "object_type", "dcim.site"),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					// Get site ID to filter assignments
					sites, _, err := client.DcimAPI.DcimSitesList(context.Background()).Slug([]string{slug + "-site"}).Execute()
					if err != nil || sites == nil || len(sites.Results) == 0 {
						t.Fatalf("Failed to find site for external deletion: %v", err)
					}
					siteID := sites.Results[0].Id

					items, _, err := client.TenancyAPI.TenancyContactAssignmentsList(context.Background()).ObjectId([]int32{siteID}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find contact_assignment for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.TenancyAPI.TenancyContactAssignmentsDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete contact_assignment: %v", err)
					}
					t.Logf("Successfully externally deleted contact_assignment with ID: %d", itemID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		}})
}

// TestAccContactAssignmentResource_removeOptionalFields tests that optional fields
// can be successfully removed from the configuration without causing inconsistent state.
func TestAccContactAssignmentResource_removeOptionalFields(t *testing.T) {
	t.Parallel()

	testutil.TestAccPreCheck(t)
	randomName := testutil.RandomName("test-contact-rem")
	randomSlug := testutil.RandomSlug("test-ca-rem")
	contactEmail := fmt.Sprintf("%s@example.com", testutil.RandomSlug("ca-rem"))

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(randomSlug + "-site")
	cleanup.RegisterContactCleanup(contactEmail)
	cleanup.RegisterContactRoleCleanup(randomSlug + "-role")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccContactAssignmentResourceWithPriorityEmail(randomName, randomSlug, contactEmail, "primary"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_contact_assignment.test", "object_type", "dcim.site"),
					resource.TestCheckResourceAttrSet("netbox_contact_assignment.test", "contact_id"),
					resource.TestCheckResourceAttrSet("netbox_contact_assignment.test", "role_id"),
					resource.TestCheckResourceAttr("netbox_contact_assignment.test", "priority", "primary"),
				),
			},
			{
				Config: testAccContactAssignmentResourceBasicWithEmail(randomName, randomSlug, contactEmail),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_contact_assignment.test", "object_type", "dcim.site"),
					resource.TestCheckResourceAttrSet("netbox_contact_assignment.test", "contact_id"),
					resource.TestCheckNoResourceAttr("netbox_contact_assignment.test", "priority"),
					// Note: role_id is present in basic config so we test removal separately
				),
			},
		},
	})
}

func TestAccContactAssignmentResource_validationErrors(t *testing.T) {
	testutil.RunMultiValidationErrorTest(t, testutil.MultiValidationErrorTestConfig{
		ResourceName: "netbox_contact_assignment",
		TestCases: map[string]testutil.ValidationErrorCase{
			"missing_object_type": {
				Config: func() string {
					return `
provider "netbox" {}

resource "netbox_contact_assignment" "test" {
  # object_type missing
  object_id = "1"
  contact_id = "1"
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
			"missing_object_id": {
				Config: func() string {
					return `
provider "netbox" {}

resource "netbox_contact_assignment" "test" {
  object_type = "dcim.site"
  # object_id missing
  contact_id = "1"
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
			"missing_contact_id": {
				Config: func() string {
					return `
provider "netbox" {}

resource "netbox_contact_assignment" "test" {
  object_type = "dcim.site"
  object_id = "1"
  # contact_id missing
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
		},
	})
}
