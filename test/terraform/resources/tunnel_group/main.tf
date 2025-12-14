# Tunnel Group Resource Test

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

# Test 1: Basic tunnel group with required fields only
resource "netbox_tunnel_group" "basic" {
  name = "test-tunnel-group-basic"
  slug = "test-tunnel-group-basic"
}

# Test 2: Tunnel group with all optional fields
resource "netbox_tunnel_group" "complete" {
  name        = "test-tunnel-group-complete"
  slug        = "test-tunnel-group-complete"
  description = "Complete tunnel group for integration testing"
}

# Test 3: Output values for verification
output "basic_tunnel_group_id" {
  value = netbox_tunnel_group.basic.id
}

output "basic_tunnel_group_name" {
  value = netbox_tunnel_group.basic.name
}

output "basic_tunnel_group_slug" {
  value = netbox_tunnel_group.basic.slug
}

output "complete_tunnel_group_name" {
  value = netbox_tunnel_group.complete.name
}

output "complete_tunnel_group_description" {
  value = netbox_tunnel_group.complete.description
}

# Validation outputs
output "basic_id_valid" {
  value = netbox_tunnel_group.basic.id != "" && netbox_tunnel_group.basic.id != null
}

output "basic_name_valid" {
  value = netbox_tunnel_group.basic.name == "test-tunnel-group-basic"
}

output "complete_description_valid" {
  value = netbox_tunnel_group.complete.description == "Complete tunnel group for integration testing"
}
