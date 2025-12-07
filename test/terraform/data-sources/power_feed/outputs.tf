# Power Feed Data Source Outputs

output "by_id_name" {
  value = data.netbox_power_feed.by_id.name
}

output "by_id_power_panel" {
  value = data.netbox_power_feed.by_id.power_panel
}

output "by_id_status" {
  value = data.netbox_power_feed.by_id.status
}

output "by_id_description" {
  value = data.netbox_power_feed.by_id.description
}
