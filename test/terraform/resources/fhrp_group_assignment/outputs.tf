# FHRP Group Assignment Resource Outputs

# Basic assignment outputs
output "basic_id" {
  value = netbox_fhrp_group_assignment.basic.id
}

output "basic_group_id" {
  value = netbox_fhrp_group_assignment.basic.group_id
}

output "basic_interface_type" {
  value = netbox_fhrp_group_assignment.basic.interface_type
}

output "basic_interface_id" {
  value = netbox_fhrp_group_assignment.basic.interface_id
}

output "basic_priority" {
  value = netbox_fhrp_group_assignment.basic.priority
}

# High priority assignment outputs
output "high_priority_id" {
  value = netbox_fhrp_group_assignment.high_priority.id
}

output "high_priority_priority" {
  value = netbox_fhrp_group_assignment.high_priority.priority
}

# Validation outputs
output "basic_id_valid" {
  value = netbox_fhrp_group_assignment.basic.id != null && netbox_fhrp_group_assignment.basic.id != ""
}

output "basic_group_matches" {
  value = netbox_fhrp_group_assignment.basic.group_id == tostring(netbox_fhrp_group.test_vrrp.id)
}

output "basic_interface_matches" {
  value = netbox_fhrp_group_assignment.basic.interface_id == netbox_interface.test.id
}

output "high_priority_valid" {
  value = netbox_fhrp_group_assignment.high_priority.priority == 255
}
