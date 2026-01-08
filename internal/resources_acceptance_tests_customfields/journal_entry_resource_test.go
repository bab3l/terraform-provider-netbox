//go:build customfields
// +build customfields

package resources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccJournalEntryResource_CustomFieldsPreservation tests that custom fields are preserved
// when updating other fields on a journal entry.
func TestAccJournalEntryResource_CustomFieldsPreservation(t *testing.T) {
	// Generate unique names
	siteName := testutil.RandomName("tf-test-site")
	siteSlug := testutil.RandomSlug("tf-test-site")

	// Custom field names
	cfText := testutil.RandomCustomFieldName("tf_text_preserve")
	cfInteger := testutil.RandomCustomFieldName("tf_int_preserve")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.ComposeCheckDestroy(testutil.CheckJournalEntryDestroy),
		Steps: []resource.TestStep{
			{
				// Step 1: Create journal entry WITH custom fields explicitly in config
				Config: testAccJournalEntryConfig_preservation_step1(
					siteName, siteSlug,
					cfText, cfInteger, "initial value", 42,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_journal_entry.test", "comments", "Initial comment"),
					resource.TestCheckResourceAttr("netbox_journal_entry.test", "custom_fields.#", "2"),
					testutil.CheckCustomFieldValue("netbox_journal_entry.test", cfText, "text", "initial value"),
					testutil.CheckCustomFieldValue("netbox_journal_entry.test", cfInteger, "integer", "42"),
				),
			},
			{
				// Step 2: Update comments WITHOUT mentioning custom_fields in config
				Config: testAccJournalEntryConfig_preservation_step2(
					siteName, siteSlug,
					cfText, cfInteger,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_journal_entry.test", "comments", "Updated comment"),
					resource.TestCheckResourceAttr("netbox_journal_entry.test", "custom_fields.#", "0"),
				),
			},
			{
				// Step 3: Import to verify custom fields still exist in NetBox
				ResourceName:            "netbox_journal_entry.test",
				ImportState:             true,
				ImportStateVerify:       false,
				ImportStateVerifyIgnore: []string{"custom_fields"},
			},
			{
				// Step 4: Add custom_fields back to config to verify they were preserved
				Config: testAccJournalEntryConfig_preservation_step1(
					siteName, siteSlug,
					cfText, cfInteger, "initial value", 42,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_journal_entry.test", "custom_fields.#", "2"),
					testutil.CheckCustomFieldValue("netbox_journal_entry.test", cfText, "text", "initial value"),
					testutil.CheckCustomFieldValue("netbox_journal_entry.test", cfInteger, "integer", "42"),
				),
			},
		},
	})
}

func testAccJournalEntryConfig_preservation_step1(
	siteName, siteSlug,
	cfTextName, cfIntName, cfTextValue string, cfIntValue int,
) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "text" {
  name         = %[3]q
  object_types = ["extras.journalentry"]
  type         = "text"
}

resource "netbox_custom_field" "integer" {
  name         = %[4]q
  object_types = ["extras.journalentry"]
  type         = "integer"
}

resource "netbox_site" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_journal_entry" "test" {
  assigned_object_type = "dcim.site"
  assigned_object_id   = netbox_site.test.id
  comments             = "Initial comment"

  custom_fields = [
    {
      name  = netbox_custom_field.text.name
      type  = "text"
      value = %[5]q
    },
    {
      name  = netbox_custom_field.integer.name
      type  = "integer"
      value = %[6]d
    }
  ]
}
`, siteName, siteSlug, cfTextName, cfIntName, cfTextValue, cfIntValue)
}

func testAccJournalEntryConfig_preservation_step2(
	siteName, siteSlug,
	cfTextName, cfIntName string,
) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "text" {
  name         = %[3]q
  object_types = ["extras.journalentry"]
  type         = "text"
}

resource "netbox_custom_field" "integer" {
  name         = %[4]q
  object_types = ["extras.journalentry"]
  type         = "integer"
}

resource "netbox_site" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_journal_entry" "test" {
  assigned_object_type = "dcim.site"
  assigned_object_id   = netbox_site.test.id
  comments             = "Updated comment"

  # custom_fields intentionally omitted
}
`, siteName, siteSlug, cfTextName, cfIntName)
}
