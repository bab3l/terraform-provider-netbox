//go:build customfields
// +build customfields

package datasources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccJournalEntryDataSource_customFields(t *testing.T) {
	customFieldName := testutil.RandomCustomFieldName("tf_test_journal_ds_cf")
	siteName := testutil.RandomName("tf-test-site-journal-ds-cf")
	siteSlug := testutil.RandomSlug("tf-test-site-journal-ds-cf")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccJournalEntryDataSourceConfig_customFields(customFieldName, siteName, siteSlug),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_journal_entry.test", "comments", "Test journal entry with custom field"),
					resource.TestCheckResourceAttr("data.netbox_journal_entry.test", "custom_fields.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_journal_entry.test", "custom_fields.0.name", customFieldName),
					resource.TestCheckResourceAttr("data.netbox_journal_entry.test", "custom_fields.0.type", "text"),
					resource.TestCheckResourceAttr("data.netbox_journal_entry.test", "custom_fields.0.value", "datasource-test-value"),
				),
			},
		},
	})
}

func testAccJournalEntryDataSourceConfig_customFields(customFieldName, siteName, siteSlug string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "test" {
  name         = %q
  object_types = ["extras.journalentry"]
  type         = "text"
}

resource "netbox_site" "test" {
  name   = %q
  slug   = %q
  status = "active"
}

resource "netbox_journal_entry" "test" {
  assigned_object_type = "dcim.site"
  assigned_object_id   = netbox_site.test.id
  comments             = "Test journal entry with custom field"
  kind                 = "info"

  custom_fields = [
    {
      name  = netbox_custom_field.test.name
      type  = "text"
      value = "datasource-test-value"
    }
  ]
}

data "netbox_journal_entry" "test" {
  id = netbox_journal_entry.test.id
}
`, customFieldName, siteName, siteSlug)
}
