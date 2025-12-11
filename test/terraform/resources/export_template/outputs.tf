# Export Template Resource Outputs

# Basic export template outputs
output "basic_id" {
  value = netbox_export_template.basic.id
}

output "basic_name" {
  value = netbox_export_template.basic.name
}

output "basic_object_types" {
  value = netbox_export_template.basic.object_types
}

# Devices export template outputs
output "devices_id" {
  value = netbox_export_template.devices.id
}

output "devices_description" {
  value = netbox_export_template.devices.description
}

# Multi type export template outputs
output "multi_type_id" {
  value = netbox_export_template.multi_type.id
}

output "multi_type_object_types_count" {
  value = length(netbox_export_template.multi_type.object_types)
}

# Complete export template outputs
output "complete_id" {
  value = netbox_export_template.complete.id
}

output "complete_name" {
  value = netbox_export_template.complete.name
}

output "complete_mime_type" {
  value = netbox_export_template.complete.mime_type
}

output "complete_file_extension" {
  value = netbox_export_template.complete.file_extension
}

output "complete_as_attachment" {
  value = netbox_export_template.complete.as_attachment
}

# JSON export template outputs
output "json_id" {
  value = netbox_export_template.json.id
}

output "json_mime_type" {
  value = netbox_export_template.json.mime_type
}

output "json_as_attachment" {
  value = netbox_export_template.json.as_attachment
}

# Validation outputs
output "basic_id_valid" {
  value = netbox_export_template.basic.id != null && netbox_export_template.basic.id != ""
}

output "multi_type_count_valid" {
  value = length(netbox_export_template.multi_type.object_types) == 2
}

output "complete_csv_valid" {
  value = netbox_export_template.complete.mime_type == "text/csv" && netbox_export_template.complete.file_extension == "csv"
}

output "json_format_valid" {
  value = netbox_export_template.json.mime_type == "application/json" && netbox_export_template.json.file_extension == "json"
}
