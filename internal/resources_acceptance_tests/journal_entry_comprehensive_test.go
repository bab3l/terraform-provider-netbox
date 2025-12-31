package resources_acceptance_tests

import (
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
)

// TestAccJournalEntryResource_Kind tests comprehensive scenarios for journal entry kind field.
// This validates that Optional+Computed string fields with proper defaults work correctly.
func TestAccJournalEntryResource_Kind(t *testing.T) {

	testutil.RunOptionalComputedFieldTestSuite(t, testutil.OptionalComputedFieldTestConfig{
		ResourceName:   "netbox_journal_entry",
		OptionalField:  "kind",
		DefaultValue:   "info",
		FieldTestValue: "warning",
		BaseConfig: func() string {
			return `
resource "netbox_site" "test" {
	name = "test-site-journal-entry"
	slug = "test-site-journal-entry"
}

resource "netbox_journal_entry" "test" {
	assigned_object_type = "dcim.site"
	assigned_object_id   = netbox_site.test.id
	comments             = "Test journal entry for kind field validation"
	# kind field intentionally omitted - should get default "info"
}
`
		},
		WithFieldConfig: func(value string) string {
			return `
resource "netbox_site" "test" {
	name = "test-site-journal-entry"
	slug = "test-site-journal-entry"
}

resource "netbox_journal_entry" "test" {
	assigned_object_type = "dcim.site"
	assigned_object_id   = netbox_site.test.id
	comments             = "Test journal entry for kind field validation"
	kind                 = "` + value + `"
}
`
		},
	})
}
