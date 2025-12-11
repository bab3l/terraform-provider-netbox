output "id_matches" {
  value = data.netbox_virtual_device_context.by_id.id == netbox_virtual_device_context.test.id
}

output "name_matches" {
  value = data.netbox_virtual_device_context.by_id.name == netbox_virtual_device_context.test.name
}

output "status_matches" {
  value = data.netbox_virtual_device_context.by_id.status == netbox_virtual_device_context.test.status
}

output "description_matches" {
  value = data.netbox_virtual_device_context.by_id.description == netbox_virtual_device_context.test.description
}
