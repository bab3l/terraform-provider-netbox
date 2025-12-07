# Power Panel Outputs

# Basic power panel outputs
output "basic_id" {
  value = netbox_power_panel.basic.id
}

output "basic_name" {
  value = netbox_power_panel.basic.name
}

# Complete power panel outputs
output "complete_id" {
  value = netbox_power_panel.complete.id
}

output "complete_name" {
  value = netbox_power_panel.complete.name
}

output "complete_description" {
  value = netbox_power_panel.complete.description
}
