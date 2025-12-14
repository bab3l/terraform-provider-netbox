# Interface Resource Test
# Tests the netbox_interface resource by creating interfaces on a device

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

# Prerequisites: Site, Manufacturer, Device Type, Device Role, and Device

resource "netbox_site" "test" {
  name        = "Interface Test Site"
  slug        = "interface-test-site"
  description = "Site for interface testing"
}

resource "netbox_manufacturer" "test" {
  name        = "Interface Test Manufacturer"
  slug        = "interface-test-manufacturer"
  description = "Manufacturer for interface testing"
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = "Interface Test Model"
  slug         = "interface-test-model"
  description  = "Device type for interface testing"
}

resource "netbox_device_role" "test" {
  name        = "Interface Test Device Role"
  slug        = "interface-test-device-role"
  description = "Device role for interface testing"
}

resource "netbox_device" "test" {
  name        = "interface-test-device"
  site        = netbox_site.test.id
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  description = "Device for interface testing"
}

# Basic interface with only required fields
resource "netbox_interface" "basic" {
  device = netbox_device.test.id
  name   = "eth0"
  type   = "1000base-t"
}

# Interface with description and label
resource "netbox_interface" "with_label" {
  device      = netbox_device.test.id
  name        = "eth1"
  type        = "1000base-t"
  label       = "Main uplink"
  description = "Primary uplink interface"
}

# Disabled interface
resource "netbox_interface" "disabled" {
  device      = netbox_device.test.id
  name        = "eth2"
  type        = "1000base-t"
  enabled     = false
  description = "Disabled interface for testing"
}

# Interface with MTU and speed
resource "netbox_interface" "with_mtu" {
  device      = netbox_device.test.id
  name        = "eth3"
  type        = "10gbase-t"
  mtu         = 9000
  speed       = 10000000 # 10 Gbps in Kbps
  description = "Jumbo frame interface"
}

# Management-only interface
resource "netbox_interface" "mgmt" {
  device      = netbox_device.test.id
  name        = "mgmt0"
  type        = "1000base-t"
  mgmt_only   = true
  description = "Out-of-band management interface"
}

# Virtual interface
resource "netbox_interface" "virtual" {
  device      = netbox_device.test.id
  name        = "vlan100"
  type        = "virtual"
  description = "VLAN 100 virtual interface"
}

# LAG (Link Aggregation Group) interface
resource "netbox_interface" "lag" {
  device      = netbox_device.test.id
  name        = "bond0"
  type        = "lag"
  description = "Link aggregation interface"
}

# Interface that's a member of the LAG
resource "netbox_interface" "lag_member" {
  device      = netbox_device.test.id
  name        = "eth4"
  type        = "1000base-t"
  lag         = netbox_interface.lag.id
  description = "LAG member interface"

  depends_on = [netbox_interface.lag]
}

# Interface with 802.1Q tagging mode
resource "netbox_interface" "tagged" {
  device      = netbox_device.test.id
  name        = "eth5"
  type        = "1000base-t"
  mode        = "tagged"
  description = "Trunk interface with VLAN tagging"
}

# Interface with mark_connected
resource "netbox_interface" "marked_connected" {
  device         = netbox_device.test.id
  name           = "eth6"
  type           = "1000base-t"
  mark_connected = true
  description    = "Interface marked as connected"
}

# Complete interface with all optional fields
resource "netbox_interface" "complete" {
  device         = netbox_device.test.id
  name           = "eth7"
  type           = "10gbase-x-sfpp"
  label          = "SFP+ Port 1"
  enabled        = true
  mtu            = 1500
  speed          = 10000000
  duplex         = "full"
  mgmt_only      = false
  description    = "Complete interface with all fields"
  mode           = "access"
  mark_connected = false
}
