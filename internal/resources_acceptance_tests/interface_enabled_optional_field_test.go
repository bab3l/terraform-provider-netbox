package resources_acceptance_tests

import (
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccInterfaceResource_EnabledOptionalField(t *testing.T) {
	// Generate unique names for this test run
	siteName := testutil.RandomName("tf-test-site-if-enabled")
	siteSlug := testutil.RandomSlug("tf-test-site-if-enabled")
	manufacturerName := testutil.RandomName("tf-test-manufacturer-if-enabled")
	manufacturerSlug := testutil.RandomSlug("tf-test-manufacturer-if-enabled")
	deviceTypeName := testutil.RandomName("tf-test-device-type-if-enabled")
	deviceTypeSlug := testutil.RandomSlug("tf-test-device-type-if-enabled")
	deviceRoleName := testutil.RandomName("tf-test-device-role-if-enabled")
	deviceRoleSlug := testutil.RandomSlug("tf-test-device-role-if-enabled")
	deviceName := testutil.RandomName("tf-test-device-if-enabled")

	config := testutil.OptionalComputedFieldTestConfig{
		ResourceName:   "netbox_interface",
		OptionalField:  "enabled",
		DefaultValue:   "true",
		FieldTestValue: "false",
		BaseConfig: func() string {
			return `
			resource "netbox_site" "test" {
				name = "` + siteName + `"
				slug = "` + siteSlug + `"
			}

			resource "netbox_device_type" "test" {
				manufacturer = netbox_manufacturer.test.id
				model        = "` + deviceTypeName + `"
				slug         = "` + deviceTypeSlug + `"
			}

			resource "netbox_manufacturer" "test" {
				name = "` + manufacturerName + `"
				slug = "` + manufacturerSlug + `"
			}

			resource "netbox_device_role" "test" {
				name = "` + deviceRoleName + `"
				slug = "` + deviceRoleSlug + `"
			}

			resource "netbox_device" "test" {
				device_type = netbox_device_type.test.id
				role        = netbox_device_role.test.id
				site        = netbox_site.test.id
				name        = "` + deviceName + `"
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
				name = "` + siteName + `"
				slug = "` + siteSlug + `"
			}

			resource "netbox_device_type" "test" {
				manufacturer = netbox_manufacturer.test.id
				model        = "` + deviceTypeName + `"
				slug         = "` + deviceTypeSlug + `"
			}

			resource "netbox_manufacturer" "test" {
				name = "` + manufacturerName + `"
				slug = "` + manufacturerSlug + `"
			}

			resource "netbox_device_role" "test" {
				name = "` + deviceRoleName + `"
				slug = "` + deviceRoleSlug + `"
			}

			resource "netbox_device" "test" {
				device_type = netbox_device_type.test.id
				role        = netbox_device_role.test.id
				site        = netbox_site.test.id
				name        = "` + deviceName + `"
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
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckInterfaceDestroy,
			testutil.CheckDeviceDestroy,
			testutil.CheckDeviceTypeDestroy,
			testutil.CheckDeviceRoleDestroy,
			testutil.CheckSiteDestroy,
			testutil.CheckManufacturerDestroy,
		),
		Steps: steps,
	})
}

func TestAccInterfaceTemplateResource_EnabledOptionalField(t *testing.T) {
	// Generate unique names for this test run
	manufacturerName := testutil.RandomName("tf-test-manufacturer-ift-enabled")
	manufacturerSlug := testutil.RandomSlug("tf-test-manufacturer-ift-enabled")
	deviceTypeName := testutil.RandomName("tf-test-device-type-ift-enabled")
	deviceTypeSlug := testutil.RandomSlug("tf-test-device-type-ift-enabled")

	config := testutil.OptionalComputedFieldTestConfig{
		ResourceName:   "netbox_interface_template",
		OptionalField:  "enabled",
		DefaultValue:   "true",
		FieldTestValue: "false",
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
				name        = "eth0"
				type        = "1000base-t"
			}
			`
		},
		WithFieldConfig: func(fieldValue string) string {
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
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckDeviceTypeDestroy,
			testutil.CheckManufacturerDestroy,
		),
		Steps: steps,
	})
}

func TestAccInterfaceTemplateResource_LabelOptionalField(t *testing.T) {
	// Generate unique names for this test run
	manufacturerName := testutil.RandomName("tf-test-manufacturer-ift-label")
	manufacturerSlug := testutil.RandomSlug("tf-test-manufacturer-ift-label")
	deviceTypeName := testutil.RandomName("tf-test-device-type-ift-label")
	deviceTypeSlug := testutil.RandomSlug("tf-test-device-type-ift-label")

	config := testutil.OptionalFieldTestConfig{
		ResourceName:   "netbox_interface_template",
		OptionalField:  "label",
		FieldTestValue: "test-label",
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
				name        = "eth0"
				type        = "1000base-t"
			}
			`
		},
		WithFieldConfig: func(fieldValue string) string {
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
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckDeviceTypeDestroy,
			testutil.CheckManufacturerDestroy,
		),
		Steps: steps,
	})
}
