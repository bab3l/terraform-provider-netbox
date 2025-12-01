package datasources_test

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/provider"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSiteDataSource_basic(t *testing.T) {
	// Generate unique names to avoid conflicts between test runs
	name := testutil.RandomName("tf-test-site-ds")
	slug := testutil.RandomSlug("tf-test-site-ds")

	// Register cleanup to ensure resource is deleted even if test fails
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckSiteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSiteDataSourceConfig(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_site.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_site.test", "name", name),
					resource.TestCheckResourceAttr("data.netbox_site.test", "slug", slug),
				),
			},
		},
	})
}

func testAccSiteDataSourceConfig(name, slug string) string {
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

data "netbox_site" "test" {
  slug = netbox_site.test.slug
}
`, name, slug)
}

func TestAccTenantDataSource_basic(t *testing.T) {
	// Generate unique names
	name := testutil.RandomName("tf-test-tenant-ds")
	slug := testutil.RandomSlug("tf-test-tenant-ds")

	// Register cleanup
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTenantCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckTenantDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTenantDataSourceConfig(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_tenant.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_tenant.test", "name", name),
					resource.TestCheckResourceAttr("data.netbox_tenant.test", "slug", slug),
				),
			},
		},
	})
}

func testAccTenantDataSourceConfig(name, slug string) string {
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

resource "netbox_tenant" "test" {
  name = %q
  slug = %q
}

data "netbox_tenant" "test" {
  slug = netbox_tenant.test.slug
}
`, name, slug)
}

func TestAccSiteGroupDataSource_basic(t *testing.T) {
	// Generate unique names
	name := testutil.RandomName("tf-test-sg-ds")
	slug := testutil.RandomSlug("tf-test-sg-ds")

	// Register cleanup
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteGroupCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckSiteGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSiteGroupDataSourceConfig(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_site_group.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_site_group.test", "name", name),
					resource.TestCheckResourceAttr("data.netbox_site_group.test", "slug", slug),
				),
			},
		},
	})
}

func testAccSiteGroupDataSourceConfig(name, slug string) string {
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

resource "netbox_site_group" "test" {
  name = %q
  slug = %q
}

data "netbox_site_group" "test" {
  slug = netbox_site_group.test.slug
}
`, name, slug)
}

func TestAccTenantGroupDataSource_basic(t *testing.T) {
	// Generate unique names
	name := testutil.RandomName("tf-test-tg-ds")
	slug := testutil.RandomSlug("tf-test-tg-ds")

	// Register cleanup
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTenantGroupCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckTenantGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTenantGroupDataSourceConfig(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_tenant_group.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_tenant_group.test", "name", name),
					resource.TestCheckResourceAttr("data.netbox_tenant_group.test", "slug", slug),
				),
			},
		},
	})
}

func testAccTenantGroupDataSourceConfig(name, slug string) string {
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

resource "netbox_tenant_group" "test" {
  name = %q
  slug = %q
}

data "netbox_tenant_group" "test" {
  slug = netbox_tenant_group.test.slug
}
`, name, slug)
}

func TestAccManufacturerDataSource_basic(t *testing.T) {
	// Generate unique names
	name := testutil.RandomName("tf-test-mfr-ds")
	slug := testutil.RandomSlug("tf-test-mfr-ds")

	// Register cleanup
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckManufacturerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccManufacturerDataSourceConfig(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_manufacturer.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_manufacturer.test", "name", name),
					resource.TestCheckResourceAttr("data.netbox_manufacturer.test", "slug", slug),
				),
			},
		},
	})
}

func testAccManufacturerDataSourceConfig(name, slug string) string {
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

resource "netbox_manufacturer" "test" {
  name = %q
  slug = %q
}

data "netbox_manufacturer" "test" {
  slug = netbox_manufacturer.test.slug
}
`, name, slug)
}

func TestAccPlatformDataSource_basic(t *testing.T) {
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

func TestAccRegionDataSource_basic(t *testing.T) {
	// Generate unique names
	name := testutil.RandomName("tf-test-region-ds")
	slug := testutil.RandomSlug("tf-test-region-ds")

	// Register cleanup
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRegionCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckRegionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRegionDataSourceConfig(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_region.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_region.test", "name", name),
					resource.TestCheckResourceAttr("data.netbox_region.test", "slug", slug),
				),
			},
		},
	})
}

func testAccRegionDataSourceConfig(name, slug string) string {
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

resource "netbox_region" "test" {
  name = %q
  slug = %q
}

data "netbox_region" "test" {
  slug = netbox_region.test.slug
}
`, name, slug)
}

func TestAccLocationDataSource_basic(t *testing.T) {
	// Generate unique names
	siteName := testutil.RandomName("tf-test-loc-ds-site")
	siteSlug := testutil.RandomSlug("tf-test-loc-ds-s")
	name := testutil.RandomName("tf-test-location-ds")
	slug := testutil.RandomSlug("tf-test-location-ds")

	// Register cleanup
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterLocationCleanup(slug)
	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.ComposeCheckDestroy(testutil.CheckLocationDestroy, testutil.CheckSiteDestroy),
		Steps: []resource.TestStep{
			{
				Config: testAccLocationDataSourceConfig(siteName, siteSlug, name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_location.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_location.test", "name", name),
					resource.TestCheckResourceAttr("data.netbox_location.test", "slug", slug),
				),
			},
		},
	})
}

func testAccLocationDataSourceConfig(siteName, siteSlug, name, slug string) string {
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

data "netbox_location" "test" {
  slug = netbox_location.test.slug
}
`, siteName, siteSlug, name, slug)
}

func TestAccRackDataSource_basic(t *testing.T) {
	// Generate unique names
	siteName := testutil.RandomName("tf-test-rack-ds-site")
	siteSlug := testutil.RandomSlug("tf-test-rack-ds-s")
	rackName := testutil.RandomName("tf-test-rack-ds")

	// Register cleanup
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRackCleanup(rackName)
	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.ComposeCheckDestroy(testutil.CheckRackDestroy, testutil.CheckSiteDestroy),
		Steps: []resource.TestStep{
			{
				Config: testAccRackDataSourceConfig(siteName, siteSlug, rackName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_rack.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_rack.test", "name", rackName),
				),
			},
		},
	})
}

func testAccRackDataSourceConfig(siteName, siteSlug, rackName string) string {
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
  name = %q
  site = netbox_site.test.id
}

data "netbox_rack" "test" {
  name = netbox_rack.test.name
}
`, siteName, siteSlug, rackName)
}
