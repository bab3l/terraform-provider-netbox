package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccInventoryItemResource_basic(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-inv-item")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInventoryItemResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_inventory_item.test", "id"),
					resource.TestCheckResourceAttr("netbox_inventory_item.test", "name", name),
					resource.TestCheckResourceAttrSet("netbox_inventory_item.test", "device"),
				),
			},
		},
	})
}

func TestAccInventoryItemResource_IDPreservation(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("inv-id")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInventoryItemResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_inventory_item.test", "id"),
					resource.TestCheckResourceAttr("netbox_inventory_item.test", "name", name),
					resource.TestCheckResourceAttrSet("netbox_inventory_item.test", "device"),
				),
			},
		},
	})
}

func TestAccInventoryItemResource_full(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-inv-item-full")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInventoryItemResourceConfig_full(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_inventory_item.test", "id"),
					resource.TestCheckResourceAttr("netbox_inventory_item.test", "name", name),
					resource.TestCheckResourceAttr("netbox_inventory_item.test", "label", "Inventory Label"),
					resource.TestCheckResourceAttr("netbox_inventory_item.test", "serial", "SN-12345"),
					resource.TestCheckResourceAttr("netbox_inventory_item.test", "asset_tag", name+"-asset-tag"),
					resource.TestCheckResourceAttr("netbox_inventory_item.test", "description", "Test inventory item"),
				),
			},
		},
	})
}

func TestAccInventoryItemResource_update(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-inv-item-update")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInventoryItemResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_inventory_item.test", "id"),
					resource.TestCheckResourceAttr("netbox_inventory_item.test", "name", name),
				),
			},
			{
				Config: testAccInventoryItemResourceConfig_full(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_inventory_item.test", "id"),
					resource.TestCheckResourceAttr("netbox_inventory_item.test", "label", "Inventory Label"),
					resource.TestCheckResourceAttr("netbox_inventory_item.test", "serial", "SN-12345"),
				),
			},
		},
	})
}

func TestAccInventoryItemResource_import(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-inv-item")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInventoryItemResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_inventory_item.test", "id"),
				),
			},
			{
				ResourceName:            "netbox_inventory_item.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"device"},
			},
		},
	})
}

func testAccInventoryItemResourceConfig_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "netbox_inventory_item" "test" {
  device = netbox_device.test.id
  name   = %q
}
`, testAccInventoryItemResourcePrereqs(name), name)
}

func testAccInventoryItemResourceConfig_full(name string) string {
	return fmt.Sprintf(`

%s

resource "netbox_inventory_item" "test" {
  device      = netbox_device.test.id
  name        = %q
  label       = "Inventory Label"
  serial      = "SN-12345"
  asset_tag   = %q
  description = "Test inventory item"
}
`, testAccInventoryItemResourcePrereqs(name), name, name+"-asset-tag")
}

func testAccInventoryItemResourcePrereqs(name string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name = %q
  slug = %q
}

resource "netbox_manufacturer" "test" {
  name = %q
  slug = %q
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = %q
  slug         = %q
}

resource "netbox_device_role" "test" {
  name = %q
  slug = %q
}

resource "netbox_device" "test" {
  site        = netbox_site.test.id
  name        = %q
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  status      = "offline"
}
`, name+"-site", testutil.RandomSlug("site"), name+"-mfr", testutil.RandomSlug("mfr"), name+"-model", testutil.RandomSlug("device"), name+"-role", testutil.RandomSlug("role"), name+"-device")
}

func TestAccConsistency_InventoryItem_LiteralNames(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-inv-item-lit")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInventoryItemConsistencyLiteralNamesConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_inventory_item.test", "id"),
					resource.TestCheckResourceAttr("netbox_inventory_item.test", "name", name),
				),
			},
			{
				Config:   testAccInventoryItemConsistencyLiteralNamesConfig(name),
				PlanOnly: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_inventory_item.test", "id"),
				),
			},
		},
	})
}

func TestAccInventoryItemResource_externalDeletion(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-inv-item-ext-del")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInventoryItemResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_inventory_item.test", "id"),
					resource.TestCheckResourceAttr("netbox_inventory_item.test", "name", name),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					items, _, err := client.DcimAPI.DcimInventoryItemsList(context.Background()).NameIc([]string{name}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find inventory_item for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.DcimAPI.DcimInventoryItemsDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete inventory_item: %v", err)
					}
					t.Logf("Successfully externally deleted inventory_item with ID: %d", itemID)
				},
				Config: testAccInventoryItemResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_inventory_item.test", "id"),
				),
			},
		},
	})
}

func testAccInventoryItemConsistencyLiteralNamesConfig(name string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name = "%s-site"
  slug = "%s-site"
}

resource "netbox_manufacturer" "test" {
  name = "%s-mfr"
  slug = "%s-mfr"
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = "%s-model"
  slug         = "%s-model"
}

resource "netbox_device_role" "test" {
  name = "%s-role"
  slug = "%s-role"
}

resource "netbox_device" "test" {
  site        = netbox_site.test.id
  name        = "%s-device"
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  status      = "offline"
}

resource "netbox_inventory_item" "test" {
  device = netbox_device.test.id
  name   = %q
}
`, name, name, name, name, name, name, name, name, name, name)
}

// NOTE: Custom field tests for inventory_item resource are in resources_acceptance_tests_customfields package
