resource "netbox_tunnel_group" "test" {
  name        = "Test Tunnel Group"
  slug        = "test-tunnel-group"
  description = "Group of related site-to-site tunnels"

  # Partial custom fields management
  # Only specified custom fields are managed, others preserved
  custom_fields = [
    {
      name  = "group_purpose"
      value = "datacenter-mesh"
    },
    {
      name  = "redundancy_level"
      value = "high"
    }
  ]

  tags = [
    "tunnel-group",
    "production"
  ]
}

# Optional: seed owned custom fields during import
import {
  to = netbox_tunnel_group.test
  id = "123"

  identity = {
    custom_fields = [
      "group_purpose:text",
      "redundancy_level:text",
    ]
  }
}
