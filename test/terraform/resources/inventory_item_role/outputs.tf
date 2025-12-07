# Outputs to verify resource creation

output "basic_id" {
  description = "ID of the basic inventory item role"
  value       = netbox_inventory_item_role.basic.id
}

output "basic_name" {
  description = "Name of the basic inventory item role"
  value       = netbox_inventory_item_role.basic.name
}

output "basic_slug" {
  description = "Slug of the basic inventory item role"
  value       = netbox_inventory_item_role.basic.slug
}

output "basic_id_valid" {
  description = "Basic inventory item role has valid ID"
  value       = netbox_inventory_item_role.basic.id != ""
}

output "complete_id" {
  description = "ID of the complete inventory item role"
  value       = netbox_inventory_item_role.complete.id
}

output "complete_name" {
  description = "Name of the complete inventory item role"
  value       = netbox_inventory_item_role.complete.name
}

output "complete_description" {
  description = "Description of the complete inventory item role"
  value       = netbox_inventory_item_role.complete.description
}
