# FHRP Group Assignment Data Source Integration Test
# Tests the netbox_fhrp_group_assignment data source

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

# Dependencies - Site
resource "netbox_site" "test" {
  name   = "Test Site for FHRP Assignment DS"
  slug   = "test-site-fhrp-assignment-ds"
  status = "active"
}

# Dependencies - Device type
resource "netbox_manufacturer" "test" {
  name = "Test Manufacturer for FHRP Assignment DS"
  slug = "test-manufacturer-fhrp-assignment-ds"
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = "Test Model for FHRP Assignment DS"
  slug         = "test-model-fhrp-assignment-ds"
}

resource "netbox_device_role" "test" {
  name  = "Test Role for FHRP Assignment DS"
  slug  = "test-role-fhrp-assignment-ds"
  color = "aabbcc"
}

# Dependencies - Device
resource "netbox_device" "test" {
  name        = "Test Device for FHRP Assignment DS"
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
  status      = "active"
}

# Dependencies - Interface
resource "netbox_interface" "test" {
  name   = "eth0"
  device = netbox_device.test.id
  type   = "1000base-t"
}

# Dependencies - FHRP Group
resource "netbox_fhrp_group" "test" {
  protocol    = "vrrp2"
  group_id    = 201
  name        = "Test VRRP Group for Assignment DS"
  description = "VRRP group for assignment data source testing"
}

# Create FHRP group assignment to look up
resource "netbox_fhrp_group_assignment" "test" {
  group_id       = netbox_fhrp_group.test.id
  interface_type = "dcim.interface"
  interface_id   = netbox_interface.test.id
  priority       = 150
}

# Look up FHRP group assignment by ID
data "netbox_fhrp_group_assignment" "by_id" {
  id = netbox_fhrp_group_assignment.test.id
}

# Outputs for verification
output "by_id_group_id" {
  value = data.netbox_fhrp_group_assignment.by_id.group_id
}

output "by_id_interface_type" {
  value = data.netbox_fhrp_group_assignment.by_id.interface_type
}

output "by_id_interface_id" {
  value = data.netbox_fhrp_group_assignment.by_id.interface_id
}

output "by_id_priority" {
  value = data.netbox_fhrp_group_assignment.by_id.priority
}

# Validation outputs
output "id_lookup_matches" {
  value = data.netbox_fhrp_group_assignment.by_id.id == netbox_fhrp_group_assignment.test.id
}

output "group_matches" {
  value = data.netbox_fhrp_group_assignment.by_id.group_id == tostring(netbox_fhrp_group.test.id)
}

output "interface_matches" {
  value = data.netbox_fhrp_group_assignment.by_id.interface_id == netbox_interface.test.id
}

output "priority_matches" {
  value = data.netbox_fhrp_group_assignment.by_id.priority == 150
}
