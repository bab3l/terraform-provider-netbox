package resources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/provider"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccProviderNetworkResource_basic(t *testing.T) {

	t.Parallel()

	providerName := testutil.RandomName("tf-test-provider")

	providerSlug := testutil.RandomSlug("tf-test-provider")

	networkName := testutil.RandomName("tf-test-network")

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterProviderCleanup(providerSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccProviderNetworkResourceConfig_basic(providerName, providerSlug, networkName),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_provider_network.test", "id"),

					resource.TestCheckResourceAttr("netbox_provider_network.test", "name", networkName),
				),
			},

			{

				ResourceName: "netbox_provider_network.test",

				ImportState: true,

				ImportStateVerify: true,
			},
		},
	})

}

func TestAccProviderNetworkResource_full(t *testing.T) {

	t.Parallel()

	providerName := testutil.RandomName("tf-test-provider-full")

	providerSlug := testutil.RandomSlug("tf-test-provider-full")

	networkName := testutil.RandomName("tf-test-network-full")

	serviceID := testutil.RandomName("svc")

	description := testutil.RandomName("description")

	updatedDescription := "Updated provider network description"

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterProviderCleanup(providerSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccProviderNetworkResourceConfig_full(providerName, providerSlug, networkName, serviceID, description),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_provider_network.test", "id"),

					resource.TestCheckResourceAttr("netbox_provider_network.test", "name", networkName),

					resource.TestCheckResourceAttr("netbox_provider_network.test", "service_id", serviceID),

					resource.TestCheckResourceAttr("netbox_provider_network.test", "description", description),
				),
			},

			{

				Config: testAccProviderNetworkResourceConfig_full(providerName, providerSlug, networkName, serviceID, updatedDescription),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_provider_network.test", "description", updatedDescription),
				),
			},
		},
	})

}

func testAccProviderNetworkResourceConfig_basic(providerName, providerSlug, networkName string) string {

	return fmt.Sprintf(`

resource "netbox_provider" "test" {

  name = %q

  slug = %q

}

resource "netbox_provider_network" "test" {

  circuit_provider = netbox_provider.test.id

  name             = %q

}

`, providerName, providerSlug, networkName)

}

func testAccProviderNetworkResourceConfig_full(providerName, providerSlug, networkName, serviceID, description string) string {

	return fmt.Sprintf(`

resource "netbox_provider" "test" {

  name = %q

  slug = %q

}

resource "netbox_provider_network" "test" {

  circuit_provider = netbox_provider.test.id

  name             = %q

  service_id       = %q

  description      = %q

}

`, providerName, providerSlug, networkName, serviceID, description)

}
