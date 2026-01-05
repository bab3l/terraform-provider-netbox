package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/provider"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRearPortTemplateResource_basic(t *testing.T) {
	t.Parallel()

	manufacturerName := testutil.RandomName("mfr")
	manufacturerSlug := testutil.RandomSlug("mfr")
	deviceTypeName := testutil.RandomName("dt")
	deviceTypeSlug := testutil.RandomSlug("dt")
	name := testutil.RandomName("rear-port")
	portType := "8p8c"

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccRearPortTemplateResourceBasic(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name, portType),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_rear_port_template.test", "id"),
					resource.TestCheckResourceAttr("netbox_rear_port_template.test", "name", name),
					resource.TestCheckResourceAttr("netbox_rear_port_template.test", "type", portType),
				),
			},
			{
				ResourceName:            "netbox_rear_port_template.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"device_type"},
			},
		},
	})
}

func TestAccRearPortTemplateResource_full(t *testing.T) {
	t.Parallel()

	manufacturerName := testutil.RandomName("mfr")
	manufacturerSlug := testutil.RandomSlug("mfr")
	deviceTypeName := testutil.RandomName("dt")
	deviceTypeSlug := testutil.RandomSlug("dt")
	name := testutil.RandomName("rear-port")
	portType := "lc"
	label := testutil.RandomName("label")
	color := "aa1409"
	description := testutil.RandomName("description")
	positions := int32(2)
	updatedName := testutil.RandomName("rear-port")
	updatedLabel := testutil.RandomName("label")
	updatedDescription := testutil.RandomName("description")
	updatedPositions := int32(4)

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccRearPortTemplateResourceFull(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name, portType, label, color, description, positions),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_rear_port_template.test", "id"),
					resource.TestCheckResourceAttr("netbox_rear_port_template.test", "name", name),
					resource.TestCheckResourceAttr("netbox_rear_port_template.test", "type", portType),
					resource.TestCheckResourceAttr("netbox_rear_port_template.test", "label", label),
					resource.TestCheckResourceAttr("netbox_rear_port_template.test", "color", color),
					resource.TestCheckResourceAttr("netbox_rear_port_template.test", "positions", fmt.Sprintf("%d", positions)),
					resource.TestCheckResourceAttr("netbox_rear_port_template.test", "description", description),
				),
			},
			{
				Config: testAccRearPortTemplateResourceFull(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, updatedName, portType, updatedLabel, color, updatedDescription, updatedPositions),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_rear_port_template.test", "name", updatedName),
					resource.TestCheckResourceAttr("netbox_rear_port_template.test", "label", updatedLabel),
					resource.TestCheckResourceAttr("netbox_rear_port_template.test", "positions", fmt.Sprintf("%d", updatedPositions)),
					resource.TestCheckResourceAttr("netbox_rear_port_template.test", "description", updatedDescription),
				),
			},
			{
				ResourceName:            "netbox_rear_port_template.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"device_type"},
			},
		},
	})
}

func TestAccConsistency_RearPortTemplate(t *testing.T) {
	t.Parallel()

	manufacturerName := testutil.RandomName("manufacturer")
	manufacturerSlug := testutil.RandomSlug("manufacturer")
	deviceTypeName := testutil.RandomName("devicetype")
	deviceTypeSlug := testutil.RandomSlug("devicetype")
	portName := testutil.RandomName("rear-port")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRearPortTemplateConsistencyConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, portName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_rear_port_template.test", "name", portName),
					resource.TestCheckResourceAttr("netbox_rear_port_template.test", "device_type", deviceTypeName),
				),
			},
			{
				PlanOnly: true,
				Config:   testAccRearPortTemplateConsistencyConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, portName),
			},
		},
	})
}

func TestAccConsistency_RearPortTemplate_LiteralNames(t *testing.T) {
	t.Parallel()

	manufacturerName := testutil.RandomName("manufacturer")
	manufacturerSlug := testutil.RandomSlug("manufacturer")
	deviceTypeName := testutil.RandomName("device-type")
	deviceTypeSlug := testutil.RandomSlug("device-type")
	resourceName := testutil.RandomName("rear_port")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRearPortTemplateConsistencyLiteralNamesConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_rear_port_template.test", "name", resourceName),
					resource.TestCheckResourceAttr("netbox_rear_port_template.test", "device_type", deviceTypeSlug),
				),
			},
			{
				PlanOnly: true,
				Config:   testAccRearPortTemplateConsistencyLiteralNamesConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, resourceName),
			},
		},
	})
}

func TestAccRearPortTemplateResource_update(t *testing.T) {
	t.Parallel()

	manufacturerName := testutil.RandomName("mfr-update")
	manufacturerSlug := testutil.RandomSlug("mfr-update")
	deviceTypeName := testutil.RandomName("dt-update")
	deviceTypeSlug := testutil.RandomSlug("dt-update")
	rearPortName := testutil.RandomName("rear-port-update")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRearPortTemplateResourceFullSimple(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, rearPortName, "Label1", testutil.Description1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_rear_port_template.test", "id"),
					resource.TestCheckResourceAttr("netbox_rear_port_template.test", "label", "Label1"),
					resource.TestCheckResourceAttr("netbox_rear_port_template.test", "description", testutil.Description1),
				),
			},
			{
				Config: testAccRearPortTemplateResourceFullSimple(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, rearPortName, "Label2", testutil.Description2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_rear_port_template.test", "label", "Label2"),
					resource.TestCheckResourceAttr("netbox_rear_port_template.test", "description", testutil.Description2),
				),
			},
		},
	})
}

func TestAccRearPortTemplateResource_externalDeletion(t *testing.T) {
	t.Parallel()

	manufacturerName := testutil.RandomName("mfr-ext-del")
	manufacturerSlug := testutil.RandomSlug("mfr-ext-del")
	deviceTypeName := testutil.RandomName("dt-ext-del")
	deviceTypeSlug := testutil.RandomSlug("dt-ext-del")
	rearPortName := testutil.RandomName("rear-port-ext-del")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRearPortTemplateResourceBasic(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, rearPortName, "8p8c"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_rear_port_template.test", "id"),
					resource.TestCheckResourceAttr("netbox_rear_port_template.test", "name", rearPortName),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					items, _, err := client.DcimAPI.DcimRearPortTemplatesList(context.Background()).NameIc([]string{rearPortName}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find rear_port_template for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.DcimAPI.DcimRearPortTemplatesDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete rear_port_template: %v", err)
					}
					t.Logf("Successfully externally deleted rear_port_template with ID: %d", itemID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccRearPortTemplateResource_IDPreservation(t *testing.T) {
	t.Parallel()

	manufacturerName := testutil.RandomName("mfr-id")
	manufacturerSlug := testutil.RandomSlug("mfr-id")
	deviceTypeName := testutil.RandomName("dt-id")
	deviceTypeSlug := testutil.RandomSlug("dt-id")
	name := testutil.RandomName("rear-port-id")
	portType := "8p8c"

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccRearPortTemplateResourceBasic(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name, portType),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_rear_port_template.test", "id"),
					resource.TestCheckResourceAttr("netbox_rear_port_template.test", "name", name),
					resource.TestCheckResourceAttr("netbox_rear_port_template.test", "type", portType),
				),
			},
		},
	})

}

func testAccRearPortTemplateResourceBasic(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name, portType string) string {
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

resource "netbox_rear_port_template" "test" {
  device_type = netbox_device_type.test.id
  name        = %q
  type        = %q
}
`, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name, portType)
}

func testAccRearPortTemplateResourceFull(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name, portType, label, color, description string, positions int32) string {
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

resource "netbox_rear_port_template" "test" {
  device_type = netbox_device_type.test.id
  name        = %q
  type        = %q
  label       = %q
  color       = %q
  positions   = %d
  description = %q
}
`, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name, portType, label, color, positions, description)
}

func testAccRearPortTemplateConsistencyConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, portName string) string {
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

resource "netbox_rear_port_template" "test" {
  device_type = netbox_device_type.test.model
  name = "%[5]s"
  type = "8p8c"
  positions = 1
}
`, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, portName)
}

func testAccRearPortTemplateConsistencyLiteralNamesConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, resourceName string) string {
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

resource "netbox_rear_port_template" "test" {
  device_type = %q
  name = %q
  type = "8p8c"
  positions = 1
  depends_on = [netbox_device_type.test]
}
`, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceTypeSlug, resourceName)
}

func testAccRearPortTemplateResourceFullSimple(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, rearPortName, label, description string) string {
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

resource "netbox_rear_port_template" "test" {
  device_type = netbox_device_type.test.id
  name        = %q
  type        = "lc"
  positions   = 4
  label       = %q
  description = %q
}
`, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, rearPortName, label, description)
}

// TestAccRearPortTemplateResource_Label tests comprehensive scenarios for rear port template label field.
// This validates that Optional+Computed string fields with empty string defaults work correctly.
func TestAccRearPortTemplateResource_Label(t *testing.T) {
	// Generate unique names for this test run
	manufacturerName := testutil.RandomName("tf-test-mfr-rear-port-tpl")
	manufacturerSlug := testutil.RandomSlug("tf-test-mfr-rear-port-tpl")
	deviceTypeName := testutil.RandomName("tf-test-dev-type-rear-port-tpl")
	deviceTypeSlug := testutil.RandomSlug("tf-test-dev-type-rear-port-tpl")
	rearPortTemplateName := testutil.RandomName("tf-test-rear-port-tpl")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)

	testutil.RunOptionalComputedFieldTestSuite(t, testutil.OptionalComputedFieldTestConfig{
		ResourceName:   "netbox_rear_port_template",
		OptionalField:  "label",
		DefaultValue:   "",
		FieldTestValue: "RP-01",
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckRearPortTemplateDestroy,
			testutil.CheckDeviceTypeDestroy,
			testutil.CheckManufacturerDestroy,
		),
		BaseConfig: func() string {
			return `
resource "netbox_manufacturer" "test" {
	name = "` + manufacturerName + `"
	slug = "` + manufacturerSlug + `"
}

resource "netbox_device_type" "test" {
	manufacturer = netbox_manufacturer.test.id
	model        = "` + deviceTypeName + `"
	slug         = "` + deviceTypeSlug + `"
}

resource "netbox_rear_port_template" "test" {
	device_type = netbox_device_type.test.id
	name        = "` + rearPortTemplateName + `"
	type        = "8p8c"
	# label field intentionally omitted - should get default ""
}
`
		},
		WithFieldConfig: func(value string) string {
			return `
resource "netbox_manufacturer" "test" {
	name = "` + manufacturerName + `"
	slug = "` + manufacturerSlug + `"
}

resource "netbox_device_type" "test" {
	manufacturer = netbox_manufacturer.test.id
	model        = "` + deviceTypeName + `"
	slug         = "` + deviceTypeSlug + `"
}

resource "netbox_rear_port_template" "test" {
	device_type = netbox_device_type.test.id
	name        = "` + rearPortTemplateName + `"
	type        = "8p8c"
	label       = "` + value + `"
}
`
		},
	})
}

// TestAccRearPortTemplateResource_Color tests comprehensive scenarios for rear port template color field.
// This validates that Optional+Computed string fields with empty string defaults work correctly.
func TestAccRearPortTemplateResource_Color(t *testing.T) {
	// Generate unique names for this test run
	manufacturerName := testutil.RandomName("tf-test-mfr-rear-port-tpl-color")
	manufacturerSlug := testutil.RandomSlug("tf-test-mfr-rear-port-tpl-color")
	deviceTypeName := testutil.RandomName("tf-test-dev-type-rear-port-tpl-color")
	deviceTypeSlug := testutil.RandomSlug("tf-test-dev-type-rear-port-tpl-color")
	rearPortTemplateName := testutil.RandomName("tf-test-rear-port-tpl-color")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)

	testutil.RunOptionalComputedFieldTestSuite(t, testutil.OptionalComputedFieldTestConfig{
		ResourceName:   "netbox_rear_port_template",
		OptionalField:  "color",
		DefaultValue:   "",
		FieldTestValue: "aa1409",
		BaseConfig: func() string {
			return `
resource "netbox_manufacturer" "test" {
	name = "` + manufacturerName + `"
	slug = "` + manufacturerSlug + `"
}

resource "netbox_device_type" "test" {
	manufacturer = netbox_manufacturer.test.id
	model        = "` + deviceTypeName + `"
	slug         = "` + deviceTypeSlug + `"
}

resource "netbox_rear_port_template" "test" {
	device_type = netbox_device_type.test.id
	name        = "` + rearPortTemplateName + `"
	type        = "8p8c"
	# color field intentionally omitted - should get default ""
}
`
		},
		WithFieldConfig: func(value string) string {
			return `
resource "netbox_manufacturer" "test" {
	name = "` + manufacturerName + `"
	slug = "` + manufacturerSlug + `"
}

resource "netbox_device_type" "test" {
	manufacturer = netbox_manufacturer.test.id
	model        = "` + deviceTypeName + `"
	slug         = "` + deviceTypeSlug + `"
}

resource "netbox_rear_port_template" "test" {
	device_type = netbox_device_type.test.id
	name        = "` + rearPortTemplateName + `"
	type        = "8p8c"
	color       = "` + value + `"
}
`
		},
	})
}

// TestAccRearPortTemplateResource_Positions tests comprehensive scenarios for rear port template positions field.
// This validates that Optional+Computed int64 fields with proper defaults work correctly.
func TestAccRearPortTemplateResource_Positions(t *testing.T) {
	// Generate unique names for this test run
	manufacturerName := testutil.RandomName("tf-test-mfr-rear-port-tpl-pos")
	manufacturerSlug := testutil.RandomSlug("tf-test-mfr-rear-port-tpl-pos")
	deviceTypeName := testutil.RandomName("tf-test-dev-type-rear-port-tpl-pos")
	deviceTypeSlug := testutil.RandomSlug("tf-test-dev-type-rear-port-tpl-pos")
	rearPortTemplateName := testutil.RandomName("tf-test-rear-port-tpl-pos")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)

	testutil.RunOptionalComputedFieldTestSuite(t, testutil.OptionalComputedFieldTestConfig{
		ResourceName:   "netbox_rear_port_template",
		OptionalField:  "positions",
		DefaultValue:   "1",
		FieldTestValue: "4",
		BaseConfig: func() string {
			return `
resource "netbox_manufacturer" "test" {
	name = "` + manufacturerName + `"
	slug = "` + manufacturerSlug + `"
}

resource "netbox_device_type" "test" {
	manufacturer = netbox_manufacturer.test.id
	model        = "` + deviceTypeName + `"
	slug         = "` + deviceTypeSlug + `"
}

resource "netbox_rear_port_template" "test" {
	device_type = netbox_device_type.test.id
	name        = "` + rearPortTemplateName + `"
	type        = "8p8c"
	# positions field intentionally omitted - should get default 1
}
`
		},
		WithFieldConfig: func(value string) string {
			return `
resource "netbox_manufacturer" "test" {
	name = "` + manufacturerName + `"
	slug = "` + manufacturerSlug + `"
}

resource "netbox_device_type" "test" {
	manufacturer = netbox_manufacturer.test.id
	model        = "` + deviceTypeName + `"
	slug         = "` + deviceTypeSlug + `"
}

resource "netbox_rear_port_template" "test" {
	device_type = netbox_device_type.test.id
	name        = "` + rearPortTemplateName + `"
	type        = "8p8c"
	positions   = ` + value + `
}
`
		},
	})
}
