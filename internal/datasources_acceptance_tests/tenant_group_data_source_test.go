package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTenantGroupDataSource_basic(t *testing.T) {

	t.Parallel()

	// Generate unique names

	name := testutil.RandomName("tf-test-tg-ds")

	slug := testutil.RandomSlug("tf-test-tg-ds")

	// Register cleanup

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterTenantGroupCleanup(slug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

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

func TestAccTenantGroupDataSource_byID(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("tf-test-tg-ds")
	slug := testutil.RandomSlug("tf-test-tg-ds")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTenantGroupCleanup(slug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		CheckDestroy: testutil.CheckTenantGroupDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccTenantGroupDataSourceConfigByID(name, slug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("data.netbox_tenant_group.test", "id"),

					resource.TestCheckResourceAttr("data.netbox_tenant_group.test", "name", name),

					resource.TestCheckResourceAttr("data.netbox_tenant_group.test", "slug", slug),
				),
			},
		},
	})

}

func TestAccTenantGroupDataSource_byName(t *testing.T) {

	t.Parallel()

	name := fmt.Sprintf("Public Cloud %s", testutil.RandomName("tf-test-tg-ds"))
	slug := testutil.RandomSlug("tf-test-tg-ds")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTenantGroupCleanup(slug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		CheckDestroy: testutil.CheckTenantGroupDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccTenantGroupDataSourceConfigByName(name, slug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("data.netbox_tenant_group.test", "id"),

					resource.TestCheckResourceAttr("data.netbox_tenant_group.test", "name", name),

					resource.TestCheckResourceAttr("data.netbox_tenant_group.test", "slug", slug),
				),
			},
		},
	})

}

func TestAccTenantGroupDataSource_bySlug(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-tg-ds")
	slug := testutil.RandomSlug("tf-test-tg-ds")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTenantGroupCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckTenantGroupDestroy,
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

func TestAccTenantGroupDataSource_IDPreservation(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-tg-id")
	slug := testutil.RandomSlug("tf-test-tg-id")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTenantGroupCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckTenantGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTenantGroupDataSourceConfig(name, slug),
				Check: resource.ComposeTestCheckFunc(
					// Verify datasource returns ID correctly
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

func testAccTenantGroupDataSourceConfigByID(name, slug string) string {

	return fmt.Sprintf(`

resource "netbox_tenant_group" "test" {

  name = %q

  slug = %q

}

data "netbox_tenant_group" "test" {

  id = netbox_tenant_group.test.id

}

`, name, slug)

}

func testAccTenantGroupDataSourceConfigByName(name, slug string) string {

	return fmt.Sprintf(`

resource "netbox_tenant_group" "test" {

  name = %q

  slug = %q

}

data "netbox_tenant_group" "test" {

  name = netbox_tenant_group.test.name

}

`, name, slug)

}
