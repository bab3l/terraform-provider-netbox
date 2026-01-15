package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func testAccPowerOutletTemplateResourceBasic(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name string) string {
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

resource "netbox_power_outlet_template" "test" {
  device_type = netbox_device_type.test.id
  name        = %q
}
`, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name)
}

func testAccPowerOutletTemplateResourceFull(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name, label, outletType, description string) string {
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

resource "netbox_power_outlet_template" "test" {
  device_type = netbox_device_type.test.id
  name        = %q
  label       = %q
  type        = %q
  description = %q
}
`, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name, label, outletType, description)
}

func TestAccPowerOutletTemplateResource_basic(t *testing.T) {
	t.Parallel()

	manufacturerName := testutil.RandomName("mfr")
	manufacturerSlug := testutil.RandomSlug("mfr")
	deviceTypeName := testutil.RandomName("dt")
	deviceTypeSlug := testutil.RandomSlug("dt")
	name := testutil.RandomName("power-outlet")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPowerOutletTemplateResourceBasic(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_power_outlet_template.test", "id"),
					resource.TestCheckResourceAttr("netbox_power_outlet_template.test", "name", name),
				),
			},
			{
				ResourceName:            "netbox_power_outlet_template.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"device_type"},
			},
		},
	})
}

func TestAccPowerOutletTemplateResource_full(t *testing.T) {
	t.Parallel()

	manufacturerName := testutil.RandomName("mfr")
	manufacturerSlug := testutil.RandomSlug("mfr")
	deviceTypeName := testutil.RandomName("dt")
	deviceTypeSlug := testutil.RandomSlug("dt")
	name := testutil.RandomName("power-outlet")
	label := testutil.RandomName("label")
	outletType := "iec-60320-c13"
	description := testutil.RandomName("description")
	updatedName := testutil.RandomName("power-outlet")
	updatedLabel := testutil.RandomName("label")
	updatedDescription := testutil.RandomName("description")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPowerOutletTemplateResourceFull(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name, label, outletType, description),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_power_outlet_template.test", "id"),
					resource.TestCheckResourceAttr("netbox_power_outlet_template.test", "name", name),
					resource.TestCheckResourceAttr("netbox_power_outlet_template.test", "label", label),
					resource.TestCheckResourceAttr("netbox_power_outlet_template.test", "type", outletType),
					resource.TestCheckResourceAttr("netbox_power_outlet_template.test", "description", description),
				),
			},
			{
				Config: testAccPowerOutletTemplateResourceFull(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, updatedName, updatedLabel, outletType, updatedDescription),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_power_outlet_template.test", "name", updatedName),
					resource.TestCheckResourceAttr("netbox_power_outlet_template.test", "label", updatedLabel),
					resource.TestCheckResourceAttr("netbox_power_outlet_template.test", "description", updatedDescription),
				),
			},
		},
	})
}

func TestAccConsistency_PowerOutletTemplate_LiteralNames(t *testing.T) {
	t.Parallel()

	manufacturerName := testutil.RandomName("manufacturer")
	manufacturerSlug := testutil.RandomSlug("manufacturer")
	deviceTypeName := testutil.RandomName("device-type")
	deviceTypeSlug := testutil.RandomSlug("device-type")
	resourceName := testutil.RandomName("power_outlet")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPowerOutletTemplateConsistencyLiteralNamesConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_power_outlet_template.test", "name", resourceName),
					resource.TestCheckResourceAttr("netbox_power_outlet_template.test", "device_type", deviceTypeSlug),
				),
			},
			{
				PlanOnly: true,
				Config:   testAccPowerOutletTemplateConsistencyLiteralNamesConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, resourceName),
			},
		},
	})
}
func TestAccPowerOutletTemplateResource_update(t *testing.T) {
	t.Parallel()

	manufacturerName := testutil.RandomName("mfr-update")
	manufacturerSlug := testutil.RandomSlug("mfr-update")
	deviceTypeName := testutil.RandomName("dt-update")
	deviceTypeSlug := testutil.RandomSlug("dt-update")
	name := testutil.RandomName("power-outlet-update")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPowerOutletTemplateResourceFull(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name, "Label1", "iec-60320-c13", testutil.Description1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_power_outlet_template.test", "id"),
					resource.TestCheckResourceAttr("netbox_power_outlet_template.test", "description", testutil.Description1),
				),
			},
			{
				Config: testAccPowerOutletTemplateResourceFull(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name, "Label2", "iec-60320-c13", testutil.Description2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_power_outlet_template.test", "description", testutil.Description2),
				),
			},
		},
	})
}

func TestAccPowerOutletTemplateResource_externalDeletion(t *testing.T) {
	t.Parallel()
	manufacturerName := testutil.RandomName("mfr-ext-del")
	manufacturerSlug := testutil.RandomSlug("mfr-ext-del")
	deviceTypeName := testutil.RandomName("dt-ext-del")
	deviceTypeSlug := testutil.RandomSlug("dt-ext-del")
	name := testutil.RandomName("power-outlet-ext-del")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPowerOutletTemplateResourceBasic(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_power_outlet_template.test", "id"),
					resource.TestCheckResourceAttr("netbox_power_outlet_template.test", "name", name),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					items, _, err := client.DcimAPI.DcimPowerOutletTemplatesList(context.Background()).NameIc([]string{name}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find power_outlet_template for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.DcimAPI.DcimPowerOutletTemplatesDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete power_outlet_template: %v", err)
					}
					t.Logf("Successfully externally deleted power_outlet_template with ID: %d", itemID)
				},
				Config: testAccPowerOutletTemplateResourceBasic(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_power_outlet_template.test", "id"),
				),
			},
		},
	})
}

func TestAccPowerOutletTemplateResource_IDPreservation(t *testing.T) {
	t.Parallel()
	manufacturerName := testutil.RandomName("mfr-id")
	manufacturerSlug := testutil.RandomSlug("mfr-id")
	deviceTypeName := testutil.RandomName("dt-id")
	deviceTypeSlug := testutil.RandomSlug("dt-id")
	name := testutil.RandomName("power-outlet-id")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPowerOutletTemplateResourceBasic(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_power_outlet_template.test", "id"),
					resource.TestCheckResourceAttr("netbox_power_outlet_template.test", "name", name),
				),
			},
		},
	})
}
func testAccPowerOutletTemplateConsistencyLiteralNamesConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, resourceName string) string {
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

resource "netbox_power_outlet_template" "test" {
  device_type = %q
  name = %q
  type = "iec-60320-c13"
  depends_on = [netbox_device_type.test]
}
`, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceTypeSlug, resourceName)
}

func TestAccPowerOutletTemplateResource_removeOptionalFields(t *testing.T) {
	t.Parallel()

	mfgName := testutil.RandomName("tf-test-mfg-rem")
	mfgSlug := testutil.RandomSlug("tf-test-mfg-rem")
	dtModel := testutil.RandomName("tf-test-dt-rem")
	dtSlug := testutil.RandomSlug("tf-test-dt-rem")
	portName := testutil.RandomName("tf-test-pot-rem")
	powerPortName := testutil.RandomName("tf-test-ppt-rem")

	// Test values for all optional fields
	const testLabel = "Test Label"
	const testType = "iec-60320-c13"
	const testFeedLeg = "A"
	const testDescription = "Test Description"

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(mfgSlug)
	cleanup.RegisterDeviceTypeCleanup(dtSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Step 1: Create with all optional fields
				Config: testAccPowerOutletTemplateResourceConfig_allOptionalFields(mfgName, mfgSlug, dtModel, dtSlug, portName, powerPortName, testLabel, testType, testFeedLeg, testDescription),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_power_outlet_template.test", "name", portName),
					resource.TestCheckResourceAttr("netbox_power_outlet_template.test", "label", testLabel),
					resource.TestCheckResourceAttr("netbox_power_outlet_template.test", "type", testType),
					resource.TestCheckResourceAttr("netbox_power_outlet_template.test", "feed_leg", testFeedLeg),
					resource.TestCheckResourceAttr("netbox_power_outlet_template.test", "description", testDescription),
					resource.TestCheckResourceAttrPair("netbox_power_outlet_template.test", "power_port", "netbox_power_port_template.feed", "id"),
				),
			},
			{
				// Step 2: Remove all optional fields
				Config: testAccPowerOutletTemplateResourceConfig_noOptionalFields(mfgName, mfgSlug, dtModel, dtSlug, portName, powerPortName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_power_outlet_template.test", "name", portName),
					resource.TestCheckNoResourceAttr("netbox_power_outlet_template.test", "label"),
					resource.TestCheckNoResourceAttr("netbox_power_outlet_template.test", "type"),
					resource.TestCheckNoResourceAttr("netbox_power_outlet_template.test", "feed_leg"),
					resource.TestCheckNoResourceAttr("netbox_power_outlet_template.test", "description"),
					resource.TestCheckNoResourceAttr("netbox_power_outlet_template.test", "power_port"),
				),
			},
		},
	})
}

func testAccPowerOutletTemplateResourceConfig_allOptionalFields(mfgName, mfgSlug, dtModel, dtSlug, portName, powerPortName, label, outletType, feedLeg, description string) string {
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

resource "netbox_power_port_template" "feed" {
	device_type = netbox_device_type.test.id
	name        = %[6]q
}

resource "netbox_power_outlet_template" "test" {
  device_type = netbox_device_type.test.id
  name = %[5]q
	label = %[7]q
	type = %[8]q
	feed_leg = %[9]q
	description = %[10]q
	power_port = netbox_power_port_template.feed.id
}
`, mfgName, mfgSlug, dtModel, dtSlug, portName, powerPortName, label, outletType, feedLeg, description)
}

func testAccPowerOutletTemplateResourceConfig_noOptionalFields(mfgName, mfgSlug, dtModel, dtSlug, portName, powerPortName string) string {
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

resource "netbox_power_port_template" "feed" {
	device_type = netbox_device_type.test.id
	name        = %[6]q
}

resource "netbox_power_outlet_template" "test" {
  device_type = netbox_device_type.test.id
  name = %[5]q
}
`, mfgName, mfgSlug, dtModel, dtSlug, portName, powerPortName)
}

func TestAccPowerOutletTemplateResource_validationErrors(t *testing.T) {
	testutil.RunMultiValidationErrorTest(t, testutil.MultiValidationErrorTestConfig{
		ResourceName: "netbox_power_outlet_template",
		TestCases: map[string]testutil.ValidationErrorCase{
			"missing_name": {
				Config: func() string {
					return `
resource "netbox_device_type" "test" {
  model = "Test Model"
  slug = "test-model"
  manufacturer = "1"
}

resource "netbox_power_outlet_template" "test" {
  device_type = netbox_device_type.test.id
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
			"invalid_device_type_reference": {
				Config: func() string {
					return `
resource "netbox_power_outlet_template" "test" {
  device_type = "99999"
  name = "Outlet1"
}
`
				},
				ExpectedError: testutil.ErrPatternNotFound,
			},
		},
	})
}
