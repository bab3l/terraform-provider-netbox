package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPowerPanelDataSource_basic(t *testing.T) {

	t.Parallel()

	cleanup := testutil.NewCleanupResource(t)

	siteName := testutil.RandomName("tf-test-power-panel-site")
	siteSlug := testutil.GenerateSlug(siteName)
	powerPanelName := testutil.RandomName("tf-test-power-panel")

	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterPowerPanelCleanup(powerPanelName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckSiteDestroy,
			testutil.CheckPowerPanelDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccPowerPanelDataSourceConfig(siteName, siteSlug, powerPanelName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_power_panel.test", "name", powerPanelName),
					resource.TestCheckResourceAttrSet("data.netbox_power_panel.test", "site"),
				),
			},
		},
	})
}

func testAccPowerPanelDataSourceConfig(siteName, siteSlug, powerPanelName string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_power_panel" "test" {
  site = netbox_site.test.id
  name = %[3]q
}

data "netbox_power_panel" "test" {
  id = netbox_power_panel.test.id
}
`, siteName, siteSlug, powerPanelName)
}
