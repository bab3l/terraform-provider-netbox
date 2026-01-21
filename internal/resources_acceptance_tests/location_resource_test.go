package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccLocationResource_basic(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-loc-site")
	siteSlug := testutil.RandomSlug("tf-test-loc-site")
	name := testutil.RandomName("tf-test-location")
	slug := testutil.RandomSlug("tf-test-location")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterLocationCleanup(slug)
	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.ComposeCheckDestroy(testutil.CheckLocationDestroy, testutil.CheckSiteDestroy),
		Steps: []resource.TestStep{
			{
				Config: testAccLocationResourceConfig_basic(siteName, siteSlug, name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_location.test", "id"),
					resource.TestCheckResourceAttr("netbox_location.test", "name", name),
					resource.TestCheckResourceAttr("netbox_location.test", "slug", slug),
					resource.TestCheckResourceAttrPair("netbox_location.test", "site", "netbox_site.test", "id"),
				),
			},
		},
	})
}

func TestAccLocationResource_full(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-loc-site-full")
	siteSlug := testutil.RandomSlug("tf-test-loc-s-full")
	name := testutil.RandomName("tf-test-location-full")
	slug := testutil.RandomSlug("tf-test-loc-full")
	description := testutil.RandomName("description")
	facility := "Building A"

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterLocationCleanup(slug)
	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.ComposeCheckDestroy(testutil.CheckLocationDestroy, testutil.CheckSiteDestroy),
		Steps: []resource.TestStep{
			{
				Config: testAccLocationResourceConfig_full(siteName, siteSlug, name, slug, description, facility),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_location.test", "id"),
					resource.TestCheckResourceAttr("netbox_location.test", "name", name),
					resource.TestCheckResourceAttr("netbox_location.test", "slug", slug),
					resource.TestCheckResourceAttr("netbox_location.test", "status", "active"),
					resource.TestCheckResourceAttr("netbox_location.test", "description", description),
					resource.TestCheckResourceAttr("netbox_location.test", "facility", facility),
				),
			},
		},
	})
}

func TestAccLocationResource_import(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-loc-site")
	siteSlug := testutil.RandomSlug("tf-test-loc-site")
	name := testutil.RandomName("tf-test-location")
	slug := testutil.RandomSlug("tf-test-location")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterLocationCleanup(slug)
	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.ComposeCheckDestroy(testutil.CheckLocationDestroy, testutil.CheckSiteDestroy),
		Steps: []resource.TestStep{
			{
				Config: testAccLocationResourceConfig_import(siteName, siteSlug, name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_location.test", "id"),
					resource.TestCheckResourceAttr("netbox_location.test", "name", name),
					resource.TestCheckResourceAttr("netbox_location.test", "slug", slug),
					resource.TestCheckResourceAttrPair("netbox_location.test", "site", "netbox_site.test", "id"),
				),
			},
			{
				ResourceName:            "netbox_location.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"site"},
				Check: resource.ComposeTestCheckFunc(
					testutil.ReferenceFieldCheck("netbox_location.test", "site"),
					testutil.ReferenceFieldCheck("netbox_location.test", "parent"),
					testutil.ReferenceFieldCheck("netbox_location.test", "tenant"),
				),
			},
			{
				Config:   testAccLocationResourceConfig_import(siteName, siteSlug, name, slug),
				PlanOnly: true,
			},
		},
	})
}

// NOTE: Custom field tests for location resource are in resources_acceptance_tests_customfields package

func TestAccLocationResource_update(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-loc-site")
	siteSlug := testutil.RandomSlug("tf-test-loc-site")
	name := testutil.RandomName("tf-test-location")
	slug := testutil.RandomSlug("tf-test-location")
	updatedName := testutil.RandomName("tf-test-location-updated")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterLocationCleanup(slug)
	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.ComposeCheckDestroy(testutil.CheckLocationDestroy, testutil.CheckSiteDestroy),
		Steps: []resource.TestStep{
			{
				Config: testAccLocationResourceConfig_basic(siteName, siteSlug, name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_location.test", "id"),
					resource.TestCheckResourceAttr("netbox_location.test", "name", name),
				),
			},
			{
				Config: testAccLocationResourceConfig_basic(siteName, siteSlug, updatedName, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_location.test", "id"),
					resource.TestCheckResourceAttr("netbox_location.test", "name", updatedName),
				),
			},
		},
	})
}

func TestAccLocationResource_tagLifecycle(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-loc-site-tags")
	siteSlug := testutil.RandomSlug("tf-test-loc-site-tags")
	name := testutil.RandomName("tf-test-location-tags")
	slug := testutil.RandomSlug("tf-test-location-tags")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2Slug := testutil.RandomSlug("tag2")
	tag3Slug := testutil.RandomSlug("tag3")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterLocationCleanup(slug)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)
	cleanup.RegisterTagCleanup(tag3Slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.ComposeCheckDestroy(testutil.CheckLocationDestroy, testutil.CheckSiteDestroy),
		Steps: []resource.TestStep{
			{
				Config: testAccLocationResourceConfig_tags(siteName, siteSlug, name, slug, tag1Slug, tag2Slug, tag3Slug, caseTag1Tag2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_location.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemAttr("netbox_location.test", "tags.*", tag1Slug),
					resource.TestCheckTypeSetElemAttr("netbox_location.test", "tags.*", tag2Slug),
				),
			},
			{
				Config: testAccLocationResourceConfig_tags(siteName, siteSlug, name, slug, tag1Slug, tag2Slug, tag3Slug, caseTag1Uscore2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_location.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemAttr("netbox_location.test", "tags.*", tag1Slug),
					resource.TestCheckTypeSetElemAttr("netbox_location.test", "tags.*", tag2Slug),
				),
			},
			{
				Config: testAccLocationResourceConfig_tags(siteName, siteSlug, name, slug, tag1Slug, tag2Slug, tag3Slug, caseTag3),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_location.test", "tags.#", "1"),
					resource.TestCheckTypeSetElemAttr("netbox_location.test", "tags.*", tag3Slug),
				),
			},
			{
				Config: testAccLocationResourceConfig_tags(siteName, siteSlug, name, slug, tag1Slug, tag2Slug, tag3Slug, tagsEmpty),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_location.test", "tags.#", "0"),
				),
			},
		},
	})
}

func TestAccLocationResource_tagOrderInvariance(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-loc-site-tag-order")
	siteSlug := testutil.RandomSlug("tf-test-loc-site-tag-order")
	name := testutil.RandomName("tf-test-location-tag-order")
	slug := testutil.RandomSlug("tf-test-location-tag-order")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2Slug := testutil.RandomSlug("tag2")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterLocationCleanup(slug)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.ComposeCheckDestroy(testutil.CheckLocationDestroy, testutil.CheckSiteDestroy),
		Steps: []resource.TestStep{
			{
				Config: testAccLocationResourceConfig_tagsOrder(siteName, siteSlug, name, slug, tag1Slug, tag2Slug, caseTag1Tag2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_location.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemAttr("netbox_location.test", "tags.*", tag1Slug),
					resource.TestCheckTypeSetElemAttr("netbox_location.test", "tags.*", tag2Slug),
				),
			},
			{
				Config: testAccLocationResourceConfig_tagsOrder(siteName, siteSlug, name, slug, tag1Slug, tag2Slug, caseTag2Uscore1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_location.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemAttr("netbox_location.test", "tags.*", tag1Slug),
					resource.TestCheckTypeSetElemAttr("netbox_location.test", "tags.*", tag2Slug),
				),
			},
		},
	})
}

func TestAccConsistency_Location_LiteralNames(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-loc-site-lit")
	siteSlug := testutil.RandomSlug("tf-test-loc-site-lit")
	name := testutil.RandomName("tf-test-location-lit")
	slug := testutil.RandomSlug("tf-test-location-lit")
	description := testutil.RandomName("description")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterLocationCleanup(slug)
	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.ComposeCheckDestroy(testutil.CheckLocationDestroy, testutil.CheckSiteDestroy),
		Steps: []resource.TestStep{
			{
				Config: testAccLocationConsistencyLiteralNamesConfig(siteName, siteSlug, name, slug, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_location.test", "id"),
					resource.TestCheckResourceAttr("netbox_location.test", "name", name),
					resource.TestCheckResourceAttr("netbox_location.test", "slug", slug),
					resource.TestCheckResourceAttr("netbox_location.test", "description", description),
				),
			},
			{
				Config:   testAccLocationConsistencyLiteralNamesConfig(siteName, siteSlug, name, slug, description),
				PlanOnly: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_location.test", "id"),
				),
			},
		},
	})
}

func testAccLocationConsistencyLiteralNamesConfig(siteName, siteSlug, name, slug, description string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name   = %q
  slug   = %q
  status = "active"
}

resource "netbox_location" "test" {
  name        = %q
  slug        = %q
  site        = netbox_site.test.id
  description = %q
}
`, siteName, siteSlug, name, slug, description)
}

func testAccLocationResourceConfig_basic(siteName, siteSlug, name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name   = %q
  slug   = %q
  status = "active"
}

resource "netbox_location" "test" {
  name   = %q
  slug   = %q
  site   = netbox_site.test.id
  status = "active"
}
`, siteName, siteSlug, name, slug)
}

func testAccLocationResourceConfig_tags(siteName, siteSlug, name, slug, tag1Slug, tag2Slug, tag3Slug, tagCase string) string {
	var tagsConfig string
	switch tagCase {
	case caseTag1Tag2:
		tagsConfig = tagsDoubleSlug
	case caseTag1Uscore2:
		tagsConfig = tagsDoubleSlug
	case caseTag3:
		tagsConfig = tagsSingleSlug
	case tagsEmpty:
		tagsConfig = tagsEmpty
	}

	return fmt.Sprintf(`
resource "netbox_tag" "tag1" {
	name = "Tag1-%[5]s"
	slug = %[5]q
}

resource "netbox_tag" "tag2" {
	name = "Tag2-%[6]s"
	slug = %[6]q
}

resource "netbox_tag" "tag3" {
	name = "Tag3-%[7]s"
	slug = %[7]q
}

resource "netbox_site" "test" {
	name = %[1]q
	slug = %[2]q
}

resource "netbox_location" "test" {
	name = %[3]q
	slug = %[4]q
	site = netbox_site.test.id
	%[8]s
}
`, siteName, siteSlug, name, slug, tag1Slug, tag2Slug, tag3Slug, tagsConfig)
}

func testAccLocationResourceConfig_tagsOrder(siteName, siteSlug, name, slug, tag1Slug, tag2Slug, tagCase string) string {
	var tagsConfig string
	switch tagCase {
	case caseTag1Tag2:
		tagsConfig = tagsDoubleSlug
	case caseTag2Uscore1:
		tagsConfig = tagsDoubleSlugReversed
	}

	return fmt.Sprintf(`
resource "netbox_tag" "tag1" {
	name = "Tag1-%[5]s"
	slug = %[5]q
}

resource "netbox_tag" "tag2" {
	name = "Tag2-%[6]s"
	slug = %[6]q
}

resource "netbox_site" "test" {
	name = %[1]q
	slug = %[2]q
}

resource "netbox_location" "test" {
	name = %[3]q
	slug = %[4]q
	site = netbox_site.test.id
	%[7]s
}
`, siteName, siteSlug, name, slug, tag1Slug, tag2Slug, tagsConfig)
}

func testAccLocationResourceConfig_full(siteName, siteSlug, name, slug, description, facility string) string {
	return fmt.Sprintf(`
terraform {
  required_providers {
    netbox = {
      source  = "bab3l/netbox"
      version = ">= 0.1.0"
    }
  }
}

provider "netbox" {}

resource "netbox_site" "test" {
  name   = %q
  slug   = %q
  status = "active"
}

resource "netbox_location" "test" {
  name        = %q
  slug        = %q
  site        = netbox_site.test.id
  status      = "active"
  description = %q
  facility    = %q
}
`, siteName, siteSlug, name, slug, description, facility)

}

func testAccLocationResourceConfig_import(siteName, siteSlug, name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name = %q
  slug = %q
}

resource "netbox_location" "test" {
  name = %q
  slug = %q
  site = netbox_site.test.id
}
`, siteName, siteSlug, name, slug)
}

func TestAccLocationResource_externalDeletion(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-loc-site-del")
	siteSlug := testutil.RandomSlug("tf-test-loc-site-del")
	name := testutil.RandomName("tf-test-location-del")
	slug := testutil.RandomSlug("tf-test-location-del")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterLocationCleanup(slug)
	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.ComposeCheckDestroy(testutil.CheckLocationDestroy, testutil.CheckSiteDestroy),
		Steps: []resource.TestStep{
			{
				Config: testAccLocationResourceConfig_basic(siteName, siteSlug, name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_location.test", "id"),
					resource.TestCheckResourceAttr("netbox_location.test", "name", name),
					resource.TestCheckResourceAttr("netbox_location.test", "slug", slug),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					items, _, err := client.DcimAPI.DcimLocationsList(context.Background()).Slug([]string{slug}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find location for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.DcimAPI.DcimLocationsDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete location: %v", err)
					}
					t.Logf("Successfully externally deleted location with ID: %d", itemID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

// TestAccLocationResource_removeOptionalFields tests that removing previously set parent and tenant fields correctly sets them to null.
// This addresses the bug where removing a nullable reference field from the configuration would not clear it in NetBox,
// causing "Provider produced inconsistent result after apply" errors.
func TestAccLocationResource_removeOptionalFields(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-site-loc-remove")
	siteSlug := testutil.RandomSlug("tf-test-site-loc-remove")
	name := testutil.RandomName("tf-test-location-remove")
	slug := testutil.RandomSlug("tf-test-location-remove")
	parentName := testutil.RandomName("tf-test-parent-location")
	parentSlug := testutil.RandomSlug("tf-test-parent-location")
	tenantName := testutil.RandomName("test-tenant-location")
	tenantSlug := testutil.GenerateSlug(tenantName)

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterLocationCleanup(slug)
	cleanup.RegisterLocationCleanup(parentSlug)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterTenantCleanup(tenantSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Step 1: Create with parent and tenant
			{
				Config: testAccLocationResourceConfig_withParentAndTenant(siteName, siteSlug, parentName, parentSlug, name, slug, tenantName, tenantSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_location.test", "parent"),
					resource.TestCheckResourceAttrSet("netbox_location.test", "tenant"),
				),
			},
			// Step 2: Remove parent and tenant - should set to null
			{
				Config: testAccLocationResourceConfig_withoutFields(siteName, siteSlug, parentName, parentSlug, name, slug, tenantName, tenantSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("netbox_location.test", "parent"),
					resource.TestCheckNoResourceAttr("netbox_location.test", "tenant"),
				),
			},
			// Step 3: Re-add parent and tenant - should work without errors
			{
				Config: testAccLocationResourceConfig_withParentAndTenant(siteName, siteSlug, parentName, parentSlug, name, slug, tenantName, tenantSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_location.test", "parent"),
					resource.TestCheckResourceAttrSet("netbox_location.test", "tenant"),
				),
			},
		},
	})
}

func testAccLocationResourceConfig_withParentAndTenant(siteName, siteSlug, parentName, parentSlug, name, slug, tenantName, tenantSlug string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name   = %[1]q
  slug   = %[2]q
  status = "active"
}

resource "netbox_tenant" "test" {
  name = %[7]q
  slug = %[8]q
}

resource "netbox_location" "parent" {
  name   = %[3]q
  slug   = %[4]q
  site   = netbox_site.test.id
  status = "active"
}

resource "netbox_location" "test" {
  name   = %[5]q
  slug   = %[6]q
  site   = netbox_site.test.id
  parent = netbox_location.parent.id
  tenant = netbox_tenant.test.id
  status = "active"
}
`, siteName, siteSlug, parentName, parentSlug, name, slug, tenantName, tenantSlug)
}

func testAccLocationResourceConfig_withoutFields(siteName, siteSlug, parentName, parentSlug, name, slug, tenantName, tenantSlug string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name   = %[1]q
  slug   = %[2]q
  status = "active"
}

resource "netbox_tenant" "test" {
  name = %[7]q
  slug = %[8]q
}

resource "netbox_location" "parent" {
  name   = %[3]q
  slug   = %[4]q
  site   = netbox_site.test.id
  status = "active"
}

resource "netbox_location" "test" {
  name   = %[5]q
  slug   = %[6]q
  site   = netbox_site.test.id
  status = "active"
  # parent and tenant removed - should set to null
}
`, siteName, siteSlug, parentName, parentSlug, name, slug, tenantName, tenantSlug)
}

func TestAccLocationResource_removeDescription(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-site-loc-desc")
	siteSlug := testutil.RandomSlug("tf-test-site-loc-desc")
	name := testutil.RandomName("tf-test-location-desc")
	slug := testutil.RandomSlug("tf-test-location-desc")
	description := testutil.RandomName("description")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterLocationCleanup(slug)
	cleanup.RegisterSiteCleanup(siteSlug)

	testutil.TestRemoveOptionalFields(t, testutil.MultiFieldOptionalTestConfig{
		ResourceName: "netbox_location",
		BaseConfig: func() string {
			return testAccLocationResourceConfig_basic(siteName, siteSlug, name, slug)
		},
		ConfigWithFields: func() string {
			return testAccLocationResourceConfig_withDescription(siteName, siteSlug, name, slug, description)
		},
		OptionalFields: map[string]string{
			"description": description,
		},
		RequiredFields: map[string]string{
			"name": name,
			"slug": slug,
		},
		CheckDestroy: testutil.CheckLocationDestroy,
	})
}

func testAccLocationResourceConfig_withDescription(siteName, siteSlug, name, slug, description string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name   = %[1]q
  slug   = %[2]q
  status = "active"
}

resource "netbox_location" "test" {
  name        = %[3]q
  slug        = %[4]q
  site        = netbox_site.test.id
  status      = "active"
  description = %[5]q
}
`, siteName, siteSlug, name, slug, description)
}

func TestAccLocationResource_validationErrors(t *testing.T) {
	testutil.RunMultiValidationErrorTest(t, testutil.MultiValidationErrorTestConfig{
		ResourceName: "netbox_location",
		TestCases: map[string]testutil.ValidationErrorCase{
			"missing_name": {
				Config: func() string {
					return `
resource "netbox_site" "test" {
  name = "Test Site"
  slug = "test-site"
  status = "active"
}

resource "netbox_location" "test" {
  slug = "test-location"
  site = netbox_site.test.id
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
			"missing_slug": {
				Config: func() string {
					return `
resource "netbox_site" "test" {
  name = "Test Site"
  slug = "test-site"
  status = "active"
}

resource "netbox_location" "test" {
  name = "Test Location"
  site = netbox_site.test.id
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
			"missing_site": {
				Config: func() string {
					return `
resource "netbox_location" "test" {
  name = "Test Location"
  slug = "test-location"
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
			"invalid_site_reference": {
				Config: func() string {
					return `
resource "netbox_location" "test" {
  name = "Test Location"
  slug = "test-location"
  site = "99999"
}
`
				},
				ExpectedError: testutil.ErrPatternNotFound,
			},
		},
	})
}
