package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPowerPanelResource_basic(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-site")
	siteSlug := testutil.RandomSlug("tf-test-site")
	panelName := testutil.RandomName("tf-test-panel")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPowerPanelResourceConfig_basic(siteName, siteSlug, panelName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_power_panel.test", "id"),
					resource.TestCheckResourceAttr("netbox_power_panel.test", "name", panelName),
				),
			},
			{
				ResourceName:            "netbox_power_panel.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"site"},
			},
		},
	})
}

func TestAccPowerPanelResource_full(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-site-full")
	siteSlug := testutil.RandomSlug("tf-test-site-full")
	panelName := testutil.RandomName("tf-test-panel-full")
	description := testutil.RandomName("description")
	updatedDescription := testutil.RandomName("description")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPowerPanelResourceConfig_full(siteName, siteSlug, panelName, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_power_panel.test", "id"),
					resource.TestCheckResourceAttr("netbox_power_panel.test", "name", panelName),
					resource.TestCheckResourceAttr("netbox_power_panel.test", "description", description),
				),
			},
			{
				Config: testAccPowerPanelResourceConfig_full(siteName, siteSlug, panelName, updatedDescription),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_power_panel.test", "description", updatedDescription),
				),
			},
		},
	})
}

func TestAccPowerPanelResource_IDPreservation(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-site-id")
	siteSlug := testutil.RandomSlug("tf-test-site-id")
	panelName := testutil.RandomName("tf-test-panel-id")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPowerPanelResourceConfig_basic(siteName, siteSlug, panelName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_power_panel.test", "id"),
					resource.TestCheckResourceAttr("netbox_power_panel.test", "name", panelName),
				),
			},
		},
	})
}

func testAccPowerPanelResourceConfig_basic(siteName, siteSlug, panelName string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name   = %q
  slug   = %q
  status = "active"
}

resource "netbox_power_panel" "test" {
  site = netbox_site.test.id
  name = %q
}
`, siteName, siteSlug, panelName)
}

func testAccPowerPanelResourceConfig_full(siteName, siteSlug, panelName, description string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name   = %[1]q
  slug   = %[2]q
  status = "active"
}

resource "netbox_power_panel" "test" {
  site        = netbox_site.test.id
  name        = %[3]q
  description = %[4]q
}
`, siteName, siteSlug, panelName, description)
}

func TestAccConsistency_PowerPanel_LiteralNames(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("site")
	siteSlug := testutil.RandomSlug("site")
	panelName := testutil.RandomName("power-panel")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPowerPanelConsistencyLiteralNamesConfig(siteName, siteSlug, panelName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_power_panel.test", "name", panelName),
					resource.TestCheckResourceAttr("netbox_power_panel.test", "site", siteName),
				),
			},
			{
				PlanOnly: true,
				Config:   testAccPowerPanelConsistencyLiteralNamesConfig(siteName, siteSlug, panelName),
			},
		},
	})
}

func testAccPowerPanelConsistencyLiteralNamesConfig(siteName, siteSlug, panelName string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name   = "%[1]s"
  slug   = "%[2]s"
  status = "active"
}

resource "netbox_power_panel" "test" {
  name = "%[3]s"
  site = "%[1]s"
  depends_on = [netbox_site.test]
}
`, siteName, siteSlug, panelName)
}

func TestAccPowerPanelResource_externalDeletion(t *testing.T) {
	t.Parallel()
	testutil.TestAccPreCheck(t)

	siteName := testutil.RandomName("tf-test-site-extdel")
	siteSlug := testutil.RandomSlug("tf-test-site-ed")
	panelName := testutil.RandomName("tf-test-panel-extdel")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPowerPanelResourceConfig_basic(siteName, siteSlug, panelName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_power_panel.test", "id"),
					resource.TestCheckResourceAttr("netbox_power_panel.test", "name", panelName),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}

					panels, _, err := client.DcimAPI.DcimPowerPanelsList(context.Background()).Name([]string{panelName}).Execute()
					if err != nil || panels == nil || len(panels.Results) == 0 {
						t.Fatalf("Failed to find power panel for external deletion: %v", err)
					}
					panelID := panels.Results[0].Id
					_, err = client.DcimAPI.DcimPowerPanelsDestroy(context.Background(), panelID).Execute()
					if err != nil {
						t.Fatalf("Failed to delete power panel: %v", err)
					}
					t.Logf("Successfully externally deleted power panel with ID: %d", panelID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

// TestAccPowerPanelResource_removeOptionalFields tests that optional fields
// can be successfully removed from the configuration without causing inconsistent state.
func TestAccPowerPanelResource_removeOptionalFields(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-site-rem")
	siteSlug := testutil.RandomSlug("tf-test-site-rem")
	locationName := testutil.RandomName("tf-test-loc-rem")
	locationSlug := testutil.RandomSlug("tf-test-loc-rem")
	panelName := testutil.RandomName("tf-test-panel-rem")
	description := testutil.RandomName("description")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPowerPanelResourceConfig_withLocation(siteName, siteSlug, locationName, locationSlug, panelName, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_power_panel.test", "name", panelName),
					resource.TestCheckResourceAttr("netbox_power_panel.test", "description", description),
					resource.TestCheckResourceAttrSet("netbox_power_panel.test", "location"),
				),
			},
			{
				Config: testAccPowerPanelResourceConfig_withLocationNoOptional(siteName, siteSlug, locationName, locationSlug, panelName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_power_panel.test", "name", panelName),
					resource.TestCheckNoResourceAttr("netbox_power_panel.test", "description"),
					resource.TestCheckNoResourceAttr("netbox_power_panel.test", "location"),
				),
			},
		},
	})
}

func testAccPowerPanelResourceConfig_withLocation(siteName, siteSlug, locationName, locationSlug, panelName, description string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name   = %[1]q
  slug   = %[2]q
  status = "active"
}

resource "netbox_location" "test" {
  name = %[3]q
  slug = %[4]q
  site = netbox_site.test.id
}

resource "netbox_power_panel" "test" {
  site        = netbox_site.test.id
  name        = %[5]q
  location    = netbox_location.test.id
  description = %[6]q
}
`, siteName, siteSlug, locationName, locationSlug, panelName, description)
}

func testAccPowerPanelResourceConfig_withLocationNoOptional(siteName, siteSlug, locationName, locationSlug, panelName string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name   = %[1]q
  slug   = %[2]q
  status = "active"
}

resource "netbox_location" "test" {
  name = %[3]q
  slug = %[4]q
  site = netbox_site.test.id
}

resource "netbox_power_panel" "test" {
  site = netbox_site.test.id
  name = %[5]q
}
`, siteName, siteSlug, locationName, locationSlug, panelName)
}
