package resources_acceptance_tests

import (
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
)

// TestAccFrontPortTemplateResource_Label tests comprehensive scenarios for front port template label field.
// This validates that Optional+Computed string fields with empty string defaults work correctly.
func TestAccFrontPortTemplateResource_Label(t *testing.T) {

	testutil.RunOptionalComputedFieldTestSuite(t, testutil.OptionalComputedFieldTestConfig{
		ResourceName:   "netbox_front_port_template",
		OptionalField:  "label",
		DefaultValue:   "",
		FieldTestValue: "FP-01",
		BaseConfig: func() string {
			return `
resource "netbox_manufacturer" "test" {
	name = "test-manufacturer-front-port-template"
	slug = "test-manufacturer-front-port-template"
}

resource "netbox_device_type" "test" {
	manufacturer = netbox_manufacturer.test.id
	model        = "test-device-type-front-port-template"
	slug         = "test-device-type-front-port-template"
}

resource "netbox_rear_port_template" "test" {
	device_type = netbox_device_type.test.id
	name        = "rear-port-for-front-port-test"
	type        = "8p8c"
}

resource "netbox_front_port_template" "test" {
	device_type = netbox_device_type.test.id
	name        = "front-port-template-label-test"
	type        = "8p8c"
	rear_port   = netbox_rear_port_template.test.name
	# label field intentionally omitted - should get default ""
}
`
		},
		WithFieldConfig: func(value string) string {
			return `
resource "netbox_manufacturer" "test" {
	name = "test-manufacturer-front-port-template"
	slug = "test-manufacturer-front-port-template"
}

resource "netbox_device_type" "test" {
	manufacturer = netbox_manufacturer.test.id
	model        = "test-device-type-front-port-template"
	slug         = "test-device-type-front-port-template"
}

resource "netbox_rear_port_template" "test" {
	device_type = netbox_device_type.test.id
	name        = "rear-port-for-front-port-test"
	type        = "8p8c"
}

resource "netbox_front_port_template" "test" {
	device_type = netbox_device_type.test.id
	name        = "front-port-template-label-test"
	type        = "8p8c"
	rear_port   = netbox_rear_port_template.test.name
	label       = "` + value + `"
}
`
		},
	})
}

// TestAccFrontPortTemplateResource_Color tests comprehensive scenarios for front port template color field.
// This validates that Optional+Computed string fields with empty string defaults work correctly.
func TestAccFrontPortTemplateResource_Color(t *testing.T) {

	testutil.RunOptionalComputedFieldTestSuite(t, testutil.OptionalComputedFieldTestConfig{
		ResourceName:   "netbox_front_port_template",
		OptionalField:  "color",
		DefaultValue:   "",
		FieldTestValue: "aa1409",
		BaseConfig: func() string {
			return `
resource "netbox_manufacturer" "test" {
	name = "test-manufacturer-front-port-template-color"
	slug = "test-manufacturer-front-port-template-color"
}

resource "netbox_device_type" "test" {
	manufacturer = netbox_manufacturer.test.id
	model        = "test-device-type-front-port-template-color"
	slug         = "test-device-type-front-port-template-color"
}

resource "netbox_rear_port_template" "test" {
	device_type = netbox_device_type.test.id
	name        = "rear-port-for-front-port-color-test"
	type        = "8p8c"
}

resource "netbox_front_port_template" "test" {
	device_type = netbox_device_type.test.id
	name        = "front-port-template-color-test"
	type        = "8p8c"
	rear_port   = netbox_rear_port_template.test.name
	# color field intentionally omitted - should get default ""
}
`
		},
		WithFieldConfig: func(value string) string {
			return `
resource "netbox_manufacturer" "test" {
	name = "test-manufacturer-front-port-template-color"
	slug = "test-manufacturer-front-port-template-color"
}

resource "netbox_device_type" "test" {
	manufacturer = netbox_manufacturer.test.id
	model        = "test-device-type-front-port-template-color"
	slug         = "test-device-type-front-port-template-color"
}

resource "netbox_rear_port_template" "test" {
	device_type = netbox_device_type.test.id
	name        = "rear-port-for-front-port-color-test"
	type        = "8p8c"
}

resource "netbox_front_port_template" "test" {
	device_type = netbox_device_type.test.id
	name        = "front-port-template-color-test"
	type        = "8p8c"
	rear_port   = netbox_rear_port_template.test.name
	color       = "` + value + `"
}
`
		},
	})
}
