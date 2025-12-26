# Look up a virtual machine by ID
data "netbox_virtual_machine" "by_id" {
  id = "1"
}

# Look up a virtual machine by name
data "netbox_virtual_machine" "by_name" {
  name = "test-vm"
}

# Use virtual machine data in outputs
output "by_id" {
  value = data.netbox_virtual_machine.by_id.name
}

output "by_name" {
  value = data.netbox_virtual_machine.by_name.id
}

output "vm_specs" {
  value = {
    status  = data.netbox_virtual_machine.by_name.status
    site    = data.netbox_virtual_machine.by_name.site
    cluster = data.netbox_virtual_machine.by_name.cluster
    role    = data.netbox_virtual_machine.by_name.role
    vcpus   = data.netbox_virtual_machine.by_name.vcpus
    memory  = data.netbox_virtual_machine.by_name.memory
    disk    = data.netbox_virtual_machine.by_name.disk
  }
}
