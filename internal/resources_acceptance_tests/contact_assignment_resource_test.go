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
	contactEmail := fmt.Sprintf("%s@example.com", testutil.RandomSlug("ca-basic"))

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

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
	contactEmail := fmt.Sprintf("%s@example.com", testutil.RandomSlug("ca-role"))

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(randomSlug + "-site")
	cleanup.RegisterContactCleanup(contactEmail)
	cleanup.RegisterContactRoleCleanup(randomSlug + "-role")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

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

func TestAccContactAssignmentResource_full(t *testing.T) {
	t.Parallel()

	randomName := testutil.RandomName("test-contact-assign-full")
	randomSlug := testutil.RandomSlug("test-ca-full")
	contactEmail := fmt.Sprintf("%s@example.com", testutil.RandomSlug("ca-full"))
	tagName := testutil.RandomName("tf-test-tag")
	tagSlug := testutil.RandomSlug("tf-test-tag")
	customFieldName := testutil.RandomCustomFieldName("tf_test_cf")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(randomSlug + "-site")
	cleanup.RegisterContactCleanup(contactEmail)
	cleanup.RegisterContactRoleCleanup(randomSlug + "-role")
	cleanup.RegisterTagCleanup(tagSlug)
	cleanup.RegisterCustomFieldCleanup(customFieldName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccContactAssignmentResourceConfig_full(randomName, randomSlug, contactEmail, tagName, tagSlug, customFieldName, "primary"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_contact_assignment.test", "id"),
					resource.TestCheckResourceAttr("netbox_contact_assignment.test", "object_type", "dcim.site"),
					resource.TestCheckResourceAttr("netbox_contact_assignment.test", "priority", "primary"),
					resource.TestCheckResourceAttrSet("netbox_contact_assignment.test", "role_id"),
					resource.TestCheckResourceAttr("netbox_contact_assignment.test", "tags.#", "1"),
					resource.TestCheckResourceAttr("netbox_contact_assignment.test", "custom_fields.#", "1"),
				),
			},
			{
				Config: testAccContactAssignmentResourceConfig_full(randomName, randomSlug, contactEmail, tagName, tagSlug, customFieldName, "secondary"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_contact_assignment.test", "priority", "secondary"),
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

		PreCheck: func() { testutil.TestAccPreCheck(t) },

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

func testAccContactAssignmentResourceWithEmail(name, slug, email string) string {

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

func testAccContactAssignmentResourceConfig_full(name, slug, email, tagName, tagSlug, customFieldName, priority string) string {
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

resource "netbox_custom_field" "test" {
  name         = %q
  object_types = ["tenancy.contactassignment"]
  type         = "text"
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
  custom_fields = [
    {
      name  = netbox_custom_field.test.name
      type  = "text"
      value = "test-value"
    }
  ]
}
`, name, slug, name, email, name, slug, tagName, tagSlug, customFieldName, priority)
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
				Config: testAccContactAssignmentResourceWithEmail(name, slug, contactEmail),
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

func TestAccContactAssignmentResource_importWithCustomFieldsAndTags(t *testing.T) {
	// NOTE: t.Parallel() intentionally omitted - this test creates/deletes global custom fields
	// that would affect other tests of the same resource type running in parallel.

	contactName := testutil.RandomName("contact")
	contactGroupName := testutil.RandomName("contact_group")
	contactGroupSlug := testutil.RandomSlug("contact_group")
	contactRoleName := testutil.RandomName("contact_role")
	contactRoleSlug := testutil.RandomSlug("contact_role")
	siteName := testutil.RandomName("site")
	siteSlug := testutil.RandomSlug("site")

	// Custom field names with underscore format
	cfText := testutil.RandomCustomFieldName("cf_text")
	cfLongtext := testutil.RandomCustomFieldName("cf_longtext")
	cfInteger := testutil.RandomCustomFieldName("cf_integer")
	cfBoolean := testutil.RandomCustomFieldName("cf_boolean")
	cfDate := testutil.RandomCustomFieldName("cf_date")
	cfUrl := testutil.RandomCustomFieldName("cf_url")
	cfJson := testutil.RandomCustomFieldName("cf_json")

	// Tag names
	tag1 := testutil.RandomName("tag1")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2 := testutil.RandomName("tag2")
	tag2Slug := testutil.RandomSlug("tag2")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccContactAssignmentResourceImportConfig_full(contactName, "", contactGroupName, contactGroupSlug, contactRoleName, contactRoleSlug, siteName, siteSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_contact_assignment.test", "id"),
					// Verify custom fields are applied
					resource.TestCheckResourceAttr("netbox_contact_assignment.test", "custom_fields.#", "7"),
					// Verify tags are applied
					resource.TestCheckResourceAttr("netbox_contact_assignment.test", "tags.#", "2"),
				),
			},
			{
				Config:                  testAccContactAssignmentResourceImportConfig_full(contactName, "", contactGroupName, contactGroupSlug, contactRoleName, contactRoleSlug, siteName, siteSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug),
				ResourceName:            "netbox_contact_assignment.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"custom_fields", "tags", "contact_id", "role_id"},
			},
			{
				Config:   testAccContactAssignmentResourceImportConfig_full(contactName, "", contactGroupName, contactGroupSlug, contactRoleName, contactRoleSlug, siteName, siteSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug),
				PlanOnly: true,
			},
		},
	})
}

func testAccContactAssignmentResourceImportConfig_full(contactName, contactSlug, contactGroupName, contactGroupSlug, contactRoleName, contactRoleSlug, siteName, siteSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug string) string {
	return fmt.Sprintf(`
# Dependencies
resource "netbox_contact_group" "test" {
  name = %q
  slug = %q
}

resource "netbox_contact" "test" {
  name  = %q
  group = netbox_contact_group.test.id
}

resource "netbox_contact_role" "test" {
  name = %q
  slug = %q
}

resource "netbox_site" "test" {
  name = %q
  slug = %q
}

# Custom Fields
resource "netbox_custom_field" "cf_text" {
  name         = %q
  type         = "text"
  object_types = ["tenancy.contactassignment"]
}

resource "netbox_custom_field" "cf_longtext" {
  name         = %q
  type         = "longtext"
  object_types = ["tenancy.contactassignment"]
}

resource "netbox_custom_field" "cf_integer" {
  name         = %q
  type         = "integer"
  object_types = ["tenancy.contactassignment"]
}

resource "netbox_custom_field" "cf_boolean" {
  name         = %q
  type         = "boolean"
  object_types = ["tenancy.contactassignment"]
}

resource "netbox_custom_field" "cf_date" {
  name         = %q
  type         = "date"
  object_types = ["tenancy.contactassignment"]
}

resource "netbox_custom_field" "cf_url" {
  name         = %q
  type         = "url"
  object_types = ["tenancy.contactassignment"]
}

resource "netbox_custom_field" "cf_json" {
  name         = %q
  type         = "json"
  object_types = ["tenancy.contactassignment"]
}

# Tags
resource "netbox_tag" "tag1" {
  name = %q
  slug = %q
}

resource "netbox_tag" "tag2" {
  name = %q
  slug = %q
}

# Main Resource
resource "netbox_contact_assignment" "test" {
  object_type = "dcim.site"
  object_id   = netbox_site.test.id
  contact_id  = netbox_contact.test.id
  role_id     = netbox_contact_role.test.id

  custom_fields = [
    {
      name  = netbox_custom_field.cf_text.name
      type  = "text"
      value = "test-value"
    },
    {
      name  = netbox_custom_field.cf_longtext.name
      type  = "longtext"
      value = "test-longtext-value"
    },
    {
      name  = netbox_custom_field.cf_integer.name
      type  = "integer"
      value = "42"
    },
    {
      name  = netbox_custom_field.cf_boolean.name
      type  = "boolean"
      value = "true"
    },
    {
      name  = netbox_custom_field.cf_date.name
      type  = "date"
      value = "2023-01-01"
    },
    {
      name  = netbox_custom_field.cf_url.name
      type  = "url"
      value = "https://example.com"
    },
    {
      name  = netbox_custom_field.cf_json.name
      type  = "json"
      value = jsonencode({"key" = "value"})
    }
  ]

  tags = [
    {
      name = netbox_tag.tag1.name
      slug = netbox_tag.tag1.slug
    },
    {
      name = netbox_tag.tag2.name
      slug = netbox_tag.tag2.slug
    }
  ]
}
`,
		contactGroupName, contactGroupSlug,
		contactName,
		contactRoleName, contactRoleSlug,
		siteName, siteSlug,
		cfText,
		cfLongtext,
		cfInteger,
		cfBoolean,
		cfDate,
		cfUrl,
		cfJson,
		tag1, tag1Slug,
		tag2, tag2Slug,
	)
}
