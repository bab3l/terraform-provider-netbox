# Look up a VLAN by ID
data "netbox_vlan" "by_id" {
  id = "1"
}

# Look up a VLAN by VID
data "netbox_vlan" "by_vid" {
  vid = 100
}

# Look up a VLAN by name
data "netbox_vlan" "by_name" {
  name = "test-vlan"
}

# Look up a VLAN by VID with optional name filter
data "netbox_vlan" "by_vid_and_name" {
  vid  = 100
  name = "test-vlan"
}

# Use VLAN data in outputs
output "by_id" {
  value = data.netbox_vlan.by_id.name
}

output "by_vid" {
  value = data.netbox_vlan.by_vid.name
}

output "vlan_info" {
  value = {
    id          = data.netbox_vlan.by_name.id
    vid         = data.netbox_vlan.by_name.vid
    name        = data.netbox_vlan.by_name.name
    site        = data.netbox_vlan.by_name.site
    group       = data.netbox_vlan.by_name.group
    status      = data.netbox_vlan.by_name.status
    role        = data.netbox_vlan.by_name.role
    description = data.netbox_vlan.by_name.description
  }
}
