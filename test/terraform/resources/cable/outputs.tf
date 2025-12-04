# Cable Resource Test Outputs

output "basic_id" {
  value       = netbox_cable.basic.id
  description = "ID of the basic cable"
}

output "basic_id_valid" {
  value       = can(tonumber(netbox_cable.basic.id)) && tonumber(netbox_cable.basic.id) > 0
  description = "Whether the basic cable ID is a valid positive number"
}

output "basic_status" {
  value       = netbox_cable.basic.status
  description = "Status of the basic cable (should default to connected)"
}

output "basic_status_valid" {
  value       = netbox_cable.basic.status == "connected"
  description = "Whether basic cable status defaults to connected"
}

output "with_type_id" {
  value       = netbox_cable.with_type.id
  description = "ID of the cable with type"
}

output "with_type_type" {
  value       = netbox_cable.with_type.type
  description = "Type of the cable"
}

output "with_type_type_valid" {
  value       = netbox_cable.with_type.type == "cat6a"
  description = "Whether cable type is correctly set"
}

output "full_id" {
  value       = netbox_cable.full.id
  description = "ID of the full cable"
}

output "full_type" {
  value       = netbox_cable.full.type
  description = "Type of the full cable"
}

output "full_type_valid" {
  value       = netbox_cable.full.type == "cat6"
  description = "Whether full cable type is correct"
}

output "full_status" {
  value       = netbox_cable.full.status
  description = "Status of the full cable"
}

output "full_status_valid" {
  value       = netbox_cable.full.status == "connected"
  description = "Whether full cable status is correct"
}

output "full_label" {
  value       = netbox_cable.full.label
  description = "Label of the full cable"
}

output "full_label_valid" {
  value       = netbox_cable.full.label == "CABLE-001"
  description = "Whether full cable label is correct"
}

output "full_color" {
  value       = netbox_cable.full.color
  description = "Color of the full cable"
}

output "full_color_valid" {
  value       = netbox_cable.full.color == "0000ff"
  description = "Whether full cable color is correct"
}

output "full_length" {
  value       = netbox_cable.full.length
  description = "Length of the full cable"
}

output "full_length_valid" {
  value       = netbox_cable.full.length == 5.5
  description = "Whether full cable length is correct"
}

output "full_length_unit" {
  value       = netbox_cable.full.length_unit
  description = "Length unit of the full cable"
}

output "full_length_unit_valid" {
  value       = netbox_cable.full.length_unit == "m"
  description = "Whether full cable length unit is correct"
}

output "full_description" {
  value       = netbox_cable.full.description
  description = "Description of the full cable"
}

output "full_description_valid" {
  value       = netbox_cable.full.description == "Ethernet cable from Device A to Device B"
  description = "Whether full cable description is correct"
}
