//go:build customfields
// +build customfields

package resources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCircuitGroupResource_CustomFieldsPreservation(t *testing.T) {
	// NOTE: t.Parallel() intentionally omitted - this test creates/deletes global custom fields

	groupName := testutil.RandomName("tf-test-group-pres")
	cfName := testutil.RandomCustomFieldName("tf_group_pres")

	cleanup := testutil.NewCleanupResource(t)
	defer cleanup.Close(t)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCircuitGroupResourcePreservationConfig_step1(groupName, cfName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_circuit_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_circuit_group.test", "custom_fields.%", "1"),
					testutil.ResourceCheckCustomFieldValue("netbox_circuit_group.test", cfName, "preserved_value"),
				),
			},
			{
				// Update without custom_fields in config - should be preserved in NetBox
				Config: testAccCircuitGroupResourcePreservationConfig_step2(groupName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_circuit_group.test", "name", groupName+"_updated"),
					// Custom fields are not in the config, so they won't appear in state
				),
			},
		},
	})
}

func testAccCircuitGroupResourcePreservationConfig_step1(groupName, cfName string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "circuit_group_pres" {
  name         = %[2]q
  type         = "text"
  object_types = ["circuits.circuitgroup"]
  required     = false
}

resource "netbox_circuit_group" "test" {
  name = %[1]q

  custom_fields = {
    (netbox_custom_field.circuit_group_pres.name) = "preserved_value"
  }

  depends_on = [netbox_custom_field.circuit_group_pres]
}
`, groupName, cfName)
}

func testAccCircuitGroupResourcePreservationConfig_step2(groupName string) string {
	return fmt.Sprintf(`
resource "netbox_circuit_group" "test" {
  name = %[1]q
}
`, groupName+"_updated")
}
