package resources_acceptance_tests

import (
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
)

// TestAccInterfaceTemplateResource_EnabledComprehensive tests comprehensive scenarios for interface template enabled field.
// This validates that Optional+Computed boolean fields work correctly across all scenarios.
func TestAccInterfaceTemplateResource_EnabledComprehensive(t *testing.T) {
	// Generate unique names for this test run
	manufacturerName := testutil.RandomName("tf-test-mfr-int-tpl")
	manufacturerSlug := testutil.RandomSlug("tf-test-mfr-int-tpl")
	deviceTypeName := testutil.RandomName("tf-test-dev-type-int-tpl")
	deviceTypeSlug := testutil.RandomSlug("tf-test-dev-type-int-tpl")
	interfaceTemplateName := testutil.RandomName("tf-test-int-tpl")

	testutil.RunOptionalComputedFieldTestSuite(t, testutil.OptionalComputedFieldTestConfig{
		ResourceName:   "netbox_interface_template",
		OptionalField:  "enabled",
		DefaultValue:   "true",
		FieldTestValue: "false",
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

resource "netbox_interface_template" "test" {
	device_type = netbox_device_type.test.id
	name        = "` + interfaceTemplateName + `"
	type        = "1000base-t"
	# enabled field intentionally omitted - should get default true
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

resource "netbox_interface_template" "test" {
	device_type = netbox_device_type.test.id
	name        = "` + interfaceTemplateName + `"
	type        = "1000base-t"
	enabled     = ` + value + `
}
`
		},
	})
}
