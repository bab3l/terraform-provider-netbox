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

# Use contact data in another resource
output "contact_phone" {
  value = data.netbox_contact.by_name.phone
}

output "contact_by_id" {
  value = data.netbox_contact.by_id
}

output "contact_by_email" {
  value = data.netbox_contact.by_email
}
