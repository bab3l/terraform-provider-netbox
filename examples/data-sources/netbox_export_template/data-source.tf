# Look up an export template by ID
data "netbox_export_template" "by_id" {
  id = "1"
}

# Look up an export template by name
data "netbox_export_template" "by_name" {
  name = "Device Inventory"
}

# Individual attribute outputs
output "export_template_id" {
  value       = data.netbox_export_template.by_id.id
  description = "The unique ID of the export template"
}

output "export_template_name" {
  value       = data.netbox_export_template.by_name.name
  description = "The name of the export template"
}

output "export_template_description" {
  value       = data.netbox_export_template.by_name.description
  description = "Description of the export template"
}

output "export_template_content_type" {
  value       = data.netbox_export_template.by_name.content_type
  description = "The content type this template exports (e.g., devices, sites)"
}

output "export_template_object_types" {
  value       = data.netbox_export_template.by_name.object_types
  description = "List of object types this template can export"
}

output "export_template_template_code" {
  value       = data.netbox_export_template.by_name.template_code
  description = "The Jinja2 template code for export formatting"
}

# Note: Export templates do not support custom fields in NetBox API
output "export_template_note" {
  value       = "Export templates are read-only template objects"
  description = "Export templates provide Jinja2 templates for exporting object data in various formats"
}
