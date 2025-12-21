package datasources_acceptance_tests

import (
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPowerFeedDataSource_basic(t *testing.T) {

	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPowerFeedDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_power_feed.test", "name", "Test Power Feed"),
					resource.TestCheckResourceAttr("data.netbox_power_feed.test", "status", "active"),
				),
			},
		},
	})
}

const testAccPowerFeedDataSourceConfig = `
resource "netbox_site" "test" {
  name = "Test Site"
  slug = "test-site"
}

resource "netbox_power_panel" "test" {
  site = netbox_site.test.id
  name = "Test Power Panel"
}

resource "netbox_power_feed" "test" {
  power_panel = netbox_power_panel.test.id
  name        = "Test Power Feed"
  status      = "active"
  type        = "primary"
  supply      = "ac"
  phase       = "single-phase"
  voltage     = 120
  amperage    = 15
}

data "netbox_power_feed" "test" {
  id = netbox_power_feed.test.id
}
`
