//go:build customfields
// +build customfields

package resources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIPSecPolicyResource_CustomFieldsPreservation(t *testing.T) {
	ipsecPolicyName := testutil.RandomName("tf-test-ipsec-policy")
	cfName := testutil.RandomCustomFieldName("tf_ipsecpol_pres")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Step 1: Create IPSec policy with custom field defined and populated
			{
				Config: testAccIPSecPolicyResourcePreservationConfig_step1(ipsecPolicyName, cfName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_ipsec_policy.test", "name", ipsecPolicyName),
					resource.TestCheckResourceAttr("netbox_ipsec_policy.test", "custom_fields.#", "1"),
					testutil.CheckCustomFieldValue("netbox_ipsec_policy.test", cfName, "text", "preserved-value"),
				),
			},
			// Step 2: Update IPSec policy without custom_fields in config (definition kept, preservation verified)
			{
				Config: testAccIPSecPolicyResourcePreservationConfig_step2(ipsecPolicyName, cfName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_ipsec_policy.test", "name", ipsecPolicyName),
					// Custom fields omitted from config, so not in state (filtered-to-owned)
					resource.TestCheckResourceAttr("netbox_ipsec_policy.test", "custom_fields.#", "0"),
				),
			},
			// Step 3: Re-add custom_fields to verify preservation in NetBox
			{
				Config: testAccIPSecPolicyResourcePreservationConfig_step1(ipsecPolicyName, cfName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_ipsec_policy.test", "name", ipsecPolicyName),
					resource.TestCheckResourceAttr("netbox_ipsec_policy.test", "custom_fields.#", "1"),
					testutil.CheckCustomFieldValue("netbox_ipsec_policy.test", cfName, "text", "preserved-value"),
				),
			},
		},
	})
}

func testAccIPSecPolicyResourcePreservationConfig_step1(
	ipsecPolicyName, cfName string,
) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "ipsecpol_pres" {
  name = %[1]q
  object_types = ["vpn.ipsecpolicy"]
  type = "text"
}

resource "netbox_ipsec_policy" "test" {
  name = %[2]q
  custom_fields = [
    {
      name = netbox_custom_field.ipsecpol_pres.name
      type = "text"
      value = "preserved-value"
    }
  ]

  depends_on = [netbox_custom_field.ipsecpol_pres]
}
`, cfName, ipsecPolicyName)
}

func testAccIPSecPolicyResourcePreservationConfig_step2(
	ipsecPolicyName, cfName string,
) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "ipsecpol_pres" {
  name = %[1]q
  object_types = ["vpn.ipsecpolicy"]
  type = "text"
}

resource "netbox_ipsec_policy" "test" {
  name = %[2]q
  # custom_fields intentionally omitted - values not managed by Terraform
  # but definition kept so field still exists in NetBox

  depends_on = [netbox_custom_field.ipsecpol_pres]
}
`, cfName, ipsecPolicyName)
}
