//go:build customfields
// +build customfields

package datasources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIKEPolicyDataSource_customFields(t *testing.T) {
	customFieldName := testutil.RandomCustomFieldName("tf_test_ikepolicy_ds_cf")
	ikePolicyName := testutil.RandomName("tf-test-ikepolicy-ds-cf")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIKEPolicyDataSourceConfig_customFields(customFieldName, ikePolicyName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_ike_policy.test", "name", ikePolicyName),
					resource.TestCheckResourceAttr("data.netbox_ike_policy.test", "custom_fields.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_ike_policy.test", "custom_fields.0.name", customFieldName),
					resource.TestCheckResourceAttr("data.netbox_ike_policy.test", "custom_fields.0.type", "text"),
					resource.TestCheckResourceAttr("data.netbox_ike_policy.test", "custom_fields.0.value", "datasource-test-value"),
				),
			},
		},
	})
}

func testAccIKEPolicyDataSourceConfig_customFields(customFieldName, ikePolicyName string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "test" {
  name         = %q
  object_types = ["vpn.ikepolicy"]
  type         = "text"
}

resource "netbox_ike_policy" "test" {
  name    = %q
  version = 2

  custom_fields = [
    {
      name  = netbox_custom_field.test.name
      type  = "text"
      value = "datasource-test-value"
    }
  ]
}

data "netbox_ike_policy" "test" {
  name = %q

  depends_on = [netbox_ike_policy.test]
}
`, customFieldName, ikePolicyName, ikePolicyName)
}
