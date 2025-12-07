# Power Port Template Outputs

# Basic power port template outputs
output "basic_id" {
  value = netbox_power_port_template.basic.id
}

output "basic_name" {
  value = netbox_power_port_template.basic.name
}

# Complete power port template outputs
output "complete_id" {
  value = netbox_power_port_template.complete.id
}

output "complete_name" {
  value = netbox_power_port_template.complete.name
}

output "complete_type" {
  value = netbox_power_port_template.complete.type
}
