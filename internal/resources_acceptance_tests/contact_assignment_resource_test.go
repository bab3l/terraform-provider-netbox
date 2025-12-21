package resources_acceptance_tests

import (
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
