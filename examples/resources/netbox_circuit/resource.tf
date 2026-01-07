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
  comments         = "Primary internet connection for datacenter"

  # Partial custom fields management
  # Only manage specific custom fields, others in NetBox are preserved
  custom_fields = [
    {
      name  = "circuit_vendor_id"
      value = "VENDOR-12345"
    },
    {
      name  = "monthly_cost"
      value = "5000"
    },
    {
      name  = "contract_end_date"
      value = "2026-12-31"
    }
  ]

  tags = [
    "production",
    "primary-circuit"
  ]
}
