# Outputs for inventory item template data source test

output "id_matches" {
  description = "Data source ID matches resource ID"
  value       = tostring(data.netbox_inventory_item_template.by_id.id) == netbox_inventory_item_template.test.id
}

output "name_matches" {
  description = "Data source name matches resource name"
  value       = data.netbox_inventory_item_template.by_id.name == netbox_inventory_item_template.test.name
}

output "label_matches" {
  description = "Data source label matches resource label"
  value       = data.netbox_inventory_item_template.by_id.label == netbox_inventory_item_template.test.label
}

output "part_id_matches" {
  description = "Data source part_id matches resource part_id"
  value       = data.netbox_inventory_item_template.by_id.part_id == netbox_inventory_item_template.test.part_id
}

output "description_matches" {
  description = "Data source description matches resource description"
  value       = data.netbox_inventory_item_template.by_id.description == netbox_inventory_item_template.test.description
}

output "device_type_id" {
  description = "Device type ID from data source"
  value       = data.netbox_inventory_item_template.by_id.device_type_id
}

output "role_id" {
  description = "Role ID from data source"
  value       = data.netbox_inventory_item_template.by_id.role_id
}

output "manufacturer_id" {
  description = "Manufacturer ID from data source"
  value       = data.netbox_inventory_item_template.by_id.manufacturer_id
}
