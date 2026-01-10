# Example: Look up a custom field by ID
data "netbox_custom_field" "by_id" {
  id = "1"
}

# Example: Look up a custom field by name
data "netbox_custom_field" "by_name" {
  name = "test_field"
}

# Example: Use custom field data in other resources
output "custom_field_id" {
  value = data.netbox_custom_field.by_id.id
}

output "custom_field_name" {
  value = data.netbox_custom_field.by_name.name
}

output "custom_field_type" {
  value = data.netbox_custom_field.by_name.type
}

output "custom_field_required" {
  value = data.netbox_custom_field.by_name.required
}

output "custom_field_object_types" {
  value = data.netbox_custom_field.by_name.object_types
}

# Note: Custom fields are metadata definitions and do not support custom fields in NetBox API
output "custom_field_note" {
  value       = "Custom field definitions are read-only metadata"
  description = "Custom fields define the structure of custom data for other objects"
}
