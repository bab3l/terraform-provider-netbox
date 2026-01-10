//go:build customfields
// +build customfields

package datasources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCircuitTypeDataSource_customFields(t *testing.T) {
	circuitTypeName := testutil.RandomName("tf-test-ct-ds-cf")
	circuitTypeSlug := testutil.GenerateSlug(circuitTypeName)
	customFieldName := testutil.RandomCustomFieldName("tf_test_ct_ds_cf")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterCircuitTypeCleanup(circuitTypeSlug)
	cleanup.RegisterCustomFieldCleanup(customFieldName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Step 1: Create circuit type with custom field and verify datasource returns it
			{
				Config: testAccCircuitTypeDataSourceConfig_withCustomFields(circuitTypeName, circuitTypeSlug, customFieldName),
				Check: resource.ComposeTestCheckFunc(
					// Verify datasource returns the custom field
					resource.TestCheckResourceAttr("data.netbox_circuit_type.test", "custom_fields.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_circuit_type.test", "custom_fields.0.name", customFieldName),
					resource.TestCheckResourceAttr("data.netbox_circuit_type.test", "custom_fields.0.type", "text"),
					resource.TestCheckResourceAttr("data.netbox_circuit_type.test", "custom_fields.0.value", "circuit-type-datasource-test"),
				),
			},
		},
	})
}

func testAccCircuitTypeDataSourceConfig_withCustomFields(circuitTypeName, circuitTypeSlug, customFieldName string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "test" {
  name         = %q
  object_types = ["circuits.circuittype"]
  type         = "text"
}

resource "netbox_circuit_type" "test" {
  name = %q
  slug = %q

  custom_fields = [
    {
      name  = netbox_custom_field.test.name
      type  = "text"
      value = "circuit-type-datasource-test"
    }
  ]
}

data "netbox_circuit_type" "test" {
  slug = netbox_circuit_type.test.slug

  depends_on = [netbox_circuit_type.test]
}
`, customFieldName, circuitTypeName, circuitTypeSlug)
}
