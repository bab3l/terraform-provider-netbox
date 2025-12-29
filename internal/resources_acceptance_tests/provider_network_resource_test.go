package resources_acceptance_tests

import (
	"context"
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

func TestAccProviderNetworkResource_IDPreservation(t *testing.T) {
	t.Parallel()

	providerName := testutil.RandomName("tf-test-provider-id")
	providerSlug := testutil.RandomSlug("tf-test-provider-id")
	networkName := testutil.RandomName("tf-test-network-id")

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
					resource.TestCheckResourceAttrSet("netbox_provider_network.test", "circuit_provider"),
				),
			},
		},
	})
}

func TestAccProviderNetworkResource_update(t *testing.T) {
	t.Parallel()

	providerName := testutil.RandomName("tf-test-provider-upd")
	providerSlug := testutil.RandomSlug("tf-test-provider-upd")
	networkName := testutil.RandomName("tf-test-network-upd")
	serviceID := "svc-12345"

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
					resource.TestCheckResourceAttr("netbox_provider_network.test", "name", networkName),
				),
			},
			{
				Config: testAccProviderNetworkResourceConfig_full(providerName, providerSlug, networkName, serviceID, testutil.Description2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_provider_network.test", "name", networkName),
					resource.TestCheckResourceAttr("netbox_provider_network.test", "service_id", serviceID),
					resource.TestCheckResourceAttr("netbox_provider_network.test", "description", testutil.Description2),
				),
			},
		},
	})
}

func TestAccProviderNetworkResource_import(t *testing.T) {
	t.Parallel()

	providerName := testutil.RandomName("tf-test-provider-imp")
	providerSlug := testutil.RandomSlug("tf-test-provider-imp")
	networkName := testutil.RandomName("tf-test-network-imp")

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
			},
			{
				ResourceName:      "netbox_provider_network.test",
				ImportState:       true,
				ImportStateVerify: true,
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

func TestAccConsistency_ProviderNetwork_LiteralNames(t *testing.T) {
	t.Parallel()

	providerName := testutil.RandomName("tf-test-provider-lit")
	providerSlug := testutil.RandomSlug("tf-test-provider-lit")
	networkName := testutil.RandomName("tf-test-network-lit")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterProviderCleanup(providerSlug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccProviderNetworkConsistencyLiteralNamesConfig(providerName, providerSlug, networkName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_provider_network.test", "id"),
					resource.TestCheckResourceAttr("netbox_provider_network.test", "name", networkName),
				),
			},
			{
				Config:   testAccProviderNetworkConsistencyLiteralNamesConfig(providerName, providerSlug, networkName),
				PlanOnly: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_provider_network.test", "id"),
				),
			},
		},
	})
}

func testAccProviderNetworkConsistencyLiteralNamesConfig(providerName, providerSlug, networkName string) string {
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

func TestAccProviderNetworkResource_externalDeletion(t *testing.T) {
	t.Parallel()
	providerName := testutil.RandomName("tf-test-provider-ext-del")
	providerSlug := testutil.RandomSlug("provider-ext-del")
	networkName := testutil.RandomName("tf-test-network-ext-del")
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "netbox_provider" "test" {
  name = %q
  slug = %q
}
resource "netbox_provider_network" "test" {
  circuit_provider = netbox_provider.test.id
  name             = %q
}
`, providerName, providerSlug, networkName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_provider_network.test", "id"),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					// List provider networks filtered by name
					items, _, err := client.CircuitsAPI.CircuitsProviderNetworksList(context.Background()).NameIc([]string{networkName}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find provider network for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.CircuitsAPI.CircuitsProviderNetworksDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete provider network: %v", err)
					}
					t.Logf("Successfully externally deleted provider network with ID: %d", itemID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}
