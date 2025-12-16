resource "netbox_provider" "test" {
  name = "Test Provider"
  slug = "test-provider"
}

resource "netbox_provider_network" "test" {
  name             = "Test Provider Network"
  circuit_provider = netbox_provider.test.name
}
