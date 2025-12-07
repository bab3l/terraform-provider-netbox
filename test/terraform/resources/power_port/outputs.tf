# Power Port Outputs

# Basic power port outputs
output "basic_id" {
  value = netbox_power_port.basic.id
}

output "basic_name" {
  value = netbox_power_port.basic.name
}

# Complete power port outputs
output "complete_id" {
  value = netbox_power_port.complete.id
}

output "complete_name" {
  value = netbox_power_port.complete.name
}

output "complete_type" {
  value = netbox_power_port.complete.type
}

output "complete_maximum_draw" {
  value = netbox_power_port.complete.maximum_draw
}
