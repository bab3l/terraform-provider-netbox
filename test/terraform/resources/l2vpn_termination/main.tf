# L2VPN Termination Resource Test

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

# Pre-requisites: Create an L2VPN for the termination
resource "netbox_l2vpn" "test" {
  name = "test-l2vpn-for-termination"
  slug = "test-l2vpn-for-termination"
  type = "vpls"
}

# Pre-requisites: Create VLAN for termination
resource "netbox_vlan" "test" {
  name   = "test-vlan-for-l2vpn-term"
  vid    = 3999
  status = "active"
}

# Test 1: Basic L2VPN termination with VLAN
resource "netbox_l2vpn_termination" "vlan" {
  l2vpn                = netbox_l2vpn.test.id
  assigned_object_type = "ipam.vlan"
  assigned_object_id   = netbox_vlan.test.id
}

# Test 2: Output values for verification
output "termination_id" {
  value = netbox_l2vpn_termination.vlan.id
}

output "termination_l2vpn" {
  value = netbox_l2vpn_termination.vlan.l2vpn
}

output "termination_object_type" {
  value = netbox_l2vpn_termination.vlan.assigned_object_type
}

output "termination_object_id" {
  value = netbox_l2vpn_termination.vlan.assigned_object_id
}

# Validation outputs
output "termination_id_valid" {
  value = netbox_l2vpn_termination.vlan.id != "" && netbox_l2vpn_termination.vlan.id != null
}

output "termination_l2vpn_valid" {
  value = netbox_l2vpn_termination.vlan.l2vpn == netbox_l2vpn.test.id
}

output "termination_object_type_valid" {
  value = netbox_l2vpn_termination.vlan.assigned_object_type == "ipam.vlan"
}
