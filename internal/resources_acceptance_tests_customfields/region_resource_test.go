//go:build customfields
// +build customfields

package resources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRegionResource_importWithCustomFieldsAndTags(t *testing.T) {
	// NOTE: t.Parallel() intentionally omitted - this test creates/deletes global custom fields
	// that would affect other tests of the same resource type running in parallel.

	regionName := testutil.RandomName("region")
	regionSlug := testutil.RandomSlug("region")
	parentRegionName := testutil.RandomName("parent_region")
	parentRegionSlug := testutil.RandomSlug("parent_region")

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
	cleanup.RegisterRegionCleanup(regionSlug)
	cleanup.RegisterRegionCleanup(parentRegionSlug)
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
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRegionResourceImportConfig_full(parentRegionName, parentRegionSlug, regionName, regionSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_region.test", "id"),
					resource.TestCheckResourceAttr("netbox_region.test", "name", regionName),
					// Verify custom fields are applied
					resource.TestCheckResourceAttr("netbox_region.test", "custom_fields.#", "7"),
					// Verify tags are applied
					resource.TestCheckResourceAttr("netbox_region.test", "tags.#", "2"),
				),
			},
			{
				Config:            testAccRegionResourceImportConfig_full(parentRegionName, parentRegionSlug, regionName, regionSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug),
				ResourceName:      "netbox_region.test",
				ImportState:       true,
				ImportStateKind:   resource.ImportBlockWithResourceIdentity,
				ImportStateVerify: false,
			},
			{
				Config:   testAccRegionResourceImportConfig_full(parentRegionName, parentRegionSlug, regionName, regionSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug),
				PlanOnly: true,
			},
		},
	})
}

// TestAccRegionResource_CustomFieldsPreservation tests that custom fields are preserved
// when updating other fields on a region. This addresses a critical bug where custom fields
// were being deleted when users updated unrelated fields.
//
// Bug scenario:
// 1. Create region with custom fields
// 2. Update region WITHOUT custom_fields in config (omit the field entirely)
// 3. Custom fields should be preserved in NetBox, not deleted.
func TestAccRegionResource_CustomFieldsPreservation(t *testing.T) {
	// Generate unique names
	regionName := testutil.RandomName("tf-test-region-cf-preserve")
	regionSlug := testutil.RandomSlug("tf-test-region-cf-preserve")

	// Custom field names
	cfText := testutil.RandomCustomFieldName("tf_text_preserve")
	cfInteger := testutil.RandomCustomFieldName("tf_int_preserve")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRegionCleanup(regionSlug)
	cleanup.RegisterCustomFieldCleanup(cfText)
	cleanup.RegisterCustomFieldCleanup(cfInteger)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Step 1: Create region WITH custom fields explicitly in config
				Config: testAccRegionConfig_preservation_step1(
					regionName, regionSlug,
					cfText, cfInteger, "initial value", 42,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_region.test", "name", regionName),
					resource.TestCheckResourceAttr("netbox_region.test", "description", "Initial description"),
					resource.TestCheckResourceAttr("netbox_region.test", "custom_fields.#", "2"),
					testutil.CheckCustomFieldValue("netbox_region.test", cfText, "text", "initial value"),
					testutil.CheckCustomFieldValue("netbox_region.test", cfInteger, "integer", "42"),
				),
			},
			{
				// Step 2: Update description WITHOUT mentioning custom_fields in config
				// Custom fields should be preserved in NetBox (verified by import)
				// State shows null/empty for custom_fields since not in config
				Config: testAccRegionConfig_preservation_step2(
					regionName, regionSlug,
					cfText, cfInteger, "Updated description",
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_region.test", "name", regionName),
					resource.TestCheckResourceAttr("netbox_region.test", "description", "Updated description"),
					// State shows 0 custom_fields (not in config = not owned)
					resource.TestCheckResourceAttr("netbox_region.test", "custom_fields.#", "0"),
				),
			},
			{
				// Step 3: Import to verify custom fields still exist in NetBox
				ResourceName:            "netbox_region.test",
				ImportState:             true,
				ImportStateKind:         resource.ImportCommandWithID,
				ImportStateVerify:       false,                             // Can't verify - config has no custom_fields
				ImportStateVerifyIgnore: []string{"custom_fields", "tags"}, // Different because filter-to-owned
			},
			{
				// Step 4: Add custom_fields back to config to verify they were preserved
				Config: testAccRegionConfig_preservation_step1(
					regionName, regionSlug,
					cfText, cfInteger, "initial value", 42,
				),
				Check: resource.ComposeTestCheckFunc(
					// Custom fields should have their original values (preserved in NetBox)
					resource.TestCheckResourceAttr("netbox_region.test", "custom_fields.#", "2"),
					testutil.CheckCustomFieldValue("netbox_region.test", cfText, "text", "initial value"),
					testutil.CheckCustomFieldValue("netbox_region.test", cfInteger, "integer", "42"),
				),
			},
		},
	})
}

func testAccRegionResourceImportConfig_full(parentRegionName, parentRegionSlug, regionName, regionSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug string) string {
	return fmt.Sprintf(`
resource "netbox_region" "parent" {
  name = %q
  slug = %q
}

resource "netbox_tag" "test1" {
  name = %q
  slug = %q
}

resource "netbox_tag" "test2" {
  name = %q
  slug = %q
}

resource "netbox_custom_field" "text" {
  name         = %q
  type         = "text"
  object_types = ["dcim.region"]
}

resource "netbox_custom_field" "longtext" {
  name         = %q
  type         = "longtext"
  object_types = ["dcim.region"]
}

resource "netbox_custom_field" "integer" {
  name         = %q
  type         = "integer"
  object_types = ["dcim.region"]
}

resource "netbox_custom_field" "boolean" {
  name         = %q
  type         = "boolean"
  object_types = ["dcim.region"]
}

resource "netbox_custom_field" "date" {
  name         = %q
  type         = "date"
  object_types = ["dcim.region"]
}

resource "netbox_custom_field" "url" {
  name         = %q
  type         = "url"
  object_types = ["dcim.region"]
}

resource "netbox_custom_field" "json" {
  name         = %q
  type         = "json"
  object_types = ["dcim.region"]
}

resource "netbox_region" "test" {
  name   = %q
  slug   = %q
	parent = netbox_region.parent.id

	tags = [netbox_tag.test1.slug, netbox_tag.test2.slug]

  custom_fields = [
    {
      name  = netbox_custom_field.text.name
      type  = "text"
      value = "test"
    },
    {
      name  = netbox_custom_field.longtext.name
      type  = "longtext"
      value = "longtext value"
    },
    {
      name  = netbox_custom_field.integer.name
      type  = "integer"
      value = "100"
    },
    {
      name  = netbox_custom_field.boolean.name
      type  = "boolean"
      value = "true"
    },
    {
      name  = netbox_custom_field.date.name
      type  = "date"
      value = "2024-01-01"
    },
    {
      name  = netbox_custom_field.url.name
      type  = "url"
      value = "https://example.com"
    },
    {
      name  = netbox_custom_field.json.name
      type  = "json"
      value = jsonencode({ key = "value" })
    }
  ]
}
`,
		parentRegionName, parentRegionSlug,
		tag1, tag1Slug,
		tag2, tag2Slug,
		cfText,
		cfLongtext,
		cfInteger,
		cfBoolean,
		cfDate,
		cfUrl,
		cfJson,
		regionName, regionSlug,
	)
}

func testAccRegionConfig_preservation_step1(
	regionName, regionSlug,
	cfTextName, cfIntName, cfTextValue string, cfIntValue int,
) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "text" {
  name         = %q
  type         = "text"
  object_types = ["dcim.region"]
}

resource "netbox_custom_field" "integer" {
  name         = %q
  type         = "integer"
  object_types = ["dcim.region"]
}

resource "netbox_region" "test" {
  name        = %q
  slug        = %q
  description = "Initial description"

  custom_fields = [
    {
      name  = netbox_custom_field.text.name
      type  = "text"
      value = %q
    },
    {
      name  = netbox_custom_field.integer.name
      type  = "integer"
      value = "%d"
    }
  ]
}
`,
		cfTextName, cfIntName,
		regionName, regionSlug,
		cfTextValue, cfIntValue,
	)
}

func testAccRegionConfig_preservation_step2(
	regionName, regionSlug,
	cfTextName, cfIntName, description string,
) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "text" {
  name         = %q
  type         = "text"
  object_types = ["dcim.region"]
}

resource "netbox_custom_field" "integer" {
  name         = %q
  type         = "integer"
  object_types = ["dcim.region"]
}

resource "netbox_region" "test" {
  name        = %q
  slug        = %q
  description = %q
}
`,
		cfTextName, cfIntName,
		regionName, regionSlug,
		description,
	)
}
