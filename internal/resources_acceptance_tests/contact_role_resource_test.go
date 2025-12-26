package resources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccContactRoleResource_basic(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("test-contact-role")

	slug := testutil.GenerateSlug(name)

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterContactRoleCleanup(slug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccContactRoleResourceConfig(name, slug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_contact_role.test", "name", name),

					resource.TestCheckResourceAttr("netbox_contact_role.test", "slug", slug),
				),
			},

			{

				ResourceName: "netbox_contact_role.test",

				ImportState: true,

				ImportStateVerify: true,
			},

			{

				Config: testAccContactRoleResourceConfig(name+"-updated", slug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_contact_role.test", "name", name+"-updated"),
				),
			},
		},
	})

}

func TestAccContactRoleResource_IDPreservation(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("cr-id")
	slug := testutil.GenerateSlug(name)

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterContactRoleCleanup(slug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		CheckDestroy: testutil.CheckContactRoleDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccContactRoleResourceConfig(name, slug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_contact_role.test", "id"),

					resource.TestCheckResourceAttr("netbox_contact_role.test", "name", name),

					resource.TestCheckResourceAttr("netbox_contact_role.test", "slug", slug),
				),
			},
		},
	})

}

func TestAccConsistency_ContactRole_LiteralNames(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("test-contact-role-lit")

	slug := testutil.GenerateSlug(name)

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterContactRoleCleanup(slug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccContactRoleConsistencyLiteralNamesConfig(name, slug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_contact_role.test", "name", name),

					resource.TestCheckResourceAttr("netbox_contact_role.test", "slug", slug),
				),
			},

			{

				Config: testAccContactRoleConsistencyLiteralNamesConfig(name, slug),

				PlanOnly: true,

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_contact_role.test", "id"),
				),
			},
		},
	})

}

func testAccContactRoleConsistencyLiteralNamesConfig(name, slug string) string {

	return fmt.Sprintf(`

resource "netbox_contact_role" "test" {

  name = %q

  slug = %q

}

`, name, slug)

}

func testAccContactRoleResourceConfig(name, slug string) string {

	return fmt.Sprintf(`

resource "netbox_contact_role" "test" {

  name = %q

  slug = %q

}

`, name, slug)

}
