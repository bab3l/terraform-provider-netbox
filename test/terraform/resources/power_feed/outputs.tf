# Outputs to verify resource creation

output "basic_id" {
  description = "ID of the basic power feed"
  value       = netbox_power_feed.basic.id
}

output "basic_name" {
  description = "Name of the basic power feed"
  value       = netbox_power_feed.basic.name
}

output "basic_id_valid" {
  description = "Basic power feed has valid ID"
  value       = netbox_power_feed.basic.id != ""
}

output "complete_id" {
  description = "ID of the complete power feed"
  value       = netbox_power_feed.complete.id
}

output "complete_name" {
  description = "Name of the complete power feed"
  value       = netbox_power_feed.complete.name
}

output "complete_status" {
  description = "Status of the complete power feed"
  value       = netbox_power_feed.complete.status
}

output "complete_voltage" {
  description = "Voltage of the complete power feed"
  value       = netbox_power_feed.complete.voltage
}

output "power_panel_id" {
  description = "ID of the power panel"
  value       = netbox_power_panel.test.id
}
