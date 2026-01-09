//go:build customfields
// +build customfields

package datasources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccProviderNetworkDataSource_customFields(t *testing.T) {
	customFieldName := testutil.RandomCustomFieldName("tf_test_providernetwork_ds_cf")
	providerName := testutil.RandomName("tf-test-provider-ds-cf")
	providerSlug := testutil.RandomSlug("tf-test-provider-ds-cf")
	networkName := testutil.RandomName("tf-test-network-ds-cf")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProviderNetworkDataSourceConfig_customFields(customFieldName, providerName, providerSlug, networkName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_provider_network.test", "name", networkName),
					resource.TestCheckResourceAttr("data.netbox_provider_network.test", "custom_fields.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_provider_network.test", "custom_fields.0.name", customFieldName),
					resource.TestCheckResourceAttr("data.netbox_provider_network.test", "custom_fields.0.type", "text"),
					resource.TestCheckResourceAttr("data.netbox_provider_network.test", "custom_fields.0.value", "datasource-test-value"),
				),
			},
		},
	})
}

func testAccProviderNetworkDataSourceConfig_customFields(customFieldName, providerName, providerSlug, networkName string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "test" {
  name         = %q
  object_types = ["circuits.providernetwork"]
  type         = "text"
}

resource "netbox_provider" "test" {
  name = %q
  slug = %q
}

resource "netbox_provider_network" "test" {
  circuit_provider = netbox_provider.test.id
  name        = %q

  custom_fields = [
    {
      name  = netbox_custom_field.test.name
      type  = "text"
      value = "datasource-test-value"
    }
  ]
}

data "netbox_provider_network" "test" {
  name = %q

  depends_on = [netbox_provider_network.test]
}
`, customFieldName, providerName, providerSlug, networkName, networkName)
}
