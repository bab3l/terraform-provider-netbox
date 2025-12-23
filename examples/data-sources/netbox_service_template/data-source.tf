# Look up service template by ID
data "netbox_service_template" "by_id" {
  id = "1"
}

# Look up service template by name
data "netbox_service_template" "by_name" {
  name = "HTTP"
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
