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
			{
				Config:   testAccSiteResourceConfig_basic(name, slug),
				PlanOnly: true,
			},
		},
	})
}

func TestAccSiteResource_full(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-site-full")
	slug := testutil.RandomSlug("tf-test-site-full")
	description := testutil.RandomName("description")
	updatedDescription := "Updated site description"
	tagName1 := testutil.RandomName("tag1")
	tagSlug1 := testutil.RandomSlug("tag1")
	tagName2 := testutil.RandomName("tag2")
	tagSlug2 := testutil.RandomSlug("tag2")
	physicalAddress := "123 Terraform Ave, Example City"
	shippingAddress := "PO Box 123, Example City"
	timeZone := "America/Los_Angeles"
	latitude := 37.7749
	longitude := -122.4194

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(slug)
	cleanup.RegisterTagCleanup(tagSlug1)
	cleanup.RegisterTagCleanup(tagSlug2)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckSiteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSiteResourceConfig_full(name, slug, description, tagName1, tagSlug1, tagName2, tagSlug2, physicalAddress, shippingAddress, timeZone, latitude, longitude),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_site.test", "id"),
					resource.TestCheckResourceAttr("netbox_site.test", "name", name),
					resource.TestCheckResourceAttr("netbox_site.test", "slug", slug),
					resource.TestCheckResourceAttr("netbox_site.test", "status", "active"),
					resource.TestCheckResourceAttr("netbox_site.test", "description", description),
					resource.TestCheckResourceAttr("netbox_site.test", "time_zone", timeZone),
					resource.TestCheckResourceAttr("netbox_site.test", "physical_address", physicalAddress),
					resource.TestCheckResourceAttr("netbox_site.test", "shipping_address", shippingAddress),
					resource.TestCheckResourceAttr("netbox_site.test", "latitude", fmt.Sprintf("%g", latitude)),
					resource.TestCheckResourceAttr("netbox_site.test", "longitude", fmt.Sprintf("%g", longitude)),
					resource.TestCheckResourceAttr("netbox_site.test", "tags.#", "2"),
				),
			},
			{
				Config:   testAccSiteResourceConfig_full(name, slug, description, tagName1, tagSlug1, tagName2, tagSlug2, physicalAddress, shippingAddress, timeZone, latitude, longitude),
				PlanOnly: true,
			},
			{
				Config: testAccSiteResourceConfig_fullUpdate(name, slug, updatedDescription, tagName1, tagSlug1, tagName2, tagSlug2, physicalAddress, shippingAddress, timeZone, latitude, longitude),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_site.test", "description", updatedDescription),
				),
			},
			{
				Config:   testAccSiteResourceConfig_fullUpdate(name, slug, updatedDescription, tagName1, tagSlug1, tagName2, tagSlug2, physicalAddress, shippingAddress, timeZone, latitude, longitude),
				PlanOnly: true,
			},
		},
	})
}

func TestAccSiteResource_tagLifecycle(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-site-tags")
	slug := testutil.RandomSlug("tf-test-site-tags")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2Slug := testutil.RandomSlug("tag2")
	tag3Slug := testutil.RandomSlug("tag3")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(slug)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)
	cleanup.RegisterTagCleanup(tag3Slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckSiteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSiteResourceConfig_tags(name, slug, tag1Slug, tag2Slug, tag3Slug, caseTag1Tag2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_site.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemAttr("netbox_site.test", "tags.*", tag1Slug),
					resource.TestCheckTypeSetElemAttr("netbox_site.test", "tags.*", tag2Slug),
				),
			},
			{
				Config: testAccSiteResourceConfig_tags(name, slug, tag1Slug, tag2Slug, tag3Slug, caseTag1Uscore2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_site.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemAttr("netbox_site.test", "tags.*", tag1Slug),
					resource.TestCheckTypeSetElemAttr("netbox_site.test", "tags.*", tag2Slug),
				),
			},
			{
				Config: testAccSiteResourceConfig_tags(name, slug, tag1Slug, tag2Slug, tag3Slug, caseTag3),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_site.test", "tags.#", "1"),
					resource.TestCheckTypeSetElemAttr("netbox_site.test", "tags.*", tag3Slug),
				),
			},
			{
				Config: testAccSiteResourceConfig_tags(name, slug, tag1Slug, tag2Slug, tag3Slug, tagsEmpty),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_site.test", "tags.#", "0"),
				),
			},
		},
	})
}

func TestAccSiteResource_tagOrderInvariance(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-site-tag-order")
	slug := testutil.RandomSlug("tf-test-site-tag-order")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2Slug := testutil.RandomSlug("tag2")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(slug)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckSiteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSiteResourceConfig_tagsOrder(name, slug, tag1Slug, tag2Slug, caseTag1Tag2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_site.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemAttr("netbox_site.test", "tags.*", tag1Slug),
					resource.TestCheckTypeSetElemAttr("netbox_site.test", "tags.*", tag2Slug),
				),
			},
			{
				Config: testAccSiteResourceConfig_tagsOrder(name, slug, tag1Slug, tag2Slug, caseTag2Uscore1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_site.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemAttr("netbox_site.test", "tags.*", tag1Slug),
					resource.TestCheckTypeSetElemAttr("netbox_site.test", "tags.*", tag2Slug),
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
				Config:   testAccSiteResourceConfig_basic(name, slug),
				PlanOnly: true,
			},
			{
				Config: testAccSiteResourceConfig_basic(updatedName, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_site.test", "id"),
					resource.TestCheckResourceAttr("netbox_site.test", "name", updatedName),
				),
			},
			{
				Config:   testAccSiteResourceConfig_basic(updatedName, slug),
				PlanOnly: true,
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
				Config:   testAccSiteResourceConfig_import(name, slug),
				PlanOnly: true,
			},
			{
				ResourceName:      "netbox_site.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:   testAccSiteResourceConfig_import(name, slug),
				PlanOnly: true,
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

func testAccSiteResourceConfig_basic(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name   = %q
  slug   = %q
  status = "active"
}
`, name, slug)
}

func testAccSiteResourceConfig_full(name, slug, description, tagName1, tagSlug1, tagName2, tagSlug2, physicalAddress, shippingAddress, timeZone string, latitude, longitude float64) string {
	return fmt.Sprintf(`
resource "netbox_tag" "tag1" {
	name = %[4]q
	slug = %[5]q
}

resource "netbox_tag" "tag2" {
	name = %[6]q
	slug = %[7]q
}

resource "netbox_site" "test" {
	name        = %[1]q
	slug        = %[2]q
	status      = "active"
	description = %[3]q
	physical_address = %[8]q
	shipping_address = %[9]q
	time_zone = %[10]q
	latitude  = %[11]f
	longitude = %[12]f

	tags = [
		netbox_tag.tag1.slug,
		netbox_tag.tag2.slug
	]
}
`, name, slug, description, tagName1, tagSlug1, tagName2, tagSlug2, physicalAddress, shippingAddress, timeZone, latitude, longitude)
}

func testAccSiteResourceConfig_fullUpdate(name, slug, description, tagName1, tagSlug1, tagName2, tagSlug2, physicalAddress, shippingAddress, timeZone string, latitude, longitude float64) string {
	return fmt.Sprintf(`
resource "netbox_tag" "tag1" {
	name = %[4]q
	slug = %[5]q
}

resource "netbox_tag" "tag2" {
	name = %[6]q
	slug = %[7]q
}

resource "netbox_site" "test" {
	name        = %[1]q
	slug        = %[2]q
	status      = "active"
	description = %[3]q
	comments    = "Updated comments"
	physical_address = %[8]q
	shipping_address = %[9]q
	time_zone = %[10]q
	latitude  = %[11]f
	longitude = %[12]f

	tags = [
		netbox_tag.tag1.slug,
		netbox_tag.tag2.slug
	]
}
`, name, slug, description, tagName1, tagSlug1, tagName2, tagSlug2, physicalAddress, shippingAddress, timeZone, latitude, longitude)
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

func testAccSiteResourceConfig_tags(name, slug, tag1Slug, tag2Slug, tag3Slug, tagCase string) string {
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
  name = "Tag1-%[3]s"
  slug = %[3]q
}

resource "netbox_tag" "tag2" {
  name = "Tag2-%[4]s"
  slug = %[4]q
}

resource "netbox_tag" "tag3" {
  name = "Tag3-%[5]s"
  slug = %[5]q
}

resource "netbox_site" "test" {
  name   = %[1]q
  slug   = %[2]q
  status = "active"
  %[6]s
}
`, name, slug, tag1Slug, tag2Slug, tag3Slug, tagsConfig)
}

func testAccSiteResourceConfig_tagsOrder(name, slug, tag1Slug, tag2Slug, tagCase string) string {
	var tagsConfig string
	switch tagCase {
	case caseTag1Tag2:
		tagsConfig = tagsDoubleSlug
	case caseTag2Uscore1:
		tagsConfig = tagsDoubleSlugReversed
	}

	return fmt.Sprintf(`
resource "netbox_tag" "tag1" {
  name = "Tag1-%[3]s"
  slug = %[3]q
}

resource "netbox_tag" "tag2" {
  name = "Tag2-%[4]s"
  slug = %[4]q
}

resource "netbox_site" "test" {
  name   = %[1]q
  slug   = %[2]q
  status = "active"
  %[5]s
}
`, name, slug, tag1Slug, tag2Slug, tagsConfig)
}
