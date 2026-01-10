# Look up tunnel termination by ID
data "netbox_tunnel_termination" "by_id" {
  id = "1"
}

# Look up tunnel termination by tunnel ID
data "netbox_tunnel_termination" "by_tunnel" {
  tunnel = "1"
}

# Look up tunnel termination by tunnel name
data "netbox_tunnel_termination" "by_tunnel_name" {
  tunnel_name = "example-tunnel"
}

output "termination_id" {
  value = data.netbox_tunnel_termination.by_id.id
}

output "termination_tunnel" {
  value = data.netbox_tunnel_termination.by_tunnel.tunnel
}

output "termination_role" {
  value = data.netbox_tunnel_termination.by_id.role
}

output "termination_type" {
  value = data.netbox_tunnel_termination.by_tunnel.termination_type
}

output "termination_outside_ip" {
  value = data.netbox_tunnel_termination.by_tunnel_name.outside_ip
}

# Access all custom fields
output "termination_custom_fields" {
  value       = data.netbox_tunnel_termination.by_id.custom_fields
  description = "All custom fields defined in NetBox for this tunnel termination"
}

# Access specific custom field by name
output "termination_location" {
  value       = try([for cf in data.netbox_tunnel_termination.by_id.custom_fields : cf.value if cf.name == "location_name"][0], null)
  description = "Example: accessing a text custom field for location"
}

output "termination_is_primary" {
  value       = try([for cf in data.netbox_tunnel_termination.by_id.custom_fields : cf.value if cf.name == "is_primary"][0], null)
  description = "Example: accessing a boolean custom field for primary status"
}
