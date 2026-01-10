//go:build customfields
// +build customfields

package datasources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCircuitTerminationDataSource_customFields(t *testing.T) {
	customFieldName := testutil.RandomCustomFieldName("tf_test_circuit_term_ds_cf")
	siteName := testutil.RandomName("tf-test-site-ct-cf")
	providerName := testutil.RandomName("tf-test-provider-ct-cf")
	circuitCID := testutil.RandomName("tf-test-cid-ct-cf")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCircuitTerminationDataSourceConfig_customFields(customFieldName, siteName, providerName, circuitCID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_circuit_termination.test", "custom_fields.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_circuit_termination.test", "custom_fields.0.name", customFieldName),
					resource.TestCheckResourceAttr("data.netbox_circuit_termination.test", "custom_fields.0.type", "text"),
					resource.TestCheckResourceAttr("data.netbox_circuit_termination.test", "custom_fields.0.value", "test-circuit-term-value"),
				),
			},
		},
	})
}

func testAccCircuitTerminationDataSourceConfig_customFields(customFieldName, siteName, providerName, circuitCID string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "test" {
  name         = %[1]q
  object_types = ["circuits.circuittermination"]
  type         = "text"
}

resource "netbox_site" "test" {
  name = %[2]q
  slug = %[2]q
}

resource "netbox_provider" "test" {
  name = %[3]q
  slug = %[3]q
}

resource "netbox_circuit_type" "test" {
  name = "Test Circuit Type"
  slug = "test-circuit-type-ct-cf"
}

resource "netbox_circuit" "test" {
  cid              = %[4]q
  circuit_provider = netbox_provider.test.slug
  type             = netbox_circuit_type.test.slug
}

resource "netbox_circuit_termination" "test" {
  circuit   = netbox_circuit.test.cid
  site      = netbox_site.test.slug
  term_side = "A"

  custom_fields = [
    {
      name  = netbox_custom_field.test.name
      type  = "text"
      value = "test-circuit-term-value"
    }
  ]
}

data "netbox_circuit_termination" "test" {
  id = netbox_circuit_termination.test.id

  depends_on = [netbox_circuit_termination.test]
}
`, customFieldName, siteName, providerName, circuitCID)
}
