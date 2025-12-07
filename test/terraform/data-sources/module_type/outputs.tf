# Module Type Data Source Outputs

output "by_id_model" {
  value = data.netbox_module_type.by_id.model
}

output "by_id_manufacturer" {
  value = data.netbox_module_type.by_id.manufacturer
}

output "by_id_description" {
  value = data.netbox_module_type.by_id.description
}
