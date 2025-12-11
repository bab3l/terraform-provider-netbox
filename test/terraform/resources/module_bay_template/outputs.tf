output "basic_id" {
  description = "ID of the basic module bay template"
  value       = netbox_module_bay_template.basic.id
}

output "basic_name" {
  description = "Name of the basic module bay template"
  value       = netbox_module_bay_template.basic.name
}

output "basic_id_valid" {
  description = "Basic module bay template has valid ID"
  value       = netbox_module_bay_template.basic.id != ""
}

output "complete_id" {
  description = "ID of the complete module bay template"
  value       = netbox_module_bay_template.complete.id
}

output "complete_name" {
  description = "Name of the complete module bay template"
  value       = netbox_module_bay_template.complete.name
}

output "complete_label" {
  description = "Label of the complete module bay template"
  value       = netbox_module_bay_template.complete.label
}

output "complete_position" {
  description = "Position of the complete module bay template"
  value       = netbox_module_bay_template.complete.position
}

output "complete_description" {
  description = "Description of the complete module bay template"
  value       = netbox_module_bay_template.complete.description
}

output "device_type_id" {
  description = "ID of the parent device type"
  value       = netbox_device_type.test.id
}

output "manufacturer_id" {
  description = "ID of the manufacturer"
  value       = netbox_manufacturer.test.id
}
