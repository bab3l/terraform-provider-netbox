# Lookup by ID
data "netbox_provider" "by_id" {
  id = "123"
}

output "by_id" {
  value = data.netbox_provider.by_id.name
}

# Lookup by slug
data "netbox_provider" "by_slug" {
  slug = "isp-provider"
}

output "by_slug" {
  value = data.netbox_provider.by_slug.name
}

# Lookup by name
data "netbox_provider" "by_name" {
  name = "ISP Provider"
}

output "by_name" {
  value = data.netbox_provider.by_name.slug
}
