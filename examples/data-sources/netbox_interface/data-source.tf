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
