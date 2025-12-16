resource "netbox_provider" "test" {
  name = "Test Provider"
  slug = "test-provider"
}

resource "netbox_provider_account" "test" {
  name             = "Test Account"
  account          = "1234567890"
  circuit_provider = netbox_provider.test.name
}
