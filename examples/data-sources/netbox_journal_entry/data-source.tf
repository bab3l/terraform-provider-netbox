# Example: Look up a journal entry by ID (only supported lookup method)
data "netbox_journal_entry" "by_id" {
  id = "123"
}

# Example: Use journal entry data in other resources
output "journal_id" {
  value = data.netbox_journal_entry.by_id.id
}

output "journal_comments" {
  value = data.netbox_journal_entry.by_id.comments
}

output "journal_kind" {
  value = data.netbox_journal_entry.by_id.kind
}

output "journal_object_type" {
  value = data.netbox_journal_entry.by_id.assigned_object_type
}

output "journal_created_by" {
  value = data.netbox_journal_entry.by_id.created_by
}

output "journal_created" {
  value = data.netbox_journal_entry.by_id.created
}

# Note: Journal entries do not support custom fields in NetBox API
output "journal_entry_note" {
  value       = "Journal entry data is read-only and does not support custom fields"
  description = "Journal entries track audit history and configuration changes"
}
