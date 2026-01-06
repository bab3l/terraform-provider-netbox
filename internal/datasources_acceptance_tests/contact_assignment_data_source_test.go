package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Acceptance tests require NETBOX_URL and NETBOX_API_TOKEN environment variables.

func TestAccContactAssignmentDataSource_IDPreservation(t *testing.T) {
	t.Parallel()
	testutil.TestAccPreCheck(t)

	randomName := testutil.RandomName("test-contact-ds-id")
	randomSlug := testutil.RandomSlug("test-ca-ds-id")
	siteSlug := testutil.RandomSlug("site-id")
	contactEmail := fmt.Sprintf("%s@example.com", testutil.RandomSlug("ca-ds-id"))

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterContactCleanup(contactEmail)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckSiteDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccContactAssignmentDataSourceConfigWithEmail(randomName, randomSlug, siteSlug, contactEmail),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_contact_assignment.test", "id"),
				),
			},
		},
	})
}

// Acceptance tests require NETBOX_URL and NETBOX_API_TOKEN environment variables.

func TestAccContactAssignmentDataSource_byID(t *testing.T) {
	t.Parallel()
	testutil.TestAccPreCheck(t)

	randomName := testutil.RandomName("test-contact-ds")
	randomSlug := testutil.RandomSlug("test-ca-ds")
	siteSlug := testutil.RandomSlug("site")
	contactEmail := fmt.Sprintf("%s@example.com", testutil.RandomSlug("ca-ds"))

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterContactCleanup(contactEmail)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckSiteDestroy,
		),
		Steps: []resource.TestStep{
			// Create resource and read via data source
			{
				Config: testAccContactAssignmentDataSourceConfigWithEmail(randomName, randomSlug, siteSlug, contactEmail),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"data.netbox_contact_assignment.test", "id",
						"netbox_contact_assignment.test", "id"),
					resource.TestCheckResourceAttrPair(
						"data.netbox_contact_assignment.test", "contact_id",
						"netbox_contact_assignment.test", "contact_id"),
					resource.TestCheckResourceAttr(
						"data.netbox_contact_assignment.test", "object_type", "dcim.site"),
				),
			},
		},
	})
}

func testAccContactAssignmentDataSourceConfigWithEmail(name, slug, siteSlug, email string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name   = "%s"
  slug   = "%s"
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

data "netbox_contact_assignment" "test" {
  id = netbox_contact_assignment.test.id
}
`, siteSlug, siteSlug, name, email, name, slug)
}
