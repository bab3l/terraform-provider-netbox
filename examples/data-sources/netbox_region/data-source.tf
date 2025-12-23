# Lookup by ID
data "netbox_region" "by_id" {
  id = "123"
}

output "by_id" {
  value = data.netbox_region.by_id.name
}

# Lookup by slug
data "netbox_region" "by_slug" {
  slug = "north-america"
}

output "by_slug" {
  value = data.netbox_region.by_slug.name
}

# Lookup by name
data "netbox_region" "by_name" {
  name = "North America"
}

output "by_name" {
  value = data.netbox_region.by_name.slug
}
