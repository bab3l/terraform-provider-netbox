//go:build customfields
// +build customfields

// Package resources_acceptance_tests_customfields contains acceptance tests for custom fields
// that require dedicated test runs to avoid conflicts with global custom field definitions.
package resources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccDeviceResource_CustomFieldsPreservation tests that custom fields are preserved
// when updating other fields on a device. This addresses a critical bug where custom fields
// were being deleted when users updated unrelated fields.
//
// Bug scenario:
// 1. Create device with custom fields
// 2. Update device WITHOUT custom_fields in config (omit the field entirely)
// 3. Custom fields should be preserved in NetBox, not deleted.
func TestAccDeviceResource_CustomFieldsPreservation(t *testing.T) {
	// Generate unique names
	deviceName := testutil.RandomName("tf-test-device-cf-preserve")
	manufacturerName := testutil.RandomName("tf-test-manufacturer")
	manufacturerSlug := testutil.RandomSlug("tf-test-mfr")
	deviceTypeModel := testutil.RandomName("tf-test-device-type")
	deviceTypeSlug := testutil.RandomSlug("tf-test-dt")
	deviceRoleName := testutil.RandomName("tf-test-device-role")
	deviceRoleSlug := testutil.RandomSlug("tf-test-dr")
	siteName := testutil.RandomName("tf-test-site")
	siteSlug := testutil.RandomSlug("tf-test-site")

	// Custom field names
	cfText := testutil.RandomCustomFieldName("tf_text_preserve")
	cfInteger := testutil.RandomCustomFieldName("tf_int_preserve")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Step 1: Create device WITH custom fields explicitly in config
				Config: testAccDeviceConfig_preservation_step1(
					deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug,
					deviceRoleName, deviceRoleSlug, siteName, siteSlug,
					cfText, cfInteger, "initial value", 42,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_device.test", "name", deviceName),
					resource.TestCheckResourceAttr("netbox_device.test", "description", "Initial description"),
					resource.TestCheckResourceAttr("netbox_device.test", "custom_fields.#", "2"),
					testutil.CheckCustomFieldValue("netbox_device.test", cfText, "text", "initial value"),
					testutil.CheckCustomFieldValue("netbox_device.test", cfInteger, "integer", "42"),
				),
			},
			{
				// Step 2: Update description WITHOUT mentioning custom_fields in config
				// Custom fields should be preserved in NetBox (verified by import)
				// State shows null/empty for custom_fields since not in config
				Config: testAccDeviceConfig_preservation_step2(
					deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug,
					deviceRoleName, deviceRoleSlug, siteName, siteSlug,
					cfText, cfInteger, "Updated description",
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_device.test", "name", deviceName),
					resource.TestCheckResourceAttr("netbox_device.test", "description", "Updated description"),
					// State shows 0 custom_fields (not in config = not owned)
					resource.TestCheckResourceAttr("netbox_device.test", "custom_fields.#", "0"),
				),
			},
			{
				// Step 3: Import to verify custom fields still exist in NetBox
				ResourceName:            "netbox_device.test",
				ImportState:             true,
				ImportStateKind:         resource.ImportCommandWithID,
				ImportStateVerify:       false,                     // Can't verify - config has no custom_fields
				ImportStateVerifyIgnore: []string{"custom_fields"}, // Different because filter-to-owned
			},
			{
				// Step 4: Add custom_fields back to config to verify they were preserved
				Config: testAccDeviceConfig_preservation_step1(
					deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug,
					deviceRoleName, deviceRoleSlug, siteName, siteSlug,
					cfText, cfInteger, "initial value", 42,
				),
				Check: resource.ComposeTestCheckFunc(
					// Custom fields should have their original values (preserved in NetBox)
					resource.TestCheckResourceAttr("netbox_device.test", "custom_fields.#", "2"),
					testutil.CheckCustomFieldValue("netbox_device.test", cfText, "text", "initial value"),
					testutil.CheckCustomFieldValue("netbox_device.test", cfInteger, "integer", "42"),
				),
			},
		},
	})
}

func testAccDeviceConfig_preservation_step1(
	deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug,
	deviceRoleName, deviceRoleSlug, siteName, siteSlug,
	cfTextName, cfIntName, cfTextValue string, cfIntValue int,
) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "text" {
  name         = %[10]q
  type         = "text"
  object_types = ["dcim.device"]
}

resource "netbox_custom_field" "integer" {
  name         = %[11]q
  type         = "integer"
  object_types = ["dcim.device"]
}

resource "netbox_manufacturer" "test" {
  name = %[2]q
  slug = %[3]q
}

resource "netbox_device_type" "test" {
  model        = %[4]q
  slug         = %[5]q
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_device_role" "test" {
  name  = %[6]q
  slug  = %[7]q
  color = "ff0000"
}

resource "netbox_site" "test" {
  name = %[8]q
  slug = %[9]q
}

resource "netbox_device" "test" {
  name        = %[1]q
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
  description = "Initial description"

  custom_fields = [
    {
      name  = netbox_custom_field.text.name
      type  = "text"
      value = %[12]q
    },
    {
      name  = netbox_custom_field.integer.name
      type  = "integer"
      value = "%[13]d"
    }
  ]
}
`,
		deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug,
		deviceRoleName, deviceRoleSlug, siteName, siteSlug,
		cfTextName, cfIntName, cfTextValue, cfIntValue,
	)
}

func testAccDeviceConfig_preservation_step2(
	deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug,
	deviceRoleName, deviceRoleSlug, siteName, siteSlug,
	cfTextName, cfIntName, description string,
) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "text" {
  name         = %[10]q
  type         = "text"
  object_types = ["dcim.device"]
}

resource "netbox_custom_field" "integer" {
  name         = %[11]q
  type         = "integer"
  object_types = ["dcim.device"]
}

resource "netbox_manufacturer" "test" {
  name = %[2]q
  slug = %[3]q
}

resource "netbox_device_type" "test" {
  model        = %[4]q
  slug         = %[5]q
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_device_role" "test" {
  name  = %[6]q
  slug  = %[7]q
  color = "ff0000"
}

resource "netbox_site" "test" {
  name = %[8]q
  slug = %[9]q
}

resource "netbox_device" "test" {
  name        = %[1]q
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
  description = %[12]q

  # custom_fields intentionally omitted - should be preserved in NetBox
  depends_on = [
    netbox_custom_field.text,
    netbox_custom_field.integer
  ]
}
`,
		deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug,
		deviceRoleName, deviceRoleSlug, siteName, siteSlug,
		cfTextName, cfIntName, description,
	)
}

// =============================================================================
// Filter-to-Owned Custom Fields Tests
// =============================================================================
//
// These tests verify the "filter-to-owned" pattern for custom fields:
//
// BEHAVIOR:
// 1. State only contains custom fields declared in config (owned fields)
// 2. Unowned fields are preserved in NetBox via merge but NOT in Terraform state
// 3. If config has custom_fields = null (omitted), state shows null
// 4. If config has custom_fields = [], state shows empty set and ALL fields cleared in NetBox
// 5. If config has specific fields, only those appear in state
//
// WHY THIS PATTERN:
// Terraform's framework requires plan and state to have the same structure for
// Optional+Computed attributes. We cannot return extra fields from the API
// that weren't in the config without causing plan/state mismatches.
//
// IMPORT BEHAVIOR:
// On import, there's no prior config, so:
// - If subsequent config declares custom_fields, those become owned
// - The import itself returns null custom_fields (no config to filter against)
// - Adding custom_fields to config after import works normally
//
// =============================================================================

// TestAccDeviceResource_CustomFieldsFilterToOwned tests the core filter-to-owned behavior:
// - Only fields declared in config appear in state
// - Unowned fields are preserved in NetBox but invisible to Terraform
// - Changing which fields are owned updates state appropriately
func TestAccDeviceResource_CustomFieldsFilterToOwned(t *testing.T) {
	// Generate unique names
	deviceName := testutil.RandomName("tf-test-device-cf-filter")
	manufacturerName := testutil.RandomName("tf-test-manufacturer")
	manufacturerSlug := testutil.RandomSlug("tf-test-mfr")
	deviceTypeModel := testutil.RandomName("tf-test-device-type")
	deviceTypeSlug := testutil.RandomSlug("tf-test-dt")
	deviceRoleName := testutil.RandomName("tf-test-device-role")
	deviceRoleSlug := testutil.RandomSlug("tf-test-dr")
	siteName := testutil.RandomName("tf-test-site")
	siteSlug := testutil.RandomSlug("tf-test-site")

	// Custom field names
	cfEnv := testutil.RandomCustomFieldName("tf_env")
	cfOwner := testutil.RandomCustomFieldName("tf_owner")
	cfCostCenter := testutil.RandomCustomFieldName("tf_cost_center")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Step 1: Create with two custom fields (env and owner)
				// State should show exactly 2 fields
				Config: testAccDeviceConfig_filterOwned_step1(
					deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug,
					deviceRoleName, deviceRoleSlug, siteName, siteSlug,
					cfEnv, cfOwner, cfCostCenter,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_device.test", "custom_fields.#", "2"),
					testutil.CheckCustomFieldValue("netbox_device.test", cfEnv, "text", "prod"),
					testutil.CheckCustomFieldValue("netbox_device.test", cfOwner, "text", "team-a"),
				),
			},
			{
				// Step 2: Remove owner from config, keep only env
				// State should show only 1 field (env with updated value)
				// owner is preserved in NetBox but NOT in state
				Config: testAccDeviceConfig_filterOwned_step2(
					deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug,
					deviceRoleName, deviceRoleSlug, siteName, siteSlug,
					cfEnv, cfOwner, cfCostCenter,
				),
				Check: resource.ComposeTestCheckFunc(
					// State shows only the owned field
					resource.TestCheckResourceAttr("netbox_device.test", "custom_fields.#", "1"),
					testutil.CheckCustomFieldValue("netbox_device.test", cfEnv, "text", "staging"),
				),
			},
			{
				// Step 3: Add cost_center to config, keep env
				// State should show 2 fields (env and cost_center)
				// owner is still preserved in NetBox but invisible
				Config: testAccDeviceConfig_filterOwned_step3(
					deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug,
					deviceRoleName, deviceRoleSlug, siteName, siteSlug,
					cfEnv, cfOwner, cfCostCenter,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_device.test", "custom_fields.#", "2"),
					testutil.CheckCustomFieldValue("netbox_device.test", cfEnv, "text", "staging"),
					testutil.CheckCustomFieldValue("netbox_device.test", cfCostCenter, "text", "123"),
				),
			},
			{
				// Step 4: Add owner back to config
				// State should show 3 fields - and owner should have its original value!
				// This verifies that unowned fields are truly preserved in NetBox
				Config: testAccDeviceConfig_filterOwned_step4(
					deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug,
					deviceRoleName, deviceRoleSlug, siteName, siteSlug,
					cfEnv, cfOwner, cfCostCenter,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_device.test", "custom_fields.#", "3"),
					testutil.CheckCustomFieldValue("netbox_device.test", cfEnv, "text", "staging"),
					testutil.CheckCustomFieldValue("netbox_device.test", cfOwner, "text", "team-a"), // Still preserved!
					testutil.CheckCustomFieldValue("netbox_device.test", cfCostCenter, "text", "123"),
				),
			},
		},
	})
}

// TestAccDeviceResource_CustomFieldsExplicitRemoval tests that when a field is removed
// from config, it's preserved in NetBox but removed from Terraform state (filter-to-owned).
func TestAccDeviceResource_CustomFieldsExplicitRemoval(t *testing.T) {
	// Generate unique names
	deviceName := testutil.RandomName("tf-test-device-cf-remove")
	manufacturerName := testutil.RandomName("tf-test-manufacturer")
	manufacturerSlug := testutil.RandomSlug("tf-test-mfr")
	deviceTypeModel := testutil.RandomName("tf-test-device-type")
	deviceTypeSlug := testutil.RandomSlug("tf-test-dt")
	deviceRoleName := testutil.RandomName("tf-test-device-role")
	deviceRoleSlug := testutil.RandomSlug("tf-test-dr")
	siteName := testutil.RandomName("tf-test-site")
	siteSlug := testutil.RandomSlug("tf-test-site")

	// Custom field names
	cfFieldA := testutil.RandomCustomFieldName("tf_field_a")
	cfFieldB := testutil.RandomCustomFieldName("tf_field_b")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Step 1: Create with two custom fields
				Config: testAccDeviceConfig_explicitRemoval_step1(
					deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug,
					deviceRoleName, deviceRoleSlug, siteName, siteSlug,
					cfFieldA, cfFieldB,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_device.test", "custom_fields.#", "2"),
					testutil.CheckCustomFieldValue("netbox_device.test", cfFieldA, "text", "value1"),
					testutil.CheckCustomFieldValue("netbox_device.test", cfFieldB, "text", "value2"),
				),
			},
			{
				// Step 2: Remove field_a from config (not owned anymore)
				// State should only show field_b
				// field_a is preserved in NetBox but not in Terraform state
				Config: testAccDeviceConfig_explicitRemoval_step2(
					deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug,
					deviceRoleName, deviceRoleSlug, siteName, siteSlug,
					cfFieldA, cfFieldB,
				),
				Check: resource.ComposeTestCheckFunc(
					// State shows only field_b (the owned field)
					resource.TestCheckResourceAttr("netbox_device.test", "custom_fields.#", "1"),
					testutil.CheckCustomFieldValue("netbox_device.test", cfFieldB, "text", "value2"),
				),
			},
		},
	})
}

// TestAccDeviceResource_CustomFieldsEmptyList tests setting custom_fields = []
// This is the explicit "clear all fields" operation.
func TestAccDeviceResource_CustomFieldsEmptyList(t *testing.T) {
	// Generate unique names
	deviceName := testutil.RandomName("tf-test-device-cf-clear")
	manufacturerName := testutil.RandomName("tf-test-manufacturer")
	manufacturerSlug := testutil.RandomSlug("tf-test-mfr")
	deviceTypeModel := testutil.RandomName("tf-test-device-type")
	deviceTypeSlug := testutil.RandomSlug("tf-test-dt")
	deviceRoleName := testutil.RandomName("tf-test-device-role")
	deviceRoleSlug := testutil.RandomSlug("tf-test-dr")
	siteName := testutil.RandomName("tf-test-site")
	siteSlug := testutil.RandomSlug("tf-test-site")

	// Custom field names
	cfFieldA := testutil.RandomCustomFieldName("tf_field_a")
	cfFieldB := testutil.RandomCustomFieldName("tf_field_b")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Step 1: Create with two custom fields
				Config: testAccDeviceConfig_emptyList_step1(
					deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug,
					deviceRoleName, deviceRoleSlug, siteName, siteSlug,
					cfFieldA, cfFieldB,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_device.test", "custom_fields.#", "2"),
				),
			},
			{
				// Step 2: Set custom_fields = [] - explicit clear all
				Config: testAccDeviceConfig_emptyList_step2(
					deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug,
					deviceRoleName, deviceRoleSlug, siteName, siteSlug,
					cfFieldA, cfFieldB,
				),
				Check: resource.ComposeTestCheckFunc(
					// Empty list means all fields cleared
					resource.TestCheckResourceAttr("netbox_device.test", "custom_fields.#", "0"),
				),
			},
		},
	})
}

// TestAccDeviceResource_importWithCustomFieldsAndTags tests importing a device
// that has custom fields and tags.
//
// With the filter-to-owned pattern:
// - Custom fields in state match what's declared in config
// - Import with config declaring custom fields should work
// - Tags work normally (not filtered)
//
// IMPORTANT: This test must NOT run in parallel with other device tests because:
// NetBox custom fields are GLOBAL - when this test creates custom fields for "dcim.device",
// they appear on ALL device objects in the system.
func TestAccDeviceResource_importWithCustomFieldsAndTags(t *testing.T) {
	// NOTE: t.Parallel() intentionally omitted - see function comment above

	// Generate unique names
	deviceName := testutil.RandomName("tf-test-device-import")
	manufacturerName := testutil.RandomName("tf-test-manufacturer")
	manufacturerSlug := testutil.RandomSlug("tf-test-mfr")
	deviceTypeModel := testutil.RandomName("tf-test-device-type")
	deviceTypeSlug := testutil.RandomSlug("tf-test-dt")
	deviceRoleName := testutil.RandomName("tf-test-device-role")
	deviceRoleSlug := testutil.RandomSlug("tf-test-dr")
	siteName := testutil.RandomName("tf-test-site")
	siteSlug := testutil.RandomSlug("tf-test-site")

	// Generate tag names
	tag1Name := testutil.RandomName("tf-test-tag1")
	tag1Slug := testutil.RandomSlug("tf-test-tag1")
	tag1Color := testutil.RandomColor()
	tag2Name := testutil.RandomName("tf-test-tag2")
	tag2Slug := testutil.RandomSlug("tf-test-tag2")
	tag2Color := testutil.RandomColor()

	// Generate custom field names and values
	cfText := testutil.RandomCustomFieldName("tf_text")
	cfTextValue := testutil.RandomName("text-value")
	cfInteger := testutil.RandomCustomFieldName("tf_integer")
	cfIntegerValue := 12345

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Step 1: Create the device with custom fields and tags
				Config: testAccDeviceConfig_importWithTags(
					deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug,
					deviceRoleName, deviceRoleSlug, siteName, siteSlug,
					tag1Name, tag1Slug, tag1Color, tag2Name, tag2Slug, tag2Color,
					cfText, cfTextValue, cfInteger, cfIntegerValue,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_device.test", "id"),
					resource.TestCheckResourceAttr("netbox_device.test", "name", deviceName),
					resource.TestCheckResourceAttr("netbox_device.test", "tags.#", "2"),
					resource.TestCheckResourceAttr("netbox_device.test", "custom_fields.#", "2"),
					testutil.CheckCustomFieldValue("netbox_device.test", cfText, "text", cfTextValue),
					testutil.CheckCustomFieldValue("netbox_device.test", cfInteger, "integer", fmt.Sprintf("%d", cfIntegerValue)),
				),
			},
			{
				// Step 2: Import the device
				// With filter-to-owned, import returns null custom_fields (no prior config to filter against)
				ResourceName:            "netbox_device.test",
				ImportState:             true,
				ImportStateKind:         resource.ImportBlockWithResourceIdentity,
				ImportStateVerify:       false,                                                    // Can't verify - import returns null custom_fields
				ImportStateVerifyIgnore: []string{"device_type", "role", "site", "custom_fields"}, // IDs vs slugs, custom_fields filtered
			},
			{
				// Step 3: Verify no drift after import
				Config: testAccDeviceConfig_importWithTags(
					deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug,
					deviceRoleName, deviceRoleSlug, siteName, siteSlug,
					tag1Name, tag1Slug, tag1Color, tag2Name, tag2Slug, tag2Color,
					cfText, cfTextValue, cfInteger, cfIntegerValue,
				),
				PlanOnly: true,
			},
		},
	})
}

// =============================================================================
// Helper Config Functions
// =============================================================================

func testAccDeviceConfig_base(
	manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug,
	deviceRoleName, deviceRoleSlug, siteName, siteSlug string,
	customFieldDefs string,
) string {
	return fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_device_type" "test" {
  model        = %[3]q
  slug         = %[4]q
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_device_role" "test" {
  name  = %[5]q
  slug  = %[6]q
  color = "ff0000"
}

resource "netbox_site" "test" {
  name = %[7]q
  slug = %[8]q
}

%s
`, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug,
		deviceRoleName, deviceRoleSlug, siteName, siteSlug, customFieldDefs)
}

// FilterOwned test configs
func testAccDeviceConfig_filterOwned_step1(
	deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug,
	deviceRoleName, deviceRoleSlug, siteName, siteSlug,
	cfEnv, cfOwner, cfCostCenter string,
) string {
	customFieldDefs := fmt.Sprintf(`
resource "netbox_custom_field" "env" {
  name         = %[1]q
  type         = "text"
  object_types = ["dcim.device"]
}

resource "netbox_custom_field" "owner" {
  name         = %[2]q
  type         = "text"
  object_types = ["dcim.device"]
}

resource "netbox_custom_field" "cost_center" {
  name         = %[3]q
  type         = "text"
  object_types = ["dcim.device"]
}
`, cfEnv, cfOwner, cfCostCenter)

	base := testAccDeviceConfig_base(
		manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug,
		deviceRoleName, deviceRoleSlug, siteName, siteSlug, customFieldDefs,
	)

	return base + fmt.Sprintf(`
resource "netbox_device" "test" {
  name        = %[1]q
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id

  custom_fields = [
    {
      name  = netbox_custom_field.env.name
      type  = "text"
      value = "prod"
    },
    {
      name  = netbox_custom_field.owner.name
      type  = "text"
      value = "team-a"
    }
  ]
}
`, deviceName)
}

func testAccDeviceConfig_filterOwned_step2(
	deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug,
	deviceRoleName, deviceRoleSlug, siteName, siteSlug,
	cfEnv, cfOwner, cfCostCenter string,
) string {
	customFieldDefs := fmt.Sprintf(`
resource "netbox_custom_field" "env" {
  name         = %[1]q
  type         = "text"
  object_types = ["dcim.device"]
}

resource "netbox_custom_field" "owner" {
  name         = %[2]q
  type         = "text"
  object_types = ["dcim.device"]
}

resource "netbox_custom_field" "cost_center" {
  name         = %[3]q
  type         = "text"
  object_types = ["dcim.device"]
}
`, cfEnv, cfOwner, cfCostCenter)

	base := testAccDeviceConfig_base(
		manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug,
		deviceRoleName, deviceRoleSlug, siteName, siteSlug, customFieldDefs,
	)

	return base + fmt.Sprintf(`
resource "netbox_device" "test" {
  name        = %[1]q
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id

  # owner removed from config - should be preserved in NetBox but not in state
  custom_fields = [
    {
      name  = netbox_custom_field.env.name
      type  = "text"
      value = "staging"
    }
  ]

  depends_on = [netbox_custom_field.owner]
}
`, deviceName)
}

func testAccDeviceConfig_filterOwned_step3(
	deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug,
	deviceRoleName, deviceRoleSlug, siteName, siteSlug,
	cfEnv, cfOwner, cfCostCenter string,
) string {
	customFieldDefs := fmt.Sprintf(`
resource "netbox_custom_field" "env" {
  name         = %[1]q
  type         = "text"
  object_types = ["dcim.device"]
}

resource "netbox_custom_field" "owner" {
  name         = %[2]q
  type         = "text"
  object_types = ["dcim.device"]
}

resource "netbox_custom_field" "cost_center" {
  name         = %[3]q
  type         = "text"
  object_types = ["dcim.device"]
}
`, cfEnv, cfOwner, cfCostCenter)

	base := testAccDeviceConfig_base(
		manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug,
		deviceRoleName, deviceRoleSlug, siteName, siteSlug, customFieldDefs,
	)

	return base + fmt.Sprintf(`
resource "netbox_device" "test" {
  name        = %[1]q
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id

  # owner still not in config, but add cost_center
  custom_fields = [
    {
      name  = netbox_custom_field.env.name
      type  = "text"
      value = "staging"
    },
    {
      name  = netbox_custom_field.cost_center.name
      type  = "text"
      value = "123"
    }
  ]

  depends_on = [netbox_custom_field.owner]
}
`, deviceName)
}

func testAccDeviceConfig_filterOwned_step4(
	deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug,
	deviceRoleName, deviceRoleSlug, siteName, siteSlug,
	cfEnv, cfOwner, cfCostCenter string,
) string {
	customFieldDefs := fmt.Sprintf(`
resource "netbox_custom_field" "env" {
  name         = %[1]q
  type         = "text"
  object_types = ["dcim.device"]
}

resource "netbox_custom_field" "owner" {
  name         = %[2]q
  type         = "text"
  object_types = ["dcim.device"]
}

resource "netbox_custom_field" "cost_center" {
  name         = %[3]q
  type         = "text"
  object_types = ["dcim.device"]
}
`, cfEnv, cfOwner, cfCostCenter)

	base := testAccDeviceConfig_base(
		manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug,
		deviceRoleName, deviceRoleSlug, siteName, siteSlug, customFieldDefs,
	)

	return base + fmt.Sprintf(`
resource "netbox_device" "test" {
  name        = %[1]q
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id

  # Add owner back - it should have preserved value from step 1
  custom_fields = [
    {
      name  = netbox_custom_field.env.name
      type  = "text"
      value = "staging"
    },
    {
      name  = netbox_custom_field.owner.name
      type  = "text"
      value = "team-a"
    },
    {
      name  = netbox_custom_field.cost_center.name
      type  = "text"
      value = "123"
    }
  ]
}
`, deviceName)
}

// ExplicitRemoval test configs
func testAccDeviceConfig_explicitRemoval_step1(
	deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug,
	deviceRoleName, deviceRoleSlug, siteName, siteSlug,
	cfFieldA, cfFieldB string,
) string {
	customFieldDefs := fmt.Sprintf(`
resource "netbox_custom_field" "field_a" {
  name         = %[1]q
  type         = "text"
  object_types = ["dcim.device"]
}

resource "netbox_custom_field" "field_b" {
  name         = %[2]q
  type         = "text"
  object_types = ["dcim.device"]
}
`, cfFieldA, cfFieldB)

	base := testAccDeviceConfig_base(
		manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug,
		deviceRoleName, deviceRoleSlug, siteName, siteSlug, customFieldDefs,
	)

	return base + fmt.Sprintf(`
resource "netbox_device" "test" {
  name        = %[1]q
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id

  custom_fields = [
    {
      name  = netbox_custom_field.field_a.name
      type  = "text"
      value = "value1"
    },
    {
      name  = netbox_custom_field.field_b.name
      type  = "text"
      value = "value2"
    }
  ]
}
`, deviceName)
}

func testAccDeviceConfig_explicitRemoval_step2(
	deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug,
	deviceRoleName, deviceRoleSlug, siteName, siteSlug,
	cfFieldA, cfFieldB string,
) string {
	customFieldDefs := fmt.Sprintf(`
resource "netbox_custom_field" "field_a" {
  name         = %[1]q
  type         = "text"
  object_types = ["dcim.device"]
}

resource "netbox_custom_field" "field_b" {
  name         = %[2]q
  type         = "text"
  object_types = ["dcim.device"]
}
`, cfFieldA, cfFieldB)

	base := testAccDeviceConfig_base(
		manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug,
		deviceRoleName, deviceRoleSlug, siteName, siteSlug, customFieldDefs,
	)

	return base + fmt.Sprintf(`
resource "netbox_device" "test" {
  name        = %[1]q
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id

  # field_a removed from config - preserved in NetBox but not in state
  custom_fields = [
    {
      name  = netbox_custom_field.field_b.name
      type  = "text"
      value = "value2"
    }
  ]

  depends_on = [netbox_custom_field.field_a]
}
`, deviceName)
}

// EmptyList test configs
func testAccDeviceConfig_emptyList_step1(
	deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug,
	deviceRoleName, deviceRoleSlug, siteName, siteSlug,
	cfFieldA, cfFieldB string,
) string {
	customFieldDefs := fmt.Sprintf(`
resource "netbox_custom_field" "field_a" {
  name         = %[1]q
  type         = "text"
  object_types = ["dcim.device"]
}

resource "netbox_custom_field" "field_b" {
  name         = %[2]q
  type         = "text"
  object_types = ["dcim.device"]
}
`, cfFieldA, cfFieldB)

	base := testAccDeviceConfig_base(
		manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug,
		deviceRoleName, deviceRoleSlug, siteName, siteSlug, customFieldDefs,
	)

	return base + fmt.Sprintf(`
resource "netbox_device" "test" {
  name        = %[1]q
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id

  custom_fields = [
    {
      name  = netbox_custom_field.field_a.name
      type  = "text"
      value = "value1"
    },
    {
      name  = netbox_custom_field.field_b.name
      type  = "text"
      value = "value2"
    }
  ]
}
`, deviceName)
}

func testAccDeviceConfig_emptyList_step2(
	deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug,
	deviceRoleName, deviceRoleSlug, siteName, siteSlug,
	cfFieldA, cfFieldB string,
) string {
	customFieldDefs := fmt.Sprintf(`
resource "netbox_custom_field" "field_a" {
  name         = %[1]q
  type         = "text"
  object_types = ["dcim.device"]
}

resource "netbox_custom_field" "field_b" {
  name         = %[2]q
  type         = "text"
  object_types = ["dcim.device"]
}
`, cfFieldA, cfFieldB)

	base := testAccDeviceConfig_base(
		manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug,
		deviceRoleName, deviceRoleSlug, siteName, siteSlug, customFieldDefs,
	)

	return base + fmt.Sprintf(`
resource "netbox_device" "test" {
  name        = %[1]q
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id

  custom_fields = [] # Explicit empty list clears all

  depends_on = [
    netbox_custom_field.field_a,
    netbox_custom_field.field_b
  ]
}
`, deviceName)
}

// Import test config (with tags for comprehensive testing)
func testAccDeviceConfig_importWithTags(
	deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug,
	deviceRoleName, deviceRoleSlug, siteName, siteSlug,
	tag1Name, tag1Slug, tag1Color, tag2Name, tag2Slug, tag2Color string,
	cfText, cfTextValue, cfInteger string, cfIntegerValue int,
) string {
	return fmt.Sprintf(`
# Create tags
resource "netbox_tag" "test1" {
  name  = %[10]q
  slug  = %[11]q
  color = %[12]q
}

resource "netbox_tag" "test2" {
  name  = %[13]q
  slug  = %[14]q
  color = %[15]q
}

# Create custom fields
resource "netbox_custom_field" "text" {
  name         = %[16]q
  type         = "text"
  object_types = ["dcim.device"]
  required     = false
}

resource "netbox_custom_field" "integer" {
  name         = %[17]q
  type         = "integer"
  object_types = ["dcim.device"]
  required     = false
}

# Create dependencies
resource "netbox_manufacturer" "test" {
  name = %[2]q
  slug = %[3]q
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = %[4]q
  slug         = %[5]q
}

resource "netbox_device_role" "test" {
  name = %[6]q
  slug = %[7]q
}

resource "netbox_site" "test" {
  name   = %[8]q
  slug   = %[9]q
  status = "active"
}

# Create device with custom fields and tags
resource "netbox_device" "test" {
  name        = %[1]q
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
  status      = "active"

	tags = [netbox_tag.test1.slug, netbox_tag.test2.slug]

  custom_fields = [
    {
      name  = netbox_custom_field.text.name
      type  = "text"
      value = %[18]q
    },
    {
      name  = netbox_custom_field.integer.name
      type  = "integer"
      value = "%[19]d"
    }
  ]
}
`, deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug,
		tag1Name, tag1Slug, tag1Color, tag2Name, tag2Slug, tag2Color,
		cfText, cfInteger, cfTextValue, cfIntegerValue)
}
