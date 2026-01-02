package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

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

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDeviceResourceConfig_full(deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug, serial, assetTag),
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
					resource.TestCheckResourceAttr("netbox_device.test", "tags.#", "0"),
					resource.TestCheckResourceAttr("netbox_device.test", "custom_fields.#", "0"),
				),
			},
		},
	})
}

func TestAccDeviceResource_update(t *testing.T) {

	t.Parallel()
	// Generate unique names
	deviceName := testutil.RandomName("tf-test-device")
	deviceNameUpdated := testutil.RandomName("tf-test-device-updated")
	manufacturerName := testutil.RandomName("tf-test-manufacturer")
	manufacturerSlug := testutil.RandomSlug("tf-test-mfr")
	deviceTypeModel := testutil.RandomName("tf-test-device-type")
	deviceTypeSlug := testutil.RandomSlug("tf-test-dt")
	deviceRoleName := testutil.RandomName("tf-test-device-role")
	deviceRoleSlug := testutil.RandomSlug("tf-test-dr")
	siteName := testutil.RandomName("tf-test-site")
	siteSlug := testutil.RandomSlug("tf-test-site")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDeviceResourceConfig_basic(deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_device.test", "name", deviceName),
					resource.TestCheckResourceAttr("netbox_device.test", "status", "active"),
				),
			},
			{
				Config: testAccDeviceResourceConfig_updated(deviceNameUpdated, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_device.test", "name", deviceNameUpdated),
					resource.TestCheckResourceAttr("netbox_device.test", "status", "staged"),
					resource.TestCheckResourceAttr("netbox_device.test", "description", "Updated description"),
					// Verify field consistency during updates
					resource.TestCheckResourceAttr("netbox_device.test", "tags.#", "0"),
					resource.TestCheckResourceAttr("netbox_device.test", "custom_fields.#", "0"),
				),
			},
		},
	})
}

func TestAccDeviceResource_IDPreservation(t *testing.T) {

	t.Parallel()
	// Generate unique names to avoid conflicts between test runs
	deviceName := testutil.RandomName("dev-id")
	manufacturerName := testutil.RandomName("mfr-dev")
	manufacturerSlug := testutil.GenerateSlug(manufacturerName)
	deviceTypeModel := testutil.RandomName("dt-dev")
	deviceTypeSlug := testutil.GenerateSlug(deviceTypeModel)
	deviceRoleName := testutil.RandomName("dr-dev")
	deviceRoleSlug := testutil.GenerateSlug(deviceRoleName)
	siteName := testutil.RandomName("site-dev")
	siteSlug := testutil.GenerateSlug(siteName)

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterDeviceCleanup(deviceName)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)
	cleanup.RegisterDeviceRoleCleanup(deviceRoleSlug)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)

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
		},
	})
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

// Helper functions to generate test configurations

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

func testAccDeviceResourceConfig_full(deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug, serial, assetTag string) string {
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
  status      = "planned"
  serial      = %[10]q
  asset_tag   = %[11]q
  description = "Test device description"
  comments    = "Test device comments"
  airflow     = "front-to-rear"
  tags        = []
  custom_fields = []
}
`, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug, deviceName, serial, assetTag)
}

func testAccDeviceResourceConfig_updated(deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug string) string {
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
  status      = "staged"
  description = "Updated description"
  tags        = []
  custom_fields = []
}
`, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug, deviceName)
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
				Config: testAccDeviceResourceConfig_basic(deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_device.test", "id"),
				),
			},
		},
	})
}
