data "netbox_provider_account" "test" {
  account = "123456789"
}

output "example" {
  value = data.netbox_provider_account.test.id
}
