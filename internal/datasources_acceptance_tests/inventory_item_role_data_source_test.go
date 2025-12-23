package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccInventoryItemRoleDataSource_basic(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("tf-test-inv-item-role")
	slug := testutil.RandomSlug("tf-test-inv-item-role")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInventoryItemRoleDataSourceConfig(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_inventory_item_role.test", "name", name),
					resource.TestCheckResourceAttr("data.netbox_inventory_item_role.test", "slug", slug),
				),
			},
		},
	})
}

func testAccInventoryItemRoleDataSourceConfig(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_inventory_item_role" "test" {
  name = %[1]q
  slug = %[2]q
  color = "ff0000"
}

data "netbox_inventory_item_role" "test" {
  id = netbox_inventory_item_role.test.id
}
`, name, slug)
}

func TestAccInventoryItemRoleDataSource_byName(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("tf-test-inv-item-role")
	slug := testutil.RandomSlug("tf-test-inv-item-role")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInventoryItemRoleDataSourceConfigByName(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_inventory_item_role.test", "name", name),
					resource.TestCheckResourceAttr("data.netbox_inventory_item_role.test", "slug", slug),
				),
			},
		},
	})
}

func testAccInventoryItemRoleDataSourceConfigByName(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_inventory_item_role" "test" {
  name = %[1]q
  slug = %[2]q
  color = "ff0000"
}

data "netbox_inventory_item_role" "test" {
  name = netbox_inventory_item_role.test.name
}
`, name, slug)
}

func TestAccInventoryItemRoleDataSource_bySlug(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("tf-test-inv-item-role")
	slug := testutil.RandomSlug("tf-test-inv-item-role")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInventoryItemRoleDataSourceConfigBySlug(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_inventory_item_role.test", "name", name),
					resource.TestCheckResourceAttr("data.netbox_inventory_item_role.test", "slug", slug),
				),
			},
		},
	})
}

func testAccInventoryItemRoleDataSourceConfigBySlug(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_inventory_item_role" "test" {
  name = %[1]q
  slug = %[2]q
  color = "ff0000"
}

data "netbox_inventory_item_role" "test" {
  slug = netbox_inventory_item_role.test.slug
}
`, name, slug)
}
