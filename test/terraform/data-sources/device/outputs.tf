# Device Data Source Outputs

output "by_id_name" {
  value = data.netbox_device.by_id.name
}

output "by_id_site" {
  value = data.netbox_device.by_id.site
}

output "by_name_id" {
  value = data.netbox_device.by_name.id
}

output "by_name_description" {
  value = data.netbox_device.by_name.description
}
