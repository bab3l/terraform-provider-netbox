# Cable Data Source Test Outputs

output "by_id_id" {
  value       = data.netbox_cable.by_id.id
  description = "ID from lookup by ID"
}

output "id_match" {
  value       = data.netbox_cable.by_id.id == netbox_cable.test.id
  description = "Whether ID lookup returns correct ID"
}

output "by_id_type" {
  value       = data.netbox_cable.by_id.type
  description = "Type from lookup by ID"
}

output "type_valid" {
  value       = data.netbox_cable.by_id.type == "cat6"
  description = "Whether type is correctly returned"
}

output "by_id_status" {
  value       = data.netbox_cable.by_id.status
  description = "Status from lookup by ID"
}

output "status_valid" {
  value       = data.netbox_cable.by_id.status == "connected"
  description = "Whether status is correctly returned"
}

output "by_id_label" {
  value       = data.netbox_cable.by_id.label
  description = "Label from lookup by ID"
}

output "label_valid" {
  value       = data.netbox_cable.by_id.label == "DS-TEST-CABLE"
  description = "Whether label is correctly returned"
}

output "by_id_color" {
  value       = data.netbox_cable.by_id.color
  description = "Color from lookup by ID"
}

output "color_valid" {
  value       = data.netbox_cable.by_id.color == "ff0000"
  description = "Whether color is correctly returned"
}

output "by_id_length" {
  value       = data.netbox_cable.by_id.length
  description = "Length from lookup by ID"
}

output "length_valid" {
  value       = data.netbox_cable.by_id.length == 10
  description = "Whether length is correctly returned"
}

output "by_id_length_unit" {
  value       = data.netbox_cable.by_id.length_unit
  description = "Length unit from lookup by ID"
}

output "length_unit_valid" {
  value       = data.netbox_cable.by_id.length_unit == "m"
  description = "Whether length unit is correctly returned"
}

output "by_id_description" {
  value       = data.netbox_cable.by_id.description
  description = "Description from lookup by ID"
}

output "description_valid" {
  value       = data.netbox_cable.by_id.description == "Cable for data source testing"
  description = "Whether description is correctly returned"
}
