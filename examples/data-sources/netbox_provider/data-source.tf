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

output "provider_slug" {
  value = data.netbox_provider.by_id.slug
}

output "provider_description" {
  value = data.netbox_provider.by_id.description
}

# Access all custom fields
output "provider_custom_fields" {
  value       = data.netbox_provider.by_id.custom_fields
  description = "All custom fields defined in NetBox for this provider"
}

# Access specific custom field by name
output "provider_account_manager" {
  value       = try([for cf in data.netbox_provider.by_id.custom_fields : cf.value if cf.name == "account_manager"][0], null)
  description = "Example: accessing a text custom field for account manager"
}

output "provider_sla_tier" {
  value       = try([for cf in data.netbox_provider.by_id.custom_fields : cf.value if cf.name == "sla_tier"][0], null)
  description = "Example: accessing a select custom field for SLA tier"
}
