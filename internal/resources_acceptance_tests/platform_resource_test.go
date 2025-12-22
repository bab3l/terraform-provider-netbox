package resources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/provider"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
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
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
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
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
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
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
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
				Config: testAccPlatformResourceConfig_basic(updatedName, platformSlug, manufacturerName, manufacturerSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_platform.test", "id"),
					resource.TestCheckResourceAttr("netbox_platform.test", "name", updatedName),
				),
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
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
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
				ImportStateVerifyIgnore: []string{"manufacturer"},
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
