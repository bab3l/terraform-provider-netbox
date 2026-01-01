package resources_acceptance_tests

import (
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
)

// TestAccRearPortTemplateResource_Label tests comprehensive scenarios for rear port template label field.
// This validates that Optional+Computed string fields with empty string defaults work correctly.
func TestAccRearPortTemplateResource_Label(t *testing.T) {
	// Generate unique names for this test run
	manufacturerName := testutil.RandomName("tf-test-mfr-rear-port-tpl")
	manufacturerSlug := testutil.RandomSlug("tf-test-mfr-rear-port-tpl")
	deviceTypeName := testutil.RandomName("tf-test-dev-type-rear-port-tpl")
	deviceTypeSlug := testutil.RandomSlug("tf-test-dev-type-rear-port-tpl")
	rearPortTemplateName := testutil.RandomName("tf-test-rear-port-tpl")

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
