# IP Range Data Source Outputs

output "by_id_start_address" {
  value = data.netbox_ip_range.by_id.start_address
}

output "by_id_end_address" {
  value = data.netbox_ip_range.by_id.end_address
}

output "by_id_status" {
  value = data.netbox_ip_range.by_id.status
}

output "by_id_description" {
  value = data.netbox_ip_range.by_id.description
}
