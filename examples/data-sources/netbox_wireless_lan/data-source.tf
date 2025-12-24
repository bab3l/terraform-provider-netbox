# Look up a wireless LAN by ID
data "netbox_wireless_lan" "by_id" {
  id = "1"
}

# Look up a wireless LAN by SSID
data "netbox_wireless_lan" "by_ssid" {
  ssid = "test-ssid"
}

# Look up a wireless LAN by group ID
data "netbox_wireless_lan" "by_group" {
  group_id = 1
}

# Use wireless LAN data in outputs
output "by_id" {
  value = data.netbox_wireless_lan.by_id.ssid
}

output "by_ssid" {
  value = data.netbox_wireless_lan.by_ssid.id
}

output "wlan_info" {
  value = {
    id          = data.netbox_wireless_lan.by_ssid.id
    ssid        = data.netbox_wireless_lan.by_ssid.ssid
    status      = data.netbox_wireless_lan.by_ssid.status
    group_name  = data.netbox_wireless_lan.by_ssid.group_name
    vlan_name   = data.netbox_wireless_lan.by_ssid.vlan_name
    auth_type   = data.netbox_wireless_lan.by_ssid.auth_type
    description = data.netbox_wireless_lan.by_ssid.description
  }
}
