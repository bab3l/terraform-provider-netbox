# Journal Entry Data Source Integration Test
# Tests the netbox_journal_entry data source for looking up existing journal entries

terraform {
  required_version = ">= 1.0"
  required_providers {
    netbox = {
      source = "bab3l/netbox"
    }
  }
}

provider "netbox" {
  # Uses NETBOX_SERVER_URL and NETBOX_API_TOKEN environment variables
}

# Create a site to attach journal entries to
resource "netbox_site" "test" {
  name        = "Journal DS Test Site"
  slug        = "journal-ds-test-site"
  description = "Site for journal entry data source testing"
}

# Create journal entries to look up
resource "netbox_journal_entry" "test_info" {
  assigned_object_type = "dcim.site"
  assigned_object_id   = netbox_site.test.id
  comments             = "Info journal entry for data source testing"
  kind                 = "info"
}

resource "netbox_journal_entry" "test_warning" {
  assigned_object_type = "dcim.site"
  assigned_object_id   = netbox_site.test.id
  comments             = "Warning journal entry for data source testing"
  kind                 = "warning"
}

# Look up journal entry by ID
data "netbox_journal_entry" "by_id_info" {
  id = netbox_journal_entry.test_info.id
}

data "netbox_journal_entry" "by_id_warning" {
  id = netbox_journal_entry.test_warning.id
}

# Outputs for verification
output "info_entry_comments" {
  value = data.netbox_journal_entry.by_id_info.comments
}

output "info_entry_kind" {
  value = data.netbox_journal_entry.by_id_info.kind
}

output "warning_entry_comments" {
  value = data.netbox_journal_entry.by_id_warning.comments
}

output "warning_entry_kind" {
  value = data.netbox_journal_entry.by_id_warning.kind
}

output "info_assigned_object_type" {
  value = data.netbox_journal_entry.by_id_info.assigned_object_type
}
