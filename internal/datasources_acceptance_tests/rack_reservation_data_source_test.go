package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRackReservationDataSource_byID(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("test-site-rr")
	siteSlug := testutil.GenerateSlug(siteName)
	rackName := testutil.RandomName("test-rack-rr")
	description := "Test Rack Reservation Description"

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterRackCleanup(rackName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRackReservationDataSourceConfig(siteName, siteSlug, rackName, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_rack_reservation.test", "description", description),
					resource.TestCheckResourceAttr("data.netbox_rack_reservation.test", "units.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_rack_reservation.test", "units.0", "1"),
				),
			},
		},
	})
}

func TestAccRackReservationDataSource_IDPreservation(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("test-site-rr-id")
	siteSlug := testutil.GenerateSlug(siteName)
	rackName := testutil.RandomName("test-rack-rr-id")
	description := "Test Rack Reservation IDPreservation"

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterRackCleanup(rackName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRackReservationDataSourceConfig(siteName, siteSlug, rackName, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_rack_reservation.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_rack_reservation.test", "description", description),
				),
			},
		},
	})
}

func testAccRackReservationDataSourceConfig(siteName, siteSlug, rackName, description string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name = %q
  slug = %q
  status = "active"
}

resource "netbox_rack" "test" {
  name = %q
  site = netbox_site.test.id
  status = "active"
  width = 19
  u_height = 42
}

resource "netbox_rack_reservation" "test" {
  rack        = netbox_rack.test.id
  units       = [1]
  description = %q
  user        = 1 # Assuming user ID 1 exists (admin)
}

data "netbox_rack_reservation" "test" {
  id = netbox_rack_reservation.test.id
}
`, siteName, siteSlug, rackName, description)
}
