# Tunnel Termination Resource Test
# This example demonstrates creating a tunnel termination in Netbox.

terraform {
  required_providers {
    netbox = {
      source = "bab3l/netbox"
    }
  }
}

provider "netbox" {
  # Uses NETBOX_SERVER_URL and NETBOX_API_TOKEN environment variables
}

# First, create a tunnel for the termination to belong to
resource "netbox_tunnel" "example" {
  name          = "test-tunnel-for-termination"
  encapsulation = "ipsec-tunnel"
  status        = "active"
  description   = "Test tunnel for demonstrating tunnel terminations"
}

# Create a basic tunnel termination as peer
resource "netbox_tunnel_termination" "example_peer" {
  tunnel           = netbox_tunnel.example.id
  termination_type = "dcim.device"
  role             = "peer"
}

# Create another termination as a hub
resource "netbox_tunnel_termination" "example_hub" {
  tunnel           = netbox_tunnel.example.id
  termination_type = "dcim.device"
  role             = "hub"
}

# Output the tunnel termination IDs
output "peer_termination_id" {
  description = "The ID of the peer tunnel termination"
  value       = netbox_tunnel_termination.example_peer.id
}

output "hub_termination_id" {
  description = "The ID of the hub tunnel termination"
  value       = netbox_tunnel_termination.example_hub.id
}

output "peer_termination_tunnel" {
  description = "The tunnel ID for the peer termination"
  value       = netbox_tunnel_termination.example_peer.tunnel
}

output "peer_termination_role" {
  description = "The role of the peer termination"
  value       = netbox_tunnel_termination.example_peer.role
}

# Validation outputs
output "peer_id_valid" {
  value = netbox_tunnel_termination.example_peer.id != "" && netbox_tunnel_termination.example_peer.id != null
}

output "peer_tunnel_valid" {
  value = netbox_tunnel_termination.example_peer.tunnel == netbox_tunnel.example.id
}

output "hub_role_valid" {
  value = netbox_tunnel_termination.example_hub.role == "hub"
}
