# Look up tenant by ID
data "netbox_tenant" "by_id" {
  id = "1"
}

# Look up tenant by slug
data "netbox_tenant" "by_slug" {
  slug = "example-tenant"
}

# Look up tenant by name
data "netbox_tenant" "by_name" {
  name = "Example Tenant"
}

# Use tenant data in other resources
output "tenant_name" {
  value = data.netbox_tenant.by_id.name
}

output "tenant_slug" {
  value = data.netbox_tenant.by_slug.slug
}

output "tenant_group" {
  value = data.netbox_tenant.by_id.group
}

output "tenant_description" {
  value = data.netbox_tenant.by_id.description
}

output "tenant_comments" {
  value = data.netbox_tenant.by_name.comments
}

# Access all custom fields
output "tenant_custom_fields" {
  value       = data.netbox_tenant.by_id.custom_fields
  description = "All custom fields defined in NetBox for this tenant"
}

# Access specific custom fields by name
output "tenant_account_number" {
  value       = try([for cf in data.netbox_tenant.by_id.custom_fields : cf.value if cf.name == "account_number"][0], null)
  description = "Example: accessing a text custom field"
}

output "tenant_billing_contact" {
  value       = try([for cf in data.netbox_tenant.by_id.custom_fields : cf.value if cf.name == "billing_contact"][0], null)
  description = "Example: accessing a text custom field"
}

output "tenant_active" {
  value       = try([for cf in data.netbox_tenant.by_id.custom_fields : cf.value if cf.name == "is_active"][0], null)
  description = "Example: accessing a boolean custom field"
}
