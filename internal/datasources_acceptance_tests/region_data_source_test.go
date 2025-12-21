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

func TestAccRegionDataSource_basic(t *testing.T) {

	t.Parallel()

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
