package datasources_acceptance_tests

import (
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccInventoryItemRoleDataSource_basic(t *testing.T) {

	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInventoryItemRoleDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_inventory_item_role.test", "name", "Test Inventory Item Role"),
					resource.TestCheckResourceAttr("data.netbox_inventory_item_role.test", "slug", "test-inventory-item-role"),
				),
			},
		},
	})
}

const testAccInventoryItemRoleDataSourceConfig = `
resource "netbox_inventory_item_role" "test" {
  name = "Test Inventory Item Role"
  slug = "test-inventory-item-role"
  color = "ff0000"
}

data "netbox_inventory_item_role" "test" {
  id = netbox_inventory_item_role.test.id
}
`
