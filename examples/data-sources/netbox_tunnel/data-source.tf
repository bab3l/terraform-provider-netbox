# Look up tunnel by ID
data "netbox_tunnel" "by_id" {
  id = "1"
}

# Look up tunnel by name
data "netbox_tunnel" "by_name" {
  name = "example-tunnel"
}

output "tunnel_id" {
  value = data.netbox_tunnel.by_id.id
}

output "tunnel_name" {
  value = data.netbox_tunnel.by_name.name
}

output "tunnel_status" {
  value = data.netbox_tunnel.by_id.status
}

output "tunnel_encapsulation" {
  value = data.netbox_tunnel.by_name.encapsulation
}

output "tunnel_group" {
  value = data.netbox_tunnel.by_id.tunnel_group
}

output "tunnel_tenant" {
  value = data.netbox_tunnel.by_id.tenant
}

output "tunnel_description" {
  value = data.netbox_tunnel.by_id.description
}

# Access all custom fields
output "tunnel_custom_fields" {
  value       = data.netbox_tunnel.by_id.custom_fields
  description = "All custom fields defined in NetBox for this tunnel"
}

# Access specific custom field by name
output "tunnel_bandwidth" {
  value       = try([for cf in data.netbox_tunnel.by_id.custom_fields : cf.value if cf.name == "bandwidth_mbps"][0], null)
  description = "Example: accessing a numeric custom field for bandwidth"
}

output "tunnel_is_encrypted" {
  value       = try([for cf in data.netbox_tunnel.by_id.custom_fields : cf.value if cf.name == "is_encrypted"][0], null)
  description = "Example: accessing a boolean custom field for encryption status"
}
