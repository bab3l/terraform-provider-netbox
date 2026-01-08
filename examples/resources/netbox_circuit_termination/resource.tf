resource "netbox_site" "test" {
  name   = "Test Site"
  slug   = "test-site"
  status = "active"
}

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
  circuit_provider = netbox_provider.test.id
  type             = netbox_circuit_type.test.id
  status           = "active"
}

resource "netbox_circuit_termination" "test_a" {
  circuit     = netbox_circuit.test.id
  term_side   = "A"
  site        = netbox_site.test.name
  port_speed  = 1000000 # 1 Gbps
  description = "Datacenter A termination"

  # Partial custom fields management
  # Manage specific custom fields while preserving others in NetBox
  custom_fields = [
    {
      name  = "demarcation_point"
      value = "DC-A-MMR-RACK-42"
    },
    {
      name  = "cross_connect_id"
      value = "XC-2024-001"
    }
  ]

  tags = [
    "datacenter-a",
    "primary"
  ]
}
