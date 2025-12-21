package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccWirelessLANDataSource_basic(t *testing.T) {

	t.Parallel()

	ssid := testutil.RandomName("wlan")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccWirelessLANDataSourceConfig(ssid),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_wireless_lan.test", "ssid", ssid),
					resource.TestCheckResourceAttrSet("data.netbox_wireless_lan.test", "id"),
				),
			},
		},
	})
}

func testAccWirelessLANDataSourceConfig(ssid string) string {
	return fmt.Sprintf(`
resource "netbox_wireless_lan" "test" {
	ssid = "%s"
}

data "netbox_wireless_lan" "test" {
	ssid = netbox_wireless_lan.test.ssid
}
`, ssid)
}
