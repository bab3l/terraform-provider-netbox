# Look up a virtual disk by ID
data "netbox_virtual_disk" "by_id" {
  id = "1"
}

# Look up a virtual disk by name (requires virtual_machine)
data "netbox_virtual_disk" "by_name" {
  name            = "disk0"
  virtual_machine = netbox_virtual_machine.example.id
}

# Use virtual disk data in outputs
output "disk_id" {
  value = data.netbox_virtual_disk.by_name.id
}

output "disk_name" {
  value = data.netbox_virtual_disk.by_name.name
}

output "disk_size" {
  value = data.netbox_virtual_disk.by_name.size
}

output "disk_virtual_machine" {
  value = data.netbox_virtual_disk.by_name.virtual_machine_name
}

output "disk_description" {
  value = data.netbox_virtual_disk.by_name.description
}

# Access all custom fields
output "disk_custom_fields" {
  value       = data.netbox_virtual_disk.by_id.custom_fields
  description = "All custom fields defined in NetBox for this virtual disk"
}

# Access specific custom field by name
output "disk_storage_type" {
  value       = try([for cf in data.netbox_virtual_disk.by_id.custom_fields : cf.value if cf.name == "storage_type"][0], null)
  description = "Example: accessing a select custom field for storage type"
}

output "disk_iops_limit" {
  value       = try([for cf in data.netbox_virtual_disk.by_id.custom_fields : cf.value if cf.name == "iops_limit"][0], null)
  description = "Example: accessing a numeric custom field for IOPS limit"
}
