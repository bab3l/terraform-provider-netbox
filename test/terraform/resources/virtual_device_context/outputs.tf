output "basic_id" {
  description = "ID of the basic virtual device context"
  value       = netbox_virtual_device_context.basic.id
}

output "basic_name" {
  description = "Name of the basic virtual device context"
  value       = netbox_virtual_device_context.basic.name
}

output "basic_status" {
  description = "Status of the basic virtual device context"
  value       = netbox_virtual_device_context.basic.status
}

output "basic_id_valid" {
  description = "Basic VDC has valid ID"
  value       = netbox_virtual_device_context.basic.id != ""
}

output "complete_id" {
  description = "ID of the complete virtual device context"
  value       = netbox_virtual_device_context.complete.id
}

output "complete_name" {
  description = "Name of the complete virtual device context"
  value       = netbox_virtual_device_context.complete.name
}

output "complete_description" {
  description = "Description of the complete virtual device context"
  value       = netbox_virtual_device_context.complete.description
}

output "device_id" {
  description = "ID of the parent device"
  value       = netbox_device.test.id
}

output "site_id" {
  description = "ID of the parent site"
  value       = netbox_site.test.id
}
