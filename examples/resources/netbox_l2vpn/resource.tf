resource "netbox_l2vpn" "test" {
  name        = "Test L2VPN"
  slug        = "test-l2vpn"
  type        = "vxlan"
  identifier  = "1000"
  description = "VXLAN overlay for datacenter interconnect"
  comments    = "Layer 2 VPN for stretched VLANs"

  # Partial custom fields management
  # Only specified custom fields are managed, others in NetBox preserved
  custom_fields = [
    {
      name  = "vni"
      value = "10000"
    },
    {
      name  = "mtu"
      value = "9000"
    },
    {
      name  = "encryption_enabled"
      value = "true"
    }
  ]

  tags = [
    "vxlan",
    "datacenter-interconnect"
  ]
}
