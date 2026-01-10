//go:build customfields
// +build customfields

package datasources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccWirelessLANGroupDataSource_customFields(t *testing.T) {
	customFieldName := testutil.RandomCustomFieldName("tf_test_wlangrp_ds_cf")
	groupName := testutil.RandomName("tf-test-wlangrp-ds-cf")
	groupSlug := testutil.RandomSlug("tf-test-wlangrp-ds-cf")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccWirelessLANGroupDataSourceConfig_customFields(customFieldName, groupName, groupSlug),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_wireless_lan_group.test", "name", groupName),
					resource.TestCheckResourceAttr("data.netbox_wireless_lan_group.test", "custom_fields.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_wireless_lan_group.test", "custom_fields.0.name", customFieldName),
					resource.TestCheckResourceAttr("data.netbox_wireless_lan_group.test", "custom_fields.0.type", "text"),
					resource.TestCheckResourceAttr("data.netbox_wireless_lan_group.test", "custom_fields.0.value", "datasource-test-value"),
				),
			},
		},
	})
}

func testAccWirelessLANGroupDataSourceConfig_customFields(customFieldName, groupName, groupSlug string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "test" {
  name         = %q
  object_types = ["wireless.wirelesslangroup"]
  type         = "text"
}

resource "netbox_wireless_lan_group" "test" {
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

data "netbox_wireless_lan_group" "test" {
  slug = netbox_wireless_lan_group.test.slug

  depends_on = [netbox_wireless_lan_group.test]
}
`, customFieldName, groupName, groupSlug)
}
