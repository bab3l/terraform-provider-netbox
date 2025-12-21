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

func TestAccSiteGroupDataSource_basic(t *testing.T) {

	// Generate unique names

	name := testutil.RandomName("tf-test-sg-ds")

	slug := testutil.RandomSlug("tf-test-sg-ds")

	// Register cleanup

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterSiteGroupCleanup(slug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckSiteGroupDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccSiteGroupDataSourceConfig(name, slug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("data.netbox_site_group.test", "id"),

					resource.TestCheckResourceAttr("data.netbox_site_group.test", "name", name),

					resource.TestCheckResourceAttr("data.netbox_site_group.test", "slug", slug),
				),
			},
		},
	})

}

func testAccSiteGroupDataSourceConfig(name, slug string) string {

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

resource "netbox_site_group" "test" {

  name = %q

  slug = %q

}

data "netbox_site_group" "test" {

  slug = netbox_site_group.test.slug

}

`, name, slug)

}
