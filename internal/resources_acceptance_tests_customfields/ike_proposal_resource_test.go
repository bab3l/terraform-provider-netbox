//go:build customfields
// +build customfields

package resources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIKEProposalResource_CustomFieldsPreservation(t *testing.T) {
	ikeProposalName := testutil.RandomName("tf-test-ike-proposal")
	cfName := testutil.RandomCustomFieldName("tf_ikeprop_pres")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Step 1: Create IKE proposal with custom field defined and populated
			{
				Config: testAccIKEProposalResourcePreservationConfig_step1(ikeProposalName, cfName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_ike_proposal.test", "name", ikeProposalName),
					resource.TestCheckResourceAttr("netbox_ike_proposal.test", "custom_fields.#", "1"),
					testutil.CheckCustomFieldValue("netbox_ike_proposal.test", cfName, "text", "preserved-value"),
				),
			},
			// Step 2: Update IKE proposal without custom_fields in config (definition kept, preservation verified)
			{
				Config: testAccIKEProposalResourcePreservationConfig_step2(ikeProposalName, cfName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_ike_proposal.test", "name", ikeProposalName),
					// Custom fields omitted from config, so not in state (filtered-to-owned)
					resource.TestCheckResourceAttr("netbox_ike_proposal.test", "custom_fields.#", "0"),
				),
			},
			// Step 3: Re-add custom_fields to verify preservation in NetBox
			{
				Config: testAccIKEProposalResourcePreservationConfig_step1(ikeProposalName, cfName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_ike_proposal.test", "name", ikeProposalName),
					resource.TestCheckResourceAttr("netbox_ike_proposal.test", "custom_fields.#", "1"),
					testutil.CheckCustomFieldValue("netbox_ike_proposal.test", cfName, "text", "preserved-value"),
				),
			},
		},
	})
}

func testAccIKEProposalResourcePreservationConfig_step1(
	ikeProposalName, cfName string,
) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "ikeprop_pres" {
  name = %[1]q
  object_types = ["vpn.ikeproposal"]
  type = "text"
}

resource "netbox_ike_proposal" "test" {
  name = %[2]q
  authentication_method = "preshared-keys"
  encryption_algorithm = "aes-128-cbc"
  group = 14
  custom_fields = [
    {
      name = netbox_custom_field.ikeprop_pres.name
      type = "text"
      value = "preserved-value"
    }
  ]

  depends_on = [netbox_custom_field.ikeprop_pres]
}
`, cfName, ikeProposalName)
}

func testAccIKEProposalResourcePreservationConfig_step2(
	ikeProposalName, cfName string,
) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "ikeprop_pres" {
  name = %[1]q
  object_types = ["vpn.ikeproposal"]
  type = "text"
}

resource "netbox_ike_proposal" "test" {
  name = %[2]q
  authentication_method = "preshared-keys"
  encryption_algorithm = "aes-128-cbc"
  group = 14
  # custom_fields intentionally omitted - values not managed by Terraform
  # but definition kept so field still exists in NetBox

  depends_on = [netbox_custom_field.ikeprop_pres]
}
`, cfName, ikeProposalName)
}
