# Provider Account Outputs

# Basic provider account outputs
output "basic_id" {
  value = netbox_provider_account.basic.id
}

output "basic_name" {
  value = netbox_provider_account.basic.name
}

# Complete provider account outputs
output "complete_id" {
  value = netbox_provider_account.complete.id
}

output "complete_name" {
  value = netbox_provider_account.complete.name
}

output "complete_account" {
  value = netbox_provider_account.complete.account
}
