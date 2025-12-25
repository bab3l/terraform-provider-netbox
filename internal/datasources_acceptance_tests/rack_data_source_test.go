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

func TestAccRackDataSource_basic(t *testing.T) {

	t.Parallel()

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

func TestAccRackDataSource_byName(t *testing.T) {

	t.Parallel()

	siteName := testutil.RandomName("tf-test-rack-ds-site")

	siteSlug := testutil.RandomSlug("tf-test-rack-ds-s")

	rackName := testutil.RandomName("tf-test-rack-ds")

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

				Config: testAccRackDataSourceConfigByName(siteName, siteSlug, rackName),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("data.netbox_rack.test", "id"),

					resource.TestCheckResourceAttr("data.netbox_rack.test", "name", rackName),
				),
			},
		},
	})

}

func TestAccRackDataSource_IDPreservation(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-rack-ds-id")
	siteSlug := testutil.RandomSlug("tf-test-rack-ds-id")
	rackName := testutil.RandomName("tf-test-rack-ds-id")

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

func testAccRackDataSourceConfigByName(siteName, siteSlug, rackName string) string {

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
