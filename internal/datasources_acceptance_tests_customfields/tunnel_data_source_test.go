//go:build customfields
// +build customfields

package datasources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTunnelDataSource_customFields(t *testing.T) {
	customFieldName := testutil.RandomCustomFieldName("tf_test_tunnel_ds_cf")
	tunnelName := testutil.RandomName("tf-test-tunnel-ds-cf")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTunnelDataSourceConfig_customFields(customFieldName, tunnelName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_tunnel.test", "name", tunnelName),
					resource.TestCheckResourceAttr("data.netbox_tunnel.test", "custom_fields.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_tunnel.test", "custom_fields.0.name", customFieldName),
					resource.TestCheckResourceAttr("data.netbox_tunnel.test", "custom_fields.0.type", "text"),
					resource.TestCheckResourceAttr("data.netbox_tunnel.test", "custom_fields.0.value", "datasource-test-value"),
				),
			},
		},
	})
}

func testAccTunnelDataSourceConfig_customFields(customFieldName, tunnelName string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "test" {
  name         = %q
  object_types = ["vpn.tunnel"]
  type         = "text"
}

resource "netbox_tunnel" "test" {
  name            = %q
  encapsulation   = "ipsec-transport"

  custom_fields = [
    {
      name  = netbox_custom_field.test.name
      type  = "text"
      value = "datasource-test-value"
    }
  ]
}

data "netbox_tunnel" "test" {
  name = %q

  depends_on = [netbox_tunnel.test]
}
`, customFieldName, tunnelName, tunnelName)
}
