package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccInventoryItemTemplateResource_basic(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-inv-template")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInventoryItemTemplateResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_inventory_item_template.test", "id"),
					resource.TestCheckResourceAttr("netbox_inventory_item_template.test", "name", name),
					resource.TestCheckResourceAttrSet("netbox_inventory_item_template.test", "device_type"),
				),
			},
		},
	})
}

func TestAccInventoryItemTemplateResource_IDPreservation(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("iit-id")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInventoryItemTemplateResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_inventory_item_template.test", "id"),
					resource.TestCheckResourceAttr("netbox_inventory_item_template.test", "name", name),
					resource.TestCheckResourceAttrSet("netbox_inventory_item_template.test", "device_type"),
				),
			},
		},
	})
}

func TestAccInventoryItemTemplateResource_full(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-inv-template-full")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInventoryItemTemplateResourceConfig_full(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_inventory_item_template.test", "id"),
					resource.TestCheckResourceAttr("netbox_inventory_item_template.test", "name", name),
					resource.TestCheckResourceAttr("netbox_inventory_item_template.test", "label", "Template Label"),
					resource.TestCheckResourceAttr("netbox_inventory_item_template.test", "part_id", "PART-001"),
					resource.TestCheckResourceAttr("netbox_inventory_item_template.test", "description", "Test template description"),
				),
			},
		},
	})
}

func TestAccInventoryItemTemplateResource_update(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-inv-template-update")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInventoryItemTemplateResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_inventory_item_template.test", "id"),
					resource.TestCheckResourceAttr("netbox_inventory_item_template.test", "name", name),
				),
			},
			{
				Config: testAccInventoryItemTemplateResourceConfig_full(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_inventory_item_template.test", "id"),
					resource.TestCheckResourceAttr("netbox_inventory_item_template.test", "label", "Template Label"),
					resource.TestCheckResourceAttr("netbox_inventory_item_template.test", "description", "Test template description"),
				),
			},
		},
	})
}

func TestAccInventoryItemTemplateResource_import(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-inv-template")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInventoryItemTemplateResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_inventory_item_template.test", "id"),
				),
			},
			{
				ResourceName:            "netbox_inventory_item_template.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"device_type"},
			},
		},
	})
}

func testAccInventoryItemTemplateResourceConfig_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "netbox_inventory_item_template" "test" {
  device_type = netbox_device_type.test.id
  name        = %q
}
`, testAccInventoryItemTemplateResourcePrereqs(name), name)
}

func testAccInventoryItemTemplateResourceConfig_full(name string) string {
	return fmt.Sprintf(`
%s

resource "netbox_inventory_item_template" "test" {
  device_type = netbox_device_type.test.id
  name        = %q
  label       = "Template Label"
  part_id     = "PART-001"
  description = "Test template description"
}
`, testAccInventoryItemTemplateResourcePrereqs(name), name)
}

func testAccInventoryItemTemplateResourcePrereqs(name string) string {
	return fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
  name = %q
  slug = %q
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = %q
  slug         = %q
}
`, name+"-mfr", testutil.RandomSlug("mfr"), name+"-model", testutil.RandomSlug("device"))
}

func TestAccConsistency_InventoryItemTemplate_LiteralNames(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-inv-template-lit")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInventoryItemTemplateConsistencyLiteralNamesConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_inventory_item_template.test", "id"),
					resource.TestCheckResourceAttr("netbox_inventory_item_template.test", "name", name),
				),
			},
			{
				Config:   testAccInventoryItemTemplateConsistencyLiteralNamesConfig(name),
				PlanOnly: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_inventory_item_template.test", "id"),
				),
			},
		},
	})
}

func TestAccInventoryItemTemplateResource_externalDeletion(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-inv-template-ext-del")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(name + "-mfr")
	cleanup.RegisterDeviceTypeCleanup(name + "-model")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInventoryItemTemplateResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_inventory_item_template.test", "id"),
					resource.TestCheckResourceAttr("netbox_inventory_item_template.test", "name", name),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					items, _, err := client.DcimAPI.DcimInventoryItemTemplatesList(context.Background()).NameIc([]string{name}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find inventory_item_template for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.DcimAPI.DcimInventoryItemTemplatesDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete inventory_item_template: %v", err)
					}
					t.Logf("Successfully externally deleted inventory_item_template with ID: %d", itemID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccInventoryItemTemplateConsistencyLiteralNamesConfig(name string) string {
	return fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
  name = "%s-mfr"
  slug = "%s-mfr"
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = "%s-model"
  slug         = "%s-model"
}

resource "netbox_inventory_item_template" "test" {
  device_type = netbox_device_type.test.id
  name        = %q
  part_id     = "PART-001"
}
`, name, name, name, name, name)
}
