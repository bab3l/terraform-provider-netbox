resource "netbox_tunnel" "test" {
  name          = "Test Tunnel"
  encapsulation = "ipsec-tunnel"
  status        = "active"
  description   = "IPSec tunnel between datacenter sites"
  comments      = "Primary site-to-site VPN tunnel"

  # Partial custom fields management
  # Only specified custom fields are managed by Terraform
  # Other custom fields set in NetBox UI or other tools are preserved
  custom_fields = [
    {
      name  = "tunnel_id"
      type  = "text"
      value = "TUN-001"
    },
    {
      name  = "encryption_algorithm"
      type  = "text"
      value = "AES-256-GCM"
    },
    {
      name  = "preshared_key_rotation_days"
      type  = "integer"
      value = "90"
    }
  ]

  tags = [
    {
      name = "ipsec"
      slug = "ipsec"
    },
    {
      name = "site-to-site"
      slug = "site-to-site"
    }
  ]
}

# Optional: seed owned custom fields during import
import {
  to = netbox_tunnel.test
  id = "123"

  identity = {
    custom_fields = [
      "tunnel_id:text",
      "encryption_algorithm:text",
      "preshared_key_rotation_days:integer",
    ]
  }
}
