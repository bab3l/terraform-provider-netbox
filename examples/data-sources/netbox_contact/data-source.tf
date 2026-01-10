# Look up a contact by ID
data "netbox_contact" "by_id" {
  id = "123"
}

# Look up a contact by name
data "netbox_contact" "by_name" {
  name = "John Doe"
}

# Look up a contact by email
data "netbox_contact" "by_email" {
  email = "john.doe@example.com"
}

# Use contact data in other resources
output "contact_name" {
  value = data.netbox_contact.by_id.name
}

output "contact_phone" {
  value = data.netbox_contact.by_name.phone
}

output "contact_email" {
  value = data.netbox_contact.by_email.email
}

output "contact_group" {
  value = data.netbox_contact.by_id.group
}

output "contact_title" {
  value = data.netbox_contact.by_id.title
}

# Access all custom fields
output "contact_custom_fields" {
  value       = data.netbox_contact.by_id.custom_fields
  description = "All custom fields defined in NetBox for this contact"
}

# Access specific custom fields by name
output "contact_department" {
  value       = try([for cf in data.netbox_contact.by_id.custom_fields : cf.value if cf.name == "department"][0], null)
  description = "Example: accessing a text custom field"
}

output "contact_on_call" {
  value       = try([for cf in data.netbox_contact.by_id.custom_fields : cf.value if cf.name == "on_call"][0], null)
  description = "Example: accessing a boolean custom field"
}
