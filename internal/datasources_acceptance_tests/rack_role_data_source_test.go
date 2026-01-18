package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRackRoleDataSource_IDPreservation(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-rackrole-ds-id")
	slug := testutil.RandomSlug("tf-test-rr-ds-id")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRackRoleCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckRackRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRackRoleDataSourceConfig(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_rack_role.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_rack_role.test", "name", name),
				),
			},
		},
	})
}

func TestAccRackRoleDataSource_byID(t *testing.T) {
	t.Parallel()

	// Generate unique names
	name := testutil.RandomName("tf-test-rackrole-ds")
	slug := testutil.RandomSlug("tf-test-rr-ds")

	// Register cleanup
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRackRoleCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckRackRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRackRoleDataSourceConfig(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_rack_role.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_rack_role.test", "name", name),
					resource.TestCheckResourceAttr("data.netbox_rack_role.test", "slug", slug),
				),
			},
		},
	})
}

func TestAccRackRoleDataSource_byName(t *testing.T) {
	t.Parallel()

	name := fmt.Sprintf("Public Cloud %s", testutil.RandomName("tf-test-rackrole-ds"))
	slug := testutil.RandomSlug("tf-test-rr-ds")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRackRoleCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckRackRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRackRoleDataSourceConfigByName(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_rack_role.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_rack_role.test", "name", name),
					resource.TestCheckResourceAttr("data.netbox_rack_role.test", "slug", slug),
				),
			},
		},
	})
}

func TestAccRackRoleDataSource_bySlug(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-rackrole-ds")
	slug := testutil.RandomSlug("tf-test-rr-ds")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRackRoleCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckRackRoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRackRoleDataSourceConfigBySlug(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_rack_role.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_rack_role.test", "name", name),
					resource.TestCheckResourceAttr("data.netbox_rack_role.test", "slug", slug),
				),
			},
		},
	})
}

func testAccRackRoleDataSourceConfig(name, slug string) string {
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

resource "netbox_rack_role" "test" {
  name = %q
  slug = %q
}

data "netbox_rack_role" "test" {
  slug = netbox_rack_role.test.slug
}
`, name, slug)
}

func testAccRackRoleDataSourceConfigByName(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_rack_role" "test" {
	name = %q
	slug = %q
}

data "netbox_rack_role" "test" {
	name = netbox_rack_role.test.name
}
`, name, slug)
}

func testAccRackRoleDataSourceConfigBySlug(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_rack_role" "test" {
	name = %q
	slug = %q
}

data "netbox_rack_role" "test" {
	slug = netbox_rack_role.test.slug
}
`, name, slug)
}
