# Lookup by ID
data "netbox_rack_type" "by_id" {
  id = "123"
}

output "by_id" {
  value = data.netbox_rack_type.by_id.model
}

# Lookup by slug
data "netbox_rack_type" "by_slug" {
  slug = "42u-2post-frame"
}

output "by_slug" {
  value = data.netbox_rack_type.by_slug.u_height
}

# Lookup by model name
data "netbox_rack_type" "by_model" {
  model = "FS-42"
}

output "by_model" {
  value = data.netbox_rack_type.by_model.manufacturer
}
