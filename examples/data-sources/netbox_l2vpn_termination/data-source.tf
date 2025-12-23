# Example: Look up an L2VPN termination by ID (only supported lookup method)
data "netbox_l2vpn_termination" "by_id" {
  id = "123"
}

# Example: Use L2VPN termination data in other resources
output "termination_id" {
  value = data.netbox_l2vpn_termination.by_id.id
}

output "termination_l2vpn" {
  value = data.netbox_l2vpn_termination.by_id.l2vpn
}

output "termination_object_type" {
  value = data.netbox_l2vpn_termination.by_id.assigned_object_type
}
