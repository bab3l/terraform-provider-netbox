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

output "wlan_id" {
  value = data.netbox_wireless_lan.by_ssid.id
}

output "wlan_ssid" {
  value = data.netbox_wireless_lan.by_ssid.ssid
}

output "wlan_status" {
  value = data.netbox_wireless_lan.by_ssid.status
}

output "wlan_group_name" {
  value = data.netbox_wireless_lan.by_ssid.group_name
}

output "wlan_vlan_name" {
  value = data.netbox_wireless_lan.by_ssid.vlan_name
}

output "wlan_auth_type" {
  value = data.netbox_wireless_lan.by_ssid.auth_type
}

output "wlan_description" {
  value = data.netbox_wireless_lan.by_ssid.description
}

# Access all custom fields
output "wlan_custom_fields" {
  value       = data.netbox_wireless_lan.by_id.custom_fields
  description = "All custom fields defined in NetBox for this wireless LAN"
}

# Access specific custom field by name
output "wlan_channel" {
  value       = try([for cf in data.netbox_wireless_lan.by_id.custom_fields : cf.value if cf.name == "channel"][0], null)
  description = "Example: accessing a numeric custom field for channel number"
}

output "wlan_is_guest" {
  value       = try([for cf in data.netbox_wireless_lan.by_id.custom_fields : cf.value if cf.name == "is_guest_network"][0], null)
  description = "Example: accessing a boolean custom field for guest network status"
}
