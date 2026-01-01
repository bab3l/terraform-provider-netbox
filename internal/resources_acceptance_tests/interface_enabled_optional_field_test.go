package resources_acceptance_tests

import (
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccInterfaceResource_EnabledOptionalField(t *testing.T) {
	config := testutil.OptionalComputedFieldTestConfig{
		ResourceName:   "netbox_interface",
		OptionalField:  "enabled",
		DefaultValue:   "true",
		FieldTestValue: "false",
		BaseConfig: func() string {
			return `
			resource "netbox_site" "test" {
				name = "test-site"
				slug = "test-site"
			}

			resource "netbox_device_type" "test" {
				manufacturer = netbox_manufacturer.test.id
				model        = "test-device-type"
				slug         = "test-device-type"
			}

			resource "netbox_manufacturer" "test" {
				name = "test-manufacturer"
				slug = "test-manufacturer"
			}

			resource "netbox_device_role" "test" {
				name = "test-device-role"
				slug = "test-device-role"
			}

			resource "netbox_device" "test" {
				device_type = netbox_device_type.test.id
				role        = netbox_device_role.test.id
				site        = netbox_site.test.id
				name        = "test-device"
			}

			resource "netbox_interface" "test" {
				device = netbox_device.test.id
				name   = "eth0"
				type   = "1000base-t"
			}
			`
		},
		WithFieldConfig: func(fieldValue string) string {
			return `
			resource "netbox_site" "test" {
				name = "test-site"
				slug = "test-site"
			}

			resource "netbox_device_type" "test" {
				manufacturer = netbox_manufacturer.test.id
				model        = "test-device-type"
				slug         = "test-device-type"
			}

			resource "netbox_manufacturer" "test" {
				name = "test-manufacturer"
				slug = "test-manufacturer"
			}

			resource "netbox_device_role" "test" {
				name = "test-device-role"
				slug = "test-device-role"
			}

			resource "netbox_device" "test" {
				device_type = netbox_device_type.test.id
				role        = netbox_device_role.test.id
				site        = netbox_site.test.id
				name        = "test-device"
			}

			resource "netbox_interface" "test" {
				device  = netbox_device.test.id
				name    = "eth0"
				type    = "1000base-t"
				enabled = ` + fieldValue + `
			}
			`
		},
	}

	steps := testutil.GenerateOptionalComputedFieldTests(t, config)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps:                    steps,
	})
}

func TestAccInterfaceTemplateResource_EnabledOptionalField(t *testing.T) {
	config := testutil.OptionalComputedFieldTestConfig{
		ResourceName:   "netbox_interface_template",
		OptionalField:  "enabled",
		DefaultValue:   "true",
		FieldTestValue: "false",
		BaseConfig: func() string {
			return `
			resource "netbox_manufacturer" "test" {
				name = "test-manufacturer"
				slug = "test-manufacturer"
			}

			resource "netbox_device_type" "test" {
				manufacturer = netbox_manufacturer.test.id
				model        = "test-device-type"
				slug         = "test-device-type"
			}

			resource "netbox_interface_template" "test" {
				device_type = netbox_device_type.test.id
				name        = "eth0"
				type        = "1000base-t"
			}
			`
		},
		WithFieldConfig: func(fieldValue string) string {
			return `
			resource "netbox_manufacturer" "test" {
				name = "test-manufacturer"
				slug = "test-manufacturer"
			}

			resource "netbox_device_type" "test" {
				manufacturer = netbox_manufacturer.test.id
				model        = "test-device-type"
				slug         = "test-device-type"
			}

			resource "netbox_interface_template" "test" {
				device_type = netbox_device_type.test.id
				name        = "eth0"
				type        = "1000base-t"
				enabled     = ` + fieldValue + `
			}
			`
		},
	}

	steps := testutil.GenerateOptionalComputedFieldTests(t, config)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps:                    steps,
	})
}

func TestAccInterfaceTemplateResource_LabelOptionalField(t *testing.T) {
	config := testutil.OptionalFieldTestConfig{
		ResourceName:   "netbox_interface_template",
		OptionalField:  "label",
		FieldTestValue: "test-label",
		BaseConfig: func() string {
			return `
			resource "netbox_manufacturer" "test" {
				name = "test-manufacturer"
				slug = "test-manufacturer"
			}

			resource "netbox_device_type" "test" {
				manufacturer = netbox_manufacturer.test.id
				model        = "test-device-type"
				slug         = "test-device-type"
			}

			resource "netbox_interface_template" "test" {
				device_type = netbox_device_type.test.id
				name        = "eth0"
				type        = "1000base-t"
			}
			`
		},
		WithFieldConfig: func(fieldValue string) string {
			return `
			resource "netbox_manufacturer" "test" {
				name = "test-manufacturer"
				slug = "test-manufacturer"
			}

			resource "netbox_device_type" "test" {
				manufacturer = netbox_manufacturer.test.id
				model        = "test-device-type"
				slug         = "test-device-type"
			}

			resource "netbox_interface_template" "test" {
				device_type = netbox_device_type.test.id
				name        = "eth0"
				type        = "1000base-t"
				label       = "` + fieldValue + `"
			}
			`
		},
	}

	steps := testutil.GenerateOptionalFieldTests(t, config)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps:                    steps,
	})
}
