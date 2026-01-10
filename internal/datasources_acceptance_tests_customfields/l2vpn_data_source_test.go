//go:build customfields
// +build customfields

package datasources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccL2VPNDataSource_customFields(t *testing.T) {
	customFieldName := testutil.RandomCustomFieldName("tf_test_l2vpn_ds_cf")
	l2vpnName := testutil.RandomName("tf-test-l2vpn-ds-cf")
	l2vpnSlug := testutil.RandomSlug("tf-test-l2vpn-ds-cf")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccL2VPNDataSourceConfig_customFields(customFieldName, l2vpnName, l2vpnSlug),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_l2vpn.test", "name", l2vpnName),
					resource.TestCheckResourceAttr("data.netbox_l2vpn.test", "custom_fields.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_l2vpn.test", "custom_fields.0.name", customFieldName),
					resource.TestCheckResourceAttr("data.netbox_l2vpn.test", "custom_fields.0.type", "text"),
					resource.TestCheckResourceAttr("data.netbox_l2vpn.test", "custom_fields.0.value", "datasource-test-value"),
				),
			},
		},
	})
}

func testAccL2VPNDataSourceConfig_customFields(customFieldName, l2vpnName, l2vpnSlug string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "test" {
  name         = %q
  object_types = ["vpn.l2vpn"]
  type         = "text"
}

resource "netbox_l2vpn" "test" {
  name = %q
  slug = %q
  type = "vpws"

  custom_fields = [
    {
      name  = netbox_custom_field.test.name
      type  = "text"
      value = "datasource-test-value"
    }
  ]
}

data "netbox_l2vpn" "test" {
  name = %q

  depends_on = [netbox_l2vpn.test]
}
`, customFieldName, l2vpnName, l2vpnSlug, l2vpnName)
}
