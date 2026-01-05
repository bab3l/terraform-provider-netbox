package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccProviderNetworkDataSource_IDPreservation(t *testing.T) {
	t.Parallel()

	providerName := testutil.RandomName("tf-test-provider-id")
	providerSlug := testutil.RandomSlug("tf-test-prov-id")
	networkName := testutil.RandomName("tf-test-network-id")
	serviceID := fmt.Sprintf("svc-%d", acctest.RandIntRange(10000, 99999))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProviderNetworkDataSourceConfig(providerName, providerSlug, networkName, serviceID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_provider_network.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_provider_network.test", "name", networkName),
				),
			},
		},
	})
}

func TestAccProviderNetworkDataSource_basic(t *testing.T) {
	t.Parallel()

	providerName := testutil.RandomName("tf-test-provider")
	providerSlug := testutil.RandomSlug("tf-test-prov")
	networkName := testutil.RandomName("tf-test-network")
	serviceID := fmt.Sprintf("svc-%d", acctest.RandIntRange(10000, 99999))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProviderNetworkDataSourceConfig(providerName, providerSlug, networkName, serviceID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_provider_network.test", "name", networkName),
					resource.TestCheckResourceAttr("data.netbox_provider_network.test", "service_id", serviceID),
				),
			},
		},
	})
}

func TestAccProviderNetworkDataSource_byName(t *testing.T) {
	t.Parallel()

	providerName := testutil.RandomName("tf-test-provider")
	providerSlug := testutil.RandomSlug("tf-test-prov")
	networkName := testutil.RandomName("tf-test-network")
	serviceID := fmt.Sprintf("svc-%d", acctest.RandIntRange(10000, 99999))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProviderNetworkDataSourceConfigByName(providerName, providerSlug, networkName, serviceID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_provider_network.test", "name", networkName),
					resource.TestCheckResourceAttr("data.netbox_provider_network.test", "service_id", serviceID),
				),
			},
		},
	})
}

func testAccProviderNetworkDataSourceConfig(providerName, providerSlug, networkName, serviceID string) string {
	return fmt.Sprintf(`
resource "netbox_provider" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_provider_network" "test" {
  circuit_provider = netbox_provider.test.id
  name             = %[3]q
  service_id       = %[4]q
}

data "netbox_provider_network" "test" {
  id = netbox_provider_network.test.id
}
`, providerName, providerSlug, networkName, serviceID)
}

func testAccProviderNetworkDataSourceConfigByName(providerName, providerSlug, networkName, serviceID string) string {
	return fmt.Sprintf(`
resource "netbox_provider" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_provider_network" "test" {
  circuit_provider = netbox_provider.test.id
  name             = %[3]q
  service_id       = %[4]q
}

data "netbox_provider_network" "test" {
  name = netbox_provider_network.test.name
}
`, providerName, providerSlug, networkName, serviceID)
}
