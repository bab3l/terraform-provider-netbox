package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccInventoryItemRoleResource_basic(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-inv-role")
	slug := testutil.RandomSlug("role")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInventoryItemRoleResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_inventory_item_role.test", "id"),
					resource.TestCheckResourceAttr("netbox_inventory_item_role.test", "name", name),
					resource.TestCheckResourceAttrSet("netbox_inventory_item_role.test", "slug"),
				),
			},
		},
	})
}

func TestAccInventoryItemRoleResource_full(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-inv-role-full")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInventoryItemRoleResourceConfig_full(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_inventory_item_role.test", "id"),
					resource.TestCheckResourceAttr("netbox_inventory_item_role.test", "name", name),
					resource.TestCheckResourceAttr("netbox_inventory_item_role.test", "color", "e41e22"),
					resource.TestCheckResourceAttr("netbox_inventory_item_role.test", "description", "Test role description"),
				),
			},
		},
	})
}

func TestAccInventoryItemRoleResource_update(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-inv-role-update")
	slug := testutil.RandomSlug("role")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInventoryItemRoleResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_inventory_item_role.test", "id"),
					resource.TestCheckResourceAttr("netbox_inventory_item_role.test", "name", name),
				),
			},
			{
				Config: testAccInventoryItemRoleResourceConfig_full(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_inventory_item_role.test", "id"),
					resource.TestCheckResourceAttr("netbox_inventory_item_role.test", "color", "e41e22"),
					resource.TestCheckResourceAttr("netbox_inventory_item_role.test", "description", "Test role description"),
				),
			},
		},
	})
}

func TestAccInventoryItemRoleResource_import(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-inv-role")
	slug := testutil.RandomSlug("role")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInventoryItemRoleResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_inventory_item_role.test", "id"),
				),
			},
			{
				ResourceName:      "netbox_inventory_item_role.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:   testAccInventoryItemRoleResourceConfig_basic(name, slug),
				PlanOnly: true,
			},
		},
	})
}

func TestAccInventoryItemRoleResource_IDPreservation(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-inv-role-id")
	slug := testutil.RandomSlug("role")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInventoryItemRoleResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_inventory_item_role.test", "id"),
					resource.TestCheckResourceAttr("netbox_inventory_item_role.test", "name", name),
				),
			},
		},
	})
}

func testAccInventoryItemRoleResourceConfig_basic(name, slug string) string {

	return fmt.Sprintf(`
resource "netbox_inventory_item_role" "test" {
  name = %q
  slug = %q
}
`, name, slug)
}

func testAccInventoryItemRoleResourceConfig_full(name string) string {
	return fmt.Sprintf(`
resource "netbox_inventory_item_role" "test" {
  name        = %q
  slug        = %q
  color       = "e41e22"
  description = "Test role description"
}
`, name, testutil.RandomSlug("role"))
}

func TestAccConsistency_InventoryItemRole_LiteralNames(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-inv-role-lit")
	slug := testutil.RandomSlug("tf-test-inv-role-lit")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInventoryItemRoleConsistencyLiteralNamesConfig(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_inventory_item_role.test", "id"),
					resource.TestCheckResourceAttr("netbox_inventory_item_role.test", "name", name),
					resource.TestCheckResourceAttr("netbox_inventory_item_role.test", "slug", slug),
				),
			},
			{
				Config:   testAccInventoryItemRoleConsistencyLiteralNamesConfig(name, slug),
				PlanOnly: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_inventory_item_role.test", "id"),
				),
			},
		},
	})
}

func TestAccInventoryItemRoleResource_externalDeletion(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-inv-role-ext-del")
	slug := testutil.RandomSlug("role")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInventoryItemRoleResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_inventory_item_role.test", "id"),
					resource.TestCheckResourceAttr("netbox_inventory_item_role.test", "name", name),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					items, _, err := client.DcimAPI.DcimInventoryItemRolesList(context.Background()).NameIc([]string{name}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find inventory_item_role for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.DcimAPI.DcimInventoryItemRolesDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete inventory_item_role: %v", err)
					}
					t.Logf("Successfully externally deleted inventory_item_role with ID: %d", itemID)
				},
				Config: testAccInventoryItemRoleResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_inventory_item_role.test", "id"),
				),
			},
		},
	})
}

func testAccInventoryItemRoleConsistencyLiteralNamesConfig(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_inventory_item_role" "test" {
  name = %q
  slug = %q
}
`, name, slug)
}

func TestAccInventoryItemRoleResource_importWithCustomFieldsAndTags(t *testing.T) {
	// NOTE: t.Parallel() intentionally omitted - this test creates/deletes global custom fields
	// that would affect other tests of the same resource type running in parallel.

	roleName := testutil.RandomName("inventory_item_role")
	roleSlug := testutil.RandomSlug("inventory_item_role")

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
		Steps: []resource.TestStep{
			{
				Config: testAccInventoryItemRoleResourceImportConfig_full(roleName, roleSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_inventory_item_role.test", "id"),
					resource.TestCheckResourceAttr("netbox_inventory_item_role.test", "name", roleName),
					resource.TestCheckResourceAttr("netbox_inventory_item_role.test", "slug", roleSlug),
					// Verify custom fields are applied
					resource.TestCheckResourceAttr("netbox_inventory_item_role.test", "custom_fields.#", "7"),
					// Verify tags are applied
					resource.TestCheckResourceAttr("netbox_inventory_item_role.test", "tags.#", "2"),
				),
			},
			{
				Config:                  testAccInventoryItemRoleResourceImportConfig_full(roleName, roleSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug),
				ResourceName:            "netbox_inventory_item_role.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"custom_fields", "tags"},
			},
			{
				Config:   testAccInventoryItemRoleResourceImportConfig_full(roleName, roleSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug),
				PlanOnly: true,
			},
		},
	})
}

func testAccInventoryItemRoleResourceImportConfig_full(roleName, roleSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug string) string {
	return fmt.Sprintf(`
# Custom Fields
resource "netbox_custom_field" "cf_text" {
  name        = %q
  type        = "text"
  object_types = ["dcim.inventoryitemrole"]
}

resource "netbox_custom_field" "cf_longtext" {
  name        = %q
  type        = "longtext"
  object_types = ["dcim.inventoryitemrole"]
}

resource "netbox_custom_field" "cf_integer" {
  name        = %q
  type        = "integer"
  object_types = ["dcim.inventoryitemrole"]
}

resource "netbox_custom_field" "cf_boolean" {
  name        = %q
  type        = "boolean"
  object_types = ["dcim.inventoryitemrole"]
}

resource "netbox_custom_field" "cf_date" {
  name        = %q
  type        = "date"
  object_types = ["dcim.inventoryitemrole"]
}

resource "netbox_custom_field" "cf_url" {
  name        = %q
  type        = "url"
  object_types = ["dcim.inventoryitemrole"]
}

resource "netbox_custom_field" "cf_json" {
  name        = %q
  type        = "json"
  object_types = ["dcim.inventoryitemrole"]
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
resource "netbox_inventory_item_role" "test" {
  name = %q
  slug = %q

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
		cfText,
		cfLongtext,
		cfInteger,
		cfBoolean,
		cfDate,
		cfUrl,
		cfJson,
		tag1, tag1Slug,
		tag2, tag2Slug,
		roleName, roleSlug,
	)
}
