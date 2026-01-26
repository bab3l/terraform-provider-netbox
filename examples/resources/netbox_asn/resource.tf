resource "netbox_rir" "test" {
  name = "RIPE"
  slug = "ripe"
}

resource "netbox_asn" "test" {
  asn         = 65001
  rir         = netbox_rir.test.id
  description = "Primary BGP ASN for datacenter"
  comments    = "Private ASN for internal BGP routing"

  # Partial custom fields management
  # Only specified custom fields are managed, others preserved
  custom_fields = [
    {
      name  = "asn_type"
      value = "private"
    },
    {
      name  = "routing_policy"
      value = "internal-only"
    },
    {
      name  = "contact_email"
      value = "neteng@example.com"
    }
  ]

  tags = [
    "private-asn",
    "bgp"
  ]
}

# Optional: seed owned custom fields during import
import {
  to = netbox_asn.test
  id = "123"

  identity = {
    custom_fields = [
      "asn_type:text",
      "routing_policy:text",
      "contact_email:text",
    ]
  }
}
