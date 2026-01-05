package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/provider"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPlatformDataSource_basic(t *testing.T) {
	t.Parallel()

	// Generate unique names for both manufacturer and platform
	// Platform requires a manufacturer, so we create both
	mfrName := testutil.RandomName("tf-test-mfr-for-plat-ds")
	mfrSlug := testutil.RandomSlug("tf-test-mfr-pds")
	platName := testutil.RandomName("tf-test-plat-ds")
	platSlug := testutil.RandomSlug("tf-test-plat-ds")

	// Register cleanup for both resources
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterPlatformCleanup(platSlug)
	cleanup.RegisterManufacturerCleanup(mfrSlug)

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
				Config: testAccPlatformDataSourceConfig(platName, platSlug, mfrName, mfrSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_platform.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_platform.test", "name", platName),
					resource.TestCheckResourceAttr("data.netbox_platform.test", "slug", platSlug),
				),
			},
		},
	})
}

func TestAccPlatformDataSource_byName(t *testing.T) {
	t.Parallel()

	mfrName := testutil.RandomName("tf-test-mfr-for-plat-ds")
	mfrSlug := testutil.RandomSlug("tf-test-mfr-pds")
	platName := testutil.RandomName("tf-test-plat-ds")
	platSlug := testutil.RandomSlug("tf-test-plat-ds")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterPlatformCleanup(platSlug)
	cleanup.RegisterManufacturerCleanup(mfrSlug)

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
				Config: testAccPlatformDataSourceConfigByName(platName, platSlug, mfrName, mfrSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_platform.by_name", "id"),
					resource.TestCheckResourceAttr("data.netbox_platform.by_name", "name", platName),
				),
			},
		},
	})
}

func testAccPlatformDataSourceConfig(platName, platSlug, mfrName, mfrSlug string) string {
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

resource "netbox_manufacturer" "test_mfr" {
  name = %q
  slug = %q
}

resource "netbox_platform" "test" {
  name         = %q
  slug         = %q
  manufacturer = netbox_manufacturer.test_mfr.slug
}

data "netbox_platform" "test" {
  slug = netbox_platform.test.slug
}
`, mfrName, mfrSlug, platName, platSlug)
}

func TestAccPlatformDataSource_IDPreservation(t *testing.T) {
	t.Parallel()

	platformName := testutil.RandomName("platform-ds-id")
	platformSlug := testutil.GenerateSlug(platformName)
	manufacturerName := testutil.RandomName("mfr-ds-id")
	manufacturerSlug := testutil.GenerateSlug(manufacturerName)

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterPlatformCleanup(platformName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckPlatformDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccPlatformDataSourceConfig(platformName, platformSlug, manufacturerName, manufacturerSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_platform.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_platform.test", "name", platformName),
				),
			},
		},
	})
}

func testAccPlatformDataSourceConfigByName(platName, platSlug, mfrName, mfrSlug string) string {
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

resource "netbox_manufacturer" "test_mfr" {
  name = %q
  slug = %q
}

resource "netbox_platform" "test" {
  name         = %q
  slug         = %q
  manufacturer = netbox_manufacturer.test_mfr.slug
}

data "netbox_platform" "by_name" {
  name = netbox_platform.test.name
}
`, mfrName, mfrSlug, platName, platSlug)
}
