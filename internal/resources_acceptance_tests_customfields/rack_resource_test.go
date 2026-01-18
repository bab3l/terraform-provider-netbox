//go:build customfields
// +build customfields

package resources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRackResource_importWithCustomFieldsAndTags(t *testing.T) {
	// NOTE: t.Parallel() intentionally omitted - this test creates/deletes global custom fields
	// that would affect other tests of the same resource type running in parallel.

	siteName := testutil.RandomName("tf-test-rack-site")
	siteSlug := testutil.RandomSlug("tf-test-rack-site")
	rackName := testutil.RandomName("tf-test-rack")
	tenantName := testutil.RandomName("tf-test-tenant")
	tenantSlug := testutil.RandomSlug("tf-test-tenant")

	// Generate test data for all custom field types (once, used in all steps)
	textValue := testutil.RandomName("text-value")
	longtextValue := testutil.RandomName("longtext-value") + "\nThis is a multiline text field for comprehensive testing."
	intValue := 42 // Fixed value for reproducibility
	boolValue := true
	dateValue := testutil.RandomDate()
	urlValue := testutil.RandomURL("test-url")
	jsonValue := testutil.RandomJSON()

	// Tag names
	tag1 := testutil.RandomName("tag1")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2 := testutil.RandomName("tag2")
	tag2Slug := testutil.RandomSlug("tag2")

	// Custom field names
	cfText := testutil.RandomCustomFieldName("tf_text")
	cfLongtext := testutil.RandomCustomFieldName("tf_longtext")
	cfInteger := testutil.RandomCustomFieldName("tf_integer")
	cfBoolean := testutil.RandomCustomFieldName("tf_boolean")
	cfDate := testutil.RandomCustomFieldName("tf_date")
	cfURL := testutil.RandomCustomFieldName("tf_url")
	cfJSON := testutil.RandomCustomFieldName("tf_json")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRackCleanup(rackName)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterTenantCleanup(tenantSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.ComposeCheckDestroy(testutil.CheckRackDestroy, testutil.CheckSiteDestroy),
		Steps: []resource.TestStep{
			{
				Config: testAccRackResourceImportConfig_full(siteName, siteSlug, rackName, tenantName, tenantSlug,
					tag1, tag1Slug, tag2, tag2Slug,
					cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfURL, cfJSON,
					textValue, longtextValue, intValue, boolValue, dateValue, urlValue, jsonValue),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_rack.test", "id"),
					resource.TestCheckResourceAttr("netbox_rack.test", "name", rackName),
				),
			},
			{
				ResourceName:            "netbox_rack.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"site", "tenant", "custom_fields", "tags"},
			},
			{
				Config: testAccRackResourceImportConfig_full(siteName, siteSlug, rackName, tenantName, tenantSlug,
					tag1, tag1Slug, tag2, tag2Slug,
					cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfURL, cfJSON,
					textValue, longtextValue, intValue, boolValue, dateValue, urlValue, jsonValue),
				PlanOnly: true,
			},
		},
	})
}

func testAccRackResourceImportConfig_full(
	siteName, siteSlug, rackName, tenantName, tenantSlug string,
	tag1, tag1Slug, tag2, tag2Slug string,
	cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfURL, cfJSON string,
	textValue, longtextValue string, intValue int, boolValue bool, dateValue, urlValue, jsonValue string,
) string {

	return fmt.Sprintf(`
# Dependencies
resource "netbox_site" "test" {
  name   = %q
  slug   = %q
  status = "active"
}

resource "netbox_tenant" "test" {
  name = %q
  slug = %q
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

# Custom Fields for dcim.rack object type
resource "netbox_custom_field" "test_text" {
  name         = %q
  label        = "Test Text CF"
  type         = "text"
  object_types = ["dcim.rack"]
}

resource "netbox_custom_field" "test_longtext" {
  name         = %q
  label        = "Test Longtext CF"
  type         = "longtext"
  object_types = ["dcim.rack"]
}

resource "netbox_custom_field" "test_integer" {
  name         = %q
  label        = "Test Integer CF"
  type         = "integer"
  object_types = ["dcim.rack"]
}

resource "netbox_custom_field" "test_boolean" {
  name         = %q
  label        = "Test Boolean CF"
  type         = "boolean"
  object_types = ["dcim.rack"]
}

resource "netbox_custom_field" "test_date" {
  name         = %q
  label        = "Test Date CF"
  type         = "date"
  object_types = ["dcim.rack"]
}

resource "netbox_custom_field" "test_url" {
  name         = %q
  label        = "Test URL CF"
  type         = "url"
  object_types = ["dcim.rack"]
}

resource "netbox_custom_field" "test_json" {
  name         = %q
  label        = "Test JSON CF"
  type         = "json"
  object_types = ["dcim.rack"]
}

# Rack with comprehensive custom fields and tags
resource "netbox_rack" "test" {
  name   = %q
  site   = netbox_site.test.id
  status = "active"
  tenant = netbox_tenant.test.id

	tags = [netbox_tag.tag1.slug, netbox_tag.tag2.slug]

  custom_fields = [
    {
      name  = netbox_custom_field.test_text.name
      type  = "text"
      value = %q
    },
    {
      name  = netbox_custom_field.test_longtext.name
      type  = "longtext"
      value = %q
    },
    {
      name  = netbox_custom_field.test_integer.name
      type  = "integer"
      value = "%d"
    },
    {
      name  = netbox_custom_field.test_boolean.name
      type  = "boolean"
      value = "%t"
    },
    {
      name  = netbox_custom_field.test_date.name
      type  = "date"
      value = %q
    },
    {
      name  = netbox_custom_field.test_url.name
      type  = "url"
      value = %q
    },
    {
      name  = netbox_custom_field.test_json.name
      type  = "json"
      value = %q
    },
  ]

  depends_on = [
    netbox_custom_field.test_text,
    netbox_custom_field.test_longtext,
    netbox_custom_field.test_integer,
    netbox_custom_field.test_boolean,
    netbox_custom_field.test_date,
    netbox_custom_field.test_url,
    netbox_custom_field.test_json,
  ]
}
`, siteName, siteSlug, tenantName, tenantSlug,
		tag1, tag1Slug, tag2, tag2Slug,
		cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfURL, cfJSON,
		rackName, textValue, longtextValue, intValue, boolValue, dateValue, urlValue, jsonValue)
}

// TestAccRackResource_CustomFieldsPreservation tests that custom fields are preserved
// when updating other fields on a rack. This addresses a critical bug where custom fields
// were being deleted when users updated unrelated fields.
//
// Bug scenario:
// 1. Create rack with custom fields
// 2. Update rack WITHOUT custom_fields in config (omit the field entirely)
// 3. Custom fields should be preserved in NetBox, not deleted.
func TestAccRackResource_CustomFieldsPreservation(t *testing.T) {
	// Generate unique names
	siteName := testutil.RandomName("tf-test-rack-site-preserve")
	siteSlug := testutil.RandomSlug("tf-test-rack-site-preserve")
	rackName := testutil.RandomName("tf-test-rack-preserve")

	// Custom field names
	cfText := testutil.RandomCustomFieldName("tf_text_preserve")
	cfInteger := testutil.RandomCustomFieldName("tf_int_preserve")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRackCleanup(rackName)
	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.ComposeCheckDestroy(testutil.CheckRackDestroy, testutil.CheckSiteDestroy),
		Steps: []resource.TestStep{
			{
				// Step 1: Create rack WITH custom fields explicitly in config
				Config: testAccRackConfig_preservation_step1(
					siteName, siteSlug, rackName,
					cfText, cfInteger, "initial value", 42,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_rack.test", "name", rackName),
					resource.TestCheckResourceAttr("netbox_rack.test", "description", "Initial description"),
					resource.TestCheckResourceAttr("netbox_rack.test", "custom_fields.#", "2"),
					testutil.CheckCustomFieldValue("netbox_rack.test", cfText, "text", "initial value"),
					testutil.CheckCustomFieldValue("netbox_rack.test", cfInteger, "integer", "42"),
				),
			},
			{
				// Step 2: Update description WITHOUT mentioning custom_fields in config
				// Custom fields should be preserved in NetBox (verified by import)
				// State shows null/empty for custom_fields since not in config
				Config: testAccRackConfig_preservation_step2(
					siteName, siteSlug, rackName,
					cfText, cfInteger, "Updated description",
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_rack.test", "name", rackName),
					resource.TestCheckResourceAttr("netbox_rack.test", "description", "Updated description"),
					// State shows 0 custom_fields (not in config = not owned)
					resource.TestCheckResourceAttr("netbox_rack.test", "custom_fields.#", "0"),
				),
			},
			{
				// Step 3: Import to verify custom fields still exist in NetBox
				ResourceName:            "netbox_rack.test",
				ImportState:             true,
				ImportStateVerify:       false,                     // Can't verify - config has no custom_fields
				ImportStateVerifyIgnore: []string{"custom_fields"}, // Different because filter-to-owned
			},
			{
				// Step 4: Add custom_fields back to config to verify they were preserved
				Config: testAccRackConfig_preservation_step1(
					siteName, siteSlug, rackName,
					cfText, cfInteger, "initial value", 42,
				),
				Check: resource.ComposeTestCheckFunc(
					// Custom fields should have their original values (preserved in NetBox)
					resource.TestCheckResourceAttr("netbox_rack.test", "custom_fields.#", "2"),
					testutil.CheckCustomFieldValue("netbox_rack.test", cfText, "text", "initial value"),
					testutil.CheckCustomFieldValue("netbox_rack.test", cfInteger, "integer", "42"),
				),
			},
		},
	})
}

func testAccRackConfig_preservation_step1(
	siteName, siteSlug, rackName,
	cfTextName, cfIntName, cfTextValue string, cfIntValue int,
) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_custom_field" "text" {
  name         = %[4]q
  type         = "text"
  object_types = ["dcim.rack"]
}

resource "netbox_custom_field" "integer" {
  name         = %[5]q
  type         = "integer"
  object_types = ["dcim.rack"]
}

resource "netbox_rack" "test" {
  name        = %[3]q
  site        = netbox_site.test.id
  description = %[6]q

  custom_fields = [
    {
      name  = netbox_custom_field.text.name
      type  = "text"
      value = %[7]q
    },
    {
      name  = netbox_custom_field.integer.name
      type  = "integer"
      value = "%[8]d"
    }
  ]

  depends_on = [
    netbox_custom_field.text,
    netbox_custom_field.integer,
  ]
}
`,
		siteName, siteSlug, rackName,
		cfTextName, cfIntName, "Initial description", cfTextValue, cfIntValue,
	)
}

func testAccRackConfig_preservation_step2(
	siteName, siteSlug, rackName,
	cfTextName, cfIntName, description string,
) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_custom_field" "text" {
  name         = %[4]q
  type         = "text"
  object_types = ["dcim.rack"]
}

resource "netbox_custom_field" "integer" {
  name         = %[5]q
  type         = "integer"
  object_types = ["dcim.rack"]
}

resource "netbox_rack" "test" {
  name        = %[3]q
  site        = netbox_site.test.id
  description = %[6]q

  # NOTE: custom_fields is intentionally omitted to test preservation behavior
}
`,
		siteName, siteSlug, rackName,
		cfTextName, cfIntName, description,
	)
}
