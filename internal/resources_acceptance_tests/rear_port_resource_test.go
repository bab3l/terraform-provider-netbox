package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRearPortResource_basic(t *testing.T) {
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
	rearPortName := testutil.RandomName("tf-test-rp")

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
				Config: testAccRearPortResourceConfig_basic(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, rearPortName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_rear_port.test", "id"),
					resource.TestCheckResourceAttr("netbox_rear_port.test", "name", rearPortName),
					resource.TestCheckResourceAttr("netbox_rear_port.test", "type", "8p8c"),
				),
			},
			{
				ResourceName:            "netbox_rear_port.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"device"},
			},
		},
	})
}

func TestAccRearPortResource_update(t *testing.T) {
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
	rearPortName := testutil.RandomName("tf-test-rp-update")
	updatedRearPortName := testutil.RandomName("tf-test-rp-updated")

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
				Config: testAccRearPortResourceConfig_forUpdate(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, rearPortName, testutil.Description1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_rear_port.test", "id"),
					resource.TestCheckResourceAttr("netbox_rear_port.test", "name", rearPortName),
					resource.TestCheckResourceAttr("netbox_rear_port.test", "type", "8p8c"),
					resource.TestCheckResourceAttr("netbox_rear_port.test", "description", testutil.Description1),
				),
			},
			{
				Config: testAccRearPortResourceConfig_forUpdate(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, updatedRearPortName, testutil.Description2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_rear_port.test", "name", updatedRearPortName),
					resource.TestCheckResourceAttr("netbox_rear_port.test", "type", "lc"),
					resource.TestCheckResourceAttr("netbox_rear_port.test", "description", testutil.Description2),
				),
			},
		},
	})
}

func TestAccRearPortResource_full(t *testing.T) {
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
	rearPortName := testutil.RandomName("tf-test-rp")

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
				Config: testAccRearPortResourceConfig_full(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, rearPortName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_rear_port.test", "id"),
					resource.TestCheckResourceAttr("netbox_rear_port.test", "name", rearPortName),
					resource.TestCheckResourceAttr("netbox_rear_port.test", "type", "lc"),
					resource.TestCheckResourceAttr("netbox_rear_port.test", "label", "Rear Port Test"),
					resource.TestCheckResourceAttr("netbox_rear_port.test", "color", "aa1409"),
					resource.TestCheckResourceAttr("netbox_rear_port.test", "positions", "4"),
					resource.TestCheckResourceAttr("netbox_rear_port.test", "description", "Test rear port"),
					resource.TestCheckResourceAttr("netbox_rear_port.test", "mark_connected", "true"),
				),
			},
			{
				ResourceName:            "netbox_rear_port.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"device"},
			},
		},
	})
}

func TestAccConsistency_RearPort(t *testing.T) {
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
	rearPortName := testutil.RandomName("rear-port")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)
	cleanup.RegisterDeviceRoleCleanup(deviceRoleSlug)
	cleanup.RegisterDeviceCleanup(deviceName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRearPortConsistencyConfig(siteName, siteSlug, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceRoleName, deviceRoleSlug, deviceName, rearPortName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_rear_port.test", "device", deviceName),
				),
			},
			{
				PlanOnly: true,
				Config:   testAccRearPortConsistencyConfig(siteName, siteSlug, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceRoleName, deviceRoleSlug, deviceName, rearPortName),
			},
		},
	})
}

func TestAccConsistency_RearPort_LiteralNames(t *testing.T) {
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
	resourceName := testutil.RandomName("rear_port")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)
	cleanup.RegisterDeviceRoleCleanup(roleSlug)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterDeviceCleanup(deviceName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRearPortConsistencyLiteralNamesConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, roleName, roleSlug, siteName, siteSlug, deviceName, resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_rear_port.test", "name", resourceName),
					resource.TestCheckResourceAttr("netbox_rear_port.test", "device", deviceName),
				),
			},
			{
				PlanOnly: true,
				Config:   testAccRearPortConsistencyLiteralNamesConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, roleName, roleSlug, siteName, siteSlug, deviceName, resourceName),
			},
		},
	})
}

func TestAccRearPortResource_externalDeletion(t *testing.T) {
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
	rearPortName := testutil.RandomName("tf-test-rp-ext-del")

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
				Config: testAccRearPortResourceConfig_basic(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, rearPortName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_rear_port.test", "id"),
					resource.TestCheckResourceAttr("netbox_rear_port.test", "name", rearPortName),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					items, _, err := client.DcimAPI.DcimRearPortsList(context.Background()).NameIc([]string{rearPortName}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find rear_port for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.DcimAPI.DcimRearPortsDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete rear_port: %v", err)
					}
					t.Logf("Successfully externally deleted rear_port with ID: %d", itemID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccRearPortResourceConfig_basic(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, rearPortName string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name = %[1]q
  slug = %[2]q
  status = "active"
}

resource "netbox_manufacturer" "test" {
  name = %[3]q
  slug = %[4]q
}

resource "netbox_device_type" "test" {
  model = %[5]q
  slug = %[6]q
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_device_role" "test" {
  name = %[7]q
  slug = %[8]q
  color = "aa1409"
}

resource "netbox_device" "test" {
  name = %[9]q
  device_type = netbox_device_type.test.id
  role = netbox_device_role.test.id
  site = netbox_site.test.id
  status = "active"
}

resource "netbox_rear_port" "test" {
  device = netbox_device.test.name
  name = %[10]q
  type = "8p8c"
}
`, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, rearPortName)
}
func testAccRearPortResourceConfig_forUpdate(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, rearPortName, description string) string {
	rearPortType := testutil.PortType8P8C
	if description == testutil.Description2 {
		rearPortType = testutil.PortTypeLC
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

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = %q
  slug         = %q
}

resource "netbox_role" "test" {
  name = %q
  slug = %q
}

resource "netbox_device" "test" {
  device_type = netbox_device_type.test.id
  role        = netbox_role.test.id
  site        = netbox_site.test.id
  name        = %q
}

resource "netbox_rear_port" "test" {
  device      = netbox_device.test.id
  name        = %q
  type        = %q
  description = %q
}
`, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, rearPortName, rearPortType, description)
}
func testAccRearPortResourceConfig_full(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, rearPortName string) string {
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
  name        = %q
  site        = netbox_site.test.id
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
}

resource "netbox_rear_port" "test" {
  device         = netbox_device.test.id
  name           = %q
  type           = "lc"
  label          = "Rear Port Test"
  color          = "aa1409"
  positions      = 4
  description    = "Test rear port"
  mark_connected = true
}
`, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, rearPortName)
}

func testAccRearPortConsistencyConfig(siteName, siteSlug, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceRoleName, deviceRoleSlug, deviceName, rearPortName string) string {
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

resource "netbox_rear_port" "test" {
  device = netbox_device.test.name
  name = "%[10]s"
  type = "8p8c"
  positions = 1
}
`, siteName, siteSlug, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceRoleName, deviceRoleSlug, deviceName, rearPortName)
}

func testAccRearPortConsistencyLiteralNamesConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, roleName, roleSlug, siteName, siteSlug, deviceName, resourceName string) string {
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

resource "netbox_rear_port" "test" {
  device = %q
  name = %q
  type = "8p8c"
  positions = 1
  depends_on = [netbox_device.test]
}
`, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, roleName, roleSlug, siteName, siteSlug, deviceName, deviceName, resourceName)
}

// TestAccRearPortResource_Positions tests comprehensive scenarios for rear port positions field.
// This validates that Optional+Computed int32 fields with proper defaults work correctly.
func TestAccRearPortResource_Positions(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-site-rear-port")
	siteSlug := testutil.RandomSlug("tf-test-site-rear-port")
	manufacturerName := testutil.RandomName("tf-test-manufacturer-rear-port")
	manufacturerSlug := testutil.RandomSlug("tf-test-manufacturer-rear-port")
	deviceTypeName := testutil.RandomName("tf-test-device-type-rear-port")
	deviceTypeSlug := testutil.RandomSlug("tf-test-device-type-rear-port")
	deviceRoleName := testutil.RandomName("tf-test-role-rear-port")
	deviceRoleSlug := testutil.RandomSlug("tf-test-role-rear-port")
	deviceName := testutil.RandomName("tf-test-device-rear-port")
	rearPortName := testutil.RandomName("rear-port-positions-test")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)
	cleanup.RegisterDeviceRoleCleanup(deviceRoleSlug)
	cleanup.RegisterDeviceCleanup(deviceName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckRearPortDestroy,
			testutil.CheckDeviceDestroy,
			testutil.CheckDeviceTypeDestroy,
			testutil.CheckDeviceRoleDestroy,
			testutil.CheckSiteDestroy,
			testutil.CheckManufacturerDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccRearPortResourceWithOptionalField(siteName, siteSlug, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceRoleName, deviceRoleSlug, deviceName, rearPortName, "positions", ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_rear_port.test", "id"),
					resource.TestCheckResourceAttr("netbox_rear_port.test", "name", rearPortName),
					resource.TestCheckResourceAttr("netbox_rear_port.test", "positions", "1"),
				),
			},
			{
				Config:   testAccRearPortResourceWithOptionalField(siteName, siteSlug, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceRoleName, deviceRoleSlug, deviceName, rearPortName, "positions", ""),
				PlanOnly: true,
			},
			{
				Config: testAccRearPortResourceWithOptionalField(siteName, siteSlug, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceRoleName, deviceRoleSlug, deviceName, rearPortName, "positions", "4"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_rear_port.test", "id"),
					resource.TestCheckResourceAttr("netbox_rear_port.test", "name", rearPortName),
					resource.TestCheckResourceAttr("netbox_rear_port.test", "positions", "4"),
				),
			},
			{
				Config:   testAccRearPortResourceWithOptionalField(siteName, siteSlug, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceRoleName, deviceRoleSlug, deviceName, rearPortName, "positions", "4"),
				PlanOnly: true,
			},
		},
	})
}

func testAccRearPortResourceWithOptionalField(siteName, siteSlug, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceRoleName, deviceRoleSlug, deviceName, rearPortName, optionalFieldName, optionalFieldValue string) string {
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
  device_type = netbox_device_type.test.id
  name        = %q
  site        = netbox_site.test.id
  role        = netbox_device_role.test.id
}

resource "netbox_rear_port" "test" {
  device = netbox_device.test.id
  name   = %q
  type   = "8p8c"%s
}
`, siteName, siteSlug, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceRoleName, deviceRoleSlug, deviceName, rearPortName, optionalField)
}

func TestAccRearPortResource_removeOptionalFields(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-site-rem")
	siteSlug := testutil.RandomSlug("tf-test-site-rem")
	mfgName := testutil.RandomName("tf-test-mfg-rem")
	mfgSlug := testutil.RandomSlug("tf-test-mfg-rem")
	dtModel := testutil.RandomName("tf-test-dt-rem")
	dtSlug := testutil.RandomSlug("tf-test-dt-rem")
	roleName := testutil.RandomName("tf-test-role-rem")
	roleSlug := testutil.RandomSlug("tf-test-role-rem")
	deviceName := testutil.RandomName("tf-test-device-rem")
	portName := testutil.RandomName("tf-test-rp-rem")
	const testLabel = "Test Label"
	const testDescription = "Test Description"

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
				Config: testAccRearPortResourceConfig_withLabel(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, portName, testLabel, testDescription),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_rear_port.test", "name", portName),
					resource.TestCheckResourceAttr("netbox_rear_port.test", "label", testLabel),
					resource.TestCheckResourceAttr("netbox_rear_port.test", "description", testDescription),
					resource.TestCheckResourceAttr("netbox_rear_port.test", "color", "ff0000"),
					resource.TestCheckResourceAttr("netbox_rear_port.test", "mark_connected", "true"),
					resource.TestCheckResourceAttr("netbox_rear_port.test", "positions", "8"),
				),
			},
			{
				Config: testAccRearPortResourceConfig_basic(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, portName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_rear_port.test", "name", portName),
					resource.TestCheckNoResourceAttr("netbox_rear_port.test", "label"),
					resource.TestCheckNoResourceAttr("netbox_rear_port.test", "description"),
					resource.TestCheckNoResourceAttr("netbox_rear_port.test", "color"),
					resource.TestCheckResourceAttr("netbox_rear_port.test", "mark_connected", "false"),
					resource.TestCheckResourceAttr("netbox_rear_port.test", "positions", "1"),
				),
			},
		},
	})
}

func testAccRearPortResourceConfig_withLabel(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, portName, label, description string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name = %[1]q
  slug = %[2]q
  status = "active"
}

resource "netbox_manufacturer" "test" {
  name = %[3]q
  slug = %[4]q
}

resource "netbox_device_type" "test" {
  model = %[5]q
  slug = %[6]q
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_device_role" "test" {
  name = %[7]q
  slug = %[8]q
  color = "aa1409"
}

resource "netbox_device" "test" {
  name = %[9]q
  device_type = netbox_device_type.test.id
  role = netbox_device_role.test.id
  site = netbox_site.test.id
  status = "active"
}

resource "netbox_rear_port" "test" {
  device        = netbox_device.test.name
  name          = %[10]q
  type          = "8p8c"
  label         = %[11]q
  description   = %[12]q
  color         = "ff0000"
  mark_connected = true
  positions     = 8
}
`, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, portName, label, description)
}
