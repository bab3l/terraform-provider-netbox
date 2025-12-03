# Interface Data Source Test
# Tests the netbox_interface data source

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

# Prerequisites: Create resources to look up

resource "netbox_site" "test" {
  name        = "Interface DS Test Site"
  slug        = "interface-ds-test-site"
  description = "Site for interface data source testing"
}

resource "netbox_manufacturer" "test" {
  name        = "Interface DS Test Manufacturer"
  slug        = "interface-ds-test-manufacturer"
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = "Interface DS Test Model"
  slug         = "interface-ds-test-model"
}

resource "netbox_device_role" "test" {
  name = "Interface DS Test Device Role"
  slug = "interface-ds-test-device-role"
}

resource "netbox_device" "test" {
  name        = "interface-ds-test-device"
  site        = netbox_site.test.id
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
}

# Create an interface to look up
resource "netbox_interface" "test" {
  device      = netbox_device.test.id
  name        = "eth0"
  type        = "1000base-t"
  label       = "Primary NIC"
  enabled     = true
  mtu         = 1500
  description = "Test interface for data source"
  mode        = "access"
}

# Test 1: Look up by ID
data "netbox_interface" "by_id" {
  id = netbox_interface.test.id

  depends_on = [netbox_interface.test]
}

# Test 2: Look up by device ID and name
data "netbox_interface" "by_device_id_and_name" {
  device = netbox_device.test.id
  name   = "eth0"

  depends_on = [netbox_interface.test]
}

# Test 3: Look up by device name and interface name
data "netbox_interface" "by_device_name_and_name" {
  device = "interface-ds-test-device"
  name   = "eth0"

  depends_on = [netbox_interface.test]
}
