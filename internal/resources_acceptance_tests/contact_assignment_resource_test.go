package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccContactAssignmentResource_basic(t *testing.T) {

	t.Parallel()

	testutil.TestAccPreCheck(t)

	randomName := testutil.RandomName("test-contact-assign")

	randomSlug := testutil.RandomSlug("test-ca")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccContactAssignmentResourceBasic(randomName, randomSlug),

				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_contact_assignment.test", "object_type", "dcim.site"),

					resource.TestCheckResourceAttrSet("netbox_contact_assignment.test", "id"),

					resource.TestCheckResourceAttrSet("netbox_contact_assignment.test", "contact_id"),

					resource.TestCheckResourceAttrSet("netbox_contact_assignment.test", "object_id"),
				),
			},

			{

				ResourceName: "netbox_contact_assignment.test",

				ImportState: true,

				ImportStateVerify: true,

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

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(randomSlug + "-site")
	cleanup.RegisterContactCleanup(randomName + "-contact")
	cleanup.RegisterContactRoleCleanup(randomSlug + "-role")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccContactAssignmentResourceWithRole(randomName, randomSlug),

				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_contact_assignment.test", "object_type", "dcim.site"),

					resource.TestCheckResourceAttr("netbox_contact_assignment.test", "priority", "primary"),

					resource.TestCheckResourceAttrSet("netbox_contact_assignment.test", "role_id"),
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

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(randomSlug + "-site")
	cleanup.RegisterContactCleanup(randomName + "-contact")
	cleanup.RegisterContactRoleCleanup(randomSlug + "-role")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccContactAssignmentResourceBasic(randomName, randomSlug),

				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_contact_assignment.test", "object_type", "dcim.site"),
				),
			},

			{

				Config: testAccContactAssignmentResourceWithPriority(randomName, randomSlug, "secondary"),

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

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(slug + "-site")
	cleanup.RegisterContactCleanup(name + "-contact")
	cleanup.RegisterContactRoleCleanup(slug + "-role")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccContactAssignmentConsistencyLiteralNamesConfig(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_contact_assignment.test", "object_type", "dcim.site"),
				),
			},
			{
				PlanOnly: true,
				Config:   testAccContactAssignmentConsistencyLiteralNamesConfig(name, slug),
			},
		},
	})
}

func TestAccContactAssignmentResource_IDPreservation(t *testing.T) {
	t.Parallel()

	testutil.TestAccPreCheck(t)

	randomName := testutil.RandomName("tf-test-contact-assign-id")
	randomSlug := testutil.RandomSlug("tf-test-ca-id")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(randomSlug + "-site")
	cleanup.RegisterContactCleanup(randomName + "-contact")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccContactAssignmentResourceBasic(randomName, randomSlug),
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

func testAccContactAssignmentResourceBasic(name, slug string) string {

	return fmt.Sprintf(`

resource "netbox_site" "test" {

  name   = "%s-site"

  slug   = "%s-site"

  status = "active"

}

resource "netbox_contact" "test" {

  name  = "%s-contact"

  email = "test@example.com"

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

`, name, slug, name, name, slug)

}

func testAccContactAssignmentResourceWithRole(name, slug string) string {

	return fmt.Sprintf(`

resource "netbox_site" "test" {

  name   = "%s-site"

  slug   = "%s-site"

  status = "active"

}

resource "netbox_contact" "test" {

  name  = "%s-contact"

  email = "test@example.com"

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

`, name, slug, name, name, slug)

}

func testAccContactAssignmentResourceWithPriority(name, slug, priority string) string {

	return fmt.Sprintf(`

resource "netbox_site" "test" {

  name   = "%s-site"

  slug   = "%s-site"

  status = "active"

}

resource "netbox_contact" "test" {

  name  = "%s-contact"

  email = "test@example.com"

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

`, name, slug, name, name, slug, priority)
}

func testAccContactAssignmentConsistencyLiteralNamesConfig(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name   = %[1]q
  slug   = %[2]q
  status = "active"
}

resource "netbox_contact" "test" {
  name  = %[1]q
  email = "test@example.com"
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
`, name, slug)
}

func TestAccContactAssignmentResource_externalDeletion(t *testing.T) {
	t.Parallel()

	testutil.TestAccPreCheck(t)

	name := testutil.RandomName("tf-test-site-del")
	slug := testutil.RandomSlug("tf-test-site-del")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(slug)
	cleanup.RegisterContactCleanup("test@example.com")
	cleanup.RegisterContactRoleCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccContactAssignmentResourceBasic(name, slug),
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
					sites, _, err := client.DcimAPI.DcimSitesList(context.Background()).Slug([]string{slug}).Execute()
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
