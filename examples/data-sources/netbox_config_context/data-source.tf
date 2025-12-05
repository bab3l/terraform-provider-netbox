# Example: Look up config context by ID
data "netbox_config_context" "by_id" {
  id = "1"
}

# Example: Look up config context by name
data "netbox_config_context" "by_name" {
  name = "default-config"
}

# Example: Use config context data in other resources
output "dns_servers" {
  value = jsondecode(data.netbox_config_context.by_name.data).dns_servers
}

output "config_context_weight" {
  value = data.netbox_config_context.by_name.weight
}

output "config_context_is_active" {
  value = data.netbox_config_context.by_name.is_active
}
