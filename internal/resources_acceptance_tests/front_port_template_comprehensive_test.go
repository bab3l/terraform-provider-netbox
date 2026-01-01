package resources_acceptance_tests

import (
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
)

// TestAccFrontPortTemplateResource_Label tests comprehensive scenarios for front port template label field.
// This validates that Optional+Computed string fields with empty string defaults work correctly.
func TestAccFrontPortTemplateResource_Label(t *testing.T) {

	// Generate unique names for this test run
	manufacturerName := testutil.RandomName("tf-test-manufacturer-fpt-label")
	manufacturerSlug := testutil.RandomSlug("tf-test-manufacturer-fpt-label")
	deviceTypeName := testutil.RandomName("tf-test-device-type-fpt-label")
	deviceTypeSlug := testutil.RandomSlug("tf-test-device-type-fpt-label")
	rearPortName := testutil.RandomName("tf-test-rear-port-fpt-label")
	frontPortName := testutil.RandomName("tf-test-front-port-fpt-label")

	testutil.RunOptionalComputedFieldTestSuite(t, testutil.OptionalComputedFieldTestConfig{
		ResourceName:   "netbox_front_port_template",
		OptionalField:  "label",
		DefaultValue:   "",
		FieldTestValue: "FP-01",
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckFrontPortTemplateDestroy,
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
	name        = "` + rearPortName + `"
	type        = "8p8c"
}

resource "netbox_front_port_template" "test" {
	device_type = netbox_device_type.test.id
	name        = "` + frontPortName + `"
	type        = "8p8c"
	rear_port   = netbox_rear_port_template.test.name
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
	name        = "` + rearPortName + `"
	type        = "8p8c"
}

resource "netbox_front_port_template" "test" {
	device_type = netbox_device_type.test.id
	name        = "` + frontPortName + `"
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

	// Generate unique names for this test run
	manufacturerName := testutil.RandomName("tf-test-manufacturer-fpt-color")
	manufacturerSlug := testutil.RandomSlug("tf-test-manufacturer-fpt-color")
	deviceTypeName := testutil.RandomName("tf-test-device-type-fpt-color")
	deviceTypeSlug := testutil.RandomSlug("tf-test-device-type-fpt-color")
	rearPortName := testutil.RandomName("tf-test-rear-port-fpt-color")
	frontPortName := testutil.RandomName("tf-test-front-port-fpt-color")

	testutil.RunOptionalComputedFieldTestSuite(t, testutil.OptionalComputedFieldTestConfig{
		ResourceName:   "netbox_front_port_template",
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
	name        = "` + rearPortName + `"
	type        = "8p8c"
}

resource "netbox_front_port_template" "test" {
	device_type = netbox_device_type.test.id
	name        = "` + frontPortName + `"
	type        = "8p8c"
	rear_port   = netbox_rear_port_template.test.name
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
	name        = "` + rearPortName + `"
	type        = "8p8c"
}

resource "netbox_front_port_template" "test" {
	device_type = netbox_device_type.test.id
	name        = "` + frontPortName + `"
	type        = "8p8c"
	rear_port   = netbox_rear_port_template.test.name
	color       = "` + value + `"
}
`
		},
	})
}
