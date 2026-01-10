# Look up a VM interface by ID
data "netbox_vm_interface" "by_id" {
  id = "1"
}

# Look up a VM interface by name and virtual machine
data "netbox_vm_interface" "by_name" {
  name            = "eth0"
  virtual_machine = "test-vm"
}

# Use VM interface data in outputs
output "by_id" {
  value = data.netbox_vm_interface.by_id.name
}

output "by_name" {
  value = data.netbox_vm_interface.by_name.id
}

output "interface_name" {
  value = data.netbox_vm_interface.by_name.name
}

output "interface_enabled" {
  value = data.netbox_vm_interface.by_name.enabled
}

output "interface_mtu" {
  value = data.netbox_vm_interface.by_name.mtu
}

output "interface_mac_address" {
  value = data.netbox_vm_interface.by_name.mac_address
}

output "interface_mode" {
  value = data.netbox_vm_interface.by_name.mode
}

output "interface_untagged_vlan" {
  value = data.netbox_vm_interface.by_name.untagged_vlan
}

# Access all custom fields
output "interface_custom_fields" {
  value       = data.netbox_vm_interface.by_id.custom_fields
  description = "All custom fields defined in NetBox for this VM interface"
}

# Access specific custom field by name
output "interface_vlan_purpose" {
  value       = try([for cf in data.netbox_vm_interface.by_id.custom_fields : cf.value if cf.name == "vlan_purpose"][0], null)
  description = "Example: accessing a text custom field for VLAN purpose"
}

output "interface_is_management" {
  value       = try([for cf in data.netbox_vm_interface.by_id.custom_fields : cf.value if cf.name == "is_management"][0], null)
  description = "Example: accessing a boolean custom field for management interface status"
}
