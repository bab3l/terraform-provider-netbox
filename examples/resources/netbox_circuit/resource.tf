resource "netbox_provider" "test" {
  name = "Test Provider"
  slug = "test-provider"
}

resource "netbox_circuit_type" "test" {
  name = "Internet Transit"
  slug = "internet-transit"
}

resource "netbox_provider_account" "test" {
  account          = "ACCT-100"
  circuit_provider = netbox_provider.test.slug
}

resource "netbox_circuit" "test" {
  cid              = "CID-12345"
  circuit_provider = netbox_provider.test.name
  provider_account = netbox_provider_account.test.account
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

# Optional: seed owned custom fields during import
import {
  to = netbox_circuit.test
  id = "123"

  identity = {
    custom_fields = [
      "circuit_vendor_id:text",
      "monthly_cost:integer",
      "contract_end_date:date",
    ]
  }
}
