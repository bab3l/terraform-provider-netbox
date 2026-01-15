package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccInventoryItemRoleResource_basic(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-inv-role")
	slug := testutil.RandomSlug("role")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterInventoryItemRoleCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInventoryItemRoleResourceConfig_basic(name, slug),
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

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterInventoryItemRoleCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
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
	slug := testutil.RandomSlug("role")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterInventoryItemRoleCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInventoryItemRoleResourceConfig_basic(name, slug),
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
	slug := testutil.RandomSlug("role")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterInventoryItemRoleCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInventoryItemRoleResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_inventory_item_role.test", "id"),
				),
			},
			{
				ResourceName:      "netbox_inventory_item_role.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:   testAccInventoryItemRoleResourceConfig_basic(name, slug),
				PlanOnly: true,
			},
		},
	})
}

func TestAccInventoryItemRoleResource_IDPreservation(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-inv-role-id")
	slug := testutil.RandomSlug("role")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterInventoryItemRoleCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInventoryItemRoleResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_inventory_item_role.test", "id"),
					resource.TestCheckResourceAttr("netbox_inventory_item_role.test", "name", name),
				),
			},
		},
	})
}

func testAccInventoryItemRoleResourceConfig_basic(name, slug string) string {

	return fmt.Sprintf(`
resource "netbox_inventory_item_role" "test" {
  name = %q
  slug = %q
}
`, name, slug)
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

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterInventoryItemRoleCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInventoryItemRoleResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_inventory_item_role.test", "id"),
					resource.TestCheckResourceAttr("netbox_inventory_item_role.test", "name", name),
					resource.TestCheckResourceAttr("netbox_inventory_item_role.test", "slug", slug),
				),
			},
			{
				Config:   testAccInventoryItemRoleResourceConfig_basic(name, slug),
				PlanOnly: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_inventory_item_role.test", "id"),
				),
			},
		},
	})
}

func TestAccInventoryItemRoleResource_removeOptionalFields(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-inv-role-rem")
	slug := testutil.RandomSlug("role")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterInventoryItemRoleCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInventoryItemRoleResourceConfig_full(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_inventory_item_role.test", "name", name),
					resource.TestCheckResourceAttr("netbox_inventory_item_role.test", "description", "Test role description"),
				),
			},
			{
				Config: testAccInventoryItemRoleResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_inventory_item_role.test", "name", name),
					resource.TestCheckNoResourceAttr("netbox_inventory_item_role.test", "description"),
				),
			},
		},
	})
}

func TestAccInventoryItemRoleResource_externalDeletion(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-inv-role-ext-del")
	slug := testutil.RandomSlug("role")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterInventoryItemRoleCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInventoryItemRoleResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_inventory_item_role.test", "id"),
					resource.TestCheckResourceAttr("netbox_inventory_item_role.test", "name", name),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					items, _, err := client.DcimAPI.DcimInventoryItemRolesList(context.Background()).NameIc([]string{name}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find inventory_item_role for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.DcimAPI.DcimInventoryItemRolesDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete inventory_item_role: %v", err)
					}
					t.Logf("Successfully externally deleted inventory_item_role with ID: %d", itemID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}
