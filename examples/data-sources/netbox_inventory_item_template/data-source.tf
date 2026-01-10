# Look up an inventory item template by ID
data "netbox_inventory_item_template" "by_id" {
  id = "1"
}

# Individual attribute outputs
output "inventory_item_template_id" {
  value       = data.netbox_inventory_item_template.by_id.id
  description = "The unique ID of the inventory item template"
}

output "inventory_item_template_name" {
  value       = data.netbox_inventory_item_template.by_id.name
  description = "The name of the inventory item template"
}

output "inventory_item_template_device_type_id" {
  value       = data.netbox_inventory_item_template.by_id.device_type_id
  description = "The device type ID this template belongs to"
}

output "inventory_item_template_description" {
  value       = data.netbox_inventory_item_template.by_id.description
  description = "Description of the inventory item template"
}

# Note: Inventory item templates do not support custom fields in NetBox API
output "inventory_item_template_note" {
  value       = "Inventory item templates are read-only device design metadata"
  description = "Template datasources do not support custom fields"
}
