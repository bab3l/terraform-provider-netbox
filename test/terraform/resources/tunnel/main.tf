# Tunnel Resource Test

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

# Test 1: Basic tunnel with required fields only
resource "netbox_tunnel" "basic" {
  name          = "test-tunnel-basic"
  status        = "active"
  encapsulation = "gre"
}

# Test 2: Tunnel with all optional fields
resource "netbox_tunnel" "complete" {
  name          = "test-tunnel-complete"
  status        = "planned"
  encapsulation = "wireguard"
  description   = "Complete tunnel for integration testing"
  comments      = "This tunnel was created for integration testing."
  tunnel_id     = 12345
}

# Test 3: IPSec tunnel with tunnel group
resource "netbox_tunnel_group" "test" {
  name = "test-tunnel-group-for-tunnel"
  slug = "test-tunnel-group-for-tunnel"
}

resource "netbox_tunnel" "with_group" {
  name          = "test-tunnel-with-group"
  status        = "active"
  encapsulation = "ipsec-tunnel"
  group         = netbox_tunnel_group.test.id
}

# Test 4: Output values for verification
output "basic_tunnel_id" {
  value = netbox_tunnel.basic.id
}

output "basic_tunnel_name" {
  value = netbox_tunnel.basic.name
}

output "basic_tunnel_status" {
  value = netbox_tunnel.basic.status
}

output "complete_tunnel_name" {
  value = netbox_tunnel.complete.name
}

output "complete_tunnel_description" {
  value = netbox_tunnel.complete.description
}

output "complete_tunnel_encapsulation" {
  value = netbox_tunnel.complete.encapsulation
}

output "with_group_tunnel_group" {
  value = netbox_tunnel.with_group.group
}

# Validation outputs
output "basic_id_valid" {
  value = netbox_tunnel.basic.id != "" && netbox_tunnel.basic.id != null
}

output "basic_name_valid" {
  value = netbox_tunnel.basic.name == "test-tunnel-basic"
}

output "complete_description_valid" {
  value = netbox_tunnel.complete.description == "Complete tunnel for integration testing"
}

output "with_group_group_valid" {
  value = netbox_tunnel.with_group.group == netbox_tunnel_group.test.id
}
