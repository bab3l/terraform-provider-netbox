# Look up inventory item role by ID
data "netbox_inventory_item_role" "by_id" {
  id = "1"
}

# Look up inventory item role by name
data "netbox_inventory_item_role" "by_name" {
  name = "Power Supply"
}

# Look up inventory item role by slug
data "netbox_inventory_item_role" "by_slug" {
  slug = "power-supply"
}

# Use inventory item role data in other resources
output "role_id" {
  value = data.netbox_inventory_item_role.by_id.id
}

output "role_name" {
  value = data.netbox_inventory_item_role.by_name.name
}

output "role_slug" {
  value = data.netbox_inventory_item_role.by_slug.slug
}

output "role_color" {
  value = data.netbox_inventory_item_role.by_slug.color
}

# Access all custom fields
output "role_custom_fields" {
  value       = data.netbox_inventory_item_role.by_id.custom_fields
  description = "All custom fields defined in NetBox for this inventory item role"
}

# Access a specific custom field by name
output "role_replacement_cycle" {
  value       = try([for cf in data.netbox_inventory_item_role.by_id.custom_fields : cf.value if cf.name == "replacement_cycle_months"][0], null)
  description = "Example: accessing a specific custom field value"
}
