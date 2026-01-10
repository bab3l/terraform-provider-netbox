# Look up a config template by ID
data "netbox_config_template" "by_id" {
  id = "123"
}

# Look up a config template by name
data "netbox_config_template" "by_name" {
  name = "device-config"
}

# Individual attribute outputs
output "config_template_id" {
  value       = data.netbox_config_template.by_id.id
  description = "The unique ID of the config template"
}

output "config_template_name" {
  value       = data.netbox_config_template.by_name.name
  description = "The name of the config template"
}

output "config_template_description" {
  value       = data.netbox_config_template.by_name.description
  description = "Description of the config template"
}

output "config_template_environment_vars" {
  value       = data.netbox_config_template.by_name.environment_vars
  description = "Environment variables available in this template"
}

output "config_template_template_code" {
  value       = data.netbox_config_template.by_name.template_code
  description = "The Jinja2 template code for device configuration"
}

# Note: Config templates do not support custom fields in NetBox API
output "config_template_note" {
  value       = "Config templates are read-only template objects"
  description = "Config templates provide Jinja2 template configurations for device provisioning"
}
