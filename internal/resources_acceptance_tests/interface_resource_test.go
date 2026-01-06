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

// NOTE: Custom field tests for interface resource are in resources_acceptance_tests_customfields package
