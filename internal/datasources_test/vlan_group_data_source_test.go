package datasources_test

import (
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVLANGroupDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "netbox_vlan_group" "test" {
					name = "test-vlan-group"
					slug = "test-vlan-group"
				}

				data "netbox_vlan_group" "test" {
					name = netbox_vlan_group.test.name
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_vlan_group.test", "name", "test-vlan-group"),
					resource.TestCheckResourceAttr("data.netbox_vlan_group.test", "slug", "test-vlan-group"),
					resource.TestCheckResourceAttrSet("data.netbox_vlan_group.test", "id"),
				),
			},
		},
	})
}
