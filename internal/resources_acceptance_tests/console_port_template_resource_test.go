package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccConsolePortTemplateResource_basic(t *testing.T) {
	t.Parallel()

	manufacturerName := testutil.RandomName("mfr")
	manufacturerSlug := testutil.RandomSlug("mfr")
	deviceTypeName := testutil.RandomName("dt")
	deviceTypeSlug := testutil.RandomSlug("dt")
	name := testutil.RandomName("console")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccConsolePortTemplateResourceBasic(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_console_port_template.test", "id"),
					resource.TestCheckResourceAttr("netbox_console_port_template.test", "name", name),
				),
			},
			{
				ResourceName:            "netbox_console_port_template.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"device_type"},
			},
			{
				Config:             testAccConsolePortTemplateResourceBasic(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func TestAccConsolePortTemplateResource_full(t *testing.T) {
	t.Parallel()

	manufacturerName := testutil.RandomName("mfr")
	manufacturerSlug := testutil.RandomSlug("mfr")
	deviceTypeName := testutil.RandomName("dt")
	deviceTypeSlug := testutil.RandomSlug("dt")
	name := testutil.RandomName("console")
	label := testutil.RandomName("label")
	portType := "rj-45"
	description := testutil.RandomName("description")
	updatedName := testutil.RandomName("console-upd")
	updatedLabel := testutil.RandomName("label-upd")
	updatedDescription := testutil.RandomName("description-upd")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccConsolePortTemplateResourceFull(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name, label, portType, description),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_console_port_template.test", "id"),
					resource.TestCheckResourceAttr("netbox_console_port_template.test", "name", name),
					resource.TestCheckResourceAttr("netbox_console_port_template.test", "label", label),
					resource.TestCheckResourceAttr("netbox_console_port_template.test", "type", portType),
					resource.TestCheckResourceAttr("netbox_console_port_template.test", "description", description),
				),
			},
			{
				Config: testAccConsolePortTemplateResourceFull(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, updatedName, updatedLabel, portType, updatedDescription),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_console_port_template.test", "name", updatedName),
					resource.TestCheckResourceAttr("netbox_console_port_template.test", "label", updatedLabel),
					resource.TestCheckResourceAttr("netbox_console_port_template.test", "description", updatedDescription),
				),
			},
		},
	})
}

func TestAccConsistency_ConsolePortTemplate(t *testing.T) {
	t.Parallel()

	manufacturerName := testutil.RandomName("manufacturer")
	manufacturerSlug := testutil.RandomSlug("manufacturer")
	deviceTypeName := testutil.RandomName("devicetype")
	deviceTypeSlug := testutil.RandomSlug("devicetype")
	portName := testutil.RandomName("console")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccConsolePortTemplateConsistencyConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, portName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_console_port_template.test", "name", portName),
					resource.TestCheckResourceAttr("netbox_console_port_template.test", "device_type", deviceTypeName),
				),
			},
			{
				PlanOnly: true,
				Config:   testAccConsolePortTemplateConsistencyConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, portName),
			},
		},
	})
}

func TestAccConsistency_ConsolePortTemplate_LiteralNames(t *testing.T) {
	t.Parallel()

	manufacturerName := testutil.RandomName("manufacturer")
	manufacturerSlug := testutil.RandomSlug("manufacturer")
	deviceTypeName := testutil.RandomName("device-type")
	deviceTypeSlug := testutil.RandomSlug("device-type")
	portName := testutil.RandomName("port")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccConsolePortTemplateConsistencyLiteralNamesConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, portName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_console_port_template.test", "name", portName),
					resource.TestCheckResourceAttr("netbox_console_port_template.test", "device_type", deviceTypeSlug),
				),
			},
			{
				// Critical: Verify no drift when refreshing state
				PlanOnly: true,
				Config:   testAccConsolePortTemplateConsistencyLiteralNamesConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, portName),
			},
		},
	})
}

func TestAccConsolePortTemplateResource_update(t *testing.T) {
	t.Parallel()

	manufacturerName := testutil.RandomName("mfr-update")
	manufacturerSlug := testutil.RandomSlug("mfr-update")
	deviceTypeName := testutil.RandomName("dt-update")
	deviceTypeSlug := testutil.RandomSlug("dt-update")
	name := testutil.RandomName("console-update")
	label1 := "Initial Label"
	label2 := "Updated Label"

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccConsolePortTemplateResourceFull(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name, label1, "rj-45", testutil.Description1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_console_port_template.test", "id"),
					resource.TestCheckResourceAttr("netbox_console_port_template.test", "label", label1),
					resource.TestCheckResourceAttr("netbox_console_port_template.test", "description", testutil.Description1),
				),
			},
			{
				Config: testAccConsolePortTemplateResourceFull(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name, label2, "rj-45", testutil.Description2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_console_port_template.test", "id"),
					resource.TestCheckResourceAttr("netbox_console_port_template.test", "label", label2),
					resource.TestCheckResourceAttr("netbox_console_port_template.test", "description", testutil.Description2),
				),
			},
		},
	})
}

func TestAccConsolePortTemplateResource_externalDeletion(t *testing.T) {
	t.Parallel()

	manufacturerName := testutil.RandomName("mfr-ext-del")
	manufacturerSlug := testutil.RandomSlug("mfr-ext-del")
	deviceTypeName := testutil.RandomName("dt-ext-del")
	deviceTypeSlug := testutil.RandomSlug("dt-ext-del")
	name := testutil.RandomName("console-ext-del")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccConsolePortTemplateResourceBasic(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_console_port_template.test", "id"),
					resource.TestCheckResourceAttr("netbox_console_port_template.test", "name", name),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					items, _, err := client.DcimAPI.DcimConsolePortTemplatesList(context.Background()).NameIc([]string{name}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find console_port_template for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.DcimAPI.DcimConsolePortTemplatesDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete console_port_template: %v", err)
					}
					t.Logf("Successfully externally deleted console_port_template with ID: %d", itemID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccConsolePortTemplateResourceBasic(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name string) string {
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

resource "netbox_console_port_template" "test" {
  device_type = netbox_device_type.test.id
  name        = %q
}
`, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name)
}

func testAccConsolePortTemplateResourceFull(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name, label, portType, description string) string {
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

resource "netbox_console_port_template" "test" {
  device_type = netbox_device_type.test.id
  name        = %q
  label       = %q
  type        = %q
  description = %q
}
`, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name, label, portType, description)
}

func testAccConsolePortTemplateConsistencyConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, portName string) string {
	return fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
  name = "%[1]s"
  slug = "%[2]s"
}

resource "netbox_device_type" "test" {
  model = "%[3]s"
  slug = "%[4]s"
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_console_port_template" "test" {
  device_type = netbox_device_type.test.model
  name = "%[5]s"
  type = "rj-45"
}
`, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, portName)
}

func testAccConsolePortTemplateConsistencyLiteralNamesConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, portName string) string {
	return fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
  name = %q
  slug = %q
}

resource "netbox_device_type" "test" {
  model        = %q
  slug         = %q
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_console_port_template" "test" {
  # Use literal string slug to mimic existing user state
  device_type = %q
  name = %q
  type = "rj-45"
  depends_on = [netbox_device_type.test]
}
`, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceTypeSlug, portName)
}

func TestAccConsolePortTemplateResource_removeOptionalFields(t *testing.T) {
	t.Parallel()

	mfgName := testutil.RandomName("tf-test-mfg-rem")
	mfgSlug := testutil.RandomSlug("tf-test-mfg-rem")
	dtModel := testutil.RandomName("tf-test-dt-rem")
	dtSlug := testutil.RandomSlug("tf-test-dt-rem")
	portName := testutil.RandomName("tf-test-cpt-rem")
	const testLabel = "Test Label"
	const testDescription = "Test Description"
	const testType = "rj-45"

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(mfgSlug)
	cleanup.RegisterDeviceTypeCleanup(dtSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccConsolePortTemplateResourceConfig_withAllOptionalFields(mfgName, mfgSlug, dtModel, dtSlug, portName, testLabel, testDescription, testType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_console_port_template.test", "name", portName),
					resource.TestCheckResourceAttr("netbox_console_port_template.test", "label", testLabel),
					resource.TestCheckResourceAttr("netbox_console_port_template.test", "description", testDescription),
					resource.TestCheckResourceAttr("netbox_console_port_template.test", "type", testType),
				),
			},
			{
				Config: testAccConsolePortTemplateResourceBasic(mfgName, mfgSlug, dtModel, dtSlug, portName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_console_port_template.test", "name", portName),
					resource.TestCheckNoResourceAttr("netbox_console_port_template.test", "label"),
					resource.TestCheckNoResourceAttr("netbox_console_port_template.test", "description"),
					resource.TestCheckNoResourceAttr("netbox_console_port_template.test", "type"),
				),
			},
		},
	})
}

func testAccConsolePortTemplateResourceConfig_withAllOptionalFields(mfgName, mfgSlug, dtModel, dtSlug, portName, label, description, portType string) string {
	return fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_device_type" "test" {
  model = %[3]q
  slug = %[4]q
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_console_port_template" "test" {
  device_type = netbox_device_type.test.id
  name = %[5]q
  label = %[6]q
  description = %[7]q
  type = %[8]q
}
`, mfgName, mfgSlug, dtModel, dtSlug, portName, label, description, portType)
}

func TestAccConsolePortTemplateResource_validationErrors(t *testing.T) {
	testutil.RunMultiValidationErrorTest(t, testutil.MultiValidationErrorTestConfig{
		ResourceName: "netbox_console_port_template",
		TestCases: map[string]testutil.ValidationErrorCase{
			"missing_name": {
				Config: func() string {
					return `
resource "netbox_device_type" "test" {
  model = "Test Model"
  slug = "test-model"
  manufacturer = "1"
}

resource "netbox_console_port_template" "test" {
  device_type = netbox_device_type.test.id
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
			"invalid_device_type_reference": {
				Config: func() string {
					return `
resource "netbox_console_port_template" "test" {
  device_type = "99999"
  name = "Console1"
}
`
				},
				ExpectedError: testutil.ErrPatternNotFound,
			},
		},
	})
}
