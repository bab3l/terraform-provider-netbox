# Device Bay Data Source Outputs

output "by_id_name" {
  value = data.netbox_device_bay.by_id.name
}

output "by_id_device" {
  value = data.netbox_device_bay.by_id.device
}

output "by_id_description" {
  value = data.netbox_device_bay.by_id.description
}
