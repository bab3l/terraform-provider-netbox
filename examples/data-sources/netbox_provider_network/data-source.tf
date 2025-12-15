data "netbox_provider_network" "test" {
  name = "test-provider-network"
}

output "example" {
  value = data.netbox_provider_network.test.id
}
