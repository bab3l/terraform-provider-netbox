# Example: Look up a location by ID
data "netbox_location" "by_id" {
  id = "1"
}

# Example: Look up a location by name
data "netbox_location" "by_name" {
  name = "Building-A-Floor-1"
}

# Example: Look up a location by slug
data "netbox_location" "by_slug" {
  slug = "building-a-floor-1"
}

# Example: Use location data in other resources
output "location_id" {
  value = data.netbox_location.by_id.id
}

output "location_name" {
  value = data.netbox_location.by_name.name
}

output "location_site" {
  value = data.netbox_location.by_slug.site
}

output "location_parent" {
  value = data.netbox_location.by_id.parent
}

output "location_status" {
  value = data.netbox_location.by_id.status
}

# Access all custom fields
output "location_custom_fields" {
  value       = data.netbox_location.by_id.custom_fields
  description = "All custom fields defined in NetBox for this location"
}

# Access specific custom fields by name
output "location_access_code" {
  value       = try([for cf in data.netbox_location.by_id.custom_fields : cf.value if cf.name == "access_code"][0], null)
  description = "Example: accessing a specific custom field value"
}

output "location_climate_controlled" {
  value       = try([for cf in data.netbox_location.by_id.custom_fields : cf.value if cf.name == "climate_controlled"][0], null)
  description = "Example: accessing a boolean custom field"
}
