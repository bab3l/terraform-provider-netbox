//go:build customfields
// +build customfields

package datasources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIPSecProposalDataSource_customFields(t *testing.T) {
	customFieldName := testutil.RandomCustomFieldName("tf_test_ipsecproposal_ds_cf")
	ipsecProposalName := testutil.RandomName("tf-test-ipsecproposal-ds-cf")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIPSecProposalDataSourceConfig_customFields(customFieldName, ipsecProposalName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_ipsec_proposal.test", "name", ipsecProposalName),
					resource.TestCheckResourceAttr("data.netbox_ipsec_proposal.test", "custom_fields.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_ipsec_proposal.test", "custom_fields.0.name", customFieldName),
					resource.TestCheckResourceAttr("data.netbox_ipsec_proposal.test", "custom_fields.0.type", "text"),
					resource.TestCheckResourceAttr("data.netbox_ipsec_proposal.test", "custom_fields.0.value", "datasource-test-value"),
				),
			},
		},
	})
}

func testAccIPSecProposalDataSourceConfig_customFields(customFieldName, ipsecProposalName string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "test" {
  name         = %q
  object_types = ["vpn.ipsecproposal"]
  type         = "text"
}

resource "netbox_ipsec_proposal" "test" {
  name                    = %q
  encryption_algorithm    = "aes-128-cbc"
  authentication_algorithm = "hmac-sha1"

  custom_fields = [
    {
      name  = netbox_custom_field.test.name
      type  = "text"
      value = "datasource-test-value"
    }
  ]
}

data "netbox_ipsec_proposal" "test" {
  name = %q

  depends_on = [netbox_ipsec_proposal.test]
}
`, customFieldName, ipsecProposalName, ipsecProposalName)
}
