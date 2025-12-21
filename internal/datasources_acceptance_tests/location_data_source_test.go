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

func TestAccLocationDataSource_basic(t *testing.T) {

	t.Parallel()

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
