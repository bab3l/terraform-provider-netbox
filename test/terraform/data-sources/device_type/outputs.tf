# Device Type Data Source Outputs

output "by_id_model" {
  value = data.netbox_device_type.by_id.model
}

output "by_id_manufacturer" {
  value = data.netbox_device_type.by_id.manufacturer
}

output "by_slug_id" {
  value = data.netbox_device_type.by_slug.id
}

output "by_slug_description" {
  value = data.netbox_device_type.by_slug.description
}
