# Outputs for Device Bay Template Resource Test

output "basic_device_bay_template_id" {
  value       = netbox_device_bay_template.basic.id
  description = "ID of the basic device bay template"
}

output "basic_device_bay_template_name" {
  value       = netbox_device_bay_template.basic.name
  description = "Name of the basic device bay template"
}

output "full_device_bay_template_id" {
  value       = netbox_device_bay_template.full.id
  description = "ID of the full device bay template"
}

output "full_device_bay_template_label" {
  value       = netbox_device_bay_template.full.label
  description = "Label of the full device bay template"
}

output "full_device_bay_template_description" {
  value       = netbox_device_bay_template.full.description
  description = "Description of the full device bay template"
}

output "multi_device_bay_template_ids" {
  value       = netbox_device_bay_template.multi[*].id
  description = "IDs of the multiple device bay templates"
}
