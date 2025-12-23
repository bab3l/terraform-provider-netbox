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

output "location_facility" {
  value = data.netbox_location.by_slug.facility
}
