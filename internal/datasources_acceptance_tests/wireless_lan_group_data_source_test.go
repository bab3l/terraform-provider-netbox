package datasources_acceptance_tests

import (
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccWirelessLANGroupDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "netbox_wireless_lan_group" "test" {
					name = "test-wlan-group"
					slug = "test-wlan-group"
				}

				data "netbox_wireless_lan_group" "test" {
					name = netbox_wireless_lan_group.test.name
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_wireless_lan_group.test", "name", "test-wlan-group"),
					resource.TestCheckResourceAttr("data.netbox_wireless_lan_group.test", "slug", "test-wlan-group"),
					resource.TestCheckResourceAttrSet("data.netbox_wireless_lan_group.test", "id"),
				),
			},
		},
	})
}
