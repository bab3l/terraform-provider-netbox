# Look up a module bay template by ID
data "netbox_module_bay_template" "by_id" {
  id = "1"
}

# Individual attribute outputs
output "module_bay_template_id" {
  value       = data.netbox_module_bay_template.by_id.id
  description = "The unique ID of the module bay template"
}

output "module_bay_template_name" {
  value       = data.netbox_module_bay_template.by_id.name
  description = "Module bay template name"
}

output "module_bay_template_device_type" {
  value       = data.netbox_module_bay_template.by_id.device_type
  description = "Device type that this module bay template belongs to"
}

output "module_bay_template_module_type" {
  value       = data.netbox_module_bay_template.by_id.module_type
  description = "Module type that this module bay template accepts"
}

output "module_bay_template_position" {
  value       = data.netbox_module_bay_template.by_id.position
  description = "The position or label of this module bay"
}

# Note: Module bay templates do not support custom fields in NetBox API
output "module_bay_template_note" {
  value       = "Module bay templates are read-only device design metadata"
  description = "Template datasources do not support custom fields"
}
