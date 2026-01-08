//go:build customfields
// +build customfields

package resources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIKEPolicyResource_CustomFieldsPreservation(t *testing.T) {
	ikePolicyName := testutil.RandomName("tf-test-ike-policy")
	cfName := testutil.RandomCustomFieldName("tf_ikepol_pres")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Step 1: Create IKE policy with custom field defined and populated
			{
				Config: testAccIKEPolicyResourcePreservationConfig_step1(ikePolicyName, cfName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_ike_policy.test", "name", ikePolicyName),
					resource.TestCheckResourceAttr("netbox_ike_policy.test", "custom_fields.#", "1"),
					testutil.CheckCustomFieldValue("netbox_ike_policy.test", cfName, "text", "preserved-value"),
				),
			},
			// Step 2: Update IKE policy without custom_fields in config (definition kept, preservation verified)
			{
				Config: testAccIKEPolicyResourcePreservationConfig_step2(ikePolicyName, cfName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_ike_policy.test", "name", ikePolicyName),
					// Custom fields omitted from config, so not in state (filtered-to-owned)
					resource.TestCheckResourceAttr("netbox_ike_policy.test", "custom_fields.#", "0"),
				),
			},
			// Step 3: Re-add custom_fields to verify preservation in NetBox
			{
				Config: testAccIKEPolicyResourcePreservationConfig_step1(ikePolicyName, cfName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_ike_policy.test", "name", ikePolicyName),
					resource.TestCheckResourceAttr("netbox_ike_policy.test", "custom_fields.#", "1"),
					testutil.CheckCustomFieldValue("netbox_ike_policy.test", cfName, "text", "preserved-value"),
				),
			},
		},
	})
}

func testAccIKEPolicyResourcePreservationConfig_step1(
	ikePolicyName, cfName string,
) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "ikepol_pres" {
  name = %[1]q
  object_types = ["vpn.ikepolicy"]
  type = "text"
}

resource "netbox_ike_policy" "test" {
  name = %[2]q
  version = 2
  custom_fields = [
    {
      name = netbox_custom_field.ikepol_pres.name
      type = "text"
      value = "preserved-value"
    }
  ]

  depends_on = [netbox_custom_field.ikepol_pres]
}
`, cfName, ikePolicyName)
}

func testAccIKEPolicyResourcePreservationConfig_step2(
	ikePolicyName, cfName string,
) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "ikepol_pres" {
  name = %[1]q
  object_types = ["vpn.ikepolicy"]
  type = "text"
}

resource "netbox_ike_policy" "test" {
  name = %[2]q
  version = 2
  # custom_fields intentionally omitted - values not managed by Terraform
  # but definition kept so field still exists in NetBox

  depends_on = [netbox_custom_field.ikepol_pres]
}
`, cfName, ikePolicyName)
}
