package datasources_acceptance_tests

import (
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRackTypeDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRackTypeDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_rack_type.test", "model", "Test Rack Type"),
					resource.TestCheckResourceAttr("data.netbox_rack_type.test", "slug", "test-rack-type"),
				),
			},
		},
	})
}

const testAccRackTypeDataSourceConfig = `
resource "netbox_manufacturer" "test" {
  name = "Test Manufacturer"
  slug = "test-manufacturer"
}

resource "netbox_rack_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = "Test Rack Type"
  slug         = "test-rack-type"
  width        = 19
  u_height     = 42
  form_factor  = "2-post-frame"
}

data "netbox_rack_type" "test" {
  id = netbox_rack_type.test.id
}
`
