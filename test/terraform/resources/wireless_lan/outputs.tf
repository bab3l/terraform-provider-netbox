# Wireless LAN Outputs

# Basic wireless LAN outputs
output "basic_id" {
  value = netbox_wireless_lan.basic.id
}

output "basic_ssid" {
  value = netbox_wireless_lan.basic.ssid
}

output "basic_status" {
  value = netbox_wireless_lan.basic.status
}

# Complete wireless LAN outputs
output "complete_id" {
  value = netbox_wireless_lan.complete.id
}

output "complete_ssid" {
  value = netbox_wireless_lan.complete.ssid
}

output "complete_auth_type" {
  value = netbox_wireless_lan.complete.auth_type
}
