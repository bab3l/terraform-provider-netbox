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

func TestAccASNRangeDataSource_basic(t *testing.T) {

	// Generate unique names

	name := testutil.RandomName("tf-test-asnrange-ds")

	slug := testutil.RandomSlug("tf-test-asnrange-ds")

	rirName := testutil.RandomName("tf-test-rir")

	rirSlug := testutil.RandomSlug("tf-test-rir")

	// Register cleanup

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterASNRangeCleanup(slug)

	cleanup.RegisterRIRCleanup(rirSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.ComposeCheckDestroy(

			testutil.CheckASNRangeDestroy,

			testutil.CheckRIRDestroy,
		),

		Steps: []resource.TestStep{

			{

				Config: testAccASNRangeDataSourceConfig(name, slug, rirName, rirSlug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("data.netbox_asn_range.test", "id"),

					resource.TestCheckResourceAttr("data.netbox_asn_range.test", "name", name),

					resource.TestCheckResourceAttr("data.netbox_asn_range.test", "slug", slug),

					resource.TestCheckResourceAttr("data.netbox_asn_range.test", "start", "64512"),

					resource.TestCheckResourceAttr("data.netbox_asn_range.test", "end", "64520"),
				),
			},
		},
	})

}

func testAccASNRangeDataSourceConfig(name, slug, rirName, rirSlug string) string {

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

resource "netbox_rir" "test" {

  name = %q

  slug = %q

}

resource "netbox_asn_range" "test" {

  name  = %q

  slug  = %q

  rir   = netbox_rir.test.id

  start = 64512

  end   = 64520

}

data "netbox_asn_range" "test" {

  slug = netbox_asn_range.test.slug

}

`, rirName, rirSlug, name, slug)

}

// Device Bay Template Data Source Tests
