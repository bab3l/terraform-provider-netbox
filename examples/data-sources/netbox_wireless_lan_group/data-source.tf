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

output "wlan_group_info" {
  value = {
    id          = data.netbox_wireless_lan_group.by_name.id
    name        = data.netbox_wireless_lan_group.by_name.name
    slug        = data.netbox_wireless_lan_group.by_name.slug
    description = data.netbox_wireless_lan_group.by_name.description
    parent_id   = data.netbox_wireless_lan_group.by_name.parent_id
    parent_name = data.netbox_wireless_lan_group.by_name.parent_name
  }
}
