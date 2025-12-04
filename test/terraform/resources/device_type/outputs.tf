output "device_type_basic_id" {
  description = "ID of the basic device type"
  value       = netbox_device_type.basic.id
}

output "device_type_basic_model" {
  description = "Model name of the basic device type"
  value       = netbox_device_type.basic.model
}

output "device_type_complete_id" {
  description = "ID of the complete device type"
  value       = netbox_device_type.complete.id
}

output "device_type_complete_u_height" {
  description = "U height of the complete device type"
  value       = netbox_device_type.complete.u_height
}

output "manufacturer_id_valid" {
  description = "Validates manufacturer was created correctly"
  value       = netbox_manufacturer.test.id != ""
}

output "basic_device_type_valid" {
  description = "Validates basic device type was created correctly"
  value       = netbox_device_type.basic.id != "" && netbox_device_type.basic.slug == "basic-test-model"
}

output "complete_device_type_valid" {
  description = "Validates complete device type was created correctly"
  value       = netbox_device_type.complete.id != "" && netbox_device_type.complete.u_height == 2
}
