# Lookup by ID
data "netbox_region" "by_id" {
  id = "123"
}

# Lookup by slug
data "netbox_region" "by_slug" {
  slug = "north-america"
}

# Lookup by name
data "netbox_region" "by_name" {
  name = "North America"
}

# Use region data in other resources
output "region_name" {
  value = data.netbox_region.by_id.name
}

output "region_slug" {
  value = data.netbox_region.by_name.slug
}

output "region_parent" {
  value = data.netbox_region.by_id.parent
}

output "region_description" {
  value = data.netbox_region.by_id.description
}

# Access all custom fields
output "region_custom_fields" {
  value       = data.netbox_region.by_id.custom_fields
  description = "All custom fields defined in NetBox for this region"
}

# Access specific custom fields by name
output "region_timezone" {
  value       = try([for cf in data.netbox_region.by_id.custom_fields : cf.value if cf.name == "timezone"][0], null)
  description = "Example: accessing a text custom field"
}

output "region_country_code" {
  value       = try([for cf in data.netbox_region.by_id.custom_fields : cf.value if cf.name == "country_code"][0], null)
  description = "Example: accessing a text custom field"
}
