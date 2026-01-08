//go:build customfields
// +build customfields

package resources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIPSecProposalResource_CustomFieldsPreservation(t *testing.T) {
	ipsecProposalName := testutil.RandomName("tf-test-ipsec-proposal")
	cfName := testutil.RandomCustomFieldName("tf_ipsecprop_pres")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Step 1: Create IPSec proposal with custom field defined and populated
			{
				Config: testAccIPSecProposalResourcePreservationConfig_step1(ipsecProposalName, cfName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_ipsec_proposal.test", "name", ipsecProposalName),
					resource.TestCheckResourceAttr("netbox_ipsec_proposal.test", "custom_fields.#", "1"),
					testutil.CheckCustomFieldValue("netbox_ipsec_proposal.test", cfName, "text", "preserved-value"),
				),
			},
			// Step 2: Update IPSec proposal without custom_fields in config (definition kept, preservation verified)
			{
				Config: testAccIPSecProposalResourcePreservationConfig_step2(ipsecProposalName, cfName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_ipsec_proposal.test", "name", ipsecProposalName),
					// Custom fields omitted from config, so not in state (filtered-to-owned)
					resource.TestCheckResourceAttr("netbox_ipsec_proposal.test", "custom_fields.#", "0"),
				),
			},
			// Step 3: Re-add custom_fields to verify preservation in NetBox
			{
				Config: testAccIPSecProposalResourcePreservationConfig_step1(ipsecProposalName, cfName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_ipsec_proposal.test", "name", ipsecProposalName),
					resource.TestCheckResourceAttr("netbox_ipsec_proposal.test", "custom_fields.#", "1"),
					testutil.CheckCustomFieldValue("netbox_ipsec_proposal.test", cfName, "text", "preserved-value"),
				),
			},
		},
	})
}

func testAccIPSecProposalResourcePreservationConfig_step1(
	ipsecProposalName, cfName string,
) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "ipsecprop_pres" {
  name = %[1]q
  object_types = ["vpn.ipsecproposal"]
  type = "text"
}

resource "netbox_ipsec_proposal" "test" {
  name = %[2]q
  encryption_algorithm = "aes-128-cbc"
  authentication_algorithm = "hmac-sha1"
  custom_fields = [
    {
      name = netbox_custom_field.ipsecprop_pres.name
      type = "text"
      value = "preserved-value"
    }
  ]

  depends_on = [netbox_custom_field.ipsecprop_pres]
}
`, cfName, ipsecProposalName)
}

func testAccIPSecProposalResourcePreservationConfig_step2(
	ipsecProposalName, cfName string,
) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "ipsecprop_pres" {
  name = %[1]q
  object_types = ["vpn.ipsecproposal"]
  type = "text"
}

resource "netbox_ipsec_proposal" "test" {
  name = %[2]q
  encryption_algorithm = "aes-128-cbc"
  authentication_algorithm = "hmac-sha1"
  # custom_fields intentionally omitted - values not managed by Terraform
  # but definition kept so field still exists in NetBox

  depends_on = [netbox_custom_field.ipsecprop_pres]
}
`, cfName, ipsecProposalName)
}
