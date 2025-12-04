# Cable Resource Test
# Tests the netbox_cable resource by creating cables between interfaces

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

# Prerequisites: Site, Manufacturer, Device Type, Device Role, Devices, and Interfaces
# Cables connect interfaces on different devices

resource "netbox_site" "test" {
  name        = "Cable Test Site"
  slug        = "cable-test-site"
  description = "Site for cable testing"
}

resource "netbox_manufacturer" "test" {
  name        = "Cable Test Manufacturer"
  slug        = "cable-test-manufacturer"
  description = "Manufacturer for cable testing"
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = "Cable Test Model"
  slug         = "cable-test-model"
  description  = "Device type for cable testing"
}

resource "netbox_device_role" "test" {
  name        = "Cable Test Device Role"
  slug        = "cable-test-device-role"
  description = "Device role for cable testing"
}

# Device A - source of cables
resource "netbox_device" "device_a" {
  name        = "cable-test-device-a"
  site        = netbox_site.test.id
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  description = "Device A for cable testing"
}

# Device B - destination of cables
resource "netbox_device" "device_b" {
  name        = "cable-test-device-b"
  site        = netbox_site.test.id
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  description = "Device B for cable testing"
}

# Interfaces on Device A
resource "netbox_interface" "device_a_eth0" {
  device      = netbox_device.device_a.id
  name        = "eth0"
  type        = "1000base-t"
  description = "Device A eth0"
}

resource "netbox_interface" "device_a_eth1" {
  device      = netbox_device.device_a.id
  name        = "eth1"
  type        = "1000base-t"
  description = "Device A eth1"
}

resource "netbox_interface" "device_a_eth2" {
  device      = netbox_device.device_a.id
  name        = "eth2"
  type        = "10gbase-t"
  description = "Device A eth2"
}

# Interfaces on Device B
resource "netbox_interface" "device_b_eth0" {
  device      = netbox_device.device_b.id
  name        = "eth0"
  type        = "1000base-t"
  description = "Device B eth0"
}

resource "netbox_interface" "device_b_eth1" {
  device      = netbox_device.device_b.id
  name        = "eth1"
  type        = "1000base-t"
  description = "Device B eth1"
}

resource "netbox_interface" "device_b_eth2" {
  device      = netbox_device.device_b.id
  name        = "eth2"
  type        = "10gbase-t"
  description = "Device B eth2"
}

# Basic cable with only required fields
resource "netbox_cable" "basic" {
  a_terminations = [{
    object_type = "dcim.interface"
    object_id   = netbox_interface.device_a_eth0.id
  }]
  b_terminations = [{
    object_type = "dcim.interface"
    object_id   = netbox_interface.device_b_eth0.id
  }]
}

# Cable with type specified
resource "netbox_cable" "with_type" {
  a_terminations = [{
    object_type = "dcim.interface"
    object_id   = netbox_interface.device_a_eth1.id
  }]
  b_terminations = [{
    object_type = "dcim.interface"
    object_id   = netbox_interface.device_b_eth1.id
  }]
  type = "cat6a"
}

# Cable with full details
resource "netbox_cable" "full" {
  a_terminations = [{
    object_type = "dcim.interface"
    object_id   = netbox_interface.device_a_eth2.id
  }]
  b_terminations = [{
    object_type = "dcim.interface"
    object_id   = netbox_interface.device_b_eth2.id
  }]
  type        = "cat6"
  status      = "connected"
  label       = "CABLE-001"
  color       = "0000ff"
  length      = 5.5
  length_unit = "m"
  description = "Ethernet cable from Device A to Device B"
  comments    = "Test cable with all optional fields"
}
