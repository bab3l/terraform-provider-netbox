package resources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccContactRoleResource_basic(t *testing.T) {

	name := testutil.RandomName("test-contact-role")

	slug := testutil.GenerateSlug(name)

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

func testAccContactRoleResourceConfig(name, slug string) string {

	return fmt.Sprintf(`

resource "netbox_contact_role" "test" {

  name = %q

  slug = %q

}

`, name, slug)

}
