//go:build customfields
// +build customfields

package datasources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIKEProposalDataSource_customFields(t *testing.T) {
	customFieldName := testutil.RandomCustomFieldName("tf_test_ikeproposal_ds_cf")
	ikeProposalName := testutil.RandomName("tf-test-ikeproposal-ds-cf")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIKEProposalDataSourceConfig_customFields(customFieldName, ikeProposalName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_ike_proposal.test", "name", ikeProposalName),
					resource.TestCheckResourceAttr("data.netbox_ike_proposal.test", "custom_fields.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_ike_proposal.test", "custom_fields.0.name", customFieldName),
					resource.TestCheckResourceAttr("data.netbox_ike_proposal.test", "custom_fields.0.type", "text"),
					resource.TestCheckResourceAttr("data.netbox_ike_proposal.test", "custom_fields.0.value", "datasource-test-value"),
				),
			},
		},
	})
}

func testAccIKEProposalDataSourceConfig_customFields(customFieldName, ikeProposalName string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "test" {
  name         = %q
  object_types = ["vpn.ikeproposal"]
  type         = "text"
}

resource "netbox_ike_proposal" "test" {
  name                    = %q
  authentication_method   = "preshared-keys"
  encryption_algorithm    = "aes-128-cbc"
  authentication_algorithm = "hmac-sha1"
  group                   = 14

  custom_fields = [
    {
      name  = netbox_custom_field.test.name
      type  = "text"
      value = "datasource-test-value"
    }
  ]
}

data "netbox_ike_proposal" "test" {
  name = %q

  depends_on = [netbox_ike_proposal.test]
}
`, customFieldName, ikeProposalName, ikeProposalName)
}
