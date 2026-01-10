# Look up a wireless LAN group by ID
data "netbox_wireless_lan_group" "by_id" {
  id = "1"
}

# Look up a wireless LAN group by slug
data "netbox_wireless_lan_group" "by_slug" {
  slug = "test-wireless-lan-group"
}

# Look up a wireless LAN group by name
data "netbox_wireless_lan_group" "by_name" {
  name = "test-wireless-lan-group"
}

# Use wireless LAN group data in outputs
output "by_id" {
  value = data.netbox_wireless_lan_group.by_id.name
}

output "by_slug" {
  value = data.netbox_wireless_lan_group.by_slug.id
}

output "wlan_group_id" {
  value = data.netbox_wireless_lan_group.by_name.id
}

output "wlan_group_name" {
  value = data.netbox_wireless_lan_group.by_name.name
}

output "wlan_group_slug" {
  value = data.netbox_wireless_lan_group.by_name.slug
}

output "wlan_group_description" {
  value = data.netbox_wireless_lan_group.by_name.description
}

output "wlan_group_parent_id" {
  value = data.netbox_wireless_lan_group.by_name.parent_id
}

output "wlan_group_parent_name" {
  value = data.netbox_wireless_lan_group.by_name.parent_name
}

# Access all custom fields
output "wlan_group_custom_fields" {
  value       = data.netbox_wireless_lan_group.by_id.custom_fields
  description = "All custom fields defined in NetBox for this wireless LAN group"
}

# Access specific custom field by name
output "wlan_group_region" {
  value       = try([for cf in data.netbox_wireless_lan_group.by_id.custom_fields : cf.value if cf.name == "region_name"][0], null)
  description = "Example: accessing a text custom field for region"
}

output "wlan_group_ap_count" {
  value       = try([for cf in data.netbox_wireless_lan_group.by_id.custom_fields : cf.value if cf.name == "ap_count"][0], null)
  description = "Example: accessing a numeric custom field for access point count"
}
