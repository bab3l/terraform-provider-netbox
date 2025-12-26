# Example 1: Lookup by ID
data "netbox_power_outlet_template" "by_id" {
  id = 1
}

output "power_outlet_template_by_id" {
  value       = data.netbox_power_outlet_template.by_id.name
  description = "Power outlet template name when looked up by ID"
}

# Example 2: Lookup by name and device_type
data "netbox_power_outlet_template" "by_device_type_and_name" {
  device_type = 10
  name        = "PSU"
}

output "power_outlet_template_by_device_type" {
  value       = data.netbox_power_outlet_template.by_device_type_and_name.type
  description = "Power outlet template type"
}

# Example 3: Lookup by name and module_type
data "netbox_power_outlet_template" "by_module_type_and_name" {
  module_type = 5
  name        = "PSU"
}

output "power_outlet_template_by_module_type" {
  value       = data.netbox_power_outlet_template.by_module_type_and_name.label
  description = "Power outlet template label"
}
