package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRackResource_basic(t *testing.T) {
	t.Parallel()

	// Generate unique names to avoid conflicts between test runs
	siteName := testutil.RandomName("tf-test-rack-site")
	siteSlug := testutil.RandomSlug("tf-test-rack-site")
	rackName := testutil.RandomName("tf-test-rack")

	// Register cleanup to ensure resources are deleted even if test fails
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRackCleanup(rackName)
	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.ComposeCheckDestroy(testutil.CheckRackDestroy, testutil.CheckSiteDestroy),
		Steps: []resource.TestStep{
			{
				Config: testAccRackResourceConfig_basic(siteName, siteSlug, rackName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_rack.test", "id"),
					resource.TestCheckResourceAttr("netbox_rack.test", "name", rackName),
					resource.TestCheckResourceAttrPair("netbox_rack.test", "site", "netbox_site.test", "id"),
				),
			},
		},
	})
}

func TestAccRackResource_full(t *testing.T) {
	t.Parallel()

	// Generate unique names
	siteName := testutil.RandomName("tf-test-rack-site-full")
	siteSlug := testutil.RandomSlug("tf-test-rack-s-full")
	rackName := testutil.RandomName("tf-test-rack-full")
	description := testutil.RandomName("description")

	// Register cleanup
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRackCleanup(rackName)
	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.ComposeCheckDestroy(testutil.CheckRackDestroy, testutil.CheckSiteDestroy),
		Steps: []resource.TestStep{
			{
				Config: testAccRackResourceConfig_full(siteName, siteSlug, rackName, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_rack.test", "id"),
					resource.TestCheckResourceAttr("netbox_rack.test", "name", rackName),
					resource.TestCheckResourceAttr("netbox_rack.test", "status", "active"),
					resource.TestCheckResourceAttr("netbox_rack.test", "description", description),
					resource.TestCheckResourceAttr("netbox_rack.test", "u_height", "42"),
					resource.TestCheckResourceAttr("netbox_rack.test", "width", "19"),
				),
			},
		},
	})
}

func TestAccRackResource_update(t *testing.T) {
	t.Parallel()

	// Generate unique names
	siteName := testutil.RandomName("tf-test-rack-site-upd")
	siteSlug := testutil.RandomSlug("tf-test-rack-s-upd")
	rackName := testutil.RandomName("tf-test-rack-upd")
	updatedName := testutil.RandomName("tf-test-rack-upd-name")

	// Register cleanup (use original name for initial cleanup, register updated name too)
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRackCleanup(rackName)
	cleanup.RegisterRackCleanup(updatedName)
	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.ComposeCheckDestroy(testutil.CheckRackDestroy, testutil.CheckSiteDestroy),
		Steps: []resource.TestStep{
			{
				Config: testAccRackResourceConfig_basic(siteName, siteSlug, rackName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_rack.test", "id"),
					resource.TestCheckResourceAttr("netbox_rack.test", "name", rackName),
				),
			},
			{
				Config: testAccRackResourceConfig_basic(siteName, siteSlug, updatedName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_rack.test", "id"),
					resource.TestCheckResourceAttr("netbox_rack.test", "name", updatedName),
				),
			},
		},
	})
}

func TestAccRackResource_withLocation(t *testing.T) {
	t.Parallel()

	// Generate unique names
	siteName := testutil.RandomName("tf-test-rack-site-loc")
	siteSlug := testutil.RandomSlug("tf-test-rack-s-loc")
	locationName := testutil.RandomName("tf-test-rack-location")
	locationSlug := testutil.RandomSlug("tf-test-rack-loc")
	rackName := testutil.RandomName("tf-test-rack-with-loc")

	// Register cleanup (rack first, then location, then site due to dependency)
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRackCleanup(rackName)
	cleanup.RegisterLocationCleanup(locationSlug)
	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.ComposeCheckDestroy(testutil.CheckRackDestroy, testutil.CheckLocationDestroy, testutil.CheckSiteDestroy),
		Steps: []resource.TestStep{
			{
				Config: testAccRackResourceConfig_withLocation(siteName, siteSlug, locationName, locationSlug, rackName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_rack.test", "id"),
					resource.TestCheckResourceAttr("netbox_rack.test", "name", rackName),
					resource.TestCheckResourceAttrPair("netbox_rack.test", "location", "netbox_location.test", "id"),
				),
			},
		},
	})
}

func TestAccRackResource_import(t *testing.T) {
	t.Parallel()

	// Generate unique names to avoid conflicts between test runs
	siteName := testutil.RandomName("tf-test-rack-site")
	siteSlug := testutil.RandomSlug("tf-test-rack-site")
	rackName := testutil.RandomName("tf-test-rack")

	// Register cleanup to ensure resources are deleted even if test fails
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRackCleanup(rackName)
	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.ComposeCheckDestroy(testutil.CheckRackDestroy, testutil.CheckSiteDestroy),
		Steps: []resource.TestStep{
			{
				Config: testAccRackResourceConfig_import(siteName, siteSlug, rackName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_rack.test", "id"),
					resource.TestCheckResourceAttr("netbox_rack.test", "name", rackName),
					resource.TestCheckResourceAttrPair("netbox_rack.test", "site", "netbox_site.test", "id"),
				),
			},
			{
				ResourceName:            "netbox_rack.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"site"},
			},
			{
				Config:   testAccRackResourceConfig_import(siteName, siteSlug, rackName),
				PlanOnly: true,
			},
		},
	})
}

func TestAccRackResource_tagLifecycle(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-rack-site-tags")
	siteSlug := testutil.RandomSlug("tf-test-rack-site-tags")
	rackName := testutil.RandomName("tf-test-rack-tags")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2Slug := testutil.RandomSlug("tag2")
	tag3Slug := testutil.RandomSlug("tag3")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRackCleanup(rackName)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)
	cleanup.RegisterTagCleanup(tag3Slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.ComposeCheckDestroy(testutil.CheckRackDestroy, testutil.CheckSiteDestroy),
		Steps: []resource.TestStep{
			{
				Config: testAccRackResourceConfig_tags(siteName, siteSlug, rackName, tag1Slug, tag2Slug, tag3Slug, caseTag1Tag2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_rack.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_rack.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag1-%s", tag1Slug),
						"slug": tag1Slug,
					}),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_rack.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag2-%s", tag2Slug),
						"slug": tag2Slug,
					}),
				),
			},
			{
				Config: testAccRackResourceConfig_tags(siteName, siteSlug, rackName, tag1Slug, tag2Slug, tag3Slug, caseTag1Uscore2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_rack.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_rack.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag1-%s", tag1Slug),
						"slug": tag1Slug,
					}),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_rack.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag2-%s", tag2Slug),
						"slug": tag2Slug,
					}),
				),
			},
			{
				Config: testAccRackResourceConfig_tags(siteName, siteSlug, rackName, tag1Slug, tag2Slug, tag3Slug, caseTag3),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_rack.test", "tags.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_rack.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag3-%s", tag3Slug),
						"slug": tag3Slug,
					}),
				),
			},
			{
				Config: testAccRackResourceConfig_tags(siteName, siteSlug, rackName, tag1Slug, tag2Slug, tag3Slug, tagsEmpty),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_rack.test", "tags.#", "0"),
				),
			},
		},
	})
}

func TestAccRackResource_tagOrderInvariance(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-rack-site-tag-order")
	siteSlug := testutil.RandomSlug("tf-test-rack-site-tag-order")
	rackName := testutil.RandomName("tf-test-rack-tag-order")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2Slug := testutil.RandomSlug("tag2")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRackCleanup(rackName)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.ComposeCheckDestroy(testutil.CheckRackDestroy, testutil.CheckSiteDestroy),
		Steps: []resource.TestStep{
			{
				Config: testAccRackResourceConfig_tagsOrder(siteName, siteSlug, rackName, tag1Slug, tag2Slug, caseTag1Tag2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_rack.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_rack.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag1-%s", tag1Slug),
						"slug": tag1Slug,
					}),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_rack.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag2-%s", tag2Slug),
						"slug": tag2Slug,
					}),
				),
			},
			{
				Config: testAccRackResourceConfig_tagsOrder(siteName, siteSlug, rackName, tag1Slug, tag2Slug, caseTag2Uscore1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_rack.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_rack.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag1-%s", tag1Slug),
						"slug": tag1Slug,
					}),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_rack.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag2-%s", tag2Slug),
						"slug": tag2Slug,
					}),
				),
			},
		},
	})
}

// NOTE: Custom field tests for rack resource are in resources_acceptance_tests_customfields package.
func TestAccConsistency_Rack(t *testing.T) {
	t.Parallel()

	rackName := testutil.RandomName("rack")
	siteName := testutil.RandomName("site")
	siteSlug := testutil.RandomSlug("site")
	tenantName := testutil.RandomName("tenant")
	tenantSlug := testutil.RandomSlug("tenant")
	roleName := testutil.RandomName("role")
	roleSlug := testutil.RandomSlug("role")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterTenantCleanup(tenantSlug)
	cleanup.RegisterRackRoleCleanup(roleSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRackConsistencyConfig(rackName, siteName, siteSlug, tenantName, tenantSlug, roleName, roleSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_rack.test", "name", rackName),
					resource.TestCheckResourceAttr("netbox_rack.test", "site", siteName),
					resource.TestCheckResourceAttr("netbox_rack.test", "tenant", tenantName),
					resource.TestCheckResourceAttr("netbox_rack.test", "role", roleName),
				),
			},
			{
				PlanOnly: true,
				Config:   testAccRackConsistencyConfig(rackName, siteName, siteSlug, tenantName, tenantSlug, roleName, roleSlug),
			},
		},
	})
}

func TestAccConsistency_Rack_LiteralNames(t *testing.T) {
	t.Parallel()

	rackName := testutil.RandomName("tf-test-rack-lit")
	siteName := testutil.RandomName("tf-test-site")
	siteSlug := testutil.RandomSlug("tf-test-site")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRackConsistencyLiteralNamesConfig(rackName, siteName, siteSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_rack.test", "id"),
					resource.TestCheckResourceAttr("netbox_rack.test", "name", rackName),
					resource.TestCheckResourceAttr("netbox_rack.test", "site", siteName),
				),
			},
			{
				Config:   testAccRackConsistencyLiteralNamesConfig(rackName, siteName, siteSlug),
				PlanOnly: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_rack.test", "id"),
				),
			},
		},
	})
}

func testAccRackConsistencyLiteralNamesConfig(rackName, siteName, siteSlug string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name = %q
  slug = %q
}

resource "netbox_rack" "test" {
  name = %q
  site = netbox_site.test.name
}
`, siteName, siteSlug, rackName)
}

// testAccRackResourceConfig_basic returns a basic test configuration.
func testAccRackResourceConfig_basic(siteName, siteSlug, rackName string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name   = %q
  slug   = %q
  status = "active"
}

resource "netbox_rack" "test" {
  name   = %q
  site   = netbox_site.test.id
  status = "active"
}
`, siteName, siteSlug, rackName)
}

// testAccRackResourceConfig_full returns a test configuration with all fields.
func testAccRackResourceConfig_full(siteName, siteSlug, rackName, description string) string {
	return fmt.Sprintf(`
terraform {
  required_providers {
    netbox = {
      source = "bab3l/netbox"
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

resource "netbox_rack" "test" {
  name        = %q
  site        = netbox_site.test.id
  status      = "active"
  u_height    = 42
  width       = 19
  description = %q
}
`, siteName, siteSlug, rackName, description)
}

func testAccRackResourceConfig_tags(siteName, siteSlug, rackName, tag1Slug, tag2Slug, tag3Slug, tagCase string) string {
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
resource "netbox_site" "test" {
	name   = %[1]q
	slug   = %[2]q
	status = "active"
}

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

resource "netbox_rack" "test" {
	name   = %[3]q
	site   = netbox_site.test.id
	status = "active"
	%[7]s
}
`, siteName, siteSlug, rackName, tag1Slug, tag2Slug, tag3Slug, tagsConfig)
}

func testAccRackResourceConfig_tagsOrder(siteName, siteSlug, rackName, tag1Slug, tag2Slug, tagCase string) string {
	var tagsConfig string
	switch tagCase {
	case caseTag1Tag2:
		tagsConfig = tagsDoubleNested
	case caseTag2Uscore1:
		tagsConfig = tagsDoubleNestedReversed
	}

	return fmt.Sprintf(`
resource "netbox_site" "test" {
	name   = %[1]q
	slug   = %[2]q
	status = "active"
}

resource "netbox_tag" "tag1" {
	name = "Tag1-%[4]s"
	slug = %[4]q
}

resource "netbox_tag" "tag2" {
	name = "Tag2-%[5]s"
	slug = %[5]q
}

resource "netbox_rack" "test" {
	name   = %[3]q
	site   = netbox_site.test.id
	status = "active"
	%[6]s
}
`, siteName, siteSlug, rackName, tag1Slug, tag2Slug, tagsConfig)
}

// testAccRackResourceConfig_withLocation returns a test configuration with location.
func testAccRackResourceConfig_withLocation(siteName, siteSlug, locationName, locationSlug, rackName string) string {
	return fmt.Sprintf(`
terraform {
  required_providers {
    netbox = {
      source = "bab3l/netbox"
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
  name = %q
  slug = %q
  site = netbox_site.test.id
}

resource "netbox_rack" "test" {
  name     = %q
  site     = netbox_site.test.id
  location = netbox_location.test.id
}
`, siteName, siteSlug, locationName, locationSlug, rackName)
}

func testAccRackResourceConfig_import(siteName, siteSlug, rackName string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_rack" "test" {
  name = %[3]q
  site = netbox_site.test.id
}
`, siteName, siteSlug, rackName)
}

func testAccRackConsistencyConfig(rackName, siteName, siteSlug, tenantName, tenantSlug, roleName, roleSlug string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name = "%[2]s"
  slug = "%[3]s"
}

resource "netbox_tenant" "test" {
  name = "%[4]s"
  slug = "%[5]s"
}

resource "netbox_rack_role" "test" {
  name = "%[6]s"
  slug = "%[7]s"
}

resource "netbox_rack" "test" {
  name = "%[1]s"
  site = netbox_site.test.name
  tenant = netbox_tenant.test.name
  role = netbox_rack_role.test.name
}
`, rackName, siteName, siteSlug, tenantName, tenantSlug, roleName, roleSlug)
}

func TestAccRackResource_externalDeletion(t *testing.T) {
	t.Parallel()
	testutil.TestAccPreCheck(t)

	siteName := testutil.RandomName("tf-test-rack-site-extdel")
	siteSlug := testutil.RandomSlug("tf-test-rack-site-ed")
	rackName := testutil.RandomName("tf-test-rack-extdel")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRackCleanup(rackName)
	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRackResourceConfig_basic(siteName, siteSlug, rackName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_rack.test", "id"),
					resource.TestCheckResourceAttr("netbox_rack.test", "name", rackName),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}

					racks, _, err := client.DcimAPI.DcimRacksList(context.Background()).Name([]string{rackName}).Execute()
					if err != nil || racks == nil || len(racks.Results) == 0 {
						t.Fatalf("Failed to find rack for external deletion: %v", err)
					}
					rackID := racks.Results[0].Id
					_, err = client.DcimAPI.DcimRacksDestroy(context.Background(), rackID).Execute()
					if err != nil {
						t.Fatalf("Failed to delete rack: %v", err)
					}
					t.Logf("Successfully externally deleted rack with ID: %d", rackID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

// TestAccRackResource_removePhysicalFields tests that optional physical fields
// can be successfully removed from the configuration without causing inconsistent state.
func TestAccRackResource_removePhysicalFields(t *testing.T) {
	t.Parallel()

	// Random names
	siteName := testutil.RandomName("tf-test-site-rack-phys")
	siteSlug := testutil.RandomSlug("tf-test-site-rack-phys")
	locationName := testutil.RandomName("tf-test-loc-rack-phys")
	locationSlug := testutil.RandomSlug("tf-test-loc-rack-phys")
	tenantName := testutil.RandomName("tf-test-tenant-rack-phys")
	tenantSlug := testutil.RandomSlug("tf-test-tenant-rack-phys")
	roleName := testutil.RandomName("tf-test-role-rack-phys")
	roleSlug := testutil.RandomSlug("tf-test-role-rack-phys")
	rackName := testutil.RandomName("tf-test-rack-phys")

	// Cleanup
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRackCleanup(rackName)
	cleanup.RegisterLocationCleanup(locationSlug)
	cleanup.RegisterTenantCleanup(tenantSlug)
	cleanup.RegisterRackRoleCleanup(roleSlug)
	cleanup.RegisterSiteCleanup(siteSlug)

	testutil.TestRemoveOptionalFields(t, testutil.MultiFieldOptionalTestConfig{
		ResourceName: "netbox_rack",
		BaseConfig: func() string {
			return testAccRackResourceConfig_physicalFieldsRemoved(
				siteName, siteSlug,
				locationName, locationSlug,
				tenantName, tenantSlug,
				roleName, roleSlug,
				rackName,
			)
		},
		ConfigWithFields: func() string {
			return testAccRackResourceConfig_physicalFields(
				siteName, siteSlug,
				locationName, locationSlug,
				tenantName, tenantSlug,
				roleName, roleSlug,
				rackName,
			)
		},
		OptionalFields: map[string]string{
			"airflow":        "front-to-rear",
			"form_factor":    "4-post-cabinet",
			"max_weight":     "1000",
			"mounting_depth": "100",
			"outer_depth":    "120",
			"outer_width":    "80",
			"weight":         "50.5",
			"serial":         "SN123456",
			"asset_tag":      "TAG123456",
			"description":    "Test description",
			"comments":       "Test comments",
		},
		RequiredFields: map[string]string{
			"name": rackName,
		},
		CheckDestroy: testutil.CheckRackDestroy,
	})
}

func testAccRackResourceConfig_physicalFields(siteName, siteSlug, locationName, locationSlug, tenantName, tenantSlug, roleName, roleSlug, rackName string) string {
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

resource "netbox_tenant" "test" {
  name = %[5]q
  slug = %[6]q
}

resource "netbox_rack_role" "test" {
  name = %[7]q
  slug = %[8]q
}

resource "netbox_rack" "test" {
  name           = %[9]q
  site           = netbox_site.test.id
  status         = "active"
  location       = netbox_location.test.id
  tenant         = netbox_tenant.test.id
  role           = netbox_rack_role.test.id
  airflow        = "front-to-rear"
  desc_units     = true
  form_factor    = "4-post-cabinet"
  max_weight     = "1000"
  mounting_depth = "100"
  outer_depth    = "120"
  outer_unit     = "mm"
  outer_width    = "80"
  starting_unit  = "1"
  u_height       = "48"
  weight         = "50.5"
  weight_unit    = "kg"
  width          = "19"
  serial         = "SN123456"
  asset_tag      = "TAG123456"
  description    = "Test description"
  comments       = "Test comments"
}
`, siteName, siteSlug, locationName, locationSlug, tenantName, tenantSlug, roleName, roleSlug, rackName)
}

func testAccRackResourceConfig_physicalFieldsRemoved(siteName, siteSlug, locationName, locationSlug, tenantName, tenantSlug, roleName, roleSlug, rackName string) string {
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

resource "netbox_tenant" "test" {
  name = %[5]q
  slug = %[6]q
}

resource "netbox_rack_role" "test" {
  name = %[7]q
  slug = %[8]q
}

resource "netbox_rack" "test" {
  name           = %[9]q
  site           = netbox_site.test.id
  status         = "active"
}
`, siteName, siteSlug, locationName, locationSlug, tenantName, tenantSlug, roleName, roleSlug, rackName)
}

func TestAccRackResource_removeReferenceFields(t *testing.T) {
	t.Parallel()

	// Random names
	siteName := testutil.RandomName("tf-test-site-rack-ref")
	siteSlug := testutil.RandomSlug("tf-test-site-rack-ref")
	mfgName := testutil.RandomName("tf-test-mfg-rack-ref")
	mfgSlug := testutil.RandomSlug("tf-test-mfg-rack-ref")
	rackTypeName := testutil.RandomName("tf-test-racktype-ref")
	rackTypeSlug := testutil.RandomSlug("tf-test-racktype-ref")
	rackName := testutil.RandomName("tf-test-rack-ref")

	// Cleanup
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRackCleanup(rackName)
	cleanup.RegisterRackTypeCleanup(rackTypeSlug)
	cleanup.RegisterManufacturerCleanup(mfgSlug)
	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckRackDestroy,
		Steps: []resource.TestStep{
			// Step 1: Create resource with rack_type reference
			{
				Config: testAccRackResourceConfig_referenceFields(
					siteName, siteSlug,
					mfgName, mfgSlug,
					rackTypeName, rackTypeSlug,
					rackName,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_rack.test", "name", rackName),
					resource.TestCheckResourceAttrPair("netbox_rack.test", "rack_type", "netbox_rack_type.test", "id"),
				),
			},
			// Step 2: Remove rack_type reference but keep rack_type resource
			{
				Config: testAccRackResourceConfig_referenceFieldsRemoved(
					siteName, siteSlug,
					mfgName, mfgSlug,
					rackTypeName, rackTypeSlug,
					rackName,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_rack.test", "name", rackName),
					resource.TestCheckNoResourceAttr("netbox_rack.test", "rack_type"),
				),
			},
			// Step 3: Re-add rack_type to verify it can be set again
			{
				Config: testAccRackResourceConfig_referenceFields(
					siteName, siteSlug,
					mfgName, mfgSlug,
					rackTypeName, rackTypeSlug,
					rackName,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_rack.test", "name", rackName),
					resource.TestCheckResourceAttrPair("netbox_rack.test", "rack_type", "netbox_rack_type.test", "id"),
				),
			},
		},
	})
}

func testAccRackResourceConfig_referenceFieldsRemoved(siteName, siteSlug, mfgName, mfgSlug, rackTypeName, rackTypeSlug, rackName string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name   = %[1]q
  slug   = %[2]q
  status = "active"
}

resource "netbox_manufacturer" "test" {
  name = %[3]q
  slug = %[4]q
}

resource "netbox_rack_type" "test" {
  model        = %[5]q
  slug         = %[6]q
  manufacturer = netbox_manufacturer.test.id
  form_factor  = "4-post-cabinet"
  u_height     = 42
  width        = 19
  weight_unit  = "kg"
}

resource "netbox_rack" "test" {
  name      = %[7]q
  site      = netbox_site.test.id
  status    = "active"
  # rack_type removed
}
`, siteName, siteSlug, mfgName, mfgSlug, rackTypeName, rackTypeSlug, rackName)
}

func testAccRackResourceConfig_referenceFields(siteName, siteSlug, mfgName, mfgSlug, rackTypeName, rackTypeSlug, rackName string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name   = %[1]q
  slug   = %[2]q
  status = "active"
}

resource "netbox_manufacturer" "test" {
  name = %[3]q
  slug = %[4]q
}

resource "netbox_rack_type" "test" {
  model        = %[5]q
  slug         = %[6]q
  manufacturer = netbox_manufacturer.test.id
  form_factor  = "4-post-cabinet"
  u_height     = 42
  width        = 19
  weight_unit  = "kg"
}

resource "netbox_rack" "test" {
  name      = %[7]q
  site      = netbox_site.test.id
  status    = "active"
  rack_type = netbox_rack_type.test.id
}
`, siteName, siteSlug, mfgName, mfgSlug, rackTypeName, rackTypeSlug, rackName)
}

// TestAccRackResource_validationErrors tests validation error scenarios.
func TestAccRackResource_validationErrors(t *testing.T) {
	testutil.RunMultiValidationErrorTest(t, testutil.MultiValidationErrorTestConfig{
		ResourceName: "netbox_rack",
		TestCases: map[string]testutil.ValidationErrorCase{
			"missing_name": {
				Config: func() string {
					return `
resource "netbox_site" "test" {
  name = "Test Site"
  slug = "test-site"
}

resource "netbox_rack" "test" {
  site = netbox_site.test.id
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
			"missing_site": {
				Config: func() string {
					return `
resource "netbox_rack" "test" {
  name = "Test Rack"
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
			"invalid_status": {
				Config: func() string {
					return `
resource "netbox_site" "test" {
  name = "Test Site"
  slug = "test-site"
}

resource "netbox_rack" "test" {
  name   = "Test Rack"
  site   = netbox_site.test.id
  status = "invalid_status"
}
`
				},
				ExpectedError: testutil.ErrPatternInvalidEnum,
			},
			"invalid_role_reference": {
				Config: func() string {
					return `
resource "netbox_site" "test" {
  name = "Test Site"
  slug = "test-site"
}

resource "netbox_rack" "test" {
  name = "Test Rack"
  site = netbox_site.test.id
  role = "99999"
}
`
				},
				ExpectedError: testutil.ErrPatternNotFound,
			},
			"invalid_tenant_reference": {
				Config: func() string {
					return `
resource "netbox_site" "test" {
  name = "Test Site"
  slug = "test-site"
}

resource "netbox_rack" "test" {
  name   = "Test Rack"
  site   = netbox_site.test.id
  tenant = "99999"
}
`
				},
				ExpectedError: testutil.ErrPatternNotFound,
			},
			"invalid_location_reference": {
				Config: func() string {
					return `
resource "netbox_site" "test" {
  name = "Test Site"
  slug = "test-site"
}

resource "netbox_rack" "test" {
  name     = "Test Rack"
  site     = netbox_site.test.id
  location = "99999"
}
`
				},
				ExpectedError: testutil.ErrPatternNotFound,
			},
		},
	})
}
