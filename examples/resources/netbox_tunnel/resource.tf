resource "netbox_tunnel" "test" {
  name          = "Test Tunnel"
  encapsulation = "ipsec-tunnel"
  status        = "active"
  description   = "IPSec tunnel between datacenter sites"
  comments      = "Primary site-to-site VPN tunnel"

  # Partial custom fields management
  # Only specified custom fields are managed, others in NetBox preserved
  custom_fields = [
    {
      name  = "tunnel_id"
      value = "TUN-001"
    },
    {
      name  = "encryption_algorithm"
      value = "AES-256-GCM"
    },
    {
      name  = "preshared_key_rotation_days"
      value = "90"
    }
  ]

  tags = [
    "ipsec",
    "site-to-site"
  ]
}
