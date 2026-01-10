# Lookup by ID
data "netbox_rir" "by_id" {
  id = "123"
}

# Lookup by name
data "netbox_rir" "by_name" {
  name = "ARIN"
}

# Lookup by slug
data "netbox_rir" "by_slug" {
  slug = "arin"
}

# Use RIR data in other resources
output "rir_name" {
  value = data.netbox_rir.by_id.name
}

output "rir_slug" {
  value = data.netbox_rir.by_slug.slug
}

output "rir_is_private" {
  value = data.netbox_rir.by_name.is_private
}

output "rir_description" {
  value = data.netbox_rir.by_id.description
}

# Access all custom fields
output "rir_custom_fields" {
  value       = data.netbox_rir.by_id.custom_fields
  description = "All custom fields defined in NetBox for this RIR"
}

# Access specific custom fields by name
output "rir_registry_url" {
  value       = try([for cf in data.netbox_rir.by_id.custom_fields : cf.value if cf.name == "registry_url"][0], null)
  description = "Example: accessing a URL custom field"
}

output "rir_contact_email" {
  value       = try([for cf in data.netbox_rir.by_id.custom_fields : cf.value if cf.name == "contact_email"][0], null)
  description = "Example: accessing a text custom field"
}
