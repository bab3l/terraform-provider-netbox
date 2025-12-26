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

// Acceptance tests require NETBOX_URL and NETBOX_API_TOKEN environment variables.

func TestAccTunnelDataSource_IDPreservation(t *testing.T) {

	t.Parallel()

	testutil.TestAccPreCheck(t)

	randomName := testutil.RandomName("tf-test-tunnel-ds-id")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTunnelCleanup(randomName)

	resource.Test(t, resource.TestCase{

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckTunnelDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccTunnelDataSourceByID(randomName),

				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttrSet("data.netbox_tunnel.test", "id"),

					resource.TestCheckResourceAttr("data.netbox_tunnel.test", "name", randomName),
				),
			},
		},
	})

}

func TestAccTunnelDataSource_byID(t *testing.T) {

	t.Parallel()

	testutil.TestAccPreCheck(t)

	randomName := testutil.RandomName("tf-test-tunnel-ds")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTunnelCleanup(randomName)

	resource.Test(t, resource.TestCase{

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckTunnelDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccTunnelDataSourceByID(randomName),

				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttr("data.netbox_tunnel.test", "name", randomName),

					resource.TestCheckResourceAttrSet("data.netbox_tunnel.test", "id"),

					resource.TestCheckResourceAttr("data.netbox_tunnel.test", "status", "active"),

					resource.TestCheckResourceAttr("data.netbox_tunnel.test", "encapsulation", "gre"),
				),
			},
		},
	})

}

func TestAccTunnelDataSource_byName(t *testing.T) {

	t.Parallel()

	testutil.TestAccPreCheck(t)

	randomName := testutil.RandomName("tf-test-tunnel-ds")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTunnelCleanup(randomName)

	resource.Test(t, resource.TestCase{

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckTunnelDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccTunnelDataSourceByName(randomName),

				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttr("data.netbox_tunnel.test", "name", randomName),

					resource.TestCheckResourceAttrSet("data.netbox_tunnel.test", "id"),

					resource.TestCheckResourceAttr("data.netbox_tunnel.test", "status", "active"),
				),
			},
		},
	})

}

func testAccTunnelDataSourceByID(name string) string {

	return fmt.Sprintf(`

resource "netbox_tunnel" "test" {

  name          = %[1]q

  status        = "active"

  encapsulation = "gre"

}

data "netbox_tunnel" "test" {

  id = netbox_tunnel.test.id

}

`, name)

}

func testAccTunnelDataSourceByName(name string) string {

	return fmt.Sprintf(`

resource "netbox_tunnel" "test" {

  name          = %[1]q

  status        = "active"

  encapsulation = "gre"

}

data "netbox_tunnel" "test" {

  name = netbox_tunnel.test.name

}

`, name)

}
