package resources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/provider"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRackReservationResource_basic(t *testing.T) {

	t.Parallel()

	siteName := testutil.RandomName("tf-test-site")

	siteSlug := testutil.RandomSlug("tf-test-site")

	rackName := testutil.RandomName("tf-test-rack")

	description := testutil.RandomName("description")

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterSiteCleanup(siteSlug)

	cleanup.RegisterRackCleanup(rackName)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckRackReservationDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccRackReservationResourceConfig_basic(siteName, siteSlug, rackName, description),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_rack_reservation.test", "id"),

					resource.TestCheckResourceAttr("netbox_rack_reservation.test", "description", description),

					resource.TestCheckResourceAttr("netbox_rack_reservation.test", "units.#", "2"),
				),
			},

			// ImportState test

			{

				ResourceName: "netbox_rack_reservation.test",

				ImportState: true,

				ImportStateVerify: true,
			},
		},
	})

}

func TestAccRackReservationResource_update(t *testing.T) {

	t.Parallel()

	siteName := testutil.RandomName("tf-test-site")

	siteSlug := testutil.RandomSlug("tf-test-site")

	rackName := testutil.RandomName("tf-test-rack")

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterSiteCleanup(siteSlug)

	cleanup.RegisterRackCleanup(rackName)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckRackReservationDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccRackReservationResourceConfig_basic(siteName, siteSlug, rackName, testutil.Description1),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_rack_reservation.test", "id"),

					resource.TestCheckResourceAttr("netbox_rack_reservation.test", "description", testutil.Description1),
				),
			},

			{

				Config: testAccRackReservationResourceConfig_basic(siteName, siteSlug, rackName, testutil.Description2),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_rack_reservation.test", "description", testutil.Description2),
				),
			},
		},
	})

}

func testAccRackReservationResourceConfig_basic(siteName, siteSlug, rackName, description string) string {

	return fmt.Sprintf(`

provider "netbox" {}

resource "netbox_site" "test" {

  name   = %[1]q

  slug   = %[2]q

  status = "active"

}

resource "netbox_rack" "test" {

  name     = %[3]q

  site     = netbox_site.test.id

  status   = "active"

  u_height = 42

}

data "netbox_user" "admin" {

  username = "admin"

}

resource "netbox_rack_reservation" "test" {

  rack        = netbox_rack.test.id

  units       = [1, 2]

  user        = data.netbox_user.admin.id

  description = %[4]q

}

`, siteName, siteSlug, rackName, description)

}

func TestAccConsistency_RackReservation_LiteralNames(t *testing.T) {
	t.Parallel()
	siteName := testutil.RandomName("site")
	siteSlug := testutil.RandomSlug("site")
	rackName := testutil.RandomName("rack")
	description := testutil.RandomName("description")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterRackCleanup(rackName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccRackReservationConsistencyLiteralNamesConfig(siteName, siteSlug, rackName, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_rack_reservation.test", "id"),
				),
			},
			{
				PlanOnly: true,
				Config:   testAccRackReservationConsistencyLiteralNamesConfig(siteName, siteSlug, rackName, description),
			},
		},
	})
}
func testAccRackReservationConsistencyLiteralNamesConfig(siteName, siteSlug, rackName, description string) string {
	return fmt.Sprintf(`
provider "netbox" {}

resource "netbox_site" "test" {
  name   = %[1]q
  slug   = %[2]q
  status = "active"
}

resource "netbox_rack" "test" {
  name     = %[3]q
  site     = netbox_site.test.id
  status   = "active"
  u_height = 42
}

data "netbox_user" "admin" {
  username = "admin"
}

resource "netbox_rack_reservation" "test" {
  rack        = netbox_rack.test.name
  units       = [1, 2]
  user        = data.netbox_user.admin.id
  description = %[4]q
}
`, siteName, siteSlug, rackName, description)
}
