//go:build customfields
// +build customfields

package resources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCircuitTypeResource_CustomFieldsPreservation(t *testing.T) {
	circuitTypeName := testutil.RandomName("tf-test-circuit-type")
	cfName := testutil.RandomCustomFieldName("tf_ct_pres")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Step 1: Create circuit type with custom field defined and populated
			{
				Config: testAccCircuitTypeResourcePreservationConfig_step1(circuitTypeName, cfName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_circuit_type.test", "name", circuitTypeName),
					resource.TestCheckResourceAttr("netbox_circuit_type.test", "custom_fields.#", "1"),
					testutil.CheckCustomFieldValue("netbox_circuit_type.test", cfName, "text", "preserved-value"),
				),
			},
			// Step 2: Update same circuit type without custom_fields in config (definition kept, preservation verified)
			{
				Config: testAccCircuitTypeResourcePreservationConfig_step2(circuitTypeName, cfName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_circuit_type.test", "name", circuitTypeName),
					// Custom fields omitted from config, so not in state (filtered-to-owned)
					resource.TestCheckResourceAttr("netbox_circuit_type.test", "custom_fields.#", "0"),
				),
			},
			// Step 3: Re-add custom_fields to verify preservation in NetBox
			{
				Config: testAccCircuitTypeResourcePreservationConfig_step1(circuitTypeName, cfName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_circuit_type.test", "name", circuitTypeName),
					resource.TestCheckResourceAttr("netbox_circuit_type.test", "custom_fields.#", "1"),
					testutil.CheckCustomFieldValue("netbox_circuit_type.test", cfName, "text", "preserved-value"),
				),
			},
		},
	})
}

func testAccCircuitTypeResourcePreservationConfig_step1(
	circuitTypeName, cfName string,
) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "ct_pres" {
  name = %[1]q
  object_types = ["circuits.circuittype"]
  type = "text"
}

resource "netbox_circuit_type" "test" {
  name = %[2]q
  slug = "tf-ct-pres-%[2]s"
  custom_fields = [
    {
      name = netbox_custom_field.ct_pres.name
      type = "text"
      value = "preserved-value"
    }
  ]

  depends_on = [netbox_custom_field.ct_pres]
}
`, cfName, circuitTypeName)
}

func testAccCircuitTypeResourcePreservationConfig_step2(
	circuitTypeName, cfName string,
) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "ct_pres" {
  name = %[1]q
  object_types = ["circuits.circuittype"]
  type = "text"
}

resource "netbox_circuit_type" "test" {
  name = %[2]q
  slug = "tf-ct-pres-%[2]s"
  # custom_fields intentionally omitted - values not managed by Terraform
  # but definition kept so field still exists in NetBox

  depends_on = [netbox_custom_field.ct_pres]
}
`, cfName, circuitTypeName)
}
