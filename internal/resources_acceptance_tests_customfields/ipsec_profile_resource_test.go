//go:build customfields
// +build customfields

package resources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIPSecProfileResource_CustomFieldsPreservation(t *testing.T) {
	ipsecProfileName := testutil.RandomName("tf-test-ipsec-profile")
	cfName := testutil.RandomCustomFieldName("tf_ipsecprof_pres")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Step 1: Create IPSec profile with custom field defined and populated
			{
				Config: testAccIPSecProfileResourcePreservationConfig_step1(ipsecProfileName, cfName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_ipsec_profile.test", "name", ipsecProfileName),
					resource.TestCheckResourceAttr("netbox_ipsec_profile.test", "custom_fields.#", "1"),
					testutil.CheckCustomFieldValue("netbox_ipsec_profile.test", cfName, "text", "preserved-value"),
				),
			},
			// Step 2: Update IPSec profile without custom_fields in config (definition kept, preservation verified)
			{
				Config: testAccIPSecProfileResourcePreservationConfig_step2(ipsecProfileName, cfName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_ipsec_profile.test", "name", ipsecProfileName),
					// Custom fields omitted from config, so not in state (filtered-to-owned)
					resource.TestCheckResourceAttr("netbox_ipsec_profile.test", "custom_fields.#", "0"),
				),
			},
			// Step 3: Re-add custom_fields to verify preservation in NetBox
			{
				Config: testAccIPSecProfileResourcePreservationConfig_step1(ipsecProfileName, cfName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_ipsec_profile.test", "name", ipsecProfileName),
					resource.TestCheckResourceAttr("netbox_ipsec_profile.test", "custom_fields.#", "1"),
					testutil.CheckCustomFieldValue("netbox_ipsec_profile.test", cfName, "text", "preserved-value"),
				),
			},
		},
	})
}

func testAccIPSecProfileResourcePreservationConfig_step1(
	ipsecProfileName, cfName string,
) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "ipsecprof_pres" {
  name = %[1]q
  object_types = ["vpn.ipsecprofile"]
  type = "text"
}

resource "netbox_ike_policy" "test" {
  name = "test-ike-policy"
  version = 1
  mode = "main"
}

resource "netbox_ipsec_policy" "test" {
  name = "test-ipsec-policy"
}

resource "netbox_ipsec_profile" "test" {
  name = %[2]q
  mode = "esp"
  ike_policy = netbox_ike_policy.test.id
  ipsec_policy = netbox_ipsec_policy.test.id
  custom_fields = [
    {
      name = netbox_custom_field.ipsecprof_pres.name
      type = "text"
      value = "preserved-value"
    }
  ]

  depends_on = [netbox_custom_field.ipsecprof_pres]
}
`, cfName, ipsecProfileName)
}

func testAccIPSecProfileResourcePreservationConfig_step2(
	ipsecProfileName, cfName string,
) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "ipsecprof_pres" {
  name = %[1]q
  object_types = ["vpn.ipsecprofile"]
  type = "text"
}

resource "netbox_ike_policy" "test" {
  name = "test-ike-policy"
  version = 1
  mode = "main"
}

resource "netbox_ipsec_policy" "test" {
  name = "test-ipsec-policy"
}

resource "netbox_ipsec_profile" "test" {
  name = %[2]q
  mode = "esp"
  ike_policy = netbox_ike_policy.test.id
  ipsec_policy = netbox_ipsec_policy.test.id
  # custom_fields intentionally omitted - values not managed by Terraform
  # but definition kept so field still exists in NetBox

  depends_on = [netbox_custom_field.ipsecprof_pres]
}
`, cfName, ipsecProfileName)
}
