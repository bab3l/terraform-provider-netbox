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

func TestAccTunnelGroupDataSource_IDPreservation(t *testing.T) {

	t.Parallel()

	testutil.TestAccPreCheck(t)

	randomName := testutil.RandomName("tf-test-tunnel-grp-ds-id")

	randomSlug := testutil.RandomSlug("tf-test-tg-ds-id")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTunnelGroupCleanup(randomName)

	resource.Test(t, resource.TestCase{

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckTunnelGroupDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccTunnelGroupDataSourceByID(randomName, randomSlug),

				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttrSet("data.netbox_tunnel_group.test", "id"),

					resource.TestCheckResourceAttr("data.netbox_tunnel_group.test", "name", randomName),
				),
			},
		},
	})

}

func TestAccTunnelGroupDataSource_byID(t *testing.T) {

	t.Parallel()

	testutil.TestAccPreCheck(t)

	randomName := testutil.RandomName("tf-test-tunnel-grp-ds")

	randomSlug := testutil.RandomSlug("tf-test-tg-ds")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTunnelGroupCleanup(randomName)

	resource.Test(t, resource.TestCase{

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckTunnelGroupDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccTunnelGroupDataSourceByID(randomName, randomSlug),

				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttr("data.netbox_tunnel_group.test", "name", randomName),

					resource.TestCheckResourceAttr("data.netbox_tunnel_group.test", "slug", randomSlug),

					resource.TestCheckResourceAttrSet("data.netbox_tunnel_group.test", "id"),
				),
			},
		},
	})

}

func TestAccTunnelGroupDataSource_byName(t *testing.T) {

	t.Parallel()

	testutil.TestAccPreCheck(t)

	randomName := testutil.RandomName("tf-test-tunnel-grp-ds")

	randomSlug := testutil.RandomSlug("tf-test-tg-ds")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTunnelGroupCleanup(randomName)

	resource.Test(t, resource.TestCase{

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckTunnelGroupDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccTunnelGroupDataSourceByName(randomName, randomSlug),

				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttr("data.netbox_tunnel_group.test", "name", randomName),

					resource.TestCheckResourceAttr("data.netbox_tunnel_group.test", "slug", randomSlug),

					resource.TestCheckResourceAttrSet("data.netbox_tunnel_group.test", "id"),
				),
			},
		},
	})

}

func TestAccTunnelGroupDataSource_bySlug(t *testing.T) {

	t.Parallel()

	testutil.TestAccPreCheck(t)

	randomName := testutil.RandomName("tf-test-tunnel-grp-ds")

	randomSlug := testutil.RandomSlug("tf-test-tg-ds")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTunnelGroupCleanup(randomName)

	resource.Test(t, resource.TestCase{

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckTunnelGroupDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccTunnelGroupDataSourceBySlug(randomName, randomSlug),

				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttr("data.netbox_tunnel_group.test", "name", randomName),

					resource.TestCheckResourceAttr("data.netbox_tunnel_group.test", "slug", randomSlug),

					resource.TestCheckResourceAttrSet("data.netbox_tunnel_group.test", "id"),
				),
			},
		},
	})

}

func testAccTunnelGroupDataSourceByID(name, slug string) string {

	return fmt.Sprintf(`

resource "netbox_tunnel_group" "test" {

  name = %[1]q

  slug = %[2]q

}

data "netbox_tunnel_group" "test" {

  id = netbox_tunnel_group.test.id

}

`, name, slug)

}

func testAccTunnelGroupDataSourceByName(name, slug string) string {

	return fmt.Sprintf(`

resource "netbox_tunnel_group" "test" {

  name = %[1]q

  slug = %[2]q

}

data "netbox_tunnel_group" "test" {

  name = netbox_tunnel_group.test.name

}

`, name, slug)

}

func testAccTunnelGroupDataSourceBySlug(name, slug string) string {

	return fmt.Sprintf(`

resource "netbox_tunnel_group" "test" {

  name = %[1]q

  slug = %[2]q

}

data "netbox_tunnel_group" "test" {

  slug = netbox_tunnel_group.test.slug

}

`, name, slug)

}
