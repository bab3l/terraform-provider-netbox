# Outputs to verify resource creation

output "basic_id" {
  description = "ID of the basic module"
  value       = netbox_module.basic.id
}

output "basic_id_valid" {
  description = "Basic module has valid ID"
  value       = netbox_module.basic.id != ""
}

output "device_id" {
  description = "ID of the parent device"
  value       = netbox_device.test.id
}

output "module_bay_id" {
  description = "ID of the module bay"
  value       = netbox_module_bay.test.id
}

output "module_type_id" {
  description = "ID of the module type"
  value       = netbox_module_type.test.id
}
