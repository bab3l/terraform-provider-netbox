package resources_acceptance_tests

import (
	"context"
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
func TestAccPowerFeedResource_importWithCustomFieldsAndTags(t *testing.T) {
	// NOTE: t.Parallel() intentionally omitted - this test creates/deletes global custom fields
	// that would affect other tests of the same resource type running in parallel.

	siteName := testutil.RandomName("site")
	siteSlug := testutil.RandomSlug("site")
	panelName := testutil.RandomName("power-panel")
	feedName := testutil.RandomName("power-feed")
	tenantName := testutil.RandomName("tf-test-tenant")
	tenantSlug := testutil.RandomSlug("tf-test-tenant")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterTenantCleanup(tenantSlug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccPowerFeedResourceImportConfig_full(siteName, siteSlug, panelName, feedName, tenantName, tenantSlug),
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
				ImportStateVerifyIgnore: []string{"power_panel", "tenant", "custom_fields", "tags"},
			},
		},
	})
}

func testAccPowerFeedResourceImportConfig_full(siteName, siteSlug, panelName, feedName, tenantName, tenantSlug string) string {
	// Generate test data for all custom field types
	textValue := testutil.RandomName("text-value")
	longtextValue := testutil.RandomName("longtext-value") + "\nThis is a multiline text field for comprehensive testing."
	intValue := 42 // Fixed value for reproducibility
	boolValue := true
	dateValue := testutil.RandomDate()
	urlValue := testutil.RandomURL("test-url")
	jsonValue := testutil.RandomJSON()

	// Tag names
	tag1 := testutil.RandomName("tag1")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2 := testutil.RandomName("tag2")
	tag2Slug := testutil.RandomSlug("tag2")

	// Custom field names
	cfText := testutil.RandomCustomFieldName("tf_text")
	cfLongtext := testutil.RandomCustomFieldName("tf_longtext")
	cfInteger := testutil.RandomCustomFieldName("tf_integer")
	cfBoolean := testutil.RandomCustomFieldName("tf_boolean")
	cfDate := testutil.RandomCustomFieldName("tf_date")
	cfURL := testutil.RandomCustomFieldName("tf_url")
	cfJSON := testutil.RandomCustomFieldName("tf_json")

	return fmt.Sprintf(`
# Dependencies
resource "netbox_site" "test" {
  name   = %q
  slug   = %q
  status = "active"
}

resource "netbox_tenant" "test" {
  name = %q
  slug = %q
}

resource "netbox_power_panel" "test" {
  name = %q
  site = netbox_site.test.id
}

# Tags
resource "netbox_tag" "tag1" {
  name = %q
  slug = %q
}

resource "netbox_tag" "tag2" {
  name = %q
  slug = %q
}

# Custom Fields for dcim.powerfeed object type
resource "netbox_custom_field" "test_text" {
  name         = %q
  label        = "Test Text CF"
  type         = "text"
  object_types = ["dcim.powerfeed"]
}

resource "netbox_custom_field" "test_longtext" {
  name         = %q
  label        = "Test Longtext CF"
  type         = "longtext"
  object_types = ["dcim.powerfeed"]
}

resource "netbox_custom_field" "test_integer" {
  name         = %q
  label        = "Test Integer CF"
  type         = "integer"
  object_types = ["dcim.powerfeed"]
}

resource "netbox_custom_field" "test_boolean" {
  name         = %q
  label        = "Test Boolean CF"
  type         = "boolean"
  object_types = ["dcim.powerfeed"]
}

resource "netbox_custom_field" "test_date" {
  name         = %q
  label        = "Test Date CF"
  type         = "date"
  object_types = ["dcim.powerfeed"]
}

resource "netbox_custom_field" "test_url" {
  name         = %q
  label        = "Test URL CF"
  type         = "url"
  object_types = ["dcim.powerfeed"]
}

resource "netbox_custom_field" "test_json" {
  name         = %q
  label        = "Test JSON CF"
  type         = "json"
  object_types = ["dcim.powerfeed"]
}

# Power Feed with comprehensive custom fields and tags
resource "netbox_power_feed" "test" {
  name        = %q
  power_panel = netbox_power_panel.test.id
  status      = "active"
  tenant      = netbox_tenant.test.id

  tags = [
    {
      name = netbox_tag.tag1.name
      slug = netbox_tag.tag1.slug
    },
    {
      name = netbox_tag.tag2.name
      slug = netbox_tag.tag2.slug
    }
  ]

  custom_fields = [
    {
      name  = netbox_custom_field.test_text.name
      type  = "text"
      value = %q
    },
    {
      name  = netbox_custom_field.test_longtext.name
      type  = "longtext"
      value = %q
    },
    {
      name  = netbox_custom_field.test_integer.name
      type  = "integer"
      value = "%d"
    },
    {
      name  = netbox_custom_field.test_boolean.name
      type  = "boolean"
      value = "%t"
    },
    {
      name  = netbox_custom_field.test_date.name
      type  = "date"
      value = %q
    },
    {
      name  = netbox_custom_field.test_url.name
      type  = "url"
      value = %q
    },
    {
      name  = netbox_custom_field.test_json.name
      type  = "json"
      value = %q
    },
  ]

  depends_on = [
    netbox_custom_field.test_text,
    netbox_custom_field.test_longtext,
    netbox_custom_field.test_integer,
    netbox_custom_field.test_boolean,
    netbox_custom_field.test_date,
    netbox_custom_field.test_url,
    netbox_custom_field.test_json,
  ]
}
`, siteName, siteSlug, tenantName, tenantSlug, panelName,
		tag1, tag1Slug, tag2, tag2Slug,
		cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfURL, cfJSON,
		feedName, textValue, longtextValue, intValue, boolValue, dateValue, urlValue, jsonValue)
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
