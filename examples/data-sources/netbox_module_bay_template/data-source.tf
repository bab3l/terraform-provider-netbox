# Lookup module bay template by ID
data "netbox_module_bay_template" "example" {
  id = "1"
}

output "module_bay_template" {
  value       = data.netbox_module_bay_template.example.name
  description = "Module bay template name"
}

output "device_type" {
  value       = data.netbox_module_bay_template.example.device_type
  description = "Device type that this module bay template belongs to"
}

output "module_type" {
  value       = data.netbox_module_bay_template.example.module_type
  description = "Module type that this module bay template accepts"
}
