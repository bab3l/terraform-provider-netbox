package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPlatformResource_basic(t *testing.T) {
	t.Parallel()

	platformName := testutil.RandomName("tf-test-platform")
	platformSlug := testutil.RandomSlug("tf-test-plat")
	manufacturerName := testutil.RandomName("tf-test-mfr-for-platform")
	manufacturerSlug := testutil.RandomSlug("tf-test-mfr-plat")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterPlatformCleanup(platformSlug)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckPlatformDestroy,
			testutil.CheckManufacturerDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccPlatformResourceConfig_basic(platformName, platformSlug, manufacturerName, manufacturerSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_platform.test", "id"),
					resource.TestCheckResourceAttr("netbox_platform.test", "name", platformName),
					resource.TestCheckResourceAttr("netbox_platform.test", "slug", platformSlug),
					resource.TestCheckResourceAttr("netbox_platform.test", "manufacturer", manufacturerSlug),
				),
			},
			{
				Config:   testAccPlatformResourceConfig_basic(platformName, platformSlug, manufacturerName, manufacturerSlug),
				PlanOnly: true,
			},
		},
	})
}

func TestAccPlatformResource_full(t *testing.T) {
	t.Parallel()

	platformName := testutil.RandomName("tf-test-platform-full")
	platformSlug := testutil.RandomSlug("tf-test-plat-full")
	manufacturerName := testutil.RandomName("tf-test-mfr-for-plat-full")
	manufacturerSlug := testutil.RandomSlug("tf-test-mfr-pf")
	description := testutil.RandomName("description")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterPlatformCleanup(platformSlug)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckPlatformDestroy,
			testutil.CheckManufacturerDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccPlatformResourceConfig_full(platformName, platformSlug, manufacturerName, manufacturerSlug, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_platform.test", "id"),
					resource.TestCheckResourceAttr("netbox_platform.test", "name", platformName),
					resource.TestCheckResourceAttr("netbox_platform.test", "slug", platformSlug),
					resource.TestCheckResourceAttr("netbox_platform.test", "manufacturer", manufacturerSlug),
					resource.TestCheckResourceAttr("netbox_platform.test", "description", description),
				),
			},
			{
				Config:   testAccPlatformResourceConfig_full(platformName, platformSlug, manufacturerName, manufacturerSlug, description),
				PlanOnly: true,
			},
		},
	})
}

func TestAccPlatformResource_update(t *testing.T) {
	t.Parallel()

	platformName := testutil.RandomName("tf-test-platform-update")
	platformSlug := testutil.RandomSlug("tf-test-plat-upd")
	manufacturerName := testutil.RandomName("tf-test-mfr-for-plat-upd")
	manufacturerSlug := testutil.RandomSlug("tf-test-mfr-pu")
	updatedName := testutil.RandomName("tf-test-platform-updated")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterPlatformCleanup(platformSlug)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckPlatformDestroy,
			testutil.CheckManufacturerDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccPlatformResourceConfig_basic(platformName, platformSlug, manufacturerName, manufacturerSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_platform.test", "id"),
					resource.TestCheckResourceAttr("netbox_platform.test", "name", platformName),
				),
			},
			{
				Config:   testAccPlatformResourceConfig_basic(platformName, platformSlug, manufacturerName, manufacturerSlug),
				PlanOnly: true,
			},
			{
				Config: testAccPlatformResourceConfig_basic(updatedName, platformSlug, manufacturerName, manufacturerSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_platform.test", "id"),
					resource.TestCheckResourceAttr("netbox_platform.test", "name", updatedName),
				),
			},
			{
				Config:   testAccPlatformResourceConfig_basic(updatedName, platformSlug, manufacturerName, manufacturerSlug),
				PlanOnly: true,
			},
		},
	})
}

func TestAccPlatformResource_import(t *testing.T) {
	t.Parallel()

	platformName := testutil.RandomName("tf-test-platform-import")
	platformSlug := testutil.RandomSlug("tf-test-plat-imp")
	manufacturerName := testutil.RandomName("tf-test-mfr-imp")
	manufacturerSlug := testutil.RandomSlug("tf-test-mfr-imp")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterPlatformCleanup(platformSlug)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckPlatformDestroy,
			testutil.CheckManufacturerDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccPlatformResourceConfig_import(platformName, platformSlug, manufacturerName, manufacturerSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_platform.test", "name", platformName),
					resource.TestCheckResourceAttr("netbox_platform.test", "slug", platformSlug),
				),
			},
			{
				ResourceName:            "netbox_platform.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"manufacturer", "display_name"},
				Check: resource.ComposeTestCheckFunc(
					testutil.ReferenceFieldCheck("netbox_platform.test", "manufacturer"),
				),
			},
			{
				Config:   testAccPlatformResourceConfig_import(platformName, platformSlug, manufacturerName, manufacturerSlug),
				PlanOnly: true,
			},
		},
	})
}

func testAccPlatformResourceConfig_basic(platformName, platformSlug, manufacturerName, manufacturerSlug string) string {
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

resource "netbox_manufacturer" "test_manufacturer" {
  name = %q
  slug = %q
}

resource "netbox_platform" "test" {
  name         = %q
  slug         = %q
  manufacturer = netbox_manufacturer.test_manufacturer.slug
}
`, manufacturerName, manufacturerSlug, platformName, platformSlug)
}

// TestAccPlatformResource_validationErrors tests validation error scenarios.
func TestAccPlatformResource_validationErrors(t *testing.T) {
	testutil.RunMultiValidationErrorTest(t, testutil.MultiValidationErrorTestConfig{
		ResourceName: "netbox_platform",
		TestCases: map[string]testutil.ValidationErrorCase{
			"missing_name": {
				Config: func() string {
					return `
resource "netbox_platform" "test" {
  slug = "test-platform"
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
			"missing_slug": {
				Config: func() string {
					return `
resource "netbox_platform" "test" {
  name = "Test Platform"
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
			"invalid_manufacturer_reference": {
				Config: func() string {
					return `
resource "netbox_platform" "test" {
  name         = "Test Platform"
  slug         = "test-platform"
  manufacturer = "99999"
}
`
				},
				ExpectedError: testutil.ErrPatternNotFound,
			},
		},
	})
}
func testAccPlatformResourceConfig_full(platformName, platformSlug, manufacturerName, manufacturerSlug, description string) string {
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

resource "netbox_manufacturer" "test_manufacturer" {
  name = %q
  slug = %q
}

resource "netbox_platform" "test" {
  name         = %q
  slug         = %q
  manufacturer = netbox_manufacturer.test_manufacturer.slug
  description  = %q
}
`, manufacturerName, manufacturerSlug, platformName, platformSlug, description)
}
func TestAccConsistency_Platform_LiteralNames(t *testing.T) {
	t.Parallel()

	platformName := testutil.RandomName("tf-test-platform-lit")
	platformSlug := testutil.RandomSlug("tf-test-plat-lit")
	manufacturerName := testutil.RandomName("tf-test-mfr-for-platform-lit")
	manufacturerSlug := testutil.RandomSlug("tf-test-mfr-plat-lit")
	description := testutil.RandomName("description")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterPlatformCleanup(platformSlug)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckPlatformDestroy,
			testutil.CheckManufacturerDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccPlatformConsistencyLiteralNamesConfig(platformName, platformSlug, manufacturerName, manufacturerSlug, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_platform.test", "id"),
					resource.TestCheckResourceAttr("netbox_platform.test", "name", platformName),
					resource.TestCheckResourceAttr("netbox_platform.test", "slug", platformSlug),
					resource.TestCheckResourceAttr("netbox_platform.test", "description", description),
				),
			},
			{
				Config:   testAccPlatformConsistencyLiteralNamesConfig(platformName, platformSlug, manufacturerName, manufacturerSlug, description),
				PlanOnly: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_platform.test", "id"),
				),
			},
		},
	})
}

func testAccPlatformConsistencyLiteralNamesConfig(platformName, platformSlug, manufacturerName, manufacturerSlug, description string) string {
	return fmt.Sprintf(`
resource "netbox_manufacturer" "test_manufacturer" {
  name = %q
  slug = %q
}

resource "netbox_platform" "test" {
  name         = %q
  slug         = %q
  manufacturer = netbox_manufacturer.test_manufacturer.slug
  description  = %q
}
`, manufacturerName, manufacturerSlug, platformName, platformSlug, description)
}
func testAccPlatformResourceConfig_import(platformName, platformSlug, manufacturerName, manufacturerSlug string) string {
	return fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
  name = %q
  slug = %q
}

resource "netbox_platform" "test" {
  name         = %q
  slug         = %q
  manufacturer = netbox_manufacturer.test.slug
}
`, manufacturerName, manufacturerSlug, platformName, platformSlug)
}

func TestAccPlatformResource_externalDeletion(t *testing.T) {
	t.Parallel()

	platformName := testutil.RandomName("test-platform-del")
	platformSlug := testutil.GenerateSlug(platformName)
	manufacturerName := testutil.RandomName("test-manufacturer")
	manufacturerSlug := testutil.GenerateSlug(manufacturerName)

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
	cleanup.RegisterPlatformCleanup(platformSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPlatformResourceConfig_basic(platformName, platformSlug, manufacturerName, manufacturerSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_platform.test", "id"),
					resource.TestCheckResourceAttr("netbox_platform.test", "name", platformName),
					resource.TestCheckResourceAttr("netbox_platform.test", "slug", platformSlug),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					items, _, err := client.DcimAPI.DcimPlatformsList(context.Background()).Slug([]string{platformSlug}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find platform for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.DcimAPI.DcimPlatformsDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete platform: %v", err)
					}
					t.Logf("Successfully externally deleted platform with ID: %d", itemID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

// TestAccPlatformResource_removeOptionalFields tests that optional nullable fields
// can be successfully removed from the configuration without causing inconsistent state.
// This verifies the bugfix for: "Provider produced inconsistent result after apply".
func TestAccPlatformResource_removeOptionalFields(t *testing.T) {
	t.Parallel()

	platformName := testutil.RandomName("tf-test-platform-rmv")
	platformSlug := testutil.RandomSlug("tf-test-plat-rmv")
	manufacturerName := testutil.RandomName("tf-test-mfr-rmv")
	manufacturerSlug := testutil.RandomSlug("tf-test-mfr-rmv")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterPlatformCleanup(platformSlug)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckPlatformDestroy,
			testutil.CheckManufacturerDestroy,
		),
		Steps: []resource.TestStep{
			// Step 1: Create platform with manufacturer
			{
				Config: testAccPlatformResourceConfig_withManufacturer(platformName, platformSlug, manufacturerName, manufacturerSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_platform.test", "id"),
					resource.TestCheckResourceAttr("netbox_platform.test", "name", platformName),
					resource.TestCheckResourceAttrSet("netbox_platform.test", "manufacturer"),
				),
			},
			// Step 2: Remove manufacturer - should set it to null
			{
				Config: testAccPlatformResourceConfig_withoutManufacturer(platformName, platformSlug, manufacturerName, manufacturerSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_platform.test", "id"),
					resource.TestCheckResourceAttr("netbox_platform.test", "name", platformName),
					resource.TestCheckNoResourceAttr("netbox_platform.test", "manufacturer"),
				),
			},
			// Step 3: Re-add manufacturer - verify it can be set again
			{
				Config: testAccPlatformResourceConfig_withManufacturer(platformName, platformSlug, manufacturerName, manufacturerSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_platform.test", "id"),
					resource.TestCheckResourceAttr("netbox_platform.test", "name", platformName),
					resource.TestCheckResourceAttrSet("netbox_platform.test", "manufacturer"),
				),
			},
		},
	})
}

func testAccPlatformResourceConfig_withManufacturer(platformName, platformSlug, manufacturerName, manufacturerSlug string) string {
	return fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
  name = %q
  slug = %q
}

resource "netbox_platform" "test" {
  name         = %q
  slug         = %q
  manufacturer = netbox_manufacturer.test.id
}
`, manufacturerName, manufacturerSlug, platformName, platformSlug)
}

func testAccPlatformResourceConfig_withoutManufacturer(platformName, platformSlug, manufacturerName, manufacturerSlug string) string {
	return fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
  name = %q
  slug = %q
}

resource "netbox_platform" "test" {
  name = %q
  slug = %q
}
`, manufacturerName, manufacturerSlug, platformName, platformSlug)
}
