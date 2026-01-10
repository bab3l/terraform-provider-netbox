//go:build customfields
// +build customfields

package datasources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccWirelessLANDataSource_customFields(t *testing.T) {
	customFieldName := testutil.RandomCustomFieldName("tf_test_wlan_ds_cf")
	ssid := testutil.RandomName("tf-test-wlan-ds-cf")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccWirelessLANDataSourceConfig_customFields(customFieldName, ssid),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_wireless_lan.test", "ssid", ssid),
					resource.TestCheckResourceAttr("data.netbox_wireless_lan.test", "custom_fields.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_wireless_lan.test", "custom_fields.0.name", customFieldName),
					resource.TestCheckResourceAttr("data.netbox_wireless_lan.test", "custom_fields.0.type", "text"),
					resource.TestCheckResourceAttr("data.netbox_wireless_lan.test", "custom_fields.0.value", "datasource-test-value"),
				),
			},
		},
	})
}

func testAccWirelessLANDataSourceConfig_customFields(customFieldName, ssid string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "test" {
  name         = %q
  object_types = ["wireless.wirelesslan"]
  type         = "text"
}

resource "netbox_wireless_lan" "test" {
  ssid = %q

  custom_fields = [
    {
      name  = netbox_custom_field.test.name
      type  = "text"
      value = "datasource-test-value"
    }
  ]
}

data "netbox_wireless_lan" "test" {
  ssid = netbox_wireless_lan.test.ssid

  depends_on = [netbox_wireless_lan.test]
}
`, customFieldName, ssid)
}
