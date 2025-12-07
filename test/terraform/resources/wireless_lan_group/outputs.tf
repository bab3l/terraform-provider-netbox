# Wireless LAN Group Outputs

# Basic wireless LAN group outputs
output "basic_id" {
  value = netbox_wireless_lan_group.basic.id
}

output "basic_name" {
  value = netbox_wireless_lan_group.basic.name
}

output "basic_slug" {
  value = netbox_wireless_lan_group.basic.slug
}

# Parent wireless LAN group outputs
output "parent_id" {
  value = netbox_wireless_lan_group.parent.id
}

output "parent_name" {
  value = netbox_wireless_lan_group.parent.name
}

# Child wireless LAN group outputs
output "child_id" {
  value = netbox_wireless_lan_group.child.id
}

output "child_name" {
  value = netbox_wireless_lan_group.child.name
}

output "child_parent" {
  value = netbox_wireless_lan_group.child.parent
}
