package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var _ = testAccPowerOutletResourceConfig_withLabel

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
				ResourceName:            "netbox_power_outlet.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"device"},
			},
		},
	})
}

func TestAccPowerOutletResource_tagLifecycle(t *testing.T) {
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
	powerOutletName := testutil.RandomName("tf-test-po-tags")
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
				Config: testAccPowerOutletResourceConfig_tags(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, powerOutletName, tag1Slug, tag2Slug, tag3Slug, caseTag1Tag2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_power_outlet.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_power_outlet.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag1-%s", tag1Slug),
						"slug": tag1Slug,
					}),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_power_outlet.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag2-%s", tag2Slug),
						"slug": tag2Slug,
					}),
				),
			},
			{
				Config: testAccPowerOutletResourceConfig_tags(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, powerOutletName, tag1Slug, tag2Slug, tag3Slug, caseTag1Uscore2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_power_outlet.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_power_outlet.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag1-%s", tag1Slug),
						"slug": tag1Slug,
					}),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_power_outlet.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag2-%s", tag2Slug),
						"slug": tag2Slug,
					}),
				),
			},
			{
				Config: testAccPowerOutletResourceConfig_tags(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, powerOutletName, tag1Slug, tag2Slug, tag3Slug, caseTag3),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_power_outlet.test", "tags.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_power_outlet.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag3-%s", tag3Slug),
						"slug": tag3Slug,
					}),
				),
			},
			{
				Config: testAccPowerOutletResourceConfig_tags(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, powerOutletName, tag1Slug, tag2Slug, tag3Slug, tagsEmpty),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_power_outlet.test", "tags.#", "0"),
				),
			},
		},
	})
}

func TestAccPowerOutletResource_tagOrderInvariance(t *testing.T) {
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
	powerOutletName := testutil.RandomName("tf-test-po-tag-order")
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
				Config: testAccPowerOutletResourceConfig_tagsOrder(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, powerOutletName, tag1Slug, tag2Slug, caseTag1Tag2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_power_outlet.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_power_outlet.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag1-%s", tag1Slug),
						"slug": tag1Slug,
					}),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_power_outlet.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag2-%s", tag2Slug),
						"slug": tag2Slug,
					}),
				),
			},
			{
				Config: testAccPowerOutletResourceConfig_tagsOrder(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, powerOutletName, tag1Slug, tag2Slug, caseTag2Uscore1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_power_outlet.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_power_outlet.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag1-%s", tag1Slug),
						"slug": tag1Slug,
					}),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_power_outlet.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag2-%s", tag2Slug),
						"slug": tag2Slug,
					}),
				),
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
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
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

func testAccPowerOutletResourceConfig_tags(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, powerOutletName, tag1Slug, tag2Slug, tag3Slug, tagCase string) string {
	var tagsConfig string
	switch tagCase {
	case caseTag1Tag2:
		tagsConfig = tagsDoubleNested
	case caseTag1Uscore2:
		tagsConfig = tagsDoubleNested
	case caseTag3:
		tagsConfig = tagsSingleNested
	case tagsEmpty:
		tagsConfig = tagsEmpty
	}

	return fmt.Sprintf(`
resource "netbox_tag" "tag1" {
	name = "Tag1-%[11]s"
	slug = %[11]q
}

resource "netbox_tag" "tag2" {
	name = "Tag2-%[12]s"
	slug = %[12]q
}

resource "netbox_tag" "tag3" {
	name = "Tag3-%[13]s"
	slug = %[13]q
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

resource "netbox_power_outlet" "test" {
	device = netbox_device.test.id
	name   = %[10]q
	%[14]s
}
`, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, powerOutletName, tag1Slug, tag2Slug, tag3Slug, tagsConfig)
}

func testAccPowerOutletResourceConfig_tagsOrder(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, powerOutletName, tag1Slug, tag2Slug, tagCase string) string {
	var tagsConfig string
	switch tagCase {
	case caseTag1Tag2:
		tagsConfig = tagsDoubleNested
	case caseTag2Uscore1:
		tagsConfig = tagsDoubleNestedReversed
	}

	return fmt.Sprintf(`
resource "netbox_tag" "tag1" {
	name = "Tag1-%[11]s"
	slug = %[11]q
}

resource "netbox_tag" "tag2" {
	name = "Tag2-%[12]s"
	slug = %[12]q
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

resource "netbox_power_outlet" "test" {
	device = netbox_device.test.id
	name   = %[10]q
	%[13]s
}
`, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, powerOutletName, tag1Slug, tag2Slug, tagsConfig)
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

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)
	cleanup.RegisterDeviceRoleCleanup(roleSlug)
	cleanup.RegisterDeviceCleanup(deviceName)

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

func TestAccPowerOutletResource_removeOptionalFields(t *testing.T) {
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
	portName := testutil.RandomName("tf-test-po-rem")
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
				Config: testAccPowerOutletResourceConfig_withAllFields(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, portName, testLabel, testDescription),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_power_outlet.test", "name", portName),
					resource.TestCheckResourceAttr("netbox_power_outlet.test", "label", testLabel),
					resource.TestCheckResourceAttr("netbox_power_outlet.test", "description", testDescription),
					resource.TestCheckResourceAttr("netbox_power_outlet.test", "type", "iec-60320-c13"),
					resource.TestCheckResourceAttrSet("netbox_power_outlet.test", "power_port"),
					resource.TestCheckResourceAttr("netbox_power_outlet.test", "feed_leg", "A"),
					resource.TestCheckResourceAttr("netbox_power_outlet.test", "mark_connected", "true"),
				),
			},
			{
				Config: testAccPowerOutletResourceConfig_basic(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, portName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_power_outlet.test", "name", portName),
					resource.TestCheckNoResourceAttr("netbox_power_outlet.test", "label"),
					resource.TestCheckNoResourceAttr("netbox_power_outlet.test", "description"),
					resource.TestCheckNoResourceAttr("netbox_power_outlet.test", "type"),
					resource.TestCheckNoResourceAttr("netbox_power_outlet.test", "power_port"),
					resource.TestCheckNoResourceAttr("netbox_power_outlet.test", "feed_leg"),
					resource.TestCheckResourceAttr("netbox_power_outlet.test", "mark_connected", "false"),
				),
			},
		},
	})
}

func testAccPowerOutletResourceConfig_withAllFields(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, portName, label, description string) string {
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

resource "netbox_power_port" "test" {
  device = netbox_device.test.name
  name   = "%[10]s-pp"
}

resource "netbox_power_outlet" "test" {
  device        = netbox_device.test.name
  name          = %[10]q
  label         = %[11]q
  description   = %[12]q
  type          = "iec-60320-c13"
  power_port    = tonumber(netbox_power_port.test.id)
  feed_leg      = "A"
  mark_connected = true
}
`, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, portName, label, description)
}

func testAccPowerOutletResourceConfig_withLabel(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, portName, label, description string) string {
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

resource "netbox_power_outlet" "test" {
  device = netbox_device.test.name
  name = %[10]q
  label = %[11]q
  description = %[12]q
}
`, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, portName, label, description)
}

func TestAccPowerOutletResource_validationErrors(t *testing.T) {
	testutil.RunMultiValidationErrorTest(t, testutil.MultiValidationErrorTestConfig{
		ResourceName: "netbox_power_outlet",
		TestCases: map[string]testutil.ValidationErrorCase{
			"missing_device": {
				Config: func() string {
					return `
resource "netbox_power_outlet" "test" {
  name = "Outlet1"
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

resource "netbox_power_outlet" "test" {
  device = netbox_device.test.id
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
			"invalid_device_reference": {
				Config: func() string {
					return `
resource "netbox_power_outlet" "test" {
  device = "99999"
  name = "Outlet1"
}
`
				},
				ExpectedError: testutil.ErrPatternNotFound,
			},
		},
	})
}

// NOTE: Custom field tests for power_outlet resource are in resources_acceptance_tests_customfields package
