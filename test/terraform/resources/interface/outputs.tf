# Interface Resource Test Outputs

output "basic_id" {
  value       = netbox_interface.basic.id
  description = "ID of the basic interface"
}

output "basic_name" {
  value       = netbox_interface.basic.name
  description = "Name of the basic interface"
}

output "basic_type" {
  value       = netbox_interface.basic.type
  description = "Type of the basic interface"
}

output "basic_enabled" {
  value       = netbox_interface.basic.enabled
  description = "Enabled status of the basic interface (should default to true)"
}

output "with_label_id" {
  value       = netbox_interface.with_label.id
  description = "ID of the interface with label"
}

output "with_label_label" {
  value       = netbox_interface.with_label.label
  description = "Label of the interface"
}

output "disabled_enabled" {
  value       = netbox_interface.disabled.enabled
  description = "Enabled status of disabled interface (should be false)"
}

output "with_mtu_mtu" {
  value       = netbox_interface.with_mtu.mtu
  description = "MTU of the interface"
}

output "with_mtu_speed" {
  value       = netbox_interface.with_mtu.speed
  description = "Speed of the interface in Kbps"
}

output "mgmt_mgmt_only" {
  value       = netbox_interface.mgmt.mgmt_only
  description = "Management-only status"
}

output "virtual_type" {
  value       = netbox_interface.virtual.type
  description = "Type of virtual interface"
}

output "lag_id" {
  value       = netbox_interface.lag.id
  description = "ID of the LAG interface"
}

output "lag_member_lag" {
  value       = netbox_interface.lag_member.lag
  description = "LAG ID that the member belongs to"
}

output "tagged_mode" {
  value       = netbox_interface.tagged.mode
  description = "802.1Q mode of the interface"
}

output "marked_connected_mark_connected" {
  value       = netbox_interface.marked_connected.mark_connected
  description = "Mark connected status"
}

output "complete_id" {
  value       = netbox_interface.complete.id
  description = "ID of the complete interface"
}

output "complete_duplex" {
  value       = netbox_interface.complete.duplex
  description = "Duplex setting of the complete interface"
}
