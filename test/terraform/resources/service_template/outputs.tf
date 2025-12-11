# Service Template Resource Outputs

# Basic service template outputs
output "basic_id" {
  value = netbox_service_template.basic.id
}

output "basic_name" {
  value = netbox_service_template.basic.name
}

output "basic_ports" {
  value = netbox_service_template.basic.ports
}

# With protocol outputs
output "with_protocol_id" {
  value = netbox_service_template.with_protocol.id
}

output "with_protocol_protocol" {
  value = netbox_service_template.with_protocol.protocol
}

# Multi port outputs
output "multi_port_id" {
  value = netbox_service_template.multi_port.id
}

output "multi_port_ports" {
  value = netbox_service_template.multi_port.ports
}

output "multi_port_ports_count" {
  value = length(netbox_service_template.multi_port.ports)
}

# Complete outputs
output "complete_id" {
  value = netbox_service_template.complete.id
}

output "complete_name" {
  value = netbox_service_template.complete.name
}

output "complete_protocol" {
  value = netbox_service_template.complete.protocol
}

output "complete_description" {
  value = netbox_service_template.complete.description
}

output "complete_comments" {
  value = netbox_service_template.complete.comments
}

# SCTP outputs
output "sctp_id" {
  value = netbox_service_template.sctp.id
}

output "sctp_protocol" {
  value = netbox_service_template.sctp.protocol
}

# Validation outputs
output "basic_id_valid" {
  value = netbox_service_template.basic.id != null && netbox_service_template.basic.id != ""
}

output "multi_port_count_valid" {
  value = length(netbox_service_template.multi_port.ports) == 3
}

output "complete_has_description" {
  value = netbox_service_template.complete.description != null && netbox_service_template.complete.description != ""
}
