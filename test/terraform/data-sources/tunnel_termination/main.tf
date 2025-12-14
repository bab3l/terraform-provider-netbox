# Tunnel Termination Data Source Test
# This example demonstrates reading tunnel termination information from Netbox.

terraform {
  required_version = ">= 1.0"
  required_providers {
    netbox = {
      source = "bab3l/netbox"
    }
  }
}

provider "netbox" {
  # Uses NETBOX_SERVER_URL and NETBOX_API_TOKEN environment variables
}

# First, create a tunnel and termination for the data source to read
resource "netbox_tunnel" "example" {
  name          = "test-tunnel-for-ds"
  encapsulation = "ipsec-tunnel"
  status        = "active"
}

resource "netbox_tunnel_termination" "example" {
  tunnel           = netbox_tunnel.example.id
  termination_type = "dcim.device"
  role             = "peer"
}

# Look up the tunnel termination by ID
data "netbox_tunnel_termination" "by_id" {
  id = netbox_tunnel_termination.example.id
}

# Look up the tunnel termination by tunnel
data "netbox_tunnel_termination" "by_tunnel" {
  tunnel = netbox_tunnel.example.id

  depends_on = [netbox_tunnel_termination.example]
}

# Output the retrieved data
output "by_id_id" {
  value = data.netbox_tunnel_termination.by_id.id
}

output "by_id_tunnel" {
  value = data.netbox_tunnel_termination.by_id.tunnel
}

output "by_id_termination_type" {
  value = data.netbox_tunnel_termination.by_id.termination_type
}

output "by_id_role" {
  value = data.netbox_tunnel_termination.by_id.role
}

output "by_tunnel_id" {
  value = data.netbox_tunnel_termination.by_tunnel.id
}

# Validation output - all lookups should return the same ID
output "all_ids_match" {
  value = data.netbox_tunnel_termination.by_id.id == data.netbox_tunnel_termination.by_tunnel.id
}
