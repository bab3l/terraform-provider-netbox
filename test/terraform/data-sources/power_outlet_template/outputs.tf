# Power Outlet Template Data Source Outputs

output "by_id_name" {
  value = data.netbox_power_outlet_template.by_id.name
}

output "by_id_device_type" {
  value = data.netbox_power_outlet_template.by_id.device_type
}

output "by_id_description" {
  value = data.netbox_power_outlet_template.by_id.description
}
