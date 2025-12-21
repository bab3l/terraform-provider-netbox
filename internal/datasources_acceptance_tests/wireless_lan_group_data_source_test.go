package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccWirelessLANGroupDataSource_basic(t *testing.T) {

	t.Parallel()

	cleanup := testutil.NewCleanupResource(t)

	name := testutil.RandomName("wlan-group")

	slug := testutil.RandomSlug("wlan-group")

	cleanup.RegisterWirelessLANGroupCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckWirelessLANGroupDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccWirelessLANGroupDataSourceConfig(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_wireless_lan_group.test", "name", name),
					resource.TestCheckResourceAttr("data.netbox_wireless_lan_group.test", "slug", slug),
					resource.TestCheckResourceAttrSet("data.netbox_wireless_lan_group.test", "id"),
				),
			},
		},
	})
}

func testAccWirelessLANGroupDataSourceConfig(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_wireless_lan_group" "test" {
	name = "%s"
	slug = "%s"
}

data "netbox_wireless_lan_group" "test" {
	name = netbox_wireless_lan_group.test.name
}
`, name, slug)
}
