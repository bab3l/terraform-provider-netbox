//go:build customfields
// +build customfields

package datasources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTunnelGroupDataSource_customFields(t *testing.T) {
	customFieldName := testutil.RandomCustomFieldName("tf_test_tunnelgroup_ds_cf")
	tunnelGroupName := testutil.RandomName("tf-test-tunnelgroup-ds-cf")
	tunnelGroupSlug := testutil.RandomSlug("tf-test-tunnelgroup-ds-cf")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTunnelGroupDataSourceConfig_customFields(customFieldName, tunnelGroupName, tunnelGroupSlug),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_tunnel_group.test", "name", tunnelGroupName),
					resource.TestCheckResourceAttr("data.netbox_tunnel_group.test", "custom_fields.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_tunnel_group.test", "custom_fields.0.name", customFieldName),
					resource.TestCheckResourceAttr("data.netbox_tunnel_group.test", "custom_fields.0.type", "text"),
					resource.TestCheckResourceAttr("data.netbox_tunnel_group.test", "custom_fields.0.value", "datasource-test-value"),
				),
			},
		},
	})
}

func testAccTunnelGroupDataSourceConfig_customFields(customFieldName, tunnelGroupName, tunnelGroupSlug string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "test" {
  name         = %q
  object_types = ["vpn.tunnelgroup"]
  type         = "text"
}

resource "netbox_tunnel_group" "test" {
  name = %q
  slug = %q

  custom_fields = [
    {
      name  = netbox_custom_field.test.name
      type  = "text"
      value = "datasource-test-value"
    }
  ]
}

data "netbox_tunnel_group" "test" {
  name = %q

  depends_on = [netbox_tunnel_group.test]
}
`, customFieldName, tunnelGroupName, tunnelGroupSlug, tunnelGroupName)
}
