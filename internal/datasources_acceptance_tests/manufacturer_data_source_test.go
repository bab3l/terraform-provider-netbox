package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccManufacturerDataSource_IDPreservation(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-mfr-ds-id")
	slug := testutil.RandomSlug("tf-test-mfr-ds-id")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckManufacturerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccManufacturerDataSourceConfig(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_manufacturer.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_manufacturer.test", "name", name),
				),
			},
		},
	})
}

func TestAccManufacturerDataSource_basic(t *testing.T) {
	t.Parallel()

	// Generate unique names
	name := testutil.RandomName("tf-test-mfr-ds")
	slug := testutil.RandomSlug("tf-test-mfr-ds")

	// Register cleanup
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckManufacturerDestroy,
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

func TestAccManufacturerDataSource_byName(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-mfr-ds")
	slug := testutil.RandomSlug("tf-test-mfr-ds")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckManufacturerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccManufacturerDataSourceConfigByName(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_manufacturer.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_manufacturer.test", "name", name),
					resource.TestCheckResourceAttr("data.netbox_manufacturer.test", "slug", slug),
				),
			},
		},
	})
}

func testAccManufacturerDataSourceConfigByName(name, slug string) string {
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
  name = netbox_manufacturer.test.name
}
`, name, slug)
}

func TestAccManufacturerDataSource_byID(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-mfr-ds")
	slug := testutil.RandomSlug("tf-test-mfr-ds")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckManufacturerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccManufacturerDataSourceConfigByID(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_manufacturer.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_manufacturer.test", "name", name),
					resource.TestCheckResourceAttr("data.netbox_manufacturer.test", "slug", slug),
				),
			},
		},
	})
}

func testAccManufacturerDataSourceConfigByID(name, slug string) string {
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
  id = netbox_manufacturer.test.id
}
`, name, slug)
}
