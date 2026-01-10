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

output "vm_status" {
  value = data.netbox_virtual_machine.by_name.status
}

output "vm_site" {
  value = data.netbox_virtual_machine.by_name.site
}

output "vm_cluster" {
  value = data.netbox_virtual_machine.by_name.cluster
}

output "vm_role" {
  value = data.netbox_virtual_machine.by_name.role
}

output "vm_vcpus" {
  value = data.netbox_virtual_machine.by_name.vcpus
}

output "vm_memory" {
  value = data.netbox_virtual_machine.by_name.memory
}

output "vm_disk" {
  value = data.netbox_virtual_machine.by_name.disk
}

# Access all custom fields
output "vm_custom_fields" {
  value       = data.netbox_virtual_machine.by_id.custom_fields
  description = "All custom fields defined in NetBox for this virtual machine"
}

# Access specific custom field by name
output "vm_os_version" {
  value       = try([for cf in data.netbox_virtual_machine.by_id.custom_fields : cf.value if cf.name == "os_version"][0], null)
  description = "Example: accessing a text custom field for OS version"
}

output "vm_backup_enabled" {
  value       = try([for cf in data.netbox_virtual_machine.by_id.custom_fields : cf.value if cf.name == "backup_enabled"][0], null)
  description = "Example: accessing a boolean custom field for backup status"
}
