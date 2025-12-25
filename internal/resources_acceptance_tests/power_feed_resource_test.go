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

func TestAccPowerFeedResource_basic(t *testing.T) {

	t.Parallel()
	siteName := testutil.RandomName("site")
	siteSlug := testutil.RandomSlug("site")
	panelName := testutil.RandomName("power-panel")
	feedName := testutil.RandomName("power-feed")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccPowerFeedResourceConfig_basic(siteName, siteSlug, panelName, feedName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_power_feed.test", "id"),
					resource.TestCheckResourceAttr("netbox_power_feed.test", "name", feedName),
					resource.TestCheckResourceAttr("netbox_power_feed.test", "status", "active"),
				),
			},
			{
				ResourceName:            "netbox_power_feed.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"power_panel"},
			},
		},
	})
}

func TestAccPowerFeedResource_full(t *testing.T) {

	t.Parallel()
	siteName := testutil.RandomName("site")
	siteSlug := testutil.RandomSlug("site")
	panelName := testutil.RandomName("power-panel")
	feedName := testutil.RandomName("power-feed")
	description := "Test power feed with all fields"
	updatedDescription := "Updated power feed description"

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccPowerFeedResourceConfig_full(siteName, siteSlug, panelName, feedName, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_power_feed.test", "id"),
					resource.TestCheckResourceAttr("netbox_power_feed.test", "name", feedName),
					resource.TestCheckResourceAttr("netbox_power_feed.test", "status", "active"),
					resource.TestCheckResourceAttr("netbox_power_feed.test", "type", "primary"),
					resource.TestCheckResourceAttr("netbox_power_feed.test", "supply", "ac"),
					resource.TestCheckResourceAttr("netbox_power_feed.test", "phase", "single-phase"),
					resource.TestCheckResourceAttr("netbox_power_feed.test", "voltage", "240"),
					resource.TestCheckResourceAttr("netbox_power_feed.test", "amperage", "30"),
					resource.TestCheckResourceAttr("netbox_power_feed.test", "description", description),
				),
			},
			{
				Config: testAccPowerFeedResourceConfig_full(siteName, siteSlug, panelName, feedName, updatedDescription),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_power_feed.test", "description", updatedDescription),
				),
			},
		},
	})
}

func TestAccPowerFeedResource_IDPreservation(t *testing.T) {
	t.Parallel()
	siteName := testutil.RandomName("site-id")
	siteSlug := testutil.RandomSlug("site-id")
	panelName := testutil.RandomName("power-panel-id")
	feedName := testutil.RandomName("power-feed-id")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccPowerFeedResourceConfig_basic(siteName, siteSlug, panelName, feedName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_power_feed.test", "id"),
					resource.TestCheckResourceAttr("netbox_power_feed.test", "name", feedName),
				),
			},
		},
	})
}

func testAccPowerFeedResourceConfig_basic(siteName, siteSlug, panelName, feedName string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name   = %[1]q
  slug   = %[2]q
  status = "active"
}

resource "netbox_power_panel" "test" {
  site = netbox_site.test.id
  name = %[3]q
}

resource "netbox_power_feed" "test" {
  power_panel = netbox_power_panel.test.id
  name        = %[4]q
}
`, siteName, siteSlug, panelName, feedName)
}

func testAccPowerFeedResourceConfig_full(siteName, siteSlug, panelName, feedName, description string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name   = %[1]q
  slug   = %[2]q
  status = "active"
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
  voltage     = 240
  amperage    = 30
  description = %[5]q
}
`, siteName, siteSlug, panelName, feedName, description)
}

func TestAccConsistency_PowerFeed(t *testing.T) {

	t.Parallel()
	siteName := testutil.RandomName("site")
	siteSlug := testutil.RandomSlug("site")
	rackName := testutil.RandomName("rack")
	locationName := testutil.RandomName("location")
	locationSlug := testutil.RandomSlug("location")
	powerPanelName := testutil.RandomName("power-panel")
	feedName := testutil.RandomName("power-feed")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPowerFeedConsistencyConfig(siteName, siteSlug, rackName, locationName, locationSlug, powerPanelName, feedName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_power_feed.test", "name", feedName),
					resource.TestCheckResourceAttr("netbox_power_feed.test", "rack", rackName),
				),
			},
			{
				PlanOnly: true,
				Config:   testAccPowerFeedConsistencyConfig(siteName, siteSlug, rackName, locationName, locationSlug, powerPanelName, feedName),
			},
		},
	})
}

func testAccPowerFeedConsistencyConfig(siteName, siteSlug, rackName, locationName, locationSlug, powerPanelName, feedName string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name = "%[1]s"
  slug = "%[2]s"
}

resource "netbox_location" "test" {
  name = "%[4]s"
  slug = "%[5]s"
  site = netbox_site.test.id
}

resource "netbox_rack" "test" {
  name     = "%[3]s"
  site     = netbox_site.test.id
  location = netbox_location.test.id
}

resource "netbox_power_panel" "test" {
  name = "%[6]s"
  site = netbox_site.test.id
}

resource "netbox_power_feed" "test" {
  name        = "%[7]s"
  power_panel = netbox_power_panel.test.id
  rack        = netbox_rack.test.name
}
`, siteName, siteSlug, rackName, locationName, locationSlug, powerPanelName, feedName)
}

func TestAccConsistency_PowerFeed_LiteralNames(t *testing.T) {

	t.Parallel()
	siteName := testutil.RandomName("site")
	siteSlug := testutil.RandomSlug("site")
	rackName := testutil.RandomName("rack")
	locationName := testutil.RandomName("location")
	locationSlug := testutil.RandomSlug("location")
	powerPanelName := testutil.RandomName("power-panel")
	feedName := testutil.RandomName("power-feed")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPowerFeedConsistencyLiteralNamesConfig(siteName, siteSlug, rackName, locationName, locationSlug, powerPanelName, feedName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_power_feed.test", "name", feedName),
					resource.TestCheckResourceAttr("netbox_power_feed.test", "rack", rackName),
				),
			},
			{
				PlanOnly: true,
				Config:   testAccPowerFeedConsistencyLiteralNamesConfig(siteName, siteSlug, rackName, locationName, locationSlug, powerPanelName, feedName),
			},
		},
	})
}

func testAccPowerFeedConsistencyLiteralNamesConfig(siteName, siteSlug, rackName, locationName, locationSlug, powerPanelName, feedName string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name = "%[1]s"
  slug = "%[2]s"
}

resource "netbox_location" "test" {
  name = "%[4]s"
  slug = "%[5]s"
  site = netbox_site.test.id
}

resource "netbox_rack" "test" {
  name     = "%[3]s"
  site     = netbox_site.test.id
  location = netbox_location.test.id
}

resource "netbox_power_panel" "test" {
  name = "%[6]s"
  site = netbox_site.test.id
}

resource "netbox_power_feed" "test" {
  name        = "%[7]s"
  power_panel = netbox_power_panel.test.id
  rack        = "%[3]s"

  depends_on = [netbox_rack.test]
}
`, siteName, siteSlug, rackName, locationName, locationSlug, powerPanelName, feedName)
}
