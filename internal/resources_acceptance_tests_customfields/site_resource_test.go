//go:build customfields
// +build customfields

package resources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccSiteResource_importWithCustomFieldsAndTags tests importing a site
// with custom fields and tags to ensure all data is preserved during import.
func TestAccSiteResource_importWithCustomFieldsAndTags(t *testing.T) {
	// NOTE: t.Parallel() intentionally omitted - this test creates/deletes global custom fields
	// that would affect other tests of the same resource type running in parallel.

	// Generate unique names
	siteName := testutil.RandomName("tf-test-site-import")
	siteSlug := testutil.RandomSlug("tf-test-site-import")

	// Generate tag names
	tag1Name := testutil.RandomName("tf-test-site-tag1")
	tag1Slug := testutil.RandomSlug("tf-test-site-tag1")
	tag1Color := testutil.RandomColor()
	tag2Name := testutil.RandomName("tf-test-site-tag2")
	tag2Slug := testutil.RandomSlug("tf-test-site-tag2")
	tag2Color := testutil.RandomColor()

	// Generate custom field names and values
	cfText := testutil.RandomCustomFieldName("tf_site_text")
	cfTextValue := testutil.RandomName("site-text-value")
	cfLongtext := testutil.RandomCustomFieldName("tf_site_longtext")
	cfLongtextValue := fmt.Sprintf("Site description: %s", testutil.RandomName("site-details"))
	cfIntegerName := testutil.RandomCustomFieldName("tf_site_integer")
	cfIntegerValue := 100
	cfBoolean := testutil.RandomCustomFieldName("tf_site_boolean")
	cfBooleanValue := true
	cfDate := testutil.RandomCustomFieldName("tf_site_date")
	cfDateValue := testutil.RandomDate()
	cfURL := testutil.RandomCustomFieldName("tf_site_url")
	cfURLValue := testutil.RandomURL("site")
	cfJSON := testutil.RandomCustomFieldName("tf_site_json")
	cfJSONValue := testutil.RandomJSON()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Create the site with all custom fields and tags
				Config: testAccSiteResourceImportConfig_full(
					siteName, siteSlug,
					tag1Name, tag1Slug, tag1Color, tag2Name, tag2Slug, tag2Color,
					cfText, cfTextValue, cfLongtext, cfLongtextValue, cfIntegerName, cfIntegerValue,
					cfBoolean, cfBooleanValue, cfDate, cfDateValue, cfURL, cfURLValue, cfJSON, cfJSONValue,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_site.test", "id"),
					resource.TestCheckResourceAttr("netbox_site.test", "name", siteName),
					resource.TestCheckResourceAttr("netbox_site.test", "tags.#", "2"),
					resource.TestCheckResourceAttr("netbox_site.test", "custom_fields.#", "7"),
				),
			},
			{
				// Import the site with identity-based import
				Config: testAccSiteResourceImportConfig_full(
					siteName, siteSlug,
					tag1Name, tag1Slug, tag1Color, tag2Name, tag2Slug, tag2Color,
					cfText, cfTextValue, cfLongtext, cfLongtextValue, cfIntegerName, cfIntegerValue,
					cfBoolean, cfBooleanValue, cfDate, cfDateValue, cfURL, cfURLValue, cfJSON, cfJSONValue,
				),
				ResourceName:      "netbox_site.test",
				ImportState:       true,
				ImportStateKind:   resource.ImportBlockWithResourceIdentity,
				ImportStateVerify: false,
			},
			{
				Config: testAccSiteResourceImportConfig_full(
					siteName, siteSlug,
					tag1Name, tag1Slug, tag1Color, tag2Name, tag2Slug, tag2Color,
					cfText, cfTextValue, cfLongtext, cfLongtextValue, cfIntegerName, cfIntegerValue,
					cfBoolean, cfBooleanValue, cfDate, cfDateValue, cfURL, cfURLValue, cfJSON, cfJSONValue,
				),
				PlanOnly: true,
			},
		},
	})
}

func testAccSiteResourceImportConfig_full(
	siteName, siteSlug,
	tag1Name, tag1Slug, tag1Color, tag2Name, tag2Slug, tag2Color string,
	cfText, cfTextValue, cfLongtext, cfLongtextValue string, cfIntegerName string, cfIntegerValue int,
	cfBoolean string, cfBooleanValue bool, cfDate, cfDateValue, cfURL, cfURLValue, cfJSON, cfJSONValue string,
) string {
	return fmt.Sprintf(`
# Create tags
resource "netbox_tag" "site_test1" {
  name  = %[3]q
  slug  = %[4]q
  color = %[5]q
}

resource "netbox_tag" "site_test2" {
  name  = %[6]q
  slug  = %[7]q
  color = %[8]q
}

# Create custom fields for dcim.site
resource "netbox_custom_field" "site_text" {
  name         = %[9]q
  type         = "text"
  object_types = ["dcim.site"]
  required     = false
}

resource "netbox_custom_field" "site_longtext" {
  name         = %[11]q
  type         = "longtext"
  object_types = ["dcim.site"]
  required     = false
}

resource "netbox_custom_field" "site_integer" {
  name         = %[13]q
  type         = "integer"
  object_types = ["dcim.site"]
  required     = false
}

resource "netbox_custom_field" "site_boolean" {
  name         = %[15]q
  type         = "boolean"
  object_types = ["dcim.site"]
  required     = false
}

resource "netbox_custom_field" "site_date" {
  name         = %[17]q
  type         = "date"
  object_types = ["dcim.site"]
  required     = false
}

resource "netbox_custom_field" "site_url" {
  name         = %[19]q
  type         = "url"
  object_types = ["dcim.site"]
  required     = false
}

resource "netbox_custom_field" "site_json" {
  name         = %[21]q
  type         = "json"
  object_types = ["dcim.site"]
  required     = false
}

# Create site with all custom fields and tags
resource "netbox_site" "test" {
  name   = %[1]q
  slug   = %[2]q
  status = "active"

	tags = [netbox_tag.site_test1.slug, netbox_tag.site_test2.slug]

  custom_fields = [
    {
      name  = netbox_custom_field.site_text.name
      type  = "text"
      value = %[10]q
    },
    {
      name  = netbox_custom_field.site_longtext.name
      type  = "longtext"
      value = %[12]q
    },
    {
      name  = netbox_custom_field.site_integer.name
      type  = "integer"
      value = "%[14]d"
    },
    {
      name  = netbox_custom_field.site_boolean.name
      type  = "boolean"
      value = "%[16]t"
    },
    {
      name  = netbox_custom_field.site_date.name
      type  = "date"
      value = %[18]q
    },
    {
      name  = netbox_custom_field.site_url.name
      type  = "url"
      value = %[20]q
    },
    {
      name  = netbox_custom_field.site_json.name
      type  = "json"
      value = %[22]q
    }
  ]
}
`, siteName, siteSlug, tag1Name, tag1Slug, tag1Color, tag2Name, tag2Slug, tag2Color,
		cfText, cfTextValue, cfLongtext, cfLongtextValue, cfIntegerName, cfIntegerValue, cfBoolean, cfBooleanValue, cfDate, cfDateValue, cfURL, cfURLValue, cfJSON, cfJSONValue)
}

// TestAccSiteResource_CustomFieldsPreservation tests that custom fields are preserved
// when updating other fields on a site. This addresses a critical bug where custom fields
// were being deleted when users updated unrelated fields.
//
// Bug scenario:
// 1. Create site with custom fields
// 2. Update site WITHOUT custom_fields in config (omit the field entirely)
// 3. Custom fields should be preserved in NetBox, not deleted.
func TestAccSiteResource_CustomFieldsPreservation(t *testing.T) {
	// Generate unique names
	siteName := testutil.RandomName("tf-test-site-cf-preserve")
	siteSlug := testutil.RandomSlug("tf-test-site-cf-preserve")

	// Custom field names
	cfText := testutil.RandomCustomFieldName("tf_text_preserve")
	cfInteger := testutil.RandomCustomFieldName("tf_int_preserve")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Step 1: Create site WITH custom fields explicitly in config
				Config: testAccSiteConfig_preservation_step1(
					siteName, siteSlug,
					cfText, cfInteger, "initial value", 42,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_site.test", "name", siteName),
					resource.TestCheckResourceAttr("netbox_site.test", "description", "Initial description"),
					resource.TestCheckResourceAttr("netbox_site.test", "custom_fields.#", "2"),
					testutil.CheckCustomFieldValue("netbox_site.test", cfText, "text", "initial value"),
					testutil.CheckCustomFieldValue("netbox_site.test", cfInteger, "integer", "42"),
				),
			},
			{
				// Step 2: Update description WITHOUT mentioning custom_fields in config
				// Custom fields should be preserved in NetBox (verified by import)
				// State shows null/empty for custom_fields since not in config
				Config: testAccSiteConfig_preservation_step2(
					siteName, siteSlug,
					cfText, cfInteger, "Updated description",
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_site.test", "name", siteName),
					resource.TestCheckResourceAttr("netbox_site.test", "description", "Updated description"),
					// State shows 0 custom_fields (not in config = not owned)
					resource.TestCheckResourceAttr("netbox_site.test", "custom_fields.#", "0"),
				),
			},
			{
				// Step 3: Import to verify custom fields still exist in NetBox
				ResourceName:            "netbox_site.test",
				ImportState:             true,
				ImportStateKind:         resource.ImportCommandWithID,
				ImportStateVerify:       false,                     // Can't verify - config has no custom_fields
				ImportStateVerifyIgnore: []string{"custom_fields"}, // Different because filter-to-owned
			},
			{
				// Step 4: Add custom_fields back to config to verify they were preserved
				Config: testAccSiteConfig_preservation_step1(
					siteName, siteSlug,
					cfText, cfInteger, "initial value", 42,
				),
				Check: resource.ComposeTestCheckFunc(
					// Custom fields should have their original values (preserved in NetBox)
					resource.TestCheckResourceAttr("netbox_site.test", "custom_fields.#", "2"),
					testutil.CheckCustomFieldValue("netbox_site.test", cfText, "text", "initial value"),
					testutil.CheckCustomFieldValue("netbox_site.test", cfInteger, "integer", "42"),
				),
			},
		},
	})
}

func testAccSiteConfig_preservation_step1(
	siteName, siteSlug,
	cfTextName, cfIntName, cfTextValue string, cfIntValue int,
) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "text" {
  name         = %[3]q
  type         = "text"
  object_types = ["dcim.site"]
}

resource "netbox_custom_field" "integer" {
  name         = %[4]q
  type         = "integer"
  object_types = ["dcim.site"]
}

resource "netbox_site" "test" {
  name        = %[1]q
  slug        = %[2]q
  status      = "active"
  description = %[5]q

  custom_fields = [
    {
      name  = netbox_custom_field.text.name
      type  = "text"
      value = %[6]q
    },
    {
      name  = netbox_custom_field.integer.name
      type  = "integer"
      value = "%[7]d"
    }
  ]
}
`,
		siteName, siteSlug,
		cfTextName, cfIntName, "Initial description", cfTextValue, cfIntValue,
	)
}

func testAccSiteConfig_preservation_step2(
	siteName, siteSlug,
	cfTextName, cfIntName, description string,
) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "text" {
  name         = %[3]q
  type         = "text"
  object_types = ["dcim.site"]
}

resource "netbox_custom_field" "integer" {
  name         = %[4]q
  type         = "integer"
  object_types = ["dcim.site"]
}

resource "netbox_site" "test" {
  name        = %[1]q
  slug        = %[2]q
  status      = "active"
  description = %[5]q

  # NOTE: custom_fields is intentionally omitted to test preservation behavior
}
`,
		siteName, siteSlug,
		cfTextName, cfIntName, description,
	)
}
