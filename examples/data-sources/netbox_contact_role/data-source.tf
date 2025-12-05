# Example: Look up contact role by ID
data "netbox_contact_role" "by_id" {
  id = "1"
}

# Example: Look up contact role by name
data "netbox_contact_role" "by_name" {
  name = "Technical"
}

# Example: Look up contact role by slug
data "netbox_contact_role" "by_slug" {
  slug = "technical"
}

# Example: Use contact role data in other resources
output "contact_role_name" {
  value = data.netbox_contact_role.by_name.name
}

output "contact_role_description" {
  value = data.netbox_contact_role.by_name.description
}
