# FHRP Group Data Source Integration Test
# Tests the netbox_fhrp_group data source for looking up existing FHRP groups

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

# Create FHRP groups to look up
resource "netbox_fhrp_group" "test_vrrp2" {
  protocol    = "vrrp2"
  group_id    = 100
  name        = "Test VRRP2 Group for Data Source"
  description = "Created for data source testing"
}

resource "netbox_fhrp_group" "test_hsrp" {
  protocol    = "hsrp"
  group_id    = 50
  name        = "Test HSRP Group for Data Source"
  description = "Created for data source testing"
  auth_type   = "plaintext"
  auth_key    = "testkey123"
}

# Look up FHRP group by ID
data "netbox_fhrp_group" "by_id" {
  id = netbox_fhrp_group.test_vrrp2.id
}

# Look up FHRP group by protocol and group_id
data "netbox_fhrp_group" "by_protocol_and_group_id" {
  protocol = netbox_fhrp_group.test_hsrp.protocol
  group_id = netbox_fhrp_group.test_hsrp.group_id
}

# Outputs for verification
output "by_id_protocol" {
  value = data.netbox_fhrp_group.by_id.protocol
}

output "by_id_group_id" {
  value = data.netbox_fhrp_group.by_id.group_id
}

output "by_id_name" {
  value = data.netbox_fhrp_group.by_id.name
}

output "by_protocol_name" {
  value = data.netbox_fhrp_group.by_protocol_and_group_id.name
}

output "by_protocol_auth_type" {
  value = data.netbox_fhrp_group.by_protocol_and_group_id.auth_type
}
