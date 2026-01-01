package resources_acceptance_tests

import (
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
)

// TestAccPowerPortTemplateResource_Label tests comprehensive scenarios for power port template label field.
// This validates that Optional+Computed string fields with empty string defaults work correctly.
func TestAccPowerPortTemplateResource_Label(t *testing.T) {
	// Generate unique names for this test run
	manufacturerName := testutil.RandomName("tf-test-mfr-pwr-port-tpl")
	manufacturerSlug := testutil.RandomSlug("tf-test-mfr-pwr-port-tpl")
	deviceTypeName := testutil.RandomName("tf-test-dev-type-pwr-port-tpl")
	deviceTypeSlug := testutil.RandomSlug("tf-test-dev-type-pwr-port-tpl")
	powerPortTemplateName := testutil.RandomName("tf-test-pwr-port-tpl")

	testutil.RunOptionalComputedFieldTestSuite(t, testutil.OptionalComputedFieldTestConfig{
		ResourceName:   "netbox_power_port_template",
		OptionalField:  "label",
		DefaultValue:   "",
		FieldTestValue: "Port-01",
		CheckDestroy: testutil.ComposeCheckDestroy(
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

resource "netbox_power_port_template" "test" {
	device_type = netbox_device_type.test.id
	name        = "` + powerPortTemplateName + `"
	type        = "iec-60320-c14"
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

resource "netbox_power_port_template" "test" {
	device_type = netbox_device_type.test.id
	name        = "` + powerPortTemplateName + `"
	type        = "iec-60320-c14"
	label       = "` + value + `"
}
`
		},
	})
}
