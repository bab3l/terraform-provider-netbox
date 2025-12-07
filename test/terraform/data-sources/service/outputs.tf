# Service Data Source Outputs

output "by_id_name" {
  value = data.netbox_service.by_id.name
}

output "by_id_device" {
  value = data.netbox_service.by_id.device
}

output "by_id_protocol" {
  value = data.netbox_service.by_id.protocol
}

output "by_id_description" {
  value = data.netbox_service.by_id.description
}
