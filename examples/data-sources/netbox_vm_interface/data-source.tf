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

output "interface_info" {
  value = {
    id            = data.netbox_vm_interface.by_name.id
    name          = data.netbox_vm_interface.by_name.name
    enabled       = data.netbox_vm_interface.by_name.enabled
    mtu           = data.netbox_vm_interface.by_name.mtu
    mac_address   = data.netbox_vm_interface.by_name.mac_address
    mode          = data.netbox_vm_interface.by_name.mode
    untagged_vlan = data.netbox_vm_interface.by_name.untagged_vlan
  }
}
