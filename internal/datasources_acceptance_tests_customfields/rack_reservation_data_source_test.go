//go:build customfields
// +build customfields

package datasources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRackReservationDataSource_customFields(t *testing.T) {
	customFieldName := testutil.RandomCustomFieldName("tf_test_rackreserv_ds_cf")
	reservationDesc := testutil.RandomName("tf-test-reservation-ds-cf")
	rackName := testutil.RandomName("tf-test-rack-ds-cf")
	siteName := testutil.RandomName("tf-test-site-ds-cf")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRackReservationDataSourceConfig_customFields(customFieldName, reservationDesc, rackName, siteName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_rack_reservation.test", "description", reservationDesc),
					resource.TestCheckResourceAttr("data.netbox_rack_reservation.test", "custom_fields.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_rack_reservation.test", "custom_fields.0.name", customFieldName),
					resource.TestCheckResourceAttr("data.netbox_rack_reservation.test", "custom_fields.0.type", "text"),
					resource.TestCheckResourceAttr("data.netbox_rack_reservation.test", "custom_fields.0.value", "datasource-test-value"),
				),
			},
		},
	})
}

func testAccRackReservationDataSourceConfig_customFields(customFieldName, reservationDesc, rackName, siteName string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "test" {
  name         = %q
  object_types = ["dcim.rackreservation"]
  type         = "text"
}

resource "netbox_site" "test" {
  name = %q
  slug = %q
}

resource "netbox_rack" "test" {
  name = %q
  site = netbox_site.test.id
}

data "netbox_user" "admin" {
  username = "admin"
}

resource "netbox_rack_reservation" "test" {
  rack        = netbox_rack.test.id
  units       = [1, 2, 3]
  description = %q
  user        = data.netbox_user.admin.id

  custom_fields = [
    {
      name  = netbox_custom_field.test.name
      type  = "text"
      value = "datasource-test-value"
    }
  ]
}

data "netbox_rack_reservation" "test" {
  id = netbox_rack_reservation.test.id

  depends_on = [netbox_rack_reservation.test]
}
`, customFieldName, siteName, siteName, rackName, reservationDesc)
}
