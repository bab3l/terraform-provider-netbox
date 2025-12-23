# Lookup by ID
data "netbox_role" "by_id" {
  id = "123"
}

output "by_id" {
  value = data.netbox_role.by_id.name
}

# Lookup by name
data "netbox_role" "by_name" {
  name = "Primary"
}

output "by_name" {
  value = data.netbox_role.by_name.slug
}

# Lookup by slug
data "netbox_role" "by_slug" {
  slug = "primary"
}

output "by_slug" {
  value = data.netbox_role.by_slug.weight
}
