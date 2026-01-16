package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccInventoryItemResource_basic(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-inv-item")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(testutil.RandomSlug("site"))
	cleanup.RegisterManufacturerCleanup(testutil.RandomSlug("mfr"))
	cleanup.RegisterDeviceTypeCleanup(testutil.RandomSlug("device"))
	cleanup.RegisterDeviceRoleCleanup(testutil.RandomSlug("role"))
	cleanup.RegisterDeviceCleanup(name + "-device")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInventoryItemResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_inventory_item.test", "id"),
					resource.TestCheckResourceAttr("netbox_inventory_item.test", "name", name),
					resource.TestCheckResourceAttrSet("netbox_inventory_item.test", "device"),
				),
			},
		},
	})
}

func TestAccInventoryItemResource_full(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-inv-item-full")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(testutil.RandomSlug("site"))
	cleanup.RegisterManufacturerCleanup(testutil.RandomSlug("mfr"))
	cleanup.RegisterDeviceTypeCleanup(testutil.RandomSlug("device"))
	cleanup.RegisterDeviceRoleCleanup(testutil.RandomSlug("role"))
	cleanup.RegisterDeviceCleanup(name + "-device")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInventoryItemResourceConfig_full(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_inventory_item.test", "id"),
					resource.TestCheckResourceAttr("netbox_inventory_item.test", "name", name),
					resource.TestCheckResourceAttr("netbox_inventory_item.test", "label", "Inventory Label"),
					resource.TestCheckResourceAttr("netbox_inventory_item.test", "serial", "SN-12345"),
					resource.TestCheckResourceAttr("netbox_inventory_item.test", "asset_tag", name+"-asset-tag"),
					resource.TestCheckResourceAttr("netbox_inventory_item.test", "description", "Test inventory item"),
				),
			},
		},
	})
}

func TestAccInventoryItemResource_update(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-inv-item-update")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(testutil.RandomSlug("site"))
	cleanup.RegisterManufacturerCleanup(testutil.RandomSlug("mfr"))
	cleanup.RegisterDeviceTypeCleanup(testutil.RandomSlug("device"))
	cleanup.RegisterDeviceRoleCleanup(testutil.RandomSlug("role"))
	cleanup.RegisterDeviceCleanup(name + "-device")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInventoryItemResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_inventory_item.test", "id"),
					resource.TestCheckResourceAttr("netbox_inventory_item.test", "name", name),
				),
			},
			{
				Config: testAccInventoryItemResourceConfig_full(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_inventory_item.test", "id"),
					resource.TestCheckResourceAttr("netbox_inventory_item.test", "label", "Inventory Label"),
					resource.TestCheckResourceAttr("netbox_inventory_item.test", "serial", "SN-12345"),
				),
			},
		},
	})
}

func TestAccInventoryItemResource_import(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-inv-item")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(testutil.RandomSlug("site"))
	cleanup.RegisterManufacturerCleanup(testutil.RandomSlug("mfr"))
	cleanup.RegisterDeviceTypeCleanup(testutil.RandomSlug("device"))
	cleanup.RegisterDeviceRoleCleanup(testutil.RandomSlug("role"))
	cleanup.RegisterDeviceCleanup(name + "-device")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInventoryItemResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_inventory_item.test", "id"),
				),
			},
			{
				ResourceName:            "netbox_inventory_item.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"device"},
			},
		},
	})
}

func TestAccInventoryItemResource_tagLifecycle(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-inv-item-tags")
	siteName := testutil.RandomName("tf-test-site-inv-item-tags")
	siteSlug := testutil.RandomSlug("tf-test-site-inv-item-tags")
	manufacturerName := testutil.RandomName("tf-test-mfr-inv-item-tags")
	manufacturerSlug := testutil.RandomSlug("tf-test-mfr-inv-item-tags")
	deviceTypeName := testutil.RandomName("tf-test-devtype-inv-item-tags")
	deviceTypeSlug := testutil.RandomSlug("tf-test-devtype-inv-item-tags")
	deviceRoleName := testutil.RandomName("tf-test-role-inv-item-tags")
	deviceRoleSlug := testutil.RandomSlug("tf-test-role-inv-item-tags")
	deviceName := testutil.RandomName("tf-test-device-inv-item-tags")

	tag1Slug := testutil.RandomSlug("tag1")
	tag2Slug := testutil.RandomSlug("tag2")
	tag3Slug := testutil.RandomSlug("tag3")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)
	cleanup.RegisterDeviceRoleCleanup(deviceRoleSlug)
	cleanup.RegisterDeviceCleanup(deviceName)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)
	cleanup.RegisterTagCleanup(tag3Slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInventoryItemResourceConfig_tags(siteName, siteSlug, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceRoleName, deviceRoleSlug, deviceName, name, tag1Slug, tag2Slug, tag3Slug, caseTag1Tag2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_inventory_item.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_inventory_item.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag1-%s", tag1Slug),
						"slug": tag1Slug,
					}),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_inventory_item.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag2-%s", tag2Slug),
						"slug": tag2Slug,
					}),
				),
			},
			{
				Config: testAccInventoryItemResourceConfig_tags(siteName, siteSlug, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceRoleName, deviceRoleSlug, deviceName, name, tag1Slug, tag2Slug, tag3Slug, caseTag1Uscore2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_inventory_item.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_inventory_item.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag1-%s", tag1Slug),
						"slug": tag1Slug,
					}),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_inventory_item.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag2-%s", tag2Slug),
						"slug": tag2Slug,
					}),
				),
			},
			{
				Config: testAccInventoryItemResourceConfig_tags(siteName, siteSlug, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceRoleName, deviceRoleSlug, deviceName, name, tag1Slug, tag2Slug, tag3Slug, caseTag3),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_inventory_item.test", "tags.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_inventory_item.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag3-%s", tag3Slug),
						"slug": tag3Slug,
					}),
				),
			},
			{
				Config: testAccInventoryItemResourceConfig_tags(siteName, siteSlug, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceRoleName, deviceRoleSlug, deviceName, name, tag1Slug, tag2Slug, tag3Slug, tagsEmpty),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_inventory_item.test", "tags.#", "0"),
				),
			},
		},
	})
}

func TestAccInventoryItemResource_tagOrderInvariance(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-inv-item-tag-order")
	siteName := testutil.RandomName("tf-test-site-inv-item-tag-order")
	siteSlug := testutil.RandomSlug("tf-test-site-inv-item-tag-order")
	manufacturerName := testutil.RandomName("tf-test-mfr-inv-item-tag-order")
	manufacturerSlug := testutil.RandomSlug("tf-test-mfr-inv-item-tag-order")
	deviceTypeName := testutil.RandomName("tf-test-devtype-inv-item-tag-order")
	deviceTypeSlug := testutil.RandomSlug("tf-test-devtype-inv-item-tag-order")
	deviceRoleName := testutil.RandomName("tf-test-role-inv-item-tag-order")
	deviceRoleSlug := testutil.RandomSlug("tf-test-role-inv-item-tag-order")
	deviceName := testutil.RandomName("tf-test-device-inv-item-tag-order")

	tag1Slug := testutil.RandomSlug("tag1")
	tag2Slug := testutil.RandomSlug("tag2")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)
	cleanup.RegisterDeviceRoleCleanup(deviceRoleSlug)
	cleanup.RegisterDeviceCleanup(deviceName)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInventoryItemResourceConfig_tagsOrder(siteName, siteSlug, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceRoleName, deviceRoleSlug, deviceName, name, tag1Slug, tag2Slug, caseTag1Tag2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_inventory_item.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_inventory_item.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag1-%s", tag1Slug),
						"slug": tag1Slug,
					}),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_inventory_item.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag2-%s", tag2Slug),
						"slug": tag2Slug,
					}),
				),
			},
			{
				Config: testAccInventoryItemResourceConfig_tagsOrder(siteName, siteSlug, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceRoleName, deviceRoleSlug, deviceName, name, tag1Slug, tag2Slug, caseTag2Uscore1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_inventory_item.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_inventory_item.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag1-%s", tag1Slug),
						"slug": tag1Slug,
					}),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_inventory_item.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag2-%s", tag2Slug),
						"slug": tag2Slug,
					}),
				),
			},
		},
	})
}

func testAccInventoryItemResourceConfig_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "netbox_inventory_item" "test" {
  device = netbox_device.test.id
  name   = %q
}
`, testAccInventoryItemResourcePrereqs(name), name)
}

func testAccInventoryItemResourceConfig_full(name string) string {
	return fmt.Sprintf(`

%s

resource "netbox_inventory_item" "test" {
  device      = netbox_device.test.id
  name        = %q
  label       = "Inventory Label"
  serial      = "SN-12345"
  asset_tag   = %q
  description = "Test inventory item"
}
`, testAccInventoryItemResourcePrereqs(name), name, name+"-asset-tag")
}

func testAccInventoryItemResourcePrereqs(name string) string {
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

func testAccInventoryItemResourceConfig_tags(siteName, siteSlug, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceRoleName, deviceRoleSlug, deviceName, itemName, tag1Slug, tag2Slug, tag3Slug, tagCase string) string {
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
	name = %[1]q
	slug = %[2]q
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
	site        = netbox_site.test.id
	name        = %[9]q
	device_type = netbox_device_type.test.id
	role        = netbox_device_role.test.id
	status      = "offline"
}

resource "netbox_inventory_item" "test" {
	device = netbox_device.test.id
	name   = %[10]q
	%[14]s
}
`, siteName, siteSlug, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceRoleName, deviceRoleSlug, deviceName, itemName, tag1Slug, tag2Slug, tag3Slug, tagsConfig)
}

func testAccInventoryItemResourceConfig_tagsOrder(siteName, siteSlug, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceRoleName, deviceRoleSlug, deviceName, itemName, tag1Slug, tag2Slug, tagCase string) string {
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
	name = %[1]q
	slug = %[2]q
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
	site        = netbox_site.test.id
	name        = %[9]q
	device_type = netbox_device_type.test.id
	role        = netbox_device_role.test.id
	status      = "offline"
}

resource "netbox_inventory_item" "test" {
	device = netbox_device.test.id
	name   = %[10]q
	%[13]s
}
`, siteName, siteSlug, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceRoleName, deviceRoleSlug, deviceName, itemName, tag1Slug, tag2Slug, tagsConfig)
}

func TestAccConsistency_InventoryItem_LiteralNames(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-inv-item-lit")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(name + "-site")
	cleanup.RegisterManufacturerCleanup(name + "-mfr")
	cleanup.RegisterDeviceTypeCleanup(name + "-model")
	cleanup.RegisterDeviceRoleCleanup(name + "-role")
	cleanup.RegisterDeviceCleanup(name + "-device")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInventoryItemConsistencyLiteralNamesConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_inventory_item.test", "id"),
					resource.TestCheckResourceAttr("netbox_inventory_item.test", "name", name),
				),
			},
			{
				Config:   testAccInventoryItemConsistencyLiteralNamesConfig(name),
				PlanOnly: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_inventory_item.test", "id"),
				),
			},
		},
	})
}

func TestAccInventoryItemResource_externalDeletion(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-inv-item-ext-del")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(name)
	cleanup.RegisterManufacturerCleanup(name + "-mfr")
	cleanup.RegisterDeviceTypeCleanup(name + "-model")
	cleanup.RegisterDeviceRoleCleanup(name)
	cleanup.RegisterDeviceCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInventoryItemResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_inventory_item.test", "id"),
					resource.TestCheckResourceAttr("netbox_inventory_item.test", "name", name),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					items, _, err := client.DcimAPI.DcimInventoryItemsList(context.Background()).NameIc([]string{name}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find inventory_item for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.DcimAPI.DcimInventoryItemsDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete inventory_item: %v", err)
					}
					t.Logf("Successfully externally deleted inventory_item with ID: %d", itemID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccInventoryItemConsistencyLiteralNamesConfig(name string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name = "%s-site"
  slug = "%s-site"
}

resource "netbox_manufacturer" "test" {
  name = "%s-mfr"
  slug = "%s-mfr"
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = "%s-model"
  slug         = "%s-model"
}

resource "netbox_device_role" "test" {
  name = "%s-role"
  slug = "%s-role"
}

resource "netbox_device" "test" {
  site        = netbox_site.test.id
  name        = "%s-device"
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  status      = "offline"
}

resource "netbox_inventory_item" "test" {
  device = netbox_device.test.id
  name   = %q
}
`, name, name, name, name, name, name, name, name, name, name)
}

// TestAccInventoryItemResource_removeOptionalFields tests that optional fields
// can be successfully added, removed, and re-added from the configuration.
func TestAccInventoryItemResource_removeOptionalFields(t *testing.T) {
	t.Parallel()

	itemName := testutil.RandomName("tf-test-invitem-rem")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(testutil.RandomSlug("site"))
	cleanup.RegisterManufacturerCleanup(testutil.RandomSlug("mfr"))
	cleanup.RegisterDeviceTypeCleanup(testutil.RandomSlug("device"))
	cleanup.RegisterDeviceRoleCleanup(testutil.RandomSlug("role"))
	cleanup.RegisterDeviceCleanup(itemName + "-device")

	testFields := map[string]string{
		"label":     "Test Label",
		"serial":    "SN-12345",
		"asset_tag": itemName + "-asset",
		"part_id":   "PART-789",
		// Note: "discovered" field is omitted as it has API default behavior that cannot be cleared
	}

	testutil.TestRemoveOptionalFields(t, testutil.MultiFieldOptionalTestConfig{
		ResourceName: "netbox_inventory_item",
		BaseConfig: func() string {
			return testAccInventoryItemResourceConfig_removeOptionalFields_base(itemName)
		},
		ConfigWithFields: func() string {
			return testAccInventoryItemResourceConfig_removeOptionalFields_withFields(itemName, testFields)
		},
		OptionalFields: testFields,
		RequiredFields: map[string]string{
			"name": itemName,
		},
	})
}

func testAccInventoryItemResourceConfig_removeOptionalFields_base(name string) string {
	return testAccInventoryItemResourcePrereqs(name) + fmt.Sprintf(`
resource "netbox_inventory_item" "test" {
  device = netbox_device.test.id
  name = %q
}
`, name)
}

func testAccInventoryItemResourceConfig_removeOptionalFields_withFields(name string, fields map[string]string) string {
	config := testAccInventoryItemResourcePrereqs(name) + fmt.Sprintf(`
resource "netbox_inventory_item" "test" {
  device = netbox_device.test.id
  name = %q
`, name)

	if label, ok := fields["label"]; ok {
		config += fmt.Sprintf("  label = %q\n", label)
	}
	if serial, ok := fields["serial"]; ok {
		config += fmt.Sprintf("  serial = %q\n", serial)
	}
	if assetTag, ok := fields["asset_tag"]; ok {
		config += fmt.Sprintf("  asset_tag = %q\n", assetTag)
	}
	if partID, ok := fields["part_id"]; ok {
		config += fmt.Sprintf("  part_id = %q\n", partID)
	}
	if discovered, ok := fields["discovered"]; ok {
		config += fmt.Sprintf("  discovered = %s\n", discovered)
	}

	config += "}\n"
	return config
}

func TestAccInventoryItemResource_validationErrors(t *testing.T) {
	testutil.RunMultiValidationErrorTest(t, testutil.MultiValidationErrorTestConfig{
		ResourceName: "netbox_inventory_item",
		TestCases: map[string]testutil.ValidationErrorCase{
			"missing_device": {
				Config: func() string {
					return `
provider "netbox" {}

resource "netbox_inventory_item" "test" {
  # device missing
  name = "test-inventory"
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
			"missing_name": {
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

resource "netbox_inventory_item" "test" {
  device = netbox_device.test.id
  # name missing
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
		},
	})
}
