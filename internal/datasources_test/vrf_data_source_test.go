package datasources_test

import (
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVRFDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "netbox_vrf" "test" {
					name = "test-vrf"
				}

				data "netbox_vrf" "test" {
					name = netbox_vrf.test.name
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_vrf.test", "name", "test-vrf"),
					resource.TestCheckResourceAttrSet("data.netbox_vrf.test", "id"),
				),
			},
		},
	})
}
