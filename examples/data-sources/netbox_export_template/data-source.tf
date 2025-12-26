# Example: Look up an export template by ID
data "netbox_export_template" "by_id" {
  id = "1"
}

# Example: Look up an export template by name
data "netbox_export_template" "by_name" {
  name = "Device Inventory"
}

# Example: Use export template data in other resources
output "export_template_id" {
  value = data.netbox_export_template.by_id.id
}

output "export_template_name" {
  value = data.netbox_export_template.by_name.name
}

output "export_template_object_types" {
  value = data.netbox_export_template.by_name.object_types
}
