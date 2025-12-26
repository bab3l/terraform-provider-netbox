package resources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccInventoryItemRoleResource_basic(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("tf-test-inv-role")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccInventoryItemRoleResourceConfig_basic(name),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_inventory_item_role.test", "id"),

					resource.TestCheckResourceAttr("netbox_inventory_item_role.test", "name", name),

					resource.TestCheckResourceAttrSet("netbox_inventory_item_role.test", "slug"),
				),
			},
		},
	})

}

func TestAccInventoryItemRoleResource_full(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("tf-test-inv-role-full")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccInventoryItemRoleResourceConfig_full(name),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_inventory_item_role.test", "id"),

					resource.TestCheckResourceAttr("netbox_inventory_item_role.test", "name", name),

					resource.TestCheckResourceAttr("netbox_inventory_item_role.test", "color", "e41e22"),

					resource.TestCheckResourceAttr("netbox_inventory_item_role.test", "description", "Test role description"),
				),
			},
		},
	})

}

func TestAccInventoryItemRoleResource_update(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("tf-test-inv-role-update")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccInventoryItemRoleResourceConfig_basic(name),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_inventory_item_role.test", "id"),

					resource.TestCheckResourceAttr("netbox_inventory_item_role.test", "name", name),
				),
			},

			{

				Config: testAccInventoryItemRoleResourceConfig_full(name),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_inventory_item_role.test", "id"),

					resource.TestCheckResourceAttr("netbox_inventory_item_role.test", "color", "e41e22"),

					resource.TestCheckResourceAttr("netbox_inventory_item_role.test", "description", "Test role description"),
				),
			},
		},
	})

}

func TestAccInventoryItemRoleResource_import(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("tf-test-inv-role")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccInventoryItemRoleResourceConfig_basic(name),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_inventory_item_role.test", "id"),
				),
			},

			{

				ResourceName: "netbox_inventory_item_role.test",

				ImportState: true,

				ImportStateVerify: true,
			},
		},
	})

}

func TestAccInventoryItemRoleResource_IDPreservation(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-inv-role-id")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInventoryItemRoleResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_inventory_item_role.test", "id"),
					resource.TestCheckResourceAttr("netbox_inventory_item_role.test", "name", name),
				),
			},
		},
	})

}

func testAccInventoryItemRoleResourceConfig_basic(name string) string {

	return fmt.Sprintf(`

resource "netbox_inventory_item_role" "test" {

  name = %q

  slug = %q

}

`, name, testutil.RandomSlug("role"))

}

func testAccInventoryItemRoleResourceConfig_full(name string) string {

	return fmt.Sprintf(`

resource "netbox_inventory_item_role" "test" {

  name        = %q

  slug        = %q

  color       = "e41e22"

  description = "Test role description"

}

`, name, testutil.RandomSlug("role"))

}

func TestAccConsistency_InventoryItemRole_LiteralNames(t *testing.T) {
	t.Parallel()
	name := testutil.RandomName("tf-test-inv-role-lit")
	slug := testutil.RandomSlug("tf-test-inv-role-lit")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInventoryItemRoleConsistencyLiteralNamesConfig(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_inventory_item_role.test", "id"),
					resource.TestCheckResourceAttr("netbox_inventory_item_role.test", "name", name),
					resource.TestCheckResourceAttr("netbox_inventory_item_role.test", "slug", slug),
				),
			},
			{
				Config:   testAccInventoryItemRoleConsistencyLiteralNamesConfig(name, slug),
				PlanOnly: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_inventory_item_role.test", "id"),
				),
			},
		},
	})
}

func testAccInventoryItemRoleConsistencyLiteralNamesConfig(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_inventory_item_role" "test" {
  name = %q
  slug = %q
}
`, name, slug)
}
