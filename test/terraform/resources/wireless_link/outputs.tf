# Wireless Link Resource Test Outputs

output "basic_wireless_link_id" {
  description = "ID of the basic wireless link"
  value       = netbox_wireless_link.basic.id
}

output "basic_wireless_link_interface_a" {
  description = "Interface A of the basic wireless link"
  value       = netbox_wireless_link.basic.interface_a
}

output "basic_wireless_link_interface_b" {
  description = "Interface B of the basic wireless link"
  value       = netbox_wireless_link.basic.interface_b
}

output "with_ssid_wireless_link_id" {
  description = "ID of the wireless link with SSID"
  value       = netbox_wireless_link.with_ssid.id
}

output "with_ssid_wireless_link_ssid" {
  description = "SSID of the wireless link"
  value       = netbox_wireless_link.with_ssid.ssid
}

output "with_ssid_wireless_link_status" {
  description = "Status of the wireless link"
  value       = netbox_wireless_link.with_ssid.status
}

output "complete_wireless_link_id" {
  description = "ID of the complete wireless link"
  value       = netbox_wireless_link.complete.id
}

output "complete_wireless_link_ssid" {
  description = "SSID of the complete wireless link"
  value       = netbox_wireless_link.complete.ssid
}

output "complete_wireless_link_status" {
  description = "Status of the complete wireless link"
  value       = netbox_wireless_link.complete.status
}

output "complete_wireless_link_tenant" {
  description = "Tenant of the complete wireless link"
  value       = netbox_wireless_link.complete.tenant
}

output "complete_wireless_link_auth_type" {
  description = "Auth type of the complete wireless link"
  value       = netbox_wireless_link.complete.auth_type
}

output "complete_wireless_link_auth_cipher" {
  description = "Auth cipher of the complete wireless link"
  value       = netbox_wireless_link.complete.auth_cipher
}

output "complete_wireless_link_distance" {
  description = "Distance of the complete wireless link"
  value       = netbox_wireless_link.complete.distance
}

output "complete_wireless_link_distance_unit" {
  description = "Distance unit of the complete wireless link"
  value       = netbox_wireless_link.complete.distance_unit
}
