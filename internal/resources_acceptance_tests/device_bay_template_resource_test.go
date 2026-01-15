package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDeviceBayTemplateResource_basic(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-dbt")
	manufacturerName := testutil.RandomName("tf-test-mfr")
	manufacturerSlug := testutil.RandomSlug("tf-test-mfr")
	deviceTypeName := testutil.RandomName("tf-test-dt")
	deviceTypeSlug := testutil.RandomSlug("tf-test-dt")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDeviceBayTemplateResourceConfig_basic(name, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_device_bay_template.test", "id"),
					resource.TestCheckResourceAttr("netbox_device_bay_template.test", "name", name),
					resource.TestCheckResourceAttrSet("netbox_device_bay_template.test", "device_type"),
				),
			},
			{
				// Test import
				ResourceName:            "netbox_device_bay_template.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"device_type"},
			},
		},
	})
}

func TestAccDeviceBayTemplateResource_IDPreservation(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("dbt-id")
	manufacturerName := testutil.RandomName("mfr-dbt")
	manufacturerSlug := testutil.GenerateSlug(manufacturerName)
	deviceTypeName := testutil.RandomName("dt-dbt")
	deviceTypeSlug := testutil.GenerateSlug(deviceTypeName)

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDeviceBayTemplateResourceConfig_basic(name, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_device_bay_template.test", "id"),
					resource.TestCheckResourceAttr("netbox_device_bay_template.test", "name", name),
					resource.TestCheckResourceAttrSet("netbox_device_bay_template.test", "device_type"),
				),
			},
		},
	})
}

func TestAccDeviceBayTemplateResource_full(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-dbt-full")
	manufacturerName := testutil.RandomName("tf-test-mfr")
	manufacturerSlug := testutil.RandomSlug("tf-test-mfr")
	deviceTypeName := testutil.RandomName("tf-test-dt")
	deviceTypeSlug := testutil.RandomSlug("tf-test-dt")
	label := testutil.RandomName("label")
	description := testutil.RandomName("description")
	updatedLabel := testutil.RandomName("label-upd")
	updatedDescription := testutil.RandomName("description-upd")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDeviceBayTemplateResourceConfig_full(name, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, label, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_device_bay_template.test", "id"),
					resource.TestCheckResourceAttr("netbox_device_bay_template.test", "name", name),
					resource.TestCheckResourceAttr("netbox_device_bay_template.test", "label", label),
					resource.TestCheckResourceAttr("netbox_device_bay_template.test", "description", description),
					resource.TestCheckResourceAttrSet("netbox_device_bay_template.test", "device_type"),
				),
			},
			{
				Config: testAccDeviceBayTemplateResourceConfig_updated(name, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, updatedLabel, updatedDescription),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_device_bay_template.test", "id"),
					resource.TestCheckResourceAttr("netbox_device_bay_template.test", "name", name),
					resource.TestCheckResourceAttr("netbox_device_bay_template.test", "label", updatedLabel),
					resource.TestCheckResourceAttr("netbox_device_bay_template.test", "description", updatedDescription),
				),
			},
		},
	})
}

func TestAccDeviceBayTemplateResource_update(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-dbt-upd")
	manufacturerName := testutil.RandomName("tf-test-mfr")
	manufacturerSlug := testutil.RandomSlug("tf-test-mfr")
	deviceTypeName := testutil.RandomName("tf-test-dt")
	deviceTypeSlug := testutil.RandomSlug("tf-test-dt")
	updatedLabel := testutil.RandomName("label-upd")
	updatedDescription := testutil.RandomName("description-upd")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDeviceBayTemplateResourceConfig_basic(name, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_device_bay_template.test", "id"),
					resource.TestCheckResourceAttr("netbox_device_bay_template.test", "name", name),
				),
			},
			{
				Config: testAccDeviceBayTemplateResourceConfig_updated(name, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, updatedLabel, updatedDescription),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_device_bay_template.test", "id"),
					resource.TestCheckResourceAttr("netbox_device_bay_template.test", "name", name),
					resource.TestCheckResourceAttr("netbox_device_bay_template.test", "label", updatedLabel),
					resource.TestCheckResourceAttr("netbox_device_bay_template.test", "description", updatedDescription),
				),
			},
		},
	})
}

func TestAccDeviceBayTemplateResource_external_deletion(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("dbt-ext-del")
	manufacturerName := testutil.RandomName("mfr-ext-del")
	manufacturerSlug := testutil.RandomSlug("mfr-ext-del")
	deviceTypeName := testutil.RandomName("dt-ext-del")
	deviceTypeSlug := testutil.RandomSlug("dt-ext-del")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDeviceBayTemplateResourceConfig_basic(name, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_device_bay_template.test", "id"),
					resource.TestCheckResourceAttr("netbox_device_bay_template.test", "name", name),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					items, _, err := client.DcimAPI.DcimDeviceBayTemplatesList(context.Background()).NameIc([]string{name}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find device_bay_template for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.DcimAPI.DcimDeviceBayTemplatesDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete device_bay_template: %v", err)
					}
					t.Logf("Successfully externally deleted device_bay_template with ID: %d", itemID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccDeviceBayTemplateResourceConfig_basic(name, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug string) string {
	return fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
  name = %q
  slug = %q
}

resource "netbox_device_type" "test" {
  model          = %q
  slug           = %q
  manufacturer   = netbox_manufacturer.test.id
  subdevice_role = "parent"
}

resource "netbox_device_bay_template" "test" {
  device_type = netbox_device_type.test.id
  name        = %q
}
`, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name)
}

func testAccDeviceBayTemplateResourceConfig_full(name, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, label, description string) string {
	return fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_device_type" "test" {
  model          = %[3]q
  slug           = %[4]q
  manufacturer   = netbox_manufacturer.test.id
  subdevice_role = "parent"
}

resource "netbox_device_bay_template" "test" {
  device_type = netbox_device_type.test.id
  name        = %[5]q
  label       = %[6]q
  description = %[7]q
}
`, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name, label, description)

}

func TestAccConsistency_DeviceBayTemplate_LiteralNames(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-dbt-lit")
	manufacturerName := testutil.RandomName("tf-test-mfr")
	manufacturerSlug := testutil.RandomSlug("tf-test-mfr")
	deviceTypeName := testutil.RandomName("tf-test-dt")
	deviceTypeSlug := testutil.RandomSlug("tf-test-dt")
	label := testutil.RandomName("label")
	description := testutil.RandomName("description")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDeviceBayTemplateResourceConfig_updated(name, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, label, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_device_bay_template.test", "id"),
					resource.TestCheckResourceAttr("netbox_device_bay_template.test", "name", name),
					resource.TestCheckResourceAttr("netbox_device_bay_template.test", "label", label),
					resource.TestCheckResourceAttr("netbox_device_bay_template.test", "description", description),
				),
			},
			{
				Config:   testAccDeviceBayTemplateResourceConfig_updated(name, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, label, description),
				PlanOnly: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_device_bay_template.test", "id"),
				),
			},
		},
	})
}

func testAccDeviceBayTemplateResourceConfig_updated(name, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, label, description string) string {
	return fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_device_type" "test" {
  model          = %[3]q
  slug           = %[4]q
  manufacturer   = netbox_manufacturer.test.id
  subdevice_role = "parent"
}

resource "netbox_device_bay_template" "test" {
  device_type = netbox_device_type.test.id
  name        = %[5]q
  label       = %[6]q
  description = %[7]q
}
`, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name, label, description)
}

// TestAccDeviceBayTemplateResource_removeOptionalFields tests that optional fields
// can be successfully removed from the configuration without causing inconsistent state.
func TestAccDeviceBayTemplateResource_removeOptionalFields(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-dbt-rem")
	manufacturerName := testutil.RandomName("tf-test-mfr")
	manufacturerSlug := testutil.RandomSlug("tf-test-mfr")
	deviceTypeName := testutil.RandomName("tf-test-dt")
	deviceTypeSlug := testutil.RandomSlug("tf-test-dt")
	label := testutil.RandomName("label")
	description := testutil.RandomName("description")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDeviceBayTemplateResourceConfig_updated(name, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, label, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_device_bay_template.test", "name", name),
					resource.TestCheckResourceAttr("netbox_device_bay_template.test", "label", label),
					resource.TestCheckResourceAttr("netbox_device_bay_template.test", "description", description),
				),
			},
			{
				Config: testAccDeviceBayTemplateResourceConfig_basic(name, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_device_bay_template.test", "name", name),
					resource.TestCheckNoResourceAttr("netbox_device_bay_template.test", "label"),
					resource.TestCheckNoResourceAttr("netbox_device_bay_template.test", "description"),
				),
			},
		},
	})
}
func TestAccDeviceBayTemplateResource_validationErrors(t *testing.T) {
	testutil.RunMultiValidationErrorTest(t, testutil.MultiValidationErrorTestConfig{
		ResourceName: "netbox_device_bay_template",
		TestCases: map[string]testutil.ValidationErrorCase{
			"missing_device_type": {
				Config: func() string {
					return `
resource "netbox_device_bay_template" "test" {
  name = "Bay1"
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
			"missing_name": {
				Config: func() string {
					return `
resource "netbox_device_type" "test" {
  model = "Test Model"
  slug = "test-model"
  manufacturer = "1"
}

resource "netbox_device_bay_template" "test" {
  device_type = netbox_device_type.test.id
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
			"invalid_device_type_reference": {
				Config: func() string {
					return `
resource "netbox_device_bay_template" "test" {
  device_type = "99999"
  name = "Bay1"
}
`
				},
				ExpectedError: testutil.ErrPatternNotFound,
			},
		},
	})
}
