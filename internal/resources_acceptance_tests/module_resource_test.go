package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccModuleResource_basic(t *testing.T) {
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
	bayName := testutil.RandomName("tf-test-mbay")
	mtModel := testutil.RandomName("tf-test-mt")

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
				Config: testAccModuleResourceConfig_basic(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, bayName, mtModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_module.test", "id"),
					resource.TestCheckResourceAttrSet("netbox_module.test", "device"),
					resource.TestCheckResourceAttrSet("netbox_module.test", "module_bay"),
					resource.TestCheckResourceAttrSet("netbox_module.test", "module_type"),
				),
			},
			{
				ResourceName:            "netbox_module.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"device", "module_bay", "module_type"},
			},
			{
				Config:             testAccModuleResourceConfig_basic(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, bayName, mtModel),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func TestAccModuleResource_full(t *testing.T) {
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
	bayName := testutil.RandomName("tf-test-mbay-full")
	mtModel := testutil.RandomName("tf-test-mt-full")
	description := "Test module with all fields"
	updatedDescription := "Updated module description"

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
				Config: testAccModuleResourceConfig_full(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, bayName, mtModel, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_module.test", "id"),
					resource.TestCheckResourceAttr("netbox_module.test", "asset_tag", mtModel+"-asset"),
					resource.TestCheckResourceAttr("netbox_module.test", "description", description),
					resource.TestCheckResourceAttr("netbox_module.test", "status", "active"),
				),
			},
			{
				Config: testAccModuleResourceConfig_full(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, bayName, mtModel, updatedDescription),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_module.test", "description", updatedDescription),
				),
			},
		},
	})
}

func TestAccModuleResource_tagLifecycle(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-site-tags")
	siteSlug := testutil.RandomSlug("tf-test-site-tags")
	mfgName := testutil.RandomName("tf-test-mfg-tags")
	mfgSlug := testutil.RandomSlug("tf-test-mfg-tags")
	dtModel := testutil.RandomName("tf-test-dt-tags")
	dtSlug := testutil.RandomSlug("tf-test-dt-tags")
	roleName := testutil.RandomName("tf-test-role-tags")
	roleSlug := testutil.RandomSlug("tf-test-role-tags")
	deviceName := testutil.RandomName("tf-test-device-tags")
	bayName := testutil.RandomName("tf-test-mbay-tags")
	mtModel := testutil.RandomName("tf-test-mt-tags")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2Slug := testutil.RandomSlug("tag2")
	tag3Slug := testutil.RandomSlug("tag3")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterManufacturerCleanup(mfgSlug)
	cleanup.RegisterDeviceTypeCleanup(dtSlug)
	cleanup.RegisterDeviceRoleCleanup(roleSlug)
	cleanup.RegisterDeviceCleanup(deviceName)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)
	cleanup.RegisterTagCleanup(tag3Slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccModuleResourceConfig_tags(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, bayName, mtModel, tag1Slug, tag2Slug, tag3Slug, caseTag1Tag2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_module.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemAttr("netbox_module.test", "tags.*", tag1Slug),
					resource.TestCheckTypeSetElemAttr("netbox_module.test", "tags.*", tag2Slug),
				),
			},
			{
				Config: testAccModuleResourceConfig_tags(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, bayName, mtModel, tag1Slug, tag2Slug, tag3Slug, caseTag1Uscore2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_module.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemAttr("netbox_module.test", "tags.*", tag1Slug),
					resource.TestCheckTypeSetElemAttr("netbox_module.test", "tags.*", tag2Slug),
				),
			},
			{
				Config: testAccModuleResourceConfig_tags(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, bayName, mtModel, tag1Slug, tag2Slug, tag3Slug, caseTag3),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_module.test", "tags.#", "1"),
					resource.TestCheckTypeSetElemAttr("netbox_module.test", "tags.*", tag3Slug),
				),
			},
			{
				Config: testAccModuleResourceConfig_tags(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, bayName, mtModel, tag1Slug, tag2Slug, tag3Slug, tagsEmpty),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_module.test", "tags.#", "0"),
				),
			},
		},
	})
}

func TestAccModuleResource_tagOrderInvariance(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-site-tag-order")
	siteSlug := testutil.RandomSlug("tf-test-site-tag-order")
	mfgName := testutil.RandomName("tf-test-mfg-tag-order")
	mfgSlug := testutil.RandomSlug("tf-test-mfg-tag-order")
	dtModel := testutil.RandomName("tf-test-dt-tag-order")
	dtSlug := testutil.RandomSlug("tf-test-dt-tag-order")
	roleName := testutil.RandomName("tf-test-role-tag-order")
	roleSlug := testutil.RandomSlug("tf-test-role-tag-order")
	deviceName := testutil.RandomName("tf-test-device-tag-order")
	bayName := testutil.RandomName("tf-test-mbay-tag-order")
	mtModel := testutil.RandomName("tf-test-mt-tag-order")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2Slug := testutil.RandomSlug("tag2")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterManufacturerCleanup(mfgSlug)
	cleanup.RegisterDeviceTypeCleanup(dtSlug)
	cleanup.RegisterDeviceRoleCleanup(roleSlug)
	cleanup.RegisterDeviceCleanup(deviceName)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccModuleResourceConfig_tagsOrder(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, bayName, mtModel, tag1Slug, tag2Slug, caseTag1Tag2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_module.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemAttr("netbox_module.test", "tags.*", tag1Slug),
					resource.TestCheckTypeSetElemAttr("netbox_module.test", "tags.*", tag2Slug),
				),
			},
			{
				Config: testAccModuleResourceConfig_tagsOrder(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, bayName, mtModel, tag1Slug, tag2Slug, caseTag2Uscore1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_module.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemAttr("netbox_module.test", "tags.*", tag1Slug),
					resource.TestCheckTypeSetElemAttr("netbox_module.test", "tags.*", tag2Slug),
				),
			},
		},
	})
}

func TestAccConsistency_Module_LiteralNames(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("site")
	siteSlug := testutil.RandomSlug("site")
	mfgName := testutil.RandomName("mfg")
	mfgSlug := testutil.RandomSlug("mfg")
	dtModel := testutil.RandomName("dt-model")
	dtSlug := testutil.RandomSlug("dt")
	roleName := testutil.RandomName("role")
	roleSlug := testutil.RandomSlug("role")
	deviceName := testutil.RandomName("device")
	bayName := testutil.RandomName("bay")
	mtModel := testutil.RandomName("mt-model")

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
				Config: testAccModuleResourceConfig_basic(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, bayName, mtModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_module.test", "id"),
				),
			},
			{
				PlanOnly: true,
				Config:   testAccModuleResourceConfig_basic(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, bayName, mtModel),
			},
		},
	})
}

func testAccModuleResourceConfig_basic(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, bayName, mtModel string) string {
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

resource "netbox_module_bay" "test" {
  device = netbox_device.test.id
  name   = %q
}

resource "netbox_module_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = %q
}

resource "netbox_module" "test" {
  device      = netbox_device.test.id
  module_bay  = netbox_module_bay.test.id
  module_type = netbox_module_type.test.id
  asset_tag   = "%s-asset-basic"
}
`, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, bayName, mtModel, mtModel)
}

func testAccModuleResourceConfig_tags(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, bayName, mtModel, tag1Slug, tag2Slug, tag3Slug, tagCase string) string {
	var tagsConfig string
	switch tagCase {
	case caseTag1Tag2:
		tagsConfig = tagsDoubleSlug
	case caseTag1Uscore2:
		tagsConfig = tagsDoubleSlug
	case caseTag3:
		tagsConfig = tagsSingleSlug
	case tagsEmpty:
		tagsConfig = tagsEmpty
	}

	return fmt.Sprintf(`
resource "netbox_tag" "tag1" {
	name = "Tag1-%[13]s"
	slug = %[13]q
}

resource "netbox_tag" "tag2" {
	name = "Tag2-%[14]s"
	slug = %[14]q
}

resource "netbox_tag" "tag3" {
	name = "Tag3-%[15]s"
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
	name  = %[7]q
	slug  = %[8]q
	color = "aa1409"
}

resource "netbox_device" "test" {
	name        = %[9]q
	device_type = netbox_device_type.test.id
	role        = netbox_device_role.test.id
	site        = netbox_site.test.id
}

resource "netbox_module_bay" "test" {
	device = netbox_device.test.id
	name   = %[10]q
}

resource "netbox_module_type" "test" {
	manufacturer = netbox_manufacturer.test.id
	model        = %[11]q
}

resource "netbox_module" "test" {
	device      = netbox_device.test.id
	module_bay  = netbox_module_bay.test.id
	module_type = netbox_module_type.test.id
	asset_tag   = "%[12]s-asset-tags"
	%[16]s
}
`, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, bayName, mtModel, mtModel, tag1Slug, tag2Slug, tag3Slug, tagsConfig)
}

func testAccModuleResourceConfig_tagsOrder(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, bayName, mtModel, tag1Slug, tag2Slug, tagCase string) string {
	var tagsConfig string
	switch tagCase {
	case caseTag1Tag2:
		tagsConfig = tagsDoubleSlug
	case caseTag2Uscore1:
		tagsConfig = tagsDoubleSlugReversed
	}

	return fmt.Sprintf(`
resource "netbox_tag" "tag1" {
	name = "Tag1-%[13]s"
	slug = %[13]q
}

resource "netbox_tag" "tag2" {
	name = "Tag2-%[14]s"
	slug = %[14]q
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
	name  = %[7]q
	slug  = %[8]q
	color = "aa1409"
}

resource "netbox_device" "test" {
	name        = %[9]q
	device_type = netbox_device_type.test.id
	role        = netbox_device_role.test.id
	site        = netbox_site.test.id
}

resource "netbox_module_bay" "test" {
	device = netbox_device.test.id
	name   = %[10]q
}

resource "netbox_module_type" "test" {
	manufacturer = netbox_manufacturer.test.id
	model        = %[11]q
}

resource "netbox_module" "test" {
	device      = netbox_device.test.id
	module_bay  = netbox_module_bay.test.id
	module_type = netbox_module_type.test.id
	asset_tag   = "%[12]s-asset-tags"
	%[15]s
}
`, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, bayName, mtModel, mtModel, tag1Slug, tag2Slug, tagsConfig)
}

// NOTE: Custom field tests for module resource are in resources_acceptance_tests_customfields package.
func testAccModuleResourceConfig_full(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, bayName, mtModel, description string) string {
	// Using mtModel as part of asset_tag to ensure uniqueness
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

resource "netbox_module_bay" "test" {
  device = netbox_device.test.id
  name   = %q
}

resource "netbox_module_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = %q
}

resource "netbox_module" "test" {
  device      = netbox_device.test.id
  module_bay  = netbox_module_bay.test.id
  module_type = netbox_module_type.test.id
  status      = "active"
  asset_tag   = "%s-asset"
  description = %q
}
`, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, bayName, mtModel, mtModel, description)
}

func TestAccModuleResource_update(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-site-upd")
	siteSlug := testutil.RandomSlug("tf-test-site-upd")
	mfgName := testutil.RandomName("tf-test-mfg-upd")
	mfgSlug := testutil.RandomSlug("tf-test-mfg-upd")
	dtModel := testutil.RandomName("tf-test-dt-upd")
	dtSlug := testutil.RandomSlug("tf-test-dt-upd")
	roleName := testutil.RandomName("tf-test-role-upd")
	roleSlug := testutil.RandomSlug("tf-test-role-upd")
	deviceName := testutil.RandomName("tf-test-device-upd")
	bayName := testutil.RandomName("tf-test-mbay-upd")
	mtModel := testutil.RandomName("tf-test-mt-upd")

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
				Config: testAccModuleResourceConfig_serial(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, bayName, mtModel, "SERIAL1"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_module.test", "id"),
					resource.TestCheckResourceAttr("netbox_module.test", "serial", "SERIAL1"),
				),
			},
			{
				Config: testAccModuleResourceConfig_serial(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, bayName, mtModel, "SERIAL2"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_module.test", "serial", "SERIAL2"),
				),
			},
		},
	})
}

func TestAccModuleResource_externalDeletion(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-site-extdel")
	siteSlug := testutil.RandomSlug("tf-test-site-extdel")
	mfgName := testutil.RandomName("tf-test-mfg-extdel")
	mfgSlug := testutil.RandomSlug("tf-test-mfg-extdel")
	dtModel := testutil.RandomName("tf-test-dt-extdel")
	dtSlug := testutil.RandomSlug("tf-test-dt-extdel")
	roleName := testutil.RandomName("tf-test-role-extdel")
	roleSlug := testutil.RandomSlug("tf-test-role-extdel")
	deviceName := testutil.RandomName("tf-test-device-extdel")
	bayName := testutil.RandomName("tf-test-mbay-extdel")
	mtModel := testutil.RandomName("tf-test-mt-extdel")

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
				Config: testAccModuleResourceConfig_basic(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, bayName, mtModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_module.test", "id"),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}

					// List modules to find the one we created
					items, _, err := client.DcimAPI.DcimModulesList(context.Background()).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find module for external deletion: %v", err)
					}

					// Find the module by checking if it belongs to our test device
					var moduleID int32
					found := false
					for _, module := range items.Results {
						if module.Device.Name.IsSet() && *module.Device.Name.Get() == deviceName {
							moduleID = module.Id
							found = true
							break
						}
					}

					if !found {
						t.Fatalf("Module not found for device %s", deviceName)
					}

					_, err = client.DcimAPI.DcimModulesDestroy(context.Background(), moduleID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete module: %v", err)
					}
					t.Logf("Successfully externally deleted module with ID: %d", moduleID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccModuleResourceConfig_serial(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, bayName, mtModel, serial string) string {
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

resource "netbox_module_bay" "test" {
  device = netbox_device.test.id
  name   = %q
}

resource "netbox_module_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = %q
}

resource "netbox_module" "test" {
  device      = netbox_device.test.id
  module_bay  = netbox_module_bay.test.id
  module_type = netbox_module_type.test.id
	asset_tag   = "%s-asset-serial"
  serial      = %q
}
`, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, bayName, mtModel, mtModel, serial)
}

func TestAccModuleResource_removeDescriptionAndComments(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-site-mod-optional")
	siteSlug := testutil.RandomSlug("tf-test-site-mod")
	mfgName := testutil.RandomName("tf-test-mfg-mod")
	mfgSlug := testutil.RandomSlug("tf-test-mfg-mod")
	dtModel := testutil.RandomName("tf-test-devtype-mod")
	dtSlug := testutil.RandomSlug("tf-test-devtype-mod")
	roleName := testutil.RandomName("tf-test-role-mod")
	roleSlug := testutil.RandomSlug("tf-test-role-mod")
	deviceName := testutil.RandomName("tf-test-device-mod")
	bayName := testutil.RandomName("tf-test-bay-mod")
	mtModel := testutil.RandomName("tf-test-modtype-mod")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterManufacturerCleanup(mfgSlug)
	cleanup.RegisterDeviceRoleCleanup(roleSlug)
	cleanup.RegisterDeviceTypeCleanup(dtSlug)
	cleanup.RegisterDeviceCleanup(deviceName)

	testutil.TestRemoveOptionalFields(t, testutil.MultiFieldOptionalTestConfig{
		ResourceName: "netbox_module",
		BaseConfig: func() string {
			return testAccModuleResourceConfig_basic(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, bayName, mtModel)
		},
		ConfigWithFields: func() string {
			return testAccModuleResourceConfig_withDescriptionAndComments(
				siteName,
				siteSlug,
				mfgName,
				mfgSlug,
				dtModel,
				dtSlug,
				roleName,
				roleSlug,
				deviceName,
				bayName,
				mtModel,
				"Test description",
				"Test comments",
			)
		},
		OptionalFields: map[string]string{
			"description": "Test description",
			"comments":    "Test comments",
		},
		CheckDestroy: testutil.CheckModuleDestroy,
	})
}

func testAccModuleResourceConfig_withDescriptionAndComments(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, bayName, mtModel, description, comments string) string {
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
  model        = %q
  slug         = %q
  manufacturer = netbox_manufacturer.test.id
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

resource "netbox_module_bay" "test" {
  device = netbox_device.test.id
  name   = %q
}

resource "netbox_module_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = %q
}

resource "netbox_module" "test" {
  device      = netbox_device.test.id
  module_bay  = netbox_module_bay.test.id
  module_type = netbox_module_type.test.id
  asset_tag   = "%s-asset-desc"
  description = %q
  comments    = %q
}
`, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, bayName, mtModel, mtModel, description, comments)
}

func TestAccModuleResource_removeOptionalFields(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-site-opt")
	siteSlug := testutil.RandomSlug("tf-test-site-opt")
	mfgName := testutil.RandomName("tf-test-mfg-opt")
	mfgSlug := testutil.RandomSlug("tf-test-mfg-opt")
	dtModel := testutil.RandomName("tf-test-dt-opt")
	dtSlug := testutil.RandomSlug("tf-test-dt-opt")
	roleName := testutil.RandomName("tf-test-role-opt")
	roleSlug := testutil.RandomSlug("tf-test-role-opt")
	deviceName := testutil.RandomName("tf-test-device-opt")
	bayName := testutil.RandomName("tf-test-mbay-opt")
	mtModel := testutil.RandomName("tf-test-mt-opt")

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
				Config: fmt.Sprintf(`
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
  model        = %[5]q
  slug         = %[6]q
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_device_role" "test" {
  name  = %[7]q
  slug  = %[8]q
  color = "aa1409"
}

resource "netbox_device" "test" {
  name        = %[9]q
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
}

resource "netbox_module_bay" "test" {
  device = netbox_device.test.id
  name   = %[10]q
}

resource "netbox_module_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = %[11]q
}

resource "netbox_module" "test" {
  device      = netbox_device.test.id
  module_bay  = netbox_module_bay.test.id
  module_type = netbox_module_type.test.id
  asset_tag   = "%s-asset-opt"
  serial      = "SN-12345"
  status      = "active"
}
`, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, bayName, mtModel, mtModel),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_module.test", "asset_tag", mtModel+"-asset-opt"),
					resource.TestCheckResourceAttr("netbox_module.test", "serial", "SN-12345"),
					resource.TestCheckResourceAttr("netbox_module.test", "status", "active"),
				),
			},
			{
				Config: fmt.Sprintf(`
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
  model        = %[5]q
  slug         = %[6]q
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_device_role" "test" {
  name  = %[7]q
  slug  = %[8]q
  color = "aa1409"
}

resource "netbox_device" "test" {
  name        = %[9]q
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
}

resource "netbox_module_bay" "test" {
  device = netbox_device.test.id
  name   = %[10]q
}

resource "netbox_module_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = %[11]q
}

resource "netbox_module" "test" {
  device      = netbox_device.test.id
  module_bay  = netbox_module_bay.test.id
  module_type = netbox_module_type.test.id
}
`, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, bayName, mtModel),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckNoResourceAttr("netbox_module.test", "asset_tag"),
					resource.TestCheckNoResourceAttr("netbox_module.test", "serial"),
					// status is Computed - API returns "active" as default even when not set
					resource.TestCheckResourceAttr("netbox_module.test", "status", "active"),
				),
			},
		},
	})
}

func TestAccModuleResource_validationErrors(t *testing.T) {
	testutil.RunMultiValidationErrorTest(t, testutil.MultiValidationErrorTestConfig{
		ResourceName: "netbox_module",
		TestCases: map[string]testutil.ValidationErrorCase{
			"missing_device": {
				Config: func() string {
					return `
provider "netbox" {}

resource "netbox_module" "test" {
  # device missing
  module_bay = 1
  module_type = "test"
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
			"missing_module_bay": {
				Config: func() string {
					return `
provider "netbox" {}

resource "netbox_site" "test" {
  name = "test-site"
  slug = "test-site"
}

resource "netbox_device_role" "test" {
  name = "test-role"
  slug = "test-role"
}

resource "netbox_device_type" "test" {
  model        = "test-device-type"
  slug         = "test-device-type"
  manufacturer = "test-manufacturer"
}

resource "netbox_device" "test" {
  name        = "test-device"
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
}

resource "netbox_module" "test" {
  device = netbox_device.test.id
  # module_bay missing
  module_type = "test"
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
			"missing_module_type": {
				Config: func() string {
					return `
provider "netbox" {}

resource "netbox_site" "test" {
  name = "test-site"
  slug = "test-site"
}

resource "netbox_device_role" "test" {
  name = "test-role"
  slug = "test-role"
}

resource "netbox_device_type" "test" {
  model        = "test-device-type"
  slug         = "test-device-type"
  manufacturer = "test-manufacturer"
}

resource "netbox_device" "test" {
  name        = "test-device"
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
}

resource "netbox_module" "test" {
  device     = netbox_device.test.id
  module_bay = 1
  # module_type missing
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
		},
	})
}
