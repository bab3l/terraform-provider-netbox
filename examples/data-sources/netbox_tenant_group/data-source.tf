# Look up tenant group by ID
data "netbox_tenant_group" "by_id" {
  id = "1"
}

# Look up tenant group by slug
data "netbox_tenant_group" "by_slug" {
  slug = "example-group"
}

# Look up tenant group by name
data "netbox_tenant_group" "by_name" {
  name = "Example Group"
}

# Use tenant group data in other resources
output "tenant_group_name" {
  value = data.netbox_tenant_group.by_name.name
}

output "tenant_group_slug" {
  value = data.netbox_tenant_group.by_slug.slug
}

output "tenant_group_parent" {
  value = data.netbox_tenant_group.by_id.parent
}

output "tenant_group_description" {
  value = data.netbox_tenant_group.by_id.description
}

# Access all custom fields
output "tenant_group_custom_fields" {
  value       = data.netbox_tenant_group.by_id.custom_fields
  description = "All custom fields defined in NetBox for this tenant group"
}

# Access specific custom fields by name
output "tenant_group_account_manager" {
  value       = try([for cf in data.netbox_tenant_group.by_id.custom_fields : cf.value if cf.name == "account_manager"][0], null)
  description = "Example: accessing a text custom field"
}

output "tenant_group_industry" {
  value       = try([for cf in data.netbox_tenant_group.by_id.custom_fields : cf.value if cf.name == "industry"][0], null)
  description = "Example: accessing a select custom field"
}
