# Power Panel Data Source Outputs

output "by_id_name" {
  value = data.netbox_power_panel.by_id.name
}

output "by_id_site" {
  value = data.netbox_power_panel.by_id.site
}

output "by_name_id" {
  value = data.netbox_power_panel.by_name.id
}

output "by_name_description" {
  value = data.netbox_power_panel.by_name.description
}
