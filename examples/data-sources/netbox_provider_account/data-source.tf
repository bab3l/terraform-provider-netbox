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
