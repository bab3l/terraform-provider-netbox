# Virtual Disk Outputs

# Basic virtual disk outputs
output "basic_id" {
  value = netbox_virtual_disk.basic.id
}

output "basic_name" {
  value = netbox_virtual_disk.basic.name
}

output "basic_size" {
  value = netbox_virtual_disk.basic.size
}

output "basic_virtual_machine" {
  value = netbox_virtual_disk.basic.virtual_machine
}

output "basic_id_valid" {
  description = "Validates that the basic virtual disk was created with an ID"
  value       = netbox_virtual_disk.basic.id != null && netbox_virtual_disk.basic.id != ""
}

output "basic_name_valid" {
  description = "Validates that the basic virtual disk name matches the input"
  value       = netbox_virtual_disk.basic.name == "disk0"
}

output "basic_size_valid" {
  description = "Validates that the basic virtual disk size matches the input"
  value       = netbox_virtual_disk.basic.size == "100"
}

# Complete virtual disk outputs
output "complete_id" {
  value = netbox_virtual_disk.complete.id
}

output "complete_name" {
  value = netbox_virtual_disk.complete.name
}

output "complete_size" {
  value = netbox_virtual_disk.complete.size
}

output "complete_description" {
  value = netbox_virtual_disk.complete.description
}

output "complete_name_valid" {
  description = "Validates that the complete virtual disk name matches the input"
  value       = netbox_virtual_disk.complete.name == "disk1"
}

output "complete_size_valid" {
  description = "Validates that the size was set correctly"
  value       = netbox_virtual_disk.complete.size == "500"
}

output "complete_description_valid" {
  description = "Validates that the description was set correctly"
  value       = netbox_virtual_disk.complete.description == "Primary data disk"
}

# Aggregate validation output
output "all_tests_passed" {
  description = "Validates all virtual disk tests passed"
  value = alltrue([
    netbox_virtual_disk.basic.id != null && netbox_virtual_disk.basic.id != "",
    netbox_virtual_disk.basic.name == "disk0",
    netbox_virtual_disk.basic.size == "100",
    netbox_virtual_disk.complete.name == "disk1",
    netbox_virtual_disk.complete.size == "500",
    netbox_virtual_disk.complete.description == "Primary data disk"
  ])
}
