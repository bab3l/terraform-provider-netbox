# Example: Look up an FHRP group assignment by ID (only supported lookup method)
data "netbox_fhrp_group_assignment" "by_id" {
  id = "1"
}

# Example: Use FHRP group assignment data in other resources
output "assignment_id" {
  value = data.netbox_fhrp_group_assignment.by_id.id
}

output "assignment_group_id" {
  value = data.netbox_fhrp_group_assignment.by_id.group_id
}

output "assignment_priority" {
  value = data.netbox_fhrp_group_assignment.by_id.priority
}
