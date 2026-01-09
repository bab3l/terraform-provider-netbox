//go:build customfields
// +build customfields

package datasources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCircuitDataSource_customFields(t *testing.T) {
	cid := testutil.RandomName("tf-test-circuit-ds-cf")
	providerName := testutil.RandomName("tf-test-provider-ds-cf")
	providerSlug := testutil.GenerateSlug(providerName)
	circuitTypeName := testutil.RandomName("tf-test-ct-ds-cf")
	circuitTypeSlug := testutil.GenerateSlug(circuitTypeName)
	customFieldName := testutil.RandomCustomFieldName("tf_test_circuit_ds_cf")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterProviderCleanup(providerSlug)
	cleanup.RegisterCircuitTypeCleanup(circuitTypeSlug)
	cleanup.RegisterCircuitCleanup(cid)
	cleanup.RegisterCustomFieldCleanup(customFieldName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Step 1: Create circuit with custom field and verify datasource returns it
			{
				Config: testAccCircuitDataSourceConfig_withCustomFields(cid, providerName, providerSlug, circuitTypeName, circuitTypeSlug, customFieldName),
				Check: resource.ComposeTestCheckFunc(
					// Verify datasource returns the custom field
					resource.TestCheckResourceAttr("data.netbox_circuit.test", "custom_fields.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_circuit.test", "custom_fields.0.name", customFieldName),
					resource.TestCheckResourceAttr("data.netbox_circuit.test", "custom_fields.0.type", "text"),
					resource.TestCheckResourceAttr("data.netbox_circuit.test", "custom_fields.0.value", "circuit-datasource-test"),
				),
			},
		},
	})
}

func testAccCircuitDataSourceConfig_withCustomFields(cid, providerName, providerSlug, circuitTypeName, circuitTypeSlug, customFieldName string) string {
	return fmt.Sprintf(`
resource "netbox_provider" "test" {
  name = %q
  slug = %q
}

resource "netbox_circuit_type" "test" {
  name = %q
  slug = %q
}

resource "netbox_custom_field" "test" {
  name         = %q
  object_types = ["circuits.circuit"]
  type         = "text"
}

resource "netbox_circuit" "test" {
  cid              = %q
  circuit_provider = netbox_provider.test.id
  type             = netbox_circuit_type.test.id

  custom_fields = [
    {
      name  = netbox_custom_field.test.name
      type  = "text"
      value = "circuit-datasource-test"
    }
  ]
}

data "netbox_circuit" "test" {
  cid = netbox_circuit.test.cid

  depends_on = [netbox_circuit.test]
}
`, providerName, providerSlug, circuitTypeName, circuitTypeSlug, customFieldName, cid)
}
