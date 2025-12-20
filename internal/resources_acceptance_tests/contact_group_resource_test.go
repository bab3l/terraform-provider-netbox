package resources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccContactGroupResource_basic(t *testing.T) {

	name := testutil.RandomName("test-contact-group")

	slug := testutil.GenerateSlug(name)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccContactGroupResourceConfig(name, slug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_contact_group.test", "name", name),

					resource.TestCheckResourceAttr("netbox_contact_group.test", "slug", slug),
				),
			},

			{

				ResourceName: "netbox_contact_group.test",

				ImportState: true,

				ImportStateVerify: true,
			},

			{

				Config: testAccContactGroupResourceConfig(name+"-updated", slug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_contact_group.test", "name", name+"-updated"),
				),
			},
		},
	})

}

func testAccContactGroupResourceConfig(name, slug string) string {

	return fmt.Sprintf(`



resource "netbox_contact_group" "test" {



  name = %q



  slug = %q



}



`, name, slug)

}
