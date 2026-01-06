package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Acceptance tests require NETBOX_URL and NETBOX_API_TOKEN environment variables.

func TestAccContactDataSource_IDPreservation(t *testing.T) {
	t.Parallel()
	testutil.TestAccPreCheck(t)

	randomName := testutil.RandomName("test-contact-ds-id")
	email := fmt.Sprintf("%s@example.com", randomName)

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterContactCleanup(email)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckContactDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccContactDataSourceByID(randomName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_contact.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_contact.test", "name", randomName),
				),
			},
		},
	})
}

func TestAccContactDataSource_byID(t *testing.T) {
	t.Parallel()
	testutil.TestAccPreCheck(t)

	randomName := testutil.RandomName("test-contact-ds")
	email := fmt.Sprintf("%s@example.com", randomName)

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterContactCleanup(email)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckContactDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccContactDataSourceByID(randomName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_contact.test", "name", randomName),
					resource.TestCheckResourceAttrSet("data.netbox_contact.test", "id"),
				),
			},
		},
	})
}

func TestAccContactDataSource_byName(t *testing.T) {
	t.Parallel()
	testutil.TestAccPreCheck(t)

	randomName := testutil.RandomName("test-contact-ds")
	email := fmt.Sprintf("%s@example.com", randomName)

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterContactCleanup(email)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckContactDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccContactDataSourceByName(randomName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_contact.test", "name", randomName),
					resource.TestCheckResourceAttrSet("data.netbox_contact.test", "id"),
				),
			},
		},
	})
}

func TestAccContactDataSource_byEmail(t *testing.T) {
	t.Parallel()
	testutil.TestAccPreCheck(t)

	randomName := testutil.RandomName("test-contact-ds")
	email := fmt.Sprintf("%s@example.com", randomName)

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterContactCleanup(email)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckContactDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccContactDataSourceByEmail(randomName, email),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_contact.test", "name", randomName),
					resource.TestCheckResourceAttr("data.netbox_contact.test", "email", email),
					resource.TestCheckResourceAttrSet("data.netbox_contact.test", "id"),
				),
			},
		},
	})
}

func testAccContactDataSourceByID(name string) string {
	return fmt.Sprintf(`
resource "netbox_contact" "test" {
  name = %[1]q
}

data "netbox_contact" "test" {
  id = netbox_contact.test.id
}
`, name)
}

func testAccContactDataSourceByName(name string) string {
	return fmt.Sprintf(`
resource "netbox_contact" "test" {
  name = %[1]q
}

data "netbox_contact" "test" {
  name = netbox_contact.test.name
}
`, name)
}

func testAccContactDataSourceByEmail(name, email string) string {
	return fmt.Sprintf(`
resource "netbox_contact" "test" {
  name  = %[1]q
  email = %[2]q
}

data "netbox_contact" "test" {
  email = netbox_contact.test.email
}
`, name, email)
}
