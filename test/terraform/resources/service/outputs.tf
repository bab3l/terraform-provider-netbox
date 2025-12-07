# Service Outputs

# Basic service outputs
output "basic_id" {
  value = netbox_service.basic.id
}

output "basic_name" {
  value = netbox_service.basic.name
}

output "basic_protocol" {
  value = netbox_service.basic.protocol
}

# Complete service outputs
output "complete_id" {
  value = netbox_service.complete.id
}

output "complete_name" {
  value = netbox_service.complete.name
}

output "complete_ports" {
  value = netbox_service.complete.ports
}
