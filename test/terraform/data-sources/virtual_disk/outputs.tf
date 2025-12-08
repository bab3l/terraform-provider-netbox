# Virtual Disk Data Source Outputs

# Look up by ID outputs
output "by_id_id" {
  value = data.netbox_virtual_disk.by_id.id
}

output "by_id_name" {
  value = data.netbox_virtual_disk.by_id.name
}

output "by_id_size" {
  value = data.netbox_virtual_disk.by_id.size
}

output "by_id_description" {
  value = data.netbox_virtual_disk.by_id.description
}

output "by_id_virtual_machine_name" {
  value = data.netbox_virtual_disk.by_id.virtual_machine_name
}

# Look up by name outputs
output "by_name_id" {
  value = data.netbox_virtual_disk.by_name.id
}

output "by_name_name" {
  value = data.netbox_virtual_disk.by_name.name
}

output "by_name_size" {
  value = data.netbox_virtual_disk.by_name.size
}

# Validation outputs
output "all_ids_match" {
  description = "Validates that all lookups return the same ID"
  value = alltrue([
    data.netbox_virtual_disk.by_id.id == netbox_virtual_disk.test.id,
    data.netbox_virtual_disk.by_name.id == netbox_virtual_disk.test.id
  ])
}

output "by_id_name_valid" {
  description = "Validates that lookup by ID returns correct name"
  value       = data.netbox_virtual_disk.by_id.name == "ds-test-disk"
}

output "by_name_name_valid" {
  description = "Validates that lookup by name returns correct name"
  value       = data.netbox_virtual_disk.by_name.name == "ds-test-disk"
}

output "sizes_match" {
  description = "Validates that sizes match the created resource"
  value       = data.netbox_virtual_disk.by_id.size == "250"
}

output "descriptions_match" {
  description = "Validates that descriptions match the created resource"
  value       = data.netbox_virtual_disk.by_id.description == "Test disk for data source"
}
