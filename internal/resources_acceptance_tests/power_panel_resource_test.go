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

func TestAccPowerPanelResource_basic(t *testing.T) {

	t.Parallel()
	siteName := testutil.RandomName("tf-test-site")
	siteSlug := testutil.RandomSlug("tf-test-site")
	panelName := testutil.RandomName("tf-test-panel")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
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
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
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
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
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
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
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
