# Look up service template by ID
data "netbox_service_template" "by_id" {
  id = "1"
}

# Look up service template by name
data "netbox_service_template" "by_name" {
  name = "HTTP"
}

output "template_id" {
  value = data.netbox_service_template.by_id.id
}

output "template_name" {
  value = data.netbox_service_template.by_id.name
}

output "template_by_id" {
  value = data.netbox_service_template.by_id.protocol
}

output "template_protocol" {
  value = data.netbox_service_template.by_name.protocol
}

output "template_ports" {
  value = data.netbox_service_template.by_name.ports
}

output "template_description" {
  value = data.netbox_service_template.by_name.description
}

# Note: Service templates do not support custom fields in NetBox API
output "service_template_note" {
  value       = "Service templates are read-only configuration objects"
  description = "Service templates define port mappings and protocols for services"
}
