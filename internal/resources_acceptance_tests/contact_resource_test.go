package resources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccContactResource_basic(t *testing.T) {

	testutil.TestAccPreCheck(t)

	randomName := testutil.RandomName("test-contact")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccContactResource(randomName, "john.doe@example.com", "+1-555-0100"),

				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_contact.test", "name", randomName),

					resource.TestCheckResourceAttr("netbox_contact.test", "email", "john.doe@example.com"),

					resource.TestCheckResourceAttr("netbox_contact.test", "phone", "+1-555-0100"),

					resource.TestCheckResourceAttrSet("netbox_contact.test", "id"),
				),
			},

			{

				ResourceName: "netbox_contact.test",

				ImportState: true,

				ImportStateVerify: true,
			},

			{

				Config: testAccContactResource(randomName, "jane.doe@example.com", "+1-555-0200"),

				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_contact.test", "name", randomName),

					resource.TestCheckResourceAttr("netbox_contact.test", "email", "jane.doe@example.com"),

					resource.TestCheckResourceAttr("netbox_contact.test", "phone", "+1-555-0200"),
				),
			},
		},
	})

}

func TestAccContactResource_full(t *testing.T) {

	testutil.TestAccPreCheck(t)

	randomName := testutil.RandomName("test-contact-full")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccContactResourceFull(randomName),

				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_contact.test", "name", randomName),

					resource.TestCheckResourceAttr("netbox_contact.test", "title", "Network Engineer"),

					resource.TestCheckResourceAttr("netbox_contact.test", "phone", "+1-555-0100"),

					resource.TestCheckResourceAttr("netbox_contact.test", "email", "engineer@example.com"),

					resource.TestCheckResourceAttr("netbox_contact.test", "address", "123 Main Street, City, Country"),

					resource.TestCheckResourceAttr("netbox_contact.test", "link", "https://example.com/profile"),

					resource.TestCheckResourceAttr("netbox_contact.test", "description", "Test contact description"),

					resource.TestCheckResourceAttr("netbox_contact.test", "comments", "Test contact comments"),

					resource.TestCheckResourceAttrSet("netbox_contact.test", "id"),
				),
			},
		},
	})

}

func TestAccConsistency_Contact(t *testing.T) {
	contactName := testutil.RandomName("contact")

	contactGroupName := testutil.RandomName("contactgroup")

	contactGroupSlug := testutil.RandomSlug("contactgroup")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccContactConsistencyConfig(contactName, contactGroupName, contactGroupSlug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_contact.test", "name", contactName),

					resource.TestCheckResourceAttr("netbox_contact.test", "group", contactGroupName),
				),
			},

			{

				PlanOnly: true,

				Config: testAccContactConsistencyConfig(contactName, contactGroupName, contactGroupSlug),
			},
		},
	})

}

func TestAccConsistency_Contact_LiteralNames(t *testing.T) {
	contactName := testutil.RandomName("contact")

	contactGroupName := testutil.RandomName("contactgroup")

	contactGroupSlug := testutil.RandomSlug("contactgroup")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccContactConsistencyLiteralNamesConfig(contactName, contactGroupName, contactGroupSlug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_contact.test", "name", contactName),

					resource.TestCheckResourceAttr("netbox_contact.test", "group", contactGroupName),
				),
			},

			{

				PlanOnly: true,

				Config: testAccContactConsistencyLiteralNamesConfig(contactName, contactGroupName, contactGroupSlug),
			},
		},
	})

}

func testAccContactResource(name, email, phone string) string {

	return fmt.Sprintf(`

resource "netbox_contact" "test" {

  name  = %[1]q

  email = %[2]q

  phone = %[3]q

}

`, name, email, phone)

}

func testAccContactResourceFull(name string) string {

	return fmt.Sprintf(`

resource "netbox_contact" "test" {

  name        = %[1]q

  title       = "Network Engineer"

  phone       = "+1-555-0100"

  email       = "engineer@example.com"

  address     = "123 Main Street, City, Country"

  link        = "https://example.com/profile"

  description = "Test contact description"

  comments    = "Test contact comments"

}

`, name)

}

func testAccContactConsistencyConfig(contactName, contactGroupName, contactGroupSlug string) string {

	return fmt.Sprintf(`

resource "netbox_contact_group" "test" {

  name = "%[2]s"

  slug = "%[3]s"

}

resource "netbox_contact" "test" {

  name = "%[1]s"

  group = netbox_contact_group.test.name

}

`, contactName, contactGroupName, contactGroupSlug)

}

func testAccContactConsistencyLiteralNamesConfig(contactName, contactGroupName, contactGroupSlug string) string {

	return fmt.Sprintf(`

resource "netbox_contact_group" "test" {

  name = "%[2]s"

  slug = "%[3]s"

}

resource "netbox_contact" "test" {

  name = "%[1]s"

  group = "%[2]s"

  depends_on = [netbox_contact_group.test]

}

`, contactName, contactGroupName, contactGroupSlug)

}
