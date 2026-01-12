package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// NOTE: Custom field tests for power_feed resource are in resources_acceptance_tests_customfields package

func TestAccPowerFeedResource_basic(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("site")
	siteSlug := testutil.RandomSlug("site")
	panelName := testutil.RandomName("power-panel")
	feedName := testutil.RandomName("power-feed")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
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
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
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

// TestAccPowerFeedResource_VoltageAmperage tests comprehensive scenarios for power feed voltage field.
// This validates that Optional+Computed numeric fields work correctly across all scenarios.
func TestAccPowerFeedResource_VoltageAmperage(t *testing.T) {
	siteName := testutil.RandomName("tf-test-site-power-feed")
	siteSlug := testutil.RandomSlug("tf-test-site-power-feed")
	powerPanelName := testutil.RandomName("tf-test-panel-power-feed")
	powerFeedName := testutil.RandomName("tf-test-power-feed-voltage")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)

	testutil.RunOptionalComputedFieldTestSuite(t, testutil.OptionalComputedFieldTestConfig{
		ResourceName:   "netbox_power_feed",
		OptionalField:  "voltage",
		DefaultValue:   "120",
		FieldTestValue: "240",
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckPowerFeedDestroy,
			testutil.CheckPowerPanelDestroy,
			testutil.CheckSiteDestroy,
		),
		BaseConfig: func() string {
			return testAccPowerFeedResourceWithOptionalField(siteName, siteSlug, powerPanelName, powerFeedName, "voltage", "")
		},
		WithFieldConfig: func(value string) string {
			return testAccPowerFeedResourceWithOptionalField(siteName, siteSlug, powerPanelName, powerFeedName, "voltage", value)
		},
	})
}

// TestAccPowerFeedResource_Amperage tests the amperage field separately.
func TestAccPowerFeedResource_Amperage(t *testing.T) {
	siteName := testutil.RandomName("tf-test-site-power-feed-amp")
	siteSlug := testutil.RandomSlug("tf-test-site-power-feed-amp")
	powerPanelName := testutil.RandomName("tf-test-panel-power-feed-amp")
	powerFeedName := testutil.RandomName("tf-test-power-feed-amperage")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)

	testutil.RunOptionalComputedFieldTestSuite(t, testutil.OptionalComputedFieldTestConfig{
		ResourceName:   "netbox_power_feed",
		OptionalField:  "amperage",
		DefaultValue:   "20",
		FieldTestValue: "30",
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckPowerFeedDestroy,
			testutil.CheckPowerPanelDestroy,
			testutil.CheckSiteDestroy,
		),
		BaseConfig: func() string {
			return testAccPowerFeedResourceWithOptionalField(siteName, siteSlug, powerPanelName, powerFeedName, "amperage", "")
		},
		WithFieldConfig: func(value string) string {
			return testAccPowerFeedResourceWithOptionalField(siteName, siteSlug, powerPanelName, powerFeedName, "amperage", value)
		},
	})
}

func testAccPowerFeedResourceWithOptionalField(siteName, siteSlug, powerPanelName, powerFeedName, optionalFieldName, optionalFieldValue string) string {
	optionalField := ""
	if optionalFieldValue != "" {
		optionalField = fmt.Sprintf("\n  %s = %s", optionalFieldName, optionalFieldValue)
	}

	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name = %q
  slug = %q
}

resource "netbox_power_panel" "test" {
  name = %q
  site = netbox_site.test.id
}

resource "netbox_power_feed" "test" {
  name        = %q
  power_panel = netbox_power_panel.test.id%s
}
`, siteName, siteSlug, powerPanelName, powerFeedName, optionalField)
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

func TestAccPowerFeedResource_externalDeletion(t *testing.T) {
	t.Parallel()
	testutil.TestAccPreCheck(t)

	siteName := testutil.RandomName("tf-test-site-extdel")
	siteSlug := testutil.RandomSlug("tf-test-site-ed")
	panelName := testutil.RandomName("tf-test-panel-extdel")
	feedName := testutil.RandomName("tf-test-feed-extdel")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPowerFeedResourceConfig_basic(siteName, siteSlug, panelName, feedName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_power_feed.test", "id"),
					resource.TestCheckResourceAttr("netbox_power_feed.test", "name", feedName),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}

					feeds, _, err := client.DcimAPI.DcimPowerFeedsList(context.Background()).Name([]string{feedName}).Execute()
					if err != nil || feeds == nil || len(feeds.Results) == 0 {
						t.Fatalf("Failed to find power feed for external deletion: %v", err)
					}
					feedID := feeds.Results[0].Id
					_, err = client.DcimAPI.DcimPowerFeedsDestroy(context.Background(), feedID).Execute()
					if err != nil {
						t.Fatalf("Failed to delete power feed: %v", err)
					}
					t.Logf("Successfully externally deleted power feed with ID: %d", feedID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccPowerFeedResource_removeDescriptionAndComments(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-site-pf-optional")
	siteSlug := testutil.RandomSlug("tf-test-site-pf")
	panelName := testutil.RandomName("tf-test-panel-optional")
	feedName := testutil.RandomName("tf-test-feed-optional")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterPowerPanelCleanup(panelName)
	cleanup.RegisterPowerFeedCleanup(feedName)

	testutil.TestRemoveOptionalFields(t, testutil.MultiFieldOptionalTestConfig{
		ResourceName: "netbox_power_feed",
		BaseConfig: func() string {
			return testAccPowerFeedResourceConfig_basic(siteName, siteSlug, panelName, feedName)
		},
		ConfigWithFields: func() string {
			return testAccPowerFeedResourceConfig_withDescriptionAndComments(
				siteName,
				siteSlug,
				panelName,
				feedName,
				"Test description",
				"Test comments",
			)
		},
		OptionalFields: map[string]string{
			"description": "Test description",
			"comments":    "Test comments",
		},
		RequiredFields: map[string]string{
			"name": feedName,
		},
		CheckDestroy: testutil.CheckPowerFeedDestroy,
	})
}

func testAccPowerFeedResourceConfig_withDescriptionAndComments(siteName, siteSlug, panelName, feedName, description, comments string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name   = %[1]q
  slug   = %[2]q
  status = "active"
}

resource "netbox_power_panel" "test" {
  name = %[3]q
  site = netbox_site.test.id
}

resource "netbox_power_feed" "test" {
  name         = %[4]q
  power_panel  = netbox_power_panel.test.id
  status       = "active"
  description  = %[5]q
  comments     = %[6]q
}
`, siteName, siteSlug, panelName, feedName, description, comments)
}
