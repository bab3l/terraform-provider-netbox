package datasources_test

import (
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVirtualChassisDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "netbox_virtual_chassis" "test" {
					name = "test-vc"
				}

				data "netbox_virtual_chassis" "test" {
					name = netbox_virtual_chassis.test.name
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_virtual_chassis.test", "name", "test-vc"),
					resource.TestCheckResourceAttrSet("data.netbox_virtual_chassis.test", "id"),
				),
			},
		},
	})
}
