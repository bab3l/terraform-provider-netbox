package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRegionDataSource_basic(t *testing.T) {
	t.Parallel()

	// Generate unique names
	name := testutil.RandomName("Public Cloud")
	slug := testutil.RandomSlug("tf-test-region-ds")

	// Register cleanup
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRegionCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckRegionDestroy,
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

func TestAccRegionDataSource_IDPreservation(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-region-ds-id")
	slug := testutil.RandomSlug("tf-test-region-ds-id")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRegionCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckRegionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRegionDataSourceConfig(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_region.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_region.test", "name", name),
				),
			},
		},
	})
}

func TestAccRegionDataSource_byID(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-region-ds")
	slug := testutil.RandomSlug("tf-test-region-ds")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRegionCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckRegionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRegionDataSourceConfigByID(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_region.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_region.test", "name", name),
					resource.TestCheckResourceAttr("data.netbox_region.test", "slug", slug),
				),
			},
		},
	})
}

func TestAccRegionDataSource_byName(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-region-ds")
	slug := testutil.RandomSlug("tf-test-region-ds")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRegionCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckRegionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRegionDataSourceConfigByName(name, slug),
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

func testAccRegionDataSourceConfigByID(name, slug string) string {
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
  id = netbox_region.test.id
}
`, name, slug)
}

func testAccRegionDataSourceConfigByName(name, slug string) string {
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
  name = netbox_region.test.name
}
`, name, slug)
}
