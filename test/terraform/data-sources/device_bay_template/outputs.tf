# Outputs for Device Bay Template Data Source Test

# By ID outputs
output "by_id_id" {
  value       = data.netbox_device_bay_template.by_id.id
  description = "ID from ID lookup"
}

output "by_id_name" {
  value       = data.netbox_device_bay_template.by_id.name
  description = "Name from ID lookup"
}

output "by_id_device_type" {
  value       = data.netbox_device_bay_template.by_id.device_type
  description = "Device type ID from ID lookup"
}

output "by_id_device_type_name" {
  value       = data.netbox_device_bay_template.by_id.device_type_name
  description = "Device type name from ID lookup"
}

output "by_id_label" {
  value       = data.netbox_device_bay_template.by_id.label
  description = "Label from ID lookup"
}

output "by_id_description" {
  value       = data.netbox_device_bay_template.by_id.description
  description = "Description from ID lookup"
}

# By name outputs
output "by_name_id" {
  value       = data.netbox_device_bay_template.by_name.id
  description = "ID from name lookup"
}

output "by_name_name" {
  value       = data.netbox_device_bay_template.by_name.name
  description = "Name from name lookup"
}
