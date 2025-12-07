# Custom Field Data Source Outputs

output "by_id_name" {
  value = data.netbox_custom_field.by_id.name
}

output "by_id_type" {
  value = data.netbox_custom_field.by_id.type
}

output "by_name_id" {
  value = data.netbox_custom_field.by_name.id
}

output "by_name_description" {
  value = data.netbox_custom_field.by_name.description
}
