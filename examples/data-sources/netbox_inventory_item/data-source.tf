# Look up inventory item by ID
data "netbox_inventory_item" "by_id" {
  id = "1"
}

# Look up inventory item by name and device
data "netbox_inventory_item" "by_name_device" {
  name      = "Module1"
  device_id = "5"
}

# Use inventory item data in other resources
output "item_id" {
  value = data.netbox_inventory_item.by_id.id
}

output "item_name" {
  value = data.netbox_inventory_item.by_name_device.name
}

output "item_device" {
  value = data.netbox_inventory_item.by_name_device.device_id
}

output "item_manufacturer" {
  value = data.netbox_inventory_item.by_id.manufacturer
}

output "item_serial" {
  value = data.netbox_inventory_item.by_id.serial
}

# Access all custom fields
output "item_custom_fields" {
  value       = data.netbox_inventory_item.by_id.custom_fields
  description = "All custom fields defined in NetBox for this inventory item"
}

# Access specific custom fields by name
output "item_purchase_date" {
  value       = try([for cf in data.netbox_inventory_item.by_id.custom_fields : cf.value if cf.name == "purchase_date"][0], null)
  description = "Example: accessing a date custom field"
}

output "item_warranty_months" {
  value       = try([for cf in data.netbox_inventory_item.by_id.custom_fields : cf.value if cf.name == "warranty_months"][0], null)
  description = "Example: accessing a numeric custom field"
}
