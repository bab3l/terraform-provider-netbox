# Example 1: Lookup by ID
data "netbox_module_bay" "by_id" {
  id = 1
}

output "module_bay_by_id" {
  value       = data.netbox_module_bay.by_id.display_name
  description = "Module bay display name when looked up by ID"
}

# Example 2: Lookup by device_id and name
data "netbox_module_bay" "by_device_and_name" {
  device_id = 5
  name      = "Bay-1"
}

output "module_bay_by_device_and_name" {
  value       = data.netbox_module_bay.by_device_and_name.label
  description = "Module bay label when looked up by device and name"
}
