# Example: Look up config context by ID
data "netbox_config_context" "by_id" {
  id = "1"
}

# Example: Look up config context by name
data "netbox_config_context" "by_name" {
  name = "default-config"
}

# Example: Use config context data in other resources
output "config_context_id" {
  value = data.netbox_config_context.by_id.id
}

output "config_context_name" {
  value = data.netbox_config_context.by_name.name
}

output "config_context_data" {
  value = data.netbox_config_context.by_name.data
}

output "dns_servers" {
  value       = try(jsondecode(data.netbox_config_context.by_name.data).dns_servers, null)
  description = "Example: Extract DNS servers from config context JSON data"
}

output "config_context_weight" {
  value = data.netbox_config_context.by_name.weight
}

output "config_context_is_active" {
  value = data.netbox_config_context.by_name.is_active
}

# Note: Config contexts do not support custom fields in NetBox API
output "config_context_note" {
  value       = "Config contexts are read-only configuration objects"
  description = "Config contexts provide JSON configuration data to objects"
}
