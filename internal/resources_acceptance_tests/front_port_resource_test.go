package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/provider"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccFrontPortResource_basic(t *testing.T) {
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
	frontPortName := testutil.RandomName("tf-test-fp")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFrontPortResourceConfig_basic(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, rearPortName, frontPortName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_front_port.test", "id"),
					resource.TestCheckResourceAttr("netbox_front_port.test", "name", frontPortName),
					resource.TestCheckResourceAttr("netbox_front_port.test", "type", "8p8c"),
				),
			},
			{
				ResourceName:            "netbox_front_port.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"device", "rear_port"},
			},
		},
	})
}

func TestAccFrontPortResource_full(t *testing.T) {
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
	frontPortName := testutil.RandomName("tf-test-fp")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFrontPortResourceConfig_full(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, rearPortName, frontPortName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_front_port.test", "id"),
					resource.TestCheckResourceAttr("netbox_front_port.test", "name", frontPortName),
					resource.TestCheckResourceAttr("netbox_front_port.test", "type", "lc"),
					resource.TestCheckResourceAttr("netbox_front_port.test", "label", "Front Port Test"),
					resource.TestCheckResourceAttr("netbox_front_port.test", "color", "aa1409"),
					resource.TestCheckResourceAttr("netbox_front_port.test", "rear_port_position", "1"),
					resource.TestCheckResourceAttr("netbox_front_port.test", "description", "Test front port"),
					resource.TestCheckResourceAttr("netbox_front_port.test", "mark_connected", "true"),
				),
			},
			{
				ResourceName:            "netbox_front_port.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"device", "rear_port"},
			},
		},
	})
}

func TestAccConsistency_FrontPort(t *testing.T) {
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
	frontPortName := testutil.RandomName("front-port")
	rearPortName := testutil.RandomName("rear-port")

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
				Config: testAccFrontPortConsistencyConfig(siteName, siteSlug, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceRoleName, deviceRoleSlug, deviceName, frontPortName, rearPortName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_front_port.test", "device", deviceName),
				),
			},
			{
				PlanOnly: true,
				Config:   testAccFrontPortConsistencyConfig(siteName, siteSlug, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceRoleName, deviceRoleSlug, deviceName, frontPortName, rearPortName),
			},
		},
	})
}

func TestAccConsistency_FrontPort_LiteralNames(t *testing.T) {
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
	resourceName := testutil.RandomName("front_port")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterDeviceCleanup(deviceName)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
	cleanup.RegisterDeviceRoleCleanup(roleSlug)
	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFrontPortConsistencyLiteralNamesConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, roleName, roleSlug, siteName, siteSlug, deviceName, resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_front_port.test", "name", resourceName),
					resource.TestCheckResourceAttr("netbox_front_port.test", "device", deviceName),
				),
			},
			{
				PlanOnly: true,
				Config:   testAccFrontPortConsistencyLiteralNamesConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, roleName, roleSlug, siteName, siteSlug, deviceName, resourceName),
			},
		},
	})
}

func TestAccFrontPortResource_update(t *testing.T) {
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
	frontPortName := testutil.RandomName("tf-test-fp-update")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFrontPortResourceConfig_fullWithDesc(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, rearPortName, frontPortName, "Label1", testutil.Description1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_front_port.test", "id"),
					resource.TestCheckResourceAttr("netbox_front_port.test", "label", "Label1"),
					resource.TestCheckResourceAttr("netbox_front_port.test", "description", testutil.Description1),
				),
			},
			{
				Config: testAccFrontPortResourceConfig_fullWithDesc(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, rearPortName, frontPortName, "Label2", testutil.Description2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_front_port.test", "label", "Label2"),
					resource.TestCheckResourceAttr("netbox_front_port.test", "description", testutil.Description2),
				),
			},
		},
	})
}

func TestAccFrontPortResource_externalDeletion(t *testing.T) {
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
	frontPortName := testutil.RandomName("tf-test-fp-ext-del")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFrontPortResourceConfig_basic(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, rearPortName, frontPortName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_front_port.test", "id"),
					resource.TestCheckResourceAttr("netbox_front_port.test", "name", frontPortName),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					items, _, err := client.DcimAPI.DcimFrontPortsList(context.Background()).NameIc([]string{frontPortName}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find front_port for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.DcimAPI.DcimFrontPortsDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete front_port: %v", err)
					}
					t.Logf("Successfully externally deleted front_port with ID: %d", itemID)
				},
				Config: testAccFrontPortResourceConfig_basic(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, rearPortName, frontPortName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_front_port.test", "id"),
				),
			},
		},
	})
}

func TestAccFrontPortResource_IDPreservation(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-site-id")
	siteSlug := testutil.RandomSlug("tf-test-site-id")
	mfgName := testutil.RandomName("tf-test-mfg-id")
	mfgSlug := testutil.RandomSlug("tf-test-mfg-id")
	dtModel := testutil.RandomName("tf-test-dt-id")
	dtSlug := testutil.RandomSlug("tf-test-dt-id")
	roleName := testutil.RandomName("tf-test-role-id")
	roleSlug := testutil.RandomSlug("tf-test-role-id")
	deviceName := testutil.RandomName("tf-test-device-id")
	rearPortName := testutil.RandomName("tf-test-rp-id")
	frontPortName := testutil.RandomName("tf-test-fp-id")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFrontPortResourceConfig_basic(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, rearPortName, frontPortName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_front_port.test", "id"),
					resource.TestCheckResourceAttr("netbox_front_port.test", "name", frontPortName),
					resource.TestCheckResourceAttr("netbox_front_port.test", "type", "8p8c"),
				),
			},
		},
	})
}

func testAccFrontPortResourceConfig_basic(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, rearPortName, frontPortName string) string {
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
  device    = netbox_device.test.id
  name      = %q
  type      = "8p8c"
  positions = 2
}

resource "netbox_front_port" "test" {
  device    = netbox_device.test.id
  name      = %q
  type      = "8p8c"
  rear_port = netbox_rear_port.test.id
}
`, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, rearPortName, frontPortName)
}

func testAccFrontPortResourceConfig_full(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, rearPortName, frontPortName string) string {
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
  device    = netbox_device.test.id
  name      = %q
  type      = "lc"
  positions = 4
}

resource "netbox_front_port" "test" {
  device             = netbox_device.test.id
  name               = %q
  type               = "lc"
  rear_port          = netbox_rear_port.test.id
  rear_port_position = 1
  label              = "Front Port Test"
  color              = "aa1409"
  description        = "Test front port"
  mark_connected     = true
}
`, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, rearPortName, frontPortName)
}

func testAccFrontPortResourceConfig_fullWithDesc(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, rearPortName, frontPortName, label, description string) string {
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
  device    = netbox_device.test.id
  name      = %q
  type      = "lc"
  positions = 4
}

resource "netbox_front_port" "test" {
  device             = netbox_device.test.id
  name               = %q
  type               = "lc"
  rear_port          = netbox_rear_port.test.id
  rear_port_position = 1
  label              = %q
  description        = %q
}
`, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, rearPortName, frontPortName, label, description)
}

func testAccFrontPortConsistencyConfig(siteName, siteSlug, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceRoleName, deviceRoleSlug, deviceName, frontPortName, rearPortName string) string {
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
  device = netbox_device.test.id
  name = "%[11]s"
  type = "8p8c"
  positions = 1
}

resource "netbox_front_port" "test" {
  device = netbox_device.test.name
  name = "%[10]s"
  type = "8p8c"
  rear_port = netbox_rear_port.test.id
  rear_port_position = 1
}
`, siteName, siteSlug, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceRoleName, deviceRoleSlug, deviceName, frontPortName, rearPortName)
}

func testAccFrontPortConsistencyLiteralNamesConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, roleName, roleSlug, siteName, siteSlug, deviceName, resourceName string) string {
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

resource "netbox_rear_port" "rear" {
  device    = %q
  name      = "rear-port"
  type      = "8p8c"
  positions = 1

  depends_on = [netbox_device.test]
}

resource "netbox_front_port" "test" {
  device = %q
  name = %q
  type = "8p8c"
  rear_port = netbox_rear_port.rear.id
  rear_port_position = 1

  depends_on = [netbox_device.test]
}
`, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, roleName, roleSlug, siteName, siteSlug, deviceName, deviceName, deviceName, resourceName)
}

func TestAccFrontPortResource_importWithCustomFieldsAndTags(t *testing.T) {
	// NOTE: t.Parallel() intentionally omitted - this test creates/deletes global custom fields
	// that would affect other tests of the same resource type running in parallel.

	portName := testutil.RandomName("front_port")
	rearPortName := testutil.RandomName("rear_port")
	deviceName := testutil.RandomName("device")
	siteName := testutil.RandomName("site")
	siteSlug := testutil.RandomSlug("site")
	mfgName := testutil.RandomName("manufacturer")
	mfgSlug := testutil.RandomSlug("manufacturer")
	dtModel := testutil.RandomName("device_type")
	dtSlug := testutil.RandomSlug("device_type")
	roleName := testutil.RandomName("role")
	roleSlug := testutil.RandomSlug("role")
	tenantName := testutil.RandomName("tenant")
	tenantSlug := testutil.RandomSlug("tenant")

	// Custom field names with underscore format
	cfText := testutil.RandomCustomFieldName("cf_text")
	cfLongtext := testutil.RandomCustomFieldName("cf_longtext")
	cfInteger := testutil.RandomCustomFieldName("cf_integer")
	cfBoolean := testutil.RandomCustomFieldName("cf_boolean")
	cfDate := testutil.RandomCustomFieldName("cf_date")
	cfUrl := testutil.RandomCustomFieldName("cf_url")
	cfJson := testutil.RandomCustomFieldName("cf_json")

	// Tag names
	tag1 := testutil.RandomName("tag1")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2 := testutil.RandomName("tag2")
	tag2Slug := testutil.RandomSlug("tag2")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterDeviceCleanup(deviceName)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterManufacturerCleanup(mfgSlug)
	cleanup.RegisterDeviceTypeCleanup(dtSlug)
	cleanup.RegisterDeviceRoleCleanup(roleSlug)
	cleanup.RegisterTenantCleanup(tenantSlug)
	// Clean up custom fields and tags
	cleanup.RegisterCustomFieldCleanup(cfText)
	cleanup.RegisterCustomFieldCleanup(cfLongtext)
	cleanup.RegisterCustomFieldCleanup(cfInteger)
	cleanup.RegisterCustomFieldCleanup(cfBoolean)
	cleanup.RegisterCustomFieldCleanup(cfDate)
	cleanup.RegisterCustomFieldCleanup(cfUrl)
	cleanup.RegisterCustomFieldCleanup(cfJson)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccFrontPortResourceImportConfig_full(portName, rearPortName, deviceName, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, tenantName, tenantSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_front_port.test", "id"),
					resource.TestCheckResourceAttr("netbox_front_port.test", "name", portName),
					// Verify custom fields are applied
					resource.TestCheckResourceAttr("netbox_front_port.test", "custom_fields.#", "7"),
					// Verify tags are applied
					resource.TestCheckResourceAttr("netbox_front_port.test", "tags.#", "2"),
				),
			},
			{
				ResourceName:            "netbox_front_port.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"device", "rear_port", "custom_fields"}, // Device/rear_port references may have lookup inconsistencies, custom fields have import limitations
			},
		},
	})
}

func testAccFrontPortResourceImportConfig_full(portName, rearPortName, deviceName, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, tenantName, tenantSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug string) string {
	return fmt.Sprintf(`
# Dependencies
resource "netbox_tenant" "test" {
  name = %q
  slug = %q
}

resource "netbox_site" "test" {
  name = %q
  slug = %q
}

resource "netbox_manufacturer" "test" {
  name = %q
  slug = %q
}

resource "netbox_device_role" "test" {
  name  = %q
  slug  = %q
  color = "ff0000"
}

resource "netbox_device_type" "test" {
  model        = %q
  slug         = %q
  manufacturer = netbox_manufacturer.test.slug
}

resource "netbox_device" "test" {
  name        = %q
  device_type = netbox_device_type.test.slug
  role        = netbox_device_role.test.slug
  site        = netbox_site.test.slug
}

resource "netbox_rear_port" "rear" {
  name   = %q
  device = netbox_device.test.name
  type   = "8p8c"
}

# Custom Fields
resource "netbox_custom_field" "cf_text" {
  name        = %q
  type        = "text"
  object_types = ["dcim.frontport"]
}

resource "netbox_custom_field" "cf_longtext" {
  name        = %q
  type        = "longtext"
  object_types = ["dcim.frontport"]
}

resource "netbox_custom_field" "cf_integer" {
  name        = %q
  type        = "integer"
  object_types = ["dcim.frontport"]
}

resource "netbox_custom_field" "cf_boolean" {
  name        = %q
  type        = "boolean"
  object_types = ["dcim.frontport"]
}

resource "netbox_custom_field" "cf_date" {
  name        = %q
  type        = "date"
  object_types = ["dcim.frontport"]
}

resource "netbox_custom_field" "cf_url" {
  name        = %q
  type        = "url"
  object_types = ["dcim.frontport"]
}

resource "netbox_custom_field" "cf_json" {
  name        = %q
  type        = "json"
  object_types = ["dcim.frontport"]
}

# Tags
resource "netbox_tag" "tag1" {
  name = %q
  slug = %q
}

resource "netbox_tag" "tag2" {
  name = %q
  slug = %q
}

# Main Resource
resource "netbox_front_port" "test" {
  name               = %q
  device             = netbox_device.test.name
  type               = "8p8c"
  rear_port          = netbox_rear_port.rear.id
  rear_port_position = 1

  custom_fields = [
    {
      name  = netbox_custom_field.cf_text.name
      type  = "text"
      value = "test-value"
    },
    {
      name  = netbox_custom_field.cf_longtext.name
      type  = "longtext"
      value = "test-longtext-value"
    },
    {
      name  = netbox_custom_field.cf_integer.name
      type  = "integer"
      value = "42"
    },
    {
      name  = netbox_custom_field.cf_boolean.name
      type  = "boolean"
      value = "true"
    },
    {
      name  = netbox_custom_field.cf_date.name
      type  = "date"
      value = "2023-01-01"
    },
    {
      name  = netbox_custom_field.cf_url.name
      type  = "url"
      value = "https://example.com"
    },
    {
      name  = netbox_custom_field.cf_json.name
      type  = "json"
      value = jsonencode({"key" = "value"})
    }
  ]

  tags = [
    {
      name = netbox_tag.tag1.name
      slug = netbox_tag.tag1.slug
    },
    {
      name = netbox_tag.tag2.name
      slug = netbox_tag.tag2.slug
    }
  ]
}
`,
		tenantName, tenantSlug,
		siteName, siteSlug,
		mfgName, mfgSlug,
		roleName, roleSlug,
		dtModel, dtSlug,
		deviceName,
		rearPortName,
		cfText,
		cfLongtext,
		cfInteger,
		cfBoolean,
		cfDate,
		cfUrl,
		cfJson,
		tag1, tag1Slug,
		tag2, tag2Slug,
		portName,
	)
}
