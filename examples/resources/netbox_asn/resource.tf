resource "netbox_rir" "test" {
  name = "RIPE"
  slug = "ripe"
}

resource "netbox_asn" "test" {
  asn         = 65001
  rir         = netbox_rir.test.name
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
