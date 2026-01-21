package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDeviceBayResource_basic(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-site")
	siteSlug := testutil.RandomSlug("tf-test-site")
	mfgName := testutil.RandomName("tf-test-mfg")
	mfgSlug := testutil.RandomSlug("tf-test-mfg")
	dtModel := testutil.RandomName("tf-test-dt")
	dtSlug := testutil.RandomSlug("tf-test-dt")
	roleName := testutil.RandomName("tf-test-role")
	roleSlug := testutil.RandomSlug("tf-test-role")
	deviceName := testutil.RandomName("tf-test-device")
	bayName := testutil.RandomName("tf-test-bay")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterManufacturerCleanup(mfgSlug)
	cleanup.RegisterDeviceTypeCleanup(dtSlug)
	cleanup.RegisterDeviceRoleCleanup(roleSlug)
	cleanup.RegisterDeviceCleanup(deviceName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDeviceBayResourceConfig_basic(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, bayName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_device_bay.test", "id"),
					resource.TestCheckResourceAttr("netbox_device_bay.test", "name", bayName),
				),
			},
			{
				ResourceName:            "netbox_device_bay.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"device"},
				Check: resource.ComposeTestCheckFunc(
					testutil.ReferenceFieldCheck("netbox_device_bay.test", "device"),
					testutil.ReferenceFieldCheck("netbox_device_bay.test", "installed_device"),
				),
			},
			{
				Config:   testAccDeviceBayResourceConfig_basic(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, bayName),
				PlanOnly: true,
			},
		},
	})
}

func TestAccDeviceBayResource_full(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-site-full")
	siteSlug := testutil.RandomSlug("tf-test-site-full")
	mfgName := testutil.RandomName("tf-test-mfg-full")
	mfgSlug := testutil.RandomSlug("tf-test-mfg-full")
	dtModel := testutil.RandomName("tf-test-dt-full")
	dtSlug := testutil.RandomSlug("tf-test-dt-full")
	roleName := testutil.RandomName("tf-test-role-full")
	roleSlug := testutil.RandomSlug("tf-test-role-full")
	deviceName := testutil.RandomName("tf-test-device-full")
	bayName := testutil.RandomName("tf-test-bay-full")
	description := testutil.RandomName("description")
	updatedDescription := "Updated device bay description"

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterManufacturerCleanup(mfgSlug)
	cleanup.RegisterDeviceTypeCleanup(dtSlug)
	cleanup.RegisterDeviceRoleCleanup(roleSlug)
	cleanup.RegisterDeviceCleanup(deviceName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDeviceBayResourceConfig_full(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, bayName, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_device_bay.test", "id"),
					resource.TestCheckResourceAttr("netbox_device_bay.test", "name", bayName),
					resource.TestCheckResourceAttr("netbox_device_bay.test", "label", "Bay Label"),
					resource.TestCheckResourceAttr("netbox_device_bay.test", "description", description),
				),
			},
			{
				Config: testAccDeviceBayResourceConfig_full(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, bayName, updatedDescription),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_device_bay.test", "description", updatedDescription),
				),
			},
		},
	})
}

func TestAccConsistency_DeviceBay(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("site")
	siteSlug := testutil.RandomSlug("site")
	manufacturerName := testutil.RandomName("manufacturer")
	manufacturerSlug := testutil.RandomSlug("manufacturer")
	deviceTypeName := testutil.RandomName("device-type")
	deviceTypeSlug := testutil.RandomSlug("device-type")
	deviceRoleName := testutil.RandomName("device-role")
	deviceRoleSlug := testutil.RandomSlug("device-role")
	deviceName := testutil.RandomName("device")
	deviceBayName := testutil.RandomName("device-bay")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterDeviceRoleCleanup(deviceRoleSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDeviceBayConsistencyConfig(siteName, siteSlug, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceRoleName, deviceRoleSlug, deviceName, deviceBayName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_device_bay.test", "device", deviceName),
				),
			},
			{
				PlanOnly: true,
				Config:   testAccDeviceBayConsistencyConfig(siteName, siteSlug, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceRoleName, deviceRoleSlug, deviceName, deviceBayName),
			},
		},
	})
}

func TestAccConsistency_DeviceBay_LiteralNames(t *testing.T) {
	t.Parallel()

	manufacturerName := testutil.RandomName("manufacturer")
	manufacturerSlug := testutil.RandomSlug("manufacturer")
	deviceTypeName := testutil.RandomName("device-type")
	deviceTypeSlug := testutil.RandomSlug("device-type")
	roleName := testutil.RandomName("role")
	roleSlug := testutil.RandomSlug("role")
	siteName := testutil.RandomName("site")
	siteSlug := testutil.RandomSlug("site")
	deviceName := testutil.RandomName("device")
	resourceName := testutil.RandomName("device_bay")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterDeviceRoleCleanup(roleSlug)
	cleanup.RegisterDeviceCleanup(deviceName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDeviceBayConsistencyLiteralNamesConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, roleName, roleSlug, siteName, siteSlug, deviceName, resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_device_bay.test", "name", resourceName),
					resource.TestCheckResourceAttr("netbox_device_bay.test", "device", deviceName),
				),
			},
			{
				PlanOnly: true,
				Config:   testAccDeviceBayConsistencyLiteralNamesConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, roleName, roleSlug, siteName, siteSlug, deviceName, resourceName),
			},
		},
	})
}

func TestAccDeviceBayResource_update(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-site-update")
	siteSlug := testutil.RandomSlug("tf-test-site-update")
	mfgName := testutil.RandomName("tf-test-mfg-update")
	mfgSlug := testutil.RandomSlug("tf-test-mfg-update")
	dtModel := testutil.RandomName("tf-test-dt-update")
	dtSlug := testutil.RandomSlug("tf-test-dt-update")
	roleName := testutil.RandomName("tf-test-role-update")
	roleSlug := testutil.RandomSlug("tf-test-role-update")
	deviceName := testutil.RandomName("tf-test-device-update")
	bayName := testutil.RandomName("tf-test-bay-update")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterManufacturerCleanup(mfgSlug)
	cleanup.RegisterDeviceTypeCleanup(dtSlug)
	cleanup.RegisterDeviceRoleCleanup(roleSlug)
	cleanup.RegisterDeviceCleanup(deviceName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDeviceBayResourceConfig_withDescription(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, bayName, testutil.Description1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_device_bay.test", "id"),
					resource.TestCheckResourceAttr("netbox_device_bay.test", "name", bayName),
					resource.TestCheckResourceAttr("netbox_device_bay.test", "description", testutil.Description1),
				),
			},
			{
				Config: testAccDeviceBayResourceConfig_withDescription(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, bayName, testutil.Description2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_device_bay.test", "description", testutil.Description2),
				),
			},
		},
	})
}

func TestAccDeviceBayResource_externalDeletion(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-site-ext-del")
	siteSlug := testutil.RandomSlug("tf-test-site-ext-del")
	mfgName := testutil.RandomName("tf-test-mfg-ext-del")
	mfgSlug := testutil.RandomSlug("tf-test-mfg-ext-del")
	dtModel := testutil.RandomName("tf-test-dt-ext-del")
	dtSlug := testutil.RandomSlug("tf-test-dt-ext-del")
	roleName := testutil.RandomName("tf-test-role-ext-del")
	roleSlug := testutil.RandomSlug("tf-test-role-ext-del")
	deviceName := testutil.RandomName("tf-test-device-ext-del")
	bayName := testutil.RandomName("tf-test-bay-ext-del")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterManufacturerCleanup(mfgSlug)
	cleanup.RegisterDeviceTypeCleanup(dtSlug)
	cleanup.RegisterDeviceRoleCleanup(roleSlug)
	cleanup.RegisterDeviceCleanup(deviceName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDeviceBayResourceConfig_basic(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, bayName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_device_bay.test", "id"),
					resource.TestCheckResourceAttr("netbox_device_bay.test", "name", bayName),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					items, _, err := client.DcimAPI.DcimDeviceBaysList(context.Background()).NameIc([]string{bayName}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find device_bay for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.DcimAPI.DcimDeviceBaysDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete device_bay: %v", err)
					}
					t.Logf("Successfully externally deleted device_bay with ID: %d", itemID)
				},
				Config: testAccDeviceBayResourceConfig_basic(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, bayName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_device_bay.test", "id"),
				),
			},
		},
	})
}

func testAccDeviceBayResourceConfig_basic(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, bayName string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name   = %q
  slug   = %q
  status = "active"
}

resource "netbox_manufacturer" "test" {
  name = %q
  slug = %q
}

resource "netbox_device_type" "test" {
  manufacturer    = netbox_manufacturer.test.id
  model           = %q
  slug            = %q
  subdevice_role  = "parent"
}

resource "netbox_device_role" "test" {
  name  = %q
  slug  = %q
  color = "aa1409"
}

resource "netbox_device" "test" {
  name        = %q
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
}

resource "netbox_device_bay" "test" {
  device = netbox_device.test.id
  name   = %q
}
`, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, bayName)
}

func testAccDeviceBayResourceConfig_withDescription(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, bayName, description string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name   = %q
  slug   = %q
  status = "active"
}

resource "netbox_manufacturer" "test" {
  name = %q
  slug = %q
}

resource "netbox_device_type" "test" {
  manufacturer    = netbox_manufacturer.test.id
  model           = %q
  slug            = %q
  subdevice_role  = "parent"
}

resource "netbox_device_role" "test" {
  name  = %q
  slug  = %q
  color = "aa1409"
}

resource "netbox_device" "test" {
  name        = %q
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
}

resource "netbox_device_bay" "test" {
  device      = netbox_device.test.id
  name        = %q
  description = %q
}
`, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, bayName, description)
}

func testAccDeviceBayResourceConfig_full(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, bayName, description string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name   = %q
  slug   = %q
  status = "active"
}

resource "netbox_manufacturer" "test" {
  name = %q
  slug = %q
}

resource "netbox_device_type" "test" {
  manufacturer   = netbox_manufacturer.test.id
  model          = %q
  slug           = %q
  subdevice_role = "parent"
}

resource "netbox_device_role" "test" {
  name  = %q
  slug  = %q
  color = "aa1409"
}

resource "netbox_device" "test" {
  name        = %q
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
}

resource "netbox_device_bay" "test" {
  device      = netbox_device.test.id
  name        = %q
  label       = "Bay Label"
  description = %q
}
`, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, bayName, description)
}

func testAccDeviceBayConsistencyConfig(siteName, siteSlug, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceRoleName, deviceRoleSlug, deviceName, deviceBayName string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name = "%[1]s"
  slug = "%[2]s"
}

resource "netbox_manufacturer" "test" {
  name = "%[3]s"
  slug = "%[4]s"
}

resource "netbox_device_type" "test" {
  model = "%[5]s"
  slug = "%[6]s"
  manufacturer = netbox_manufacturer.test.id
  subdevice_role = "parent"
}

resource "netbox_device_role" "test" {
  name = "%[7]s"
  slug = "%[8]s"
}

resource "netbox_device" "test" {
  name = "%[9]s"
  device_type = netbox_device_type.test.id
  role = netbox_device_role.test.id
  site = netbox_site.test.id
}

resource "netbox_device_bay" "test" {
  device = netbox_device.test.name
  name = "%[10]s"
}
`, siteName, siteSlug, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceRoleName, deviceRoleSlug, deviceName, deviceBayName)
}

func testAccDeviceBayConsistencyLiteralNamesConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, roleName, roleSlug, siteName, siteSlug, deviceName, resourceName string) string {
	return fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
  name = %q
  slug = %q
}

resource "netbox_device_type" "test" {
  model          = %q
  slug           = %q
  manufacturer   = netbox_manufacturer.test.id
  subdevice_role = "parent"
}

resource "netbox_site" "test" {
  name = %q
  slug = %q
}

resource "netbox_device_role" "test" {
  name = %q
  slug = %q
  color = "ff0000"
}

resource "netbox_device" "test" {
  name        = %q
  device_type = netbox_device_type.test.id
  site        = netbox_site.test.id
  role        = netbox_device_role.test.id
  status      = "active"
}

resource "netbox_device_bay" "test" {
  device = %q
  name = %q
  depends_on = [netbox_device.test]
}
`, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, roleName, roleSlug, siteName, siteSlug, deviceName, deviceName, resourceName)
}

// NOTE: Custom field tests for device_bay resource are in resources_acceptance_tests_customfields package

// TestAccDeviceBayResource_removeOptionalFields tests that optional nullable fields
// can be successfully removed from the configuration without causing inconsistent state.
// This verifies the bugfix for: "Provider produced inconsistent result after apply".
func TestAccDeviceBayResource_removeOptionalFields(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-site-db")
	siteSlug := testutil.RandomSlug("tf-test-site-db")
	mfgName := testutil.RandomName("tf-test-mfg-db")
	mfgSlug := testutil.RandomSlug("tf-test-mfg-db")
	dtModel := testutil.RandomName("tf-test-dt-db")
	dtSlug := testutil.RandomSlug("tf-test-dt-db")
	roleName := testutil.RandomName("tf-test-role-db")
	roleSlug := testutil.RandomSlug("tf-test-role-db")
	parentDeviceName := testutil.RandomName("tf-test-parent-db")
	childDeviceName := testutil.RandomName("tf-test-child-db")
	bayName := testutil.RandomName("tf-test-bay-db")
	const testLabel = "Test Bay Label"

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterManufacturerCleanup(mfgSlug)
	cleanup.RegisterDeviceTypeCleanup(dtSlug)
	cleanup.RegisterDeviceRoleCleanup(roleSlug)
	cleanup.RegisterDeviceCleanup(parentDeviceName)
	cleanup.RegisterDeviceCleanup(childDeviceName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Step 1: Create device bay with all optional fields
			{
				Config: testAccDeviceBayResourceConfig_withAllOptionalFields(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, parentDeviceName, childDeviceName, bayName, testLabel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_device_bay.test", "id"),
					resource.TestCheckResourceAttr("netbox_device_bay.test", "name", bayName),
					resource.TestCheckResourceAttr("netbox_device_bay.test", "label", testLabel),
					resource.TestCheckResourceAttrSet("netbox_device_bay.test", "installed_device"),
				),
			},
			// Step 2: Remove optional fields - should set them to null
			{
				Config: testAccDeviceBayResourceConfig_minimal(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, parentDeviceName, childDeviceName, bayName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_device_bay.test", "id"),
					resource.TestCheckResourceAttr("netbox_device_bay.test", "name", bayName),
					resource.TestCheckNoResourceAttr("netbox_device_bay.test", "label"),
					resource.TestCheckNoResourceAttr("netbox_device_bay.test", "installed_device"),
				),
			},
			// Step 3: Re-add optional fields - verify they can be set again
			{
				Config: testAccDeviceBayResourceConfig_withAllOptionalFields(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, parentDeviceName, childDeviceName, bayName, testLabel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_device_bay.test", "id"),
					resource.TestCheckResourceAttr("netbox_device_bay.test", "name", bayName),
					resource.TestCheckResourceAttr("netbox_device_bay.test", "label", testLabel),
					resource.TestCheckResourceAttrSet("netbox_device_bay.test", "installed_device"),
				),
			},
		},
	})
}

func testAccDeviceBayResourceConfig_withAllOptionalFields(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, parentDeviceName, childDeviceName, bayName, label string) string {
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
  model        = %q
  slug         = %q
  manufacturer = netbox_manufacturer.test.id
  subdevice_role = "parent"
}

resource "netbox_device_role" "test" {
  name  = %q
  slug  = %q
  color = "ff0000"
}

resource "netbox_device" "parent" {
  name        = %q
  device_type = netbox_device_type.test.id
  site        = netbox_site.test.id
  role        = netbox_device_role.test.id
  status      = "active"
}

resource "netbox_device" "child" {
  name        = %q
  device_type = netbox_device_type.test.id
  site        = netbox_site.test.id
  role        = netbox_device_role.test.id
  status      = "active"
}

resource "netbox_device_bay" "test" {
  device           = netbox_device.parent.id
  name             = %q
  label            = %q
  installed_device = netbox_device.child.id
}
`, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, parentDeviceName, childDeviceName, bayName, label)
}

func testAccDeviceBayResourceConfig_minimal(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, parentDeviceName, childDeviceName, bayName string) string {
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
  model        = %q
  slug         = %q
  manufacturer = netbox_manufacturer.test.id
  subdevice_role = "parent"
}

resource "netbox_device_role" "test" {
  name  = %q
  slug  = %q
  color = "ff0000"
}

resource "netbox_device" "parent" {
  name        = %q
  device_type = netbox_device_type.test.id
  site        = netbox_site.test.id
  role        = netbox_device_role.test.id
  status      = "active"
}

resource "netbox_device" "child" {
  name        = %q
  device_type = netbox_device_type.test.id
  site        = netbox_site.test.id
  role        = netbox_device_role.test.id
  status      = "active"
}

resource "netbox_device_bay" "test" {
  device = netbox_device.parent.id
  name   = %q
}
`, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, parentDeviceName, childDeviceName, bayName)
}

func TestAccDeviceBayResource_validationErrors(t *testing.T) {
	testutil.RunMultiValidationErrorTest(t, testutil.MultiValidationErrorTestConfig{
		ResourceName: "netbox_device_bay",
		TestCases: map[string]testutil.ValidationErrorCase{
			"missing_device": {
				Config: func() string {
					return `
resource "netbox_device_bay" "test" {
  name = "Bay1"
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
			"missing_name": {
				Config: func() string {
					return `
resource "netbox_device" "test" {
  name = "test-device"
  device_type = "1"
  role = "1"
  site = "1"
  status = "active"
}

resource "netbox_device_bay" "test" {
  device = netbox_device.test.id
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
			"invalid_device_reference": {
				Config: func() string {
					return `
resource "netbox_device_bay" "test" {
  device = "99999"
  name = "Bay1"
}
`
				},
				ExpectedError: testutil.ErrPatternNotFound,
			},
		},
	})
}
