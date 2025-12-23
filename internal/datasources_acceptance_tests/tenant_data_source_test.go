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

func TestAccTenantDataSource_basic(t *testing.T) {

	t.Parallel()

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

func TestAccTenantDataSource_byID(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("tf-test-tenant-ds")
	slug := testutil.RandomSlug("tf-test-tenant-ds")

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

				Config: testAccTenantDataSourceConfigByID(name, slug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("data.netbox_tenant.test", "id"),

					resource.TestCheckResourceAttr("data.netbox_tenant.test", "name", name),

					resource.TestCheckResourceAttr("data.netbox_tenant.test", "slug", slug),
				),
			},
		},
	})

}

func TestAccTenantDataSource_byName(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("tf-test-tenant-ds")
	slug := testutil.RandomSlug("tf-test-tenant-ds")

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

				Config: testAccTenantDataSourceConfigByName(name, slug),

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

func testAccTenantDataSourceConfigByID(name, slug string) string {

	return fmt.Sprintf(`

resource "netbox_tenant" "test" {

  name = %q

  slug = %q

}

data "netbox_tenant" "test" {

  id = netbox_tenant.test.id

}

`, name, slug)

}

func testAccTenantDataSourceConfigByName(name, slug string) string {

	return fmt.Sprintf(`

resource "netbox_tenant" "test" {

  name = %q

  slug = %q

}

data "netbox_tenant" "test" {

  name = netbox_tenant.test.name

}

`, name, slug)

}
