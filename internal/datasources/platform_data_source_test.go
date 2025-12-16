package datasources_test

import (
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPlatformDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "netbox_platform" "test" {
				  name = "Test Platform DS"
				  slug = "test-platform-ds"
				}

				data "netbox_platform" "test" {
				  name = netbox_platform.test.name
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_platform.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_platform.test", "name", "Test Platform DS"),
					resource.TestCheckResourceAttr("data.netbox_platform.test", "slug", "test-platform-ds"),
				),
			},
		},
	})
}
