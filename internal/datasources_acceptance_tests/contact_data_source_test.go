package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/provider"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Acceptance tests require NETBOX_URL and NETBOX_API_TOKEN environment variables.

func TestAccContactDataSource_byID(t *testing.T) {

	t.Parallel()

	testutil.TestAccPreCheck(t)

	cleanup := testutil.NewCleanupResource(t)

	randomName := testutil.RandomName("test-contact-ds")

	email := fmt.Sprintf("%s@example.com", randomName)

	cleanup.RegisterContactCleanup(email)

	resource.Test(t, resource.TestCase{

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

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

	cleanup := testutil.NewCleanupResource(t)

	randomName := testutil.RandomName("test-contact-ds")

	email := fmt.Sprintf("%s@example.com", randomName)

	cleanup.RegisterContactCleanup(email)

	resource.Test(t, resource.TestCase{

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

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

	cleanup := testutil.NewCleanupResource(t)

	randomName := testutil.RandomName("test-contact-ds")

	email := fmt.Sprintf("%s@example.com", randomName)

	cleanup.RegisterContactCleanup(email)

	resource.Test(t, resource.TestCase{

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

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
