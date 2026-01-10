# Look up a virtual device context by ID
data "netbox_virtual_device_context" "by_id" {
  id = "1"
}

# Use virtual device context data in outputs
output "vdc_id" {
  value = data.netbox_virtual_device_context.by_id.id
}

output "vdc_name" {
  value = data.netbox_virtual_device_context.by_id.name
}

output "vdc_device" {
  value = data.netbox_virtual_device_context.by_id.device
}

output "vdc_identifier" {
  value = data.netbox_virtual_device_context.by_id.identifier
}

output "vdc_status" {
  value = data.netbox_virtual_device_context.by_id.status
}

output "vdc_primary_ip4" {
  value = data.netbox_virtual_device_context.by_id.primary_ip4
}

output "vdc_primary_ip6" {
  value = data.netbox_virtual_device_context.by_id.primary_ip6
}

# Access all custom fields
output "vdc_custom_fields" {
  value       = data.netbox_virtual_device_context.by_id.custom_fields
  description = "All custom fields defined in NetBox for this virtual device context"
}

# Access specific custom field by name
output "vdc_tenant_name" {
  value       = try([for cf in data.netbox_virtual_device_context.by_id.custom_fields : cf.value if cf.name == "tenant_name"][0], null)
  description = "Example: accessing a text custom field for tenant name"
}

output "vdc_cpu_limit" {
  value       = try([for cf in data.netbox_virtual_device_context.by_id.custom_fields : cf.value if cf.name == "cpu_limit"][0], null)
  description = "Example: accessing a numeric custom field for CPU limit"
}
