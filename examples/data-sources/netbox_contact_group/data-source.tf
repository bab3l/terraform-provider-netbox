# Example: Look up contact group by ID
data "netbox_contact_group" "by_id" {
  id = "1"
}

# Example: Look up contact group by name
data "netbox_contact_group" "by_name" {
  name = "IT Department"
}

# Example: Look up contact group by slug
data "netbox_contact_group" "by_slug" {
  slug = "it-department"
}

# Example: Use contact group data in other resources
output "contact_group_name" {
  value = data.netbox_contact_group.by_name.name
}

output "contact_group_parent_id" {
  value = data.netbox_contact_group.by_name.parent_id
}
