# Look up an interface template by ID
data "netbox_interface_template" "by_id" {
  id = 1
}

# Look up an interface template by name and device type
data "netbox_interface_template" "by_device_type" {
  name        = "GigabitEthernet"
  device_type = 5
}

# Look up an interface template by name and module type
data "netbox_interface_template" "by_module_type" {
  name        = "SFP"
  module_type = 10
}

# Individual attribute outputs
output "interface_template_id" {
  value       = data.netbox_interface_template.by_id.id
  description = "The unique ID of the interface template"
}

output "interface_template_name" {
  value       = data.netbox_interface_template.by_device_type.name
  description = "The name of the interface template"
}

output "interface_template_type" {
  value       = data.netbox_interface_template.by_device_type.type
  description = "The interface type (e.g., 1000base-t, sfp-plus)"
}

output "interface_template_mtu" {
  value       = data.netbox_interface_template.by_device_type.mtu
  description = "Maximum transmission unit (MTU) for this interface"
}

output "interface_template_enabled" {
  value       = data.netbox_interface_template.by_device_type.enabled
  description = "Whether this interface is enabled by default"
}

output "interface_template_device_type" {
  value       = data.netbox_interface_template.by_device_type.device_type
  description = "The device type this template belongs to"
}

# Note: Interface templates do not support custom fields in NetBox API
output "interface_template_note" {
  value       = "Interface templates are read-only device design metadata"
  description = "Template datasources do not support custom fields"
}
