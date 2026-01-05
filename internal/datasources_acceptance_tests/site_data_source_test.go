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

func TestAccSiteDataSource_basic(t *testing.T) {
	t.Parallel()

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

func TestAccSiteDataSource_byID(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-site-ds")
	slug := testutil.RandomSlug("tf-test-site-ds")

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
				Config: testAccSiteDataSourceConfigByID(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_site.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_site.test", "name", name),
					resource.TestCheckResourceAttr("data.netbox_site.test", "slug", slug),
				),
			},
		},
	})
}

func TestAccSiteDataSource_byName(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-site-ds")
	slug := testutil.RandomSlug("tf-test-site-ds")

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
				Config: testAccSiteDataSourceConfigByName(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_site.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_site.test", "name", name),
					resource.TestCheckResourceAttr("data.netbox_site.test", "slug", slug),
				),
			},
		},
	})
}

func TestAccSiteDataSource_IDPreservation(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-site-id")
	slug := testutil.RandomSlug("tf-test-site-id")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckSiteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSiteDataSourceConfig(name, slug),
				Check: resource.ComposeTestCheckFunc(
					// Verify datasource returns ID correctly
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

func testAccSiteDataSourceConfigByID(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name   = %q
  slug   = %q
  status = "active"
}

data "netbox_site" "test" {
  id = netbox_site.test.id
}
`, name, slug)
}

func testAccSiteDataSourceConfigByName(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name   = %q
  slug   = %q
  status = "active"
}

data "netbox_site" "test" {
  name = netbox_site.test.name
}
`, name, slug)
}
