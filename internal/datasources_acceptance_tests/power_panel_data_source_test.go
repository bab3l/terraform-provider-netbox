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

	siteName := testutil.RandomName("test-power-panel-site-ds")
	siteSlug := testutil.GenerateSlug(siteName)

	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterPowerPanelCleanup("Test Power Panel")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckSiteDestroy,
			testutil.CheckPowerPanelDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccPowerPanelDataSourceConfig(siteName, siteSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_power_panel.test", "name", "Test Power Panel"),
					resource.TestCheckResourceAttrSet("data.netbox_power_panel.test", "site"),
				),
			},
		},
	})
}

func testAccPowerPanelDataSourceConfig(siteName, siteSlug string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_power_panel" "test" {
  site = netbox_site.test.id
  name = "Test Power Panel"
}

data "netbox_power_panel" "test" {
  id = netbox_power_panel.test.id
}
`, siteName, siteSlug)
}
