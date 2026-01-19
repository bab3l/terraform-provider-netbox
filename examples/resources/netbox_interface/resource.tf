# Basic interface
resource "netbox_interface" "example" {
  device      = netbox_device.example.name
  name        = "eth0"
  type        = "1000base-t"
  description = "Main network interface"

  # Partial custom fields management
  # Only specified custom fields are managed, others preserved
  custom_fields = [
    {
      name  = "circuit_id"
      value = "CKT-001"
    },
    {
      name  = "vlan_purpose"
      value = "management"
    }
  ]

  tags = [
    "management-interface"
  ]
}

# Optional: seed owned custom fields during import
import {
  to = netbox_interface.example
  id = "123"

  identity = {
    custom_fields = [
      "circuit_id:text",
      "vlan_purpose:text",
    ]
  }
}

# Interface with full configuration
resource "netbox_interface" "complete" {
  device         = netbox_device.example.name
  name           = "eth1"
  type           = "10gbase-x-sfpp"
  label          = "SFP+ Port 1"
  enabled        = true
  mtu            = 9000
  speed          = 10000000
  duplex         = "full"
  mgmt_only      = false
  description    = "High-speed data interface"
  mode           = "tagged"
  mark_connected = true

  # Partial custom fields management
  custom_fields = [
    {
      name  = "transceiver_type"
      value = "SFP-10G-SR"
    },
    {
      name  = "uplink_provider"
      value = "internal"
    },
    {
      name  = "monitoring_enabled"
      value = "true"
    }
  ]

  tags = [
    "uplink",
    "high-speed"
  ]
}

# Optional: seed owned custom fields during import
import {
  to = netbox_interface.complete
  id = "124"

  identity = {
    custom_fields = [
      "transceiver_type:text",
      "uplink_provider:text",
      "monitoring_enabled:boolean",
    ]
  }
}

# Virtual interface
resource "netbox_interface" "vlan100" {
  device      = netbox_device.example.name
  name        = "vlan100"
  type        = "virtual"
  description = "VLAN 100 virtual interface"
}

# LAG interface
resource "netbox_interface" "bond0" {
  device      = netbox_device.example.name
  name        = "bond0"
  type        = "lag"
  description = "Link aggregation group"
}

# LAG member interface
resource "netbox_interface" "lag_member" {
  device = netbox_device.example.name
  name   = "eth2"
  type   = "1000base-t"
  lag    = netbox_interface.bond0.name
}

# Management-only interface
resource "netbox_interface" "mgmt" {
  device    = netbox_device.example.id
  name      = "mgmt0"
  type      = "1000base-t"
  mgmt_only = true
}
