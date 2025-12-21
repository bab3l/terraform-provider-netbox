package datasources_acceptance_tests

import (
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVLANDataSource_basic(t *testing.T) {

	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "netbox_vlan" "test" {
					name = "test-vlan"
					vid  = 100
				}

				data "netbox_vlan" "test" {
					name = netbox_vlan.test.name
					vid  = netbox_vlan.test.vid
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_vlan.test", "name", "test-vlan"),
					resource.TestCheckResourceAttr("data.netbox_vlan.test", "vid", "100"),
					resource.TestCheckResourceAttrSet("data.netbox_vlan.test", "id"),
				),
			},
		},
	})
}
