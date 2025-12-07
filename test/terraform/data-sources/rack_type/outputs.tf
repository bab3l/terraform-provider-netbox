# Rack Type Data Source Outputs

output "by_id_model" {
  value = data.netbox_rack_type.by_id.model
}

output "by_id_manufacturer" {
  value = data.netbox_rack_type.by_id.manufacturer
}

output "by_id_u_height" {
  value = data.netbox_rack_type.by_id.u_height
}

output "by_id_description" {
  value = data.netbox_rack_type.by_id.description
}
