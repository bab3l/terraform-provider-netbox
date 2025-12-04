# Cable Data Source Test
# Tests the netbox_cable data source by looking up cables

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

# Prerequisites: Site, Manufacturer, Device Type, Device Role, Devices, Interfaces, and Cable

resource "netbox_site" "test" {
  name        = "Cable DS Test Site"
  slug        = "cable-ds-test-site"
  description = "Site for cable data source testing"
}

resource "netbox_manufacturer" "test" {
  name        = "Cable DS Test Manufacturer"
  slug        = "cable-ds-test-manufacturer"
  description = "Manufacturer for cable data source testing"
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = "Cable DS Test Model"
  slug         = "cable-ds-test-model"
  description  = "Device type for cable data source testing"
}

resource "netbox_device_role" "test" {
  name        = "Cable DS Test Device Role"
  slug        = "cable-ds-test-device-role"
  description = "Device role for cable data source testing"
}

# Device A
resource "netbox_device" "device_a" {
  name        = "cable-ds-test-device-a"
  site        = netbox_site.test.id
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  description = "Device A for cable data source testing"
}

# Device B
resource "netbox_device" "device_b" {
  name        = "cable-ds-test-device-b"
  site        = netbox_site.test.id
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  description = "Device B for cable data source testing"
}

# Interface on Device A
resource "netbox_interface" "device_a_eth0" {
  device      = netbox_device.device_a.id
  name        = "eth0"
  type        = "1000base-t"
  description = "Device A eth0 for cable DS test"
}

# Interface on Device B
resource "netbox_interface" "device_b_eth0" {
  device      = netbox_device.device_b.id
  name        = "eth0"
  type        = "1000base-t"
  description = "Device B eth0 for cable DS test"
}

# Create a cable to look up
resource "netbox_cable" "test" {
  a_terminations = [{
    object_type = "dcim.interface"
    object_id   = netbox_interface.device_a_eth0.id
  }]
  b_terminations = [{
    object_type = "dcim.interface"
    object_id   = netbox_interface.device_b_eth0.id
  }]
  type        = "cat6"
  status      = "connected"
  label       = "DS-TEST-CABLE"
  color       = "ff0000"
  length      = 10
  length_unit = "m"
  description = "Cable for data source testing"
}

# Look up by ID
data "netbox_cable" "by_id" {
  id = netbox_cable.test.id
}
