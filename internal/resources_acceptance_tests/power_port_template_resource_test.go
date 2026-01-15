package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func testAccPowerPortTemplateResourceBasic(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name string) string {
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

resource "netbox_power_port_template" "test" {
  device_type = netbox_device_type.test.id
  name        = %q
}
`, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name)
}

func testAccPowerPortTemplateResourceFull(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name, label, portType, description string, maxDraw, allocDraw int) string {
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

resource "netbox_power_port_template" "test" {
  device_type     = netbox_device_type.test.id
  name            = %q
  label           = %q
  type            = %q
  maximum_draw    = %d
  allocated_draw  = %d
  description     = %q
}
`, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name, label, portType, maxDraw, allocDraw, description)
}

func testAccPowerPortTemplateResourceConfig_forUpdate(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name, description string) string {
	maximumDraw := "100"
	if description == testutil.Description2 {
		maximumDraw = "200"
	}

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

resource "netbox_power_port_template" "test" {
  device_type  = netbox_device_type.test.id
  name         = %q
  maximum_draw = %s
  description  = %q
}
`, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name, maximumDraw, description)
}

func TestAccPowerPortTemplateResource_basic(t *testing.T) {

	t.Parallel()
	manufacturerName := testutil.RandomName("mfr")
	manufacturerSlug := testutil.RandomSlug("mfr")
	deviceTypeName := testutil.RandomName("dt")
	deviceTypeSlug := testutil.RandomSlug("dt")
	name := testutil.RandomName("power-port")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPowerPortTemplateResourceBasic(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_power_port_template.test", "id"),
					resource.TestCheckResourceAttr("netbox_power_port_template.test", "name", name),
				),
			},
			{
				ResourceName:            "netbox_power_port_template.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"device_type"},
			},
		},
	})
}

func TestAccPowerPortTemplateResource_update(t *testing.T) {
	t.Parallel()

	manufacturerName := testutil.RandomName("mfr-update")
	manufacturerSlug := testutil.RandomSlug("mfr-update")
	deviceTypeName := testutil.RandomName("dt-update")
	deviceTypeSlug := testutil.RandomSlug("dt-update")
	name := testutil.RandomName("power-port-update")
	updatedName := testutil.RandomName("power-port-updated")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPowerPortTemplateResourceConfig_forUpdate(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name, testutil.Description1),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_power_port_template.test", "id"),
					resource.TestCheckResourceAttr("netbox_power_port_template.test", "name", name),
					resource.TestCheckResourceAttr("netbox_power_port_template.test", "maximum_draw", "100"),
					resource.TestCheckResourceAttr("netbox_power_port_template.test", "description", testutil.Description1),
				),
			},
			{
				Config: testAccPowerPortTemplateResourceConfig_forUpdate(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, updatedName, testutil.Description2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_power_port_template.test", "name", updatedName),
					resource.TestCheckResourceAttr("netbox_power_port_template.test", "maximum_draw", "200"),
					resource.TestCheckResourceAttr("netbox_power_port_template.test", "description", testutil.Description2),
				),
			},
		},
	})
}

func TestAccPowerPortTemplateResource_full(t *testing.T) {

	t.Parallel()
	manufacturerName := testutil.RandomName("mfr")
	manufacturerSlug := testutil.RandomSlug("mfr")
	deviceTypeName := testutil.RandomName("dt")
	deviceTypeSlug := testutil.RandomSlug("dt")
	name := testutil.RandomName("power-port")
	label := testutil.RandomName("label")
	portType := "iec-60320-c14"
	description := testutil.RandomName("description")
	maxDraw := 500
	allocDraw := 400
	updatedName := testutil.RandomName("power-port")
	updatedLabel := testutil.RandomName("label")
	updatedDescription := testutil.RandomName("description")
	updatedMaxDraw := 600
	updatedAllocDraw := 450

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPowerPortTemplateResourceFull(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name, label, portType, description, maxDraw, allocDraw),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_power_port_template.test", "id"),
					resource.TestCheckResourceAttr("netbox_power_port_template.test", "name", name),
					resource.TestCheckResourceAttr("netbox_power_port_template.test", "label", label),
					resource.TestCheckResourceAttr("netbox_power_port_template.test", "type", portType),
					resource.TestCheckResourceAttr("netbox_power_port_template.test", "maximum_draw", "500"),
					resource.TestCheckResourceAttr("netbox_power_port_template.test", "allocated_draw", "400"),
					resource.TestCheckResourceAttr("netbox_power_port_template.test", "description", description),
				),
			},
			{
				Config: testAccPowerPortTemplateResourceFull(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, updatedName, updatedLabel, portType, updatedDescription, updatedMaxDraw, updatedAllocDraw),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_power_port_template.test", "name", updatedName),
					resource.TestCheckResourceAttr("netbox_power_port_template.test", "label", updatedLabel),
					resource.TestCheckResourceAttr("netbox_power_port_template.test", "maximum_draw", "600"),
					resource.TestCheckResourceAttr("netbox_power_port_template.test", "allocated_draw", "450"),
					resource.TestCheckResourceAttr("netbox_power_port_template.test", "description", updatedDescription),
				),
			},
		},
	})
}

func TestAccConsistency_PowerPortTemplate(t *testing.T) {

	t.Parallel()
	manufacturerName := testutil.RandomName("manufacturer")
	manufacturerSlug := testutil.RandomSlug("manufacturer")
	deviceTypeName := testutil.RandomName("devicetype")
	deviceTypeSlug := testutil.RandomSlug("devicetype")
	portName := testutil.RandomName("power-port")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPowerPortTemplateConsistencyConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, portName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_power_port_template.test", "name", portName),
					resource.TestCheckResourceAttr("netbox_power_port_template.test", "device_type", deviceTypeName),
				),
			},
			{
				PlanOnly: true,
				Config:   testAccPowerPortTemplateConsistencyConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, portName),
			},
		},
	})
}

func TestAccPowerPortTemplateResource_externalDeletion(t *testing.T) {
	t.Parallel()
	manufacturerName := testutil.RandomName("mfr-ext-del")
	manufacturerSlug := testutil.RandomSlug("mfr-ext-del")
	deviceTypeName := testutil.RandomName("dt-ext-del")
	deviceTypeSlug := testutil.RandomSlug("dt-ext-del")
	name := testutil.RandomName("power-port-ext-del")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPowerPortTemplateResourceBasic(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_power_port_template.test", "id"),
					resource.TestCheckResourceAttr("netbox_power_port_template.test", "name", name),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					items, _, err := client.DcimAPI.DcimPowerPortTemplatesList(context.Background()).NameIc([]string{name}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find power_port_template for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.DcimAPI.DcimPowerPortTemplatesDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete power_port_template: %v", err)
					}
					t.Logf("Successfully externally deleted power_port_template with ID: %d", itemID)
				},
				Config: testAccPowerPortTemplateResourceBasic(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_power_port_template.test", "id"),
				),
			},
		},
	})
}

func testAccPowerPortTemplateConsistencyConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, portName string) string {
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

resource "netbox_power_port_template" "test" {
  device_type = netbox_device_type.test.model
  name = "%[5]s"
  type = "iec-60320-c14"
}
`, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, portName)
}

func TestAccConsistency_PowerPortTemplate_LiteralNames(t *testing.T) {

	t.Parallel()
	manufacturerName := testutil.RandomName("manufacturer")
	manufacturerSlug := testutil.RandomSlug("manufacturer")
	deviceTypeName := testutil.RandomName("device-type")
	deviceTypeSlug := testutil.RandomSlug("device-type")
	resourceName := testutil.RandomName("power_port")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPowerPortTemplateConsistencyLiteralNamesConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_power_port_template.test", "name", resourceName),
					resource.TestCheckResourceAttr("netbox_power_port_template.test", "device_type", deviceTypeSlug),
				),
			},
			{
				PlanOnly: true,
				Config:   testAccPowerPortTemplateConsistencyLiteralNamesConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, resourceName),
			},
		},
	})
}

func testAccPowerPortTemplateConsistencyLiteralNamesConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, resourceName string) string {
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

resource "netbox_power_port_template" "test" {
  device_type = %q
  name = %q
  type = "iec-60320-c14"
  depends_on = [netbox_device_type.test]
}
`, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceTypeSlug, resourceName)
}

func TestAccPowerPortTemplateResource_removeOptionalFields(t *testing.T) {
	t.Parallel()

	mfgName := testutil.RandomName("tf-test-mfg-rem")
	mfgSlug := testutil.RandomSlug("tf-test-mfg-rem")
	dtModel := testutil.RandomName("tf-test-dt-rem")
	dtSlug := testutil.RandomSlug("tf-test-dt-rem")
	portName := testutil.RandomName("tf-test-ppt-rem")

	// Test values for all optional fields
	const testLabel = "Test Label"
	const testType = "iec-60320-c14"
	const testMaxDraw = 500
	const testAllocDraw = 400
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
				Config: testAccPowerPortTemplateResourceConfig_allOptionalFields(mfgName, mfgSlug, dtModel, dtSlug, portName, testLabel, testType, testMaxDraw, testAllocDraw, testDescription),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_power_port_template.test", "name", portName),
					resource.TestCheckResourceAttr("netbox_power_port_template.test", "label", testLabel),
					resource.TestCheckResourceAttr("netbox_power_port_template.test", "type", testType),
					resource.TestCheckResourceAttr("netbox_power_port_template.test", "maximum_draw", "500"),
					resource.TestCheckResourceAttr("netbox_power_port_template.test", "allocated_draw", "400"),
					resource.TestCheckResourceAttr("netbox_power_port_template.test", "description", testDescription),
				),
			},
			{
				// Step 2: Remove all optional fields
				Config: testAccPowerPortTemplateResourceConfig_noOptionalFields(mfgName, mfgSlug, dtModel, dtSlug, portName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_power_port_template.test", "name", portName),
					resource.TestCheckNoResourceAttr("netbox_power_port_template.test", "label"),
					resource.TestCheckNoResourceAttr("netbox_power_port_template.test", "type"),
					resource.TestCheckNoResourceAttr("netbox_power_port_template.test", "maximum_draw"),
					resource.TestCheckNoResourceAttr("netbox_power_port_template.test", "allocated_draw"),
					resource.TestCheckNoResourceAttr("netbox_power_port_template.test", "description"),
				),
			},
		},
	})
}

func testAccPowerPortTemplateResourceConfig_allOptionalFields(mfgName, mfgSlug, dtModel, dtSlug, portName, label, portType string, maxDraw, allocDraw int, description string) string {
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

resource "netbox_power_port_template" "test" {
  device_type = netbox_device_type.test.id
  name = %[5]q
  label = %[6]q
  type = %[7]q
  maximum_draw = %[8]d
  allocated_draw = %[9]d
  description = %[10]q
}
`, mfgName, mfgSlug, dtModel, dtSlug, portName, label, portType, maxDraw, allocDraw, description)
}

func testAccPowerPortTemplateResourceConfig_noOptionalFields(mfgName, mfgSlug, dtModel, dtSlug, portName string) string {
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

resource "netbox_power_port_template" "test" {
  device_type = netbox_device_type.test.id
  name = %[5]q
}
`, mfgName, mfgSlug, dtModel, dtSlug, portName)
}
