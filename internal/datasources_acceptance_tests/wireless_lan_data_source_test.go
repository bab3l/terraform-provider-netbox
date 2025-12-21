package datasources_acceptance_tests

import (
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccWirelessLANDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "netbox_wireless_lan" "test" {
					ssid = "test-ssid"
				}

				data "netbox_wireless_lan" "test" {
					ssid = netbox_wireless_lan.test.ssid
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_wireless_lan.test", "ssid", "test-ssid"),
					resource.TestCheckResourceAttrSet("data.netbox_wireless_lan.test", "id"),
				),
			},
		},
	})
}
