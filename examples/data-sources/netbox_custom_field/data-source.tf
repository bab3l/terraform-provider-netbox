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
