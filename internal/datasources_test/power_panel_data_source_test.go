package datasources_test

import (
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPowerPanelDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPowerPanelDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_power_panel.test", "name", "Test Power Panel"),
					resource.TestCheckResourceAttrSet("data.netbox_power_panel.test", "site"),
				),
			},
		},
	})
}

const testAccPowerPanelDataSourceConfig = `
resource "netbox_site" "test" {
  name = "Test Site"
  slug = "test-site"
}

resource "netbox_power_panel" "test" {
  site = netbox_site.test.id
  name = "Test Power Panel"
}

data "netbox_power_panel" "test" {
  id = netbox_power_panel.test.id
}
`
