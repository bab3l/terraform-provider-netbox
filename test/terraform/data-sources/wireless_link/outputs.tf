# Wireless Link Data Source Test Outputs

output "wireless_link_id" {
  description = "ID of the wireless link found by ID lookup"
  value       = data.netbox_wireless_link.by_id.id
}

output "wireless_link_interface_a" {
  description = "Interface A of the wireless link"
  value       = data.netbox_wireless_link.by_id.interface_a
}

output "wireless_link_interface_b" {
  description = "Interface B of the wireless link"
  value       = data.netbox_wireless_link.by_id.interface_b
}

output "wireless_link_ssid" {
  description = "SSID of the wireless link"
  value       = data.netbox_wireless_link.by_id.ssid
}

output "wireless_link_status" {
  description = "Status of the wireless link"
  value       = data.netbox_wireless_link.by_id.status
}

output "wireless_link_description" {
  description = "Description of the wireless link"
  value       = data.netbox_wireless_link.by_id.description
}

output "wireless_link_comments" {
  description = "Comments of the wireless link"
  value       = data.netbox_wireless_link.by_id.comments
}
