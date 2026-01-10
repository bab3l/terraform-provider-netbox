# Lookup by ID
data "netbox_provider_account" "by_id" {
  id = "123"
}

output "by_id" {
  value = data.netbox_provider_account.by_id.account
}

# Lookup by account identifier
data "netbox_provider_account" "by_account" {
  account = "ACC-123456789"
}

output "by_account" {
  value = data.netbox_provider_account.by_account.name
}

output "account_provider" {
  value = data.netbox_provider_account.by_id.provider
}

output "account_description" {
  value = data.netbox_provider_account.by_id.description
}

# Access all custom fields
output "account_custom_fields" {
  value       = data.netbox_provider_account.by_id.custom_fields
  description = "All custom fields defined in NetBox for this provider account"
}

# Access specific custom field by name
output "account_billing_email" {
  value       = try([for cf in data.netbox_provider_account.by_id.custom_fields : cf.value if cf.name == "billing_email"][0], null)
  description = "Example: accessing a text custom field for billing email"
}

output "account_monthly_cost" {
  value       = try([for cf in data.netbox_provider_account.by_id.custom_fields : cf.value if cf.name == "monthly_cost_usd"][0], null)
  description = "Example: accessing a numeric custom field for monthly cost"
}
