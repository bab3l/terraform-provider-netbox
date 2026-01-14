package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var _ = testAccInventoryItemTemplateResourceConfig_removeOptionalFields_base

func TestAccInventoryItemTemplateResource_basic(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-inv-template")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(testutil.RandomSlug("mfr"))
	cleanup.RegisterDeviceTypeCleanup(testutil.RandomSlug("device"))

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

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(testutil.RandomSlug("mfr"))
	cleanup.RegisterDeviceTypeCleanup(testutil.RandomSlug("device"))

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

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(testutil.RandomSlug("mfr"))
	cleanup.RegisterDeviceTypeCleanup(testutil.RandomSlug("device"))

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

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(testutil.RandomSlug("mfr"))
	cleanup.RegisterDeviceTypeCleanup(testutil.RandomSlug("device"))

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

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(testutil.RandomSlug("mfr"))
	cleanup.RegisterDeviceTypeCleanup(testutil.RandomSlug("device"))

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

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(name + "-mfr")
	cleanup.RegisterDeviceTypeCleanup(name + "-model")

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

// TestAccInventoryItemTemplateResource_removeOptionalFields tests that optional fields
// can be successfully removed from the configuration without causing inconsistent state.
// This verifies the bugfix for: "Provider produced inconsistent result after apply".
func TestAccInventoryItemTemplateResource_removeOptionalFields(t *testing.T) {
	t.Parallel()

	templateName := testutil.RandomName("tf-test-invtempl-rem")
	manufacturerName := testutil.RandomName("tf-test-mfr-rem")
	manufacturerSlug := testutil.RandomSlug("tf-test-mfr-rem")
	roleName := testutil.RandomName("tf-test-role-rem")
	roleSlug := testutil.RandomSlug("tf-test-role-rem")

	cleanup := testutil.NewCleanupResource(t)
	// We only strictly track the extra resources here; prereqs are handled by their own implicit cleanup or best-effort
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
	cleanup.RegisterInventoryItemRoleCleanup(roleSlug)
	// Note: RegisterInventoryItemTemplateCleanup might not be available or needed if we rely on CheckDestroy of the TestStep (default behavior destroys at end)

	resourceName := "netbox_inventory_item_template.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInventoryItemTemplateResourceConfig_removeOptionalFields_withFields(
					templateName, manufacturerName, manufacturerSlug, roleName, roleSlug,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", templateName),
					resource.TestCheckResourceAttr(resourceName, "label", "Test Label"),
					resource.TestCheckResourceAttr(resourceName, "description", "Test Description"),
					// Verify references are set to the IDs of the extra resources
					resource.TestCheckResourceAttrPair(resourceName, "role", "netbox_inventory_item_role.extra", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "manufacturer", "netbox_manufacturer.extra", "id"),
				),
			},
			{
				Config: testAccInventoryItemTemplateResourceConfig_removeOptionalFields_detached(
					templateName, manufacturerName, manufacturerSlug, roleName, roleSlug,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", templateName),
					// Verify optional fields are cleared
					resource.TestCheckNoResourceAttr(resourceName, "label"),
					resource.TestCheckNoResourceAttr(resourceName, "description"),
					resource.TestCheckNoResourceAttr(resourceName, "role"),
					resource.TestCheckNoResourceAttr(resourceName, "manufacturer"),
				),
			},
		},
	})
}

func testAccInventoryItemTemplateResourceConfig_removeOptionalFields_base(name string) string {
	return testAccInventoryItemTemplateResourcePrereqs(name) + fmt.Sprintf(`
resource "netbox_inventory_item_template" "test" {
  device_type = netbox_device_type.test.id
  name        = %q
}
`, name)
}

func testAccInventoryItemTemplateResourceConfig_removeOptionalFields_withFields(
	name, manufacturerName, manufacturerSlug, roleName, roleSlug string) string {
	return testAccInventoryItemTemplateResourcePrereqs(name) + fmt.Sprintf(`
resource "netbox_manufacturer" "extra" {
  name = %[2]q
  slug = %[3]q
}

resource "netbox_inventory_item_role" "extra" {
  name  = %[4]q
  slug  = %[5]q
  color = "ff0000"
}

resource "netbox_inventory_item_template" "test" {
  device_type  = netbox_device_type.test.id
  name         = %[1]q
  label        = "Test Label"
  description  = "Test Description"
  role         = netbox_inventory_item_role.extra.id
  manufacturer = netbox_manufacturer.extra.id
}
`, name, manufacturerName, manufacturerSlug, roleName, roleSlug)
}

func testAccInventoryItemTemplateResourceConfig_removeOptionalFields_detached(
	name, manufacturerName, manufacturerSlug, roleName, roleSlug string) string {
	return testAccInventoryItemTemplateResourcePrereqs(name) + fmt.Sprintf(`
resource "netbox_manufacturer" "extra" {
  name = %[2]q
  slug = %[3]q
}

resource "netbox_inventory_item_role" "extra" {
  name  = %[4]q
  slug  = %[5]q
  color = "ff0000"
}

resource "netbox_inventory_item_template" "test" {
  device_type = netbox_device_type.test.id
  name        = %[1]q
}
`, name, manufacturerName, manufacturerSlug, roleName, roleSlug)
}
