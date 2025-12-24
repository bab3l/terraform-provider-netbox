# Look up a VLAN group by ID
data "netbox_vlan_group" "by_id" {
  id = "1"
}

# Look up a VLAN group by slug
data "netbox_vlan_group" "by_slug" {
  slug = "test-vlan-group"
}

# Look up a VLAN group by name
data "netbox_vlan_group" "by_name" {
  name = "test-vlan-group"
}

# Use VLAN group data in outputs
output "by_id" {
  value = data.netbox_vlan_group.by_id.name
}

output "by_slug" {
  value = data.netbox_vlan_group.by_slug.id
}

output "by_name" {
  value = data.netbox_vlan_group.by_name.slug
}

output "vlan_group_info" {
  value = {
    id          = data.netbox_vlan_group.by_name.id
    name        = data.netbox_vlan_group.by_name.name
    slug        = data.netbox_vlan_group.by_name.slug
    description = data.netbox_vlan_group.by_name.description
  }
}
