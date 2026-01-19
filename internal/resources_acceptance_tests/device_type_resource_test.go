package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDeviceTypeResource_basic(t *testing.T) {
	t.Parallel()

	model := testutil.RandomName("tf-test-device-type")
	slug := testutil.RandomSlug("tf-test-dt")
	manufacturerName := testutil.RandomName("tf-test-manufacturer")
	manufacturerSlug := testutil.RandomSlug("tf-test-mfr")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterDeviceTypeCleanup(slug)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDeviceTypeResourceConfig_basic(model, slug, manufacturerName, manufacturerSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_device_type.test", "id"),
					resource.TestCheckResourceAttr("netbox_device_type.test", "model", model),
					resource.TestCheckResourceAttr("netbox_device_type.test", "slug", slug),
					resource.TestCheckResourceAttrSet("netbox_device_type.test", "manufacturer"),
					resource.TestCheckResourceAttr("netbox_device_type.test", "u_height", "1"),
				),
			},
			{
				// Test import
				ResourceName:            "netbox_device_type.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"manufacturer"},
			},
			{
				Config:   testAccDeviceTypeResourceConfig_basic(model, slug, manufacturerName, manufacturerSlug),
				PlanOnly: true,
			},
		},
	})
}

func TestAccDeviceTypeResource_full(t *testing.T) {
	t.Parallel()

	model := testutil.RandomName("tf-test-device-type-full")
	slug := testutil.RandomSlug("tf-test-dt-full")
	manufacturerName := testutil.RandomName("tf-test-manufacturer")
	manufacturerSlug := testutil.RandomSlug("tf-test-mfr")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterDeviceTypeCleanup(slug)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDeviceTypeResourceConfig_full(model, slug, manufacturerName, manufacturerSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_device_type.test", "id"),
					resource.TestCheckResourceAttr("netbox_device_type.test", "model", model),
					resource.TestCheckResourceAttr("netbox_device_type.test", "slug", slug),
					resource.TestCheckResourceAttrSet("netbox_device_type.test", "manufacturer"),
					resource.TestCheckResourceAttr("netbox_device_type.test", "part_number", "TEST-PART-001"),
					resource.TestCheckResourceAttr("netbox_device_type.test", "u_height", "2"),
					resource.TestCheckResourceAttr("netbox_device_type.test", "is_full_depth", "true"),
					resource.TestCheckResourceAttr("netbox_device_type.test", "airflow", "front-to-rear"),
					resource.TestCheckResourceAttr("netbox_device_type.test", "description", "Test device type with full options"),
					resource.TestCheckResourceAttr("netbox_device_type.test", "comments", "Test comments for device type"),
				),
			},
		},
	})
}

func TestAccDeviceTypeResource_update(t *testing.T) {
	t.Parallel()

	model := testutil.RandomName("tf-test-device-type-update")
	slug := testutil.RandomSlug("tf-test-dt-upd")
	manufacturerName := testutil.RandomName("tf-test-manufacturer")
	manufacturerSlug := testutil.RandomSlug("tf-test-mfr")
	updatedModel := testutil.RandomName("tf-test-device-type-updated")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterDeviceTypeCleanup(slug)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDeviceTypeResourceConfig_basic(model, slug, manufacturerName, manufacturerSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_device_type.test", "id"),
					resource.TestCheckResourceAttr("netbox_device_type.test", "model", model),
					resource.TestCheckResourceAttr("netbox_device_type.test", "slug", slug),
				),
			},
			{
				Config: testAccDeviceTypeResourceConfig_updated(updatedModel, slug, manufacturerName, manufacturerSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_device_type.test", "id"),
					resource.TestCheckResourceAttr("netbox_device_type.test", "model", updatedModel),
					resource.TestCheckResourceAttr("netbox_device_type.test", "slug", slug),
					resource.TestCheckResourceAttr("netbox_device_type.test", "description", "Updated description"),
					resource.TestCheckResourceAttr("netbox_device_type.test", "u_height", "4"),
				),
			},
		},
	})
}

func testAccDeviceTypeResourceConfig_basic(model, slug, manufacturerName, manufacturerSlug string) string {
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
`, manufacturerName, manufacturerSlug, model, slug)
}

func testAccDeviceTypeResourceConfig_full(model, slug, manufacturerName, manufacturerSlug string) string {
	return fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
  name = %q
  slug = %q
}

resource "netbox_device_type" "test" {
  manufacturer      = netbox_manufacturer.test.id
  model             = %q
  slug              = %q
  part_number       = "TEST-PART-001"
  u_height          = 2
  is_full_depth     = true
  airflow           = "front-to-rear"
  description       = "Test device type with full options"
  comments          = "Test comments for device type"
}
`, manufacturerName, manufacturerSlug, model, slug)
}

func testAccDeviceTypeResourceConfig_updated(model, slug, manufacturerName, manufacturerSlug string) string {
	return fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
  name = %q
  slug = %q
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = %q
  slug         = %q
  u_height     = 4
  description  = "Updated description"
}
`, manufacturerName, manufacturerSlug, model, slug)
}

func TestAccConsistency_DeviceType_LiteralNames(t *testing.T) {
	t.Parallel()

	model := testutil.RandomName("tf-test-device-type-lit")
	slug := testutil.RandomSlug("tf-test-dt-lit")
	manufacturerName := testutil.RandomName("tf-test-manufacturer-lit")
	manufacturerSlug := testutil.RandomSlug("tf-test-mfr-lit")
	description := testutil.RandomName("description")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterDeviceTypeCleanup(slug)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.ComposeCheckDestroy(testutil.CheckDeviceTypeDestroy, testutil.CheckManufacturerDestroy),
		Steps: []resource.TestStep{
			{
				Config: testAccDeviceTypeConsistencyLiteralNamesConfig(model, slug, manufacturerName, manufacturerSlug, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_device_type.test", "id"),
					resource.TestCheckResourceAttr("netbox_device_type.test", "model", model),
					resource.TestCheckResourceAttr("netbox_device_type.test", "slug", slug),
					resource.TestCheckResourceAttr("netbox_device_type.test", "description", description),
				),
			},
			{
				Config:   testAccDeviceTypeConsistencyLiteralNamesConfig(model, slug, manufacturerName, manufacturerSlug, description),
				PlanOnly: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_device_type.test", "id"),
				),
			},
		},
	})
}

func testAccDeviceTypeConsistencyLiteralNamesConfig(model, slug, manufacturerName, manufacturerSlug, description string) string {
	return fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
  name = %q
  slug = %q
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = %q
  slug         = %q
  description  = %q
}
`, manufacturerName, manufacturerSlug, model, slug, description)
}

func TestAccDeviceTypeResource_externalDeletion(t *testing.T) {
	t.Parallel()

	model := testutil.RandomName("test-device-type-del")
	slug := testutil.GenerateSlug(model)
	manufacturerName := testutil.RandomName("test-manufacturer")
	manufacturerSlug := testutil.GenerateSlug(manufacturerName)

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterDeviceTypeCleanup(slug)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDeviceTypeResourceConfig_basic(model, slug, manufacturerName, manufacturerSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_device_type.test", "id"),
					resource.TestCheckResourceAttr("netbox_device_type.test", "model", model),
					resource.TestCheckResourceAttr("netbox_device_type.test", "slug", slug),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					items, _, err := client.DcimAPI.DcimDeviceTypesList(context.Background()).Slug([]string{slug}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find device_type for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.DcimAPI.DcimDeviceTypesDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete device_type: %v", err)
					}
					t.Logf("Successfully externally deleted device_type with ID: %d", itemID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccDeviceTypeResource_removeDescriptionAndComments(t *testing.T) {
	t.Parallel()

	model := testutil.RandomName("tf-test-devtype-optional")
	slug := testutil.RandomSlug("tf-test-devtype-optional")
	manufacturerName := testutil.RandomName("tf-test-manufacturer")
	manufacturerSlug := testutil.RandomSlug("tf-test-mfr")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterDeviceTypeCleanup(slug)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)

	testutil.TestRemoveOptionalFields(t, testutil.MultiFieldOptionalTestConfig{
		ResourceName: "netbox_device_type",
		BaseConfig: func() string {
			return testAccDeviceTypeResourceConfig_basic(model, slug, manufacturerName, manufacturerSlug)
		},
		ConfigWithFields: func() string {
			return testAccDeviceTypeResourceConfig_withDescriptionAndComments(
				model,
				slug,
				manufacturerName,
				manufacturerSlug,
				"Test description",
				"Test comments",
			)
		},
		OptionalFields: map[string]string{
			"description": "Test description",
			"comments":    "Test comments",
		},
		RequiredFields: map[string]string{
			"model": model,
			"slug":  slug,
		},
		CheckDestroy: testutil.CheckDeviceTypeDestroy,
	})
}

func testAccDeviceTypeResourceConfig_withDescriptionAndComments(model, slug, manufacturerName, manufacturerSlug, description, comments string) string {
	return fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
  name = %[3]q
  slug = %[4]q
}

resource "netbox_device_type" "test" {
  model        = %[1]q
  slug         = %[2]q
  manufacturer = netbox_manufacturer.test.id
  description  = %[5]q
  comments     = %[6]q
}
`, model, slug, manufacturerName, manufacturerSlug, description, comments)
}

func TestAccDeviceTypeResource_removeOptionalFields_part_number_u_height_weight(t *testing.T) {
	model := testutil.RandomName("tf-test-dt-opt")
	slug := testutil.RandomSlug("tf-test-dt-opt")
	mfgName := testutil.RandomName("tf-test-mfg-opt")
	mfgSlug := testutil.RandomSlug("tf-test-mfg-opt")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterDeviceTypeCleanup(slug)
	cleanup.RegisterManufacturerCleanup(mfgSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Step 1: Create with optional fields
			{
				Config: fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_device_type" "test" {
  model        = %[3]q
  slug         = %[4]q
  manufacturer = netbox_manufacturer.test.id
  part_number  = "PN-12345"
  u_height     = 2.0
  weight       = 10.5
  weight_unit  = "kg"
}
`, mfgName, mfgSlug, model, slug),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_device_type.test", "part_number", "PN-12345"),
					resource.TestCheckResourceAttr("netbox_device_type.test", "u_height", "2"),
					resource.TestCheckResourceAttr("netbox_device_type.test", "weight", "10.5"),
					resource.TestCheckResourceAttr("netbox_device_type.test", "weight_unit", "kg"),
				),
			},
			// Step 2: Remove optional fields
			{
				Config: fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_device_type" "test" {
  model        = %[3]q
  slug         = %[4]q
  manufacturer = netbox_manufacturer.test.id
  weight_unit  = "kg"
}
`, mfgName, mfgSlug, model, slug),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckNoResourceAttr("netbox_device_type.test", "part_number"),
					resource.TestCheckResourceAttr("netbox_device_type.test", "u_height", "1"),
					resource.TestCheckNoResourceAttr("netbox_device_type.test", "weight"),
					resource.TestCheckResourceAttr("netbox_device_type.test", "weight_unit", "kg"),
				),
			},
		},
	})
}

// TestAccDeviceTypeResource_validationErrors tests validation error scenarios.
func TestAccDeviceTypeResource_validationErrors(t *testing.T) {
	testutil.RunMultiValidationErrorTest(t, testutil.MultiValidationErrorTestConfig{
		ResourceName: "netbox_device_type",
		TestCases: map[string]testutil.ValidationErrorCase{
			"missing_manufacturer": {
				Config: func() string {
					return `
resource "netbox_device_type" "test" {
  model = "Test Model"
  slug  = "test-model"
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
			"missing_model": {
				Config: func() string {
					return `
resource "netbox_manufacturer" "test" {
  name = "Test Manufacturer"
  slug = "test-mfg"
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  slug         = "test-model"
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
			"missing_slug": {
				Config: func() string {
					return `
resource "netbox_manufacturer" "test" {
  name = "Test Manufacturer"
  slug = "test-mfg"
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = "Test Model"
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
			"invalid_manufacturer_reference": {
				Config: func() string {
					return `
resource "netbox_device_type" "test" {
  manufacturer = "99999"
  model        = "Test Model"
  slug         = "test-model"
}
`
				},
				ExpectedError: testutil.ErrPatternNotFound,
			},
			"invalid_airflow": {
				Config: func() string {
					return `
resource "netbox_manufacturer" "test" {
  name = "Test Manufacturer"
  slug = "test-mfg"
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = "Test Model"
  slug         = "test-model"
  airflow      = "invalid_airflow"
}
`
				},
				ExpectedError: testutil.ErrPatternInvalidEnum,
			},
			"invalid_weight_unit": {
				Config: func() string {
					return `
resource "netbox_manufacturer" "test" {
  name = "Test Manufacturer"
  slug = "test-mfg"
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = "Test Model"
  slug         = "test-model"
  weight_unit  = "invalid_unit"
}
`
				},
				ExpectedError: testutil.ErrPatternInvalidEnum,
			},
		},
	})
}
