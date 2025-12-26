package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPowerFeedDataSource_IDPreservation(t *testing.T) {
	t.Parallel()

	cleanup := testutil.NewCleanupResource(t)

	siteName := testutil.RandomName("tf-test-power-feed-site-id")
	siteSlug := testutil.GenerateSlug(siteName)
	powerPanelName := testutil.RandomName("tf-test-power-panel-id")
	powerFeedName := testutil.RandomName("tf-test-power-feed-id")

	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterPowerFeedCleanup(powerFeedName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckSiteDestroy,
			testutil.CheckPowerFeedDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccPowerFeedDataSourceConfig(siteName, siteSlug, powerPanelName, powerFeedName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_power_feed.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_power_feed.test", "name", powerFeedName),
				),
			},
		},
	})
}

func TestAccPowerFeedDataSource_basic(t *testing.T) {

	t.Parallel()

	cleanup := testutil.NewCleanupResource(t)

	siteName := testutil.RandomName("tf-test-power-feed-site")
	siteSlug := testutil.GenerateSlug(siteName)
	powerPanelName := testutil.RandomName("tf-test-power-panel")
	powerFeedName := testutil.RandomName("tf-test-power-feed")

	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterPowerFeedCleanup(powerFeedName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckSiteDestroy,
			testutil.CheckPowerFeedDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccPowerFeedDataSourceConfig(siteName, siteSlug, powerPanelName, powerFeedName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_power_feed.test", "name", powerFeedName),
					resource.TestCheckResourceAttr("data.netbox_power_feed.test", "status", "active"),
				),
			},
		},
	})
}

func TestAccPowerFeedDataSource_byPanelAndName(t *testing.T) {

	t.Parallel()

	cleanup := testutil.NewCleanupResource(t)

	siteName := testutil.RandomName("tf-test-power-feed-site")
	siteSlug := testutil.GenerateSlug(siteName)
	powerPanelName := testutil.RandomName("tf-test-power-panel")
	powerFeedName := testutil.RandomName("tf-test-power-feed")

	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterPowerFeedCleanup(powerFeedName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckSiteDestroy,
			testutil.CheckPowerFeedDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccPowerFeedDataSourceConfigByPanelAndName(siteName, siteSlug, powerPanelName, powerFeedName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_power_feed.by_panel", "name", powerFeedName),
					resource.TestCheckResourceAttrSet("data.netbox_power_feed.by_panel", "power_panel"),
				),
			},
		},
	})
}

func testAccPowerFeedDataSourceConfig(siteName, siteSlug, powerPanelName, powerFeedName string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_power_panel" "test" {
  site = netbox_site.test.id
  name = %[3]q
}

resource "netbox_power_feed" "test" {
  power_panel = netbox_power_panel.test.id
  name        = %[4]q
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
`, siteName, siteSlug, powerPanelName, powerFeedName)
}

func testAccPowerFeedDataSourceConfigByPanelAndName(siteName, siteSlug, powerPanelName, powerFeedName string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_power_panel" "test" {
  site = netbox_site.test.id
  name = %[3]q
}

resource "netbox_power_feed" "test" {
  power_panel = netbox_power_panel.test.id
  name        = %[4]q
  status      = "active"
  type        = "primary"
  supply      = "ac"
  phase       = "single-phase"
  voltage     = 120
  amperage    = 15
}

data "netbox_power_feed" "by_panel" {
  power_panel = netbox_power_panel.test.id
  name        = netbox_power_feed.test.name
}
`, siteName, siteSlug, powerPanelName, powerFeedName)
}
