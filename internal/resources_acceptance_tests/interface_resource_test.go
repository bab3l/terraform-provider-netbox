package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccInterfaceResource_basic(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-interface")
	siteSlug := testutil.RandomSlug("site")
	mfrSlug := testutil.RandomSlug("mfr")
	deviceSlug := testutil.RandomSlug("device")
	roleSlug := testutil.RandomSlug("role")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterManufacturerCleanup(mfrSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceSlug)
	cleanup.RegisterDeviceRoleCleanup(roleSlug)
	cleanup.RegisterDeviceCleanup(name + "-device")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInterfaceResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_interface.test", "id"),
					resource.TestCheckResourceAttr("netbox_interface.test", "name", name),
					resource.TestCheckResourceAttr("netbox_interface.test", "type", "1000base-t"),
					resource.TestCheckResourceAttrSet("netbox_interface.test", "device"),
				),
			},
			{
				ResourceName:            "netbox_interface.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"device"},
			},
		},
	})
}

func TestAccInterfaceResource_update(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-interface-update")
	updatedName := testutil.RandomName("tf-test-interface-updated")
	siteSlug := testutil.RandomSlug("site")
	mfrSlug := testutil.RandomSlug("mfr")
	deviceSlug := testutil.RandomSlug("device")
	roleSlug := testutil.RandomSlug("role")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterManufacturerCleanup(mfrSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceSlug)
	cleanup.RegisterDeviceRoleCleanup(roleSlug)
	cleanup.RegisterDeviceCleanup(name + "-device")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInterfaceResourceConfig_forUpdate(name, testutil.Description1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_interface.test", "id"),
					resource.TestCheckResourceAttr("netbox_interface.test", "name", name),
					resource.TestCheckResourceAttr("netbox_interface.test", "type", "1000base-t"),
					resource.TestCheckResourceAttr("netbox_interface.test", "enabled", "true"),
					resource.TestCheckResourceAttr("netbox_interface.test", "mtu", "1500"),
					resource.TestCheckResourceAttr("netbox_interface.test", "description", testutil.Description1),
				),
			},
			{
				Config: testAccInterfaceResourceConfig_forUpdate(updatedName, testutil.Description2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_interface.test", "name", updatedName),
					resource.TestCheckResourceAttr("netbox_interface.test", "type", "10gbase-x-sfpp"),
					resource.TestCheckResourceAttr("netbox_interface.test", "enabled", "false"),
					resource.TestCheckResourceAttr("netbox_interface.test", "mtu", "9000"),
					resource.TestCheckResourceAttr("netbox_interface.test", "description", testutil.Description2),
				),
			},
		},
	})
}

func TestAccInterfaceResource_full(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-interface-full")
	siteSlug := testutil.RandomSlug("site")
	mfrSlug := testutil.RandomSlug("mfr")
	deviceSlug := testutil.RandomSlug("device")
	roleSlug := testutil.RandomSlug("role")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterManufacturerCleanup(mfrSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceSlug)
	cleanup.RegisterDeviceRoleCleanup(roleSlug)
	cleanup.RegisterDeviceCleanup(name + "-device")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInterfaceResourceConfig_full(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_interface.test", "id"),
					resource.TestCheckResourceAttr("netbox_interface.test", "name", name),
					resource.TestCheckResourceAttr("netbox_interface.test", "type", "1000base-t"),
					resource.TestCheckResourceAttr("netbox_interface.test", "enabled", "false"),
					resource.TestCheckResourceAttr("netbox_interface.test", "mtu", "1500"),
					resource.TestCheckResourceAttr("netbox_interface.test", "mgmt_only", "true"),
					resource.TestCheckResourceAttr("netbox_interface.test", "description", "Test interface with full options"),
				),
			},
		},
	})
}

func TestAccConsistency_Interface(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-interface-consistency")
	siteSlug := testutil.RandomSlug("site")
	mfrSlug := testutil.RandomSlug("mfr")
	deviceSlug := testutil.RandomSlug("device")
	roleSlug := testutil.RandomSlug("role")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterManufacturerCleanup(mfrSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceSlug)
	cleanup.RegisterDeviceRoleCleanup(roleSlug)
	cleanup.RegisterDeviceCleanup(name + "-device")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInterfaceResourceConfig_consistency_device_name(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_interface.test", "id"),
					resource.TestCheckResourceAttr("netbox_interface.test", "name", name),
					resource.TestCheckResourceAttr("netbox_interface.test", "type", "1000base-t"),
				),
			},
		},
	})
}

func TestAccConsistency_Interface_LiteralNames(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-interface-literal")
	siteSlug := testutil.RandomSlug("site")
	mfrSlug := testutil.RandomSlug("mfr")
	deviceSlug := testutil.RandomSlug("device")
	roleSlug := testutil.RandomSlug("role")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterManufacturerCleanup(mfrSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceSlug)
	cleanup.RegisterDeviceRoleCleanup(roleSlug)
	cleanup.RegisterDeviceCleanup(name + "-device")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInterfaceResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_interface.test", "id"),
					resource.TestCheckResourceAttr("netbox_interface.test", "name", name),
					resource.TestCheckResourceAttr("netbox_interface.test", "type", "1000base-t"),
				),
			},
		},
	})
}

// TestAccInterfaceResource_EnabledComprehensive tests comprehensive scenarios for interface enabled field.
// This validates that Optional+Computed boolean fields work correctly across all scenarios.
func TestAccInterfaceResource_EnabledComprehensive(t *testing.T) {
	siteName := testutil.RandomName("tf-test-site-interface")
	siteSlug := testutil.RandomSlug("tf-test-site-interface")
	manufacturerName := testutil.RandomName("tf-test-manufacturer-interface")
	manufacturerSlug := testutil.RandomSlug("tf-test-manufacturer-interface")
	deviceRoleName := testutil.RandomName("tf-test-device-role-interface")
	deviceRoleSlug := testutil.RandomSlug("tf-test-device-role-interface")
	deviceTypeName := testutil.RandomName("tf-test-device-type-interface")
	deviceTypeSlug := testutil.RandomSlug("tf-test-device-type-interface")
	deviceName := testutil.RandomName("tf-test-device-interface")
	interfaceName := testutil.RandomName("eth0-enabled-test")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)
	cleanup.RegisterDeviceRoleCleanup(deviceRoleSlug)
	cleanup.RegisterDeviceCleanup(deviceName)

	testutil.RunOptionalComputedFieldTestSuite(t, testutil.OptionalComputedFieldTestConfig{
		ResourceName:   "netbox_interface",
		OptionalField:  "enabled",
		DefaultValue:   "true",
		FieldTestValue: "false",
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckInterfaceDestroy,
			testutil.CheckDeviceDestroy,
			testutil.CheckDeviceTypeDestroy,
			testutil.CheckDeviceRoleDestroy,
			testutil.CheckSiteDestroy,
			testutil.CheckManufacturerDestroy,
		),
		BaseConfig: func() string {
			return testAccInterfaceResourceWithOptionalField(siteName, siteSlug, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceRoleName, deviceRoleSlug, deviceName, interfaceName, "enabled", "")
		},
		WithFieldConfig: func(value string) string {
			return testAccInterfaceResourceWithOptionalField(siteName, siteSlug, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceRoleName, deviceRoleSlug, deviceName, interfaceName, "enabled", value)
		},
	})
}

func testAccInterfaceResourceWithOptionalField(siteName, siteSlug, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceRoleName, deviceRoleSlug, deviceName, interfaceName, optionalFieldName, optionalFieldValue string) string {
	optionalField := ""
	if optionalFieldValue != "" {
		optionalField = fmt.Sprintf("\n  %s = %s", optionalFieldName, optionalFieldValue)
	}

	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name = %q
  slug = %q
}

resource "netbox_manufacturer" "test" {
  name = %q
  slug = %q
}

resource "netbox_device_role" "test" {
  name     = %q
  slug     = %q
  color    = "aa1409"
  vm_role  = false
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = %q
  slug         = %q
}

resource "netbox_device" "test" {
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
  name        = %q
}

resource "netbox_interface" "test" {
  device = netbox_device.test.id
  name   = %q
  type   = "1000base-t"%s
}
`, siteName, siteSlug, manufacturerName, manufacturerSlug, deviceRoleName, deviceRoleSlug, deviceTypeName, deviceTypeSlug, deviceName, interfaceName, optionalField)
}

func testAccInterfaceResourceConfig_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "netbox_interface" "test" {
  device = netbox_device.test.id
  name   = %q
  type   = "1000base-t"
}
`, testAccInterfaceResourcePrereqs(name), name)

}

func testAccInterfaceResourceConfig_forUpdate(name, description string) string {
	// Toggle between different types and settings based on description
	interfaceType := testutil.InterfaceType1000BaseT
	enabled := "true" //nolint:goconst // Boolean value specific to test configuration
	mtu := 1500

	if description == testutil.Description2 {
		interfaceType = testutil.InterfaceType10GBaseSFPP
		enabled = "false" //nolint:goconst // Boolean value specific to test configuration
		mtu = 9000
	}

	return fmt.Sprintf(`
%s

resource "netbox_interface" "test" {
  device      = netbox_device.test.id
  name        = %q
  type        = %q
  enabled     = %s
  mtu         = %d
  description = %q
}
`, testAccInterfaceResourcePrereqs(name), name, interfaceType, enabled, mtu, description)
}

func testAccInterfaceResourceConfig_full(name string) string {
	return fmt.Sprintf(`
%s

resource "netbox_interface" "test" {
  device       = netbox_device.test.id
  name         = %q
  type         = "1000base-t"
  enabled      = false
  mtu          = 1500
  mgmt_only    = true
  description  = "Test interface with full options"
}
`, testAccInterfaceResourcePrereqs(name), name)

}

func testAccInterfaceResourceConfig_consistency_device_name(name string) string {
	return fmt.Sprintf(`
%s

resource "netbox_interface" "test" {
  device = netbox_device.test.name
  name   = %q
  type   = "1000base-t"
}
`, testAccInterfaceResourcePrereqs(name), name)
}

func testAccInterfaceResourcePrereqs(name string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name = %q
  slug = %q
}

resource "netbox_manufacturer" "test" {
  name = %q
  slug = %q
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = %q
  slug         = %q
}

resource "netbox_device_role" "test" {
  name = %q
  slug = %q
}

resource "netbox_device" "test" {
  site        = netbox_site.test.id
  name        = %q
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  status      = "offline"
}
`, name+"-site", testutil.RandomSlug("site"), name+"-mfr", testutil.RandomSlug("mfr"), name+"-model", testutil.RandomSlug("device"), name+"-role", testutil.RandomSlug("role"), name+"-device")
}

func TestAccInterfaceResource_externalDeletion(t *testing.T) {
	t.Parallel()
	testutil.TestAccPreCheck(t)

	name := testutil.RandomName("tf-test-interface-ext-del")
	siteSlug := testutil.RandomSlug("site")
	mfrSlug := testutil.RandomSlug("mfr")
	deviceSlug := testutil.RandomSlug("device")
	roleSlug := testutil.RandomSlug("role")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterManufacturerCleanup(mfrSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceSlug)
	cleanup.RegisterDeviceRoleCleanup(roleSlug)
	cleanup.RegisterDeviceCleanup(name + "-device")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInterfaceResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_interface.test", "id"),
					resource.TestCheckResourceAttr("netbox_interface.test", "name", name),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					items, _, err := client.DcimAPI.DcimInterfacesList(context.Background()).NameIc([]string{name}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find interface for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.DcimAPI.DcimInterfacesDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete interface: %v", err)
					}
					t.Logf("Successfully externally deleted interface with ID: %d", itemID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccInterfaceResource_removeDescriptionAndLabel(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-interface-optional")
	siteSlug := testutil.RandomSlug("site")
	mfrSlug := testutil.RandomSlug("mfr")
	deviceSlug := testutil.RandomSlug("device")
	roleSlug := testutil.RandomSlug("role")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterManufacturerCleanup(mfrSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceSlug)
	cleanup.RegisterDeviceRoleCleanup(roleSlug)
	cleanup.RegisterDeviceCleanup(name + "-device")

	testutil.TestRemoveOptionalFields(t, testutil.MultiFieldOptionalTestConfig{
		ResourceName: "netbox_interface",
		BaseConfig: func() string {
			return testAccInterfaceResourceConfig_basic(name)
		},
		ConfigWithFields: func() string {
			return testAccInterfaceResourceConfig_withDescriptionAndLabel(
				name,
				"Test description",
				"Test label",
			)
		},
		OptionalFields: map[string]string{
			"description": "Test description",
			"label":       "Test label",
		},
		RequiredFields: map[string]string{
			"name": name,
			"type": "1000base-t",
		},
	})
}

func testAccInterfaceResourceConfig_withDescriptionAndLabel(name, description, label string) string {
	siteSlug := testutil.RandomSlug("site")
	mfrSlug := testutil.RandomSlug("mfr")
	deviceSlug := testutil.RandomSlug("device")
	roleSlug := testutil.RandomSlug("role")
	deviceName := name + "-device"

	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name = "Test Site"
  slug = %[2]q
}

resource "netbox_manufacturer" "test" {
  name = "Test Manufacturer"
  slug = %[3]q
}

resource "netbox_device_type" "test" {
  model        = "Test Device Type"
  slug         = %[4]q
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_device_role" "test" {
  name = "Test Role"
  slug = %[5]q
  color = "ff0000"
}

resource "netbox_device" "test" {
  name        = %[8]q
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
}

resource "netbox_interface" "test" {
  name        = %[1]q
  device      = netbox_device.test.id
  type        = "1000base-t"
  description = %[6]q
  label       = %[7]q
}
`, name, siteSlug, mfrSlug, deviceSlug, roleSlug, description, label, deviceName)
}

func TestAccInterfaceResource_removeOptionalFields(t *testing.T) {
	t.Parallel()

	interfaceName := testutil.RandomName("tf-test-iface-optional")
	siteName := testutil.RandomName("tf-test-site")
	siteSlug := testutil.RandomSlug("tf-test-site")
	manufacturerName := testutil.RandomName("tf-test-mfr")
	manufacturerSlug := testutil.RandomSlug("tf-test-mfr")
	deviceTypeName := testutil.RandomName("tf-test-devtype")
	deviceTypeSlug := testutil.RandomSlug("tf-test-devtype")
	roleName := testutil.RandomName("tf-test-role")
	roleSlug := testutil.RandomSlug("tf-test-role")
	deviceName := testutil.RandomName("tf-test-device")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)
	cleanup.RegisterDeviceRoleCleanup(roleSlug)
	cleanup.RegisterDeviceCleanup(deviceName)

	testutil.TestRemoveOptionalFields(t, testutil.MultiFieldOptionalTestConfig{
		ResourceName: "netbox_interface",
		BaseConfig: func() string {
			return testAccInterfaceResourceConfig_removeOptionalFields_base(
				interfaceName, siteName, siteSlug, manufacturerName, manufacturerSlug,
				deviceTypeName, deviceTypeSlug, roleName, roleSlug, deviceName,
			)
		},
		ConfigWithFields: func() string {
			return testAccInterfaceResourceConfig_removeOptionalFields_withFields(
				interfaceName, siteName, siteSlug, manufacturerName, manufacturerSlug,
				deviceTypeName, deviceTypeSlug, roleName, roleSlug, deviceName,
			)
		},
		OptionalFields: map[string]string{
			"duplex":      "full",
			"label":       "Test Label",
			"mac_address": "00:11:22:33:44:55",
			"mode":        "access",
			"mtu":         "1500",
			"speed":       "1000000",
		},
		RequiredFields: map[string]string{
			"name": interfaceName,
		},
		CheckDestroy: testutil.CheckInterfaceDestroy,
	})
}

func testAccInterfaceResourceConfig_removeOptionalFields_base(interfaceName, siteName, siteSlug, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, roleName, roleSlug, deviceName string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name = %[2]q
  slug = %[3]q
}

resource "netbox_manufacturer" "test" {
  name = %[4]q
  slug = %[5]q
}

resource "netbox_device_type" "test" {
  model        = %[6]q
  slug         = %[7]q
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_device_role" "test" {
  name  = %[8]q
  slug  = %[9]q
  color = "ff0000"
}

resource "netbox_device" "test" {
  name        = %[10]q
  site        = netbox_site.test.id
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
}

resource "netbox_interface" "test" {
  name   = %[1]q
  device = netbox_device.test.id
  type   = "1000base-t"
}
`, interfaceName, siteName, siteSlug, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, roleName, roleSlug, deviceName)
}

func testAccInterfaceResourceConfig_removeOptionalFields_withFields(interfaceName, siteName, siteSlug, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, roleName, roleSlug, deviceName string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name = %[2]q
  slug = %[3]q
}

resource "netbox_manufacturer" "test" {
  name = %[4]q
  slug = %[5]q
}

resource "netbox_device_type" "test" {
  model        = %[6]q
  slug         = %[7]q
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_device_role" "test" {
  name  = %[8]q
  slug  = %[9]q
  color = "ff0000"
}

resource "netbox_device" "test" {
  name        = %[10]q
  site        = netbox_site.test.id
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
}

resource "netbox_interface" "test" {
  name        = %[1]q
  device      = netbox_device.test.id
  type        = "1000base-t"
  duplex      = "full"
  label       = "Test Label"
  mac_address = "00:11:22:33:44:55"
  mode        = "access"
  mtu         = 1500
  speed       = 1000000
}
`, interfaceName, siteName, siteSlug, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, roleName, roleSlug, deviceName)
}

// NOTE: Custom field tests for interface resource are in resources_acceptance_tests_customfields package
