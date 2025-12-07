# Wireless LAN Data Source Outputs

output "by_id_ssid" {
  value = data.netbox_wireless_lan.by_id.ssid
}

output "by_id_status" {
  value = data.netbox_wireless_lan.by_id.status
}

output "by_id_description" {
  value = data.netbox_wireless_lan.by_id.description
}
