package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// NOTE: Custom field tests for contact assignment resource are in resources_acceptance_tests_customfields package

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

func TestAccContactAssignmentResource_withTags(t *testing.T) {
	t.Parallel()

	randomName := testutil.RandomName("test-contact-assign-tags")
	randomSlug := testutil.RandomSlug("test-ca-tags")
	contactEmail := fmt.Sprintf("%s@example.com", testutil.RandomSlug("ca-tags"))
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
				Config: testAccContactAssignmentResourceConfig_withTags(randomName, randomSlug, contactEmail, tagName, tagSlug, "primary"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_contact_assignment.test", "id"),
					resource.TestCheckResourceAttr("netbox_contact_assignment.test", "object_type", "dcim.site"),
					resource.TestCheckResourceAttr("netbox_contact_assignment.test", "priority", "primary"),
					resource.TestCheckResourceAttrSet("netbox_contact_assignment.test", "role_id"),
					resource.TestCheckResourceAttr("netbox_contact_assignment.test", "tags.#", "1"),
				),
			},
			{
				Config: testAccContactAssignmentResourceConfig_withTags(randomName, randomSlug, contactEmail, tagName, tagSlug, "secondary"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_contact_assignment.test", "priority", "secondary"),
					resource.TestCheckResourceAttr("netbox_contact_assignment.test", "tags.#", "1"),
				),
			},
		},
	})
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

func TestAccContactAssignmentResource_IDPreservation(t *testing.T) {
	t.Parallel()
	testutil.TestAccPreCheck(t)

	randomName := testutil.RandomName("tf-test-contact-assign-id")
	randomSlug := testutil.RandomSlug("tf-test-ca-id")
	contactEmail := fmt.Sprintf("%s@example.com", testutil.RandomSlug("ca-id"))

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(randomSlug + "-site")
	cleanup.RegisterContactCleanup(contactEmail)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccContactAssignmentResourceBasicWithEmail(randomName, randomSlug, contactEmail),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_contact_assignment.test", "id"),
					resource.TestCheckResourceAttr("netbox_contact_assignment.test", "object_type", "dcim.site"),
					resource.TestCheckResourceAttrSet("netbox_contact_assignment.test", "contact_id"),
					resource.TestCheckResourceAttrSet("netbox_contact_assignment.test", "object_id"),
				),
			},
		},
	})
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

func testAccContactAssignmentResourceConfig_withTags(name, slug, email, tagName, tagSlug, priority string) string {
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
  priority    = %q
  tags = [
    {
      name = netbox_tag.test.name
      slug = netbox_tag.test.slug
    }
  ]
}
`, name, slug, name, email, name, slug, tagName, tagSlug, priority)
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
