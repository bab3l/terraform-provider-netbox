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

output "fhrp_group_auth_type" {
  value = data.netbox_fhrp_group.by_id.auth_type
}

output "fhrp_group_auth_key" {
  value     = data.netbox_fhrp_group.by_id.auth_key
  sensitive = true
}

# Access all custom fields
output "fhrp_group_custom_fields" {
  value       = data.netbox_fhrp_group.by_id.custom_fields
  description = "All custom fields defined in NetBox for this FHRP group"
}

# Access specific custom field by name
output "fhrp_group_vip" {
  value       = try([for cf in data.netbox_fhrp_group.by_id.custom_fields : cf.value if cf.name == "virtual_ip"][0], null)
  description = "Example: accessing a text custom field for virtual IP"
}

output "fhrp_group_preempt" {
  value       = try([for cf in data.netbox_fhrp_group.by_id.custom_fields : cf.value if cf.name == "preempt_enabled"][0], null)
  description = "Example: accessing a boolean custom field for preempt status"
}
