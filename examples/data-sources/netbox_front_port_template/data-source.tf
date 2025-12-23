# Example: Look up a front port template by ID
data "netbox_front_port_template" "by_id" {
  id = 1
}

# Example: Look up a front port template by name (with optional device_type)
data "netbox_front_port_template" "by_name" {
  name        = "GigabitEthernet"
  device_type = 5
}

# Example: Look up a front port template by name (with optional module_type)
data "netbox_front_port_template" "by_name_module" {
  name        = "SFP"
  module_type = 10
}

# Example: Use front port template data in other resources
output "template_id" {
  value = data.netbox_front_port_template.by_id.id
}

output "template_name" {
  value = data.netbox_front_port_template.by_name.name
}

output "template_type" {
  value = data.netbox_front_port_template.by_name.type
}
