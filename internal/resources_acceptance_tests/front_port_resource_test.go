package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
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

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterDeviceCleanup(deviceName)
	cleanup.RegisterDeviceTypeCleanup(dtSlug)
	cleanup.RegisterManufacturerCleanup(mfgSlug)
	cleanup.RegisterDeviceRoleCleanup(roleSlug)
	cleanup.RegisterSiteCleanup(siteSlug)

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
			{
				Config:             testAccFrontPortResourceConfig_basic(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, rearPortName, frontPortName),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
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

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterDeviceCleanup(deviceName)
	cleanup.RegisterDeviceTypeCleanup(dtSlug)
	cleanup.RegisterManufacturerCleanup(mfgSlug)
	cleanup.RegisterDeviceRoleCleanup(roleSlug)
	cleanup.RegisterSiteCleanup(siteSlug)

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
			{
				Config:             testAccFrontPortResourceConfig_full(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, rearPortName, frontPortName),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
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

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterDeviceCleanup(deviceName)
	cleanup.RegisterDeviceTypeCleanup(dtSlug)
	cleanup.RegisterManufacturerCleanup(mfgSlug)
	cleanup.RegisterDeviceRoleCleanup(roleSlug)
	cleanup.RegisterSiteCleanup(siteSlug)

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
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccFrontPortResource_tagLifecycle(t *testing.T) {
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
	tag1Name := testutil.RandomName("tag1")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2Name := testutil.RandomName("tag2")
	tag2Slug := testutil.RandomSlug("tag2")
	tag3Name := testutil.RandomName("tag3")
	tag3Slug := testutil.RandomSlug("tag3")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterDeviceCleanup(deviceName)
	cleanup.RegisterDeviceTypeCleanup(dtSlug)
	cleanup.RegisterManufacturerCleanup(mfgSlug)
	cleanup.RegisterDeviceRoleCleanup(roleSlug)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)
	cleanup.RegisterTagCleanup(tag3Slug)

	testutil.RunTagLifecycleTest(t, testutil.TagLifecycleTestConfig{
		ResourceName: "netbox_front_port",
		ConfigWithoutTags: func() string {
			return testAccFrontPortResourceConfig_tags(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, rearPortName, frontPortName, tag1Name, tag1Slug, tag2Name, tag2Slug, tag3Name, tag3Slug, "none")
		},
		ConfigWithTags: func() string {
			return testAccFrontPortResourceConfig_tags(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, rearPortName, frontPortName, tag1Name, tag1Slug, tag2Name, tag2Slug, tag3Name, tag3Slug, "tag1_tag2")
		},
		ConfigWithDifferentTags: func() string {
			return testAccFrontPortResourceConfig_tags(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, rearPortName, frontPortName, tag1Name, tag1Slug, tag2Name, tag2Slug, tag3Name, tag3Slug, "tag3")
		},
		ExpectedTagCount:          2,
		ExpectedDifferentTagCount: 1,
		CheckDestroy:              testutil.CheckFrontPortDestroy,
	})
}

func TestAccFrontPortResource_tagOrderInvariance(t *testing.T) {
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
	tag1Name := testutil.RandomName("tag1")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2Name := testutil.RandomName("tag2")
	tag2Slug := testutil.RandomSlug("tag2")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterDeviceCleanup(deviceName)
	cleanup.RegisterDeviceTypeCleanup(dtSlug)
	cleanup.RegisterManufacturerCleanup(mfgSlug)
	cleanup.RegisterDeviceRoleCleanup(roleSlug)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)

	testutil.RunTagOrderTest(t, testutil.TagOrderTestConfig{
		ResourceName: "netbox_front_port",
		ConfigWithTagsOrderA: func() string {
			return testAccFrontPortResourceConfig_tagsOrder(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, rearPortName, frontPortName, tag1Name, tag1Slug, tag2Name, tag2Slug, "tag1_tag2")
		},
		ConfigWithTagsOrderB: func() string {
			return testAccFrontPortResourceConfig_tagsOrder(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, rearPortName, frontPortName, tag1Name, tag1Slug, tag2Name, tag2Slug, "tag2_tag1")
		},
		ExpectedTagCount: 2,
		CheckDestroy:     testutil.CheckFrontPortDestroy,
	})
}

func testAccFrontPortResourceConfig_basic(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, rearPortName, frontPortName string) string {
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

resource "netbox_rear_port" "rear" {
  device = netbox_device.test.name
  name = %[10]q
  type = "8p8c"
  positions = 1
}

resource "netbox_front_port" "test" {
  device = netbox_device.test.name
  name = %[11]q
  type = "8p8c"
  rear_port = netbox_rear_port.rear.id
  rear_port_position = 1
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

func testAccFrontPortResourceConfig_tags(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, rearPortName, frontPortName, tag1Name, tag1Slug, tag2Name, tag2Slug, tag3Name, tag3Slug, tagCase string) string {
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
  name = %[12]q
  slug = %[13]q
}

resource "netbox_tag" "tag2" {
  name = %[14]q
  slug = %[15]q
}

resource "netbox_tag" "tag3" {
  name = %[16]q
  slug = %[17]q
}

resource "netbox_site" "test" {
  name   = %[1]q
  slug   = %[2]q
  status = "active"
}

resource "netbox_manufacturer" "test" {
  name = %[3]q
  slug = %[4]q
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = %[5]q
  slug         = %[6]q
}

resource "netbox_device_role" "test" {
  name = %[7]q
  slug = %[8]q
}

resource "netbox_device" "test" {
  name        = %[9]q
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
}

resource "netbox_rear_port" "test" {
  device = netbox_device.test.id
  name   = %[10]q
  type   = "8p8c"
}

resource "netbox_front_port" "test" {
  device             = netbox_device.test.id
  name               = %[11]q
  type               = "8p8c"
  rear_port          = netbox_rear_port.test.id
  rear_port_position = 1
  %[18]s
}
`, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, rearPortName, frontPortName, tag1Name, tag1Slug, tag2Name, tag2Slug, tag3Name, tag3Slug, tagsList)
}

func testAccFrontPortResourceConfig_tagsOrder(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, rearPortName, frontPortName, tag1Name, tag1Slug, tag2Name, tag2Slug, tagOrder string) string {
	var tagsOrder string
	switch tagOrder {
	case caseTag1Uscore2:
		tagsOrder = tagsDoubleSlug
	case caseTag2Uscore1:
		tagsOrder = tagsDoubleSlugReversed
	}

	return fmt.Sprintf(`
resource "netbox_tag" "tag1" {
  name = %[12]q
  slug = %[13]q
}

resource "netbox_tag" "tag2" {
  name = %[14]q
  slug = %[15]q
}

resource "netbox_site" "test" {
  name   = %[1]q
  slug   = %[2]q
  status = "active"
}

resource "netbox_manufacturer" "test" {
  name = %[3]q
  slug = %[4]q
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = %[5]q
  slug         = %[6]q
}

resource "netbox_device_role" "test" {
  name = %[7]q
  slug = %[8]q
}

resource "netbox_device" "test" {
  name        = %[9]q
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
}

resource "netbox_rear_port" "test" {
  device = netbox_device.test.id
  name   = %[10]q
  type   = "8p8c"
}

resource "netbox_front_port" "test" {
  device             = netbox_device.test.id
  name               = %[11]q
  type               = "8p8c"
  rear_port          = netbox_rear_port.test.id
  rear_port_position = 1
  %[16]s
}
`, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, rearPortName, frontPortName, tag1Name, tag1Slug, tag2Name, tag2Slug, tagsOrder)
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

func TestAccFrontPortResource_removeOptionalFields(t *testing.T) {
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
	rearPortName := testutil.RandomName("tf-test-rear-rem")
	portName := testutil.RandomName("tf-test-fp-rem")
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
				Config: testAccFrontPortResourceConfig_withLabel(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, rearPortName, portName, testLabel, testDescription),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_front_port.test", "name", portName),
					resource.TestCheckResourceAttr("netbox_front_port.test", "label", testLabel),
					resource.TestCheckResourceAttr("netbox_front_port.test", "description", testDescription),
					resource.TestCheckResourceAttr("netbox_front_port.test", "color", "ff0000"),
					resource.TestCheckResourceAttr("netbox_front_port.test", "mark_connected", "true"),
					resource.TestCheckResourceAttr("netbox_front_port.test", "rear_port_position", "2"),
				),
			},
			{
				Config: testAccFrontPortResourceConfig_basic(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, rearPortName, portName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_front_port.test", "name", portName),
					resource.TestCheckNoResourceAttr("netbox_front_port.test", "label"),
					resource.TestCheckNoResourceAttr("netbox_front_port.test", "description"),
					resource.TestCheckNoResourceAttr("netbox_front_port.test", "color"),
					resource.TestCheckResourceAttr("netbox_front_port.test", "mark_connected", "false"),
					resource.TestCheckResourceAttr("netbox_front_port.test", "rear_port_position", "1"),
				),
			},
		},
	})
}

func testAccFrontPortResourceConfig_withLabel(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, rearPortName, portName, label, description string) string {
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

resource "netbox_rear_port" "rear" {
  device = netbox_device.test.name
  name = %[10]q
  type = "8p8c"
  positions = 2
}

resource "netbox_front_port" "test" {
  device = netbox_device.test.name
  name = %[11]q
  type = "8p8c"
  rear_port = netbox_rear_port.rear.id
  rear_port_position = 2
  label = %[12]q
  description = %[13]q
  color = "ff0000"
  mark_connected = true
}
`, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, rearPortName, portName, label, description)
}

func TestAccFrontPortResource_validationErrors(t *testing.T) {
	testutil.RunMultiValidationErrorTest(t, testutil.MultiValidationErrorTestConfig{
		ResourceName: "netbox_front_port",
		TestCases: map[string]testutil.ValidationErrorCase{
			"missing_device": {
				Config: func() string {
					return `
resource "netbox_front_port" "test" {
  name = "Front1"
  type = "8p8c"
  rear_port = "1"
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

resource "netbox_front_port" "test" {
  device = netbox_device.test.id
  type = "8p8c"
  rear_port = "1"
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
			"missing_type": {
				Config: func() string {
					return `
resource "netbox_device" "test" {
  name = "test-device"
  device_type = "1"
  role = "1"
  site = "1"
  status = "active"
}

resource "netbox_front_port" "test" {
  device = netbox_device.test.id
  name = "Front1"
  rear_port = "1"
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
			"missing_rear_port": {
				Config: func() string {
					return `
resource "netbox_device" "test" {
  name = "test-device"
  device_type = "1"
  role = "1"
  site = "1"
  status = "active"
}

resource "netbox_front_port" "test" {
  device = netbox_device.test.id
  name = "Front1"
  type = "8p8c"
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
			"invalid_device_reference": {
				Config: func() string {
					return `
resource "netbox_front_port" "test" {
  device = "99999"
  name = "Front1"
  type = "8p8c"
  rear_port = "1"
}
`
				},
				ExpectedError: testutil.ErrPatternNotFound,
			},
			"invalid_rear_port_reference": {
				Config: func() string {
					return `
resource "netbox_device" "test" {
  name = "test-device"
  device_type = "1"
  role = "1"
  site = "1"
  status = "active"
}

resource "netbox_front_port" "test" {
  device = netbox_device.test.id
  name = "Front1"
  type = "8p8c"
  rear_port = "99999"
}
`
				},
				ExpectedError: testutil.ErrPatternNotFound,
			},
		},
	})
}

// NOTE: Custom field tests for front_port resource are in resources_acceptance_tests_customfields package
