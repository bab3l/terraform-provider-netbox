data "netbox_provider" "test" {
  name = "test-provider"
}

output "example" {
  value = data.netbox_provider.test.id
}
