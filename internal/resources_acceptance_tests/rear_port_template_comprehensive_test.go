package resources_acceptance_tests

import (
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
)

// TestAccRearPortTemplateResource_Label tests comprehensive scenarios for rear port template label field.
// This validates that Optional+Computed string fields with empty string defaults work correctly.
func TestAccRearPortTemplateResource_Label(t *testing.T) {

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
	name = "test-manufacturer-rear-port-template"
	slug = "test-manufacturer-rear-port-template"
}

resource "netbox_device_type" "test" {
	manufacturer = netbox_manufacturer.test.id
	model        = "test-device-type-rear-port-template"
	slug         = "test-device-type-rear-port-template"
}

resource "netbox_rear_port_template" "test" {
	device_type = netbox_device_type.test.id
	name        = "rear-port-template-label-test"
	type        = "8p8c"
	# label field intentionally omitted - should get default ""
}
`
		},
		WithFieldConfig: func(value string) string {
			return `
resource "netbox_manufacturer" "test" {
	name = "test-manufacturer-rear-port-template"
	slug = "test-manufacturer-rear-port-template"
}

resource "netbox_device_type" "test" {
	manufacturer = netbox_manufacturer.test.id
	model        = "test-device-type-rear-port-template"
	slug         = "test-device-type-rear-port-template"
}

resource "netbox_rear_port_template" "test" {
	device_type = netbox_device_type.test.id
	name        = "rear-port-template-label-test"
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

	testutil.RunOptionalComputedFieldTestSuite(t, testutil.OptionalComputedFieldTestConfig{
		ResourceName:   "netbox_rear_port_template",
		OptionalField:  "color",
		DefaultValue:   "",
		FieldTestValue: "aa1409",
		BaseConfig: func() string {
			return `
resource "netbox_manufacturer" "test" {
	name = "test-manufacturer-rear-port-template-color"
	slug = "test-manufacturer-rear-port-template-color"
}

resource "netbox_device_type" "test" {
	manufacturer = netbox_manufacturer.test.id
	model        = "test-device-type-rear-port-template-color"
	slug         = "test-device-type-rear-port-template-color"
}

resource "netbox_rear_port_template" "test" {
	device_type = netbox_device_type.test.id
	name        = "rear-port-template-color-test"
	type        = "8p8c"
	# color field intentionally omitted - should get default ""
}
`
		},
		WithFieldConfig: func(value string) string {
			return `
resource "netbox_manufacturer" "test" {
	name = "test-manufacturer-rear-port-template-color"
	slug = "test-manufacturer-rear-port-template-color"
}

resource "netbox_device_type" "test" {
	manufacturer = netbox_manufacturer.test.id
	model        = "test-device-type-rear-port-template-color"
	slug         = "test-device-type-rear-port-template-color"
}

resource "netbox_rear_port_template" "test" {
	device_type = netbox_device_type.test.id
	name        = "rear-port-template-color-test"
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

	testutil.RunOptionalComputedFieldTestSuite(t, testutil.OptionalComputedFieldTestConfig{
		ResourceName:   "netbox_rear_port_template",
		OptionalField:  "positions",
		DefaultValue:   "1",
		FieldTestValue: "4",
		BaseConfig: func() string {
			return `
resource "netbox_manufacturer" "test" {
	name = "test-manufacturer-rear-port-template-positions"
	slug = "test-manufacturer-rear-port-template-positions"
}

resource "netbox_device_type" "test" {
	manufacturer = netbox_manufacturer.test.id
	model        = "test-device-type-rear-port-template-positions"
	slug         = "test-device-type-rear-port-template-positions"
}

resource "netbox_rear_port_template" "test" {
	device_type = netbox_device_type.test.id
	name        = "rear-port-template-positions-test"
	type        = "8p8c"
	# positions field intentionally omitted - should get default 1
}
`
		},
		WithFieldConfig: func(value string) string {
			return `
resource "netbox_manufacturer" "test" {
	name = "test-manufacturer-rear-port-template-positions"
	slug = "test-manufacturer-rear-port-template-positions"
}

resource "netbox_device_type" "test" {
	manufacturer = netbox_manufacturer.test.id
	model        = "test-device-type-rear-port-template-positions"
	slug         = "test-device-type-rear-port-template-positions"
}

resource "netbox_rear_port_template" "test" {
	device_type = netbox_device_type.test.id
	name        = "rear-port-template-positions-test"
	type        = "8p8c"
	positions   = ` + value + `
}
`
		},
	})
}
