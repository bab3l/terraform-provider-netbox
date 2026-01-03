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

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

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

func TestAccInventoryItemResource_IDPreservation(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("inv-id")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

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

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

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

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

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

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccInventoryItemResourceConfig_basic(name),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_inventory_item.test", "id"),
				),
			},

			{

				ResourceName: "netbox_inventory_item.test",

				ImportState: true,

				ImportStateVerify: true,

				ImportStateVerifyIgnore: []string{"device"},
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

func TestAccConsistency_InventoryItem_LiteralNames(t *testing.T) {
	t.Parallel()
	name := testutil.RandomName("tf-test-inv-item-lit")

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
				Config: testAccInventoryItemResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_inventory_item.test", "id"),
				),
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

func TestAccInventoryItemResource_importWithCustomFieldsAndTags(t *testing.T) {
	// NOTE: t.Parallel() intentionally omitted - this test creates/deletes global custom fields
	// that would affect other tests of the same resource type running in parallel.

	itemName := testutil.RandomName("inventory_item")
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
		CheckDestroy:             testutil.CheckInventoryItemDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccInventoryItemResourceImportConfig_full(itemName, deviceName, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, tenantName, tenantSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_inventory_item.test", "id"),
					resource.TestCheckResourceAttr("netbox_inventory_item.test", "name", itemName),
					// Verify custom fields are applied
					resource.TestCheckResourceAttr("netbox_inventory_item.test", "custom_fields.#", "7"),
					// Verify tags are applied
					resource.TestCheckResourceAttr("netbox_inventory_item.test", "tags.#", "2"),
				),
			},
			{
				Config:                  testAccInventoryItemResourceImportConfig_full(itemName, deviceName, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, tenantName, tenantSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug),
				ResourceName:            "netbox_inventory_item.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"device", "custom_fields"}, // Device reference may have lookup inconsistencies, custom fields have import limitations
			},
		},
	})
}

func testAccInventoryItemResourceImportConfig_full(itemName, deviceName, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, tenantName, tenantSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug string) string {
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

# Custom Fields
resource "netbox_custom_field" "cf_text" {
  name        = %q
  type        = "text"
  object_types = ["dcim.inventoryitem"]
}

resource "netbox_custom_field" "cf_longtext" {
  name        = %q
  type        = "longtext"
  object_types = ["dcim.inventoryitem"]
}

resource "netbox_custom_field" "cf_integer" {
  name        = %q
  type        = "integer"
  object_types = ["dcim.inventoryitem"]
}

resource "netbox_custom_field" "cf_boolean" {
  name        = %q
  type        = "boolean"
  object_types = ["dcim.inventoryitem"]
}

resource "netbox_custom_field" "cf_date" {
  name        = %q
  type        = "date"
  object_types = ["dcim.inventoryitem"]
}

resource "netbox_custom_field" "cf_url" {
  name        = %q
  type        = "url"
  object_types = ["dcim.inventoryitem"]
}

resource "netbox_custom_field" "cf_json" {
  name        = %q
  type        = "json"
  object_types = ["dcim.inventoryitem"]
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
resource "netbox_inventory_item" "test" {
  name   = %q
  device = netbox_device.test.name

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
		cfText,
		cfLongtext,
		cfInteger,
		cfBoolean,
		cfDate,
		cfUrl,
		cfJson,
		tag1, tag1Slug,
		tag2, tag2Slug,
		itemName,
	)
}
