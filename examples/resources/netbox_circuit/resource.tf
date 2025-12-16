resource "netbox_provider" "test" {
  name = "Test Provider"
  slug = "test-provider"
}

resource "netbox_circuit_type" "test" {
  name = "Internet Transit"
  slug = "internet-transit"
}

resource "netbox_circuit" "test" {
  cid              = "CID-12345"
  circuit_provider = netbox_provider.test.name
  type             = netbox_circuit_type.test.name
  status           = "active"
  description      = "Main Internet Circuit"
}
