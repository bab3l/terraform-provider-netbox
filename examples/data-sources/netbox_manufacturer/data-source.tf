# Look up manufacturer by ID
data "netbox_manufacturer" "by_id" {
  id = "1"
}

# Look up manufacturer by name
data "netbox_manufacturer" "by_name" {
  name = "Cisco"
}

# Look up manufacturer by slug
data "netbox_manufacturer" "by_slug" {
  slug = "cisco"
}

# Use manufacturer data in other resources
output "manufacturer_id" {
  value = data.netbox_manufacturer.by_id.id
}

output "manufacturer_name" {
  value = data.netbox_manufacturer.by_name.name
}

output "manufacturer_slug" {
  value = data.netbox_manufacturer.by_slug.slug
}

output "manufacturer_description" {
  value = data.netbox_manufacturer.by_id.description
}

# Access all custom fields
output "manufacturer_custom_fields" {
  value       = data.netbox_manufacturer.by_id.custom_fields
  description = "All custom fields defined in NetBox for this manufacturer"
}

# Access specific custom fields by name
output "manufacturer_support_url" {
  value       = try([for cf in data.netbox_manufacturer.by_id.custom_fields : cf.value if cf.name == "support_url"][0], null)
  description = "Example: accessing a URL custom field"
}

output "manufacturer_preferred_vendor" {
  value       = try([for cf in data.netbox_manufacturer.by_id.custom_fields : cf.value if cf.name == "preferred_vendor"][0], null)
  description = "Example: accessing a boolean custom field"
}
