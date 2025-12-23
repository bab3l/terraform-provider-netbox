# Example: Look up an FHRP group by ID
data "netbox_fhrp_group" "by_id" {
  id = "1"
}

# Example: Look up an FHRP group by protocol and group_id
data "netbox_fhrp_group" "by_protocol_and_group" {
  protocol = "vrrp3"
  group_id = 10
}

# Example: Use FHRP group data in other resources
output "fhrp_group_id" {
  value = data.netbox_fhrp_group.by_id.id
}

output "fhrp_group_protocol" {
  value = data.netbox_fhrp_group.by_protocol_and_group.protocol
}

output "fhrp_group_description" {
  value = data.netbox_fhrp_group.by_protocol_and_group.description
}
