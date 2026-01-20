resource "netbox_provider" "test" {
  name = "Test Provider"
  slug = "test-provider"
}

resource "netbox_provider_network" "test" {
  name             = "Test Provider Network"
  circuit_provider = netbox_provider.test.name
  description      = "Provider backbone network"
  comments         = "Redundant provider network with multiple POPs"

  # Partial custom fields management
  # Manage specific custom fields while preserving others
  custom_fields = [
    {
      name  = "network_tier"
      value = "tier-1"
    },
    {
      name  = "asn"
      value = "65001"
    },
    {
      name  = "ipv6_enabled"
      value = "true"
    }
  ]

  tags = [
    "backbone",
    "ipv6-capable"
  ]
}

# Optional: seed owned custom fields during import
import {
  to = netbox_provider_network.test
  id = "123"

  identity = {
    custom_fields = [
      "network_tier:text",
      "asn:integer",
      "ipv6_enabled:boolean",
    ]
  }
}
