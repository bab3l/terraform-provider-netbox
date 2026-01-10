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

output "contact_group_slug" {
  value = data.netbox_contact_group.by_slug.slug
}

output "contact_group_parent" {
  value = data.netbox_contact_group.by_id.parent
}

output "contact_group_description" {
  value = data.netbox_contact_group.by_id.description
}

# Access all custom fields
output "contact_group_custom_fields" {
  value       = data.netbox_contact_group.by_id.custom_fields
  description = "All custom fields defined in NetBox for this contact group"
}

# Access specific custom fields by name
output "contact_group_manager" {
  value       = try([for cf in data.netbox_contact_group.by_id.custom_fields : cf.value if cf.name == "manager_name"][0], null)
  description = "Example: accessing a text custom field"
}

output "contact_group_budget_code" {
  value       = try([for cf in data.netbox_contact_group.by_id.custom_fields : cf.value if cf.name == "budget_code"][0], null)
  description = "Example: accessing a text custom field"
}
