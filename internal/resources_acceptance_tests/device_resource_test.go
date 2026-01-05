package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDeviceResource_basic(t *testing.T) {
	t.Parallel()

	// Generate unique names to avoid conflicts between test runs
	deviceName := testutil.RandomName("tf-test-device")
	manufacturerName := testutil.RandomName("tf-test-manufacturer")
	manufacturerSlug := testutil.RandomSlug("tf-test-mfr")
	deviceTypeModel := testutil.RandomName("tf-test-device-type")
	deviceTypeSlug := testutil.RandomSlug("tf-test-dt")
	deviceRoleName := testutil.RandomName("tf-test-device-role")
	deviceRoleSlug := testutil.RandomSlug("tf-test-dr")
	siteName := testutil.RandomName("tf-test-site")
	siteSlug := testutil.RandomSlug("tf-test-site")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
	cleanup.RegisterDeviceRoleCleanup(deviceRoleSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)
	cleanup.RegisterDeviceCleanup(deviceName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDeviceResourceConfig_basic(deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_device.test", "id"),
					resource.TestCheckResourceAttr("netbox_device.test", "name", deviceName),
					resource.TestCheckResourceAttrSet("netbox_device.test", "device_type"),
					resource.TestCheckResourceAttrSet("netbox_device.test", "role"),
					resource.TestCheckResourceAttrSet("netbox_device.test", "site"),
					resource.TestCheckResourceAttr("netbox_device.test", "status", "active"),
				),
			},
			{
				ResourceName:      "netbox_device.test",
				ImportState:       true,
				ImportStateVerify: true,
				// Note: some fields may use slugs in config but IDs in state after import
				ImportStateVerifyIgnore: []string{"device_type", "role", "site"},
			},
			{
				Config:   testAccDeviceResourceConfig_basic(deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug),
				PlanOnly: true,
			},
		},
	})
}

func TestAccDeviceResource_full(t *testing.T) {

	t.Parallel()
	// Generate unique names
	deviceName := testutil.RandomName("tf-test-device")
	manufacturerName := testutil.RandomName("tf-test-manufacturer")
	manufacturerSlug := testutil.RandomSlug("tf-test-mfr")
	deviceTypeModel := testutil.RandomName("tf-test-device-type")
	deviceTypeSlug := testutil.RandomSlug("tf-test-dt")
	deviceRoleName := testutil.RandomName("tf-test-device-role")
	deviceRoleSlug := testutil.RandomSlug("tf-test-dr")
	siteName := testutil.RandomName("tf-test-site")
	siteSlug := testutil.RandomSlug("tf-test-site")
	serial := testutil.RandomName("SN")
	assetTag := testutil.RandomName("AT")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)
	cleanup.RegisterDeviceRoleCleanup(deviceRoleSlug)
	cleanup.RegisterDeviceCleanup(deviceName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDeviceResourceConfig_full(deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug, serial, assetTag),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_device.test", "id"),
					resource.TestCheckResourceAttr("netbox_device.test", "name", deviceName),
					resource.TestCheckResourceAttrSet("netbox_device.test", "device_type"),
					resource.TestCheckResourceAttrSet("netbox_device.test", "role"),
					resource.TestCheckResourceAttrSet("netbox_device.test", "site"),
					resource.TestCheckResourceAttr("netbox_device.test", "status", "planned"),
					resource.TestCheckResourceAttr("netbox_device.test", "serial", serial),
					resource.TestCheckResourceAttr("netbox_device.test", "asset_tag", assetTag),
					resource.TestCheckResourceAttr("netbox_device.test", "description", "Test device description"),
					resource.TestCheckResourceAttr("netbox_device.test", "comments", "Test device comments"),
					resource.TestCheckResourceAttr("netbox_device.test", "airflow", "front-to-rear"),
					resource.TestCheckResourceAttr("netbox_device.test", "tags.#", "0"),
					resource.TestCheckResourceAttr("netbox_device.test", "custom_fields.#", "0"),
				),
			},
		},
	})
}

// TestAccDeviceResource_StatusOptionalField tests comprehensive scenarios for the device status optional field.
// This validates that Optional+Computed fields work correctly across all scenarios.
func TestAccDeviceResource_StatusOptionalField(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-site-device-status")
	siteSlug := testutil.RandomSlug("tf-test-site-device-status")
	manufacturerName := testutil.RandomName("tf-test-manufacturer-device-status")
	manufacturerSlug := testutil.RandomSlug("tf-test-manufacturer-device-status")
	deviceRoleName := testutil.RandomName("tf-test-device-role-status")
	deviceRoleSlug := testutil.RandomSlug("tf-test-device-role-status")
	deviceTypeName := testutil.RandomName("tf-test-device-type-status")
	deviceTypeSlug := testutil.RandomSlug("tf-test-device-type-status")
	deviceName := testutil.RandomName("tf-test-device-status")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)
	cleanup.RegisterDeviceRoleCleanup(deviceRoleSlug)
	cleanup.RegisterDeviceCleanup(deviceName)

	testutil.RunOptionalComputedFieldTestSuite(t, testutil.OptionalComputedFieldTestConfig{
		ResourceName:   "netbox_device",
		OptionalField:  "status",
		DefaultValue:   "active",
		FieldTestValue: "planned",
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckDeviceDestroy,
			testutil.CheckDeviceTypeDestroy,
			testutil.CheckDeviceRoleDestroy,
			testutil.CheckSiteDestroy,
			testutil.CheckManufacturerDestroy,
		),
		BaseConfig: func() string {
			return testAccDeviceResourceWithOptionalField(siteName, siteSlug, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceRoleName, deviceRoleSlug, deviceName, "status", "")
		},
		WithFieldConfig: func(value string) string {
			return testAccDeviceResourceWithOptionalField(siteName, siteSlug, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceRoleName, deviceRoleSlug, deviceName, "status", value)
		},
	})
}

func testAccDeviceResourceWithOptionalField(siteName, siteSlug, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceRoleName, deviceRoleSlug, deviceName, optionalFieldName, optionalFieldValue string) string {
	optionalField := ""
	if optionalFieldValue != "" {
		optionalField = fmt.Sprintf("\n  %s = %q", optionalFieldName, optionalFieldValue)
	}

	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name = %q
  slug = %q
}

resource "netbox_manufacturer" "test" {
  name = %q
  slug = %q
}

resource "netbox_device_role" "test" {
  name    = %q
  slug    = %q
  color   = "aa1409"
  vm_role = false
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = %q
  slug         = %q
}

resource "netbox_device" "test" {
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
  name        = %q%s
}
`, siteName, siteSlug, manufacturerName, manufacturerSlug, deviceRoleName, deviceRoleSlug, deviceTypeName, deviceTypeSlug, deviceName, optionalField)
}

func TestAccConsistency_Device(t *testing.T) {
	t.Parallel()

	deviceName := testutil.RandomName("device")
	deviceTypeName := testutil.RandomName("device-type")
	deviceTypeSlug := testutil.RandomSlug("device-type")
	manufacturerName := testutil.RandomName("manufacturer")
	manufacturerSlug := testutil.RandomSlug("manufacturer")
	roleName := testutil.RandomName("role")
	roleSlug := testutil.RandomSlug("role")
	siteName := testutil.RandomName("site")
	siteSlug := testutil.RandomSlug("site")
	tenantName := testutil.RandomName("tenant")
	tenantSlug := testutil.RandomSlug("tenant")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterDeviceRoleCleanup(roleSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDeviceConsistencyConfig(deviceName, deviceTypeName, deviceTypeSlug, manufacturerName, manufacturerSlug, roleName, roleSlug, siteName, siteSlug, tenantName, tenantSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_device.test", "name", deviceName),
					resource.TestCheckResourceAttrSet("netbox_device.test", "device_type"),
					resource.TestCheckResourceAttrSet("netbox_device.test", "role"),
					resource.TestCheckResourceAttrSet("netbox_device.test", "site"),
					resource.TestCheckResourceAttrSet("netbox_device.test", "tenant"),
				),
			},
			{
				PlanOnly: true,
				Config:   testAccDeviceConsistencyConfig(deviceName, deviceTypeName, deviceTypeSlug, manufacturerName, manufacturerSlug, roleName, roleSlug, siteName, siteSlug, tenantName, tenantSlug),
			},
		},
	})
}

func TestAccConsistency_Device_LiteralNames(t *testing.T) {
	t.Parallel()

	deviceName := testutil.RandomName("device")
	deviceTypeName := testutil.RandomName("device-type")
	deviceTypeSlug := testutil.RandomSlug("device-type")
	manufacturerName := testutil.RandomName("manufacturer")
	manufacturerSlug := testutil.RandomSlug("manufacturer")
	roleName := testutil.RandomName("role")
	roleSlug := testutil.RandomSlug("role")
	siteName := testutil.RandomName("site")
	siteSlug := testutil.RandomSlug("site")
	tenantName := testutil.RandomName("tenant")
	tenantSlug := testutil.RandomSlug("tenant")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterDeviceRoleCleanup(roleSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDeviceConsistencyLiteralNamesConfig(deviceName, deviceTypeName, deviceTypeSlug, manufacturerName, manufacturerSlug, roleName, roleSlug, siteName, siteSlug, tenantName, tenantSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_device.test", "name", deviceName),
					resource.TestCheckResourceAttr("netbox_device.test", "device_type", deviceTypeSlug),
					resource.TestCheckResourceAttr("netbox_device.test", "role", roleSlug),
					resource.TestCheckResourceAttr("netbox_device.test", "site", siteName),
					resource.TestCheckResourceAttr("netbox_device.test", "tenant", tenantName),
				),
			},
			{
				// Critical: Verify no drift when refreshing state
				PlanOnly: true,
				Config:   testAccDeviceConsistencyLiteralNamesConfig(deviceName, deviceTypeName, deviceTypeSlug, manufacturerName, manufacturerSlug, roleName, roleSlug, siteName, siteSlug, tenantName, tenantSlug),
			},
		},
	})
}

// Helper functions to generate test configurations.
func testAccDeviceResourceConfig_basic(deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug string) string {
	return fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = %[3]q
  slug         = %[4]q
}

resource "netbox_device_role" "test" {
  name = %[5]q
  slug = %[6]q
}

resource "netbox_site" "test" {
  name   = %[7]q
  slug   = %[8]q
  status = "active"
}

resource "netbox_device" "test" {
  name        = %[9]q
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
}
`, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug, deviceName)
}

func testAccDeviceResourceConfig_full(deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug, serial, assetTag string) string {
	return fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = %[3]q
  slug         = %[4]q
}

resource "netbox_device_role" "test" {
  name = %[5]q
  slug = %[6]q
}

resource "netbox_site" "test" {
  name   = %[7]q
  slug   = %[8]q
  status = "active"
}

resource "netbox_device" "test" {
  name        = %[9]q
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
  status      = "planned"
  serial      = %[10]q
  asset_tag   = %[11]q
  description = "Test device description"
  comments    = "Test device comments"
  airflow     = "front-to-rear"
  tags        = []
  custom_fields = []
}
`, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug, deviceName, serial, assetTag)
}

func testAccDeviceConsistencyConfig(deviceName, deviceTypeName, deviceTypeSlug, manufacturerName, manufacturerSlug, roleName, roleSlug, siteName, siteSlug, tenantName, tenantSlug string) string {
	return fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
  name = "%[4]s"
  slug = "%[5]s"
}

resource "netbox_device_type" "test" {
  model = "%[2]s"
  slug = "%[3]s"
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_device_role" "test" {
  name = "%[6]s"
  slug = "%[7]s"
}

resource "netbox_site" "test" {
  name = "%[8]s"
  slug = "%[9]s"
}

resource "netbox_tenant" "test" {
  name = "%[10]s"
  slug = "%[11]s"
}

resource "netbox_device" "test" {
  name = "%[1]s"
  device_type = netbox_device_type.test.id
  role = netbox_device_role.test.id
  site = netbox_site.test.id
  tenant = netbox_tenant.test.id
}
`, deviceName, deviceTypeName, deviceTypeSlug, manufacturerName, manufacturerSlug, roleName, roleSlug, siteName, siteSlug, tenantName, tenantSlug)
}

func testAccDeviceConsistencyLiteralNamesConfig(deviceName, deviceTypeName, deviceTypeSlug, manufacturerName, manufacturerSlug, roleName, roleSlug, siteName, siteSlug, tenantName, tenantSlug string) string {
	return fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
  name = "%[4]s"
  slug = "%[5]s"
}

resource "netbox_device_type" "test" {
  model = "%[2]s"
  slug = "%[3]s"
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_device_role" "test" {
  name = "%[6]s"
  slug = "%[7]s"
}

resource "netbox_site" "test" {
  name = "%[8]s"
  slug = "%[9]s"
}

resource "netbox_tenant" "test" {
  name = "%[10]s"
  slug = "%[11]s"
}

resource "netbox_device" "test" {
  name = "%[1]s"
  # Use literal string names to mimic existing user state
  device_type = "%[3]s"
  role = "%[7]s"
  site = "%[8]s"
  tenant = "%[10]s"

  depends_on = [netbox_device_type.test, netbox_device_role.test, netbox_site.test, netbox_tenant.test]
}
`, deviceName, deviceTypeName, deviceTypeSlug, manufacturerName, manufacturerSlug, roleName, roleSlug, siteName, siteSlug, tenantName, tenantSlug)
}

// TestAccDeviceResource_fieldConsistency tests that optional fields like airflow, tags, and custom_fields
// maintain consistency between plan and apply, addressing the "Provider produced inconsistent result after apply" bug.
func TestAccDeviceResource_fieldConsistency(t *testing.T) {
	t.Parallel()

	// Generate unique names to avoid conflicts
	deviceName := testutil.RandomName("tf-test-device-consistency")
	manufacturerName := testutil.RandomName("tf-test-manufacturer")
	manufacturerSlug := testutil.RandomSlug("tf-test-mfr")
	deviceTypeModel := testutil.RandomName("tf-test-device-type")
	deviceTypeSlug := testutil.RandomSlug("tf-test-dt")
	deviceRoleName := testutil.RandomName("tf-test-device-role")
	deviceRoleSlug := testutil.RandomSlug("tf-test-dr")
	siteName := testutil.RandomName("tf-test-site")
	siteSlug := testutil.RandomSlug("tf-test-site")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
	cleanup.RegisterDeviceRoleCleanup(deviceRoleSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)
	cleanup.RegisterDeviceCleanup(deviceName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Test creating device without specifying optional fields like airflow, tags, custom_fields
				Config: testAccDeviceResourceConfig_minimalOptionalFields(
					deviceName, manufacturerName, manufacturerSlug, deviceTypeModel,
					deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_device.test", "id"),
					resource.TestCheckResourceAttr("netbox_device.test", "name", deviceName),
					resource.TestCheckResourceAttrSet("netbox_device.test", "device_type"),
					resource.TestCheckResourceAttrSet("netbox_device.test", "role"),
					resource.TestCheckResourceAttrSet("netbox_device.test", "site"),
					resource.TestCheckResourceAttr("netbox_device.test", "status", "active"),
					// These fields should be handled consistently (null or computed default)
					// The exact values don't matter as much as consistency
				),
			},
			{
				// Test with empty sets for tags and custom_fields - should remain empty, not become null
				Config: testAccDeviceResourceConfig_emptySets(
					deviceName, manufacturerName, manufacturerSlug, deviceTypeModel,
					deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_device.test", "id"),
					resource.TestCheckResourceAttr("netbox_device.test", "name", deviceName),
					// tags and custom_fields should remain as empty sets, not null
					resource.TestCheckResourceAttr("netbox_device.test", "tags.#", "0"),
					resource.TestCheckResourceAttr("netbox_device.test", "custom_fields.#", "0"),
				),
			},
			{
				// Test with explicit airflow value
				Config: testAccDeviceResourceConfig_explicitAirflow(
					deviceName, manufacturerName, manufacturerSlug, deviceTypeModel,
					deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_device.test", "id"),
					resource.TestCheckResourceAttr("netbox_device.test", "airflow", "front-to-rear"),
				),
			},
		},
	})
}

// TestAccDeviceResource_importWithCustomFieldsAndTags tests importing a pre-existing device
// that has various custom field types and tags properly imports all data.
//
// IMPORTANT: This test must NOT run in parallel with other device tests because:
// NetBox custom fields are GLOBAL - when this test creates custom fields for "dcim.device",
// they appear on ALL device objects in the system. When this test deletes its custom fields
// during cleanup, other tests' devices may still have references to them, causing
// "Unknown field name" errors when those tests try to update their devices.
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
	cfLongtext := testutil.RandomCustomFieldName("tf_longtext")
	cfLongtextValue := fmt.Sprintf("This is a longer text field with multiple words: %s", testutil.RandomName("longtext"))
	cfIntegerName := testutil.RandomCustomFieldName("tf_integer")
	cfIntegerValue := 12345
	cfBoolean := testutil.RandomCustomFieldName("tf_boolean")
	cfBooleanValue := true
	cfDate := testutil.RandomCustomFieldName("tf_date")
	cfDateValue := testutil.RandomDate()
	cfURL := testutil.RandomCustomFieldName("tf_url")
	cfURLValue := testutil.RandomURL("device")
	cfJSON := testutil.RandomCustomFieldName("tf_json")
	cfJSONValue := testutil.RandomJSON()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Create the device with all custom fields and tags
				Config: testAccDeviceResourceImportConfig_full(
					deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug,
					deviceRoleName, deviceRoleSlug, siteName, siteSlug,
					tag1Name, tag1Slug, tag1Color, tag2Name, tag2Slug, tag2Color,
					cfText, cfTextValue, cfLongtext, cfLongtextValue, cfIntegerName, cfIntegerValue,
					cfBoolean, cfBooleanValue, cfDate, cfDateValue, cfURL, cfURLValue, cfJSON, cfJSONValue,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_device.test", "id"),
					resource.TestCheckResourceAttr("netbox_device.test", "name", deviceName),
					resource.TestCheckResourceAttr("netbox_device.test", "tags.#", "2"),
					resource.TestCheckResourceAttr("netbox_device.test", "custom_fields.#", "7"),
				),
			},
			{
				// Import the device and verify all fields are preserved
				ResourceName:            "netbox_device.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"device_type", "role", "site"}, // These use IDs after import but slugs in config
				// The import should preserve all custom fields and tags
				Check: resource.ComposeTestCheckFunc(
					// Verify basic fields
					resource.TestCheckResourceAttr("netbox_device.test", "name", deviceName),

					// Verify tags are imported
					resource.TestCheckResourceAttr("netbox_device.test", "tags.#", "2"),

					// Verify custom fields are imported
					resource.TestCheckResourceAttr("netbox_device.test", "custom_fields.#", "7"),

					// Verify specific custom field values
					resource.TestCheckResourceAttr("netbox_device.test", "custom_fields.0.name", cfText),
					resource.TestCheckResourceAttr("netbox_device.test", "custom_fields.0.type", "text"),
					resource.TestCheckResourceAttr("netbox_device.test", "custom_fields.0.value", cfTextValue),

					resource.TestCheckResourceAttr("netbox_device.test", "custom_fields.1.name", cfLongtext),
					resource.TestCheckResourceAttr("netbox_device.test", "custom_fields.1.type", "longtext"),
					resource.TestCheckResourceAttr("netbox_device.test", "custom_fields.1.value", cfLongtextValue),

					resource.TestCheckResourceAttr("netbox_device.test", "custom_fields.2.name", cfIntegerName),
					resource.TestCheckResourceAttr("netbox_device.test", "custom_fields.2.type", "integer"),
					resource.TestCheckResourceAttr("netbox_device.test", "custom_fields.2.value", fmt.Sprintf("%d", cfIntegerValue)),

					resource.TestCheckResourceAttr("netbox_device.test", "custom_fields.3.name", cfBoolean),
					resource.TestCheckResourceAttr("netbox_device.test", "custom_fields.3.type", "boolean"),
					resource.TestCheckResourceAttr("netbox_device.test", "custom_fields.3.value", fmt.Sprintf("%t", cfBooleanValue)),

					resource.TestCheckResourceAttr("netbox_device.test", "custom_fields.4.name", cfDate),
					resource.TestCheckResourceAttr("netbox_device.test", "custom_fields.4.type", "date"),
					resource.TestCheckResourceAttr("netbox_device.test", "custom_fields.4.value", cfDateValue),

					resource.TestCheckResourceAttr("netbox_device.test", "custom_fields.5.name", cfURL),
					resource.TestCheckResourceAttr("netbox_device.test", "custom_fields.5.type", "url"),
					resource.TestCheckResourceAttr("netbox_device.test", "custom_fields.5.value", cfURLValue),

					resource.TestCheckResourceAttr("netbox_device.test", "custom_fields.6.name", cfJSON),
					resource.TestCheckResourceAttr("netbox_device.test", "custom_fields.6.type", "json"),
					resource.TestCheckResourceAttr("netbox_device.test", "custom_fields.6.value", cfJSONValue),
				),
			},
			{
				// Verify no drift after import
				Config: testAccDeviceResourceImportConfig_full(
					deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug,
					deviceRoleName, deviceRoleSlug, siteName, siteSlug,
					tag1Name, tag1Slug, tag1Color, tag2Name, tag2Slug, tag2Color,
					cfText, cfTextValue, cfLongtext, cfLongtextValue, cfIntegerName, cfIntegerValue,
					cfBoolean, cfBooleanValue, cfDate, cfDateValue, cfURL, cfURLValue, cfJSON, cfJSONValue,
				),
				PlanOnly: true,
			},
		},
	})
}

func testAccDeviceResourceConfig_minimalOptionalFields(deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug string) string {
	return fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = %[3]q
  slug         = %[4]q
}

resource "netbox_device_role" "test" {
  name = %[5]q
  slug = %[6]q
}

resource "netbox_site" "test" {
  name   = %[7]q
  slug   = %[8]q
  status = "active"
}

resource "netbox_device" "test" {
  name        = %[9]q
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
  # Deliberately omitting airflow, tags, custom_fields to test consistency
}
`, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug, deviceName)
}

func testAccDeviceResourceConfig_emptySets(deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug string) string {
	return fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = %[3]q
  slug         = %[4]q
}

resource "netbox_device_role" "test" {
  name = %[5]q
  slug = %[6]q
}

resource "netbox_site" "test" {
  name   = %[7]q
  slug   = %[8]q
  status = "active"
}

resource "netbox_device" "test" {
  name         = %[9]q
  device_type  = netbox_device_type.test.id
  role         = netbox_device_role.test.id
  site         = netbox_site.test.id
  tags         = []
  custom_fields = []
  # Testing empty sets should remain empty, not become null
}
`, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug, deviceName)
}

func testAccDeviceResourceConfig_explicitAirflow(deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug string) string {
	return fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = %[3]q
  slug         = %[4]q
}

resource "netbox_device_role" "test" {
  name = %[5]q
  slug = %[6]q
}

resource "netbox_site" "test" {
  name   = %[7]q
  slug   = %[8]q
  status = "active"
}

resource "netbox_device" "test" {
  name        = %[9]q
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
  airflow     = "front-to-rear"
}
`, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug, deviceName)
}

func testAccDeviceResourceImportConfig_full(
	deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug,
	deviceRoleName, deviceRoleSlug, siteName, siteSlug,
	tag1Name, tag1Slug, tag1Color, tag2Name, tag2Slug, tag2Color string,
	cfText, cfTextValue, cfLongtext, cfLongtextValue string, cfIntegerName string, cfIntegerValue int,
	cfBoolean string, cfBooleanValue bool, cfDate, cfDateValue, cfURL, cfURLValue, cfJSON, cfJSONValue string,
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

resource "netbox_custom_field" "longtext" {
  name         = %[18]q
  type         = "longtext"
  object_types = ["dcim.device"]
  required     = false
}

resource "netbox_custom_field" "integer" {
  name         = %[20]q
  type         = "integer"
  object_types = ["dcim.device"]
  required     = false
}

resource "netbox_custom_field" "boolean" {
  name         = %[22]q
  type         = "boolean"
  object_types = ["dcim.device"]
  required     = false
}

resource "netbox_custom_field" "date" {
  name         = %[24]q
  type         = "date"
  object_types = ["dcim.device"]
  required     = false
}

resource "netbox_custom_field" "url" {
  name         = %[26]q
  type         = "url"
  object_types = ["dcim.device"]
  required     = false
}

resource "netbox_custom_field" "json" {
  name         = %[28]q
  type         = "json"
  object_types = ["dcim.device"]
  required     = false
}

# Create dependencies
resource "netbox_manufacturer" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = %[3]q
  slug         = %[4]q
}

resource "netbox_device_role" "test" {
  name = %[5]q
  slug = %[6]q
}

resource "netbox_site" "test" {
  name   = %[7]q
  slug   = %[8]q
  status = "active"
}

# Create device with all custom fields and tags
resource "netbox_device" "test" {
  name        = %[9]q
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
  status      = "active"

  tags = [
    {
      name = netbox_tag.test1.name
      slug = netbox_tag.test1.slug
    },
    {
      name = netbox_tag.test2.name
      slug = netbox_tag.test2.slug
    }
  ]

  custom_fields = [
    {
      name  = netbox_custom_field.text.name
      type  = "text"
      value = %[17]q
    },
    {
      name  = netbox_custom_field.longtext.name
      type  = "longtext"
      value = %[19]q
    },
    {
      name  = netbox_custom_field.integer.name
      type  = "integer"
      value = "%[21]d"
    },
    {
      name  = netbox_custom_field.boolean.name
      type  = "boolean"
      value = "%[23]t"
    },
    {
      name  = netbox_custom_field.date.name
      type  = "date"
      value = %[25]q
    },
    {
      name  = netbox_custom_field.url.name
      type  = "url"
      value = %[27]q
    },
    {
      name  = netbox_custom_field.json.name
      type  = "json"
      value = %[29]q
    }
  ]
}
`, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug, deviceName,
		tag1Name, tag1Slug, tag1Color, tag2Name, tag2Slug, tag2Color,
		cfText, cfTextValue, cfLongtext, cfLongtextValue, cfIntegerName, cfIntegerValue, cfBoolean, cfBooleanValue, cfDate, cfDateValue, cfURL, cfURLValue, cfJSON, cfJSONValue)
}

func TestAccDeviceResource_externalDeletion(t *testing.T) {
	t.Parallel()

	testutil.TestAccPreCheck(t)
	deviceName := testutil.RandomName("tf-test-device-ext-del")
	manufacturerName := testutil.RandomName("tf-test-manufacturer")
	manufacturerSlug := testutil.RandomSlug("tf-test-mfr")
	deviceTypeModel := testutil.RandomName("tf-test-device-type")
	deviceTypeSlug := testutil.RandomSlug("tf-test-dt")
	deviceRoleName := testutil.RandomName("tf-test-device-role")
	deviceRoleSlug := testutil.RandomSlug("tf-test-dr")
	siteName := testutil.RandomName("tf-test-site")
	siteSlug := testutil.RandomSlug("tf-test-site")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterDeviceCleanup(deviceName)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
	cleanup.RegisterDeviceRoleCleanup(deviceRoleSlug)
	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDeviceResourceConfig_basic(deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_device.test", "id"),
					resource.TestCheckResourceAttr("netbox_device.test", "name", deviceName),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					items, _, err := client.DcimAPI.DcimDevicesList(context.Background()).NameIc([]string{deviceName}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find device for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.DcimAPI.DcimDevicesDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete device: %v", err)
					}
					t.Logf("Successfully externally deleted device with ID: %d", itemID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}
