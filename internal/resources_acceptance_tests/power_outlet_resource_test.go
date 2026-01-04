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

func TestAccPowerOutletResource_basic(t *testing.T) {
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
	powerOutletName := testutil.RandomName("tf-test-po")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterManufacturerCleanup(mfgSlug)
	cleanup.RegisterDeviceTypeCleanup(dtSlug)
	cleanup.RegisterDeviceRoleCleanup(roleSlug)
	cleanup.RegisterDeviceCleanup(deviceName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccPowerOutletResourceConfig_basic(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, powerOutletName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_power_outlet.test", "id"),
					resource.TestCheckResourceAttr("netbox_power_outlet.test", "name", powerOutletName),
				),
			},
			{
				ResourceName:            "netbox_power_outlet.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"device"},
			},
		},
	})
}

func TestAccPowerOutletResource_full(t *testing.T) {
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
	powerOutletName := testutil.RandomName("tf-test-po-full")
	description := "Test power outlet with all fields"
	updatedDescription := "Updated power outlet description"

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterManufacturerCleanup(mfgSlug)
	cleanup.RegisterDeviceTypeCleanup(dtSlug)
	cleanup.RegisterDeviceRoleCleanup(roleSlug)
	cleanup.RegisterDeviceCleanup(deviceName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccPowerOutletResourceConfig_full(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, powerOutletName, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_power_outlet.test", "id"),
					resource.TestCheckResourceAttr("netbox_power_outlet.test", "name", powerOutletName),
					resource.TestCheckResourceAttr("netbox_power_outlet.test", "description", description),
					resource.TestCheckResourceAttr("netbox_power_outlet.test", "type", "iec-60320-c13"),
				),
			},
			{
				Config: testAccPowerOutletResourceConfig_full(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, powerOutletName, updatedDescription),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_power_outlet.test", "description", updatedDescription),
				),
			},
		},
	})
}

func TestAccPowerOutletResource_update(t *testing.T) {
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
	powerOutletName := testutil.RandomName("tf-test-po-update")

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
				Config: testAccPowerOutletResourceConfig_full(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, powerOutletName, testutil.Description1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_power_outlet.test", "id"),
					resource.TestCheckResourceAttr("netbox_power_outlet.test", "description", testutil.Description1),
				),
			},
			{
				Config: testAccPowerOutletResourceConfig_full(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, powerOutletName, testutil.Description2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_power_outlet.test", "description", testutil.Description2),
				),
			},
		},
	})
}

func TestAccPowerOutletResource_externalDeletion(t *testing.T) {
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
	powerOutletName := testutil.RandomName("tf-test-po-ext-del")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPowerOutletResourceConfig_basic(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, powerOutletName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_power_outlet.test", "id"),
					resource.TestCheckResourceAttr("netbox_power_outlet.test", "name", powerOutletName),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					items, _, err := client.DcimAPI.DcimPowerOutletsList(context.Background()).NameIc([]string{powerOutletName}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find power_outlet for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.DcimAPI.DcimPowerOutletsDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete power_outlet: %v", err)
					}
					t.Logf("Successfully externally deleted power_outlet with ID: %d", itemID)
				},
				Config: testAccPowerOutletResourceConfig_basic(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, powerOutletName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_power_outlet.test", "id"),
				),
			},
		},
	})
}

func TestAccPowerOutletResource_IDPreservation(t *testing.T) {
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
	powerOutletName := testutil.RandomName("tf-test-po-id")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterManufacturerCleanup(mfgSlug)
	cleanup.RegisterDeviceTypeCleanup(dtSlug)
	cleanup.RegisterDeviceRoleCleanup(roleSlug)
	cleanup.RegisterDeviceCleanup(deviceName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccPowerOutletResourceConfig_basic(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, powerOutletName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_power_outlet.test", "id"),
					resource.TestCheckResourceAttr("netbox_power_outlet.test", "name", powerOutletName),
				),
			},
		},
	})
}

func testAccPowerOutletResourceConfig_basic(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, powerOutletName string) string {
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
  manufacturer = netbox_manufacturer.test.id
  model        = %q
  slug         = %q
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

resource "netbox_power_outlet" "test" {
  device = netbox_device.test.id
  name   = %q
}
`, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, powerOutletName)
}

func testAccPowerOutletResourceConfig_full(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, powerOutletName, description string) string {
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
  manufacturer = netbox_manufacturer.test.id
  model        = %q
  slug         = %q
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

resource "netbox_power_outlet" "test" {
  device      = netbox_device.test.id
  name        = %q
  type        = "iec-60320-c13"
  description = %q
}
`, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, powerOutletName, description)
}

func TestAccConsistency_PowerOutlet(t *testing.T) {
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
	powerOutletName := testutil.RandomName("power-outlet")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPowerOutletConsistencyConfig(siteName, siteSlug, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceRoleName, deviceRoleSlug, deviceName, powerOutletName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_power_outlet.test", "device", deviceName),
				),
			},
			{
				PlanOnly: true,
				Config:   testAccPowerOutletConsistencyConfig(siteName, siteSlug, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceRoleName, deviceRoleSlug, deviceName, powerOutletName),
			},
		},
	})
}

func testAccPowerOutletConsistencyConfig(siteName, siteSlug, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceRoleName, deviceRoleSlug, deviceName, powerOutletName string) string {
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

resource "netbox_power_outlet" "test" {
  device = netbox_device.test.name
  name = "%[10]s"
  type = "iec-60320-c13"
}
`, siteName, siteSlug, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceRoleName, deviceRoleSlug, deviceName, powerOutletName)
}

func TestAccConsistency_PowerOutlet_LiteralNames(t *testing.T) {

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
	resourceName := testutil.RandomName("power_outlet")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPowerOutletConsistencyLiteralNamesConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, roleName, roleSlug, siteName, siteSlug, deviceName, resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_power_outlet.test", "name", resourceName),
					resource.TestCheckResourceAttr("netbox_power_outlet.test", "device", deviceName),
				),
			},
			{
				PlanOnly: true,
				Config:   testAccPowerOutletConsistencyLiteralNamesConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, roleName, roleSlug, siteName, siteSlug, deviceName, resourceName),
			},
		},
	})
}

func testAccPowerOutletConsistencyLiteralNamesConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, roleName, roleSlug, siteName, siteSlug, deviceName, resourceName string) string {
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

resource "netbox_power_outlet" "test" {
  device = %q
  name = %q
  type = "iec-60320-c13"
  depends_on = [netbox_device.test]
}
`, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, roleName, roleSlug, siteName, siteSlug, deviceName, deviceName, resourceName)
}

func TestAccPowerOutletResource_importWithCustomFieldsAndTags(t *testing.T) {
	// NOTE: t.Parallel() intentionally omitted - this test creates/deletes global custom fields
	// that would affect other tests of the same resource type running in parallel.

	powerOutletName := testutil.RandomName("power_outlet")
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

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckPowerOutletDestroy,
		Steps: []resource.TestStep{
			// First create the resource
			{
				Config: testAccPowerOutletResourceImportConfig_full(powerOutletName, deviceName, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, tenantName, tenantSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_power_outlet.test", "id"),
					resource.TestCheckResourceAttr("netbox_power_outlet.test", "name", powerOutletName),
				),
			},
			// Then test import
			{
				Config:                  testAccPowerOutletResourceImportConfig_full(powerOutletName, deviceName, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, tenantName, tenantSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug),
				ResourceName:            "netbox_power_outlet.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"device", "custom_fields"}, // Device reference may have lookup inconsistencies
			},
		},
	})
}

func testAccPowerOutletResourceImportConfig_full(powerOutletName, deviceName, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, tenantName, tenantSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug string) string {
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
  name = %q
  slug = %q
  color = "9e9e9e"
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model = %q
  slug = %q
}

resource "netbox_device" "test" {
  device_type = netbox_device_type.test.id
  role = netbox_device_role.test.id
  site = netbox_site.test.id
  name = %q
}

# Custom field definitions with different object types
resource "netbox_custom_field" "test_text" {
  name = %q
  type = "text"
  object_types = ["dcim.poweroutlet"]
}

resource "netbox_custom_field" "test_longtext" {
  name = %q
  type = "longtext"
  object_types = ["dcim.poweroutlet"]
}

resource "netbox_custom_field" "test_integer" {
  name = %q
  type = "integer"
  object_types = ["dcim.poweroutlet"]
}

resource "netbox_custom_field" "test_boolean" {
  name = %q
  type = "boolean"
  object_types = ["dcim.poweroutlet"]
}

resource "netbox_custom_field" "test_date" {
  name = %q
  type = "date"
  object_types = ["dcim.poweroutlet"]
}

resource "netbox_custom_field" "test_url" {
  name = %q
  type = "url"
  object_types = ["dcim.poweroutlet"]
}

resource "netbox_custom_field" "test_json" {
  name = %q
  type = "json"
  object_types = ["dcim.poweroutlet"]
}

# Tag definitions
resource "netbox_tag" "test_1" {
  name = %q
  slug = %q
  color = "ff0000"
}

resource "netbox_tag" "test_2" {
  name = %q
  slug = %q
  color = "00ff00"
}

# Power outlet with custom fields and tags
resource "netbox_power_outlet" "test" {
  device = netbox_device.test.id
  name = %q
  type = "iec-60320-c13"

  custom_fields = [
    {
      name  = netbox_custom_field.test_text.name
      type  = "text"
      value = "custom text value"
    },
    {
      name  = netbox_custom_field.test_longtext.name
      type  = "longtext"
      value = "custom longtext value"
    },
    {
      name  = netbox_custom_field.test_integer.name
      type  = "integer"
      value = "123"
    },
    {
      name  = netbox_custom_field.test_boolean.name
      type  = "boolean"
      value = "false"
    },
    {
      name  = netbox_custom_field.test_date.name
      type  = "date"
      value = "2023-12-25"
    },
    {
      name  = netbox_custom_field.test_url.name
      type  = "url"
      value = "https://custom.example.com"
    },
    {
      name  = netbox_custom_field.test_json.name
      type  = "json"
      value = jsonencode({"custom": "json", "array": [1, 2, 3]})
    }
  ]

  tags = [
    {
      name = netbox_tag.test_1.name
      slug = netbox_tag.test_1.slug
    },
    {
      name = netbox_tag.test_2.name
      slug = netbox_tag.test_2.slug
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
		cfText,
		cfLongtext,
		cfInteger,
		cfBoolean,
		cfDate,
		cfUrl,
		cfJson,
		tag1, tag1Slug,
		tag2, tag2Slug,
		powerOutletName,
	)
}
