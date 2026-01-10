# Look up interface by ID
data "netbox_interface" "by_id" {
  id = "123"
}

# Look up interface by device name and interface name
data "netbox_interface" "by_name" {
  device = "my-server"
  name   = "eth0"
}

# Look up interface by device ID and interface name
data "netbox_interface" "by_device_id" {
  device = "42"
  name   = "GigabitEthernet0/0"
}

# Use interface data in other resources
output "interface_type" {
  value = data.netbox_interface.by_name.type
}

output "interface_mac" {
  value = data.netbox_interface.by_name.mac_address
}

output "interface_mtu" {
  value = data.netbox_interface.by_id.mtu
}

output "interface_enabled" {
  value = data.netbox_interface.by_id.enabled
}

# Access all custom fields
output "interface_custom_fields" {
  value       = data.netbox_interface.by_id.custom_fields
  description = "All custom fields defined in NetBox for this interface"
}

# Access specific custom fields by name
output "interface_vlan_mode" {
  value       = try([for cf in data.netbox_interface.by_id.custom_fields : cf.value if cf.name == "vlan_mode"][0], null)
  description = "Example: accessing a select custom field"
}

output "interface_uplink" {
  value       = try([for cf in data.netbox_interface.by_id.custom_fields : cf.value if cf.name == "is_uplink"][0], null)
  description = "Example: accessing a boolean custom field"
}
