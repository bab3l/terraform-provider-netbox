# Inventory Item Data Source Outputs

output "by_id_name" {
  value = data.netbox_inventory_item.by_id.name
}

output "by_id_device_id" {
  value = data.netbox_inventory_item.by_id.device_id
}

output "by_id_description" {
  value = data.netbox_inventory_item.by_id.description
}
