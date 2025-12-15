resource "netbox_circuit_type" "test" {
  name = "Test Type"
  slug = "test-type"
}

resource "netbox_provider" "test" {
  name = "Test Provider"
  slug = "test-provider"
}

resource "netbox_circuit" "test" {
  cid             = "123456789"
  provider_id     = netbox_provider.test.id
  circuit_type_id = netbox_circuit_type.test.id
  status          = "active"
}

resource "netbox_circuit_group" "test" {
  name = "Test Circuit Group"
  slug = "test-circuit-group"
}

resource "netbox_circuit_group_assignment" "test" {
  circuit_id = netbox_circuit.test.id
  group_id   = netbox_circuit_group.test.id
}
