package datasources_test

import (
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRoleDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRoleDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_role.test", "name", "Test Role"),
					resource.TestCheckResourceAttr("data.netbox_role.test", "slug", "test-role"),
				),
			},
		},
	})
}

const testAccRoleDataSourceConfig = `
resource "netbox_role" "test" {
  name = "Test Role"
  slug = "test-role"
}

data "netbox_role" "test" {
  id = netbox_role.test.id
}
`
