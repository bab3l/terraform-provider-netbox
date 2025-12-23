# Example: Look up a journal entry by ID (only supported lookup method)
data "netbox_journal_entry" "by_id" {
  id = 123
}

# Example: Use journal entry data in other resources
output "journal_comments" {
  value = data.netbox_journal_entry.by_id.comments
}

output "journal_kind" {
  value = data.netbox_journal_entry.by_id.kind
}

output "journal_object_type" {
  value = data.netbox_journal_entry.by_id.assigned_object_type
}
