# FHRP Group Assignment Resource Integration Test
# Tests the netbox_fhrp_group_assignment resource

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
  name   = "Test Site for FHRP Assignment"
  slug   = "test-site-fhrp-assignment"
  status = "active"
}

# Dependencies - Device type
resource "netbox_manufacturer" "test" {
  name = "Test Manufacturer for FHRP Assignment"
  slug = "test-manufacturer-fhrp-assignment"
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.slug
  model        = "Test Model for FHRP Assignment"
  slug         = "test-model-fhrp-assignment"
}

resource "netbox_device_role" "test" {
  name  = "Test Role for FHRP Assignment"
  slug  = "test-role-fhrp-assignment"
  color = "aabbcc"
}

# Dependencies - Device
resource "netbox_device" "test" {
  name        = "Test Device for FHRP Assignment"
  device_type = netbox_device_type.test.model
  role        = netbox_device_role.test.slug
  site        = netbox_site.test.slug
  status      = "active"
}

# Dependencies - Interface
resource "netbox_interface" "test" {
  name   = "eth0"
  device = netbox_device.test.name
  type   = "1000base-t"
}

resource "netbox_interface" "test2" {
  name   = "eth1"
  device = netbox_device.test.name
  type   = "1000base-t"
}

# Dependencies - FHRP Group
resource "netbox_fhrp_group" "test_vrrp" {
  protocol    = "vrrp2"
  group_id    = 200
  name        = "Test VRRP Group for Assignment"
  description = "VRRP group for assignment testing"
}

resource "netbox_fhrp_group" "test_hsrp" {
  protocol    = "hsrp"
  group_id    = 100
  name        = "Test HSRP Group for Assignment"
  description = "HSRP group for assignment testing"
}

# Test 1: Basic FHRP group assignment
resource "netbox_fhrp_group_assignment" "basic" {
  group_id       = netbox_fhrp_group.test_vrrp.id
  interface_type = "dcim.interface"
  interface_id   = netbox_interface.test.id
  priority       = 100
}

# Test 2: FHRP group assignment with different priority
resource "netbox_fhrp_group_assignment" "high_priority" {
  group_id       = netbox_fhrp_group.test_hsrp.id
  interface_type = "dcim.interface"
  interface_id   = netbox_interface.test2.id
  priority       = 255
}
