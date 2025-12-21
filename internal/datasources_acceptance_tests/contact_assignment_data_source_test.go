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

func TestAccContactAssignmentDataSource_basic(t *testing.T) {

	t.Parallel()

	testutil.TestAccPreCheck(t)

	randomName := testutil.RandomName("test-contact-ds")

	randomSlug := testutil.RandomSlug("test-ca-ds")

	siteSlug := testutil.RandomSlug("site")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.ComposeCheckDestroy(

			testutil.CheckSiteDestroy,
		),

		Steps: []resource.TestStep{

			// Create resource and read via data source

			{

				Config: testAccContactAssignmentDataSourceConfig(randomName, randomSlug, siteSlug),

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

func testAccContactAssignmentDataSourceConfig(name, slug, siteSlug string) string {

	return fmt.Sprintf(`

resource "netbox_site" "test" {

  name   = "%s"

  slug   = "%s"

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

data "netbox_contact_assignment" "test" {

  id = netbox_contact_assignment.test.id

}

`, siteSlug, siteSlug, name, name, slug)

}
