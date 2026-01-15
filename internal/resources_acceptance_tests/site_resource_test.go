package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// NOTE: Custom field tests for site resource are in resources_acceptance_tests_customfields package

func TestAccSiteResource_basic(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-site")
	slug := testutil.RandomSlug("tf-test-site")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckSiteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSiteResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_site.test", "id"),
					resource.TestCheckResourceAttr("netbox_site.test", "name", name),
					resource.TestCheckResourceAttr("netbox_site.test", "slug", slug),
					resource.TestCheckResourceAttr("netbox_site.test", "status", "active"),
				),
			},
		},
	})
}

func TestAccSiteResource_full(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-site-full")
	slug := testutil.RandomSlug("tf-test-site-full")
	description := testutil.RandomName("description")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckSiteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSiteResourceConfig_full(name, slug, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_site.test", "id"),
					resource.TestCheckResourceAttr("netbox_site.test", "name", name),
					resource.TestCheckResourceAttr("netbox_site.test", "slug", slug),
					resource.TestCheckResourceAttr("netbox_site.test", "status", "active"),
					resource.TestCheckResourceAttr("netbox_site.test", "description", description),
				),
			},
		},
	})
}

func TestAccSiteResource_update(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-site-update")
	slug := testutil.RandomSlug("tf-test-site-upd")
	updatedName := testutil.RandomName("tf-test-site-updated")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckSiteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSiteResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_site.test", "id"),
					resource.TestCheckResourceAttr("netbox_site.test", "name", name),
				),
			},
			{
				Config: testAccSiteResourceConfig_basic(updatedName, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_site.test", "id"),
					resource.TestCheckResourceAttr("netbox_site.test", "name", updatedName),
				),
			},
		},
	})
}

func TestAccSiteResource_import(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-site-import")
	slug := testutil.RandomSlug("tf-test-site-imp")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckSiteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSiteResourceConfig_import(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_site.test", "name", name),
					resource.TestCheckResourceAttr("netbox_site.test", "slug", slug),
				),
			},
			{
				ResourceName:      "netbox_site.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccConsistency_Site(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("site")
	siteSlug := testutil.RandomSlug("site")
	regionName := testutil.RandomName("region")
	regionSlug := testutil.RandomSlug("region")
	groupName := testutil.RandomName("group")
	groupSlug := testutil.RandomSlug("group")
	tenantName := testutil.RandomName("tenant")
	tenantSlug := testutil.RandomSlug("tenant")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSiteConsistencyConfig(siteName, siteSlug, regionName, regionSlug, groupName, groupSlug, tenantName, tenantSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_site.test", "name", siteName),
					resource.TestCheckResourceAttr("netbox_site.test", "region", regionName),
					resource.TestCheckResourceAttr("netbox_site.test", "group", groupName),
					resource.TestCheckResourceAttr("netbox_site.test", "tenant", tenantName),
				),
			},
			{
				PlanOnly: true,
				Config:   testAccSiteConsistencyConfig(siteName, siteSlug, regionName, regionSlug, groupName, groupSlug, tenantName, tenantSlug),
			},
		},
	})
}

func TestAccConsistency_Site_LiteralNames(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-site-lit")
	siteSlug := testutil.RandomSlug("tf-test-site-lit")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSiteResourceConfig_basic(siteName, siteSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_site.test", "id"),
					resource.TestCheckResourceAttr("netbox_site.test", "name", siteName),
					resource.TestCheckResourceAttr("netbox_site.test", "slug", siteSlug),
				),
			},
			{
				Config:   testAccSiteResourceConfig_basic(siteName, siteSlug),
				PlanOnly: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_site.test", "id"),
				),
			},
		},
	})
}

func TestAccSiteResource_IDPreservation(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-site-id")
	slug := testutil.RandomSlug("tf-test-site-id")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckSiteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSiteResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_site.test", "id"),
					resource.TestCheckResourceAttr("netbox_site.test", "name", name),
					resource.TestCheckResourceAttr("netbox_site.test", "slug", slug),
				),
			},
		},
	})
}

func testAccSiteResourceConfig_basic(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name   = %q
  slug   = %q
  status = "active"
}
`, name, slug)
}

func testAccSiteResourceConfig_full(name, slug, description string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name        = %q
  slug        = %q
  status      = "active"
  description = %q
}
`, name, slug, description)
}

func testAccSiteResourceConfig_import(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name   = %q
  slug   = %q
  status = "active"
}
`, name, slug)
}

func testAccSiteConsistencyConfig(siteName, siteSlug, regionName, regionSlug, groupName, groupSlug, tenantName, tenantSlug string) string {
	return fmt.Sprintf(`
resource "netbox_region" "test" {
  name = "%[3]s"
  slug = "%[4]s"
}

resource "netbox_site_group" "test" {
  name = "%[5]s"
  slug = "%[6]s"
}

resource "netbox_tenant" "test" {
  name = "%[7]s"
  slug = "%[8]s"
}

resource "netbox_site" "test" {
  name = "%[1]s"
  slug = "%[2]s"
  region = netbox_region.test.name
  group = netbox_site_group.test.name
  tenant = netbox_tenant.test.name
}
`, siteName, siteSlug, regionName, regionSlug, groupName, groupSlug, tenantName, tenantSlug)
}

func TestAccSiteResource_externalDeletion(t *testing.T) {
	t.Parallel()
	name := testutil.RandomName("tf-test-site-ext-del")
	slug := testutil.RandomSlug("site")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSiteResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_site.test", "id"),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					items, _, err := client.DcimAPI.DcimSitesList(context.Background()).SlugIc([]string{slug}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find site for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.DcimAPI.DcimSitesDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete site: %v", err)
					}
					t.Logf("Successfully externally deleted site with ID: %d", itemID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

// TestAccSiteResource_removeOptionalFields tests that removing previously set tenant, region, and group fields correctly sets them to null.
// This addresses the bug where removing a nullable reference field from the configuration would not clear it in NetBox,
// causing "Provider produced inconsistent result after apply" errors.
func TestAccSiteResource_removeOptionalFields(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-site-remove")
	slug := testutil.RandomSlug("tf-test-site-remove")
	tenantName := testutil.RandomName("test-tenant-site")
	tenantSlug := testutil.GenerateSlug(tenantName)
	regionName := testutil.RandomName("test-region-site")
	regionSlug := testutil.GenerateSlug(regionName)
	groupName := testutil.RandomName("test-group-site")
	groupSlug := testutil.GenerateSlug(groupName)

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(slug)
	cleanup.RegisterTenantCleanup(tenantSlug)
	cleanup.RegisterRegionCleanup(regionSlug)
	cleanup.RegisterSiteGroupCleanup(groupSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Step 1: Create with tenant, region, and group
			{
				Config: testAccSiteResourceConfig_withAllFields(name, slug, tenantName, tenantSlug, regionName, regionSlug, groupName, groupSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_site.test", "tenant"),
					resource.TestCheckResourceAttrSet("netbox_site.test", "region"),
					resource.TestCheckResourceAttrSet("netbox_site.test", "group"),
				),
			},
			// Step 2: Remove tenant, region, and group - should set to null
			{
				Config: testAccSiteResourceConfig_withoutFields(name, slug, tenantName, tenantSlug, regionName, regionSlug, groupName, groupSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("netbox_site.test", "tenant"),
					resource.TestCheckNoResourceAttr("netbox_site.test", "region"),
					resource.TestCheckNoResourceAttr("netbox_site.test", "group"),
				),
			},
			// Step 3: Re-add tenant, region, and group - should work without errors
			{
				Config: testAccSiteResourceConfig_withAllFields(name, slug, tenantName, tenantSlug, regionName, regionSlug, groupName, groupSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_site.test", "tenant"),
					resource.TestCheckResourceAttrSet("netbox_site.test", "region"),
					resource.TestCheckResourceAttrSet("netbox_site.test", "group"),
				),
			},
		},
	})
}

func testAccSiteResourceConfig_withAllFields(name, slug, tenantName, tenantSlug, regionName, regionSlug, groupName, groupSlug string) string {
	return fmt.Sprintf(`
resource "netbox_tenant" "test" {
  name = %[3]q
  slug = %[4]q
}

resource "netbox_region" "test" {
  name = %[5]q
  slug = %[6]q
}

resource "netbox_site_group" "test" {
  name = %[7]q
  slug = %[8]q
}

resource "netbox_site" "test" {
  name   = %[1]q
  slug   = %[2]q
  status = "active"
  tenant = netbox_tenant.test.id
  region = netbox_region.test.id
  group  = netbox_site_group.test.id
}
`, name, slug, tenantName, tenantSlug, regionName, regionSlug, groupName, groupSlug)
}

func testAccSiteResourceConfig_withoutFields(name, slug, tenantName, tenantSlug, regionName, regionSlug, groupName, groupSlug string) string {
	return fmt.Sprintf(`
resource "netbox_tenant" "test" {
  name = %[3]q
  slug = %[4]q
}

resource "netbox_region" "test" {
  name = %[5]q
  slug = %[6]q
}

resource "netbox_site_group" "test" {
  name = %[7]q
  slug = %[8]q
}

resource "netbox_site" "test" {
  name   = %[1]q
  slug   = %[2]q
  status = "active"
  # tenant, region, and group removed - should set to null
}
`, name, slug, tenantName, tenantSlug, regionName, regionSlug, groupName, groupSlug)
}

func TestAccSiteResource_removeDescriptionAndComments(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-site-optional")
	siteSlug := testutil.RandomName("tf-test-site-optional")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteName)

	testutil.TestRemoveOptionalFields(t, testutil.MultiFieldOptionalTestConfig{
		ResourceName: "netbox_site",
		BaseConfig: func() string {
			return testAccSiteResourceConfig_minimal(siteName, siteSlug)
		},
		ConfigWithFields: func() string {
			return testAccSiteResourceConfig_withDescriptionAndComments(
				siteName,
				siteSlug,
				"Test description",
				"Test comments",
			)
		},
		OptionalFields: map[string]string{
			"description": "Test description",
			"comments":    "Test comments",
		},
		RequiredFields: map[string]string{
			"name": siteName,
			"slug": siteSlug,
		},
		CheckDestroy: testutil.CheckSiteDestroy,
	})
}

func testAccSiteResourceConfig_minimal(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name   = %[1]q
  slug   = %[2]q
  status = "active"
}
`, name, slug)
}

func testAccSiteResourceConfig_withDescriptionAndComments(name, slug, description, comments string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name        = %[1]q
  slug        = %[2]q
  status      = "active"
  description = %[3]q
  comments    = %[4]q
}
`, name, slug, description, comments)
}
func TestAccSiteResource_validationErrors(t *testing.T) {
	t.Parallel()

	testutil.RunMultiValidationErrorTest(t, testutil.MultiValidationErrorTestConfig{
		ResourceName: "netbox_site",
		TestCases: map[string]testutil.ValidationErrorCase{
			"missing_name": {
				Config: func() string {
					return `
resource "netbox_site" "test" {
  slug = "test-site"
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
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
			"invalid_status": {
				Config: func() string {
					return `
resource "netbox_site" "test" {
  name   = "Test Site"
  slug   = "test-site"
  status = "invalid_status"
}
`
				},
				ExpectedError: testutil.ErrPatternInvalidEnum,
			},
			"invalid_region_reference": {
				Config: func() string {
					return `
resource "netbox_site" "test" {
  name   = "Test Site"
  slug   = "test-site"
  region = "99999999"
}
`
				},
				ExpectedError: testutil.ErrPatternNotFound,
			},
			"invalid_group_reference": {
				Config: func() string {
					return `
resource "netbox_site" "test" {
  name  = "Test Site"
  slug  = "test-site"
  group = "99999999"
}
`
				},
				ExpectedError: testutil.ErrPatternNotFound,
			},
			"invalid_tenant_reference": {
				Config: func() string {
					return `
resource "netbox_site" "test" {
  name   = "Test Site"
  slug   = "test-site"
  tenant = "99999999"
}
`
				},
				ExpectedError: testutil.ErrPatternNotFound,
			},
		},
		CheckDestroy: testutil.CheckSiteDestroy,
	})
}
