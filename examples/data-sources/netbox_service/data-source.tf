# Look up service by ID
data "netbox_service" "by_id" {
  id = "1"
}

# Look up service by name and device
data "netbox_service" "by_name_and_device" {
  name   = "http"
  device = "1"
}

# Look up service by name and virtual machine
data "netbox_service" "by_name_and_vm" {
  name            = "ssh"
  virtual_machine = "1"
}

# Use service data in other resources
output "service_name" {
  value = data.netbox_service.by_id.name
}

output "service_protocol" {
  value = data.netbox_service.by_name_and_device.protocol
}

output "service_ports" {
  value = data.netbox_service.by_name_and_device.ports
}

output "service_description" {
  value = data.netbox_service.by_id.description
}

output "service_ipaddresses" {
  value = data.netbox_service.by_id.ipaddresses
}

# Access all custom fields
output "service_custom_fields" {
  value       = data.netbox_service.by_id.custom_fields
  description = "All custom fields defined in NetBox for this service"
}

# Access specific custom fields by name
output "service_monitoring_enabled" {
  value       = try([for cf in data.netbox_service.by_id.custom_fields : cf.value if cf.name == "monitoring_enabled"][0], null)
  description = "Example: accessing a boolean custom field"
}

output "service_health_check_url" {
  value       = try([for cf in data.netbox_service.by_id.custom_fields : cf.value if cf.name == "health_check_url"][0], null)
  description = "Example: accessing a URL custom field"
}
