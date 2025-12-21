package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPowerFeedDataSource_basic(t *testing.T) {

	t.Parallel()

	cleanup := testutil.NewCleanupResource(t)

	siteName := testutil.RandomName("test-power-feed-site-ds")
	siteSlug := testutil.GenerateSlug(siteName)

	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterPowerFeedCleanup("Test Power Feed")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckSiteDestroy,
			testutil.CheckPowerFeedDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccPowerFeedDataSourceConfig(siteName, siteSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_power_feed.test", "name", "Test Power Feed"),
					resource.TestCheckResourceAttr("data.netbox_power_feed.test", "status", "active"),
				),
			},
		},
	})
}

func testAccPowerFeedDataSourceConfig(siteName, siteSlug string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name = %[1]q
  slug = %[2]q
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
`, siteName, siteSlug)
}
