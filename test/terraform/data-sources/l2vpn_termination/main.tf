# L2VPN Termination Data Source Test

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

# Create pre-requisites
resource "netbox_l2vpn" "test" {
  name = "test-l2vpn-for-term-ds"
  slug = "test-l2vpn-for-term-ds"
  type = "vpls"
}

resource "netbox_vlan" "test" {
  name   = "test-vlan-for-l2vpn-term-ds"
  vid    = 3998
  status = "active"
}

# Create termination resource
resource "netbox_l2vpn_termination" "test" {
  l2vpn                = netbox_l2vpn.test.id
  assigned_object_type = "ipam.vlan"
  assigned_object_id   = netbox_vlan.test.id
}

# Test: Look up by ID
data "netbox_l2vpn_termination" "by_id" {
  id = netbox_l2vpn_termination.test.id
}

# Outputs for verification
output "termination_id" {
  value = data.netbox_l2vpn_termination.by_id.id
}

output "termination_l2vpn" {
  value = data.netbox_l2vpn_termination.by_id.l2vpn
}

output "termination_object_type" {
  value = data.netbox_l2vpn_termination.by_id.assigned_object_type
}

output "termination_object_id" {
  value = data.netbox_l2vpn_termination.by_id.assigned_object_id
}

# Validation outputs
output "id_valid" {
  value = data.netbox_l2vpn_termination.by_id.id == netbox_l2vpn_termination.test.id
}

output "l2vpn_valid" {
  value = data.netbox_l2vpn_termination.by_id.l2vpn == netbox_l2vpn.test.id
}

output "object_type_valid" {
  value = data.netbox_l2vpn_termination.by_id.assigned_object_type == "ipam.vlan"
}
