resource "netbox_circuit_type" "test" {
  name = "Test Type"
  slug = "test-type"
}

resource "netbox_provider" "test" {
  name = "Test Provider"
  slug = "test-provider"
}

resource "netbox_circuit" "test" {
  cid              = "123456789"
  circuit_provider = netbox_provider.test.name
  type             = netbox_circuit_type.test.name
  status           = "active"
}

resource "netbox_circuit_group" "test" {
  name = "Test Circuit Group"
  slug = "test-circuit-group"
}

resource "netbox_circuit_group_assignment" "test" {
  circuit_id = netbox_circuit.test.id
  group      = netbox_circuit_group.test.name
  priority   = "primary"

  # Note: circuit_group_assignment is tags-only (no custom_fields support)
  # Tags use replace-all semantics (not merge like custom fields)
  tags = [
    "circuit-assignment",
    "primary-link"
  ]
}
