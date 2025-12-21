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

func TestAccRouteTargetDataSource_basic(t *testing.T) {

	t.Parallel()

	// Generate unique name - route targets have 21 char max, use format like "65000:400-<random>"

	name := fmt.Sprintf("65000:400-%s", testutil.RandomSlug("ds")[:8])

	// Register cleanup

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterRouteTargetCleanup(name)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckRouteTargetDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccRouteTargetDataSourceConfig(name),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("data.netbox_route_target.test", "id"),

					resource.TestCheckResourceAttr("data.netbox_route_target.test", "name", name),
				),
			},
		},
	})

}

func testAccRouteTargetDataSourceConfig(name string) string {

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

resource "netbox_route_target" "test" {

  name = %q

}

data "netbox_route_target" "test" {

  name = netbox_route_target.test.name

}

`, name)

}

// Virtual Disk Data Source Tests
