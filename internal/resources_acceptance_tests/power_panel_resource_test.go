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

func TestAccPowerPanelResource_tagLifecycle(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-site-tags")
	siteSlug := testutil.RandomSlug("tf-test-site-tags")
	panelName := testutil.RandomName("tf-test-panel-tags")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2Slug := testutil.RandomSlug("tag2")
	tag3Slug := testutil.RandomSlug("tag3")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)
	cleanup.RegisterTagCleanup(tag3Slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPowerPanelResourceConfig_tags(siteName, siteSlug, panelName, tag1Slug, tag2Slug, tag3Slug, caseTag1Tag2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_power_panel.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_power_panel.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag1-%s", tag1Slug),
						"slug": tag1Slug,
					}),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_power_panel.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag2-%s", tag2Slug),
						"slug": tag2Slug,
					}),
				),
			},
			{
				Config: testAccPowerPanelResourceConfig_tags(siteName, siteSlug, panelName, tag1Slug, tag2Slug, tag3Slug, caseTag1Uscore2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_power_panel.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_power_panel.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag1-%s", tag1Slug),
						"slug": tag1Slug,
					}),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_power_panel.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag2-%s", tag2Slug),
						"slug": tag2Slug,
					}),
				),
			},
			{
				Config: testAccPowerPanelResourceConfig_tags(siteName, siteSlug, panelName, tag1Slug, tag2Slug, tag3Slug, caseTag3),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_power_panel.test", "tags.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_power_panel.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag3-%s", tag3Slug),
						"slug": tag3Slug,
					}),
				),
			},
			{
				Config: testAccPowerPanelResourceConfig_tags(siteName, siteSlug, panelName, tag1Slug, tag2Slug, tag3Slug, tagsEmpty),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_power_panel.test", "tags.#", "0"),
				),
			},
		},
	})
}

func TestAccPowerPanelResource_tagOrderInvariance(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-site-tag-order")
	siteSlug := testutil.RandomSlug("tf-test-site-tag-order")
	panelName := testutil.RandomName("tf-test-panel-tag-order")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2Slug := testutil.RandomSlug("tag2")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPowerPanelResourceConfig_tagsOrder(siteName, siteSlug, panelName, tag1Slug, tag2Slug, caseTag1Tag2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_power_panel.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_power_panel.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag1-%s", tag1Slug),
						"slug": tag1Slug,
					}),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_power_panel.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag2-%s", tag2Slug),
						"slug": tag2Slug,
					}),
				),
			},
			{
				Config: testAccPowerPanelResourceConfig_tagsOrder(siteName, siteSlug, panelName, tag1Slug, tag2Slug, caseTag2Uscore1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_power_panel.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_power_panel.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag1-%s", tag1Slug),
						"slug": tag1Slug,
					}),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_power_panel.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag2-%s", tag2Slug),
						"slug": tag2Slug,
					}),
				),
			},
		},
	})
}

func TestAccPowerPanelResource_update(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-site-update")
	siteSlug := testutil.RandomSlug("tf-test-site-update")
	panelName := testutil.RandomName("tf-test-panel-update")
	updatedPanelName := testutil.RandomName("tf-test-panel-updated")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPowerPanelResourceConfig_forUpdate(siteName, siteSlug, panelName, testutil.Description1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_power_panel.test", "id"),
					resource.TestCheckResourceAttr("netbox_power_panel.test", "name", panelName),
					resource.TestCheckResourceAttr("netbox_power_panel.test", "description", testutil.Description1),
				),
			},
			{
				Config: testAccPowerPanelResourceConfig_forUpdate(siteName, siteSlug, updatedPanelName, testutil.Description2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_power_panel.test", "name", updatedPanelName),
					resource.TestCheckResourceAttr("netbox_power_panel.test", "description", testutil.Description2),
				),
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

func testAccPowerPanelResourceConfig_forUpdate(siteName, siteSlug, panelName, description string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name = %q
  slug = %q
}

resource "netbox_power_panel" "test" {
  site        = netbox_site.test.id
  name        = %q
  description = %q
}
`, siteName, siteSlug, panelName, description)
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

func testAccPowerPanelResourceConfig_tags(siteName, siteSlug, panelName, tag1Slug, tag2Slug, tag3Slug, tagCase string) string {
	var tagsConfig string
	switch tagCase {
	case caseTag1Tag2:
		tagsConfig = tagsDoubleNested
	case caseTag1Uscore2:
		tagsConfig = tagsDoubleNested
	case caseTag3:
		tagsConfig = tagsSingleNested
	case tagsEmpty:
		tagsConfig = tagsEmpty
	}

	return fmt.Sprintf(`
resource "netbox_tag" "tag1" {
	name = "Tag1-%[4]s"
	slug = %[4]q
}

resource "netbox_tag" "tag2" {
	name = "Tag2-%[5]s"
	slug = %[5]q
}

resource "netbox_tag" "tag3" {
	name = "Tag3-%[6]s"
	slug = %[6]q
}

resource "netbox_site" "test" {
	name   = %[1]q
	slug   = %[2]q
	status = "active"
}

resource "netbox_power_panel" "test" {
	site = netbox_site.test.id
	name = %[3]q
	%[7]s
}
`, siteName, siteSlug, panelName, tag1Slug, tag2Slug, tag3Slug, tagsConfig)
}

func testAccPowerPanelResourceConfig_tagsOrder(siteName, siteSlug, panelName, tag1Slug, tag2Slug, tagCase string) string {
	var tagsConfig string
	switch tagCase {
	case caseTag1Tag2:
		tagsConfig = tagsDoubleNested
	case caseTag2Uscore1:
		tagsConfig = tagsDoubleNestedReversed
	}

	return fmt.Sprintf(`
resource "netbox_tag" "tag1" {
	name = "Tag1-%[4]s"
	slug = %[4]q
}

resource "netbox_tag" "tag2" {
	name = "Tag2-%[5]s"
	slug = %[5]q
}

resource "netbox_site" "test" {
	name   = %[1]q
	slug   = %[2]q
	status = "active"
}

resource "netbox_power_panel" "test" {
	site = netbox_site.test.id
	name = %[3]q
	%[6]s
}
`, siteName, siteSlug, panelName, tag1Slug, tag2Slug, tagsConfig)
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

func TestAccPowerPanelResource_validationErrors(t *testing.T) {
	testutil.RunMultiValidationErrorTest(t, testutil.MultiValidationErrorTestConfig{
		ResourceName: "netbox_power_panel",
		TestCases: map[string]testutil.ValidationErrorCase{
			"missing_site": {
				Config: func() string {
					return `
provider "netbox" {}

resource "netbox_power_panel" "test" {
  # site missing
  name = "Test Panel"
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
			"missing_name": {
				Config: func() string {
					return `
provider "netbox" {}

resource "netbox_site" "test" {
  name = "Test Site"
  slug = "test-site"
}

resource "netbox_power_panel" "test" {
  site = netbox_site.test.id
  # name missing
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
		},
	})
}
