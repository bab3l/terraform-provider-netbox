# Look up a journal entry by ID
data "netbox_journal_entry" "example" {
  id = 123
}

# Use the journal entry data
output "journal_comments" {
  value = data.netbox_journal_entry.example.comments
}

output "journal_kind" {
  value = data.netbox_journal_entry.example.kind
}
