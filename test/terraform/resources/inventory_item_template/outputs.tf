output "basic_id" {
  description = "ID of the basic inventory item template"
  value       = netbox_inventory_item_template.basic.id
}

output "basic_name" {
  description = "Name of the basic inventory item template"
  value       = netbox_inventory_item_template.basic.name
}

output "basic_id_valid" {
  description = "Basic inventory item template has valid ID"
  value       = netbox_inventory_item_template.basic.id != ""
}

output "parent_id" {
  description = "ID of the parent inventory item template"
  value       = netbox_inventory_item_template.parent.id
}

output "child_id" {
  description = "ID of the child inventory item template"
  value       = netbox_inventory_item_template.child.id
}

output "child_parent_id" {
  description = "Parent ID of the child inventory item template"
  value       = netbox_inventory_item_template.child.parent
}

output "parent_child_match" {
  description = "Child parent ID matches parent ID"
  value       = netbox_inventory_item_template.child.parent == netbox_inventory_item_template.parent.id
}

output "complete_id" {
  description = "ID of the complete inventory item template"
  value       = netbox_inventory_item_template.complete.id
}

output "complete_name" {
  description = "Name of the complete inventory item template"
  value       = netbox_inventory_item_template.complete.name
}

output "complete_label" {
  description = "Label of the complete inventory item template"
  value       = netbox_inventory_item_template.complete.label
}

output "complete_part_id" {
  description = "Part ID of the complete inventory item template"
  value       = netbox_inventory_item_template.complete.part_id
}

output "complete_description" {
  description = "Description of the complete inventory item template"
  value       = netbox_inventory_item_template.complete.description
}

output "device_type_id" {
  description = "ID of the parent device type"
  value       = netbox_device_type.test.id
}

output "manufacturer_id" {
  description = "ID of the manufacturer"
  value       = netbox_manufacturer.test.id
}
