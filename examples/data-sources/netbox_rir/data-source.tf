# Lookup by ID
data "netbox_rir" "by_id" {
  id = "123"
}

output "by_id" {
  value = data.netbox_rir.by_id.name
}

# Lookup by name
data "netbox_rir" "by_name" {
  name = "ARIN"
}

output "by_name" {
  value = data.netbox_rir.by_name.is_private
}

# Lookup by slug
data "netbox_rir" "by_slug" {
  slug = "arin"
}

output "by_slug" {
  value = data.netbox_rir.by_slug.name
}
