# Look up site group by ID
data "netbox_site_group" "by_id" {
  id = "1"
}

# Look up site group by slug
data "netbox_site_group" "by_slug" {
  slug = "north-america"
}

# Look up site group by name
data "netbox_site_group" "by_name" {
  name = "North America"
}

# Use site group data in other resources
output "site_group_name" {
  value = data.netbox_site_group.by_name.name
}

output "site_group_slug" {
  value = data.netbox_site_group.by_slug.slug
}

output "site_group_parent" {
  value = data.netbox_site_group.by_id.parent
}

output "site_group_description" {
  value = data.netbox_site_group.by_id.description
}

# Access all custom fields
output "site_group_custom_fields" {
  value       = data.netbox_site_group.by_id.custom_fields
  description = "All custom fields defined in NetBox for this site group"
}

# Access specific custom fields by name
output "site_group_regional_manager" {
  value       = try([for cf in data.netbox_site_group.by_id.custom_fields : cf.value if cf.name == "regional_manager"][0], null)
  description = "Example: accessing a text custom field"
}

output "site_group_cost_center" {
  value       = try([for cf in data.netbox_site_group.by_id.custom_fields : cf.value if cf.name == "cost_center"][0], null)
  description = "Example: accessing a text custom field"
}
