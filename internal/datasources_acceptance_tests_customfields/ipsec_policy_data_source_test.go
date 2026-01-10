//go:build customfields
// +build customfields

package datasources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIPSecPolicyDataSource_customFields(t *testing.T) {
	customFieldName := testutil.RandomCustomFieldName("tf_test_ipsecpolicy_ds_cf")
	ipsecPolicyName := testutil.RandomName("tf-test-ipsecpolicy-ds-cf")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIPSecPolicyDataSourceConfig_customFields(customFieldName, ipsecPolicyName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_ipsec_policy.test", "name", ipsecPolicyName),
					resource.TestCheckResourceAttr("data.netbox_ipsec_policy.test", "custom_fields.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_ipsec_policy.test", "custom_fields.0.name", customFieldName),
					resource.TestCheckResourceAttr("data.netbox_ipsec_policy.test", "custom_fields.0.type", "text"),
					resource.TestCheckResourceAttr("data.netbox_ipsec_policy.test", "custom_fields.0.value", "datasource-test-value"),
				),
			},
		},
	})
}

func testAccIPSecPolicyDataSourceConfig_customFields(customFieldName, ipsecPolicyName string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "test" {
  name         = %q
  object_types = ["vpn.ipsecpolicy"]
  type         = "text"
}

resource "netbox_ipsec_policy" "test" {
  name = %q

  custom_fields = [
    {
      name  = netbox_custom_field.test.name
      type  = "text"
      value = "datasource-test-value"
    }
  ]
}

data "netbox_ipsec_policy" "test" {
  name = %q

  depends_on = [netbox_ipsec_policy.test]
}
`, customFieldName, ipsecPolicyName, ipsecPolicyName)
}
