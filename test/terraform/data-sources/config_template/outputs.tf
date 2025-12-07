# Config Template Data Source Outputs

output "by_id_name" {
  value = data.netbox_config_template.by_id.name
}

output "by_id_template_code" {
  value = data.netbox_config_template.by_id.template_code
}

output "by_name_id" {
  value = data.netbox_config_template.by_name.id
}

output "by_name_description" {
  value = data.netbox_config_template.by_name.description
}
