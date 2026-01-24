package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const configTemplateCode = "{{ device.name }}"

func TestAccDeviceResource_basic(t *testing.T) {
	t.Parallel()

	// Generate unique names to avoid conflicts between test runs
	deviceName := testutil.RandomName("tf-test-device")
	manufacturerName := testutil.RandomName("tf-test-manufacturer")
	manufacturerSlug := testutil.RandomSlug("tf-test-mfr")
	deviceTypeModel := testutil.RandomName("tf-test-device-type")
	deviceTypeSlug := testutil.RandomSlug("tf-test-dt")
	deviceRoleName := testutil.RandomName("tf-test-device-role")
	deviceRoleSlug := testutil.RandomSlug("tf-test-dr")
	siteName := testutil.RandomName("tf-test-site")
	siteSlug := testutil.RandomSlug("tf-test-site")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
	cleanup.RegisterDeviceRoleCleanup(deviceRoleSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)
	cleanup.RegisterDeviceCleanup(deviceName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDeviceResourceConfig_basic(deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_device.test", "id"),
					resource.TestCheckResourceAttr("netbox_device.test", "name", deviceName),
					resource.TestCheckResourceAttrSet("netbox_device.test", "device_type"),
					resource.TestCheckResourceAttrSet("netbox_device.test", "role"),
					resource.TestCheckResourceAttrSet("netbox_device.test", "site"),
					resource.TestCheckResourceAttr("netbox_device.test", "status", "active"),
				),
			},
			{
				ResourceName:      "netbox_device.test",
				ImportState:       true,
				ImportStateVerify: true,
				// Note: some fields may use slugs in config but IDs in state after import
				ImportStateVerifyIgnore: []string{"device_type", "role", "site"},
			},
			{
				Config:   testAccDeviceResourceConfig_basic(deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug),
				PlanOnly: true,
			},
		},
	})
}

func TestAccDeviceResource_import(t *testing.T) {
	t.Parallel()

	deviceName := testutil.RandomName("tf-test-device-import")
	manufacturerName := testutil.RandomName("tf-test-manufacturer-import")
	manufacturerSlug := testutil.RandomSlug("tf-test-mfr-import")
	deviceTypeModel := testutil.RandomName("tf-test-device-type-import")
	deviceTypeSlug := testutil.RandomSlug("tf-test-dt-import")
	deviceRoleName := testutil.RandomName("tf-test-device-role-import")
	deviceRoleSlug := testutil.RandomSlug("tf-test-dr-import")
	siteName := testutil.RandomName("tf-test-site-import")
	siteSlug := testutil.RandomSlug("tf-test-site-import")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
	cleanup.RegisterDeviceRoleCleanup(deviceRoleSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)
	cleanup.RegisterDeviceCleanup(deviceName)

	testutil.RunImportTest(t, testutil.ImportTestConfig{
		ResourceName: "netbox_device",
		Config: func() string {
			return testAccDeviceResourceConfig_basic(
				deviceName,
				manufacturerName,
				manufacturerSlug,
				deviceTypeModel,
				deviceTypeSlug,
				deviceRoleName,
				deviceRoleSlug,
				siteName,
				siteSlug,
			)
		},
		ImportStateVerifyIgnore: []string{"device_type", "role", "site"},
		AdditionalChecks:        testutil.ValidateReferenceIDs("netbox_device.test", "device_type", "role", "site"),
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckDeviceDestroy,
			testutil.CheckDeviceTypeDestroy,
			testutil.CheckDeviceRoleDestroy,
			testutil.CheckSiteDestroy,
			testutil.CheckManufacturerDestroy,
		),
	})
}

func TestAccDeviceResource_update(t *testing.T) {
	t.Parallel()

	deviceName := testutil.RandomName("tf-test-device-update")
	updatedDeviceName := testutil.RandomName("tf-test-device-updated")
	manufacturerName := testutil.RandomName("tf-test-manufacturer")
	manufacturerSlug := testutil.RandomSlug("tf-test-mfr")
	deviceTypeModel := testutil.RandomName("tf-test-device-type")
	deviceTypeSlug := testutil.RandomSlug("tf-test-dt")
	deviceRoleName := testutil.RandomName("tf-test-device-role")
	deviceRoleSlug := testutil.RandomSlug("tf-test-dr")
	siteName := testutil.RandomName("tf-test-site")
	siteSlug := testutil.RandomSlug("tf-test-site")
	serial := testutil.RandomName("SN")
	updatedSerial := testutil.RandomName("SN-UPD")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
	cleanup.RegisterDeviceRoleCleanup(deviceRoleSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)
	cleanup.RegisterDeviceCleanup(deviceName)
	cleanup.RegisterDeviceCleanup(updatedDeviceName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDeviceResourceConfig_forUpdate(deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug, serial, testutil.Description1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_device.test", "id"),
					resource.TestCheckResourceAttr("netbox_device.test", "name", deviceName),
					resource.TestCheckResourceAttr("netbox_device.test", "status", "active"),
					resource.TestCheckResourceAttr("netbox_device.test", "serial", serial),
					resource.TestCheckResourceAttr("netbox_device.test", "description", testutil.Description1),
				),
			},
			{
				Config: testAccDeviceResourceConfig_forUpdate(updatedDeviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug, updatedSerial, testutil.Description2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_device.test", "name", updatedDeviceName),
					resource.TestCheckResourceAttr("netbox_device.test", "status", "planned"),
					resource.TestCheckResourceAttr("netbox_device.test", "serial", updatedSerial),
					resource.TestCheckResourceAttr("netbox_device.test", "description", testutil.Description2),
				),
			},
		},
	})
}

func TestAccDeviceResource_full(t *testing.T) {

	t.Parallel()
	// Generate unique names
	deviceName := testutil.RandomName("tf-test-device")
	manufacturerName := testutil.RandomName("tf-test-manufacturer")
	manufacturerSlug := testutil.RandomSlug("tf-test-mfr")
	deviceTypeModel := testutil.RandomName("tf-test-device-type")
	deviceTypeSlug := testutil.RandomSlug("tf-test-dt")
	deviceRoleName := testutil.RandomName("tf-test-device-role")
	deviceRoleSlug := testutil.RandomSlug("tf-test-dr")
	siteName := testutil.RandomName("tf-test-site")
	siteSlug := testutil.RandomSlug("tf-test-site")
	serial := testutil.RandomName("SN")
	assetTag := testutil.RandomName("AT")
	clusterTypeName := testutil.RandomName("tf-test-cluster-type")
	clusterTypeSlug := testutil.RandomSlug("tf-test-cluster-type")
	clusterName := testutil.RandomName("tf-test-cluster")
	configTemplateName := testutil.RandomName("tf-test-config-template")
	configTemplateCode := configTemplateCode

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)
	cleanup.RegisterDeviceRoleCleanup(deviceRoleSlug)
	cleanup.RegisterDeviceCleanup(deviceName)
	cleanup.RegisterClusterTypeCleanup(clusterTypeSlug)
	cleanup.RegisterClusterCleanup(clusterName)
	cleanup.RegisterConfigTemplateCleanup(configTemplateName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDeviceResourceConfig_full(deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug, serial, assetTag, clusterTypeName, clusterTypeSlug, clusterName, configTemplateName, configTemplateCode),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_device.test", "id"),
					resource.TestCheckResourceAttr("netbox_device.test", "name", deviceName),
					resource.TestCheckResourceAttrSet("netbox_device.test", "device_type"),
					resource.TestCheckResourceAttrSet("netbox_device.test", "role"),
					resource.TestCheckResourceAttrSet("netbox_device.test", "site"),
					resource.TestCheckResourceAttr("netbox_device.test", "status", "planned"),
					resource.TestCheckResourceAttr("netbox_device.test", "serial", serial),
					resource.TestCheckResourceAttr("netbox_device.test", "asset_tag", assetTag),
					resource.TestCheckResourceAttr("netbox_device.test", "description", "Test device description"),
					resource.TestCheckResourceAttr("netbox_device.test", "comments", "Test device comments"),
					resource.TestCheckResourceAttr("netbox_device.test", "airflow", "front-to-rear"),
					resource.TestCheckResourceAttrSet("netbox_device.test", "cluster"),
					resource.TestCheckResourceAttrSet("netbox_device.test", "config_template"),
					resource.TestCheckResourceAttr("netbox_device.test", "tags.#", "0"),
				),
			},
		},
	})
}

// TestAccDeviceResource_StatusOptionalField tests comprehensive scenarios for the device status optional field.
// This validates that Optional+Computed fields work correctly across all scenarios.
func TestAccDeviceResource_StatusOptionalField(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-site-device-status")
	siteSlug := testutil.RandomSlug("tf-test-site-device-status")
	manufacturerName := testutil.RandomName("tf-test-manufacturer-device-status")
	manufacturerSlug := testutil.RandomSlug("tf-test-manufacturer-device-status")
	deviceRoleName := testutil.RandomName("tf-test-device-role-status")
	deviceRoleSlug := testutil.RandomSlug("tf-test-device-role-status")
	deviceTypeName := testutil.RandomName("tf-test-device-type-status")
	deviceTypeSlug := testutil.RandomSlug("tf-test-device-type-status")
	deviceName := testutil.RandomName("tf-test-device-status")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)
	cleanup.RegisterDeviceRoleCleanup(deviceRoleSlug)
	cleanup.RegisterDeviceCleanup(deviceName)

	testutil.RunOptionalComputedFieldTestSuite(t, testutil.OptionalComputedFieldTestConfig{
		ResourceName:   "netbox_device",
		OptionalField:  "status",
		DefaultValue:   "active",
		FieldTestValue: "planned",
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckDeviceDestroy,
			testutil.CheckDeviceTypeDestroy,
			testutil.CheckDeviceRoleDestroy,
			testutil.CheckSiteDestroy,
			testutil.CheckManufacturerDestroy,
		),
		BaseConfig: func() string {
			return testAccDeviceResourceWithOptionalField(siteName, siteSlug, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceRoleName, deviceRoleSlug, deviceName, "status", "")
		},
		WithFieldConfig: func(value string) string {
			return testAccDeviceResourceWithOptionalField(siteName, siteSlug, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceRoleName, deviceRoleSlug, deviceName, "status", value)
		},
	})
}

func testAccDeviceResourceWithOptionalField(siteName, siteSlug, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceRoleName, deviceRoleSlug, deviceName, optionalFieldName, optionalFieldValue string) string {
	optionalField := ""
	if optionalFieldValue != "" {
		optionalField = fmt.Sprintf("\n  %s = %q", optionalFieldName, optionalFieldValue)
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
  name    = %q
  slug    = %q
  color   = "aa1409"
  vm_role = false
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
  name        = %q%s
}
`, siteName, siteSlug, manufacturerName, manufacturerSlug, deviceRoleName, deviceRoleSlug, deviceTypeName, deviceTypeSlug, deviceName, optionalField)
}

func TestAccConsistency_Device(t *testing.T) {
	t.Parallel()

	deviceName := testutil.RandomName("device")
	deviceTypeName := testutil.RandomName("device-type")
	deviceTypeSlug := testutil.RandomSlug("device-type")
	manufacturerName := testutil.RandomName("manufacturer")
	manufacturerSlug := testutil.RandomSlug("manufacturer")
	roleName := testutil.RandomName("role")
	roleSlug := testutil.RandomSlug("role")
	siteName := testutil.RandomName("site")
	siteSlug := testutil.RandomSlug("site")
	tenantName := testutil.RandomName("tenant")
	tenantSlug := testutil.RandomSlug("tenant")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterDeviceRoleCleanup(roleSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDeviceConsistencyConfig(deviceName, deviceTypeName, deviceTypeSlug, manufacturerName, manufacturerSlug, roleName, roleSlug, siteName, siteSlug, tenantName, tenantSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_device.test", "name", deviceName),
					resource.TestCheckResourceAttrSet("netbox_device.test", "device_type"),
					resource.TestCheckResourceAttrSet("netbox_device.test", "role"),
					resource.TestCheckResourceAttrSet("netbox_device.test", "site"),
					resource.TestCheckResourceAttrSet("netbox_device.test", "tenant"),
				),
			},
			{
				PlanOnly: true,
				Config:   testAccDeviceConsistencyConfig(deviceName, deviceTypeName, deviceTypeSlug, manufacturerName, manufacturerSlug, roleName, roleSlug, siteName, siteSlug, tenantName, tenantSlug),
			},
		},
	})
}

func TestAccConsistency_Device_LiteralNames(t *testing.T) {
	t.Parallel()

	deviceName := testutil.RandomName("device")
	deviceTypeName := testutil.RandomName("device-type")
	deviceTypeSlug := testutil.RandomSlug("device-type")
	manufacturerName := testutil.RandomName("manufacturer")
	manufacturerSlug := testutil.RandomSlug("manufacturer")
	roleName := testutil.RandomName("role")
	roleSlug := testutil.RandomSlug("role")
	siteName := testutil.RandomName("site")
	siteSlug := testutil.RandomSlug("site")
	tenantName := testutil.RandomName("tenant")
	tenantSlug := testutil.RandomSlug("tenant")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterDeviceRoleCleanup(roleSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDeviceConsistencyLiteralNamesConfig(deviceName, deviceTypeName, deviceTypeSlug, manufacturerName, manufacturerSlug, roleName, roleSlug, siteName, siteSlug, tenantName, tenantSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_device.test", "name", deviceName),
					resource.TestCheckResourceAttr("netbox_device.test", "device_type", deviceTypeSlug),
					resource.TestCheckResourceAttr("netbox_device.test", "role", roleSlug),
					resource.TestCheckResourceAttr("netbox_device.test", "site", siteName),
					resource.TestCheckResourceAttr("netbox_device.test", "tenant", tenantName),
				),
			},
			{
				// Critical: Verify no drift when refreshing state
				PlanOnly: true,
				Config:   testAccDeviceConsistencyLiteralNamesConfig(deviceName, deviceTypeName, deviceTypeSlug, manufacturerName, manufacturerSlug, roleName, roleSlug, siteName, siteSlug, tenantName, tenantSlug),
			},
		},
	})
}

// Helper functions to generate test configurations.
func testAccDeviceResourceConfig_basic(deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug string) string {
	return fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = %[3]q
  slug         = %[4]q
}

resource "netbox_device_role" "test" {
  name = %[5]q
  slug = %[6]q
}

resource "netbox_site" "test" {
  name   = %[7]q
  slug   = %[8]q
  status = "active"
}

resource "netbox_device" "test" {
  name        = %[9]q
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
}
`, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug, deviceName)
}

func testAccDeviceResourceConfig_forUpdate(deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug, serial, description string) string {
	// Toggle status between active and planned to test updates
	status := testutil.StatusActive
	if description == testutil.Description2 {
		status = testutil.StatusPlanned
	}

	return fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = %[3]q
  slug         = %[4]q
}

resource "netbox_device_role" "test" {
  name = %[5]q
  slug = %[6]q
}

resource "netbox_site" "test" {
  name   = %[7]q
  slug   = %[8]q
  status = "active"
}

resource "netbox_device" "test" {
  name        = %[9]q
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
  status      = %[12]q
  serial      = %[10]q
  description = %[11]q
}
`, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug, deviceName, serial, description, status)
}

func testAccDeviceResourceConfig_full(deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug, serial, assetTag, clusterTypeName, clusterTypeSlug, clusterName, configTemplateName, configTemplateCode string) string {
	return fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = %[3]q
  slug         = %[4]q
}

resource "netbox_device_role" "test" {
  name = %[5]q
  slug = %[6]q
}

resource "netbox_site" "test" {
  name   = %[7]q
  slug   = %[8]q
  status = "active"
}

resource "netbox_cluster_type" "test" {
	name = %[12]q
	slug = %[13]q
}

resource "netbox_cluster" "test" {
	name = %[14]q
	type = netbox_cluster_type.test.id
}

resource "netbox_config_template" "test" {
	name          = %[15]q
	template_code = %[16]q
}

resource "netbox_device" "test" {
  name        = %[9]q
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
  status      = "planned"
  serial      = %[10]q
  asset_tag   = %[11]q
  description = "Test device description"
  comments    = "Test device comments"
  airflow     = "front-to-rear"
	cluster     = netbox_cluster.test.id
	config_template = netbox_config_template.test.id
  tags        = []
}
`, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug, deviceName, serial, assetTag, clusterTypeName, clusterTypeSlug, clusterName, configTemplateName, configTemplateCode)
}

func testAccDeviceConsistencyConfig(deviceName, deviceTypeName, deviceTypeSlug, manufacturerName, manufacturerSlug, roleName, roleSlug, siteName, siteSlug, tenantName, tenantSlug string) string {
	return fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
  name = "%[4]s"
  slug = "%[5]s"
}

resource "netbox_device_type" "test" {
  model = "%[2]s"
  slug = "%[3]s"
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_device_role" "test" {
  name = "%[6]s"
  slug = "%[7]s"
}

resource "netbox_site" "test" {
  name = "%[8]s"
  slug = "%[9]s"
}

resource "netbox_tenant" "test" {
  name = "%[10]s"
  slug = "%[11]s"
}

resource "netbox_device" "test" {
  name = "%[1]s"
  device_type = netbox_device_type.test.id
  role = netbox_device_role.test.id
  site = netbox_site.test.id
  tenant = netbox_tenant.test.id
}
`, deviceName, deviceTypeName, deviceTypeSlug, manufacturerName, manufacturerSlug, roleName, roleSlug, siteName, siteSlug, tenantName, tenantSlug)
}

func testAccDeviceConsistencyLiteralNamesConfig(deviceName, deviceTypeName, deviceTypeSlug, manufacturerName, manufacturerSlug, roleName, roleSlug, siteName, siteSlug, tenantName, tenantSlug string) string {
	return fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
  name = "%[4]s"
  slug = "%[5]s"
}

resource "netbox_device_type" "test" {
  model = "%[2]s"
  slug = "%[3]s"
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_device_role" "test" {
  name = "%[6]s"
  slug = "%[7]s"
}

resource "netbox_site" "test" {
  name = "%[8]s"
  slug = "%[9]s"
}

resource "netbox_tenant" "test" {
  name = "%[10]s"
  slug = "%[11]s"
}

resource "netbox_device" "test" {
  name = "%[1]s"
  # Use literal string names to mimic existing user state
  device_type = "%[3]s"
  role = "%[7]s"
  site = "%[8]s"
  tenant = "%[10]s"

  depends_on = [netbox_device_type.test, netbox_device_role.test, netbox_site.test, netbox_tenant.test]
}
`, deviceName, deviceTypeName, deviceTypeSlug, manufacturerName, manufacturerSlug, roleName, roleSlug, siteName, siteSlug, tenantName, tenantSlug)
}

func TestAccDeviceResource_externalDeletion(t *testing.T) {
	t.Parallel()

	testutil.TestAccPreCheck(t)
	deviceName := testutil.RandomName("tf-test-device-ext-del")
	manufacturerName := testutil.RandomName("tf-test-manufacturer")
	manufacturerSlug := testutil.RandomSlug("tf-test-mfr")
	deviceTypeModel := testutil.RandomName("tf-test-device-type")
	deviceTypeSlug := testutil.RandomSlug("tf-test-dt")
	deviceRoleName := testutil.RandomName("tf-test-device-role")
	deviceRoleSlug := testutil.RandomSlug("tf-test-dr")
	siteName := testutil.RandomName("tf-test-site")
	siteSlug := testutil.RandomSlug("tf-test-site")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterDeviceCleanup(deviceName)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
	cleanup.RegisterDeviceRoleCleanup(deviceRoleSlug)
	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDeviceResourceConfig_basic(deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_device.test", "id"),
					resource.TestCheckResourceAttr("netbox_device.test", "name", deviceName),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					items, _, err := client.DcimAPI.DcimDevicesList(context.Background()).NameIc([]string{deviceName}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find device for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.DcimAPI.DcimDevicesDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete device: %v", err)
					}
					t.Logf("Successfully externally deleted device with ID: %d", itemID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccDeviceResource_removeDescriptionAndComments(t *testing.T) {
	t.Parallel()

	deviceName := testutil.RandomName("tf-test-device-optional")
	manufacturerName := testutil.RandomName("tf-test-manufacturer")
	manufacturerSlug := testutil.RandomSlug("tf-test-mfr")
	deviceTypeModel := testutil.RandomName("tf-test-device-type")
	deviceTypeSlug := testutil.RandomSlug("tf-test-dt")
	deviceRoleName := testutil.RandomName("tf-test-device-role")
	deviceRoleSlug := testutil.RandomSlug("tf-test-dr")
	siteName := testutil.RandomName("tf-test-site")
	siteSlug := testutil.RandomSlug("tf-test-site")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterDeviceCleanup(deviceName)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
	cleanup.RegisterDeviceRoleCleanup(deviceRoleSlug)
	cleanup.RegisterSiteCleanup(siteSlug)

	testutil.TestRemoveOptionalFields(t, testutil.MultiFieldOptionalTestConfig{
		ResourceName: "netbox_device",
		BaseConfig: func() string {
			return testAccDeviceResourceConfig_minimal(
				deviceName,
				manufacturerName,
				manufacturerSlug,
				deviceTypeModel,
				deviceTypeSlug,
				deviceRoleName,
				deviceRoleSlug,
				siteName,
				siteSlug,
			)
		},
		ConfigWithFields: func() string {
			return testAccDeviceResourceConfig_withDescriptionAndComments(
				deviceName,
				manufacturerName,
				manufacturerSlug,
				deviceTypeModel,
				deviceTypeSlug,
				deviceRoleName,
				deviceRoleSlug,
				siteName,
				siteSlug,
				"Test description",
				"Test comments",
			)
		},
		OptionalFields: map[string]string{
			"description": "Test description",
			"comments":    "Test comments",
		},
		RequiredFields: map[string]string{
			"name": deviceName,
		},
		CheckDestroy: testutil.CheckDeviceDestroy,
	})
}

func testAccDeviceResourceConfig_minimal(deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug string) string {
	return fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
  name = %[2]q
  slug = %[3]q
}

resource "netbox_device_type" "test" {
  model         = %[4]q
  slug          = %[5]q
  manufacturer  = netbox_manufacturer.test.id
}

resource "netbox_device_role" "test" {
  name = %[6]q
  slug = %[7]q
  color = "aa1409"
}

resource "netbox_site" "test" {
  name   = %[8]q
  slug   = %[9]q
  status = "active"
}

resource "netbox_device" "test" {
  name        = %[1]q
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
  status      = "active"
}
`, deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug)
}

func testAccDeviceResourceConfig_withDescriptionAndComments(deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug, description, comments string) string {
	return fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
  name = %[2]q
  slug = %[3]q
}

resource "netbox_device_type" "test" {
  model         = %[4]q
  slug          = %[5]q
  manufacturer  = netbox_manufacturer.test.id
}

resource "netbox_device_role" "test" {
  name = %[6]q
  slug = %[7]q
  color = "aa1409"
}

resource "netbox_site" "test" {
  name   = %[8]q
  slug   = %[9]q
  status = "active"
}

resource "netbox_device" "test" {
  name        = %[1]q
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
  status      = "active"
  description = %[10]q
  comments    = %[11]q
}
`, deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug, description, comments)
}

func TestAccDeviceResource_removeOptionalFields(t *testing.T) {
	t.Parallel()

	deviceName := testutil.RandomName("tf-test-device-optional")
	siteName := testutil.RandomName("tf-test-site")
	siteSlug := testutil.RandomSlug("tf-test-site")
	manufacturerName := testutil.RandomName("tf-test-mfr")
	manufacturerSlug := testutil.RandomSlug("tf-test-mfr")
	deviceTypeName := testutil.RandomName("tf-test-devtype")
	deviceTypeSlug := testutil.RandomSlug("tf-test-devtype")
	roleName := testutil.RandomName("tf-test-role")
	roleSlug := testutil.RandomSlug("tf-test-role")
	clusterTypeName := testutil.RandomName("tf-test-cluster-type")
	clusterTypeSlug := testutil.RandomSlug("tf-test-cluster-type")
	clusterName := testutil.RandomName("tf-test-cluster")
	configTemplateName := testutil.RandomName("tf-test-config-template")
	configTemplateCode := "{{ device.name }}"

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)
	cleanup.RegisterDeviceRoleCleanup(roleSlug)
	cleanup.RegisterDeviceCleanup(deviceName)
	cleanup.RegisterClusterTypeCleanup(clusterTypeSlug)
	cleanup.RegisterClusterCleanup(clusterName)
	cleanup.RegisterConfigTemplateCleanup(configTemplateName)

	testutil.TestRemoveOptionalFields(t, testutil.MultiFieldOptionalTestConfig{
		ResourceName: "netbox_device",
		BaseConfig: func() string {
			return testAccDeviceResourceConfig_removeOptionalFields_base(
				deviceName, siteName, siteSlug, manufacturerName, manufacturerSlug,
				deviceTypeName, deviceTypeSlug, roleName, roleSlug, clusterTypeName, clusterTypeSlug, clusterName, configTemplateName, configTemplateCode,
			)
		},
		ConfigWithFields: func() string {
			return testAccDeviceResourceConfig_removeOptionalFields_withFields(
				deviceName, siteName, siteSlug, manufacturerName, manufacturerSlug,
				deviceTypeName, deviceTypeSlug, roleName, roleSlug, clusterTypeName, clusterTypeSlug, clusterName, configTemplateName, configTemplateCode,
			)
		},
		OptionalFields: map[string]string{
			"latitude":        "37.7749",
			"longitude":       "-122.4194",
			"vc_position":     "1",
			"vc_priority":     "100",
			"cluster":         clusterName,
			"config_template": configTemplateName,
		},
		RequiredFields: map[string]string{
			"name": deviceName,
		},
		CheckDestroy: testutil.CheckDeviceDestroy,
	})
}

func testAccDeviceResourceConfig_removeOptionalFields_base(deviceName, siteName, siteSlug, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, roleName, roleSlug, clusterTypeName, clusterTypeSlug, clusterName, configTemplateName, configTemplateCode string) string {
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

resource "netbox_cluster_type" "test" {
	name = %[10]q
	slug = %[11]q
}

resource "netbox_cluster" "test" {
	name = %[12]q
	type = netbox_cluster_type.test.id
}

resource "netbox_config_template" "test" {
	name          = %[13]q
	template_code = %[14]q
}

resource "netbox_device" "test" {
  name        = %[1]q
  site        = netbox_site.test.id
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
}
`, deviceName, siteName, siteSlug, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, roleName, roleSlug, clusterTypeName, clusterTypeSlug, clusterName, configTemplateName, configTemplateCode)
}

func testAccDeviceResourceConfig_removeOptionalFields_withFields(deviceName, siteName, siteSlug, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, roleName, roleSlug, clusterTypeName, clusterTypeSlug, clusterName, configTemplateName, configTemplateCode string) string {
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

resource "netbox_cluster_type" "test" {
	name = %[10]q
	slug = %[11]q
}

resource "netbox_cluster" "test" {
	name = %[12]q
	type = netbox_cluster_type.test.id
}

resource "netbox_config_template" "test" {
	name          = %[13]q
	template_code = %[14]q
}

resource "netbox_device" "test" {
  name        = %[1]q
  site        = netbox_site.test.id
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  latitude    = 37.7749
  longitude   = -122.4194
  vc_position = 1
  vc_priority = 100
	cluster     = netbox_cluster.test.name
	config_template = netbox_config_template.test.name
}
`, deviceName, siteName, siteSlug, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, roleName, roleSlug, clusterTypeName, clusterTypeSlug, clusterName, configTemplateName, configTemplateCode)
}

// TestAccDeviceResource_validationErrors tests validation error scenarios.
func TestAccDeviceResource_validationErrors(t *testing.T) {
	testutil.RunMultiValidationErrorTest(t, testutil.MultiValidationErrorTestConfig{
		ResourceName: "netbox_device",
		TestCases: map[string]testutil.ValidationErrorCase{
			"missing_device_type": {
				Config: func() string {
					return `
resource "netbox_site" "test" {
  name = "Test Site"
  slug = "test-site"
}

resource "netbox_device_role" "test" {
  name  = "Test Role"
  slug  = "test-role"
  color = "ff0000"
}

resource "netbox_device" "test" {
  name = "Test Device"
  site = netbox_site.test.id
  role = netbox_device_role.test.id
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
			"missing_role": {
				Config: func() string {
					return `
resource "netbox_site" "test" {
  name = "Test Site"
  slug = "test-site"
}

resource "netbox_manufacturer" "test" {
  name = "Test Manufacturer"
  slug = "test-mfg"
}

resource "netbox_device_type" "test" {
  model        = "Test Model"
  slug         = "test-model"
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_device" "test" {
  name        = "Test Device"
  site        = netbox_site.test.id
  device_type = netbox_device_type.test.id
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
			"missing_site": {
				Config: func() string {
					return `
resource "netbox_manufacturer" "test" {
  name = "Test Manufacturer"
  slug = "test-mfg"
}

resource "netbox_device_type" "test" {
  model        = "Test Model"
  slug         = "test-model"
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_device_role" "test" {
  name  = "Test Role"
  slug  = "test-role"
  color = "ff0000"
}

resource "netbox_device" "test" {
  name        = "Test Device"
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
			"invalid_status": {
				Config: func() string {
					return `
resource "netbox_site" "test" {
  name = "Test Site"
  slug = "test-site"
}

resource "netbox_manufacturer" "test" {
  name = "Test Manufacturer"
  slug = "test-mfg"
}

resource "netbox_device_type" "test" {
  model        = "Test Model"
  slug         = "test-model"
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_device_role" "test" {
  name  = "Test Role"
  slug  = "test-role"
  color = "ff0000"
}

resource "netbox_device" "test" {
  name        = "Test Device"
  site        = netbox_site.test.id
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  status      = "invalid_status"
}
`
				},
				ExpectedError: testutil.ErrPatternInvalidEnum,
			},
			"invalid_tenant_reference": {
				Config: func() string {
					return `
resource "netbox_site" "test" {
  name = "Test Site"
  slug = "test-site"
}

resource "netbox_manufacturer" "test" {
  name = "Test Manufacturer"
  slug = "test-mfg"
}

resource "netbox_device_type" "test" {
  model        = "Test Model"
  slug         = "test-model"
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_device_role" "test" {
  name  = "Test Role"
  slug  = "test-role"
  color = "ff0000"
}

resource "netbox_device" "test" {
  name        = "Test Device"
  site        = netbox_site.test.id
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  tenant      = "99999"
}
`
				},
				ExpectedError: testutil.ErrPatternNotFound,
			},
			"invalid_platform_reference": {
				Config: func() string {
					return `
resource "netbox_site" "test" {
  name = "Test Site"
  slug = "test-site"
}

resource "netbox_manufacturer" "test" {
  name = "Test Manufacturer"
  slug = "test-mfg"
}

resource "netbox_device_type" "test" {
  model        = "Test Model"
  slug         = "test-model"
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_device_role" "test" {
  name  = "Test Role"
  slug  = "test-role"
  color = "ff0000"
}

resource "netbox_device" "test" {
  name        = "Test Device"
  site        = netbox_site.test.id
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  platform    = "99999"
}
`
				},
				ExpectedError: testutil.ErrPatternNotFound,
			},
		},
	})
}

func TestAccDeviceResource_tagLifecycle(t *testing.T) {
	t.Parallel()

	deviceName := testutil.RandomName("tf-test-device-tags")
	manufacturerName := testutil.RandomName("tf-test-manufacturer-tags")
	manufacturerSlug := testutil.RandomSlug("tf-test-mfr-tags")
	deviceTypeModel := testutil.RandomName("tf-test-device-type-tags")
	deviceTypeSlug := testutil.RandomSlug("tf-test-dt-tags")
	deviceRoleName := testutil.RandomName("tf-test-device-role-tags")
	deviceRoleSlug := testutil.RandomSlug("tf-test-dr-tags")
	siteName := testutil.RandomName("tf-test-site-tags")
	siteSlug := testutil.RandomSlug("tf-test-site-tags")
	tag1Name := testutil.RandomName("tag1")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2Name := testutil.RandomName("tag2")
	tag2Slug := testutil.RandomSlug("tag2")
	tag3Name := testutil.RandomName("tag3")
	tag3Slug := testutil.RandomSlug("tag3")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
	cleanup.RegisterDeviceRoleCleanup(deviceRoleSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)
	cleanup.RegisterDeviceCleanup(deviceName)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)
	cleanup.RegisterTagCleanup(tag3Slug)

	testutil.RunTagLifecycleTest(t, testutil.TagLifecycleTestConfig{
		ResourceName: "netbox_device",
		ConfigWithoutTags: func() string {
			return testAccDeviceResourceConfig_tagLifecycle(deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug, tag1Name, tag1Slug, tag2Name, tag2Slug, tag3Name, tag3Slug, "none")
		},
		ConfigWithTags: func() string {
			return testAccDeviceResourceConfig_tagLifecycle(deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug, tag1Name, tag1Slug, tag2Name, tag2Slug, tag3Name, tag3Slug, caseTag1Uscore2)
		},
		ConfigWithDifferentTags: func() string {
			return testAccDeviceResourceConfig_tagLifecycle(deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug, tag1Name, tag1Slug, tag2Name, tag2Slug, tag3Name, tag3Slug, "tag3")
		},
		ExpectedTagCount:          2,
		ExpectedDifferentTagCount: 1,
		CheckDestroy:              testutil.CheckDeviceDestroy,
	})
}

func TestAccDeviceResource_tagOrderInvariance(t *testing.T) {
	t.Parallel()

	deviceName := testutil.RandomName("tf-test-device-tagorder")
	manufacturerName := testutil.RandomName("tf-test-manufacturer-tagorder")
	manufacturerSlug := testutil.RandomSlug("tf-test-mfr-tagorder")
	deviceTypeModel := testutil.RandomName("tf-test-device-type-tagorder")
	deviceTypeSlug := testutil.RandomSlug("tf-test-dt-tagorder")
	deviceRoleName := testutil.RandomName("tf-test-device-role-tagorder")
	deviceRoleSlug := testutil.RandomSlug("tf-test-dr-tagorder")
	siteName := testutil.RandomName("tf-test-site-tagorder")
	siteSlug := testutil.RandomSlug("tf-test-site-tagorder")
	tag1Name := testutil.RandomName("tag1")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2Name := testutil.RandomName("tag2")
	tag2Slug := testutil.RandomSlug("tag2")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
	cleanup.RegisterDeviceRoleCleanup(deviceRoleSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)
	cleanup.RegisterDeviceCleanup(deviceName)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)

	testutil.RunTagOrderTest(t, testutil.TagOrderTestConfig{
		ResourceName: "netbox_device",
		ConfigWithTagsOrderA: func() string {
			return testAccDeviceResourceConfig_tagOrder(deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug, tag1Name, tag1Slug, tag2Name, tag2Slug, caseTag1Uscore2)
		},
		ConfigWithTagsOrderB: func() string {
			return testAccDeviceResourceConfig_tagOrder(deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug, tag1Name, tag1Slug, tag2Name, tag2Slug, "tag2_tag1")
		},
		ExpectedTagCount: 2,
		CheckDestroy:     testutil.CheckDeviceDestroy,
	})
}

// Config helper for tag lifecycle testing.
func testAccDeviceResourceConfig_tagLifecycle(deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug, tag1Name, tag1Slug, tag2Name, tag2Slug, tag3Name, tag3Slug, tagCase string) string {
	var tagsList string
	switch tagCase {
	case caseTag1Uscore2:
		tagsList = tagsDoubleSlug
	case caseTag3:
		tagsList = tagsSingleSlug
	default:
		tagsList = tagsEmpty
	}

	return fmt.Sprintf(`
resource "netbox_tag" "tag1" {
  name = %[10]q
  slug = %[11]q
}

resource "netbox_tag" "tag2" {
  name = %[12]q
  slug = %[13]q
}

resource "netbox_tag" "tag3" {
  name = %[14]q
  slug = %[15]q
}

resource "netbox_manufacturer" "test" {
  name = %[2]q
  slug = %[3]q
}

resource "netbox_device_type" "test" {
  model = %[4]q
  slug = %[5]q
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_device_role" "test" {
  name = %[6]q
  slug = %[7]q
}

resource "netbox_site" "test" {
  name = %[8]q
  slug = %[9]q
}

resource "netbox_device" "test" {
  name        = %[1]q
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
  %[16]s
}
`, deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug, tag1Name, tag1Slug, tag2Name, tag2Slug, tag3Name, tag3Slug, tagsList)
}

// Config helper for tag order testing.
func testAccDeviceResourceConfig_tagOrder(deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug, tag1Name, tag1Slug, tag2Name, tag2Slug, tagOrder string) string {
	var tagsOrder string
	switch tagOrder {
	case caseTag1Uscore2:
		tagsOrder = tagsDoubleSlug
	case "tag2_tag1":
		tagsOrder = tagsDoubleSlugReversed
	}

	return fmt.Sprintf(`
resource "netbox_tag" "tag1" {
  name = %[10]q
  slug = %[11]q
}

resource "netbox_tag" "tag2" {
  name = %[12]q
  slug = %[13]q
}

resource "netbox_manufacturer" "test" {
  name = %[2]q
  slug = %[3]q
}

resource "netbox_device_type" "test" {
  model = %[4]q
  slug = %[5]q
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_device_role" "test" {
  name = %[6]q
  slug = %[7]q
}

resource "netbox_site" "test" {
  name = %[8]q
  slug = %[9]q
}

resource "netbox_device" "test" {
  name        = %[1]q
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
  %[14]s
}
`, deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug, tag1Name, tag1Slug, tag2Name, tag2Slug, tagsOrder)
}
